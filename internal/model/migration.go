package model

import (
	"github.com/penndev/galite/internal/admin/model/system"
	modelWafCdn "github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/orm"
)

// 注册表结构。
func Migration() {
	// 后台系统用户权限管理模块
	orm.DB.AutoMigrate(&system.SysAdmin{})
	orm.DB.AutoMigrate(&system.SysRole{})
	orm.DB.AutoMigrate(&system.SysAccessLog{})
	// 注册wafcdn的数据库
	modelWafCdn.Migration(orm.DB)
}

// 进行后台任务处理
func Runner() {
	go modelWafCdn.CacheDeleteRunner()
}
