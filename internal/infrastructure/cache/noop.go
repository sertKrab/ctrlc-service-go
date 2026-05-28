package cache

import (
	"context"
	"time"
)

type NoopCache struct{}

func (n *NoopCache) Set(_ context.Context, _ string, _ interface{}, _ time.Duration) error {
	return nil
}
func (n *NoopCache) Get(_ context.Context, _ string) (string, error) {
	return "", ErrCacheMiss
}
func (n *NoopCache) Delete(_ context.Context, _ string) error { return nil }
func (n *NoopCache) Exists(_ context.Context, _ string) (bool, error) { return false, nil }
