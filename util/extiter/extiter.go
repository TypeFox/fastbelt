// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package extiter

import (
	"fmt"
	"iter"
	"reflect"
	"slices"
	"strings"

	"typefox.dev/fastbelt/util/collections"
)

// Empty returns an empty [iter.Seq] of type T.
//
// The returned sequence yields no elements. It can be reused as a stable
// sentinel, for example when an API must return a sequence but has nothing
// to enumerate.
func Empty[T any]() iter.Seq[T] {
	return func(yield func(T) bool) {}
}

// Of returns a sequence over the given elements in order.
func Of[T any](elements ...T) iter.Seq[T] {
	return slices.Values(elements)
}

// Count returns the number of elements in seq.
//
// Count fully consumes seq.
func Count[T any](seq iter.Seq[T]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

// IsEmpty reports whether seq yields no elements.
//
// A nil sequence is treated as empty.
func IsEmpty[T any](seq iter.Seq[T]) bool {
	if seq == nil {
		return true
	}
	for range seq {
		return false
	}
	return true
}

// Join concatenates the string form of each element in seq, separated by separator.
//
// Strings pass through unchanged. Values implementing [fmt.Stringer] use
// String(). Numeric and boolean types use default formatting. Nil is rendered
// as "nil"; other types use a default %v representation.
func Join[T any](seq iter.Seq[T], separator string) string {
	var parts []string
	for value := range seq {
		parts = append(parts, toString(value))
	}
	return strings.Join(parts, separator)
}

// IndexOf returns the index of the first element equal to searchElement, or -1
// if seq does not contain it.
func IndexOf[T comparable](seq iter.Seq[T], searchElement T) int {
	index := 0
	for value := range seq {
		if value == searchElement {
			return index
		}
		index++
	}
	return -1
}

// Every reports whether predicate holds for every element in seq.
//
// An empty sequence satisfies Every.
func Every[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	for value := range seq {
		if !predicate(value) {
			return false
		}
	}
	return true
}

// Any reports whether predicate holds for at least one element in seq.
//
// An empty sequence does not satisfy Any.
func Any[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	for value := range seq {
		if predicate(value) {
			return true
		}
	}
	return false
}

// ForEach calls action once for each element in seq.
//
// The second argument to action is the zero-based index of the element.
func ForEach[T any](seq iter.Seq[T], action func(T, int)) {
	index := 0
	for value := range seq {
		action(value, index)
		index++
	}
}

// Map returns a lazy sequence that applies fn to each element of seq.
//
// Iteration stops early when the consumer stops the sequence.
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for value := range seq {
			if !yield(fn(value)) {
				return
			}
		}
	}
}

// Filter returns a lazy sequence of elements in seq for which predicate is true.
//
// Iteration stops early when the consumer stops the sequence.
func Filter[T any](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for value := range seq {
			if predicate(value) {
				if !yield(value) {
					return
				}
			}
		}
	}
}

// FilterType returns a lazy sequence of elements in seq that have dynamic type U.
//
// Each element is type-asserted to U; elements that do not assert are skipped.
// Iteration stops early when the consumer stops the sequence.
func FilterType[T, U any](seq iter.Seq[T]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for value := range seq {
			if castValue, ok := any(value).(U); ok {
				if !yield(castValue) {
					return
				}
			}
		}
	}
}

// NonZero returns a lazy sequence of elements in seq that are not the zero
// value of their type.
//
// Numeric zero, false booleans, empty strings, nil pointers and interfaces,
// and other types' zero values as determined by [reflect.Value.IsZero] are
// removed.
func NonZero[T any](seq iter.Seq[T]) iter.Seq[T] {
	return Filter(seq, func(value T) bool {
		return !isZero(value)
	})
}

// Reduce folds seq with fn, using the first element as the initial accumulator.
//
// The returned bool is false when seq is empty; otherwise it is true and the
// result is the folded value. For a single element, that element is returned
// unchanged.
func Reduce[T any](seq iter.Seq[T], fn func(T, T) T) (T, bool) {
	var result T
	found := false

	for value := range seq {
		if !found {
			result = value
			found = true
		} else {
			result = fn(result, value)
		}
	}
	return result, found
}

// ReduceWithInitial folds seq into initialValue using fn.
//
// An empty sequence returns initialValue unchanged.
func ReduceWithInitial[T, U any](seq iter.Seq[T], fn func(U, T) U, initialValue U) U {
	result := initialValue
	for value := range seq {
		result = fn(result, value)
	}
	return result
}

// ReduceRight folds seq from right to left with fn.
//
// The sequence is materialized before folding. The last element becomes the
// initial accumulator; fn is then applied to the accumulator and each
// preceding element in turn. The returned bool is false when seq is empty.
func ReduceRight[T any](seq iter.Seq[T], fn func(T, T) T) (T, bool) {
	// Collect all elements first since we need to iterate in reverse
	elements := slices.Collect(seq)
	if len(elements) == 0 {
		var zero T
		return zero, false
	}

	result := elements[len(elements)-1]
	for i := len(elements) - 2; i >= 0; i-- {
		result = fn(result, elements[i])
	}
	return result, true
}

// ReduceRightWithInitial folds seq from right to left into initialValue using fn.
//
// The sequence is materialized before folding. Elements are combined starting
// with the last element in seq.
func ReduceRightWithInitial[T, U any](seq iter.Seq[T], fn func(U, T) U, initialValue U) U {
	elements := slices.Collect(seq)
	result := initialValue
	for i := len(elements) - 1; i >= 0; i-- {
		result = fn(result, elements[i])
	}
	return result
}

// Find returns the first element in seq for which predicate is true, its
// zero-based index, and a bool indicating whether such an element exists.
//
// When no element matches, Find returns the zero value of T, index -1, and false.
func Find[T any](seq iter.Seq[T], predicate func(T) bool) (T, int, bool) {
	index := 0
	for value := range seq {
		if predicate(value) {
			return value, index, true
		}
		index++
	}
	var zero T
	return zero, -1, false
}

// Contains reports whether seq includes an element equal to element.
func Contains[T comparable](seq iter.Seq[T], element T) bool {
	for value := range seq {
		if value == element {
			return true
		}
	}
	return false
}

// FlatMap returns a lazy sequence formed by applying fn to each element of seq
// and concatenating the resulting sequences in order.
//
// Iteration stops early when the consumer stops the sequence.
func FlatMap[T, U any](seq iter.Seq[T], fn func(T) iter.Seq[U]) iter.Seq[U] {
	return func(yield func(U) bool) {
		for value := range seq {
			for nestedValue := range fn(value) {
				if !yield(nestedValue) {
					return
				}
			}
		}
	}
}

// Head returns the first element of seq and whether one exists.
//
// When seq is empty, Head returns the zero value of T and false.
func Head[T any](seq iter.Seq[T]) (T, bool) {
	for value := range seq {
		return value, true
	}
	var zero T
	return zero, false
}

// Tail returns a lazy sequence of the elements in seq after skipping the first n.
//
// If n is greater than or equal to the length of seq, the result is empty.
func Tail[T any](seq iter.Seq[T], n int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for value := range seq {
			if count >= n {
				if !yield(value) {
					return
				}
			}
			count++
		}
	}
}

// Limit returns a lazy sequence of at most maxSize elements from seq.
//
// If maxSize is zero or negative, the result is empty.
func Limit[T any](seq iter.Seq[T], maxSize int) iter.Seq[T] {
	return func(yield func(T) bool) {
		count := 0
		for value := range seq {
			if count >= maxSize {
				return
			}
			if !yield(value) {
				return
			}
			count++
		}
	}
}

// Distinct returns a lazy sequence of the first occurrence of each distinct
// element in seq.
//
// When keyFn is nil, elements are compared by value. Otherwise keyFn defines
// the key used to decide whether two elements are the same; the first element
// for each key is kept.
func Distinct[T any](seq iter.Seq[T], keyFn func(T) any) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(collections.Set[any])
		for value := range seq {
			var key any
			if keyFn != nil {
				key = keyFn(value)
			} else {
				key = value
			}

			if seen.Add(key) && !yield(value) {
				return
			}
		}
	}
}

// Exclude returns a lazy sequence of elements in seq whose key is not present in other.
//
// The other sequence is fully consumed before filtering begins. When keyFn is
// nil, elements are compared by value; otherwise keyFn defines the comparison key.
func Exclude[T any](seq iter.Seq[T], other iter.Seq[T], keyFn func(T) any) iter.Seq[T] {
	// Collect keys from the other sequence
	otherKeySet := make(collections.Set[any])
	for value := range other {
		var key any
		if keyFn != nil {
			key = keyFn(value)
		} else {
			key = value
		}
		otherKeySet.Add(key)
	}

	return Filter(seq, func(e T) bool {
		var ownKey any
		if keyFn != nil {
			ownKey = keyFn(e)
		} else {
			ownKey = e
		}
		return !otherKeySet.Has(ownKey)
	})
}

// Concat returns a lazy sequence that yields every element of sequences in order.
//
// Iteration stops early when the consumer stops the sequence.
func Concat[T any](sequences ...iter.Seq[T]) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, seq := range sequences {
			for value := range seq {
				if !yield(value) {
					return
				}
			}
		}
	}
}

// Helper functions

func toString(item any) string {
	if item == nil {
		return "nil"
	}
	if str, ok := item.(string); ok {
		return str
	}
	if stringer, ok := item.(fmt.Stringer); ok {
		return stringer.String()
	}
	// Handle basic types
	switch v := item.(type) {
	case int:
		return fmt.Sprintf("%d", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case uint:
		return fmt.Sprintf("%d", v)
	case uint8:
		return fmt.Sprintf("%d", v)
	case uint16:
		return fmt.Sprintf("%d", v)
	case uint32:
		return fmt.Sprintf("%d", v)
	case uint64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%g", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case bool:
		return fmt.Sprintf("%t", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func isZero(value any) bool {
	if value == nil {
		return true
	}

	// Use reflection to check if value is the zero value of its type
	switch v := value.(type) {
	case int:
		return v == 0
	case int8:
		return v == 0
	case int16:
		return v == 0
	case int32:
		return v == 0
	case int64:
		return v == 0
	case uint:
		return v == 0
	case uint8:
		return v == 0
	case uint16:
		return v == 0
	case uint32:
		return v == 0
	case uint64:
		return v == 0
	case float32:
		return v == 0.0
	case float64:
		return v == 0.0
	case string:
		return v == ""
	case bool:
		return !v
	default:
		// For other types, use reflection
		val := reflect.ValueOf(value)
		return val.IsZero()
	}
}
