package respond

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"html/template"
	"net/http"
	"os"
)

type responseType int

const (
	jsonResponse responseType = iota
	htmlResponse
	textResponse
	fileResponse
	base64Response
)

// Response 自定义响应
type Response struct {
	out          http.ResponseWriter // 输出目标
	responseType responseType        // 默认值是 0
	headers      map[string]string   // http headers
	status       int                 // 响应状态码，默认 200
	data         interface{}         // 响应的数据
	htmlTemplate *template.Template
}

func New(out http.ResponseWriter) *Response {
	return &Response{
		out:          out,
		responseType: jsonResponse,
		headers:      map[string]string{},
		status:       http.StatusOK,
	}
}

// Header 添加一个 header
func (r *Response) Header(key string, value string) *Response {
	if r.headers == nil {
		r.headers = make(map[string]string)
	}
	r.headers[key] = value
	return r
}

// Status 设置状态码
func (r *Response) Status(status int) *Response {
	r.status = status
	return r
}

// JSON 响应JSON数据
func (r *Response) JSON(data interface{}) error {
	r.data = data
	r.responseType = jsonResponse
	return r.write()
}

// JSON 执行HTML模版
func (r *Response) HTML(htmlTemplate *template.Template, data interface{}) error {
	r.htmlTemplate = htmlTemplate
	r.data = data
	r.responseType = htmlResponse
	return r.write()
}

// JSON 响应Text数据
func (r *Response) Text(data string) error {
	r.data = data
	r.responseType = textResponse
	return r.write()
}

// File 响应文件
func (r *Response) File(filepath string) error {
	r.data = filepath
	r.responseType = fileResponse
	return r.write()
}

// Base64 响应base64数据
func (r *Response) Base64(data []byte) error {
	r.data = data
	r.responseType = base64Response
	return r.write()
}

// write 写入响应
func (r *Response) write() error {
	switch r.responseType {
	case htmlResponse:
		r.Header("Content-Type", "text/html")
		r.writeHeader()
		return r.htmlHandler()
	case textResponse:
		r.Header("Content-Type", "text/plain")
		r.writeHeader()
		return r.textHandler()
	case fileResponse:
		// TODO: 待定
		// r.Header("Content-Type", "")
		r.writeHeader()
		return r.fileHandler()
	case base64Response:
		// TODO: 待定
		// r.Header("Content-Type", "")
		r.writeHeader()
		return r.base64Handler()
	default:
		r.Header("Content-Type", "application/json")
		r.writeHeader()
		return r.jsonHandler()
	}
}

// writeHeader 添加响应头
func (r *Response) writeHeader() *Response {
	for k, v := range r.headers {
		r.out.Header().Add(k, v)
	}
	return r
}

// jsonHandler 响应 JSON 数据
func (r *Response) jsonHandler() error {
	r.out.WriteHeader(r.status)
	return json.NewEncoder(r.out).Encode(r.data)
}

// htmlHandler 响应 HTML 模版
func (r *Response) htmlHandler() error {
	var out = new(bytes.Buffer)
	err := r.htmlTemplate.Execute(out, r.data)
	if err != nil {
		return err
	}

	r.out.WriteHeader(r.status)
	_, err = out.WriteTo(r.out)
	return err
}

// textHandler 响应文本
func (r *Response) textHandler() error {
	text, _ := r.data.(string)
	r.out.WriteHeader(r.status)
	_, err := r.out.Write([]byte(text))
	return err
}

// fileHandler 响应文件
func (r *Response) fileHandler() error {
	filepath, ok := r.data.(string)
	if !ok {
		return errors.New("not file path")
	}

	buf, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	r.out.WriteHeader(r.status)
	_, err = r.out.Write(buf)
	return err
}

// base64Handler 响应 base64 文件格式
func (r *Response) base64Handler() error {
	dataBytes, ok := r.data.([]byte)
	if !ok {
		return errors.New("not bytes")
	}

	var result string
	mimeType := http.DetectContentType(dataBytes)

	switch mimeType {
	case "image/jpeg":
		result += "data:image/jpeg;base64,"
	case "image/png":
		result += "data:image/png;base64,"
		// TODO: 其他类型
		// ...
	}
	// w.Header().Add("Content-Type", mimeType)
	result += base64.StdEncoding.EncodeToString(dataBytes)
	r.out.WriteHeader(r.status)
	_, err := r.out.Write([]byte(result))
	return err
}
