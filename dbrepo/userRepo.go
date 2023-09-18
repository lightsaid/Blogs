package dbrepo

import (
	"context"

	"github.com/lightsaid/blogs/models"
)

// UserRepo 定义user表操作方法
type UserRepo interface {
	// Insert 添加一个用户
	Insert(ctx context.Context, user *models.User) (int64, error)
	// GetByEmail 获取用户
	Get(ctx context.Context, id int64) (user *models.User, err error)
	// GetByEmail 获取用户
	GetByEmail(ctx context.Context, email string) (user *models.User, err error)
	// 激活用户
	Activate(ctx context.Context, id int64) error
	// 更新用户
	Update(ctx context.Context, user *models.User) error
}

// 接口检查
var _ UserRepo = (*userRepo)(nil)

// userRepo 实现 UserRepo 接口
type userRepo struct {
	DB Queryable
}

// Insert 添加一个用户
func (store *userRepo) Insert(ctx context.Context, user *models.User) (int64, error) {
	query := `insert into users(
		email,
		password,
		username,
		avatar
	) values($1, $2, $3, $4)`

	result, err := store.DB.ExecContext(ctx, query, user.Email, user.Password, user.UserName, user.Avatar)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetByEmail 获取用户
func (store *userRepo) Get(ctx context.Context, id int64) (user *models.User, err error) {
	user = new(models.User)
	query := `select 
		id, email, password, username, avatar, role, activated_at, created_at, updated_at
	from users where id=$1 and deleted_at is null`
	err = store.DB.GetContext(ctx, user, query, id)

	return
}

// GetByEmail 根据获取用户
func (store *userRepo) GetByEmail(ctx context.Context, email string) (user *models.User, err error) {
	user = new(models.User)
	query := `select 
		id, email, password, username, avatar, role, activated_at, created_at, updated_at
	from users where email=$1 and deleted_at is null`
	err = store.DB.GetContext(ctx, user, query, email)

	return
}

// Update 更新用户
func (store *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `
	update users 
	set
		username=$1,
		avatar=$2
	where id=$3 and deleted_at is null;
	`
	result, err := store.DB.ExecContext(ctx, query, user.UserName, user.Avatar, user.ID)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

// Activate 激活用户
func (store *userRepo) Activate(ctx context.Context, id int64) error {
	query := `update users set activated_at=datetime('now', 'localtime') where id=$1`
	result, err := store.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	return err
}
