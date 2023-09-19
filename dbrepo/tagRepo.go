package dbrepo

import (
	"context"
	"database/sql"

	"github.com/lightsaid/blogs/models"
)

// TagRepo 定义Tag表操作方法
type TagRepo interface {
	// Insert 添加一个tag
	Insert(ctx context.Context, tag *models.Tag) (int64, error)
	// GetAll 获取所有tag
	GetAll(ctx context.Context) (list []*models.Tag, err error)
	// Get 获取一个分类
	Get(ctx context.Context, id int64) (tag *models.Tag, err error)
	// Update 更新一个tag
	Update(ctx context.Context, tag *models.Tag) error
	// 删除一个tag
	Delete(ctx context.Context, id int64) error
}

// 接口检查
var _ TagRepo = (*tagRepo)(nil)

// tagRepo 实现 TagRepo 接口
type tagRepo struct {
	DB Queryable
}

// Insert 添加一个tag
func (store *tagRepo) Insert(ctx context.Context, tag *models.Tag) (int64, error) {
	querySQL := `insert into tags(title, slug) values($1, $2);`
	result, err := store.DB.ExecContext(ctx, querySQL, tag.Title, tag.Slug)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// GetAll 获取所有tag
func (store *tagRepo) GetAll(ctx context.Context) (list []*models.Tag, err error) {
	querySQL := `select id, title, slug, created_at, updated_at from tags where deleted_at is null`
	err = store.DB.SelectContext(ctx, &list, querySQL)
	if err == sql.ErrNoRows {
		err = nil
	}
	return
}

// Get 获取一个分类
func (store *tagRepo) Get(ctx context.Context, id int64) (tag *models.Tag, err error) {
	query := `select id, title, slug, created_at, updated_at from tags where id=$1 and deleted_at is null`
	tag = new(models.Tag)
	err = store.DB.GetContext(ctx, tag, query, id)
	return
}

// Update 更新一个tag
func (store *tagRepo) Update(ctx context.Context, tag *models.Tag) error {
	query := `
	update tags 
	set
		title=$1,
		slug=$2,
		updated_at=datetime('now', 'localtime')
	where id=$3 and deleted_at is null;
	`

	result, err := store.DB.ExecContext(ctx, query, tag.Title, tag.Slug, tag.ID)
	if err != nil {
		return err
	}
	return handleRowsAffected(result.RowsAffected())
}

// 删除一个tag
func (store *tagRepo) Delete(ctx context.Context, id int64) error {
	query := `
	update tags 
	set
		deleted_at=datetime('now', 'localtime')
	where id=$1 and deleted_at is null;
	`

	result, err := store.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return handleRowsAffected(result.RowsAffected())
}
