package cache

import (
	"TTCS/src/common/log"
	"context"
	"github.com/redis/go-redis/v9"
	"os"
	"strconv"
	"time"
)

type RedisCache struct {
	c *redis.Client
}

func NewRedisClient() *RedisCache {
	db, _ := strconv.ParseInt(os.Getenv("REDIS_DB"), 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       int(db),
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("ping redis error, err:[%s]", err.Error())
	}
	return &RedisCache{
		c: client,
	}
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.c.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.c.Get(context.Background(), key).Result()
}

func (r *RedisCache) Del(key string) error {
	return r.c.Del(context.Background(), key).Err()
}

func (r *RedisCache) Expire(key string, expiration time.Duration) error {
	return r.c.Expire(context.Background(), key, expiration).Err()
}

func (r *RedisCache) Exists(key string) bool {
	_, err := r.c.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}
	return true
}
