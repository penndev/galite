package cache

import (
	"bytes"
	"context"
	"encoding/gob"
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

func Encode(data any) (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

func Decode(str string, target any) error {
	buf := bytes.NewBufferString(str)
	dec := gob.NewDecoder(buf)
	return dec.Decode(target)
}
