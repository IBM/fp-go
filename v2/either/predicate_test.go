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
	"strings"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestExists_Success(t *testing.T) {
	t.Run("Right value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		input := Right[string](5)

		// Act
		result := hasPositive(input)

		// Assert
		assert.True(t, result, "should return true for Right value that passes predicate")
	})

	t.Run("Right value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		hasNonNegative := Exists[string](isNonNegative)
		input := Right[string](0)

		// Act
		result := hasNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Right value at boundary that passes predicate")
	})

	t.Run("Right value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		hasLongString := Exists[int](isLongString)
		input := Right[int]("hello world")

		// Act
		result := hasLongString(input)

		// Assert
		assert.True(t, result, "should return true for Right string that passes predicate")
	})

	t.Run("Right value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		hasEvenPositive := Exists[string](isEvenAndPositive)
		input := Right[string](4)

		// Act
		result := hasEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Right value that passes complex predicate")
	})
}

func TestExists_Failure(t *testing.T) {
	t.Run("Right value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		input := Right[string](-3)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Right value that fails predicate")
	})

	t.Run("Right value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		input := Right[string](0)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Right value at boundary that fails predicate")
	})

	t.Run("Left value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		input := Left[int]("error")

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Left value regardless of predicate")
	})

	t.Run("Left value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		input := Left[int]("error")

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Left value regardless of predicate")
	})

	t.Run("Right value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		hasLongString := Exists[int](isLongString)
		input := Right[int]("short")

		// Act
		result := hasLongString(input)

		// Assert
		assert.False(t, result, "should return false for Right string that fails predicate")
	})
}

func TestExists_EdgeCases(t *testing.T) {
	t.Run("Right with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		hasZero := Exists[string](isZero)
		input := Right[string](0)

		// Act
		result := hasZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Right with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		hasEmpty := Exists[int](isEmpty)
		input := Right[int]("")

		// Act
		result := hasEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Right with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		hasNil := Exists[string](isNil)
		input := Right[string]([]int(nil))

		// Act
		result := hasNil(input)

		// Assert
		assert.True(t, result, "should handle nil slice correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		hasAny := Exists[string](alwaysTrue)

		// Act & Assert
		assert.True(t, hasAny(Right[string](42)), "should return true for Right with always-true predicate")
		assert.False(t, hasAny(Left[int]("error")), "should return false for Left even with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		hasNone := Exists[string](alwaysFalse)

		// Act & Assert
		assert.False(t, hasNone(Right[string](42)), "should return false for Right with always-false predicate")
		assert.False(t, hasNone(Left[int]("error")), "should return false for Left with always-false predicate")
	})
}

func TestExists_Integration(t *testing.T) {
	t.Run("use in filtering slice of Eithers", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists[string](isPositive)
		values := []Either[string, int]{
			Right[string](5),
			Left[int]("error1"),
			Right[string](-3),
			Right[string](10),
			Left[int]("error2"),
			Right[string](0),
		}

		// Act
		var filtered []Either[string, int]
		for _, v := range values {
			if hasPositive(v) {
				filtered = append(filtered, v)
			}
		}

		// Assert
		assert.Len(t, filtered, 2, "should filter to only Right values with positive numbers")
		assert.Equal(t, Right[string](5), filtered[0])
		assert.Equal(t, Right[string](10), filtered[1])
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		hasEven := Exists[string](isEven)
		hasPositive := Exists[string](isPositive)

		input1 := Right[string](4)
		input2 := Right[string](3)
		input3 := Right[string](-4)

		// Act & Assert
		assert.True(t, hasEven(input1) && hasPositive(input1), "should pass both predicates")
		assert.False(t, hasEven(input2) && hasPositive(input2), "should fail even predicate")
		assert.False(t, hasEven(input3) && hasPositive(input3), "should fail positive predicate")
	})
}

func TestExistsLeft_Success(t *testing.T) {
	t.Run("Left value that passes predicate", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		input := Left[int]("validation: invalid input")

		// Act
		result := hasValidationError(input)

		// Assert
		assert.True(t, result, "should return true for Left value that passes predicate")
	})

	t.Run("Left value with numeric predicate", func(t *testing.T) {
		// Arrange
		isNegativeErrorCode := func(n int) bool { return n < 0 }
		hasNegativeCode := ExistsLeft[string](isNegativeErrorCode)
		input := Left[string](-404)

		// Act
		result := hasNegativeCode(input)

		// Assert
		assert.True(t, result, "should return true for Left numeric value that passes predicate")
	})

	t.Run("Left value with complex predicate", func(t *testing.T) {
		// Arrange
		isLongError := func(s string) bool { return len(s) > 10 && strings.Contains(s, "error") }
		hasLongError := ExistsLeft[int](isLongError)
		input := Left[int]("this is a long error message")

		// Act
		result := hasLongError(input)

		// Assert
		assert.True(t, result, "should return true for Left value that passes complex predicate")
	})
}

func TestExistsLeft_Failure(t *testing.T) {
	t.Run("Left value that fails predicate", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		input := Left[int]("network: connection failed")

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Left value that fails predicate")
	})

	t.Run("Right value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		input := Right[string](42)

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Right value regardless of predicate")
	})

	t.Run("Right value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		input := Right[string](42)

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Right value regardless of predicate")
	})

	t.Run("Left value with empty string predicate", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		hasEmpty := ExistsLeft[int](isEmpty)
		input := Left[int]("not empty")

		// Act
		result := hasEmpty(input)

		// Assert
		assert.False(t, result, "should return false for Left value that fails predicate")
	})
}

func TestExistsLeft_EdgeCases(t *testing.T) {
	t.Run("Left with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		hasEmpty := ExistsLeft[int](isEmpty)
		input := Left[int]("")

		// Act
		result := hasEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Left with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		hasZero := ExistsLeft[string](isZero)
		input := Left[string](0)

		// Act
		result := hasZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Left with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		hasNil := ExistsLeft[string](isNil)
		input := Left[string]([]int(nil))

		// Act
		result := hasNil(input)

		// Assert
		assert.True(t, result, "should handle nil slice correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(string) bool { return true }
		hasAny := ExistsLeft[int](alwaysTrue)

		// Act & Assert
		assert.True(t, hasAny(Left[int]("error")), "should return true for Left with always-true predicate")
		assert.False(t, hasAny(Right[string](42)), "should return false for Right even with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(string) bool { return false }
		hasNone := ExistsLeft[int](alwaysFalse)

		// Act & Assert
		assert.False(t, hasNone(Left[int]("error")), "should return false for Left with always-false predicate")
		assert.False(t, hasNone(Right[string](42)), "should return false for Right with always-false predicate")
	})
}

func TestExistsLeft_Integration(t *testing.T) {
	t.Run("use in error categorization", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		results := []Either[string, int]{
			Left[int]("validation: empty field"),
			Right[string](100),
			Left[int]("network: timeout"),
			Left[int]("validation: invalid format"),
			Right[string](200),
		}

		// Act
		var validationErrors []Either[string, int]
		for _, r := range results {
			if hasValidationError(r) {
				validationErrors = append(validationErrors, r)
			}
		}

		// Assert
		assert.Len(t, validationErrors, 2, "should filter to only validation errors")
		assert.Equal(t, Left[int]("validation: empty field"), validationErrors[0])
		assert.Equal(t, Left[int]("validation: invalid format"), validationErrors[1])
	})

	t.Run("combine with other error predicates", func(t *testing.T) {
		// Arrange
		isNetworkError := func(s string) bool {
			return strings.Contains(s, "network")
		}
		isTimeoutError := func(s string) bool {
			return strings.Contains(s, "timeout")
		}
		hasNetworkError := ExistsLeft[int](isNetworkError)
		hasTimeoutError := ExistsLeft[int](isTimeoutError)

		input1 := Left[int]("network: timeout")
		input2 := Left[int]("network: connection refused")
		input3 := Left[int]("validation: error")

		// Act & Assert
		assert.True(t, hasNetworkError(input1) && hasTimeoutError(input1), "should pass both predicates")
		assert.True(t, hasNetworkError(input2) && !hasTimeoutError(input2), "should pass network but not timeout")
		assert.False(t, hasNetworkError(input3) || hasTimeoutError(input3), "should fail both predicates")
	})

	t.Run("count errors by type", func(t *testing.T) {
		// Arrange
		isValidationError := func(s string) bool {
			return strings.HasPrefix(s, "validation:")
		}
		isNetworkError := func(s string) bool {
			return strings.HasPrefix(s, "network:")
		}
		hasValidationError := ExistsLeft[int](isValidationError)
		hasNetworkError := ExistsLeft[int](isNetworkError)

		results := []Either[string, int]{
			Left[int]("validation: error1"),
			Right[string](1),
			Left[int]("network: error2"),
			Left[int]("validation: error3"),
			Right[string](2),
			Left[int]("other: error4"),
		}

		// Act
		validationCount := 0
		networkCount := 0
		for _, r := range results {
			if hasValidationError(r) {
				validationCount++
			}
			if hasNetworkError(r) {
				networkCount++
			}
		}

		// Assert
		assert.Equal(t, 2, validationCount, "should count validation errors correctly")
		assert.Equal(t, 1, networkCount, "should count network errors correctly")
	})
}

func BenchmarkExists(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists[string](isPositive)
	input := Right[string](42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsLeft(b *testing.B) {
	isValidationError := func(s string) bool {
		return strings.HasPrefix(s, "validation:")
	}
	hasValidationError := ExistsLeft[int](isValidationError)
	input := Left[int]("validation: error")

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func BenchmarkExistsPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists[string](isPositive)
	input := Right[string](-42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsLeftPredicateFails(b *testing.B) {
	isValidationError := func(s string) bool {
		return strings.HasPrefix(s, "validation:")
	}
	hasValidationError := ExistsLeft[int](isValidationError)
	input := Left[int]("network: error")

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func BenchmarkExistsOnLeft(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists[string](isPositive)
	input := Left[int]("error")

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsLeftOnRight(b *testing.B) {
	isValidationError := func(s string) bool {
		return strings.HasPrefix(s, "validation:")
	}
	hasValidationError := ExistsLeft[int](isValidationError)
	input := Right[string](42)

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func TestForAll_Success(t *testing.T) {
	t.Run("Right value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		input := Right[string](5)

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Right value that passes predicate")
	})

	t.Run("Right value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		allNonNegative := ForAll[string](isNonNegative)
		input := Right[string](0)

		// Act
		result := allNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Right value at boundary that passes predicate")
	})

	t.Run("Left value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		input := Left[int]("error")

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Left value (vacuous truth)")
	})

	t.Run("Left value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		input := Left[int]("error")

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Left value regardless of predicate (vacuous truth)")
	})

	t.Run("Right value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		allLongString := ForAll[int](isLongString)
		input := Right[int]("hello world")

		// Act
		result := allLongString(input)

		// Assert
		assert.True(t, result, "should return true for Right string that passes predicate")
	})

	t.Run("Right value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		allEvenPositive := ForAll[string](isEvenAndPositive)
		input := Right[string](4)

		// Act
		result := allEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Right value that passes complex predicate")
	})
}

func TestForAll_Failure(t *testing.T) {
	t.Run("Right value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		input := Right[string](-3)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Right value that fails predicate")
	})

	t.Run("Right value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		input := Right[string](0)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Right value at boundary that fails predicate")
	})

	t.Run("Right value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		allLongString := ForAll[int](isLongString)
		input := Right[int]("short")

		// Act
		result := allLongString(input)

		// Assert
		assert.False(t, result, "should return false for Right string that fails predicate")
	})

	t.Run("Right value with complex predicate that fails", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		allEvenPositive := ForAll[string](isEvenAndPositive)
		input := Right[string](3)

		// Act
		result := allEvenPositive(input)

		// Assert
		assert.False(t, result, "should return false for Right value that fails complex predicate")
	})
}

func TestForAll_EdgeCases(t *testing.T) {
	t.Run("Right with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		allZero := ForAll[string](isZero)
		input := Right[string](0)

		// Act
		result := allZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Right with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		allEmpty := ForAll[int](isEmpty)
		input := Right[int]("")

		// Act
		result := allEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Right with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		allNil := ForAll[string](isNil)
		input := Right[string]([]int(nil))

		// Act
		result := allNil(input)

		// Assert
		assert.True(t, result, "should handle nil slice correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(int) bool { return true }
		allTrue := ForAll[string](alwaysTrue)

		// Act & Assert
		assert.True(t, allTrue(Right[string](42)), "should return true for Right with always-true predicate")
		assert.True(t, allTrue(Left[int]("error")), "should return true for Left with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		allFalse := ForAll[string](alwaysFalse)

		// Act & Assert
		assert.False(t, allFalse(Right[string](42)), "should return false for Right with always-false predicate")
		assert.True(t, allFalse(Left[int]("error")), "should return true for Left even with always-false predicate (vacuous truth)")
	})

	t.Run("Left with various error types", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)

		// Act & Assert
		assert.True(t, allPositive(Left[int]("")), "should return true for Left with empty string error")
		assert.True(t, allPositive(Left[int]("error")), "should return true for Left with string error")
		assert.True(t, allPositive(Left[int]("validation: failed")), "should return true for Left with any error")
	})
}

func TestForAll_Integration(t *testing.T) {
	t.Run("validate all successful results meet criteria", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		values := []Either[string, int]{
			Right[string](5),
			Left[int]("error1"),
			Right[string](10),
			Left[int]("error2"),
			Right[string](3),
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
		assert.True(t, allValid, "should return true when all Right values pass predicate (Left values ignored)")
	})

	t.Run("detect when any Right value fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		values := []Either[string, int]{
			Right[string](5),
			Left[int]("error1"),
			Right[string](-3), // This fails
			Right[string](10),
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
		assert.False(t, allValid, "should return false when any Right value fails predicate")
	})

	t.Run("all Left values pass vacuously", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		values := []Either[string, int]{
			Left[int]("error1"),
			Left[int]("error2"),
			Left[int]("error3"),
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
		assert.True(t, allValid, "should return true for all Left values (vacuous truth)")
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		allEven := ForAll[string](isEven)
		allPositive := ForAll[string](isPositive)

		input1 := Right[string](4)
		input2 := Right[string](3)
		input3 := Right[string](-4)
		input4 := Left[int]("error")

		// Act & Assert
		assert.True(t, allEven(input1) && allPositive(input1), "should pass both predicates")
		assert.False(t, allEven(input2) && allPositive(input2), "should fail even predicate")
		assert.False(t, allEven(input3) && allPositive(input3), "should fail positive predicate")
		assert.True(t, allEven(input4) && allPositive(input4), "Left should pass both predicates (vacuous truth)")
	})

	t.Run("contrast with Exists", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll[string](isPositive)
		hasPositive := Exists[string](isPositive)

		leftValue := Left[int]("error")
		rightPositive := Right[string](5)
		rightNegative := Right[string](-3)

		// Act & Assert - ForAll behavior
		assert.True(t, allPositive(leftValue), "ForAll: Left is vacuously true")
		assert.True(t, allPositive(rightPositive), "ForAll: Right with passing predicate is true")
		assert.False(t, allPositive(rightNegative), "ForAll: Right with failing predicate is false")

		// Act & Assert - Exists behavior (contrast)
		assert.False(t, hasPositive(leftValue), "Exists: Left is false")
		assert.True(t, hasPositive(rightPositive), "Exists: Right with passing predicate is true")
		assert.False(t, hasPositive(rightNegative), "Exists: Right with failing predicate is false")
	})

	t.Run("De Morgan's law: ForAll(p) ≡ not(Exists(not(p)))", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isNotPositive := func(n int) bool { return !isPositive(n) }
		allPositive := ForAll[string](isPositive)
		hasNotPositive := Exists[string](isNotPositive)

		testCases := []Either[string, int]{
			Right[string](5),
			Right[string](-3),
			Left[int]("error"),
		}

		// Act & Assert - Verify De Morgan's law
		for _, tc := range testCases {
			forAllResult := allPositive(tc)
			notExistsNotResult := !hasNotPositive(tc)
			assert.Equal(t, forAllResult, notExistsNotResult,
				"ForAll(p) should equal not(Exists(not(p))) for %v", tc)
		}
	})
}

func TestForAll_WithComplexTypes(t *testing.T) {
	t.Run("with struct type", func(t *testing.T) {
		// Arrange
		type Point struct {
			X, Y int
		}
		isInFirstQuadrant := func(p Point) bool { return p.X > 0 && p.Y > 0 }
		allInFirstQuadrant := ForAll[string](isInFirstQuadrant)

		// Act & Assert
		assert.True(t, allInFirstQuadrant(Right[string](Point{1, 1})), "point in first quadrant passes")
		assert.False(t, allInFirstQuadrant(Right[string](Point{-1, 1})), "point outside first quadrant fails")
		assert.True(t, allInFirstQuadrant(Left[Point]("error")), "Left passes vacuously")
	})

	t.Run("with slice type", func(t *testing.T) {
		// Arrange
		hasElements := func(s []int) bool { return len(s) > 0 }
		allNonEmpty := ForAll[string](hasElements)

		// Act & Assert
		assert.True(t, allNonEmpty(Right[string]([]int{1, 2, 3})), "non-empty slice passes")
		assert.False(t, allNonEmpty(Right[string]([]int{})), "empty slice fails")
		assert.True(t, allNonEmpty(Left[[]int]("error")), "Left passes vacuously")
	})

	t.Run("with map type", func(t *testing.T) {
		// Arrange
		hasKey := func(key string) func(map[string]int) bool {
			return func(m map[string]int) bool {
				_, exists := m[key]
				return exists
			}
		}
		allHaveAgeKey := ForAll[string](hasKey("age"))

		// Act & Assert
		assert.True(t, allHaveAgeKey(Right[string](map[string]int{"age": 25})), "map with key passes")
		assert.False(t, allHaveAgeKey(Right[string](map[string]int{"name": 1})), "map without key fails")
		assert.True(t, allHaveAgeKey(Left[map[string]int]("error")), "Left passes vacuously")
	})
}

func BenchmarkForAll(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll[string](isPositive)
	input := Right[string](42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll[string](isPositive)
	input := Right[string](-42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllOnLeft(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll[string](isPositive)
	input := Left[int]("error")

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllComplexPredicate(b *testing.B) {
	isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
	allEvenPositive := ForAll[string](isEvenAndPositive)
	input := Right[string](42)

	b.ResetTimer()
	for range b.N {
		_ = allEvenPositive(input)
	}
}

// Made with Bob
