package system

import (
	// "github.com/penndev/galite/admin"
	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/model/orm"
)

type SysRole struct {
	orm.Model
	Name   string           `json:"name"`
	Status uint8            `json:"status"`
	Menu   []string         `gorm:"serializer:json" json:"menu"`
	Route  []bind.RouteItem `gorm:"serializer:json" json:"route"`
	Remark string           `json:"remark"`
}
