package models

// Model 数据库表公共字段，提供给其他model使用
type Model struct {
	ID        int64   `db:"id" json:"id"`
	CreatedAt string  `db:"created_at" json:"created_at"`
	UpdatedAt string  `db:"updated_at" json:"updated_at"`
	DeletedAt *string `db:"deleted_at" json:"deleted_at"`
}
