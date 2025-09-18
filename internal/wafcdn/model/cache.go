package model

import (
	"os"

	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/pkg/orm"
	"go.uber.org/zap"
	"gorm.io/gorm"
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
	db := (&Cache{}).DB()
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&caches).Error; err != nil {
			logger.L.Error("clearCache", zap.Error(err))
			return err
		}
		for _, cache := range caches {
			if err := os.Remove(cache.Path); err != nil && !os.IsNotExist(err) {
				logger.L.Error("clearCache", zap.Error(err))
			}
			if err := os.Remove(cache.Path + ".head"); err != nil && !os.IsNotExist(err) {
				logger.L.Error("clearCache", zap.Error(err))
			}
		}
		return nil
	})
}
