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

package readerioeither

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("sequences parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns IOEither[ReaderIOEither[string, error, int]]
		original := func(x int) IOEither[error, ReaderIOEither[string, error, int]] {
			return ioeither.Right[error](func(s string) IOEither[error, int] {
				return ioeither.Right[error](x + len(s))
			})
		}

		// Sequenced: takes string, returns ReaderIOEither[int, error, int]
		sequenced := Sequence(original)

		// Test original
		innerFunc1 := original(10)()
		assert.True(t, either.IsRight(innerFunc1))
		readerFunc1, _ := either.Unwrap(innerFunc1)
		result1 := readerFunc1("hello")()
		value1, err1 := either.Unwrap(result1)
		assert.NoError(t, err1)
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced("hello")
		result2 := innerFunc2(10)()
		value2, err2 := either.Unwrap(result2)
		assert.NoError(t, err2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) IOEither[error, ReaderIOEither[string, error, int]] {
			if x < 0 {
				return ioeither.Left[ReaderIOEither[string, error, int]](expectedError)
			}
			return ioeither.Right[error](func(s string) IOEither[error, int] {
				return ioeither.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test with error
		innerFunc := sequenced("test")
		result := innerFunc(-1)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		original := func(x int) IOEither[error, ReaderIOEither[string, error, int]] {
			return ioeither.Right[error](func(s string) IOEither[error, int] {
				if S.IsEmpty(s) {
					return ioeither.Left[int](expectedError)
				}
				return ioeither.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test with inner error
		innerFunc := sequenced("")
		result := innerFunc(10)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("sequences parameter order for Reader inner type", func(t *testing.T) {
		// Original: takes int, returns IOEither[Reader[string, int]]
		original := func(x int) IOEither[error, reader.Reader[string, int]] {
			return ioeither.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		// Sequenced: takes string, returns ReaderIOEither[int, error, int]
		sequenced := SequenceReader(original)

		// Test original
		readerFunc := original(10)()
		assert.True(t, either.IsRight(readerFunc))
		innerReader, _ := either.Unwrap(readerFunc)
		value1 := innerReader("hello")
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc := sequenced("hello")
		result := innerFunc(10)()
		value2, err := either.Unwrap(result)
		assert.NoError(t, err)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) IOEither[error, reader.Reader[string, int]] {
			if x < 0 {
				return ioeither.Left[reader.Reader[string, int]](expectedError)
			}
			return ioeither.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test with error
		innerFunc := sequenced("test")
		result := innerFunc(-1)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSequenceReaderIO(t *testing.T) {
	t.Run("sequences parameter order for ReaderIO inner type", func(t *testing.T) {
		// Original: takes int, returns IOEither[ReaderIO[string, int]]
		original := func(x int) IOEither[error, readerio.ReaderIO[string, int]] {
			return ioeither.Right[error](func(s string) IO[int] {
				return func() int {
					return x + len(s)
				}
			})
		}

		// Sequenced: takes string, returns ReaderIOEither[int, error, int]
		sequenced := SequenceReaderIO(original)

		// Test original
		readerIOFunc := original(10)()
		assert.True(t, either.IsRight(readerIOFunc))
		innerReaderIO, _ := either.Unwrap(readerIOFunc)
		value1 := innerReaderIO("hello")()
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc := sequenced("hello")
		result := innerFunc(10)()
		value2, err := either.Unwrap(result)
		assert.NoError(t, err)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) IOEither[error, readerio.ReaderIO[string, int]] {
			if x < 0 {
				return ioeither.Left[readerio.ReaderIO[string, int]](expectedError)
			}
			return ioeither.Right[error](func(s string) IO[int] {
				return func() int {
					return x + len(s)
				}
			})
		}

		sequenced := SequenceReaderIO(original)

		// Test with error
		innerFunc := sequenced("test")
		result := innerFunc(-1)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves IO effects", func(t *testing.T) {
		counter := 0

		original := func(x int) IOEither[error, readerio.ReaderIO[string, int]] {
			return ioeither.Right[error](func(s string) IO[int] {
				return func() int {
					counter++
					return x + len(s)
				}
			})
		}

		sequenced := SequenceReaderIO(original)

		// Execute multiple times to verify IO effects
		innerFunc := sequenced("test")
		innerFunc(10)()
		innerFunc(10)()
		assert.Equal(t, 2, counter)
	})
}

func TestSequenceReaderEither(t *testing.T) {
	t.Run("sequences parameter order for ReaderEither inner type", func(t *testing.T) {
		// Original: takes int, returns IOEither[ReaderEither[string, error, int]]
		original := func(x int) IOEither[error, ReaderEither[string, error, int]] {
			return ioeither.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x + len(s))
			})
		}

		// Sequenced: takes string, returns ReaderIOEither[int, error, int]
		sequenced := SequenceReaderEither(original)

		// Test original
		readerEitherFunc := original(10)()
		assert.True(t, either.IsRight(readerEitherFunc))
		innerReaderEither, _ := either.Unwrap(readerEitherFunc)
		result1 := innerReaderEither("hello")
		value1, err1 := either.Unwrap(result1)
		assert.NoError(t, err1)
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc := sequenced("hello")
		result2 := innerFunc(10)()
		value2, err2 := either.Unwrap(result2)
		assert.NoError(t, err2)
		assert.Equal(t, 15, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) IOEither[error, ReaderEither[string, error, int]] {
			if x < 0 {
				return ioeither.Left[ReaderEither[string, error, int]](expectedError)
			}
			return ioeither.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x + len(s))
			})
		}

		sequenced := SequenceReaderEither(original)

		// Test with error
		innerFunc := sequenced("test")
		result := innerFunc(-1)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		original := func(x int) IOEither[error, ReaderEither[string, error, int]] {
			return ioeither.Right[error](func(s string) either.Either[error, int] {
				if S.IsEmpty(s) {
					return either.Left[int](expectedError)
				}
				return either.Right[error](x + len(s))
			})
		}

		sequenced := SequenceReaderEither(original)

		// Test with inner error
		innerFunc := sequenced("")
		result := innerFunc(10)()
		_, err := either.Unwrap(result)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		original := func(cfg Config) IOEither[error, ReaderIOEither[Database, error, string]] {
			if cfg.Timeout <= 0 {
				return ioeither.Left[ReaderIOEither[Database, error, string]](errors.New("invalid timeout"))
			}
			return ioeither.Right[error](func(db Database) IOEither[error, string] {
				if S.IsEmpty(db.ConnectionString) {
					return ioeither.Left[string](errors.New("empty connection string"))
				}
				return ioeither.Right[error](fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout))
			})
		}

		sequenced := Sequence(original)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		innerFunc := sequenced(db)
		result := innerFunc(cfg)()
		value, err := either.Unwrap(result)
		assert.NoError(t, err)
		assert.Equal(t, "Query on localhost:5432 with timeout 30", value)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) IOEither[error, ReaderIOEither[string, error, int]] {
			return ioeither.Right[error](func(s string) IOEither[error, int] {
				return ioeither.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		innerFunc := sequenced("")
		result := innerFunc(0)()
		value, err := either.Unwrap(result)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})
}

func TestTraverse(t *testing.T) {
	t.Run("basic transformation with environment swap", func(t *testing.T) {
		// Original: ReaderIOEither[int, error, int] - takes int environment, produces int
		original := func(x int) IOEither[error, int] {
			return ioeither.Right[error](x * 2)
		}

		// Kleisli function: func(int) ReaderIOEither[string, error, int]
		kleisli := func(a int) ReaderIOEither[string, error, int] {
			return func(s string) IOEither[error, int] {
				return ioeither.Right[error](a + len(s))
			}
		}

		// Traverse returns: func(ReaderIOEither[int, error, int]) func(string) ReaderIOEither[int, error, int]
		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// result is func(string) ReaderIOEither[int, error, int]
		// Provide string first ("hello"), then int (10)
		innerFunc := result("hello")
		finalResult := innerFunc(10)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 25, value) // (10 * 2) + len("hello") = 20 + 5 = 25
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		original := func(x int) IOEither[error, int] {
			if x < 0 {
				return ioeither.Left[int](expectedError)
			}
			return ioeither.Right[error](x)
		}

		kleisli := func(a int) ReaderIOEither[string, error, int] {
			return func(s string) IOEither[error, int] {
				return ioeither.Right[error](a + len(s))
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Test with negative value to trigger error
		innerFunc := result("test")
		finalResult := innerFunc(-1)()
		_, err := either.Unwrap(finalResult)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("preserves inner error from Kleisli", func(t *testing.T) {
		expectedError := errors.New("inner error")

		original := Ask[int, error]()

		kleisli := func(a int) ReaderIOEither[string, error, int] {
			return func(s string) IOEither[error, int] {
				if S.IsEmpty(s) {
					return ioeither.Left[int](expectedError)
				}
				return ioeither.Right[error](a + len(s))
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Test with empty string to trigger inner error
		innerFunc := result("")
		finalResult := innerFunc(10)()
		_, err := either.Unwrap(finalResult)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		// Transform int to string using environment-dependent logic
		original := Ask[int, error]()

		kleisli := func(a int) ReaderIOEither[string, error, string] {
			return func(prefix string) IOEither[error, string] {
				return ioeither.Right[error](fmt.Sprintf("%s-%d", prefix, a))
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		innerFunc := result("ID")
		finalResult := innerFunc(42)()
		value, err := either.Unwrap(finalResult)
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

		original := func(cfg Config) IOEither[error, int] {
			if cfg.Multiplier <= 0 {
				return ioeither.Left[int](errors.New("invalid multiplier"))
			}
			return ioeither.Right[error](10 * cfg.Multiplier)
		}

		kleisli := func(value int) ReaderIOEither[Database, error, string] {
			return func(db Database) IOEither[error, string] {
				return ioeither.Right[error](fmt.Sprintf("%s:%d", db.Prefix, value))
			}
		}

		traversed := Traverse[Config](kleisli)
		result := traversed(original)

		cfg := Config{Multiplier: 5}
		db := Database{Prefix: "result"}

		innerFunc := result(db)
		finalResult := innerFunc(cfg)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, "result:50", value)
	})

	t.Run("preserves IO effects", func(t *testing.T) {
		outerCounter := 0
		innerCounter := 0

		original := func(x int) IOEither[error, int] {
			return func() either.Either[error, int] {
				outerCounter++
				return either.Right[error](x)
			}
		}

		kleisli := func(a int) ReaderIOEither[string, error, int] {
			return func(s string) IOEither[error, int] {
				return func() either.Either[error, int] {
					innerCounter++
					return either.Right[error](a + len(s))
				}
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		// Execute multiple times to verify IO effects
		innerFunc := result("test")
		innerFunc(10)()
		innerFunc(10)()

		assert.Equal(t, 2, outerCounter)
		assert.Equal(t, 2, innerCounter)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		original := Ask[int, error]()

		// First transformation: multiply by environment value
		kleisli1 := func(a int) ReaderIOEither[int, error, int] {
			return func(multiplier int) IOEither[error, int] {
				return ioeither.Right[error](a * multiplier)
			}
		}

		traversed := Traverse[int](kleisli1)
		result := traversed(original)

		innerFunc := result(3)
		finalResult := innerFunc(5)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 15, value) // 5 * 3 = 15
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := Ask[int, error]()

		kleisli := func(a int) ReaderIOEither[string, error, int] {
			return func(s string) IOEither[error, int] {
				return ioeither.Right[error](a + len(s))
			}
		}

		traversed := Traverse[int](kleisli)
		result := traversed(original)

		innerFunc := result("")
		finalResult := innerFunc(0)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})
}

func TestTraverseReader(t *testing.T) {
	t.Run("basic transformation with Reader dependency", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := F.Pipe1(
			Ask[int, error](),
			Map[int, error](N.Mul(2)),
		)

		// Reader-based transformation
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](multiply)
		result := traversed(original)

		// Provide Config first, then int
		cfg := Config{Multiplier: 5}
		innerFunc := result(cfg)
		finalResult := innerFunc(10)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 100, value) // (10 * 2) * 5 = 100
	})

	t.Run("preserves outer error", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		expectedError := errors.New("outer error")

		// Original computation that fails
		original := func(x int) IOEither[error, int] {
			if x < 0 {
				return ioeither.Left[int](expectedError)
			}
			return ioeither.Right[error](x)
		}

		// Reader-based transformation (won't be called)
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](multiply)
		result := traversed(original)

		// Provide Config and negative value
		cfg := Config{Multiplier: 5}
		innerFunc := result(cfg)
		finalResult := innerFunc(-1)()
		_, err := either.Unwrap(finalResult)
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different types", func(t *testing.T) {
		type Database struct {
			Prefix string
		}

		// Original computation producing an int
		original := Ask[int, error]()

		// Reader-based transformation: int -> string using Database
		format := func(a int) func(Database) string {
			return func(db Database) string {
				return fmt.Sprintf("%s:%d", db.Prefix, a)
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Database, error](format)
		result := traversed(original)

		// Provide Database first, then int
		db := Database{Prefix: "ID"}
		innerFunc := result(db)
		finalResult := innerFunc(42)()
		value, err := either.Unwrap(finalResult)
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
		original := func(ctx Context) IOEither[error, string] {
			return ioeither.Right[error](fmt.Sprintf("value:%d", ctx.Value))
		}

		// Reader-based transformation using Settings
		decorate := func(s string) func(Settings) string {
			return func(settings Settings) string {
				return settings.Prefix + s + settings.Suffix
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[Context, Settings, error](decorate)
		result := traversed(original)

		// Provide Settings first, then Context
		settings := Settings{Prefix: "[", Suffix: "]"}
		ctx := Context{Value: 100}
		innerFunc := result(settings)
		finalResult := innerFunc(ctx)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, "[value:100]", value)
	})

	t.Run("enables partial application", func(t *testing.T) {
		type Config struct {
			Factor int
		}

		// Original computation
		original := Ask[int, error]()

		// Reader-based transformation
		scale := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Factor
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](scale)
		result := traversed(original)

		// Partially apply Config
		cfg := Config{Factor: 3}
		withConfig := result(cfg)

		// Can now use with different inputs
		finalResult1 := withConfig(10)()
		value1, err1 := either.Unwrap(finalResult1)
		assert.NoError(t, err1)
		assert.Equal(t, 30, value1)

		// Reuse with different input
		finalResult2 := withConfig(20)()
		value2, err2 := either.Unwrap(finalResult2)
		assert.NoError(t, err2)
		assert.Equal(t, 60, value2)
	})

	t.Run("preserves IO effects", func(t *testing.T) {
		type Config struct {
			Value int
		}

		outerCounter := 0
		innerCounter := 0

		// Original computation with IO effects
		original := func(x int) IOEither[error, int] {
			return func() either.Either[error, int] {
				outerCounter++
				return either.Right[error](x)
			}
		}

		// Reader-based transformation
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				innerCounter++
				return a * cfg.Value
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](multiply)
		result := traversed(original)

		// Execute multiple times to verify IO effects
		cfg := Config{Value: 5}
		innerFunc := result(cfg)
		innerFunc(10)()
		innerFunc(10)()

		assert.Equal(t, 2, outerCounter)
		assert.Equal(t, 2, innerCounter)
	})

	t.Run("works with zero values", func(t *testing.T) {
		type Config struct {
			Offset int
		}

		// Original computation with zero value
		original := Ask[int, error]()

		// Reader-based transformation
		add := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a + cfg.Offset
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](add)
		result := traversed(original)

		// Provide Config with zero offset and zero input
		cfg := Config{Offset: 0}
		innerFunc := result(cfg)
		finalResult := innerFunc(0)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 0, value)
	})

	t.Run("chains multiple transformations", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original computation
		original := func(x int) IOEither[error, int] {
			return ioeither.Right[error](x * 2)
		}

		// Reader-based transformation
		multiply := func(a int) func(Config) int {
			return func(cfg Config) int {
				return a * cfg.Multiplier
			}
		}

		// Apply TraverseReader
		traversed := TraverseReader[int, Config, error](multiply)
		result := traversed(original)

		// Provide Config and execute
		cfg := Config{Multiplier: 4}
		innerFunc := result(cfg)
		finalResult := innerFunc(5)()
		value, err := either.Unwrap(finalResult)
		assert.NoError(t, err)
		assert.Equal(t, 40, value) // (5 * 2) * 4 = 40
	})

	t.Run("works with complex Reader logic", func(t *testing.T) {
		type ValidationRules struct {
			MinValue int
			MaxValue int
		}

		// Original computation
		original := Ask[int, error]()

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
		traversed := TraverseReader[int, ValidationRules, error](validate)
		result := traversed(original)

		// Test with value within range
		rules1 := ValidationRules{MinValue: 0, MaxValue: 100}
		innerFunc1 := result(rules1)
		finalResult1 := innerFunc1(50)()
		value1, err1 := either.Unwrap(finalResult1)
		assert.NoError(t, err1)
		assert.Equal(t, 50, value1)

		// Test with value above max
		rules2 := ValidationRules{MinValue: 0, MaxValue: 30}
		innerFunc2 := result(rules2)
		finalResult2 := innerFunc2(50)()
		value2, err2 := either.Unwrap(finalResult2)
		assert.NoError(t, err2)
		assert.Equal(t, 30, value2) // Clamped to max

		// Test with value below min
		rules3 := ValidationRules{MinValue: 60, MaxValue: 100}
		innerFunc3 := result(rules3)
		finalResult3 := innerFunc3(50)()
		value3, err3 := either.Unwrap(finalResult3)
		assert.NoError(t, err3)
		assert.Equal(t, 60, value3) // Clamped to min
	})
}
