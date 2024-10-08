package system

import "github.com/penndev/galite/model/orm"

type SysAccessLog struct {
	orm.Model
	SysAdminID uint     `json:"SysAdminId"`
	SysAdmin   SysAdmin `json:"SysAdmin"`
	Method     string   `json:"method"`
	Payload    string   `json:"payload"`
	Status     int      `json:"status"`
	Path       string   `json:"path"`
	IP         string   `json:"ip"`
}
