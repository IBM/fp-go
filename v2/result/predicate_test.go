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
	"strings"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestExists_Success(t *testing.T) {
	t.Run("Ok value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Of(5)

		// Act
		result := hasPositive(input)

		// Assert
		assert.True(t, result, "should return true for Ok value that passes predicate")
	})

	t.Run("Ok value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		hasNonNegative := Exists(isNonNegative)
		input := Of(0)

		// Act
		result := hasNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Ok value at boundary that passes predicate")
	})

	t.Run("Ok value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		hasLongString := Exists(isLongString)
		input := Of("hello world")

		// Act
		result := hasLongString(input)

		// Assert
		assert.True(t, result, "should return true for Ok string that passes predicate")
	})

	t.Run("Ok value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		hasEvenPositive := Exists(isEvenAndPositive)
		input := Of(4)

		// Act
		result := hasEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Ok value that passes complex predicate")
	})
}

func TestExists_Failure(t *testing.T) {
	t.Run("Ok value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Of(-3)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Ok value that fails predicate")
	})

	t.Run("Ok value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Of(0)

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Ok value at boundary that fails predicate")
	})

	t.Run("Error value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Left[int](errors.New("error"))

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Error value regardless of predicate")
	})

	t.Run("Error value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		input := Left[int](errors.New("error"))

		// Act
		result := hasPositive(input)

		// Assert
		assert.False(t, result, "should return false for Error value regardless of predicate")
	})

	t.Run("Ok value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		hasLongString := Exists(isLongString)
		input := Of("short")

		// Act
		result := hasLongString(input)

		// Assert
		assert.False(t, result, "should return false for Ok string that fails predicate")
	})
}

func TestExists_EdgeCases(t *testing.T) {
	t.Run("Ok with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		hasZero := Exists(isZero)
		input := Of(0)

		// Act
		result := hasZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Ok with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		hasEmpty := Exists(isEmpty)
		input := Of("")

		// Act
		result := hasEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Ok with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		hasNil := Exists(isNil)
		input := Of([]int(nil))

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
		assert.True(t, hasAny(Of(42)), "should return true for Ok with always-true predicate")
		assert.False(t, hasAny(Left[int](errors.New("error"))), "should return false for Error even with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		hasNone := Exists(alwaysFalse)

		// Act & Assert
		assert.False(t, hasNone(Of(42)), "should return false for Ok with always-false predicate")
		assert.False(t, hasNone(Left[int](errors.New("error"))), "should return false for Error with always-false predicate")
	})
}

func TestExists_Integration(t *testing.T) {
	t.Run("use in filtering slice of Results", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		hasPositive := Exists(isPositive)
		values := []Result[int]{
			Of(5),
			Left[int](errors.New("error1")),
			Of(-3),
			Of(10),
			Left[int](errors.New("error2")),
			Of(0),
		}

		// Act
		var filtered []Result[int]
		for _, v := range values {
			if hasPositive(v) {
				filtered = append(filtered, v)
			}
		}

		// Assert
		assert.Len(t, filtered, 2, "should filter to only Ok values with positive numbers")
		assert.Equal(t, Of(5), filtered[0])
		assert.Equal(t, Of(10), filtered[1])
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		hasEven := Exists(isEven)
		hasPositive := Exists(isPositive)

		input1 := Of(4)
		input2 := Of(3)
		input3 := Of(-4)

		// Act & Assert
		assert.True(t, hasEven(input1) && hasPositive(input1), "should pass both predicates")
		assert.False(t, hasEven(input2) && hasPositive(input2), "should fail even predicate")
		assert.False(t, hasEven(input3) && hasPositive(input3), "should fail positive predicate")
	})
}

func TestExistsError_Success(t *testing.T) {
	t.Run("Error value that passes predicate", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.Contains(err.Error(), "validation")
		}
		hasValidationError := ExistsError[int](isValidationError)
		input := Left[int](errors.New("validation: invalid input"))

		// Act
		result := hasValidationError(input)

		// Assert
		assert.True(t, result, "should return true for Error value that passes predicate")
	})

	t.Run("Error value with complex predicate", func(t *testing.T) {
		// Arrange
		isLongError := func(err error) bool {
			msg := err.Error()
			return len(msg) > 10 && strings.Contains(msg, "error")
		}
		hasLongError := ExistsError[int](isLongError)
		input := Left[int](errors.New("this is a long error message"))

		// Act
		result := hasLongError(input)

		// Assert
		assert.True(t, result, "should return true for Error value that passes complex predicate")
	})

	t.Run("Error value with prefix check", func(t *testing.T) {
		// Arrange
		hasPrefix := func(prefix string) func(error) bool {
			return func(err error) bool {
				return strings.HasPrefix(err.Error(), prefix)
			}
		}
		hasNetworkError := ExistsError[string](hasPrefix("network:"))
		input := Left[string](errors.New("network: connection failed"))

		// Act
		result := hasNetworkError(input)

		// Assert
		assert.True(t, result, "should return true for Error with matching prefix")
	})
}

func TestExistsError_Failure(t *testing.T) {
	t.Run("Error value that fails predicate", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.Contains(err.Error(), "validation")
		}
		hasValidationError := ExistsError[int](isValidationError)
		input := Left[int](errors.New("network: connection failed"))

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Error value that fails predicate")
	})

	t.Run("Ok value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.Contains(err.Error(), "validation")
		}
		hasValidationError := ExistsError[int](isValidationError)
		input := Of(42)

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Ok value regardless of predicate")
	})

	t.Run("Ok value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.Contains(err.Error(), "validation")
		}
		hasValidationError := ExistsError[int](isValidationError)
		input := Of(42)

		// Act
		result := hasValidationError(input)

		// Assert
		assert.False(t, result, "should return false for Ok value regardless of predicate")
	})

	t.Run("Error value with empty string check", func(t *testing.T) {
		// Arrange
		isEmpty := func(err error) bool { return len(err.Error()) == 0 }
		hasEmpty := ExistsError[int](isEmpty)
		input := Left[int](errors.New("not empty"))

		// Act
		result := hasEmpty(input)

		// Assert
		assert.False(t, result, "should return false for Error value that fails predicate")
	})
}

func TestExistsError_EdgeCases(t *testing.T) {
	t.Run("Error with empty message", func(t *testing.T) {
		// Arrange
		isEmpty := func(err error) bool { return len(err.Error()) == 0 }
		hasEmpty := ExistsError[int](isEmpty)
		input := Left[int](errors.New(""))

		// Act
		result := hasEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty error message correctly")
	})

	t.Run("predicate always returns true", func(t *testing.T) {
		// Arrange
		alwaysTrue := func(error) bool { return true }
		hasAny := ExistsError[int](alwaysTrue)

		// Act & Assert
		assert.True(t, hasAny(Left[int](errors.New("error"))), "should return true for Error with always-true predicate")
		assert.False(t, hasAny(Of(42)), "should return false for Ok even with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(error) bool { return false }
		hasNone := ExistsError[int](alwaysFalse)

		// Act & Assert
		assert.False(t, hasNone(Left[int](errors.New("error"))), "should return false for Error with always-false predicate")
		assert.False(t, hasNone(Of(42)), "should return false for Ok with always-false predicate")
	})

	t.Run("nil error check", func(t *testing.T) {
		// Arrange
		isNil := func(err error) bool { return err == nil }
		hasNil := ExistsError[int](isNil)

		// Note: In practice, Result[T] should never contain a nil error in the Error case,
		// but we test the predicate behavior
		input := Left[int](errors.New("not nil"))

		// Act
		result := hasNil(input)

		// Assert
		assert.False(t, result, "should return false for non-nil error")
	})
}

func TestExistsError_Integration(t *testing.T) {
	t.Run("use in error categorization", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.Contains(err.Error(), "validation:")
		}
		hasValidationError := ExistsError[int](isValidationError)
		results := []Result[int]{
			Left[int](errors.New("validation: empty field")),
			Of(100),
			Left[int](errors.New("network: timeout")),
			Left[int](errors.New("validation: invalid format")),
			Of(200),
		}

		// Act
		var validationErrors []Result[int]
		for _, r := range results {
			if hasValidationError(r) {
				validationErrors = append(validationErrors, r)
			}
		}

		// Assert
		assert.Len(t, validationErrors, 2, "should filter to only validation errors")
		assert.Equal(t, Left[int](errors.New("validation: empty field")), validationErrors[0])
		assert.Equal(t, Left[int](errors.New("validation: invalid format")), validationErrors[1])
	})

	t.Run("combine with other error predicates", func(t *testing.T) {
		// Arrange
		isNetworkError := func(err error) bool {
			return strings.Contains(err.Error(), "network")
		}
		isTimeoutError := func(err error) bool {
			return strings.Contains(err.Error(), "timeout")
		}
		hasNetworkError := ExistsError[int](isNetworkError)
		hasTimeoutError := ExistsError[int](isTimeoutError)

		input1 := Left[int](errors.New("network: timeout"))
		input2 := Left[int](errors.New("network: connection refused"))
		input3 := Left[int](errors.New("validation: error"))

		// Act & Assert
		assert.True(t, hasNetworkError(input1) && hasTimeoutError(input1), "should pass both predicates")
		assert.True(t, hasNetworkError(input2) && !hasTimeoutError(input2), "should pass network but not timeout")
		assert.False(t, hasNetworkError(input3) || hasTimeoutError(input3), "should fail both predicates")
	})

	t.Run("count errors by type", func(t *testing.T) {
		// Arrange
		isValidationError := func(err error) bool {
			return strings.HasPrefix(err.Error(), "validation:")
		}
		isNetworkError := func(err error) bool {
			return strings.HasPrefix(err.Error(), "network:")
		}
		hasValidationError := ExistsError[int](isValidationError)
		hasNetworkError := ExistsError[int](isNetworkError)

		results := []Result[int]{
			Left[int](errors.New("validation: error1")),
			Of(1),
			Left[int](errors.New("network: error2")),
			Left[int](errors.New("validation: error3")),
			Of(2),
			Left[int](errors.New("other: error4")),
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

	t.Run("error severity classification", func(t *testing.T) {
		// Arrange
		isCritical := func(err error) bool {
			return strings.Contains(err.Error(), "critical") || strings.Contains(err.Error(), "fatal")
		}
		hasCriticalError := ExistsError[string](isCritical)

		results := []Result[string]{
			Left[string](errors.New("critical: system failure")),
			Of("success"),
			Left[string](errors.New("warning: low memory")),
			Left[string](errors.New("fatal: database connection lost")),
		}

		// Act
		var criticalErrors []Result[string]
		for _, r := range results {
			if hasCriticalError(r) {
				criticalErrors = append(criticalErrors, r)
			}
		}

		// Assert
		assert.Len(t, criticalErrors, 2, "should identify critical errors")
	})
}

func BenchmarkExists(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := Of(42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsError(b *testing.B) {
	isValidationError := func(err error) bool {
		return strings.Contains(err.Error(), "validation:")
	}
	hasValidationError := ExistsError[int](isValidationError)
	input := Left[int](errors.New("validation: error"))

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func BenchmarkExistsPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := Of(-42)

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsErrorPredicateFails(b *testing.B) {
	isValidationError := func(err error) bool {
		return strings.Contains(err.Error(), "validation:")
	}
	hasValidationError := ExistsError[int](isValidationError)
	input := Left[int](errors.New("network: error"))

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func BenchmarkExistsOnError(b *testing.B) {
	isPositive := N.MoreThan(0)
	hasPositive := Exists(isPositive)
	input := Left[int](errors.New("error"))

	b.ResetTimer()
	for range b.N {
		_ = hasPositive(input)
	}
}

func BenchmarkExistsErrorOnOk(b *testing.B) {
	isValidationError := func(err error) bool {
		return strings.Contains(err.Error(), "validation:")
	}
	hasValidationError := ExistsError[int](isValidationError)
	input := Of(42)

	b.ResetTimer()
	for range b.N {
		_ = hasValidationError(input)
	}
}

func TestForAll_Success(t *testing.T) {
	t.Run("Ok value that passes predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Of(5)

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Ok value that passes predicate")
	})

	t.Run("Ok value at boundary that passes predicate", func(t *testing.T) {
		// Arrange
		isNonNegative := func(n int) bool { return n >= 0 }
		allNonNegative := ForAll(isNonNegative)
		input := Of(0)

		// Act
		result := allNonNegative(input)

		// Assert
		assert.True(t, result, "should return true for Ok value at boundary that passes predicate")
	})

	t.Run("Error value with predicate that would pass", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Left[int](errors.New("error"))

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Error value (vacuous truth)")
	})

	t.Run("Error value with predicate that would fail", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Left[int](errors.New("error"))

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for Error value regardless of predicate (vacuous truth)")
	})

	t.Run("Ok value with string predicate", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 5 }
		allLongString := ForAll(isLongString)
		input := Of("hello world")

		// Act
		result := allLongString(input)

		// Assert
		assert.True(t, result, "should return true for Ok string that passes predicate")
	})

	t.Run("Ok value with complex predicate", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		allEvenPositive := ForAll(isEvenAndPositive)
		input := Of(4)

		// Act
		result := allEvenPositive(input)

		// Assert
		assert.True(t, result, "should return true for Ok value that passes complex predicate")
	})
}

func TestForAll_Failure(t *testing.T) {
	t.Run("Ok value that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Of(-3)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Ok value that fails predicate")
	})

	t.Run("Ok value at boundary that fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		input := Of(0)

		// Act
		result := allPositive(input)

		// Assert
		assert.False(t, result, "should return false for Ok value at boundary that fails predicate")
	})

	t.Run("Ok value with string predicate that fails", func(t *testing.T) {
		// Arrange
		isLongString := func(s string) bool { return len(s) > 10 }
		allLongString := ForAll(isLongString)
		input := Of("short")

		// Act
		result := allLongString(input)

		// Assert
		assert.False(t, result, "should return false for Ok string that fails predicate")
	})

	t.Run("Ok value with complex predicate that fails", func(t *testing.T) {
		// Arrange
		isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
		allEvenPositive := ForAll(isEvenAndPositive)
		input := Of(3)

		// Act
		result := allEvenPositive(input)

		// Assert
		assert.False(t, result, "should return false for Ok value that fails complex predicate")
	})
}

func TestForAll_EdgeCases(t *testing.T) {
	t.Run("Ok with zero value", func(t *testing.T) {
		// Arrange
		isZero := func(n int) bool { return n == 0 }
		allZero := ForAll(isZero)
		input := Of(0)

		// Act
		result := allZero(input)

		// Assert
		assert.True(t, result, "should handle zero value correctly")
	})

	t.Run("Ok with empty string", func(t *testing.T) {
		// Arrange
		isEmpty := func(s string) bool { return len(s) == 0 }
		allEmpty := ForAll(isEmpty)
		input := Of("")

		// Act
		result := allEmpty(input)

		// Assert
		assert.True(t, result, "should handle empty string correctly")
	})

	t.Run("Ok with nil slice", func(t *testing.T) {
		// Arrange
		isNil := func(s []int) bool { return s == nil }
		allNil := ForAll(isNil)
		input := Of([]int(nil))

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
		assert.True(t, allTrue(Of(42)), "should return true for Ok with always-true predicate")
		assert.True(t, allTrue(Left[int](errors.New("error"))), "should return true for Error with always-true predicate")
	})

	t.Run("predicate always returns false", func(t *testing.T) {
		// Arrange
		alwaysFalse := func(int) bool { return false }
		allFalse := ForAll(alwaysFalse)

		// Act & Assert
		assert.False(t, allFalse(Of(42)), "should return false for Ok with always-false predicate")
		assert.True(t, allFalse(Left[int](errors.New("error"))), "should return true for Error even with always-false predicate (vacuous truth)")
	})

	t.Run("Error with various error types", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)

		// Act & Assert
		assert.True(t, allPositive(Left[int](errors.New(""))), "should return true for Error with empty message")
		assert.True(t, allPositive(Left[int](errors.New("error"))), "should return true for Error with message")
		assert.True(t, allPositive(Left[int](errors.New("validation: failed"))), "should return true for Error with any message")
	})
}

func TestForAll_Integration(t *testing.T) {
	t.Run("validate all successful results meet criteria", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		values := []Result[int]{
			Of(5),
			Left[int](errors.New("error1")),
			Of(10),
			Left[int](errors.New("error2")),
			Of(3),
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
		assert.True(t, allValid, "should return true when all Ok values pass predicate (Error values ignored)")
	})

	t.Run("detect when any Ok value fails predicate", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		values := []Result[int]{
			Of(5),
			Left[int](errors.New("error1")),
			Of(-3), // This fails
			Of(10),
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
		assert.False(t, allValid, "should return false when any Ok value fails predicate")
	})

	t.Run("all Error values pass vacuously", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		values := []Result[int]{
			Left[int](errors.New("error1")),
			Left[int](errors.New("error2")),
			Left[int](errors.New("error3")),
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
		assert.True(t, allValid, "should return true for all Error values (vacuous truth)")
	})

	t.Run("combine with other predicates", func(t *testing.T) {
		// Arrange
		isEven := func(n int) bool { return n%2 == 0 }
		isPositive := N.MoreThan(0)
		allEven := ForAll(isEven)
		allPositive := ForAll(isPositive)

		input1 := Of(4)
		input2 := Of(3)
		input3 := Of(-4)
		input4 := Left[int](errors.New("error"))

		// Act & Assert
		assert.True(t, allEven(input1) && allPositive(input1), "should pass both predicates")
		assert.False(t, allEven(input2) && allPositive(input2), "should fail even predicate")
		assert.False(t, allEven(input3) && allPositive(input3), "should fail positive predicate")
		assert.True(t, allEven(input4) && allPositive(input4), "Error should pass both predicates (vacuous truth)")
	})

	t.Run("contrast with Exists", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)
		hasPositive := Exists(isPositive)

		errorValue := Left[int](errors.New("error"))
		okPositive := Of(5)
		okNegative := Of(-3)

		// Act & Assert - ForAll behavior
		assert.True(t, allPositive(errorValue), "ForAll: Error is vacuously true")
		assert.True(t, allPositive(okPositive), "ForAll: Ok with passing predicate is true")
		assert.False(t, allPositive(okNegative), "ForAll: Ok with failing predicate is false")

		// Act & Assert - Exists behavior (contrast)
		assert.False(t, hasPositive(errorValue), "Exists: Error is false")
		assert.True(t, hasPositive(okPositive), "Exists: Ok with passing predicate is true")
		assert.False(t, hasPositive(okNegative), "Exists: Ok with failing predicate is false")
	})

	t.Run("De Morgan's law: ForAll(p) ≡ not(Exists(not(p)))", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		isNotPositive := func(n int) bool { return !isPositive(n) }
		allPositive := ForAll(isPositive)
		hasNotPositive := Exists(isNotPositive)

		testCases := []Result[int]{
			Of(5),
			Of(-3),
			Left[int](errors.New("error")),
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
		allInFirstQuadrant := ForAll(isInFirstQuadrant)

		// Act & Assert
		assert.True(t, allInFirstQuadrant(Of(Point{1, 1})), "point in first quadrant passes")
		assert.False(t, allInFirstQuadrant(Of(Point{-1, 1})), "point outside first quadrant fails")
		assert.True(t, allInFirstQuadrant(Left[Point](errors.New("error"))), "Error passes vacuously")
	})

	t.Run("with slice type", func(t *testing.T) {
		// Arrange
		hasElements := func(s []int) bool { return len(s) > 0 }
		allNonEmpty := ForAll(hasElements)

		// Act & Assert
		assert.True(t, allNonEmpty(Of([]int{1, 2, 3})), "non-empty slice passes")
		assert.False(t, allNonEmpty(Of([]int{})), "empty slice fails")
		assert.True(t, allNonEmpty(Left[[]int](errors.New("error"))), "Error passes vacuously")
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
		assert.True(t, allHaveAgeKey(Of(map[string]int{"age": 25})), "map with key passes")
		assert.False(t, allHaveAgeKey(Of(map[string]int{"name": 1})), "map without key fails")
		assert.True(t, allHaveAgeKey(Left[map[string]int](errors.New("error"))), "Error passes vacuously")
	})
}

func TestForAll_ErrorHandling(t *testing.T) {
	t.Run("with various error messages", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)

		testCases := []error{
			errors.New("validation: must be positive"),
			errors.New("network: connection failed"),
			errors.New(""),
		}

		// Act & Assert
		for _, err := range testCases {
			input := Left[int](err)
			result := allPositive(input)
			assert.True(t, result, "should return true for error regardless of message (vacuous truth)")
		}
	})

	t.Run("with wrapped errors", func(t *testing.T) {
		// Arrange
		isPositive := N.MoreThan(0)
		allPositive := ForAll(isPositive)

		baseErr := errors.New("base error")
		wrappedErr := errors.New("wrapped: " + baseErr.Error())
		input := Left[int](wrappedErr)

		// Act
		result := allPositive(input)

		// Assert
		assert.True(t, result, "should return true for wrapped error (vacuous truth)")
	})
}

func BenchmarkForAll(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := Of(42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllPredicateFails(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := Of(-42)

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllOnError(b *testing.B) {
	isPositive := N.MoreThan(0)
	allPositive := ForAll(isPositive)
	input := Left[int](errors.New("error"))

	b.ResetTimer()
	for range b.N {
		_ = allPositive(input)
	}
}

func BenchmarkForAllComplexPredicate(b *testing.B) {
	isEvenAndPositive := func(n int) bool { return n > 0 && n%2 == 0 }
	allEvenPositive := ForAll(isEvenAndPositive)
	input := Of(42)

	b.ResetTimer()
	for range b.N {
		_ = allEvenPositive(input)
	}
}
