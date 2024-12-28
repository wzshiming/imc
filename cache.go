package imc

import (
	"sync"
	"time"
)

// Cache wraps UnsafeCache with mutex for thread safety
type Cache[K comparable, T any] struct {
	mu    sync.RWMutex
	cache *UnsafeCache[K, T]
}

// NewCache creates a new thread-safe cache
func NewCache[K comparable, T any]() *Cache[K, T] {
	return &Cache[K, T]{
		cache: NewUnsafeCache[K, T](),
	}
}

// Set adds or updates a key-value pair in the cache
func (c *Cache[K, T]) Set(key K, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Set(key, value)
}

// SetWithTTL adds or updates a key-value pair with an expiration time
func (c *Cache[K, T]) SetWithTTL(key K, value T, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.SetWithTTL(key, value, ttl)
}

// Get retrieves a value from the cache
func (c *Cache[K, T]) Get(key K) (T, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Get(key)
}

// Remove deletes a key from the cache
func (c *Cache[K, T]) Remove(key K) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.cache.Remove(key)
}

// Len returns the number of items in the cache.
func (c *Cache[K, T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cache.Len()
}

// Evict removes expired items from the cache
func (c *Cache[K, T]) Evict(yield func(key K, value T) bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache.Evict(yield)
}

// Iter iterates over all items in the cache
func (c *Cache[K, T]) Iter(yield func(key K, value T) bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	c.cache.Iter(yield)
}
