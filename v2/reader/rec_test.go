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

package reader

import (
	"testing"

	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// Test types for TailRec tests
type (
	// FactorialState represents the state for factorial computation
	FactorialState struct {
		n   int
		acc int
	}

	// CountdownState represents a simple countdown state
	CountdownState struct {
		count int
	}

	// SumState represents state for summing numbers
	SumState struct {
		current int
		total   int
	}

	// TestConfig is a simple environment for testing
	TestConfig struct {
		maxIterations int
		multiplier    int
	}
)

func TestTailRec_Factorial(t *testing.T) {
	t.Run("computes factorial correctly", func(t *testing.T) {
		// Define factorial step function
		factorialStep := func(state FactorialState) Reader[TestConfig, Trampoline[FactorialState, int]] {
			return func(cfg TestConfig) Trampoline[FactorialState, int] {
				if state.n <= 1 {
					return TR.Land[FactorialState](state.acc)
				}
				return TR.Bounce[int](FactorialState{
					n:   state.n - 1,
					acc: state.acc * state.n,
				})
			}
		}

		factorial := TailRec(factorialStep)
		config := TestConfig{maxIterations: 1000}

		// Test various factorial values
		testCases := []struct {
			n        int
			expected int
		}{
			{0, 1},
			{1, 1},
			{5, 120},
			{6, 720},
			{10, 3628800},
		}

		for _, tc := range testCases {
			result := factorial(FactorialState{n: tc.n, acc: 1})(config)
			assert.Equal(t, tc.expected, result, "factorial(%d) should equal %d", tc.n, tc.expected)
		}
	})

	t.Run("handles large recursion depth without stack overflow", func(t *testing.T) {
		factorialStep := func(state FactorialState) Reader[TestConfig, Trampoline[FactorialState, int]] {
			return func(cfg TestConfig) Trampoline[FactorialState, int] {
				if state.n <= 1 {
					return TR.Land[FactorialState](state.acc)
				}
				return TR.Bounce[int](FactorialState{
					n:   state.n - 1,
					acc: state.acc * state.n,
				})
			}
		}

		factorial := TailRec(factorialStep)
		config := TestConfig{maxIterations: 10000}

		// This would cause stack overflow with regular recursion
		// We're just testing it doesn't crash, not the actual value (which would overflow int)
		result := factorial(FactorialState{n: 1000, acc: 1})(config)
		assert.NotNil(t, result)
	})
}

func TestTailRec_Countdown(t *testing.T) {
	t.Run("counts down to zero", func(t *testing.T) {
		countdownStep := func(state CountdownState) Reader[TestConfig, Trampoline[CountdownState, string]] {
			return func(cfg TestConfig) Trampoline[CountdownState, string] {
				if state.count <= 0 {
					return TR.Land[CountdownState]("Done!")
				}
				return TR.Bounce[string](CountdownState{count: state.count - 1})
			}
		}

		countdown := TailRec(countdownStep)
		config := TestConfig{}

		result := countdown(CountdownState{count: 10})(config)
		assert.Equal(t, "Done!", result)
	})

	t.Run("handles immediate landing", func(t *testing.T) {
		countdownStep := func(state CountdownState) Reader[TestConfig, Trampoline[CountdownState, string]] {
			return func(cfg TestConfig) Trampoline[CountdownState, string] {
				if state.count <= 0 {
					return TR.Land[CountdownState]("Done!")
				}
				return TR.Bounce[string](CountdownState{count: state.count - 1})
			}
		}

		countdown := TailRec(countdownStep)
		config := TestConfig{}

		// Start at 0, should land immediately
		result := countdown(CountdownState{count: 0})(config)
		assert.Equal(t, "Done!", result)
	})
}

func TestTailRec_Sum(t *testing.T) {
	t.Run("sums numbers from 1 to n", func(t *testing.T) {
		sumStep := func(state SumState) Reader[TestConfig, Trampoline[SumState, int]] {
			return func(cfg TestConfig) Trampoline[SumState, int] {
				if state.current <= 0 {
					return TR.Land[SumState](state.total)
				}
				return TR.Bounce[int](SumState{
					current: state.current - 1,
					total:   state.total + state.current,
				})
			}
		}

		sum := TailRec(sumStep)
		config := TestConfig{}

		testCases := []struct {
			n        int
			expected int
		}{
			{0, 0},
			{1, 1},
			{5, 15},     // 1+2+3+4+5
			{10, 55},    // 1+2+...+10
			{100, 5050}, // 1+2+...+100
		}

		for _, tc := range testCases {
			result := sum(SumState{current: tc.n, total: 0})(config)
			assert.Equal(t, tc.expected, result, "sum(1..%d) should equal %d", tc.n, tc.expected)
		}
	})
}

func TestTailRec_WithEnvironment(t *testing.T) {
	t.Run("uses environment in computation", func(t *testing.T) {
		// Multiply accumulator by environment multiplier at each step
		multiplyStep := func(state FactorialState) Reader[TestConfig, Trampoline[FactorialState, int]] {
			return func(cfg TestConfig) Trampoline[FactorialState, int] {
				if state.n <= 0 {
					return TR.Land[FactorialState](state.acc)
				}
				return TR.Bounce[int](FactorialState{
					n:   state.n - 1,
					acc: state.acc * cfg.multiplier,
				})
			}
		}

		multiply := TailRec(multiplyStep)

		// Test with different multipliers
		config2 := TestConfig{multiplier: 2}
		result2 := multiply(FactorialState{n: 5, acc: 1})(config2)
		assert.Equal(t, 32, result2, "1 * 2^5 should equal 32")

		config3 := TestConfig{multiplier: 3}
		result3 := multiply(FactorialState{n: 4, acc: 1})(config3)
		assert.Equal(t, 81, result3, "1 * 3^4 should equal 81")
	})

	t.Run("respects environment limits", func(t *testing.T) {
		limitedStep := func(state CountdownState) Reader[TestConfig, Trampoline[CountdownState, int]] {
			return func(cfg TestConfig) Trampoline[CountdownState, int] {
				if state.count <= 0 {
					return TR.Land[CountdownState](state.count)
				}
				if state.count > cfg.maxIterations {
					// Hit limit, return error value
					return TR.Land[CountdownState](-1)
				}
				return TR.Bounce[int](CountdownState{count: state.count - 1})
			}
		}

		limitedCountdown := TailRec(limitedStep)

		// Within limit
		config := TestConfig{maxIterations: 100}
		result := limitedCountdown(CountdownState{count: 50})(config)
		assert.Equal(t, 0, result)

		// Exceeds limit
		result2 := limitedCountdown(CountdownState{count: 200})(config)
		assert.Equal(t, -1, result2)
	})
}

func TestTailRec_DifferentTypes(t *testing.T) {
	t.Run("works with string accumulator", func(t *testing.T) {
		type StringState struct {
			count int
			str   string
		}

		buildStringStep := func(state StringState) Reader[TestConfig, Trampoline[StringState, string]] {
			return func(cfg TestConfig) Trampoline[StringState, string] {
				if state.count <= 0 {
					return TR.Land[StringState](state.str)
				}
				return TR.Bounce[string](StringState{
					count: state.count - 1,
					str:   state.str + "x",
				})
			}
		}

		buildString := TailRec(buildStringStep)
		config := TestConfig{}

		result := buildString(StringState{count: 5, str: ""})(config)
		assert.Equal(t, "xxxxx", result)
	})

	t.Run("works with slice accumulator", func(t *testing.T) {
		type SliceState struct {
			count int
			items []int
		}

		collectStep := func(state SliceState) Reader[TestConfig, Trampoline[SliceState, []int]] {
			return func(cfg TestConfig) Trampoline[SliceState, []int] {
				if state.count <= 0 {
					return TR.Land[SliceState](state.items)
				}
				newItems := append(state.items, state.count)
				return TR.Bounce[[]int](SliceState{
					count: state.count - 1,
					items: newItems,
				})
			}
		}

		collect := TailRec(collectStep)
		config := TestConfig{}

		result := collect(SliceState{count: 5, items: []int{}})(config)
		assert.Equal(t, []int{5, 4, 3, 2, 1}, result)
	})
}

func TestTailRec_EdgeCases(t *testing.T) {
	t.Run("handles single bounce", func(t *testing.T) {
		singleBounceStep := func(n int) Reader[TestConfig, Trampoline[int, string]] {
			return func(cfg TestConfig) Trampoline[int, string] {
				if n == 0 {
					return TR.Land[int]("landed")
				}
				return TR.Bounce[string](0)
			}
		}

		singleBounce := TailRec(singleBounceStep)
		config := TestConfig{}

		result := singleBounce(1)(config)
		assert.Equal(t, "landed", result)
	})

	t.Run("handles zero iterations", func(t *testing.T) {
		immediateStep := func(n int) Reader[TestConfig, Trampoline[int, int]] {
			return func(cfg TestConfig) Trampoline[int, int] {
				return TR.Land[int](n * 2)
			}
		}

		immediate := TailRec(immediateStep)
		config := TestConfig{}

		result := immediate(42)(config)
		assert.Equal(t, 84, result)
	})
}

func TestTailRec_ComplexState(t *testing.T) {
	t.Run("handles complex state transitions", func(t *testing.T) {
		type ComplexState struct {
			phase int
			value int
			done  bool
		}

		// Multi-phase computation
		complexStep := func(state ComplexState) Reader[TestConfig, Trampoline[ComplexState, int]] {
			return func(cfg TestConfig) Trampoline[ComplexState, int] {
				if state.done {
					return TR.Land[ComplexState](state.value)
				}

				switch state.phase {
				case 0:
					// Phase 0: multiply by 2
					return TR.Bounce[int](ComplexState{
						phase: 1,
						value: state.value * 2,
						done:  false,
					})
				case 1:
					// Phase 1: add 10
					return TR.Bounce[int](ComplexState{
						phase: 2,
						value: state.value + 10,
						done:  false,
					})
				case 2:
					// Phase 2: multiply by multiplier from config
					return TR.Bounce[int](ComplexState{
						phase: 3,
						value: state.value * cfg.multiplier,
						done:  true,
					})
				default:
					return TR.Land[ComplexState](state.value)
				}
			}
		}

		complex := TailRec(complexStep)
		config := TestConfig{multiplier: 3}

		// (5 * 2 + 10) * 3 = 60
		result := complex(ComplexState{phase: 0, value: 5, done: false})(config)
		assert.Equal(t, 60, result)
	})
}
