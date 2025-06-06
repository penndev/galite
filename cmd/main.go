package main

import (
	"log"
	"net/http"
	"time"

	"github.com/penndev/galite/config"
	"github.com/penndev/galite/route"
)

func main() {

	// 初始化各种组件
	//config.init()

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
