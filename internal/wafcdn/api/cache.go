package api

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/internal/wafcdn/model"
	"github.com/shirou/gopsutil/v4/disk"
	"go.uber.org/zap"
)

var cacheNumber uint64 = 0

// 缓存文件入库，并清盘
func HandlePutCache(c *gin.Context) {
	param := &model.Cache{}
	if err := c.BindJSON(param); err != nil {
		c.JSON(400, gin.H{"message": "参数错误" + err.Error()})
		return
	}

	// 保存文件并清盘
	if cacheNumber%500 == 0 {
		dir, maxUsed, n := filepath.Dir(param.Path), 95.00, 1000
		stat, err := disk.Usage(dir)
		if err != nil {
			log.Println("clearCache", zap.Error(err))
		}
		if stat.UsedPercent > maxUsed {
			var caches []model.Cache
			(&model.Cache{}).DB().Order("id asc").Limit(n).Find(&caches)
			if err := model.CacheDeleteList(caches); err != nil {
				log.Println("clearCache", zap.Error(err))
			}
		}
	}
	cacheNumber++

	err := param.DB().Where(
		"site_id = ? and method = ? and uri = ?",
		param.SiteID, param.Method, param.Uri,
	).Assign(*param).FirstOrCreate(param).Error
	if err != nil {
		logger.L.Warn("保存错误", zap.Error(err))
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
