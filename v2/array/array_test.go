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

package array

import (
	"fmt"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestMap1(t *testing.T) {

	src := []string{"a", "b", "c"}

	up := Map(strings.ToUpper)(src)

	var up1 = []string{}
	for _, s := range src {
		up1 = append(up1, strings.ToUpper(s))
	}

	var up2 = []string{}
	for i := range src {
		up2 = append(up2, strings.ToUpper(src[i]))
	}

	assert.Equal(t, up, up1)
	assert.Equal(t, up, up2)
}

func TestMap(t *testing.T) {

	mapper := Map(utils.Upper)

	src := []string{"a", "b", "c"}

	dst := mapper(src)

	assert.Equal(t, dst, []string{"A", "B", "C"})
}

func TestReduceRight(t *testing.T) {
	values := From("a", "b", "c")
	f := func(a, acc string) string {
		return fmt.Sprintf("%s%s", acc, a)
	}
	b := ""

	assert.Equal(t, "cba", ReduceRight(f, b)(values))
	assert.Equal(t, "", ReduceRight(f, b)(Empty[string]()))
}

func TestReduce(t *testing.T) {

	values := MakeBy(101, F.Identity[int])

	sum := func(val int, current int) int {
		return val + current
	}
	reducer := Reduce(sum, 0)

	result := reducer(values)
	assert.Equal(t, result, 5050)

}

func TestEmpty(t *testing.T) {
	assert.True(t, IsNonEmpty(MakeBy(101, F.Identity[int])))
	assert.True(t, IsEmpty([]int{}))
}

func TestAp(t *testing.T) {
	assert.Equal(t,
		[]int{2, 4, 6, 3, 6, 9},
		F.Pipe1(
			[]func(int) int{
				utils.Double,
				utils.Triple,
			},
			Ap[int]([]int{1, 2, 3}),
		),
	)
}

func TestIntercalate(t *testing.T) {
	is := Intercalate(S.Monoid)("-")

	assert.Equal(t, "", is(Empty[string]()))
	assert.Equal(t, "a", is([]string{"a"}))
	assert.Equal(t, "a-b-c", is([]string{"a", "b", "c"}))
	assert.Equal(t, "a--c", is([]string{"a", "", "c"}))
	assert.Equal(t, "a-b", is([]string{"a", "b"}))
	assert.Equal(t, "a-b-c-d", is([]string{"a", "b", "c", "d"}))
}

func TestIntersperse(t *testing.T) {
	// Test with empty array
	assert.Equal(t, []int{}, Intersperse(0)([]int{}))

	// Test with single element
	assert.Equal(t, []int{1}, Intersperse(0)([]int{1}))

	// Test with multiple elements
	assert.Equal(t, []int{1, 0, 2, 0, 3}, Intersperse(0)([]int{1, 2, 3}))
}

func TestPrependAll(t *testing.T) {
	empty := Empty[int]()
	prep := PrependAll(0)
	assert.Equal(t, empty, prep(empty))
	assert.Equal(t, []int{0, 1, 0, 2, 0, 3}, prep([]int{1, 2, 3}))
	assert.Equal(t, []int{0, 1}, prep([]int{1}))
	assert.Equal(t, []int{0, 1, 0, 2, 0, 3, 0, 4}, prep([]int{1, 2, 3, 4}))
}

func TestFlatten(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, Flatten([][]int{{1}, {2}, {3}}))
}

func TestLookup(t *testing.T) {
	data := []int{0, 1, 2}
	none := O.None[int]()

	assert.Equal(t, none, Lookup[int](-1)(data))
	assert.Equal(t, none, Lookup[int](10)(data))
	assert.Equal(t, O.Some(1), Lookup[int](1)(data))
}

func TestSlice(t *testing.T) {
	data := []int{0, 1, 2, 3}

	assert.Equal(t, []int{1, 2}, Slice[int](1, 3)(data))
}

func TestFrom(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3}, From(1, 2, 3))
}

func TestPartition(t *testing.T) {

	pred := func(n int) bool {
		return n > 2
	}

	assert.Equal(t, T.MakeTuple2(Empty[int](), Empty[int]()), Partition(pred)(Empty[int]()))
	assert.Equal(t, T.MakeTuple2(From(1), From(3)), Partition(pred)(From(1, 3)))
}

func TestFilterChain(t *testing.T) {
	src := From(1, 2, 3)

	f := func(i int) O.Option[[]string] {
		if i%2 != 0 {
			return O.Of(From(fmt.Sprintf("a%d", i), fmt.Sprintf("b%d", i)))
		}
		return O.None[[]string]()
	}

	res := FilterChain(f)(src)

	assert.Equal(t, From("a1", "b1", "a3", "b3"), res)
}

func TestFilterMap(t *testing.T) {
	src := From(1, 2, 3)

	f := func(i int) O.Option[string] {
		if i%2 != 0 {
			return O.Of(fmt.Sprintf("a%d", i))
		}
		return O.None[string]()
	}

	res := FilterMap(f)(src)

	assert.Equal(t, From("a1", "a3"), res)
}

func TestFoldMap(t *testing.T) {
	src := From("a", "b", "c")

	fold := FoldMap[string](S.Monoid)(strings.ToUpper)

	assert.Equal(t, "ABC", fold(src))
}

func ExampleFoldMap() {
	src := From("a", "b", "c")

	fold := FoldMap[string](S.Monoid)(strings.ToUpper)

	fmt.Println(fold(src))

	// Output: ABC

}

// TestReverse tests the Reverse function
func TestReverse(t *testing.T) {
	t.Run("Reverse integers", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Reverse(input)
		expected := []int{5, 4, 3, 2, 1}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse strings", func(t *testing.T) {
		input := []string{"hello", "world", "foo", "bar"}
		result := Reverse(input)
		expected := []string{"bar", "foo", "world", "hello"}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse empty slice", func(t *testing.T) {
		input := []int{}
		result := Reverse(input)
		assert.Equal(t, []int{}, result)
	})

	t.Run("Reverse single element", func(t *testing.T) {
		input := []string{"only"}
		result := Reverse(input)
		assert.Equal(t, []string{"only"}, result)
	})

	t.Run("Reverse two elements", func(t *testing.T) {
		input := []int{1, 2}
		result := Reverse(input)
		assert.Equal(t, []int{2, 1}, result)
	})

	t.Run("Does not modify original slice", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		originalCopy := []int{1, 2, 3, 4, 5}
		_ = Reverse(original)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("Reverse with floats", func(t *testing.T) {
		input := []float64{1.1, 2.2, 3.3}
		result := Reverse(input)
		expected := []float64{3.3, 2.2, 1.1}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse with structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		input := []Person{
			{"Alice", 30},
			{"Bob", 25},
			{"Charlie", 35},
		}
		result := Reverse(input)
		expected := []Person{
			{"Charlie", 35},
			{"Bob", 25},
			{"Alice", 30},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse with pointers", func(t *testing.T) {
		a, b, c := 1, 2, 3
		input := []*int{&a, &b, &c}
		result := Reverse(input)
		assert.Equal(t, []*int{&c, &b, &a}, result)
	})

	t.Run("Double reverse returns original order", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		reversed := Reverse(original)
		doubleReversed := Reverse(reversed)
		assert.Equal(t, original, doubleReversed)
	})

	t.Run("Reverse with large slice", func(t *testing.T) {
		input := MakeBy(1000, F.Identity[int])
		result := Reverse(input)

		// Check first and last elements
		assert.Equal(t, 999, result[0])
		assert.Equal(t, 0, result[999])

		// Check length
		assert.Equal(t, 1000, len(result))
	})

	t.Run("Reverse palindrome", func(t *testing.T) {
		input := []int{1, 2, 3, 2, 1}
		result := Reverse(input)
		assert.Equal(t, input, result)
	})
}

// TestReverseComposition tests Reverse with other array operations
func TestReverseComposition(t *testing.T) {
	t.Run("Reverse after Map", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := F.Pipe2(
			input,
			Map(N.Mul(2)),
			Reverse[int],
		)
		expected := []int{10, 8, 6, 4, 2}
		assert.Equal(t, expected, result)
	})

	t.Run("Map after Reverse", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := F.Pipe2(
			input,
			Reverse[int],
			Map(N.Mul(2)),
		)
		expected := []int{10, 8, 6, 4, 2}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse with Filter", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := F.Pipe2(
			input,
			Filter(func(n int) bool { return n%2 == 0 }),
			Reverse[int],
		)
		expected := []int{6, 4, 2}
		assert.Equal(t, expected, result)
	})

	t.Run("Reverse with Reduce", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		reversed := Reverse(input)
		result := Reduce(func(acc, val string) string {
			return acc + val
		}, "")(reversed)
		assert.Equal(t, "cba", result)
	})

	t.Run("Reverse with Flatten", func(t *testing.T) {
		input := [][]int{{1, 2}, {3, 4}, {5, 6}}
		result := F.Pipe2(
			input,
			Reverse[[]int],
			Flatten[int],
		)
		expected := []int{5, 6, 3, 4, 1, 2}
		assert.Equal(t, expected, result)
	})
}

// TestReverseUseCases demonstrates practical use cases for Reverse
func TestReverseUseCases(t *testing.T) {
	t.Run("Process events in reverse chronological order", func(t *testing.T) {
		events := []string{"2024-01-01", "2024-01-02", "2024-01-03"}
		reversed := Reverse(events)

		// Most recent first
		assert.Equal(t, "2024-01-03", reversed[0])
		assert.Equal(t, "2024-01-01", reversed[2])
	})

	t.Run("Implement stack behavior (LIFO)", func(t *testing.T) {
		stack := []int{1, 2, 3, 4, 5}
		reversed := Reverse(stack)

		// Pop from reversed (LIFO)
		assert.Equal(t, 5, reversed[0])
		assert.Equal(t, 4, reversed[1])
	})

	t.Run("Reverse string characters", func(t *testing.T) {
		chars := []rune("hello")
		reversed := Reverse(chars)
		result := string(reversed)
		assert.Equal(t, "olleh", result)
	})

	t.Run("Check palindrome", func(t *testing.T) {
		word := []rune("racecar")
		reversed := Reverse(word)
		assert.Equal(t, word, reversed)

		notPalindrome := []rune("hello")
		reversedNot := Reverse(notPalindrome)
		assert.NotEqual(t, notPalindrome, reversedNot)
	})

	t.Run("Reverse transformation pipeline", func(t *testing.T) {
		// Apply transformations in reverse order
		numbers := []int{1, 2, 3}

		// Normal: add 10, then multiply by 2
		normal := F.Pipe2(
			numbers,
			Map(N.Add(10)),
			Map(N.Mul(2)),
		)

		// Reversed order of operations
		reversed := F.Pipe2(
			numbers,
			Map(N.Mul(2)),
			Map(N.Add(10)),
		)

		assert.NotEqual(t, normal, reversed)
		assert.Equal(t, []int{22, 24, 26}, normal)
		assert.Equal(t, []int{12, 14, 16}, reversed)
	})
}

// TestReverseProperties tests mathematical properties of Reverse
func TestReverseProperties(t *testing.T) {
	t.Run("Involution property: Reverse(Reverse(x)) == x", func(t *testing.T) {
		testCases := [][]int{
			{1, 2, 3, 4, 5},
			{1},
			{},
			{1, 2},
			{5, 4, 3, 2, 1},
		}

		for _, original := range testCases {
			result := Reverse(Reverse(original))
			assert.Equal(t, original, result)
		}
	})

	t.Run("Length preservation: len(Reverse(x)) == len(x)", func(t *testing.T) {
		testCases := [][]int{
			{1, 2, 3, 4, 5},
			{1},
			{},
			MakeBy(100, F.Identity[int]),
		}

		for _, input := range testCases {
			result := Reverse(input)
			assert.Equal(t, len(input), len(result))
		}
	})

	t.Run("First element becomes last", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Reverse(input)

		if len(input) > 0 {
			assert.Equal(t, input[0], result[len(result)-1])
			assert.Equal(t, input[len(input)-1], result[0])
		}
	})
}
