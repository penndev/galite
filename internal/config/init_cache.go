package config

import (
	"errors"
	"strings"

	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/pkg/cache"
)

func InitCache() error {
	var err error
	var cache cache.Interface
	dsn := cfg.Cache.Dsn
	switch {
	case strings.HasPrefix(dsn, "redis://"):
		if cache, err = lib.InitRedis(dsn); err != nil {
			return err
		}
		closeList = append(closeList, cache.(closeInterface))
	case strings.HasPrefix(dsn, "ttlmap://"):
		if cache, err = lib.InitTTLMap(dsn); err != nil {
			return err
		}
	default:
		return errors.New("not set cache env")
	}
	lib.SetCache(cache)
	return nil
}
