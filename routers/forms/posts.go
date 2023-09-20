package forms

import "github.com/lightsaid/blogs/dbrepo"

type NewPostsRequest struct {
	AuthorID    int64   `json:"author_id" validate:"required,gte=1"`
	Title       string  `json:"title" validate:"required"`
	Content     string  `json:"content" validate:"required"`
	Keyword     string  `json:"keyword"`
	Slug        string  `json:"slug"`
	Abstract    string  `json:"abstract"`
	CoverID     int64   `json:"cover_id"`
	TagIDs      []int64 `json:"tag_ids,omitempty"`
	CategoryIDs []int64 `json:"category_ids,omitempty"`
}

type UpdatePostsRequest struct {
	ID          int64   `json:"id" validate:"required,gte=1"`
	Title       string  `json:"title" validate:"required"`
	Content     string  `json:"content" validate:"required"`
	Keyword     string  `json:"keyword"`
	Slug        string  `json:"slug"`
	Abstract    string  `json:"abstract"`
	CoverID     int64   `json:"cover_id"`
	TagIDs      []int64 `json:"tag_ids,omitempty"`
	CategoryIDs []int64 `json:"category_ids,omitempty"`
}

type ListRequest struct {
	Meta dbrepo.Metadata `json:"meta"`
	List interface{}     `json:"list"`
}
