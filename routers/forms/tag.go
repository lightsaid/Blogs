package forms

type AddTagRequest struct {
	Title string `json:"title" validate:"required"`
}

type UpdateTagRequest struct {
	ID    int64  `json:"id" validate:"required,min=1"`
	Title string `json:"title" validate:"required"`
}
