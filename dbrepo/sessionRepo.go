package dbrepo

import "github.com/lightsaid/blogs/models"

// SessionRepo 定义Session表操作方法
type SessionRepo interface {
	Insert(Session *models.Session) (int64, error)
}

// 接口检查
var _ SessionRepo = (*sessionRepo)(nil)

// sessionRepo 实现 SessionRepo 接口
type sessionRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *sessionRepo) Insert(Session *models.Session) (int64, error) {
	return 0, nil
}
