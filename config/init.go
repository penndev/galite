package config

import (
	"errors"
	"log"
	"os"
)

// 包应该被引入就立即执行初始化操作
func init() {
	if err := InitEnv(); err != nil {
		log.Panic(err)
	}
	if err := InitZap(); err != nil {
		log.Panic(err)
	}
	if err := InitGorm(); err != nil {
		log.Panic(err)
	}
}

func CacheURL() string {
	if os.Getenv("CACHE_URL") == "" {
		log.Panic(errors.New("env CACHE_URL [redis://default:@127.0.0.1:6379/1] not found"))
	}
	return os.Getenv("CACHE_URL")
}
