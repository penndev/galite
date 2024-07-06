package system

import (
	"github.com/penndev/wga/config"
	"github.com/penndev/wga/model/sugar"
)

type SysAdmin struct {
	sugar.Model
	Email     string  `gorm:"uniqueIndex,size=256" json:"email"`
	Passwd    string  `json:"-"`
	SysRoleID uint    `json:"adminRoleId"`
	SysRole   SysRole `json:"AdminRole"`
	NickName  string  `json:"nickname"`
	Status    uint8   `json:"status"`
	Remark    string  `json:"remark"`
}

func SysAdminGetByEmail(email string) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := config.DB.Where(&SysAdmin{Email: email}).Preload("SysRole").First(&sysAdmin)
	return &sysAdmin, result.Error
}

func SysAdminGetByID(id string) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := config.DB.Preload("SysRole").Find(&sysAdmin, id)
	return &sysAdmin, result.Error
}
