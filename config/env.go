package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// 根据不同的mode设置不同的方式
// gin.mode 跟随本变量配置
// gorm.logger 方式跟随本变量
var Mode string

// 开发模式： dev - 控制台输出优化
// 开发模式： prod - zap日志输出内容
const ModeDEV = "dev"
const ModePROD = "prod"

// 服务器监听的端口
// gin 默认监听的端口为8000
var Listen = ":8000"

// Jwt加解密的盐
// 在分布式app上盐必须一致
var JWTSecret []byte

// 加载环境变量
// 设置app mode
// 设置app secret
// 设置app listen
func InitEnv() error {
	godotenv.Load(".env")

	// 设置应用模式
	if os.Getenv("APP_MODE") != "dev" {
		Mode = ModePROD
	} else {
		Mode = ModeDEV
	}

	// 设置app密钥
	if os.Getenv("APP_SECRET") == "" {
		return errors.New("env APP_SECRET not found")
	} else {
		JWTSecret = []byte(os.Getenv("APP_SECRET"))
	}

	// 设置默认的监听端口
	if os.Getenv("APP_LISTEN") != "" {
		Listen = os.Getenv("APP_LISTEN")
	}

	return nil
}
