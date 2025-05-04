package http

import (
	"log/slog"
	"net/http"

	"github.com/TeoPlow/online-music-service/src/musical/internal/logger"
)

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func logMW(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Logger.Info("HTTP_Request",
			slog.String("Method", r.Method),
			slog.String("URL", r.RequestURI),
			slog.Any("Content-Type", r.Header["Content-Type"]),
		)

		wrappedWriter := &responseWriter{ResponseWriter: w, statusCode: 200}
		h.ServeHTTP(wrappedWriter, r)

		if wrappedWriter.statusCode >= 500 {
			logger.Logger.Error("HTTP_Request",
				slog.String("Method", r.Method),
				slog.String("URL", r.RequestURI),
				slog.Int("Status", wrappedWriter.statusCode),
			)
		} else {
			logger.Logger.Info("HTTP_Request",
				slog.String("Method", r.Method),
				slog.String("URL", r.RequestURI),
				slog.Int("Status", wrappedWriter.statusCode),
			)
		}
	})
}
