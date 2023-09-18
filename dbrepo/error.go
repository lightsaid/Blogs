package dbrepo

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// HTTP 状态码参考 https://developer.mozilla.org/zh-CN/docs/Web/HTTP/Status/200

/*
	HTTP 状态码约定：
		参数基本验证不通过：400
		参数基本验证通过, 但是实体服务器无法处理：422
		参数基本验证通过, 但是出发数据库的唯一约束：409
		查询不到数据： 404，查找
		没有登录或者token格式错误，过期：401
		已登陆但是没有权限访问：403
		没有影响到数据行： 418
*/

// TODO: 待处理 errors.Is 检查错误

var errNotRowsAffected = errors.New("not rows affected")

// 定义预期错误
var (
	ErrNotFound     = NewDBError("没有记录", http.StatusNotFound)                  // 404
	ErrRecordExists = NewDBError("记录已存在", http.StatusConflict)                 // 409
	ErrRowsAffected = NewDBError("没有操作", http.StatusTeapot)                    // 418
	ErrDBInternal   = NewDBError("服务错误，请稍后重试", http.StatusInternalServerError) // 500
)

// DBError 数据库错误类型，为了将异常信息和提示信息分开
type DBError struct {
	err    error
	msg    string
	status int
}

// NewDBError 创建一个db错误
func NewDBError(msg string, status int) *DBError {
	return &DBError{
		msg:    msg,
		err:    errors.New(msg),
		status: status,
	}
}

// Error 实现错误接口，返回具体错误 err
func (de *DBError) Error() string {
	return fmt.Sprintf("%v", de.err)
}

// Message 返回错误 msg
func (de *DBError) Message() string {
	return de.msg
}

// Message 返回错误 msg
func (de *DBError) StatusCode() int {
	return de.status
}

// AsMessage 使用旧的DBError创建一个新的DBError，并重新定义msg
func (de *DBError) AsMessage(msg string) *DBError {
	return &DBError{
		msg:    msg,
		err:    de.err,
		status: de.status,
	}
}

// AsError 追加异常，返回一个新的 *DBError
func (de *DBError) AsError(err error, msgs ...string) *DBError {
	var newErr error
	var newMsg = de.msg
	if len(msgs) > 0 {
		newMsg = msgs[0]
		newErr = fmt.Errorf("%v: %w", newMsg, err)
	} else {
		newErr = fmt.Errorf("%w", err)
	}

	if de.err == nil {
		de.err = newErr
	} else {
		newErr = fmt.Errorf("%v | %w", de.err, newErr)
	}

	return &DBError{
		msg:    newMsg,
		err:    newErr,
		status: de.status,
	}
}

// Unwrap 解开，提供给 errors.Is 和 errors.As 使用
func (de *DBError) Unwrap() error {
	return de.err
}

// CheckError 检查错误，处理常见错误
func CheckError(err error) *DBError {
	var sqliteErr sqlite3.Error
	if errors.As(err, &sqliteErr) {
		if errors.Is(sqliteErr.Code, sqlite3.ErrConstraint) {
			var errMsg = sqliteErr.Error()
			if strings.Contains(errMsg, "tags.title") {
				return ErrRecordExists.AsError(err, "标签已存在")
			}
			if strings.Contains(errMsg, "users.email") {
				return ErrRecordExists.AsError(err, "邮箱已存在")
			}
			if strings.Contains(errMsg, "category.title") {
				return ErrRecordExists.AsError(err, "分类已存在")
			}
			// TODO: 这里不一定是唯一约束冲突，还有可能是非空约束
			// return ErrRecordExists.AsError(err)
		}

		return ErrDBInternal.AsError(err)
	}

	if errors.Is(err, sql.ErrNoRows) || errors.Is(err, errNotRowsAffected) {
		return ErrNotFound.AsError(err)
	}

	// TODO: 更多错误处理

	return ErrDBInternal.AsError(err)
}

func handleRowsAffected(rowsAff int64, err error) error {
	if err != nil {
		return err
	}

	if rowsAff == 0 {
		return errNotRowsAffected
	}

	return nil
}
