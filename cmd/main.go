package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/internal"
	"github.com/penndev/galite/internal/config"
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
	if err := config.InitCache(); err != nil {
		log.Panic(err)
	}

	// 启动Http服务器 高性能版
	httpServe := &http.Server{
		Addr:           config.Listen(),
		Handler:        internal.InitRoute(),
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Listening Serve http://%s \n", config.Listen())
	log.Panic(httpServe.ListenAndServe())
}
