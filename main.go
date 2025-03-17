package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/cache"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model"
	"github.com/penndev/galite/route"
)

func main() {

	config.Init()
	// 初始化redis
	cache.InitRedis(config.CacheURL())

	// 数据库处理
	model.InitGorm(config.GormDial(), config.GormLogger()) // 初始化
	model.Migration()                                      // 表自动迁移

	// 启动Http服务器 高性能版
	httpServe := &http.Server{
		Addr:           config.Listen,
		Handler:        route.Init(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening Serve http://%s \n", config.Listen)
	log.Panic(httpServe.ListenAndServe())
}
