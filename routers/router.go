package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jmoiron/sqlx"
	"github.com/lightsaid/blogs/dbrepo"
	v1 "github.com/lightsaid/blogs/routers/controllers/v1"
	"github.com/lightsaid/blogs/routers/middlewares"
	"github.com/lightsaid/blogs/service"
)

// NewRouter 创建一个 HTTP API 路由器
func NewRouter(db *sqlx.DB) http.Handler {
	// 创建 dbrepo，实现 Controller
	store := dbrepo.NewRepository(db)
	postsCtlr := v1.NewPostsController(store.PostsRepo)
	categoryCtrl := v1.NewCategoryController(store.CategoryRepo)
	tagCtrl := v1.NewTagController(store.TagRepo)

	userCtrl := v1.NewUserController(store.UserRepo, store.SessionRepo)

	// 设置server提供给中间件使用
	userServer := service.NewUserServer(store.UserRepo, store.SessionRepo)
	middlewares.SetUserServer(userServer)

	// 主路由
	mux := chi.NewRouter()

	// chi 中间件
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)

	// 处理跨域中间件
	mux.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// 自定义中间件
	mux.Use(middlewares.Logger)
	mux.Use(middlewares.Recoverer)

	// v1 路由
	apiV1 := chi.NewRouter()
	apiV1.Get("/healthz", v1.HealthZ)

	apiV1.Post("/register", userCtrl.Register)
	apiV1.Post("/login", userCtrl.Login)
	apiV1.Post("/refresh", userCtrl.Refresh) // 刷新 token

	apiV1.Route("/blogs", func(r chi.Router) {
		r.Get("/category", categoryCtrl.List)             // 获取分类列表
		r.Get("/category/{id:^[0-9]+}", categoryCtrl.Get) // 获取单个分类

		r.Get("/tags", tagCtrl.List)             // 获取标签列表
		r.Get("/tags/{id:^[0-9]+}", tagCtrl.Get) // 获取单个标签

		r.Get("/posts", nil)                       // 获取文章列表
		r.Get("/posts/{id:^[0-9]+}", nil)          // 获取文章详情
		r.Get("/posts/category/{id:^[0-9]+}", nil) // 根据分类获取文章列表
		r.Get("/posts/tag/{id:^[0-9]+}", nil)      // 根据tag获取文章列表
		r.Get("/posts/search/{keyword}", nil)      // 查找文章，获取列表
	})

	apiV1.Route("/auth", func(r chi.Router) {
		// 中间件
		r.Use(middlewares.RequireAuth)

		// 路由
		r.Post("/logout", userCtrl.Logout)        // 注销
		r.Get("/profile", userCtrl.GetProfile)    // 获取个人信息
		r.Put("/profile", userCtrl.UpdateProfile) // 更新个人信息
	})

	apiV1.Route("/category", func(r chi.Router) {
		r.Post("/", middlewares.RequirePermission(categoryCtrl.Add, 1))                  // 添加
		r.Put("/", middlewares.RequirePermission(categoryCtrl.Update, 1))                // 更新
		r.Delete("/{id:^[0-9]+}", middlewares.RequirePermission(categoryCtrl.Delete, 1)) // 删除
	})

	apiV1.Route("/tags", func(r chi.Router) {
		r.Post("/", middlewares.RequirePermission(tagCtrl.Add, 1))                  // 添加
		r.Put("/", middlewares.RequirePermission(tagCtrl.Update, 1))                // 更新
		r.Delete("/{id:^[0-9]+}", middlewares.RequirePermission(tagCtrl.Delete, 1)) // 删除
	})

	apiV1.Route("/posts", func(r chi.Router) {
		r.Post("/", middlewares.RequirePermission(postsCtlr.Add, 1)) // 新增
		r.Put("/", nil)                                              // 更新
		r.Delete("/{id:^[0-9]+}", nil)                               // 删除
	})

	apiV1.Route("/assets", func(r chi.Router) {
		// 中间件
		r.Use(middlewares.RequireAuth)

		// 路由
		r.Post("/", nil)               // 添加
		r.Put("/", nil)                // 更新
		r.Delete("/{id:^[0-9]+}", nil) // 删除
		r.Get("/", nil)                // 获取列表
		r.Get("/{id:^[0-9]+}", nil)    // 获取单个
	})

	// 评论路由，谁都可以，无需权限
	apiV1.Route("/comment", func(r chi.Router) {
		r.Post("/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("ok")) }) // 添加
		r.Put("/", nil)                                                                     // 更新
		r.Delete("/{id:^[0-9]+}", nil)                                                      // 删除
		r.Get("/", nil)                                                                     // 获取列表
		r.Get("/{id:^[0-9]+}", nil)                                                         // 获取单个
	})

	// 将v1路由附加到主路由上
	mux.Mount("/api/v1", apiV1)

	return mux
}
