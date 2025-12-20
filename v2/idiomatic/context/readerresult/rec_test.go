// Copyright (c) 2025 IBM Corp.
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
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// TestTailRecFactorial tests factorial computation with context
func TestTailRecFactorial(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	factorialStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.n <= 0 {
				return TR.Land[State](state.acc), nil
			}
			return TR.Bounce[int](State{state.n - 1, state.acc * state.n}), nil
		}
	}

	factorial := TailRec(factorialStep)
	result, err := factorial(State{5, 1})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 120, result)
}

// TestTailRecFibonacci tests Fibonacci computation
func TestTailRecFibonacci(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	fibStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.n <= 0 {
				return TR.Land[State](state.curr), nil
			}
			return TR.Bounce[int](State{state.n - 1, state.curr, state.prev + state.curr}), nil
		}
	}

	fib := TailRec(fibStep)
	result, err := fib(State{10, 0, 1})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 89, result) // 10th Fibonacci number
}

// TestTailRecCountdown tests countdown computation
func TestTailRecCountdown(t *testing.T) {
	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(10)(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}

// TestTailRecImmediateTermination tests immediate termination (Land on first call)
func TestTailRecImmediateTermination(t *testing.T) {
	immediateStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			return TR.Land[int](n * 2), nil
		}
	}

	immediate := TailRec(immediateStep)
	result, err := immediate(42)(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 84, result)
}

// TestTailRecStackSafety tests that TailRec handles large iterations without stack overflow
func TestTailRecStackSafety(t *testing.T) {
	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(10000)(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}

// TestTailRecSumList tests summing a list
func TestTailRecSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	sumStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if A.IsEmpty(state.list) {
				return TR.Land[State](state.sum), nil
			}
			return TR.Bounce[int](State{state.list[1:], state.sum + state.list[0]}), nil
		}
	}

	sumList := TailRec(sumStep)
	result, err := sumList(State{[]int{1, 2, 3, 4, 5}, 0})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

// TestTailRecCollatzConjecture tests the Collatz conjecture
func TestTailRecCollatzConjecture(t *testing.T) {
	collatzStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n <= 1 {
				return TR.Land[int](n), nil
			}
			if n%2 == 0 {
				return TR.Bounce[int](n / 2), nil
			}
			return TR.Bounce[int](3*n + 1), nil
		}
	}

	collatz := TailRec(collatzStep)
	result, err := collatz(10)(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 1, result)
}

// TestTailRecGCD tests greatest common divisor
func TestTailRecGCD(t *testing.T) {
	type State struct {
		a int
		b int
	}

	gcdStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.b == 0 {
				return TR.Land[State](state.a), nil
			}
			return TR.Bounce[int](State{state.b, state.a % state.b}), nil
		}
	}

	gcd := TailRec(gcdStep)
	result, err := gcd(State{48, 18})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 6, result)
}

// TestTailRecErrorPropagation tests that errors are properly propagated
func TestTailRecErrorPropagation(t *testing.T) {
	expectedErr := errors.New("computation error")

	errorStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n == 5 {
				return TR.Trampoline[int, int]{}, expectedErr
			}
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	computation := TailRec(errorStep)
	result, err := computation(10)(context.Background())

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextCancellationImmediate tests short circuit when context is already canceled
func TestTailRecContextCancellationImmediate(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately before execution

	stepExecuted := false
	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			stepExecuted = true
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(10)(ctx)

	// Should short circuit without executing any steps
	assert.False(t, stepExecuted, "Step should not be executed when context is already canceled")
	assert.Error(t, err)
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextCancellationDuringExecution tests short circuit when context is canceled during execution
func TestTailRecContextCancellationDuringExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	executionCount := 0
	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			executionCount++
			// Cancel after 3 iterations
			if executionCount == 3 {
				cancel()
			}
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(100)(ctx)

	// Should stop after cancellation
	assert.Error(t, err)
	assert.LessOrEqual(t, executionCount, 4, "Should stop shortly after cancellation")
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextWithTimeout tests behavior with timeout context
func TestTailRecContextWithTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	executionCount := 0
	slowStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			executionCount++
			// Simulate slow computation
			time.Sleep(20 * time.Millisecond)
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	computation := TailRec(slowStep)
	result, err := computation(100)(ctx)

	// Should timeout and return error
	assert.Error(t, err)
	assert.Less(t, executionCount, 100, "Should not complete all iterations due to timeout")
	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextWithCause tests that context.Cause is properly returned
func TestTailRecContextWithCause(t *testing.T) {
	customErr := errors.New("custom cancellation reason")
	ctx, cancel := context.WithCancelCause(context.Background())
	cancel(customErr)

	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(10)(ctx)

	assert.Error(t, err)
	assert.Equal(t, customErr, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextCancellationMultipleIterations tests that cancellation is checked on each iteration
func TestTailRecContextCancellationMultipleIterations(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	executionCount := 0
	maxExecutions := 5

	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			executionCount++
			if executionCount == maxExecutions {
				cancel()
			}
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(1000)(ctx)

	// Should detect cancellation on next iteration check
	assert.Error(t, err)
	// Should stop within 1-2 iterations after cancellation
	assert.LessOrEqual(t, executionCount, maxExecutions+2)
	assert.Equal(t, context.Canceled, err)
	assert.Equal(t, 0, result) // zero value
}

// TestTailRecContextNotCanceled tests normal execution when context is not canceled
func TestTailRecContextNotCanceled(t *testing.T) {
	ctx := context.Background()

	executionCount := 0
	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			executionCount++
			if n <= 0 {
				return TR.Land[int](n), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(10)(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 11, executionCount) // 10, 9, 8, ..., 1, 0
	assert.Equal(t, 0, result)
}

// TestTailRecPowerOfTwo tests computing power of 2
func TestTailRecPowerOfTwo(t *testing.T) {
	type State struct {
		exponent int
		result   int
		target   int
	}

	powerStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.exponent >= state.target {
				return TR.Land[State](state.result), nil
			}
			return TR.Bounce[int](State{state.exponent + 1, state.result * 2, state.target}), nil
		}
	}

	power := TailRec(powerStep)
	result, err := power(State{0, 1, 10})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 1024, result) // 2^10
}

// TestTailRecFindInRange tests finding a value in a range
func TestTailRecFindInRange(t *testing.T) {
	type State struct {
		current int
		max     int
		target  int
	}

	findStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.current >= state.max {
				return TR.Land[State](-1), nil // Not found
			}
			if state.current == state.target {
				return TR.Land[State](state.current), nil // Found
			}
			return TR.Bounce[int](State{state.current + 1, state.max, state.target}), nil
		}
	}

	find := TailRec(findStep)
	result, err := find(State{0, 100, 42})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

// TestTailRecFindNotInRange tests finding a value not in range
func TestTailRecFindNotInRange(t *testing.T) {
	type State struct {
		current int
		max     int
		target  int
	}

	findStep := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.current >= state.max {
				return TR.Land[State](-1), nil // Not found
			}
			if state.current == state.target {
				return TR.Land[State](state.current), nil // Found
			}
			return TR.Bounce[int](State{state.current + 1, state.max, state.target}), nil
		}
	}

	find := TailRec(findStep)
	result, err := find(State{0, 100, 200})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, -1, result)
}

// TestTailRecWithContextValue tests that context values are accessible
func TestTailRecWithContextValue(t *testing.T) {
	type contextKey string
	const multiplierKey contextKey = "multiplier"

	ctx := context.WithValue(context.Background(), multiplierKey, 3)

	countdownStep := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			if n <= 0 {
				multiplier := ctx.Value(multiplierKey).(int)
				return TR.Land[int](n * multiplier), nil
			}
			return TR.Bounce[int](n - 1), nil
		}
	}

	countdown := TailRec(countdownStep)
	result, err := countdown(5)(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 0, result) // 0 * 3 = 0
}

// TestTailRecComplexState tests with complex state structure
func TestTailRecComplexState(t *testing.T) {
	type ComplexState struct {
		counter   int
		sum       int
		product   int
		completed bool
	}

	complexStep := func(state ComplexState) ReaderResult[TR.Trampoline[ComplexState, string]] {
		return func(ctx context.Context) (TR.Trampoline[ComplexState, string], error) {
			if state.counter <= 0 || state.completed {
				result := fmt.Sprintf("sum=%d, product=%d", state.sum, state.product)
				return TR.Land[ComplexState](result), nil
			}
			newState := ComplexState{
				counter:   state.counter - 1,
				sum:       state.sum + state.counter,
				product:   state.product * state.counter,
				completed: state.counter == 1,
			}
			return TR.Bounce[string](newState), nil
		}
	}

	computation := TailRec(complexStep)
	result, err := computation(ComplexState{5, 0, 1, false})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "sum=15, product=120", result)
}

// TestTailRecZeroIterations tests when computation terminates immediately
func TestTailRecZeroIterations(t *testing.T) {
	step := func(n int) ReaderResult[TR.Trampoline[int, string]] {
		return func(ctx context.Context) (TR.Trampoline[int, string], error) {
			return TR.Land[int]("immediate"), nil
		}
	}

	computation := TailRec(step)
	result, err := computation(0)(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, "immediate", result)
}

// TestTailRecErrorInFirstIteration tests error on first iteration
func TestTailRecErrorInFirstIteration(t *testing.T) {
	expectedErr := errors.New("first iteration error")

	step := func(n int) ReaderResult[TR.Trampoline[int, int]] {
		return func(ctx context.Context) (TR.Trampoline[int, int], error) {
			return TR.Trampoline[int, int]{}, expectedErr
		}
	}

	computation := TailRec(step)
	result, err := computation(10)(context.Background())

	assert.Error(t, err)
	assert.Equal(t, expectedErr, err)
	assert.Equal(t, 0, result)
}

// TestTailRecAlternatingBounce tests alternating between different values
func TestTailRecAlternatingBounce(t *testing.T) {
	type State struct {
		value     int
		alternate bool
		count     int
	}

	step := func(state State) ReaderResult[TR.Trampoline[State, int]] {
		return func(ctx context.Context) (TR.Trampoline[State, int], error) {
			if state.count >= 10 {
				return TR.Land[State](state.value), nil
			}
			newValue := state.value
			if state.alternate {
				newValue += 1
			} else {
				newValue -= 1
			}
			return TR.Bounce[int](State{newValue, !state.alternate, state.count + 1}), nil
		}
	}

	computation := TailRec(step)
	result, err := computation(State{0, true, 0})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, 0, result) // Should alternate +1, -1 and end at 0
}

// TestTailRecLargeAccumulation tests accumulating large values
func TestTailRecLargeAccumulation(t *testing.T) {
	type State struct {
		n   int
		sum int64
	}

	step := func(state State) ReaderResult[TR.Trampoline[State, int64]] {
		return func(ctx context.Context) (TR.Trampoline[State, int64], error) {
			if state.n <= 0 {
				return TR.Land[State](state.sum), nil
			}
			return TR.Bounce[int64](State{state.n - 1, state.sum + int64(state.n)}), nil
		}
	}

	computation := TailRec(step)
	result, err := computation(State{1000, 0})(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, int64(500500), result) // Sum of 1 to 1000
}
