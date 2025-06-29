package config

import (
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
