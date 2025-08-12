package system

import (
	"github.com/penndev/galite/pkg/orm"
)

type SysAdmin struct {
	orm.ModelBase
	Email     string  `gorm:"uniqueIndex,size=256" json:"email"`
	Passwd    string  `json:"-"`
	SysRoleID *uint   `json:"SysRoleId"` // 必须用指针因为外键关联问题 foreign key constraint
	SysRole   SysRole `json:"SysRole"`
	Nickname  string  `json:"nickname"`
	Status    uint8   `json:"status"`
	OtpStatus uint8   `json:"otpStatus"`
	OtpTitle  string  `json:"otpTitle"`
	OtpSecret string  `json:"otpSecret"`
	Remark    string  `json:"remark"`
}

func SysAdminGetByEmail(email string) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := sysAdmin.Bind(&SysAdmin{Email: email}).BindGorm().Preload("SysRole").First(&sysAdmin)
	return &sysAdmin, result.Error
}

func SysAdminGetByID(id int) (*SysAdmin, error) {
	var sysAdmin SysAdmin
	result := sysAdmin.DB().Preload("SysRole").First(&sysAdmin, id)
	return &sysAdmin, result.Error
}
