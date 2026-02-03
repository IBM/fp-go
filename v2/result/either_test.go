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

package result

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError := Left[string](err)

	assert.True(t, IsLeft(withError))
	assert.False(t, IsRight(withError))
}

func TestIsRight(t *testing.T) {
	noError := Right("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {
	e := errors.New("s")
	assert.Equal(t, F.Pipe1(Right("abc"), Map(utils.StringLen)), Right(3))

	val2 := F.Pipe1(Left[string](e), Map(utils.StringLen))
	exp2 := Left[int](e)

	assert.Equal(t, val2, exp2)
}

func TestUnwrapError(t *testing.T) {
	a := ""
	err := errors.New("Some error")
	withError := Left[string](err)

	res, extracted := UnwrapError(withError)
	assert.Equal(t, a, res)
	assert.Equal(t, extracted, err)

}

func TestReduce(t *testing.T) {

	s := S.Semigroup

	assert.Equal(t, "foobar", F.Pipe1(Right("bar"), Reduce(s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string](errors.New("bar")), Reduce(s.Concat, "foo")))

}
func TestAp(t *testing.T) {
	f := S.Size

	maError := errors.New("maError")
	mabError := errors.New("mabError")

	assert.Equal(t, Right(3), F.Pipe1(Right(f), Ap[int](Right("abc"))))
	assert.Equal(t, Left[int](maError), F.Pipe1(Right(f), Ap[int](Left[string](maError))))
	assert.Equal(t, Left[int](mabError), F.Pipe1(Left[func(string) int](mabError), Ap[int](Left[string](maError))))
}

func TestAlt(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Left[int](a)))))
	assert.Equal(t, Right(2), F.Pipe1(Left[int](b), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Left[int](b), F.Pipe1(Left[int](a), Alt(F.Constant(Left[int](b)))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[int])
	maError := errors.New("maError")

	assert.Equal(t, Right("abc"), F.Pipe1(Right("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string](maError), F.Pipe1(Left[string](maError), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")

	f := ChainOptionK[int, int](F.Constant(a))(func(n int) Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right(1), f(Right(1)))
	assert.Equal(t, Left[int](a), f(Right(-1)))
	assert.Equal(t, Left[int](b), f(Left[int](b)))
}

func TestFromOption(t *testing.T) {
	none := errors.New("none")

	assert.Equal(t, Left[int](none), FromOption[int](F.Constant(none))(O.None[int]()))
	assert.Equal(t, Right(1), FromOption[int](F.Constant(none))(O.Some(1)))
}

func TestStringer(t *testing.T) {
	e := Of("foo")
	exp := "Right[string](foo)"

	assert.Equal(t, exp, e.String())

	var s fmt.Stringer = &e
	assert.Equal(t, exp, s.String())
}

func TestFromIO(t *testing.T) {
	f := IO.Of("abc")
	e := FromIO(f)

	assert.Equal(t, Right("abc"), e)
}

// TestOrElse tests recovery from error
func TestOrElse(t *testing.T) {
	// Test basic recovery from Left
	recover := OrElse(func(err error) Result[int] {
		return Right(0) // default value
	})

	leftResult := Left[int](errors.New("fail"))
	assert.Equal(t, Right(0), recover(leftResult))

	// Test that Right values pass through unchanged
	rightResult := Right(42)
	assert.Equal(t, Right(42), recover(rightResult))

	// Test conditional recovery
	recoverSpecific := OrElse(func(err error) Result[int] {
		if err.Error() == "not found" {
			return Right(0) // default for not found
		}
		return Left[int](err) // propagate other errors
	})

	notFoundErr := errors.New("not found")
	assert.Equal(t, Right(0), recoverSpecific(Left[int](notFoundErr)))

	otherErr := errors.New("other error")
	assert.Equal(t, Left[int](otherErr), recoverSpecific(Left[int](otherErr)))
}

// TestZeroEqualsDefaultInitialization tests that Zero returns the same value as default initialization
func TestZeroEqualsDefaultInitialization(t *testing.T) {
	// Default initialization of Result
	var defaultInit Result[int]

	// Zero function
	zero := Zero[int]()

	// They should be equal
	assert.Equal(t, defaultInit, zero, "Zero should equal default initialization")
	assert.Equal(t, IsRight(defaultInit), IsRight(zero), "Both should be Right")
	assert.Equal(t, IsLeft(defaultInit), IsLeft(zero), "Both should not be Left")
}

// TestInstanceOf tests the InstanceOf function for type assertions
func TestInstanceOf(t *testing.T) {
	// Test successful type assertion with int
	t.Run("successful int assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[int](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(42), result)
	})

	// Test successful type assertion with string
	t.Run("successful string assertion", func(t *testing.T) {
		var value any = "hello"
		result := InstanceOf[string](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right("hello"), result)
	})

	// Test successful type assertion with float64
	t.Run("successful float64 assertion", func(t *testing.T) {
		var value any = 3.14
		result := InstanceOf[float64](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(3.14), result)
	})

	// Test successful type assertion with pointer
	t.Run("successful pointer assertion", func(t *testing.T) {
		val := 42
		var value any = &val
		result := InstanceOf[*int](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, 42, *v)
	})

	// Test successful type assertion with struct
	t.Run("successful struct assertion", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		var value any = Person{Name: "Alice", Age: 30}
		result := InstanceOf[Person](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(Person{Name: "Alice", Age: 30}), result)
	})

	// Test failed type assertion - int to string
	t.Run("failed int to string assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[string](value)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected")
		assert.Contains(t, err.Error(), "got")
	})

	// Test failed type assertion - string to int
	t.Run("failed string to int assertion", func(t *testing.T) {
		var value any = "hello"
		result := InstanceOf[int](value)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Error(t, err)
	})

	// Test failed type assertion - int to float64
	t.Run("failed int to float64 assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[float64](value)
		assert.True(t, IsLeft(result))
	})

	// Test with nil value
	t.Run("nil value assertion", func(t *testing.T) {
		var value any = nil
		result := InstanceOf[string](value)
		assert.True(t, IsLeft(result))
	})

	// Test chaining with Map
	t.Run("chaining with Map", func(t *testing.T) {
		var value any = 10
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Map(func(n int) int { return n * 2 }),
		)
		assert.Equal(t, Right(20), result)
	})

	// Test chaining with Map on failed assertion
	t.Run("chaining with Map on failed assertion", func(t *testing.T) {
		var value any = "not a number"
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Map(func(n int) int { return n * 2 }),
		)
		assert.True(t, IsLeft(result))
	})

	// Test with Chain for dependent operations
	t.Run("chaining with Chain", func(t *testing.T) {
		var value any = 5
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Chain(func(n int) Result[string] {
				if n > 0 {
					return Right(fmt.Sprintf("positive: %d", n))
				}
				return Left[string](errors.New("not positive"))
			}),
		)
		assert.Equal(t, Right("positive: 5"), result)
	})

	// Test with GetOrElse for default value
	t.Run("GetOrElse with failed assertion", func(t *testing.T) {
		var value any = "not an int"
		result := F.Pipe2(
			value,
			InstanceOf[int],
			GetOrElse(func(err error) int { return -1 }),
		)
		assert.Equal(t, -1, result)
	})

	// Test with GetOrElse for successful assertion
	t.Run("GetOrElse with successful assertion", func(t *testing.T) {
		var value any = 42
		result := F.Pipe2(
			value,
			InstanceOf[int],
			GetOrElse(func(err error) int { return -1 }),
		)
		assert.Equal(t, 42, result)
	})

	// Test with interface type
	t.Run("interface type assertion", func(t *testing.T) {
		var value any = errors.New("test error")
		result := InstanceOf[error](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, "test error", v.Error())
	})

	// Test with slice type
	t.Run("slice type assertion", func(t *testing.T) {
		var value any = []int{1, 2, 3}
		result := InstanceOf[[]int](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right([]int{1, 2, 3}), result)
	})

	// Test with map type
	t.Run("map type assertion", func(t *testing.T) {
		var value any = map[string]int{"a": 1, "b": 2}
		result := InstanceOf[map[string]int](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, 1, v["a"])
		assert.Equal(t, 2, v["b"])
	})
}

// TestMonadChainLeft tests the MonadChainLeft function with various scenarios
func TestMonadChainLeft(t *testing.T) {
	t.Run("Left value is transformed by function", func(t *testing.T) {
		// Transform error to success
		result := MonadChainLeft(
			Left[int](errors.New("not found")),
			func(err error) Result[int] {
				if err.Error() == "not found" {
					return Right(0) // default value
				}
				return Left[int](err)
			},
		)
		assert.Equal(t, Of(0), result)
	})

	t.Run("Left value error is transformed", func(t *testing.T) {
		// Transform error with additional context
		result := MonadChainLeft(
			Left[int](errors.New("database error")),
			func(err error) Result[int] {
				return Left[int](fmt.Errorf("wrapped: %w", err))
			},
		)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Contains(t, err.Error(), "wrapped:")
		assert.Contains(t, err.Error(), "database error")
	})

	t.Run("Right value passes through unchanged", func(t *testing.T) {
		// Right value should not be affected
		result := MonadChainLeft(
			Right(42),
			func(err error) Result[int] {
				return Left[int](errors.New("should not be called"))
			},
		)
		assert.Equal(t, Of(42), result)
	})

	t.Run("Chain multiple error transformations", func(t *testing.T) {
		// First transformation
		step1 := MonadChainLeft(
			Left[int](errors.New("error1")),
			func(err error) Result[int] {
				return Left[int](errors.New("error2"))
			},
		)
		// Second transformation
		step2 := MonadChainLeft(
			step1,
			func(err error) Result[int] {
				return Left[int](fmt.Errorf("final: %s", err.Error()))
			},
		)
		assert.True(t, IsLeft(step2))
		_, err := UnwrapError(step2)
		assert.Equal(t, "final: error2", err.Error())
	})

	t.Run("Error recovery with fallback", func(t *testing.T) {
		// Recover from specific errors
		result := MonadChainLeft(
			Left[int](errors.New("timeout")),
			func(err error) Result[int] {
				if err.Error() == "timeout" {
					return Right(999) // fallback value
				}
				return Left[int](err)
			},
		)
		assert.Equal(t, Of(999), result)
	})

	t.Run("Conditional error handling", func(t *testing.T) {
		// Handle different error types differently
		handleError := func(err error) Result[string] {
			switch err.Error() {
			case "not found":
				return Right("default")
			case "timeout":
				return Right("retry")
			default:
				return Left[string](err)
			}
		}

		result1 := MonadChainLeft(Left[string](errors.New("not found")), handleError)
		assert.Equal(t, Of("default"), result1)

		result2 := MonadChainLeft(Left[string](errors.New("timeout")), handleError)
		assert.Equal(t, Of("retry"), result2)

		result3 := MonadChainLeft(Left[string](errors.New("other")), handleError)
		assert.True(t, IsLeft(result3))
	})

	t.Run("Type preservation", func(t *testing.T) {
		// Ensure type is preserved through transformation
		result := MonadChainLeft(
			Left[string](errors.New("error")),
			func(err error) Result[string] {
				return Right("recovered")
			},
		)
		assert.Equal(t, Of("recovered"), result)
	})
}

// TestChainLeft tests the curried ChainLeft function
func TestChainLeft(t *testing.T) {
	t.Run("Curried function transforms Left value", func(t *testing.T) {
		// Create a reusable error handler
		handleNotFound := ChainLeft(func(err error) Result[int] {
			if err.Error() == "not found" {
				return Right(0)
			}
			return Left[int](err)
		})

		result := handleNotFound(Left[int](errors.New("not found")))
		assert.Equal(t, Of(0), result)
	})

	t.Run("Curried function with Right value", func(t *testing.T) {
		handler := ChainLeft(func(err error) Result[int] {
			return Left[int](errors.New("should not be called"))
		})

		result := handler(Right(42))
		assert.Equal(t, Of(42), result)
	})

	t.Run("Use in pipeline with Pipe", func(t *testing.T) {
		// Create error transformer
		wrapError := ChainLeft(func(err error) Result[string] {
			return Left[string](fmt.Errorf("Error: %w", err))
		})

		result := F.Pipe1(
			Left[string](errors.New("failed")),
			wrapError,
		)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Contains(t, err.Error(), "Error:")
		assert.Contains(t, err.Error(), "failed")
	})

	t.Run("Compose multiple ChainLeft operations", func(t *testing.T) {
		// First handler: convert error to string representation
		handler1 := ChainLeft(func(err error) Result[int] {
			return Left[int](errors.New(err.Error()))
		})

		// Second handler: add prefix to error
		handler2 := ChainLeft(func(err error) Result[int] {
			return Left[int](fmt.Errorf("Handled: %w", err))
		})

		result := F.Pipe2(
			Left[int](errors.New("original")),
			handler1,
			handler2,
		)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Contains(t, err.Error(), "Handled:")
		assert.Contains(t, err.Error(), "original")
	})

	t.Run("Error recovery in pipeline", func(t *testing.T) {
		// Handler that recovers from specific errors
		recoverFromTimeout := ChainLeft(func(err error) Result[int] {
			if err.Error() == "timeout" {
				return Right(0) // recovered value
			}
			return Left[int](err) // propagate other errors
		})

		// Test with timeout error
		result1 := F.Pipe1(
			Left[int](errors.New("timeout")),
			recoverFromTimeout,
		)
		assert.Equal(t, Of(0), result1)

		// Test with other error
		result2 := F.Pipe1(
			Left[int](errors.New("other error")),
			recoverFromTimeout,
		)
		assert.True(t, IsLeft(result2))
	})

	t.Run("ChainLeft with Map combination", func(t *testing.T) {
		// Combine ChainLeft with Map to handle both channels
		errorHandler := ChainLeft(func(err error) Result[int] {
			return Left[int](fmt.Errorf("Error: %w", err))
		})

		valueMapper := Map(func(n int) string {
			return fmt.Sprintf("Value: %d", n)
		})

		// Test with Left
		result1 := F.Pipe2(
			Left[int](errors.New("fail")),
			errorHandler,
			valueMapper,
		)
		assert.True(t, IsLeft(result1))

		// Test with Right
		result2 := F.Pipe2(
			Right(42),
			errorHandler,
			valueMapper,
		)
		assert.Equal(t, Of("Value: 42"), result2)
	})

	t.Run("Reusable error handlers", func(t *testing.T) {
		// Create a library of reusable error handlers
		recoverNotFound := ChainLeft(func(err error) Result[string] {
			if err.Error() == "not found" {
				return Right("default")
			}
			return Left[string](err)
		})

		recoverTimeout := ChainLeft(func(err error) Result[string] {
			if err.Error() == "timeout" {
				return Right("retry")
			}
			return Left[string](err)
		})

		// Apply handlers in sequence
		result := F.Pipe2(
			Left[string](errors.New("not found")),
			recoverNotFound,
			recoverTimeout,
		)
		assert.Equal(t, Of("default"), result)
	})

	t.Run("Error transformation pipeline", func(t *testing.T) {
		// Build a pipeline that transforms errors step by step
		addContext := ChainLeft(func(err error) Result[int] {
			return Left[int](fmt.Errorf("context: %w", err))
		})

		addTimestamp := ChainLeft(func(err error) Result[int] {
			return Left[int](fmt.Errorf("[2024-01-01] %w", err))
		})

		result := F.Pipe2(
			Left[int](errors.New("base error")),
			addContext,
			addTimestamp,
		)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Contains(t, err.Error(), "[2024-01-01]")
		assert.Contains(t, err.Error(), "context:")
		assert.Contains(t, err.Error(), "base error")
	})

	t.Run("Conditional recovery based on error content", func(t *testing.T) {
		// Recover from errors matching specific patterns
		smartRecover := ChainLeft(func(err error) Result[int] {
			msg := err.Error()
			if msg == "not found" {
				return Right(0)
			}
			if msg == "timeout" {
				return Right(-1)
			}
			if msg == "unauthorized" {
				return Right(-2)
			}
			return Left[int](err)
		})

		assert.Equal(t, Of(0), smartRecover(Left[int](errors.New("not found"))))
		assert.Equal(t, Of(-1), smartRecover(Left[int](errors.New("timeout"))))
		assert.Equal(t, Of(-2), smartRecover(Left[int](errors.New("unauthorized"))))
		assert.True(t, IsLeft(smartRecover(Left[int](errors.New("unknown")))))
	})
}
