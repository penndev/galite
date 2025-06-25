package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/penndev/galite/pkg/logger"
)

// 从env赋值为Cfg实例
func InitEnv() error {
	godotenv.Load(".env")

	// 设置app密钥
	secret := os.Getenv("APP_SECRET")
	if secret == "" {
		return errors.New("env APP_SECRET not found")
	}

	// 设置默认的监听端口
	listen := os.Getenv("APP_LISTEN")
	if listen == "" {
		listen = ":8000"
	}

	// 设置应用模式
	var mode CfgMode
	if os.Getenv("APP_MODE") != "dev" {
		mode = ModeProd
	} else {
		mode = ModeDev
	}

	//
	cfg = Config{
		Listen: listen,
		Secret: secret,
		Mode:   mode,
		Logger: Logger{
			Dsn:   os.Getenv("APP_LOGGER_FILE"),
			Level: logger.ParseLogLevel(os.Getenv("APP_LOGGER_LEVEL")),
		},
		Database: Database{
			Dsn: os.Getenv("DB_URL"),
		},
		Cache: Cache{
			Dsn: os.Getenv("CACHE_URL"),
		},
		Badger: Badger{
			Dsn: os.Getenv("BADGER_URL"),
		},
	}
	return nil
}
