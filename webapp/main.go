package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"claw-destine.com/camboose/webapp/projects"
	"claw-destine.com/camboose/webapp/tasks"
	"github.com/a-h/templ"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(status int) {
	lrw.status = status
	lrw.ResponseWriter.WriteHeader(status)
}

func withRequestLogging(next http.Handler) http.Handler {
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

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("./static")))

	mux.Handle("/components/projects", templ.Handler(projects.Projects()))
	mux.Handle("/components/tasks", templ.Handler(tasks.Tasks()))

	server := &http.Server{
		Addr:    ":3000",
		Handler: withRequestLogging(mux),
	}

	slog.Info("Starting HTTP server", "address", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("HTTP server stopped", "error", err)
	}

}
