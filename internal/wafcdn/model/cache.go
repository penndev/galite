package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dgraph-io/badger/v4"
	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/internal/logger"
	"github.com/penndev/galite/pkg/orm"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 原本的header map[string]string
// 无法适用重复的header头 如果设置两次set cookie
type Header http.Header

func (h Header) MarshalJSON() ([]byte, error) {
	m := make(map[string]any)
	for k, v := range h {
		if len(v) == 1 {
			m[k] = v[0]
		} else {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

func (h *Header) UnmarshalJSON(data []byte) error {
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	result := make(map[string][]string)
	for k, v := range m {
		switch val := v.(type) {
		case string:
			result[k] = []string{val}
		case []string:
			strVal := make([]string, len(val))
			for i, item := range val {
				strVal[i] = fmt.Sprintf("%v", item)
			}
			result[k] = strVal
		default:
			result[k] = []string{fmt.Sprintf("%v", val)}
		}
	}
	*h = result
	return nil
}

type Cache struct {
	orm.Model
	SiteID uint   `json:"site_id" form:"site_id" binding:"required"`
	Uri    string `json:"uri" form:"uri" binding:"required"`
	Method string `json:"method" form:"method" binding:"required"`
	Header Header `gorm:"serializer:json" json:"header" form:"header"` // 代理返回头
	Path   string `json:"path" form:"path"`
	Time   int    `json:"time" form:"time"`
}

func (cache *Cache) CacheKey() string {
	return fmt.Sprintf("%d%s%s", cache.SiteID, cache.Uri, cache.Method)
}

// 删除一个文件应该如何删除呢
// 首先肯定要开启事务来达成一个原子性的操作
// 因为使用了 Badger 所以要用这个来先操作锁
func DeleteCaches(ids []uint) error {
	// 开启数据库事务。来删除数据库中的缓存
	// 首先批量删除Badger和缓存文件，保证请求不会出现存在缓存则404的情况。
	var cache Cache
	var caches []Cache
	cache.Bind(cache, func(db *gorm.DB) *gorm.DB {
		db.Where("id IN ?", ids)
		return db
	}).Find(&caches)
	// 需要优化批量删除性能不会太好
	for _, c := range caches {
		err := lib.Badger.Update(func(txn *badger.Txn) error {
			err := txn.Delete([]byte(c.CacheKey()))
			return err
		})
		if err != nil {
			logger.L.Warn("cache/delete", zap.Error(err))
			continue
		}
		os.Remove(c.Path)
		c.Gorm().Delete(&c)
	}
	return nil
}
