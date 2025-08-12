package model

import (
	"github.com/penndev/galite/internal/admin/model/system"
	"gorm.io/gorm"
)

func Migration(db *gorm.DB) {
	db.AutoMigrate(&system.SysAdmin{})
	db.AutoMigrate(&system.SysRole{})
	db.AutoMigrate(&system.SysAccessLog{})
}
