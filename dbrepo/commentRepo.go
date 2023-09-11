package dbrepo

import "github.com/lightsaid/blogs/models"

// CommentRepo 定义Comment表操作方法
type CommentRepo interface {
	Insert(Comment *models.Comment) (int64, error)
}

// 接口检查
var _ CommentRepo = (*commentRepo)(nil)

// commentRepo 实现 CommentRepo 接口
type commentRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *commentRepo) Insert(Comment *models.Comment) (int64, error) {
	return 0, nil
}
