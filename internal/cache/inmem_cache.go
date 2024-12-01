package cache

import (
	"fmt"
	"sync"
	"time"
)

type InMemoryCache struct {
	mu    sync.RWMutex
	store map[string]cacheItem
}

type cacheItem struct {
	value      string
	expiration time.Time
}

func NewInMemoryCache() *InMemoryCache {
	return &InMemoryCache{
		store: make(map[string]cacheItem),
	}
}

func (c *InMemoryCache) Get(key string) (string, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.store[key]
	if !found || time.Now().After(item.expiration) {
		return "", fmt.Errorf("key not found or expired")
	}
	return item.value, nil
}

func (c *InMemoryCache) Set(key, value string, ttl int) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.store[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(time.Duration(ttl) * time.Second),
	}
	return nil
}

func (c *InMemoryCache) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.store, key)
	return nil
}
