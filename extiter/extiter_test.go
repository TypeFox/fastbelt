// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package extiter

import (
	"iter"
	"reflect"
	"slices"
	"testing"
)

func TestCount(t *testing.T) {
	tests := []struct {
		name     string
		slice    []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{1}, 1},
		{"multiple elements", []int{1, 2, 3}, 3},
		{"many elements", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			seq := slices.Values(tt.slice)
			count := Count(seq)
			if count != tt.expected {
				t.Errorf("expected count %d, got %d", tt.expected, count)
			}
		})
	}
}

func TestIsEmpty(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		if !IsEmpty(seq) {
			t.Error("expected sequence to be empty")
		}
	})

	t.Run("non-empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		if IsEmpty(seq) {
			t.Error("expected sequence to not be empty")
		}
	})
}

func TestToSet(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		result := ToSet(seq)
		expected := make(map[any]struct{})
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("non-empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := ToSet(seq)
		expected := map[any]struct{}{
			1: {},
			2: {},
			3: {},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("duplicate elements", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 2, 3, 1})
		result := ToSet(seq)
		expected := map[any]struct{}{
			1: {},
			2: {},
			3: {},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestToMap(t *testing.T) {
	type testStruct struct {
		A int
		B string
	}

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]testStruct{})
		result := ToMap(seq, nil, nil)
		expected := make(map[any]any)
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("key and value unmapped", func(t *testing.T) {
		data := []testStruct{{1, "foo"}, {2, "bar"}}
		seq := slices.Values(data)
		result := ToMap(seq, nil, nil)
		expected := map[any]any{
			testStruct{1, "foo"}: testStruct{1, "foo"},
			testStruct{2, "bar"}: testStruct{2, "bar"},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("key mapped", func(t *testing.T) {
		data := []testStruct{{1, "foo"}, {2, "bar"}}
		seq := slices.Values(data)
		result := ToMap(seq, func(e testStruct) any { return e.B }, nil)
		expected := map[any]any{
			"foo": testStruct{1, "foo"},
			"bar": testStruct{2, "bar"},
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("value mapped", func(t *testing.T) {
		data := []testStruct{{1, "foo"}, {2, "bar"}}
		seq := slices.Values(data)
		result := ToMap(seq, nil, func(e testStruct) any { return e.A })
		expected := map[any]any{
			testStruct{1, "foo"}: 1,
			testStruct{2, "bar"}: 2,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("key and value mapped", func(t *testing.T) {
		data := []testStruct{{1, "foo"}, {2, "bar"}}
		seq := slices.Values(data)
		result := ToMap(seq, func(e testStruct) any { return e.B }, func(e testStruct) any { return e.A })
		expected := map[any]any{
			"foo": 1,
			"bar": 2,
		}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestJoin(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]string{})
		result := Join(seq, ",")
		if result != "" {
			t.Errorf("expected empty string, got %q", result)
		}
	})

	t.Run("string sequence with default separator", func(t *testing.T) {
		seq := slices.Values([]string{"a", "b"})
		result := Join(seq, "")
		expected := "a,b"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("string sequence with custom separator", func(t *testing.T) {
		seq := slices.Values([]string{"a", "b"})
		result := Join(seq, " & ")
		expected := "a & b"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})

	t.Run("number sequence", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := Join(seq, ",")
		expected := "1,2,3"
		if result != expected {
			t.Errorf("expected %q, got %q", expected, result)
		}
	})
}

func TestIndexOf(t *testing.T) {
	t.Run("number sequence present", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		index := IndexOf(seq, 2)
		if index != 1 {
			t.Errorf("expected index 1, got %d", index)
		}
	})

	t.Run("number sequence absent", func(t *testing.T) {
		seq := slices.Values([]int{1, 3})
		index := IndexOf(seq, 2)
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})

	t.Run("with fromIndex", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 2, 4})
		index := IndexOf(seq, 2, 2)
		if index != 3 {
			t.Errorf("expected index 3, got %d", index)
		}
	})
}

func TestEvery(t *testing.T) {
	t.Run("all true", func(t *testing.T) {
		seq := slices.Values([]int{2, 4, 6})
		result := Every(seq, func(value int) bool { return value%2 == 0 })
		if !result {
			t.Error("expected true")
		}
	})

	t.Run("not all true", func(t *testing.T) {
		seq := slices.Values([]int{2, 3, 6})
		result := Every(seq, func(value int) bool { return value%2 == 0 })
		if result {
			t.Error("expected false")
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		result := Every(seq, func(value int) bool { return value%2 == 0 })
		if !result {
			t.Error("expected true for empty sequence")
		}
	})
}

func TestAny(t *testing.T) {
	t.Run("some true", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := Any(seq, func(value int) bool { return value%2 == 0 })
		if !result {
			t.Error("expected true")
		}
	})

	t.Run("none true", func(t *testing.T) {
		seq := slices.Values([]int{1, 3, 5})
		result := Any(seq, func(value int) bool { return value%2 == 0 })
		if result {
			t.Error("expected false")
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		result := Any(seq, func(value int) bool { return value%2 == 0 })
		if result {
			t.Error("expected false for empty sequence")
		}
	})
}

func TestForEach(t *testing.T) {
	t.Run("sum values and indices", func(t *testing.T) {
		seq := slices.Values([]int{2, 4, 6})
		sumValue := 0
		sumIndex := 0
		ForEach(seq, func(value int, index int) {
			sumValue += value
			sumIndex += index
		})
		if sumValue != 12 {
			t.Errorf("expected sumValue 12, got %d", sumValue)
		}
		if sumIndex != 3 {
			t.Errorf("expected sumIndex 3, got %d", sumIndex)
		}
	})
}

func TestMap(t *testing.T) {
	t.Run("increment numbers", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		mapped := Map(seq, func(value int) int { return value + 1 })
		result := slices.Collect(mapped)
		expected := []int{2, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("convert to string", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		mapped := Map(seq, func(value int) string { return toString(value) })
		result := slices.Collect(mapped)
		expected := []string{"1", "2", "3"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestFilter(t *testing.T) {
	t.Run("filter even numbers", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5, 6})
		filtered := Filter(seq, func(value int) bool { return value%2 == 0 })
		result := slices.Collect(filtered)
		expected := []int{2, 4, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("filter greater than value", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		filtered := Filter(seq, func(value int) bool { return value >= 3 })
		result := slices.Collect(filtered)
		expected := []int{3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestNonNullable(t *testing.T) {
	t.Run("filter zero values", func(t *testing.T) {
		seq := slices.Values([]int{0, 1, 0, 2, 3})
		filtered := NonNullable(seq)
		result := slices.Collect(filtered)
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestReduce(t *testing.T) {
	t.Run("sum numbers", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4})
		result, ok := Reduce(seq, func(a, b int) int { return a + b })
		if !ok {
			t.Error("expected reduce to succeed")
		}
		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		_, ok := Reduce(seq, func(a, b int) int { return a + b })
		if ok {
			t.Error("expected reduce to fail for empty sequence")
		}
	})

	t.Run("single element", func(t *testing.T) {
		seq := slices.Values([]int{42})
		result, ok := Reduce(seq, func(a, b int) int { return a + b })
		if !ok {
			t.Error("expected reduce to succeed")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})
}

func TestReduceWithInitial(t *testing.T) {
	t.Run("sum with initial value", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := ReduceWithInitial(seq, func(acc int, val int) int {
			return acc + val
		}, 10)
		if result != 16 {
			t.Errorf("expected 16, got %v", result)
		}
	})

	t.Run("concatenate to slice", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := ReduceWithInitial(seq, func(acc []int, val int) []int {
			return append(acc, val)
		}, []int{})
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestReduceRight(t *testing.T) {
	t.Run("sum numbers right to left", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4})
		result, ok := ReduceRight(seq, func(a, b int) int { return a + b })
		if !ok {
			t.Error("expected reduce to succeed")
		}
		if result != 10 {
			t.Errorf("expected 10, got %d", result)
		}
	})

	t.Run("subtract numbers right to left", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result, ok := ReduceRight(seq, func(a, b int) int { return a - b })
		if !ok {
			t.Error("expected reduce to succeed")
		}
		// ((3 - 2) - 1) = 0
		if result != 0 {
			t.Errorf("expected 0, got %d", result)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		_, ok := ReduceRight(seq, func(a, b int) int { return a + b })
		if ok {
			t.Error("expected reduce to fail for empty sequence")
		}
	})

	t.Run("single element", func(t *testing.T) {
		seq := slices.Values([]int{42})
		result, ok := ReduceRight(seq, func(a, b int) int { return a + b })
		if !ok {
			t.Error("expected reduce to succeed")
		}
		if result != 42 {
			t.Errorf("expected 42, got %d", result)
		}
	})
}

func TestReduceRightWithInitial(t *testing.T) {
	t.Run("concatenate strings right to left", func(t *testing.T) {
		seq := slices.Values([]string{"a", "b", "c"})
		result := ReduceRightWithInitial(seq, func(acc string, val string) string {
			return acc + val
		}, "")
		if result != "cba" {
			t.Errorf("expected 'cba', got %v", result)
		}
	})

	t.Run("build slice in reverse", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		result := ReduceRightWithInitial(seq, func(acc []int, val int) []int {
			return append(acc, val)
		}, []int{})
		expected := []int{3, 2, 1}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestFind(t *testing.T) {
	t.Run("number found", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		result, found := Find(seq, func(value int) bool { return value > 3 })
		if !found {
			t.Error("expected to find value")
		}
		if result != 4 {
			t.Errorf("expected 4, got %d", result)
		}
	})

	t.Run("number not found", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		_, found := Find(seq, func(value int) bool { return value > 5 })
		if found {
			t.Error("expected not to find value")
		}
	})
}

func TestFindIndex(t *testing.T) {
	t.Run("found at index", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		index := FindIndex(seq, func(value int) bool { return value > 3 })
		if index != 3 {
			t.Errorf("expected index 3, got %d", index)
		}
	})

	t.Run("not found", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		index := FindIndex(seq, func(value int) bool { return value > 5 })
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("element present", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		if !Contains(seq, 2) {
			t.Error("expected sequence to contain 2")
		}
	})

	t.Run("element absent", func(t *testing.T) {
		seq := slices.Values([]int{1, 3})
		if Contains(seq, 2) {
			t.Error("expected sequence not to contain 2")
		}
	})
}

func TestHead(t *testing.T) {
	t.Run("non-empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		head, ok := Head(seq)
		if !ok {
			t.Error("expected to get head")
		}
		if head != 1 {
			t.Errorf("expected 1, got %d", head)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		_, ok := Head(seq)
		if ok {
			t.Error("expected not to get head from empty sequence")
		}
	})
}

func TestTail(t *testing.T) {
	t.Run("skip one element", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		tail := Tail(seq, 1)
		result := slices.Collect(tail)
		expected := []int{2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("skip three elements", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		tail := Tail(seq, 3)
		result := slices.Collect(tail)
		expected := []int{4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})
		tail := Tail(seq, 1)
		result := slices.Collect(tail)
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})
}

func TestLimit(t *testing.T) {
	t.Run("limit to 3 elements", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5})
		limited := Limit(seq, 3)
		result := slices.Collect(limited)
		expected := []int{1, 2, 3}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("limit to 0 elements", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		limited := Limit(seq, 0)
		result := slices.Collect(limited)
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})
}

func TestDistinct(t *testing.T) {
	t.Run("distinct numbers", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 2, 3, 1, 4})
		distinct := Distinct(seq, nil)
		result := slices.Collect(distinct)
		expected := []int{1, 2, 3, 4}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("distinct with custom key function", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		people := []Person{
			{"Alice", 25},
			{"Bob", 30},
			{"Alice", 26}, // Different age, same name
		}

		seq := slices.Values(people)
		distinct := Distinct(seq, func(p Person) any { return p.Name })
		result := slices.Collect(distinct)
		expected := []Person{{"Alice", 25}, {"Bob", 30}}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestExclude(t *testing.T) {
	t.Run("exclude overlapping string values", func(t *testing.T) {
		seq1 := slices.Values([]string{"a", "b", "c"})
		seq2 := slices.Values([]string{"b", "d"})
		result := Exclude(seq1, seq2, func(s string) any { return s })
		actual := slices.Collect(result)
		expected := []string{"a", "c"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})

	t.Run("no overlapping values", func(t *testing.T) {
		seq1 := slices.Values([]string{"a", "b"})
		seq2 := slices.Values([]string{"c", "d"})
		result := Exclude(seq1, seq2, func(s string) any { return s })
		actual := slices.Collect(result)
		expected := []string{"a", "b"}
		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("expected %v, got %v", expected, actual)
		}
	})
}

func TestConcat(t *testing.T) {
	t.Run("concatenate two sequences", func(t *testing.T) {
		seq1 := slices.Values([]string{"a", "b"})
		seq2 := slices.Values([]string{"c", "d"})
		concatenated := Concat(seq1, seq2)
		result := slices.Collect(concatenated)
		expected := []string{"a", "b", "c", "d"}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("concatenate three sequences", func(t *testing.T) {
		seq1 := slices.Values([]int{1, 2})
		seq2 := slices.Values([]int{3, 4})
		seq3 := slices.Values([]int{5, 6})
		concatenated := Concat(seq1, seq2, seq3)
		result := slices.Collect(concatenated)
		expected := []int{1, 2, 3, 4, 5, 6}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

func TestFlatMap(t *testing.T) {
	t.Run("flatten nested sequences", func(t *testing.T) {
		seq := slices.Values([][]int{{1, 2}, {3, 4}, {5}})
		flattened := FlatMap(seq, func(slice []int) iter.Seq[int] {
			return slices.Values(slice)
		})
		result := slices.Collect(flattened)
		expected := []int{1, 2, 3, 4, 5}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// Test chaining operations
func TestChaining(t *testing.T) {
	t.Run("filter then map", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5, 6})
		result := slices.Collect(
			Map(
				Filter(seq, func(x int) bool { return x%2 == 0 }),
				func(x int) int { return x * 2 },
			),
		)
		expected := []int{4, 8, 12}
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})

	t.Run("complex chaining", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})

		// Filter even numbers, map to their squares, limit to 3
		result := slices.Collect(
			Limit(
				Map(
					Filter(seq, func(x int) bool { return x%2 == 0 }),
					func(x int) int { return x * x },
				),
				3,
			),
		)
		expected := []int{4, 16, 36} // 2^2, 4^2, 6^2
		if !reflect.DeepEqual(result, expected) {
			t.Errorf("expected %v, got %v", expected, result)
		}
	})
}

// Test edge cases
func TestEdgeCases(t *testing.T) {
	t.Run("operations on empty sequence", func(t *testing.T) {
		seq := slices.Values([]int{})

		// All these should work without panicking
		if !IsEmpty(seq) {
			t.Error("expected empty sequence")
		}
		if Count(seq) != 0 {
			t.Error("expected count 0")
		}
		if len(slices.Collect(seq)) != 0 {
			t.Error("expected empty slice")
		}
		if Join(seq, ",") != "" {
			t.Error("expected empty string")
		}
		if IndexOf(seq, 1) != -1 {
			t.Error("expected -1 for indexOf on empty sequence")
		}
		if !Every(seq, func(int) bool { return false }) {
			t.Error("expected Every to return true for empty sequence")
		}
		if Any(seq, func(int) bool { return true }) {
			t.Error("expected Any to return false for empty sequence")
		}
	})

	t.Run("tail with large skip count", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		tail := Tail(seq, 10) // Skip more than available
		result := slices.Collect(tail)
		if len(result) != 0 {
			t.Errorf("expected empty slice, got %v", result)
		}
	})

	t.Run("indexOf with fromIndex beyond sequence", func(t *testing.T) {
		seq := slices.Values([]int{1, 2, 3})
		index := IndexOf(seq, 1, 10)
		if index != -1 {
			t.Errorf("expected -1, got %d", index)
		}
	})
}