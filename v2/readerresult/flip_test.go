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

	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Result[ReaderResult[string, int]]
		original := func(x int) result.Result[ReaderResult[string, int]] {
			return result.Right(func(s string) result.Result[int] {
				return result.Right(x + len(s))
			})
		}

		// Sequenceped: takes string, returns Result[ReaderResult[int, int]]
		sequenced := Sequence(original)

		// Test original
		result1, err1 := result.Unwrap(original(10))
		assert.NoError(t, err1)
		innerResult1, innerErr1 := result.Unwrap(result1("hello"))
		assert.NoError(t, innerErr1)
		assert.Equal(t, 15, innerResult1)

		// Test sequenced - sequenced returns func(string) func(int) Result[int]
		innerFunc := sequenced("hello")
		innerResult2, innerErr2 := result.Unwrap(innerFunc(10))
		assert.NoError(t, innerErr2)
		assert.Equal(t, 15, innerResult2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Config, returns Result[ReaderResult[Database, string]]
		query := func(cfg Config) result.Result[ReaderResult[Database, string]] {
			if cfg.Timeout <= 0 {
				return result.Left[ReaderResult[Database, string]](errors.New("invalid timeout"))
			}
			return result.Right(func(db Database) result.Result[string] {
				if S.IsEmpty(db.ConnectionString) {
					return result.Left[string](errors.New("empty connection string"))
				}
				return result.Right(fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout))
			})
		}

		// Sequenceped: takes Database, returns Result[ReaderResult[Config, string]]
		sequenced := Sequence(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original with valid inputs
		result1, err1 := result.Unwrap(query(cfg))
		assert.NoError(t, err1)
		innerResult1, innerErr1 := result.Unwrap(result1(db))
		assert.NoError(t, innerErr1)
		assert.Equal(t, expected, innerResult1)

		// Test sequenced with valid inputs - sequenced returns func(Database) func(Config) Result[string]
		innerFunc2 := sequenced(db)
		innerResult2, innerErr2 := result.Unwrap(innerFunc2(cfg))
		assert.NoError(t, innerErr2)
		assert.Equal(t, expected, innerResult2)
	})

	t.Run("preserves outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")

		// Original that fails at outer level
		original := func(x int) result.Result[ReaderResult[string, int]] {
			if x < 0 {
				return result.Left[ReaderResult[string, int]](expectedError)
			}
			return result.Right(func(s string) result.Result[int] {
				return result.Right(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with error
		_, err1 := result.Unwrap(original(-1))
		assert.Error(t, err1)
		assert.Equal(t, expectedError, err1)

		// Test sequenced - the outer error becomes an inner error after flip
		// sequenced returns func(string) func(int) Result[int]
		innerFunc := sequenced("test")
		_, innerErr2 := result.Unwrap(innerFunc(-1))
		assert.Error(t, innerErr2)
		assert.Equal(t, expectedError, innerErr2)
	})

	t.Run("preserves inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")

		// Original that fails at inner level
		original := func(x int) result.Result[ReaderResult[string, int]] {
			return result.Right(func(s string) result.Result[int] {
				if S.IsEmpty(s) {
					return result.Left[int](expectedError)
				}
				return result.Right(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with inner error
		result1, err1 := result.Unwrap(original(10))
		assert.NoError(t, err1)
		_, innerErr1 := result.Unwrap(result1(""))
		assert.Error(t, innerErr1)
		assert.Equal(t, expectedError, innerErr1)

		// Test sequenced with inner error - sequenced returns func(string) func(int) Result[int]
		innerFunc := sequenced("")
		_, innerErr2 := result.Unwrap(innerFunc(10))
		assert.Error(t, innerErr2)
		assert.Equal(t, expectedError, innerErr2)
	})

	t.Run("works with multiple flips", func(t *testing.T) {
		// Original function
		original := func(x int) result.Result[ReaderResult[string, int]] {
			return result.Right(func(s string) result.Result[int] {
				return result.Right(x * len(s))
			})
		}

		// Sequence once
		sequenced1 := Sequence(original)

		// Sequence twice (should be equivalent to original)
		// Note: sequenced1 has type func(string) func(int) Result[int]
		// To flip it again, we need to wrap it back into the expected form
		sequenced2 := func(s string) result.Result[ReaderResult[int, int]] {
			return result.Right(sequenced1(s))
		}
		sequenced2Sequenceped := Sequence(sequenced2)

		// Test that double flip returns to original behavior
		result1, _ := result.Unwrap(original(3))
		innerResult1, _ := result.Unwrap(result1("test"))

		innerFunc := sequenced2Sequenceped(3)
		innerResult2, _ := result.Unwrap(innerFunc("test"))

		assert.Equal(t, innerResult1, innerResult2)
		assert.Equal(t, 12, innerResult2) // 3 * 4
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) result.Result[ReaderResult[string, int]] {
			return result.Right(func(s string) result.Result[int] {
				return result.Right(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test with zero values
		result1, _ := result.Unwrap(original(0))
		value1, _ := result.Unwrap(result1(""))
		assert.Equal(t, 0, value1)

		innerFunc := sequenced("")
		value2, _ := result.Unwrap(innerFunc(0))
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

		// Original: takes User, returns Result[ReaderResult[Session, string]]
		authenticate := func(user User) result.Result[ReaderResult[Session, string]] {
			if user.ID <= 0 {
				return result.Left[ReaderResult[Session, string]](errors.New("invalid user ID"))
			}
			return result.Right(func(session Session) result.Result[string] {
				if S.IsEmpty(session.Token) {
					return result.Left[string](errors.New("empty token"))
				}
				return result.Right(fmt.Sprintf("User %s (ID: %d) authenticated with token %s",
					user.Name, user.ID, session.Token))
			})
		}

		sequenced := Sequence(authenticate)

		user := User{ID: 42, Name: "Alice"}
		session := Session{Token: "abc123"}

		expected := "User Alice (ID: 42) authenticated with token abc123"

		// Test original
		result1, err1 := result.Unwrap(authenticate(user))
		assert.NoError(t, err1)
		value1, innerErr1 := result.Unwrap(result1(session))
		assert.NoError(t, innerErr1)
		assert.Equal(t, expected, value1)

		// Test sequenced - sequenced returns func(Session) func(User) Result[string]
		innerFunc := sequenced(session)
		value2, innerErr2 := result.Unwrap(innerFunc(user))
		assert.NoError(t, innerErr2)
		assert.Equal(t, expected, value2)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) result.Result[ReaderResult[Empty, int]] {
			return result.Right(func(e Empty) result.Result[int] {
				return result.Right(x * 2)
			})
		}

		sequenced := Sequence(original)

		empty := Empty{}

		result1, _ := result.Unwrap(original(10))
		value1, _ := result.Unwrap(result1(empty))
		assert.Equal(t, 20, value1)

		innerFunc := sequenced(empty)
		value2, _ := result.Unwrap(innerFunc(10))
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) result.Result[ReaderResult[*Data, int]] {
			return result.Right(func(d *Data) result.Result[int] {
				if d == nil {
					return result.Right(x)
				}
				return result.Right(x + d.Value)
			})
		}

		sequenced := Sequence(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		result1, _ := result.Unwrap(original(42))
		value1, _ := result.Unwrap(result1(data))
		assert.Equal(t, 142, value1)

		innerFunc2 := sequenced(data)
		value2, _ := result.Unwrap(innerFunc2(42))
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3, _ := result.Unwrap(original(42))
		value3, _ := result.Unwrap(result3(nil))
		assert.Equal(t, 42, value3)

		innerFunc3 := sequenced(nil)
		value4, _ := result.Unwrap(innerFunc3(42))
		assert.Equal(t, 42, value4)
	})

	t.Run("handles errors at different levels", func(t *testing.T) {
		// Original that can fail at both levels
		original := func(x int) result.Result[ReaderResult[string, int]] {
			if x < 0 {
				return result.Left[ReaderResult[string, int]](errors.New("outer: negative value"))
			}
			return result.Right(func(s string) result.Result[int] {
				if S.IsEmpty(s) {
					return result.Left[int](errors.New("inner: empty string"))
				}
				return result.Right(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test outer error
		_, err1 := result.Unwrap(original(-1))
		assert.Error(t, err1)
		assert.Contains(t, err1.Error(), "outer")

		// Test inner error
		result2, err2 := result.Unwrap(original(10))
		assert.NoError(t, err2)
		_, innerErr2 := result.Unwrap(result2(""))
		assert.Error(t, innerErr2)
		assert.Contains(t, innerErr2.Error(), "inner")

		// Test sequenced with inner error - sequenced returns func(string) func(int) Result[int]
		innerFunc := sequenced("")
		_, innerErr3 := result.Unwrap(innerFunc(10))
		assert.Error(t, innerErr3)
		assert.Contains(t, innerErr3.Error(), "inner")
	})

	t.Run("validates error propagation with outer failure", func(t *testing.T) {
		type Config struct {
			Valid bool
		}

		original := func(cfg Config) result.Result[ReaderResult[int, string]] {
			if !cfg.Valid {
				return result.Left[ReaderResult[int, string]](errors.New("invalid config"))
			}
			return result.Right(func(x int) result.Result[string] {
				return result.Right(fmt.Sprintf("value: %d", x))
			})
		}

		sequenced := Sequence(original)

		invalidCfg := Config{Valid: false}

		// Original with invalid config
		_, err1 := result.Unwrap(original(invalidCfg))
		assert.Error(t, err1)

		// Sequenceped returns success at outer level, but fails at inner level
		// sequenced returns func(int) func(Config) Result[string]
		innerFunc := sequenced(42)
		_, innerErr2 := result.Unwrap(innerFunc(invalidCfg))
		assert.Error(t, innerErr2)
	})

	t.Run("validates error propagation with inner failure", func(t *testing.T) {
		original := func(threshold int) result.Result[ReaderResult[int, string]] {
			return result.Right(func(value int) result.Result[string] {
				if value < threshold {
					return result.Left[string](fmt.Errorf("value %d below threshold %d", value, threshold))
				}
				return result.Right(fmt.Sprintf("valid: %d", value))
			})
		}

		sequenced := Sequence(original)

		// Test with value below threshold
		result1, err1 := result.Unwrap(original(10))
		assert.NoError(t, err1)
		_, innerErr1 := result.Unwrap(result1(5))
		assert.Error(t, innerErr1)

		// Sequenceped should also fail - sequenced returns func(int) func(int) Result[string]
		innerFunc := sequenced(5)
		_, innerErr2 := result.Unwrap(innerFunc(10))
		assert.Error(t, innerErr2)
	})

	t.Run("works with idiomatic Go error handling", func(t *testing.T) {
		// Demonstrates using Result with standard Go error patterns
		original := func(divisor int) result.Result[ReaderResult[int, int]] {
			if divisor == 0 {
				return result.Left[ReaderResult[int, int]](errors.New("division by zero"))
			}
			return result.Right(func(dividend int) result.Result[int] {
				return result.Right(dividend / divisor)
			})
		}

		sequenced := Sequence(original)

		// Test successful division
		result1, err1 := result.Unwrap(original(2))
		assert.NoError(t, err1)
		value1, innerErr1 := result.Unwrap(result1(10))
		assert.NoError(t, innerErr1)
		assert.Equal(t, 5, value1)

		// Test sequenced successful division - sequenced returns func(int) func(int) Result[int]
		innerFunc := sequenced(10)
		value2, innerErr2 := result.Unwrap(innerFunc(2))
		assert.NoError(t, innerErr2)
		assert.Equal(t, 5, value2)

		// Test division by zero
		_, err3 := result.Unwrap(original(0))
		assert.Error(t, err3)
		assert.Equal(t, "division by zero", err3.Error())
	})
}
