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

package effect

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// TestEitherize_Success tests successful conversion with Eitherize
func TestEitherize_Success(t *testing.T) {
	t.Run("converts successful function to Effect", func(t *testing.T) {
		// Arrange
		successFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return cfg.Prefix + "-success", nil
		}
		eff := Eitherize(successFunc)

		// Act
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "LOG-success", result)
	})

	t.Run("preserves context values", func(t *testing.T) {
		// Arrange
		type ctxKey string
		key := ctxKey("testKey")
		expectedValue := "contextValue"

		contextFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			value := ctx.Value(key)
			if value == nil {
				return "", errors.New("context value not found")
			}
			return value.(string), nil
		}
		eff := Eitherize(contextFunc)

		// Act
		ioResult := Provide[string](testConfig)(eff)
		readerResult := RunSync(ioResult)
		ctx := context.WithValue(context.Background(), key, expectedValue)
		result, err := readerResult(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedValue, result)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Arrange
		intFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return cfg.Multiplier, nil
		}
		eff := Eitherize(intFunc)

		// Act
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 3, result)
	})
}

// TestEitherize_Failure tests error handling with Eitherize
func TestEitherize_Failure(t *testing.T) {
	t.Run("converts error to failure", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("operation failed")
		failFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return "", expectedErr
		}
		eff := Eitherize(failFunc)

		// Act
		_, err := runEffect(eff, testConfig)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("preserves error message", func(t *testing.T) {
		// Arrange
		expectedErr := fmt.Errorf("validation error: field is required")
		failFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return 0, expectedErr
		}
		eff := Eitherize(failFunc)

		// Act
		_, err := runEffect(eff, testConfig)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

// TestEitherize_EdgeCases tests edge cases for Eitherize
func TestEitherize_EdgeCases(t *testing.T) {
	t.Run("handles nil context", func(t *testing.T) {
		// Arrange
		nilCtxFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			if ctx == nil {
				return "nil-context", nil
			}
			return "non-nil-context", nil
		}
		eff := Eitherize(nilCtxFunc)

		// Act
		ioResult := Provide[string](testConfig)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(nil)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "nil-context", result)
	})

	t.Run("handles zero value config", func(t *testing.T) {
		// Arrange
		zeroFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return cfg.Prefix, nil
		}
		eff := Eitherize(zeroFunc)

		// Act
		result, err := runEffect(eff, TestConfig{})

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "", result)
	})

	t.Run("handles pointer types", func(t *testing.T) {
		// Arrange
		type User struct {
			Name string
		}
		ptrFunc := func(cfg TestConfig, ctx context.Context) (*User, error) {
			return &User{Name: cfg.Prefix}, nil
		}
		eff := Eitherize(ptrFunc)

		// Act
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "LOG", result.Name)
	})
}

// TestEitherize_Integration tests integration with other operations
func TestEitherize_Integration(t *testing.T) {
	t.Run("composes with Map", func(t *testing.T) {
		// Arrange
		baseFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return cfg.Multiplier, nil
		}
		eff := Eitherize(baseFunc)

		// Act
		pipeline := F.Pipe1(
			eff,
			Map[TestConfig](func(n int) string { return strconv.Itoa(n) }),
		)
		result, err := runEffect(pipeline, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "3", result)
	})

	t.Run("composes with Chain", func(t *testing.T) {
		// Arrange
		firstFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return cfg.Multiplier, nil
		}
		secondFunc := func(n int) Effect[TestConfig, string] {
			return Succeed[TestConfig](fmt.Sprintf("value: %d", n))
		}

		// Act
		pipeline := F.Pipe1(
			Eitherize(firstFunc),
			Chain[TestConfig](secondFunc),
		)
		result, err := runEffect(pipeline, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "value: 3", result)
	})
}

// TestEitherize1_Success tests successful conversion with Eitherize1
func TestEitherize1_Success(t *testing.T) {
	t.Run("converts successful function to Kleisli", func(t *testing.T) {
		// Arrange
		multiplyFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			return n * cfg.Multiplier, nil
		}
		kleisli := Eitherize1(multiplyFunc)

		// Act
		eff := kleisli(10)
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 30, result)
	})

	t.Run("works with string input", func(t *testing.T) {
		// Arrange
		concatFunc := func(cfg TestConfig, ctx context.Context, s string) (string, error) {
			return cfg.Prefix + "-" + s, nil
		}
		kleisli := Eitherize1(concatFunc)

		// Act
		eff := kleisli("input")
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "LOG-input", result)
	})

	t.Run("preserves context in Kleisli", func(t *testing.T) {
		// Arrange
		type ctxKey string
		key := ctxKey("factor")

		scaleFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			factor := ctx.Value(key)
			if factor == nil {
				return n * cfg.Multiplier, nil
			}
			return n * factor.(int), nil
		}
		kleisli := Eitherize1(scaleFunc)

		// Act
		eff := kleisli(5)
		ioResult := Provide[int](testConfig)(eff)
		readerResult := RunSync(ioResult)
		ctx := context.WithValue(context.Background(), key, 7)
		result, err := readerResult(ctx)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 35, result)
	})
}

// TestEitherize1_Failure tests error handling with Eitherize1
func TestEitherize1_Failure(t *testing.T) {
	t.Run("converts error to failure in Kleisli", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("division by zero")
		divideFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			if n == 0 {
				return 0, expectedErr
			}
			return 100 / n, nil
		}
		kleisli := Eitherize1(divideFunc)

		// Act
		eff := kleisli(0)
		_, err := runEffect(eff, testConfig)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("preserves error context", func(t *testing.T) {
		// Arrange
		validateFunc := func(cfg TestConfig, ctx context.Context, s string) (string, error) {
			if len(s) > 10 {
				return "", fmt.Errorf("string too long: %d > 10", len(s))
			}
			return s, nil
		}
		kleisli := Eitherize1(validateFunc)

		// Act
		eff := kleisli("this-string-is-too-long")
		_, err := runEffect(eff, testConfig)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "string too long")
	})
}

// TestEitherize1_EdgeCases tests edge cases for Eitherize1
func TestEitherize1_EdgeCases(t *testing.T) {
	t.Run("handles zero value input", func(t *testing.T) {
		// Arrange
		zeroFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			return n, nil
		}
		kleisli := Eitherize1(zeroFunc)

		// Act
		eff := kleisli(0)
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})

	t.Run("handles pointer input", func(t *testing.T) {
		// Arrange
		type Input struct {
			Value int
		}
		ptrFunc := func(cfg TestConfig, ctx context.Context, in *Input) (int, error) {
			if in == nil {
				return 0, errors.New("nil input")
			}
			return in.Value * cfg.Multiplier, nil
		}
		kleisli := Eitherize1(ptrFunc)

		// Act
		eff := kleisli(&Input{Value: 7})
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 21, result)
	})

	t.Run("handles nil pointer input", func(t *testing.T) {
		// Arrange
		type Input struct {
			Value int
		}
		ptrFunc := func(cfg TestConfig, ctx context.Context, in *Input) (int, error) {
			if in == nil {
				return 0, errors.New("nil input")
			}
			return in.Value, nil
		}
		kleisli := Eitherize1(ptrFunc)

		// Act
		eff := kleisli((*Input)(nil))
		_, err := runEffect(eff, testConfig)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nil input")
	})
}

// TestEitherize1_Integration tests integration with other operations
func TestEitherize1_Integration(t *testing.T) {
	t.Run("composes with Chain", func(t *testing.T) {
		// Arrange
		parseFunc := func(cfg TestConfig, ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}
		doubleFunc := func(n int) Effect[TestConfig, int] {
			return Succeed[TestConfig](n * 2)
		}

		parseKleisli := Eitherize1(parseFunc)

		// Act
		pipeline := F.Pipe2(
			Succeed[TestConfig]("42"),
			Chain[TestConfig](parseKleisli),
			Chain[TestConfig](doubleFunc),
		)
		result, err := runEffect(pipeline, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 84, result)
	})

	t.Run("handles error in chain", func(t *testing.T) {
		// Arrange
		parseFunc := func(cfg TestConfig, ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}
		parseKleisli := Eitherize1(parseFunc)

		// Act
		pipeline := F.Pipe1(
			Succeed[TestConfig]("not-a-number"),
			Chain[TestConfig](parseKleisli),
		)
		_, err := runEffect(pipeline, testConfig)

		// Assert
		assert.Error(t, err)
	})

	t.Run("composes multiple Kleisli arrows", func(t *testing.T) {
		// Arrange
		parseFunc := func(cfg TestConfig, ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}
		formatFunc := func(cfg TestConfig, ctx context.Context, n int) (string, error) {
			return fmt.Sprintf("%s-%d", cfg.Prefix, n), nil
		}

		parseKleisli := Eitherize1(parseFunc)
		formatKleisli := Eitherize1(formatFunc)

		// Act
		pipeline := F.Pipe2(
			Succeed[TestConfig]("123"),
			Chain[TestConfig](parseKleisli),
			Chain[TestConfig](formatKleisli),
		)
		result, err := runEffect(pipeline, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "LOG-123", result)
	})
}

// TestEitherize_TypeSafety tests type safety across different scenarios
func TestEitherize_TypeSafety(t *testing.T) {
	t.Run("Eitherize with complex types", func(t *testing.T) {
		// Arrange
		type ComplexResult struct {
			Data  map[string]int
			Count int
		}

		complexFunc := func(cfg TestConfig, ctx context.Context) (ComplexResult, error) {
			return ComplexResult{
				Data:  map[string]int{cfg.Prefix: cfg.Multiplier},
				Count: cfg.Multiplier,
			}, nil
		}
		eff := Eitherize(complexFunc)

		// Act
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 3, result.Data["LOG"])
		assert.Equal(t, 3, result.Count)
	})

	t.Run("Eitherize1 with different input and output types", func(t *testing.T) {
		// Arrange
		type Input struct {
			ID int
		}
		type Output struct {
			Name string
		}

		convertFunc := func(cfg TestConfig, ctx context.Context, in Input) (Output, error) {
			return Output{Name: fmt.Sprintf("%s-%d", cfg.Prefix, in.ID)}, nil
		}
		kleisli := Eitherize1(convertFunc)

		// Act
		eff := kleisli(Input{ID: 99})
		result, err := runEffect(eff, testConfig)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, "LOG-99", result.Name)
	})
}
