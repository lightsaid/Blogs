package initializers

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/lightsaid/blogs/config"
	"gopkg.in/natefinch/lumberjack.v2"
)

// LogHandler 定义一个LogHandler，方便扩展和自定义一些功能
type LogHandler struct {
	slog.Handler
}

// Handle 如果ctx有 request_id，记录到日志输出
func (lh *LogHandler) Handle(ctx context.Context, r slog.Record) error {
	// requestID 是有chi中间件创建和获取的
	requestID := middleware.GetReqID(ctx)

	if requestID != "" {
		r.AddAttrs(slog.String("request_id", requestID))
	}

	return lh.Handler.Handle(ctx, r)
}

// InitLogger 创建JSON格式日志，并设置为全局默认的slog
func InitLogger() {
	logLevel := toSlogLevel(config.AppConf.Logger.Level)

	var out io.Writer = os.Stdout
	var addSource bool
	if config.AppConf.Server.Env == config.EnvProd {
		out = &lumberjack.Logger{
			Filename:   config.AppConf.Logger.Filename,
			MaxSize:    40, // megabytes
			MaxBackups: 20,
			MaxAge:     30,   //days
			Compress:   true, // disabled by default
		}
		addSource = true
	}

	handler := slog.NewJSONHandler(
		out,
		&slog.HandlerOptions{
			AddSource: addSource,
			Level:     logLevel,
		})

	myHandler := LogHandler{Handler: handler}

	l := slog.New(&myHandler)

	// 设为全局默认 slog 实例
	slog.SetDefault(l)
}

// toSlogLevel 从字符串转换成slog.Leveler
func toSlogLevel(level string) slog.Leveler {
	switch level {
	case "DEBUG":
		return slog.LevelDebug
	case "INFO":
		return slog.LevelInfo
	case "WARN":
		return slog.LevelWarn
	case "ERROR":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

// LevelDebug Level = -4
// LevelInfo  Level = 0
// LevelWarn  Level = 4
// LevelError Level = 8
