package dbrepo

import "github.com/lightsaid/blogs/models"

// UserRepo 定义user表操作方法
type UserRepo interface {
	Insert(user *models.User) (int64, error)
}

// 接口检查
var _ UserRepo = (*userRepo)(nil)

// userRepo 实现 UserRepo 接口
type userRepo struct {
	BD Queryable
}

// Insert 添加一个用户
func (store *userRepo) Insert(user *models.User) (int64, error) {
	return 0, nil
}
