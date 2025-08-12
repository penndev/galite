package lib

import (
	"log"

	"github.com/penndev/galite/pkg/cache"
	"github.com/penndev/galite/pkg/orm"
	"gorm.io/gorm"
)

// 设置cache实例
// 通过 lib.SetCache() 设置
// 通过 lib.Cache 获取
var Cache cache.Interface

func SetCache(c cache.Interface) {
	if c == nil {
		log.Panic("cache is nil")
	}
	Cache = c
}

// GormDB 数据库连接实例
// 将多个数据库连接进行抽象，减少不同数据库配置依赖
// 不论什么数据库最终提供给gram的都为 gorm.Dialector
// 通过 lib.SetGorm() 设置
// 通过 lib.GormDB 获取
var GormDB *gorm.DB

func SetGorm(db *gorm.DB) {
	if db == nil {
		log.Panic("db is nil")
	}
	orm.SetDB(db)
	GormDB = db
}

// NginxManager 管理nginx的结构体
// 通过 lib.SetNginx() 设置
// 通过 lib.Nginx 获取
var Nginx NginxManager

func SetNginx(binary, prefix, output string) {
	Nginx = NginxManager{
		Binary:     binary,
		Prefix:     prefix,
		OutputFile: output,
	}

	if err := Nginx.Reload(); err != nil {
		if err := Nginx.Start(); err != nil {
			log.Panic(err)
		}
	}
}
