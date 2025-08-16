package cache

import (
	"bytes"
	"context"
	"time"

	"github.com/penndev/galite/pkg/util"
	"github.com/redis/go-redis/v9"
)

// Redis实例
type RedisClient struct {
	*redis.Client
}

// var Redis *redisClient

func (r *RedisClient) GetAny(k string, s any) error {
	result, err := r.Get(context.TODO(), k).Result()
	if err != nil {
		return err
	}
	if err = util.Decode(bytes.NewBufferString(result), s); err != nil {
		return err
	}
	return nil
}

func (r *RedisClient) SetAny(k string, s any, exp time.Duration) error {
	buf, err := util.Encode(s)
	if err != nil {
		return err
	}
	return r.Set(context.TODO(), k, buf.String(), exp).Err()
}

func (r *RedisClient) Delete(k string) error {
	return r.Del(context.TODO(), k).Err()
}

func InitRedis(redisURL string) (Interface, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	// redis 默认已经存在连接池
	redisClient := &RedisClient{
		Client: redis.NewClient(opt),
	}
	if _, err := redisClient.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}
	return redisClient, nil
}
