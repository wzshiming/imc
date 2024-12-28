package heap

import (
	"cmp"
	"container/heap"
)

// Heap is a generic min-heap implementation that allows efficient access to elements by their data value.
type Heap[K cmp.Ordered, T comparable] struct {
	entries waitEntries[K, T]
	indexes map[T]*waitEntry[K, T]
}

// NewHeap creates a new empty heap.
func NewHeap[K cmp.Ordered, T comparable]() *Heap[K, T] {
	return &Heap[K, T]{
		indexes: map[T]*waitEntry[K, T]{},
	}
}

// Push adds an item to the heap with the given key and data.
func (h *Heap[K, T]) Push(key K, data T) {
	entry := &waitEntry[K, T]{key: key, data: data}
	heap.Push(&h.entries, entry)
	h.indexes[data] = entry
}

// Pop removes and returns the item with the smallest key from the heap.
// Returns zero values and false if the heap is empty.
func (h *Heap[K, T]) Pop() (k K, v T, ok bool) {
	if len(h.entries) == 0 {
		return k, v, false
	}
	item := heap.Pop(&h.entries).(*waitEntry[K, T])
	delete(h.indexes, item.data)
	return item.key, item.data, true
}

// Peek returns the item with the smallest key without removing it.
// Returns zero values and false if the heap is empty.
func (h *Heap[K, T]) Peek() (k K, v T, ok bool) {
	if len(h.entries) == 0 {
		return k, v, false
	}
	item := h.entries[0]
	return item.key, item.data, true
}

// Remove removes the item with the given data value if it exists.
// Returns true if an item was removed, false if not found.
func (h *Heap[K, T]) Remove(data T) bool {
	if item, ok := h.indexes[data]; ok && item.index >= 0 {
		heap.Remove(&h.entries, item.index)
		delete(h.indexes, data)
		return true
	}
	return false
}

// Len returns the number of items in the heap.
func (h *Heap[K, T]) Len() int {
	return h.entries.Len()
}

// waitEntry represents an item in the heap with its key, data and index.
type waitEntry[K cmp.Ordered, T any] struct {
	key   K
	data  T
	index int
}

// waitEntries implements heap.Interface for a slice of waitEntry.
type waitEntries[K cmp.Ordered, T any] []*waitEntry[K, T]

func (w waitEntries[K, T]) Len() int {
	return len(w)
}

func (w waitEntries[K, T]) Less(i, j int) bool {
	return cmp.Less[K](w[i].key, w[j].key)
}

func (w waitEntries[K, T]) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
	w[i].index = i
	w[j].index = j
}

// Push implements heap.Interface. Do not call directly, use heap.Push instead.
func (w *waitEntries[K, T]) Push(x any) {
	n := len(*w)
	item := x.(*waitEntry[K, T])
	item.index = n
	*w = append(*w, item)
}

// Pop implements heap.Interface. Do not call directly, use heap.Pop instead.
func (w *waitEntries[K, T]) Pop() any {
	n := len(*w)
	item := (*w)[n-1]
	item.index = -1
	*w = (*w)[0:(n - 1)]
	return item
}
