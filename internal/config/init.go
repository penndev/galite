package config

import (
	"log"
	"time"

	"github.com/penndev/galite/internal/lib"
)

var StartTime time.Time = time.Now()

// 包应该被引入就立即执行初始化操作
// 方便阅读所以在main中显式执行
// 某些场景下有些业务不需要
func Init() {
	if err := InitEnv(); err != nil {
		log.Panic(err)
	}
	if err := InitLogger(); err != nil {
		log.Panic(err)
	}
	if err := InitGorm(); err != nil {
		log.Panic(err)
	}
	if err := InitRedis(); err != nil {
		log.Panic(err)
	}
	if err := InitBadger(); err != nil {
		log.Panic(err)
	}
}

func InitBadger() error {
	// 初始化keyVal文件存储系统
	if err := lib.InitBadger(cfg.Badger.Dsn); err != nil {
		return err
	}
	closeList = append(closeList, lib.Badger)
	return nil
}

func InitRedis() error {
	if err := lib.InitRedis(cfg.Cache.Dsn); err != nil {
		return err
	}
	closeList = append(closeList, lib.Redis)
	return nil
}
