package dbrepo

import "github.com/lightsaid/blogs/models"

// TagRepo 定义Tag表操作方法
type TagRepo interface {
	Insert(Tag *models.Tag) (int64, error)
}

// 接口检查
var _ TagRepo = (*tagRepo)(nil)

// tagRepo 实现 TagRepo 接口
type tagRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *tagRepo) Insert(Tag *models.Tag) (int64, error) {
	return 0, nil
}
