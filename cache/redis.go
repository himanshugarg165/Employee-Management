package cache

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

type redisClient struct {
	client *redis.Client
}

func NewRedisClient(config CacheConfig) (Cache, error) {
	rds := &redisClient{}
	rds.client = redis.NewClient(&redis.Options{
		Addr:     config.URL,
		Password: config.Password,
		DB:       0,
	})
	err := rds.client.Ping().Err()
	if err != nil {
		panic(err)
	}
	return rds, err
}

func (c *redisClient) Set(key string, value interface{}, expiry time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := c.client.Set(key, data, expiry).Err(); err != nil {
		return err
	}
	return nil
}

func (c *redisClient) Get(key string, value interface{}) error {
	data, err := c.client.Get(key).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	return nil
}

func (c *redisClient) Delete(key string) error {
	return c.client.Del(key).Err()
}

func (c *redisClient) Flush() error {
	return c.client.FlushAll().Err()
}
