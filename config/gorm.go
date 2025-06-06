package config

import (
	"github.com/penndev/galite/model"
	"gorm.io/gorm/logger"
)

// 初始化数据库连接
// 开发模式使用默认的gorm logger
// 生产模式使用配置的zap gorm logger
func InitGorm() error {
	gormDial, err := GormDial()
	if err != nil {
		return err
	}
	var gormLogger logger.Interface
	if Mode == ModeDEV {
		gormLogger, err = GormDefaultLogger()
	} else {
		gormLogger, err = GormZapLogger()
	}
	if err != nil {
		return err
	}

	model.InitGorm(gormDial, gormLogger) // 初始化
	model.Migration()                    // 表自动迁移
	return nil
}
