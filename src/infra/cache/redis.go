package cache

import (
	"TTCS/src/common/configs"
	"TTCS/src/common/log"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisCache struct {
	c *redis.Client
}

func NewRedisClient() *redis.Client {
	redisConfigs := configs.GetConfig().Redis
	client := redis.NewClient(&redis.Options{
		Addr:     redisConfigs.Address,
		Password: redisConfigs.Password,
		DB:       redisConfigs.DB,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		log.Fatal("ping redis error, err:[%s]", err.Error())
	}
	return client
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	return r.c.Set(context.Background(), key, value, expiration).Err()
}

func (r *RedisCache) Get(key string) (interface{}, error) {
	return r.c.Get(context.Background(), key).Result()
}
