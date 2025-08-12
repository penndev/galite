package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/internal"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/model"
	"github.com/penndev/galite/internal/wafcdn"
	"github.com/penndev/galite/pkg/util"
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
	model.Migration() // 每次启动都需要同步数据库结构 表自动迁移
	if err := config.InitCache(); err != nil {
		log.Panic(err)
	}
	model.Runner() // 异步执行定时任务

	//
	lib.SetNginx(
		util.GetEnv("NGINX_BINARY", "openresty"),
		util.GetEnv("NGINX_PREFIX", "./"),
		util.GetEnv("NGINX_OUTPUT", ""),
	)

	// 启动Http服务器 高性能版
	route := internal.InitRoute()                    // 挂载默认接口路由
	route.Static("/-admin", "./dist")                // 挂载后台管理UI
	route.GroupPush("/@wafcdn", wafcdn.InitApiRoute) // openresty通讯接口
	httpServe := &http.Server{
		Addr:           config.Listen(),
		Handler:        route,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 15,
	}
	log.Printf("[HTTP] Listening on (%s): http://%s\n", time.Since(config.StartTime), config.Listen())
	log.Panic(httpServe.ListenAndServe())
}
