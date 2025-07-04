package admin

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/admin/bind"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/penndev/galite/pkg/util"
)

// 获取站点列表
func CacheList(c *gin.Context) {
	param := &bindCacheParam{}
	if err := c.BindQuery(&param); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误"})
		return
	}
	var total int64
	var list []model.Cache

	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

// 删除文件如何保持 原子性。
func CacheDelete(c *gin.Context) {
	ids := c.QueryArray("ids")
	idInt, err := util.StrConvArr[uint](ids)
	model.CacheDeleteByIds(idInt)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, bind.Message{
		Message: "完成",
	})
}
