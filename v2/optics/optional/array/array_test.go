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
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestAt_GetOption tests the GetOption functionality
func TestAt_GetOption(t *testing.T) {
	t.Run("returns Some for valid index", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40}
		optional := At[int](1)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.Some(20), result)
	})

	t.Run("returns Some for first element", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](0)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.Some(10), result)
	})

	t.Run("returns Some for last element", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](2)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.Some(30), result)
	})

	t.Run("returns None for negative index", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](-1)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for out of bounds index", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](10)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for empty array", func(t *testing.T) {
		numbers := []int{}
		optional := At[int](0)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.None[int](), result)
	})

	t.Run("returns None for nil array", func(t *testing.T) {
		var numbers []int
		optional := At[int](0)

		result := optional.GetOption(numbers)

		assert.Equal(t, O.None[int](), result)
	})
}

// TestAt_Set tests the Set functionality
func TestAt_Set(t *testing.T) {
	t.Run("updates element at valid index", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40}
		optional := At[int](1)

		result := optional.Set(25)(numbers)

		assert.Equal(t, []int{10, 25, 30, 40}, result)
		assert.Equal(t, []int{10, 20, 30, 40}, numbers) // Original unchanged
	})

	t.Run("updates first element", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](0)

		result := optional.Set(5)(numbers)

		assert.Equal(t, []int{5, 20, 30}, result)
	})

	t.Run("updates last element", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](2)

		result := optional.Set(35)(numbers)

		assert.Equal(t, []int{10, 20, 35}, result)
	})

	t.Run("is no-op for negative index", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](-1)

		result := optional.Set(99)(numbers)

		assert.Equal(t, numbers, result)
	})

	t.Run("is no-op for out of bounds index", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](10)

		result := optional.Set(99)(numbers)

		assert.Equal(t, numbers, result)
	})

	t.Run("is no-op for empty array", func(t *testing.T) {
		numbers := []int{}
		optional := At[int](0)

		result := optional.Set(99)(numbers)

		assert.Equal(t, numbers, result)
	})

	t.Run("is no-op for nil array", func(t *testing.T) {
		var numbers []int
		optional := At[int](0)

		result := optional.Set(99)(numbers)

		assert.Equal(t, numbers, result)
	})
}

// TestAt_OptionalLaw1_GetSetNoOp tests Optional Law 1: GetSet Law (No-op on None)
// If GetOption(s) returns None, then Set(a)(s) must return s unchanged (no-op).
func TestAt_OptionalLaw1_GetSetNoOp(t *testing.T) {
	t.Run("out of bounds index - set is no-op", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](10)

		// Verify GetOption returns None
		assert.Equal(t, O.None[int](), optional.GetOption(numbers))

		// Set should be a no-op
		result := optional.Set(99)(numbers)
		assert.Equal(t, numbers, result)
	})

	t.Run("negative index - set is no-op", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](-1)

		// Verify GetOption returns None
		assert.Equal(t, O.None[int](), optional.GetOption(numbers))

		// Set should be a no-op
		result := optional.Set(99)(numbers)
		assert.Equal(t, numbers, result)
	})

	t.Run("empty array - set is no-op", func(t *testing.T) {
		numbers := []int{}
		optional := At[int](0)

		// Verify GetOption returns None
		assert.Equal(t, O.None[int](), optional.GetOption(numbers))

		// Set should be a no-op
		result := optional.Set(99)(numbers)
		assert.Equal(t, numbers, result)
	})

	t.Run("nil array - set is no-op", func(t *testing.T) {
		var numbers []int
		optional := At[int](0)

		// Verify GetOption returns None
		assert.Equal(t, O.None[int](), optional.GetOption(numbers))

		// Set should be a no-op
		result := optional.Set(99)(numbers)
		assert.Equal(t, numbers, result)
	})
}

// TestAt_OptionalLaw2_SetGet tests Optional Law 2: SetGet Law (Get what you Set)
// If GetOption(s) returns Some(_), then GetOption(Set(a)(s)) must return Some(a).
func TestAt_OptionalLaw2_SetGet(t *testing.T) {
	t.Run("set then get returns the set value", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40}
		optional := At[int](1)

		// Verify GetOption returns Some (precondition)
		assert.True(t, O.IsSome(optional.GetOption(numbers)))

		// Set a new value
		newValue := 25
		updated := optional.Set(newValue)(numbers)

		// GetOption on updated should return Some(newValue)
		result := optional.GetOption(updated)
		assert.Equal(t, O.Some(newValue), result)
	})

	t.Run("set first element then get", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](0)

		assert.True(t, O.IsSome(optional.GetOption(numbers)))

		newValue := 5
		updated := optional.Set(newValue)(numbers)

		result := optional.GetOption(updated)
		assert.Equal(t, O.Some(newValue), result)
	})

	t.Run("set last element then get", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](2)

		assert.True(t, O.IsSome(optional.GetOption(numbers)))

		newValue := 35
		updated := optional.Set(newValue)(numbers)

		result := optional.GetOption(updated)
		assert.Equal(t, O.Some(newValue), result)
	})

	t.Run("multiple indices satisfy law", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40, 50}

		for i := range 5 {
			optional := At[int](i)

			assert.True(t, O.IsSome(optional.GetOption(numbers)))

			newValue := i * 100
			updated := optional.Set(newValue)(numbers)

			result := optional.GetOption(updated)
			assert.Equal(t, O.Some(newValue), result)
		}
	})
}

// TestAt_OptionalLaw3_SetSet tests Optional Law 3: SetSet Law (Last Set Wins)
// Setting twice is the same as setting once with the final value.
// Formally: Set(b)(Set(a)(s)) = Set(b)(s)
func TestAt_OptionalLaw3_SetSet(t *testing.T) {
	eqSlice := EQ.FromEquals(func(a, b []int) bool {
		if len(a) != len(b) {
			return false
		}
		for i := range len(a) {
			if a[i] != b[i] {
				return false
			}
		}
		return true
	})

	t.Run("setting twice equals setting once with final value", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40}
		optional := At[int](1)

		// Set twice: first to 25, then to 99
		setTwice := F.Pipe2(
			numbers,
			optional.Set(25),
			optional.Set(99),
		)

		// Set once with final value
		setOnce := optional.Set(99)(numbers)

		assert.True(t, eqSlice.Equals(setTwice, setOnce))
	})

	t.Run("multiple sets - last one wins", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](0)

		// Set multiple times
		result := F.Pipe4(
			numbers,
			optional.Set(1),
			optional.Set(2),
			optional.Set(3),
			optional.Set(4),
		)

		// Should equal setting once with final value
		expected := optional.Set(4)(numbers)

		assert.True(t, eqSlice.Equals(result, expected))
	})

	t.Run("set twice on out of bounds - both no-ops", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](10)

		// Set twice on out of bounds
		setTwice := F.Pipe2(
			numbers,
			optional.Set(25),
			optional.Set(99),
		)

		// Set once on out of bounds
		setOnce := optional.Set(99)(numbers)

		// Both should be no-ops, returning original
		assert.True(t, eqSlice.Equals(setTwice, numbers))
		assert.True(t, eqSlice.Equals(setOnce, numbers))
		assert.True(t, eqSlice.Equals(setTwice, setOnce))
	})
}

// TestAt_EdgeCases tests edge cases and boundary conditions
func TestAt_EdgeCases(t *testing.T) {
	t.Run("single element array", func(t *testing.T) {
		numbers := []int{42}
		optional := At[int](0)

		// Get
		assert.Equal(t, O.Some(42), optional.GetOption(numbers))

		// Set
		updated := optional.Set(99)(numbers)
		assert.Equal(t, []int{99}, updated)

		// Out of bounds
		outOfBounds := At[int](1)
		assert.Equal(t, O.None[int](), outOfBounds.GetOption(numbers))
		assert.Equal(t, numbers, outOfBounds.Set(99)(numbers))
	})

	t.Run("large array", func(t *testing.T) {
		numbers := make([]int, 1000)
		for i := range 1000 {
			numbers[i] = i
		}

		optional := At[int](500)

		// Get
		assert.Equal(t, O.Some(500), optional.GetOption(numbers))

		// Set
		updated := optional.Set(9999)(numbers)
		assert.Equal(t, 9999, updated[500])
		assert.Equal(t, 500, numbers[500]) // Original unchanged
	})

	t.Run("works with different types", func(t *testing.T) {
		// String array
		strings := []string{"a", "b", "c"}
		strOptional := At[string](1)
		assert.Equal(t, O.Some("b"), strOptional.GetOption(strings))
		assert.Equal(t, []string{"a", "x", "c"}, strOptional.Set("x")(strings))

		// Bool array
		bools := []bool{true, false, true}
		boolOptional := At[bool](1)
		assert.Equal(t, O.Some(false), boolOptional.GetOption(bools))
		assert.Equal(t, []bool{true, true, true}, boolOptional.Set(true)(bools))
	})

	t.Run("preserves array capacity", func(t *testing.T) {
		numbers := make([]int, 3, 10)
		numbers[0], numbers[1], numbers[2] = 10, 20, 30

		optional := At[int](1)
		updated := optional.Set(25)(numbers)

		assert.Equal(t, []int{10, 25, 30}, updated)
		assert.Equal(t, 3, len(updated))
	})
}

// TestAt_Integration tests integration scenarios
func TestAt_Integration(t *testing.T) {
	t.Run("multiple optionals on same array", func(t *testing.T) {
		numbers := []int{10, 20, 30, 40}

		first := At[int](0)
		second := At[int](1)
		third := At[int](2)

		// Update multiple indices
		result := F.Pipe3(
			numbers,
			first.Set(1),
			second.Set(2),
			third.Set(3),
		)

		assert.Equal(t, []int{1, 2, 3, 40}, result)
		assert.Equal(t, []int{10, 20, 30, 40}, numbers) // Original unchanged
	})

	t.Run("chaining operations", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](1)

		// Get, verify, set, get again
		original := optional.GetOption(numbers)
		assert.Equal(t, O.Some(20), original)

		updated := optional.Set(25)(numbers)
		newValue := optional.GetOption(updated)
		assert.Equal(t, O.Some(25), newValue)

		// Original still unchanged
		assert.Equal(t, O.Some(20), optional.GetOption(numbers))
	})

	t.Run("conditional update based on current value", func(t *testing.T) {
		numbers := []int{10, 20, 30}
		optional := At[int](1)

		// Get current value and conditionally update
		result := F.Pipe1(
			optional.GetOption(numbers),
			O.Fold(
				func() []int { return numbers },
				func(current int) []int {
					if current > 15 {
						return optional.Set(current * 2)(numbers)
					}
					return numbers
				},
			),
		)

		assert.Equal(t, []int{10, 40, 30}, result)
	})
}
