package middlewares

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger 请求logger
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		ss := time.Since(startTime).Milliseconds()
		slog.InfoContext(
			r.Context(),
			r.Method+" "+r.URL.RequestURI(),
			slog.String("startTime", startTime.Format(time.RFC3339)),
			slog.Int64("duration(ms)", ss),
		)
	})
}
