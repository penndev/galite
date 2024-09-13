package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/cache"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model"
	"github.com/penndev/galite/route"
	"gorm.io/gorm/logger"
)

func main() {

	config.Init()
	// 初始化redis
	cache.InitRedis(config.CacheRedisURL)
	// 初始化数据库
	if config.Mode == config.ModeDEV {
		model.InitGorm(config.GormDial, logger.Default)
	} else {
		model.InitGorm(config.GormDial, config.GormZapLogger)
	}
	model.Migration() // 表自动迁移

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
