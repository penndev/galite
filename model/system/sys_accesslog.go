package system

import "github.com/penndev/wga/model/sugar"

type SysAccessLog struct {
	sugar.Model
	SysAdminID uint     `json:"adminId"`
	SysAdmin   SysAdmin `json:"AdminUser"`
	Method     string   `json:"method"`
	Payload    string   `json:"payload"`
	Status     int      `json:"status"`
	Path       string   `json:"path"`
	IP         string   `json:"ip"`
}
