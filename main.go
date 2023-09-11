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
)

func main() {
	// 配置文件存放路径
	var conf string
	flag.StringVar(&conf, "conf", "config.json", "配置文件路径")
	flag.Parse()

	if err := config.LoadingToAppConf(conf); err != nil {
		log.Fatal("加载配置错误: ", err)
	}

	addr := fmt.Sprintf("%s:%d", config.AppConf.Server.Host, config.AppConf.Server.Port)

	server := http.Server{
		Addr:         addr,
		Handler:      nil,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  time.Minute,
	}

	errCh := make(chan error)

	go func() {
		done := make(chan os.Signal, 1)
		signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
		<-done

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		errCh <- server.Shutdown(ctx)
	}()

	slog.Info("server start on " + addr)
	if err := server.ListenAndServe(); err != nil {
		slog.Error("server stopped", slog.String("err", err.Error()))
	}

	err := <-errCh
	if err != nil {
		slog.Error("server stopped", slog.String("errCh", err.Error()))
	}
}
