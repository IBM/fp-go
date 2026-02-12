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

package either

import (
	"math"
	"strconv"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	t.Run("Right value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, "not positive")
		input := Right[string](5)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be Left")
		assert.True(t, IsRight(right), "right should be Right")

		_, leftVal := Unwrap(left)
		rightVal, _ := Unwrap(right)

		assert.Equal(t, "not positive", leftVal)
		assert.Equal(t, 5, rightVal)
	})

	t.Run("Right value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, "not positive")
		input := Right[string](-3)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Right (failed predicate)")
		assert.True(t, IsLeft(right), "right should be Left")

		leftVal, _ := Unwrap(left)
		_, rightVal := Unwrap(right)

		assert.Equal(t, -3, leftVal)
		assert.Equal(t, "not positive", rightVal)
	})

	t.Run("Right value at boundary (zero)", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, "not positive")
		input := Right[string](0)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Right (zero fails predicate)")
		assert.True(t, IsLeft(right), "right should be Left")

		leftVal, _ := Unwrap(left)
		_, rightVal := Unwrap(right)

		assert.Equal(t, 0, leftVal)
		assert.Equal(t, "not positive", rightVal)
	})

	t.Run("Left value passes through unchanged", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, "not positive")
		input := Left[int]("original error")

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be Left")
		assert.True(t, IsLeft(right), "right should be Left")

		_, leftVal := Unwrap(left)
		_, rightVal := Unwrap(right)

		assert.Equal(t, "original error", leftVal)
		assert.Equal(t, "original error", rightVal)
	})

	t.Run("String predicate - even length strings", func(t *testing.T) {
		// Arrange
		isEvenLength := func(s string) bool { return len(s)%2 == 0 }
		partition := Partition(isEvenLength, 0)

		// Act & Assert - passes predicate
		result1 := partition(Right[int]("test"))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, "test", rightVal1)

		// Act & Assert - fails predicate
		result2 := partition(Right[int]("hello"))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, "hello", leftVal2)
	})

	t.Run("Boolean predicate - identity function", func(t *testing.T) {
		// Arrange
		identity := func(b bool) bool { return b }
		partition := Partition(identity, "false value")

		// Act & Assert - true passes
		result1 := partition(Right[string](true))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, true, rightVal1)

		// Act & Assert - false fails
		result2 := partition(Right[string](false))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, false, leftVal2)
	})

	t.Run("Complex type predicate - struct field check", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		isAdult := func(p Person) bool { return p.Age >= 18 }
		partition := Partition(isAdult, "minor")

		// Act & Assert - adult passes
		adult := Person{Name: "Alice", Age: 25}
		result1 := partition(Right[string](adult))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, adult, rightVal1)

		// Act & Assert - minor fails
		minor := Person{Name: "Bob", Age: 15}
		result2 := partition(Right[string](minor))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, minor, leftVal2)
	})

	t.Run("Predicate always true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		partition := Partition(alwaysTrue, "never used")

		// Act
		result := partition(Right[string](42))
		left, right := P.Unpack(result)

		// Assert - all values pass to right
		assert.True(t, IsLeft(left))
		assert.True(t, IsRight(right))
		rightVal, _ := Unwrap(right)
		assert.Equal(t, 42, rightVal)
	})

	t.Run("Predicate always false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		partition := Partition(alwaysFalse, "always used")

		// Act
		result := partition(Right[string](42))
		left, right := P.Unpack(result)

		// Assert - all values fail to left
		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		leftVal, _ := Unwrap(left)
		assert.Equal(t, 42, leftVal)
	})

	t.Run("Multiple calls with same partition function", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		partition := Partition(isEven, -1)

		// Act & Assert - multiple values
		values := []int{2, 3, 4, 5, 6}
		for _, val := range values {
			result := partition(Right[int](val))
			left, right := P.Unpack(result)

			if val%2 == 0 {
				assert.True(t, IsRight(right), "even value %d should be in right", val)
				rightVal, _ := Unwrap(right)
				assert.Equal(t, val, rightVal)
			} else {
				assert.True(t, IsRight(left), "odd value %d should be in left", val)
				leftVal, _ := Unwrap(left)
				assert.Equal(t, val, leftVal)
			}
		}
	})

	t.Run("Empty value types", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		partition := Partition(isEmpty, 0)

		// Act & Assert - empty string passes
		result1 := partition(Right[int](""))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, "", rightVal1)

		// Act & Assert - non-empty string fails
		result2 := partition(Right[int]("hello"))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, "hello", leftVal2)
	})
}

func TestPartitionWithDifferentErrorTypes(t *testing.T) {
	t.Run("String error type", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, "error")

		result := partition(Right[string](5))
		left, right := P.Unpack(result)

		assert.True(t, IsLeft(left))
		assert.True(t, IsRight(right))
	})

	t.Run("Int error type", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, -999)

		result := partition(Right[int](-3))
		left, right := P.Unpack(result)

		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		_, rightVal := Unwrap(right)
		assert.Equal(t, -999, rightVal)
	})

	t.Run("Struct error type", func(t *testing.T) {
		type CustomError struct {
			Code    int
			Message string
		}

		isPositive := N.MoreThan(0)
		defaultError := CustomError{Code: 400, Message: "Invalid"}
		partition := Partition(isPositive, defaultError)

		result := partition(Right[CustomError](-3))
		left, right := P.Unpack(result)

		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		_, rightVal := Unwrap(right)
		assert.Equal(t, defaultError, rightVal)
	})
}

// Benchmark tests
func BenchmarkPartition(b *testing.B) {
	isPositive := N.MoreThan(0)
	partition := Partition(isPositive, "not positive")
	input := Right[string](42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partition(input)
	}
}

func BenchmarkPartitionLeft(b *testing.B) {
	isPositive := N.MoreThan(0)
	partition := Partition(isPositive, "not positive")
	input := Left[int]("error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partition(input)
	}
}

func BenchmarkPartitionPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	partition := Partition(isPositive, "not positive")
	input := Right[string](-42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partition(input)
	}
}

func TestFilter(t *testing.T) {
	t.Run("Right value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "not positive")
		input := Right[string](5)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsRight(result), "result should be Right")
		val, _ := Unwrap(result)
		assert.Equal(t, 5, val)
	})

	t.Run("Right value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "not positive")
		input := Right[string](-3)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be Left")
		_, err := Unwrap(result)
		assert.Equal(t, "not positive", err)
	})

	t.Run("Right value at boundary (zero)", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "not positive")
		input := Right[string](0)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "zero should fail predicate")
		_, err := Unwrap(result)
		assert.Equal(t, "not positive", err)
	})

	t.Run("Left value passes through unchanged", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "not positive")
		input := Left[int]("original error")

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be Left")
		_, err := Unwrap(result)
		assert.Equal(t, "original error", err)
	})

	t.Run("String predicate - even length strings", func(t *testing.T) {
		// Arrange
		isEvenLength := func(s string) bool { return len(s)%2 == 0 }
		filter := Filter(isEvenLength, 0)

		// Act & Assert - passes predicate
		result1 := filter(Right[int]("test"))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, "test", val1)

		// Act & Assert - fails predicate
		result2 := filter(Right[int]("hello"))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, 0, err2)
	})

	t.Run("Boolean predicate - identity function", func(t *testing.T) {
		// Arrange
		identity := func(b bool) bool { return b }
		filter := Filter(identity, "false value")

		// Act & Assert - true passes
		result1 := filter(Right[string](true))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, true, val1)

		// Act & Assert - false fails
		result2 := filter(Right[string](false))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "false value", err2)
	})

	t.Run("Complex type predicate - struct field check", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		isAdult := func(p Person) bool { return p.Age >= 18 }
		filter := Filter(isAdult, "minor")

		// Act & Assert - adult passes
		adult := Person{Name: "Alice", Age: 25}
		result1 := filter(Right[string](adult))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, adult, val1)

		// Act & Assert - minor fails
		minor := Person{Name: "Bob", Age: 15}
		result2 := filter(Right[string](minor))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "minor", err2)
	})

	t.Run("Predicate always true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		filter := Filter(alwaysTrue, "never used")

		// Act
		result := filter(Right[string](42))

		// Assert - all values pass through
		assert.True(t, IsRight(result))
		val, _ := Unwrap(result)
		assert.Equal(t, 42, val)
	})

	t.Run("Predicate always false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		filter := Filter(alwaysFalse, "always filtered")

		// Act
		result := filter(Right[string](42))

		// Assert - all values filtered out
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, "always filtered", err)
	})

	t.Run("Multiple calls with same filter function", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		filter := Filter(isEven, -1)

		// Act & Assert - multiple values
		values := []int{2, 3, 4, 5, 6}
		for _, val := range values {
			result := filter(Right[int](val))

			if val%2 == 0 {
				assert.True(t, IsRight(result), "even value %d should pass", val)
				resultVal, _ := Unwrap(result)
				assert.Equal(t, val, resultVal)
			} else {
				assert.True(t, IsLeft(result), "odd value %d should be filtered", val)
				_, err := Unwrap(result)
				assert.Equal(t, -1, err)
			}
		}
	})

	t.Run("Empty value types", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		filter := Filter(isEmpty, 0)

		// Act & Assert - empty string passes
		result1 := filter(Right[int](""))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, "", val1)

		// Act & Assert - non-empty string fails
		result2 := filter(Right[int]("hello"))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, 0, err2)
	})

	t.Run("Chaining multiple filters", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isEven := func(n int) bool { return n%2 == 0 }
		filterPositive := Filter(isPositive, "not positive")
		filterEven := Filter(isEven, "not even")

		// Act & Assert - passes both filters
		result1 := filterEven(filterPositive(Right[string](4)))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 4, val1)

		// Act & Assert - passes first, fails second
		result2 := filterEven(filterPositive(Right[string](3)))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "not even", err2)

		// Act & Assert - fails first filter
		result3 := filterEven(filterPositive(Right[string](-2)))
		assert.True(t, IsLeft(result3))
		_, err3 := Unwrap(result3)
		assert.Equal(t, "not positive", err3)

		// Act & Assert - Left passes through both
		result4 := filterEven(filterPositive(Left[int]("original")))
		assert.True(t, IsLeft(result4))
		_, err4 := Unwrap(result4)
		assert.Equal(t, "original", err4)
	})

	t.Run("Filter preserves Left error type", func(t *testing.T) {
		// Arrange
		type CustomError struct {
			Code    int
			Message string
		}
		isPositive := N.MoreThan(0)
		defaultError := CustomError{Code: 400, Message: "Invalid"}
		filter := Filter(isPositive, defaultError)

		// Act - Left with different error
		originalError := CustomError{Code: 500, Message: "Server Error"}
		result := filter(Left[int](originalError))

		// Assert - original error preserved
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, originalError, err)
		assert.NotEqual(t, defaultError, err)
	})

	t.Run("Filter with nil predicate behavior", func(t *testing.T) {
		// Arrange
		isNonNil := func(p *int) bool { return p != nil }
		filter := Filter(isNonNil, "nil pointer")

		// Act & Assert - non-nil passes
		val := 42
		result1 := filter(Right[string](&val))
		assert.True(t, IsRight(result1))

		// Act & Assert - nil fails
		result2 := filter(Right[string]((*int)(nil)))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "nil pointer", err2)
	})
}

func TestFilterWithDifferentErrorTypes(t *testing.T) {
	t.Run("String error type", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "error")

		result := filter(Right[string](5))
		assert.True(t, IsRight(result))
	})

	t.Run("Int error type", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, -999)

		result := filter(Right[int](-3))
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, -999, err)
	})

	t.Run("Struct error type", func(t *testing.T) {
		type CustomError struct {
			Code    int
			Message string
		}

		isPositive := N.MoreThan(0)
		defaultError := CustomError{Code: 400, Message: "Invalid"}
		filter := Filter(isPositive, defaultError)

		result := filter(Right[CustomError](-3))
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, defaultError, err)
	})

	t.Run("Error interface type", func(t *testing.T) {
		type ValidationError struct {
			Field string
			Issue string
		}

		isPositive := N.MoreThan(0)
		defaultError := ValidationError{Field: "value", Issue: "must be positive"}
		filter := Filter(isPositive, defaultError)

		result := filter(Right[ValidationError](0))
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, defaultError, err)
	})
}

func TestFilterEdgeCases(t *testing.T) {
	t.Run("Filter with zero value as empty", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, 0)

		// Act
		result := filter(Right[int](-5))

		// Assert
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, 0, err)
	})

	t.Run("Filter with empty string as empty", func(t *testing.T) {
		// Arrange
		isLongEnough := func(s string) bool { return len(s) >= 5 }
		filter := Filter(isLongEnough, "")

		// Act
		result := filter(Right[string]("hi"))

		// Assert
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, "", err)
	})

	t.Run("Filter reuses same Left instance", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, "error")

		// Act - multiple calls that fail predicate
		result1 := filter(Right[string](-1))
		result2 := filter(Right[string](-2))

		// Assert - both should be Left with same error
		assert.True(t, IsLeft(result1))
		assert.True(t, IsLeft(result2))
		_, err1 := Unwrap(result1)
		_, err2 := Unwrap(result2)
		assert.Equal(t, err1, err2)
	})
}

// Benchmark tests for Filter
func BenchmarkFilter(b *testing.B) {
	isPositive := N.MoreThan(0)
	filter := Filter(isPositive, "not positive")
	input := Right[string](42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(input)
	}
}

func BenchmarkFilterLeft(b *testing.B) {
	isPositive := N.MoreThan(0)
	filter := Filter(isPositive, "not positive")
	input := Left[int]("error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(input)
	}
}

func BenchmarkFilterPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	filter := Filter(isPositive, "not positive")
	input := Right[string](-42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(input)
	}
}

func BenchmarkFilterChained(b *testing.B) {
	isPositive := N.MoreThan(0)
	isEven := func(n int) bool { return n%2 == 0 }
	filterPositive := Filter(isPositive, "not positive")
	filterEven := Filter(isEven, "not even")
	input := Right[string](42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterEven(filterPositive(input))
	}
}

func TestFilterMap(t *testing.T) {
	t.Run("Right value with Some result", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) option.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return option.Some(n)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(parseInt, "invalid number")
		input := Right[string]("42")

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsRight(result), "result should be Right")
		val, _ := Unwrap(result)
		assert.Equal(t, 42, val)
	})

	t.Run("Right value with None result", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) option.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return option.Some(n)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(parseInt, "invalid number")
		input := Right[string]("abc")

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be Left")
		_, err := Unwrap(result)
		assert.Equal(t, "invalid number", err)
	})

	t.Run("Left value passes through", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) option.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return option.Some(n)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(parseInt, "invalid number")
		input := Left[string]("original error")

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be Left")
		_, err := Unwrap(result)
		assert.Equal(t, "original error", err)
	})

	t.Run("Extract optional field from struct", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name  string
			Email option.Option[string]
		}
		extractEmail := func(p Person) option.Option[string] { return p.Email }
		filterMap := FilterMap(extractEmail, "no email")

		// Act & Assert - has email
		result1 := filterMap(Right[string](Person{Name: "Alice", Email: option.Some("alice@example.com")}))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, "alice@example.com", val1)

		// Act & Assert - no email
		result2 := filterMap(Right[string](Person{Name: "Bob", Email: option.None[string]()}))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "no email", err2)
	})

	t.Run("Transform and filter numbers", func(t *testing.T) {
		// Arrange
		sqrtIfPositive := func(n int) option.Option[float64] {
			if n >= 0 {
				return option.Some(math.Sqrt(float64(n)))
			}
			return option.None[float64]()
		}
		filterMap := FilterMap(sqrtIfPositive, "negative number")

		// Act & Assert - positive number
		result1 := filterMap(Right[string](16))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 4.0, val1)

		// Act & Assert - negative number
		result2 := filterMap(Right[string](-4))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "negative number", err2)

		// Act & Assert - zero
		result3 := filterMap(Right[string](0))
		assert.True(t, IsRight(result3))
		val3, _ := Unwrap(result3)
		assert.Equal(t, 0.0, val3)
	})

	t.Run("Lookup in map", func(t *testing.T) {
		// Arrange
		data := map[string]int{"a": 1, "b": 2, "c": 3}
		lookup := func(key string) option.Option[int] {
			if val, ok := data[key]; ok {
				return option.Some(val)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(lookup, "key not found")

		// Act & Assert - existing key
		result1 := filterMap(Right[string]("b"))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 2, val1)

		// Act & Assert - non-existing key
		result2 := filterMap(Right[string]("z"))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "key not found", err2)
	})

	t.Run("Chain multiple FilterMap operations", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) option.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return option.Some(n)
			}
			return option.None[int]()
		}
		doubleIfEven := func(n int) option.Option[int] {
			if n%2 == 0 {
				return option.Some(n * 2)
			}
			return option.None[int]()
		}
		filterMap1 := FilterMap(parseInt, "invalid number")
		filterMap2 := FilterMap(doubleIfEven, "not even")

		// Act & Assert - valid even number
		result1 := filterMap2(filterMap1(Right[string]("4")))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 8, val1)

		// Act & Assert - valid odd number
		result2 := filterMap2(filterMap1(Right[string]("3")))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "not even", err2)

		// Act & Assert - invalid number
		result3 := filterMap2(filterMap1(Right[string]("abc")))
		assert.True(t, IsLeft(result3))
		_, err3 := Unwrap(result3)
		assert.Equal(t, "invalid number", err3)
	})

	t.Run("Always returns Some", func(t *testing.T) {
		// Arrange
		alwaysSome := func(n int) option.Option[string] {
			return option.Some(strconv.Itoa(n))
		}
		filterMap := FilterMap(alwaysSome, "never used")

		// Act
		result := filterMap(Right[string](42))

		// Assert
		assert.True(t, IsRight(result))
		val, _ := Unwrap(result)
		assert.Equal(t, "42", val)
	})

	t.Run("Always returns None", func(t *testing.T) {
		// Arrange
		alwaysNone := func(int) option.Option[string] {
			return option.None[string]()
		}
		filterMap := FilterMap(alwaysNone, "always filtered")

		// Act
		result := filterMap(Right[string](42))

		// Assert
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, "always filtered", err)
	})

	t.Run("Type transformation with validation", func(t *testing.T) {
		// Arrange
		type Input struct{ Value string }
		type Output struct{ Number int }

		parseInput := func(in Input) option.Option[Output] {
			if n, err := strconv.Atoi(in.Value); err == nil && n > 0 {
				return option.Some(Output{Number: n})
			}
			return option.None[Output]()
		}
		filterMap := FilterMap(parseInput, "invalid input")

		// Act & Assert - valid input
		result1 := filterMap(Right[string](Input{Value: "42"}))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, Output{Number: 42}, val1)

		// Act & Assert - invalid input
		result2 := filterMap(Right[string](Input{Value: "abc"}))
		assert.True(t, IsLeft(result2))

		// Act & Assert - zero (invalid per predicate)
		result3 := filterMap(Right[string](Input{Value: "0"}))
		assert.True(t, IsLeft(result3))
	})
}

func TestFilterMapWithDifferentErrorTypes(t *testing.T) {
	t.Run("String error type", func(t *testing.T) {
		parseInt := func(s string) option.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return option.Some(n)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(parseInt, "error")

		result := filterMap(Right[string]("42"))
		assert.True(t, IsRight(result))
	})

	t.Run("Int error type", func(t *testing.T) {
		doubleIfPositive := func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n * 2)
			}
			return option.None[int]()
		}
		filterMap := FilterMap(doubleIfPositive, -999)

		result := filterMap(Right[int](-5))
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, -999, err)
	})

	t.Run("Struct error type", func(t *testing.T) {
		type CustomError struct {
			Code    int
			Message string
		}

		validate := func(n int) option.Option[int] {
			if n > 0 {
				return option.Some(n)
			}
			return option.None[int]()
		}
		defaultError := CustomError{Code: 400, Message: "Invalid"}
		filterMap := FilterMap(validate, defaultError)

		result := filterMap(Right[CustomError](-3))
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, defaultError, err)
	})
}

func TestPartitionMap(t *testing.T) {
	t.Run("Right value that maps to Left", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) Either[string, int] {
			if n < 0 {
				return Left[int]("negative: " + strconv.Itoa(n))
			}
			return Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, "not classified")
		input := Right[string](-3)

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Right (contains error from f)")
		assert.True(t, IsLeft(right), "right should be Left")

		leftVal, _ := Unwrap(left)
		_, rightErr := Unwrap(right)

		assert.Equal(t, "negative: -3", leftVal)
		assert.Equal(t, "not classified", rightErr)
	})

	t.Run("Right value that maps to Right", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) Either[string, int] {
			if n < 0 {
				return Left[int]("negative: " + strconv.Itoa(n))
			}
			return Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, "not classified")
		input := Right[string](5)

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be Left")
		assert.True(t, IsRight(right), "right should be Right (contains value from f)")

		_, leftErr := Unwrap(left)
		rightVal, _ := Unwrap(right)

		assert.Equal(t, "not classified", leftErr)
		assert.Equal(t, 25, rightVal)
	})

	t.Run("Left value passes through to both sides", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) Either[string, int] {
			if n < 0 {
				return Left[int]("negative")
			}
			return Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, "not classified")
		input := Left[int]("original error")

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be Left")
		assert.True(t, IsLeft(right), "right should be Left")

		_, leftErr := Unwrap(left)
		_, rightErr := Unwrap(right)

		assert.Equal(t, "original error", leftErr)
		assert.Equal(t, "original error", rightErr)
	})

	t.Run("Validate and transform user input", func(t *testing.T) {
		// Arrange
		type ValidationError struct {
			Field   string
			Message string
		}
		type User struct {
			Name string
			Age  int
		}

		validateUser := func(input map[string]string) Either[ValidationError, User] {
			name, hasName := input["name"]
			ageStr, hasAge := input["age"]
			if !hasName {
				return Left[User](ValidationError{"name", "missing"})
			}
			if !hasAge {
				return Left[User](ValidationError{"age", "missing"})
			}
			age, err := strconv.Atoi(ageStr)
			if err != nil {
				return Left[User](ValidationError{"age", "invalid"})
			}
			return Right[ValidationError](User{name, age})
		}
		partitionMap := PartitionMap(validateUser, ValidationError{"", "not processed"})

		// Act & Assert - valid input
		validInput := map[string]string{"name": "Alice", "age": "30"}
		result1 := partitionMap(Right[ValidationError](validInput))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, User{"Alice", 30}, rightVal1)

		// Act & Assert - invalid input (missing age)
		invalidInput := map[string]string{"name": "Bob"}
		result2 := partitionMap(Right[ValidationError](invalidInput))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, ValidationError{"age", "missing"}, leftVal2)
	})

	t.Run("Classify strings by length", func(t *testing.T) {
		// Arrange
		classifyString := func(s string) Either[string, int] {
			if len(s) < 5 {
				return Left[int]("too short: " + s)
			}
			return Right[string](len(s))
		}
		partitionMap := PartitionMap(classifyString, "not classified")

		// Act & Assert - short string
		result1 := partitionMap(Right[string]("hi"))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsRight(left1))
		assert.True(t, IsLeft(right1))
		leftVal1, _ := Unwrap(left1)
		assert.Equal(t, "too short: hi", leftVal1)

		// Act & Assert - long string
		result2 := partitionMap(Right[string]("hello world"))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsLeft(left2))
		assert.True(t, IsRight(right2))
		rightVal2, _ := Unwrap(right2)
		assert.Equal(t, 11, rightVal2)
	})

	t.Run("Always maps to Left", func(t *testing.T) {
		// Arrange
		alwaysLeft := func(n int) Either[string, int] {
			return Left[int]("error: " + strconv.Itoa(n))
		}
		partitionMap := PartitionMap(alwaysLeft, "not classified")

		// Act
		result := partitionMap(Right[string](42))
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		leftVal, _ := Unwrap(left)
		assert.Equal(t, "error: 42", leftVal)
	})

	t.Run("Always maps to Right", func(t *testing.T) {
		// Arrange
		alwaysRight := func(n int) Either[string, int] {
			return Right[string](n * 2)
		}
		partitionMap := PartitionMap(alwaysRight, "not classified")

		// Act
		result := partitionMap(Right[string](42))
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left))
		assert.True(t, IsRight(right))
		rightVal, _ := Unwrap(right)
		assert.Equal(t, 84, rightVal)
	})

	t.Run("Complex type transformation", func(t *testing.T) {
		// Arrange
		type Request struct{ ID int }
		type Success struct{ Data string }
		type Failure struct{ Code int }

		processRequest := func(req Request) Either[Failure, Success] {
			if req.ID <= 0 {
				return Left[Success](Failure{Code: 400})
			}
			return Right[Failure](Success{Data: "processed-" + strconv.Itoa(req.ID)})
		}
		partitionMap := PartitionMap(processRequest, Failure{Code: 0})

		// Act & Assert - valid request
		result1 := partitionMap(Right[Failure](Request{ID: 123}))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, Success{Data: "processed-123"}, rightVal1)

		// Act & Assert - invalid request
		result2 := partitionMap(Right[Failure](Request{ID: -1}))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, Failure{Code: 400}, leftVal2)
	})

	t.Run("Multiple calls with same partition function", func(t *testing.T) {
		// Arrange
		classify := func(n int) Either[string, string] {
			if n%2 == 0 {
				return Right[string]("even")
			}
			return Left[string]("odd")
		}
		partitionMap := PartitionMap(classify, "unclassified")

		// Act & Assert - multiple values
		values := []int{2, 3, 4, 5, 6}
		for _, val := range values {
			result := partitionMap(Right[string](val))
			left, right := P.Unpack(result)

			if val%2 == 0 {
				assert.True(t, IsRight(right), "even value %d should be in right", val)
				rightVal, _ := Unwrap(right)
				assert.Equal(t, "even", rightVal)
			} else {
				assert.True(t, IsRight(left), "odd value %d should be in left", val)
				leftVal, _ := Unwrap(left)
				assert.Equal(t, "odd", leftVal)
			}
		}
	})
}

func TestPartitionMapWithDifferentErrorTypes(t *testing.T) {
	t.Run("String error types", func(t *testing.T) {
		classify := func(n int) Either[string, int] {
			if n < 0 {
				return Left[int]("negative")
			}
			return Right[string](n)
		}
		partitionMap := PartitionMap(classify, "not classified")

		result := partitionMap(Right[string](5))
		left, right := P.Unpack(result)
		assert.True(t, IsLeft(left))
		assert.True(t, IsRight(right))
	})

	t.Run("Int error types", func(t *testing.T) {
		classify := func(n int) Either[int, int] {
			if n < 0 {
				return Left[int](-1)
			}
			return Right[int](n)
		}
		partitionMap := PartitionMap(classify, -999)

		result := partitionMap(Right[int](-5))
		left, right := P.Unpack(result)
		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		leftVal, _ := Unwrap(left)
		assert.Equal(t, -1, leftVal)
	})

	t.Run("Struct error types", func(t *testing.T) {
		type ErrorA struct{ Code int }
		type ErrorB struct{ Message string }

		classify := func(n int) Either[ErrorA, int] {
			if n < 0 {
				return Left[int](ErrorA{Code: 400})
			}
			return Right[ErrorA](n)
		}
		defaultError := ErrorB{Message: "not classified"}
		partitionMap := PartitionMap(classify, defaultError)

		result := partitionMap(Right[ErrorB](-3))
		left, right := P.Unpack(result)
		assert.True(t, IsRight(left))
		assert.True(t, IsLeft(right))
		leftVal, _ := Unwrap(left)
		assert.Equal(t, ErrorA{Code: 400}, leftVal)
	})
}

// Benchmark tests for FilterMap
func BenchmarkFilterMap(b *testing.B) {
	parseInt := func(s string) option.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return option.Some(n)
		}
		return option.None[int]()
	}
	filterMap := FilterMap(parseInt, "invalid")
	input := Right[string]("42")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterMap(input)
	}
}

func BenchmarkFilterMapLeft(b *testing.B) {
	parseInt := func(s string) option.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return option.Some(n)
		}
		return option.None[int]()
	}
	filterMap := FilterMap(parseInt, "invalid")
	input := Left[string]("error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterMap(input)
	}
}

func BenchmarkFilterMapNone(b *testing.B) {
	parseInt := func(s string) option.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return option.Some(n)
		}
		return option.None[int]()
	}
	filterMap := FilterMap(parseInt, "invalid")
	input := Right[string]("abc")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterMap(input)
	}
}

// Benchmark tests for PartitionMap
func BenchmarkPartitionMap(b *testing.B) {
	classify := func(n int) Either[string, int] {
		if n < 0 {
			return Left[int]("negative")
		}
		return Right[string](n * n)
	}
	partitionMap := PartitionMap(classify, "not classified")
	input := Right[string](42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partitionMap(input)
	}
}

func BenchmarkPartitionMapLeft(b *testing.B) {
	classify := func(n int) Either[string, int] {
		if n < 0 {
			return Left[int]("negative")
		}
		return Right[string](n * n)
	}
	partitionMap := PartitionMap(classify, "not classified")
	input := Left[int]("error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partitionMap(input)
	}
}

func BenchmarkPartitionMapToLeft(b *testing.B) {
	classify := func(n int) Either[string, int] {
		if n < 0 {
			return Left[int]("negative")
		}
		return Right[string](n * n)
	}
	partitionMap := PartitionMap(classify, "not classified")
	input := Right[string](-42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partitionMap(input)
	}
}
