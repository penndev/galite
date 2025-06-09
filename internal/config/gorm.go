package config

import (
	"log"
	"os"
	"time"

	loggerNal "github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/model"
	loggerPkg "github.com/penndev/galite/pkg/logger"
	"gorm.io/gorm/logger"
)

// 设置gorm数据库的实例
func InitGorm() error {
	gormDial, err := model.Dial(cfg.Database.Dsn)
	if err != nil {
		return err
	}
	if cfg.Mode == ModeDev {
		gormLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  cfg.Logger.Level.GormLevel(),
			IgnoreRecordNotFoundError: true,
			Colorful:                  true, //颜色控制
		})
		model.InitGorm(gormDial, gormLogger) // 初始化
		model.Migration()                    // 表自动迁移
	} else {
		gormLogger := loggerPkg.ZapGormLogger(
			loggerNal.GormZapLogger,
			200*time.Millisecond,
		)
		model.InitGorm(gormDial, gormLogger) // 初始化
	}
	return nil
}
