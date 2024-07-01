package config

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func initRedis() {
	if os.Getenv("REDIS_DSN") == "" {
		log.Panic(errors.New("REDIS_DSN .env not found"))
	}
	opt, err := redis.ParseURL("redis://" + os.Getenv("REDIS_DSN"))
	if err != nil {
		log.Panic(err)
	}
	// redis 默认已经存在连接池
	Redis = redis.NewClient(opt)
	if _, err := Redis.Ping(context.TODO()).Result(); err != nil {
		log.Panic(err)
	}
}
