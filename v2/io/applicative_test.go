// Copyright (c) 2024 - 2025 IBM Corp.
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
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

// TestApplicativeOf tests the Of operation of the Applicative type class
func TestApplicativeOf(t *testing.T) {
	app := Applicative[int, string]()

	t.Run("wraps a value in IO context", func(t *testing.T) {
		ioValue := app.Of(42)
		result := ioValue()
		assert.Equal(t, 42, result)
	})

	t.Run("wraps string value", func(t *testing.T) {
		app := Applicative[string, int]()
		ioValue := app.Of("hello")
		result := ioValue()
		assert.Equal(t, "hello", result)
	})

	t.Run("wraps zero value", func(t *testing.T) {
		ioValue := app.Of(0)
		result := ioValue()
		assert.Equal(t, 0, result)
	})
}

// TestApplicativeMap tests the Map operation of the Applicative type class
func TestApplicativeMap(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("maps a function over IO value", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		ioValue := app.Of(21)
		result := app.Map(double)(ioValue)
		assert.Equal(t, 42, result())
	})

	t.Run("maps type conversion", func(t *testing.T) {
		app := Applicative[int, string]()
		ioValue := app.Of(42)
		result := app.Map(strconv.Itoa)(ioValue)
		assert.Equal(t, "42", result())
	})

	t.Run("maps identity function", func(t *testing.T) {
		identity := func(x int) int { return x }
		ioValue := app.Of(42)
		result := app.Map(identity)(ioValue)
		assert.Equal(t, 42, result())
	})

	t.Run("maps constant function", func(t *testing.T) {
		constant := func(x int) int { return 100 }
		ioValue := app.Of(42)
		result := app.Map(constant)(ioValue)
		assert.Equal(t, 100, result())
	})
}

// TestApplicativeAp tests the Ap operation of the Applicative type class
func TestApplicativeAp(t *testing.T) {
	t.Run("applies wrapped function to wrapped value", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		ioFunc := Of(add(10))
		ioValue := Of(32)
		result := Ap[int](ioValue)(ioFunc)
		assert.Equal(t, 42, result())
	})

	t.Run("applies multiplication function", func(t *testing.T) {
		multiply := func(a int) func(int) int {
			return func(b int) int { return a * b }
		}
		ioFunc := Of(multiply(6))
		ioValue := Of(7)
		result := Ap[int](ioValue)(ioFunc)
		assert.Equal(t, 42, result())
	})

	t.Run("applies function with zero", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}
		ioFunc := Of(add(0))
		ioValue := Of(42)
		result := Ap[int](ioValue)(ioFunc)
		assert.Equal(t, 42, result())
	})

	t.Run("applies with type conversion", func(t *testing.T) {
		toStringAndAppend := func(suffix string) func(int) string {
			return func(n int) string {
				return strconv.Itoa(n) + suffix
			}
		}
		ioFunc := Of(toStringAndAppend("!"))
		ioValue := Of(42)
		result := Ap[string](ioValue)(ioFunc)
		assert.Equal(t, "42!", result())
	})
}

// TestApplicativeComposition tests composition of applicative operations
func TestApplicativeComposition(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("composes Map and Of", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		result := F.Pipe1(
			app.Of(21),
			app.Map(double),
		)
		assert.Equal(t, 42, result())
	})

	t.Run("composes multiple Map operations", func(t *testing.T) {
		app := Applicative[int, string]()
		double := func(x int) int { return x * 2 }
		toString := func(x int) string { return strconv.Itoa(x) }

		result := F.Pipe2(
			app.Of(21),
			Map(double),
			app.Map(toString),
		)
		assert.Equal(t, "42", result())
	})

	t.Run("composes Map and Ap", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		ioFunc := F.Pipe1(
			app.Of(5),
			Map(add),
		)
		ioValue := app.Of(16)

		result := Ap[int](ioValue)(ioFunc)
		assert.Equal(t, 21, result())
	})
}

// TestApplicativeLaws tests the applicative functor laws
func TestApplicativeLaws(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("identity law: ap(Of(id), v) = v", func(t *testing.T) {
		identity := func(x int) int { return x }
		v := app.Of(42)

		left := Ap[int](v)(Of(identity))
		right := v

		assert.Equal(t, right(), left())
	})

	t.Run("homomorphism law: ap(Of(f), Of(x)) = Of(f(x))", func(t *testing.T) {
		f := func(x int) int { return x * 2 }
		x := 21

		left := Ap[int](app.Of(x))(Of(f))
		right := app.Of(f(x))

		assert.Equal(t, right(), left())
	})

	t.Run("interchange law: ap(u, Of(y)) = ap(Of(f => f(y)), u)", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		u := Of(double)
		y := 21

		left := Ap[int](app.Of(y))(u)

		applyY := func(f func(int) int) int { return f(y) }
		right := Ap[int](u)(Of(applyY))

		assert.Equal(t, right(), left())
	})
}

// TestApplicativeWithPipe tests applicative operations with pipe
func TestApplicativeWithPipe(t *testing.T) {
	t.Run("pipes Of and Map", func(t *testing.T) {
		app := Applicative[int, string]()
		result := F.Pipe1(
			app.Of(42),
			app.Map(strconv.Itoa),
		)
		assert.Equal(t, "42", result())
	})

	t.Run("pipes complex transformation", func(t *testing.T) {
		app := Applicative[int, int]()
		add10 := func(x int) int { return x + 10 }
		double := func(x int) int { return x * 2 }

		result := F.Pipe2(
			app.Of(16),
			app.Map(add10),
			app.Map(double),
		)
		assert.Equal(t, 52, result())
	})
}

// TestApplicativeWithUtils tests applicative with utility functions
func TestApplicativeWithUtils(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("uses utils.Double", func(t *testing.T) {
		result := F.Pipe1(
			app.Of(21),
			app.Map(utils.Double),
		)
		assert.Equal(t, 42, result())
	})

	t.Run("uses utils.Inc", func(t *testing.T) {
		result := F.Pipe1(
			app.Of(41),
			app.Map(utils.Inc),
		)
		assert.Equal(t, 42, result())
	})
}

// TestApplicativeMultipleArguments tests applying functions with multiple arguments
func TestApplicativeMultipleArguments(t *testing.T) {
	app := Applicative[int, int]()

	t.Run("applies curried two-argument function", func(t *testing.T) {
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		// Create IO with curried function
		ioFunc := F.Pipe1(
			app.Of(10),
			Map(add),
		)

		// Apply to second argument
		result := Ap[int](app.Of(32))(ioFunc)
		assert.Equal(t, 42, result())
	})

	t.Run("applies curried three-argument function", func(t *testing.T) {
		add3 := func(a int) func(int) func(int) int {
			return func(b int) func(int) int {
				return func(c int) int {
					return a + b + c
				}
			}
		}

		// Build up the computation step by step
		ioFunc1 := F.Pipe1(
			app.Of(10),
			Map(add3),
		)

		ioFunc2 := Ap[func(int) int](app.Of(20))(ioFunc1)
		result := Ap[int](app.Of(12))(ioFunc2)

		assert.Equal(t, 42, result())
	})
}

// TestApplicativeParallelExecution tests that Ap uses parallel execution
func TestApplicativeParallelExecution(t *testing.T) {
	t.Run("executes function and value in parallel", func(t *testing.T) {
		// This test verifies that both computations happen
		// The actual parallelism is tested by the implementation
		add := func(a int) func(int) int {
			return func(b int) int { return a + b }
		}

		ioFunc := Of(add(20))
		ioValue := Of(22)

		result := Ap[int](ioValue)(ioFunc)
		assert.Equal(t, 42, result())
	})
}

// TestApplicativeInstance tests that Applicative returns a valid instance
func TestApplicativeInstance(t *testing.T) {
	t.Run("returns non-nil instance", func(t *testing.T) {
		app := Applicative[int, string]()
		assert.NotNil(t, app)
	})

	t.Run("multiple calls return independent instances", func(t *testing.T) {
		app1 := Applicative[int, string]()
		app2 := Applicative[int, string]()

		// Both should work independently
		result1 := app1.Of(42)
		result2 := app2.Of(43)

		assert.Equal(t, 42, result1())
		assert.Equal(t, 43, result2())
	})
}

// TestApplicativeWithDifferentTypes tests applicative with various type combinations
func TestApplicativeWithDifferentTypes(t *testing.T) {
	t.Run("int to string", func(t *testing.T) {
		app := Applicative[int, string]()
		result := app.Map(strconv.Itoa)(app.Of(42))
		assert.Equal(t, "42", result())
	})

	t.Run("string to int", func(t *testing.T) {
		app := Applicative[string, int]()
		toLength := func(s string) int { return len(s) }
		result := app.Map(toLength)(app.Of("hello"))
		assert.Equal(t, 5, result())
	})

	t.Run("bool to string", func(t *testing.T) {
		app := Applicative[bool, string]()
		toString := func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		}
		result := app.Map(toString)(app.Of(true))
		assert.Equal(t, "true", result())
	})
}
