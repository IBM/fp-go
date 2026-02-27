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

package readerreaderioresult

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Prefix string
	MaxLen int
}

var testConfig = TestConfig{
	Prefix: "test",
	MaxLen: 100,
}

// TestEitherize_Success tests successful conversion with Eitherize
func TestEitherize_Success(t *testing.T) {
	t.Run("converts successful function to ReaderReaderIOResult", func(t *testing.T) {
		// Arrange
		successFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return cfg.Prefix + "-success", nil
		}
		rr := Eitherize(successFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of("test-success"), outcome)
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
		rr := Eitherize(contextFunc)

		ctx := context.WithValue(context.Background(), key, expectedValue)

		// Act
		outcome := rr(testConfig)(ctx)()

		// Assert
		assert.Equal(t, result.Of(expectedValue), outcome)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Arrange
		intFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return cfg.MaxLen, nil
		}
		rr := Eitherize(intFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(100), outcome)
	})
}

// TestEitherize_Failure tests error handling with Eitherize
func TestEitherize_Failure(t *testing.T) {
	t.Run("converts error to Left", func(t *testing.T) {
		// Arrange
		expectedErr := errors.New("operation failed")
		failFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return "", expectedErr
		}
		rr := Eitherize(failFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
		assert.Equal(t, result.Left[string](expectedErr), outcome)
	})

	t.Run("preserves error message", func(t *testing.T) {
		// Arrange
		expectedErr := fmt.Errorf("validation error: field is required")
		failFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return 0, expectedErr
		}
		rr := Eitherize(failFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
		leftValue := result.MonadFold(outcome,
			F.Identity[error],
			func(int) error { return nil },
		)
		assert.Equal(t, expectedErr, leftValue)
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
		rr := Eitherize(nilCtxFunc)

		// Act
		outcome := rr(testConfig)(nil)()

		// Assert
		assert.Equal(t, result.Of("nil-context"), outcome)
	})

	t.Run("handles zero value config", func(t *testing.T) {
		// Arrange
		zeroFunc := func(cfg TestConfig, ctx context.Context) (string, error) {
			return cfg.Prefix, nil
		}
		rr := Eitherize(zeroFunc)

		// Act
		outcome := rr(TestConfig{})(context.Background())()

		// Assert
		assert.Equal(t, result.Of(""), outcome)
	})

	t.Run("handles pointer types", func(t *testing.T) {
		// Arrange
		type User struct {
			Name string
		}
		ptrFunc := func(cfg TestConfig, ctx context.Context) (*User, error) {
			return &User{Name: "Alice"}, nil
		}
		rr := Eitherize(ptrFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsRight(outcome))
		user := result.MonadFold(outcome,
			func(error) *User { return nil },
			F.Identity[*User],
		)
		assert.NotNil(t, user)
		assert.Equal(t, "Alice", user.Name)
	})
}

// TestEitherize_Integration tests integration with other operations
func TestEitherize_Integration(t *testing.T) {
	t.Run("composes with Map", func(t *testing.T) {
		// Arrange
		baseFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return 42, nil
		}
		rr := Eitherize(baseFunc)

		// Act
		pipeline := F.Pipe1(
			rr,
			Map[TestConfig](func(n int) string { return strconv.Itoa(n) }),
		)
		outcome := pipeline(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of("42"), outcome)
	})

	t.Run("composes with Chain", func(t *testing.T) {
		// Arrange
		firstFunc := func(cfg TestConfig, ctx context.Context) (int, error) {
			return 10, nil
		}
		secondFunc := func(n int) ReaderReaderIOResult[TestConfig, string] {
			return Of[TestConfig](fmt.Sprintf("value: %d", n))
		}

		// Act
		pipeline := F.Pipe1(
			Eitherize(firstFunc),
			Chain[TestConfig](secondFunc),
		)
		outcome := pipeline(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of("value: 10"), outcome)
	})
}

// TestEitherize1_Success tests successful conversion with Eitherize1
func TestEitherize1_Success(t *testing.T) {
	t.Run("converts successful function to Kleisli", func(t *testing.T) {
		// Arrange
		addFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			return n + cfg.MaxLen, nil
		}
		kleisli := Eitherize1(addFunc)

		// Act
		outcome := kleisli(10)(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(110), outcome)
	})

	t.Run("works with string input", func(t *testing.T) {
		// Arrange
		concatFunc := func(cfg TestConfig, ctx context.Context, s string) (string, error) {
			return cfg.Prefix + "-" + s, nil
		}
		kleisli := Eitherize1(concatFunc)

		// Act
		outcome := kleisli("input")(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of("test-input"), outcome)
	})

	t.Run("preserves context in Kleisli", func(t *testing.T) {
		// Arrange
		type ctxKey string
		key := ctxKey("multiplier")

		multiplyFunc := func(cfg TestConfig, ctx context.Context, n int) (int, error) {
			multiplier := ctx.Value(key)
			if multiplier == nil {
				return n, nil
			}
			return n * multiplier.(int), nil
		}
		kleisli := Eitherize1(multiplyFunc)

		ctx := context.WithValue(context.Background(), key, 3)

		// Act
		outcome := kleisli(5)(testConfig)(ctx)()

		// Assert
		assert.Equal(t, result.Of(15), outcome)
	})
}

// TestEitherize1_Failure tests error handling with Eitherize1
func TestEitherize1_Failure(t *testing.T) {
	t.Run("converts error to Left in Kleisli", func(t *testing.T) {
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
		outcome := kleisli(0)(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
		assert.Equal(t, result.Left[int](expectedErr), outcome)
	})

	t.Run("preserves error context", func(t *testing.T) {
		// Arrange
		validateFunc := func(cfg TestConfig, ctx context.Context, s string) (string, error) {
			if len(s) > cfg.MaxLen {
				return "", fmt.Errorf("string too long: %d > %d", len(s), cfg.MaxLen)
			}
			return s, nil
		}
		kleisli := Eitherize1(validateFunc)

		longString := string(make([]byte, 200))

		// Act
		outcome := kleisli(longString)(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
		leftValue := result.MonadFold(outcome,
			F.Identity[error],
			func(string) error { return nil },
		)
		assert.Contains(t, leftValue.Error(), "string too long")
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
		outcome := kleisli(0)(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(0), outcome)
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
			return in.Value, nil
		}
		kleisli := Eitherize1(ptrFunc)

		// Act
		outcome := kleisli(&Input{Value: 42})(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(42), outcome)
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
		outcome := kleisli((*Input)(nil))(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
	})
}

// TestEitherize1_Integration tests integration with other operations
func TestEitherize1_Integration(t *testing.T) {
	t.Run("composes with Chain", func(t *testing.T) {
		// Arrange
		parseFunc := func(cfg TestConfig, ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}
		doubleFunc := func(n int) ReaderReaderIOResult[TestConfig, int] {
			return Of[TestConfig](n * 2)
		}

		parseKleisli := Eitherize1(parseFunc)

		// Act
		pipeline := F.Pipe2(
			Of[TestConfig]("42"),
			Chain[TestConfig](parseKleisli),
			Chain[TestConfig](doubleFunc),
		)
		outcome := pipeline(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(84), outcome)
	})

	t.Run("handles error in chain", func(t *testing.T) {
		// Arrange
		parseFunc := func(cfg TestConfig, ctx context.Context, s string) (int, error) {
			return strconv.Atoi(s)
		}
		parseKleisli := Eitherize1(parseFunc)

		// Act
		pipeline := F.Pipe1(
			Of[TestConfig]("not-a-number"),
			Chain[TestConfig](parseKleisli),
		)
		outcome := pipeline(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsLeft(outcome))
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
			Of[TestConfig]("123"),
			Chain[TestConfig](parseKleisli),
			Chain[TestConfig](formatKleisli),
		)
		outcome := pipeline(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of("test-123"), outcome)
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
				Data:  map[string]int{"key": 42},
				Count: 1,
			}, nil
		}
		rr := Eitherize(complexFunc)

		// Act
		outcome := rr(testConfig)(context.Background())()

		// Assert
		assert.True(t, result.IsRight(outcome))
		value := result.MonadFold(outcome,
			func(error) ComplexResult { return ComplexResult{} },
			F.Identity[ComplexResult],
		)
		assert.Equal(t, 42, value.Data["key"])
		assert.Equal(t, 1, value.Count)
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
		outcome := kleisli(Input{ID: 99})(testConfig)(context.Background())()

		// Assert
		assert.Equal(t, result.Of(Output{Name: "test-99"}), outcome)
	})
}
