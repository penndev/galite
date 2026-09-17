package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/penndev/galite/internal"
	"github.com/penndev/galite/internal/config"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/model"
	"github.com/penndev/galite/pkg/util"
	"github.com/penndev/gopkg/ipregion"
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
	if config.Mode() == config.ModeDev {
		model.Migration() // 表自动迁移
	}
	if err := config.InitCache(); err != nil {
		log.Panic(err)
	}
	ipregionDB := util.GetEnv("IPREGION_DB", "")
	if ipregionDB != "" {
		if _, err := os.Stat(ipregionDB); err == nil {
			searcher, err := ipregion.Open(ipregionDB)
			if err != nil {
				log.Panic(err)
			}
			lib.SetIPRegion(searcher)
		}
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
