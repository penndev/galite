package system

import (
	"log"
	"net/http"
	"strconv"

	"github.com/penndev/galite/admin/bind"
	"github.com/penndev/galite/config"
	"github.com/penndev/galite/model/system"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"golang.org/x/crypto/bcrypt"
)

func AdminList(c *gin.Context) {
	param := &bindSystemAdminParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysAdmin

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func AdminAdd(c *gin.Context) {
	param := &system.SysAdmin{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误" + err.Error()})
		return
	}
	if param.Passwd == "" {
		str, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
		if err != nil {
			config.Logger.Error("创建管理员密码失败", zap.Error(err))
			c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "初始化管理员失败，请查看错误日志"})
			return
		}
		param.Passwd = string(str)
	}
	if err := param.Bind(param).Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func AdminUpdate(c *gin.Context) {
	param := &system.SysAdmin{}
	if err := c.BindJSON(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}

	if err := param.Bind(param).Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func AdminDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param := &system.SysAdmin{}
	param.ID = uint(id)
	if err := param.Bind(param).Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

func AdminAccessLog(c *gin.Context) {
	param := &bindSysAccessParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysAccessLog

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}
