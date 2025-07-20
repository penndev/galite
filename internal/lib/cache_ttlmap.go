package lib

import (
	"errors"
	"reflect"
	"time"

	"github.com/penndev/galite/pkg/cache"
	"github.com/penndev/gopkg/ttlmap"
)

// TTLMap实例
type ttlmapClient struct {
	*ttlmap.TTLMap
}

// var TTLMap *ttlmapClient

func (t *ttlmapClient) GetAny(k string, s any) error {
	result, ok := t.Get(k)
	if !ok {
		return errors.New(k + " not found")
	}
	// s 应该是指针，否则不能赋值
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Ptr {
		return errors.New("s must be a pointer")
	}

	// 将结果赋值给 s 指向的值
	v.Elem().Set(reflect.ValueOf(result))
	return nil
}

func (t *ttlmapClient) SetAny(k string, s any, exp time.Duration) error {
	t.Set(k, s, exp)
	return nil
}

func (r *ttlmapClient) Delete(k string) error {
	r.TTLMap.Delete(k)
	return nil
}

func InitTTLMap(dsn string) (cache.Interface, error) {
	TTLMap := &ttlmapClient{
		TTLMap: ttlmap.New(),
	}
	return TTLMap, nil
}
