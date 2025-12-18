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

package ioeither

import (
	"errors"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// TestTailRecFactorial tests computing factorial using tail recursion
func TestTailRecFactorial(t *testing.T) {
	type FactState struct {
		n      int
		result int
	}

	factorial := TailRec(func(state FactState) IOEither[error, TR.Trampoline[FactState, int]] {
		if state.n <= 1 {
			// Terminate with final result
			return Of[error](TR.Land[FactState](state.result))
		}
		// Continue with next iteration
		return Of[error](TR.Bounce[int](FactState{
			n:      state.n - 1,
			result: state.result * state.n,
		}))
	})

	t.Run("factorial of 5", func(t *testing.T) {
		result := factorial(FactState{n: 5, result: 1})()
		assert.Equal(t, E.Right[error](120), result)
	})

	t.Run("factorial of 0", func(t *testing.T) {
		result := factorial(FactState{n: 0, result: 1})()
		assert.Equal(t, E.Right[error](1), result)
	})

	t.Run("factorial of 1", func(t *testing.T) {
		result := factorial(FactState{n: 1, result: 1})()
		assert.Equal(t, E.Right[error](1), result)
	})

	t.Run("factorial of 10", func(t *testing.T) {
		result := factorial(FactState{n: 10, result: 1})()
		assert.Equal(t, E.Right[error](3628800), result)
	})
}

// TestTailRecFibonacci tests computing Fibonacci numbers using tail recursion
func TestTailRecFibonacci(t *testing.T) {
	type FibState struct {
		n    int
		prev int
		curr int
	}

	fibonacci := TailRec(func(state FibState) IOEither[error, TR.Trampoline[FibState, int]] {
		if state.n == 0 {
			return Of[error](TR.Land[FibState](state.curr))
		}
		return Of[error](TR.Bounce[int](FibState{
			n:    state.n - 1,
			prev: state.curr,
			curr: state.prev + state.curr,
		}))
	})

	t.Run("fibonacci of 0", func(t *testing.T) {
		result := fibonacci(FibState{n: 0, prev: 0, curr: 1})()
		assert.Equal(t, E.Right[error](1), result)
	})

	t.Run("fibonacci of 1", func(t *testing.T) {
		result := fibonacci(FibState{n: 1, prev: 0, curr: 1})()
		assert.Equal(t, E.Right[error](1), result)
	})

	t.Run("fibonacci of 10", func(t *testing.T) {
		result := fibonacci(FibState{n: 10, prev: 0, curr: 1})()
		assert.Equal(t, E.Right[error](89), result)
	})
}

// TestTailRecSumList tests summing a list with tail recursion
func TestTailRecSumList(t *testing.T) {
	type SumState struct {
		items []int
		sum   int
	}

	sumList := TailRec(func(state SumState) IOEither[error, TR.Trampoline[SumState, int]] {
		if A.IsEmpty(state.items) {
			return Of[error](TR.Land[SumState](state.sum))
		}
		return Of[error](TR.Bounce[int](SumState{
			items: state.items[1:],
			sum:   state.sum + state.items[0],
		}))
	})

	t.Run("sum empty list", func(t *testing.T) {
		result := sumList(SumState{items: []int{}, sum: 0})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("sum single element", func(t *testing.T) {
		result := sumList(SumState{items: []int{42}, sum: 0})()
		assert.Equal(t, E.Right[error](42), result)
	})

	t.Run("sum multiple elements", func(t *testing.T) {
		result := sumList(SumState{items: []int{1, 2, 3, 4, 5}, sum: 0})()
		assert.Equal(t, E.Right[error](15), result)
	})
}

// TestTailRecWithError tests tail recursion that can fail
func TestTailRecWithError(t *testing.T) {
	type DivState struct {
		n      int
		result int
	}

	// Divide n by 2 repeatedly until it reaches 1, fail if we encounter an odd number > 1
	divideByTwo := TailRec(func(state DivState) IOEither[error, TR.Trampoline[DivState, int]] {
		if state.n == 1 {
			return Of[error](TR.Land[DivState](state.result))
		}
		if state.n%2 != 0 {
			return Left[TR.Trampoline[DivState, int]](fmt.Errorf("cannot divide odd number %d", state.n))
		}
		return Of[error](TR.Bounce[int](DivState{
			n:      state.n / 2,
			result: state.result + 1,
		}))
	})

	t.Run("success with power of 2", func(t *testing.T) {
		result := divideByTwo(DivState{n: 8, result: 0})()
		assert.Equal(t, E.Right[error](3), result) // 8 -> 4 -> 2 -> 1 (3 divisions)
	})

	t.Run("success with 1", func(t *testing.T) {
		result := divideByTwo(DivState{n: 1, result: 0})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("failure with odd number", func(t *testing.T) {
		result := divideByTwo(DivState{n: 5, result: 0})()
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot divide odd number 5")
	})

	t.Run("failure after some iterations", func(t *testing.T) {
		result := divideByTwo(DivState{n: 12, result: 0})()
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot divide odd number 3")
	})
}

// TestTailRecCountdown tests a simple countdown
func TestTailRecCountdown(t *testing.T) {
	countdown := TailRec(func(n int) IOEither[error, TR.Trampoline[int, string]] {
		if n <= 0 {
			return Of[error](TR.Land[int]("Done!"))
		}
		return Of[error](TR.Bounce[string](n - 1))
	})

	t.Run("countdown from 5", func(t *testing.T) {
		result := countdown(5)()
		assert.Equal(t, E.Right[error]("Done!"), result)
	})

	t.Run("countdown from 0", func(t *testing.T) {
		result := countdown(0)()
		assert.Equal(t, E.Right[error]("Done!"), result)
	})

	t.Run("countdown from negative", func(t *testing.T) {
		result := countdown(-5)()
		assert.Equal(t, E.Right[error]("Done!"), result)
	})
}

// TestTailRecStackSafety tests that TailRec doesn't overflow the stack with large iterations
func TestTailRecStackSafety(t *testing.T) {
	// Count down from a large number - this would overflow the stack with regular recursion
	largeCountdown := TailRec(func(n int) IOEither[error, TR.Trampoline[int, int]] {
		if n <= 0 {
			return Of[error](TR.Land[int](0))
		}
		return Of[error](TR.Bounce[int](n - 1))
	})

	t.Run("large iteration count", func(t *testing.T) {
		// This should complete without stack overflow
		result := largeCountdown(10000)()
		assert.Equal(t, E.Right[error](0), result)
	})
}

// TestTailRecFindInList tests searching for an element in a list
func TestTailRecFindInList(t *testing.T) {
	type FindState struct {
		items  []string
		target string
		index  int
	}

	findInList := TailRec(func(state FindState) IOEither[error, TR.Trampoline[FindState, int]] {
		if A.IsEmpty(state.items) {
			return Left[TR.Trampoline[FindState, int]](errors.New("not found"))
		}
		if state.items[0] == state.target {
			return Of[error](TR.Land[FindState](state.index))
		}
		return Of[error](TR.Bounce[int](FindState{
			items:  state.items[1:],
			target: state.target,
			index:  state.index + 1,
		}))
	})

	t.Run("find existing element", func(t *testing.T) {
		result := findInList(FindState{
			items:  []string{"a", "b", "c", "d"},
			target: "c",
			index:  0,
		})()
		assert.Equal(t, E.Right[error](2), result)
	})

	t.Run("find first element", func(t *testing.T) {
		result := findInList(FindState{
			items:  []string{"a", "b", "c"},
			target: "a",
			index:  0,
		})()
		assert.Equal(t, E.Right[error](0), result)
	})

	t.Run("element not found", func(t *testing.T) {
		result := findInList(FindState{
			items:  []string{"a", "b", "c"},
			target: "z",
			index:  0,
		})()
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Error(t, err)
		assert.Equal(t, "not found", err.Error())
	})

	t.Run("empty list", func(t *testing.T) {
		result := findInList(FindState{
			items:  []string{},
			target: "a",
			index:  0,
		})()
		assert.True(t, E.IsLeft(result))
	})
}
