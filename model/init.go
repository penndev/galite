package model

import (
	"github.com/penndev/wga/config"
	"github.com/penndev/wga/model/system"
)

// 注册表结构。
func Migration() {
	config.DB.AutoMigrate(system.SysUser{})
}
