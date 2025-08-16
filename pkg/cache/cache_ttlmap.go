package cache

import (
	"errors"
	"reflect"
	"time"

	"github.com/penndev/gopkg/ttlmap"
)

// TTLMap实例
type TTLMap struct {
	*ttlmap.Map
}

// var TTLMap *ttlmapClient

func (t *TTLMap) GetAny(k string, s any) error {
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

func (t *TTLMap) SetAny(k string, s any, exp time.Duration) error {
	if exp < 1 {
		exp := 24 * time.Hour
		t.Set(k, s, exp)
	} else {
		t.Set(k, s, exp)
	}
	return nil
}

func (r *TTLMap) Delete(k string) error {
	r.Map.Delete(k)
	return nil
}

func InitTTLMap(dsn string) (Interface, error) {
	TTLMap := &TTLMap{
		Map: ttlmap.New(),
	}
	return TTLMap, nil
}
