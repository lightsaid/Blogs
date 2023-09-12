package dbrepo

type Repository struct {
	UserRepo     UserRepo
	PostsRepo    PostsRepo
	CategoryRepo CategoryRepo
	TagRepo      TagRepo
	SessionRepo  SessionRepo
	AssetsRepo   AssetsRepo
	CommentRepo  CommentRepo
}

// NewRepository 创建一个Repository仓库，使用 Queryable 接口，同时兼容 sql.DB 和 sql.Tx 接口
func NewRepository(db Queryable) *Repository {
	return &Repository{
		&userRepo{DB: db},
		&postsRepo{DB: db},
		&categoryRepo{DB: db},
		&tagRepo{DB: db},
		&sessionRepo{DB: db},
		&assetsRepo{DB: db},
		&commentRepo{DB: db},
	}
}
