package v1

import (
	"log/slog"
	"net/http"
)

func HealthZ(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "v1 health check")
	w.Write([]byte("ok"))
}
