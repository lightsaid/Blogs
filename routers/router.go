package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/blogs/dbrepo"
	v1 "github.com/lightsaid/blogs/routers/controllers/v1"
	"github.com/lightsaid/blogs/routers/middlewares"
)

// NewRouter 创建一个 HTTP API 路由器
func NewRouter(db *sqlx.DB) http.Handler {

	store := dbrepo.NewRepository(db)
	postsCtlr := v1.NewPostsController(store.PostsRepo)
	categoryCtrl := v1.NewCategoryController(store.CategoryRepo)
	tagCtrl := v1.NewTagController(store.TagRepo)

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

	apiV1.Route("/auth", func(r chi.Router) {
		r.Post("/login", nil)   // 登录
		r.Post("/logout", nil)  // 注销
		r.Post("/refresh", nil) // 刷新 token
	})

	apiV1.Route("/profile", func(r chi.Router) {
		r.Get("/", nil) // 获取个人信息
		r.Put("/", nil) // 更新个人信息
	})

	apiV1.Route("/posts", func(r chi.Router) {
		r.Post("/", postsCtlr.Add)           // 新增
		r.Put("/", nil)                      // 更新
		r.Delete("/{id:^[0-9]+}", nil)       // 删除
		r.Get("/", nil)                      // 获取列表
		r.Get("/{id:^[0-9]+}", nil)          // 获取详情
		r.Get("/category/{id:^[0-9]+}", nil) // 根据分类获取文章列表
		r.Get("/tag/{id:^[0-9]+}", nil)      // 根据tag获取文章列表
		r.Get("/search/{keyword}", nil)      // 查询文章，获取列表
	})

	apiV1.Route("/category", func(r chi.Router) {
		r.Post("/", categoryCtrl.Add)                  // 添加
		r.Put("/", categoryCtrl.Update)                // 更新
		r.Delete("/{id:^[0-9]+}", categoryCtrl.Delete) // 删除
		r.Get("/", categoryCtrl.List)                  // 获取列表
		r.Get("/{id:^[0-9]+}", categoryCtrl.Get)       // 获取单个
	})

	apiV1.Route("/tags", func(r chi.Router) {
		r.Post("/", tagCtrl.Add)                  // 添加
		r.Put("/", tagCtrl.Update)                // 更新
		r.Delete("/{id:^[0-9]+}", tagCtrl.Delete) // 删除
		r.Get("/", tagCtrl.List)                  // 获取列表
		r.Get("/{id:^[0-9]+}", tagCtrl.Get)       // 获取单个
	})

	apiV1.Route("/assets", func(r chi.Router) {
		r.Post("/", nil)               // 添加
		r.Put("/", nil)                // 更新
		r.Delete("/{id:^[0-9]+}", nil) // 删除
		r.Get("/", nil)                // 获取列表
		r.Get("/{id:^[0-9]+}", nil)    // 获取单个
	})

	apiV1.Route("/comment", func(r chi.Router) {
		r.Post("/", nil)               // 添加
		r.Put("/", nil)                // 更新
		r.Delete("/{id:^[0-9]+}", nil) // 删除
		r.Get("/", nil)                // 获取列表
		r.Get("/{id:^[0-9]+}", nil)    // 获取单个
	})

	// 将v1路由附加到主路由上
	mux.Mount("/api/v1", apiV1)

	return mux
}
