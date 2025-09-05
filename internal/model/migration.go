package model

import (
	modelAdmin "github.com/penndev/galite/internal/admin/model"
	"github.com/penndev/galite/internal/lib"
	modelWafCdn "github.com/penndev/galite/internal/wafcdn/model"
)

// 注册表结构。
func Migration() {
	// 后台系统用户权限管理模块
	modelAdmin.Migration(lib.GormDB)
	// 注册wafcdn的数据库
	modelWafCdn.Migration(lib.GormDB)
}
