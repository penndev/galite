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

func CacheDeleteTruncate() error {
	return (&CacheDelete{}).DB().Exec("truncate table cache_deletes").Error
}

// 运行锁 防止删除并发
var CacheDeleteActionRuning = false

func CacheDeleteAction() {
	if CacheDeleteActionRuning {
		return
	} else {
		CacheDeleteActionRuning = true
	}

	defer func() {
		CacheDeleteActionRuning = false
		if r := recover(); r != nil {
			log.Println("CacheDeleteAction panic:", r)
		}
	}()

	for {
		cacheDelete := &CacheDelete{
			Status: false,
		}
		cacheDelete.Bind(cacheDelete, func(orm *gorm.DB) *gorm.DB {
			return orm.Where("status = false").Order("id asc").Limit(1)
		}).Find(cacheDelete)
		if cacheDelete.ID < 1 {
			return
		}
		cache := &Cache{}
		query := cache.Bind(cache, func(orm *gorm.DB) *gorm.DB {
			return orm.Where("site_id = ? and uri like ?", cacheDelete.SiteID, cacheDelete.Uri+"%")
		}).BindGorm()

		var total int64
		query.Count(&total)

		cacheDelete.Log += "匹配缓存数为：" + strconv.FormatInt(total, 10) + "\n"

		deleteCacheTotal := 0
		for {
			time.Sleep(10 * time.Millisecond)
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
