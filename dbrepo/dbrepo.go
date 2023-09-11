package dbrepo

type Repository struct {
}

// NewRepository 创建一个Repository仓库，使用 Queryable 接口，同时兼容 sql.DB 和 sql.Tx 接口
func NewRepository(db Queryable) *Repository {
	return &Repository{}
}
