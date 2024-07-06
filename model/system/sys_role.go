package system

import (
	"github.com/penndev/wga/admin/middle"
	"github.com/penndev/wga/model/sugar"
)

type SysRole struct {
	sugar.Model
	Name   string             `json:"name"`
	Status uint8              `json:"status"`
	Menu   []string           `gorm:"serializer:json" json:"menu"`
	Route  []middle.RouteItem `gorm:"serializer:json" json:"route"`
	Remark string             `json:"remark"`
}
