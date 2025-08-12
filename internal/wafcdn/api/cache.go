package api

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func HandleGetCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindQuery(param); err != nil {
		logger.L.Warn("参数错误", zap.Error(err))
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	if err := param.GetCache(); err == nil {
		c.JSON(http.StatusOK, param)
		return
	}

	err := param.DB().Where("site_id = ? and method = ? and uri = ?", param.SiteID, param.Method, param.Uri).First(param).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Message": "查询失败(" + err.Error() + ")"})
		return
	} else { // 保存缓存结果 永久缓存
		if err := param.SetCache(); err != nil {
			logger.L.Warn("保存错误", zap.Error(err))
		}
		// param.GetCache()
		c.JSON(http.StatusOK, param)
	}
}

func HandlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}
	err := param.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Where(
			"site_id = ? and method = ? and uri = ?",
			param.SiteID, param.Method, param.Uri,
		).Assign(*param).FirstOrCreate(param).Error; err != nil {
			return err
		}
		if err := os.Rename(param.Path+".lock", param.Path); err != nil {
			return err
		}
		if err := param.SetCache(); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logger.L.Warn("保存错误", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"Message": "创建失败(" + err.Error() + ")"})
	} else {
		c.JSON(http.StatusOK, gin.H{"Message": "完成"})
	}
}
