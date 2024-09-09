package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model"
	"github.com/penndev/galite/route"
)

func main() {
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
