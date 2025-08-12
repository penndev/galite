package system

import "github.com/penndev/galite/pkg/orm"

type SysAccessLog struct {
	orm.ModelBase
	SysAdminID uint     `json:"SysAdminId"`
	SysAdmin   SysAdmin `json:"SysAdmin"`
	Method     string   `json:"method"`
	Payload    string   `json:"payload"`
	Status     int      `json:"status"`
	Path       string   `json:"path"`
	IP         string   `json:"ip"`
}
