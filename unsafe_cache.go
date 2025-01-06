package imc

import (
	"time"

	"github.com/wzshiming/imc/internal/heap"
)

// Default function to get current time
var nowFunc = time.Now

// UnsafeCache represents an in-memory cache with optional TTL expiration.
type UnsafeCache[K comparable, T any] struct {
	data       map[K]T              // Map storing key-value pairs
	heap       *heap.Heap[int64, K] // Min-heap for tracking expiration times
	nextExpiry int64                // Next expiration timestamp, -1 if no expiring items
}

// NewUnsafeCache creates a new empty cache instance.
func NewUnsafeCache[K comparable, T any]() *UnsafeCache[K, T] {
	return &UnsafeCache[K, T]{
		data:       map[K]T{},
		heap:       heap.NewHeap[int64, K](),
		nextExpiry: -1,
	}
}

// Set adds or updates a key-value pair in the cache without expiration.
func (c *UnsafeCache[K, T]) Set(key K, value T) {
	c.data[key] = value
}

// SetWithTTL adds or updates a key-value pair with a time-to-live duration.
func (c *UnsafeCache[K, T]) SetWithTTL(key K, value T, ttl time.Duration) {
	expiry := nowFunc().Add(ttl).Unix()
	_, ok := c.data[key]
	if ok {
		_ = c.heap.Remove(key)
	}

	if c.nextExpiry < 0 || expiry < c.nextExpiry {
		c.nextExpiry = expiry
	}
	c.heap.Push(expiry, key)
	c.data[key] = value
}

// Get retrieves a value from the cache by key.
// Returns the value and true if found, zero value and false if not found.
func (c *UnsafeCache[K, T]) Get(key K) (T, bool) {
	item, ok := c.data[key]
	return item, ok
}

// Remove deletes a key-value pair from the cache.
// Returns true if the key was found and removed, false otherwise.
func (c *UnsafeCache[K, T]) Remove(key K) bool {
	_, ok := c.data[key]
	if !ok {
		return false
	}
	delete(c.data, key)
	_ = c.heap.Remove(key)
	return true
}

// Len returns the number of items in the cache.
func (c *UnsafeCache[K, T]) Len() int {
	return len(c.data)
}

// NextExpiry returns the next expiry time in Unix timestamp format.
func (c *UnsafeCache[K, T]) NextExpiry() int64 {
	return c.nextExpiry
}

// Evict removes expired items from the cache.
// The yield function, if provided, is called for each evicted item.
// If yield returns false, eviction stops.
func (c *UnsafeCache[K, T]) Evict(yield func(key K, value T) bool) {
	if c.nextExpiry < 0 {
		return
	}

	currentTime := nowFunc().Unix()
	if currentTime < c.nextExpiry {
		return
	}

	for {
		expiry, _, ok := c.heap.Peek()
		if !ok {
			c.nextExpiry = -1
			return
		}
		if expiry > currentTime {
			c.nextExpiry = expiry
			return
		}
		_, key, _ := c.heap.Pop()
		if yield == nil {
			delete(c.data, key)
			continue
		}
		value := c.data[key]
		delete(c.data, key)
		if yield(key, value) {
			continue
		}
		expiry, _, ok = c.heap.Peek()
		if !ok {
			c.nextExpiry = -1
			return
		}
		c.nextExpiry = expiry
		return
	}
}

// Iter iterates over all items in the cache.
// The yield function is called for each item.
// If yield returns false, iteration stops.
func (c *UnsafeCache[K, T]) Iter(yield func(key K, value T) bool) {
	for k, v := range c.data {
		if !yield(k, v) {
			return
		}
	}
}
