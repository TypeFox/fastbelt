// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package collections

import (
	"iter"
	"maps"
)

// A Set is an unordered collection of distinct values. It is a named map type
// that wraps the map[T]struct{} idiom behind [Set.Add], [Set.Has], and
// [Set.Remove] so call sites read as membership operations rather than map
// bookkeeping.
//
// Obtain a writable Set from [NewSet] or a composite literal: the nil zero value
// can be read but panics on [Set.Add], like any nil map. Set is not safe for
// concurrent use by multiple goroutines.
type Set[T comparable] map[T]struct{}

const setMinCapacity = 8

// NewSet returns a Set containing values, with duplicates collapsed.
// The backing map is pre-sized for at least len(values) entries, but never
// below setMinCapacity, so empty sets and small seeds avoid early growth.
func NewSet[T comparable](values ...T) Set[T] {
	s := make(Set[T], max(len(values), setMinCapacity))
	for _, value := range values {
		s[value] = struct{}{}
	}
	return s
}

// Add inserts value into the set. Adding a value already present is a no-op.
func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

// Has reports whether value is a member of the set.
func (s Set[T]) Has(value T) bool {
	_, ok := s[value]
	return ok
}

// Remove deletes value from the set. Removing an absent value is a no-op.
func (s Set[T]) Remove(value T) {
	delete(s, value)
}

// Len returns the number of values in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// All iterates the values in the set. Iteration order is unspecified.
func (s Set[T]) All() iter.Seq[T] {
	return maps.Keys(s)
}
