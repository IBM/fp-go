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
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestMonadExtendWithRight tests MonadExtend with Right values
func TestMonadExtendWithRight(t *testing.T) {
	t.Run("applies function to Right value", func(t *testing.T) {
		input := Right[error](42)

		// Function that extracts and doubles the value if Right
		f := func(e Either[error, int]) int {
			return Fold(
				F.Constant1[error](0),
				N.Mul(2),
			)(e)
		}

		result := MonadExtend(input, f)

		assert.True(t, IsRight(result))
		assert.Equal(t, 84, GetOrElse(F.Constant1[error](0))(result))
	})

	t.Run("function receives entire Either context", func(t *testing.T) {
		input := Right[error]("hello")

		// Function that creates metadata about the Either
		f := func(e Either[error, string]) string {
			return Fold(
				func(err error) string { return "error: " + err.Error() },
				S.Prepend("value: "),
			)(e)
		}

		result := MonadExtend(input, f)

		assert.True(t, IsRight(result))
		assert.Equal(t, "value: hello", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("can count Right occurrences", func(t *testing.T) {
		input := Right[error](100)

		counter := func(e Either[error, int]) int {
			return Fold(
				F.Constant1[error](0),
				F.Constant1[int](1),
			)(e)
		}

		result := MonadExtend(input, counter)

		assert.True(t, IsRight(result))
		assert.Equal(t, 1, GetOrElse(func(error) int { return -1 })(result))
	})
}

// TestMonadExtendWithLeft tests MonadExtend with Left values
func TestMonadExtendWithLeft(t *testing.T) {
	t.Run("returns Left without applying function", func(t *testing.T) {
		testErr := errors.New("test error")
		input := Left[int](testErr)

		// Function should not be called
		called := false
		f := func(e Either[error, int]) int {
			called = true
			return 42
		}

		result := MonadExtend(input, f)

		assert.False(t, called, "function should not be called for Left")
		assert.True(t, IsLeft(result))
		_, leftVal := Unwrap(result)
		assert.Equal(t, testErr, leftVal)
	})

	t.Run("preserves Left error type", func(t *testing.T) {
		input := Left[string](errors.New("original error"))

		f := func(e Either[error, string]) string {
			return "should not be called"
		}

		result := MonadExtend(input, f)

		assert.True(t, IsLeft(result))
		_, leftVal := Unwrap(result)
		assert.Equal(t, "original error", leftVal.Error())
	})
}

// TestMonadExtendEdgeCases tests edge cases for MonadExtend
func TestMonadExtendEdgeCases(t *testing.T) {
	t.Run("function returns zero value", func(t *testing.T) {
		input := Right[error](42)

		f := func(e Either[error, int]) int {
			return 0
		}

		result := MonadExtend(input, f)

		assert.True(t, IsRight(result))
		assert.Equal(t, 0, GetOrElse(func(error) int { return -1 })(result))
	})

	t.Run("function changes type", func(t *testing.T) {
		input := Right[error](42)

		f := func(e Either[error, int]) string {
			return Fold(
				F.Constant1[error]("error"),
				S.Format[int]("number: %d"),
			)(e)
		}

		result := MonadExtend(input, f)

		assert.True(t, IsRight(result))
		assert.Equal(t, "number: 42", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("nested Either handling", func(t *testing.T) {
		inner := Right[error](10)
		outer := Right[error](inner)

		// Extract the inner value
		f := func(e Either[error, Either[error, int]]) int {
			return Fold(
				F.Constant1[error](-1),
				func(innerEither Either[error, int]) int {
					return GetOrElse(F.Constant1[error](-2))(innerEither)
				},
			)(e)
		}

		result := MonadExtend(outer, f)

		assert.True(t, IsRight(result))
		assert.Equal(t, 10, GetOrElse(F.Constant1[error](-3))(result))
	})
}

// TestExtendWithRight tests Extend (curried version) with Right values
func TestExtendWithRight(t *testing.T) {
	t.Run("creates reusable extender", func(t *testing.T) {
		// Create a reusable extender
		doubler := Extend(func(e Either[error, int]) int {
			return Fold(
				F.Constant1[error](0),
				N.Mul(2),
			)(e)
		})

		result1 := doubler(Right[error](21))
		result2 := doubler(Right[error](50))

		assert.True(t, IsRight(result1))
		assert.Equal(t, 42, GetOrElse(F.Constant1[error](0))(result1))

		assert.True(t, IsRight(result2))
		assert.Equal(t, 100, GetOrElse(F.Constant1[error](0))(result2))
	})

	t.Run("metadata extractor", func(t *testing.T) {
		getMetadata := Extend(func(e Either[error, string]) string {
			return Fold(
				func(err error) string { return "error: " + err.Error() },
				S.Prepend("value: "),
			)(e)
		})

		result := getMetadata(Right[error]("test"))

		assert.True(t, IsRight(result))
		assert.Equal(t, "value: test", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("composition with other operations", func(t *testing.T) {
		// Create an extender that counts characters
		charCounter := Extend(func(e Either[error, string]) int {
			return Fold(
				F.Constant1[error](0),
				S.Size,
			)(e)
		})

		// Apply to a Right value
		input := Right[error]("hello")
		result := charCounter(input)

		assert.True(t, IsRight(result))
		assert.Equal(t, 5, GetOrElse(func(error) int { return -1 })(result))
	})
}

// TestExtendWithLeft tests Extend with Left values
func TestExtendWithLeft(t *testing.T) {
	t.Run("returns Left without calling function", func(t *testing.T) {
		testErr := errors.New("test error")

		called := false
		extender := Extend(func(e Either[error, int]) int {
			called = true
			return 42
		})

		result := extender(Left[int](testErr))

		assert.False(t, called, "function should not be called for Left")
		assert.True(t, IsLeft(result))
		_, leftVal := Unwrap(result)
		assert.Equal(t, testErr, leftVal)
	})

	t.Run("preserves error through multiple applications", func(t *testing.T) {
		originalErr := errors.New("original")

		extender := Extend(func(e Either[error, string]) string {
			return "transformed"
		})

		result := extender(Left[string](originalErr))

		assert.True(t, IsLeft(result))
		_, leftVal := Unwrap(result)
		assert.Equal(t, originalErr, leftVal)
	})
}

// TestExtendChaining tests chaining multiple Extend operations
func TestExtendChaining(t *testing.T) {
	t.Run("chain multiple extenders", func(t *testing.T) {
		// First extender: double the value
		doubler := Extend(func(e Either[error, int]) int {
			return Fold(
				F.Constant1[error](0),
				N.Mul(2),
			)(e)
		})

		// Second extender: add 10
		adder := Extend(func(e Either[error, int]) int {
			return Fold(
				F.Constant1[error](0),
				N.Add(10),
			)(e)
		})

		input := Right[error](5)
		result := adder(doubler(input))

		assert.True(t, IsRight(result))
		assert.Equal(t, 20, GetOrElse(F.Constant1[error](0))(result))
	})

	t.Run("short-circuits on Left", func(t *testing.T) {
		testErr := errors.New("error")

		extender1 := Extend(func(e Either[error, int]) int { return 1 })
		extender2 := Extend(func(e Either[error, int]) int { return 2 })

		input := Left[int](testErr)
		result := extender2(extender1(input))

		assert.True(t, IsLeft(result))
		_, leftVal := Unwrap(result)
		assert.Equal(t, testErr, leftVal)
	})
}

// TestExtendTypeTransformations tests type transformations with Extend
func TestExtendTypeTransformations(t *testing.T) {
	t.Run("int to string transformation", func(t *testing.T) {
		toString := Extend(func(e Either[error, int]) string {
			return Fold(
				F.Constant1[error]("error"),
				strconv.Itoa,
			)(e)
		})

		result := toString(Right[error](42))

		assert.True(t, IsRight(result))
		assert.Equal(t, "42", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("string to bool transformation", func(t *testing.T) {
		isEmpty := Extend(func(e Either[error, string]) bool {
			return Fold(
				F.Constant1[error](true),
				S.IsEmpty,
			)(e)
		})

		result1 := isEmpty(Right[error](""))
		result2 := isEmpty(Right[error]("hello"))

		assert.True(t, IsRight(result1))
		assert.True(t, GetOrElse(F.Constant1[error](false))(result1))

		assert.True(t, IsRight(result2))
		assert.False(t, GetOrElse(F.Constant1[error](true))(result2))
	})
}

// TestExtendWithComplexTypes tests Extend with complex types
func TestExtendWithComplexTypes(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	t.Run("extract field from struct", func(t *testing.T) {
		getName := Extend(func(e Either[error, User]) string {
			return Fold(
				func(err error) string { return "unknown" },
				func(u User) string { return u.Name },
			)(e)
		})

		user := User{Name: "Alice", Age: 30}
		result := getName(Right[error](user))

		assert.True(t, IsRight(result))
		assert.Equal(t, "Alice", GetOrElse(func(error) string { return "" })(result))
	})

	t.Run("compute derived value", func(t *testing.T) {
		isAdult := Extend(func(e Either[error, User]) bool {
			return Fold(
				func(err error) bool { return false },
				func(u User) bool { return u.Age >= 18 },
			)(e)
		})

		user1 := User{Name: "Bob", Age: 25}
		user2 := User{Name: "Charlie", Age: 15}

		result1 := isAdult(Right[error](user1))
		result2 := isAdult(Right[error](user2))

		assert.True(t, IsRight(result1))
		assert.True(t, GetOrElse(F.Constant1[error](false))(result1))

		assert.True(t, IsRight(result2))
		assert.False(t, GetOrElse(F.Constant1[error](true))(result2))
	})
}
