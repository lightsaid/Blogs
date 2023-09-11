package middlewares

import (
	"log/slog"
	"net/http"
)

// Recoverer 恐慌恢复中间件
func Recoverer(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if rvr == http.ErrAbortHandler {
					// 参考 chi Recoverer 中间件
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rvr)
				}

				// logEntry := GetLogEntry(r)
				// if logEntry != nil {
				// 	logEntry.Panic(rvr, debug.Stack())
				// } else {
				// 	PrintPrettyStack(rvr)
				// }

				// TODO: 详细log
				slog.ErrorContext(r.Context(), "Request Panic", slog.Any("error", rvr))

				if r.Header.Get("Connection") != "Upgrade" {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		}()

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
