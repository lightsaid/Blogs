package dbrepo

import "github.com/lightsaid/blogs/models"

// CategoryRepo 定义Category表操作方法
type CategoryRepo interface {
	Insert(Category *models.Category) (int64, error)
}

// 接口检查
var _ CategoryRepo = (*categoryRepo)(nil)

// categoryRepo 实现 CategoryRepo 接口
type categoryRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *categoryRepo) Insert(Category *models.Category) (int64, error) {
	return 0, nil
}
