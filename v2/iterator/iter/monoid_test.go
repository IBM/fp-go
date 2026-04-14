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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatMonoid_Identity(t *testing.T) {
	t.Run("left identity", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		seq := From(1, 2, 3)

		result := monoid.Concat(monoid.Empty(), seq)
		collected := slices.Collect(result)

		assert.Equal(t, []int{1, 2, 3}, collected)
	})

	t.Run("right identity", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		seq := From(1, 2, 3)

		result := monoid.Concat(seq, monoid.Empty())
		collected := slices.Collect(result)

		assert.Equal(t, []int{1, 2, 3}, collected)
	})
}

func TestConcatMonoid_Associativity(t *testing.T) {
	monoid := ConcatMonoid[int]()
	seq1 := From(1, 2)
	seq2 := From(3, 4)
	seq3 := From(5, 6)

	// (a + b) + c
	left := monoid.Concat(monoid.Concat(seq1, seq2), seq3)
	leftResult := slices.Collect(left)

	// a + (b + c)
	right := monoid.Concat(seq1, monoid.Concat(seq2, seq3))
	rightResult := slices.Collect(right)

	assert.Equal(t, leftResult, rightResult)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, leftResult)
}

func TestConcatMonoid_DeterministicOrder(t *testing.T) {
	t.Run("concatenates in deterministic order", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)
		seq3 := From(7, 8, 9)

		result := monoid.Concat(monoid.Concat(seq1, seq2), seq3)
		collected := slices.Collect(result)

		// Order is always deterministic
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, collected)
	})

	t.Run("multiple runs produce same order", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		// Run multiple times
		results := make([][]int, 5)
		for i := range 5 {
			result := monoid.Concat(seq1, seq2)
			results[i] = slices.Collect(result)
		}

		// All results should be identical
		expected := []int{1, 2, 3, 4, 5, 6}
		for i, result := range results {
			assert.Equal(t, expected, result, "run %d should match", i)
		}
	})
}

func TestConcatMonoid_WithReduce(t *testing.T) {
	monoid := ConcatMonoid[int]()
	sequences := []Seq[int]{
		From(1, 2, 3),
		From(4, 5, 6),
		From(7, 8, 9),
	}

	result := MonadReduce(From(sequences...), monoid.Concat, monoid.Empty())
	collected := slices.Collect(result)

	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, collected)
}

func TestConcatMonoid_WithFoldMap(t *testing.T) {
	monoid := ConcatMonoid[int]()
	numbers := From(1, 2, 3)

	result := MonadFoldMap(numbers, func(n int) Seq[int] {
		return From(n, n*10, n*100)
	}, monoid)
	collected := slices.Collect(result)

	// Deterministic order: each number's expansion in sequence
	assert.Equal(t, []int{1, 10, 100, 2, 20, 200, 3, 30, 300}, collected)
}

func TestConcatMonoid_ComparisonWithMergeMonoid(t *testing.T) {
	t.Run("ConcatMonoid is deterministic", func(t *testing.T) {
		concatMonoid := ConcatMonoid[int]()
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		result := concatMonoid.Concat(seq1, seq2)
		collected := slices.Collect(result)

		// Always the same order
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected)
	})

	t.Run("MergeMonoid may be non-deterministic", func(t *testing.T) {
		mergeMonoid := MergeMonoid[int](10)
		seq1 := From(1, 2, 3)
		seq2 := From(4, 5, 6)

		result := mergeMonoid.Concat(seq1, seq2)
		collected := slices.Collect(result)

		// Contains all elements but order may vary
		assert.ElementsMatch(t, []int{1, 2, 3, 4, 5, 6}, collected)
		// Note: We can't assert exact order as it's non-deterministic
	})
}

func TestConcatMonoid_EmptySequences(t *testing.T) {
	t.Run("concatenating empty sequences", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		empty1 := Empty[int]()
		empty2 := Empty[int]()

		result := monoid.Concat(empty1, empty2)
		collected := slices.Collect(result)

		assert.Empty(t, collected)
	})

	t.Run("concatenating with empty in middle", func(t *testing.T) {
		monoid := ConcatMonoid[int]()
		seq1 := From(1, 2)
		empty := Empty[int]()
		seq2 := From(3, 4)

		result := monoid.Concat(monoid.Concat(seq1, empty), seq2)
		collected := slices.Collect(result)

		assert.Equal(t, []int{1, 2, 3, 4}, collected)
	})
}

func TestConcatMonoid_WithComplexTypes(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	monoid := ConcatMonoid[Person]()
	seq1 := From(Person{"Alice", 30}, Person{"Bob", 25})
	seq2 := From(Person{"Charlie", 35}, Person{"Diana", 28})

	result := monoid.Concat(seq1, seq2)
	collected := slices.Collect(result)

	expected := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
		{"Diana", 28},
	}
	assert.Equal(t, expected, collected)
}

func BenchmarkConcatMonoid_TwoSequences(b *testing.B) {
	monoid := ConcatMonoid[int]()
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)

	b.ResetTimer()
	for range b.N {
		result := monoid.Concat(seq1, seq2)
		for range result {
		}
	}
}

func BenchmarkConcatMonoid_Reduce(b *testing.B) {
	monoid := ConcatMonoid[int]()
	sequences := []Seq[int]{
		From(1, 2, 3),
		From(4, 5, 6),
		From(7, 8, 9),
		From(10, 11, 12),
	}

	b.ResetTimer()
	for range b.N {
		result := MonadReduce(From(sequences...), monoid.Concat, monoid.Empty())
		for range result {
		}
	}
}

func BenchmarkConcatMonoid_VsMergeMonoid(b *testing.B) {
	seq1 := From(1, 2, 3, 4, 5)
	seq2 := From(6, 7, 8, 9, 10)

	b.Run("ConcatMonoid", func(b *testing.B) {
		monoid := ConcatMonoid[int]()
		b.ResetTimer()
		for range b.N {
			result := monoid.Concat(seq1, seq2)
			for range result {
			}
		}
	})

	b.Run("MergeMonoid", func(b *testing.B) {
		monoid := MergeMonoid[int](10)
		b.ResetTimer()
		for range b.N {
			result := monoid.Concat(seq1, seq2)
			for range result {
			}
		}
	})
}

func ExampleConcatMonoid() {
	monoid := ConcatMonoid[int]()
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)

	result := monoid.Concat(seq1, seq2)
	for v := range result {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}

func ExampleConcatMonoid_identity() {
	monoid := ConcatMonoid[int]()
	seq := From(1, 2, 3)

	// Left identity
	result1 := monoid.Concat(monoid.Empty(), seq)
	for v := range result1 {
		fmt.Println(v)
	}

	// Right identity
	result2 := monoid.Concat(seq, monoid.Empty())
	for v := range result2 {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
	// 1
	// 2
	// 3
}

func ExampleConcatMonoid_reduce() {
	monoid := ConcatMonoid[int]()
	sequences := []Seq[int]{
		From(1, 2, 3),
		From(4, 5, 6),
		From(7, 8, 9),
	}

	result := MonadReduce(From(sequences...), monoid.Concat, monoid.Empty())
	for v := range result {
		fmt.Println(v)
	}
	// Output:
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	// 7
	// 8
	// 9
}

func ExampleConcatMonoid_comparison() {
	seq1 := From(1, 2, 3)
	seq2 := From(4, 5, 6)

	// ConcatMonoid: Sequential, deterministic
	concatMonoid := ConcatMonoid[int]()
	concat := concatMonoid.Concat(seq1, seq2)
	fmt.Println("ConcatMonoid (always same order):")
	for v := range concat {
		fmt.Println(v)
	}

	// MergeMonoid: Concurrent, non-deterministic
	// Note: Output order may vary in actual runs
	mergeMonoid := MergeMonoid[int](10)
	merged := mergeMonoid.Concat(seq1, seq2)
	fmt.Println("\nMergeMonoid (order may vary):")
	collected := slices.Collect(merged)
	// Sort for consistent test output
	slices.Sort(collected)
	for _, v := range collected {
		fmt.Println(v)
	}
	// Output:
	// ConcatMonoid (always same order):
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
	//
	// MergeMonoid (order may vary):
	// 1
	// 2
	// 3
	// 4
	// 5
	// 6
}

// Made with Bob
