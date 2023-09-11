package models

type Session struct {
	ID           int64  `db:"id" json:"id"`
	UserID       int64  `db:"user_id" json:"user_id"`
	RefreshToken string `db:"refresh_token" json:"refresh_token"`
	ClientIP     string `db:"client_ip" json:"client_ip"`
	CreatedAt    string `db:"created_at" json:"created_at"`
	ExpiredAt    string `db:"expired_at" json:"expired_at"`
}
