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

package tailrec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBounce verifies that Bounce creates a trampoline in the bounce state
func TestBounce(t *testing.T) {
	t.Run("creates bounce state with integer", func(t *testing.T) {
		tramp := Bounce[string](42)

		assert.False(t, tramp.Landed, "should be in bounce state")
		assert.Equal(t, 42, tramp.Bounce, "bounce value should match")
		assert.Equal(t, "", tramp.Land, "land value should be zero value")
	})

	t.Run("creates bounce state with string", func(t *testing.T) {
		tramp := Bounce[int]("hello")

		assert.False(t, tramp.Landed, "should be in bounce state")
		assert.Equal(t, "hello", tramp.Bounce, "bounce value should match")
		assert.Equal(t, 0, tramp.Land, "land value should be zero value")
	})

	t.Run("creates bounce state with struct", func(t *testing.T) {
		type State struct {
			n   int
			acc int
		}
		state := State{n: 5, acc: 10}
		tramp := Bounce[int](state)

		assert.False(t, tramp.Landed, "should be in bounce state")
		assert.Equal(t, state, tramp.Bounce, "bounce value should match")
		assert.Equal(t, 0, tramp.Land, "land value should be zero value")
	})
}

// TestLand verifies that Land creates a trampoline in the land state
func TestLand(t *testing.T) {
	t.Run("creates land state with integer", func(t *testing.T) {
		tramp := Land[string](42)

		assert.True(t, tramp.Landed, "should be in land state")
		assert.Equal(t, 42, tramp.Land, "land value should match")
		assert.Equal(t, "", tramp.Bounce, "bounce value should be zero value")
	})

	t.Run("creates land state with string", func(t *testing.T) {
		tramp := Land[int]("result")

		assert.True(t, tramp.Landed, "should be in land state")
		assert.Equal(t, "result", tramp.Land, "land value should match")
		assert.Equal(t, 0, tramp.Bounce, "bounce value should be zero value")
	})

	t.Run("creates land state with struct", func(t *testing.T) {
		type Result struct {
			value int
			done  bool
		}
		result := Result{value: 100, done: true}
		tramp := Land[int](result)

		assert.True(t, tramp.Landed, "should be in land state")
		assert.Equal(t, result, tramp.Land, "land value should match")
		assert.Equal(t, 0, tramp.Bounce, "bounce value should be zero value")
	})
}

// TestFieldAccess verifies that trampoline fields can be accessed directly
func TestFieldAccess(t *testing.T) {
	t.Run("accesses bounce state fields", func(t *testing.T) {
		tramp := Bounce[string](42)

		assert.False(t, tramp.Landed)
		assert.Equal(t, 42, tramp.Bounce)
		assert.Equal(t, "", tramp.Land)
	})

	t.Run("accesses land state fields", func(t *testing.T) {
		tramp := Land[int]("done")

		assert.True(t, tramp.Landed)
		assert.Equal(t, "done", tramp.Land)
		assert.Equal(t, 0, tramp.Bounce)
	})
}

// TestFactorial demonstrates a complete factorial implementation using trampolines
func TestFactorial(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	factorialStep := func(state State) Trampoline[State, int] {
		if state.n <= 1 {
			return Land[State](state.acc)
		}
		return Bounce[int](State{state.n - 1, state.acc * state.n})
	}

	factorial := func(n int) int {
		current := Bounce[int](State{n, 1})
		for {
			if current.Landed {
				return current.Land
			}
			current = factorialStep(current.Bounce)
		}
	}

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"factorial of 0", 0, 1},
		{"factorial of 1", 1, 1},
		{"factorial of 5", 5, 120},
		{"factorial of 10", 10, 3628800},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := factorial(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFibonacci demonstrates Fibonacci sequence using trampolines
func TestFibonacci(t *testing.T) {
	type State struct {
		n    int
		curr int
		prev int
	}

	fibStep := func(state State) Trampoline[State, int] {
		if state.n <= 0 {
			return Land[State](state.curr)
		}
		return Bounce[int](State{
			n:    state.n - 1,
			curr: state.prev + state.curr,
			prev: state.curr,
		})
	}

	fibonacci := func(n int) int {
		current := Bounce[int](State{n, 1, 0})
		for {
			if current.Landed {
				return current.Land
			}
			current = fibStep(current.Bounce)
		}
	}

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"fib(0)", 0, 1},
		{"fib(1)", 1, 1},
		{"fib(5)", 5, 8},
		{"fib(10)", 10, 89},
		{"fib(20)", 20, 10946},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := fibonacci(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestSumList demonstrates list processing using trampolines
func TestSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	sumStep := func(state State) Trampoline[State, int] {
		if len(state.list) == 0 {
			return Land[State](state.sum)
		}
		return Bounce[int](State{
			list: state.list[1:],
			sum:  state.sum + state.list[0],
		})
	}

	sumList := func(list []int) int {
		current := Bounce[int](State{list, 0})
		for {
			if current.Landed {
				return current.Land
			}
			current = sumStep(current.Bounce)
		}
	}

	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"empty list", []int{}, 0},
		{"single element", []int{42}, 42},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 15},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed numbers", []int{10, -5, 3, -2}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sumList(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestCountdown demonstrates a simple countdown using trampolines
func TestCountdown(t *testing.T) {
	countdownStep := func(n int) Trampoline[int, int] {
		if n <= 0 {
			return Land[int](0)
		}
		return Bounce[int](n - 1)
	}

	countdown := func(n int) int {
		current := Bounce[int](n)
		for {
			if current.Landed {
				return current.Land
			}
			current = countdownStep(current.Bounce)
		}
	}

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"countdown from 0", 0, 0},
		{"countdown from 1", 1, 0},
		{"countdown from 10", 10, 0},
		{"countdown from 1000", 1000, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := countdown(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestGCD demonstrates greatest common divisor using trampolines
func TestGCD(t *testing.T) {
	type State struct {
		a int
		b int
	}

	gcdStep := func(state State) Trampoline[State, int] {
		if state.b == 0 {
			return Land[State](state.a)
		}
		return Bounce[int](State{state.b, state.a % state.b})
	}

	gcd := func(a, b int) int {
		current := Bounce[int](State{a, b})
		for {
			if current.Landed {
				return current.Land
			}
			current = gcdStep(current.Bounce)
		}
	}

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"gcd(48, 18)", 48, 18, 6},
		{"gcd(100, 50)", 100, 50, 50},
		{"gcd(17, 13)", 17, 13, 1},
		{"gcd(0, 5)", 0, 5, 5},
		{"gcd(5, 0)", 5, 0, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gcd(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDeepRecursion verifies that trampolines can handle very deep recursion
// without stack overflow
func TestDeepRecursion(t *testing.T) {
	countdownStep := func(n int) Trampoline[int, int] {
		if n <= 0 {
			return Land[int](0)
		}
		return Bounce[int](n - 1)
	}

	countdown := func(n int) int {
		current := Bounce[int](n)
		for {
			if current.Landed {
				return current.Land
			}
			current = countdownStep(current.Bounce)
		}
	}

	// This would cause stack overflow with regular recursion
	result := countdown(100000)
	assert.Equal(t, 0, result, "should handle deep recursion without stack overflow")
}

// TestDifferentTypes verifies trampolines work with various type combinations
func TestDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		tramp := Land[int]("result")
		assert.True(t, tramp.Landed)
		assert.Equal(t, "result", tramp.Land)
	})

	t.Run("string to bool", func(t *testing.T) {
		tramp := Bounce[bool]("state")
		assert.False(t, tramp.Landed)
		assert.Equal(t, "state", tramp.Bounce)
	})

	t.Run("struct to struct", func(t *testing.T) {
		type Input struct{ x int }
		type Output struct{ y string }

		tramp := Land[Input](Output{y: "done"})
		assert.True(t, tramp.Landed)
		assert.Equal(t, Output{y: "done"}, tramp.Land)
	})

	t.Run("slice to map", func(t *testing.T) {
		tramp := Bounce[map[string]int]([]string{"a", "b"})
		assert.False(t, tramp.Landed)
		assert.Equal(t, []string{"a", "b"}, tramp.Bounce)
	})
}
