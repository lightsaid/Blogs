package models

type Assets struct {
	Model
	UserID  int64  `db:"user_id" json:"user_id"`
	PostsID int64  `db:"posts_id" json:"posts_id"`
	Data    []byte `db:"data" json:"data"`
	Ext     string `db:"ext" json:"ext"`
	Name    string `db:"name" json:"name"`
	Size    int64  `db:"size" json:"size"`
}
