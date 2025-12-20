// Copyright (c) 2025 IBM Corp.
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

package iter

import (
	"fmt"
	"testing"

	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// TestCompress tests the Compress function
func TestCompress(t *testing.T) {
	t.Run("filters with alternating selectors", func(t *testing.T) {
		data := From(1, 2, 3, 4, 5)
		selectors := From(true, false, true, false, true)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []int{1, 3, 5}, result)
	})

	t.Run("filters strings with boolean selectors", func(t *testing.T) {
		data := From("a", "b", "c", "d", "e")
		selectors := From(true, true, false, false, true)
		filtered := Compress[string](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []string{"a", "b", "e"}, result)
	})

	t.Run("all true selectors pass all elements", func(t *testing.T) {
		data := From(10, 20, 30)
		selectors := From(true, true, true)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []int{10, 20, 30}, result)
	})

	t.Run("all false selectors produce empty sequence", func(t *testing.T) {
		data := From(1, 2, 3)
		selectors := From(false, false, false)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Empty(t, result)
	})

	t.Run("shorter selector sequence stops early", func(t *testing.T) {
		data := From("a", "b", "c", "d", "e")
		selectors := From(true, true, false)
		filtered := Compress[string](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []string{"a", "b"}, result)
	})

	t.Run("shorter data sequence stops early", func(t *testing.T) {
		data := From(1, 2, 3)
		selectors := From(true, false, true, true, true)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []int{1, 3}, result)
	})

	t.Run("empty data sequence produces empty result", func(t *testing.T) {
		data := Empty[int]()
		selectors := From(true, true, true)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Empty(t, result)
	})

	t.Run("empty selector sequence produces empty result", func(t *testing.T) {
		data := From(1, 2, 3)
		selectors := Empty[bool]()
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Empty(t, result)
	})

	t.Run("both empty sequences produce empty result", func(t *testing.T) {
		data := Empty[int]()
		selectors := Empty[bool]()
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Empty(t, result)
	})

	t.Run("single element with true selector", func(t *testing.T) {
		data := From(42)
		selectors := From(true)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []int{42}, result)
	})

	t.Run("single element with false selector", func(t *testing.T) {
		data := From(42)
		selectors := From(false)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Empty(t, result)
	})
}

// TestCompressWithComplexTypes tests Compress with complex data types
func TestCompressWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("filters struct values", func(t *testing.T) {
		data := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
			Person{"Charlie", 35},
			Person{"David", 28},
		)
		selectors := From(true, false, true, false)
		filtered := Compress[Person](selectors)(data)
		result := toSlice(filtered)
		expected := []Person{
			{"Alice", 30},
			{"Charlie", 35},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("filters pointer values", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		p3 := &Person{"Charlie", 35}
		data := From(p1, p2, p3)
		selectors := From(false, true, true)
		filtered := Compress[*Person](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []*Person{p2, p3}, result)
	})
}

// TestCompressWithChainedOperations tests Compress with other operations
func TestCompressWithChainedOperations(t *testing.T) {
	t.Run("compress then map", func(t *testing.T) {
		data := From(1, 2, 3, 4, 5)
		selectors := From(true, false, true, false, true)
		result := toSlice(
			MonadMap(
				Compress[int](selectors)(data),
				func(x int) int { return x * 10 },
			),
		)
		assert.Equal(t, []int{10, 30, 50}, result)
	})

	t.Run("map then compress", func(t *testing.T) {
		data := From(1, 2, 3, 4, 5)
		mapped := MonadMap(data, func(x int) int { return x * 2 })
		selectors := From(true, true, false, false, true)
		result := toSlice(Compress[int](selectors)(mapped))
		assert.Equal(t, []int{2, 4, 10}, result)
	})

	t.Run("compress with filtered data", func(t *testing.T) {
		data := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		evens := MonadFilter(data, func(x int) bool { return x%2 == 0 })
		selectors := From(true, false, true, false, true)
		result := toSlice(Compress[int](selectors)(evens))
		assert.Equal(t, []int{2, 6, 10}, result)
	})
}

// TestToSeqPair tests the ToSeqPair function
func TestToSeqPair(t *testing.T) {
	t.Run("converts Seq2 to sequence of pairs", func(t *testing.T) {
		seq2 := MonadZip(From("a", "b", "c"), From(1, 2, 3))
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)
		expected := []Pair[string, int]{
			P.MakePair("a", 1),
			P.MakePair("b", 2),
			P.MakePair("c", 3),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("converts empty Seq2", func(t *testing.T) {
		seq2 := MonadZip(Empty[int](), Empty[string]())
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)
		assert.Empty(t, result)
	})

	t.Run("converts single pair", func(t *testing.T) {
		seq2 := MonadZip(From(42), From("answer"))
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)
		expected := []Pair[int, string]{
			P.MakePair(42, "answer"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("stops at shorter sequence", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3, 4, 5), From("a", "b"))
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)
		expected := []Pair[int, string]{
			P.MakePair(1, "a"),
			P.MakePair(2, "b"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestToSeqPairWithOperations tests ToSeqPair with other operations
func TestToSeqPairWithOperations(t *testing.T) {
	t.Run("map over pairs", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3), From(10, 20, 30))
		pairs := ToSeqPair(seq2)
		sums := MonadMap(pairs, func(p Pair[int, int]) int {
			return P.Head(p) + P.Tail(p)
		})
		result := toSlice(sums)
		assert.Equal(t, []int{11, 22, 33}, result)
	})

	t.Run("filter pairs", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3, 4, 5), From(10, 20, 30, 40, 50))
		pairs := ToSeqPair(seq2)
		filtered := MonadFilter(pairs, func(p Pair[int, int]) bool {
			return P.Head(p)%2 == 0
		})
		result := toSlice(filtered)
		expected := []Pair[int, int]{
			P.MakePair(2, 20),
			P.MakePair(4, 40),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("extract first elements from pairs", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3), From("x", "y", "z"))
		pairs := ToSeqPair(seq2)
		firsts := MonadMap(pairs, func(p Pair[int, string]) int {
			return P.Head(p)
		})
		result := toSlice(firsts)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("extract second elements from pairs", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		pairs := ToSeqPair(seq2)
		seconds := MonadMap(pairs, func(p Pair[int, string]) string {
			return P.Tail(p)
		})
		result := toSlice(seconds)
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})
}

// TestCompressAndToSeqPairTogether tests using both functions together
func TestCompressAndToSeqPairTogether(t *testing.T) {
	t.Run("compress uses ToSeqPair internally", func(t *testing.T) {
		// This test verifies the integration works correctly
		data := From(10, 20, 30, 40, 50)
		selectors := From(true, false, true, true, false)
		filtered := Compress[int](selectors)(data)
		result := toSlice(filtered)
		assert.Equal(t, []int{10, 30, 40}, result)
	})
}

// Benchmark tests
func BenchmarkCompress(b *testing.B) {
	data := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	selectors := From(true, false, true, false, true, false, true, false, true, false)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := Compress[int](selectors)(data)
		for range filtered {
		}
	}
}

func BenchmarkToSeqPair(b *testing.B) {
	seq2 := MonadZip(From(1, 2, 3, 4, 5), From("a", "b", "c", "d", "e"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pairs := ToSeqPair(seq2)
		for range pairs {
		}
	}
}

// Example tests for documentation
func ExampleCompress() {
	data := From(1, 2, 3, 4, 5)
	selectors := From(true, false, true, false, true)
	filtered := Compress[int](selectors)(data)

	for v := range filtered {
		fmt.Printf("%d ", v)
	}
	// Output: 1 3 5
}

func ExampleCompress_allTrue() {
	data := From(10, 20, 30)
	selectors := From(true, true, true)
	filtered := Compress[int](selectors)(data)

	for v := range filtered {
		fmt.Printf("%d ", v)
	}
	// Output: 10 20 30
}

func ExampleCompress_allFalse() {
	data := From(1, 2, 3)
	selectors := From(false, false, false)
	filtered := Compress[int](selectors)(data)

	count := 0
	for range filtered {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}

func ExampleToSeqPair() {
	seq2 := MonadZip(From(1, 2, 3), From("a", "b", "c"))
	pairs := ToSeqPair(seq2)

	for p := range pairs {
		fmt.Printf("(%d, %s) ", P.Head(p), P.Tail(p))
	}
	// Output: (1, a) (2, b) (3, c)
}

func ExampleToSeqPair_withMap() {
	seq2 := MonadZip(From(1, 2, 3), From(10, 20, 30))
	pairs := ToSeqPair(seq2)
	sums := MonadMap(pairs, func(p Pair[int, int]) int {
		return P.Head(p) + P.Tail(p)
	})

	for sum := range sums {
		fmt.Printf("%d ", sum)
	}
	// Output: 11 22 33
}
