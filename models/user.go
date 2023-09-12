package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Model
	Email       string  `db:"email" json:"email"`
	Password    string  `db:"password" json:"password"`
	UserName    string  `db:"username" json:"username"`
	Avatar      string  `db:"avatar" json:"avatar"`
	Role        int     `db:"role" json:"role"`
	ActivatedAt *string `db:"activated_at" json:"activated_at"`
}

func (user *User) SetHashedPassword(plainText string) error {
	buf, err := bcrypt.GenerateFromPassword([]byte(plainText), 12)
	if err != nil {
		return err
	}
	user.Password = string(buf)
	return nil
}

func (user *User) MatchesPassword(plainText, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))
	return err == nil
}
