package dbrepo

import (
	"context"
	"database/sql"

	"github.com/lightsaid/blogs/models"
)

// CategoryRepo 定义Category表操作方法
type CategoryRepo interface {
	// Insert 添加一个分类
	Insert(ctx context.Context, category *models.Category) (int64, error)
	// GetAll 获取所有分类
	GetAll(ctx context.Context) (list []*models.Category, err error)
	// Get 获取一个分类
	Get(ctx context.Context, id int64) (category *models.Category, err error)
	// Update 更新一个分类
	Update(ctx context.Context, category *models.Category) error
	// 删除一个分类
	Delete(ctx context.Context, id int64) error
}

// 接口检查
var _ CategoryRepo = (*categoryRepo)(nil)

// categoryRepo 实现 CategoryRepo 接口
type categoryRepo struct {
	DB Queryable
}

// Insert 添加一个分类
func (store *categoryRepo) Insert(ctx context.Context, category *models.Category) (int64, error) {
	querySQL := `insert into category(title, slug) values();`
	result, err := store.DB.ExecContext(ctx, querySQL, category.Title, category.Slug)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// GetAll 获取所有分类
func (store *categoryRepo) GetAll(ctx context.Context) (list []*models.Category, err error) {
	querySQL := `select id, title, slug, created_at, updated_at from category where deleted_at is null`
	err = store.DB.SelectContext(ctx, &list, querySQL)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// Get 获取一个分类
func (store *categoryRepo) Get(ctx context.Context, id int64) (category *models.Category, err error) {
	query := `select id, title, slug, created_at, updated_at from category where id=$1 and deleted_at is null`
	category = new(models.Category)
	store.DB.GetContext(ctx, category, query, id)
	return
}

// Update 更新一个分类
func (store *categoryRepo) Update(ctx context.Context, category *models.Category) error {
	query := `
	update category 
	set
		title=$1,
		slug=$2,
		updated_at=datetime('now', 'localtime')
	where id=$3 and deleted_at is null;
	`

	result, err := store.DB.ExecContext(ctx, query, category.Title, category.Slug, category.ID)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	return err
}

// 删除一个分类
func (store *categoryRepo) Delete(ctx context.Context, id int64) error {
	query := `
	update category 
	set
		deleted_at=datetime('now', 'localtime')
	where id=$1;
	`

	result, err := store.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	return err
}
