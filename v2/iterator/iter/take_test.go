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
	"github.com/stretchr/testify/assert"
)

// TestTake tests basic Take functionality
func TestTake(t *testing.T) {
	t.Run("takes first n elements from sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(Take[int](3)(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takes first element", func(t *testing.T) {
		seq := From(10, 20, 30)
		result := toSlice(Take[int](1)(seq))
		assert.Equal(t, []int{10}, result)
	})

	t.Run("takes all elements when n equals length", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := toSlice(Take[int](3)(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takes all elements when n exceeds length", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := toSlice(Take[int](10)(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takes from string sequence", func(t *testing.T) {
		seq := From("a", "b", "c", "d", "e")
		result := toSlice(Take[string](3)(seq))
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("takes from single element sequence", func(t *testing.T) {
		seq := From(42)
		result := toSlice(Take[int](1)(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("takes from large sequence", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := toSlice(Take[int](5)(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// TestTakeZeroOrNegative tests Take with zero or negative values
func TestTakeZeroOrNegative(t *testing.T) {
	t.Run("returns empty sequence when n is zero", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(Take[int](0)(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty sequence when n is negative", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(Take[int](-1)(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty sequence when n is large negative", func(t *testing.T) {
		seq := From("a", "b", "c")
		result := toSlice(Take[string](-100)(seq))
		assert.Empty(t, result)
	})
}

// TestTakeEmpty tests Take with empty sequences
func TestTakeEmpty(t *testing.T) {
	t.Run("returns empty from empty integer sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := toSlice(Take[int](5)(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty from empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		result := toSlice(Take[string](3)(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty when taking zero from empty", func(t *testing.T) {
		seq := Empty[int]()
		result := toSlice(Take[int](0)(seq))
		assert.Empty(t, result)
	})
}

// TestTakeWithComplexTypes tests Take with complex data types
func TestTakeWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("takes structs", func(t *testing.T) {
		seq := From(
			Person{"Alice", 30},
			Person{"Bob", 25},
			Person{"Charlie", 35},
			Person{"David", 28},
		)
		result := toSlice(Take[Person](2)(seq))
		expected := []Person{
			{"Alice", 30},
			{"Bob", 25},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("takes pointers", func(t *testing.T) {
		p1 := &Person{"Alice", 30}
		p2 := &Person{"Bob", 25}
		p3 := &Person{"Charlie", 35}
		seq := From(p1, p2, p3)
		result := toSlice(Take[*Person](2)(seq))
		assert.Equal(t, []*Person{p1, p2}, result)
	})

	t.Run("takes slices", func(t *testing.T) {
		seq := From([]int{1, 2}, []int{3, 4}, []int{5, 6}, []int{7, 8})
		result := toSlice(Take[[]int](2)(seq))
		expected := [][]int{{1, 2}, {3, 4}}
		assert.Equal(t, expected, result)
	})
}

// TestTakeWithChainedOperations tests Take with other sequence operations
func TestTakeWithChainedOperations(t *testing.T) {
	t.Run("take after map", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, N.Mul(2))
		result := toSlice(Take[int](3)(mapped))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("take after filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		result := toSlice(Take[int](3)(filtered))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("map after take", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		taken := Take[int](3)(seq)
		result := toSlice(MonadMap(taken, N.Mul(10)))
		assert.Equal(t, []int{10, 20, 30}, result)
	})

	t.Run("filter after take", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8)
		taken := Take[int](6)(seq)
		result := toSlice(MonadFilter(taken, func(x int) bool { return x%2 == 0 }))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("take after chain", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		result := toSlice(Take[int](4)(chained))
		assert.Equal(t, []int{1, 10, 2, 20}, result)
	})

	t.Run("multiple takes", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		taken1 := Take[int](7)(seq)
		taken2 := Take[int](3)(taken1)
		result := toSlice(taken2)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestTakeWithReplicate tests Take with Replicate
func TestTakeWithReplicate(t *testing.T) {
	t.Run("takes from replicated sequence", func(t *testing.T) {
		seq := Replicate(10, 42)
		result := toSlice(Take[int](3)(seq))
		assert.Equal(t, []int{42, 42, 42}, result)
	})

	t.Run("takes all from short replicate", func(t *testing.T) {
		seq := Replicate(2, "hello")
		result := toSlice(Take[string](5)(seq))
		assert.Equal(t, []string{"hello", "hello"}, result)
	})

	t.Run("takes zero from replicate", func(t *testing.T) {
		seq := Replicate(5, 100)
		result := toSlice(Take[int](0)(seq))
		assert.Empty(t, result)
	})
}

// TestTakeWithMakeBy tests Take with MakeBy
func TestTakeWithMakeBy(t *testing.T) {
	t.Run("takes from generated sequence", func(t *testing.T) {
		seq := MakeBy(10, func(i int) int { return i * i })
		result := toSlice(Take[int](5)(seq))
		assert.Equal(t, []int{0, 1, 4, 9, 16}, result)
	})

	t.Run("takes more than generated", func(t *testing.T) {
		seq := MakeBy(3, func(i int) int { return i + 1 })
		result := toSlice(Take[int](10)(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestTakeWithPrependAppend tests Take with Prepend and Append
func TestTakeWithPrependAppend(t *testing.T) {
	t.Run("take from prepended sequence", func(t *testing.T) {
		seq := From(2, 3, 4, 5)
		prepended := Prepend(1)(seq)
		result := toSlice(Take[int](3)(prepended))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("take from appended sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(4)(seq)
		result := toSlice(Take[int](2)(appended))
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("take includes appended element", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(4)(seq)
		result := toSlice(Take[int](4)(appended))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

// TestTakeWithFlatten tests Take with Flatten
func TestTakeWithFlatten(t *testing.T) {
	t.Run("takes from flattened sequence", func(t *testing.T) {
		nested := From(From(1, 2), From(3, 4), From(5, 6))
		flattened := Flatten(nested)
		result := toSlice(Take[int](4)(flattened))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("takes from flattened with empty inner sequences", func(t *testing.T) {
		nested := From(From(1, 2), Empty[int](), From(3, 4))
		flattened := Flatten(nested)
		result := toSlice(Take[int](3)(flattened))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestTakeDoesNotConsumeEntireSequence tests that Take is lazy
func TestTakeDoesNotConsumeEntireSequence(t *testing.T) {
	t.Run("only consumes needed elements", func(t *testing.T) {
		callCount := 0
		seq := MonadMap(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), func(x int) int {
			callCount++
			return x * 2
		})

		taken := Take[int](3)(seq)

		// Manually iterate to verify lazy evaluation
		result := []int{}
		for v := range taken {
			result = append(result, v)
		}

		assert.Equal(t, []int{2, 4, 6}, result)
		// The map function may be called one extra time to check if there are more elements
		// This is expected behavior with Go's range over iterators
		assert.LessOrEqual(t, callCount, 4, "should not consume significantly more than needed")
		assert.GreaterOrEqual(t, callCount, 3, "should consume at least the needed elements")
	})

	t.Run("stops early with filter", func(t *testing.T) {
		callCount := 0
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool {
			callCount++
			return x%2 == 0
		})

		taken := Take[int](2)(filtered)

		// Manually iterate to verify lazy evaluation
		result := []int{}
		for v := range taken {
			result = append(result, v)
		}

		assert.Equal(t, []int{2, 4}, result)
		// Should stop after finding 2 even numbers, may check a few more elements
		assert.LessOrEqual(t, callCount, 7, "should not consume significantly more than needed")
		assert.GreaterOrEqual(t, callCount, 4, "should consume at least enough to find 2 evens")
	})
}

// TestTakeEdgeCases tests edge cases
func TestTakeEdgeCases(t *testing.T) {
	t.Run("take 1 from single element", func(t *testing.T) {
		seq := From(42)
		result := toSlice(Take[int](1)(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("take 0 from single element", func(t *testing.T) {
		seq := From(42)
		result := toSlice(Take[int](0)(seq))
		assert.Empty(t, result)
	})

	t.Run("take large number from small sequence", func(t *testing.T) {
		seq := From(1, 2)
		result := toSlice(Take[int](1000000)(seq))
		assert.Equal(t, []int{1, 2}, result)
	})

	t.Run("take with very large n", func(t *testing.T) {
		seq := From(1, 2, 3)
		result := toSlice(Take[int](int(^uint(0) >> 1))(seq)) // max int
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// Benchmark tests
func BenchmarkTake(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taken := Take[int](5)(seq)
		for range taken {
		}
	}
}

func BenchmarkTakeLargeSequence(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	seq := From(data...)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taken := Take[int](100)(seq)
		for range taken {
		}
	}
}

func BenchmarkTakeWithMap(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapped := MonadMap(seq, N.Mul(2))
		taken := Take[int](5)(mapped)
		for range taken {
		}
	}
}

func BenchmarkTakeWithFilter(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		taken := Take[int](3)(filtered)
		for range taken {
		}
	}
}

// Example tests for documentation
func ExampleTake() {
	seq := From(1, 2, 3, 4, 5)
	taken := Take[int](3)(seq)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3
}

func ExampleTake_moreThanAvailable() {
	seq := From(1, 2, 3)
	taken := Take[int](10)(seq)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3
}

func ExampleTake_zero() {
	seq := From(1, 2, 3, 4, 5)
	taken := Take[int](0)(seq)

	count := 0
	for range taken {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}

func ExampleTake_withFilter() {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	evens := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
	taken := Take[int](3)(evens)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6
}

func ExampleTake_withMap() {
	seq := From(1, 2, 3, 4, 5)
	doubled := MonadMap(seq, N.Mul(2))
	taken := Take[int](3)(doubled)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6
}

func ExampleTake_chained() {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	result := Take[int](5)(
		MonadFilter(seq, func(x int) bool { return x > 3 }),
	)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 4 5 6 7 8
}

// TestSkipWhile tests basic SkipWhile functionality
func TestSkipWhile(t *testing.T) {
	t.Run("skips while predicate is true", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 2, 1)
		result := toSlice(SkipWhile(func(x int) bool { return x < 4 })(seq))
		assert.Equal(t, []int{4, 5, 2, 1}, result)
	})

	t.Run("skips none when first element fails", func(t *testing.T) {
		seq := From(5, 1, 2, 3)
		result := toSlice(SkipWhile(func(x int) bool { return x < 5 })(seq))
		assert.Equal(t, []int{5, 1, 2, 3}, result)
	})

	t.Run("skips all when predicate always true", func(t *testing.T) {
		seq := From(2, 4, 6, 8)
		result := toSlice(SkipWhile(func(x int) bool { return x%2 == 0 })(seq))
		assert.Empty(t, result)
	})

	t.Run("skips from string sequence", func(t *testing.T) {
		seq := From("a", "b", "c", "1", "d", "e")
		isLetter := func(s string) bool { return s >= "a" && s <= "z" }
		result := toSlice(SkipWhile(isLetter)(seq))
		assert.Equal(t, []string{"1", "d", "e"}, result)
	})

	t.Run("continues after predicate fails", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 1, 2, 3)
		result := toSlice(SkipWhile(func(x int) bool { return x < 4 })(seq))
		assert.Equal(t, []int{4, 1, 2, 3}, result)
	})

	t.Run("skips single element", func(t *testing.T) {
		seq := From(1, 10, 2, 3)
		result := toSlice(SkipWhile(func(x int) bool { return x < 10 })(seq))
		assert.Equal(t, []int{10, 2, 3}, result)
	})
}

// TestSkipWhileEmpty tests SkipWhile with empty sequences
func TestSkipWhileEmpty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := toSlice(SkipWhile(func(x int) bool { return x > 0 })(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty when predicate always satisfied", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(SkipWhile(func(x int) bool { return x < 10 })(seq))
		assert.Empty(t, result)
	})
}

// TestSkipWhileWithComplexTypes tests SkipWhile with complex data types
func TestSkipWhileWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("skips structs while condition met", func(t *testing.T) {
		seq := From(
			Person{"Alice", 25},
			Person{"Bob", 30},
			Person{"Charlie", 35},
			Person{"David", 28},
		)
		result := toSlice(SkipWhile(func(p Person) bool { return p.Age < 35 })(seq))
		expected := []Person{
			{"Charlie", 35},
			{"David", 28},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("skips pointers while condition met", func(t *testing.T) {
		p1 := &Person{"Alice", 25}
		p2 := &Person{"Bob", 30}
		p3 := &Person{"Charlie", 35}
		p4 := &Person{"David", 28}
		seq := From(p1, p2, p3, p4)
		result := toSlice(SkipWhile(func(p *Person) bool { return p.Age < 35 })(seq))
		assert.Equal(t, []*Person{p3, p4}, result)
	})

	t.Run("skips slices while condition met", func(t *testing.T) {
		seq := From([]int{1}, []int{1, 2}, []int{1, 2, 3}, []int{1})
		result := toSlice(SkipWhile(func(s []int) bool { return len(s) < 3 })(seq))
		expected := [][]int{{1, 2, 3}, {1}}
		assert.Equal(t, expected, result)
	})
}

// TestSkipWhileWithChainedOperations tests SkipWhile with other sequence operations
func TestSkipWhileWithChainedOperations(t *testing.T) {
	t.Run("skipWhile after map", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, N.Mul(2))
		result := toSlice(SkipWhile(func(x int) bool { return x < 8 })(mapped))
		assert.Equal(t, []int{8, 10}, result)
	})

	t.Run("skipWhile after filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		result := toSlice(SkipWhile(func(x int) bool { return x < 6 })(filtered))
		assert.Equal(t, []int{6, 8, 10}, result)
	})

	t.Run("map after skipWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		skipped := SkipWhile(func(x int) bool { return x < 4 })(seq)
		result := toSlice(MonadMap(skipped, N.Mul(10)))
		assert.Equal(t, []int{40, 50}, result)
	})

	t.Run("filter after skipWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8)
		skipped := SkipWhile(func(x int) bool { return x < 4 })(seq)
		result := toSlice(MonadFilter(skipped, func(x int) bool { return x%2 == 0 }))
		assert.Equal(t, []int{4, 6, 8}, result)
	})

	t.Run("skipWhile after chain", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		result := toSlice(SkipWhile(func(x int) bool { return x < 20 })(chained))
		assert.Equal(t, []int{20, 3, 30}, result)
	})

	t.Run("skip after skipWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		skipped1 := SkipWhile(func(x int) bool { return x < 4 })(seq)
		skipped2 := Skip[int](2)(skipped1)
		result := toSlice(skipped2)
		assert.Equal(t, []int{6, 7, 8, 9, 10}, result)
	})

	t.Run("skipWhile after skip", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		skipped := Skip[int](3)(seq)
		result := toSlice(SkipWhile(func(x int) bool { return x < 7 })(skipped))
		assert.Equal(t, []int{7, 8, 9, 10}, result)
	})

	t.Run("takeWhile after skipWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		skipped := SkipWhile(func(x int) bool { return x < 4 })(seq)
		taken := TakeWhile(func(x int) bool { return x < 8 })(skipped)
		result := toSlice(taken)
		assert.Equal(t, []int{4, 5, 6, 7}, result)
	})

	t.Run("skipWhile after takeWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		taken := TakeWhile(func(x int) bool { return x < 8 })(seq)
		skipped := SkipWhile(func(x int) bool { return x < 4 })(taken)
		result := toSlice(skipped)
		assert.Equal(t, []int{4, 5, 6, 7}, result)
	})
}

// TestSkipWhileWithReplicate tests SkipWhile with Replicate
func TestSkipWhileWithReplicate(t *testing.T) {
	t.Run("skips all from replicated sequence", func(t *testing.T) {
		seq := Replicate(10, 5)
		result := toSlice(SkipWhile(func(x int) bool { return x == 5 })(seq))
		assert.Empty(t, result)
	})

	t.Run("skips none when predicate fails on replicate", func(t *testing.T) {
		seq := Replicate(5, 10)
		result := toSlice(SkipWhile(func(x int) bool { return x < 10 })(seq))
		assert.Equal(t, []int{10, 10, 10, 10, 10}, result)
	})
}

// TestSkipWhileWithMakeBy tests SkipWhile with MakeBy
func TestSkipWhileWithMakeBy(t *testing.T) {
	t.Run("skips from generated sequence", func(t *testing.T) {
		seq := MakeBy(10, func(i int) int { return i * i })
		result := toSlice(SkipWhile(func(x int) bool { return x < 25 })(seq))
		assert.Equal(t, []int{25, 36, 49, 64, 81}, result)
	})

	t.Run("skips all from generated sequence", func(t *testing.T) {
		seq := MakeBy(5, func(i int) int { return i + 1 })
		result := toSlice(SkipWhile(func(x int) bool { return x < 100 })(seq))
		assert.Empty(t, result)
	})
}

// TestSkipWhileWithPrependAppend tests SkipWhile with Prepend and Append
func TestSkipWhileWithPrependAppend(t *testing.T) {
	t.Run("skipWhile from prepended sequence", func(t *testing.T) {
		seq := From(2, 3, 4, 5)
		prepended := Prepend(1)(seq)
		result := toSlice(SkipWhile(func(x int) bool { return x < 4 })(prepended))
		assert.Equal(t, []int{4, 5}, result)
	})

	t.Run("skipWhile from appended sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(10)(seq)
		result := toSlice(SkipWhile(func(x int) bool { return x < 10 })(appended))
		assert.Equal(t, []int{10}, result)
	})

	t.Run("skipWhile includes appended element", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(4)(seq)
		result := toSlice(SkipWhile(func(x int) bool { return x < 3 })(appended))
		assert.Equal(t, []int{3, 4}, result)
	})
}

// TestSkipWhileWithFlatten tests SkipWhile with Flatten
func TestSkipWhileWithFlatten(t *testing.T) {
	t.Run("skips from flattened sequence", func(t *testing.T) {
		nested := From(From(1, 2), From(3, 4), From(5, 6))
		flattened := Flatten(nested)
		result := toSlice(SkipWhile(func(x int) bool { return x < 4 })(flattened))
		assert.Equal(t, []int{4, 5, 6}, result)
	})

	t.Run("skips from flattened with empty inner sequences", func(t *testing.T) {
		nested := From(From(1, 2), Empty[int](), From(3, 4))
		flattened := Flatten(nested)
		result := toSlice(SkipWhile(func(x int) bool { return x < 3 })(flattened))
		assert.Equal(t, []int{3, 4}, result)
	})
}

// TestSkipWhileDoesNotConsumeEntireSequence tests that SkipWhile is lazy
func TestSkipWhileDoesNotConsumeEntireSequence(t *testing.T) {
	t.Run("only consumes needed elements", func(t *testing.T) {
		callCount := 0
		seq := MonadMap(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), func(x int) int {
			callCount++
			return x * 2
		})

		skipped := SkipWhile(func(x int) bool { return x < 8 })(seq)

		result := []int{}
		for v := range skipped {
			result = append(result, v)
		}

		assert.Equal(t, []int{8, 10, 12, 14, 16, 18, 20}, result)
		// Should process all elements since we iterate through all remaining
		assert.Equal(t, 10, callCount, "should process all elements")
	})

	t.Run("stops early when consumer stops", func(t *testing.T) {
		callCount := 0
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool {
			callCount++
			return x%2 == 0
		})

		skipped := SkipWhile(func(x int) bool { return x < 6 })(filtered)

		result := []int{}
		count := 0
		for v := range skipped {
			result = append(result, v)
			count++
			if count == 2 {
				break
			}
		}

		assert.Equal(t, []int{6, 8}, result)
		// Should stop after getting 2 elements
		assert.LessOrEqual(t, callCount, 9, "should not consume all elements")
	})
}

// TestSkipWhileEdgeCases tests edge cases
func TestSkipWhileEdgeCases(t *testing.T) {
	t.Run("skipWhile with always false predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(SkipWhile(func(x int) bool { return false })(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("skipWhile with always true predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(SkipWhile(func(x int) bool { return true })(seq))
		assert.Empty(t, result)
	})

	t.Run("skipWhile from single element that passes", func(t *testing.T) {
		seq := From(42)
		result := toSlice(SkipWhile(func(x int) bool { return x > 0 })(seq))
		assert.Empty(t, result)
	})

	t.Run("skipWhile from single element that fails", func(t *testing.T) {
		seq := From(42)
		result := toSlice(SkipWhile(func(x int) bool { return x < 0 })(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("skipWhile with complex predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := toSlice(SkipWhile(func(x int) bool {
			return x%2 == 1 || x < 5
		})(seq))
		assert.Equal(t, []int{6, 7, 8, 9, 10}, result)
	})

	t.Run("skipWhile yields elements that satisfy predicate after first failure", func(t *testing.T) {
		seq := From(1, 2, 3, 10, 1, 2, 3)
		result := toSlice(SkipWhile(func(x int) bool { return x < 10 })(seq))
		assert.Equal(t, []int{10, 1, 2, 3}, result)
	})
}

// Benchmark tests for SkipWhile
func BenchmarkSkipWhile(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		skipped := SkipWhile(func(x int) bool { return x < 6 })(seq)
		for range skipped {
		}
	}
}

func BenchmarkSkipWhileLargeSequence(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	seq := From(data...)
	b.ResetTimer()
	for range b.N {
		skipped := SkipWhile(func(x int) bool { return x < 100 })(seq)
		for range skipped {
		}
	}
}

func BenchmarkSkipWhileWithMap(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		mapped := MonadMap(seq, N.Mul(2))
		skipped := SkipWhile(func(x int) bool { return x < 12 })(mapped)
		for range skipped {
		}
	}
}

func BenchmarkSkipWhileWithFilter(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		skipped := SkipWhile(func(x int) bool { return x < 6 })(filtered)
		for range skipped {
		}
	}
}

// Example tests for documentation
func ExampleSkipWhile() {
	seq := From(1, 2, 3, 4, 5, 2, 1)
	skipped := SkipWhile(func(x int) bool { return x < 4 })(seq)

	for v := range skipped {
		fmt.Printf("%d ", v)
	}
	// Output: 4 5 2 1
}

func ExampleSkipWhile_allSatisfy() {
	seq := From(2, 4, 6, 8)
	skipped := SkipWhile(func(x int) bool { return x%2 == 0 })(seq)

	count := 0
	for range skipped {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}

func ExampleSkipWhile_firstFails() {
	seq := From(5, 1, 2, 3)
	skipped := SkipWhile(func(x int) bool { return x < 5 })(seq)

	for v := range skipped {
		fmt.Printf("%d ", v)
	}
	// Output: 5 1 2 3
}

func ExampleSkipWhile_withMap() {
	seq := From(1, 2, 3, 4, 5)
	doubled := MonadMap(seq, N.Mul(2))
	skipped := SkipWhile(func(x int) bool { return x < 8 })(doubled)

	for v := range skipped {
		fmt.Printf("%d ", v)
	}
	// Output: 8 10
}

func ExampleSkipWhile_strings() {
	seq := From("a", "b", "c", "1", "d", "e")
	isLetter := func(s string) bool { return s >= "a" && s <= "z" }
	skipped := SkipWhile(isLetter)(seq)

	for v := range skipped {
		fmt.Printf("%s ", v)
	}
	// Output: 1 d e
}

// TestTakeWhile tests basic TakeWhile functionality
func TestTakeWhile(t *testing.T) {
	t.Run("takes while predicate is true", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 2, 1)
		result := toSlice(TakeWhile(func(x int) bool { return x < 4 })(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takes all when predicate always true", func(t *testing.T) {
		seq := From(2, 4, 6, 8)
		result := toSlice(TakeWhile(func(x int) bool { return x%2 == 0 })(seq))
		assert.Equal(t, []int{2, 4, 6, 8}, result)
	})

	t.Run("takes none when first element fails", func(t *testing.T) {
		seq := From(5, 1, 2, 3)
		result := toSlice(TakeWhile(func(x int) bool { return x < 5 })(seq))
		assert.Empty(t, result)
	})

	t.Run("takes from string sequence", func(t *testing.T) {
		seq := From("a", "b", "c", "1", "d", "e")
		isLetter := func(s string) bool { return s >= "a" && s <= "z" }
		result := toSlice(TakeWhile(isLetter)(seq))
		assert.Equal(t, []string{"a", "b", "c"}, result)
	})

	t.Run("takes single element", func(t *testing.T) {
		seq := From(1, 10, 2, 3)
		result := toSlice(TakeWhile(func(x int) bool { return x < 10 })(seq))
		assert.Equal(t, []int{1}, result)
	})

	t.Run("stops at first false predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 1, 2, 3)
		result := toSlice(TakeWhile(func(x int) bool { return x < 4 })(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestTakeWhileEmpty tests TakeWhile with empty sequences
func TestTakeWhileEmpty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := toSlice(TakeWhile(func(x int) bool { return x > 0 })(seq))
		assert.Empty(t, result)
	})

	t.Run("returns empty when predicate never satisfied", func(t *testing.T) {
		seq := From(10, 20, 30)
		result := toSlice(TakeWhile(func(x int) bool { return x < 5 })(seq))
		assert.Empty(t, result)
	})
}

// TestTakeWhileWithComplexTypes tests TakeWhile with complex data types
func TestTakeWhileWithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("takes structs while condition met", func(t *testing.T) {
		seq := From(
			Person{"Alice", 25},
			Person{"Bob", 30},
			Person{"Charlie", 35},
			Person{"David", 28},
		)
		result := toSlice(TakeWhile(func(p Person) bool { return p.Age < 35 })(seq))
		expected := []Person{
			{"Alice", 25},
			{"Bob", 30},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("takes pointers while condition met", func(t *testing.T) {
		p1 := &Person{"Alice", 25}
		p2 := &Person{"Bob", 30}
		p3 := &Person{"Charlie", 35}
		seq := From(p1, p2, p3)
		result := toSlice(TakeWhile(func(p *Person) bool { return p.Age < 35 })(seq))
		assert.Equal(t, []*Person{p1, p2}, result)
	})

	t.Run("takes slices while condition met", func(t *testing.T) {
		seq := From([]int{1}, []int{1, 2}, []int{1, 2, 3}, []int{1})
		result := toSlice(TakeWhile(func(s []int) bool { return len(s) < 3 })(seq))
		expected := [][]int{{1}, {1, 2}}
		assert.Equal(t, expected, result)
	})
}

// TestTakeWhileWithChainedOperations tests TakeWhile with other sequence operations
func TestTakeWhileWithChainedOperations(t *testing.T) {
	t.Run("takeWhile after map", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, N.Mul(2))
		result := toSlice(TakeWhile(func(x int) bool { return x < 8 })(mapped))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("takeWhile after filter", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		result := toSlice(TakeWhile(func(x int) bool { return x < 7 })(filtered))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("map after takeWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		taken := TakeWhile(func(x int) bool { return x < 4 })(seq)
		result := toSlice(MonadMap(taken, N.Mul(10)))
		assert.Equal(t, []int{10, 20, 30}, result)
	})

	t.Run("filter after takeWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8)
		taken := TakeWhile(func(x int) bool { return x < 7 })(seq)
		result := toSlice(MonadFilter(taken, func(x int) bool { return x%2 == 0 }))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("takeWhile after chain", func(t *testing.T) {
		seq := From(1, 2, 3)
		chained := MonadChain(seq, func(x int) Seq[int] {
			return From(x, x*10)
		})
		result := toSlice(TakeWhile(func(x int) bool { return x < 20 })(chained))
		assert.Equal(t, []int{1, 10, 2}, result)
	})

	t.Run("take after takeWhile", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		taken1 := TakeWhile(func(x int) bool { return x < 8 })(seq)
		taken2 := Take[int](3)(taken1)
		result := toSlice(taken2)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takeWhile after take", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		taken := Take[int](7)(seq)
		result := toSlice(TakeWhile(func(x int) bool { return x < 5 })(taken))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

// TestTakeWhileWithReplicate tests TakeWhile with Replicate
func TestTakeWhileWithReplicate(t *testing.T) {
	t.Run("takes from replicated sequence", func(t *testing.T) {
		seq := Replicate(10, 5)
		result := toSlice(TakeWhile(func(x int) bool { return x == 5 })(seq))
		assert.Equal(t, []int{5, 5, 5, 5, 5, 5, 5, 5, 5, 5}, result)
	})

	t.Run("takes none when predicate fails on replicate", func(t *testing.T) {
		seq := Replicate(5, 10)
		result := toSlice(TakeWhile(func(x int) bool { return x < 10 })(seq))
		assert.Empty(t, result)
	})
}

// TestTakeWhileWithMakeBy tests TakeWhile with MakeBy
func TestTakeWhileWithMakeBy(t *testing.T) {
	t.Run("takes from generated sequence", func(t *testing.T) {
		seq := MakeBy(10, func(i int) int { return i * i })
		result := toSlice(TakeWhile(func(x int) bool { return x < 25 })(seq))
		assert.Equal(t, []int{0, 1, 4, 9, 16}, result)
	})

	t.Run("takes all from generated sequence", func(t *testing.T) {
		seq := MakeBy(5, func(i int) int { return i + 1 })
		result := toSlice(TakeWhile(func(x int) bool { return x < 100 })(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// TestTakeWhileWithPrependAppend tests TakeWhile with Prepend and Append
func TestTakeWhileWithPrependAppend(t *testing.T) {
	t.Run("takeWhile from prepended sequence", func(t *testing.T) {
		seq := From(2, 3, 4, 5)
		prepended := Prepend(1)(seq)
		result := toSlice(TakeWhile(func(x int) bool { return x < 4 })(prepended))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takeWhile from appended sequence", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(10)(seq)
		result := toSlice(TakeWhile(func(x int) bool { return x < 10 })(appended))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("takeWhile includes appended element", func(t *testing.T) {
		seq := From(1, 2, 3)
		appended := Append(4)(seq)
		result := toSlice(TakeWhile(func(x int) bool { return x < 5 })(appended))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})
}

// TestTakeWhileWithFlatten tests TakeWhile with Flatten
func TestTakeWhileWithFlatten(t *testing.T) {
	t.Run("takes from flattened sequence", func(t *testing.T) {
		nested := From(From(1, 2), From(3, 4), From(5, 6))
		flattened := Flatten(nested)
		result := toSlice(TakeWhile(func(x int) bool { return x < 5 })(flattened))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("takes from flattened with empty inner sequences", func(t *testing.T) {
		nested := From(From(1, 2), Empty[int](), From(3, 4))
		flattened := Flatten(nested)
		result := toSlice(TakeWhile(func(x int) bool { return x < 4 })(flattened))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestTakeWhileDoesNotConsumeEntireSequence tests that TakeWhile is lazy
func TestTakeWhileDoesNotConsumeEntireSequence(t *testing.T) {
	t.Run("only consumes needed elements", func(t *testing.T) {
		callCount := 0
		seq := MonadMap(From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), func(x int) int {
			callCount++
			return x * 2
		})

		taken := TakeWhile(func(x int) bool { return x < 8 })(seq)

		result := []int{}
		for v := range taken {
			result = append(result, v)
		}

		assert.Equal(t, []int{2, 4, 6}, result)
		// Should stop after finding element that fails predicate
		assert.LessOrEqual(t, callCount, 5, "should not consume significantly more than needed")
		assert.GreaterOrEqual(t, callCount, 4, "should consume at least enough to find failure")
	})

	t.Run("stops early with filter", func(t *testing.T) {
		callCount := 0
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		filtered := MonadFilter(seq, func(x int) bool {
			callCount++
			return x%2 == 0
		})

		taken := TakeWhile(func(x int) bool { return x < 7 })(filtered)

		result := []int{}
		for v := range taken {
			result = append(result, v)
		}

		assert.Equal(t, []int{2, 4, 6}, result)
		// Should stop after finding even number >= 7
		assert.LessOrEqual(t, callCount, 9, "should not consume significantly more than needed")
		assert.GreaterOrEqual(t, callCount, 7, "should consume at least enough to find 8")
	})
}

// TestTakeWhileEdgeCases tests edge cases
func TestTakeWhileEdgeCases(t *testing.T) {
	t.Run("takeWhile with always false predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(TakeWhile(func(x int) bool { return false })(seq))
		assert.Empty(t, result)
	})

	t.Run("takeWhile with always true predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(TakeWhile(func(x int) bool { return true })(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("takeWhile from single element that passes", func(t *testing.T) {
		seq := From(42)
		result := toSlice(TakeWhile(func(x int) bool { return x > 0 })(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("takeWhile from single element that fails", func(t *testing.T) {
		seq := From(42)
		result := toSlice(TakeWhile(func(x int) bool { return x < 0 })(seq))
		assert.Empty(t, result)
	})

	t.Run("takeWhile with complex predicate", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := toSlice(TakeWhile(func(x int) bool {
			return x%2 == 1 || x < 5
		})(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// Benchmark tests for TakeWhile
func BenchmarkTakeWhile(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		taken := TakeWhile(func(x int) bool { return x < 6 })(seq)
		for range taken {
		}
	}
}

func BenchmarkTakeWhileLargeSequence(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}
	seq := From(data...)
	b.ResetTimer()
	for range b.N {
		taken := TakeWhile(func(x int) bool { return x < 100 })(seq)
		for range taken {
		}
	}
}

func BenchmarkTakeWhileWithMap(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		mapped := MonadMap(seq, N.Mul(2))
		taken := TakeWhile(func(x int) bool { return x < 12 })(mapped)
		for range taken {
		}
	}
}

func BenchmarkTakeWhileWithFilter(b *testing.B) {
	seq := From(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	b.ResetTimer()
	for range b.N {
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		taken := TakeWhile(func(x int) bool { return x < 7 })(filtered)
		for range taken {
		}
	}
}

// Example tests for documentation
func ExampleTakeWhile() {
	seq := From(1, 2, 3, 4, 5, 2, 1)
	taken := TakeWhile(func(x int) bool { return x < 4 })(seq)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3
}

func ExampleTakeWhile_allSatisfy() {
	seq := From(2, 4, 6, 8)
	taken := TakeWhile(func(x int) bool { return x%2 == 0 })(seq)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6 8
}

func ExampleTakeWhile_firstFails() {
	seq := From(5, 1, 2, 3)
	taken := TakeWhile(func(x int) bool { return x < 5 })(seq)

	count := 0
	for range taken {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}

func ExampleTakeWhile_withMap() {
	seq := From(1, 2, 3, 4, 5)
	doubled := MonadMap(seq, N.Mul(2))
	taken := TakeWhile(func(x int) bool { return x < 8 })(doubled)

	for v := range taken {
		fmt.Printf("%d ", v)
	}
	// Output: 2 4 6
}

func ExampleTakeWhile_strings() {
	seq := From("a", "b", "c", "1", "d", "e")
	isLetter := func(s string) bool { return s >= "a" && s <= "z" }
	taken := TakeWhile(isLetter)(seq)

	for v := range taken {
		fmt.Printf("%s ", v)
	}
	// Output: a b c
}
