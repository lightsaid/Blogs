package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"log/slog"
)

// ReadJSON 读取入参，绑定到 dst 上
func ReadJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	maxBytes := 2 << 20 // 2 MB
	// 限制请求体大小为2MB以内
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// 使用请求体创建一个解码器
	dec := json.NewDecoder(r.Body)

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		switch {
		case errors.As(err, &syntaxError):
			return errors.New("请输入JSON格式请求体")
		case errors.Is(err, io.EOF):
			return errors.New("请求体不能为空")
		case err.Error() == "http: request body too large":
			return errors.New("请求体过大")
		default:
			slog.ErrorContext(r.Context(), "未知错误，请检查参数类型是否匹配", slog.String("error", err.Error()))
			return errors.New("未知错误，请检查参数类型是否匹配")
		}
	}

	// 输出入参日志
	slog.InfoContext(
		r.Context(),
		"请求入参",
		slog.Any("json", dst),
	)

	return nil
}
