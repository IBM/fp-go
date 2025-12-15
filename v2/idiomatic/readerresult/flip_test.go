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

package readerresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("sequences parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns ReaderResult[string, int]
		original := func(x int) (ReaderResult[string, int], error) {
			if x < 0 {
				return nil, errors.New("negative value")
			}
			return func(s string) (int, error) {
				return x + len(s), nil
			}, nil
		}

		// Sequenced: takes string first, then int
		sequenced := Sequence(original)

		// Test original
		innerFunc1, err1 := original(10)
		assert.NoError(t, err1)
		result1, err2 := innerFunc1("hello")
		assert.NoError(t, err2)
		assert.Equal(t, 15, result1)

		// Test sequenced
		result2, err3 := sequenced("hello")(10)
		assert.NoError(t, err3)
		assert.Equal(t, 15, result2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) (ReaderResult[string, int], error) {
			if x < 0 {
				return nil, expectedError
			}
			return func(s string) (int, error) {
				return x + len(s), nil
			}, nil
		}

		sequenced := Sequence(original)

		// Test with error
		_, err := sequenced("test")(-1)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		original := func(x int) (ReaderResult[string, int], error) {
			return func(s string) (int, error) {
				if S.IsEmpty(s) {
					return 0, expectedError
				}
				return x + len(s), nil
			}, nil
		}

		sequenced := Sequence(original)

		// Test with inner error
		_, err := sequenced("")(10)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Transform int to string
		original := func(x int) (ReaderResult[string, string], error) {
			return func(prefix string) (string, error) {
				return fmt.Sprintf("%s-%d", prefix, x), nil
			}, nil
		}

		sequenced := Sequence(original)

		result, err := sequenced("ID")(42)
		assert.NoError(t, err)
		assert.Equal(t, "ID-42", result)
	})

	t.Run("works with struct environments", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		original := func(cfg Config) (ReaderResult[Database, string], error) {
			if cfg.Timeout <= 0 {
				return nil, errors.New("invalid timeout")
			}
			return func(db Database) (string, error) {
				if S.IsEmpty(db.ConnectionString) {
					return "", errors.New("empty connection string")
				}
				return fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout), nil
			}, nil
		}

		sequenced := Sequence(original)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		result, err := sequenced(db)(cfg)
		assert.NoError(t, err)
		assert.Equal(t, "Query on localhost:5432 with timeout 30", result)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) (ReaderResult[string, int], error) {
			return func(s string) (int, error) {
				return x + len(s), nil
			}, nil
		}

		sequenced := Sequence(original)

		result, err := sequenced("")(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("sequences parameter order for Reader inner type", func(t *testing.T) {
		// Original: takes int, returns Reader[string, int]
		original := func(x int) (reader.Reader[string, int], error) {
			if x < 0 {
				return nil, errors.New("negative value")
			}
			return func(s string) int {
				return x + len(s)
			}, nil
		}

		// Sequenced: takes string first, then int
		sequenced := SequenceReader(original)

		// Test original
		readerFunc, err1 := original(10)
		assert.NoError(t, err1)
		value1 := readerFunc("hello")
		assert.Equal(t, 15, value1)

		// Test sequenced
		value2, err2 := sequenced("hello")(10)
		assert.NoError(t, err2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) (reader.Reader[string, int], error) {
			if x < 0 {
				return nil, expectedError
			}
			return func(s string) int {
				return x + len(s)
			}, nil
		}

		sequenced := SequenceReader(original)

		// Test with error
		_, err := sequenced("test")(-1)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Transform int to string using Reader
		original := func(x int) (reader.Reader[string, string], error) {
			return func(prefix string) string {
				return fmt.Sprintf("%s-%d", prefix, x)
			}, nil
		}

		sequenced := SequenceReader(original)

		result, err := sequenced("ID")(42)
		assert.NoError(t, err)
		assert.Equal(t, "ID-42", result)
	})

	t.Run("works with struct environments", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		original := func(x int) (reader.Reader[Config, int], error) {
			if x < 0 {
				return nil, errors.New("negative value")
			}
			return func(cfg Config) int {
				return x * cfg.Multiplier
			}, nil
		}

		sequenced := SequenceReader(original)

		cfg := Config{Multiplier: 5}
		result, err := sequenced(cfg)(10)
		assert.NoError(t, err)
		assert.Equal(t, 50, result)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) (reader.Reader[string, int], error) {
			return func(s string) int {
				return x + len(s)
			}, nil
		}

		sequenced := SequenceReader(original)

		result, err := sequenced("")(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, result)
	})
}

func TestTraverse(t *testing.T) {
	t.Run("basic transformation with environment swap", func(t *testing.T) {
		// Original: ReaderResult[int, int] - takes int environment, produces int
		original := func(x int) (int, error) {
			if x < 0 {
				return 0, errors.New("negative value")
			}
			return x * 2, nil
		}

		// Kleisli function: func(int) ReaderResult[string, int]
		kleisli := func(a int) ReaderResult[string, int] {
			return func(s string) (int, error) {
				return a + len(s), nil
			}
		}

		// Traverse returns: func(ReaderResult[int, int]) func(string) ReaderResult[int, int]
		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// result is func(string) ReaderResult[int, int]
		// Provide string first ("hello"), then int (10)
		value, err := result("hello")(10)
		assert.NoError(t, err)
		assert.Equal(t, 25, value) // (10 * 2) + len("hello") = 20 + 5 = 25
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) (int, error) {
			if x < 0 {
				return 0, expectedError
			}
			return x, nil
		}

		kleisli := func(a int) ReaderResult[string, int] {
			return func(s string) (int, error) {
				return a + len(s), nil
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Test with negative value to trigger error
		_, err := result("test")(-1)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves inner error from Kleisli", func(t *testing.T) {
		expectedError := errors.New("inner error")

		original := Ask[int]()

		kleisli := func(a int) ReaderResult[string, int] {
			return func(s string) (int, error) {
				if S.IsEmpty(s) {
					return 0, expectedError
				}
				return a + len(s), nil
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Test with empty string to trigger inner error
		_, err := result("")(10)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Transform int to string using environment-dependent logic
		original := Ask[int]()

		kleisli := func(a int) ReaderResult[string, string] {
			return func(prefix string) (string, error) {
				return fmt.Sprintf("%s-%d", prefix, a), nil
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		value, err := result("ID")(42)
		assert.NoError(t, err)
		assert.Equal(t, "ID-42", value)
	})

	t.Run("works with struct environments", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}
		type Database struct {
			Prefix string
		}

		original := func(cfg Config) (int, error) {
			if cfg.Multiplier <= 0 {
				return 0, errors.New("invalid multiplier")
			}
			return 10 * cfg.Multiplier, nil
		}

		kleisli := func(value int) ReaderResult[Database, string] {
			return func(db Database) (string, error) {
				return fmt.Sprintf("%s:%d", db.Prefix, value), nil
			}
		}

		traversed := Traverse[Config](kleisli)
		result := traversed(original)

		cfg := Config{Multiplier: 5}
		db := Database{Prefix: "result"}

		value, err := result(db)(cfg)
		assert.NoError(t, err)
		assert.Equal(t, "result:50", value)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		original := Ask[int]()

		// First transformation: multiply by environment value
		kleisli1 := func(a int) ReaderResult[int, int] {
			return func(multiplier int) (int, error) {
				return a * multiplier, nil
			}
		}

		traversed := Traverse[int](kleisli1)
		result := traversed(original)

		value, err := result(3)(5)
		assert.NoError(t, err)
		assert.Equal(t, 15, value) // 5 * 3 = 15
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := Ask[int]()

		kleisli := func(a int) ReaderResult[string, int] {
			return func(s string) (int, error) {
				return a + len(s), nil
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		value, err := result("")(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("enables partial application", func(t *testing.T) {
		original := Ask[int]()

		kleisli := func(a int) ReaderResult[int, int] {
			return func(factor int) (int, error) {
				return a * factor, nil
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Partially apply factor
		withFactor := result(3)

		// Can now use with different inputs
		value1, err1 := withFactor(10)
		assert.NoError(t, err1)
		assert.Equal(t, 30, value1)

		// Reuse with different input
		value2, err2 := withFactor(20)
		assert.NoError(t, err2)
		assert.Equal(t, 60, value2)
	})
}

func TestTraverseReader(t *testing.T) {
	t.Run("basic transformation with Reader dependency", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := F.Pipe1(
			Ask[int](),
			Map[int](N.Mul(2)),
		)

		// Reader-based transformation
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](multiply)
		result := traversed(original)

		// Provide Config first, then int
		cfg := Config{Multiplier: 5}
		value, err := result(cfg)(10)
		assert.NoError(t, err)
		assert.Equal(t, 100, value) // (10 * 2) * 5 = 100
	})

	t.Run("preserves outer error", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		expectedError := errors.New("outer error")

		// Original computation that fails
		original := func(x int) (int, error) {
			if x < 0 {
				return 0, expectedError
			}
			return x, nil
		}

		// Reader-based transformation (won't be called)
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](multiply)
		result := traversed(original)

		// Provide Config and negative value
		cfg := Config{Multiplier: 5}
		_, err := result(cfg)(-1)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		type Database struct {
			Prefix string
		}

		// Original computation producing an int
		original := Ask[int]()

		// Reader-based transformation: int -> string using Database
		format := func(a int) func(Database) string {
			return func(db Database) string {
				return fmt.Sprintf("%s:%d", db.Prefix, a)
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](format)
		result := traversed(original)

		// Provide Database first, then int
		db := Database{Prefix: "ID"}
		value, err := result(db)(42)
		assert.NoError(t, err)
		assert.Equal(t, "ID:42", value)
	})

	t.Run("works with struct environments", func(t *testing.T) {
		type Settings struct {
			Prefix string
			Suffix string
		}
		type Context struct {
			Value int
		}

		// Original computation
		original := func(ctx Context) (string, error) {
			return fmt.Sprintf("value:%d", ctx.Value), nil
		}

		// Reader-based transformation using Settings
		decorate := func(s string) func(Settings) string {
			return func(settings Settings) string {
				return settings.Prefix + s + settings.Suffix
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[Context](decorate)
		result := traversed(original)

		// Provide Settings first, then Context
		settings := Settings{Prefix: "[", Suffix: "]"}
		ctx := Context{Value: 100}
		value, err := result(settings)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, "[value:100]", value)
	})

	t.Run("enables partial application", func(t *testing.T) {
		type Config struct {
			Factor int
		}

		// Original computation
		original := Ask[int]()

		// Reader-based transformation
		scale := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Factor
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](scale)
		result := traversed(original)

		// Partially apply Config
		cfg := Config{Factor: 3}
		withConfig := result(cfg)

		// Can now use with different inputs
		value1, err1 := withConfig(10)
		assert.NoError(t, err1)
		assert.Equal(t, 30, value1)

		// Reuse with different input
		value2, err2 := withConfig(20)
		assert.NoError(t, err2)
		assert.Equal(t, 60, value2)
	})

	t.Run("works with zero values", func(t *testing.T) {
		type Config struct {
			Offset int
		}

		// Original computation with zero value
		original := Ask[int]()

		// Reader-based transformation
		add := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a + cfg.Offset
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](add)
		result := traversed(original)

		// Provide Config with zero offset and zero input
		cfg := Config{Offset: 0}
		value, err := result(cfg)(0)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := func(x int) (int, error) {
			return x * 2, nil
		}

		// Reader-based transformation
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](multiply)
		result := traversed(original)

		// Provide Config and execute
		cfg := Config{Multiplier: 4}
		value, err := result(cfg)(5)
		assert.NoError(t, err)
		assert.Equal(t, 40, value) // (5 * 2) * 4 = 40
	})

	t.Run("works with complex Reader logic", func(t *testing.T) {
		type ValidationRules struct {
			MinValue int
			MaxValue int
		}

		// Original computation
		original := Ask[int]()

		// Reader-based transformation with validation logic
		validate := func(a int) func(ValidationRules) int {
			return func(rules ValidationRules) int {
				if a < rules.MinValue {
					return rules.MinValue
				}
				if a > rules.MaxValue {
					return rules.MaxValue
				}
				return a
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int](validate)
		result := traversed(original)

		// Test with value within range
		rules1 := ValidationRules{MinValue: 0, MaxValue: 100}
		value1, err1 := result(rules1)(50)
		assert.NoError(t, err1)
		assert.Equal(t, 50, value1)

		// Test with value above max
		rules2 := ValidationRules{MinValue: 0, MaxValue: 30}
		value2, err2 := result(rules2)(50)
		assert.NoError(t, err2)
		assert.Equal(t, 30, value2) // Clamped to max

		// Test with value below min
		rules3 := ValidationRules{MinValue: 60, MaxValue: 100}
		value3, err3 := result(rules3)(50)
		assert.NoError(t, err3)
		assert.Equal(t, 60, value3) // Clamped to min
	})
}
