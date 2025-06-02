// Package storage реализует доступ к хранилищу данных и обработку ошибок уровня хранилища.
package storage

import (
	"context"
	"errors"
	"log/slog"

	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
	"github.com/TeoPlow/online-music-service/src/musical/internal/db/cache"
	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

var (
	ErrNotExists = errors.New("not exists")
	ErrSQL       = errors.New("error in SQL query")
)

// Порог, который используется для текстового поиска
const threshold = 0.1

type Executor interface {
	Select(ctx context.Context, dst any, query string, args ...any) error
	Get(ctx context.Context, dst any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type Cache interface {
	Set(сtx context.Context, model, id string, value any) error
	Get(сtx context.Context, model, id string, dst any) error
	Delete(сtx context.Context, model, id string) error
}

type repository struct {
	db Executor
	c  Cache
}

func (r *repository) getExecutor(ctx context.Context) Executor {
	if tx, ok := db.GetTxFromContext(ctx); ok {
		return &tx
	}
	return r.db
}

func (r *repository) checkCache(ctx context.Context, model string, id uuid.UUID, dst any) bool {
	if r.c == nil {
		return false
	}
	if err := r.c.Get(ctx, model, id.String(), dst); err != nil {
		if errors.Is(err, cache.ErrNotFound) {
			return false
		}
		logger.Logger.Error("failed get from cache",
			slog.String("error", err.Error()))
		return false
	}
	return true
}

func (r *repository) updateCache(ctx context.Context, model string, id uuid.UUID, data any) {
	if r.c == nil {
		return
	}
	if err := r.c.Set(ctx, model, id.String(), data); err != nil {
		logger.Logger.Error("failed set cache",
			slog.String("error", err.Error()))
	}
}

func (r *repository) clearCache(ctx context.Context, model string, id uuid.UUID) {
	if r.c == nil {
		return
	}
	if err := r.c.Delete(ctx, model, id.String()); err != nil {
		logger.Logger.Error("failed delete from cache",
			slog.String("error", err.Error()))
	}
}