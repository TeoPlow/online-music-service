// Package storage реализует доступ к хранилищу данных и обработку ошибок уровня хранилища.
package storage

import (
	"context"
	"errors"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"

	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
)

var (
	ErrNotExists = errors.New("not exists")
	ErrSQL       = errors.New("error in SQL query")
)

type Executor interface {
	Select(ctx context.Context, dst any, query string, args ...any) error
	Get(ctx context.Context, dst any, query string, args ...any) error
	Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error)
	ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row
}

type repository struct {
	db Executor
}

func (r *repository) getExecutor(ctx context.Context) Executor {
	if tx, ok := db.GetTxFromContext(ctx); ok {
		return &tx
	}
	return r.db
}
