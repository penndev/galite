package config

import (
	"errors"
	"strings"

	"github.com/penndev/galite/internal/lib"
	"github.com/penndev/galite/pkg/cache"
)

func InitCache() error {
	var err error
	var cacheClient cache.Interface
	dsn := cfg.Cache.Dsn
	switch {
	case strings.HasPrefix(dsn, "redis://"):
		if cacheClient, err = cache.InitRedis(dsn); err != nil {
			return err
		}
		closeList = append(closeList, cacheClient.(closeInterface))
	case strings.HasPrefix(dsn, "ttlmap://"):
		if cacheClient, err = cache.InitTTLMap(dsn); err != nil {
			return err
		}
	default:
		return errors.New("not set cache env")
	}
	lib.SetCache(cacheClient)
	return nil
}
