package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"go.uber.org/zap"
)

// 缓存文件入库，并清盘
func HandlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	// 处理异步清理文件。
	go model.CacheDeleteMaxUsed(param.Path)

	if err := param.DB().Create(param).Error; err != nil {
		logger.L.Error("HandlePutCache", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}

func HandlePurgeCache(c *gin.Context) {
	param := struct {
		SiteID uint   `form:"site_id" binding:"required"`
		Uri    string `form:"uri" binding:"required"`
	}{}
	if err := c.BindQuery(&param); err != nil {
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	var caches []model.Cache
	(&model.Cache{}).DB().Where("site_id = ? and uri like ?", param.SiteID, param.Uri+"%").Find(&caches)
	if err := model.CacheDeleteList(caches); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "删除失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "已删除" + strconv.Itoa(len(caches)) + "条"})
	}
}
