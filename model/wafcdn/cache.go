package wafcdn

import (
	"github.com/penndev/galite/model/orm"
)

type Cache struct {
	orm.Model
	SiteID string            `json:"site_id" form:"site_id" binding:"required"`
	Method string            `json:"method" form:"method" binding:"required"`
	Uri    string            `json:"uri" form:"uri" binding:"required"`
	Header map[string]string `gorm:"serializer:json" json:"header" form:"header"` // 代理返回头
	Path   string            `json:"path" form:"path"`
	Time   int               `json:"time" form:"time"`
}
