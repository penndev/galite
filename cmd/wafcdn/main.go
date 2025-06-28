package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/internal"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/wafcdn"
)

func main() {
	// 初始化各种组件
	defer config.Defer()
	if err := config.InitEnv(); err != nil {
		log.Panic(err)
	}
	if err := config.InitLogger(); err != nil {
		log.Panic(err)
	}
	if err := config.InitGorm(); err != nil {
		log.Panic(err)
	}
	if err := config.InitBadger(); err != nil {
		log.Panic(err)
	}

	route := internal.InitRoute()
	// 挂载接口路由
	wafcdn.InitApiRoute(route.Group("/@wafcdn"))
	// 挂载后台管理UI
	route.Static("/-admin", "./dist")

	// 启动Http服务器 高性能版
	httpServe := &http.Server{
		Addr:           config.Listen(),
		Handler:        route,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening Serve http://%s \n", config.Listen())
	log.Panic(httpServe.ListenAndServe())
}
