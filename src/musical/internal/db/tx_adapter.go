package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type txAdapter struct {
	tx *pgxpool.Tx
}

func (a *txAdapter) Select(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Select(ctx, a.tx, dst, query, args...)
}

func (a *txAdapter) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, a.tx, dst, query, args...)
}

func (a *txAdapter) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return a.tx.Exec(ctx, query, args...)
}

func (a *txAdapter) ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return a.tx.QueryRow(ctx, query, args...)
}
