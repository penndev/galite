package model

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/pkg/orm"
	"github.com/shirou/gopsutil/v4/disk"
	"go.uber.org/zap"
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

// 运行锁 防止删除并发 暂时不加锁
var CacheDeleteActionRuning = false

// 后台运行删除任务
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

var cacheNumber uint64 = 0

// 缓存阈值删除任务
func CacheDeleteMaxUsed(path string) {
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
	// 保存文件并清盘
	allow := false
	if cacheNumber%100 == 0 {
		dir, maxUsed, n := filepath.Dir(path), 95.00, 200
		stat, err := disk.Usage(dir)
		if err != nil {
			logger.L.Error("clearCache", zap.Error(err))
		}
		if stat.UsedPercent > maxUsed {
			var caches []Cache
			(&Cache{}).DB().Order("id asc").Limit(n).Find(&caches)
			if err := CacheDeleteList(caches); err != nil {
				logger.L.Error("clearCache", zap.Error(err))
			}
		}
		if stat.UsedPercent > 97 {
			allow = true
		} else {
			allow = false
		}

	}
	if allow {
		cacheNumber = 0
	} else {
		cacheNumber++
	}
}
