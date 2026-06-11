// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package collections

import (
	"iter"
	"maps"
)

// A MultiMap is a map from keys to ordered lists of values.
// Each key can hold zero or more values; [Put] and [PutAll] append to a key's
// list and preserve insertion order within that list.
//
// MultiMap is not safe for concurrent use by multiple goroutines.
type MultiMap[K comparable, V any] interface {
	// Get returns the values stored for key in insertion order.
	// The result is nil when key is absent.
	// Callers must not modify the returned slice.
	Get(key K) []V
	// Has reports whether key has at least one stored value.
	Has(key K) bool
	// TryGet returns the values stored for key and reports whether key is present.
	// The result is nil when key is absent.
	// Callers must not modify the returned slice.
	TryGet(key K) ([]V, bool)
	// Put appends value to the list for key.
	Put(key K, value V)
	// PutAll appends values to the list for key.
	PutAll(key K, values []V)
	// Remove deletes key and all of its values.
	Remove(key K)
	// RemoveWhen deletes values for key that satisfy predicate.
	// If no values remain for key, the key is removed.
	RemoveWhen(key K, predicate func(V) bool)
	// Clear removes all keys and values.
	Clear()
	// Keys iterates the distinct keys in the multimap.
	// Iteration order is unspecified.
	Keys() iter.Seq[K]
	// Values iterates every stored value, flattened across keys.
	// Order within a key follows insertion order; order across keys is unspecified.
	Values() iter.Seq[V]
	// All iterates every stored key-value pair.
	// Each value is yielded separately with its key.
	// Order within a key follows insertion order; order across keys is unspecified.
	All() iter.Seq2[K, V]
	// Groups iterates each key with its full value slice.
	// Iteration order is unspecified.
	Groups() iter.Seq2[K, []V]
	// Size returns the total number of stored values across all keys.
	Size() int
}

// NewMultiMap returns an empty [MultiMap].
func NewMultiMap[K comparable, V any]() MultiMap[K, V] {
	return &multiMap[K, V]{
		data: make(map[K][]V),
	}
}

type multiMap[K comparable, V any] struct {
	data map[K][]V
	// Store size to allow constant time retrieval
	size int
}

func (m *multiMap[K, V]) Get(key K) []V {
	return m.data[key]
}

func (m *multiMap[K, V]) Has(key K) bool {
	_, ok := m.data[key]
	return ok
}

func (m *multiMap[K, V]) TryGet(key K) ([]V, bool) {
	values, exists := m.data[key]
	return values, exists
}

func (m *multiMap[K, V]) Put(key K, value V) {
	m.data[key] = append(m.data[key], value)
	m.size++
}

func (m *multiMap[K, V]) PutAll(key K, values []V) {
	m.data[key] = append(m.data[key], values...)
	m.size += len(values)
}

func (m *multiMap[K, V]) Remove(key K) {
	if values, exists := m.data[key]; exists {
		m.size -= len(values)
		delete(m.data, key)
	}
}

func (m *multiMap[K, V]) RemoveWhen(key K, predicate func(V) bool) {
	if values, exists := m.data[key]; exists {
		// Generate a new slice with values that do not match the predicate
		newValues := make([]V, 0, len(values))
		for _, value := range values {
			if !predicate(value) {
				newValues = append(newValues, value)
			}
		}
		m.size -= len(values) - len(newValues)
		if len(newValues) == 0 {
			delete(m.data, key)
		} else {
			m.data[key] = newValues
		}
	}
}

func (m *multiMap[K, V]) Clear() {
	m.data = make(map[K][]V)
	m.size = 0
}

func (m *multiMap[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.data)
}

func (m *multiMap[K, V]) Values() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, values := range m.data {
			for _, value := range values {
				if !yield(value) {
					return
				}
			}
		}
	}
}

func (m *multiMap[K, V]) All() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for key, values := range m.data {
			for _, value := range values {
				if !yield(key, value) {
					return
				}
			}
		}
	}
}

func (m *multiMap[K, V]) Groups() iter.Seq2[K, []V] {
	return maps.All(m.data)
}

func (m *multiMap[K, V]) Size() int {
	return m.size
}
