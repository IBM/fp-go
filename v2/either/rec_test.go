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

package either

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// TestTailRecFactorial tests factorial computation with error handling
func TestTailRecFactorial(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	factorialStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.n < 0 {
			return Left[TR.Trampoline[State, int]]("negative input not allowed")
		}
		if state.n <= 1 {
			return Right[string](TR.Land[State](state.acc))
		}
		return Right[string](TR.Bounce[int](State{state.n - 1, state.acc * state.n}))
	}

	factorial := TailRec(factorialStep)

	// Test successful computation
	result := factorial(State{5, 1})
	assert.Equal(t, Of[string](120), result)

	// Test base case
	result = factorial(State{0, 1})
	assert.Equal(t, Of[string](1), result)

	// Test error case
	result = factorial(State{-1, 1})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Equal(t, "negative input not allowed", err)
}

// TestTailRecFibonacci tests Fibonacci computation with validation
func TestTailRecFibonacci(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	fibStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.n < 0 {
			return Left[TR.Trampoline[State, int]]("negative index")
		}
		if state.curr > 1000 {
			return Left[TR.Trampoline[State, int]](fmt.Sprintf("value too large: %d", state.curr))
		}
		if state.n <= 0 {
			return Right[string](TR.Land[State](state.curr))
		}
		return Right[string](TR.Bounce[int](State{state.n - 1, state.curr, state.prev + state.curr}))
	}

	fib := TailRec(fibStep)

	// Test successful computation
	result := fib(State{10, 0, 1})
	assert.Equal(t, Of[string](89), result) // 10th Fibonacci number

	// Test base case
	result = fib(State{0, 0, 1})
	assert.Equal(t, Of[string](1), result)

	// Test error case - negative
	result = fib(State{-1, 0, 1})
	assert.True(t, IsLeft(result))

	// Test error case - value too large
	result = fib(State{20, 0, 1})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "value too large")
}

// TestTailRecCountdown tests countdown with validation
func TestTailRecCountdown(t *testing.T) {
	countdownStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		if n < 0 {
			return Left[TR.Trampoline[int, int]]("already negative")
		}
		if n == 0 {
			return Right[string](TR.Land[int](0))
		}
		return Right[string](TR.Bounce[int](n - 1))
	}

	countdown := TailRec(countdownStep)

	// Test successful countdown
	result := countdown(10)
	assert.Equal(t, Of[string](0), result)

	// Test immediate termination
	result = countdown(0)
	assert.Equal(t, Of[string](0), result)

	// Test error case
	result = countdown(-5)
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Equal(t, "already negative", err)
}

// TestTailRecSumList tests summing a list with error handling
func TestTailRecSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	sumStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.sum > 100 {
			return Left[TR.Trampoline[State, int]](fmt.Sprintf("sum exceeds limit: %d", state.sum))
		}
		if A.IsEmpty(state.list) {
			return Right[string](TR.Land[State](state.sum))
		}
		return Right[string](TR.Bounce[int](State{state.list[1:], state.sum + state.list[0]}))
	}

	sumList := TailRec(sumStep)

	// Test successful sum
	result := sumList(State{[]int{1, 2, 3, 4, 5}, 0})
	assert.Equal(t, Of[string](15), result)

	// Test empty list
	result = sumList(State{[]int{}, 0})
	assert.Equal(t, Of[string](0), result)

	// Test error case - sum too large
	result = sumList(State{[]int{50, 60}, 0})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "sum exceeds limit")
}

// TestTailRecImmediateTermination tests immediate termination (Land on first call)
func TestTailRecImmediateTermination(t *testing.T) {
	immediateStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		return Right[string](TR.Land[int](n * 2))
	}

	immediate := TailRec(immediateStep)
	result := immediate(21)

	assert.Equal(t, Of[string](42), result)
}

// TestTailRecImmediateError tests immediate error (Left on first call)
func TestTailRecImmediateError(t *testing.T) {
	immediateErrorStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		return Left[TR.Trampoline[int, int]]("immediate error")
	}

	immediateError := TailRec(immediateErrorStep)
	result := immediateError(42)

	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Equal(t, "immediate error", err)
}

// TestTailRecStackSafety tests that TailRec handles large iterations without stack overflow
func TestTailRecStackSafety(t *testing.T) {
	countdownStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		if n <= 0 {
			return Right[string](TR.Land[int](n))
		}
		return Right[string](TR.Bounce[int](n - 1))
	}

	countdown := TailRec(countdownStep)
	result := countdown(10000)

	assert.Equal(t, Of[string](0), result)
}

// TestTailRecFindInRange tests finding a value in a range
func TestTailRecFindInRange(t *testing.T) {
	type State struct {
		current int
		max     int
		target  int
	}

	findStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.current > 1000 {
			return Left[TR.Trampoline[State, int]]("search exceeded maximum iterations")
		}
		if state.current >= state.max {
			return Right[string](TR.Land[State](-1)) // Not found
		}
		if state.current == state.target {
			return Right[string](TR.Land[State](state.current)) // Found
		}
		return Right[string](TR.Bounce[int](State{state.current + 1, state.max, state.target}))
	}

	find := TailRec(findStep)

	// Test found
	result := find(State{0, 100, 42})
	assert.Equal(t, Of[string](42), result)

	// Test not found
	result = find(State{0, 100, 200})
	assert.Equal(t, Of[string](-1), result)

	// Test error - exceeded iterations
	result = find(State{0, 2000, 1500})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "exceeded maximum")
}

// TestTailRecCollatzConjecture tests the Collatz conjecture
func TestTailRecCollatzConjecture(t *testing.T) {
	collatzStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		if n <= 0 {
			return Left[TR.Trampoline[int, int]]("invalid input: must be positive")
		}
		if n == 1 {
			return Right[string](TR.Land[int](1))
		}
		if n%2 == 0 {
			return Right[string](TR.Bounce[int](n / 2))
		}
		return Right[string](TR.Bounce[int](3*n + 1))
	}

	collatz := TailRec(collatzStep)

	// Test various starting points
	result := collatz(10)
	assert.Equal(t, Of[string](1), result)

	result = collatz(27)
	assert.Equal(t, Of[string](1), result)

	// Test error case
	result = collatz(0)
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "invalid input")
}

// TestTailRecGCD tests greatest common divisor computation
func TestTailRecGCD(t *testing.T) {
	type State struct {
		a int
		b int
	}

	gcdStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.a < 0 || state.b < 0 {
			return Left[TR.Trampoline[State, int]]("negative values not allowed")
		}
		if state.b == 0 {
			return Right[string](TR.Land[State](state.a))
		}
		return Right[string](TR.Bounce[int](State{state.b, state.a % state.b}))
	}

	gcd := TailRec(gcdStep)

	// Test successful GCD
	result := gcd(State{48, 18})
	assert.Equal(t, Of[string](6), result)

	result = gcd(State{100, 35})
	assert.Equal(t, Of[string](5), result)

	// Test error case
	result = gcd(State{-10, 5})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "negative values")
}

// TestTailRecPowerOfTwo tests computing powers of 2
func TestTailRecPowerOfTwo(t *testing.T) {
	type State struct {
		exponent int
		result   int
		target   int
	}

	powerStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.target < 0 {
			return Left[TR.Trampoline[State, int]]("negative exponent not supported")
		}
		if state.exponent >= state.target {
			return Right[string](TR.Land[State](state.result))
		}
		return Right[string](TR.Bounce[int](State{state.exponent + 1, state.result * 2, state.target}))
	}

	power := TailRec(powerStep)

	// Test 2^10
	result := power(State{0, 1, 10})
	assert.Equal(t, Of[string](1024), result)

	// Test 2^0
	result = power(State{0, 1, 0})
	assert.Equal(t, Of[string](1), result)

	// Test error case
	result = power(State{0, 1, -1})
	assert.True(t, IsLeft(result))
}

// TestTailRecErrorInMiddle tests error occurring in the middle of recursion
func TestTailRecErrorInMiddle(t *testing.T) {
	countdownStep := func(n int) Either[string, TR.Trampoline[int, int]] {
		if n == 5 {
			return Left[TR.Trampoline[int, int]]("error at 5")
		}
		if n <= 0 {
			return Right[string](TR.Land[int](n))
		}
		return Right[string](TR.Bounce[int](n - 1))
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)

	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Equal(t, "error at 5", err)
}

// TestTailRecMultipleErrorConditions tests multiple error conditions
func TestTailRecMultipleErrorConditions(t *testing.T) {
	type State struct {
		value int
		steps int
	}

	step := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.steps > 100 {
			return Left[TR.Trampoline[State, int]]("too many steps")
		}
		if state.value < 0 {
			return Left[TR.Trampoline[State, int]]("negative value encountered")
		}
		if state.value == 0 {
			return Right[string](TR.Land[State](state.steps))
		}
		return Right[string](TR.Bounce[int](State{state.value - 1, state.steps + 1}))
	}

	counter := TailRec(step)

	// Test successful case
	result := counter(State{10, 0})
	assert.Equal(t, Of[string](10), result)

	// Test too many steps error
	result = counter(State{200, 0})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "too many steps")
}

// TestTailRecWithComplexState tests recursion with complex state
func TestTailRecWithComplexState(t *testing.T) {
	type State struct {
		numbers []int
		sum     int
		product int
	}

	processStep := func(state State) Either[string, TR.Trampoline[State, State]] {
		if state.product > 10000 {
			return Left[TR.Trampoline[State, State]]("product overflow")
		}
		if A.IsEmpty(state.numbers) {
			return Right[string](TR.Land[State](state))
		}
		head := state.numbers[0]
		tail := state.numbers[1:]
		return Right[string](TR.Bounce[State](State{
			numbers: tail,
			sum:     state.sum + head,
			product: state.product * head,
		}))
	}

	process := TailRec(processStep)

	// Test successful processing
	result := process(State{[]int{2, 3, 4}, 0, 1})
	assert.True(t, IsRight(result))
	finalState, _ := Unwrap(result)
	assert.Equal(t, 9, finalState.sum)
	assert.Equal(t, 24, finalState.product)

	// Test overflow error
	result = process(State{[]int{100, 200, 300}, 0, 1})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Contains(t, err, "product overflow")
}

// TestTailRecDivisionByZeroProtection tests protection against division by zero
func TestTailRecDivisionByZeroProtection(t *testing.T) {
	type State struct {
		numerator   int
		denominator int
		result      int
	}

	divideStep := func(state State) Either[string, TR.Trampoline[State, int]] {
		if state.denominator == 0 {
			return Left[TR.Trampoline[State, int]]("division by zero")
		}
		if state.numerator < state.denominator {
			return Right[string](TR.Land[State](state.result))
		}
		return Right[string](TR.Bounce[int](State{
			numerator:   state.numerator - state.denominator,
			denominator: state.denominator,
			result:      state.result + 1,
		}))
	}

	divide := TailRec(divideStep)

	// Test successful division
	result := divide(State{10, 3, 0})
	assert.Equal(t, Of[string](3), result) // 10 / 3 = 3 (integer division)

	// Test division by zero
	result = divide(State{10, 0, 0})
	assert.True(t, IsLeft(result))
	_, err := Unwrap(result)
	assert.Equal(t, "division by zero", err)
}

// TestTailRecStringProcessing tests recursion with string processing
func TestTailRecStringProcessing(t *testing.T) {
	type State struct {
		remaining string
		count     int
	}

	countVowels := func(state State) Either[string, TR.Trampoline[State, int]] {
		if len(state.remaining) == 0 {
			return Right[string](TR.Land[State](state.count))
		}
		char := state.remaining[0]
		isVowel := char == 'a' || char == 'e' || char == 'i' || char == 'o' || char == 'u' ||
			char == 'A' || char == 'E' || char == 'I' || char == 'O' || char == 'U'
		newCount := state.count
		if isVowel {
			newCount++
		}
		return Right[string](TR.Bounce[int](State{state.remaining[1:], newCount}))
	}

	counter := TailRec(countVowels)

	result := counter(State{"hello world", 0})
	assert.Equal(t, Of[string](3), result) // e, o, o
}
