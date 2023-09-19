package v1

import (
	"errors"
	"net/http"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/dbrepo"
	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/routers/contexts"
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

	wErr := writeCookie(w, r, config.AppConf.Cookie.Name, res.AccessToken, 0)
	if wErr != nil {
		errorResponse(w, r, http.StatusInternalServerError, wErr, "服务错误，登录失败！")
		return
	}
	data := envelop{"data": res, "msg": successText}
	successResponse(w, r, data)
}

func (uc *UserController) Refresh(w http.ResponseWriter, r *http.Request) {
	var req forms.RefreshRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	aToken, err := uc.server.RenewAccessToken(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	wErr := writeCookie(w, r, config.AppConf.Cookie.Name, aToken, 0)
	if wErr != nil {
		errorResponse(w, r, http.StatusInternalServerError, wErr, "服务错误，刷新失败")
		return
	}

	data := envelop{"msg": successText}
	successResponse(w, r, data)
}

func (uc *UserController) Logout(w http.ResponseWriter, r *http.Request) {
	// 删除cookie
	err := writeCookie(w, r, config.AppConf.Cookie.Name, "", -1)
	if err != nil {
		e := errs.ErrInternalServer.AsException(err)
		errorResponse(w, r, e.StatusCode(), e, e.Message())
		return
	}

	successResponse(w, r, envelop{"msg": successText})
}

func (uc *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	user := contexts.ContextGetAuthUser(r)
	if user == nil {
		e := errs.ErrInternalServer.AsException(errors.New("contexts.ContextGetAuthUser(r) -> nil"))
		errorResponse(w, r, e.StatusCode(), e, e.Message())
		return
	}

	res := forms.LoginResponse{
		ID:       user.ID,
		UserName: user.UserName,
		Avatar:   user.Avatar,
	}

	data := envelop{"data": res, "msg": successText}

	successResponse(w, r, data)
}

func (uc *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req forms.UpdateUserRequest
	if ok := bindRequest(w, r, &req); !ok {
		return
	}

	err := uc.server.Update(r.Context(), req)
	if err != nil {
		errorResponse(w, r, err.StatusCode(), err, err.Message())
		return
	}

	successResponse(w, r, envelop{"msg": successText})
}
