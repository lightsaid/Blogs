package dbrepo

import (
	"context"

	"github.com/lightsaid/blogs/models"
)

// SessionRepo 定义Session表操作方法
type SessionRepo interface {
	// Insert 添加一个session
	Insert(ctx context.Context, session *models.Session) (int64, error)
	// Get 获取一个session
	Get(ctx context.Context, id int64) (session *models.Session, err error)
}

// 接口检查
var _ SessionRepo = (*sessionRepo)(nil)

// sessionRepo 实现 SessionRepo 接口
type sessionRepo struct {
	DB Queryable
}

// Insert 添加一个session
func (store *sessionRepo) Insert(ctx context.Context, session *models.Session) (int64, error) {
	query := `insert into sessions(user_id, refresh_token, client_ip, expired_at) 
	values($1, $2, $3, $4);`
	result, err := store.DB.ExecContext(ctx, query, session.UserID, session.RefreshToken, session.ClientIP, session.ExpiredAt)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// Get 获取一个session
func (store *sessionRepo) Get(ctx context.Context, id int64) (session *models.Session, err error) {
	query := `select id, user_id, refresh_token, client_ip, created_at, expired_at where id=$1`
	session = new(models.Session)
	err = store.DB.GetContext(ctx, session, query, id)
	return
}
