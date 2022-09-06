package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type MemoryCache struct {
	items map[string]interface{}
	lock  sync.RWMutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]interface{}, 1024),
	}
}

func (c *MemoryCache) Get(key string, ctx context.Context) interface{} {
	return getWithOtel(ctx, "MemoryCache.get", key, func() (bool, interface{}) {
		// 3% with 200ms sleep and return nil
		if rand.Intn(100) <= 3 {
			time.Sleep(200 * time.Millisecond)
			return false, nil
		}

		c.lock.RLock()
		defer c.lock.RUnlock()
		return true, c.items[key]
	})
}

func (c *MemoryCache) Set(key string, item interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[key] = item
	return nil
}
