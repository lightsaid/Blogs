package v1

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/cookie"
	"github.com/lightsaid/blogs/errs"
	"github.com/lightsaid/blogs/request"
	"github.com/lightsaid/blogs/respond"
)

const (
	successText = "successful"
)

type envelop map[string]interface{}

func bindRequest(w http.ResponseWriter, r *http.Request, req interface{}) bool {
	err := request.ReadJSON(w, r, &req)
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, err.Error())
		return false
	}

	if err = config.Validate.Struct(req); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			slog.ErrorContext(r.Context(), "请使用结构体指针作为参数", slog.String("error", err.Error()))
			errorResponse(w, r, http.StatusBadRequest, err, "参数格式错误")
			return false
		}

		errs := err.(validator.ValidationErrors)

		var msgs []string
		for _, e := range errs.Translate(config.Trans) {
			msgs = append(msgs, e)
		}

		errorResponse(w, r, http.StatusBadRequest, err, strings.Join(msgs, ", "))
		return false
	}
	return true
}

func bindParamInt64(w http.ResponseWriter, r *http.Request, key string) (int64, bool) {
	val, err := strconv.Atoi(chi.URLParam(r, key))
	if err != nil {
		errorResponse(w, r, http.StatusBadRequest, err, "请输入合法参数")
		return 0, false
	}

	return int64(val), true
}

func errorResponse(w http.ResponseWriter, r *http.Request, status int, err error, msg string) {
	slog.ErrorContext(r.Context(), msg, slog.String("error", err.Error()))
	data := envelop{"msg": msg}
	if err := respond.New(w).Status(status).JSON(data); err != nil {
		slog.ErrorContext(r.Context(), msg, slog.String("error", err.Error()), slog.Any("data", data))
	}
}

func successResponse(w http.ResponseWriter, r *http.Request, data interface{}) {
	err := respond.New(w).JSON(data)
	if err != nil {
		slog.ErrorContext(r.Context(), "successResponse", slog.String("error", err.Error()), slog.Any("data", data))
	}
}

func writeCookie(w http.ResponseWriter, r *http.Request, cookieName, cookieValue string, maxAge int, expires ...time.Time) error {
	var secure bool
	if config.AppConf.Server.Env == config.EnvProd {
		secure = true
	}

	expiresAt := time.Now().Add(config.ParseDuration(config.AppConf.Token.RefreshExpire, 72*time.Hour))
	if len(expires) > 0 {
		expiresAt = expires[0]
	}

	// 写入 cookie
	httpCookie := http.Cookie{
		Name:     cookieName,
		Value:    cookieValue,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true, // 仅在HTTP请求中传递，禁止JavaScript访问
		SameSite: http.SameSiteLaxMode,
		Secure:   secure, // 是否启用https
		MaxAge:   maxAge,
	}

	return cookie.WriteSigned(w, httpCookie, config.AppConf.Cookie.SecretKey)
}

// Write 定义一个公开方法，提供给外部包使用（middleware）
// data 只有 err.StatusCode() == 200 才会使用
func Write(w http.ResponseWriter, r *http.Request, data interface{}, err *errs.AppError) {
	if err.StatusCode() == http.StatusOK {
		successResponse(w, r, data)
		return
	}
	errorResponse(w, r, err.StatusCode(), err, err.Message())
}
