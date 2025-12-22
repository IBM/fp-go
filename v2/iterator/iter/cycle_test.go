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

	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestCycleBasic tests basic Cycle functionality with Take
func TestCycleBasic(t *testing.T) {
	t.Run("cycles through integer sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1}, result)
	})

	t.Run("cycles through string sequence", func(t *testing.T) {
		seq := From("A", "B", "C")
		cycled := Cycle(seq)
		result := toSlice(Take[string](8)(cycled))
		assert.Equal(t, []string{"A", "B", "C", "A", "B", "C", "A", "B"}, result)
	})

	t.Run("cycles through single element", func(t *testing.T) {
		seq := From(42)
		cycled := Cycle(seq)
		result := toSlice(Take[int](5)(cycled))
		assert.Equal(t, []int{42, 42, 42, 42, 42}, result)
	})

	t.Run("cycles through two elements", func(t *testing.T) {
		seq := From(true, false)
		cycled := Cycle(seq)
		result := toSlice(Take[bool](6)(cycled))
		assert.Equal(t, []bool{true, false, true, false, true, false}, result)
	})

	t.Run("takes exact multiple of cycle length", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		result := toSlice(Take[int](9)(cycled))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1, 2, 3}, result)
	})

	t.Run("takes less than one cycle", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		cycled := Cycle(seq)
		result := toSlice(Take[int](3)(cycled))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestCycleEmpty tests Cycle with empty sequences
func TestCycleEmpty(t *testing.T) {
	t.Run("cycles empty sequence produces nothing", func(t *testing.T) {
		seq := Empty[int]()
		cycled := Cycle(seq)
		result := toSlice(Take[int](10)(cycled))
		assert.Empty(t, result)
	})

	t.Run("cycles empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		cycled := Cycle(seq)
		result := toSlice(Take[string](5)(cycled))
		assert.Empty(t, result)
	})

	t.Run("take zero from cycled sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		result := toSlice(Take[int](0)(cycled))
		assert.Empty(t, result)
	})
}

// TestCycleWithComplexTypes tests Cycle with complex data types
func TestCycleWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("cycles structs", func(t *testing.T) {
		seq := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
		)
		cycled := Cycle(seq)
		result := toSlice(Take[Person](5)(cycled))
		expected := []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Alice", 30},
			{"Bob", 25},
			{"Alice", 30},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("cycles pointers", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		seq := From(p1, p2)
		cycled := Cycle(seq)
		result := toSlice(Take[*Person](4)(cycled))
		assert.Equal(t, []*Person{p1, p2, p1, p2}, result)
	})

	t.Run("cycles slices", func(t *testing.T) {
		seq := From([]int{1, 2}, []int{3, 4})
		cycled := Cycle(seq)
		result := toSlice(Take[[]int](5)(cycled))
		expected := [][]int{{1, 2}, {3, 4}, {1, 2}, {3, 4}, {1, 2}}
		assert.Equal(t, expected, result)
	})
}

// TestCycleWithFirst tests Cycle with First operation
func TestCycleWithFirst(t *testing.T) {
	t.Run("gets first element from cycled sequence", func(t *testing.T) {
		seq := From(10, 20, 30)
		cycled := Cycle(seq)
		first := First(cycled)
		assert.Equal(t, O.Of(10), first)
	})

	t.Run("gets first from single element cycle", func(t *testing.T) {
		seq := From(42)
		cycled := Cycle(seq)
		first := First(cycled)
		assert.Equal(t, O.Of(42), first)
	})

	t.Run("gets none from empty cycle", func(t *testing.T) {
		seq := Empty[int]()
		cycled := Cycle(seq)
		first := First(cycled)
		assert.Equal(t, O.None[int](), first)
	})
}

// TestCycleWithChainedOperations tests Cycle with other operations
func TestCycleWithChainedOperations(t *testing.T) {
	t.Run("cycle then map", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		mapped := MonadMap(cycled, N.Mul(10))
		result := toSlice(Take[int](7)(mapped))
		assert.Equal(t, []int{10, 20, 30, 10, 20, 30, 10}, result)
	})

	t.Run("cycle then filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		cycled := Cycle(seq)
		filtered := MonadFilter(cycled, func(x int) bool { return x%2 == 0 })
		result := toSlice(Take[int](6)(filtered))
		assert.Equal(t, []int{2, 4, 2, 4, 2, 4}, result)
	})

	t.Run("map then cycle", func(t *testing.T) {
		seq := From(1, 2, 3)
		mapped := MonadMap(seq, N.Mul(2))
		cycled := Cycle(mapped)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{2, 4, 6, 2, 4, 6, 2}, result)
	})

	t.Run("filter then cycle", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		cycled := Cycle(filtered)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{2, 4, 6, 2, 4, 6, 2}, result)
	})

	t.Run("cycle with multiple takes", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		taken1 := Take[int](10)(cycled)
		taken2 := Take[int](5)(taken1)
		result := toSlice(taken2)
		assert.Equal(t, []int{1, 2, 3, 1, 2}, result)
	})
}

// TestCycleWithReplicate tests Cycle with Replicate
func TestCycleWithReplicate(t *testing.T) {
	t.Run("cycles replicated values", func(t *testing.T) {
		seq := Replicate(3, "X")
		cycled := Cycle(seq)
		result := toSlice(Take[string](7)(cycled))
		assert.Equal(t, []string{"X", "X", "X", "X", "X", "X", "X"}, result)
	})

	t.Run("cycles single replicated value", func(t *testing.T) {
		seq := Replicate(1, 99)
		cycled := Cycle(seq)
		result := toSlice(Take[int](5)(cycled))
		assert.Equal(t, []int{99, 99, 99, 99, 99}, result)
	})
}

// TestCycleWithMakeBy tests Cycle with MakeBy
func TestCycleWithMakeBy(t *testing.T) {
	t.Run("cycles generated sequence", func(t *testing.T) {
		seq := MakeBy(3, func(i int) int { return i * i })
		cycled := Cycle(seq)
		result := toSlice(Take[int](8)(cycled))
		assert.Equal(t, []int{0, 1, 4, 0, 1, 4, 0, 1}, result)
	})

	t.Run("cycles single generated element", func(t *testing.T) {
		seq := MakeBy(1, func(i int) int { return i + 10 })
		cycled := Cycle(seq)
		result := toSlice(Take[int](4)(cycled))
		assert.Equal(t, []int{10, 10, 10, 10}, result)
	})
}

// TestCycleWithPrependAppend tests Cycle with Prepend and Append
func TestCycleWithPrependAppend(t *testing.T) {
	t.Run("cycle prepended sequence", func(t *testing.T) {
		seq := From(2, 3)
		prepended := Prepend(1)(seq)
		cycled := Cycle(prepended)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1}, result)
	})

	t.Run("cycle appended sequence", func(t *testing.T) {
		seq := From(1, 2)
		appended := Append(3)(seq)
		cycled := Cycle(appended)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1}, result)
	})
}

// TestCycleWithFlatten tests Cycle with Flatten
func TestCycleWithFlatten(t *testing.T) {
	t.Run("cycles flattened sequence", func(t *testing.T) {
		nested := From(From(1, 2), From(3))
		flattened := Flatten(nested)
		cycled := Cycle(flattened)
		result := toSlice(Take[int](7)(cycled))
		assert.Equal(t, []int{1, 2, 3, 1, 2, 3, 1}, result)
	})
}

// TestCycleWithChain tests Cycle with Chain
func TestCycleWithChain(t *testing.T) {
	t.Run("cycles chained sequence", func(t *testing.T) {
		seq := From(1, 2)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		cycled := Cycle(chained)
		result := toSlice(Take[int](10)(cycled))
		assert.Equal(t, []int{1, 10, 2, 20, 1, 10, 2, 20, 1, 10}, result)
	})
}

// TestCycleEarlyTermination tests that Cycle respects early termination
func TestCycleEarlyTermination(t *testing.T) {
	t.Run("terminates when yield returns false", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)

		count := 0
		for v := range cycled {
			count++
			if v == 2 && count > 2 {
				break
			}
		}

		// Should have stopped at the second occurrence of 2
		assert.Equal(t, 5, count) // 1, 2, 3, 1, 2
	})

	t.Run("take limits infinite cycle", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		taken := Take[int](100)(cycled)

		result := toSlice(taken)
		assert.Len(t, result, 100)

		// Verify pattern repeats correctly
		for i := 0; i < 100; i++ {
			expected := (i % 3) + 1
			assert.Equal(t, expected, result[i])
		}
	})
}

// TestCycleLargeSequence tests Cycle with larger sequences
func TestCycleLargeSequence(t *testing.T) {
	t.Run("cycles large sequence", func(t *testing.T) {
		data := make([]int, 10)
		for i := range data {
			data[i] = i
		}
		seq := From(data...)
		cycled := Cycle(seq)
		result := toSlice(Take[int](25)(cycled))

		assert.Len(t, result, 25)
		// Verify first cycle
		for i := 0; i < 10; i++ {
			assert.Equal(t, i, result[i])
		}
		// Verify second cycle
		for i := 10; i < 20; i++ {
			assert.Equal(t, i-10, result[i])
		}
		// Verify partial third cycle
		for i := 20; i < 25; i++ {
			assert.Equal(t, i-20, result[i])
		}
	})
}

// TestCycleWithReduce tests Cycle with Reduce (limited by Take)
func TestCycleWithReduce(t *testing.T) {
	t.Run("reduces limited cycled sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		cycled := Cycle(seq)
		limited := Take[int](10)(cycled)
		sum := MonadReduce(limited, func(acc, x int) int { return acc + x }, 0)
		// 1+2+3+1+2+3+1+2+3+1 = 19
		assert.Equal(t, 19, sum)
	})
}

// TestCycleEdgeCases tests edge cases
func TestCycleEdgeCases(t *testing.T) {
	t.Run("cycle with very long take", func(t *testing.T) {
		seq := From(1, 2)
		cycled := Cycle(seq)
		result := toSlice(Take[int](1000)(cycled))
		assert.Len(t, result, 1000)

		// Verify pattern
		for i := 0; i < 1000; i++ {
			expected := (i % 2) + 1
			assert.Equal(t, expected, result[i])
		}
	})

	t.Run("cycle single element many times", func(t *testing.T) {
		seq := From(7)
		cycled := Cycle(seq)
		result := toSlice(Take[int](100)(cycled))
		assert.Len(t, result, 100)
		for _, v := range result {
			assert.Equal(t, 7, v)
		}
	})
}

// Benchmark tests
func BenchmarkCycle(b *testing.B) {
	seq := From(1, 2, 3, 4, 5)
	cycled := Cycle(seq)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taken := Take[int](100)(cycled)
		for range taken {
		}
	}
}

func BenchmarkCycleSingleElement(b *testing.B) {
	seq := From(42)
	cycled := Cycle(seq)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taken := Take[int](100)(cycled)
		for range taken {
		}
	}
}

func BenchmarkCycleWithMap(b *testing.B) {
	seq := From(1, 2, 3, 4, 5)
	cycled := Cycle(seq)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := MonadMap(cycled, N.Mul(2))
		taken := Take[int](100)(mapped)
		for range taken {
		}
	}
}

func BenchmarkCycleWithFilter(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	cycled := Cycle(seq)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := MonadFilter(cycled, func(x int) bool { return x%2 == 0 })
		taken := Take[int](50)(filtered)
		for range taken {
		}
	}
}

// Example tests for documentation
func ExampleCycle() {
	seq := From(1, 2, 3)
	cycled := Cycle(seq)
	result := Take[int](7)(cycled)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 1 2 3 1
}

func ExampleCycle_singleElement() {
	seq := From("X")
	cycled := Cycle(seq)
	result := Take[string](5)(cycled)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: X X X X X
}

func ExampleCycle_withFirst() {
	seq := From(10, 20, 30)
	cycled := Cycle(seq)
	first := First(cycled)

	if value, ok := O.Unwrap(first); ok {
		fmt.Printf("First: %d\n", value)
	}
	// Output: First: 10
}

func ExampleCycle_withFilter() {
	seq := From(1, 2, 3, 4, 5)
	cycled := Cycle(seq)
	evens := MonadFilter(cycled, func(x int) bool { return x%2 == 0 })
	result := Take[int](6)(evens)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 2 4 2 4
}

func ExampleCycle_withMap() {
	seq := From(1, 2, 3)
	cycled := Cycle(seq)
	doubled := MonadMap(cycled, N.Mul(2))
	result := Take[int](7)(doubled)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6 2 4 6 2
}

func ExampleCycle_empty() {
	seq := Empty[int]()
	cycled := Cycle(seq)
	result := Take[int](5)(cycled)

	count := 0
	for range result {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}

func ExampleCycle_exactMultiple() {
	seq := From("A", "B", "C")
	cycled := Cycle(seq)
	result := Take[string](9)(cycled)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: A B C A B C A B C
}

// TestCycleWithZip tests Cycle combined with Zip operator
func TestCycleWithZip(t *testing.T) {
	t.Run("zip infinite cycled sequence with finite sequence", func(t *testing.T) {
		// Create an infinite sequence by cycling
		infinite := Cycle(From(1, 2, 3))
		// Create a finite sequence
		finite := From("a", "b", "c", "d", "e")

		// Zip them together - should stop when finite sequence ends
		zipped := MonadZip(infinite, finite)

		// Convert to slice for verification
		result := make([]struct {
			num int
			str string
		}, 0)
		for num, str := range zipped {
			result = append(result, struct {
				num int
				str string
			}{num, str})
		}

		// Should have 5 pairs (limited by finite sequence)
		assert.Len(t, result, 5)
		assert.Equal(t, 1, result[0].num)
		assert.Equal(t, "a", result[0].str)
		assert.Equal(t, 2, result[1].num)
		assert.Equal(t, "b", result[1].str)
		assert.Equal(t, 3, result[2].num)
		assert.Equal(t, "c", result[2].str)
		assert.Equal(t, 1, result[3].num) // Cycle repeats
		assert.Equal(t, "d", result[3].str)
		assert.Equal(t, 2, result[4].num)
		assert.Equal(t, "e", result[4].str)
	})

	t.Run("zip finite sequence with infinite cycled sequence", func(t *testing.T) {
		// Reverse order: finite first, infinite second
		finite := From(10, 20, 30)
		infinite := Cycle(From("X", "Y"))

		zipped := MonadZip(finite, infinite)

		result := make([]struct {
			num int
			str string
		}, 0)
		for num, str := range zipped {
			result = append(result, struct {
				num int
				str string
			}{num, str})
		}

		// Should have 3 pairs (limited by finite sequence)
		assert.Len(t, result, 3)
		assert.Equal(t, 10, result[0].num)
		assert.Equal(t, "X", result[0].str)
		assert.Equal(t, 20, result[1].num)
		assert.Equal(t, "Y", result[1].str)
		assert.Equal(t, 30, result[2].num)
		assert.Equal(t, "X", result[2].str) // Cycle repeats
	})

	t.Run("zip two cycled sequences with take", func(t *testing.T) {
		// Both sequences are infinite, so we need Take to limit
		cycle1 := Cycle(From(1, 2))
		cycle2 := Cycle(From("a", "b", "c"))

		zipped := MonadZip(cycle1, cycle2)

		// Use Take to limit the infinite result
		count := 0
		result := make([]struct {
			num int
			str string
		}, 0)
		for num, str := range zipped {
			result = append(result, struct {
				num int
				str string
			}{num, str})
			count++
			if count >= 7 {
				break
			}
		}

		assert.Len(t, result, 7)
		// Verify the pattern
		assert.Equal(t, 1, result[0].num)
		assert.Equal(t, "a", result[0].str)
		assert.Equal(t, 2, result[1].num)
		assert.Equal(t, "b", result[1].str)
		assert.Equal(t, 1, result[2].num) // cycle1 repeats
		assert.Equal(t, "c", result[2].str)
		assert.Equal(t, 2, result[3].num)
		assert.Equal(t, "a", result[3].str) // cycle2 repeats
	})
}
