package handlers

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.Info("new request", "method", r.Method, "path", r.URL.Path, "duration", time.Since(start), "from", r.RemoteAddr)

	})
}
