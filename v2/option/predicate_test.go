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

package option

import (
	"strings"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestExists_Success(t *testing.T) {
	t.Run("Some value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Some(5)

		// Act
		result := hasPositive(input)

		// Assert
		assert.True(t, result, "should return true for Some value that passes predicate")
	})

	t.Run("Some value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		hasNonNegative := Exists(isNonNegative)
		input := Some(0)

		// Act
		result := hasNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Some value at boundary that passes predicate")
	})

	t.Run("Some value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		hasLongString := Exists(isLongString)
		input := Some("hello world")

		// Act
		result := hasLongString(input)

		// Assert
		assert.True(t, result, "should return true for Some string that passes predicate")
	})

	t.Run("Some value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		hasEvenPositive := Exists(isEvenAndPositive)
		input := Some(4)

		// Act
		result := hasEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Some value that passes complex predicate")
	})
}

func TestExists_Failure(t *testing.T) {
	t.Run("Some value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Some(-3)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Some value that fails predicate")
	})

	t.Run("Some value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Some(0)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Some value at boundary that fails predicate")
	})

	t.Run("None value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := None[int]()

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for None value regardless of predicate")
	})

	t.Run("None value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := None[int]()

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for None value regardless of predicate")
	})

	t.Run("Some value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		hasLongString := Exists(isLongString)
		input := Some("short")

		// Act
		result := hasLongString(input)

		// Assert
		assert.False(t, result, "should return false for Some string that fails predicate")
	})
}

func TestExists_EdgeCases(t *testing.T) {
	t.Run("Some with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		hasZero := Exists(isZero)
		input := Some(0)

		// Act
		result := hasZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Some with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		hasEmpty := Exists(isEmpty)
		input := Some("")

		// Act
		result := hasEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Some with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		hasNil := Exists(isNil)
		input := Some([]int(nil))

		// Act
		result := hasNil(input)

		// Assert
		assert.True(t, result, "should handle nil slice correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		hasAny := Exists(alwaysTrue)

		// Act & Assert
		assert.True(t, hasAny(Some(42)), "should return true for Some with always-true predicate")
		assert.False(t, hasAny(None[int]()), "should return false for None even with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		hasNone := Exists(alwaysFalse)

		// Act & Assert
		assert.False(t, hasNone(Some(42)), "should return false for Some with always-false predicate")
		assert.False(t, hasNone(None[int]()), "should return false for None with always-false predicate")
	})

	t.Run("Some with false boolean value", func(t *testing.T) {
		// Arrange
		isFalse := func(b bool) bool { return !b }
		hasFalse := Exists(isFalse)
		input := Some(false)

		// Act
		result := hasFalse(input)

		// Assert
		assert.True(t, result, "should handle false boolean value correctly")
	})
}

func TestExists_Integration(t *testing.T) {
	t.Run("use in filtering slice of Options", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		values := []Option[int]{
			Some(5),
			None[int](),
			Some(-3),
			Some(10),
			None[int](),
			Some(0),
		}

		// Act
		var filtered []Option[int]
		for _, v := range values {
			if hasPositive(v) {
				filtered = append(filtered, v)
			}
		}

		// Assert
		assert.Len(t, filtered, 2, "should filter to only Some values with positive numbers")
		assert.Equal(t, Some(5), filtered[0])
		assert.Equal(t, Some(10), filtered[1])
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		hasEven := Exists(isEven)
		hasPositive := Exists(isPositive)

		input1 := Some(4)
		input2 := Some(3)
		input3 := Some(-4)

		// Act & Assert
		assert.True(t, hasEven(input1) && hasPositive(input1), "should pass both predicates")
		assert.False(t, hasEven(input2) && hasPositive(input2), "should fail even predicate")
		assert.False(t, hasEven(input3) && hasPositive(input3), "should fail positive predicate")
	})

	t.Run("use with string operations", func(t *testing.T) {
		// Arrange
		hasPrefix := func(prefix string) func(string) bool {
			return func(s string) bool {
				return strings.HasPrefix(s, prefix)
			}
		}
		hasHelloPrefix := Exists(hasPrefix("hello"))

		values := []Option[string]{
			Some("hello world"),
			None[string](),
			Some("goodbye"),
			Some("hello there"),
		}

		// Act
		var filtered []Option[string]
		for _, v := range values {
			if hasHelloPrefix(v) {
				filtered = append(filtered, v)
			}
		}

		// Assert
		assert.Len(t, filtered, 2, "should filter strings with 'hello' prefix")
	})

	t.Run("count values matching predicate", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		hasEven := Exists(isEven)

		values := []Option[int]{
			Some(2),
			Some(3),
			None[int](),
			Some(4),
			Some(5),
			None[int](),
			Some(6),
		}

		// Act
		count := 0
		for _, v := range values {
			if hasEven(v) {
				count++
			}
		}

		// Assert
		assert.Equal(t, 3, count, "should count even numbers correctly")
	})

	t.Run("use in validation chain", func(t *testing.T) {
		// Arrange
		type User struct {
			Name string
			Age  int
		}

		isAdult := func(u User) bool { return u.Age >= 18 }
		hasValidName := func(u User) bool { return len(u.Name) > 0 }

		hasAdult := Exists(isAdult)
		hasName := Exists(hasValidName)

		validUser := Some(User{Name: "Alice", Age: 25})
		minorUser := Some(User{Name: "Bob", Age: 15})
		noNameUser := Some(User{Name: "", Age: 30})
		noneUser := None[User]()

		// Act & Assert
		assert.True(t, hasAdult(validUser) && hasName(validUser), "valid user passes all checks")
		assert.False(t, hasAdult(minorUser), "minor fails adult check")
		assert.False(t, hasName(noNameUser), "user without name fails name check")
		assert.False(t, hasAdult(noneUser) || hasName(noneUser), "None fails all checks")
	})
}

func TestExists_WithComplexTypes(t *testing.T) {
	t.Run("with struct type", func(t *testing.T) {
		// Arrange
		type Point struct {
			X, Y int
		}
		isOrigin := func(p Point) bool { return p.X == 0 && p.Y == 0 }
		hasOrigin := Exists(isOrigin)

		// Act & Assert
		assert.True(t, hasOrigin(Some(Point{0, 0})), "origin point passes")
		assert.False(t, hasOrigin(Some(Point{1, 0})), "non-origin point fails")
		assert.False(t, hasOrigin(None[Point]()), "None fails")
	})

	t.Run("with slice type", func(t *testing.T) {
		// Arrange
		hasElements := func(s []int) bool { return len(s) > 0 }
		hasNonEmpty := Exists(hasElements)

		// Act & Assert
		assert.True(t, hasNonEmpty(Some([]int{1, 2, 3})), "non-empty slice passes")
		assert.False(t, hasNonEmpty(Some([]int{})), "empty slice fails")
		assert.False(t, hasNonEmpty(None[[]int]()), "None fails")
	})

	t.Run("with map type", func(t *testing.T) {
		// Arrange
		hasKey := func(key string) func(map[string]int) bool {
			return func(m map[string]int) bool {
				_, exists := m[key]
				return exists
			}
		}
		hasAgeKey := Exists(hasKey("age"))

		// Act & Assert
		assert.True(t, hasAgeKey(Some(map[string]int{"age": 25})), "map with key passes")
		assert.False(t, hasAgeKey(Some(map[string]int{"name": 1})), "map without key fails")
		assert.False(t, hasAgeKey(None[map[string]int]()), "None fails")
	})
}

func BenchmarkExists(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := Some(42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := Some(-42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsOnNone(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := None[int]()

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsComplexPredicate(b *testing.B) {
	isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
	hasEvenPositive := Exists(isEvenAndPositive)
	input := Some(42)

	b.ResetTimer()
	for range b.N {
		_ = hasEvenPositive(input)
	}
}

func BenchmarkExistsStringPredicate(b *testing.B) {
	isLongString := func(s string) bool { return len(s) > 10 }
	hasLongString := Exists(isLongString)
	input := Some("hello world from benchmark")

	b.ResetTimer()
	for range b.N {
		_ = hasLongString(input)
	}
}

func TestForAll_Success(t *testing.T) {
	t.Run("Some value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Some(5)

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Some value that passes predicate")
	})

	t.Run("Some value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		allNonNegative := ForAll(isNonNegative)
		input := Some(0)

		// Act
		result := allNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Some value at boundary that passes predicate")
	})

	t.Run("None value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := None[int]()

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for None value (vacuous truth)")
	})

	t.Run("None value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := None[int]()

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for None value regardless of predicate (vacuous truth)")
	})

	t.Run("Some value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		allLongString := ForAll(isLongString)
		input := Some("hello world")

		// Act
		result := allLongString(input)

		// Assert
		assert.True(t, result, "should return true for Some string that passes predicate")
	})

	t.Run("Some value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		allEvenPositive := ForAll(isEvenAndPositive)
		input := Some(4)

		// Act
		result := allEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Some value that passes complex predicate")
	})
}

func TestForAll_Failure(t *testing.T) {
	t.Run("Some value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Some(-3)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Some value that fails predicate")
	})

	t.Run("Some value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Some(0)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Some value at boundary that fails predicate")
	})

	t.Run("Some value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		allLongString := ForAll(isLongString)
		input := Some("short")

		// Act
		result := allLongString(input)

		// Assert
		assert.False(t, result, "should return false for Some string that fails predicate")
	})
}

func TestForAll_EdgeCases(t *testing.T) {
	t.Run("Some with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		allZero := ForAll(isZero)
		input := Some(0)

		// Act
		result := allZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Some with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		allEmpty := ForAll(isEmpty)
		input := Some("")

		// Act
		result := allEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Some with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		allNil := ForAll(isNil)
		input := Some([]int(nil))

		// Act
		result := allNil(input)

		// Assert
		assert.True(t, result, "should handle nil slice correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		allTrue := ForAll(alwaysTrue)

		// Act & Assert
		assert.True(t, allTrue(Some(42)), "should return true for Some with always-true predicate")
		assert.True(t, allTrue(None[int]()), "should return true for None with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		allFalse := ForAll(alwaysFalse)

		// Act & Assert
		assert.False(t, allFalse(Some(42)), "should return false for Some with always-false predicate")
		assert.True(t, allFalse(None[int]()), "should return true for None even with always-false predicate")
	})

	t.Run("Some with false boolean value", func(t *testing.T) {
		// Arrange
		isFalse := func(b bool) bool { return !b }
		allFalse := ForAll(isFalse)
		input := Some(false)

		// Act
		result := allFalse(input)

		// Assert
		assert.True(t, result, "should handle false boolean value correctly")
	})
}

func TestForAll_Integration(t *testing.T) {
	t.Run("validate all values in slice of Options", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		values := []Option[int]{
			Some(5),
			None[int](),
			Some(10),
			None[int](),
			Some(3),
		}

		// Act
		allValid := true
		for _, v := range values {
			if !allPositive(v) {
				allValid = false
				break
			}
		}

		// Assert
		assert.True(t, allValid, "should return true when all Some values pass predicate")
	})

	t.Run("detect invalid value in slice of Options", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		values := []Option[int]{
			Some(5),
			None[int](),
			Some(-3),
			Some(10),
		}

		// Act
		allValid := true
		for _, v := range values {
			if !allPositive(v) {
				allValid = false
				break
			}
		}

		// Assert
		assert.False(t, allValid, "should return false when any Some value fails predicate")
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		allEven := ForAll(isEven)
		allPositive := ForAll(isPositive)

		input1 := Some(4)
		input2 := Some(3)
		input3 := Some(-4)
		input4 := None[int]()

		// Act & Assert
		assert.True(t, allEven(input1) && allPositive(input1), "should pass both predicates")
		assert.False(t, allEven(input2) && allPositive(input2), "should fail even predicate")
		assert.False(t, allEven(input3) && allPositive(input3), "should fail positive predicate")
		assert.True(t, allEven(input4) && allPositive(input4), "None passes all predicates")
	})

	t.Run("use with string operations", func(t *testing.T) {
		// Arrange
		hasPrefix := func(prefix string) func(string) bool {
			return func(s string) bool {
				return strings.HasPrefix(s, prefix)
			}
		}
		allHaveHelloPrefix := ForAll(hasPrefix("hello"))

		values := []Option[string]{
			Some("hello world"),
			None[string](),
			Some("hello there"),
			None[string](),
		}

		// Act
		allValid := true
		for _, v := range values {
			if !allHaveHelloPrefix(v) {
				allValid = false
				break
			}
		}

		// Assert
		assert.True(t, allValid, "should return true when all Some strings have prefix")
	})

	t.Run("count values failing predicate", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		allEven := ForAll(isEven)

		values := []Option[int]{
			Some(2),
			Some(3),
			None[int](),
			Some(4),
			Some(5),
			None[int](),
			Some(6),
		}

		// Act
		failCount := 0
		for _, v := range values {
			if !allEven(v) {
				failCount++
			}
		}

		// Assert
		assert.Equal(t, 2, failCount, "should count odd numbers correctly")
	})

	t.Run("use in validation chain", func(t *testing.T) {
		// Arrange
		type User struct {
			Name string
			Age  int
		}

		isAdult := func(u User) bool { return u.Age >= 18 }
		hasValidName := func(u User) bool { return len(u.Name) > 0 }

		allAdult := ForAll(isAdult)
		allHaveName := ForAll(hasValidName)

		validUser := Some(User{Name: "Alice", Age: 25})
		minorUser := Some(User{Name: "Bob", Age: 15})
		noNameUser := Some(User{Name: "", Age: 30})
		noneUser := None[User]()

		// Act & Assert
		assert.True(t, allAdult(validUser) && allHaveName(validUser), "valid user passes all checks")
		assert.False(t, allAdult(minorUser), "minor fails adult check")
		assert.False(t, allHaveName(noNameUser), "user without name fails name check")
		assert.True(t, allAdult(noneUser) && allHaveName(noneUser), "None passes all checks")
	})
}

func TestForAll_WithComplexTypes(t *testing.T) {
	t.Run("with struct type", func(t *testing.T) {
		// Arrange
		type Point struct {
			X, Y int
		}
		isOrigin := func(p Point) bool { return p.X == 0 && p.Y == 0 }
		allOrigin := ForAll(isOrigin)

		// Act & Assert
		assert.True(t, allOrigin(Some(Point{0, 0})), "origin point passes")
		assert.False(t, allOrigin(Some(Point{1, 0})), "non-origin point fails")
		assert.True(t, allOrigin(None[Point]()), "None passes (vacuous truth)")
	})

	t.Run("with slice type", func(t *testing.T) {
		// Arrange
		hasElements := func(s []int) bool { return len(s) > 0 }
		allNonEmpty := ForAll(hasElements)

		// Act & Assert
		assert.True(t, allNonEmpty(Some([]int{1, 2, 3})), "non-empty slice passes")
		assert.False(t, allNonEmpty(Some([]int{})), "empty slice fails")
		assert.True(t, allNonEmpty(None[[]int]()), "None passes (vacuous truth)")
	})

	t.Run("with map type", func(t *testing.T) {
		// Arrange
		hasKey := func(key string) func(map[string]int) bool {
			return func(m map[string]int) bool {
				_, exists := m[key]
				return exists
			}
		}
		allHaveAgeKey := ForAll(hasKey("age"))

		// Act & Assert
		assert.True(t, allHaveAgeKey(Some(map[string]int{"age": 25})), "map with key passes")
		assert.False(t, allHaveAgeKey(Some(map[string]int{"name": 1})), "map without key fails")
		assert.True(t, allHaveAgeKey(None[map[string]int]()), "None passes (vacuous truth)")
	})
}

func TestForAll_DeMorganLaws(t *testing.T) {
	t.Run("ForAll(p) ≡ not(Exists(not(p)))", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isNotPositive := func(n int) bool { return !isPositive(n) }

		allPositive := ForAll(isPositive)
		hasNotPositive := Exists(isNotPositive)

		testCases := []Option[int]{
			Some(5),
			Some(-3),
			Some(0),
			None[int](),
		}

		// Act & Assert
		for _, tc := range testCases {
			result1 := allPositive(tc)
			result2 := !hasNotPositive(tc)
			assert.Equal(t, result1, result2, "ForAll(p) should equal not(Exists(not(p)))")
		}
	})

	t.Run("Exists(p) ≡ not(ForAll(not(p)))", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isNotPositive := func(n int) bool { return !isPositive(n) }

		hasPositive := Exists(isPositive)
		allNotPositive := ForAll(isNotPositive)

		testCases := []Option[int]{
			Some(5),
			Some(-3),
			Some(0),
			None[int](),
		}

		// Act & Assert
		for _, tc := range testCases {
			result1 := hasPositive(tc)
			result2 := !allNotPositive(tc)
			assert.Equal(t, result1, result2, "Exists(p) should equal not(ForAll(not(p)))")
		}
	})
}

func BenchmarkForAll(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := Some(42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := Some(-42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllOnNone(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := None[int]()

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllComplexPredicate(b *testing.B) {
	isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
	allEvenPositive := ForAll(isEvenAndPositive)
	input := Some(42)

	b.ResetTimer()
	for range b.N {
		_ = allEvenPositive(input)
	}
}

func BenchmarkForAllStringPredicate(b *testing.B) {
	isLongString := func(s string) bool { return len(s) > 10 }
	allLongString := ForAll(isLongString)
	input := Some("hello world from benchmark")

	b.ResetTimer()
	for range b.N {
		_ = allLongString(input)
	}
}
