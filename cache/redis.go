package cache

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

// Redis实例
var Redis *redis.Client

func InitRedis(redisURL string) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Panic(err)
	}
	// redis 默认已经存在连接池
	Redis = redis.NewClient(opt)
	if _, err := Redis.Ping(context.TODO()).Result(); err != nil {
		log.Panic(err)
	}
}
