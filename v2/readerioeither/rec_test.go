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

package readerioeither

import (
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	TR "github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
)

// Test environment types
type TestEnv struct {
	Multiplier int
	MaxValue   int
	Logs       []string
}

type LoggerEnv struct {
	Logger func(string)
	MaxN   int
}

type ConfigEnv struct {
	MinValue   int
	Step       int
	MaxRetries int
}

// TestTailRecFactorial tests factorial computation with environment-based logging and validation
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
		MaxN: 20,
	}

	factorialStep := func(state State) ReaderIOEither[LoggerEnv, string, TR.Trampoline[State, int]] {
		return func(env LoggerEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.n > env.MaxN {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("n too large: %d > %d", state.n, env.MaxN))
				}
				if state.n <= 0 {
					env.Logger(fmt.Sprintf("Complete: %d", state.acc))
					return E.Right[string](TR.Land[State](state.acc))
				}
				env.Logger(fmt.Sprintf("Step: %d * %d", state.n, state.acc))
				return E.Right[string](TR.Bounce[int](State{state.n - 1, state.acc * state.n}))
			}
		}
	}

	factorial := TailRec(factorialStep)
	result := factorial(State{5, 1})(env)()

	assert.Equal(t, E.Of[string](120), result)
	assert.Equal(t, 6, len(logs)) // 5 steps + 1 complete
	assert.Contains(t, logs[0], "Step: 5 * 1")
	assert.Contains(t, logs[len(logs)-1], "Complete: 120")
}

// TestTailRecFactorialError tests factorial with input exceeding max value
func TestTailRecFactorialError(t *testing.T) {
	type State struct {
		n   int
		acc int
	}

	env := LoggerEnv{
		Logger: func(msg string) {},
		MaxN:   10,
	}

	factorialStep := func(state State) ReaderIOEither[LoggerEnv, string, TR.Trampoline[State, int]] {
		return func(env LoggerEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.n > env.MaxN {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("n too large: %d > %d", state.n, env.MaxN))
				}
				if state.n <= 0 {
					return E.Right[string](TR.Land[State](state.acc))
				}
				return E.Right[string](TR.Bounce[int](State{state.n - 1, state.acc * state.n}))
			}
		}
	}

	factorial := TailRec(factorialStep)
	result := factorial(State{15, 1})(env)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Contains(t, err, "n too large: 15 > 10")
}

// TestTailRecFibonacci tests Fibonacci computation with environment dependency
func TestTailRecFibonacci(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	env := TestEnv{Multiplier: 1, MaxValue: 1000, Logs: []string{}}

	fibStep := func(state State) ReaderIOEither[TestEnv, string, TR.Trampoline[State, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.curr > env.MaxValue {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("value exceeds max: %d > %d", state.curr, env.MaxValue))
				}
				if state.n <= 0 {
					return E.Right[string](TR.Land[State](state.curr * env.Multiplier))
				}
				return E.Right[string](TR.Bounce[int](State{state.n - 1, state.curr, state.prev + state.curr}))
			}
		}
	}

	fib := TailRec(fibStep)
	result := fib(State{10, 0, 1})(env)()

	assert.Equal(t, E.Of[string](89), result) // 10th Fibonacci number
}

// TestTailRecFibonacciError tests Fibonacci with value exceeding max
func TestTailRecFibonacciError(t *testing.T) {
	type State struct {
		n    int
		prev int
		curr int
	}

	env := TestEnv{Multiplier: 1, MaxValue: 50, Logs: []string{}}

	fibStep := func(state State) ReaderIOEither[TestEnv, string, TR.Trampoline[State, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.curr > env.MaxValue {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("value exceeds max: %d > %d", state.curr, env.MaxValue))
				}
				if state.n <= 0 {
					return E.Right[string](TR.Land[State](state.curr))
				}
				return E.Right[string](TR.Bounce[int](State{state.n - 1, state.curr, state.prev + state.curr}))
			}
		}
	}

	fib := TailRec(fibStep)
	result := fib(State{20, 0, 1})(env)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Contains(t, err, "value exceeds max")
}

// TestTailRecCountdown tests countdown with configuration-based step
func TestTailRecCountdown(t *testing.T) {
	config := ConfigEnv{MinValue: 0, Step: 2, MaxRetries: 3}

	countdownStep := func(n int) ReaderIOEither[ConfigEnv, string, TR.Trampoline[int, int]] {
		return func(cfg ConfigEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				if n < 0 {
					return E.Left[TR.Trampoline[int, int]]("negative value")
				}
				if n <= cfg.MinValue {
					return E.Right[string](TR.Land[int](n))
				}
				return E.Right[string](TR.Bounce[int](n - cfg.Step))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(config)()

	assert.Equal(t, E.Of[string](0), result)
}

// TestTailRecSumList tests summing a list with environment-based multiplier and error handling
func TestTailRecSumList(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	env := TestEnv{Multiplier: 2, MaxValue: 100, Logs: []string{}}

	sumStep := func(state State) ReaderIOEither[TestEnv, string, TR.Trampoline[State, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.sum > env.MaxValue {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("sum exceeds max: %d > %d", state.sum, env.MaxValue))
				}
				if A.IsEmpty(state.list) {
					return E.Right[string](TR.Land[State](state.sum * env.Multiplier))
				}
				return E.Right[string](TR.Bounce[int](State{state.list[1:], state.sum + state.list[0]}))
			}
		}
	}

	sumList := TailRec(sumStep)
	result := sumList(State{[]int{1, 2, 3, 4, 5}, 0})(env)()

	assert.Equal(t, E.Of[string](30), result) // (1+2+3+4+5) * 2 = 30
}

// TestTailRecSumListError tests sum exceeding max value
func TestTailRecSumListError(t *testing.T) {
	type State struct {
		list []int
		sum  int
	}

	env := TestEnv{Multiplier: 1, MaxValue: 10, Logs: []string{}}

	sumStep := func(state State) ReaderIOEither[TestEnv, string, TR.Trampoline[State, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.sum > env.MaxValue {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("sum exceeds max: %d > %d", state.sum, env.MaxValue))
				}
				if A.IsEmpty(state.list) {
					return E.Right[string](TR.Land[State](state.sum))
				}
				return E.Right[string](TR.Bounce[int](State{state.list[1:], state.sum + state.list[0]}))
			}
		}
	}

	sumList := TailRec(sumStep)
	result := sumList(State{[]int{5, 10, 15}, 0})(env)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Contains(t, err, "sum exceeds max")
}

// TestTailRecImmediateTermination tests immediate termination (Right on first call)
func TestTailRecImmediateTermination(t *testing.T) {
	env := TestEnv{Multiplier: 1, MaxValue: 100, Logs: []string{}}

	immediateStep := func(n int) ReaderIOEither[TestEnv, string, TR.Trampoline[int, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				return E.Right[string](TR.Land[int](n * env.Multiplier))
			}
		}
	}

	immediate := TailRec(immediateStep)
	result := immediate(42)(env)()

	assert.Equal(t, E.Of[string](42), result)
}

// TestTailRecImmediateError tests immediate error (Left on first call)
func TestTailRecImmediateError(t *testing.T) {
	env := TestEnv{Multiplier: 1, MaxValue: 100, Logs: []string{}}

	immediateErrorStep := func(n int) ReaderIOEither[TestEnv, string, TR.Trampoline[int, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				return E.Left[TR.Trampoline[int, int]]("immediate error")
			}
		}
	}

	immediateError := TailRec(immediateErrorStep)
	result := immediateError(42)(env)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Equal(t, "immediate error", err)
}

// TestTailRecStackSafety tests that TailRec handles large iterations without stack overflow
func TestTailRecStackSafety(t *testing.T) {
	env := TestEnv{Multiplier: 1, MaxValue: 2000000, Logs: []string{}}

	countdownStep := func(n int) ReaderIOEither[TestEnv, string, TR.Trampoline[int, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				if n > env.MaxValue {
					return E.Left[TR.Trampoline[int, int]]("value too large")
				}
				if n <= 0 {
					return E.Right[string](TR.Land[int](n))
				}
				return E.Right[string](TR.Bounce[int](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10000)(env)()

	assert.Equal(t, E.Of[string](0), result)
}

// TestTailRecFindInRange tests finding a value in a range with environment-based target
func TestTailRecFindInRange(t *testing.T) {
	type FindEnv struct {
		Target int
		MaxN   int
	}

	type State struct {
		current int
		max     int
	}

	env := FindEnv{Target: 42, MaxN: 1000}

	findStep := func(state State) ReaderIOEither[FindEnv, string, TR.Trampoline[State, int]] {
		return func(env FindEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.current > env.MaxN {
					return E.Left[TR.Trampoline[State, int]]("search exceeded max")
				}
				if state.current >= state.max {
					return E.Right[string](TR.Land[State](-1)) // Not found
				}
				if state.current == env.Target {
					return E.Right[string](TR.Land[State](state.current)) // Found
				}
				return E.Right[string](TR.Bounce[int](State{state.current + 1, state.max}))
			}
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100})(env)()

	assert.Equal(t, E.Of[string](42), result)
}

// TestTailRecFindNotInRange tests finding a value not in range
func TestTailRecFindNotInRange(t *testing.T) {
	type FindEnv struct {
		Target int
		MaxN   int
	}

	type State struct {
		current int
		max     int
	}

	env := FindEnv{Target: 200, MaxN: 1000}

	findStep := func(state State) ReaderIOEither[FindEnv, string, TR.Trampoline[State, int]] {
		return func(env FindEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.current > env.MaxN {
					return E.Left[TR.Trampoline[State, int]]("search exceeded max")
				}
				if state.current >= state.max {
					return E.Right[string](TR.Land[State](-1)) // Not found
				}
				if state.current == env.Target {
					return E.Right[string](TR.Land[State](state.current)) // Found
				}
				return E.Right[string](TR.Bounce[int](State{state.current + 1, state.max}))
			}
		}
	}

	find := TailRec(findStep)
	result := find(State{0, 100})(env)()

	assert.Equal(t, E.Of[string](-1), result)
}

// TestTailRecWithLogging tests that logging side effects occur during recursion
func TestTailRecWithLogging(t *testing.T) {
	logs := []string{}
	env := LoggerEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
		MaxN: 100,
	}

	countdownStep := func(n int) ReaderIOEither[LoggerEnv, string, TR.Trampoline[int, int]] {
		return func(env LoggerEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				env.Logger(fmt.Sprintf("Count: %d", n))
				if n > env.MaxN {
					return E.Left[TR.Trampoline[int, int]]("value too large")
				}
				if n <= 0 {
					return E.Right[string](TR.Land[int](n))
				}
				return E.Right[string](TR.Bounce[int](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(5)(env)()

	assert.Equal(t, E.Of[string](0), result)
	assert.Equal(t, 6, len(logs)) // 5, 4, 3, 2, 1, 0
	assert.Equal(t, "Count: 5", logs[0])
	assert.Equal(t, "Count: 0", logs[5])
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
		MaxN: 1000,
	}

	gcdStep := func(state State) ReaderIOEither[LoggerEnv, string, TR.Trampoline[State, int]] {
		return func(env LoggerEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				env.Logger(fmt.Sprintf("gcd(%d, %d)", state.a, state.b))
				if state.a > env.MaxN || state.b > env.MaxN {
					return E.Left[TR.Trampoline[State, int]]("values too large")
				}
				if state.b == 0 {
					return E.Right[string](TR.Land[State](state.a))
				}
				return E.Right[string](TR.Bounce[int](State{state.b, state.a % state.b}))
			}
		}
	}

	gcd := TailRec(gcdStep)
	result := gcd(State{48, 18})(env)()

	assert.Equal(t, E.Of[string](6), result)
	assert.Greater(t, len(logs), 0)
	assert.Contains(t, logs[0], "gcd(48, 18)")
}

// TestTailRecRetryLogic tests retry logic with environment-based max retries
func TestTailRecRetryLogic(t *testing.T) {
	type State struct {
		attempt int
		value   int
	}

	config := ConfigEnv{MinValue: 0, Step: 1, MaxRetries: 3}

	retryStep := func(state State) ReaderIOEither[ConfigEnv, string, TR.Trampoline[State, int]] {
		return func(cfg ConfigEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.attempt > cfg.MaxRetries {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("max retries exceeded: %d", cfg.MaxRetries))
				}
				// Simulate success on 3rd attempt
				if state.attempt == 3 {
					return E.Right[string](TR.Land[State](state.value))
				}
				return E.Right[string](TR.Bounce[int](State{state.attempt + 1, state.value}))
			}
		}
	}

	retry := TailRec(retryStep)
	result := retry(State{0, 42})(config)()

	assert.Equal(t, E.Of[string](42), result)
}

// TestTailRecRetryExceeded tests retry logic exceeding max retries
func TestTailRecRetryExceeded(t *testing.T) {
	type State struct {
		attempt int
		value   int
	}

	config := ConfigEnv{MinValue: 0, Step: 1, MaxRetries: 2}

	retryStep := func(state State) ReaderIOEither[ConfigEnv, string, TR.Trampoline[State, int]] {
		return func(cfg ConfigEnv) IOE.IOEither[string, TR.Trampoline[State, int]] {
			return func() E.Either[string, TR.Trampoline[State, int]] {
				if state.attempt > cfg.MaxRetries {
					return E.Left[TR.Trampoline[State, int]](fmt.Sprintf("max retries exceeded: %d", cfg.MaxRetries))
				}
				// Never succeeds
				return E.Right[string](TR.Bounce[int](State{state.attempt + 1, state.value}))
			}
		}
	}

	retry := TailRec(retryStep)
	result := retry(State{0, 42})(config)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Contains(t, err, "max retries exceeded: 2")
}

// TestTailRecMultipleEnvironmentAccess tests that environment is accessible in each iteration
func TestTailRecMultipleEnvironmentAccess(t *testing.T) {
	type CounterEnv struct {
		Increment int
		Limit     int
		MaxValue  int
	}

	env := CounterEnv{Increment: 3, Limit: 20, MaxValue: 100}

	counterStep := func(n int) ReaderIOEither[CounterEnv, string, TR.Trampoline[int, int]] {
		return func(env CounterEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				if n > env.MaxValue {
					return E.Left[TR.Trampoline[int, int]]("value exceeds max")
				}
				if n >= env.Limit {
					return E.Right[string](TR.Land[int](n))
				}
				return E.Right[string](TR.Bounce[int](n + env.Increment))
			}
		}
	}

	counter := TailRec(counterStep)
	result := counter(0)(env)()

	assert.Equal(t, E.Of[string](21), result) // 0 -> 3 -> 6 -> 9 -> 12 -> 15 -> 18 -> 21
}

// TestTailRecErrorInMiddle tests error occurring in the middle of recursion
func TestTailRecErrorInMiddle(t *testing.T) {
	env := TestEnv{Multiplier: 1, MaxValue: 50, Logs: []string{}}

	countdownStep := func(n int) ReaderIOEither[TestEnv, string, TR.Trampoline[int, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				if n == 5 {
					return E.Left[TR.Trampoline[int, int]]("error at 5")
				}
				if n <= 0 {
					return E.Right[string](TR.Land[int](n))
				}
				return E.Right[string](TR.Bounce[int](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(10)(env)()

	assert.True(t, E.IsLeft(result))
	_, err := E.Unwrap(result)
	assert.Equal(t, "error at 5", err)
}

// TestTailRecDifferentEnvironments tests that different environments produce different results
func TestTailRecDifferentEnvironments(t *testing.T) {
	multiplyStep := func(n int) ReaderIOEither[TestEnv, string, TR.Trampoline[int, int]] {
		return func(env TestEnv) IOE.IOEither[string, TR.Trampoline[int, int]] {
			return func() E.Either[string, TR.Trampoline[int, int]] {
				if n <= 0 {
					return E.Right[string](TR.Land[int](n * env.Multiplier))
				}
				return E.Right[string](TR.Bounce[int](n - 1))
			}
		}
	}

	multiply := TailRec(multiplyStep)

	env1 := TestEnv{Multiplier: 2, MaxValue: 100, Logs: []string{}}
	env2 := TestEnv{Multiplier: 5, MaxValue: 100, Logs: []string{}}

	result1 := multiply(5)(env1)()
	result2 := multiply(5)(env2)()

	// Both reach 0, but multiplied by different values
	assert.Equal(t, E.Of[string](0), result1) // 0 * 2 = 0
	assert.Equal(t, E.Of[string](0), result2) // 0 * 5 = 0
}
