package model

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

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
	return fmt.Sprintf("%d:%s:%s", cache.SiteID, cache.Uri, cache.Method)
}

func (cache *Cache) GetCache() error {
	err := lib.Cache.GetAny(cache.CacheKey(), cache)
	return err
}

// 缓存需要的数据，缓存到数据库，因为多次调用所以放置model中
func (cache *Cache) SetCache() error {
	return lib.Cache.SetAny(
		cache.CacheKey(),
		*cache,
		-1,
	)
}

func CacheGetByIds(ids []uint) []Cache {
	m := &Cache{}
	var caches []Cache
	m.Bind(m, func(db *gorm.DB) *gorm.DB {
		db.Where("id IN ?", ids)
		return db
	}).Find(&caches)
	return caches
}

// 删除一个文件应该如何删除呢
// 首先肯定要开启事务来达成一个原子性的操作
// 然后在数据库插入成功后再修改完整的文件名
// 不然文件名称修改后但是最新的文件再次被删除了。
func CacheDeleteByData(caches []Cache) error {
	m := &Cache{}
	m.Gorm().Transaction(func(tx *gorm.DB) error {
		for _, cache := range caches {
			tx.Delete(&cache)
			lib.Cache.Delete(cache.CacheKey())
			if err := os.Remove(cache.Path); err != nil {
				logger.L.Warn("删除失败", zap.Error(err))
			}
		}
		return nil
	})
	return nil
}
