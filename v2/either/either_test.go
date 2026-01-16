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

	assert.True(t, IsRight(e), "Zero should create a Right value")
	assert.False(t, IsLeft(e), "Zero should not create a Left value")

	value, err := Unwrap(e)
	assert.Equal(t, 0, value, "Right value should be zero for int")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithStrings tests Zero function with string types
func TestZeroWithStrings(t *testing.T) {
	e := Zero[error, string]()

	assert.True(t, IsRight(e), "Zero should create a Right value")
	assert.False(t, IsLeft(e), "Zero should not create a Left value")

	value, err := Unwrap(e)
	assert.Equal(t, "", value, "Right value should be empty string")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithBooleans tests Zero function with boolean types
func TestZeroWithBooleans(t *testing.T) {
	e := Zero[error, bool]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Equal(t, false, value, "Right value should be false for bool")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithFloats tests Zero function with float types
func TestZeroWithFloats(t *testing.T) {
	e := Zero[error, float64]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Equal(t, 0.0, value, "Right value should be 0.0 for float64")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithPointers tests Zero function with pointer types
func TestZeroWithPointers(t *testing.T) {
	e := Zero[error, *int]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Nil(t, value, "Right value should be nil for pointer type")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithSlices tests Zero function with slice types
func TestZeroWithSlices(t *testing.T) {
	e := Zero[error, []int]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Nil(t, value, "Right value should be nil for slice type")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithMaps tests Zero function with map types
func TestZeroWithMaps(t *testing.T) {
	e := Zero[error, map[string]int]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Nil(t, value, "Right value should be nil for map type")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithStructs tests Zero function with struct types
func TestZeroWithStructs(t *testing.T) {
	type TestStruct struct {
		Field1 int
		Field2 string
	}

	e := Zero[error, TestStruct]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	expected := TestStruct{Field1: 0, Field2: ""}
	assert.Equal(t, expected, value, "Right value should be zero value for struct")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithInterfaces tests Zero function with interface types
func TestZeroWithInterfaces(t *testing.T) {
	e := Zero[error, interface{}]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.Nil(t, value, "Right value should be nil for interface type")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithCustomErrorType tests Zero function with custom error types
func TestZeroWithCustomErrorType(t *testing.T) {
	type CustomError struct {
		Code    int
		Message string
	}

	e := Zero[CustomError, string]()

	assert.True(t, IsRight(e), "Zero should create a Right value")
	assert.False(t, IsLeft(e), "Zero should not create a Left value")

	value, err := Unwrap(e)
	assert.Equal(t, "", value, "Right value should be empty string")
	assert.Equal(t, CustomError{Code: 0, Message: ""}, err, "Error should be zero value for CustomError")
}

// TestZeroCanBeUsedWithOtherFunctions tests that Zero Eithers work with other either functions
func TestZeroCanBeUsedWithOtherFunctions(t *testing.T) {
	e := Zero[error, int]()

	// Test with Map
	mapped := MonadMap(e, func(n int) string {
		return fmt.Sprintf("%d", n)
	})
	assert.True(t, IsRight(mapped), "Mapped Zero should still be Right")
	value, _ := Unwrap(mapped)
	assert.Equal(t, "0", value, "Mapped value should be '0'")

	// Test with Chain
	chained := MonadChain(e, func(n int) Either[error, string] {
		return Right[error](fmt.Sprintf("value: %d", n))
	})
	assert.True(t, IsRight(chained), "Chained Zero should still be Right")
	chainedValue, _ := Unwrap(chained)
	assert.Equal(t, "value: 0", chainedValue, "Chained value should be 'value: 0'")

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

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	expected := ComplexType{Nested: nil, Ptr: nil}
	assert.Equal(t, expected, value, "Right value should be zero value for complex struct")
	assert.Nil(t, err, "Error should be nil for Right value")
}

// TestZeroWithOption tests Zero with Option type
func TestZeroWithOption(t *testing.T) {
	e := Zero[error, O.Option[int]]()

	assert.True(t, IsRight(e), "Zero should create a Right value")

	value, err := Unwrap(e)
	assert.True(t, O.IsNone(value), "Right value should be None for Option type")
	assert.Nil(t, err, "Error should be nil for Right value")
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
