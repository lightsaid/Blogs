package errs

import "net/http"

// 定义公共错误码
var (
	ErrSuccess               = NewAppError("请求成功", http.StatusOK)                           // 200 表示请求成功
	ErrBadRequest            = NewAppError("入参错误", http.StatusBadRequest)                   // 400  表示客户端发送的请求有误，服务器无法理解。
	ErrUnauthorized          = NewAppError("验证失败", http.StatusUnauthorized)                 // 401 表示客户端未经身份验证或身份验证失败。
	ErrForbidden             = NewAppError("禁止访问", http.StatusForbidden)                    // 403 表示客户端未经授权访问资源。
	ErrNotFound              = NewAppError("没有记录", http.StatusNotFound)                     // 404 表示请求的资源不存在
	ErrRequestTimeout        = NewAppError("请求超时", http.StatusRequestTimeout)               // 408 表示客户端请求超时。
	ErrRecordExists          = NewAppError("记录已存在", http.StatusConflict)                    // 409 表示记录已存在
	ErrRequestEntityTooLarge = NewAppError("入参过大", http.StatusRequestEntityTooLarge)        // 413 示客户端请求体太大，服务器无法处理。
	ErrRowsAffected          = NewAppError("没有操作", http.StatusTeapot)                       // 418 啥也没做
	ErrUnprocessableEntity   = NewAppError("入参有误", http.StatusUnprocessableEntity)          // 422 表示客户端发送的请求格式正确，但服务器无法处理，通常是由于请求体中缺少必要的字段或参数错误。
	ErrTooManyRequests       = NewAppError("请求繁忙", http.StatusTooManyRequests)              // 429 表示客户端发送的请求过多，超出了限制。
	ErrInternalServer        = NewAppError("务器内部错误，请稍后重试！", http.StatusInternalServerError) // 500  表示服务器内部错误。
)
