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

package validate

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/stretchr/testify/assert"
)

// TestValidateType tests the Validate type structure
func TestValidateType(t *testing.T) {
	t.Run("basic validate function", func(t *testing.T) {
		// Create a simple validator that checks if a number is positive
		validatePositive := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				if n > 0 {
					return validation.Success(n)
				}
				return validation.FailureWithMessage[int](n, "must be positive")(ctx)
			}
		}

		// Test with positive number
		result := validatePositive(42)(nil)
		assert.Equal(t, validation.Of(42), result)

		// Test with negative number
		result = validatePositive(-5)(nil)
		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "must be positive", errors[0].Messsage)
	})

	t.Run("validate with context", func(t *testing.T) {
		validateWithContext := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				if s == "" {
					return validation.FailureWithMessage[string](s, "empty string")(ctx)
				}
				return validation.Success(s)
			}
		}

		ctx := validation.Context{
			{Key: "username", Type: "string"},
		}

		result := validateWithContext("")(ctx)
		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, ctx, errors[0].Context)
	})
}

// TestValidateComposition tests composing validators
func TestValidateComposition(t *testing.T) {
	t.Run("sequential validation", func(t *testing.T) {
		// First validator: check if string is not empty
		validateNotEmpty := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				if s == "" {
					return validation.FailureWithMessage[string](s, "must not be empty")(ctx)
				}
				return validation.Success(s)
			}
		}

		// Second validator: check if string has minimum length
		validateMinLength := func(minLen int) func(string) Reader[validation.Context, validation.Validation[string]] {
			return func(s string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					if len(s) < minLen {
						return validation.FailureWithMessage[string](s, "too short")(ctx)
					}
					return validation.Success(s)
				}
			}
		}

		// Test with valid input
		input := "hello"
		result1 := validateNotEmpty(input)(nil)
		assert.Equal(t, validation.Of("hello"), result1)

		result2 := validateMinLength(3)(input)(nil)
		assert.Equal(t, validation.Of("hello"), result2)

		// Test with invalid input
		shortInput := "hi"
		result3 := validateMinLength(5)(shortInput)(nil)
		assert.True(t, E.IsLeft(result3))
	})
}

// TestValidateWithDifferentTypes tests validators with various input/output types
func TestValidateWithDifferentTypes(t *testing.T) {
	t.Run("string to int conversion", func(t *testing.T) {
		// Validator that parses string to int
		validateParseInt := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				// Simple parsing logic for testing
				if s == "42" {
					return validation.Success(42)
				}
				return validation.FailureWithMessage[int](s, "invalid integer")(ctx)
			}
		}

		result := validateParseInt("42")(nil)
		assert.Equal(t, validation.Of(42), result)

		result = validateParseInt("abc")(nil)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("struct validation", func(t *testing.T) {
		type User struct {
			Name  string
			Age   int
			Email string
		}

		validateUser := func(u User) Reader[validation.Context, validation.Validation[User]] {
			return func(ctx validation.Context) validation.Validation[User] {
				if u.Name == "" {
					return validation.FailureWithMessage[User](u, "name is required")(ctx)
				}
				if u.Age < 0 {
					return validation.FailureWithMessage[User](u, "age must be non-negative")(ctx)
				}
				if u.Email == "" {
					return validation.FailureWithMessage[User](u, "email is required")(ctx)
				}
				return validation.Success(u)
			}
		}

		validUser := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
		result := validateUser(validUser)(nil)
		assert.Equal(t, validation.Of(validUser), result)

		invalidUser := User{Name: "", Age: 30, Email: "alice@example.com"}
		result = validateUser(invalidUser)(nil)
		assert.True(t, E.IsLeft(result))
	})
}

// TestValidateContextTracking tests context tracking through nested structures
func TestValidateContextTracking(t *testing.T) {
	t.Run("nested context", func(t *testing.T) {
		validateField := func(value string, fieldName string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				// Add field to context
				newCtx := append(ctx, validation.ContextEntry{
					Key:  fieldName,
					Type: "string",
				})

				if value == "" {
					return validation.FailureWithMessage[string](value, "field is empty")(newCtx)
				}
				return validation.Success(value)
			}
		}

		baseCtx := validation.Context{
			{Key: "user", Type: "User"},
		}

		result := validateField("", "email")(baseCtx)
		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)

		// Check that context includes both user and email
		assert.Len(t, errors[0].Context, 2)
		assert.Equal(t, "user", errors[0].Context[0].Key)
		assert.Equal(t, "email", errors[0].Context[1].Key)
	})
}

// TestValidateErrorMessages tests error message generation
func TestValidateErrorMessages(t *testing.T) {
	t.Run("custom error messages", func(t *testing.T) {
		validateRange := func(min, max int) func(int) Reader[validation.Context, validation.Validation[int]] {
			return func(n int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					if n < min {
						return validation.FailureWithMessage[int](n, "value too small")(ctx)
					}
					if n > max {
						return validation.FailureWithMessage[int](n, "value too large")(ctx)
					}
					return validation.Success(n)
				}
			}
		}

		result := validateRange(0, 100)(150)(nil)
		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Equal(t, "value too large", errors[0].Messsage)

		result = validateRange(0, 100)(-10)(nil)
		assert.True(t, E.IsLeft(result))
		_, errors = E.Unwrap(result)
		assert.Equal(t, "value too small", errors[0].Messsage)
	})
}

// TestValidateTransformations tests validators that transform values
func TestValidateTransformations(t *testing.T) {
	t.Run("normalize and validate", func(t *testing.T) {
		// Validator that normalizes (trims) and validates
		validateAndNormalize := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				// Simple trim simulation - trim all leading and trailing spaces
				normalized := s
				// Trim leading spaces
				for len(normalized) > 0 && normalized[0] == ' ' {
					normalized = normalized[1:]
				}
				// Trim trailing spaces
				for len(normalized) > 0 && normalized[len(normalized)-1] == ' ' {
					normalized = normalized[:len(normalized)-1]
				}

				if normalized == "" {
					return validation.FailureWithMessage[string](s, "empty after normalization")(ctx)
				}
				return validation.Success(normalized)
			}
		}

		result := validateAndNormalize(" hello ")(nil)
		assert.Equal(t, validation.Of("hello"), result)

		result = validateAndNormalize("   ")(nil)
		assert.True(t, E.IsLeft(result))
	})
}

// TestValidateChaining tests chaining multiple validators
func TestValidateChaining(t *testing.T) {
	t.Run("chain validators manually", func(t *testing.T) {
		// First validator
		v1 := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				if n < 0 {
					return validation.FailureWithMessage[int](n, "must be non-negative")(ctx)
				}
				return validation.Success(n)
			}
		}

		// Second validator (depends on first)
		v2 := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				if n > 100 {
					return validation.FailureWithMessage[int](n, "must be <= 100")(ctx)
				}
				return validation.Success(n)
			}
		}

		// Test valid value
		input := 50
		result1 := v1(input)(nil)
		assert.Equal(t, validation.Of(50), result1)

		result2 := v2(input)(nil)
		assert.Equal(t, validation.Of(50), result2)

		// Test invalid value (too large)
		input = 150
		result1 = v1(input)(nil)
		assert.Equal(t, validation.Of(150), result1)

		result2 = v2(input)(nil)
		assert.True(t, E.IsLeft(result2))
	})
}

// TestValidateComplexScenarios tests real-world validation scenarios
func TestValidateComplexScenarios(t *testing.T) {
	t.Run("email validation", func(t *testing.T) {
		validateEmail := func(email string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				// Simple email validation for testing
				hasAt := false
				hasDot := false
				for _, c := range email {
					if c == '@' {
						hasAt = true
					}
					if c == '.' {
						hasDot = true
					}
				}

				if !hasAt || !hasDot {
					return validation.FailureWithMessage[string](email, "invalid email format")(ctx)
				}
				return validation.Success(email)
			}
		}

		result := validateEmail("user@example.com")(nil)
		assert.Equal(t, validation.Of("user@example.com"), result)

		result = validateEmail("invalid-email")(nil)
		assert.True(t, E.IsLeft(result))

		result = validateEmail("no-domain@")(nil)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("password strength validation", func(t *testing.T) {
		validatePassword := func(pwd string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				if len(pwd) < 8 {
					return validation.FailureWithMessage[string](pwd, "password too short")(ctx)
				}

				hasUpper := false
				hasLower := false
				hasDigit := false

				for _, c := range pwd {
					if c >= 'A' && c <= 'Z' {
						hasUpper = true
					}
					if c >= 'a' && c <= 'z' {
						hasLower = true
					}
					if c >= '0' && c <= '9' {
						hasDigit = true
					}
				}

				if !hasUpper || !hasLower || !hasDigit {
					return validation.FailureWithMessage[string](pwd, "password must contain upper, lower, and digit")(ctx)
				}

				return validation.Success(pwd)
			}
		}

		result := validatePassword("StrongPass123")(nil)
		assert.Equal(t, validation.Of("StrongPass123"), result)

		result = validatePassword("weak")(nil)
		assert.True(t, E.IsLeft(result))

		result = validatePassword("nouppercase123")(nil)
		assert.True(t, E.IsLeft(result))
	})
}

// Benchmark tests
func BenchmarkValidate_Success(b *testing.B) {
	validate := func(n int) Reader[validation.Context, validation.Validation[int]] {
		return func(ctx validation.Context) validation.Validation[int] {
			if n > 0 {
				return validation.Success(n)
			}
			return validation.FailureWithMessage[int](n, "must be positive")(ctx)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate(42)(nil)
	}
}

func BenchmarkValidate_Failure(b *testing.B) {
	validate := func(n int) Reader[validation.Context, validation.Validation[int]] {
		return func(ctx validation.Context) validation.Validation[int] {
			if n > 0 {
				return validation.Success(n)
			}
			return validation.FailureWithMessage[int](n, "must be positive")(ctx)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate(-1)(nil)
	}
}

func BenchmarkValidate_WithContext(b *testing.B) {
	validate := func(s string) Reader[validation.Context, validation.Validation[string]] {
		return func(ctx validation.Context) validation.Validation[string] {
			if s == "" {
				return validation.FailureWithMessage[string](s, "empty string")(ctx)
			}
			return validation.Success(s)
		}
	}

	ctx := validation.Context{
		{Key: "field1", Type: "string"},
		{Key: "field2", Type: "string"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validate("test")(ctx)
	}
}

// TestOf tests the Of function
func TestOf(t *testing.T) {
	t.Run("creates successful validation with value", func(t *testing.T) {
		validator := Of[string](42)
		result := validator("any input")(nil)

		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("ignores input value", func(t *testing.T) {
		validator := Of[string]("success")

		result1 := validator("input1")(nil)
		result2 := validator("input2")(nil)
		result3 := validator("")(nil)

		assert.Equal(t, validation.Of("success"), result1)
		assert.Equal(t, validation.Of("success"), result2)
		assert.Equal(t, validation.Of("success"), result3)
	})

	t.Run("works with different types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		user := User{Name: "Alice", Age: 30}
		validator := Of[int](user)
		result := validator(123)(nil)

		assert.Equal(t, validation.Of(user), result)
	})
}

// TestMonadMap tests the MonadMap function
func TestMonadMap(t *testing.T) {
	t.Run("transforms successful validation", func(t *testing.T) {
		validator := Of[string](21)
		doubled := MonadMap(validator, N.Mul(2))

		result := doubled("input")(nil)
		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("preserves validation errors", func(t *testing.T) {
		failingValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "validation failed")(ctx)
			}
		}

		mapped := MonadMap(failingValidator, N.Mul(2))
		result := mapped("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "validation failed", errors[0].Messsage)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		validator := Of[string](10)
		transformed := MonadMap(
			MonadMap(
				MonadMap(validator, N.Add(5)),
				N.Mul(2),
			),
			N.Sub(10),
		)

		result := transformed("input")(nil)
		assert.Equal(t, validation.Of(20), result) // (10 + 5) * 2 - 10 = 20
	})

	t.Run("transforms between different types", func(t *testing.T) {
		validator := Of[string](42)
		toString := MonadMap(validator, func(x int) string {
			return "value: " + string(rune(x+'0'))
		})

		result := toString("input")(nil)
		assert.True(t, E.IsRight(result))
		if E.IsRight(result) {
			value, _ := E.Unwrap(result)
			assert.Contains(t, value, "value:")
		}
	})
}

// TestMap tests the Map function
func TestMap(t *testing.T) {
	t.Run("creates reusable transformation", func(t *testing.T) {
		double := Map[string](N.Mul(2))

		validator1 := Of[string](21)
		validator2 := Of[string](10)

		result1 := double(validator1)("input")(nil)
		result2 := double(validator2)("input")(nil)

		assert.Equal(t, validation.Of(42), result1)
		assert.Equal(t, validation.Of(20), result2)
	})

	t.Run("preserves errors in transformation", func(t *testing.T) {
		increment := Map[string](func(x int) int { return x + 1 })

		failingValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "error")(ctx)
			}
		}

		result := increment(failingValidator)("input")(nil)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("composes with other operators", func(t *testing.T) {
		addFive := Map[string](N.Add(5))
		double := Map[string](N.Mul(2))

		validator := Of[string](10)
		composed := double(addFive(validator))

		result := composed("input")(nil)
		assert.Equal(t, validation.Of(30), result) // (10 + 5) * 2 = 30
	})
}

// TestChain tests the Chain function
func TestChain(t *testing.T) {
	t.Run("sequences dependent validations", func(t *testing.T) {
		// First validator: parse string to int
		parseValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				if s == "42" {
					return validation.Success(42)
				}
				return validation.FailureWithMessage[int](s, "invalid number")(ctx)
			}
		}

		// Second validator: check if number is positive
		checkPositive := func(n int) Validate[string, string] {
			return func(input string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					if n > 0 {
						return validation.Success("positive")
					}
					return validation.FailureWithMessage[string](n, "not positive")(ctx)
				}
			}
		}

		chained := Chain(checkPositive)(parseValidator)
		result := chained("42")(nil)

		assert.Equal(t, validation.Of("positive"), result)
	})

	t.Run("stops on first validation failure", func(t *testing.T) {
		failingValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "first failed")(ctx)
			}
		}

		neverCalled := func(n int) Validate[string, string] {
			return func(input string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					// This should never be reached
					t.Error("Second validator should not be called")
					return validation.Success("should not reach")
				}
			}
		}

		chained := Chain(neverCalled)(failingValidator)
		result := chained("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Equal(t, "first failed", errors[0].Messsage)
	})

	t.Run("propagates second validation failure", func(t *testing.T) {
		successValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.Success(42)
			}
		}

		failingSecond := func(n int) Validate[string, string] {
			return func(input string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					return validation.FailureWithMessage[string](n, "second failed")(ctx)
				}
			}
		}

		chained := Chain(failingSecond)(successValidator)
		result := chained("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Equal(t, "second failed", errors[0].Messsage)
	})
}

// TestMonadAp tests the MonadAp function
func TestMonadAp(t *testing.T) {
	t.Run("applies function to value when both succeed", func(t *testing.T) {
		funcValidator := Of[string](N.Mul(2))
		valueValidator := Of[string](21)

		result := MonadAp(funcValidator, valueValidator)("input")(nil)

		assert.Equal(t, validation.Of(42), result)
	})

	t.Run("accumulates errors when function validator fails", func(t *testing.T) {
		failingFunc := func(s string) Reader[validation.Context, validation.Validation[func(int) int]] {
			return func(ctx validation.Context) validation.Validation[func(int) int] {
				return validation.FailureWithMessage[func(int) int](s, "func failed")(ctx)
			}
		}
		valueValidator := Of[string](21)

		result := MonadAp(failingFunc, valueValidator)("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "func failed", errors[0].Messsage)
	})

	t.Run("accumulates errors when value validator fails", func(t *testing.T) {
		funcValidator := Of[string](N.Mul(2))
		failingValue := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "value failed")(ctx)
			}
		}

		result := MonadAp(funcValidator, failingValue)("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 1)
		assert.Equal(t, "value failed", errors[0].Messsage)
	})

	t.Run("returns error when both validators fail", func(t *testing.T) {
		failingFunc := func(s string) Reader[validation.Context, validation.Validation[func(int) int]] {
			return func(ctx validation.Context) validation.Validation[func(int) int] {
				return validation.FailureWithMessage[func(int) int](s, "func failed")(ctx)
			}
		}
		failingValue := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "value failed")(ctx)
			}
		}

		result := MonadAp(failingFunc, failingValue)("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		// Note: The current implementation returns the first error encountered
		assert.GreaterOrEqual(t, len(errors), 1)
		// At least one of the errors should be present
		hasError := false
		for _, err := range errors {
			if err.Messsage == "func failed" || err.Messsage == "value failed" {
				hasError = true
				break
			}
		}
		assert.True(t, hasError, "Should contain at least one validation error")
	})
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	t.Run("creates reusable applicative operator", func(t *testing.T) {
		valueValidator := Of[string](21)
		applyTo21 := Ap[int](valueValidator)

		double := Of[string](N.Mul(2))
		triple := Of[string](func(x int) int { return x * 3 })

		result1 := applyTo21(double)("input")(nil)
		result2 := applyTo21(triple)("input")(nil)

		assert.Equal(t, validation.Of(42), result1)
		assert.Equal(t, validation.Of(63), result2)
	})

	t.Run("preserves errors from value validator", func(t *testing.T) {
		failingValue := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "value error")(ctx)
			}
		}

		applyToFailing := Ap[int](failingValue)
		funcValidator := Of[string](N.Mul(2))

		result := applyToFailing(funcValidator)("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Equal(t, "value error", errors[0].Messsage)
	})

	t.Run("preserves errors from function validator", func(t *testing.T) {
		valueValidator := Of[string](21)
		applyTo21 := Ap[int](valueValidator)

		failingFunc := func(s string) Reader[validation.Context, validation.Validation[func(int) int]] {
			return func(ctx validation.Context) validation.Validation[func(int) int] {
				return validation.FailureWithMessage[func(int) int](s, "func error")(ctx)
			}
		}

		result := applyTo21(failingFunc)("input")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Equal(t, "func error", errors[0].Messsage)
	})
}

// TestMonadLaws tests that the monad laws hold for Validate
func TestMonadLaws(t *testing.T) {
	t.Run("left identity: Of(a) >>= f === f(a)", func(t *testing.T) {
		a := 42
		f := func(x int) Validate[string, string] {
			return Of[string]("value: " + string(rune(x+'0')))
		}

		// Of(a) >>= f
		left := Chain(f)(Of[string](a))
		// f(a)
		right := f(a)

		leftResult := left("input")(nil)
		rightResult := right("input")(nil)

		assert.Equal(t, E.IsRight(leftResult), E.IsRight(rightResult))
		if E.IsRight(leftResult) {
			leftVal, _ := E.Unwrap(leftResult)
			rightVal, _ := E.Unwrap(rightResult)
			assert.Equal(t, leftVal, rightVal)
		}
	})

	t.Run("right identity: m >>= Of === m", func(t *testing.T) {
		m := Of[string](42)

		// m >>= Of
		chained := Chain(func(x int) Validate[string, int] {
			return Of[string](x)
		})(m)

		mResult := m("input")(nil)
		chainedResult := chained("input")(nil)

		assert.Equal(t, E.IsRight(mResult), E.IsRight(chainedResult))
		if E.IsRight(mResult) {
			mVal, _ := E.Unwrap(mResult)
			chainedVal, _ := E.Unwrap(chainedResult)
			assert.Equal(t, mVal, chainedVal)
		}
	})
}

// TestFunctorLaws tests that the functor laws hold for Validate
func TestFunctorLaws(t *testing.T) {
	t.Run("identity: map(id) === id", func(t *testing.T) {
		validator := Of[string](42)
		identity := func(x int) int { return x }

		mapped := MonadMap(validator, identity)

		origResult := validator("input")(nil)
		mappedResult := mapped("input")(nil)

		assert.Equal(t, E.IsRight(origResult), E.IsRight(mappedResult))
		if E.IsRight(origResult) {
			origVal, _ := E.Unwrap(origResult)
			mappedVal, _ := E.Unwrap(mappedResult)
			assert.Equal(t, origVal, mappedVal)
		}
	})

	t.Run("composition: map(f . g) === map(f) . map(g)", func(t *testing.T) {
		validator := Of[string](10)
		f := N.Mul(2)
		g := N.Add(5)

		// map(f . g)
		composed := MonadMap(validator, func(x int) int { return f(g(x)) })

		// map(f) . map(g)
		separate := MonadMap(MonadMap(validator, g), f)

		composedResult := composed("input")(nil)
		separateResult := separate("input")(nil)

		assert.Equal(t, E.IsRight(composedResult), E.IsRight(separateResult))
		if E.IsRight(composedResult) {
			composedVal, _ := E.Unwrap(composedResult)
			separateVal, _ := E.Unwrap(separateResult)
			assert.Equal(t, composedVal, separateVal)
		}
	})
}

// TestChainLeft tests the ChainLeft function
func TestChainLeft(t *testing.T) {
	t.Run("transforms failures while preserving successes", func(t *testing.T) {
		// Create a failing validator
		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "validation failed")(ctx)
			}
		}

		// Handler that recovers from specific errors
		handler := ChainLeft(func(errs Errors) Validate[int, int] {
			for _, err := range errs {
				if err.Messsage == "validation failed" {
					return Of[int](0) // recover with default
				}
			}
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return E.Left[int](errs)
				}
			}
		})

		validator := handler(failingValidator)
		result := validator(-5)(nil)

		assert.Equal(t, validation.Of(0), result, "Should recover from failure")
	})

	t.Run("preserves success values unchanged", func(t *testing.T) {
		successValidator := Of[int](42)

		handler := ChainLeft(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return validation.FailureWithMessage[int](input, "should not be called")(ctx)
				}
			}
		})

		validator := handler(successValidator)
		result := validator(100)(nil)

		assert.Equal(t, validation.Of(42), result, "Success should pass through unchanged")
	})

	t.Run("aggregates errors when transformation also fails", func(t *testing.T) {
		failingValidator := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				return validation.FailureWithMessage[string](s, "original error")(ctx)
			}
		}

		handler := ChainLeft(func(errs Errors) Validate[string, string] {
			return func(input string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					return validation.FailureWithMessage[string](input, "additional error")(ctx)
				}
			}
		})

		validator := handler(failingValidator)
		result := validator("test")(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 2, "Should aggregate both errors")

		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "original error")
		assert.Contains(t, messages, "additional error")
	})

	t.Run("adds context to errors", func(t *testing.T) {
		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "invalid value")(ctx)
			}
		}

		addContext := ChainLeft(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return E.Left[int](validation.Errors{
						{
							Context:  validation.Context{{Key: "user", Type: "User"}, {Key: "age", Type: "int"}},
							Messsage: "failed to validate user age",
						},
					})
				}
			}
		})

		validator := addContext(failingValidator)
		result := validator(150)(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 2, "Should have both original and context errors")
	})

	t.Run("can be composed in pipeline", func(t *testing.T) {
		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "error1")(ctx)
			}
		}

		handler1 := ChainLeft(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return validation.FailureWithMessage[int](input, "error2")(ctx)
				}
			}
		})

		handler2 := ChainLeft(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return validation.FailureWithMessage[int](input, "error3")(ctx)
				}
			}
		})

		validator := handler2(handler1(failingValidator))
		result := validator(42)(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.GreaterOrEqual(t, len(errors), 2, "Should accumulate errors through pipeline")
	})

	t.Run("provides access to original input", func(t *testing.T) {
		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "failed")(ctx)
			}
		}

		// Handler uses input to determine recovery strategy
		handler := ChainLeft(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					// Use input value to decide on recovery
					if input < 0 {
						return validation.Of(0)
					}
					if input > 100 {
						return validation.Of(100)
					}
					return E.Left[int](errs)
				}
			}
		})

		validator := handler(failingValidator)

		result1 := validator(-10)(nil)
		assert.Equal(t, validation.Of(0), result1, "Should recover negative to 0")

		result2 := validator(150)(nil)
		assert.Equal(t, validation.Of(100), result2, "Should recover large to 100")
	})

	t.Run("works with different input and output types", func(t *testing.T) {
		// Validator that converts string to int
		parseValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "parse failed")(ctx)
			}
		}

		// Handler that provides default based on input string
		handler := ChainLeft(func(errs Errors) Validate[string, int] {
			return func(input string) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					if input == "default" {
						return validation.Of(42)
					}
					return E.Left[int](errs)
				}
			}
		})

		validator := handler(parseValidator)
		result := validator("default")(nil)

		assert.Equal(t, validation.Of(42), result)
	})
}

// TestOrElse tests the OrElse function
func TestOrElse(t *testing.T) {
	t.Run("provides fallback for failing validation", func(t *testing.T) {
		// Primary validator that fails
		primaryValidator := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				return validation.FailureWithMessage[string](s, "not found")(ctx)
			}
		}

		// Use OrElse to provide fallback
		withFallback := OrElse(func(errs Errors) Validate[string, string] {
			return Of[string]("default value")
		})

		validator := withFallback(primaryValidator)
		result := validator("missing")(nil)

		assert.Equal(t, validation.Of("default value"), result)
	})

	t.Run("preserves success values unchanged", func(t *testing.T) {
		successValidator := Of[string]("success")

		withFallback := OrElse(func(errs Errors) Validate[string, string] {
			return Of[string]("fallback")
		})

		validator := withFallback(successValidator)
		result := validator("input")(nil)

		assert.Equal(t, validation.Of("success"), result, "Should not use fallback for success")
	})

	t.Run("aggregates errors when fallback also fails", func(t *testing.T) {
		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "primary failed")(ctx)
			}
		}

		withFallback := OrElse(func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					return validation.FailureWithMessage[int](input, "fallback failed")(ctx)
				}
			}
		})

		validator := withFallback(failingValidator)
		result := validator(42)(nil)

		assert.True(t, E.IsLeft(result))
		_, errors := E.Unwrap(result)
		assert.Len(t, errors, 2, "Should aggregate both errors")

		messages := make([]string, len(errors))
		for i, err := range errors {
			messages[i] = err.Messsage
		}
		assert.Contains(t, messages, "primary failed")
		assert.Contains(t, messages, "fallback failed")
	})

	t.Run("supports multiple fallback strategies", func(t *testing.T) {
		failingValidator := func(s string) Reader[validation.Context, validation.Validation[string]] {
			return func(ctx validation.Context) validation.Validation[string] {
				return validation.FailureWithMessage[string](s, "not in database")(ctx)
			}
		}

		// First fallback: try cache
		tryCache := OrElse(func(errs Errors) Validate[string, string] {
			return func(input string) Reader[validation.Context, validation.Validation[string]] {
				return func(ctx validation.Context) validation.Validation[string] {
					if input == "cached" {
						return validation.Of("from cache")
					}
					return E.Left[string](errs)
				}
			}
		})

		// Second fallback: use default
		useDefault := OrElse(func(errs Errors) Validate[string, string] {
			return Of[string]("default")
		})

		// Compose fallbacks
		validator := useDefault(tryCache(failingValidator))

		// Test with cached value
		result1 := validator("cached")(nil)
		assert.Equal(t, validation.Of("from cache"), result1)

		// Test with non-cached value (should use default)
		result2 := validator("other")(nil)
		assert.Equal(t, validation.Of("default"), result2)
	})

	t.Run("provides input-dependent fallback", func(t *testing.T) {
		failingValidator := func(s string) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](s, "parse failed")(ctx)
			}
		}

		// Fallback with different defaults based on input
		smartFallback := OrElse(func(errs Errors) Validate[string, int] {
			return func(input string) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					// Provide context-aware defaults
					if input == "http" {
						return validation.Of(80)
					}
					if input == "https" {
						return validation.Of(443)
					}
					return validation.Of(8080)
				}
			}
		})

		validator := smartFallback(failingValidator)

		result1 := validator("http")(nil)
		assert.Equal(t, validation.Of(80), result1)

		result2 := validator("https")(nil)
		assert.Equal(t, validation.Of(443), result2)

		result3 := validator("other")(nil)
		assert.Equal(t, validation.Of(8080), result3)
	})

	t.Run("is equivalent to ChainLeft", func(t *testing.T) {
		// Create identical handlers
		handler := func(errs Errors) Validate[int, int] {
			return func(input int) Reader[validation.Context, validation.Validation[int]] {
				return func(ctx validation.Context) validation.Validation[int] {
					if input < 0 {
						return validation.Of(0)
					}
					return E.Left[int](errs)
				}
			}
		}

		failingValidator := func(n int) Reader[validation.Context, validation.Validation[int]] {
			return func(ctx validation.Context) validation.Validation[int] {
				return validation.FailureWithMessage[int](n, "failed")(ctx)
			}
		}

		// Apply with ChainLeft
		withChainLeft := ChainLeft(handler)(failingValidator)

		// Apply with OrElse
		withOrElse := OrElse(handler)(failingValidator)

		// Test with same inputs
		inputs := []int{-10, 0, 10, -5, 100}
		for _, input := range inputs {
			result1 := withChainLeft(input)(nil)
			result2 := withOrElse(input)(nil)

			// Results should be identical
			assert.Equal(t, E.IsLeft(result1), E.IsLeft(result2))
			if E.IsRight(result1) {
				val1, _ := E.Unwrap(result1)
				val2, _ := E.Unwrap(result2)
				assert.Equal(t, val1, val2, "OrElse and ChainLeft should produce identical results")
			}
		}
	})

	t.Run("works in complex validation pipeline", func(t *testing.T) {
		type Config struct {
			Port int
			Host string
		}

		// Validator that tries to parse config
		parseConfig := func(s string) Reader[validation.Context, validation.Validation[Config]] {
			return func(ctx validation.Context) validation.Validation[Config] {
				return validation.FailureWithMessage[Config](s, "invalid config")(ctx)
			}
		}

		// Fallback to environment variables
		tryEnv := OrElse(func(errs Errors) Validate[string, Config] {
			return func(input string) Reader[validation.Context, validation.Validation[Config]] {
				return func(ctx validation.Context) validation.Validation[Config] {
					// Simulate env var lookup
					if input == "from_env" {
						return validation.Of(Config{Port: 8080, Host: "localhost"})
					}
					return E.Left[Config](errs)
				}
			}
		})

		// Final fallback to defaults
		useDefaults := OrElse(func(errs Errors) Validate[string, Config] {
			return Of[string](Config{Port: 3000, Host: "0.0.0.0"})
		})

		// Build pipeline
		validator := useDefaults(tryEnv(parseConfig))

		// Test with env fallback
		result1 := validator("from_env")(nil)
		assert.True(t, E.IsRight(result1))
		if E.IsRight(result1) {
			cfg, _ := E.Unwrap(result1)
			assert.Equal(t, 8080, cfg.Port)
			assert.Equal(t, "localhost", cfg.Host)
		}

		// Test with default fallback
		result2 := validator("other")(nil)
		assert.True(t, E.IsRight(result2))
		if E.IsRight(result2) {
			cfg, _ := E.Unwrap(result2)
			assert.Equal(t, 3000, cfg.Port)
			assert.Equal(t, "0.0.0.0", cfg.Host)
		}
	})
}
