package forms

type AddUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	UserName string `json:"username" validate:"required,min=2"`
	Avatar   string `json:"avatar"`
}

type UpdateUserRequest struct {
	ID       int64  `json:"id" validate:"required,gte=1"`
	UserName string `json:"username" validate:"required,min=2"`
	Avatar   string `json:"avatar"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	ClientIP string
}

type LoginResponse struct {
	ID           int64  `json:"id"`
	UserName     string `json:"username"`
	Avatar       string `json:"avatar"`
	RefreshToken string `json:"refresh_token"`
	LoginAt      string `json:"login_at"`
	AccessToken  string `json:"-"`
}
