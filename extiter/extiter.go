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
)

// Count returns the number of elements in the sequence
func Count[T any](seq iter.Seq[T]) int {
	count := 0
	for range seq {
		count++
	}
	return count
}

// IsEmpty returns true if the sequence contains no elements
func IsEmpty[T any](seq iter.Seq[T]) bool {
	for range seq {
		return false
	}
	return true
}

// ToSet collects all elements of the sequence into a set (map[any]struct{})
func ToSet[T any](seq iter.Seq[T]) map[any]struct{} {
	result := make(map[any]struct{})
	for value := range seq {
		result[value] = struct{}{}
	}
	return result
}

// ToMap collects all elements into a map, using the provided functions to determine keys and values
// If keyFn is nil, the sequence elements are used as keys
// If valueFn is nil, the sequence elements are used as values
func ToMap[T any](seq iter.Seq[T], keyFn func(T) any, valueFn func(T) any) map[any]any {
	result := make(map[any]any)
	for value := range seq {
		var key, mapValue any
		if keyFn != nil {
			key = keyFn(value)
		} else {
			key = value
		}
		if valueFn != nil {
			mapValue = valueFn(value)
		} else {
			mapValue = value
		}
		result[key] = mapValue
	}
	return result
}

// Join concatenates all elements into a string, separated by the specified separator
// If separator is empty, elements are separated by commas
func Join[T any](seq iter.Seq[T], separator string) string {
	if separator == "" {
		separator = ","
	}

	var parts []string
	for value := range seq {
		parts = append(parts, toString(value))
	}
	return strings.Join(parts, separator)
}

// IndexOf returns the index of the first occurrence of a value, or -1 if not found
// If fromIndex is provided, the search starts from that position
func IndexOf[T comparable](seq iter.Seq[T], searchElement T, fromIndex ...int) int {
	startIndex := 0
	if len(fromIndex) > 0 {
		startIndex = fromIndex[0]
	}

	index := 0
	for value := range seq {
		if index >= startIndex && value == searchElement {
			return index
		}
		index++
	}
	return -1
}

// Every returns true if all elements satisfy the predicate
func Every[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	for value := range seq {
		if !predicate(value) {
			return false
		}
	}
	return true
}

// Any returns true if any element satisfies the predicate
func Any[T any](seq iter.Seq[T], predicate func(T) bool) bool {
	for value := range seq {
		if predicate(value) {
			return true
		}
	}
	return false
}

// ForEach performs the specified action for each element
func ForEach[T any](seq iter.Seq[T], action func(T, int)) {
	index := 0
	for value := range seq {
		action(value, index)
		index++
	}
}

// Map returns a sequence that yields the results of applying the function to each element
func Map[T, U any](seq iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for value := range seq {
			if !yield(fn(value)) {
				return
			}
		}
	}
}

// Filter returns a sequence containing only elements that satisfy the predicate
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

// NonNullable returns a sequence containing only non-zero values
func NonNullable[T any](seq iter.Seq[T]) iter.Seq[T] {
	return Filter(seq, func(value T) bool {
		return !isZero(value)
	})
}

// Reduce applies a function to each element, accumulating the result
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

// ReduceWithInitial applies a function to each element with an initial value
func ReduceWithInitial[T, U any](seq iter.Seq[T], fn func(U, T) U, initialValue U) U {
	result := initialValue
	for value := range seq {
		result = fn(result, value)
	}
	return result
}

// ReduceRight applies a function to each element in reverse order
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

// ReduceRightWithInitial applies a function to each element in reverse order with an initial value
func ReduceRightWithInitial[T, U any](seq iter.Seq[T], fn func(U, T) U, initialValue U) U {
	elements := slices.Collect(seq)
	result := initialValue
	for i := len(elements) - 1; i >= 0; i-- {
		result = fn(result, elements[i])
	}
	return result
}

// Find returns the first element that satisfies the predicate, or the zero value and false
func Find[T any](seq iter.Seq[T], predicate func(T) bool) (T, bool) {
	for value := range seq {
		if predicate(value) {
			return value, true
		}
	}
	var zero T
	return zero, false
}

// FindIndex returns the index of the first element that satisfies the predicate, or -1
func FindIndex[T any](seq iter.Seq[T], predicate func(T) bool) int {
	index := 0
	for value := range seq {
		if predicate(value) {
			return index
		}
		index++
	}
	return -1
}

// Contains returns true if the sequence contains the specified element
func Contains[T comparable](seq iter.Seq[T], element T) bool {
	for value := range seq {
		if value == element {
			return true
		}
	}
	return false
}

// FlatMap applies a function to each element and flattens the results
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

// Head returns the first element, or the zero value and false if empty
func Head[T any](seq iter.Seq[T]) (T, bool) {
	for value := range seq {
		return value, true
	}
	var zero T
	return zero, false
}

// Tail returns a sequence that skips the first n elements
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

// Limit returns a sequence limited to the specified number of elements
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

// Distinct returns a sequence containing only unique elements
// If keyFn is provided, it's used to determine uniqueness
func Distinct[T any](seq iter.Seq[T], keyFn func(T) any) iter.Seq[T] {
	return func(yield func(T) bool) {
		seen := make(map[any]struct{})
		for value := range seq {
			var key any
			if keyFn != nil {
				key = keyFn(value)
			} else {
				key = value
			}
			
			if _, exists := seen[key]; !exists {
				seen[key] = struct{}{}
				if !yield(value) {
					return
				}
			}
		}
	}
}

// Exclude returns a sequence containing elements that don't exist in the other sequence
// If keyFn is provided, it's used to determine equality
func Exclude[T any](seq iter.Seq[T], other iter.Seq[T], keyFn func(T) any) iter.Seq[T] {
	// Collect keys from the other sequence
	otherKeySet := make(map[any]struct{})
	for value := range other {
		key := keyFn(value)
		otherKeySet[key] = struct{}{}
	}

	return Filter(seq, func(e T) bool {
		ownKey := keyFn(e)
		_, exists := otherKeySet[ownKey]
		return !exists
	})
}

// Concat combines multiple sequences into one
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
	if stringer, ok := item.(interface{ String() string }); ok {
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