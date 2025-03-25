package model

import (
	"github.com/penndev/galite/model/orm"
	"github.com/penndev/galite/model/system"
	"github.com/penndev/galite/model/wafcdn"
)

// 注册表结构。
func Migration() {
	// 后台系统用户权限管理模块
	orm.DB.AutoMigrate(&system.SysAdmin{})
	orm.DB.AutoMigrate(&system.SysRole{})
	orm.DB.AutoMigrate(&system.SysAccessLog{})

	// WAFCDN模块
	orm.DB.AutoMigrate(&wafcdn.Cache{})
	orm.DB.AutoMigrate(&wafcdn.Site{})
	orm.DB.AutoMigrate(&wafcdn.Domain{})

}
