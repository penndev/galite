package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
)

// 获取站点列表
func LogList(c *gin.Context) {
	param := &bindLogParam{}
	if err := c.BindQuery(param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}

	var total int64
	var list []model.Log

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func LogClear(c *gin.Context) {
	err := model.LogTruncate()
	if err == nil {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	} else {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "执行错误" + err.Error()})
	}
}
