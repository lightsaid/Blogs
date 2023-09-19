package contexts

import (
	"context"
	"net/http"

	"github.com/lightsaid/blogs/models"
)

type contextKey string

const (
	authUserContextKey = contextKey("user")
)

// contextSetAuthUser 设置认证用户基础信息
func ContextSetAuthUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), authUserContextKey, user)
	return r.WithContext(ctx)
}

// contextGetAuthUser 获取认证用户基础信息, 如果用户信息不存在，返回nil
func ContextGetAuthUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(authUserContextKey).(*models.User)
	if !ok {
		return nil
	}
	return user
}
