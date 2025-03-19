package wafcdn

import (
	"github.com/penndev/galite/model/orm"
)

type Cache struct {
	orm.Model
	SiteID string            `json:"site_id"`
	Method string            `json:"method"`
	Uri    string            `json:"uri"`
	Header map[string]string `gorm:"serializer:json" json:"header"` // 代理返回头
	Path   string            `json:"path"`
	Time   int               `json:"time"`
}
