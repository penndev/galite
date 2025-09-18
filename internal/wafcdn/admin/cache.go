package admin

import (
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
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
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
	cIDs, err := util.StrArrConv[uint](ids)
	if err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: err.Error()})
		return
	}
	var caches []model.Cache
	(&model.Cache{}).DB().Where("id in ?", cIDs).Find(&caches)
	model.CacheDeleteList(caches)
	c.JSON(http.StatusOK, bind.Message{
		Message: "完成",
	})
}

func CacheDeleteList(c *gin.Context) {
	param := &bindCacheDeleteParam{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	var total int64
	var list []model.CacheDelete
	m := param.Param() //处理筛选
	m.List(&total, &list)
	c.JSON(http.StatusOK, bind.DataList{Total: total, Data: list})
}

func CacheDeleteAdd(c *gin.Context) {
	param := &model.CacheDelete{}
	if err := c.BindJSON(&param); err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "参数错误" + err.Error()})
		return
	}
	if err := param.DB().Create(param).Error; err != nil {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "创建失败(" + err.Error() + ")"})
	} else {
		go model.CacheDeleteAction()
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	}
}

func CacheDeleteClear(c *gin.Context) {
	err := model.CacheDeleteTruncate()
	if err == nil {
		c.JSON(http.StatusOK, bind.Message{Message: "完成"})
	} else {
		c.JSON(http.StatusBadRequest, bind.Message{Message: "执行错误" + err.Error()})
	}
}
