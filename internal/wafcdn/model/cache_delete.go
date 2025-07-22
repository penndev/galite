package model

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/penndev/galite/pkg/orm"
	"gorm.io/gorm"
)

type CacheDelete struct {
	orm.Model
	SiteID uint   `json:"site_id" form:"site_id"`
	Uri    string `json:"uri" form:"uri" binding:"required"`
	Log    string `json:"log" form:"log"`
	Status bool   `json:"status" form:"status"` // 是否完成了任务。
}

func CacheDeleteRunner() {
	timeSleep := 10 * time.Second
	for {
		cacheDelete := &CacheDelete{}
		if err := cacheDelete.Bind(cacheDelete, func(db *gorm.DB) *gorm.DB {
			db.Where("status = false").Order("id ASC")
			return db
		}).First(cacheDelete).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				time.Sleep(timeSleep)
				// log.Println("没找到记录重新循环(正式环境关掉): ", err)
				continue
			}
			log.Println("Runner: ", err)
		}

		// 筛选缓存条件
		cache := &Cache{}
		action := cache.Bind(cache, func(db *gorm.DB) *gorm.DB {
			db.Where("site_id = ? and uri like ?", cacheDelete.SiteID, cacheDelete.Uri+"%")
			return db
		})
		// 查询总数
		var total int64
		action.Count(&total)
		actionLog := "匹配缓存数为：" + strconv.FormatInt(total, 10) + "\n"
		cacheDelete.Log += actionLog
		totalDeleted := 0
		for {
			var caches []Cache
			if err := action.Limit(1000).Find(&caches).Error; err != nil {
				log.Println(err)
				time.Sleep(timeSleep)
				break
			}
			CacheDeleteByData(caches)
			totalCurrent := len(caches)
			totalDeleted += totalCurrent
			actionLog = fmt.Sprintf("缓存总数为 %d  已经删除总数为 %d  本次删除数量为 %d \r\n", total, totalDeleted, totalCurrent)
			cacheDelete.Log += actionLog
			cacheDelete.Save()
			if totalCurrent == 0 {
				break //全部缓存完成
			}
			// 删除文件 需要增加休眠时间
			time.Sleep(500 * time.Millisecond)
		}
		cacheDelete.Status = true
		cacheDelete.Save()
	}
}
