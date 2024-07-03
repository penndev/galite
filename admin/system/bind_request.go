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
type bindSystemUserParam struct {
	sugar.BindListParam
	Name string `form:"name" binding:"omitempty,min=4,max=64"`
}

func (b *bindSystemUserParam) Param() *system.SysUser {
	m := &system.SysUser{
		Name: b.Name,
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m)
	}
	m.Bind(m, w, b)
	return m
}
