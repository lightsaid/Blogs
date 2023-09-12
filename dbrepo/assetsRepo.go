package dbrepo

import (
	"context"

	"github.com/lightsaid/blogs/models"
)

// AssetsRepo 定义Assets表操作方法
type AssetsRepo interface {
	// Insert 添加一个资源
	Insert(ctx context.Context, assets *models.Assets) (int64, error)
	// Get 获取一个资源
	Get(id int64) (ctx context.Context, assets *models.Assets, err error)
	// GetPostsID 获取一个资源, 根据posts id
	GetPostsID(postsID int64) (ctx context.Context, assets *models.Assets, err error)
	// GetUserID 获取一个资源, 根据users id
	GetUserID(userID int64) (ctx context.Context, assets *models.Assets, err error)
	// GetListByUserID 根据用户id获取资源列表
	GetListByUserID(ctx context.Context, userID int64, filter Filters) (list []*models.Assets, err error)
	// GetListByPostsID 根据posts id获取资源列表
	GetListByPostsID(ctx context.Context, postsID int64, filter Filters) (list []*models.Assets, err error)
}

// 接口检查
var _ AssetsRepo = (*assetsRepo)(nil)

// assetsRepo 实现 AssetsRepo 接口
type assetsRepo struct {
	DB Queryable
}

// Insert 添加一个资源
func (store *assetsRepo) Insert(ctx context.Context, assets *models.Assets) (int64, error) {
	query := `insert into assets(user_id, posts_id, data, ext, name, size)
	values($1, $2, $3, $4, $5, $6)`

	result, err := store.DB.ExecContext(ctx, query, assets.UserID,
		assets.PostsID, assets.Data, assets.Ext, assets.Name, assets.Size)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

// Get 获取一个资源
func (store *assetsRepo) Get(id int64) (ctx context.Context, assets *models.Assets, err error) {
	query := `select id, user_id, posts_id, data, ext, name, size, created_at, updated_at where id=$1 and deleted_at is null`
	assets = new(models.Assets)
	err = store.DB.GetContext(ctx, assets, query, id)
	return
}

// GetPostsID 获取一个资源, 根据posts id
func (store *assetsRepo) GetPostsID(postsID int64) (ctx context.Context, assets *models.Assets, err error) {
	query := `select id, user_id, posts_id, data, ext, name, size, created_at, updated_at where posts_id=$1 and deleted_at is null`
	assets = new(models.Assets)
	err = store.DB.GetContext(ctx, assets, query, postsID)
	return
}

// GetUserID 获取一个资源, 根据users id
func (store *assetsRepo) GetUserID(userID int64) (ctx context.Context, assets *models.Assets, err error) {
	query := `select id, user_id, posts_id, data, ext, name, size, created_at, updated_at where user_id=$1 and deleted_at is null`
	assets = new(models.Assets)
	err = store.DB.GetContext(ctx, assets, query, userID)
	return
}

// GetListByUserID 根据用户id获取资源列表
func (store *assetsRepo) GetListByUserID(ctx context.Context, userID int64, filter Filters) (list []*models.Assets, err error) {
	query := `
	select 
		id, user_id, posts_id, data, ext, name, size, created_at, updated_at 
	where user_id=$1 and deleted_at is null
	order by $2
	limit $3 offet $4
	`
	err = store.DB.SelectContext(ctx, &list, query, userID, filter.defaultSort(), filter.limit(), filter.offset())
	return
}

// GetListByPostsID 根据posts id获取资源列表
func (store *assetsRepo) GetListByPostsID(ctx context.Context, postsID int64, filter Filters) (list []*models.Assets, err error) {
	query := `
	select 
		id, user_id, posts_id, data, ext, name, size, created_at, updated_at 
	where posts_id=$1 and deleted_at is null
	order by $2
	limit $3 offet $4
	`
	err = store.DB.SelectContext(ctx, &list, query, postsID, filter.defaultSort(), filter.limit(), filter.offset())
	return
}
