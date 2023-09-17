package v1

import (
	"log/slog"
	"net/http"

	"github.com/lightsaid/blogs/config"
)

func HealthZ(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "v1 health check")
	data := envelop{"version": config.AppConf.Server.Version, "status": "ok"}
	successResponse(w, r, data)
}
