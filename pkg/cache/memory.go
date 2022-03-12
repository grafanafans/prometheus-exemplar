package cache

import (
	"context"
	"math/rand"
	"sync"
	"time"

	"github.com/opentracing/opentracing-go"
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
	var (
		hit        bool
		parentSpan = opentracing.SpanFromContext(ctx)
	)

	if parentSpan != nil {
		childSpan := parentSpan.Tracer().StartSpan(
			"MemoryCache.get",
			opentracing.ChildOf(parentSpan.Context()))
		childSpan.SetTag("key", key)
		defer func() {
			childSpan.SetTag("hit", hit)
			childSpan.Finish()
		}()
	}

	// 3% with 200ms sleep and return nil
	if rand.Intn(100) <= 3 {
		time.Sleep(200 * time.Millisecond)
		hit = false
		return nil
	}
	hit = true

	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.items[key]
}

func (c *MemoryCache) Set(key string, item interface{}) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.items[key] = item
	return nil
}
