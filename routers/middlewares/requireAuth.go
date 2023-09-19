package middlewares

import (
	"fmt"
	"net/http"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/cookie"
	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/routers/contexts"
	v1 "github.com/lightsaid/blogs/routers/controllers/v1"
	"github.com/lightsaid/blogs/service"
)

var userServer *service.UserServer

// 设置一个userServer提供给本包使用
func SetUserServer(server *service.UserServer) {
	userServer = server
}

// RequireAuth 必须登录认证才能访问
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("-->> RequireAuth")
		accessToken, err := cookie.ReadSigned(r, config.AppConf.Cookie.Name, config.AppConf.Cookie.SecretKey)
		if err != nil && err != http.ErrNoCookie {
			appErr := errs.ErrInternalServer.AsException(err)
			v1.Write(w, r, nil, appErr)
			return
		}

		tokenPayload, err := config.TokenMaker.ParseToken(accessToken)
		if err != nil {
			appErr := errs.ErrUnauthorized.AsException(err)
			v1.Write(w, r, nil, appErr)
			return
		}

		user, appErr := userServer.Get(r.Context(), tokenPayload.UserID)
		if appErr != nil {
			v1.Write(w, r, nil, appErr)
			return
		}

		r = contexts.ContextSetAuthUser(r, user)

		next.ServeHTTP(w, r)
	})
}
