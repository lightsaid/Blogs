package forms

type AddCategoryRequest struct {
	Title string `json:"title" validate:"required"`
}

type UpdateCategoryRequest struct {
	ID    int64  `json:"id" validate:"required,gte=1"`
	Title string `json:"title" validate:"required"`
}
