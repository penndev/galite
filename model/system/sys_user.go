package system

import (
	"github.com/penndev/wga/config"
	"github.com/penndev/wga/model/sugar"
)

type SysUser struct {
	sugar.Model
	Name     string `json:"name"`
	Password string `json:"-"`
	RoleID   uint   `json:"roleId"`
	Remark   string `json:"remark"`
}

func GetSysUserByName(name string) (SysUser, error) {
	var sysuser SysUser
	result := config.DB.Where(&SysUser{Name: name}).First(&sysuser)
	return sysuser, result.Error
}

func CreateSysUser(sysuser *SysUser) error {
	result := config.DB.Save(sysuser)
	return result.Error
}
