package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var Mode = "dev"

// 模式： dev - 控制台输出优化
// 模式： prod - zap日志输出内容
const ModeDEV = "dev"
const ModePROD = "prod"

// http listen addr
var Listen = ":8000"

// 处理 jwt secret
var JWTSecret []byte

// // 给gin的 prod模式下的日志输出文件
// var GinZapLogger *zap.Logger

// 给用户使用的 log文件。
var Logger *zap.Logger

// 初始化所有的env变量。避免不存在env引起的程序异常
func Init() {
	godotenv.Load(".env")
	// 处理app启动参数
	var err error
	if os.Getenv("APP_MODE") != "dev" {
		Mode = "prod"
		if os.Getenv("APP_LOGGER_FILE") == "" {
			log.Panic(errors.New("env APP_LOGGER_FILE [gin.log] not found"))
		}
		// GinZapLogger 只收集prod模式下的日志
		Logger, err = ZapLogger(os.Getenv("APP_LOGGER_FILE"), ParseLogLevel(os.Getenv("APP_LOGGER_LEVEL")), 1024, 30)
		if err != nil {
			log.Panic(err)
		}
	} else {
		Logger, err = zap.NewDevelopment()
		if err != nil {
			log.Panic(err)
		}
	}

	if os.Getenv("APP_LISTEN") != "" {
		Listen = os.Getenv("APP_LISTEN")
	}

	if os.Getenv("APP_SECRET") == "" {
		log.Panic(errors.New("env APP_SECRET [secret] not found"))
	}
	JWTSecret = []byte(os.Getenv("APP_SECRET"))
}

func CacheURL() string {
	if os.Getenv("CACHE_URL") == "" {
		log.Panic(errors.New("env CACHE_URL [redis://default:@127.0.0.1:6379/1] not found"))
	}
	return os.Getenv("CACHE_URL")
}
