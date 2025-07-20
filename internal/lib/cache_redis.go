package lib

import (
	"bytes"
	"context"
	"time"

	"github.com/penndev/galite/pkg/cache"
	"github.com/penndev/galite/pkg/util"
	"github.com/redis/go-redis/v9"
)

// Redis实例
type redisClient struct {
	*redis.Client
}

// var Redis *redisClient

func (r *redisClient) GetAny(k string, s any) error {
	result, err := r.Get(context.TODO(), k).Result()
	if err != nil {
		return err
	}
	if err = util.Decode(bytes.NewBufferString(result), s); err != nil {
		return err
	}
	return nil
}

func (r *redisClient) SetAny(k string, s any, exp time.Duration) error {
	buf, err := util.Encode(s)
	if err != nil {
		return err
	}
	return r.Set(context.TODO(), k, buf.String(), exp).Err()
}

func (r *redisClient) Delete(k string) error {
	return r.Del(context.TODO(), k).Err()
}

func InitRedis(redisURL string) (cache.Interface, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	// redis 默认已经存在连接池
	Redis := &redisClient{
		Client: redis.NewClient(opt),
	}
	if _, err := Redis.Ping(context.TODO()).Result(); err != nil {
		return nil, err
	}
	return Redis, nil
}
