package errs

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lightsaid/blogs/dbrepo"
	"github.com/mattn/go-sqlite3"
)

// 在此统一处理系统公共错误，如何：数据库错误

// HandleSQLError 处理SQL错误, 返回 *AppError
func HandleSQLError(err error) *AppError {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			var errMsg = sqliteErr.Error()
			if strings.Contains(errMsg, "tags.title") {
				return ErrRecordExists.AsException(err, "标签已存在")
			}
			if strings.Contains(errMsg, "users.email") {
				return ErrRecordExists.AsException(err, "邮箱已存在")
			}
			if strings.Contains(errMsg, "category.title") {
				return ErrRecordExists.AsException(err, "分类已存在")
			}

			// TODO: 这里不一定是唯一约束冲突，还有可能是非空约束
			// return ErrRecordExists.AsError(err)
		}

		return ErrInternalServer.AsException(err)
	}

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, dbrepo.ErrNotRowsAffected) {
		return ErrNotFound.AsException(err)
	}

	// TODO: 更多错误处理

	return ErrInternalServer.AsException(err)
}
