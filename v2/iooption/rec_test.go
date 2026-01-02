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

package iooption

import (
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// TestTailRecFactorial tests computing factorial using tail recursion with optional result
func TestTailRecFactorial(t *testing.T) {
	type FactState struct {
		n      int
		result int
	}

	factorial := TailRec(func(state FactState) IOOption[TR.Trampoline[FactState, int]] {
		if state.n < 0 {
			// Negative numbers have no factorial
			return None[TR.Trampoline[FactState, int]]()
		}
		if state.n <= 1 {
			// Terminate with final result
			return Of(TR.Land[FactState](state.result))
		}
		// Continue with next iteration
		return Of(TR.Bounce[int](FactState{
			n:      state.n - 1,
			result: state.result * state.n,
		}))
	})

	t.Run("factorial of 5", func(t *testing.T) {
		result := factorial(FactState{n: 5, result: 1})()
		assert.Equal(t, O.Some(120), result)
	})

	t.Run("factorial of 0", func(t *testing.T) {
		result := factorial(FactState{n: 0, result: 1})()
		assert.Equal(t, O.Some(1), result)
	})

	t.Run("factorial of 1", func(t *testing.T) {
		result := factorial(FactState{n: 1, result: 1})()
		assert.Equal(t, O.Some(1), result)
	})

	t.Run("factorial of negative number returns None", func(t *testing.T) {
		result := factorial(FactState{n: -5, result: 1})()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("factorial of 10", func(t *testing.T) {
		result := factorial(FactState{n: 10, result: 1})()
		assert.Equal(t, O.Some(3628800), result)
	})
}

// TestTailRecSafeDivision tests integer division with optional result
func TestTailRecSafeDivision(t *testing.T) {
	type DivState struct {
		numerator   int
		denominator int
		steps       int
	}

	safeDivide := TailRec(func(state DivState) IOOption[TR.Trampoline[DivState, int]] {
		if state.denominator == 0 {
			return None[TR.Trampoline[DivState, int]]() // Division by zero
		}
		if state.numerator < state.denominator {
			return Of(TR.Land[DivState](state.steps))
		}
		return Of(TR.Bounce[int](DivState{
			numerator:   state.numerator - state.denominator,
			denominator: state.denominator,
			steps:       state.steps + 1,
		}))
	})

	t.Run("10 divided by 3", func(t *testing.T) {
		result := safeDivide(DivState{numerator: 10, denominator: 3, steps: 0})()
		assert.Equal(t, O.Some(3), result)
	})

	t.Run("division by zero returns None", func(t *testing.T) {
		result := safeDivide(DivState{numerator: 10, denominator: 0, steps: 0})()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("exact division", func(t *testing.T) {
		result := safeDivide(DivState{numerator: 15, denominator: 5, steps: 0})()
		assert.Equal(t, O.Some(3), result)
	})

	t.Run("numerator less than denominator", func(t *testing.T) {
		result := safeDivide(DivState{numerator: 2, denominator: 5, steps: 0})()
		assert.Equal(t, O.Some(0), result)
	})
}

// TestTailRecFindInRange tests finding a value in a range with optional result
func TestTailRecFindInRange(t *testing.T) {
	type FindState struct {
		current int
		target  int
		max     int
	}

	findInRange := TailRec(func(state FindState) IOOption[TR.Trampoline[FindState, int]] {
		if state.current > state.max {
			return None[TR.Trampoline[FindState, int]]() // Not found
		}
		if state.current == state.target {
			return Of(TR.Land[FindState](state.current))
		}
		return Of(TR.Bounce[int](FindState{
			current: state.current + 1,
			target:  state.target,
			max:     state.max,
		}))
	})

	t.Run("find existing value", func(t *testing.T) {
		result := findInRange(FindState{current: 1, target: 5, max: 10})()
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("value not in range returns None", func(t *testing.T) {
		result := findInRange(FindState{current: 1, target: 15, max: 10})()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("find first value", func(t *testing.T) {
		result := findInRange(FindState{current: 1, target: 1, max: 10})()
		assert.Equal(t, O.Some(1), result)
	})

	t.Run("find last value", func(t *testing.T) {
		result := findInRange(FindState{current: 1, target: 10, max: 10})()
		assert.Equal(t, O.Some(10), result)
	})
}

// TestTailRecSumUntilLimit tests summing numbers until a limit with optional result
func TestTailRecSumUntilLimit(t *testing.T) {
	type SumState struct {
		current int
		sum     int
		limit   int
	}

	sumUntilLimit := TailRec(func(state SumState) IOOption[TR.Trampoline[SumState, int]] {
		if state.sum > state.limit {
			return None[TR.Trampoline[SumState, int]]() // Exceeded limit
		}
		if state.current <= 0 {
			return Of(TR.Land[SumState](state.sum))
		}
		return Of(TR.Bounce[int](SumState{
			current: state.current - 1,
			sum:     state.sum + state.current,
			limit:   state.limit,
		}))
	})

	t.Run("sum within limit", func(t *testing.T) {
		result := sumUntilLimit(SumState{current: 5, sum: 0, limit: 100})()
		assert.Equal(t, O.Some(15), result) // 5+4+3+2+1 = 15
	})

	t.Run("sum exceeds limit returns None", func(t *testing.T) {
		result := sumUntilLimit(SumState{current: 10, sum: 0, limit: 20})()
		assert.Equal(t, O.None[int](), result) // Would exceed 20
	})

	t.Run("sum of zero", func(t *testing.T) {
		result := sumUntilLimit(SumState{current: 0, sum: 0, limit: 100})()
		assert.Equal(t, O.Some(0), result)
	})
}

// TestTailRecCountdown tests a simple countdown with optional result
func TestTailRecCountdown(t *testing.T) {
	countdown := TailRec(func(n int) IOOption[TR.Trampoline[int, string]] {
		if n < 0 {
			return None[TR.Trampoline[int, string]]() // Negative not allowed
		}
		if n == 0 {
			return Of(TR.Land[int]("Done!"))
		}
		return Of(TR.Bounce[string](n - 1))
	})

	t.Run("countdown from 5", func(t *testing.T) {
		result := countdown(5)()
		assert.Equal(t, O.Some("Done!"), result)
	})

	t.Run("countdown from 0", func(t *testing.T) {
		result := countdown(0)()
		assert.Equal(t, O.Some("Done!"), result)
	})

	t.Run("countdown from negative returns None", func(t *testing.T) {
		result := countdown(-5)()
		assert.Equal(t, O.None[string](), result)
	})
}

// TestTailRecStackSafety tests that TailRec doesn't overflow the stack with large iterations
func TestTailRecStackSafety(t *testing.T) {
	// Count down from a large number - this would overflow the stack with regular recursion
	largeCountdown := TailRec(func(n int) IOOption[TR.Trampoline[int, int]] {
		if n < 0 {
			return None[TR.Trampoline[int, int]]()
		}
		if n == 0 {
			return Of(TR.Land[int](0))
		}
		return Of(TR.Bounce[int](n - 1))
	})

	t.Run("large iteration count", func(t *testing.T) {
		// This should complete without stack overflow
		result := largeCountdown(10000)()
		assert.Equal(t, O.Some(0), result)
	})
}

// TestTailRecValidation tests validation with early termination
func TestTailRecValidation(t *testing.T) {
	type ValidationState struct {
		items []int
		index int
	}

	// Validate all items are positive, return count if valid
	validatePositive := TailRec(func(state ValidationState) IOOption[TR.Trampoline[ValidationState, int]] {
		if state.index >= len(state.items) {
			return Of(TR.Land[ValidationState](state.index))
		}
		if state.items[state.index] <= 0 {
			return None[TR.Trampoline[ValidationState, int]]() // Invalid item
		}
		return Of(TR.Bounce[int](ValidationState{
			items: state.items,
			index: state.index + 1,
		}))
	})

	t.Run("all items valid", func(t *testing.T) {
		result := validatePositive(ValidationState{items: []int{1, 2, 3, 4, 5}, index: 0})()
		assert.Equal(t, O.Some(5), result)
	})

	t.Run("invalid item returns None", func(t *testing.T) {
		result := validatePositive(ValidationState{items: []int{1, 2, -3, 4, 5}, index: 0})()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("empty list", func(t *testing.T) {
		result := validatePositive(ValidationState{items: []int{}, index: 0})()
		assert.Equal(t, O.Some(0), result)
	})

	t.Run("first item invalid", func(t *testing.T) {
		result := validatePositive(ValidationState{items: []int{-1, 2, 3}, index: 0})()
		assert.Equal(t, O.None[int](), result)
	})
}

// TestTailRecCollatzConjecture tests the Collatz conjecture with optional result
func TestTailRecCollatzConjecture(t *testing.T) {
	type CollatzState struct {
		n     int
		steps int
	}

	// Count steps to reach 1 in Collatz sequence
	collatz := TailRec(func(state CollatzState) IOOption[TR.Trampoline[CollatzState, int]] {
		if state.n <= 0 {
			return None[TR.Trampoline[CollatzState, int]]() // Invalid input
		}
		if state.n == 1 {
			return Of(TR.Land[CollatzState](state.steps))
		}
		if state.n%2 == 0 {
			return Of(TR.Bounce[int](CollatzState{n: state.n / 2, steps: state.steps + 1}))
		}
		return Of(TR.Bounce[int](CollatzState{n: 3*state.n + 1, steps: state.steps + 1}))
	})

	t.Run("collatz for 1", func(t *testing.T) {
		result := collatz(CollatzState{n: 1, steps: 0})()
		assert.Equal(t, O.Some(0), result)
	})

	t.Run("collatz for 2", func(t *testing.T) {
		result := collatz(CollatzState{n: 2, steps: 0})()
		assert.Equal(t, O.Some(1), result) // 2 -> 1
	})

	t.Run("collatz for 3", func(t *testing.T) {
		result := collatz(CollatzState{n: 3, steps: 0})()
		assert.Equal(t, O.Some(7), result) // 3 -> 10 -> 5 -> 16 -> 8 -> 4 -> 2 -> 1
	})

	t.Run("collatz for negative returns None", func(t *testing.T) {
		result := collatz(CollatzState{n: -5, steps: 0})()
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("collatz for zero returns None", func(t *testing.T) {
		result := collatz(CollatzState{n: 0, steps: 0})()
		assert.Equal(t, O.None[int](), result)
	})
}
