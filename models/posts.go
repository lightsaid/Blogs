package models

type Posts struct {
	Model
	AuthorID int64  `db:"author_id" json:"author_id"`
	Title    string `db:"title" json:"title"`
	Content  string `db:"content" json:"content"`
	Keyword  string `db:"keyword" json:"keyword"`
	Slug     string `db:"slug" json:"slug"`
	Abstract string `db:"abstract" json:"abstract"`
	CoverID  int64  `db:"cover_image_id" json:"cover_id"`
	Views    int    `db:"views" json:"views"`
	Comments int    `db:"comments" json:"comments"`
	Likes    int    `db:"likes" json:"likes"`
}

type PostsCategory struct {
	PostsID    int64 `db:"posts_id" json:"posts_id"`
	CategoryID int64 `db:"category_id" json:"category_id"`
}

type PostsTag struct {
	PostsID int64 `db:"posts_id" json:"posts_id"`
	TagID   int64 `db:"tag_id" json:"tag_id"`
}
