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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
)

func TestSequenceReader(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: ReaderIOResult[Reader[string, int]]
		// = func(context.Context) func() Either[error, func(string) int]
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				return either.Right[error](func(s string) int {
					return 10 + len(s)
				})
			}
		}

		// Sequenced: func(string) func(context.Context) IOResult[int]
		// The Reader environment (string) is now the first parameter
		sequenced := SequenceReader(original)

		ctx := t.Context()

		// Test original
		result1 := original(ctx)()
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1("hello")
		assert.Equal(t, 15, value1)

		// Test sequenced - note the flipped order: string first, then context
		result2 := sequenced("hello")(ctx)()
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}

		// Original: ReaderIOResult[Reader[Database, string]]
		query := func(ctx context.Context) func() Either[Reader[Database, string]] {
			return func() Either[Reader[Database, string]] {
				if ctx.Err() != nil {
					return either.Left[Reader[Database, string]](ctx.Err())
				}
				return either.Right[error](func(db Database) string {
					return fmt.Sprintf("Query on %s", db.ConnectionString)
				})
			}
		}

		db := Database{ConnectionString: "localhost:5432"}
		ctx := t.Context()

		expected := "Query on localhost:5432"

		// Sequence it
		sequenced := SequenceReader(query)

		// Test original with valid inputs
		result1 := query(ctx)()
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(db)
		assert.Equal(t, expected, value1)

		// Test sequenced with valid inputs - Database first, then context
		result2 := sequenced(db)(ctx)()
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, expected, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				return either.Left[Reader[string, int]](expectedError)
			}
		}

		ctx := t.Context()

		// Test original with error
		result1 := original(ctx)()
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error is preserved
		sequenced := SequenceReader(original)
		result2 := sequenced("test")(ctx)()
		assert.True(t, either.IsLeft(result2))
		_, err2 := either.Unwrap(result2)
		assert.Equal(t, expectedError, err2)
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original function
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				return either.Right[error](func(s string) int {
					return 3 * len(s)
				})
			}
		}

		ctx := t.Context()

		// Sequence
		sequenced := SequenceReader(original)

		// Test that sequence produces correct results
		result1 := original(ctx)()
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1("test")

		result2 := sequenced("test")(ctx)()
		value2, _ := either.Unwrap(result2)

		assert.Equal(t, value1, value2)
		assert.Equal(t, 12, value2) // 3 * 4
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				return either.Right[error](func(s string) int {
					return len(s)
				})
			}
		}

		ctx := t.Context()
		sequenced := SequenceReader(original)

		// Test with zero values
		result1 := original(ctx)()
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1("")
		assert.Equal(t, 0, value1)

		result2 := sequenced("")(ctx)()
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 0, value2)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				if ctx.Err() != nil {
					return either.Left[Reader[string, int]](ctx.Err())
				}
				return either.Right[error](func(s string) int {
					return len(s)
				})
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		sequenced := SequenceReader(original)

		result := sequenced("test")(ctx)()
		assert.True(t, either.IsLeft(result))
		_, err := either.Unwrap(result)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("enables point-free style with partial application", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := func(ctx context.Context) func() Either[Reader[Config, int]] {
			return func() Either[Reader[Config, int]] {
				return either.Right[error](func(cfg Config) int {
					return cfg.Multiplier * 10
				})
			}
		}

		// Sequence to enable partial application
		sequenced := SequenceReader(original)

		// Partially apply the Config
		cfg := Config{Multiplier: 5}
		withConfig := sequenced(cfg)

		// Now we have a ReaderIOResult[int] that can be used in different contexts
		ctx1 := t.Context()
		result1 := withConfig(ctx1)()
		assert.True(t, either.IsRight(result1))
		value1, _ := either.Unwrap(result1)
		assert.Equal(t, 50, value1)

		// Can reuse with different context
		ctx2 := t.Context()
		result2 := withConfig(ctx2)()
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 50, value2)
	})
}

func TestSequenceReaderIO(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: ReaderIOResult[ReaderIO[int]]
		// = func(context.Context) func() Either[error, func(context.Context) func() int]
		original := func(ctx context.Context) func() Either[ReaderIO[int]] {
			return func() Either[ReaderIO[int]] {
				return either.Right[error](func(innerCtx context.Context) func() int {
					return func() int {
						return 20
					}
				})
			}
		}

		ctx := t.Context()
		sequenced := SequenceReaderIO(original)

		// Test original
		result1 := original(ctx)()
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(ctx)()
		assert.Equal(t, 20, value1)

		// Test sequenced - context first, then context again for inner ReaderIO
		result2 := sequenced(ctx)(ctx)()
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 20, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(ctx context.Context) func() Either[ReaderIO[int]] {
			return func() Either[ReaderIO[int]] {
				return either.Left[ReaderIO[int]](expectedError)
			}
		}

		ctx := t.Context()

		// Test original with error
		result1 := original(ctx)()
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error is preserved
		sequenced := SequenceReaderIO(original)
		result2 := sequenced(ctx)(ctx)()
		assert.True(t, either.IsLeft(result2))
		_, err2 := either.Unwrap(result2)
		assert.Equal(t, expectedError, err2)
	})

	t.Run("respects context cancellation in outer context", func(t *testing.T) {
		original := func(ctx context.Context) func() Either[ReaderIO[int]] {
			return func() Either[ReaderIO[int]] {
				if ctx.Err() != nil {
					return either.Left[ReaderIO[int]](ctx.Err())
				}
				return either.Right[error](func(innerCtx context.Context) func() int {
					return func() int {
						return 20
					}
				})
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		sequenced := SequenceReaderIO(original)

		result := sequenced(ctx)(ctx)()
		assert.True(t, either.IsLeft(result))
		_, err := either.Unwrap(result)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestSequenceReaderResult(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: ReaderIOResult[ReaderResult[int]]
		// = func(context.Context) func() Either[error, func(context.Context) Either[error, int]]
		original := func(ctx context.Context) func() Either[ReaderResult[int]] {
			return func() Either[ReaderResult[int]] {
				return either.Right[error](func(innerCtx context.Context) Either[int] {
					return either.Right[error](20)
				})
			}
		}

		ctx := t.Context()
		sequenced := SequenceReaderResult(original)

		// Test original
		result1 := original(ctx)()
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(ctx)
		assert.True(t, either.IsRight(innerResult1))
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, 20, value1)

		// Test sequenced
		result2 := sequenced(ctx)(ctx)()
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 20, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(ctx context.Context) func() Either[ReaderResult[int]] {
			return func() Either[ReaderResult[int]] {
				return either.Left[ReaderResult[int]](expectedError)
			}
		}

		ctx := t.Context()

		// Test original with error
		result1 := original(ctx)()
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error is preserved
		sequenced := SequenceReaderResult(original)
		result2 := sequenced(ctx)(ctx)()
		assert.True(t, either.IsLeft(result2))
		_, err2 := either.Unwrap(result2)
		assert.Equal(t, expectedError, err2)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		// Original that fails at inner level
		original := func(ctx context.Context) func() Either[ReaderResult[int]] {
			return func() Either[ReaderResult[int]] {
				return either.Right[error](func(innerCtx context.Context) Either[int] {
					return either.Left[int](expectedError)
				})
			}
		}

		ctx := t.Context()

		// Test original with inner error
		result1 := original(ctx)()
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(ctx)
		assert.True(t, either.IsLeft(innerResult1))
		_, innerErr1 := either.Unwrap(innerResult1)
		assert.Equal(t, expectedError, innerErr1)

		// Test sequenced with inner error
		sequenced := SequenceReaderResult(original)
		result2 := sequenced(ctx)(ctx)()
		assert.True(t, either.IsLeft(result2))
		_, innerErr2 := either.Unwrap(result2)
		assert.Equal(t, expectedError, innerErr2)
	})

	t.Run("handles errors at different levels", func(t *testing.T) {
		// Original that can fail at both levels
		makeOriginal := func(x int) ReaderIOResult[ReaderResult[int]] {
			return func(ctx context.Context) func() Either[ReaderResult[int]] {
				return func() Either[ReaderResult[int]] {
					if x < -10 {
						return either.Left[ReaderResult[int]](errors.New("outer: too negative"))
					}
					return either.Right[error](func(innerCtx context.Context) Either[int] {
						if x < 0 {
							return either.Left[int](errors.New("inner: negative value"))
						}
						return either.Right[error](x * 2)
					})
				}
			}
		}

		ctx := t.Context()

		// Test outer error
		sequenced1 := SequenceReaderResult(makeOriginal(-20))
		result1 := sequenced1(ctx)(ctx)()
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Contains(t, err1.Error(), "outer")

		// Test inner error
		sequenced2 := SequenceReaderResult(makeOriginal(-5))
		result2 := sequenced2(ctx)(ctx)()
		assert.True(t, either.IsLeft(result2))
		_, err2 := either.Unwrap(result2)
		assert.Contains(t, err2.Error(), "inner")

		// Test success
		sequenced3 := SequenceReaderResult(makeOriginal(10))
		result3 := sequenced3(ctx)(ctx)()
		assert.True(t, either.IsRight(result3))
		value3, _ := either.Unwrap(result3)
		assert.Equal(t, 20, value3)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		original := func(ctx context.Context) func() Either[ReaderResult[int]] {
			return func() Either[ReaderResult[int]] {
				if ctx.Err() != nil {
					return either.Left[ReaderResult[int]](ctx.Err())
				}
				return either.Right[error](func(innerCtx context.Context) Either[int] {
					if innerCtx.Err() != nil {
						return either.Left[int](innerCtx.Err())
					}
					return either.Right[error](20)
				})
			}
		}

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		sequenced := SequenceReaderResult(original)

		result := sequenced(ctx)(ctx)()
		assert.True(t, either.IsLeft(result))
		_, err := either.Unwrap(result)
		assert.Equal(t, context.Canceled, err)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(ctx context.Context) func() Either[Reader[Empty, int]] {
			return func() Either[Reader[Empty, int]] {
				return either.Right[error](func(e Empty) int {
					return 20
				})
			}
		}

		ctx := t.Context()
		empty := Empty{}
		sequenced := SequenceReader(original)

		result1 := original(ctx)()
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(empty)
		assert.Equal(t, 20, value1)

		result2 := sequenced(empty)(ctx)()
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(ctx context.Context) func() Either[Reader[*Data, int]] {
			return func() Either[Reader[*Data, int]] {
				return either.Right[error](func(d *Data) int {
					if d == nil {
						return 42
					}
					return 42 + d.Value
				})
			}
		}

		ctx := t.Context()
		data := &Data{Value: 100}
		sequenced := SequenceReader(original)

		// Test with non-nil pointer
		result1 := original(ctx)()
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(data)
		assert.Equal(t, 142, value1)

		result2 := sequenced(data)(ctx)()
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3 := sequenced(nil)(ctx)()
		value3, _ := either.Unwrap(result3)
		assert.Equal(t, 42, value3)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		// The same inputs should always produce the same outputs
		original := func(ctx context.Context) func() Either[Reader[string, int]] {
			return func() Either[Reader[string, int]] {
				return either.Right[error](func(s string) int {
					return 10 + len(s)
				})
			}
		}

		ctx := t.Context()
		sequenced := SequenceReader(original)

		// Call multiple times with same inputs
		for range 5 {
			result1 := original(ctx)()
			innerFunc1, _ := either.Unwrap(result1)
			value1 := innerFunc1("hello")
			assert.Equal(t, 15, value1)

			result2 := sequenced("hello")(ctx)()
			value2, _ := either.Unwrap(result2)
			assert.Equal(t, 15, value2)
		}
	})
}

func TestTraverseReader(t *testing.T) {
	t.Run("basic transformation with Reader dependency", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := Right(10)

		// Reader-based transformation
		multiply := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(multiply)
		result := traversed(original)

		// Provide Config and execute
		cfg := Config{Multiplier: 5}
		ctx := t.Context()
		finalResult := result(cfg)(ctx)()

		assert.True(t, either.IsRight(finalResult))
		value, _ := either.Unwrap(finalResult)
		assert.Equal(t, 50, value)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		expectedError := errors.New("computation failed")

		// Original computation that fails
		original := Left[int](expectedError)

		// Reader-based transformation (won't be called)
		multiply := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(multiply)
		result := traversed(original)

		// Provide Config and execute
		cfg := Config{Multiplier: 5}
		ctx := t.Context()
		finalResult := result(cfg)(ctx)()

		assert.True(t, either.IsLeft(finalResult))
		_, err := either.Unwrap(finalResult)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		type Database struct {
			Prefix string
		}

		// Original computation producing an int
		original := Right(42)

		// Reader-based transformation: int -> string using Database
		format := func(x int) func(Database) string {
			return func(db Database) string {
				return fmt.Sprintf("%s:%d", db.Prefix, x)
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(format)
		result := traversed(original)

		// Provide Database and execute
		db := Database{Prefix: "ID"}
		ctx := t.Context()
		finalResult := result(db)(ctx)()

		assert.True(t, either.IsRight(finalResult))
		value, _ := either.Unwrap(finalResult)
		assert.Equal(t, "ID:42", value)
	})

	t.Run("works with struct environments", func(t *testing.T) {
		type Settings struct {
			Prefix string
			Suffix string
		}

		// Original computation
		original := Right("value")

		// Reader-based transformation using Settings
		decorate := func(s string) func(Settings) string {
			return func(settings Settings) string {
				return settings.Prefix + s + settings.Suffix
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(decorate)
		result := traversed(original)

		// Provide Settings and execute
		settings := Settings{Prefix: "[", Suffix: "]"}
		ctx := t.Context()
		finalResult := result(settings)(ctx)()

		assert.True(t, either.IsRight(finalResult))
		value, _ := either.Unwrap(finalResult)
		assert.Equal(t, "[value]", value)
	})

	t.Run("enables partial application", func(t *testing.T) {
		type Config struct {
			Factor int
		}

		// Original computation
		original := Right(10)

		// Reader-based transformation
		scale := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x * cfg.Factor
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(scale)
		result := traversed(original)

		// Partially apply Config
		cfg := Config{Factor: 3}
		withConfig := result(cfg)

		// Can now use with different contexts
		ctx1 := t.Context()
		finalResult1 := withConfig(ctx1)()
		assert.True(t, either.IsRight(finalResult1))
		value1, _ := either.Unwrap(finalResult1)
		assert.Equal(t, 30, value1)

		// Reuse with different context
		ctx2 := t.Context()
		finalResult2 := withConfig(ctx2)()
		assert.True(t, either.IsRight(finalResult2))
		value2, _ := either.Unwrap(finalResult2)
		assert.Equal(t, 30, value2)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		type Config struct {
			Value int
		}

		// Original computation that checks context
		original := func(ctx context.Context) func() Either[int] {
			return func() Either[int] {
				if ctx.Err() != nil {
					return either.Left[int](ctx.Err())
				}
				return either.Right[error](10)
			}
		}

		// Reader-based transformation
		multiply := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x * cfg.Value
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(multiply)
		result := traversed(original)

		// Use canceled context
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		cfg := Config{Value: 5}
		finalResult := result(cfg)(ctx)()

		assert.True(t, either.IsLeft(finalResult))
		_, err := either.Unwrap(finalResult)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("works with zero values", func(t *testing.T) {
		type Config struct {
			Offset int
		}

		// Original computation with zero value
		original := Right(0)

		// Reader-based transformation
		add := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x + cfg.Offset
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(add)
		result := traversed(original)

		// Provide Config with zero offset
		cfg := Config{Offset: 0}
		ctx := t.Context()
		finalResult := result(cfg)(ctx)()

		assert.True(t, either.IsRight(finalResult))
		value, _ := either.Unwrap(finalResult)
		assert.Equal(t, 0, value)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := Right(5)

		// First Reader-based transformation
		multiply := func(x int) Reader[Config, int] {
			return func(cfg Config) int {
				return x * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(multiply)
		result := traversed(original)

		// Provide Config and execute
		cfg := Config{Multiplier: 4}
		ctx := t.Context()
		finalResult := result(cfg)(ctx)()

		assert.True(t, either.IsRight(finalResult))
		value, _ := either.Unwrap(finalResult)
		assert.Equal(t, 20, value) // 5 * 4 = 20
	})

	t.Run("works with complex Reader logic", func(t *testing.T) {
		type ValidationRules struct {
			MinValue int
			MaxValue int
		}

		// Original computation
		original := Right(50)

		// Reader-based transformation with validation logic
		validate := func(x int) func(ValidationRules) int {
			return func(rules ValidationRules) int {
				if x < rules.MinValue {
					return rules.MinValue
				}
				if x > rules.MaxValue {
					return rules.MaxValue
				}
				return x
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader(validate)
		result := traversed(original)

		// Test with value within range
		rules1 := ValidationRules{MinValue: 0, MaxValue: 100}
		ctx := t.Context()
		finalResult1 := result(rules1)(ctx)()
		assert.True(t, either.IsRight(finalResult1))
		value1, _ := either.Unwrap(finalResult1)
		assert.Equal(t, 50, value1)

		// Test with value above max
		rules2 := ValidationRules{MinValue: 0, MaxValue: 30}
		finalResult2 := result(rules2)(ctx)()
		assert.True(t, either.IsRight(finalResult2))
		value2, _ := either.Unwrap(finalResult2)
		assert.Equal(t, 30, value2) // Clamped to max

		// Test with value below min
		rules3 := ValidationRules{MinValue: 60, MaxValue: 100}
		finalResult3 := result(rules3)(ctx)()
		assert.True(t, either.IsRight(finalResult3))
		value3, _ := either.Unwrap(finalResult3)
		assert.Equal(t, 60, value3) // Clamped to min
	})
}
