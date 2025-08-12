package system

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/admin/model/system"
	"github.com/penndev/galite/internal/logger"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func AdminList(c *gin.Context) {
	param := &bindSystemAdminParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysAdmin
	log.Println("i am  here")
	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func AdminAdd(c *gin.Context) {
	param := &system.SysAdmin{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if param.Passwd == "" {
		str, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
		if err != nil {
			logger.L.Error("创建管理员密码失败", zap.Error(err))
			c.JSON(http.StatusBadRequest, bind.Message{Message: "初始化管理员失败，请查看错误日志"})
			return
		}
		param.Passwd = string(str)
	}
	if err := param.DB().Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

func AdminUpdate(c *gin.Context) {
	param := &system.SysAdmin{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	if param.Passwd == "" {
		str, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.MinCost)
		if err != nil {
			logger.L.Error("创建管理员密码失败", zap.Error(err))
			c.JSON(http.StatusBadRequest, bind.Message{Message: "初始化管理员失败，请查看错误日志"})
			return
		}
		param.Passwd = string(str)
	}
	if err := param.DB().Updates(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "更新失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

func AdminDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	param := &system.SysAdmin{}
	param.ID = uint(id)
	if err := param.DB().Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

func AdminAccessLog(c *gin.Context) {
	c.Set("accessLog", false) // 设置访问日志标志
	param := &bindSysAccessParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	var total int64
	var list []system.SysAccessLog

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}
