package models

type Category struct {
	Model
	Title string `db:"title" json:"title"`
	Slug  string `db:"slug" json:"slug"`
}
