package model

import (
	"github.com/penndev/galite/model/orm"
	"github.com/penndev/galite/model/system"
)

// 注册表结构。
func Migration() {
	orm.DB.AutoMigrate(&system.SysAdmin{})
	orm.DB.AutoMigrate(&system.SysRole{})
	orm.DB.AutoMigrate(&system.SysAccessLog{})
}
