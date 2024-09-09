package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// 创建一个切割日志
func NewLogger(logFileName string) (*lumberjack.Logger, error) {
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logFileName,
		MaxSize:    1024,
		MaxBackups: 30,
	}
	_, err := lumberjackLogger.Write([]byte(fmt.Sprintf("%s Create lumberjack logger file[%s] \n", time.Now(), logFileName)))
	return lumberjackLogger, err
}

func initLogger() {

	if gin.Mode() == gin.ReleaseMode {
		// 初始化日志
		if os.Getenv("LOGGER_FILE") == "" {
			log.Panic(errors.New("LOGGER_FILE .env not found"))
		}

		lgf, err := NewLogger(os.Getenv("LOGGER_FILE"))
		if err != nil {
			log.Panic(err)
		}
		var level zapcore.Level
		switch os.Getenv("LOGGER_LEVEL") {
		case "debug":
			level = zapcore.DebugLevel
		case "info":
			level = zapcore.InfoLevel
		case "warn":
			level = zapcore.WarnLevel
		case "error":
			level = zapcore.ErrorLevel
		default:
			level = zapcore.InfoLevel
		}
		core := zapcore.NewCore(zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()), zapcore.AddSync(lgf), level)
		Logger = zap.New(core)
	} else {
		var err error
		Logger, err = zap.NewDevelopment()
		if err != nil {
			log.Panic(err)
		}
	}

	Logger.Info("WGA Starting")
}
