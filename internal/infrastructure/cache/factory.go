package cache

import (
	"fmt"

	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
)

func NewCache(cfg *config.Config) (Cache, error) {
	if !cfg.EnableRedis {
		return nil, fmt.Errorf("redis is required for durable runtime state; set ENABLE_REDIS=true")
	}
	client, err := newRedisClient(cfg)
	if err != nil {
		return nil, err
	}
	return &RedisCache{client: client}, nil
}
