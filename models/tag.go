package models

type Tag struct {
	Model
	Title string `db:"title" json:"title"`
	Slug  string `db:"slug" json:"slug"`
}
