// Package db отвечает за инициализацию и работу с базой данных.
package db

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
)

type TxManager struct {
	db Database
}

func NewTxManager(ctx context.Context) (*TxManager, error) {
	pool, err := pgxpool.Connect(ctx, config.Config.DBConn)
	if err != nil {
		return nil, err
	}
	return &TxManager{Database{cluster: pool}}, nil
}

func (m *TxManager) GetDatabase() Database {
	return m.db
}

func (m *TxManager) RunSerializable(ctx context.Context, f func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}
	return m.beginFunc(ctx, opts, f)
}

func (m *TxManager) RunReadUncommited(ctx context.Context, f func(context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   pgx.ReadUncommitted,
		AccessMode: pgx.ReadOnly,
	}
	return m.beginFunc(ctx, opts, f)
}

func (m *TxManager) beginFunc(ctx context.Context, opts pgx.TxOptions, f func(context.Context) error) error {
	tx, err := m.db.cluster.BeginTx(ctx, opts)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := f(ctx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
