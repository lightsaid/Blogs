package dbrepo

import "github.com/lightsaid/blogs/models"

// SQLColumn SQL 查询结果 Scan
type SQLColumn struct {
	ID        *int64  `db:"id"`
	Title     *string `db:"title"`
	Slug      *string `db:"slug"`
	CreatedAt *string `db:"created_at"`
	UpdatedAt *string `db:"updated_at"`
}

func (sc *SQLColumn) ToTag() models.Tag {
	var tag models.Tag
	if sc.ID != nil {
		tag.ID = *sc.ID
	}
	if sc.Title != nil {
		tag.Title = *sc.Title
	}
	if sc.Slug != nil {
		tag.Slug = *sc.Slug
	}
	if sc.CreatedAt != nil {
		tag.CreatedAt = *sc.CreatedAt
	}
	if sc.UpdatedAt != nil {
		tag.UpdatedAt = *sc.UpdatedAt
	}

	return tag
}

func (sc *SQLColumn) ToCategory() models.Category {
	var category models.Category
	if sc.ID != nil {
		category.ID = *sc.ID
	}
	if sc.Title != nil {
		category.Title = *sc.Title
	}
	if sc.Slug != nil {
		category.Slug = *sc.Slug
	}
	if sc.CreatedAt != nil {
		category.CreatedAt = *sc.CreatedAt
	}
	if sc.UpdatedAt != nil {
		category.UpdatedAt = *sc.UpdatedAt
	}

	return category
}
