# extiter (Extended Iterators)

The `extiter` package provides extended functionality for Go's built-in `iter.Seq` type, offering functional programming utilities for working with sequences.

## Overview

This package has been refactored from the previous `stream` package to work directly with Go's built-in `iter.Seq` type (available in Go 1.23+), removing the need for custom interfaces and leveraging the standard library.

## Key Changes from `stream` Package

- **No more interfaces**: Works directly with `iter.Seq[T]`
- **Standard library integration**: Use `slices.Values` to create sequences from slices, `slices.Collect` to collect sequences into slices
- **Standalone functions**: All operations are now standalone functions instead of methods
- **Removed redundant functions**: Functions like `FromSlice` and `ToSlice` are replaced by `slices.Values` and `slices.Collect` respectively

## Usage Examples

```go
package main

import (
    "fmt"
    "slices"
    "github.com/TypeFox/langium-to-go/extiter"
)

func main() {
    // Create a sequence from a slice
    numbers := slices.Values([]int{1, 2, 3, 4, 5, 6})
    
    // Chain operations: filter even numbers, map to squares, limit to 3
    result := slices.Collect(
        extiter.Limit(
            extiter.Map(
                extiter.Filter(numbers, func(x int) bool { return x%2 == 0 }),
                func(x int) int { return x * x },
            ),
            3,
        ),
    )
    
    fmt.Println(result) // [4, 16, 36]
    
    // Other useful operations
    seq := slices.Values([]string{"apple", "banana", "cherry"})
    
    // Count elements
    count := extiter.Count(seq)
    fmt.Println("Count:", count) // Count: 3
    
    // Join elements
    joined := extiter.Join(seq, ", ")
    fmt.Println("Joined:", joined) // Joined: apple, banana, cherry
    
    // Find first element matching condition
    result, found := extiter.Find(seq, func(s string) bool { 
        return len(s) > 5 
    })
    if found {
        fmt.Println("Found:", result) // Found: banana
    }
}
```

## Available Functions

### Collecting Operations
- `Count[T](seq iter.Seq[T]) int` - Count elements
- `IsEmpty[T](seq iter.Seq[T]) bool` - Check if sequence is empty
- `ToSet[T](seq iter.Seq[T]) map[any]struct{}` - Collect into set
- `ToMap[T](seq iter.Seq[T], keyFn, valueFn func(T) any) map[any]any` - Collect into map

### Transformation Operations
- `Map[T, U](seq iter.Seq[T], fn func(T) U) iter.Seq[U]` - Transform elements
- `Filter[T](seq iter.Seq[T], predicate func(T) bool) iter.Seq[T]` - Filter elements
- `FlatMap[T, U](seq iter.Seq[T], fn func(T) iter.Seq[U]) iter.Seq[U]` - Flatten nested sequences
- `Distinct[T](seq iter.Seq[T], keyFn func(T) any) iter.Seq[T]` - Remove duplicates

### Sequence Operations
- `Limit[T](seq iter.Seq[T], maxSize int) iter.Seq[T]` - Limit number of elements
- `Tail[T](seq iter.Seq[T], n int) iter.Seq[T]` - Skip first n elements
- `Concat[T](sequences ...iter.Seq[T]) iter.Seq[T]` - Concatenate sequences
- `Exclude[T](seq iter.Seq[T], other iter.Seq[T], keyFn func(T) any) iter.Seq[T]` - Exclude elements

### Search Operations
- `Find[T](seq iter.Seq[T], predicate func(T) bool) (T, bool)` - Find first matching element
- `FindIndex[T](seq iter.Seq[T], predicate func(T) bool) int` - Find index of first match
- `Contains[T comparable](seq iter.Seq[T], element T) bool` - Check if contains element
- `IndexOf[T comparable](seq iter.Seq[T], searchElement T, fromIndex ...int) int` - Find index of element

### Aggregation Operations
- `Reduce[T](seq iter.Seq[T], fn func(T, T) T) (T, bool)` - Reduce to single value
- `ReduceWithInitial[T, U](seq iter.Seq[T], fn func(U, T) U, initialValue U) U` - Reduce with initial value
- `Every[T](seq iter.Seq[T], predicate func(T) bool) bool` - Check if all elements match
- `Any[T](seq iter.Seq[T], predicate func(T) bool) bool` - Check if any element matches

### Other Operations
- `ForEach[T](seq iter.Seq[T], action func(T, int))` - Execute action for each element
- `Head[T](seq iter.Seq[T]) (T, bool)` - Get first element
- `Join[T](seq iter.Seq[T], separator string) string` - Join elements into string
- `NonNullable[T](seq iter.Seq[T]) iter.Seq[T]` - Filter out zero values

## Migration from `stream` Package

```go
// Old stream package
stream := stream.FromSlice([]int{1, 2, 3})
result := stream.Map(func(x int) any { return x * 2 }).ToSlice()

// New extiter package
seq := slices.Values([]int{1, 2, 3})
result := slices.Collect(extiter.Map(seq, func(x int) int { return x * 2 }))
```

## Requirements

- Go 1.23+ (for `iter` package support)