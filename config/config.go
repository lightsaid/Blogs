package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/lightsaid/blogs/token"
)

// config package 存放配置和全局变量

// 全局常量
const (
	EnvDev  = "dev"  // dev 开发模式
	EnvProd = "prod" // prod 生产模式
)

// 全局变量
var (
	// 解析配置文件到此变量，提供给全局使用
	AppConf = new(Config)

	// 验证器的翻译器
	Trans ut.Translator

	// 验证器
	Validate *validator.Validate

	// Token 管理者
	TokenMaker token.Maker
)

type Config struct {
	Server struct {
		Host    string
		Port    int
		Env     string
		Version string
	}

	Database struct {
		Source       string
		MaxOpenConns int
		MaxIdleConns int
		MaxIdleTime  string
	}

	Token struct {
		SecretKey     string
		TokenExpire   string
		RefreshExpire string // Refresh Token 有效时长
	}

	Cookie struct {
		SecretKey string
		Name      string
	}

	Logger struct {
		Filename string // 存储位置
		Level    string // DEBUG|INFO|WARN|ERROR
	}
}

// LoadingToAppConf 加载配置文件到AppConf指针中
func LoadingToAppConf(filename string) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataBytes, AppConf)
}

// ParseDuration 将v解析为 time.Duration，解析失败使用默认值
func ParseDuration(v string, defaultVal time.Duration) time.Duration {
	dur, err := time.ParseDuration(v)
	if err != nil {
		slog.Error("解析time.Duration失败", slog.String("v", v))
		return defaultVal
	}
	return dur
}
