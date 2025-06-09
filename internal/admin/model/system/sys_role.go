package system

import (
	"github.com/penndev/galite/pkg/orm"
)

type RouteItem struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

type SysRole struct {
	orm.Model
	Name   string      `json:"name"`
	Status uint8       `json:"status"`
	Menu   []string    `gorm:"serializer:json" json:"menu"`
	Route  []RouteItem `gorm:"serializer:json" json:"route"`
	Remark string      `json:"remark"`
}
