package admin

import (
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/orm"
	"gorm.io/gorm"
)

// =========================================================== //
// 获取用户列表
type bindSiteParam struct {
	orm.BindListParam
	Remark string `form:"remark" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindSiteParam) Param() *model.Site {
	m := &model.Site{}
	if b.Remark != "" {
		m.Remark = "%" + b.Remark + "%"
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m).Preload("Domains")
	}
	m.Bind(m, w, b)
	return m
}

// =========================================================== //
// 获取用户列表
type bindDomainParam struct {
	orm.BindListParam
	Name string `form:"name" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindDomainParam) Param() *model.Domain {
	m := &model.Domain{}
	if b.Name != "" {
		m.Name = "%" + b.Name + "%"
	}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m)
	}
	m.Bind(m, w, b)
	return m
}

// =========================================================== //
// 获取用户列表
type bindCacheParam struct {
	orm.BindListParam
	SiteID uint   `form:"site_id"`
	Uri    string `form:"uri" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindCacheParam) Param() *model.Cache {
	m := &model.Cache{}
	w := func(orm *gorm.DB) *gorm.DB {
		if b.Uri != "" {
			orm = orm.Where("uri like ?", b.Uri+"%")
		}
		return orm
	}
	m.Bind(m, w, b)
	return m
}

// =========================================================== //
// 获取用户列表
type bindLogParam struct {
	orm.BindListParam
	Name string `form:"name" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindLogParam) Param() *model.Log {
	m := &model.Log{}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where(m)
	}
	m.Bind(m, w, b)
	return m
}

// =========================================================== //
// 获取用户列表
type bindCacheDeleteParam struct {
	orm.BindListParam
	SiteID uint   `form:"site_id"`
	Uri    string `form:"uri" binding:"omitempty,min=2,max=64"`
}

// 处理列表请求数据。
func (b *bindCacheDeleteParam) Param() *model.CacheDelete {
	m := &model.CacheDelete{}
	w := func(orm *gorm.DB) *gorm.DB {
		return orm.Where("uri like ?", b.Uri+"%")
	}
	m.Bind(m, w, b)
	return m
}
