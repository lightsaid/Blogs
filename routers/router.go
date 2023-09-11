package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	v1 "github.com/lightsaid/blogs/routers/controllers/v1"
	"github.com/lightsaid/blogs/routers/middlewares"
)

// NewRouter 创建一个 HTTP API 路由器
func NewRouter() http.Handler {
	// 主路由
	mux := chi.NewRouter()

	// chi 中间件
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	// 自定义中间件
	mux.Use(middlewares.Logger)
	mux.Use(middlewares.Recoverer)

	// v1 路由
	apiV1 := chi.NewRouter()
	apiV1.Get("/healthz", v1.HealthZ)

	// 将v1路由附加到主路由上
	mux.Mount("/api/v1", apiV1)

	return mux
}
