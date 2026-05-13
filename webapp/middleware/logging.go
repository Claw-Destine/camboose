package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func WithRequestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &loggingResponseWriter{ResponseWriter: w, status: http.StatusOK}

		defer func() {
			if recovered := recover(); recovered != nil {
				lrw.status = http.StatusInternalServerError
				http.Error(lrw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				slog.Error("HTTP panic",
					"method", r.Method,
					"path", r.URL.Path,
					"error", fmt.Sprint(recovered),
				)
			}

			attrs := []any{
				"method", r.Method,
				"path", r.URL.Path,
				"status", lrw.status,
				"duration", time.Since(start),
				"remote_addr", r.RemoteAddr,
			}

			switch {
			case lrw.status >= http.StatusInternalServerError:
				slog.Error("HTTP request error", attrs...)
			case lrw.status >= http.StatusBadRequest:
				slog.Warn("HTTP request client error", attrs...)
			default:
				slog.Info("HTTP request", attrs...)
			}
		}()

		next.ServeHTTP(lrw, r)
	})
}
