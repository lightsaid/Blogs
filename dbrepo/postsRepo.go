package dbrepo

import "github.com/lightsaid/blogs/models"

// PostsRepo 定义Posts表操作方法
type PostsRepo interface {
	Insert(Posts *models.Posts) (int64, error)
}

// 接口检查
var _ PostsRepo = (*postsRepo)(nil)

// postsRepo 实现 PostsRepo 接口
type postsRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *postsRepo) Insert(Posts *models.Posts) (int64, error) {
	return 0, nil
}
