package cache

import (
	"log"

	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
)

func NewCache(cfg *config.Config) Cache {
	if !cfg.EnableRedis {
		log.Println("[Cache] Redis disabled — using noop cache")
		return &NoopCache{}
	}
	client, err := newRedisClient(cfg)
	if err != nil {
		log.Printf("[Cache] Redis connect failed, falling back to noop: %v", err)
		return &NoopCache{}
	}
	log.Println("[Cache] Redis connected")
	return &RedisCache{client: client}
}
