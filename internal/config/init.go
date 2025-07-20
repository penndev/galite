package config

import (
	"log"
	"time"
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
	if err := InitCache(); err != nil {
		log.Panic(err)
	}
}
