package api

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"go.uber.org/zap"
)

func HandleGetCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindQuery(param); err != nil {
		logger.L.Warn("参数错误", zap.Error(err))
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	// 设置看门口是否。
	cacheKey := param.CacheKey()
	if err := lib.Redis.GetStruct(cacheKey, param); err != nil {
		c.JSON(http.StatusOK, param)
		return
	}
	err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).First(param).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "查询失败(" + err.Error() + ")"})
		return
	} else { // 保存缓存结果 永久缓存
		if err := param.SetCache(cacheKey); err != nil {
			logger.L.Warn("保存错误", zap.Error(err))
		}
		c.JSON(http.StatusOK, param)
	}
}

func HandlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	if err := param.Bind(param).Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).Assign(*param).FirstOrCreate(&model.Cache{
		SiteID: param.SiteID,
		Method: param.Method,
		Uri:    param.Uri,
		Path:   param.Path,
	}).Error; err != nil {
		logger.L.Warn("保存错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		if err := param.SetCache(param.CacheKey()); err != nil {
			logger.L.Warn("保存错误", zap.Error(err))
		}
		if err := os.Rename(param.Path+".lock", param.Path); err != nil {
			logger.L.Warn("保存错误", zap.Error(err))
		}
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}

func HandlePutLog(c *gin.Context) {
	param := &model.Log{}
	if err := c.BindJSON(param); err != nil {
		log.Println("参数错误", err.Error())
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	param.Gorm().Create(param)
}
