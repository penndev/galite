package system

import (
	"github.com/penndev/galite/model/orm"
)

type SysAdmin struct {
	orm.Model
	Email     string  `gorm:"uniqueIndex,size=256" json:"email"`
	Passwd    string  `json:"-"`
	SysRoleID *uint   `json:"SysRoleId"`
	SysRole   SysRole `json:"SysRole"`
	Nickname  string  `json:"nickname"`
	Status    uint8   `json:"status"`
	Remark    string  `json:"remark"`
}

func SysAdminGetByEmail(email string) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := orm.DB.Where(&SysAdmin{Email: email}).Preload("SysRole").First(&sysAdmin)
	return &sysAdmin, result.Error
}

func SysAdminGetByID(id string) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := orm.DB.Preload("SysRole").Where("id = ?", id).First(&sysAdmin)
	return &sysAdmin, result.Error
}
