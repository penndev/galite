package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/penndev/wga/config"
	"github.com/penndev/wga/model"
	"github.com/penndev/wga/route"
)

func main() {
	/**
	 * 初始化所有配置
	 * zap日志 config.Logger
	 * gorm连接池实例 config.DB
	 * redis连接池实例 config.Redis
	 */
	config.Init()
	defer config.Defer()
	model.Migration()

	// 加载所有路由
	handle := route.Init()

	// 启动Http服务器 高性能版
	httpAddr := os.Getenv("GIN_LISTEN")
	if httpAddr == "" {
		httpAddr = ":80"
	}

	httpServe := &http.Server{
		Addr:           httpAddr,
		Handler:        handle,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening and serving HTTP on [%s]\n", httpAddr)
	log.Panic(httpServe.ListenAndServe())
}
