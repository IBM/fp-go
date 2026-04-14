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

package traversal

import (
	"testing"

	AR "github.com/IBM/fp-go/v2/array"
	C "github.com/IBM/fp-go/v2/constant"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	AT "github.com/IBM/fp-go/v2/optics/traversal/array/const"
	AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
	"github.com/stretchr/testify/assert"
)

func TestGetAll(t *testing.T) {

	as := AR.From(1, 2, 3)

	tr := AT.FromArray[int](AR.Monoid[int]())

	sa := F.Pipe1(
		Id[[]int, C.Const[[]int, []int]](),
		Compose[[]int, C.Const[[]int, []int], []int, int](tr),
	)

	getall := GetAll[int](as)(sa)

	assert.Equal(t, AR.From(1, 2, 3), getall)
}

func TestFold(t *testing.T) {

	monoidSum := N.MonoidSum[int]()

	as := AR.From(1, 2, 3)

	tr := AT.FromArray[int, int](monoidSum)

	sa := F.Pipe1(
		Id[[]int, C.Const[int, []int]](),
		Compose[[]int, C.Const[int, []int], []int, int](tr),
	)

	folded := Fold(sa)(as)

	assert.Equal(t, 6, folded)
}

func TestTraverse(t *testing.T) {

	as := AR.From(1, 2, 3)

	tr := AI.FromArray[int]()

	sa := F.Pipe1(
		Id[[]int, []int](),
		Compose[[]int, []int](tr),
	)

	res := sa(utils.Double)(as)

	assert.Equal(t, AR.From(2, 4, 6), res)
}

func TestFilter_Success(t *testing.T) {
	t.Run("filters and modifies only matching elements", func(t *testing.T) {
		// Arrange
		numbers := []int{-2, -1, 0, 1, 2, 3}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		// Filter to only positive numbers
		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act - double only positive numbers
		result := filteredTraversal(func(n int) int { return n * 2 })(numbers)

		// Assert
		assert.Equal(t, []int{-2, -1, 0, 2, 4, 6}, result)
	})

	t.Run("filters even numbers and triples them", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5, 6}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		// Filter to only even numbers
		isEven := func(n int) bool { return n%2 == 0 }
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isEven),
		)

		// Act
		result := filteredTraversal(func(n int) int { return n * 3 })(numbers)

		// Assert
		assert.Equal(t, []int{1, 6, 3, 12, 5, 18}, result)
	})

	t.Run("filters strings by length", func(t *testing.T) {
		// Arrange
		words := []string{"a", "ab", "abc", "abcd", "abcde"}
		arrayTraversal := AI.FromArray[string]()
		baseTraversal := F.Pipe1(
			Id[[]string, []string](),
			Compose[[]string, []string, []string, string](arrayTraversal),
		)

		// Filter strings with length > 2
		longerThanTwo := func(s string) bool { return len(s) > 2 }
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]string, []string, string, string](F.Identity[string], F.Identity[func(string) string])(longerThanTwo),
		)

		// Act - convert to uppercase
		result := filteredTraversal(func(s string) string {
			return s + "!"
		})(words)

		// Assert
		assert.Equal(t, []string{"a", "ab", "abc!", "abcd!", "abcde!"}, result)
	})
}

func TestFilter_EdgeCases(t *testing.T) {
	t.Run("empty array returns empty array", func(t *testing.T) {
		// Arrange
		numbers := []int{}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act
		result := filteredTraversal(utils.Double)(numbers)

		// Assert
		assert.Equal(t, []int{}, result)
	})

	t.Run("no elements match predicate", func(t *testing.T) {
		// Arrange
		numbers := []int{-5, -4, -3, -2, -1}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act
		result := filteredTraversal(utils.Double)(numbers)

		// Assert - all elements unchanged
		assert.Equal(t, []int{-5, -4, -3, -2, -1}, result)
	})

	t.Run("all elements match predicate", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act
		result := filteredTraversal(utils.Double)(numbers)

		// Assert - all elements doubled
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})

	t.Run("single element matching", func(t *testing.T) {
		// Arrange
		numbers := []int{42}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act
		result := filteredTraversal(utils.Double)(numbers)

		// Assert
		assert.Equal(t, []int{84}, result)
	})

	t.Run("single element not matching", func(t *testing.T) {
		// Arrange
		numbers := []int{-42}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isPositive := N.MoreThan(0)
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isPositive),
		)

		// Act
		result := filteredTraversal(utils.Double)(numbers)

		// Assert
		assert.Equal(t, []int{-42}, result)
	})
}

func TestFilter_Integration(t *testing.T) {
	t.Run("multiple filters composed", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		// Filter to only even numbers, then only those > 4
		isEven := func(n int) bool { return n%2 == 0 }
		greaterThanFour := N.MoreThan(4)

		filteredTraversal := F.Pipe2(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isEven),
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(greaterThanFour),
		)

		// Act - add 100 to matching elements
		result := filteredTraversal(func(n int) int { return n + 100 })(numbers)

		// Assert - only 6, 8, 10 should be modified
		assert.Equal(t, []int{1, 2, 3, 4, 5, 106, 7, 108, 9, 110}, result)
	})

	t.Run("filter with identity transformation", func(t *testing.T) {
		// Arrange
		numbers := []int{1, 2, 3, 4, 5}
		arrayTraversal := AI.FromArray[int]()
		baseTraversal := F.Pipe1(
			Id[[]int, []int](),
			Compose[[]int, []int](arrayTraversal),
		)

		isEven := func(n int) bool { return n%2 == 0 }
		filteredTraversal := F.Pipe1(
			baseTraversal,
			Filter[[]int, []int](F.Identity[int], F.Identity[func(int) int])(isEven),
		)

		// Act - identity transformation
		result := filteredTraversal(F.Identity[int])(numbers)

		// Assert - array unchanged
		assert.Equal(t, []int{1, 2, 3, 4, 5}, result)
	})
}
