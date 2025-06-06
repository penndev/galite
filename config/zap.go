package config

import (
	"errors"
	"os"

	"go.uber.org/zap"
)

// Zap的logger实例
// 所有的日志输出均由此记录器处理
var Logger *zap.Logger

// 初始化Logger日志文件
// 生产模式，日志输出到文件或者三方流
//
// 开发模式，日志输出到控制台
func InitZap() error {
	if Mode == ModePROD {
		if os.Getenv("APP_LOGGER_FILE") == "" {
			return errors.New("env APP_LOGGER_FILE not found")
		}
		var err error
		Logger, err = ZapLogger(os.Getenv("APP_LOGGER_FILE"), ParseLogLevel(os.Getenv("APP_LOGGER_LEVEL")).toZapLoggerLevel(), 1024, 30)
		if err != nil {
			return err
		}
	} else {
		var err error
		Logger, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	}
	return nil
}
