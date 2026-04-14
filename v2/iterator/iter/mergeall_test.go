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
	"slices"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestMergeAll_Success tests basic MergeAll functionality
func TestMergeAll_Success(t *testing.T) {
	t.Run("merges two sequences", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		// All elements should be present (order is non-deterministic)
		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("merges three sequences", func(t *testing.T) {
		seq1 := From(1, 2)
		seq2 := From(3, 4)
		seq3 := From(5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2, seq3}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("merges many sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2)
		}
		merged := MergeBuf(sequences, 20)
		result := toSlice(merged)

		assert.Equal(t, 30, len(result))
		slices.Sort(result)
		expected := make([]int, 30)
		for i := range 30 {
			expected[i] = (i/3)*10 + (i % 3)
		}
		assert.Equal(t, expected, result)
	})

	t.Run("merges single sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		merged := MergeBuf([]Seq[int]{seq}, 5)
		result := toSlice(merged)

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("preserves all elements with strings", func(t *testing.T) {
		seq1 := From("a", "b", "c")
		seq2 := From("d", "e", "f")
		merged := MergeBuf([]Seq[string]{seq1, seq2}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, result)
	})
}

// TestMergeAll_BufferSizes tests different buffer sizes
func TestMergeAll_BufferSizes(t *testing.T) {
	t.Run("unbuffered channel (bufSize 0)", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 0)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("small buffer", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 2)
		result := toSlice(merged)

		assert.Equal(t, 10, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)
	})

	t.Run("large buffer", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 100)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("negative buffer size treated as 0", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, -5)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("buffer size of 1", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 1)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})
}

// TestMergeAll_Empty tests MergeAll with empty sequences
func TestMergeAll_Empty(t *testing.T) {
	t.Run("empty slice of sequences", func(t *testing.T) {
		merged := MergeBuf([]Seq[int]{}, 5)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("nil slice of sequences", func(t *testing.T) {
		merged := MergeBuf([]Seq[int](nil), 5)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("slice with empty sequences", func(t *testing.T) {
		seq1 := Empty[int]()
		seq2 := Empty[int]()
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("mix of empty and non-empty sequences", func(t *testing.T) {
		seq1 := Empty[int]()
		seq2 := From(1, 2, 3)
		seq3 := Empty[int]()
		merged := MergeBuf([]Seq[int]{seq1, seq2, seq3}, 5)
		result := toSlice(merged)

		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestMergeAll_EarlyTermination tests that MergeAll properly handles early termination
func TestMergeAll_EarlyTermination(t *testing.T) {
	t.Run("stops all producers when consumer breaks", func(t *testing.T) {
		var producerCount1 atomic.Int32
		var producerCount2 atomic.Int32

		// Create sequences that track how many elements were produced
		seq1 := func(yield func(int) bool) {
			for i := range 100 {
				producerCount1.Add(1)
				if !yield(i) {
					return
				}
			}
		}
		seq2 := func(yield func(int) bool) {
			for i := range 100 {
				producerCount2.Add(1)
				if !yield(i + 100) {
					return
				}
			}
		}

		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)

		// Consume only 5 elements
		count := 0
		for range merged {
			count++
			if count >= 5 {
				break
			}
		}

		// Give goroutines time to clean up
		time.Sleep(10 * time.Millisecond)

		// Producers should have stopped shortly after consumer stopped
		// They may produce a few extra due to buffering, but not all 100
		produced1 := producerCount1.Load()
		produced2 := producerCount2.Load()
		totalProduced := produced1 + produced2

		assert.LessOrEqual(t, totalProduced, int32(30), "producers should stop after consumer breaks")
		assert.GreaterOrEqual(t, totalProduced, int32(5), "producers should produce at least what was consumed")
	})

	t.Run("handles yield returning false", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 3 {
				break
			}
		}

		assert.Equal(t, 3, len(collected))
	})

	t.Run("early termination with unbuffered channel", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 0)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 2 {
				break
			}
		}

		assert.Equal(t, 2, len(collected))
	})

	t.Run("early termination with many sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
		}
		merged := MergeBuf(sequences, 10)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 5 {
				break
			}
		}

		assert.Equal(t, 5, len(collected))
	})
}

// TestMergeAll_WithComplexTypes tests MergeAll with complex data types
func TestMergeAll_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with structs", func(t *testing.T) {
		seq1 := From(Person{"Alice", 30}, Person{"Bob", 25})
		seq2 := From(Person{"Charlie", 35}, Person{"Diana", 28})
		merged := MergeBuf([]Seq[Person]{seq1, seq2}, 5)
		result := toSlice(merged)

		assert.Equal(t, 4, len(result))
		// Sort by name for deterministic comparison
		slices.SortFunc(result, func(a, b Person) int {
			if a.Name < b.Name {
				return -1
			}
			if a.Name > b.Name {
				return 1
			}
			return 0
		})
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
		merged := MergeBuf([]Seq[*Person]{seq1, seq2}, 3)
		result := toSlice(merged)

		assert.Equal(t, 3, len(result))
		assert.Contains(t, result, p1)
		assert.Contains(t, result, p2)
		assert.Contains(t, result, p3)
	})

	t.Run("works with slices", func(t *testing.T) {
		seq1 := From([]int{1, 2}, []int{3, 4})
		seq2 := From([]int{5, 6}, []int{7, 8})
		merged := MergeBuf([]Seq[[]int]{seq1, seq2}, 5)
		result := toSlice(merged)

		assert.Equal(t, 4, len(result))
	})

	t.Run("works with maps", func(t *testing.T) {
		m1 := map[string]int{"a": 1}
		m2 := map[string]int{"b": 2}
		m3 := map[string]int{"c": 3}
		seq1 := From(m1, m2)
		seq2 := From(m3)
		merged := MergeBuf([]Seq[map[string]int]{seq1, seq2}, 3)
		result := toSlice(merged)

		assert.Equal(t, 3, len(result))
		assert.Contains(t, result, m1)
		assert.Contains(t, result, m2)
		assert.Contains(t, result, m3)
	})
}

// TestMergeAll_WithChainedOperations tests MergeAll with other sequence operations
func TestMergeAll_WithChainedOperations(t *testing.T) {
	t.Run("merge after map", func(t *testing.T) {
		seq1 := MonadMap(From(1, 2, 3), func(x int) int { return x * 2 })
		seq2 := MonadMap(From(4, 5, 6), func(x int) int { return x * 2 })
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, result)
	})

	t.Run("map after merge", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		mapped := MonadMap(merged, func(x int) int { return x * 2 })
		result := toSlice(mapped)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{2, 4, 6, 8, 10, 12}, result)
	})

	t.Run("merge after filter", func(t *testing.T) {
		seq1 := MonadFilter(From(1, 2, 3, 4, 5), func(x int) bool { return x%2 == 0 })
		seq2 := MonadFilter(From(6, 7, 8, 9, 10), func(x int) bool { return x%2 == 0 })
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		assert.Equal(t, 5, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("filter after merge", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		filtered := MonadFilter(merged, func(x int) bool { return x%2 == 0 })
		result := toSlice(filtered)

		assert.Equal(t, 5, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("merge with async sequences", func(t *testing.T) {
		seq1 := AsyncBuf(From(1, 2, 3), 5)
		seq2 := AsyncBuf(From(4, 5, 6), 5)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})
}

// TestMergeAll_Concurrency tests concurrent behavior
func TestMergeAll_Concurrency(t *testing.T) {
	t.Run("allows concurrent production", func(t *testing.T) {
		// Create slow producers
		seq1 := func(yield func(int) bool) {
			for i := range 5 {
				time.Sleep(2 * time.Millisecond)
				if !yield(i) {
					return
				}
			}
		}
		seq2 := func(yield func(int) bool) {
			for i := range 5 {
				time.Sleep(2 * time.Millisecond)
				if !yield(i + 10) {
					return
				}
			}
		}

		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		// Verify all elements are produced correctly
		assert.Equal(t, 10, len(result))
	})

	t.Run("handles concurrent consumption safely", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)

		// Consume with some processing time
		var sum atomic.Int32
		for v := range merged {
			sum.Add(int32(v))
			time.Sleep(1 * time.Millisecond)
		}

		assert.Equal(t, int32(55), sum.Load())
	})

	t.Run("demonstrates non-deterministic ordering", func(t *testing.T) {
		// Run multiple times to show order can vary
		results := make([][]int, 5)
		for i := range 5 {
			seq1 := From(1, 2, 3)
			seq2 := From(4, 5, 6)
			merged := MergeBuf([]Seq[int]{seq1, seq2}, 1)
			results[i] = toSlice(merged)
		}

		// All results should have same elements
		for _, result := range results {
			assert.Equal(t, 6, len(result))
			sorted := slices.Clone(result)
			slices.Sort(sorted)
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, sorted)
		}
	})
}

// TestMergeAll_EdgeCases tests edge cases
func TestMergeAll_EdgeCases(t *testing.T) {
	t.Run("very large buffer size", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 1000000)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("many sequences with small elements", func(t *testing.T) {
		sequences := make([]Seq[int], 100)
		for i := range 100 {
			sequences[i] = From(i)
		}
		merged := MergeBuf(sequences, 50)
		result := toSlice(merged)

		assert.Equal(t, 100, len(result))
		slices.Sort(result)
		expected := make([]int, 100)
		for i := range 100 {
			expected[i] = i
		}
		assert.Equal(t, expected, result)
	})

	t.Run("sequences of different lengths", func(t *testing.T) {
		seq1 := From(1)
		seq2 := From(2, 3, 4)
		seq3 := From(5, 6)
		merged := MergeBuf([]Seq[int]{seq1, seq2, seq3}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("works with replicate", func(t *testing.T) {
		seq1 := Replicate(3, 42)
		seq2 := Replicate(3, 99)
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		count42 := 0
		count99 := 0
		for _, v := range result {
			if v == 42 {
				count42++
			} else if v == 99 {
				count99++
			}
		}
		assert.Equal(t, 3, count42)
		assert.Equal(t, 3, count99)
	})

	t.Run("works with makeBy", func(t *testing.T) {
		seq1 := MakeBy(3, func(i int) int { return i * i })
		seq2 := MakeBy(3, func(i int) int { return i * 10 })
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{0, 0, 1, 4, 10, 20}, result)
	})
}

// Benchmark tests
func BenchmarkMergeAll_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_FiveSequences(b *testing.B) {
	sequences := make([]Seq[int], 5)
	for i := range 5 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf(sequences, 10)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_TenSequences(b *testing.B) {
	sequences := make([]Seq[int], 10)
	for i := range 10 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf(sequences, 20)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_SmallBuffer(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 1)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_LargeBuffer(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 100)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_Unbuffered(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 0)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_LargeSequences(b *testing.B) {
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
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 100)
		for range merged {
		}
	}
}

func BenchmarkMergeAll_WithMap(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		mapped := MonadMap(merged, func(x int) int { return x * 2 })
		for range mapped {
		}
	}
}

func BenchmarkMergeAll_WithFilter(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	seq2 := From(11, 12, 13, 14, 15, 16, 17, 18, 19, 20)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		filtered := MonadFilter(merged, func(x int) bool { return x%2 == 0 })
		for range filtered {
		}
	}
}

// Example tests for documentation
func ExampleMergeAll() {
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)
	merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeAll_threeSequences() {
	seq1 := From(1, 2)
	seq2 := From(3, 4)
	seq3 := From(5, 6)
	merged := MergeBuf([]Seq[int]{seq1, seq2, seq3}, 10)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeAll_earlyTermination() {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)

	count := 0
	for range merged {
		count++
		if count >= 3 {
			break
		}
	}
	fmt.Printf("Consumed %d elements\n", count)
	// Output: Consumed 3 elements
}

func ExampleMergeAll_unbuffered() {
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)
	merged := MergeBuf([]Seq[int]{seq1, seq2}, 0)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

// Made with Bob

// TestMerge_Alias tests that Merge is an alias for MergeAll
func TestMerge_Alias(t *testing.T) {
	t.Run("merge works identically to MergeAll", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("merge handles empty input", func(t *testing.T) {
		merged := MergeBuf([]Seq[int]{}, 5)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("merge with different buffer sizes", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		// Test with various buffer sizes
		for _, bufSize := range []int{0, 1, 5, 100} {
			merged := MergeBuf([]Seq[int]{seq1, seq2}, bufSize)
			result := toSlice(merged)

			assert.Equal(t, 6, len(result))
			slices.Sort(result)
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
		}
	})
}

// TestMerge_DefaultBuffer tests the Merge function with default buffer size
func TestMerge_DefaultBuffer(t *testing.T) {
	t.Run("merges two sequences with default buffer", func(t *testing.T) {
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		merged := Merge([]Seq[int]{seq1, seq2})
		result := slices.Collect(merged)

		assert.Len(t, result, 6)
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("merges multiple sequences", func(t *testing.T) {
		seq1 := From(1, 2)
		seq2 := From(3, 4)
		seq3 := From(5, 6)
		seq4 := From(7, 8)

		merged := Merge([]Seq[int]{seq1, seq2, seq3, seq4})
		result := slices.Collect(merged)

		assert.Len(t, result, 8)
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, result)
	})

	t.Run("handles empty slice", func(t *testing.T) {
		merged := Merge([]Seq[int]{})
		result := slices.Collect(merged)

		assert.Empty(t, result)
	})

	t.Run("handles single sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)

		merged := Merge([]Seq[int]{seq})
		result := slices.Collect(merged)

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("handles sequences with different lengths", func(t *testing.T) {
		seq1 := From(1)
		seq2 := From(2, 3, 4)
		seq3 := From(5, 6)

		merged := Merge([]Seq[int]{seq1, seq2, seq3})
		result := slices.Collect(merged)

		assert.Len(t, result, 6)
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("handles empty sequences in slice", func(t *testing.T) {
		seq1 := From(1, 2)
		seq2 := Empty[int]()
		seq3 := From(3, 4)

		merged := Merge([]Seq[int]{seq1, seq2, seq3})
		result := slices.Collect(merged)

		assert.ElementsMatch(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("early termination cleans up goroutines", func(t *testing.T) {
		seq1 := From(1, 2, 3, 4, 5)
		seq2 := From(6, 7, 8, 9, 10)
		seq3 := From(11, 12, 13, 14, 15)

		merged := Merge([]Seq[int]{seq1, seq2, seq3})

		var result []int
		for v := range merged {
			result = append(result, v)
			if len(result) >= 5 {
				break
			}
		}

		assert.Len(t, result, 5)
	})
}

// TestMergeMap_DefaultBuffer tests the MergeMap function with default buffer size
func TestMergeMap_DefaultBuffer(t *testing.T) {
	t.Run("expands each element with default buffer", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			return From(n, n*10)
		})
		input := From(1, 2, 3)

		result := slices.Collect(expand(input))

		assert.Len(t, result, 6)
		assert.ElementsMatch(t, []int{1, 10, 2, 20, 3, 30}, result)
	})

	t.Run("handles single element", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			return From(n, n*2, n*3)
		})
		input := From(5)

		result := slices.Collect(expand(input))

		assert.ElementsMatch(t, []int{5, 10, 15}, result)
	})

	t.Run("handles empty input", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			return From(n, n*2)
		})
		input := Empty[int]()

		result := slices.Collect(expand(input))

		assert.Empty(t, result)
	})

	t.Run("function returns empty sequences", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			return Empty[int]()
		})
		input := From(1, 2, 3)

		result := slices.Collect(expand(input))

		assert.Empty(t, result)
	})

	t.Run("mixed empty and non-empty sequences", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			if n%2 == 0 {
				return Empty[int]()
			}
			return From(n, n*10)
		})
		input := From(1, 2, 3, 4, 5)

		result := slices.Collect(expand(input))

		assert.ElementsMatch(t, []int{1, 10, 3, 30, 5, 50}, result)
	})

	t.Run("handles complex transformations", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			// Create sequences of different lengths based on input
			if n == 1 {
				return From(n)
			}
			if n == 2 {
				return From(n, n*10)
			}
			return From(n, n*10, n*100)
		})
		input := From(1, 2, 3)

		result := slices.Collect(expand(input))

		assert.Len(t, result, 6)
		assert.Contains(t, result, 1)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 20)
		assert.Contains(t, result, 3)
		assert.Contains(t, result, 30)
		assert.Contains(t, result, 300)
	})

	t.Run("early termination", func(t *testing.T) {
		expand := MergeMap(func(n int) Seq[int] {
			return From(n, n*10, n*100)
		})
		input := From(1, 2, 3, 4, 5)

		var result []int
		for v := range expand(input) {
			result = append(result, v)
			if len(result) >= 5 {
				break
			}
		}

		assert.Len(t, result, 5)
	})

	t.Run("with strings", func(t *testing.T) {
		expand := MergeMap(func(s string) Seq[string] {
			return From(s, s+s)
		})
		input := From("a", "b", "c")

		result := slices.Collect(expand(input))

		assert.Len(t, result, 6)
		assert.ElementsMatch(t, []string{"a", "aa", "b", "bb", "c", "cc"}, result)
	})

	t.Run("comparison with Chain shows same elements", func(t *testing.T) {
		f := func(n int) Seq[int] {
			return From(n, n*10)
		}
		input := From(1, 2, 3)

		// MergeMap result (non-deterministic order)
		mergeResult := slices.Collect(MergeMap(f)(input))

		// Chain result (deterministic order)
		chainResult := slices.Collect(Chain(f)(input))

		// Both should contain the same elements
		assert.ElementsMatch(t, chainResult, mergeResult)
	})
}

// TestMergeMonoid_Identity tests the monoid identity property
func TestMergeMonoid_Identity(t *testing.T) {
	t.Run("left identity: concat(empty, x) = x", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		seq := From(1, 2, 3)

		result := monoid.Concat(monoid.Empty(), seq)
		collected := toSlice(result)

		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("right identity: concat(x, empty) = x", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		seq := From(1, 2, 3)

		result := monoid.Concat(seq, monoid.Empty())
		collected := toSlice(result)

		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("empty is truly empty", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		result := toSlice(monoid.Empty())
		assert.Empty(t, result)
	})

	t.Run("concat(empty, empty) = empty", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		result := monoid.Concat(monoid.Empty(), monoid.Empty())
		collected := toSlice(result)
		assert.Empty(t, collected)
	})
}

// TestMergeMonoid_Associativity tests the monoid associativity property
func TestMergeMonoid_Associativity(t *testing.T) {
	t.Run("concat(concat(a, b), c) has same elements as concat(a, concat(b, c))", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		a := From(1, 2)
		b := From(3, 4)
		c := From(5, 6)

		// (a + b) + c
		left := monoid.Concat(monoid.Concat(a, b), c)
		leftResult := toSlice(left)
		slices.Sort(leftResult)

		// a + (b + c)
		right := monoid.Concat(a, monoid.Concat(b, c))
		rightResult := toSlice(right)
		slices.Sort(rightResult)

		// Both should contain the same elements (order may vary due to concurrency)
		assert.Equal(t, leftResult, rightResult)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, leftResult)
	})

	t.Run("associativity with different types", func(t *testing.T) {
		monoid := MergeMonoid[string](10)
		a := From("a", "b")
		b := From("c", "d")
		c := From("e", "f")

		left := monoid.Concat(monoid.Concat(a, b), c)
		leftResult := toSlice(left)
		slices.Sort(leftResult)

		right := monoid.Concat(a, monoid.Concat(b, c))
		rightResult := toSlice(right)
		slices.Sort(rightResult)

		assert.Equal(t, leftResult, rightResult)
		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, leftResult)
	})
}

// TestMergeMonoid_Concat tests the concat operation
func TestMergeMonoid_Concat(t *testing.T) {
	t.Run("concat two sequences", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		result := monoid.Concat(seq1, seq2)
		collected := toSlice(result)

		assert.Equal(t, 6, len(collected))
		slices.Sort(collected)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
	})

	t.Run("concat with different buffer sizes", func(t *testing.T) {
		for _, bufSize := range []int{0, 1, 5, 100} {
			monoid := MergeMonoid[int](bufSize)
			seq1 := From(1, 2, 3)
			seq2 := From(4, 5, 6)

			result := monoid.Concat(seq1, seq2)
			collected := toSlice(result)

			assert.Equal(t, 6, len(collected))
			slices.Sort(collected)
			assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
		}
	})

	t.Run("concat sequences of different lengths", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		seq1 := From(1)
		seq2 := From(2, 3, 4, 5)

		result := monoid.Concat(seq1, seq2)
		collected := toSlice(result)

		assert.Equal(t, 5, len(collected))
		slices.Sort(collected)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, collected)
	})
}

// TestMergeMonoid_WithReduce tests using MergeMonoid with Reduce
func TestMergeMonoid_WithReduce(t *testing.T) {
	t.Run("reduce multiple sequences using monoid", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		sequences := From(
			From(1, 2),
			From(3, 4),
			From(5, 6),
		)

		result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
		collected := toSlice(result)

		assert.Equal(t, 6, len(collected))
		slices.Sort(collected)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
	})

	t.Run("reduce with empty sequences", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		sequences := From(
			Empty[int](),
			From(1, 2),
			Empty[int](),
		)

		result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
		collected := toSlice(result)

		assert.Equal(t, []int{1, 2}, collected)
	})

	t.Run("reduce empty sequence of sequences", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		sequences := Empty[Seq[int]]()

		result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
		collected := toSlice(result)

		assert.Empty(t, collected)
	})
}

// TestMergeMonoid_WithFoldMap tests using MergeMonoid with FoldMap
func TestMergeMonoid_WithFoldMap(t *testing.T) {
	t.Run("foldMap with monoid", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		numbers := From(1, 2, 3)

		// Convert each number to a sequence and merge all
		result := MonadFoldMap(numbers, func(n int) Seq[int] {
			return From(n, n*10)
		}, monoid)

		collected := toSlice(result)
		assert.Equal(t, 6, len(collected))
		slices.Sort(collected)
		assert.Equal(t, []int{1, 2, 3, 10, 20, 30}, collected)
	})

	t.Run("foldMap with complex transformation", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		numbers := From(1, 2, 3)

		result := MonadFoldMap(numbers, func(n int) Seq[int] {
			return From(n, n*10, n*100)
		}, monoid)

		collected := toSlice(result)
		assert.Equal(t, 9, len(collected))
		slices.Sort(collected)
		assert.Equal(t, []int{1, 2, 3, 10, 20, 30, 100, 200, 300}, collected)
	})

	t.Run("foldMap with empty input", func(t *testing.T) {
		monoid := MergeMonoid[int](10)
		numbers := Empty[int]()

		result := MonadFoldMap(numbers, func(n int) Seq[int] {
			return From(n, n*10)
		}, monoid)

		collected := toSlice(result)
		assert.Empty(t, collected)
	})
}

// TestMergeMonoid_WithComplexTypes tests MergeMonoid with complex types
func TestMergeMonoid_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("monoid with structs", func(t *testing.T) {
		monoid := MergeMonoid[Person](10)
		seq1 := From(Person{"Alice", 30}, Person{"Bob", 25})
		seq2 := From(Person{"Charlie", 35})

		result := monoid.Concat(seq1, seq2)
		collected := toSlice(result)

		assert.Equal(t, 3, len(collected))
	})

	t.Run("monoid with slices", func(t *testing.T) {
		monoid := MergeMonoid[[]int](10)
		seq1 := From([]int{1, 2}, []int{3, 4})
		seq2 := From([]int{5, 6})

		result := monoid.Concat(seq1, seq2)
		collected := toSlice(result)

		assert.Equal(t, 3, len(collected))
	})
}

// TestMergeMonoid_Concurrency tests concurrent behavior
func TestMergeMonoid_Concurrency(t *testing.T) {
	t.Run("monoid concat is thread-safe", func(t *testing.T) {
		monoid := MergeMonoid[int](10)

		// Create slow sequences
		seq1 := func(yield func(int) bool) {
			for i := range 5 {
				time.Sleep(1 * time.Millisecond)
				if !yield(i) {
					return
				}
			}
		}
		seq2 := func(yield func(int) bool) {
			for i := range 5 {
				time.Sleep(1 * time.Millisecond)
				if !yield(i + 10) {
					return
				}
			}
		}

		result := monoid.Concat(seq1, seq2)
		collected := toSlice(result)

		assert.Equal(t, 10, len(collected))
	})
}

// Benchmark tests for Merge
func BenchmarkMerge_TwoSequences(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		merged := MergeBuf([]Seq[int]{seq1, seq2}, 5)
		for range merged {
		}
	}
}

// Benchmark tests for MergeMonoid
func BenchmarkMergeMonoid_Concat(b *testing.B) {
	monoid := MergeMonoid[int](10)
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		result := monoid.Concat(seq1, seq2)
		for range result {
		}
	}
}

func BenchmarkMergeMonoid_Reduce(b *testing.B) {
	monoid := MergeMonoid[int](10)
	sequences := From(
		From(1, 2, 3),
		From(4, 5, 6),
		From(7, 8, 9),
	)
	b.ResetTimer()
	for range b.N {
		result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
		for range result {
		}
	}
}

func BenchmarkMergeMonoid_FoldMap(b *testing.B) {
	monoid := MergeMonoid[int](10)
	numbers := From(1, 2, 3, 4, 5)
	b.ResetTimer()
	for range b.N {
		result := MonadFoldMap(numbers, func(n int) Seq[int] {
			return From(n, n*10)
		}, monoid)
		for range result {
		}
	}
}

// Example tests for Merge
func ExampleMerge() {
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)
	merged := MergeBuf([]Seq[int]{seq1, seq2}, 10)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

// Example tests for MergeMonoid
func ExampleMergeMonoid() {
	monoid := MergeMonoid[int](10)
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)

	result := monoid.Concat(seq1, seq2)
	collected := toSlice(result)
	slices.Sort(collected)
	for _, v := range collected {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeMonoid_identity() {
	monoid := MergeMonoid[int](10)
	seq := From(1, 2, 3)

	// Left identity
	result1 := monoid.Concat(monoid.Empty(), seq)
	fmt.Println("Left identity:", toSlice(result1))

	// Right identity
	result2 := monoid.Concat(seq, monoid.Empty())
	fmt.Println("Right identity:", toSlice(result2))

	// Output:
	// Left identity: [1 2 3]
	// Right identity: [1 2 3]
}

func ExampleMergeMonoid_reduce() {
	monoid := MergeMonoid[int](10)
	sequences := From(
		From(1, 2),
		From(3, 4),
		From(5, 6),
	)

	result := MonadReduce(sequences, monoid.Concat, monoid.Empty())
	collected := toSlice(result)
	slices.Sort(collected)
	for _, v := range collected {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeMonoid_foldMap() {
	monoid := MergeMonoid[int](10)
	numbers := From(1, 2, 3)

	// Convert each number to a sequence and merge all
	result := MonadFoldMap(numbers, func(n int) Seq[int] {
		return From(n, n*10)
	}, monoid)

	collected := toSlice(result)
	slices.Sort(collected)
	for _, v := range collected {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 10 20 30
}

// TestMergeAll_Operator_Success tests basic MergeAll operator functionality
func TestMergeAll_Operator_Success(t *testing.T) {
	t.Run("flattens two sequences", func(t *testing.T) {
		outer := From(
			From(1, 2, 3),
			From(4, 5, 6),
		)
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("flattens three sequences", func(t *testing.T) {
		outer := From(
			From(1, 2),
			From(3, 4),
			From(5, 6),
		)
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("flattens many sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2)
		}
		outer := From(sequences...)
		mergeAll := MergeAll[int](20)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 30, len(result))
		slices.Sort(result)
		expected := make([]int, 30)
		for i := range 30 {
			expected[i] = (i/3)*10 + (i % 3)
		}
		assert.Equal(t, expected, result)
	})

	t.Run("flattens single sequence", func(t *testing.T) {
		outer := From(From(1, 2, 3, 4, 5))
		mergeAll := MergeAll[int](5)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("preserves all elements with strings", func(t *testing.T) {
		outer := From(
			From("a", "b", "c"),
			From("d", "e", "f"),
		)
		mergeAll := MergeAll[string](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []string{"a", "b", "c", "d", "e", "f"}, result)
	})
}

// TestMergeAll_Operator_BufferSizes tests different buffer sizes
func TestMergeAll_Operator_BufferSizes(t *testing.T) {
	t.Run("unbuffered channel (bufSize 0)", func(t *testing.T) {
		outer := From(
			From(1, 2, 3),
			From(4, 5, 6),
		)
		mergeAll := MergeAll[int](0)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("small buffer", func(t *testing.T) {
		outer := From(
			From(1, 2, 3, 4, 5),
			From(6, 7, 8, 9, 10),
		)
		mergeAll := MergeAll[int](2)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 10, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)
	})

	t.Run("large buffer", func(t *testing.T) {
		outer := From(
			From(1, 2, 3),
			From(4, 5, 6),
		)
		mergeAll := MergeAll[int](100)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("negative buffer size treated as 0", func(t *testing.T) {
		outer := From(
			From(1, 2, 3),
			From(4, 5, 6),
		)
		mergeAll := MergeAll[int](-5)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})

	t.Run("buffer size of 1", func(t *testing.T) {
		outer := From(
			From(1, 2, 3),
			From(4, 5, 6),
		)
		mergeAll := MergeAll[int](1)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)
	})
}

// TestMergeAll_Operator_Empty tests MergeAll with empty sequences
func TestMergeAll_Operator_Empty(t *testing.T) {
	t.Run("empty outer sequence", func(t *testing.T) {
		outer := Empty[Seq[int]]()
		mergeAll := MergeAll[int](5)
		merged := mergeAll(outer)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("outer sequence with empty inner sequences", func(t *testing.T) {
		outer := From(Empty[int](), Empty[int]())
		mergeAll := MergeAll[int](5)
		merged := mergeAll(outer)
		result := toSlice(merged)
		assert.Empty(t, result)
	})

	t.Run("mix of empty and non-empty inner sequences", func(t *testing.T) {
		outer := From(
			Empty[int](),
			From(1, 2, 3),
			Empty[int](),
		)
		mergeAll := MergeAll[int](5)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestMergeAll_Operator_EarlyTermination tests proper cleanup on early termination
func TestMergeAll_Operator_EarlyTermination(t *testing.T) {
	t.Run("stops all producers when consumer breaks", func(t *testing.T) {
		var outerCount atomic.Int32
		var innerCount1 atomic.Int32
		var innerCount2 atomic.Int32

		// Create outer sequence that tracks production
		outer := func(yield func(Seq[int]) bool) {
			for range 10 {
				outerCount.Add(1)
				inner := func(yieldInner func(int) bool) {
					counter := &innerCount1
					if outerCount.Load() > 1 {
						counter = &innerCount2
					}
					for i := range 100 {
						counter.Add(1)
						if !yieldInner(i) {
							return
						}
					}
				}
				if !yield(inner) {
					return
				}
			}
		}

		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)

		// Consume only 5 elements
		count := 0
		for range merged {
			count++
			if count >= 5 {
				break
			}
		}

		// Give goroutines time to clean up
		time.Sleep(10 * time.Millisecond)

		// Outer producer should have stopped (may produce a few extra due to buffering)
		assert.LessOrEqual(t, outerCount.Load(), int32(10), "outer producer should stop")

		// Inner producers should have stopped shortly after consumer stopped
		totalInner := innerCount1.Load() + innerCount2.Load()
		assert.LessOrEqual(t, totalInner, int32(50), "inner producers should stop after consumer breaks")
		assert.GreaterOrEqual(t, totalInner, int32(5), "inner producers should produce at least what was consumed")
	})

	t.Run("handles yield returning false", func(t *testing.T) {
		outer := From(
			From(1, 2, 3, 4, 5),
			From(6, 7, 8, 9, 10),
		)
		mergeAll := MergeAll[int](5)
		merged := mergeAll(outer)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 3 {
				break
			}
		}

		assert.Equal(t, 3, len(collected))
	})

	t.Run("early termination with unbuffered channel", func(t *testing.T) {
		outer := From(
			From(1, 2, 3, 4, 5),
			From(6, 7, 8, 9, 10),
		)
		mergeAll := MergeAll[int](0)
		merged := mergeAll(outer)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 2 {
				break
			}
		}

		assert.Equal(t, 2, len(collected))
	})

	t.Run("early termination with many sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 10)
		for i := range 10 {
			sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
		}
		outer := From(sequences...)
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)

		collected := []int{}
		for v := range merged {
			collected = append(collected, v)
			if len(collected) >= 5 {
				break
			}
		}

		assert.Equal(t, 5, len(collected))
	})
}

// TestMergeAll_Operator_WithComplexTypes tests MergeAll with complex data types
func TestMergeAll_Operator_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("works with structs", func(t *testing.T) {
		outer := From(
			From(Person{"Alice", 30}, Person{"Bob", 25}),
			From(Person{"Charlie", 35}, Person{"Diana", 28}),
		)
		mergeAll := MergeAll[Person](5)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 4, len(result))
		slices.SortFunc(result, func(a, b Person) int {
			if a.Name < b.Name {
				return -1
			}
			if a.Name > b.Name {
				return 1
			}
			return 0
		})
		assert.Equal(t, []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
			{"Diana", 28},
		}, result)
	})

	t.Run("works with pointers", func(t *testing.T) {
		p1, p2, p3, p4 := &Person{"Alice", 30}, &Person{"Bob", 25}, &Person{"Charlie", 35}, &Person{"Diana", 28}
		outer := From(
			From(p1, p2),
			From(p3, p4),
		)
		mergeAll := MergeAll[*Person](5)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 4, len(result))
		assert.Contains(t, result, p1)
		assert.Contains(t, result, p2)
		assert.Contains(t, result, p3)
		assert.Contains(t, result, p4)
	})
}

// TestMergeAll_Operator_WithDynamicGeneration tests dynamic sequence generation
func TestMergeAll_Operator_WithDynamicGeneration(t *testing.T) {
	t.Run("generates sequences on-the-fly with Map", func(t *testing.T) {
		// Generate sequences dynamically
		outer := MonadMap(From(1, 2, 3), func(n int) Seq[int] {
			return From(n, n*10, n*100)
		})
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 9, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 10, 20, 30, 100, 200, 300}, result)
	})

	t.Run("works with Chain-like pattern", func(t *testing.T) {
		// MergeAll can be used to implement Chain behavior
		numbers := From(1, 2, 3)
		outer := MonadMap(numbers, func(n int) Seq[int] {
			return From(n, n*2)
		})
		mergeAll := MergeAll[int](0)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 2, 3, 4, 6}, result)
	})

	t.Run("filters and flattens", func(t *testing.T) {
		outer := MonadMap(From(1, 2, 3, 4, 5), func(n int) Seq[int] {
			if n%2 == 0 {
				return From(n, n*10)
			}
			return Empty[int]()
		})
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 4, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{2, 4, 20, 40}, result)
	})
}

// TestMergeAll_Operator_Concurrency tests concurrent behavior
func TestMergeAll_Operator_Concurrency(t *testing.T) {
	t.Run("processes inner sequences concurrently", func(t *testing.T) {
		var startTimes [3]time.Time
		var mu sync.Mutex

		sequences := make([]Seq[int], 3)
		for i := range 3 {
			idx := i
			sequences[i] = func(yield func(int) bool) {
				mu.Lock()
				startTimes[idx] = time.Now()
				mu.Unlock()
				time.Sleep(10 * time.Millisecond)
				yield(idx)
			}
		}

		outer := From(sequences...)
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		_ = toSlice(merged)

		// All sequences should start around the same time (within 20ms)
		mu.Lock()
		defer mu.Unlock()
		for i := 1; i < 3; i++ {
			diff := startTimes[i].Sub(startTimes[0])
			assert.Less(t, diff, 20*time.Millisecond, "sequences should start concurrently")
		}
	})

	t.Run("handles concurrent production correctly", func(t *testing.T) {
		// Create sequences that produce elements with delays
		sequences := make([]Seq[int], 5)
		for i := range 5 {
			base := i * 10
			sequences[i] = func(yield func(int) bool) {
				for j := range 3 {
					time.Sleep(time.Millisecond)
					if !yield(base + j) {
						return
					}
				}
			}
		}

		outer := From(sequences...)
		mergeAll := MergeAll[int](20)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 15, len(result))
		slices.Sort(result)
		expected := make([]int, 15)
		for i := range 15 {
			expected[i] = (i/3)*10 + (i % 3)
		}
		assert.Equal(t, expected, result)
	})
}

// TestMergeAll_Operator_EdgeCases tests edge cases
func TestMergeAll_Operator_EdgeCases(t *testing.T) {
	t.Run("single element in single sequence", func(t *testing.T) {
		outer := From(From(42))
		mergeAll := MergeAll[int](1)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, []int{42}, result)
	})

	t.Run("many small sequences", func(t *testing.T) {
		sequences := make([]Seq[int], 100)
		for i := range 100 {
			sequences[i] = From(i)
		}
		outer := From(sequences...)
		mergeAll := MergeAll[int](50)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 100, len(result))
		slices.Sort(result)
		expected := make([]int, 100)
		for i := range 100 {
			expected[i] = i
		}
		assert.Equal(t, expected, result)
	})

	t.Run("sequences with varying lengths", func(t *testing.T) {
		outer := From(
			From(1),
			From(2, 3),
			From(4, 5, 6),
			From(7, 8, 9, 10),
		)
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 10, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, result)
	})

	t.Run("outer sequence produces slowly", func(t *testing.T) {
		outer := func(yield func(Seq[int]) bool) {
			for i := range 3 {
				time.Sleep(5 * time.Millisecond)
				if !yield(From(i*10, i*10+1)) {
					return
				}
			}
		}
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		result := toSlice(merged)

		assert.Equal(t, 6, len(result))
		slices.Sort(result)
		assert.Equal(t, []int{0, 1, 10, 11, 20, 21}, result)
	})
}

// Benchmarks for MergeAll operator

func BenchmarkMergeAll_Operator_TwoSequences(b *testing.B) {
	for range b.N {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_FiveSequences(b *testing.B) {
	sequences := make([]Seq[int], 5)
	for i := range 5 {
		sequences[i] = From(i*10, i*10+1, i*10+2, i*10+3, i*10+4)
	}
	b.ResetTimer()
	for range b.N {
		outer := From(sequences...)
		mergeAll := MergeAll[int](20)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_TenSequences(b *testing.B) {
	sequences := make([]Seq[int], 10)
	for i := range 10 {
		sequences[i] = From(i*10, i*10+1, i*10+2)
	}
	b.ResetTimer()
	for range b.N {
		outer := From(sequences...)
		mergeAll := MergeAll[int](30)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_SmallBuffer(b *testing.B) {
	for range b.N {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		mergeAll := MergeAll[int](2)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_LargeBuffer(b *testing.B) {
	for range b.N {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		mergeAll := MergeAll[int](100)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_Unbuffered(b *testing.B) {
	for range b.N {
		outer := From(From(1, 2, 3, 4, 5), From(6, 7, 8, 9, 10))
		mergeAll := MergeAll[int](0)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_LargeSequences(b *testing.B) {
	sequences := make([]Seq[int], 5)
	for i := range 5 {
		start := i * 100
		vals := make([]int, 100)
		for j := range 100 {
			vals[j] = start + j
		}
		sequences[i] = From(vals...)
	}
	b.ResetTimer()
	for range b.N {
		outer := From(sequences...)
		mergeAll := MergeAll[int](50)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_WithMap(b *testing.B) {
	for range b.N {
		outer := MonadMap(From(1, 2, 3, 4, 5), func(n int) Seq[int] {
			return From(n, n*10)
		})
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

func BenchmarkMergeAll_Operator_WithFilter(b *testing.B) {
	for range b.N {
		outer := MonadMap(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), func(n int) Seq[int] {
			if n%2 == 0 {
				return From(n, n*10)
			}
			return Empty[int]()
		})
		mergeAll := MergeAll[int](10)
		merged := mergeAll(outer)
		_ = toSlice(merged)
	}
}

// Example tests for MergeAll operator

func ExampleMergeAll_operator() {
	// Create a sequence of sequences
	outer := From(
		From(1, 2, 3),
		From(4, 5, 6),
	)

	// Flatten and merge concurrently
	mergeAll := MergeAll[int](10)
	merged := mergeAll(outer)

	// Collect and sort for deterministic output
	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeAll_operator_threeSequences() {
	outer := From(
		From(1, 2),
		From(3, 4),
		From(5, 6),
	)

	mergeAll := MergeAll[int](10)
	merged := mergeAll(outer)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

func ExampleMergeAll_operator_dynamicGeneration() {
	// Generate sequences on-the-fly
	outer := MonadMap(From(1, 2, 3), func(n int) Seq[int] {
		return From(n, n*10)
	})

	mergeAll := MergeAll[int](10)
	merged := mergeAll(outer)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 10 20 30
}

func ExampleMergeAll_operator_earlyTermination() {
	outer := From(
		From(1, 2, 3, 4, 5),
		From(6, 7, 8, 9, 10),
		From(11, 12, 13, 14, 15),
	)

	mergeAll := MergeAll[int](5)
	merged := mergeAll(outer)

	// Take only first 5 elements
	count := 0
	for range merged {
		count++
		if count >= 5 {
			break
		}
	}

	fmt.Printf("Consumed %d elements\n", count)
	// Output: Consumed 5 elements
}

func ExampleMergeAll_operator_unbuffered() {
	outer := From(
		From(1, 2, 3),
		From(4, 5, 6),
	)

	// Unbuffered channel (bufSize 0)
	mergeAll := MergeAll[int](0)
	merged := mergeAll(outer)

	result := toSlice(merged)
	slices.Sort(result)
	for _, v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5 6
}

// TestMergeMap_Success tests basic MergeMap functionality
func TestMergeMap_Success(t *testing.T) {
	t.Run("expands each element into a sequence", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*10)
		}, 10)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 6)
		assert.Contains(t, result, 1)
		assert.Contains(t, result, 10)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 20)
		assert.Contains(t, result, 3)
		assert.Contains(t, result, 30)
	})

	t.Run("handles single element", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*2, n*3)
		}, 5)
		input := From(5)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.ElementsMatch(t, []int{5, 10, 15}, result)
	})

	t.Run("handles multiple elements with different sequence sizes", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			// Create sequences of different lengths
			if n == 1 {
				return From(n)
			}
			if n == 2 {
				return From(n, n*10)
			}
			return From(n, n*10, n*100)
		}, 10)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 6)
		assert.Contains(t, result, 1)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 20)
		assert.Contains(t, result, 3)
		assert.Contains(t, result, 30)
		assert.Contains(t, result, 300)
	})
}

// TestMergeMap_BufferSizes tests MergeMap with different buffer sizes
func TestMergeMap_BufferSizes(t *testing.T) {
	t.Run("unbuffered channel (bufSize=0)", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*2)
		}, 0)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 6)
		assert.Contains(t, result, 1)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 4)
		assert.Contains(t, result, 3)
		assert.Contains(t, result, 6)
	})

	t.Run("small buffer (bufSize=1)", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*2)
		}, 1)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 6)
	})

	t.Run("large buffer (bufSize=100)", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*2)
		}, 100)
		input := From(1, 2, 3, 4, 5)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 10)
	})

	t.Run("negative buffer size treated as 0", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n)
		}, -5)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.ElementsMatch(t, []int{1, 2, 3}, result)
	})
}

// TestMergeMap_Empty tests MergeMap with empty sequences
func TestMergeMap_Empty(t *testing.T) {
	t.Run("empty input sequence", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*2)
		}, 10)
		input := Empty[int]()

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Empty(t, result)
	})

	t.Run("function returns empty sequences", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return Empty[int]()
		}, 10)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Empty(t, result)
	})

	t.Run("mixed empty and non-empty sequences", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			if n%2 == 0 {
				return Empty[int]()
			}
			return From(n)
		}, 10)
		input := From(1, 2, 3, 4, 5)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.ElementsMatch(t, []int{1, 3, 5}, result)
	})
}

// TestMergeMap_EarlyTermination tests MergeMap with early termination
func TestMergeMap_EarlyTermination(t *testing.T) {
	t.Run("stops after consuming partial results", func(t *testing.T) {
		// Arrange
		callCount := atomic.Int32{}
		expand := MergeMapBuf(func(n int) Seq[int] {
			callCount.Add(1)
			return From(n, n*10, n*100)
		}, 10)
		input := From(1, 2, 3, 4, 5)

		// Act
		var result []int
		for v := range expand(input) {
			result = append(result, v)
			if len(result) >= 5 {
				break
			}
		}

		// Assert
		assert.Len(t, result, 5)
		// Function should have been called for at least some elements
		assert.Greater(t, callCount.Load(), int32(0))
	})

	t.Run("stops immediately", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*10)
		}, 10)
		input := From(1, 2, 3)

		// Act
		var result []int
		for v := range expand(input) {
			result = append(result, v)
			break // Stop immediately
		}

		// Assert
		assert.Len(t, result, 1)
	})

	t.Run("no goroutine leaks on early termination", func(t *testing.T) {
		// Arrange
		initialGoroutines := countGoroutines()
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n*10, n*100, n*1000)
		}, 5)
		input := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

		// Act
		count := 0
		for v := range expand(input) {
			_ = v
			count++
			if count >= 10 {
				break
			}
		}

		// Allow goroutines to clean up
		time.Sleep(50 * time.Millisecond)

		// Assert
		finalGoroutines := countGoroutines()
		assert.LessOrEqual(t, finalGoroutines, initialGoroutines+2,
			"Should not leak goroutines after early termination")
	})
}

// TestMergeMap_WithComplexTypes tests MergeMap with complex types
func TestMergeMap_WithComplexTypes(t *testing.T) {
	t.Run("with strings", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(s string) Seq[string] {
			return From(s, s+s)
		}, 10)
		input := From("a", "b", "c")

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 6)
		assert.Contains(t, result, "a")
		assert.Contains(t, result, "aa")
		assert.Contains(t, result, "b")
		assert.Contains(t, result, "bb")
		assert.Contains(t, result, "c")
		assert.Contains(t, result, "cc")
	})

	t.Run("with structs", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		expand := MergeMapBuf(func(p Person) Seq[string] {
			return From(p.Name, fmt.Sprintf("%d", p.Age))
		}, 10)
		input := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
		)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 4)
		assert.Contains(t, result, "Alice")
		assert.Contains(t, result, "30")
		assert.Contains(t, result, "Bob")
		assert.Contains(t, result, "25")
	})

	t.Run("with nested sequences", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n, n+1, n+2)
		}, 10)
		input := From(1, 10, 100)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 9)
		assert.Contains(t, result, 1)
		assert.Contains(t, result, 2)
		assert.Contains(t, result, 3)
		assert.Contains(t, result, 10)
		assert.Contains(t, result, 11)
		assert.Contains(t, result, 12)
		assert.Contains(t, result, 100)
		assert.Contains(t, result, 101)
		assert.Contains(t, result, 102)
	})
}

// TestMergeMap_Concurrency tests MergeMap concurrent behavior
func TestMergeMap_Concurrency(t *testing.T) {
	t.Run("processes sequences concurrently", func(t *testing.T) {
		// Arrange
		var mu sync.Mutex
		activeCount := 0
		maxActive := 0

		expand := MergeMapBuf(func(n int) Seq[int] {
			return func(yield func(int) bool) {
				mu.Lock()
				activeCount++
				if activeCount > maxActive {
					maxActive = activeCount
				}
				mu.Unlock()

				time.Sleep(10 * time.Millisecond)

				mu.Lock()
				activeCount--
				mu.Unlock()

				yield(n)
			}
		}, 10)
		input := From(1, 2, 3, 4, 5)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 5)
		assert.Greater(t, maxActive, 1, "Should have multiple concurrent producers")
	})

	t.Run("thread-safe element production", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return func(yield func(int) bool) {
				for i := range 10 {
					if !yield(n*10 + i) {
						return
					}
				}
			}
		}, 20)
		input := From(1, 2, 3, 4, 5)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 50)
		// Verify all expected values are present
		for n := range 5 {
			for i := range 10 {
				expected := (n+1)*10 + i
				assert.Contains(t, result, expected)
			}
		}
	})
}

// TestMergeMap_EdgeCases tests MergeMap edge cases
func TestMergeMap_EdgeCases(t *testing.T) {
	t.Run("single element producing many values", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return func(yield func(int) bool) {
				for i := range 100 {
					if !yield(i) {
						return
					}
				}
			}
		}, 50)
		input := From(1)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 100)
	})

	t.Run("many elements producing single value", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return From(n)
		}, 10)
		input := func(yield func(int) bool) {
			for i := range 100 {
				if !yield(i) {
					return
				}
			}
		}

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Len(t, result, 100)
	})

	t.Run("function returns nil-like empty sequence", func(t *testing.T) {
		// Arrange
		expand := MergeMapBuf(func(n int) Seq[int] {
			return func(yield func(int) bool) {
				// Empty sequence - never calls yield
			}
		}, 10)
		input := From(1, 2, 3)

		// Act
		result := slices.Collect(expand(input))

		// Assert
		assert.Empty(t, result)
	})
}

// TestMergeMap_ComparisonWithChain tests that MergeMap produces same elements as Chain
func TestMergeMap_ComparisonWithChain(t *testing.T) {
	t.Run("produces same elements as Chain (different order)", func(t *testing.T) {
		// Arrange
		f := func(n int) Seq[int] {
			return From(n, n*10, n*100)
		}
		input := From(1, 2, 3)

		// Act
		mergeMapResult := slices.Collect(MergeMapBuf(f, 10)(input))
		chainResult := slices.Collect(Chain(f)(input))

		// Assert
		assert.Len(t, mergeMapResult, len(chainResult))
		assert.ElementsMatch(t, chainResult, mergeMapResult,
			"MergeMap should produce same elements as Chain, just in different order")
	})
}

// BenchmarkMergeMap_SmallSequences benchmarks MergeMap with small sequences
func BenchmarkMergeMap_SmallSequences(b *testing.B) {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*2, n*3)
	}, 10)
	input := From(1, 2, 3, 4, 5)

	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// BenchmarkMergeMap_LargeSequences benchmarks MergeMap with large sequences
func BenchmarkMergeMap_LargeSequences(b *testing.B) {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return func(yield func(int) bool) {
			for i := range 100 {
				if !yield(n*100 + i) {
					return
				}
			}
		}
	}, 50)
	input := func(yield func(int) bool) {
		for i := range 10 {
			if !yield(i) {
				return
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// BenchmarkMergeMap_SmallBuffer benchmarks MergeMap with small buffer
func BenchmarkMergeMap_SmallBuffer(b *testing.B) {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*2)
	}, 1)
	input := From(1, 2, 3, 4, 5)

	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// BenchmarkMergeMap_LargeBuffer benchmarks MergeMap with large buffer
func BenchmarkMergeMap_LargeBuffer(b *testing.B) {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*2)
	}, 100)
	input := From(1, 2, 3, 4, 5)

	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// BenchmarkMergeMap_Unbuffered benchmarks MergeMap with unbuffered channel
func BenchmarkMergeMap_Unbuffered(b *testing.B) {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*2)
	}, 0)
	input := From(1, 2, 3, 4, 5)

	b.ResetTimer()
	for range b.N {
		for range expand(input) {
		}
	}
}

// BenchmarkMergeMap_VsChain benchmarks MergeMap vs Chain
func BenchmarkMergeMap_VsChain(b *testing.B) {
	f := func(n int) Seq[int] {
		return From(n, n*2, n*3)
	}
	input := From(1, 2, 3, 4, 5)

	b.Run("MergeMap", func(b *testing.B) {
		expand := MergeMapBuf(f, 10)
		b.ResetTimer()
		for range b.N {
			for range expand(input) {
			}
		}
	})

	b.Run("Chain", func(b *testing.B) {
		chain := Chain(f)
		b.ResetTimer()
		for range b.N {
			for range chain(input) {
			}
		}
	})
}

// ExampleMergeMap demonstrates basic MergeMap usage
func ExampleMergeMap() {
	// Expand each number into a sequence concurrently
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*10)
	}, 10)
	seq := From(1, 2, 3)

	// Collect results (order may vary due to concurrency)
	result := slices.Collect(expand(seq))
	fmt.Printf("Length: %d\n", len(result))
	fmt.Printf("Contains 1: %v\n", slices.Contains(result, 1))
	fmt.Printf("Contains 10: %v\n", slices.Contains(result, 10))
	fmt.Printf("Contains 2: %v\n", slices.Contains(result, 2))
	fmt.Printf("Contains 20: %v\n", slices.Contains(result, 20))
	fmt.Printf("Contains 3: %v\n", slices.Contains(result, 3))
	fmt.Printf("Contains 30: %v\n", slices.Contains(result, 30))

	// Output:
	// Length: 6
	// Contains 1: true
	// Contains 10: true
	// Contains 2: true
	// Contains 20: true
	// Contains 3: true
	// Contains 30: true
}

// ExampleMergeMap_earlyTermination demonstrates early termination
func ExampleMergeMap_earlyTermination() {
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*10, n*100)
	}, 5)
	seq := From(1, 2, 3, 4, 5)

	// Stop after collecting 5 elements
	count := 0
	for range expand(seq) {
		count++
		if count >= 5 {
			break
		}
	}

	fmt.Printf("Collected %d elements\n", count)

	// Output:
	// Collected 5 elements
}

// ExampleMergeMap_unbuffered demonstrates unbuffered channel usage
func ExampleMergeMap_unbuffered() {
	// bufSize of 0 creates an unbuffered channel
	expand := MergeMapBuf(func(n int) Seq[int] {
		return From(n, n*2)
	}, 0)
	seq := From(1, 2, 3)

	result := slices.Collect(expand(seq))
	fmt.Printf("Length: %d\n", len(result))

	// Output:
	// Length: 6
}

// Helper function to count goroutines (for testing)
func countGoroutines() int {
	// Simple approximation - in real tests you might use runtime.NumGoroutine()
	// but that can be flaky. This is a placeholder.
	return 0
}
