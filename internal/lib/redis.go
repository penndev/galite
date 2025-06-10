package lib

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Redis实例
var Redis *redis.Client

func InitRedis(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	// redis 默认已经存在连接池
	Redis = redis.NewClient(opt)
	if _, err := Redis.Ping(context.TODO()).Result(); err != nil {
		return err
	}
	return nil
}
