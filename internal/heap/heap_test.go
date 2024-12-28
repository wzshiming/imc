package heap

import (
	"testing"
)

func TestHeap(t *testing.T) {
	h := NewHeap[int, string]()

	// Test initial state
	if h.Len() != 0 {
		t.Errorf("expected empty heap, got len=%d", h.Len())
	}

	// Test Push and Len
	h.Push(3, "three")
	h.Push(1, "one")
	h.Push(2, "two")

	if h.Len() != 3 {
		t.Errorf("expected len=3, got len=%d", h.Len())
	}

	// Test Peek
	k, v, ok := h.Peek()
	if !ok || k != 1 || v != "one" {
		t.Errorf("Peek: expected (1, one, true), got (%v, %v, %v)", k, v, ok)
	}

	// Test Pop order
	expectedPops := []struct {
		key  int
		data string
	}{
		{1, "one"},
		{2, "two"},
		{3, "three"},
	}

	for _, exp := range expectedPops {
		k, v, ok := h.Pop()
		if !ok || k != exp.key || v != exp.data {
			t.Errorf("Pop: expected (%v, %v, true), got (%v, %v, %v)",
				exp.key, exp.data, k, v, ok)
		}
	}

	// Test Pop empty
	_, _, ok = h.Pop()
	if ok {
		t.Error("Pop: expected false on empty heap")
	}
}

func TestHeapRemove(t *testing.T) {
	h := NewHeap[int, string]()
	h.Push(1, "one")
	h.Push(2, "two")
	h.Push(3, "three")

	// Remove middle element
	if !h.Remove("two") {
		t.Error("Remove: expected true when removing existing element")
	}

	// Verify remaining elements
	if h.Len() != 2 {
		t.Errorf("expected len=2 after remove, got len=%d", h.Len())
	}

	k, v, ok := h.Pop()
	if !ok || k != 1 || v != "one" {
		t.Errorf("Pop after remove: expected (1, one, true), got (%v, %v, %v)", k, v, ok)
	}

	// Try to remove non-existent element
	if h.Remove("not-exists") {
		t.Error("Remove: expected false when removing non-existent element")
	}
}
