package model

import "gorm.io/gorm"

// 注册表结构。
func Migration(db *gorm.DB) {
	// 后台系统用户权限管理模块
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Cache{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&CacheDelete{})
}
