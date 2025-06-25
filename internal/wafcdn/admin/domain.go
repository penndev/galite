package admin

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
)

// 添加新的站点

func DomainAdd(c *gin.Context) {
	param := &model.Domain{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误" + err.Error()})
		return
	}
	if err := param.Bind(param).Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}

// 获取站点列表
func DomainList(c *gin.Context) {
	param := &bindDomainParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	var total int64
	var list []model.Domain

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

// 更新资料
func DomainUpdate(c *gin.Context) {
	param := &model.Domain{}
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

// 删除资料
func DomainDelete(c *gin.Context) {
	id, err := strconv.Atoi(c.Query("id"))
	if id < 1 || err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "参数错误"})
		return
	}
	param := &model.Domain{}
	param.ID = uint(id)
	if err := param.Bind(param).Delete(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.ErrorMessage{Message: "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, bind.ErrorMessage{Message: "完成"})
	}
}
