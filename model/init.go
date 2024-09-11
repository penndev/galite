package model

import (
	"log"

	"github.com/penndev/galite/model/orm"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Database 在中间件中初始化mysql链接
func InitGorm(dialector gorm.Dialector, Logger logger.Interface) {
	dataBase, err := gorm.Open(dialector, &gorm.Config{
		Logger: Logger, // 重写日志
		// DisableForeignKeyConstraintWhenMigrating: true, // 禁止物理外键约束
	})
	if err != nil {
		log.Panic(err)
	}

	// sqlDB, err := dataBase.DB()
	// if err != nil {
	// 	log.Panic(err)
	// }

	// 最大空闲数
	// sqlDB.SetMaxIdleConns(DBMaxIdleConns)
	// 最大连接数
	// sqlDB.SetMaxOpenConns(DBMaxOpenConns)
	// 最大存活时间
	// sqlDB.SetConnMaxLifetime(time.Hour)

	orm.DB = dataBase
}
