package db

import (
	"context"

	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
	cluster *pgxpool.Pool
}

func (db Database) Close() {
	db.cluster.Close()
}

func (db Database) Get(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Get(ctx, db.cluster, dst, query, args...)
}

func (db Database) Select(ctx context.Context, dst any, query string, args ...any) error {
	return pgxscan.Select(ctx, db.cluster, dst, query, args...)
}

func (db Database) Exec(ctx context.Context, query string, args ...any) (pgconn.CommandTag, error) {
	return db.cluster.Exec(ctx, query, args...)
}

func (db Database) ExecQueryRow(ctx context.Context, query string, args ...any) pgx.Row {
	return db.cluster.QueryRow(ctx, query, args...)
}

func (db Database) ExecQuery(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return db.cluster.Query(ctx, query, args...)
}
