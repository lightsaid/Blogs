package v1

import (
	"net/http"
	"time"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/cookie"
	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/routers/forms"
	"github.com/lightsaid/blogs/service"
)

type UserController struct {
	server *service.UserServer
}

func NewUserController(store dbrepo.UserRepo, session dbrepo.SessionRepo) *UserController {
	return &UserController{
		server: service.NewUserServer(store, session),
	}
}

func (uc *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var req forms.AddUserRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	newID, err := uc.server.Create(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}
	data := envelop{"id": newID, "msg": successText}
	successResponse(w, r, data)
}

func (uc *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var req forms.LoginRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}
	req.ClientIP = r.RemoteAddr
	res, err := uc.server.Login(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	var secure bool
	if config.AppConf.Server.Env == config.EnvProd {
		secure = true
	}

	// 写入 cookie
	httpCookie := http.Cookie{
		Name:     config.AppConf.Cookie.Name,
		Value:    res.AccessToken,
		Path:     "/",
		Expires:  time.Now().Add(config.ParseDuration(config.AppConf.Token.RefreshExpire, 72*time.Hour)),
		HttpOnly: true, // 仅在HTTP请求中传递，禁止JavaScript访问
		SameSite: http.SameSiteLaxMode,
		Secure:   secure, // 是否启用https
	}

	wErr := cookie.WriteSigned(w, httpCookie, config.AppConf.Cookie.SecretKey)
	if wErr != nil {
		errorResponse(w, r, http.StatusInternalServerError, wErr, "服务错误，登录失败！")
		return
	}
	data := envelop{"data": res, "msg": successText}
	successResponse(w, r, data)
}

func (uc *UserController) Refresh(w http.ResponseWriter, r *http.Request) {

}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {

}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {

}

func (uc *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {

}
