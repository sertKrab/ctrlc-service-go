package cache

import (
	"strings"
	"testing"

	"git.trovefin.com/poc/ctrlc-service-go/internal/config"
)

func TestNewCacheRejectsDisabledRedis(t *testing.T) {
	_, err := NewCache(&config.Config{EnableRedis: false})
	if err == nil || !strings.Contains(err.Error(), "ENABLE_REDIS=true") {
		t.Fatalf("NewCache err = %v, want ENABLE_REDIS guidance", err)
	}
}

func TestNoopCacheIsNotDurable(t *testing.T) {
	if (&NoopCache{}).Durable() {
		t.Fatal("NoopCache Durable() = true, want false")
	}
}

func TestNewCacheRejectsRedisConnectFailure(t *testing.T) {
	_, err := NewCache(&config.Config{
		EnableRedis: true,
		RedisHost:   "127.0.0.1",
		RedisPort:   "1",
	})
	if err == nil || !strings.Contains(err.Error(), "redis ping") {
		t.Fatalf("NewCache err = %v, want redis ping failure", err)
	}
}
