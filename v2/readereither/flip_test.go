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

package readereither

import (
	"errors"
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Either[error, ReaderEither[string, error, int]]
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			return either.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x + len(s))
			})
		}

		// Sequenceped: takes string, returns Reader[int, Either[error, int]]
		sequenced := Sequence(original)

		// Test original
		result1 := original(10)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1("hello")
		assert.True(t, either.IsRight(innerResult1))
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, 15, value1)

		// Test sequenced - note it returns a Reader directly
		innerFunc2 := sequenced("hello")
		innerResult2 := innerFunc2(10)
		assert.True(t, either.IsRight(innerResult2))
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, 15, value2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Config, returns Either[error, ReaderEither[Database, error, string]]
		query := func(cfg Config) either.Either[error, ReaderEither[Database, error, string]] {
			if cfg.Timeout <= 0 {
				return either.Left[ReaderEither[Database, error, string]](errors.New("invalid timeout"))
			}
			return either.Right[error](func(db Database) either.Either[error, string] {
				if S.IsEmpty(db.ConnectionString) {
					return either.Left[string](errors.New("empty connection string"))
				}
				return either.Right[error](fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout))
			})
		}

		// Sequenceped: takes Database, returns Reader[Config, Either[error, string]]
		sequenced := Sequence(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original with valid inputs
		result1 := query(cfg)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(db)
		assert.True(t, either.IsRight(innerResult1))
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, expected, value1)

		// Test sequenced with valid inputs
		innerFunc2 := sequenced(db)
		innerResult2 := innerFunc2(cfg)
		assert.True(t, either.IsRight(innerResult2))
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, expected, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			if x < 0 {
				return either.Left[ReaderEither[string, error, int]](expectedError)
			}
			return either.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with error
		result1 := original(-1)
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error becomes an inner error after flip
		innerFunc2 := sequenced("test")
		innerResult2 := innerFunc2(-1)
		assert.True(t, either.IsLeft(innerResult2))
		_, innerErr2 := either.Unwrap(innerResult2)
		assert.Equal(t, expectedError, innerErr2)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		// Original that fails at inner level
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			return either.Right[error](func(s string) either.Either[error, int] {
				if S.IsEmpty(s) {
					return either.Left[int](expectedError)
				}
				return either.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with inner error
		result1 := original(10)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1("")
		assert.True(t, either.IsLeft(innerResult1))
		_, innerErr1 := either.Unwrap(innerResult1)
		assert.Equal(t, expectedError, innerErr1)

		// Test sequenced with inner error
		innerFunc2 := sequenced("")
		innerResult2 := innerFunc2(10)
		assert.True(t, either.IsLeft(innerResult2))
		_, innerErr2 := either.Unwrap(innerResult2)
		assert.Equal(t, expectedError, innerErr2)
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original function
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			return either.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x * len(s))
			})
		}

		// Sequence once
		sequenced := Sequence(original)

		// Test that flip produces correct results
		result1 := original(3)
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1("test")
		value1, _ := either.Unwrap(innerResult1)

		innerFunc2 := sequenced("test")
		innerResult2 := innerFunc2(3)
		value2, _ := either.Unwrap(innerResult2)

		assert.Equal(t, value1, value2)
		assert.Equal(t, 12, value2) // 3 * 4
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			return either.Right[error](func(s string) either.Either[error, int] {
				return either.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test with zero values
		result1 := original(0)
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1("")
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, 0, value1)

		innerFunc2 := sequenced("")
		innerResult2 := innerFunc2(0)
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, 0, value2)
	})

	t.Run("works with complex computations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		type Session struct {
			Token string
		}

		// Original: takes User, returns Either[error, ReaderEither[Session, error, string]]
		authenticate := func(user User) either.Either[error, ReaderEither[Session, error, string]] {
			if user.ID <= 0 {
				return either.Left[ReaderEither[Session, error, string]](errors.New("invalid user ID"))
			}
			return either.Right[error](func(session Session) either.Either[error, string] {
				if S.IsEmpty(session.Token) {
					return either.Left[string](errors.New("empty token"))
				}
				return either.Right[error](fmt.Sprintf("User %s (ID: %d) authenticated with token %s",
					user.Name, user.ID, session.Token))
			})
		}

		sequenced := Sequence(authenticate)

		user := User{ID: 42, Name: "Alice"}
		session := Session{Token: "abc123"}

		expected := "User Alice (ID: 42) authenticated with token abc123"

		// Test original
		result1 := authenticate(user)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(session)
		assert.True(t, either.IsRight(innerResult1))
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, expected, value1)

		// Test sequenced
		innerFunc2 := sequenced(session)
		innerResult2 := innerFunc2(user)
		assert.True(t, either.IsRight(innerResult2))
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, expected, value2)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) either.Either[error, ReaderEither[Empty, error, int]] {
			return either.Right[error](func(e Empty) either.Either[error, int] {
				return either.Right[error](x * 2)
			})
		}

		sequenced := Sequence(original)

		empty := Empty{}

		result1 := original(10)
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(empty)
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, 20, value1)

		innerFunc2 := sequenced(empty)
		innerResult2 := innerFunc2(10)
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) either.Either[error, ReaderEither[*Data, error, int]] {
			return either.Right[error](func(d *Data) either.Either[error, int] {
				if d == nil {
					return either.Right[error](x)
				}
				return either.Right[error](x + d.Value)
			})
		}

		sequenced := Sequence(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		result1 := original(42)
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(data)
		value1, _ := either.Unwrap(innerResult1)
		assert.Equal(t, 142, value1)

		innerFunc2 := sequenced(data)
		innerResult2 := innerFunc2(42)
		value2, _ := either.Unwrap(innerResult2)
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3 := original(42)
		innerFunc3, _ := either.Unwrap(result3)
		innerResult3 := innerFunc3(nil)
		value3, _ := either.Unwrap(innerResult3)
		assert.Equal(t, 42, value3)

		innerFunc4 := sequenced(nil)
		innerResult4 := innerFunc4(42)
		value4, _ := either.Unwrap(innerResult4)
		assert.Equal(t, 42, value4)
	})

	t.Run("handles errors at different levels", func(t *testing.T) {
		// Original that can fail at both levels
		original := func(x int) either.Either[error, ReaderEither[string, error, int]] {
			if x < 0 {
				return either.Left[ReaderEither[string, error, int]](errors.New("outer: negative value"))
			}
			return either.Right[error](func(s string) either.Either[error, int] {
				if S.IsEmpty(s) {
					return either.Left[int](errors.New("inner: empty string"))
				}
				return either.Right[error](x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test outer error
		result1 := original(-1)
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Contains(t, err1.Error(), "outer")

		// Test inner error
		result2 := original(10)
		assert.True(t, either.IsRight(result2))
		innerFunc2, _ := either.Unwrap(result2)
		innerResult2 := innerFunc2("")
		assert.True(t, either.IsLeft(innerResult2))
		_, innerErr2 := either.Unwrap(innerResult2)
		assert.Contains(t, innerErr2.Error(), "inner")

		// Test sequenced with inner error
		innerFunc3 := sequenced("")
		innerResult3 := innerFunc3(10)
		assert.True(t, either.IsLeft(innerResult3))
		_, innerErr3 := either.Unwrap(innerResult3)
		assert.Contains(t, innerErr3.Error(), "inner")
	})

	t.Run("validates error propagation with outer failure", func(t *testing.T) {
		type Config struct {
			Valid bool
		}

		original := func(cfg Config) either.Either[error, ReaderEither[int, error, string]] {
			if !cfg.Valid {
				return either.Left[ReaderEither[int, error, string]](errors.New("invalid config"))
			}
			return either.Right[error](func(x int) either.Either[error, string] {
				return either.Right[error](fmt.Sprintf("value: %d", x))
			})
		}

		sequenced := Sequence(original)

		invalidCfg := Config{Valid: false}

		// Original with invalid config
		result1 := original(invalidCfg)
		assert.True(t, either.IsLeft(result1))

		// Sequenceped returns a Reader, but when called with invalid config, it returns Left
		innerFunc2 := sequenced(42)
		innerResult2 := innerFunc2(invalidCfg)
		assert.True(t, either.IsLeft(innerResult2))
	})

	t.Run("validates error propagation with inner failure", func(t *testing.T) {
		original := func(threshold int) either.Either[error, ReaderEither[int, error, string]] {
			return either.Right[error](func(value int) either.Either[error, string] {
				if value < threshold {
					return either.Left[string](fmt.Errorf("value %d below threshold %d", value, threshold))
				}
				return either.Right[error](fmt.Sprintf("valid: %d", value))
			})
		}

		sequenced := Sequence(original)

		// Test with value below threshold
		result1 := original(10)
		innerFunc1, _ := either.Unwrap(result1)
		innerResult1 := innerFunc1(5)
		assert.True(t, either.IsLeft(innerResult1))

		// Sequenceped should also fail
		innerFunc2 := sequenced(5)
		innerResult2 := innerFunc2(10)
		assert.True(t, either.IsLeft(innerResult2))
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Either[error, Reader[string, int]]
		original := func(x int) either.Either[error, Reader[string, int]] {
			return either.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		// Sequenceped: takes string, returns Reader[int, Either[error, int]]
		sequenced := SequenceReader(original)

		// Test original
		result1 := original(10)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1("hello")
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced("hello")
		result2 := innerFunc2(10)
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Config, returns Either[error, Reader[Database, string]]
		query := func(cfg Config) either.Either[error, Reader[Database, string]] {
			if cfg.Timeout <= 0 {
				return either.Left[Reader[Database, string]](errors.New("invalid timeout"))
			}
			return either.Right[error](func(db Database) string {
				return fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout)
			})
		}

		// Sequenceped: takes Database, returns Reader[Config, Either[error, string]]
		sequenced := SequenceReader(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original with valid inputs
		result1 := query(cfg)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(db)
		assert.Equal(t, expected, value1)

		// Test sequenced with valid inputs
		innerFunc2 := sequenced(db)
		result2 := innerFunc2(cfg)
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, expected, value2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(x int) either.Either[error, Reader[string, int]] {
			if x < 0 {
				return either.Left[Reader[string, int]](expectedError)
			}
			return either.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test original with error
		result1 := original(-1)
		assert.True(t, either.IsLeft(result1))
		_, err1 := either.Unwrap(result1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error becomes an inner error after flip
		innerFunc2 := sequenced("test")
		result2 := innerFunc2(-1)
		assert.True(t, either.IsLeft(result2))
		_, err2 := either.Unwrap(result2)
		assert.Equal(t, expectedError, err2)
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original: takes bool, returns Either[error, Reader[int, string]]
		original := func(flag bool) either.Either[error, Reader[int, string]] {
			if !flag {
				return either.Left[Reader[int, string]](errors.New("flag is false"))
			}
			return either.Right[error](func(n int) string {
				if flag {
					return fmt.Sprintf("positive: %d", n)
				}
				return fmt.Sprintf("negative: %d", -n)
			})
		}

		sequenced := SequenceReader(original)

		// Test with true flag
		result1 := original(true)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(42)
		assert.Equal(t, "positive: 42", value1)

		innerFunc2 := sequenced(42)
		result2 := innerFunc2(true)
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, "positive: 42", value2)

		// Test with false flag
		result3 := original(false)
		assert.True(t, either.IsLeft(result3))

		innerFunc4 := sequenced(42)
		result4 := innerFunc4(false)
		assert.True(t, either.IsLeft(result4))
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) either.Either[error, Reader[string, int]] {
			return either.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test with zero values
		result1 := original(0)
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1("")
		assert.Equal(t, 0, value1)

		innerFunc2 := sequenced("")
		result2 := innerFunc2(0)
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 0, value2)
	})

	t.Run("works with complex computations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		type Session struct {
			Token string
		}

		// Original: takes User, returns Either[error, Reader[Session, string]]
		authenticate := func(user User) either.Either[error, Reader[Session, string]] {
			if user.ID <= 0 {
				return either.Left[Reader[Session, string]](errors.New("invalid user ID"))
			}
			return either.Right[error](func(session Session) string {
				return fmt.Sprintf("User %s (ID: %d) authenticated with token %s",
					user.Name, user.ID, session.Token)
			})
		}

		sequenced := SequenceReader(authenticate)

		user := User{ID: 42, Name: "Alice"}
		session := Session{Token: "abc123"}

		expected := "User Alice (ID: 42) authenticated with token abc123"

		// Test original
		result1 := authenticate(user)
		assert.True(t, either.IsRight(result1))
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(session)
		assert.Equal(t, expected, value1)

		// Test sequenced
		innerFunc2 := sequenced(session)
		result2 := innerFunc2(user)
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, expected, value2)
	})
}

func TestSequenceReaderEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) either.Either[error, Reader[Empty, int]] {
			return either.Right[error](func(e Empty) int {
				return x * 2
			})
		}

		sequenced := SequenceReader(original)

		empty := Empty{}

		result1 := original(10)
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(empty)
		assert.Equal(t, 20, value1)

		innerFunc2 := sequenced(empty)
		result2 := innerFunc2(10)
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) either.Either[error, Reader[*Data, int]] {
			return either.Right[error](func(d *Data) int {
				if d == nil {
					return x
				}
				return x + d.Value
			})
		}

		sequenced := SequenceReader(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		result1 := original(42)
		innerFunc1, _ := either.Unwrap(result1)
		value1 := innerFunc1(data)
		assert.Equal(t, 142, value1)

		innerFunc2 := sequenced(data)
		result2 := innerFunc2(42)
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3 := original(42)
		innerFunc3, _ := either.Unwrap(result3)
		value3 := innerFunc3(nil)
		assert.Equal(t, 42, value3)

		innerFunc4 := sequenced(nil)
		result4 := innerFunc4(42)
		value4, _ := either.Unwrap(result4)
		assert.Equal(t, 42, value4)
	})

	t.Run("validates error propagation with outer failure", func(t *testing.T) {
		type Config struct {
			Valid bool
		}

		original := func(cfg Config) either.Either[error, Reader[int, string]] {
			if !cfg.Valid {
				return either.Left[Reader[int, string]](errors.New("invalid config"))
			}
			return either.Right[error](func(x int) string {
				return fmt.Sprintf("value: %d", x)
			})
		}

		sequenced := SequenceReader(original)

		invalidCfg := Config{Valid: false}

		// Original with invalid config
		result1 := original(invalidCfg)
		assert.True(t, either.IsLeft(result1))

		// Sequenceped returns a Reader, but when called with invalid config, it returns Left
		innerFunc2 := sequenced(42)
		result2 := innerFunc2(invalidCfg)
		assert.True(t, either.IsLeft(result2))
	})

	t.Run("can be used in composition", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original function
		multiply := func(x int) either.Either[error, Reader[Config, int]] {
			if x < 0 {
				return either.Left[Reader[Config, int]](errors.New("negative value"))
			}
			return either.Right[error](func(c Config) int {
				return x * c.Multiplier
			})
		}

		// Sequence it
		sequenced := SequenceReader(multiply)

		// Use sequenced version to partially apply config first
		config := Config{Multiplier: 10}
		multiplyBy10 := sequenced(config)

		// Now we have a function int -> Either[error, int]
		result1 := multiplyBy10(5)
		assert.True(t, either.IsRight(result1))
		value1, _ := either.Unwrap(result1)
		assert.Equal(t, 50, value1)

		result2 := multiplyBy10(10)
		assert.True(t, either.IsRight(result2))
		value2, _ := either.Unwrap(result2)
		assert.Equal(t, 100, value2)

		// Test with negative value
		result3 := multiplyBy10(-5)
		assert.True(t, either.IsLeft(result3))
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		// The same inputs should always produce the same outputs
		original := func(x int) either.Either[error, Reader[string, int]] {
			return either.Right[error](func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Call multiple times with same inputs
		for range 5 {
			result1 := original(10)
			innerFunc1, _ := either.Unwrap(result1)
			value1 := innerFunc1("hello")
			assert.Equal(t, 15, value1)

			innerFunc2 := sequenced("hello")
			result2 := innerFunc2(10)
			value2, _ := either.Unwrap(result2)
			assert.Equal(t, 15, value2)
		}
	})
}
