package errs

import (
	"fmt"
)

// AppError 一个错误结构体，统一返回错误结构，包含客户消息、异常消息、HTTP 状态码
type AppError struct {
	message    string // 客户消息
	exception  error  // 异常消息，排查问题
	statusCode int    // HTTP 状态码
}

// NewAppError 创建一个 AppError
func NewAppError(msg string, statusCode int) *AppError {
	return &AppError{
		message:    msg,
		statusCode: statusCode,
	}
}

// Error 实现 error 接口
func (a *AppError) Error() string {
	err := fmt.Sprintf("code: %d, message: %s", a.statusCode, a.message)
	if a.exception != nil {
		err = fmt.Sprintf("%s, exception: %v", err, a.exception.Error())
	}
	return err
}

// Unwrap 解开，提供给 errors.Is 和 errors.As 使用
func (a *AppError) Unwrap() error {
	return a.exception
}

// StatusCode 返回 HTTP StatusCode
func (a *AppError) StatusCode() int {
	return a.statusCode
}

// Message 返回客户端可读错误信息
func (a *AppError) Message() string {
	return a.message
}

// AsMessage 修改客户端消息，并返回一个新的 AppError 指针
func (a *AppError) AsMessage(msg string) *AppError {
	return &AppError{
		statusCode: a.statusCode,
		message:    msg,
		exception:  a.exception,
	}
}

// AsException 添加/追加错误, 返回一个新的 AppError 指针
func (a *AppError) AsException(err error, msgs ...string) *AppError {
	var e error
	if a.exception == nil {
		e = fmt.Errorf("%w", err)
	} else {
		e = fmt.Errorf("%v | %w", a.exception, err)
	}
	newErr := &AppError{
		statusCode: a.statusCode,
		message:    a.message,
		exception:  e,
	}
	if len(msgs) > 0 {
		newErr.message = msgs[0]
	}
	return newErr
}
