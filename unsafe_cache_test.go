package imc

import (
	"testing"
	"time"
)

func TestUnsafeCacheSetAndGet(t *testing.T) {
	c := NewUnsafeCache[string, int]()
	c.Set("key1", 100)

	val, exists := c.Get("key1")
	if !exists {
		t.Error("expected key to exist")
	}
	if val != 100 {
		t.Errorf("expected value 100, got %d", val)
	}
}

func TestUnsafeCacheNonExistentKey(t *testing.T) {
	c := NewUnsafeCache[string, int]()
	_, exists := c.Get("missing")
	if exists {
		t.Error("expected key to not exist")
	}
}

func TestUnsafeCacheRemoveKey(t *testing.T) {
	c := NewUnsafeCache[string, int]()
	c.Set("key1", 100)

	removed := c.Remove("key1")
	if !removed {
		t.Error("expected Remove to return true")
	}

	_, exists := c.Get("key1")
	if exists {
		t.Error("expected key to be removed")
	}

	// Test removing non-existent key
	removed = c.Remove("missing")
	if removed {
		t.Error("expected Remove to return false for non-existent key")
	}
}

func TestUnsafeCacheEviction(t *testing.T) {
	c := NewUnsafeCache[string, int]()
	mockTime := time.Now()
	nowFunc = func() time.Time { return mockTime }

	c.SetWithTTL("expire1", 100, time.Second)
	c.SetWithTTL("expire2", 200, time.Hour)

	// Advance mock time by 2 seconds
	nowFunc = func() time.Time { return mockTime.Add(2 * time.Second) }
	c.Evict(nil)

	_, exists := c.Get("expire1")
	if exists {
		t.Error("expected expired key to be evicted")
	}

	_, exists = c.Get("expire2")
	if !exists {
		t.Error("expected non-expired key to remain")
	}
}

func TestUnsafeCacheUpdateExistingKey(t *testing.T) {
	c := NewUnsafeCache[string, int]()
	c.SetWithTTL("key1", 100, time.Hour)
	c.SetWithTTL("key1", 200, time.Hour)

	val, exists := c.Get("key1")
	if !exists {
		t.Error("expected key to exist")
	}
	if val != 200 {
		t.Errorf("expected value 200, got %d", val)
	}
}
