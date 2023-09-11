package models

type User struct {
	Model
	Email       string  `db:"email" json:"email"`
	Password    string  `db:"password" json:"password"`
	UserName    string  `db:"username" json:"username"`
	Avatar      string  `db:"avatar" json:"avatar"`
	Role        int     `db:"role" json:"role"`
	ActivatedAt *string `db:"activated_at" json:"activated_at"`
}
