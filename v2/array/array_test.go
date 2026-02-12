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
	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
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

	assert.Equal(t, pair.MakePair(Empty[int](), Empty[int]()), Partition(pred)(Empty[int]()))
	assert.Equal(t, pair.MakePair(From(1), From(3)), Partition(pred)(From(1, 3)))
}

func TestChainOptionK(t *testing.T) {
	src := From(1, 2, 3)

	f := func(i int) O.Option[[]string] {
		if i%2 != 0 {
			return O.Of(From(fmt.Sprintf("a%d", i), fmt.Sprintf("b%d", i)))
		}
		return O.None[[]string]()
	}

	res := ChainOptionK(f)(src)

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

// TestExtract tests the Extract function
func TestExtract(t *testing.T) {
	t.Run("Extract from non-empty array", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Extract(input)
		assert.Equal(t, 1, result)
	})

	t.Run("Extract from single element array", func(t *testing.T) {
		input := []string{"hello"}
		result := Extract(input)
		assert.Equal(t, "hello", result)
	})

	t.Run("Extract from empty array returns zero value", func(t *testing.T) {
		input := []int{}
		result := Extract(input)
		assert.Equal(t, 0, result)
	})

	t.Run("Extract from empty string array returns empty string", func(t *testing.T) {
		input := []string{}
		result := Extract(input)
		assert.Equal(t, "", result)
	})

	t.Run("Extract does not modify original array", func(t *testing.T) {
		original := []int{1, 2, 3}
		originalCopy := []int{1, 2, 3}
		_ = Extract(original)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("Extract with floats", func(t *testing.T) {
		input := []float64{3.14, 2.71, 1.41}
		result := Extract(input)
		assert.Equal(t, 3.14, result)
	})

	t.Run("Extract with structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		input := []Person{
			{"Alice", 30},
			{"Bob", 25},
		}
		result := Extract(input)
		assert.Equal(t, Person{"Alice", 30}, result)
	})
}

// TestExtractComonadLaws tests comonad laws for Extract
func TestExtractComonadLaws(t *testing.T) {
	t.Run("Extract ∘ Of == Identity", func(t *testing.T) {
		value := 42
		result := Extract(Of(value))
		assert.Equal(t, value, result)
	})

	t.Run("Extract ∘ Extend(f) == f", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		f := func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		}

		// Extract(Extend(f)(input)) should equal f(input)
		extended := Extend(f)(input)
		result := Extract(extended)
		expected := f(input)

		assert.Equal(t, expected, result)
	})
}

// TestExtend tests the Extend function
func TestExtend(t *testing.T) {
	t.Run("Extend with sum of suffixes", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		sumSuffix := Extend(func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		})
		result := sumSuffix(input)
		expected := []int{10, 9, 7, 4} // [1+2+3+4, 2+3+4, 3+4, 4]
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with length of suffixes", func(t *testing.T) {
		input := []int{10, 20, 30}
		lengths := Extend(Size[int])
		result := lengths(input)
		expected := []int{3, 2, 1}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with head extraction", func(t *testing.T) {
		input := []int{1, 2, 3}
		duplicate := Extend(func(as []int) int {
			return F.Pipe2(as, Head[int], O.GetOrElse(F.Constant(0)))
		})
		result := duplicate(input)
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with empty array", func(t *testing.T) {
		input := []int{}
		result := Extend(Size[int])(input)
		assert.Equal(t, []int{}, result)
	})

	t.Run("Extend with single element", func(t *testing.T) {
		input := []string{"hello"}
		result := Extend(func(as []string) int { return len(as) })(input)
		expected := []int{1}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend does not modify original array", func(t *testing.T) {
		original := []int{1, 2, 3}
		originalCopy := []int{1, 2, 3}
		_ = Extend(Size[int])(original)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("Extend with string concatenation", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		concat := Extend(func(as []string) string {
			return MonadReduce(as, func(acc, s string) string { return acc + s }, "")
		})
		result := concat(input)
		expected := []string{"abc", "bc", "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with max of suffixes", func(t *testing.T) {
		input := []int{3, 1, 4, 1, 5}
		maxSuffix := Extend(func(as []int) int {
			if len(as) == 0 {
				return 0
			}
			max := as[0]
			for _, v := range as[1:] {
				if v > max {
					max = v
				}
			}
			return max
		})
		result := maxSuffix(input)
		expected := []int{5, 5, 5, 5, 5}
		assert.Equal(t, expected, result)
	})
}

// TestExtendComonadLaws tests comonad laws for Extend
func TestExtendComonadLaws(t *testing.T) {
	t.Run("Left identity: Extend(Extract) == Identity", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Extend(Extract[int])(input)
		assert.Equal(t, input, result)
	})

	t.Run("Right identity: Extract ∘ Extend(f) == f", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		f := func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		}

		// Extract(Extend(f)(input)) should equal f(input)
		result := F.Pipe2(input, Extend(f), Extract[int])
		expected := f(input)

		assert.Equal(t, expected, result)
	})

	t.Run("Associativity: Extend(f) ∘ Extend(g) == Extend(f ∘ Extend(g))", func(t *testing.T) {
		input := []int{1, 2, 3}

		// f: sum of array
		f := func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		}

		// g: length of array
		g := func(as []int) int {
			return len(as)
		}

		// Left side: Extend(f) ∘ Extend(g)
		left := F.Pipe2(input, Extend(g), Extend(f))

		// Right side: Extend(f ∘ Extend(g))
		right := Extend(func(as []int) int {
			return f(Extend(g)(as))
		})(input)

		assert.Equal(t, left, right)
	})
}

// TestExtendComposition tests Extend with other array operations
func TestExtendComposition(t *testing.T) {
	t.Run("Extend after Map", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := F.Pipe2(
			input,
			Map(N.Mul(2)),
			Extend(func(as []int) int {
				return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
			}),
		)
		expected := []int{12, 10, 6} // [2+4+6, 4+6, 6]
		assert.Equal(t, expected, result)
	})

	t.Run("Map after Extend", func(t *testing.T) {
		input := []int{1, 2, 3}
		result := F.Pipe2(
			input,
			Extend(Size[int]),
			Map(N.Mul(10)),
		)
		expected := []int{30, 20, 10}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with Filter", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := F.Pipe2(
			input,
			Filter(func(n int) bool { return n%2 == 0 }),
			Extend(Size[int]),
		)
		expected := []int{3, 2, 1} // lengths of [2,4,6], [4,6], [6]
		assert.Equal(t, expected, result)
	})
}

// TestExtendUseCases demonstrates practical use cases for Extend
func TestExtendUseCases(t *testing.T) {
	t.Run("Running sum (cumulative sum from each position)", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		runningSum := Extend(func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		})
		result := runningSum(input)
		expected := []int{15, 14, 12, 9, 5}
		assert.Equal(t, expected, result)
	})

	t.Run("Sliding window average", func(t *testing.T) {
		input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		windowAvg := Extend(func(as []float64) float64 {
			if len(as) == 0 {
				return 0
			}
			sum := MonadReduce(as, func(acc, x float64) float64 { return acc + x }, 0.0)
			return sum / float64(len(as))
		})
		result := windowAvg(input)
		expected := []float64{3.0, 3.5, 4.0, 4.5, 5.0}
		assert.Equal(t, expected, result)
	})

	t.Run("Check if suffix is sorted", func(t *testing.T) {
		input := []int{1, 2, 3, 2, 1}
		isSorted := Extend(func(as []int) bool {
			for i := 1; i < len(as); i++ {
				if as[i] < as[i-1] {
					return false
				}
			}
			return true
		})
		result := isSorted(input)
		expected := []bool{false, false, false, false, true}
		assert.Equal(t, expected, result)
	})

	t.Run("Count remaining elements", func(t *testing.T) {
		events := []string{"start", "middle", "end"}
		remaining := Extend(Size[string])
		result := remaining(events)
		expected := []int{3, 2, 1}
		assert.Equal(t, expected, result)
	})
}

// TestConcat tests the Concat function
func TestConcat(t *testing.T) {
	t.Run("Concat two non-empty arrays", func(t *testing.T) {
		base := []int{1, 2, 3}
		toAppend := []int{4, 5, 6}
		result := Concat(toAppend)(base)
		expected := []int{1, 2, 3, 4, 5, 6}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with empty array to append", func(t *testing.T) {
		base := []int{1, 2, 3}
		empty := []int{}
		result := Concat(empty)(base)
		assert.Equal(t, base, result)
	})

	t.Run("Concat to empty base array", func(t *testing.T) {
		empty := []int{}
		toAppend := []int{1, 2, 3}
		result := Concat(toAppend)(empty)
		assert.Equal(t, toAppend, result)
	})

	t.Run("Concat two empty arrays", func(t *testing.T) {
		empty1 := []int{}
		empty2 := []int{}
		result := Concat(empty2)(empty1)
		assert.Equal(t, []int{}, result)
	})

	t.Run("Concat strings", func(t *testing.T) {
		words1 := []string{"hello", "world"}
		words2 := []string{"foo", "bar"}
		result := Concat(words2)(words1)
		expected := []string{"hello", "world", "foo", "bar"}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat single element arrays", func(t *testing.T) {
		arr1 := []int{1}
		arr2 := []int{2}
		result := Concat(arr2)(arr1)
		expected := []int{1, 2}
		assert.Equal(t, expected, result)
	})

	t.Run("Does not modify original arrays", func(t *testing.T) {
		base := []int{1, 2, 3}
		toAppend := []int{4, 5, 6}
		baseCopy := []int{1, 2, 3}
		toAppendCopy := []int{4, 5, 6}

		_ = Concat(toAppend)(base)

		assert.Equal(t, baseCopy, base)
		assert.Equal(t, toAppendCopy, toAppend)
	})

	t.Run("Concat with floats", func(t *testing.T) {
		arr1 := []float64{1.1, 2.2}
		arr2 := []float64{3.3, 4.4}
		result := Concat(arr2)(arr1)
		expected := []float64{1.1, 2.2, 3.3, 4.4}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with structs", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		arr1 := []Person{{"Alice", 30}, {"Bob", 25}}
		arr2 := []Person{{"Charlie", 35}}
		result := Concat(arr2)(arr1)
		expected := []Person{{"Alice", 30}, {"Bob", 25}, {"Charlie", 35}}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat large arrays", func(t *testing.T) {
		arr1 := MakeBy(500, F.Identity[int])
		arr2 := MakeBy(500, func(i int) int { return i + 500 })
		result := Concat(arr2)(arr1)

		assert.Equal(t, 1000, len(result))
		assert.Equal(t, 0, result[0])
		assert.Equal(t, 499, result[499])
		assert.Equal(t, 500, result[500])
		assert.Equal(t, 999, result[999])
	})

	t.Run("Concat multiple times", func(t *testing.T) {
		arr1 := []int{1}
		arr2 := []int{2}
		arr3 := []int{3}

		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Concat(arr3),
		)

		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
	})
}

// TestConcatComposition tests Concat with other array operations
func TestConcatComposition(t *testing.T) {
	t.Run("Concat after Map", func(t *testing.T) {
		numbers := []int{1, 2, 3}
		result := F.Pipe2(
			numbers,
			Map(N.Mul(2)),
			Concat([]int{10, 20}),
		)
		expected := []int{2, 4, 6, 10, 20}
		assert.Equal(t, expected, result)
	})

	t.Run("Map after Concat", func(t *testing.T) {
		arr1 := []int{1, 2}
		arr2 := []int{3, 4}
		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Map(N.Mul(2)),
		)
		expected := []int{2, 4, 6, 8}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with Filter", func(t *testing.T) {
		arr1 := []int{1, 2, 3, 4}
		arr2 := []int{5, 6, 7, 8}
		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Filter(func(n int) bool { return n%2 == 0 }),
		)
		expected := []int{2, 4, 6, 8}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with Reduce", func(t *testing.T) {
		arr1 := []int{1, 2, 3}
		arr2 := []int{4, 5, 6}
		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Reduce(func(acc, x int) int { return acc + x }, 0),
		)
		expected := 21 // 1+2+3+4+5+6
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with Reverse", func(t *testing.T) {
		arr1 := []int{1, 2, 3}
		arr2 := []int{4, 5, 6}
		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Reverse[int],
		)
		expected := []int{6, 5, 4, 3, 2, 1}
		assert.Equal(t, expected, result)
	})

	t.Run("Concat with Flatten", func(t *testing.T) {
		arr1 := [][]int{{1, 2}, {3, 4}}
		arr2 := [][]int{{5, 6}}
		result := F.Pipe2(
			arr1,
			Concat(arr2),
			Flatten[int],
		)
		expected := []int{1, 2, 3, 4, 5, 6}
		assert.Equal(t, expected, result)
	})

	t.Run("Multiple Concat operations", func(t *testing.T) {
		arr1 := []int{1}
		arr2 := []int{2}
		arr3 := []int{3}
		arr4 := []int{4}

		result := Concat(arr4)(Concat(arr3)(Concat(arr2)(arr1)))

		expected := []int{1, 2, 3, 4}
		assert.Equal(t, expected, result)
	})
}

// TestConcatUseCases demonstrates practical use cases for Concat
func TestConcatUseCases(t *testing.T) {
	t.Run("Building array incrementally", func(t *testing.T) {
		header := []string{"Name", "Age"}
		data := []string{"Alice", "30"}
		footer := []string{"Total: 1"}

		result := F.Pipe2(
			header,
			Concat(data),
			Concat(footer),
		)

		expected := []string{"Name", "Age", "Alice", "30", "Total: 1"}
		assert.Equal(t, expected, result)
	})

	t.Run("Merging results from multiple operations", func(t *testing.T) {
		evens := Filter(func(n int) bool { return n%2 == 0 })([]int{1, 2, 3, 4, 5, 6})
		odds := Filter(func(n int) bool { return n%2 != 0 })([]int{1, 2, 3, 4, 5, 6})

		result := Concat(odds)(evens)
		expected := []int{2, 4, 6, 1, 3, 5}
		assert.Equal(t, expected, result)
	})

	t.Run("Combining prefix and suffix", func(t *testing.T) {
		prefix := []string{"Mr.", "Dr."}
		names := []string{"Smith", "Jones"}

		addPrefix := func(name string) []string {
			return Map(func(p string) string { return p + " " + name })(prefix)
		}

		result := F.Pipe2(
			names,
			Chain(addPrefix),
			F.Identity[[]string],
		)

		expected := []string{"Mr. Smith", "Dr. Smith", "Mr. Jones", "Dr. Jones"}
		assert.Equal(t, expected, result)
	})

	t.Run("Queue-like behavior", func(t *testing.T) {
		queue := []int{1, 2, 3}
		newItems := []int{4, 5}

		// Add items to end of queue
		updatedQueue := Concat(newItems)(queue)

		assert.Equal(t, []int{1, 2, 3, 4, 5}, updatedQueue)
		assert.Equal(t, 1, updatedQueue[0])                   // Front of queue
		assert.Equal(t, 5, updatedQueue[len(updatedQueue)-1]) // Back of queue
	})

	t.Run("Combining configuration arrays", func(t *testing.T) {
		defaultConfig := []string{"--verbose", "--color"}
		userConfig := []string{"--output=file.txt", "--format=json"}

		finalConfig := Concat(userConfig)(defaultConfig)

		expected := []string{"--verbose", "--color", "--output=file.txt", "--format=json"}
		assert.Equal(t, expected, finalConfig)
	})
}

// TestConcatProperties tests mathematical properties of Concat
func TestConcatProperties(t *testing.T) {
	t.Run("Associativity: (a + b) + c == a + (b + c)", func(t *testing.T) {
		a := []int{1, 2}
		b := []int{3, 4}
		c := []int{5, 6}

		// (a + b) + c
		left := Concat(c)(Concat(b)(a))

		// a + (b + c)
		right := Concat(Concat(c)(b))(a)

		assert.Equal(t, left, right)
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, left)
	})

	t.Run("Identity: a + [] == a and [] + a == a", func(t *testing.T) {
		arr := []int{1, 2, 3}
		empty := []int{}

		// Right identity
		rightResult := Concat(empty)(arr)
		assert.Equal(t, arr, rightResult)

		// Left identity
		leftResult := Concat(arr)(empty)
		assert.Equal(t, arr, leftResult)
	})

	t.Run("Length property: len(a + b) == len(a) + len(b)", func(t *testing.T) {
		testCases := []struct {
			arr1 []int
			arr2 []int
		}{
			{[]int{1, 2, 3}, []int{4, 5}},
			{[]int{1}, []int{2, 3, 4, 5}},
			{[]int{}, []int{1, 2, 3}},
			{[]int{1, 2, 3}, []int{}},
			{MakeBy(100, F.Identity[int]), MakeBy(50, F.Identity[int])},
		}

		for _, tc := range testCases {
			result := Concat(tc.arr2)(tc.arr1)
			expectedLen := len(tc.arr1) + len(tc.arr2)
			assert.Equal(t, expectedLen, len(result))
		}
	})

	t.Run("Order preservation: elements maintain their relative order", func(t *testing.T) {
		arr1 := []int{1, 2, 3}
		arr2 := []int{4, 5, 6}
		result := Concat(arr2)(arr1)

		// Check arr1 elements are in order
		assert.Equal(t, 1, result[0])
		assert.Equal(t, 2, result[1])
		assert.Equal(t, 3, result[2])

		// Check arr2 elements are in order after arr1
		assert.Equal(t, 4, result[3])
		assert.Equal(t, 5, result[4])
		assert.Equal(t, 6, result[5])
	})

	t.Run("Immutability: original arrays are not modified", func(t *testing.T) {
		original1 := []int{1, 2, 3}
		original2 := []int{4, 5, 6}
		copy1 := []int{1, 2, 3}
		copy2 := []int{4, 5, 6}

		_ = Concat(original2)(original1)

		assert.Equal(t, copy1, original1)
		assert.Equal(t, copy2, original2)
	})
}
