package wafcdn

import (
	"github.com/penndev/galite/model/orm"
	"github.com/penndev/galite/model/wafcdn"
	"gorm.io/gorm"
)

// 获取用户列表
type bindSiteParam struct {
	orm.BindListParam
	Remark string `form:"remark" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindSiteParam) Param() *wafcdn.Site {
	m := &wafcdn.Site{}
	if b.Remark != "" {
		m.Remark = "%" + b.Remark + "%"
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m).Preload("Domains")
	}
	m.Bind(m, w, b)
	return m
}

// 获取用户列表
type bindDomainParam struct {
	orm.BindListParam
	Name string `form:"name" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindDomainParam) Param() *wafcdn.Domain {
	m := &wafcdn.Domain{}
	if b.Name != "" {
		m.Name = "%" + b.Name + "%"
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m)
	}
	m.Bind(m, w, b)
	return m
}
