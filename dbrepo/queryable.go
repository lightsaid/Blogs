package dbrepo

import (
	"context"
	"database/sql"
)

// Queryable 将 sqlx.DB 和 sqlx.Tx 常用公共方法提取为一个接口，
// 为了在实现事务时，sqlx.Tx 可以重用 “使用sqlx.DB操作数据库” 的方法
type Queryable interface {
	Get(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...any) (sql.Result, error)
	Query(query string, args ...any) (*sql.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
	QueryRow(query string, args ...any) *sql.Row
	Prepare(query string) (*sql.Stmt, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
