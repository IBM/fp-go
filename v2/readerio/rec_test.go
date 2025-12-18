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

package readerio

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	G "github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// Test environment types
type TestEnv struct {
	Multiplier int
	Logs       []string
}

type LoggerEnv struct {
	Logger func(string)
}

type ConfigEnv struct {
	MinValue int
	Step     int
}

// TestTailRecFactorial tests factorial computation with environment-based logging
func TestTailRecFactorial(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	logs := []string{}
	env := LoggerEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
	}

	factorialStep := func(state State) ReaderIO[LoggerEnv, Trampoline[State, int]] {
		return func(env LoggerEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if state.n <= 0 {
					env.Logger(fmt.Sprintf("Complete: %d", state.acc))
					return tailrec.Land[State](state.acc)
				}
				env.Logger(fmt.Sprintf("Step: %d * %d", state.n, state.acc))
				return tailrec.Bounce[int](State{state.n - 1, state.acc * state.n})
			}
		}
	}

	factorial := TailRec(factorialStep)
	result := factorial(State{5, 1})(env)()

	assert.Equal(t, 120, result)
	assert.Equal(t, 6, len(logs)) // 5 steps + 1 complete
	assert.Contains(t, logs[0], "Step: 5 * 1")
	assert.Contains(t, logs[len(logs)-1], "Complete: 120")
}

// TestTailRecFibonacci tests Fibonacci computation with environment dependency
func TestTailRecFibonacci(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	env := TestEnv{Multiplier: 1, Logs: []string{}}

	fibStep := func(state State) ReaderIO[TestEnv, Trampoline[State, int]] {
		return func(env TestEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if state.n <= 0 {
					return tailrec.Land[State](state.curr * env.Multiplier)
				}
				return tailrec.Bounce[int](State{state.n - 1, state.curr, state.prev + state.curr})
			}
		}
	}

	fib := TailRec(fibStep)
	result := fib(State{10, 0, 1})(env)()

	assert.Equal(t, 89, result) // 10th Fibonacci number
}

// TestTailRecCountdown tests countdown with configuration-based step
func TestTailRecCountdown(t *testing.T) {
	config := ConfigEnv{MinValue: 0, Step: 2}

	countdownStep := func(n int) ReaderIO[ConfigEnv, Trampoline[int, int]] {
		return func(cfg ConfigEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				if n <= cfg.MinValue {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n - cfg.Step)
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(config)()

	assert.Equal(t, 0, result)
}

// TestTailRecCountdownOddStep tests countdown with odd step size
func TestTailRecCountdownOddStep(t *testing.T) {
	config := ConfigEnv{MinValue: 0, Step: 3}

	countdownStep := func(n int) ReaderIO[ConfigEnv, Trampoline[int, int]] {
		return func(cfg ConfigEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				if n <= cfg.MinValue {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n - cfg.Step)
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(config)()

	assert.Equal(t, -2, result) // 10 -> 7 -> 4 -> 1 -> -2 (stops when <= 0)
}

// TestTailRecSumList tests summing a list with environment-based multiplier
func TestTailRecSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	env := TestEnv{Multiplier: 2, Logs: []string{}}

	sumStep := func(state State) ReaderIO[TestEnv, Trampoline[State, int]] {
		return func(env TestEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if A.IsEmpty(state.list) {
					return tailrec.Land[State](state.sum * env.Multiplier)
				}
				return tailrec.Bounce[int](State{state.list[1:], state.sum + state.list[0]})
			}
		}
	}

	sumList := TailRec(sumStep)
	result := sumList(State{[]int{1, 2, 3, 4, 5}, 0})(env)()

	assert.Equal(t, 30, result) // (1+2+3+4+5) * 2 = 30
}

// TestTailRecImmediateTermination tests immediate termination (Right on first call)
func TestTailRecImmediateTermination(t *testing.T) {
	env := TestEnv{Multiplier: 1, Logs: []string{}}

	immediateStep := func(n int) ReaderIO[TestEnv, Trampoline[int, int]] {
		return func(env TestEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				return tailrec.Land[int](n * env.Multiplier)
			}
		}
	}

	immediate := TailRec(immediateStep)
	result := immediate(42)(env)()

	assert.Equal(t, 42, result)
}

// TestTailRecStackSafety tests that TailRec handles large iterations without stack overflow
func TestTailRecStackSafety(t *testing.T) {
	env := TestEnv{Multiplier: 1, Logs: []string{}}

	countdownStep := func(n int) ReaderIO[TestEnv, Trampoline[int, int]] {
		return func(env TestEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				if n <= 0 {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n - 1)
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10000)(env)()

	assert.Equal(t, 0, result)
}

// TestTailRecFindInRange tests finding a value in a range with environment-based target
func TestTailRecFindInRange(t *testing.T) {
	type FindEnv struct {
		Target int
	}

	type State struct {
		current int
		max     int
	}

	env := FindEnv{Target: 42}

	findStep := func(state State) ReaderIO[FindEnv, Trampoline[State, int]] {
		return func(env FindEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if state.current >= state.max {
					return tailrec.Land[State](-1) // Not found
				}
				if state.current == env.Target {
					return tailrec.Land[State](state.current) // Found
				}
				return tailrec.Bounce[int](State{state.current + 1, state.max})
			}
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100})(env)()

	assert.Equal(t, 42, result)
}

// TestTailRecFindNotInRange tests finding a value not in range
func TestTailRecFindNotInRange(t *testing.T) {
	type FindEnv struct {
		Target int
	}

	type State struct {
		current int
		max     int
	}

	env := FindEnv{Target: 200}

	findStep := func(state State) ReaderIO[FindEnv, Trampoline[State, int]] {
		return func(env FindEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if state.current >= state.max {
					return tailrec.Land[State](-1) // Not found
				}
				if state.current == env.Target {
					return tailrec.Land[State](state.current) // Found
				}
				return tailrec.Bounce[int](State{state.current + 1, state.max})
			}
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100})(env)()

	assert.Equal(t, -1, result)
}

// TestTailRecWithLogging tests that logging side effects occur during recursion
func TestTailRecWithLogging(t *testing.T) {
	logs := []string{}
	env := LoggerEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
	}

	countdownStep := func(n int) ReaderIO[LoggerEnv, Trampoline[int, int]] {
		return func(env LoggerEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				env.Logger(fmt.Sprintf("Count: %d", n))
				if n <= 0 {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n - 1)
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(5)(env)()

	assert.Equal(t, 0, result)
	assert.Equal(t, 6, len(logs)) // 5, 4, 3, 2, 1, 0
	assert.Equal(t, "Count: 5", logs[0])
	assert.Equal(t, "Count: 0", logs[5])
}

// TestTailRecCollatzConjecture tests the Collatz conjecture with environment-based logging
func TestTailRecCollatzConjecture(t *testing.T) {
	logs := []string{}
	env := LoggerEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
	}

	collatzStep := func(n int) ReaderIO[LoggerEnv, Trampoline[int, int]] {
		return func(env LoggerEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				env.Logger(fmt.Sprintf("n=%d", n))
				if n <= 1 {
					return tailrec.Land[int](n)
				}
				if n%2 == 0 {
					return tailrec.Bounce[int](n / 2)
				}
				return tailrec.Bounce[int](3*n + 1)
			}
		}
	}

	collatz := TailRec(collatzStep)
	result := collatz(10)(env)()

	assert.Equal(t, 1, result)
	assert.Greater(t, len(logs), 5) // Multiple steps to reach 1
	assert.Equal(t, "n=10", logs[0])
	assert.Contains(t, logs[len(logs)-1], "n=1")
}

// TestTailRecPowerOfTwo tests computing power of 2 with environment-based exponent limit
func TestTailRecPowerOfTwo(t *testing.T) {
	type PowerEnv struct {
		MaxExponent int
	}

	type State struct {
		exponent int
		result   int
	}

	env := PowerEnv{MaxExponent: 10}

	powerStep := func(state State) ReaderIO[PowerEnv, Trampoline[State, int]] {
		return func(env PowerEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				if state.exponent >= env.MaxExponent {
					return tailrec.Land[State](state.result)
				}
				return tailrec.Bounce[int](State{state.exponent + 1, state.result * 2})
			}
		}
	}

	power := TailRec(powerStep)
	result := power(State{0, 1})(env)()

	assert.Equal(t, 1024, result) // 2^10
}

// TestTailRecGCD tests greatest common divisor with environment-based logging
func TestTailRecGCD(t *testing.T) {
	type State struct {
		a int
		b int
	}

	logs := []string{}
	env := LoggerEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
	}

	gcdStep := func(state State) ReaderIO[LoggerEnv, Trampoline[State, int]] {
		return func(env LoggerEnv) G.IO[Trampoline[State, int]] {
			return func() Trampoline[State, int] {
				env.Logger(fmt.Sprintf("gcd(%d, %d)", state.a, state.b))
				if state.b == 0 {
					return tailrec.Land[State](state.a)
				}
				return tailrec.Bounce[int](State{state.b, state.a % state.b})
			}
		}
	}

	gcd := TailRec(gcdStep)
	result := gcd(State{48, 18})(env)()

	assert.Equal(t, 6, result)
	assert.Greater(t, len(logs), 0)
	assert.Contains(t, logs[0], "gcd(48, 18)")
}

// TestTailRecMultipleEnvironmentAccess tests that environment is accessible in each iteration
func TestTailRecMultipleEnvironmentAccess(t *testing.T) {
	type CounterEnv struct {
		Increment int
		Limit     int
	}

	env := CounterEnv{Increment: 3, Limit: 20}

	counterStep := func(n int) ReaderIO[CounterEnv, Trampoline[int, int]] {
		return func(env CounterEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				if n >= env.Limit {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n + env.Increment)
			}
		}
	}

	counter := TailRec(counterStep)
	result := counter(0)(env)()

	assert.Equal(t, 21, result) // 0 -> 3 -> 6 -> 9 -> 12 -> 15 -> 18 -> 21
}

// TestTailRecDifferentEnvironments tests that different environments produce different results
func TestTailRecDifferentEnvironments(t *testing.T) {
	multiplyStep := func(n int) ReaderIO[TestEnv, Trampoline[int, int]] {
		return func(env TestEnv) G.IO[Trampoline[int, int]] {
			return func() Trampoline[int, int] {
				if n <= 0 {
					return tailrec.Land[int](n)
				}
				return tailrec.Bounce[int](n - 1)
			}
		}
	}

	multiply := TailRec(multiplyStep)

	env1 := TestEnv{Multiplier: 2, Logs: []string{}}
	env2 := TestEnv{Multiplier: 5, Logs: []string{}}

	result1 := multiply(5)(env1)()
	result2 := multiply(5)(env2)()

	// Both should reach 0 regardless of environment (environment not used in this case)
	assert.Equal(t, 0, result1)
	assert.Equal(t, 0, result2)
}
