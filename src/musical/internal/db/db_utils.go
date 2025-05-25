package db

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

func GetTxFromContext(ctx context.Context) (txAdapter, bool) {
	tx, ok := ctx.Value(txKey{}).(*pgxpool.Tx)
	return txAdapter{tx}, ok
}
