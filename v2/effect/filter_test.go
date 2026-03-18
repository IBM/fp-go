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

package effect

import (
	"errors"
	"fmt"
	"slices"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type FilterTestConfig struct {
	MaxValue int
	MinValue int
}

// Helper to collect iterator results from an effect
func collectSeqEffect[C, A any](eff Effect[C, Seq[A]], cfg C) []A {
	result, err := runEffect(eff, cfg)
	if err != nil {
		return nil
	}
	return slices.Collect(result)
}

func TestFilterArray_Success(t *testing.T) {
	t.Run("filters array keeping matching elements", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterArray[FilterTestConfig](isPositive)
		input := Succeed[FilterTestConfig]([]int{1, -2, 3, -4, 5})

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 3, 5}, result)
	})

	t.Run("returns empty array when no elements match", func(t *testing.T) {
		// Arrange
		isNegative := N.LessThan(0)
		filterOp := FilterArray[FilterTestConfig](isNegative)
		input := Succeed[FilterTestConfig]([]int{1, 2, 3})

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{}, result)
	})

	t.Run("returns all elements when all match", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(n int) bool { return true }
		filterOp := FilterArray[FilterTestConfig](alwaysTrue)
		input := Succeed[FilterTestConfig]([]int{1, 2, 3})

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{1, 2, 3}, result)
	})
}

func TestFilterIter_Success(t *testing.T) {
	t.Run("filters iterator keeping matching elements", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		filterOp := FilterIter[FilterTestConfig](isEven)
		input := Succeed[FilterTestConfig](slices.Values([]int{1, 2, 3, 4, 5, 6}))

		// Act
		collected := collectSeqEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Equal(t, []int{2, 4, 6}, collected)
	})

	t.Run("returns empty iterator when no elements match", func(t *testing.T) {
		// Arrange
		isNegative := N.LessThan(0)
		filterOp := FilterIter[FilterTestConfig](isNegative)
		input := Succeed[FilterTestConfig](slices.Values([]int{1, 2, 3}))

		// Act
		collected := collectSeqEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Empty(t, collected)
	})
}

func TestFilterArray_WithContext(t *testing.T) {
	t.Run("uses context for filtering", func(t *testing.T) {
		// Arrange
		cfg := FilterTestConfig{MaxValue: 100, MinValue: 0}
		inRange := func(n int) bool { return n >= cfg.MinValue && n <= cfg.MaxValue }
		filterOp := FilterArray[FilterTestConfig](inRange)
		input := Succeed[FilterTestConfig]([]int{-10, 50, 150, 75})

		// Act
		result, err := runEffect(filterOp(input), cfg)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{50, 75}, result)
	})
}

func TestFilterArray_EdgeCases(t *testing.T) {
	t.Run("handles empty array", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterArray[FilterTestConfig](isPositive)
		input := Succeed[FilterTestConfig]([]int{})

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{}, result)
	})

	t.Run("preserves error from input", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterArray[FilterTestConfig](isPositive)
		inputErr := errors.New("input error")
		input := Fail[FilterTestConfig, []int](inputErr)

		// Act
		_, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, inputErr, err)
	})
}

func TestFilterIter_EdgeCases(t *testing.T) {
	t.Run("handles empty iterator", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterIter[FilterTestConfig](isPositive)
		input := Succeed[FilterTestConfig](slices.Values([]int{}))

		// Act
		collected := collectSeqEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Empty(t, collected)
	})

	t.Run("preserves error from input", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterIter[FilterTestConfig](isPositive)
		inputErr := errors.New("input error")
		input := Fail[FilterTestConfig, Seq[int]](inputErr)

		// Act
		_, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, inputErr, err)
	})
}

func TestFilter_GenericFilter(t *testing.T) {
	t.Run("works with custom filter function", func(t *testing.T) {
		// Arrange
		customFilter := func(p Predicate[int]) Endomorphism[[]int] {
			return A.Filter(p)
		}
		filterOp := Filter[FilterTestConfig](customFilter)
		isEven := func(n int) bool { return n%2 == 0 }
		input := Succeed[FilterTestConfig]([]int{1, 2, 3, 4, 5})

		// Act
		result, err := runEffect(filterOp(isEven)(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4}, result)
	})
}

func TestFilterMapArray_Success(t *testing.T) {
	t.Run("filters and maps array elements", func(t *testing.T) {
		// Arrange
		parsePositive := func(n int) O.Option[string] {
			if n > 0 {
				return O.Some(fmt.Sprintf("positive:%d", n))
			}
			return O.None[string]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](parsePositive)
		input := Succeed[FilterTestConfig]([]int{-1, 2, -3, 4, 5})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"positive:2", "positive:4", "positive:5"}, result)
	})

	t.Run("returns empty when no elements match", func(t *testing.T) {
		// Arrange
		neverMatch := func(n int) O.Option[int] {
			return O.None[int]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](neverMatch)
		input := Succeed[FilterTestConfig]([]int{1, 2, 3})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{}, result)
	})

	t.Run("maps all elements when all match", func(t *testing.T) {
		// Arrange
		double := func(n int) O.Option[int] {
			return O.Some(n * 2)
		}
		filterMapOp := FilterMapArray[FilterTestConfig](double)
		input := Succeed[FilterTestConfig]([]int{1, 2, 3})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4, 6}, result)
	})
}

func TestFilterMapIter_Success(t *testing.T) {
	t.Run("filters and maps iterator elements", func(t *testing.T) {
		// Arrange
		doubleEven := func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.Some(n * 2)
			}
			return O.None[int]()
		}
		filterMapOp := FilterMapIter[FilterTestConfig](doubleEven)
		input := Succeed[FilterTestConfig](slices.Values([]int{1, 2, 3, 4, 5}))

		// Act
		collected := collectSeqEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.Equal(t, []int{4, 8}, collected)
	})
}

func TestFilterMapArray_TypeConversion(t *testing.T) {
	t.Run("converts int to string", func(t *testing.T) {
		// Arrange
		intToString := func(n int) O.Option[string] {
			if n > 0 {
				return O.Some(fmt.Sprintf("%d", n))
			}
			return O.None[string]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](intToString)
		input := Succeed[FilterTestConfig]([]int{-1, 2, -3, 4})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"2", "4"}, result)
	})

	t.Run("converts string to int", func(t *testing.T) {
		// Arrange
		parseEven := func(s string) O.Option[int] {
			var n int
			if _, err := fmt.Sscanf(s, "%d", &n); err == nil && n%2 == 0 {
				return O.Some(n)
			}
			return O.None[int]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](parseEven)
		input := Succeed[FilterTestConfig]([]string{"1", "2", "3", "4", "invalid"})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4}, result)
	})
}

func TestFilterMapArray_EdgeCases(t *testing.T) {
	t.Run("handles empty array", func(t *testing.T) {
		// Arrange
		double := func(n int) O.Option[int] {
			return O.Some(n * 2)
		}
		filterMapOp := FilterMapArray[FilterTestConfig](double)
		input := Succeed[FilterTestConfig]([]int{})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{}, result)
	})

	t.Run("preserves error from input", func(t *testing.T) {
		// Arrange
		double := func(n int) O.Option[int] {
			return O.Some(n * 2)
		}
		filterMapOp := FilterMapArray[FilterTestConfig](double)
		inputErr := errors.New("input error")
		input := Fail[FilterTestConfig, []int](inputErr)

		// Act
		_, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, inputErr, err)
	})
}

func TestFilterMapIter_EdgeCases(t *testing.T) {
	t.Run("handles empty iterator", func(t *testing.T) {
		// Arrange
		double := func(n int) O.Option[int] {
			return O.Some(n * 2)
		}
		filterMapOp := FilterMapIter[FilterTestConfig](double)
		input := Succeed[FilterTestConfig](slices.Values([]int{}))

		// Act
		collected := collectSeqEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.Empty(t, collected)
	})
}

func TestFilterMap_GenericFilterMap(t *testing.T) {
	t.Run("works with custom filterMap function", func(t *testing.T) {
		// Arrange
		customFilterMap := func(f O.Kleisli[int, string]) Reader[[]int, []string] {
			return A.FilterMap(f)
		}
		filterMapOp := FilterMap[FilterTestConfig](customFilterMap)
		intToString := func(n int) O.Option[string] {
			if n > 0 {
				return O.Some(fmt.Sprintf("%d", n))
			}
			return O.None[string]()
		}
		input := Succeed[FilterTestConfig]([]int{-1, 2, -3, 4})

		// Act
		result, err := runEffect(filterMapOp(intToString)(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"2", "4"}, result)
	})
}

func TestFilter_Composition(t *testing.T) {
	t.Run("chains multiple filters", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isEven := func(n int) bool { return n%2 == 0 }
		filterPositive := FilterArray[FilterTestConfig](isPositive)
		filterEven := FilterArray[FilterTestConfig](isEven)
		input := Succeed[FilterTestConfig]([]int{-2, -1, 0, 1, 2, 3, 4, 5, 6})

		// Act
		result, err := runEffect(filterEven(filterPositive(input)), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("chains filter and filterMap", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		doubleEven := func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.Some(n * 2)
			}
			return O.None[int]()
		}
		filterOp := FilterArray[FilterTestConfig](isPositive)
		filterMapOp := FilterMapArray[FilterTestConfig](doubleEven)
		input := Succeed[FilterTestConfig]([]int{-2, 1, 2, 3, 4, 5})

		// Act
		result, err := runEffect(filterMapOp(filterOp(input)), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{4, 8}, result)
	})
}

func TestFilter_WithComplexTypes(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	t.Run("filters structs", func(t *testing.T) {
		// Arrange
		isAdult := func(u User) bool { return u.Age >= 18 }
		filterOp := FilterArray[FilterTestConfig](isAdult)
		users := []User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 16},
			{Name: "Charlie", Age: 30},
		}
		input := Succeed[FilterTestConfig](users)

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		expected := []User{
			{Name: "Alice", Age: 25},
			{Name: "Charlie", Age: 30},
		}
		assert.Equal(t, expected, result)
	})

	t.Run("filterMaps structs to different type", func(t *testing.T) {
		// Arrange
		extractAdultName := func(u User) O.Option[string] {
			if u.Age >= 18 {
				return O.Some(u.Name)
			}
			return O.None[string]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](extractAdultName)
		users := []User{
			{Name: "Alice", Age: 25},
			{Name: "Bob", Age: 16},
			{Name: "Charlie", Age: 30},
		}
		input := Succeed[FilterTestConfig](users)

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []string{"Alice", "Charlie"}, result)
	})
}

func TestFilter_BoundaryConditions(t *testing.T) {
	t.Run("filters with boundary predicate", func(t *testing.T) {
		// Arrange
		inRange := func(n int) bool { return n >= 0 && n <= 100 }
		filterOp := FilterArray[FilterTestConfig](inRange)
		input := Succeed[FilterTestConfig]([]int{-1, 0, 50, 100, 101})

		// Act
		result, err := runEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{0, 50, 100}, result)
	})

	t.Run("filterMap with boundary conditions", func(t *testing.T) {
		// Arrange
		clampToRange := func(n int) O.Option[int] {
			if n >= 0 && n <= 100 {
				return O.Some(n)
			}
			return O.None[int]()
		}
		filterMapOp := FilterMapArray[FilterTestConfig](clampToRange)
		input := Succeed[FilterTestConfig]([]int{-1, 0, 50, 100, 101})

		// Act
		result, err := runEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, []int{0, 50, 100}, result)
	})
}

func TestFilter_WithIterators(t *testing.T) {
	t.Run("filters large iterator efficiently", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		filterOp := FilterIter[FilterTestConfig](isEven)

		// Create iterator for range 0-99
		makeSeq := func(yield func(int) bool) {
			for i := range 100 {
				if !yield(i) {
					return
				}
			}
		}
		input := Succeed[FilterTestConfig](Seq[int](makeSeq))

		// Act
		collected := collectSeqEffect(filterOp(input), FilterTestConfig{})

		// Assert
		assert.Equal(t, 50, len(collected))
		assert.Equal(t, 0, collected[0])
		assert.Equal(t, 98, collected[49])
	})

	t.Run("filterMap with iterator", func(t *testing.T) {
		// Arrange
		squareEven := func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.Some(n * n)
			}
			return O.None[int]()
		}
		filterMapOp := FilterMapIter[FilterTestConfig](squareEven)
		input := Succeed[FilterTestConfig](slices.Values([]int{1, 2, 3, 4, 5}))

		// Act
		collected := collectSeqEffect(filterMapOp(input), FilterTestConfig{})

		// Assert
		assert.Equal(t, []int{4, 16}, collected)
	})
}

func TestFilter_ErrorPropagation(t *testing.T) {
	t.Run("filter propagates Left through chain", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filterOp := FilterArray[FilterTestConfig](isPositive)
		originalErr := errors.New("original error")

		// Create an effect that fails
		failedEffect := F.Pipe1(
			Succeed[FilterTestConfig]([]int{1, 2, 3}),
			Chain(func([]int) Effect[FilterTestConfig, []int] {
				return Fail[FilterTestConfig, []int](originalErr)
			}),
		)

		// Act
		_, err := runEffect(filterOp(failedEffect), FilterTestConfig{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
	})

	t.Run("filterMap propagates Left through chain", func(t *testing.T) {
		// Arrange
		double := func(n int) O.Option[int] {
			return O.Some(n * 2)
		}
		filterMapOp := FilterMapArray[FilterTestConfig](double)
		originalErr := errors.New("original error")

		// Create an effect that fails
		failedEffect := F.Pipe1(
			Succeed[FilterTestConfig]([]int{1, 2, 3}),
			Chain(func([]int) Effect[FilterTestConfig, []int] {
				return Fail[FilterTestConfig, []int](originalErr)
			}),
		)

		// Act
		_, err := runEffect(filterMapOp(failedEffect), FilterTestConfig{})

		// Assert
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
	})
}

func TestFilter_Integration(t *testing.T) {
	t.Run("complex filtering pipeline", func(t *testing.T) {
		// Arrange: Filter positive numbers, then double evens, then filter > 5
		isPositive := N.MoreThan(0)
		doubleEven := func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.Some(n * 2)
			}
			return O.None[int]()
		}
		isGreaterThan5 := N.MoreThan(5)

		pipeline := F.Pipe3(
			Succeed[FilterTestConfig]([]int{-2, -1, 0, 1, 2, 3, 4, 5, 6}),
			FilterArray[FilterTestConfig](isPositive),
			FilterMapArray[FilterTestConfig](doubleEven),
			FilterArray[FilterTestConfig](isGreaterThan5),
		)

		// Act
		result, err := runEffect(pipeline, FilterTestConfig{})

		// Assert
		assert.NoError(t, err)
		// Positive: [1,2,3,4,5,6] -> DoubleEven: [4,8,12] -> >5: [8,12]
		assert.Equal(t, []int{8, 12}, result)
	})
}


