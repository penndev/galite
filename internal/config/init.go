package config

import (
	"log"

	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/pkg/logger"
)

type CfgMode int

const (
	ModeDev CfgMode = iota
	ModeProd
)

type Database struct {
	Dsn string // 数据库连接信息
}

type Cache struct {
	Dsn string // 缓存连接信息
}

type Badger struct {
	Dsn string // 缓存连接信息
}

type Logger struct {
	Dsn   string          // 日志处理方式，可以写入文件或者队列
	Level logger.LogLevel // 数据库日志级别
}

type Config struct {
	Listen string  // Http Listen Addr
	Secret string  // 加密中的盐
	Mode   CfgMode // 配置的app模式

	Logger   Logger   // 日志的处理
	Database Database // 数据库配置信息
	Cache    Cache    // 缓存配置信息
	Badger   Badger   // 缓存配置信息
}

var cfg Config

func Mode() CfgMode {
	return cfg.Mode
}

func Listen() string {
	return cfg.Listen
}

func Secret() []byte {
	return []byte(cfg.Secret)
}

// 执行程序清理
func Defer() {

}

// 包应该被引入就立即执行初始化操作
// 方便阅读所以在main中显式执行
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
	// defer lib.Badger.Close()
	return nil
}

func InitRedis() error {
	return lib.InitRedis(cfg.Cache.Dsn)
}
