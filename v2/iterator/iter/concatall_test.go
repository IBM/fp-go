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

package iter

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConcatBuf_Success tests basic ConcatBuf functionality
func TestConcatBuf_Success(t *testing.T) {
	t.Run("concatenates two sequences in order", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("concatenates three sequences in order", func(t *testing.T) {
		seq1 := From(1, 2)
		seq2 := From(3, 4)
		seq3 := From(5, 6)
		result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2, seq3}, defaultBufferSize))

		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("concatenates many sequences in order", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2)
		}
		result := toSlice(ConcatBuf(sequences, defaultBufferSize))

		assert.Equal(t, 30, len(result))
		expected := make([]int, 30)
		for i := range 30 {
			expected[i] = (i/3)*10 + (i % 3)
		}
		assert.Equal(t, expected, result)
	})

	t.Run("concatenates single sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(ConcatBuf([]Seq[int]{seq}, defaultBufferSize))

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("preserves order with strings", func(t *testing.T) {
		seq1 := From("a", "b", "c")
		seq2 := From("d", "e", "f")
		result := toSlice(ConcatBuf([]Seq[string]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, result)
	})
}

// TestConcatBuf_Empty tests ConcatBuf with empty inputs
func TestConcatBuf_Empty(t *testing.T) {
	t.Run("empty slice of sequences", func(t *testing.T) {
		result := toSlice(Concat([]Seq[int]{}))
		assert.Empty(t, result)
	})

	t.Run("nil slice of sequences", func(t *testing.T) {
		result := toSlice(Concat([]Seq[int](nil)))
		assert.Empty(t, result)
	})

	t.Run("slice with empty sequences", func(t *testing.T) {
		seq1 := Empty[int]()
		seq2 := Empty[int]()
		result := toSlice(Concat([]Seq[int]{seq1, seq2}))
		assert.Empty(t, result)
	})

	t.Run("mix of empty and non-empty sequences preserves order", func(t *testing.T) {
		seq1 := Empty[int]()
		seq2 := From(1, 2, 3)
		seq3 := Empty[int]()
		seq4 := From(4, 5)
		result := toSlice(Concat([]Seq[int]{seq1, seq2, seq3, seq4}))

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// TestConcatBuf_OrderPreservation verifies deterministic ordering
func TestConcatBuf_OrderPreservation(t *testing.T) {
	t.Run("order is always deterministic", func(t *testing.T) {
		// Run multiple times to confirm determinism
		expected := []int{1, 2, 3, 4, 5, 6}
		for range 10 {
			seq1 := From(1, 2, 3)
			seq2 := From(4, 5, 6)
			result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize))
			assert.Equal(t, expected, result)
		}
	})

	t.Run("seq1 elements always precede seq2 elements", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize))

		// All elements from seq1 must come before any element from seq2
		maxSeq1 := 0
		for _, v := range result[:3] {
			if v > maxSeq1 {
				maxSeq1 = v
			}
		}
		minSeq2 := result[3]
		assert.Less(t, maxSeq1, minSeq2, "all seq1 elements must precede seq2 elements")
	})
}

// TestConcatBuf_EarlyTermination tests that ConcatBuf handles early termination
func TestConcatBuf_EarlyTermination(t *testing.T) {
	t.Run("stops at correct element when consumer breaks", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)

		result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)

		var consumed []int
		for v := range result {
			consumed = append(consumed, v)
			if len(consumed) >= 3 {
				break
			}
		}

		// Consumer receives the first 3 elements in order.
		// Producer goroutines may run ahead due to buffering, but
		// the consumer always observes a deterministic prefix.
		assert.Equal(t, []int{1, 2, 3}, consumed)
	})

	t.Run("handles yield returning false mid-sequence", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)

		collected := []int{}
		for v := range result {
			collected = append(collected, v)
			if len(collected) >= 4 {
				break
			}
		}

		assert.Equal(t, []int{1, 2, 3, 4}, collected)
	})

	t.Run("consumer receives only seq1 elements when breaking after seq1", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)

		var consumed []int
		for v := range result {
			consumed = append(consumed, v)
			if len(consumed) >= 3 {
				break
			}
		}

		// Consumer receives exactly seq1's elements; seq2 elements are never yielded.
		// Producer goroutines may start concurrently, but the drainer forwards
		// seq2 elements only after seq1 is fully drained.
		assert.Equal(t, []int{1, 2, 3}, consumed)
	})
}

// TestConcatBuf_WithComplexTypes tests ConcatBuf with complex data types
func TestConcatBuf_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with structs, preserving order", func(t *testing.T) {
		seq1 := From(Person{"Alice", 30}, Person{"Bob", 25})
		seq2 := From(Person{"Charlie", 35}, Person{"Diana", 28})
		result := toSlice(ConcatBuf([]Seq[Person]{seq1, seq2}, defaultBufferSize))

		expected := []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
			{"Diana", 28},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with pointers", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		p3 := &Person{"Charlie", 35}
		seq1 := From(p1, p2)
		seq2 := From(p3)
		result := toSlice(ConcatBuf([]Seq[*Person]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, []*Person{p1, p2, p3}, result)
	})

	t.Run("works with slices", func(t *testing.T) {
		seq1 := From([]int{1, 2}, []int{3, 4})
		seq2 := From([]int{5, 6}, []int{7, 8})
		result := toSlice(ConcatBuf([]Seq[[]int]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}, result)
	})
}

// TestConcatBuf_WithChainedOperations tests ConcatBuf with other sequence operations
func TestConcatBuf_WithChainedOperations(t *testing.T) {
	t.Run("concat after map", func(t *testing.T) {
		seq1 := MonadMap(From(1, 2, 3), N.Mul(2))
		seq2 := MonadMap(From(4, 5, 6), N.Mul(2))
		result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, result)
	})

	t.Run("map after concat preserves order", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		concat := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)
		result := toSlice(MonadMap(concat, N.Mul(2)))

		assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, result)
	})

	t.Run("concat after filter", func(t *testing.T) {
		seq1 := MonadFilter(From(1, 2, 3, 4, 5), func(x int) bool { return x%2 == 0 })
		seq2 := MonadFilter(From(6, 7, 8, 9, 10), func(x int) bool { return x%2 == 0 })
		result := toSlice(ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize))

		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("filter after concat preserves order", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		concat := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)
		result := toSlice(MonadFilter(concat, func(x int) bool { return x%2 == 0 }))

		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})
}

// TestConcat_Alias tests that Concat is an alias for ConcatBuf
func TestConcat_Alias(t *testing.T) {
	t.Run("concat works identically to ConcatBuf", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		result := toSlice(Concat([]Seq[int]{seq1, seq2}))

		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("concat handles empty input", func(t *testing.T) {
		result := toSlice(Concat([]Seq[int]{}))
		assert.Empty(t, result)
	})

	t.Run("concat preserves order deterministically", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		for range 5 {
			result := toSlice(Concat([]Seq[int]{seq1, seq2}))
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
		}
	})
}

// TestConcatMapBuf_Success tests basic ConcatMapBuf functionality
func TestConcatMapBuf_Success(t *testing.T) {
	t.Run("expands each element in order", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return From(n, n*10)
		}, defaultBufferSize)
		input := From(1, 2, 3)

		result := toSlice(expand(input))

		// Order must be deterministic: all of f(1), then all of f(2), then all of f(3)
		assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, result)
	})

	t.Run("handles single element", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return From(n, n*2, n*3)
		}, defaultBufferSize)
		input := From(5)

		result := toSlice(expand(input))

		assert.Equal(t, []int{5, 10, 15}, result)
	})

	t.Run("handles empty input", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return From(n, n*2)
		}, defaultBufferSize)
		input := Empty[int]()

		result := toSlice(expand(input))

		assert.Empty(t, result)
	})

	t.Run("function returns empty sequences", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return Empty[int]()
		}, defaultBufferSize)
		input := From(1, 2, 3)

		result := toSlice(expand(input))

		assert.Empty(t, result)
	})

	t.Run("mixed empty and non-empty sequences preserve order", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			if n%2 == 0 {
				return Empty[int]()
			}
			return From(n, n*10)
		}, defaultBufferSize)
		input := From(1, 2, 3, 4, 5)

		result := toSlice(expand(input))

		assert.Equal(t, []int{1, 10, 3, 30, 5, 50}, result)
	})

	t.Run("order is deterministic across runs", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return From(n, n*10, n*100)
		}, defaultBufferSize)
		expected := []int{1, 10, 100, 2, 20, 200, 3, 30, 300}

		for range 5 {
			result := toSlice(expand(From(1, 2, 3)))
			assert.Equal(t, expected, result)
		}
	})
}

// TestConcatMapBuf_EarlyTermination tests early termination behavior
func TestConcatMapBuf_EarlyTermination(t *testing.T) {
	t.Run("stops after count elements", func(t *testing.T) {
		expand := ConcatMapBuf(func(n int) Seq[int] {
			return From(n, n*10, n*100)
		}, defaultBufferSize)
		input := From(1, 2, 3, 4, 5)

		var result []int
		for v := range expand(input) {
			result = append(result, v)
			if len(result) >= 4 {
				break
			}
		}

		assert.Equal(t, []int{1, 10, 100, 2}, result)
	})
}

// TestConcatMapBuf_vs_Chain verifies ConcatMapBuf produces the same results as Chain
func TestConcatMapBuf_vs_Chain(t *testing.T) {
	t.Run("produces same elements and order as Chain", func(t *testing.T) {
		f := func(n int) Seq[int] {
			return From(n, n*10)
		}
		input := From(1, 2, 3)

		concatResult := toSlice(ConcatMapBuf(f, defaultBufferSize)(input))
		chainResult := toSlice(Chain(f)(input))

		assert.Equal(t, chainResult, concatResult)
	})

	t.Run("produces same elements and order as ConcatMap", func(t *testing.T) {
		f := func(n int) Seq[int] {
			return From(n, n*10, n*100)
		}
		input := From(1, 2, 3)

		concatMapBufResult := toSlice(ConcatMapBuf(f, defaultBufferSize)(input))
		concatMapResult := toSlice(ConcatMap(f)(input))

		assert.Equal(t, concatMapResult, concatMapBufResult)
	})

	t.Run("comparison with MergeMap shows same elements but guaranteed order", func(t *testing.T) {
		f := func(n int) Seq[int] {
			return From(n, n*10)
		}
		input := From(1, 2, 3)

		concatResult := toSlice(ConcatMapBuf(f, defaultBufferSize)(input))

		// ConcatMapBuf always produces this exact order
		assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, concatResult)
	})
}

// TestMonadConcatMap_Success tests basic MonadConcatMap functionality
func TestMonadConcatMap_Success(t *testing.T) {
	t.Run("flatMaps in order", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := toSlice(MonadConcatMap(seq, func(n int) Seq[int] {
			return From(n, n*10)
		}))

		assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, result)
	})

	t.Run("handles empty input", func(t *testing.T) {
		result := toSlice(MonadConcatMap(Empty[int](), func(n int) Seq[int] {
			return From(n, n*10)
		}))

		assert.Empty(t, result)
	})

	t.Run("handles function returning empty sequences", func(t *testing.T) {
		result := toSlice(MonadConcatMap(From(1, 2, 3), func(n int) Seq[int] {
			return Empty[int]()
		}))

		assert.Empty(t, result)
	})

	t.Run("order is deterministic", func(t *testing.T) {
		seq := From(1, 2, 3)
		expected := []int{1, 10, 100, 2, 20, 200, 3, 30, 300}

		for range 5 {
			result := toSlice(MonadConcatMap(seq, func(n int) Seq[int] {
				return From(n, n*10, n*100)
			}))
			assert.Equal(t, expected, result)
		}
	})
}

// TestMonadConcatMap_vs_MonadChain verifies MonadConcatMap is equivalent to MonadChain
func TestMonadConcatMap_vs_MonadChain(t *testing.T) {
	t.Run("produces same result as MonadChain", func(t *testing.T) {
		f := func(n int) Seq[int] {
			return From(n, n*10)
		}
		input := From(1, 2, 3)

		concatMapResult := toSlice(MonadConcatMap(input, f))
		chainResult := toSlice(MonadChain(input, f))

		assert.Equal(t, chainResult, concatMapResult)
	})
}

// TestConcatMonoid_Sequential tests the ConcatMonoid from monoid.go via ConcatBuf
func TestConcatBuf_WithConcatMonoid(t *testing.T) {
	t.Run("ConcatMonoid concat equals ConcatBuf", func(t *testing.T) {
		monoid := ConcatMonoid[int](defaultBufferSize)
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		monoidResult := toSlice(monoid.Concat(seq1, seq2))
		concatResult := toSlice(Concat([]Seq[int]{seq1, seq2}))

		assert.Equal(t, monoidResult, concatResult)
	})

	t.Run("reduce multiple sequences using ConcatMonoid preserves order", func(t *testing.T) {
		monoid := ConcatMonoid[int](defaultBufferSize)
		sequences := From(
			From(1, 2),
			From(3, 4),
			From(5, 6),
		)

		result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
		collected := toSlice(result())

		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
	})
}

// Benchmark tests

func BenchmarkConcatBuf_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)
		for range result {
		}
	}
}

func BenchmarkConcatBuf_FiveSequences(b *testing.B) {
	sequences := make([]Seq[int], 5)
	for i := range 5 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		result := ConcatBuf(sequences, defaultBufferSize)
		for range result {
		}
	}
}

func BenchmarkConcatBuf_TenSequences(b *testing.B) {
	sequences := make([]Seq[int], 10)
	for i := range 10 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		result := ConcatBuf(sequences, defaultBufferSize)
		for range result {
		}
	}
}

func BenchmarkConcatBuf_LargeSequences(b *testing.B) {
	data1 := make([]int, 1000)
	data2 := make([]int, 1000)
	for i := range 1000 {
		data1[i] = i
		data2[i] = i + 1000
	}
	seq1 := From(data1...)
	seq2 := From(data2...)
	b.ResetTimer()
	for range b.N {
		result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)
		for range result {
		}
	}
}

func BenchmarkConcatMapBuf_ThreeElements(b *testing.B) {
	expand := ConcatMapBuf(func(n int) Seq[int] {
		return From(n, n*10, n*100)
	}, defaultBufferSize)
	input := From(1, 2, 3)
	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

func BenchmarkConcatMapBuf_WithFilter(b *testing.B) {
	expand := ConcatMapBuf(func(n int) Seq[int] {
		if n%2 == 0 {
			return From(n, n*10)
		}
		return Empty[int]()
	}, defaultBufferSize)
	input := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// Example tests

func ExampleConcatBuf() {
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)
	result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleConcatBuf_threeSequences() {
	seq1 := From(1, 2)
	seq2 := From(3, 4)
	seq3 := From(5, 6)
	result := ConcatBuf([]Seq[int]{seq1, seq2, seq3}, defaultBufferSize)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleConcatBuf_earlyTermination() {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	result := ConcatBuf([]Seq[int]{seq1, seq2}, defaultBufferSize)

	count := 0
	for range result {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Printf("Consumed %d elements\n", count)
	// Output: Consumed 3 elements
}

func ExampleConcatMapBuf() {
	expand := ConcatMapBuf(func(n int) Seq[int] {
		return From(n, n*10)
	}, defaultBufferSize)
	seq := From(1, 2, 3)

	for v := range expand(seq) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 10 2 20 3 30
}

func ExampleMonadConcatMap() {
	seq := From(1, 2, 3)
	result := MonadConcatMap(seq, func(n int) Seq[int] {
		return From(n, n*10)
	})

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 10 2 20 3 30
}

// ---------------------------------------------------------------------------
// ConcatAll tests
// ---------------------------------------------------------------------------

// TestConcatAll_Success verifies that all elements are produced and in the
// correct order for the basic happy-path cases.
func TestConcatAll_Success(t *testing.T) {
	t.Run("flattens two sequences in order", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("flattens three sequences in order", func(t *testing.T) {
		outer := From(From(1, 2), From(3, 4), From(5, 6))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("flattens single sequence", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("flattens many sequences in order", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2)
		}
		result := toSlice(ConcatAll[int](8)(From(sequences...)))
		assert.Equal(t, 30, len(result))
		expected := make([]int, 30)
		for i := range 30 {
			expected[i] = (i/3)*10 + (i % 3)
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with strings", func(t *testing.T) {
		outer := From(From("a", "b"), From("c", "d"), From("e"))
		result := toSlice(ConcatAll[string](8)(outer))
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, result)
	})

	t.Run("works with single-element inner sequences", func(t *testing.T) {
		outer := From(From(1), From(2), From(3))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestConcatAll_OrderPreservation is the key differentiator from MergeAll:
// ConcatAll always drains seq[n] completely before yielding any element of seq[n+1].
func TestConcatAll_OrderPreservation(t *testing.T) {
	t.Run("all elements from seq[0] precede all elements of seq[1]", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
		// Every element from the first group (1-3) must come before any from the second (4-6)
		maxFirst := 0
		for _, v := range result[:3] {
			if v > maxFirst {
				maxFirst = v
			}
		}
		assert.Less(t, maxFirst, result[3])
	})

	t.Run("order is deterministic across repeated runs", func(t *testing.T) {
		expected := make([]int, 0, 30)
		for i := range 10 {
			expected = append(expected, i*10, i*10+1, i*10+2)
		}
		for range 10 {
			sequences := make([]Seq[int], 10)
			for i := range 10 {
				sequences[i] = From(i*10, i*10+1, i*10+2)
			}
			result := toSlice(ConcatAll[int](8)(From(sequences...)))
			assert.Equal(t, expected, result)
		}
	})

	t.Run("order preserved even when later producers finish first", func(t *testing.T) {
		// seq[1] completes immediately; seq[0] is slow — seq[0] elements must still come first.
		seq0 := func(yield func(int) bool) {
			for _, v := range []int{1, 2, 3} {
				time.Sleep(5 * time.Millisecond)
				if !yield(v) {
					return
				}
			}
		}
		seq1 := From(4, 5, 6)
		result := toSlice(ConcatAll[int](8)(From(seq0, seq1)))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("order preserved with three slow sequences", func(t *testing.T) {
		make1ms := func(start int) Seq[int] {
			return func(yield func(int) bool) {
				for i := range 3 {
					time.Sleep(time.Millisecond)
					if !yield(start + i) {
						return
					}
				}
			}
		}
		outer := From(make1ms(0), make1ms(10), make1ms(20))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{0, 1, 2, 10, 11, 12, 20, 21, 22}, result)
	})
}

// TestConcatAll_BufferSizes verifies correctness across the full range of buffer sizes.
func TestConcatAll_BufferSizes(t *testing.T) {
	check := func(t *testing.T, bufSize int) {
		t.Helper()
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](bufSize)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	}

	t.Run("unbuffered (bufSize 0)", func(t *testing.T) { check(t, 0) })
	t.Run("buffer size 1", func(t *testing.T) { check(t, 1) })
	t.Run("small buffer (2)", func(t *testing.T) { check(t, 2) })
	t.Run("default buffer (8)", func(t *testing.T) { check(t, 8) })
	t.Run("large buffer (100)", func(t *testing.T) { check(t, 100) })
	t.Run("negative buffer treated as 0", func(t *testing.T) { check(t, -5) })
}

// TestConcatAll_Empty verifies correct handling of empty inputs.
func TestConcatAll_Empty(t *testing.T) {
	t.Run("empty outer sequence", func(t *testing.T) {
		result := toSlice(ConcatAll[int](8)(Empty[Seq[int]]()))
		assert.Empty(t, result)
	})

	t.Run("outer sequence with all empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), Empty[int](), Empty[int]())
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Empty(t, result)
	})

	t.Run("leading empty inner sequence", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2, 3))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("trailing empty inner sequence", func(t *testing.T) {
		outer := From(From(1, 2, 3), Empty[int]())
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("empty inner sequences interspersed preserve order", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2), Empty[int](), From(3, 4), Empty[int]())
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

// TestConcatAll_EarlyTermination verifies that all goroutines are cleaned up
// properly when the consumer stops before the sequences are exhausted.
func TestConcatAll_EarlyTermination(t *testing.T) {
	t.Run("stops after N elements", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		var collected []int
		for v := range ConcatAll[int](8)(outer) {
			collected = append(collected, v)
			if len(collected) >= 3 {
				break
			}
		}
		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("producers stop after consumer breaks — production bounded", func(t *testing.T) {
		var produced atomic.Int32
		makeSeq := func(start int) Seq[int] {
			return func(yield func(int) bool) {
				for i := range 100 {
					produced.Add(1)
					if !yield(start + i) {
						return
					}
				}
			}
		}
		outer := From(makeSeq(0), makeSeq(100), makeSeq(200))

		count := 0
		for range ConcatAll[int](8)(outer) {
			count++
			if count >= 5 {
				break
			}
		}
		time.Sleep(20 * time.Millisecond) // let goroutines notice done
		assert.LessOrEqual(t, produced.Load(), int32(50),
			"producers should stop shortly after consumer breaks")
		assert.GreaterOrEqual(t, produced.Load(), int32(5))
	})

	t.Run("early termination with unbuffered channel", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		var collected []int
		for v := range ConcatAll[int](0)(outer) {
			collected = append(collected, v)
			if len(collected) >= 2 {
				break
			}
		}
		assert.Equal(t, 2, len(collected))
	})

	t.Run("elements from later sequences are never yielded after early exit", func(t *testing.T) {
		// seq2's producer goroutine may start concurrently (ConcatAll pre-fetches),
		// but the drainer only forwards seq1's elements to the consumer before done
		// is closed, so the consumer never observes seq2 values.
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8)
		outer := From(seq1, seq2)
		var consumed []int
		for v := range ConcatAll[int](8)(outer) {
			consumed = append(consumed, v)
			if len(consumed) >= 3 {
				break
			}
		}
		assert.Equal(t, []int{1, 2, 3}, consumed, "consumer should see exactly 3 elements")
	})

	t.Run("early termination with many sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 20)
		for i := range 20 {
			sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
		}
		var collected []int
		for v := range ConcatAll[int](8)(From(sequences...)) {
			collected = append(collected, v)
			if len(collected) >= 7 {
				break
			}
		}
		assert.Equal(t, 7, len(collected))
	})
}

// TestConcatAll_Concurrency verifies that all inner producers start concurrently,
// even though their output is drained in order.
func TestConcatAll_Concurrency(t *testing.T) {
	t.Run("all inner producers start concurrently", func(t *testing.T) {
		const n = 4
		startTimes := make([]time.Time, n)
		var mu sync.Mutex

		sequences := make([]Seq[int], n)
		for i := range n {
			idx := i
			sequences[i] = func(yield func(int) bool) {
				mu.Lock()
				startTimes[idx] = time.Now()
				mu.Unlock()
				time.Sleep(20 * time.Millisecond)
				yield(idx)
			}
		}

		_ = toSlice(ConcatAll[int](8)(From(sequences...)))

		mu.Lock()
		defer mu.Unlock()
		// All goroutines should have started within 20 ms of each other.
		for i := 1; i < n; i++ {
			diff := startTimes[i].Sub(startTimes[0])
			if diff < 0 {
				diff = -diff
			}
			assert.Less(t, diff, 20*time.Millisecond,
				"producer %d should start concurrently with producer 0", i)
		}
	})

	t.Run("slow outer sequence does not block draining of first inner channel", func(t *testing.T) {
		// The outer goroutine sleeps between items; the drainer should begin
		// yielding from the first inner channel without waiting for the rest.
		firstItemAt := make(chan time.Time, 1)
		outerStartedAt := time.Now()

		outer := func(yield func(Seq[int]) bool) {
			outerStartedAt = time.Now()
			if !yield(From(1, 2, 3)) {
				return
			}
			time.Sleep(50 * time.Millisecond)
			yield(From(4, 5, 6))
		}

		var collected []int
		for v := range ConcatAll[int](8)(outer) {
			if len(collected) == 0 {
				firstItemAt <- time.Now()
			}
			collected = append(collected, v)
		}

		elapsed := (<-firstItemAt).Sub(outerStartedAt)
		assert.Less(t, elapsed, 30*time.Millisecond,
			"first item should arrive quickly without waiting for the whole outer sequence")
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
	})

	t.Run("concurrent production is faster than sequential would be", func(t *testing.T) {
		sleepPerItem := 5 * time.Millisecond
		nSeqs, nItems := 4, 3
		makeSlowSeq := func(start int) Seq[int] {
			return func(yield func(int) bool) {
				for i := range nItems {
					time.Sleep(sleepPerItem)
					if !yield(start + i) {
						return
					}
				}
			}
		}
		sequences := make([]Seq[int], nSeqs)
		for i := range nSeqs {
			sequences[i] = makeSlowSeq(i * 10)
		}

		start := time.Now()
		result := toSlice(ConcatAll[int](8)(From(sequences...)))
		elapsed := time.Since(start)

		assert.Equal(t, nSeqs*nItems, len(result))
		// Sequential would take nSeqs * nItems * sleepPerItem; concurrent should be
		// roughly nItems * sleepPerItem (all producers running in parallel).
		sequential := time.Duration(nSeqs*nItems) * sleepPerItem
		concurrent := time.Duration(nItems) * sleepPerItem
		assert.Less(t, elapsed, sequential,
			"concurrent production should be faster than sequential")
		// Allow 3× slack for scheduling jitter.
		assert.Less(t, elapsed, 3*concurrent+30*time.Millisecond,
			"elapsed should be close to single-sequence time")
	})
}

// TestConcatAll_WithComplexTypes verifies that ConcatAll works with non-primitive element types.
func TestConcatAll_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with structs, preserving order", func(t *testing.T) {
		outer := From(
			From(Person{"Alice", 30}, Person{"Bob", 25}),
			From(Person{"Charlie", 35}, Person{"Diana", 28}),
		)
		result := toSlice(ConcatAll[Person](8)(outer))
		assert.Equal(t, []Person{
			{"Alice", 30}, {"Bob", 25}, {"Charlie", 35}, {"Diana", 28},
		}, result)
	})

	t.Run("works with pointers, identity preserved", func(t *testing.T) {
		p1, p2, p3 := &Person{"Alice", 30}, &Person{"Bob", 25}, &Person{"Charlie", 35}
		outer := From(From(p1, p2), From(p3))
		result := toSlice(ConcatAll[*Person](8)(outer))
		assert.Equal(t, []*Person{p1, p2, p3}, result)
	})

	t.Run("works with slice elements", func(t *testing.T) {
		outer := From(
			From([]int{1, 2}, []int{3, 4}),
			From([]int{5, 6}),
		)
		result := toSlice(ConcatAll[[]int](8)(outer))
		assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}}, result)
	})
}

// TestConcatAll_WithDynamicGeneration tests using ConcatAll with sequences produced on-the-fly.
func TestConcatAll_WithDynamicGeneration(t *testing.T) {
	t.Run("sequences generated by Map, order preserved", func(t *testing.T) {
		outer := MonadMap(From(1, 2, 3), func(n int) Seq[int] {
			return From(n, n*10, n*100)
		})
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 10, 100, 2, 20, 200, 3, 30, 300}, result)
	})

	t.Run("filters zeros via empty inner sequences, order preserved", func(t *testing.T) {
		outer := MonadMap(From(1, 2, 3, 4, 5), func(n int) Seq[int] {
			if n%2 == 0 {
				return Empty[int]()
			}
			return From(n, n*10)
		})
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 10, 3, 30, 5, 50}, result)
	})

	t.Run("works with map then concatAll", func(t *testing.T) {
		double := ConcatAll[int](8)(MonadMap(From(1, 2, 3), func(n int) Seq[int] {
			return From(n, n)
		}))
		assert.Equal(t, []int{1, 1, 2, 2, 3, 3}, toSlice(double))
	})
}

// TestConcatAll_vs_MergeAll confirms that ConcatAll and MergeAll contain the same
// elements but that ConcatAll guarantees deterministic order while MergeAll does not.
func TestConcatAll_vs_MergeAll(t *testing.T) {
	t.Run("same elements as MergeAll", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		seq3 := From(7, 8, 9)

		concatResult := toSlice(ConcatAll[int](8)(From(seq1, seq2, seq3)))
		mergeResult := toSlice(MergeAll[int](8)(From(seq1, seq2, seq3)))

		assert.ElementsMatch(t, mergeResult, concatResult)
	})

	t.Run("ConcatAll order is always deterministic", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5, 6}
		for range 20 {
			result := toSlice(ConcatAll[int](8)(From(From(1, 2, 3), From(4, 5, 6))))
			assert.Equal(t, expected, result)
		}
	})
}

// TestConcatAll_WithChainedOperations verifies composability with other operators.
func TestConcatAll_WithChainedOperations(t *testing.T) {
	t.Run("map after concatAll", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(MonadMap(ConcatAll[int](8)(outer), N.Mul(2)))
		assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, result)
	})

	t.Run("filter after concatAll preserves order", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4), From(5, 6, 7, 8))
		result := toSlice(MonadFilter(ConcatAll[int](8)(outer), func(x int) bool { return x%2 == 0 }))
		assert.Equal(t, []int{2, 4, 6, 8}, result)
	})

	t.Run("concatAll after map preserves order", func(t *testing.T) {
		outer := MonadMap(From(From(1, 2), From(3, 4)), func(s Seq[int]) Seq[int] {
			return MonadMap(s, N.Mul(10))
		})
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{10, 20, 30, 40}, result)
	})
}

// TestConcatAll_EdgeCases covers less obvious scenarios.
func TestConcatAll_EdgeCases(t *testing.T) {
	t.Run("single element in single inner sequence", func(t *testing.T) {
		result := toSlice(ConcatAll[int](8)(From(From(42))))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("many single-element inner sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 100)
		for i := range 100 {
			sequences[i] = From(i)
		}
		result := toSlice(ConcatAll[int](8)(From(sequences...)))
		expected := make([]int, 100)
		for i := range 100 {
			expected[i] = i
		}
		assert.Equal(t, expected, result)
	})

	t.Run("inner sequences of varying lengths", func(t *testing.T) {
		outer := From(From(1), From(2, 3), From(4, 5, 6), From(7, 8, 9, 10))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)
	})

	t.Run("outer sequence produced slowly", func(t *testing.T) {
		outer := func(yield func(Seq[int]) bool) {
			for i := range 3 {
				time.Sleep(5 * time.Millisecond)
				if !yield(From(i*10, i*10+1)) {
					return
				}
			}
		}
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{0, 1, 10, 11, 20, 21}, result)
	})

	t.Run("large inner sequences", func(t *testing.T) {
		data1 := make([]int, 500)
		data2 := make([]int, 500)
		for i := range 500 {
			data1[i] = i
			data2[i] = i + 500
		}
		result := toSlice(ConcatAll[int](16)(From(From(data1...), From(data2...))))
		assert.Equal(t, 1000, len(result))
		for i, v := range result {
			assert.Equal(t, i, v, "element at index %d should be %d", i, i)
		}
	})
}

// Benchmarks

func BenchmarkConcatAll_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](8)(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAll_FiveSequences(b *testing.B) {
	sequences := make([]Seq[int], 5)
	for i := range 5 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](8)(From(sequences...)) {
		}
	}
}

func BenchmarkConcatAll_TenSequences(b *testing.B) {
	sequences := make([]Seq[int], 10)
	for i := range 10 {
		sequences[i] = From(i*10, i*10+1, i*10+2)
	}
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](8)(From(sequences...)) {
		}
	}
}

func BenchmarkConcatAll_Unbuffered(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](0)(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAll_LargeBuffer(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](100)(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAll_LargeSequences(b *testing.B) {
	data1 := make([]int, 500)
	data2 := make([]int, 500)
	for i := range 500 {
		data1[i] = i
		data2[i] = i + 500
	}
	seq1 := From(data1...)
	seq2 := From(data2...)
	b.ResetTimer()
	for range b.N {
		for range ConcatAll[int](16)(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAll_DynamicGeneration(b *testing.B) {
	b.ResetTimer()
	for range b.N {
		outer := MonadMap(From(1, 2, 3, 4, 5), func(n int) Seq[int] {
			return From(n, n*10, n*100)
		})
		for range ConcatAll[int](8)(outer) {
		}
	}
}

// Examples

func ExampleConcatAll() {
	outer := From(From(1, 2, 3), From(4, 5, 6))
	for v := range ConcatAll[int](8)(outer) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleConcatAll_threeSequences() {
	outer := From(From(1, 2), From(3, 4), From(5, 6))
	for v := range ConcatAll[int](8)(outer) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleConcatAll_earlyTermination() {
	outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
	count := 0
	for range ConcatAll[int](8)(outer) {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Printf("Consumed %d elements\n", count)
	// Output: Consumed 3 elements
}

func ExampleConcatAll_dynamicGeneration() {
	outer := MonadMap(From(1, 2, 3), func(n int) Seq[int] {
		return From(n, n*10)
	})
	for v := range ConcatAll[int](8)(outer) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 10 2 20 3 30
}

// ---------------------------------------------------------------------------
// ConcatAllSeq tests
// ---------------------------------------------------------------------------

// TestConcatAllSeq_Success verifies basic output correctness.
func TestConcatAllSeq_Success(t *testing.T) {
	t.Run("flattens two inner sequences", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("flattens three inner sequences", func(t *testing.T) {
		outer := From(From(1, 2), From(3, 4), From(5, 6))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("flattens single inner sequence", func(t *testing.T) {
		outer := From(From(42, 43, 44))
		assert.Equal(t, []int{42, 43, 44}, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("works with strings", func(t *testing.T) {
		outer := From(From("a", "b"), From("c"))
		assert.Equal(t, []string{"a", "b", "c"}, toSlice(ConcatAllSeq[string]()(outer)))
	})
}

// TestConcatAllSeq_Empty verifies correct handling of empty inputs.
func TestConcatAllSeq_Empty(t *testing.T) {
	t.Run("empty outer sequence", func(t *testing.T) {
		assert.Empty(t, toSlice(ConcatAllSeq[int]()(Empty[Seq[int]]())))
	})

	t.Run("all empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), Empty[int]())
		assert.Empty(t, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("leading empty inner sequence", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2))
		assert.Equal(t, []int{1, 2}, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("trailing empty inner sequence", func(t *testing.T) {
		outer := From(From(1, 2), Empty[int]())
		assert.Equal(t, []int{1, 2}, toSlice(ConcatAllSeq[int]()(outer)))
	})

	t.Run("interspersed empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2), Empty[int](), From(3))
		assert.Equal(t, []int{1, 2, 3}, toSlice(ConcatAllSeq[int]()(outer)))
	})
}

// TestConcatAllSeq_EarlyTermination verifies that the consumer can break out
// of the range loop and that iteration stops immediately.
func TestConcatAllSeq_EarlyTermination(t *testing.T) {
	t.Run("stops mid-first-inner-sequence", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8))
		var got []int
		for v := range ConcatAllSeq[int]()(outer) {
			got = append(got, v)
			if len(got) == 2 {
				break
			}
		}
		assert.Equal(t, []int{1, 2}, got)
	})

	t.Run("stops exactly at inner sequence boundary", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		var got []int
		for v := range ConcatAllSeq[int]()(outer) {
			got = append(got, v)
			if len(got) == 3 {
				break
			}
		}
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("second inner sequence is never started after early exit", func(t *testing.T) {
		seq2Started := false
		seq2 := func(yield func(int) bool) {
			seq2Started = true
			yield(4)
		}
		outer := From(From(1, 2, 3), seq2)
		for v := range ConcatAllSeq[int]()(outer) {
			if v == 2 {
				break
			}
		}
		assert.False(t, seq2Started, "seq2 should not be started after early exit")
	})
}

// TestConcatAllSeq_IsSequential verifies that ConcatAllSeq processes inner
// sequences one after another — seq[n+1] does not start until seq[n] is done.
// This is the key semantic difference from ConcatAllPar.
func TestConcatAllSeq_IsSequential(t *testing.T) {
	var order []string
	seq1 := func(yield func(int) bool) {
		order = append(order, "seq1-start")
		yield(1)
		yield(2)
		order = append(order, "seq1-end")
	}
	seq2 := func(yield func(int) bool) {
		order = append(order, "seq2-start")
		yield(3)
		order = append(order, "seq2-end")
	}

	_ = toSlice(ConcatAllSeq[int]()(From[Seq[int]](seq1, seq2)))

	assert.Equal(t, []string{"seq1-start", "seq1-end", "seq2-start", "seq2-end"}, order,
		"seq2 must start only after seq1 is fully drained")
}

// TestConcatAllSeq_MatchesConcatAll_BufSize1 verifies that ConcatAllSeq and
// ConcatAll(1) produce identical output — they must be the same code path.
func TestConcatAllSeq_MatchesConcatAll_BufSize1(t *testing.T) {
	cases := [][]Seq[int]{
		{From(1, 2, 3), From(4, 5, 6)},
		{Empty[int](), From(1), Empty[int]()},
		{From(10, 20)},
	}
	for _, seqs := range cases {
		outer := From(seqs...)
		direct := toSlice(ConcatAllSeq[int]()(outer))
		viaAll := toSlice(ConcatAll[int](1)(outer))
		assert.Equal(t, direct, viaAll)
	}
}

// TestConcatAllSeq_ManySequences verifies correct ordering across a large
// number of inner sequences, ensuring each group appears before the next.
func TestConcatAllSeq_ManySequences(t *testing.T) {
	const n = 20
	sequences := make([]Seq[int], n)
	for i := range n {
		sequences[i] = From(i*10, i*10+1, i*10+2)
	}
	result := toSlice(ConcatAllSeq[int]()(From(sequences...)))

	require.Len(t, result, n*3)
	expected := make([]int, n*3)
	for i := range n {
		expected[i*3] = i * 10
		expected[i*3+1] = i*10 + 1
		expected[i*3+2] = i*10 + 2
	}
	assert.Equal(t, expected, result)
}

// TestConcatAllSeq_SingleElementInnerSeqs verifies that sequences containing
// exactly one element each are concatenated in the correct order.
func TestConcatAllSeq_SingleElementInnerSeqs(t *testing.T) {
	seqs := make([]Seq[int], 5)
	for i := range 5 {
		seqs[i] = From(i)
	}
	result := toSlice(ConcatAllSeq[int]()(From(seqs...)))
	assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
}

// TestConcatAllSeq_ComplexTypes verifies that ConcatAllSeq works correctly
// with struct elements, preserving both order and identity.
func TestConcatAllSeq_ComplexTypes(t *testing.T) {
	type Point struct{ X, Y int }
	seq1 := From(Point{1, 2}, Point{3, 4})
	seq2 := From(Point{5, 6})
	result := toSlice(ConcatAllSeq[Point]()(From(seq1, seq2)))
	assert.Equal(t, []Point{{1, 2}, {3, 4}, {5, 6}}, result)
}

// TestConcatAllSeq_OperatorReuse verifies that the Operator returned by
// ConcatAllSeq can be applied more than once without interference.
func TestConcatAllSeq_OperatorReuse(t *testing.T) {
	op := ConcatAllSeq[int]()

	r1 := toSlice(op(From(From(1, 2), From(3, 4))))
	r2 := toSlice(op(From(From(5, 6), From(7, 8))))

	assert.Equal(t, []int{1, 2, 3, 4}, r1)
	assert.Equal(t, []int{5, 6, 7, 8}, r2)
}

// TestConcatAllSeq_DynamicOuterSeq verifies that ConcatAllSeq correctly
// handles an outer sequence whose elements are produced lazily.
func TestConcatAllSeq_DynamicOuterSeq(t *testing.T) {
	callCount := 0
	outer := func(yield func(Seq[int]) bool) {
		for i := range 3 {
			callCount++
			if !yield(From(i*10, i*10+1)) {
				return
			}
		}
	}
	result := toSlice(ConcatAllSeq[int]()(outer))
	assert.Equal(t, []int{0, 1, 10, 11, 20, 21}, result)
	assert.Equal(t, 3, callCount)
}

// TestConcatAllSeq_EarlyTerminationCountsProducerCalls verifies that
// ConcatAllSeq never calls the outer producer after the consumer breaks,
// confirming there is no pre-fetching or buffering.
func TestConcatAllSeq_EarlyTerminationStopsOuter(t *testing.T) {
	outerCalls := 0
	outer := func(yield func(Seq[int]) bool) {
		for i := range 10 {
			outerCalls++
			if !yield(From(i)) {
				return
			}
		}
	}
	var got []int
	for v := range ConcatAllSeq[int]()(outer) {
		got = append(got, v)
		if len(got) == 3 {
			break
		}
	}
	assert.Equal(t, []int{0, 1, 2}, got)
	assert.Equal(t, 3, outerCalls,
		"outer producer should be called exactly once per consumed element")
}

// ---------------------------------------------------------------------------
// ConcatAllPar direct tests
// ---------------------------------------------------------------------------

// TestConcatAllPar_Direct exercises ConcatAllPar at several buffer sizes,
// bypassing the ConcatAll dispatch table, to confirm that correctness is not
// an artefact of the bufSize==1 sequential path.
func TestConcatAllPar_Direct(t *testing.T) {
	for _, bufSize := range []int{-1, 0, 2, 8, 100} {
		bufSize := bufSize
		t.Run(fmt.Sprintf("bufSize=%d", bufSize), func(t *testing.T) {
			outer := From(From(1, 2, 3), From(4, 5, 6))
			result := toSlice(ConcatAllPar[int](bufSize)(outer))
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
		})
	}
}

// TestConcatAllPar_Empty verifies empty-input edge cases for ConcatAllPar.
func TestConcatAllPar_Empty(t *testing.T) {
	t.Run("empty outer sequence", func(t *testing.T) {
		assert.Empty(t, toSlice(ConcatAllPar[int](8)(Empty[Seq[int]]())))
	})
	t.Run("all empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), Empty[int](), Empty[int]())
		assert.Empty(t, toSlice(ConcatAllPar[int](8)(outer)))
	})
	t.Run("leading empty inner sequence", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2, 3))
		assert.Equal(t, []int{1, 2, 3}, toSlice(ConcatAllPar[int](8)(outer)))
	})
	t.Run("trailing empty inner sequence", func(t *testing.T) {
		outer := From(From(1, 2, 3), Empty[int]())
		assert.Equal(t, []int{1, 2, 3}, toSlice(ConcatAllPar[int](8)(outer)))
	})
	t.Run("interspersed empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), From(1, 2), Empty[int](), From(3))
		assert.Equal(t, []int{1, 2, 3}, toSlice(ConcatAllPar[int](8)(outer)))
	})
}

// TestConcatAllPar_OrderPreservation verifies the core guarantee: output is
// always in input order even when inner sequences complete at different rates.
func TestConcatAllPar_OrderPreservation(t *testing.T) {
	t.Run("slow first, fast second: first still comes first", func(t *testing.T) {
		slow := func(yield func(int) bool) {
			time.Sleep(20 * time.Millisecond)
			yield(1)
		}
		result := toSlice(ConcatAllPar[int](8)(From[Seq[int]](slow, From(2))))
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("deterministic across 20 runs", func(t *testing.T) {
		expected := []int{1, 2, 3, 4, 5, 6}
		for range 20 {
			result := toSlice(ConcatAllPar[int](8)(From(From(1, 2, 3), From(4, 5, 6))))
			assert.Equal(t, expected, result)
		}
	})

	t.Run("10 sequences, all buffer sizes, correct group order", func(t *testing.T) {
		expected := make([]int, 30)
		for i := range 30 {
			expected[i] = i
		}
		for _, bufSize := range []int{0, 2, 8} {
			seqs := make([]Seq[int], 10)
			for i := range 10 {
				seqs[i] = From(i*3, i*3+1, i*3+2)
			}
			result := toSlice(ConcatAllPar[int](bufSize)(From(seqs...)))
			assert.Equal(t, expected, result, "bufSize=%d", bufSize)
		}
	})
}

// TestConcatAllPar_EarlyTermination verifies that all goroutines stop promptly
// when the consumer breaks out of the range loop.
func TestConcatAllPar_EarlyTermination(t *testing.T) {
	t.Run("consumer stops after first element", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		var got []int
		for v := range ConcatAllPar[int](8)(outer) {
			got = append(got, v)
			break
		}
		assert.Equal(t, []int{1}, got)
	})

	t.Run("infinite inner producer stops after consumer breaks", func(t *testing.T) {
		var produced atomic.Int64
		inf := func(yield func(int) bool) {
			for i := 0; ; i++ {
				produced.Add(1)
				if !yield(i) {
					return
				}
			}
		}
		var got []int
		for v := range ConcatAllPar[int](8)(From[Seq[int]](inf, inf)) {
			got = append(got, v)
			if len(got) == 5 {
				break
			}
		}
		assert.Len(t, got, 5)
		// Production should stop once done is closed. Verify with two snapshots.
		time.Sleep(30 * time.Millisecond)
		snap1 := produced.Load()
		time.Sleep(30 * time.Millisecond)
		snap2 := produced.Load()
		assert.Equal(t, snap1, snap2,
			"infinite producer should have stopped after consumer broke out")
	})

	t.Run("elements observed by consumer are a strict prefix", func(t *testing.T) {
		// The drainer forwards values in order; breaking after 3 elements must
		// yield exactly [1, 2, 3] regardless of how far producers have run.
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8))
		var got []int
		for v := range ConcatAllPar[int](8)(outer) {
			got = append(got, v)
			if len(got) == 3 {
				break
			}
		}
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("unbuffered channels: early termination does not deadlock", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8))
		var got []int
		for v := range ConcatAllPar[int](0)(outer) {
			got = append(got, v)
			if len(got) == 2 {
				break
			}
		}
		assert.Equal(t, 2, len(got))
	})
}

// TestConcatAllPar_UnbufferedChannels confirms that bufSize=0 (unbuffered ch,
// inners, and out channels) neither deadlocks nor corrupts output order.
func TestConcatAllPar_UnbufferedChannels(t *testing.T) {
	t.Run("two non-empty sequences", func(t *testing.T) {
		result := toSlice(ConcatAllPar[int](0)(From(From(1, 2, 3), From(4, 5, 6))))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})
	t.Run("single element per inner sequence", func(t *testing.T) {
		seqs := make([]Seq[int], 5)
		for i := range 5 {
			seqs[i] = From(i)
		}
		assert.Equal(t, []int{0, 1, 2, 3, 4}, toSlice(ConcatAllPar[int](0)(From(seqs...))))
	})
	t.Run("empty outer", func(t *testing.T) {
		assert.Empty(t, toSlice(ConcatAllPar[int](0)(Empty[Seq[int]]())))
	})
	t.Run("all empty inners", func(t *testing.T) {
		outer := From(Empty[int](), Empty[int]())
		assert.Empty(t, toSlice(ConcatAllPar[int](0)(outer)))
	})
	t.Run("slow first, fast second", func(t *testing.T) {
		slow := func(yield func(int) bool) {
			time.Sleep(10 * time.Millisecond)
			yield(1)
		}
		result := toSlice(ConcatAllPar[int](0)(From[Seq[int]](slow, From(2))))
		assert.Equal(t, []int{1, 2}, result)
	})
}

// TestConcatAllPar_OperatorReuse verifies that the Operator value returned by
// ConcatAllPar is safe to invoke more than once sequentially.
func TestConcatAllPar_OperatorReuse(t *testing.T) {
	op := ConcatAllPar[int](8)
	r1 := toSlice(op(From(From(1, 2), From(3, 4))))
	r2 := toSlice(op(From(From(5, 6), From(7, 8))))
	assert.Equal(t, []int{1, 2, 3, 4}, r1)
	assert.Equal(t, []int{5, 6, 7, 8}, r2)
}

// ---------------------------------------------------------------------------
// ConcatAllPar race-detection tests
//
// These tests exercise concurrent access patterns that the race detector
// (go test -race) would flag if the implementation has unsynchronised state.
// Each test creates multiple goroutines that invoke ConcatAllPar independently;
// all shared state inside ConcatAllPar is local per invocation, so there must
// be no races.
// ---------------------------------------------------------------------------

// TestConcatAllPar_Race_ConcurrentInvocations calls the same Operator value
// from many goroutines simultaneously.  Each call creates independent channels
// and goroutines, so no sharing across invocations should occur.
func TestConcatAllPar_Race_ConcurrentInvocations(t *testing.T) {
	op := ConcatAllPar[int](8)
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := toSlice(op(From(From(1, 2, 3), From(4, 5, 6))))
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
		}()
	}
	wg.Wait()
}

// TestConcatAllPar_Race_ConcurrentEarlyTermination breaks early from many
// concurrent invocations; close(done) must execute exactly once per invocation
// without racing against the goroutines that read from done.
func TestConcatAllPar_Race_ConcurrentEarlyTermination(t *testing.T) {
	op := ConcatAllPar[int](8)
	var wg sync.WaitGroup
	for range 20 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var got []int
			for v := range op(From(From(1, 2, 3, 4, 5), From(6, 7, 8))) {
				got = append(got, v)
				if len(got) == 2 {
					break
				}
			}
			assert.Equal(t, 2, len(got))
		}()
	}
	wg.Wait()
}

// TestConcatAllPar_Race_ProducerConsumerContention stresses the channel
// handoff between inner producers (Iₙ), the drainer (D), and the consumer
// under high throughput to surface any unsynchronised shared state.
func TestConcatAllPar_Race_ProducerConsumerContention(t *testing.T) {
	var totalEmitted atomic.Int64
	makeSeq := func(n int) Seq[int] {
		return func(yield func(int) bool) {
			for i := range n {
				totalEmitted.Add(1)
				if !yield(i) {
					return
				}
			}
		}
	}
	seqs := make([]Seq[int], 8)
	for i := range 8 {
		seqs[i] = makeSeq(50)
	}
	count := 0
	for range ConcatAllPar[int](8)(From(seqs...)) {
		count++
	}
	assert.Equal(t, 8*50, count)
	assert.Equal(t, int64(8*50), totalEmitted.Load())
}

// TestConcatAllPar_Race_DoneClosedConcurrently verifies that goroutines safely
// observe done being closed while they are in the middle of channel operations.
func TestConcatAllPar_Race_DoneClosedConcurrently(t *testing.T) {
	var wg sync.WaitGroup
	for range 10 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			var produced atomic.Int64
			makeInf := func(yield func(int) bool) {
				for i := 0; ; i++ {
					produced.Add(1)
					if !yield(i) {
						return
					}
				}
			}
			count := 0
			for range ConcatAllPar[int](4)(From[Seq[int]](makeInf, makeInf, makeInf)) {
				count++
				if count == 10 {
					break
				}
			}
			assert.Equal(t, 10, count)
		}()
	}
	wg.Wait()
}

// ---------------------------------------------------------------------------
// ConcatAllPar benchmarks
// ---------------------------------------------------------------------------

func BenchmarkConcatAllPar_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range ConcatAllPar[int](8)(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAllPar_LargeSequences(b *testing.B) {
	data1 := make([]int, 500)
	data2 := make([]int, 500)
	for i := range 500 {
		data1[i] = i
		data2[i] = i + 500
	}
	seq1 := From(data1...)
	seq2 := From(data2...)
	b.ResetTimer()
	for range b.N {
		for range ConcatAllPar[int](8)(From(seq1, seq2)) {
		}
	}
}

// BenchmarkConcatAllPar_vs_Seq compares the parallel and sequential paths
// for CPU-bound (non-sleeping) sequences of moderate size.
func BenchmarkConcatAllPar_vs_Seq(b *testing.B) {
	data := make([]int, 200)
	for i := range data {
		data[i] = i
	}
	seqs := []Seq[int]{From(data...), From(data...), From(data...), From(data...)}

	b.Run("par-bufSize8", func(b *testing.B) {
		for range b.N {
			for range ConcatAllPar[int](8)(From(seqs...)) {
			}
		}
	})

	b.Run("par-bufSize0", func(b *testing.B) {
		for range b.N {
			for range ConcatAllPar[int](0)(From(seqs...)) {
			}
		}
	})
}

// ---------------------------------------------------------------------------
// ConcatAll routing / dispatch tests
// ---------------------------------------------------------------------------

// TestConcatAll_RoutingSemantics documents the non-monotonic dispatch:
// bufSize == 1 → sequential (no goroutines); bufSize == 0 → concurrent+unbuffered.
func TestConcatAll_RoutingSemantics(t *testing.T) {
	// With bufSize=1 (ConcatAllSeq), inner sequences are started strictly in order.
	t.Run("bufSize=1 is sequential: seq[n+1] starts only after seq[n] ends", func(t *testing.T) {
		var order []string
		seq1 := func(yield func(int) bool) {
			order = append(order, "seq1-start")
			yield(1)
			order = append(order, "seq1-end")
		}
		seq2 := func(yield func(int) bool) {
			order = append(order, "seq2-start")
			yield(2)
		}
		_ = toSlice(ConcatAll[int](1)(From[Seq[int]](seq1, seq2)))
		assert.Equal(t, []string{"seq1-start", "seq1-end", "seq2-start"}, order)
	})

	// With bufSize=0 (ConcatAllPar), inner producers start concurrently so
	// seq2-start can appear before seq1-end.
	t.Run("bufSize=0 output is still ordered even with concurrent producers", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](0)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	// Larger bufSizes also preserve deterministic output order.
	t.Run("bufSize=8 output is deterministic", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](8)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	// Negative bufSize is forwarded to ConcatAllPar (clamped to 0 internally).
	t.Run("negative bufSize produces correct ordered output", func(t *testing.T) {
		outer := From(From(1, 2, 3), From(4, 5, 6))
		result := toSlice(ConcatAll[int](-4)(outer))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})
}

// ---------------------------------------------------------------------------
// ConcatAllSeq benchmarks
// ---------------------------------------------------------------------------

func BenchmarkConcatAllSeq_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		for range ConcatAllSeq[int]()(From(seq1, seq2)) {
		}
	}
}

func BenchmarkConcatAllSeq_LargeSequences(b *testing.B) {
	data1 := make([]int, 500)
	data2 := make([]int, 500)
	for i := range 500 {
		data1[i] = i
		data2[i] = i + 500
	}
	seq1 := From(data1...)
	seq2 := From(data2...)
	b.ResetTimer()
	for range b.N {
		for range ConcatAllSeq[int]()(From(seq1, seq2)) {
		}
	}
}

// BenchmarkConcatAll_SeqVsPar compares the sequential (bufSize=1) and
// concurrent (bufSize=8) paths for CPU-bound sequences.
func BenchmarkConcatAll_SeqVsPar(b *testing.B) {
	data := make([]int, 200)
	for i := range data {
		data[i] = i
	}
	seqs := []Seq[int]{From(data...), From(data...), From(data...), From(data...)}

	b.Run("sequential-bufSize1", func(b *testing.B) {
		for range b.N {
			for range ConcatAll[int](1)(From(seqs...)) {
			}
		}
	})

	b.Run("concurrent-bufSize8", func(b *testing.B) {
		for range b.N {
			for range ConcatAll[int](8)(From(seqs...)) {
			}
		}
	})
}

// ---------------------------------------------------------------------------
// ConcatAllSeq examples
// ---------------------------------------------------------------------------

func ExampleConcatAllSeq() {
	outer := From(From(1, 2, 3), From(4, 5, 6))
	for v := range ConcatAllSeq[int]()(outer) {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleConcatAllSeq_earlyTermination() {
	outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8))
	count := 0
	for range ConcatAllSeq[int]()(outer) {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Printf("Consumed %d elements\n", count)
	// Output: Consumed 3 elements
}
