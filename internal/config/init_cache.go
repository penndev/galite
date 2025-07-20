package config

import (
	"errors"
	"strings"

	"github.com/penndev/galite/internal/lib"
)

func InitCache() error {
	var err error
	dsn := cfg.Cache.Dsn
	switch {
	case strings.HasPrefix(dsn, "redis://"):
		if lib.Cache, err = lib.InitRedis(dsn); err != nil {
			return err
		}
		closeList = append(closeList, lib.Cache.(closeInterface))
	case strings.HasPrefix(dsn, "ttlmap://"):
		if lib.Cache, err = lib.InitTTLMap(dsn); err != nil {
			return err
		}
	default:
		return errors.New("not set cache env")
	}
	return nil
}
