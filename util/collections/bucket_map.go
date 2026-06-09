// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package collections

type Comparable[T any] interface {
	Hash() uint64
	Equals(other T) bool
}

type pair[K Comparable[K], V any] struct {
	key   K
	value V
}

type BucketMap[K Comparable[K], V any] struct {
	data map[uint64][]pair[K, V]
}

func NewBucketMap[K Comparable[K], V any]() *BucketMap[K, V] {
	return &BucketMap[K, V]{
		data: make(map[uint64][]pair[K, V]),
	}
}

// Len returns the number of key-value pairs in the map.
func (m *BucketMap[K, V]) Len() int {
	length := 0
	for _, bucket := range m.data {
		length += len(bucket)
	}
	return length
}

// Get returns the value for key if it exists, or the zero value of V if not.
// The boolean is true if the key exists, and false if it does not.
func (m *BucketMap[K, V]) Get(key K) (V, bool) {
	hash := key.Hash()
	bucket, exists := m.data[hash]
	var zero V
	if !exists {
		return zero, false
	}
	for _, pair := range bucket {
		if pair.key.Equals(key) {
			return pair.value, true
		}
	}
	return zero, false
}

// GetOrInsert returns the existing value for key if it exists, or inserts and
// returns value if not. The boolean is true if the value was inserted, and false
// if the key already existed.
func (m *BucketMap[K, V]) GetOrInsert(key K, value V) (V, bool) {
	hash := key.Hash()
	bucket := m.data[hash]
	for _, pair := range bucket {
		if pair.key.Equals(key) {
			return pair.value, false // already exists, return existing value
		}
	}
	m.data[hash] = append(bucket, pair[K, V]{key: key, value: value})
	return value, true
}

// Set updates the value for key if it exists, or inserts it if not. Returns true
// if the value was successfully set, and false if the key already existed.
func (m *BucketMap[K, V]) Set(key K, value V) bool {
	_, exists := m.GetOrInsert(key, value)
	return exists
}

// Has returns true if key exists in the map, and false otherwise.
func (m *BucketMap[K, V]) Has(key K) bool {
	hash := key.Hash()
	bucket, exists := m.data[hash]
	if !exists {
		return false
	}
	for _, pair := range bucket {
		if pair.key.Equals(key) {
			return true
		}
	}
	return false
}

// Remove deletes key from the map if it exists, returning true if the key was removed
func (m *BucketMap[K, V]) Remove(key K) bool {
	hash := key.Hash()
	bucket, exists := m.data[hash]
	if !exists {
		return false
	}
	for i, pair := range bucket {
		if pair.key.Equals(key) {
			// Remove value from bucket
			m.data[hash] = append(bucket[:i], bucket[i+1:]...)
			if len(m.data[hash]) == 0 {
				delete(m.data, hash) // clean up empty bucket
			}
			return true
		}
	}
	return false
}

type BucketSet[T Comparable[T]] struct {
	impl *BucketMap[T, struct{}]
}

func NewBucketSet[T Comparable[T]]() *BucketSet[T] {
	return &BucketSet[T]{
		impl: NewBucketMap[T, struct{}](),
	}
}

// Len returns the number of elements in the set.
func (m *BucketSet[T]) Len() int {
	return m.impl.Len()
}

// Add inserts value into the set, returning true if it was added and false if it already existed.
func (m *BucketSet[T]) Add(value T) bool {
	return m.impl.Set(value, struct{}{})
}

// Has returns true if value exists in the set, and false otherwise.
func (m *BucketSet[T]) Has(value T) bool {
	return m.impl.Has(value)
}

// Remove deletes value from the set if it exists, returning true if the value was removed.
func (m *BucketSet[T]) Remove(value T) bool {
	return m.impl.Remove(value)
}
