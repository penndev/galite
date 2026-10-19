package model

import "gorm.io/gorm"

func Migration(db *gorm.DB) {
	db.AutoMigrate(&Log{})
	db.AutoMigrate(&Cache{})
	db.AutoMigrate(&Site{})
	db.AutoMigrate(&Cert{})
	db.AutoMigrate(&Domain{})
	db.AutoMigrate(&CacheDelete{})
}
