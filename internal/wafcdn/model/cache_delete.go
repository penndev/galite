package model

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/penndev/galite/pkg/orm"
	"gorm.io/gorm"
)

type CacheDelete struct {
	orm.ModelBase
	SiteID uint   `json:"site_id" form:"site_id"`
	Uri    string `json:"uri" form:"uri" binding:"required"`
	Log    string `json:"log" form:"log"`
	Status bool   `json:"status" form:"status"` // 是否完成了任务。
}

func CacheDeleteAction() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("CacheDeleteAction panic:", r)
		}
	}()
	time.Sleep(10 * time.Second)
	for {
		var cacheDelete CacheDelete
		cacheDelete.DB().Where("status = false").Limit(1).Order("id ASC").Find(&cacheDelete)
		if cacheDelete.ID < 1 {
			time.Sleep(10 * time.Second)
			continue
		}

		cache := &Cache{}
		cacheWhere := func(db *gorm.DB) *gorm.DB {
			db.Where("site_id = ? and uri like ?", cacheDelete.SiteID, cacheDelete.Uri+"%")
			return db
		}
		var total int64
		query := cache.Bind(cache, cacheWhere).BindGorm()
		query.Count(&total)
		cacheDelete.Log += "匹配缓存数为：" + strconv.FormatInt(total, 10) + "\n"

		deleteCacheTotal := 0
		for {
			time.Sleep(100 * time.Millisecond)
			var caches []Cache
			if err := query.Limit(1000).Find(&caches).Error; err != nil {
				log.Println(err)
				time.Sleep(10 * time.Second)
				break
			}
			CacheDeleteList(caches)

			deleteCacheTotal += len(caches)
			cacheDelete.Log = fmt.Sprintf("缓存总数为 %d  已经删除总数为 %d\r\n", total, deleteCacheTotal)
			if len(caches) > 0 {
				cacheDelete.Save()
			} else {
				cacheDelete.Status = true
				cacheDelete.Save()
				break
			}
		}
	}
}
