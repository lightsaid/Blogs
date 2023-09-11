package models

type Comment struct {
	Model
	UserID   int64  `db:"user_id" json:"user_id"`
	ParentID int64  `db:"parent_id" json:"parent_id"`
	Content  string `db:"content" json:"content"`
	Nickname string `db:"nickname" json:"nickname"`
	Email    string `db:"email" json:"email"`
}
