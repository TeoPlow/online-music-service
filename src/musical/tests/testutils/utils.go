// Package testutils содержит различные утилиты, которые используются для тестирования
package testutils

import (
	"context"
	"fmt"
	"time"

	"github.com/TeoPlow/online-music-service/src/musical/internal/db"
)

func ToPtr[T any](val T) *T {
	return &val
}

func DateOnly(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func TruncateTables(database db.Database, tables ...string) {
	for _, table := range tables {
		query := fmt.Sprintf(`TRUNCATE TABLE "%s" RESTART IDENTITY CASCADE`, table)
		_, err := database.Exec(context.Background(), query)
		if err != nil {
			panic(fmt.Errorf("failed to truncate table %s: %w", table, err))
		}
	}
}
