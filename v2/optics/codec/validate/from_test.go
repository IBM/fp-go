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
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/codec/validation"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFromReaderResult_Success tests that FromReaderResult correctly converts
// a successful ReaderResult into a successful Validate
func TestFromReaderResult_Success(t *testing.T) {
	t.Run("converts successful ReaderResult with integer", func(t *testing.T) {
		// Create a ReaderResult that always succeeds
		successRR := func(input int) result.Result[string] {
			return result.Of(fmt.Sprintf("value: %d", input))
		}

		// Convert to Validate
		validator := FromReaderResult[int, string](successRR)

		// Execute the validator
		validationResult := validator(42)(nil)

		// Verify success
		assert.Equal(t, validation.Success("value: 42"), validationResult)
	})

	t.Run("converts successful ReaderResult with string input", func(t *testing.T) {
		// Create a ReaderResult that parses a string to int
		parseIntRR := result.Eitherize1(strconv.Atoi)

		// Convert to Validate
		validator := FromReaderResult[string, int](parseIntRR)

		// Execute with valid input
		validationResult := validator("123")(nil)

		// Verify success
		assert.Equal(t, validation.Success(123), validationResult)
	})

	t.Run("converts successful ReaderResult with complex type", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		// Create a ReaderResult that creates a User
		createUserRR := func(input string) result.Result[User] {
			return result.Of(User{Name: input, Age: 25})
		}

		// Convert to Validate
		validator := FromReaderResult[string, User](createUserRR)

		// Execute the validator
		validationResult := validator("Alice")(nil)

		// Verify success
		assert.Equal(t, validation.Success(User{Name: "Alice", Age: 25}), validationResult)
	})

	t.Run("preserves success with empty context", func(t *testing.T) {
		successRR := func(input int) result.Result[int] {
			return result.Of(input * 2)
		}

		validator := FromReaderResult[int, int](successRR)
		validationResult := validator(21)(Context{})

		assert.Equal(t, validation.Success(42), validationResult)
	})

	t.Run("preserves success with non-empty context", func(t *testing.T) {
		successRR := func(input string) result.Result[string] {
			return result.Of(input + " processed")
		}

		validator := FromReaderResult[string, string](successRR)
		ctx := Context{
			{Key: "user", Type: "User"},
			{Key: "name", Type: "string"},
		}
		validationResult := validator("test")(ctx)

		assert.Equal(t, validation.Success("test processed"), validationResult)
	})
}

// TestFromReaderResult_Failure tests that FromReaderResult correctly converts
// a failed ReaderResult into a failed Validate with proper error information
func TestFromReaderResult_Failure(t *testing.T) {
	t.Run("converts failed ReaderResult to validation error", func(t *testing.T) {
		expectedErr := errors.New("parse error")

		// Create a ReaderResult that always fails
		failureRR := func(input string) result.Result[int] {
			return result.Left[int](expectedErr)
		}

		// Convert to Validate
		validator := FromReaderResult[string, int](failureRR)

		// Execute the validator
		validationResult := validator("invalid")(nil)

		// Verify failure
		assert.True(t, either.IsLeft(validationResult))
		errors := either.MonadFold(validationResult,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)

		require.Len(t, errors, 1)
		assert.Equal(t, "unable to decode", errors[0].Messsage)
		assert.Equal(t, "invalid", errors[0].Value)
		assert.Equal(t, expectedErr, errors[0].Cause)
	})

	t.Run("preserves original error as cause", func(t *testing.T) {
		originalErr := fmt.Errorf("original error: %w", errors.New("root cause"))

		failureRR := func(input int) result.Result[string] {
			return result.Left[string](originalErr)
		}

		validator := FromReaderResult[int, string](failureRR)
		validationResult := validator(42)(nil)

		assert.True(t, either.IsLeft(validationResult))
		errors := either.MonadFold(validationResult,
			F.Identity[Errors],
			func(string) Errors { return nil },
		)

		require.Len(t, errors, 1)
		assert.Equal(t, originalErr, errors[0].Cause)
		assert.ErrorIs(t, errors[0].Cause, originalErr)
	})

	t.Run("includes context in validation error", func(t *testing.T) {
		failureRR := func(input string) result.Result[int] {
			return result.Left[int](errors.New("conversion failed"))
		}

		validator := FromReaderResult[string, int](failureRR)
		ctx := Context{
			{Key: "user", Type: "User"},
			{Key: "age", Type: "int"},
		}
		validationResult := validator("abc")(ctx)

		assert.True(t, either.IsLeft(validationResult))
		errors := either.MonadFold(validationResult,
			F.Identity[Errors],
			func(int) Errors { return nil },
		)

		require.Len(t, errors, 1)
		assert.Equal(t, ctx, errors[0].Context)
		assert.Equal(t, "abc", errors[0].Value)
	})

	t.Run("handles different error types", func(t *testing.T) {
		testCases := []struct {
			name  string
			err   error
			input string
		}{
			{
				name:  "simple error",
				err:   errors.New("simple error"),
				input: "test1",
			},
			{
				name:  "formatted error",
				err:   fmt.Errorf("formatted error: %s", "details"),
				input: "test2",
			},
			{
				name:  "wrapped error",
				err:   fmt.Errorf("wrapped: %w", errors.New("inner")),
				input: "test3",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				failureRR := func(input string) result.Result[int] {
					return result.Left[int](tc.err)
				}

				validator := FromReaderResult[string, int](failureRR)
				validationResult := validator(tc.input)(nil)

				assert.True(t, either.IsLeft(validationResult))
				errors := either.MonadFold(validationResult,
					F.Identity[Errors],
					func(int) Errors { return nil },
				)

				require.Len(t, errors, 1)
				assert.Equal(t, tc.err, errors[0].Cause)
				assert.Equal(t, tc.input, errors[0].Value)
			})
		}
	})
}

// TestFromReaderResult_Integration tests FromReaderResult in combination with
// other validation operations
func TestFromReaderResult_Integration(t *testing.T) {
	t.Run("chains with other validators", func(t *testing.T) {
		// Parse string to int
		parseIntRR := result.Eitherize1(strconv.Atoi)

		// Validate positive
		validatePositive := func(n int) Validate[string, int] {
			return func(input string) Reader[Context, Validation[int]] {
				return func(ctx Context) Validation[int] {
					if n > 0 {
						return validation.Success(n)
					}
					return validation.FailureWithMessage[int](n, "must be positive")(ctx)
				}
			}
		}

		// Combine validators
		validator := F.Pipe1(
			FromReaderResult[string, int](parseIntRR),
			Chain(validatePositive),
		)

		// Test with valid positive number
		result1 := validator("42")(nil)
		assert.True(t, either.IsRight(result1))

		// Test with valid negative number (should fail positive check)
		result2 := validator("-5")(nil)
		assert.True(t, either.IsLeft(result2))

		// Test with invalid string (should fail parsing)
		result3 := validator("abc")(nil)
		assert.True(t, either.IsLeft(result3))
	})

	t.Run("maps successful result", func(t *testing.T) {
		parseIntRR := result.Eitherize1(strconv.Atoi)

		// Convert and map to double the value
		validator := F.Pipe1(
			FromReaderResult[string, int](parseIntRR),
			Map[string, int, int](func(n int) int { return n * 2 }),
		)

		validationResult := validator("21")(nil)
		assert.Equal(t, validation.Success(42), validationResult)
	})

	t.Run("composes with Do and Bind", func(t *testing.T) {
		type State struct {
			parsed int
			valid  bool
		}

		parseIntRR := result.Eitherize1(strconv.Atoi)

		validator := F.Pipe2(
			Do[string](State{}),
			Bind(func(p int) func(State) State {
				return func(s State) State { s.parsed = p; return s }
			}, func(s State) Validate[string, int] {
				return FromReaderResult[string, int](parseIntRR)
			}),
			Let[string](func(v bool) func(State) State {
				return func(s State) State { s.valid = v; return s }
			}, func(s State) bool {
				return s.parsed > 0
			}),
		)

		result := validator("42")(nil)
		assert.Equal(t, validation.Success(State{parsed: 42, valid: true}), result)
	})
}

// TestFromReaderResult_EdgeCases tests edge cases and boundary conditions
func TestFromReaderResult_EdgeCases(t *testing.T) {
	t.Run("handles nil context", func(t *testing.T) {
		successRR := func(input int) result.Result[int] {
			return result.Of(input)
		}

		validator := FromReaderResult[int, int](successRR)
		validationResult := validator(42)(nil)

		assert.True(t, either.IsRight(validationResult))
	})

	t.Run("handles empty input", func(t *testing.T) {
		identityRR := func(input string) result.Result[string] {
			return result.Of(input)
		}

		validator := FromReaderResult[string, string](identityRR)
		validationResult := validator("")(nil)

		assert.Equal(t, validation.Success(""), validationResult)
	})

	t.Run("handles zero values", func(t *testing.T) {
		identityRR := func(input int) result.Result[int] {
			return result.Of(input)
		}

		validator := FromReaderResult[int, int](identityRR)
		validationResult := validator(0)(nil)

		assert.Equal(t, validation.Success(0), validationResult)
	})

	t.Run("handles pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		createDataRR := func(input int) result.Result[*Data] {
			return result.Of(&Data{Value: input})
		}

		validator := FromReaderResult[int, *Data](createDataRR)
		validationResult := validator(42)(nil)

		assert.True(t, either.IsRight(validationResult))
		data := either.MonadFold(validationResult,
			func(Errors) *Data { return nil },
			F.Identity[*Data],
		)
		require.NotNil(t, data)
		assert.Equal(t, 42, data.Value)
	})

	t.Run("handles slice types", func(t *testing.T) {
		splitRR := func(input string) result.Result[[]string] {
			if input == "" {
				return result.Left[[]string](errors.New("empty input"))
			}
			return result.Of([]string{input, input})
		}

		validator := FromReaderResult[string, []string](splitRR)
		validationResult := validator("test")(nil)

		assert.Equal(t, validation.Success([]string{"test", "test"}), validationResult)
	})

	t.Run("handles map types", func(t *testing.T) {
		createMapRR := func(input string) result.Result[map[string]int] {
			return result.Of(map[string]int{input: len(input)})
		}

		validator := FromReaderResult[string, map[string]int](createMapRR)
		validationResult := validator("hello")(nil)

		assert.Equal(t, validation.Success(map[string]int{"hello": 5}), validationResult)
	})
}

// TestFromReaderResult_TypeSafety tests that the function maintains type safety
func TestFromReaderResult_TypeSafety(t *testing.T) {
	t.Run("maintains input type", func(t *testing.T) {
		// This test verifies that the input type is preserved
		intToStringRR := func(input int) result.Result[string] {
			return result.Of(fmt.Sprintf("%d", input))
		}

		validator := FromReaderResult[int, string](intToStringRR)

		// This should compile and work correctly
		validationResult := validator(42)(nil)
		assert.Equal(t, validation.Success("42"), validationResult)
	})

	t.Run("maintains output type", func(t *testing.T) {
		// This test verifies that the output type is preserved
		stringToIntRR := result.Eitherize1(strconv.Atoi)

		validator := FromReaderResult[string, int](stringToIntRR)
		validationResult := validator("42")(nil)

		// The result should be Validation[int]
		assert.Equal(t, validation.Success(42), validationResult)
	})

	t.Run("works with different type combinations", func(t *testing.T) {
		type Input struct{ Value string }
		type Output struct{ Result int }

		transformRR := result.Eitherize1(func(input Input) (Output, error) {
			val, err := strconv.Atoi(input.Value)
			if err != nil {
				return Output{}, err
			}
			return Output{Result: val}, nil
		})

		validator := FromReaderResult[Input, Output](transformRR)
		validationResult := validator(Input{Value: "42"})(nil)

		assert.Equal(t, validation.Success(Output{Result: 42}), validationResult)
	})
}

// BenchmarkFromReaderResult_Success benchmarks the success path
func BenchmarkFromReaderResult_Success(b *testing.B) {
	successRR := func(input int) result.Result[int] {
		return result.Of(input * 2)
	}

	validator := FromReaderResult[int, int](successRR)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator(42)(nil)
	}
}

// BenchmarkFromReaderResult_Failure benchmarks the failure path
func BenchmarkFromReaderResult_Failure(b *testing.B) {
	failureRR := func(input int) result.Result[int] {
		return result.Left[int](errors.New("error"))
	}

	validator := FromReaderResult[int, int](failureRR)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator(42)(nil)
	}
}

// BenchmarkFromReaderResult_WithContext benchmarks with context
func BenchmarkFromReaderResult_WithContext(b *testing.B) {
	successRR := func(input int) result.Result[int] {
		return result.Of(input * 2)
	}

	validator := FromReaderResult[int, int](successRR)
	ctx := Context{
		{Key: "user", Type: "User"},
		{Key: "age", Type: "int"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = validator(42)(ctx)
	}
}
