// Package logger предоставляет функции для инициализации и использования логгера в приложении.
package logger

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TeoPlow/online-music-service/src/musical/internal/config"
)

var Logger = slog.Default()

func InitLogger() {
	var lvl slog.Level
	switch config.Config.Log.Level {
	case "error":
		lvl = slog.LevelError
	case "warn":
		lvl = slog.LevelWarn
	case "info":
		lvl = slog.LevelInfo
	case "debug":
		lvl = slog.LevelDebug
	}

	var h slog.Handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: lvl})
	if config.Config.Log.File != "" {
		f, err := os.OpenFile(config.Config.Log.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o666)
		if err != nil {
			Logger.Error(fmt.Sprint(err))
		}
		h = slog.NewJSONHandler(f, &slog.HandlerOptions{Level: lvl})
	}
	Logger = slog.New(h)
}
