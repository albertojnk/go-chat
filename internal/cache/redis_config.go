package cache

import "github.com/albertojnk/go-chat/common"

// RedisConfig contains the configurations for cache.
type RedisConfig struct {
	Addr     string
	Port     string
	Password string
}

// NewRedisConfig load the configuration for database.
func NewRedisConfig() *RedisConfig {
	password := common.GetEnv("REDIS_PASSWORD", "abcd1234")
	if password == "empty" {
		password = ""
	}

	config := &RedisConfig{
		Addr:     common.GetEnv("REDIS_HOST", "localhost:6378"),
		Password: password,
	}

	return config
}
