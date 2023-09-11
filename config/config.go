package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"time"
)

// config package 存放配置和全局变量

// 全局常量
const (
	EnvDev  = "dev"  // dev 开发模式
	EnvProd = "prod" // prod 生成模式
)

// 解析配置文件到此变量，提供给全局使用
var AppConf = new(Config)

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
		SecretKey   string
		TokenExpire string
	}

	Cookie struct {
		SecretKey string
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
