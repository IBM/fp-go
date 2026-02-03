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
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
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
	noError := Right[error]("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {

	assert.Equal(t, F.Pipe1(Right[error]("abc"), Map[error](utils.StringLen)), Right[error](3))

	val2 := F.Pipe1(Left[string]("s"), Map[string](utils.StringLen))
	exp2 := Left[int]("s")

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

	assert.Equal(t, "foobar", F.Pipe1(Right[string]("bar"), Reduce[string](s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string]("bar"), Reduce[string](s.Concat, "foo")))

}
func TestAp(t *testing.T) {
	f := S.Size

	assert.Equal(t, Right[string](3), F.Pipe1(Right[string](f), Ap[int](Right[string]("abc"))))
	assert.Equal(t, Left[int]("maError"), F.Pipe1(Right[string](f), Ap[int](Left[string]("maError"))))
	assert.Equal(t, Left[int]("mabError"), F.Pipe1(Left[func(string) int]("mabError"), Ap[int](Left[string]("maError"))))
}

func TestAlt(t *testing.T) {
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Left[int]("a")))))
	assert.Equal(t, Right[string](2), F.Pipe1(Left[int]("b"), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Left[int]("b"), F.Pipe1(Left[int]("a"), Alt(F.Constant(Left[int]("b")))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[string, int])

	assert.Equal(t, Right[string]("abc"), F.Pipe1(Right[string]("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string]("maError"), F.Pipe1(Left[string]("maError"), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("a"))(func(n int) Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right[string](1), f(Right[string](1)))
	assert.Equal(t, Left[int]("a"), f(Right[string](-1)))
	assert.Equal(t, Left[int]("b"), f(Left[int]("b")))
}

func TestFromOption(t *testing.T) {
	assert.Equal(t, Left[int]("none"), FromOption[int](F.Constant("none"))(O.None[int]()))
	assert.Equal(t, Right[string](1), FromOption[int](F.Constant("none"))(O.Some(1)))
}

func TestStringer(t *testing.T) {
	e := Of[error]("foo")
	exp := "Right[string](foo)"

	assert.Equal(t, exp, e.String())

	var s fmt.Stringer = &e
	assert.Equal(t, exp, s.String())
}

// TestZeroWithIntegers tests Zero function with integer types
func TestZeroWithIntegers(t *testing.T) {
	e := Zero[error, int]()

	assert.Equal(t, Of[error](0), e, "Zero should create a Right value with zero for int")
}

// TestZeroWithStrings tests Zero function with string types
func TestZeroWithStrings(t *testing.T) {
	e := Zero[error, string]()

	assert.Equal(t, Of[error](""), e, "Zero should create a Right value with empty string")
}

// TestZeroWithBooleans tests Zero function with boolean types
func TestZeroWithBooleans(t *testing.T) {
	e := Zero[error, bool]()

	assert.Equal(t, Of[error](false), e, "Zero should create a Right value with false for bool")
}

// TestZeroWithFloats tests Zero function with float types
func TestZeroWithFloats(t *testing.T) {
	e := Zero[error, float64]()

	assert.Equal(t, Of[error](0.0), e, "Zero should create a Right value with 0.0 for float64")
}

// TestZeroWithPointers tests Zero function with pointer types
func TestZeroWithPointers(t *testing.T) {
	e := Zero[error, *int]()

	var nilPtr *int
	assert.Equal(t, Of[error](nilPtr), e, "Zero should create a Right value with nil pointer")
}

// TestZeroWithSlices tests Zero function with slice types
func TestZeroWithSlices(t *testing.T) {
	e := Zero[error, []int]()

	var nilSlice []int
	assert.Equal(t, Of[error](nilSlice), e, "Zero should create a Right value with nil slice")
}

// TestZeroWithMaps tests Zero function with map types
func TestZeroWithMaps(t *testing.T) {
	e := Zero[error, map[string]int]()

	var nilMap map[string]int
	assert.Equal(t, Of[error](nilMap), e, "Zero should create a Right value with nil map")
}

// TestZeroWithStructs tests Zero function with struct types
func TestZeroWithStructs(t *testing.T) {
	type TestStruct struct {
		Field1 int
		Field2 string
	}

	e := Zero[error, TestStruct]()

	expected := TestStruct{Field1: 0, Field2: ""}
	assert.Equal(t, Of[error](expected), e, "Zero should create a Right value with zero value for struct")
}

// TestZeroWithInterfaces tests Zero function with interface types
func TestZeroWithInterfaces(t *testing.T) {
	e := Zero[error, interface{}]()

	var nilInterface interface{}
	assert.Equal(t, Of[error](nilInterface), e, "Zero should create a Right value with nil interface")
}

// TestZeroWithCustomErrorType tests Zero function with custom error types
func TestZeroWithCustomErrorType(t *testing.T) {
	type CustomError struct {
		Code    int
		Message string
	}

	e := Zero[CustomError, string]()

	assert.Equal(t, Of[CustomError](""), e, "Zero should create a Right value with empty string")
}

// TestZeroCanBeUsedWithOtherFunctions tests that Zero Eithers work with other either functions
func TestZeroCanBeUsedWithOtherFunctions(t *testing.T) {
	e := Zero[error, int]()

	// Test with Map
	mapped := MonadMap(e, func(n int) string {
		return fmt.Sprintf("%d", n)
	})
	assert.Equal(t, Of[error]("0"), mapped, "Mapped Zero should be Right with '0'")

	// Test with Chain
	chained := MonadChain(e, func(n int) Either[error, string] {
		return Right[error](fmt.Sprintf("value: %d", n))
	})
	assert.Equal(t, Of[error]("value: 0"), chained, "Chained Zero should be Right with 'value: 0'")

	// Test with Fold
	folded := MonadFold(e,
		func(err error) string { return "error" },
		func(n int) string { return fmt.Sprintf("success: %d", n) },
	)
	assert.Equal(t, "success: 0", folded, "Folded value should be 'success: 0'")
}

// TestZeroEquality tests that multiple Zero calls produce equal Eithers
func TestZeroEquality(t *testing.T) {
	e1 := Zero[error, int]()
	e2 := Zero[error, int]()

	assert.Equal(t, IsRight(e1), IsRight(e2), "Both should be Right")
	assert.Equal(t, IsLeft(e1), IsLeft(e2), "Both should not be Left")

	v1, err1 := Unwrap(e1)
	v2, err2 := Unwrap(e2)
	assert.Equal(t, v1, v2, "Values should be equal")
	assert.Equal(t, err1, err2, "Errors should be equal")
}

// TestZeroWithComplexTypes tests Zero with more complex nested types
func TestZeroWithComplexTypes(t *testing.T) {
	type ComplexType struct {
		Nested map[string][]int
		Ptr    *string
	}

	e := Zero[error, ComplexType]()

	expected := ComplexType{Nested: nil, Ptr: nil}
	assert.Equal(t, Of[error](expected), e, "Zero should create a Right value with zero value for complex struct")
}

// TestZeroWithOption tests Zero with Option type
func TestZeroWithOption(t *testing.T) {
	e := Zero[error, O.Option[int]]()

	assert.Equal(t, Of[error](O.None[int]()), e, "Zero should create a Right value with None option")
}

// TestZeroIsNotLeft tests that Zero never creates a Left value
func TestZeroIsNotLeft(t *testing.T) {
	// Test with various type combinations
	e1 := Zero[string, int]()
	e2 := Zero[error, string]()
	e3 := Zero[int, bool]()

	assert.False(t, IsLeft(e1), "Zero should never create a Left value")
	assert.False(t, IsLeft(e2), "Zero should never create a Left value")
	assert.False(t, IsLeft(e3), "Zero should never create a Left value")

	assert.True(t, IsRight(e1), "Zero should always create a Right value")
	assert.True(t, IsRight(e2), "Zero should always create a Right value")
	assert.True(t, IsRight(e3), "Zero should always create a Right value")
}

// TestZeroEqualsDefaultInitialization tests that Zero returns the same value as default initialization
func TestZeroEqualsDefaultInitialization(t *testing.T) {
	// Default initialization of Either
	var defaultInit Either[error, int]

	// Zero function
	zero := Zero[error, int]()

	// They should be equal
	assert.Equal(t, defaultInit, zero, "Zero should equal default initialization")
	assert.Equal(t, IsRight(defaultInit), IsRight(zero), "Both should be Right")
	assert.Equal(t, IsLeft(defaultInit), IsLeft(zero), "Both should not be Left")
}

// TestMonadChainLeft tests the MonadChainLeft function with various scenarios
func TestMonadChainLeft(t *testing.T) {
	t.Run("Left value is transformed by function", func(t *testing.T) {
		// Transform error to success
		result := MonadChainLeft(
			Left[int](errors.New("not found")),
			func(err error) Either[string, int] {
				if err.Error() == "not found" {
					return Right[string](0) // default value
				}
				return Left[int](err.Error())
			},
		)
		assert.Equal(t, Of[string](0), result)
	})

	t.Run("Left value error type is transformed", func(t *testing.T) {
		// Transform error type from int to string
		result := MonadChainLeft(
			Left[int](404),
			func(code int) Either[string, int] {
				return Left[int](fmt.Sprintf("Error code: %d", code))
			},
		)
		assert.Equal(t, Left[int]("Error code: 404"), result)
	})

	t.Run("Right value passes through unchanged", func(t *testing.T) {
		// Right value should not be affected
		result := MonadChainLeft(
			Right[error](42),
			func(err error) Either[string, int] {
				return Left[int]("should not be called")
			},
		)
		assert.Equal(t, Of[string](42), result)
	})

	t.Run("Chain multiple error transformations", func(t *testing.T) {
		// First transformation
		step1 := MonadChainLeft(
			Left[int](errors.New("error1")),
			func(err error) Either[error, int] {
				return Left[int](errors.New("error2"))
			},
		)
		// Second transformation
		step2 := MonadChainLeft(
			step1,
			func(err error) Either[string, int] {
				return Left[int](err.Error())
			},
		)
		assert.Equal(t, Left[int]("error2"), step2)
	})

	t.Run("Error recovery with fallback", func(t *testing.T) {
		// Recover from specific errors
		result := MonadChainLeft(
			Left[int](errors.New("timeout")),
			func(err error) Either[error, int] {
				if err.Error() == "timeout" {
					return Right[error](999) // fallback value
				}
				return Left[int](err)
			},
		)
		assert.Equal(t, Of[error](999), result)
	})

	t.Run("Transform error to different Left", func(t *testing.T) {
		// Transform one error to another
		result := MonadChainLeft(
			Left[string]("original error"),
			func(s string) Either[int, string] {
				return Left[string](len(s))
			},
		)
		assert.Equal(t, Left[string](14), result) // length of "original error"
	})
}

// TestChainLeft tests the curried ChainLeft function
func TestChainLeft(t *testing.T) {
	t.Run("Curried function transforms Left value", func(t *testing.T) {
		// Create a reusable error handler
		handleNotFound := ChainLeft[error, string](func(err error) Either[string, int] {
			if err.Error() == "not found" {
				return Right[string](0)
			}
			return Left[int](err.Error())
		})

		result := handleNotFound(Left[int](errors.New("not found")))
		assert.Equal(t, Of[string](0), result)
	})

	t.Run("Curried function with Right value", func(t *testing.T) {
		handler := ChainLeft[error, string](func(err error) Either[string, int] {
			return Left[int]("should not be called")
		})

		result := handler(Right[error](42))
		assert.Equal(t, Of[string](42), result)
	})

	t.Run("Use in pipeline with Pipe", func(t *testing.T) {
		// Create error transformer
		toStringError := ChainLeft[int, string](func(code int) Either[string, string] {
			return Left[string](fmt.Sprintf("Error: %d", code))
		})

		result := F.Pipe1(
			Left[string](404),
			toStringError,
		)
		assert.Equal(t, Left[string]("Error: 404"), result)
	})

	t.Run("Compose multiple ChainLeft operations", func(t *testing.T) {
		// First handler: convert error to string
		handler1 := ChainLeft[error, string](func(err error) Either[string, int] {
			return Left[int](err.Error())
		})

		// Second handler: add prefix to string error
		handler2 := ChainLeft[string, string](func(s string) Either[string, int] {
			return Left[int]("Handled: " + s)
		})

		result := F.Pipe2(
			Left[int](errors.New("original")),
			handler1,
			handler2,
		)
		assert.Equal(t, Left[int]("Handled: original"), result)
	})

	t.Run("Error recovery in pipeline", func(t *testing.T) {
		// Handler that recovers from specific errors
		recoverFromTimeout := ChainLeft(func(err error) Either[error, int] {
			if err.Error() == "timeout" {
				return Right[error](0) // recovered value
			}
			return Left[int](err) // propagate other errors
		})

		// Test with timeout error
		result1 := F.Pipe1(
			Left[int](errors.New("timeout")),
			recoverFromTimeout,
		)
		assert.Equal(t, Of[error](0), result1)

		// Test with other error
		result2 := F.Pipe1(
			Left[int](errors.New("other error")),
			recoverFromTimeout,
		)
		assert.True(t, IsLeft(result2))
	})

	t.Run("Transform error type in pipeline", func(t *testing.T) {
		// Convert numeric error codes to descriptive strings
		codeToMessage := ChainLeft(func(code int) Either[string, string] {
			messages := map[int]string{
				404: "Not Found",
				500: "Internal Server Error",
			}
			if msg, ok := messages[code]; ok {
				return Left[string](msg)
			}
			return Left[string](fmt.Sprintf("Unknown error: %d", code))
		})

		result := F.Pipe1(
			Left[string](404),
			codeToMessage,
		)
		assert.Equal(t, Left[string]("Not Found"), result)
	})

	t.Run("ChainLeft with Map combination", func(t *testing.T) {
		// Combine ChainLeft with Map to handle both channels
		errorHandler := ChainLeft(func(err error) Either[string, int] {
			return Left[int]("Error: " + err.Error())
		})

		valueMapper := Map[string](S.Format[int]("Value: %d"))

		// Test with Left
		result1 := F.Pipe2(
			Left[int](errors.New("fail")),
			errorHandler,
			valueMapper,
		)
		assert.Equal(t, Left[string]("Error: fail"), result1)

		// Test with Right
		result2 := F.Pipe2(
			Right[error](42),
			errorHandler,
			valueMapper,
		)
		assert.Equal(t, Of[string]("Value: 42"), result2)
	})
}
