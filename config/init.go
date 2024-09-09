package config

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// JWT签名密钥
// var JWTSecret []byte

// 日志实例
var Logger *zap.Logger

// Redis实例
var Redis *redis.Client

// DB 数据库链接单例
var Gorm *gorm.DB

/**
 * 初始化所有配置
 * gin mode设置
 * redis连接池实例 config.Redis
 * zap日志 config.Logger
 * gorm连接池实例 config.DB
 */
func Init() {
	godotenv.Load()
	if os.Getenv("APP_MODE" == "dev") {

	}

	gin.SetMode(os.Getenv("GIN_MODE"))
	// JWTSecret = []byte(os.Getenv("SECRET"))

	initLogger()
	initGorm()
	initRedis()
}

// 配置项回收
func Defer() {
	Logger.Sync()
}
