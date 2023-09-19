package middlewares

import (
	"errors"
	"net/http"

	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/routers/contexts"
	v1 "github.com/lightsaid/blogs/routers/controllers/v1"
)

// RequirePermission 指定权限通行，默认会先执行RequireAuth中间件
func RequirePermission(next http.HandlerFunc, roles ...int) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := contexts.ContextGetAuthUser(r)
		if user == nil {
			v1.Write(w, r, nil, errs.ErrInternalServer.AsException(errors.New("contexts.ContextGetAuthUser -> nil")))
			return
		}

		for _, role := range roles {
			if role == user.Role {
				next.ServeHTTP(w, r)
				return
			}
		}

		v1.Write(w, r, nil, errs.ErrForbidden)
	})

	// 加一层 http.HandlerFunc 为了使用者方便，少写代码
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 先执行 登录信息验证 在执行权限认证
		RequireAuth(fn).ServeHTTP(w, r)
	})
}
