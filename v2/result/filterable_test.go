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

package result

import (
	"errors"
	"math"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestPartition(t *testing.T) {
	t.Run("Ok value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, errors.New("not positive"))
		input := Of(5)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be error")
		assert.True(t, IsRight(right), "right should be Ok")

		rightVal, _ := Unwrap(right)
		assert.Equal(t, 5, rightVal)
	})

	t.Run("Ok value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, errors.New("not positive"))
		input := Of(-3)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Ok (failed predicate)")
		assert.True(t, IsLeft(right), "right should be error")

		leftVal, _ := Unwrap(left)
		assert.Equal(t, -3, leftVal)
	})

	t.Run("Ok value at boundary (zero)", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, errors.New("not positive"))
		input := Of(0)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Ok (zero fails predicate)")
		assert.True(t, IsLeft(right), "right should be error")

		leftVal, _ := Unwrap(left)
		assert.Equal(t, 0, leftVal)
	})

	t.Run("Error passes through unchanged", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		partition := Partition(isPositive, errors.New("not positive"))
		originalError := errors.New("original error")
		input := Left[int](originalError)

		// Act
		result := partition(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be error")
		assert.True(t, IsLeft(right), "right should be error")

		_, leftErr := Unwrap(left)
		_, rightErr := Unwrap(right)

		assert.Equal(t, originalError, leftErr)
		assert.Equal(t, originalError, rightErr)
	})

	t.Run("String predicate - even length strings", func(t *testing.T) {
		// Arrange
		isEvenLength := func(s string) bool { return len(s)%2 == 0 }
		partition := Partition(isEvenLength, errors.New("odd length"))

		// Act & Assert - passes predicate
		result1 := partition(Of("test"))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, "test", rightVal1)

		// Act & Assert - fails predicate
		result2 := partition(Of("hello"))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, "hello", leftVal2)
	})

	t.Run("Complex type predicate - struct field check", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name string
			Age  int
		}
		isAdult := func(p Person) bool { return p.Age >= 18 }
		partition := Partition(isAdult, errors.New("minor"))

		// Act & Assert - adult passes
		adult := Person{Name: "Alice", Age: 25}
		result1 := partition(Of(adult))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, adult, rightVal1)

		// Act & Assert - minor fails
		minor := Person{Name: "Bob", Age: 15}
		result2 := partition(Of(minor))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, minor, leftVal2)
	})
}

func TestFilter(t *testing.T) {
	t.Run("Ok value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, errors.New("not positive"))
		input := Of(5)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsRight(result), "result should be Ok")
		val, _ := Unwrap(result)
		assert.Equal(t, 5, val)
	})

	t.Run("Ok value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, errors.New("not positive"))
		input := Of(-3)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be error")
		_, err := Unwrap(result)
		assert.Equal(t, "not positive", err.Error())
	})

	t.Run("Ok value at boundary (zero)", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, errors.New("not positive"))
		input := Of(0)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "zero should fail predicate")
		_, err := Unwrap(result)
		assert.Equal(t, "not positive", err.Error())
	})

	t.Run("Error passes through unchanged", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, errors.New("not positive"))
		originalError := errors.New("original error")
		input := Left[int](originalError)

		// Act
		result := filter(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be error")
		_, err := Unwrap(result)
		assert.Equal(t, originalError, err)
	})

	t.Run("Chaining multiple filters", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isEven := func(n int) bool { return n%2 == 0 }
		filterPositive := Filter(isPositive, errors.New("not positive"))
		filterEven := Filter(isEven, errors.New("not even"))

		// Act & Assert - passes both filters
		result1 := filterEven(filterPositive(Of(4)))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 4, val1)

		// Act & Assert - passes first, fails second
		result2 := filterEven(filterPositive(Of(3)))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "not even", err2.Error())

		// Act & Assert - fails first filter
		result3 := filterEven(filterPositive(Of(-2)))
		assert.True(t, IsLeft(result3))
		_, err3 := Unwrap(result3)
		assert.Equal(t, "not positive", err3.Error())

		// Act & Assert - error passes through both
		originalErr := errors.New("original")
		result4 := filterEven(filterPositive(Left[int](originalErr)))
		assert.True(t, IsLeft(result4))
		_, err4 := Unwrap(result4)
		assert.Equal(t, originalErr, err4)
	})

	t.Run("Filter preserves error", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		filter := Filter(isPositive, errors.New("default error"))

		// Act - error with different message
		originalError := errors.New("server error")
		result := filter(Left[int](originalError))

		// Assert - original error preserved
		assert.True(t, IsLeft(result))
		_, err := Unwrap(result)
		assert.Equal(t, originalError, err)
	})
}

func TestFilterMap(t *testing.T) {
	t.Run("Ok value with Some result", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) O.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return O.Some(n)
			}
			return O.None[int]()
		}
		filterMap := FilterMap(parseInt, errors.New("invalid number"))
		input := Of("42")

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsRight(result), "result should be Ok")
		val, _ := Unwrap(result)
		assert.Equal(t, 42, val)
	})

	t.Run("Ok value with None result", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) O.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return O.Some(n)
			}
			return O.None[int]()
		}
		filterMap := FilterMap(parseInt, errors.New("invalid number"))
		input := Of("abc")

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be error")
		_, err := Unwrap(result)
		assert.Equal(t, "invalid number", err.Error())
	})

	t.Run("Error passes through", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) O.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return O.Some(n)
			}
			return O.None[int]()
		}
		filterMap := FilterMap(parseInt, errors.New("invalid number"))
		originalError := errors.New("original error")
		input := Left[string](originalError)

		// Act
		result := filterMap(input)

		// Assert
		assert.True(t, IsLeft(result), "result should be error")
		_, err := Unwrap(result)
		assert.Equal(t, originalError, err)
	})

	t.Run("Extract optional field from struct", func(t *testing.T) {
		// Arrange
		type Person struct {
			Name  string
			Email O.Option[string]
		}
		extractEmail := func(p Person) O.Option[string] { return p.Email }
		filterMap := FilterMap(extractEmail, errors.New("no email"))

		// Act & Assert - has email
		result1 := filterMap(Of(Person{Name: "Alice", Email: O.Some("alice@example.com")}))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, "alice@example.com", val1)

		// Act & Assert - no email
		result2 := filterMap(Of(Person{Name: "Bob", Email: O.None[string]()}))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "no email", err2.Error())
	})

	t.Run("Transform and filter numbers", func(t *testing.T) {
		// Arrange
		sqrtIfPositive := func(n int) O.Option[float64] {
			if n >= 0 {
				return O.Some(math.Sqrt(float64(n)))
			}
			return O.None[float64]()
		}
		filterMap := FilterMap(sqrtIfPositive, errors.New("negative number"))

		// Act & Assert - positive number
		result1 := filterMap(Of(16))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 4.0, val1)

		// Act & Assert - negative number
		result2 := filterMap(Of(-4))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "negative number", err2.Error())

		// Act & Assert - zero
		result3 := filterMap(Of(0))
		assert.True(t, IsRight(result3))
		val3, _ := Unwrap(result3)
		assert.Equal(t, 0.0, val3)
	})

	t.Run("Chain multiple FilterMap operations", func(t *testing.T) {
		// Arrange
		parseInt := func(s string) O.Option[int] {
			if n, err := strconv.Atoi(s); err == nil {
				return O.Some(n)
			}
			return O.None[int]()
		}
		doubleIfEven := func(n int) O.Option[int] {
			if n%2 == 0 {
				return O.Some(n * 2)
			}
			return O.None[int]()
		}
		filterMap1 := FilterMap(parseInt, errors.New("invalid number"))
		filterMap2 := FilterMap(doubleIfEven, errors.New("not even"))

		// Act & Assert - valid even number
		result1 := filterMap2(filterMap1(Of("4")))
		assert.True(t, IsRight(result1))
		val1, _ := Unwrap(result1)
		assert.Equal(t, 8, val1)

		// Act & Assert - valid odd number
		result2 := filterMap2(filterMap1(Of("3")))
		assert.True(t, IsLeft(result2))
		_, err2 := Unwrap(result2)
		assert.Equal(t, "not even", err2.Error())

		// Act & Assert - invalid number
		result3 := filterMap2(filterMap1(Of("abc")))
		assert.True(t, IsLeft(result3))
		_, err3 := Unwrap(result3)
		assert.Equal(t, "invalid number", err3.Error())
	})
}

func TestPartitionMap(t *testing.T) {
	t.Run("Ok value that maps to Left", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) E.Either[string, int] {
			if n < 0 {
				return E.Left[int]("negative: " + strconv.Itoa(n))
			}
			return E.Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, errors.New("not classified"))
		input := Of(-3)

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsRight(left), "left should be Ok (contains error from f)")
		assert.True(t, IsLeft(right), "right should be error")

		leftVal, _ := Unwrap(left)
		assert.Equal(t, "negative: -3", leftVal)
	})

	t.Run("Ok value that maps to Right", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) E.Either[string, int] {
			if n < 0 {
				return E.Left[int]("negative: " + strconv.Itoa(n))
			}
			return E.Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, errors.New("not classified"))
		input := Of(5)

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be error")
		assert.True(t, IsRight(right), "right should be Ok (contains value from f)")

		rightVal, _ := Unwrap(right)
		assert.Equal(t, 25, rightVal)
	})

	t.Run("Error passes through to both sides", func(t *testing.T) {
		// Arrange
		classifyNumber := func(n int) E.Either[string, int] {
			if n < 0 {
				return E.Left[int]("negative")
			}
			return E.Right[string](n * n)
		}
		partitionMap := PartitionMap(classifyNumber, errors.New("not classified"))
		originalError := errors.New("original error")
		input := Left[int](originalError)

		// Act
		result := partitionMap(input)
		left, right := P.Unpack(result)

		// Assert
		assert.True(t, IsLeft(left), "left should be error")
		assert.True(t, IsLeft(right), "right should be error")

		_, leftErr := Unwrap(left)
		_, rightErr := Unwrap(right)

		assert.Equal(t, originalError, leftErr)
		assert.Equal(t, originalError, rightErr)
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

		validateUser := func(input map[string]string) E.Either[ValidationError, User] {
			name, hasName := input["name"]
			ageStr, hasAge := input["age"]
			if !hasName {
				return E.Left[User](ValidationError{"name", "missing"})
			}
			if !hasAge {
				return E.Left[User](ValidationError{"age", "missing"})
			}
			age, err := strconv.Atoi(ageStr)
			if err != nil {
				return E.Left[User](ValidationError{"age", "invalid"})
			}
			return E.Right[ValidationError](User{name, age})
		}
		partitionMap := PartitionMap(validateUser, errors.New("not processed"))

		// Act & Assert - valid input
		validInput := map[string]string{"name": "Alice", "age": "30"}
		result1 := partitionMap(Of(validInput))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsLeft(left1))
		assert.True(t, IsRight(right1))
		rightVal1, _ := Unwrap(right1)
		assert.Equal(t, User{"Alice", 30}, rightVal1)

		// Act & Assert - invalid input (missing age)
		invalidInput := map[string]string{"name": "Bob"}
		result2 := partitionMap(Of(invalidInput))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsRight(left2))
		assert.True(t, IsLeft(right2))
		leftVal2, _ := Unwrap(left2)
		assert.Equal(t, ValidationError{"age", "missing"}, leftVal2)
	})

	t.Run("Classify strings by length", func(t *testing.T) {
		// Arrange
		classifyString := func(s string) E.Either[string, int] {
			if len(s) < 5 {
				return E.Left[int]("too short: " + s)
			}
			return E.Right[string](len(s))
		}
		partitionMap := PartitionMap(classifyString, errors.New("not classified"))

		// Act & Assert - short string
		result1 := partitionMap(Of("hi"))
		left1, right1 := P.Unpack(result1)
		assert.True(t, IsRight(left1))
		assert.True(t, IsLeft(right1))
		leftVal1, _ := Unwrap(left1)
		assert.Equal(t, "too short: hi", leftVal1)

		// Act & Assert - long string
		result2 := partitionMap(Of("hello world"))
		left2, right2 := P.Unpack(result2)
		assert.True(t, IsLeft(left2))
		assert.True(t, IsRight(right2))
		rightVal2, _ := Unwrap(right2)
		assert.Equal(t, 11, rightVal2)
	})
}

// Benchmark tests
func BenchmarkPartition(b *testing.B) {
	isPositive := N.MoreThan(0)
	partition := Partition(isPositive, errors.New("not positive"))
	input := Of(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partition(input)
	}
}

func BenchmarkPartitionError(b *testing.B) {
	isPositive := N.MoreThan(0)
	partition := Partition(isPositive, errors.New("not positive"))
	input := Left[int](errors.New("error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partition(input)
	}
}

func BenchmarkFilter(b *testing.B) {
	isPositive := N.MoreThan(0)
	filter := Filter(isPositive, errors.New("not positive"))
	input := Of(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(input)
	}
}

func BenchmarkFilterError(b *testing.B) {
	isPositive := N.MoreThan(0)
	filter := Filter(isPositive, errors.New("not positive"))
	input := Left[int](errors.New("error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter(input)
	}
}

func BenchmarkFilterChained(b *testing.B) {
	isPositive := N.MoreThan(0)
	isEven := func(n int) bool { return n%2 == 0 }
	filterPositive := Filter(isPositive, errors.New("not positive"))
	filterEven := Filter(isEven, errors.New("not even"))
	input := Of(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterEven(filterPositive(input))
	}
}

func BenchmarkFilterMap(b *testing.B) {
	parseInt := func(s string) O.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return O.Some(n)
		}
		return O.None[int]()
	}
	filterMap := FilterMap(parseInt, errors.New("invalid"))
	input := Of("42")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterMap(input)
	}
}

func BenchmarkFilterMapError(b *testing.B) {
	parseInt := func(s string) O.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return O.Some(n)
		}
		return O.None[int]()
	}
	filterMap := FilterMap(parseInt, errors.New("invalid"))
	input := Left[string](errors.New("error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filterMap(input)
	}
}

func BenchmarkPartitionMap(b *testing.B) {
	classify := func(n int) E.Either[string, int] {
		if n < 0 {
			return E.Left[int]("negative")
		}
		return E.Right[string](n * n)
	}
	partitionMap := PartitionMap(classify, errors.New("not classified"))
	input := Of(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partitionMap(input)
	}
}

func BenchmarkPartitionMapError(b *testing.B) {
	classify := func(n int) E.Either[string, int] {
		if n < 0 {
			return E.Left[int]("negative")
		}
		return E.Right[string](n * n)
	}
	partitionMap := PartitionMap(classify, errors.New("not classified"))
	input := Left[int](errors.New("error"))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = partitionMap(input)
	}
}
