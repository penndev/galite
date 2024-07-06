package system

import "github.com/penndev/wga/model/sugar"

type SysAccessLog struct {
	sugar.Model
	SysAdminID uint     `json:"adminId"`
	SysAdmin   SysAdmin `json:"AdminUser"`
	Method     string   `json:"method"`
	Path       string   `json:"path"`
	IP         string   `json:"ip"`
}
