package lib

import (
	"bytes"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis实例
type redisClient struct {
	*redis.Client
}

var Redis *redisClient

func (r *redisClient) GetStruct(k string, s any) error {
	result, err := Redis.Get(context.TODO(), k).Result()
	if err != nil {
		return err
	}
	if err = Decode(bytes.NewBufferString(result), s); err != nil {
		return err
	}
	return nil
}

func (r *redisClient) SetStruct(k string, s any, exp time.Duration) error {
	buf, err := Encode(s)
	if err != nil {
		return err
	}
	return Redis.Set(context.TODO(), k, buf.String(), exp).Err()
}

func InitRedis(redisURL string) error {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return err
	}
	// redis 默认已经存在连接池
	Redis = &redisClient{
		Client: redis.NewClient(opt),
	}
	if _, err := Redis.Ping(context.TODO()).Result(); err != nil {
		return err
	}
	return nil
}
