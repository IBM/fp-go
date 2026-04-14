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
	"sync/atomic"
	"testing"
	"time"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// TestAsync_Success tests basic Async functionality
func TestAsync_Success(t *testing.T) {
	t.Run("converts sequence to async with buffer", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 10)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("preserves element order", func(t *testing.T) {
		seq := From("a", "b", "c", "d", "e")
		async := AsyncBuf(seq, 5)
		result := toSlice(async)
		assert.Equal(t, []string{"a", "b", "c", "d", "e"}, result)
	})

	t.Run("works with single element", func(t *testing.T) {
		seq := From(42)
		async := AsyncBuf(seq, 1)
		result := toSlice(async)
		assert.Equal(t, []int{42}, result)
	})

	t.Run("works with large sequence", func(t *testing.T) {
		data := make([]int, 100)
		for i := range data {
			data[i] = i
		}
		seq := From(data...)
		async := AsyncBuf(seq, 20)
		result := toSlice(async)
		assert.Equal(t, data, result)
	})
}

// TestAsync_BufferSizes tests different buffer sizes
func TestAsync_BufferSizes(t *testing.T) {
	t.Run("unbuffered channel (bufSize 0)", func(t *testing.T) {
		seq := From(1, 2, 3)
		async := AsyncBuf(seq, 0)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("small buffer", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 2)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("large buffer", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 100)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("negative buffer size treated as 0", func(t *testing.T) {
		seq := From(1, 2, 3)
		async := AsyncBuf(seq, -5)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("buffer size equals sequence length", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 5)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("buffer size larger than sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		async := AsyncBuf(seq, 10)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestAsync_Empty tests Async with empty sequences
func TestAsync_Empty(t *testing.T) {
	t.Run("empty integer sequence", func(t *testing.T) {
		seq := Empty[int]()
		async := AsyncBuf(seq, 5)
		result := toSlice(async)
		assert.Empty(t, result)
	})

	t.Run("empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		async := AsyncBuf(seq, 10)
		result := toSlice(async)
		assert.Empty(t, result)
	})

	t.Run("empty with zero buffer", func(t *testing.T) {
		seq := Empty[int]()
		async := AsyncBuf(seq, 0)
		result := toSlice(async)
		assert.Empty(t, result)
	})
}

// TestAsync_EarlyTermination tests that Async properly handles early termination
func TestAsync_EarlyTermination(t *testing.T) {
	t.Run("stops producer when consumer breaks", func(t *testing.T) {
		var producerCount atomic.Int32

		// Create a sequence that tracks how many elements were produced
		seq := func(yield func(int) bool) {
			for i := range 100 {
				producerCount.Add(1)
				if !yield(i) {
					return
				}
			}
		}

		async := AsyncBuf(seq, 10)

		// Consume only 5 elements
		count := 0
		for range async {
			count++
			if count >= 5 {
				break
			}
		}

		// Give goroutine time to clean up
		time.Sleep(10 * time.Millisecond)

		// Producer should have stopped shortly after consumer stopped
		// It may produce a few extra due to buffering, but not all 100
		produced := producerCount.Load()
		assert.LessOrEqual(t, produced, int32(20), "producer should stop after consumer breaks")
		assert.GreaterOrEqual(t, produced, int32(5), "producer should produce at least what was consumed")
	})

	t.Run("handles yield returning false", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		async := AsyncBuf(seq, 5)

		collected := []int{}
		for v := range async {
			collected = append(collected, v)
			if v == 3 {
				break
			}
		}

		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("early termination with unbuffered channel", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 0)

		collected := []int{}
		for v := range async {
			collected = append(collected, v)
			if v == 2 {
				break
			}
		}

		assert.Equal(t, []int{1, 2}, collected)
	})
}

// TestAsync_WithComplexTypes tests Async with complex data types
func TestAsync_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with structs", func(t *testing.T) {
		seq := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
			Person{"Charlie", 35},
		)
		async := AsyncBuf(seq, 5)
		result := toSlice(async)
		expected := []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with pointers", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		p3 := &Person{"Charlie", 35}
		seq := From(p1, p2, p3)
		async := AsyncBuf(seq, 3)
		result := toSlice(async)
		assert.Equal(t, []*Person{p1, p2, p3}, result)
	})

	t.Run("works with slices", func(t *testing.T) {
		seq := From([]int{1, 2}, []int{3, 4}, []int{5, 6})
		async := AsyncBuf(seq, 2)
		result := toSlice(async)
		expected := [][]int{{1, 2}, {3, 4}, {5, 6}}
		assert.Equal(t, expected, result)
	})

	t.Run("works with maps", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"b": 2}
		m3 := map[string]int{"c": 3}
		seq := From(m1, m2, m3)
		async := AsyncBuf(seq, 3)
		result := toSlice(async)
		assert.Equal(t, []map[string]int{m1, m2, m3}, result)
	})
}

// TestAsync_WithChainedOperations tests Async with other sequence operations
func TestAsync_WithChainedOperations(t *testing.T) {
	t.Run("async after map", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, N.Mul(2))
		async := AsyncBuf(mapped, 5)
		result := toSlice(async)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("map after async", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 5)
		mapped := MonadMap(async, N.Mul(2))
		result := toSlice(mapped)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("async after filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		async := AsyncBuf(filtered, 5)
		result := toSlice(async)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("filter after async", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		async := AsyncBuf(seq, 5)
		filtered := MonadFilter(async, func(x int) bool { return x%2 == 0 })
		result := toSlice(filtered)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("async after chain", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		async := AsyncBuf(chained, 10)
		result := toSlice(async)
		assert.Equal(t, []int{1, 10, 2, 20, 3, 30}, result)
	})

	t.Run("multiple async operations", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async1 := AsyncBuf(seq, 3)
		async2 := AsyncBuf(async1, 2)
		result := toSlice(async2)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// TestAsync_Concurrency tests concurrent behavior
func TestAsync_Concurrency(t *testing.T) {
	t.Run("allows concurrent production and consumption", func(t *testing.T) {
		// Create a slow producer
		seq := func(yield func(int) bool) {
			for i := range 5 {
				time.Sleep(5 * time.Millisecond)
				if !yield(i) {
					return
				}
			}
		}

		async := AsyncBuf(seq, 10)

		result := toSlice(async)

		// Verify all elements are produced correctly
		assert.Equal(t, []int{0, 1, 2, 3, 4}, result)
	})

	t.Run("handles concurrent consumption safely", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		async := AsyncBuf(seq, 5)

		// Consume with some processing time
		var sum atomic.Int32
		for v := range async {
			sum.Add(int32(v))
			time.Sleep(1 * time.Millisecond)
		}

		assert.Equal(t, int32(55), sum.Load())
	})
}

// TestAsync_EdgeCases tests edge cases
func TestAsync_EdgeCases(t *testing.T) {
	t.Run("very large buffer size", func(t *testing.T) {
		seq := From(1, 2, 3)
		async := AsyncBuf(seq, 1000000)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("buffer size of 1", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		async := AsyncBuf(seq, 1)
		result := toSlice(async)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("works with replicate", func(t *testing.T) {
		seq := Replicate(5, 42)
		async := AsyncBuf(seq, 3)
		result := toSlice(async)
		assert.Equal(t, []int{42, 42, 42, 42, 42}, result)
	})

	t.Run("works with makeBy", func(t *testing.T) {
		seq := MakeBy(5, func(i int) int { return i * i })
		async := AsyncBuf(seq, 3)
		result := toSlice(async)
		assert.Equal(t, []int{0, 1, 4, 9, 16}, result)
	})
}

// Benchmark tests
func BenchmarkAsync(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 5)
		for range async {
		}
	}
}

func BenchmarkAsync_LargeSequence(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	seq := From(data...)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 100)
		for range async {
		}
	}
}

func BenchmarkAsync_SmallBuffer(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 1)
		for range async {
		}
	}
}

func BenchmarkAsync_LargeBuffer(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 100)
		for range async {
		}
	}
}

func BenchmarkAsync_Unbuffered(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 0)
		for range async {
		}
	}
}

func BenchmarkAsync_WithMap(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 5)
		mapped := MonadMap(async, N.Mul(2))
		for range mapped {
		}
	}
}

func BenchmarkAsync_WithFilter(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		async := AsyncBuf(seq, 5)
		filtered := MonadFilter(async, func(x int) bool { return x%2 == 0 })
		for range filtered {
		}
	}
}

// Example tests for documentation
func ExampleAsync() {
	seq := From(1, 2, 3, 4, 5)
	async := AsyncBuf(seq, 10)

	for v := range async {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5
}

func ExampleAsync_unbuffered() {
	seq := From(1, 2, 3)
	async := AsyncBuf(seq, 0)

	for v := range async {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3
}

func ExampleAsync_earlyTermination() {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	async := AsyncBuf(seq, 5)

	count := 0
	for v := range async {
		fmt.Printf("%d ", v)
		count++
		if count >= 3 {
			break
		}
	}
	// Output: 1 2 3
}

func ExampleAsync_withMap() {
	seq := From(1, 2, 3, 4, 5)
	async := AsyncBuf(seq, 5)
	doubled := MonadMap(async, N.Mul(2))

	for v := range doubled {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6 8 10
}

func ExampleAsync_withFilter() {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	async := AsyncBuf(seq, 5)
	evens := MonadFilter(async, func(x int) bool { return x%2 == 0 })

	for v := range evens {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6 8 10
}

// TestAsync2_Success tests basic Async2 functionality
func TestAsync2_Success(t *testing.T) {
	t.Run("converts Seq2 to async with buffer", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		async := Async2Buf(seq, 10)
		result := toMap(async)
		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("preserves key-value pairs order", func(t *testing.T) {
		seq := MonadZip(From("x", "y", "z"), From(10, 20, 30))
		async := Async2Buf(seq, 5)

		keys := []string{}
		values := []int{}
		for k, v := range async {
			keys = append(keys, k)
			values = append(values, v)
		}

		assert.Equal(t, []string{"x", "y", "z"}, keys)
		assert.Equal(t, []int{10, 20, 30}, values)
	})

	t.Run("works with single pair", func(t *testing.T) {
		seq := Of2("key", 42)
		async := Async2Buf(seq, 1)
		result := toMap(async)
		assert.Equal(t, map[string]int{"key": 42}, result)
	})

	t.Run("works with large Seq2", func(t *testing.T) {
		keys := make([]int, 100)
		values := make([]string, 100)
		for i := range keys {
			keys[i] = i
			values[i] = fmt.Sprintf("val%d", i)
		}
		seq := MonadZip(From(keys...), From(values...))
		async := Async2Buf(seq, 20)
		result := toMap(async)
		assert.Equal(t, 100, len(result))
		for i := range 100 {
			assert.Equal(t, fmt.Sprintf("val%d", i), result[i])
		}
	})
}

// TestAsync2_BufferSizes tests different buffer sizes
func TestAsync2_BufferSizes(t *testing.T) {
	t.Run("unbuffered channel (bufSize 0)", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		async := Async2Buf(seq, 0)
		result := toMap(async)
		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("negative buffer size treated as 0", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		async := Async2Buf(seq, -5)
		result := toMap(async)
		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("large buffer", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		async := Async2Buf(seq, 100)
		result := toMap(async)
		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})
}

// TestAsync2_Empty tests Async2 with empty sequences
func TestAsync2_Empty(t *testing.T) {
	t.Run("empty Seq2", func(t *testing.T) {
		seq := MonadZip(Empty[int](), Empty[string]())
		async := Async2Buf(seq, 5)
		result := toMap(async)
		assert.Empty(t, result)
	})
}

// TestAsync2_EarlyTermination tests that Async2 properly handles early termination
func TestAsync2_EarlyTermination(t *testing.T) {
	t.Run("stops producer when consumer breaks", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), From("a", "b", "c", "d", "e", "f", "g", "h", "i", "j"))
		async := Async2Buf(seq, 5)

		count := 0
		for range async {
			count++
			if count >= 3 {
				break
			}
		}

		assert.Equal(t, 3, count)
	})
}

// TestAsync2_WithChainedOperations tests Async2 with other operations
func TestAsync2_WithChainedOperations(t *testing.T) {
	t.Run("async2 after map", func(t *testing.T) {
		seq := MonadZip(From(1, 2, 3), From(10, 20, 30))
		mapped := MonadMapWithKey(seq, func(k, v int) int { return k + v })
		async := Async2Buf(mapped, 5)
		result := toMap(async)
		expected := map[int]int{1: 11, 2: 22, 3: 33}
		assert.Equal(t, expected, result)
	})
}

// TestToSeqPair_Success tests basic ToSeqPair functionality
func TestToSeqPair_Success(t *testing.T) {
	t.Run("converts Seq2 to Seq of Pairs", func(t *testing.T) {
		seq2 := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, 1, pair.Head(result[0]))
		assert.Equal(t, "a", pair.Tail(result[0]))
		assert.Equal(t, 2, pair.Head(result[1]))
		assert.Equal(t, "b", pair.Tail(result[1]))
		assert.Equal(t, 3, pair.Head(result[2]))
		assert.Equal(t, "c", pair.Tail(result[2]))
	})

	t.Run("preserves order", func(t *testing.T) {
		seq2 := MonadZip(From("x", "y", "z"), From(10, 20, 30))
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)

		assert.Equal(t, 3, len(result))
		for i, p := range result {
			expectedKey := string(rune('x' + i))
			expectedVal := (i + 1) * 10
			assert.Equal(t, expectedKey, pair.Head(p))
			assert.Equal(t, expectedVal, pair.Tail(p))
		}
	})

	t.Run("works with single pair", func(t *testing.T) {
		seq2 := Of2("key", 42)
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)

		assert.Equal(t, 1, len(result))
		assert.Equal(t, "key", pair.Head(result[0]))
		assert.Equal(t, 42, pair.Tail(result[0]))
	})
}

// TestToSeqPair_Empty tests ToSeqPair with empty sequences
func TestToSeqPair_Empty(t *testing.T) {
	t.Run("empty Seq2 produces empty Seq", func(t *testing.T) {
		seq2 := MonadZip(Empty[int](), Empty[string]())
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)
		assert.Empty(t, result)
	})
}

// TestToSeqPair_WithComplexTypes tests ToSeqPair with complex types
func TestToSeqPair_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with struct values", func(t *testing.T) {
		seq2 := MonadZip(
			From(1, 2, 3),
			From(Person{"Alice", 30}, Person{"Bob", 25}, Person{"Charlie", 35}),
		)
		pairs := ToSeqPair(seq2)
		result := toSlice(pairs)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, 1, pair.Head(result[0]))
		assert.Equal(t, Person{"Alice", 30}, pair.Tail(result[0]))
	})
}

// TestFromSeqPair_Success tests basic FromSeqPair functionality
func TestFromSeqPair_Success(t *testing.T) {
	t.Run("converts Seq of Pairs to Seq2", func(t *testing.T) {
		pairs := From(
			pair.MakePair(1, "a"),
			pair.MakePair(2, "b"),
			pair.MakePair(3, "c"),
		)
		seq2 := FromSeqPair(pairs)
		result := toMap(seq2)

		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("preserves order", func(t *testing.T) {
		pairs := From(
			pair.MakePair("x", 10),
			pair.MakePair("y", 20),
			pair.MakePair("z", 30),
		)
		seq2 := FromSeqPair(pairs)

		keys := []string{}
		values := []int{}
		for k, v := range seq2 {
			keys = append(keys, k)
			values = append(values, v)
		}

		assert.Equal(t, []string{"x", "y", "z"}, keys)
		assert.Equal(t, []int{10, 20, 30}, values)
	})

	t.Run("works with single pair", func(t *testing.T) {
		pairs := From(pair.MakePair("key", 42))
		seq2 := FromSeqPair(pairs)
		result := toMap(seq2)

		assert.Equal(t, map[string]int{"key": 42}, result)
	})
}

// TestFromSeqPair_Empty tests FromSeqPair with empty sequences
func TestFromSeqPair_Empty(t *testing.T) {
	t.Run("empty Seq produces empty Seq2", func(t *testing.T) {
		pairs := Empty[Pair[int, string]]()
		seq2 := FromSeqPair(pairs)
		result := toMap(seq2)
		assert.Empty(t, result)
	})
}

// TestFromSeqPair_WithComplexTypes tests FromSeqPair with complex types
func TestFromSeqPair_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with struct values", func(t *testing.T) {
		pairs := From(
			pair.MakePair(1, Person{"Alice", 30}),
			pair.MakePair(2, Person{"Bob", 25}),
			pair.MakePair(3, Person{"Charlie", 35}),
		)
		seq2 := FromSeqPair(pairs)
		result := toMap(seq2)

		expected := map[int]Person{
			1: {"Alice", 30},
			2: {"Bob", 25},
			3: {"Charlie", 35},
		}
		assert.Equal(t, expected, result)
	})
}

// TestRoundTrip tests that ToSeqPair and FromSeqPair are inverses
func TestRoundTrip(t *testing.T) {
	t.Run("ToSeqPair then FromSeqPair", func(t *testing.T) {
		original := MonadZip(From(1, 2, 3), From("a", "b", "c"))
		pairs := ToSeqPair(original)
		restored := FromSeqPair(pairs)
		result := toMap(restored)

		expected := map[int]string{1: "a", 2: "b", 3: "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("FromSeqPair then ToSeqPair", func(t *testing.T) {
		original := From(
			pair.MakePair(1, "a"),
			pair.MakePair(2, "b"),
			pair.MakePair(3, "c"),
		)
		seq2 := FromSeqPair(original)
		restored := ToSeqPair(seq2)
		result := toSlice(restored)

		assert.Equal(t, 3, len(result))
		assert.Equal(t, 1, pair.Head(result[0]))
		assert.Equal(t, "a", pair.Tail(result[0]))
	})
}

// Benchmark tests for Async2
func BenchmarkAsync2(b *testing.B) {
	seq := MonadZip(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), From("a", "b", "c", "d", "e", "f", "g", "h", "i", "j"))
	b.ResetTimer()
	for range b.N {
		async := Async2Buf(seq, 5)
		for range async {
		}
	}
}

func BenchmarkAsync2_LargeSequence(b *testing.B) {
	keys := make([]int, 1000)
	values := make([]string, 1000)
	for i := range keys {
		keys[i] = i
		values[i] = fmt.Sprintf("val%d", i)
	}
	seq := MonadZip(From(keys...), From(values...))
	b.ResetTimer()
	for range b.N {
		async := Async2Buf(seq, 100)
		for range async {
		}
	}
}

// Benchmark tests for FromSeqPair
func BenchmarkFromSeqPair(b *testing.B) {
	pairs := From(
		pair.MakePair(1, "a"),
		pair.MakePair(2, "b"),
		pair.MakePair(3, "c"),
		pair.MakePair(4, "d"),
		pair.MakePair(5, "e"),
	)
	b.ResetTimer()
	for range b.N {
		seq2 := FromSeqPair(pairs)
		for range seq2 {
		}
	}
}

func BenchmarkRoundTrip(b *testing.B) {
	seq := MonadZip(From(1, 2, 3, 4, 5), From("a", "b", "c", "d", "e"))
	b.ResetTimer()
	for range b.N {
		pairs := ToSeqPair(seq)
		restored := FromSeqPair(pairs)
		for range restored {
		}
	}
}

// Example tests for Async2
func ExampleAsync2() {
	seq := MonadZip(From(1, 2, 3), From("a", "b", "c"))
	async := Async2Buf(seq, 10)

	for k, v := range async {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 1: a
	// 2: b
	// 3: c
}

func ExampleAsync2_earlyTermination() {
	seq := MonadZip(From(1, 2, 3, 4, 5), From("a", "b", "c", "d", "e"))
	async := Async2Buf(seq, 5)

	count := 0
	for k, v := range async {
		fmt.Printf("%d: %s\n", k, v)
		count++
		if count >= 2 {
			break
		}
	}
	// Output:
	// 1: a
	// 2: b
}

// Example tests for FromSeqPair
func ExampleFromSeqPair() {
	pairs := From(
		pair.MakePair(1, "a"),
		pair.MakePair(2, "b"),
		pair.MakePair(3, "c"),
	)
	seq2 := FromSeqPair(pairs)

	for k, v := range seq2 {
		fmt.Printf("%d: %s\n", k, v)
	}
	// Output:
	// 1: a
	// 2: b
	// 3: c
}
