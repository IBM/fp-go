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

package readerresult

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestTailRecFactorial tests factorial computation with context
func TestTailRecFactorial(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	factorialStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.n <= 0 {
				return R.Of(E.Right[State](state.acc))
			}
			return R.Of(E.Left[int](State{state.n - 1, state.acc * state.n}))
		}
	}

	factorial := TailRec(factorialStep)
	result := factorial(State{5, 1})(context.Background())

	assert.Equal(t, R.Of(120), result)
}

// TestTailRecFibonacci tests Fibonacci computation
func TestTailRecFibonacci(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	fibStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.n <= 0 {
				return R.Of(E.Right[State](state.curr))
			}
			return R.Of(E.Left[int](State{state.n - 1, state.curr, state.prev + state.curr}))
		}
	}

	fib := TailRec(fibStep)
	result := fib(State{10, 0, 1})(context.Background())

	assert.Equal(t, R.Of(89), result) // 10th Fibonacci number
}

// TestTailRecCountdown tests countdown computation
func TestTailRecCountdown(t *testing.T) {
	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(context.Background())

	assert.Equal(t, R.Of(0), result)
}

// TestTailRecImmediateTermination tests immediate termination (Right on first call)
func TestTailRecImmediateTermination(t *testing.T) {
	immediateStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			return R.Of(E.Right[int](n * 2))
		}
	}

	immediate := TailRec(immediateStep)
	result := immediate(42)(context.Background())

	assert.Equal(t, R.Of(84), result)
}

// TestTailRecStackSafety tests that TailRec handles large iterations without stack overflow
func TestTailRecStackSafety(t *testing.T) {
	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10000)(context.Background())

	assert.Equal(t, R.Of(0), result)
}

// TestTailRecSumList tests summing a list
func TestTailRecSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	sumStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if A.IsEmpty(state.list) {
				return R.Of(E.Right[State](state.sum))
			}
			return R.Of(E.Left[int](State{state.list[1:], state.sum + state.list[0]}))
		}
	}

	sumList := TailRec(sumStep)
	result := sumList(State{[]int{1, 2, 3, 4, 5}, 0})(context.Background())

	assert.Equal(t, R.Of(15), result)
}

// TestTailRecCollatzConjecture tests the Collatz conjecture
func TestTailRecCollatzConjecture(t *testing.T) {
	collatzStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n <= 1 {
				return R.Of(E.Right[int](n))
			}
			if n%2 == 0 {
				return R.Of(E.Left[int](n / 2))
			}
			return R.Of(E.Left[int](3*n + 1))
		}
	}

	collatz := TailRec(collatzStep)
	result := collatz(10)(context.Background())

	assert.Equal(t, R.Of(1), result)
}

// TestTailRecGCD tests greatest common divisor
func TestTailRecGCD(t *testing.T) {
	type State struct {
		a int
		b int
	}

	gcdStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.b == 0 {
				return R.Of(E.Right[State](state.a))
			}
			return R.Of(E.Left[int](State{state.b, state.a % state.b}))
		}
	}

	gcd := TailRec(gcdStep)
	result := gcd(State{48, 18})(context.Background())

	assert.Equal(t, R.Of(6), result)
}

// TestTailRecErrorPropagation tests that errors are properly propagated
func TestTailRecErrorPropagation(t *testing.T) {
	expectedErr := errors.New("computation error")

	errorStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n == 5 {
				return R.Left[E.Either[int, int]](expectedErr)
			}
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	computation := TailRec(errorStep)
	result := computation(10)(context.Background())

	assert.True(t, R.IsLeft(result))
	_, err := R.Unwrap(result)
	assert.Equal(t, expectedErr, err)
}

// TestTailRecContextCancellationImmediate tests short circuit when context is already canceled
func TestTailRecContextCancellationImmediate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately before execution

	stepExecuted := false
	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			stepExecuted = true
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(ctx)

	// Should short circuit without executing any steps
	assert.False(t, stepExecuted, "Step should not be executed when context is already canceled")
	assert.True(t, R.IsLeft(result))
	_, err := R.Unwrap(result)
	assert.Equal(t, context.Canceled, err)
}

// TestTailRecContextCancellationDuringExecution tests short circuit when context is canceled during execution
func TestTailRecContextCancellationDuringExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	executionCount := 0
	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			executionCount++
			// Cancel after 3 iterations
			if executionCount == 3 {
				cancel()
			}
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(100)(ctx)

	// Should stop after cancellation
	assert.True(t, R.IsLeft(result))
	assert.LessOrEqual(t, executionCount, 4, "Should stop shortly after cancellation")
	_, err := R.Unwrap(result)
	assert.Equal(t, context.Canceled, err)
}

// TestTailRecContextWithTimeout tests behavior with timeout context
func TestTailRecContextWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	executionCount := 0
	slowStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			executionCount++
			// Simulate slow computation
			time.Sleep(20 * time.Millisecond)
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	computation := TailRec(slowStep)
	result := computation(100)(ctx)

	// Should timeout and return error
	assert.True(t, R.IsLeft(result))
	assert.Less(t, executionCount, 100, "Should not complete all iterations due to timeout")
	_, err := R.Unwrap(result)
	assert.Equal(t, context.DeadlineExceeded, err)
}

// TestTailRecContextWithCause tests that context.Cause is properly returned
func TestTailRecContextWithCause(t *testing.T) {
	customErr := errors.New("custom cancellation reason")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(customErr)

	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(ctx)

	assert.True(t, R.IsLeft(result))
	_, err := R.Unwrap(result)
	assert.Equal(t, customErr, err)
}

// TestTailRecContextCancellationMultipleIterations tests that cancellation is checked on each iteration
func TestTailRecContextCancellationMultipleIterations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	executionCount := 0
	maxExecutions := 5

	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			executionCount++
			if executionCount == maxExecutions {
				cancel()
			}
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(1000)(ctx)

	// Should detect cancellation on next iteration check
	assert.True(t, R.IsLeft(result))
	// Should stop within 1-2 iterations after cancellation
	assert.LessOrEqual(t, executionCount, maxExecutions+2)
	_, err := R.Unwrap(result)
	assert.Equal(t, context.Canceled, err)
}

// TestTailRecContextNotCanceled tests normal execution when context is not canceled
func TestTailRecContextNotCanceled(t *testing.T) {
	ctx := context.Background()

	executionCount := 0
	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			executionCount++
			if n <= 0 {
				return R.Of(E.Right[int](n))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(ctx)

	assert.Equal(t, 11, executionCount) // 10, 9, 8, ..., 1, 0
	assert.Equal(t, R.Of(0), result)
}

// TestTailRecPowerOfTwo tests computing power of 2
func TestTailRecPowerOfTwo(t *testing.T) {
	type State struct {
		exponent int
		result   int
		target   int
	}

	powerStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.exponent >= state.target {
				return R.Of(E.Right[State](state.result))
			}
			return R.Of(E.Left[int](State{state.exponent + 1, state.result * 2, state.target}))
		}
	}

	power := TailRec(powerStep)
	result := power(State{0, 1, 10})(context.Background())

	assert.Equal(t, R.Of(1024), result) // 2^10
}

// TestTailRecFindInRange tests finding a value in a range
func TestTailRecFindInRange(t *testing.T) {
	type State struct {
		current int
		max     int
		target  int
	}

	findStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.current >= state.max {
				return R.Of(E.Right[State](-1)) // Not found
			}
			if state.current == state.target {
				return R.Of(E.Right[State](state.current)) // Found
			}
			return R.Of(E.Left[int](State{state.current + 1, state.max, state.target}))
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100, 42})(context.Background())

	assert.Equal(t, R.Of(42), result)
}

// TestTailRecFindNotInRange tests finding a value not in range
func TestTailRecFindNotInRange(t *testing.T) {
	type State struct {
		current int
		max     int
		target  int
	}

	findStep := func(state State) ReaderResult[E.Either[State, int]] {
		return func(ctx context.Context) Result[E.Either[State, int]] {
			if state.current >= state.max {
				return R.Of(E.Right[State](-1)) // Not found
			}
			if state.current == state.target {
				return R.Of(E.Right[State](state.current)) // Found
			}
			return R.Of(E.Left[int](State{state.current + 1, state.max, state.target}))
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100, 200})(context.Background())

	assert.Equal(t, R.Of(-1), result)
}

// TestTailRecWithContextValue tests that context values are accessible
func TestTailRecWithContextValue(t *testing.T) {
	type contextKey string
	const multiplierKey contextKey = "multiplier"

	ctx := context.WithValue(context.Background(), multiplierKey, 3)

	countdownStep := func(n int) ReaderResult[E.Either[int, int]] {
		return func(ctx context.Context) Result[E.Either[int, int]] {
			if n <= 0 {
				multiplier := ctx.Value(multiplierKey).(int)
				return R.Of(E.Right[int](n * multiplier))
			}
			return R.Of(E.Left[int](n - 1))
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(5)(ctx)

	assert.Equal(t, R.Of(0), result) // 0 * 3 = 0
}

// TestTailRecComplexState tests with complex state structure
func TestTailRecComplexState(t *testing.T) {
	type ComplexState struct {
		counter   int
		sum       int
		product   int
		completed bool
	}

	complexStep := func(state ComplexState) ReaderResult[E.Either[ComplexState, string]] {
		return func(ctx context.Context) Result[E.Either[ComplexState, string]] {
			if state.counter <= 0 || state.completed {
				result := fmt.Sprintf("sum=%d, product=%d", state.sum, state.product)
				return R.Of(E.Right[ComplexState](result))
			}
			newState := ComplexState{
				counter:   state.counter - 1,
				sum:       state.sum + state.counter,
				product:   state.product * state.counter,
				completed: state.counter == 1,
			}
			return R.Of(E.Left[string](newState))
		}
	}

	computation := TailRec(complexStep)
	result := computation(ComplexState{5, 0, 1, false})(context.Background())

	assert.Equal(t, R.Of("sum=15, product=120"), result)
}
