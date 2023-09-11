package dbrepo

import "github.com/lightsaid/blogs/models"

// AssetsRepo 定义Assets表操作方法
type AssetsRepo interface {
	Insert(Assets *models.Assets) (int64, error)
}

// 接口检查
var _ AssetsRepo = (*assetsRepo)(nil)

// assetsRepo 实现 AssetsRepo 接口
type assetsRepo struct {
	BD Queryable
}

// Insert 添加一个分类
func (store *assetsRepo) Insert(Assets *models.Assets) (int64, error) {
	return 0, nil
}
