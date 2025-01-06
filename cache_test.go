package imc

import (
	"context"
	"testing"
	"time"
)

func TestRunEvict(t *testing.T) {
	cache := NewCache[string, int]()
	cache.Set("key1", 1)
	cache.SetWithTTL("key2", 2, 1*time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go cache.RunEvict(ctx, func(key string, value int) bool {
		return true
	})

	time.Sleep(2 * time.Second) // Wait for eviction to occur

	if _, found := cache.Get("key1"); !found {
		t.Errorf("Expected key1 to be present, but it was evicted")
	}

	if _, found := cache.Get("key2"); found {
		t.Errorf("Expected key2 to be evicted, but it was still present")
	}
}
