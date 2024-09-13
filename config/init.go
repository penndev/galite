package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// 模式： dev - 控制台输出优化
// 模式： prod - zap日志输出内容
var Mode = "dev"

const ModeDEV = "dev"
const ModePROD = "prod"

// http listen addr
var Listen = ":8000"

// 处理 jwt secret
var JWTSecret = []byte("secret")

// redis连接URL
var CacheRedisURL = "redis://default:@127.0.0.1:6379/1"

// 数据库连接实例，
// 将多个数据库连接进行抽象，减少不同数据库配置依赖
// 不论什么数据库最终提供给gram的都为 gorm.Dialector
var GormDial gorm.Dialector
var GormZapLogger logger.Interface

// 给gin的 prod模式下的日志输出文件
var GinZapLogger *zap.Logger

// 给用户使用的 log文件。
var Logger *zap.Logger

// 初始化所有的env变量。避免不存在env引起的程序异常
func Init() {
	godotenv.Load(".env")
	// 处理app启动参数
	var err error
	if os.Getenv("APP_MODE") != "dev" {
		Mode = "prod"
		if os.Getenv("DB_LOGGER_FILE") == "" {
			os.Setenv("DB_LOGGER_FILE", "gorm.log")
		}
		if os.Getenv("APP_LOGGER_FILE") == "" {
			os.Setenv("APP_LOGGER_FILE", "gin.log")
		}
		// GinZapLogger 只收集prod模式下的日志
		GinZapLogger, err = ZapLogger(os.Getenv("APP_LOGGER_FILE"), ParseLogLevel(os.Getenv("DB_LOGGER_LEVEL")), 1024, 30)
		if err != nil {
			log.Panic(err)
		}
		Logger = GinZapLogger
	} else {
		GinZapLogger, err = zap.NewDevelopment()
		if err != nil {
			log.Panic(err)
		}
		Logger = GinZapLogger
	}

	if os.Getenv("APP_LISTEN") != "" {
		Listen = os.Getenv("APP_LISTEN")
	}
	if os.Getenv("APP_SECRET") != "" {
		JWTSecret = []byte(os.Getenv("APP_SECRET"))
	}
	if os.Getenv("CACHE_REDIS_URL") != "" {
		CacheRedisURL = os.Getenv("CACHE_REDIS_URL")
	}
	// 处理数据库
	switch {
	case strings.HasPrefix(os.Getenv("DB_URL"), "mariadb://"):
		GormDial = mysql.Open(strings.TrimPrefix(os.Getenv("DB_URL"), "mariadb://"))
	case strings.HasPrefix(os.Getenv("DB_URL"), "mysql://"):
		GormDial = mysql.Open(strings.TrimPrefix(os.Getenv("DB_URL"), "mysql://"))
	default:
		log.Panic(errors.New("env DB_URL err"))
	}
	GormZapLogger, err = GormLogger(os.Getenv("DB_LOGGER_FILE"), ParseLogLevel(os.Getenv("DB_LOGGER_LEVEL")), 1024, 30, 200)
	if err != nil {
		log.Panic(err)
	}

}
