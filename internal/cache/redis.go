package cache

import (
	"time"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq" // postgres
)

// Redis contains objects for database communication.
type Redis struct {
	*redis.Client
}

// NewRedis create new cache instance.
func NewRedis(config *RedisConfig) (*Redis, error) {

	option := &redis.Options{
		Addr: config.Addr,
		DB:   1, // use default DB
	}

	if config.Password != "" {
		option.Password = config.Password
	}

	rdb := redis.NewClient(option)

	err := rdb.Ping().Err()
	if err != nil {
		return nil, err
	}

	return &Redis{rdb}, nil
}

//SetValue run sql query
func (cache *Redis) SetValue(key string, value interface{}, expiration time.Duration) error {
	err := cache.Set(key, value, expiration).Err()
	return err
}

//GetValue run sql query
func (cache *Redis) GetValue(key string) (string, error) {
	val, err := cache.Get(key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

//ExpireKey .
func (cache *Redis) ExpireKey(key string) error {
	err := cache.Expire(key, 1*time.Second).Err()
	if err != nil {
		return err
	}

	return nil
}

//GetInstance .
func (cache *Redis) GetInstance() *redis.Client {
	return cache.Client
}
