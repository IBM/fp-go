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

package io

import (
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// TestTailRec_Factorial tests computing factorial using tail recursion
func TestTailRec_Factorial(t *testing.T) {
	type FactState struct {
		n      int
		result int
	}

	factorial := TailRec(func(state FactState) IO[TR.Trampoline[FactState, int]] {
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
		assert.Equal(t, 120, result)
	})

	t.Run("factorial of 0", func(t *testing.T) {
		result := factorial(FactState{n: 0, result: 1})()
		assert.Equal(t, 1, result)
	})

	t.Run("factorial of 1", func(t *testing.T) {
		result := factorial(FactState{n: 1, result: 1})()
		assert.Equal(t, 1, result)
	})

	t.Run("factorial of 10", func(t *testing.T) {
		result := factorial(FactState{n: 10, result: 1})()
		assert.Equal(t, 3628800, result)
	})

	t.Run("factorial of 12", func(t *testing.T) {
		result := factorial(FactState{n: 12, result: 1})()
		assert.Equal(t, 479001600, result)
	})
}

// TestTailRec_Fibonacci tests computing Fibonacci numbers using tail recursion
func TestTailRec_Fibonacci(t *testing.T) {
	type FibState struct {
		n    int
		prev int
		curr int
	}

	fibonacci := TailRec(func(state FibState) IO[TR.Trampoline[FibState, int]] {
		if state.n == 0 {
			return Of(TR.Land[FibState](state.curr))
		}
		return Of(TR.Bounce[int](FibState{
			n:    state.n - 1,
			prev: state.curr,
			curr: state.prev + state.curr,
		}))
	})

	t.Run("fibonacci of 0", func(t *testing.T) {
		result := fibonacci(FibState{n: 0, prev: 0, curr: 1})()
		assert.Equal(t, 1, result)
	})

	t.Run("fibonacci of 1", func(t *testing.T) {
		result := fibonacci(FibState{n: 1, prev: 0, curr: 1})()
		assert.Equal(t, 1, result)
	})

	t.Run("fibonacci of 10", func(t *testing.T) {
		result := fibonacci(FibState{n: 10, prev: 0, curr: 1})()
		assert.Equal(t, 89, result)
	})

	t.Run("fibonacci of 20", func(t *testing.T) {
		result := fibonacci(FibState{n: 20, prev: 0, curr: 1})()
		assert.Equal(t, 10946, result)
	})
}

// TestTailRec_SumList tests summing a list with tail recursion
func TestTailRec_SumList(t *testing.T) {
	type SumState struct {
		items []int
		sum   int
	}

	sumList := TailRec(func(state SumState) IO[TR.Trampoline[SumState, int]] {
		if A.IsEmpty(state.items) {
			return Of(TR.Land[SumState](state.sum))
		}
		return Of(TR.Bounce[int](SumState{
			items: state.items[1:],
			sum:   state.sum + state.items[0],
		}))
	})

	t.Run("sum empty list", func(t *testing.T) {
		result := sumList(SumState{items: []int{}, sum: 0})()
		assert.Equal(t, 0, result)
	})

	t.Run("sum single element", func(t *testing.T) {
		result := sumList(SumState{items: []int{42}, sum: 0})()
		assert.Equal(t, 42, result)
	})

	t.Run("sum multiple elements", func(t *testing.T) {
		result := sumList(SumState{items: []int{1, 2, 3, 4, 5}, sum: 0})()
		assert.Equal(t, 15, result)
	})

	t.Run("sum with negative numbers", func(t *testing.T) {
		result := sumList(SumState{items: []int{-1, 2, -3, 4, -5}, sum: 0})()
		assert.Equal(t, -3, result)
	})
}

// TestTailRec_Countdown tests a simple countdown
func TestTailRec_Countdown(t *testing.T) {
	countdown := TailRec(func(n int) IO[TR.Trampoline[int, string]] {
		if n <= 0 {
			return Of(TR.Land[int]("Done!"))
		}
		return Of(TR.Bounce[string](n - 1))
	})

	t.Run("countdown from 5", func(t *testing.T) {
		result := countdown(5)()
		assert.Equal(t, "Done!", result)
	})

	t.Run("countdown from 0", func(t *testing.T) {
		result := countdown(0)()
		assert.Equal(t, "Done!", result)
	})

	t.Run("countdown from negative", func(t *testing.T) {
		result := countdown(-5)()
		assert.Equal(t, "Done!", result)
	})

	t.Run("countdown from 100", func(t *testing.T) {
		result := countdown(100)()
		assert.Equal(t, "Done!", result)
	})
}

// TestTailRec_StackSafety tests that TailRec doesn't overflow the stack with large iterations
func TestTailRec_StackSafety(t *testing.T) {
	// Count down from a large number - this would overflow the stack with regular recursion
	largeCountdown := TailRec(func(n int) IO[TR.Trampoline[int, int]] {
		if n <= 0 {
			return Of(TR.Land[int](0))
		}
		return Of(TR.Bounce[int](n - 1))
	})

	t.Run("large iteration count", func(t *testing.T) {
		// This should complete without stack overflow
		result := largeCountdown(10000)()
		assert.Equal(t, 0, result)
	})

	t.Run("very large iteration count", func(t *testing.T) {
		// Even larger - would definitely overflow with regular recursion
		result := largeCountdown(100000)()
		assert.Equal(t, 0, result)
	})
}

// TestTailRec_SumToN tests summing numbers from 1 to N
func TestTailRec_SumToN(t *testing.T) {
	type SumState struct {
		current int
		limit   int
		sum     int
	}

	sumToN := TailRec(func(state SumState) IO[TR.Trampoline[SumState, int]] {
		if state.current > state.limit {
			return Of(TR.Land[SumState](state.sum))
		}
		return Of(TR.Bounce[int](SumState{
			current: state.current + 1,
			limit:   state.limit,
			sum:     state.sum + state.current,
		}))
	})

	t.Run("sum to 10", func(t *testing.T) {
		result := sumToN(SumState{current: 1, limit: 10, sum: 0})()
		assert.Equal(t, 55, result) // 1+2+3+4+5+6+7+8+9+10 = 55
	})

	t.Run("sum to 100", func(t *testing.T) {
		result := sumToN(SumState{current: 1, limit: 100, sum: 0})()
		assert.Equal(t, 5050, result) // n*(n+1)/2 = 100*101/2 = 5050
	})

	t.Run("sum to 0", func(t *testing.T) {
		result := sumToN(SumState{current: 1, limit: 0, sum: 0})()
		assert.Equal(t, 0, result)
	})

	t.Run("sum to 1", func(t *testing.T) {
		result := sumToN(SumState{current: 1, limit: 1, sum: 0})()
		assert.Equal(t, 1, result)
	})
}

// TestTailRec_DoubleList tests doubling all elements in a list
func TestTailRec_DoubleList(t *testing.T) {
	type ListState struct {
		items []int
		acc   []int
	}

	doubleAll := TailRec(func(state ListState) IO[TR.Trampoline[ListState, []int]] {
		if A.IsEmpty(state.items) {
			return Of(TR.Land[ListState](state.acc))
		}
		doubled := append(state.acc, state.items[0]*2)
		return Of(TR.Bounce[[]int](ListState{
			items: state.items[1:],
			acc:   doubled,
		}))
	})

	t.Run("double empty list", func(t *testing.T) {
		result := doubleAll(ListState{items: []int{}, acc: []int{}})()
		assert.Equal(t, []int{}, result)
	})

	t.Run("double single element", func(t *testing.T) {
		result := doubleAll(ListState{items: []int{5}, acc: []int{}})()
		assert.Equal(t, []int{10}, result)
	})

	t.Run("double multiple elements", func(t *testing.T) {
		result := doubleAll(ListState{items: []int{1, 2, 3}, acc: []int{}})()
		assert.Equal(t, []int{2, 4, 6}, result)
	})

	t.Run("double with negative numbers", func(t *testing.T) {
		result := doubleAll(ListState{items: []int{-1, 0, 1}, acc: []int{}})()
		assert.Equal(t, []int{-2, 0, 2}, result)
	})
}

// TestTailRec_ReverseList tests reversing a list using tail recursion
func TestTailRec_ReverseList(t *testing.T) {
	type ReverseState struct {
		items []string
		acc   []string
	}

	reverseList := TailRec(func(state ReverseState) IO[TR.Trampoline[ReverseState, []string]] {
		if A.IsEmpty(state.items) {
			return Of(TR.Land[ReverseState](state.acc))
		}
		return Of(TR.Bounce[[]string](ReverseState{
			items: state.items[1:],
			acc:   append([]string{state.items[0]}, state.acc...),
		}))
	})

	t.Run("reverse empty list", func(t *testing.T) {
		result := reverseList(ReverseState{items: []string{}, acc: []string{}})()
		assert.Equal(t, []string{}, result)
	})

	t.Run("reverse single element", func(t *testing.T) {
		result := reverseList(ReverseState{items: []string{"a"}, acc: []string{}})()
		assert.Equal(t, []string{"a"}, result)
	})

	t.Run("reverse multiple elements", func(t *testing.T) {
		result := reverseList(ReverseState{items: []string{"a", "b", "c"}, acc: []string{}})()
		assert.Equal(t, []string{"c", "b", "a"}, result)
	})
}

// TestTailRec_Power tests computing power using tail recursion
func TestTailRec_Power(t *testing.T) {
	type PowerState struct {
		base   int
		exp    int
		result int
	}

	power := TailRec(func(state PowerState) IO[TR.Trampoline[PowerState, int]] {
		if state.exp == 0 {
			return Of(TR.Land[PowerState](state.result))
		}
		return Of(TR.Bounce[int](PowerState{
			base:   state.base,
			exp:    state.exp - 1,
			result: state.result * state.base,
		}))
	})

	t.Run("2^0", func(t *testing.T) {
		result := power(PowerState{base: 2, exp: 0, result: 1})()
		assert.Equal(t, 1, result)
	})

	t.Run("2^3", func(t *testing.T) {
		result := power(PowerState{base: 2, exp: 3, result: 1})()
		assert.Equal(t, 8, result)
	})

	t.Run("3^4", func(t *testing.T) {
		result := power(PowerState{base: 3, exp: 4, result: 1})()
		assert.Equal(t, 81, result)
	})

	t.Run("5^5", func(t *testing.T) {
		result := power(PowerState{base: 5, exp: 5, result: 1})()
		assert.Equal(t, 3125, result)
	})
}

// TestTailRec_GCD tests computing greatest common divisor using Euclidean algorithm
func TestTailRec_GCD(t *testing.T) {
	type GCDState struct {
		a int
		b int
	}

	gcd := TailRec(func(state GCDState) IO[TR.Trampoline[GCDState, int]] {
		if state.b == 0 {
			return Of(TR.Land[GCDState](state.a))
		}
		return Of(TR.Bounce[int](GCDState{
			a: state.b,
			b: state.a % state.b,
		}))
	})

	t.Run("gcd(48, 18)", func(t *testing.T) {
		result := gcd(GCDState{a: 48, b: 18})()
		assert.Equal(t, 6, result)
	})

	t.Run("gcd(100, 50)", func(t *testing.T) {
		result := gcd(GCDState{a: 100, b: 50})()
		assert.Equal(t, 50, result)
	})

	t.Run("gcd(17, 19)", func(t *testing.T) {
		result := gcd(GCDState{a: 17, b: 19})()
		assert.Equal(t, 1, result) // coprime numbers
	})

	t.Run("gcd(1071, 462)", func(t *testing.T) {
		result := gcd(GCDState{a: 1071, b: 462})()
		assert.Equal(t, 21, result)
	})
}

// TestTailRec_CountOccurrences tests counting occurrences of a value in a list
func TestTailRec_CountOccurrences(t *testing.T) {
	type CountState struct {
		items  []int
		target int
		count  int
	}

	countOccurrences := TailRec(func(state CountState) IO[TR.Trampoline[CountState, int]] {
		if A.IsEmpty(state.items) {
			return Of(TR.Land[CountState](state.count))
		}
		newCount := state.count
		if state.items[0] == state.target {
			newCount++
		}
		return Of(TR.Bounce[int](CountState{
			items:  state.items[1:],
			target: state.target,
			count:  newCount,
		}))
	})

	t.Run("count in empty list", func(t *testing.T) {
		result := countOccurrences(CountState{items: []int{}, target: 5, count: 0})()
		assert.Equal(t, 0, result)
	})

	t.Run("count with no matches", func(t *testing.T) {
		result := countOccurrences(CountState{items: []int{1, 2, 3}, target: 5, count: 0})()
		assert.Equal(t, 0, result)
	})

	t.Run("count with single match", func(t *testing.T) {
		result := countOccurrences(CountState{items: []int{1, 2, 3, 4, 5}, target: 3, count: 0})()
		assert.Equal(t, 1, result)
	})

	t.Run("count with multiple matches", func(t *testing.T) {
		result := countOccurrences(CountState{items: []int{1, 2, 2, 3, 2, 4}, target: 2, count: 0})()
		assert.Equal(t, 3, result)
	})

	t.Run("count all same", func(t *testing.T) {
		result := countOccurrences(CountState{items: []int{5, 5, 5, 5}, target: 5, count: 0})()
		assert.Equal(t, 4, result)
	})
}

// TestTailRec_ImmediateTermination tests that immediate termination works correctly
func TestTailRec_ImmediateTermination(t *testing.T) {
	immediate := TailRec(func(n int) IO[TR.Trampoline[int, string]] {
		return Of(TR.Land[int]("immediate"))
	})

	result := immediate(42)()
	assert.Equal(t, "immediate", result)
}

// TestTailRec_SingleBounce tests a single bounce before landing
func TestTailRec_SingleBounce(t *testing.T) {
	singleBounce := TailRec(func(n int) IO[TR.Trampoline[int, int]] {
		if n == 0 {
			return Of(TR.Land[int](100))
		}
		return Of(TR.Bounce[int](0))
	})

	result := singleBounce(1)()
	assert.Equal(t, 100, result)
}
