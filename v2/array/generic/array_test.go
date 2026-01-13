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

package generic

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

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

	t.Run("Extract with custom slice type", func(t *testing.T) {
		type IntSlice []int
		input := IntSlice{10, 20, 30}
		result := Extract(input)
		assert.Equal(t, 10, result)
	})
}

// TestExtractComonadLaws tests comonad laws for Extract
func TestExtractComonadLaws(t *testing.T) {
	t.Run("Extract ∘ Of == Identity", func(t *testing.T) {
		value := 42
		result := Extract(Of[[]int](value))
		assert.Equal(t, value, result)
	})

	t.Run("Extract ∘ Extend(f) == f", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		f := func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		}

		// Extract(Extend(f)(input)) should equal f(input)
		extended := Extend[[]int, []int](f)(input)
		result := Extract(extended)
		expected := f(input)

		assert.Equal(t, expected, result)
	})
}

// TestExtend tests the Extend function
func TestExtend(t *testing.T) {
	t.Run("Extend with sum of suffixes", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		sumSuffix := Extend[[]int, []int](func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		})
		result := sumSuffix(input)
		expected := []int{10, 9, 7, 4} // [1+2+3+4, 2+3+4, 3+4, 4]
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with length of suffixes", func(t *testing.T) {
		input := []int{10, 20, 30}
		lengths := Extend[[]int, []int](Size[[]int, int])
		result := lengths(input)
		expected := []int{3, 2, 1}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with head extraction", func(t *testing.T) {
		input := []int{1, 2, 3}
		duplicate := Extend[[]int, []int](Extract[[]int, int])
		result := duplicate(input)
		expected := []int{1, 2, 3}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with empty array", func(t *testing.T) {
		input := []int{}
		result := Extend[[]int, []int](Size[[]int, int])(input)
		assert.Equal(t, []int{}, result)
	})

	t.Run("Extend with single element", func(t *testing.T) {
		input := []string{"hello"}
		result := Extend[[]string, []int](func(as []string) int { return len(as) })(input)
		expected := []int{1}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend does not modify original array", func(t *testing.T) {
		original := []int{1, 2, 3}
		originalCopy := []int{1, 2, 3}
		_ = Extend[[]int, []int](Size[[]int, int])(original)
		assert.Equal(t, originalCopy, original)
	})

	t.Run("Extend with string concatenation", func(t *testing.T) {
		input := []string{"a", "b", "c"}
		concat := Extend[[]string, []string](func(as []string) string {
			return MonadReduce(as, func(acc, s string) string { return acc + s }, "")
		})
		result := concat(input)
		expected := []string{"abc", "bc", "c"}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with custom slice types", func(t *testing.T) {
		type IntSlice []int
		type ResultSlice []int
		input := IntSlice{1, 2, 3}
		sumSuffix := Extend[IntSlice, ResultSlice](func(as IntSlice) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		})
		result := sumSuffix(input)
		expected := ResultSlice{6, 5, 3}
		assert.Equal(t, expected, result)
	})
}

// TestExtendComonadLaws tests comonad laws for Extend
func TestExtendComonadLaws(t *testing.T) {
	t.Run("Left identity: Extend(Extract) == Identity", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		result := Extend[[]int, []int](Extract[[]int, int])(input)
		assert.Equal(t, input, result)
	})

	t.Run("Right identity: Extract ∘ Extend(f) == f", func(t *testing.T) {
		input := []int{1, 2, 3, 4}
		f := func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		}

		// Extract(Extend(f)(input)) should equal f(input)
		result := F.Pipe2(input, Extend[[]int, []int](f), Extract[[]int, int])
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
		left := F.Pipe2(input, Extend[[]int, []int](g), Extend[[]int, []int](f))

		// Right side: Extend(f ∘ Extend(g))
		right := Extend[[]int, []int](func(as []int) int {
			return f(Extend[[]int, []int](g)(as))
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
			Map[[]int, []int](func(x int) int { return x * 2 }),
			Extend[[]int, []int](func(as []int) int {
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
			Extend[[]int, []int](Size[[]int, int]),
			Map[[]int, []int](func(x int) int { return x * 10 }),
		)
		expected := []int{30, 20, 10}
		assert.Equal(t, expected, result)
	})

	t.Run("Extend with Filter", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5, 6}
		result := F.Pipe2(
			input,
			Filter[[]int](func(n int) bool { return n%2 == 0 }),
			Extend[[]int, []int](Size[[]int, int]),
		)
		expected := []int{3, 2, 1} // lengths of [2,4,6], [4,6], [6]
		assert.Equal(t, expected, result)
	})
}

// TestExtendUseCases demonstrates practical use cases for Extend
func TestExtendUseCases(t *testing.T) {
	t.Run("Running sum (cumulative sum from each position)", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		runningSum := Extend[[]int, []int](func(as []int) int {
			return MonadReduce(as, func(acc, x int) int { return acc + x }, 0)
		})
		result := runningSum(input)
		expected := []int{15, 14, 12, 9, 5}
		assert.Equal(t, expected, result)
	})

	t.Run("Sliding window average", func(t *testing.T) {
		input := []float64{1.0, 2.0, 3.0, 4.0, 5.0}
		windowAvg := Extend[[]float64, []float64](func(as []float64) float64 {
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
		isSorted := Extend[[]int, []bool](func(as []int) bool {
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
		remaining := Extend[[]string, []int](Size[[]string, string])
		result := remaining(events)
		expected := []int{3, 2, 1}
		assert.Equal(t, expected, result)
	})
}
