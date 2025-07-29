# Langium to Go

[![CI](https://github.com/TypeFox/langium-to-go/actions/workflows/ci.yml/badge.svg)](https://github.com/TypeFox/langium-to-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/TypeFox/langium-to-go)](https://goreportcard.com/report/github.com/TypeFox/langium-to-go)
[![Coverage](https://codecov.io/gh/TypeFox/langium-to-go/branch/main/graph/badge.svg)](https://codecov.io/gh/TypeFox/langium-to-go)

A Go implementation of Langium's streaming utilities, providing functional programming patterns for processing sequences of data.

## Features

### Stream Package

The `stream` package provides a lazy evaluation stream implementation with methods similar to TypeScript/JavaScript arrays but optimized for Go:

- **Lazy Evaluation**: Operations are only executed when the stream is consumed
- **Method Chaining**: Fluent interface for combining operations
- **Type Safety**: Full Go generics support
- **Memory Efficient**: No intermediate collections unless explicitly materialized

## Installation

```bash
go get github.com/TypeFox/langium-to-go
```

## Usage

### Basic Stream Operations

```go
package main

import (
    "fmt"
    "github.com/TypeFox/langium-to-go/stream"
)

func main() {
    // Create a stream from a slice
    numbers := stream.FromSlice([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
    
    // Chain operations: filter even numbers, square them, take first 3
    result := numbers.
        Filter(func(x int) bool { return x%2 == 0 }).
        Map(func(x int) any { return x * x }).
        ToSlice()
    
    fmt.Println(result) // Output: [4, 16, 36, 64, 100]
}
```

### Advanced Examples

```go
// Working with custom types
type Person struct {
    Name string
    Age  int
}

people := stream.FromSlice([]Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Charlie", 35},
})

// Find adults over 30
adults := people.
    Filter(func(p Person) bool { return p.Age > 30 }).
    ToSlice()

// Create a name-to-age mapping
nameToAge := people.ToMap(
    func(p Person) any { return p.Name },  // key function
    func(p Person) any { return p.Age },   // value function
)
```

### Available Methods

#### Core Operations
- `FromSlice([]T)` - Create stream from slice
- `Iterator()` - Get iterator for manual traversal
- `ToSlice()` - Materialize stream to slice
- `ToSet()` - Materialize stream to set (map[any]struct{})
- `ToMap(keyFn, valueFn)` - Materialize stream to map

#### Querying
- `IsEmpty()` - Check if stream has no elements
- `Count()` - Get number of elements
- `Contains(element)` - Check if element exists
- `IndexOf(element, fromIndex...)` - Find index of element
- `Find(predicate)` - Find first matching element
- `FindIndex(predicate)` - Find index of first matching element

#### Transformation
- `Map(fn)` - Transform each element
- `Filter(predicate)` - Keep elements matching predicate
- `NonNullable()` - Remove zero-value elements
- `Distinct(keyFn...)` - Remove duplicates
- `Exclude(other, keyFn...)` - Remove elements present in other collection

#### Aggregation
- `Reduce(fn)` - Reduce to single value
- `ReduceWithInitial(fn, initial)` - Reduce with initial value
- `ReduceRight(fn)` - Reduce from right to left
- `Every(predicate)` - Check if all elements match
- `Any(predicate)` - Check if any element matches
- `ForEach(action)` - Execute action for each element

#### Manipulation
- `Head()` - Get first element
- `Tail(n)` - Skip first n elements
- `Limit(maxSize)` - Limit to maximum number of elements
- `Concat(other)` - Concatenate with another iterable
- `Join(separator)` - Join elements into string

## Development

### Prerequisites

- Go 1.21 or later
- Git

### Building

```bash
go build ./...
```

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with race detection
go test -race ./...
```

### Linting

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
golangci-lint run
```

## CI/CD

This project uses GitHub Actions for continuous integration:

- **Multi-platform testing**: Linux, Windows, macOS
- **Multi-version testing**: Go 1.21, 1.22, 1.23
- **Code quality**: Linting with golangci-lint
- **Security scanning**: gosec security analysis
- **Coverage reporting**: Codecov integration
- **Dependency verification**: go mod verify and tidy checks

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Inspired by [Langium](https://langium.org/) TypeScript utilities
- Stream API design influenced by Java Streams and JavaScript Array methods