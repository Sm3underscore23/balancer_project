package inmemorycache

import (
	"sync"

	"balancer/internal/model"
)

type InMemoryTokenBucketCache struct {
	mu    sync.Mutex
	store map[string]*model.TokenBucket
}

func NewInMemoryTokenBucketCache() *InMemoryTokenBucketCache {
	return &InMemoryTokenBucketCache{
		store: make(map[string]*model.TokenBucket),
	}
}

func (c *InMemoryTokenBucketCache) Get(clientID string) (*model.TokenBucket, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.store[clientID]
	return val, ok
}

func (c *InMemoryTokenBucketCache) Set(clientID string, bucket *model.TokenBucket) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.store[clientID] = bucket
}

func (c *InMemoryTokenBucketCache) Delete(clientID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, clientID)
}

func (c *InMemoryTokenBucketCache) Range(f func(clientID string, bucket *model.TokenBucket)) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.store {
		f(k, v)
	}
}
