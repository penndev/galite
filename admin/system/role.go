package system

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/wga/admin/bind"
	"github.com/penndev/wga/model/system"
)

func RoleList(c *gin.Context) {
	param := &bindSystemRoleParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysRole

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func RoleAdd(c *gin.Context) {
	param := &system.SysRole{}
	if err := c.BindJSON(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param.Bind(param)
	if err := param.Create(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func RoleUpdate(c *gin.Context) {
	param := &system.SysRole{}
	if err := c.BindJSON(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param.Bind(param)
	if err := param.Update(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func RoleDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param := &system.SysRole{}
	param.ID = uint(id)
	param.Bind(param)
	if err := param.Delete(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}
