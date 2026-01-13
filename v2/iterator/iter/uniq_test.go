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
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestUniqBasic tests basic Uniq functionality
func TestUniqBasic(t *testing.T) {
	t.Run("removes duplicate integers", func(t *testing.T) {
		seq := From(1, 2, 3, 2, 4, 1, 5)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("removes duplicate strings", func(t *testing.T) {
		seq := From("apple", "banana", "apple", "cherry", "banana")
		unique := Uniq(F.Identity[string])
		result := toSlice(unique(seq))
		assert.Equal(t, []string{"apple", "banana", "cherry"}, result)
	})

	t.Run("keeps first occurrence", func(t *testing.T) {
		seq := From(1, 2, 1, 3, 2, 4)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{1, 2, 3, 4}, result)
	})

	t.Run("single element", func(t *testing.T) {
		seq := From(42)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("all duplicates", func(t *testing.T) {
		seq := From(5, 5, 5, 5)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{5}, result)
	})

	t.Run("already unique", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}

// TestUniqEmpty tests Uniq with empty sequences
func TestUniqEmpty(t *testing.T) {
	t.Run("empty integer sequence", func(t *testing.T) {
		seq := Empty[int]()
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Empty(t, result)
	})

	t.Run("empty string sequence", func(t *testing.T) {
		seq := Empty[string]()
		unique := Uniq(F.Identity[string])
		result := toSlice(unique(seq))
		assert.Empty(t, result)
	})
}

// TestUniqWithKeyExtractor tests Uniq with custom key extraction
func TestUniqWithKeyExtractor(t *testing.T) {
	t.Run("unique by string length", func(t *testing.T) {
		seq := From("a", "bb", "c", "dd", "eee", "f")
		uniqueByLength := Uniq(S.Size)
		result := toSlice(uniqueByLength(seq))
		assert.Equal(t, []string{"a", "bb", "eee"}, result)
	})

	t.Run("unique by absolute value", func(t *testing.T) {
		seq := From(1, -1, 2, -2, 3, 1, -3)
		uniqueByAbs := Uniq(func(x int) int {
			if x < 0 {
				return -x
			}
			return x
		})
		result := toSlice(uniqueByAbs(seq))
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("case-insensitive unique strings", func(t *testing.T) {
		seq := From("Hello", "world", "HELLO", "World", "test")
		uniqueCaseInsensitive := Uniq(strings.ToLower)
		result := toSlice(uniqueCaseInsensitive(seq))
		assert.Equal(t, []string{"Hello", "world", "test"}, result)
	})

	t.Run("unique by modulo", func(t *testing.T) {
		seq := From(1, 4, 7, 2, 5, 8, 3)
		uniqueByMod3 := Uniq(func(x int) int { return x % 3 })
		result := toSlice(uniqueByMod3(seq))
		assert.Equal(t, []int{1, 2, 3}, result) // 1%3=1, 4%3=1 (dup), 7%3=1 (dup), 2%3=2, 5%3=2 (dup), 8%3=2 (dup), 3%3=0
	})
}

// TestUniqWithComplexTypes tests Uniq with structs and complex types
func TestUniqWithComplexTypes(t *testing.T) {
	type Person struct {
		ID   int
		Name string
	}

	t.Run("unique structs by ID", func(t *testing.T) {
		seq := From(
			Person{1, "Alice"},
			Person{2, "Bob"},
			Person{1, "Alice2"}, // duplicate ID
			Person{3, "Charlie"},
			Person{2, "Bob2"}, // duplicate ID
		)
		uniqueByID := Uniq(func(p Person) int { return p.ID })
		result := toSlice(uniqueByID(seq))
		assert.Equal(t, []Person{
			{1, "Alice"},
			{2, "Bob"},
			{3, "Charlie"},
		}, result)
	})

	t.Run("unique structs by name", func(t *testing.T) {
		seq := From(
			Person{1, "Alice"},
			Person{2, "Bob"},
			Person{3, "Alice"}, // duplicate name
		)
		uniqueByName := Uniq(func(p Person) string { return p.Name })
		result := toSlice(uniqueByName(seq))
		assert.Equal(t, []Person{
			{1, "Alice"},
			{2, "Bob"},
		}, result)
	})

	t.Run("unique slices by length", func(t *testing.T) {
		seq := From([]int{1, 2}, []int{3}, []int{4, 5}, []int{6})
		uniqueByLength := Uniq(func(s []int) int { return len(s) })
		result := toSlice(uniqueByLength(seq))
		assert.Len(t, result, 2)
		assert.Equal(t, 2, len(result[0]))
		assert.Equal(t, 1, len(result[1]))
	})
}

// TestStrictUniq tests StrictUniq functionality
func TestStrictUniq(t *testing.T) {
	t.Run("removes duplicate integers", func(t *testing.T) {
		seq := From(1, 2, 3, 2, 4, 1, 5)
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("removes duplicate strings", func(t *testing.T) {
		seq := From("apple", "banana", "apple", "cherry", "banana")
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []string{"apple", "banana", "cherry"}, result)
	})

	t.Run("single element", func(t *testing.T) {
		seq := From(42)
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("all duplicates", func(t *testing.T) {
		seq := From("x", "x", "x")
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []string{"x"}, result)
	})

	t.Run("empty sequence", func(t *testing.T) {
		seq := Empty[int]()
		result := toSlice(StrictUniq(seq))
		assert.Empty(t, result)
	})

	t.Run("already unique", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})

	t.Run("boolean values", func(t *testing.T) {
		seq := From(true, false, true, false, true)
		result := toSlice(StrictUniq(seq))
		assert.Equal(t, []bool{true, false}, result)
	})
}

// TestUniqWithChainedOperations tests Uniq combined with other operations
func TestUniqWithChainedOperations(t *testing.T) {
	t.Run("uniq then map", func(t *testing.T) {
		seq := From(1, 2, 3, 2, 4, 1)
		unique := Uniq(F.Identity[int])
		mapped := MonadMap(unique(seq), func(x int) int { return x * 2 })
		result := toSlice(mapped)
		assert.Equal(t, []int{2, 4, 6, 8}, result)
	})

	t.Run("map then uniq", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5)
		mapped := MonadMap(seq, func(x int) int { return x % 3 })
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(mapped))
		assert.Equal(t, []int{1, 2, 0}, result)
	})

	t.Run("filter then uniq", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 6, 2, 4, 6)
		filtered := MonadFilter(seq, func(x int) bool { return x%2 == 0 })
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(filtered))
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("uniq then filter", func(t *testing.T) {
		seq := From(1, 2, 3, 2, 4, 1, 5, 6)
		unique := Uniq(F.Identity[int])
		filtered := MonadFilter(unique(seq), func(x int) bool { return x%2 == 0 })
		result := toSlice(filtered)
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("uniq then take", func(t *testing.T) {
		seq := From(1, 2, 3, 2, 4, 1, 5)
		unique := Uniq(F.Identity[int])
		taken := Take[int](3)(unique(seq))
		result := toSlice(taken)
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("take then uniq", func(t *testing.T) {
		seq := From(1, 2, 1, 3, 2, 4, 5)
		taken := Take[int](5)(seq)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(taken))
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

// TestUniqEarlyTermination tests that Uniq respects early termination
func TestUniqEarlyTermination(t *testing.T) {
	t.Run("terminates when yield returns false", func(t *testing.T) {
		seq := From(1, 2, 3, 4, 5, 2, 6, 7)
		unique := Uniq(F.Identity[int])

		count := 0
		for v := range unique(seq) {
			count++
			if v >= 4 {
				break
			}
		}

		assert.Equal(t, 4, count) // Should stop at 4
	})
}

// TestUniqLargeSequence tests Uniq with larger sequences
func TestUniqLargeSequence(t *testing.T) {
	t.Run("uniq large sequence with many duplicates", func(t *testing.T) {
		// Create sequence with repeating pattern
		data := make([]int, 1000)
		for i := range data {
			data[i] = i % 10 // Only 10 unique values
		}
		seq := From(data...)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))

		assert.Len(t, result, 10)
		for i := 0; i < 10; i++ {
			assert.Equal(t, i, result[i])
		}
	})

	t.Run("uniq large sequence all unique", func(t *testing.T) {
		data := make([]int, 100)
		for i := range data {
			data[i] = i
		}
		seq := From(data...)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))

		assert.Len(t, result, 100)
		for i := 0; i < 100; i++ {
			assert.Equal(t, i, result[i])
		}
	})
}

// TestUniqPreservesOrder tests that Uniq maintains element order
func TestUniqPreservesOrder(t *testing.T) {
	t.Run("maintains order of first occurrences", func(t *testing.T) {
		seq := From(5, 3, 5, 1, 3, 2, 1, 4)
		unique := Uniq(F.Identity[int])
		result := toSlice(unique(seq))
		assert.Equal(t, []int{5, 3, 1, 2, 4}, result)
	})
}

// Benchmark tests
func BenchmarkUniq(b *testing.B) {
	seq := From(1, 2, 3, 2, 4, 1, 5, 3, 6, 4, 7, 5)
	unique := Uniq(F.Identity[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range unique(seq) {
		}
	}
}

func BenchmarkStrictUniq(b *testing.B) {
	seq := From(1, 2, 3, 2, 4, 1, 5, 3, 6, 4, 7, 5)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range StrictUniq(seq) {
		}
	}
}

func BenchmarkUniqLarge(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i % 100
	}
	seq := From(data...)
	unique := Uniq(F.Identity[int])
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for range unique(seq) {
		}
	}
}

// Example tests for documentation
func ExampleUniq() {
	seq := From(1, 2, 3, 2, 4, 1, 5)
	unique := Uniq(F.Identity[int])
	result := unique(seq)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5
}

func ExampleUniq_byLength() {
	seq := From("a", "bb", "c", "dd", "eee")
	uniqueByLength := Uniq(S.Size)
	result := uniqueByLength(seq)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: a bb eee
}

func ExampleUniq_caseInsensitive() {
	seq := From("Hello", "world", "HELLO", "World", "test")
	uniqueCaseInsensitive := Uniq(func(s string) string {
		return strings.ToLower(s)
	})
	result := uniqueCaseInsensitive(seq)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: Hello world test
}

func ExampleStrictUniq() {
	seq := From(1, 2, 3, 2, 4, 1, 5)
	result := StrictUniq(seq)

	for v := range result {
		fmt.Printf("%d ", v)
	}
	// Output: 1 2 3 4 5
}

func ExampleStrictUniq_strings() {
	seq := From("apple", "banana", "apple", "cherry", "banana")
	result := StrictUniq(seq)

	for v := range result {
		fmt.Printf("%s ", v)
	}
	// Output: apple banana cherry
}

func ExampleUniq_empty() {
	seq := Empty[int]()
	unique := Uniq(F.Identity[int])
	result := unique(seq)

	count := 0
	for range result {
		count++
	}
	fmt.Printf("Count: %d\n", count)
	// Output: Count: 0
}
