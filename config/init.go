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
var JWTSecret []byte

// 日志实例
var Logger *zap.Logger

// Redis实例
var Redis *redis.Client

// DB 数据库链接单例
var DB *gorm.DB

func Defer() {
	Logger.Sync()
}

// Init 初始化配置项
func Init() {
	// 加载env
	godotenv.Load()
	// 设置gin.mode
	gin.SetMode(os.Getenv("GIN_MODE"))
	JWTSecret = []byte(os.Getenv("SECRET"))

	initLoger()
	initGorm()
	initRedis()
}
