// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package stateless

import (
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestFirst tests getting the first element from a non-empty iterator
func TestFirst(t *testing.T) {
	t.Run("returns first element from integer iterator", func(t *testing.T) {
		seq := From(1, 2, 3)
		fst := First(seq)
		assert.Equal(t, O.Of(1), fst)
	})

	t.Run("returns first element from string iterator", func(t *testing.T) {
		seq := From("a", "b", "c")
		fst := First(seq)
		assert.Equal(t, O.Of("a"), fst)
	})

	t.Run("returns first element from single element iterator", func(t *testing.T) {
		seq := From(42)
		fst := First(seq)
		assert.Equal(t, O.Of(42), fst)
	})

	t.Run("returns first element from large iterator", func(t *testing.T) {
		seq := From(100, 200, 300, 400, 500)
		fst := First(seq)
		assert.Equal(t, O.Of(100), fst)
	})
}

// TestNoFirst tests getting the first element from an empty iterator
func TestNoFirst(t *testing.T) {
	t.Run("returns None for empty integer iterator", func(t *testing.T) {
		seq := Empty[int]()
		fst := First(seq)
		assert.Equal(t, O.None[int](), fst)
	})

	t.Run("returns None for empty string iterator", func(t *testing.T) {
		seq := Empty[string]()
		fst := First(seq)
		assert.Equal(t, O.None[string](), fst)
	})

	t.Run("returns None for empty struct iterator", func(t *testing.T) {
		type TestStruct struct {
			Value int
		}
		seq := Empty[TestStruct]()
		fst := First(seq)
		assert.Equal(t, O.None[TestStruct](), fst)
	})
}

// TestFirstWithFiltered tests First with filtered iterators
func TestFirstWithFiltered(t *testing.T) {
	t.Run("returns first element matching filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		filtered := Filter(N.MoreThan(3))(seq)
		fst := First(filtered)
		assert.Equal(t, O.Of(4), fst)
	})

	t.Run("returns None when no elements match filter", func(t *testing.T) {
		seq := From(1, 2, 3)
		filtered := Filter(N.MoreThan(10))(seq)
		fst := First(filtered)
		assert.Equal(t, O.None[int](), fst)
	})

	t.Run("returns first even number", func(t *testing.T) {
		seq := From(1, 3, 5, 6, 7, 8)
		filtered := Filter(func(x int) bool { return x%2 == 0 })(seq)
		fst := First(filtered)
		assert.Equal(t, O.Of(6), fst)
	})
}

// TestFirstWithMapped tests First with mapped iterators
func TestFirstWithMapped(t *testing.T) {
	t.Run("returns first element after mapping", func(t *testing.T) {
		seq := From(1, 2, 3)
		mapped := Map(N.Mul(2))(seq)
		fst := First(mapped)
		assert.Equal(t, O.Of(2), fst)
	})

	t.Run("returns first string after mapping", func(t *testing.T) {
		seq := From(1, 2, 3)
		mapped := Map(S.Format[int]("num-%d"))(seq)
		fst := First(mapped)
		assert.Equal(t, O.Of("num-1"), fst)
	})
}

// TestFirstWithTake tests First with Take
func TestFirstWithTake(t *testing.T) {
	t.Run("returns first element from taken subset", func(t *testing.T) {
		seq := From(10, 20, 30, 40, 50)
		taken := Take[int](3)(seq)
		fst := First(taken)
		assert.Equal(t, O.Of(10), fst)
	})

	t.Run("returns None when taking zero elements", func(t *testing.T) {
		seq := From(1, 2, 3)
		taken := Take[int](0)(seq)
		fst := First(taken)
		assert.Equal(t, O.None[int](), fst)
	})
}

// TestFirstWithComplex tests First with complex types
func TestFirstWithComplex(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("returns first person", func(t *testing.T) {
		seq := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
			Person{"Charlie", 35},
		)
		fst := First(seq)
		expected := O.Of(Person{"Alice", 30})
		assert.Equal(t, expected, fst)
	})

	t.Run("returns first pointer", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		seq := From(p1, p2)
		fst := First(seq)
		assert.Equal(t, O.Of(p1), fst)
	})
}

// TestFirstDoesNotConsumeEntireIterator tests that First only consumes the first element
func TestFirstDoesNotConsumeEntireIterator(t *testing.T) {
	t.Run("only consumes first element", func(t *testing.T) {
		callCount := 0
		seq := Map(func(x int) int {
			callCount++
			return x * 2
		})(From(1, 2, 3, 4, 5))

		fst := First(seq)

		assert.Equal(t, O.Of(2), fst)
		// Should only have called the map function once for the first element
		assert.Equal(t, 1, callCount)
	})
}

// TestFirstWithChainedOperations tests First with multiple chained operations
func TestFirstWithChainedOperations(t *testing.T) {
	t.Run("chains filter, map, and first", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := First(
			Map(N.Mul(10))(
				Filter(N.MoreThan(5))(seq),
			),
		)
		assert.Equal(t, O.Of(60), result)
	})

	t.Run("chains take, filter, and first", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := First(
			Filter(N.MoreThan(3))(
				Take[int](7)(seq),
			),
		)
		assert.Equal(t, O.Of(4), result)
	})
}

// Benchmark tests
func BenchmarkFirst(b *testing.B) {
	seq := From(1, 2, 3, 4, 5)
	b.ResetTimer()
	for b.Loop() {
		First(seq)
	}
}

func BenchmarkFirstLargeIterator(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	seq := FromArray(data)

	for b.Loop() {
		First(seq)
	}
}

// Example tests for documentation
func ExampleFirst() {
	iter := From(1, 2, 3, 4, 5)
	first := First(iter)

	if value, ok := O.Unwrap(first); ok {
		fmt.Printf("First element: %d\n", value)
	}
	// Output: First element: 1
}

func ExampleFirst_empty() {
	iter := Empty[int]()
	first := First(iter)

	if _, ok := O.Unwrap(first); !ok {
		fmt.Println("Iterator is empty")
	}
	// Output: Iterator is empty
}

func ExampleFirst_filtered() {
	iter := From(1, 2, 3, 4, 5)
	filtered := Filter(N.MoreThan(3))(iter)
	first := First(filtered)

	if value, ok := O.Unwrap(first); ok {
		fmt.Printf("First element > 3: %d\n", value)
	}
	// Output: First element > 3: 4
}
