package system

import (
	"github.com/penndev/wga/model/sugar"
	"github.com/penndev/wga/model/system"
	"gorm.io/gorm"
)

// 用户登录请求体
type bindLoginInput struct {
	Username  string `binding:"required,min=4,max=64"`   // 用户名
	Password  string `binding:"required,min=6,max=64"`   // 密码
	Captcha   string `binding:"required,alphanum,len=4"` // 验证码
	CaptchaId string `binding:"required,uuid"`           // 验证码ID
}

// 获取用户列表
type bindSystemAdminParam struct {
	sugar.BindListParam
	Email string `form:"email" binding:"omitempty,min=4,max=64"`
}

// 处理列表请求数据。
func (b *bindSystemAdminParam) Param() *system.SysAdmin {
	m := &system.SysAdmin{
		Email: b.Email,
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m).Preload("AdminRole", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		})
	}
	m.Bind(m, w, b)
	return m
}

// 获取权限列表
type bindSystemRoleParam struct {
	sugar.BindListParam
	Status uint8 `form:"status" binding:"omitempty,min=0,max=1"`
}

func (b *bindSystemRoleParam) Param() *system.SysRole {
	m := &system.SysRole{
		Status: b.Status,
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m)
	}
	m.Bind(m, w, b)
	return m
}

// 获取访问日志列表
type bindSysAccessParam struct {
	sugar.BindListParam
}

// 处理列表请求数据。
func (b *bindSysAccessParam) Param() *system.SysAccessLog {
	m := &system.SysAccessLog{}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m).Preload("SysAdmin", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, email")
		})
	}
	m.Bind(m, w, b)
	return m
}
