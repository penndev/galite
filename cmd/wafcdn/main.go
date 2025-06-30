package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/penndev/galite/internal"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/lib"
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
	if err := config.InitRedis(); err != nil {
		log.Panic(err)
	}
	// 跟随程序启动openresty
	lib.Nginx = lib.NginxManager{
		Binary:     os.Getenv("NGINX_BINARY"),
		Prefix:     os.Getenv("NGINX_PREFIX"),
		OutputFile: os.Getenv("NGINX_OUTPUT"),
	}
	if err := lib.Nginx.Reload(); err != nil {
		if err := lib.Nginx.Start(); err != nil {
			log.Panic(err)
		}
	}
	defer lib.Nginx.Stop()
	// 启动Http服务器 高性能版
	route := internal.InitRoute()                    // 挂载默认接口路由
	route.Static("/-admin", "./dist")                // 挂载后台管理UI
	route.GroupPush("/@wafcdn", wafcdn.InitApiRoute) // openresty通讯接口
	httpServe := &http.Server{
		Addr:           config.Listen(),
		Handler:        route,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("[HTTP] Listening on (%s): http://%s\n", time.Since(config.StartTime), config.Listen())
	log.Panic(httpServe.ListenAndServe())
}
