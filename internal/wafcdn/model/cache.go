package model

import (
	"os"

	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/pkg/orm"
	"go.uber.org/zap"
)

type Cache struct {
	orm.ModelBase
	SiteID uint   `json:"site_id" form:"site_id" binding:"required"`
	Uri    string `json:"uri" form:"uri" binding:"required" gorm:"index"`
	Method string `json:"method" form:"method" binding:"required"`
	Path   string `json:"path" form:"path"`
	Time   int    `json:"time" form:"time"`
}

func CacheDeleteList(caches []Cache) error {
	for _, cache := range caches {
		err := os.Remove(cache.Path)
		if err != nil {
			logger.L.Error("clearCache", zap.Error(err))
		}
		err = os.Remove(cache.Path + ".head")
		if err != nil {
			logger.L.Error("clearCache", zap.Error(err))
		}
	}
	return (&Cache{}).DB().Delete(&caches).Error
}
