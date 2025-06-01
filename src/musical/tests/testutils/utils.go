// Package testutils содержит различные утилиты, которые используются для тестирования
package testutils

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"

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

func TrackFromFile(name string) *bytes.Buffer {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic(os.ErrNotExist)
	}

	path := filepath.Join(filepath.Dir(filename), "../testdata/", name+".mp3")
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			panic(err)
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	res := make([]byte, stat.Size())
	if _, err = file.Read(res); err != nil {
		panic(err)
	}

	return bytes.NewBuffer(res)
}

func CompareTracks(expected, actual []byte) bool {
	hashExpected := sha256.Sum256(expected)
	hashActual := sha256.Sum256(actual)

	if hashExpected != hashActual {
		_ = os.WriteFile("failed_test.mp3", actual, 0o644)
		return false
	}
	return true
}

func ContextWithUserID(ctx context.Context, userID uuid.UUID) context.Context {
	md := metadata.Pairs("user_id", userID.String())
	return metadata.NewOutgoingContext(ctx, md)
}
