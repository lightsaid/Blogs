package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/lightsaid/blogs/config"
	"github.com/lightsaid/blogs/initializers"
	"github.com/lightsaid/blogs/routers"
)

func main() {
	// 配置文件存放路径
	var conf string
	flag.StringVar(&conf, "conf", "config.json", "配置文件路径")
	flag.Parse()

	// 加载JSON配置到 config.AppConf 中
	err := config.LoadingToAppConf(conf)
	logFatal(err, "加载配置错误")

	// 初始化logger（slog）
	initializers.InitLogger()

	// 连接数据（SQLite）
	db, err := initializers.InitSQLite()
	logFatal(err, "连接SQLite错误")
	defer func() {
		slog.Info("正在释放SQLite资源...")
		if err := db.Close(); err != nil {
			slog.Error("关闭SQLite错误", slog.String("error", err.Error()))
		}
	}()

	// 初始化验证器, 并设置为全局的
	trans, validate := initializers.InitValidator()
	config.Trans = trans
	config.Validate = validate

	// 初始化路由 (routers)
	mux := routers.NewRouter(db)

	addr := fmt.Sprintf("%s:%d", config.AppConf.Server.Host, config.AppConf.Server.Port)

	// 创建 web 服务
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	errCh := make(chan error)

	// 监听系统信号（ctrl+c）并关机
	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done
		close(done)

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		errCh <- server.Shutdown(ctx)
	}()

	// 启动 web 服务
	slog.Info("server start on " + addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server stopped", slog.String("err", err.Error()))
	}

	err = <-errCh
	if err != nil {
		slog.Error("server stopped", slog.String("errCh", err.Error()))
	}
	close(errCh)
}

// 打印致命错误/退出程序
func logFatal(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}
}
