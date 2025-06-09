package config

import (
	loggerNal "github.com/penndev/galite/internal/logger"
	loggerPkg "github.com/penndev/galite/pkg/logger"
	"go.uber.org/zap"
)

// 处理日志实例
func InitLogger() error {
	var err error
	var mainLog *zap.Logger
	if cfg.Mode == ModeDev {
		mainLog, err = zap.NewDevelopment()
		if err != nil {
			return err
		}
	} else {
		mainLog, err = loggerPkg.ZapFileLogger(
			cfg.Logger.Dsn,
			cfg.Logger.Level.ZapLevel(),
			1024,
			30,
		)
		if err != nil {
			return err
		}
	}
	loggerNal.GinZapLogger = mainLog
	loggerNal.GormZapLogger = mainLog
	loggerNal.ZapLogger = mainLog
	return nil
}
