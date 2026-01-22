package utils

import (
	"iter"
	"maps"
)

type MultiMap[K comparable, V any] interface {
	Get(key K) []V
	Has(key K) bool
	TryGet(key K) ([]V, bool)
	Put(key K, value V)
	PutAll(key K, values []V)
	Remove(key K)
	RemoveWhen(key K, predicate func(V) bool)
	Clear()
	Keys() iter.Seq[K]
	Values() iter.Seq[V]
	All() iter.Seq2[K, V]
	Groups() iter.Seq2[K, []V]
	Size() int
}

func NewMultiMap[K comparable, V any]() MultiMap[K, V] {
	return &MultiMapImpl[K, V]{
		data: make(map[K][]V),
	}
}

type MultiMapImpl[K comparable, V any] struct {
	data map[K][]V
	// Store size to allow constant time retrieval
	size int
}

func (m *MultiMapImpl[K, V]) Get(key K) []V {
	return m.data[key]
}

func (m *MultiMapImpl[K, V]) Has(key K) bool {
	_, ok := m.data[key]
	return ok
}

func (m *MultiMapImpl[K, V]) TryGet(key K) ([]V, bool) {
	values, exists := m.data[key]
	return values, exists
}

func (m *MultiMapImpl[K, V]) Put(key K, value V) {
	m.data[key] = append(m.data[key], value)
	m.size++
}

func (m *MultiMapImpl[K, V]) PutAll(key K, values []V) {
	m.data[key] = append(m.data[key], values...)
	m.size += len(values)
}

func (m *MultiMapImpl[K, V]) Remove(key K) {
	if values, exists := m.data[key]; exists {
		m.size -= len(values)
		delete(m.data, key)
	}
}

func (m *MultiMapImpl[K, V]) RemoveWhen(key K, predicate func(V) bool) {
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

func (m *MultiMapImpl[K, V]) Clear() {
	m.data = make(map[K][]V)
	m.size = 0
}

func (m *MultiMapImpl[K, V]) Keys() iter.Seq[K] {
	return maps.Keys(m.data)
}

func (m *MultiMapImpl[K, V]) Values() iter.Seq[V] {
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

func (m *MultiMapImpl[K, V]) All() iter.Seq2[K, V] {
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

func (m *MultiMapImpl[K, V]) Groups() iter.Seq2[K, []V] {
	return maps.All(m.data)
}

func (m *MultiMapImpl[K, V]) Size() int {
	return m.size
}
