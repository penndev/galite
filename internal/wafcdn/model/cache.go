package model

import (
	"os"

	"github.com/penndev/galite/pkg/orm"
)

type Cache struct {
	orm.ModelBase
	SiteID uint   `json:"site_id" form:"site_id" binding:"required"`
	Uri    string `json:"uri" form:"uri" binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
	Path   string `json:"path" form:"path"`
	Time   int    `json:"time" form:"time"`
}

func CacheDeleteList(caches []Cache) error {
	for _, cache := range caches {
		err := os.Remove(cache.Path)
		if err != nil {
			return err
		}
		err = os.Remove(cache.Path + ".head")
		if err != nil {
			return err
		}
	}
	return (&Cache{}).DB().Delete(&caches).Error
}
