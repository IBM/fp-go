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

package readeroption

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestSequence(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Option[ReaderOption[string, int]]
		original := func(x int) option.Option[ReaderOption[string, int]] {
			return option.Some(func(s string) option.Option[int] {
				return option.Some(x + len(s))
			})
		}

		// Sequenceped: takes string, returns Option[ReaderOption[int, int]]
		sequenced := Sequence(original)

		// Test original
		result1 := original(10)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1("hello")
		assert.True(t, option.IsSome(innerResult1))
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, 15, value1)

		// Test sequenced - returns func(string) func(int) Option[int]
		innerFunc2 := sequenced("hello")
		innerResult2 := innerFunc2(10)
		assert.True(t, option.IsSome(innerResult2))
		value2, _ := option.Unwrap(innerResult2)
		assert.Equal(t, 15, value2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Config, returns Option[ReaderOption[Database, string]]
		query := func(cfg Config) option.Option[ReaderOption[Database, string]] {
			if cfg.Timeout <= 0 {
				return option.None[ReaderOption[Database, string]]()
			}
			return option.Some(func(db Database) option.Option[string] {
				if S.IsEmpty(db.ConnectionString) {
					return option.None[string]()
				}
				return option.Some(fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout))
			})
		}

		// Sequenceped: takes Database, returns Option[ReaderOption[Config, string]]
		sequenced := Sequence(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original with valid inputs
		result1 := query(cfg)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1(db)
		assert.True(t, option.IsSome(innerResult1))
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, expected, value1)

		// Test sequenced with valid inputs
		innerFunc2 := sequenced(db)
		innerResult2 := innerFunc2(cfg)
		assert.True(t, option.IsSome(innerResult2))
		value2, _ := option.Unwrap(innerResult2)
		assert.Equal(t, expected, value2)
	})

	t.Run("preserves outer None", func(t *testing.T) {
		// Original that returns None at outer level
		original := func(x int) option.Option[ReaderOption[string, int]] {
			if x < 0 {
				return option.None[ReaderOption[string, int]]()
			}
			return option.Some(func(s string) option.Option[int] {
				return option.Some(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with None
		result1 := original(-1)
		assert.True(t, option.IsNone(result1))

		// Test sequenced - the outer None becomes an inner None after flip
		innerFunc2 := sequenced("test")
		innerResult2 := innerFunc2(-1)
		assert.True(t, option.IsNone(innerResult2))
	})

	t.Run("preserves inner None", func(t *testing.T) {
		// Original that returns None at inner level
		original := func(x int) option.Option[ReaderOption[string, int]] {
			return option.Some(func(s string) option.Option[int] {
				if S.IsEmpty(s) {
					return option.None[int]()
				}
				return option.Some(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test original with inner None
		result1 := original(10)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1("")
		assert.True(t, option.IsNone(innerResult1))

		// Test sequenced with inner None
		innerFunc2 := sequenced("")
		innerResult2 := innerFunc2(10)
		assert.True(t, option.IsNone(innerResult2))
	})

	t.Run("works with multiple flips", func(t *testing.T) {
		// Original function
		original := func(x int) option.Option[ReaderOption[string, int]] {
			return option.Some(func(s string) option.Option[int] {
				return option.Some(x * len(s))
			})
		}

		// Sequence once
		sequenced1 := Sequence(original)

		// Test that flip produces correct results
		result1 := original(3)
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1("test")
		value1, _ := option.Unwrap(innerResult1)

		innerFunc2 := sequenced1("test")
		innerResult2 := innerFunc2(3)
		value2, _ := option.Unwrap(innerResult2)

		assert.Equal(t, value1, value2)
		assert.Equal(t, 12, value2) // 3 * 4
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) option.Option[ReaderOption[string, int]] {
			return option.Some(func(s string) option.Option[int] {
				return option.Some(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test with zero values
		result1 := original(0)
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1("")
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, 0, value1)

		innerFunc2 := sequenced("")
		innerResult2 := innerFunc2(0)
		value2, _ := option.Unwrap(innerResult2)
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

		// Original: takes User, returns Option[ReaderOption[Session, string]]
		authenticate := func(user User) option.Option[ReaderOption[Session, string]] {
			if user.ID <= 0 {
				return option.None[ReaderOption[Session, string]]()
			}
			return option.Some(func(session Session) option.Option[string] {
				if S.IsEmpty(session.Token) {
					return option.None[string]()
				}
				return option.Some(fmt.Sprintf("User %s (ID: %d) authenticated with token %s",
					user.Name, user.ID, session.Token))
			})
		}

		sequenced := Sequence(authenticate)

		user := User{ID: 42, Name: "Alice"}
		session := Session{Token: "abc123"}

		expected := "User Alice (ID: 42) authenticated with token abc123"

		// Test original
		result1 := authenticate(user)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1(session)
		assert.True(t, option.IsSome(innerResult1))
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, expected, value1)

		// Test sequenced
		innerFunc2 := sequenced(session)
		innerResult2 := innerFunc2(user)
		assert.True(t, option.IsSome(innerResult2))
		value2, _ := option.Unwrap(innerResult2)
		assert.Equal(t, expected, value2)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) option.Option[ReaderOption[Empty, int]] {
			return option.Some(func(e Empty) option.Option[int] {
				return option.Some(x * 2)
			})
		}

		sequenced := Sequence(original)

		empty := Empty{}

		result1 := original(10)
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1(empty)
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, 20, value1)

		innerFunc2 := sequenced(empty)
		innerResult2 := innerFunc2(10)
		value2, _ := option.Unwrap(innerResult2)
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) option.Option[ReaderOption[*Data, int]] {
			return option.Some(func(d *Data) option.Option[int] {
				if d == nil {
					return option.Some(x)
				}
				return option.Some(x + d.Value)
			})
		}

		sequenced := Sequence(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		result1 := original(42)
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1(data)
		value1, _ := option.Unwrap(innerResult1)
		assert.Equal(t, 142, value1)

		innerFunc2 := sequenced(data)
		innerResult2 := innerFunc2(42)
		value2, _ := option.Unwrap(innerResult2)
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3 := original(42)
		innerFunc3, _ := option.Unwrap(result3)
		innerResult3 := innerFunc3(nil)
		value3, _ := option.Unwrap(innerResult3)
		assert.Equal(t, 42, value3)

		innerFunc4 := sequenced(nil)
		innerResult4 := innerFunc4(42)
		value4, _ := option.Unwrap(innerResult4)
		assert.Equal(t, 42, value4)
	})

	t.Run("handles None at different levels", func(t *testing.T) {
		// Original that can return None at both levels
		original := func(x int) option.Option[ReaderOption[string, int]] {
			if x < 0 {
				return option.None[ReaderOption[string, int]]()
			}
			return option.Some(func(s string) option.Option[int] {
				if S.IsEmpty(s) {
					return option.None[int]()
				}
				return option.Some(x + len(s))
			})
		}

		sequenced := Sequence(original)

		// Test outer None
		result1 := original(-1)
		assert.True(t, option.IsNone(result1))

		// Test inner None
		result2 := original(10)
		assert.True(t, option.IsSome(result2))
		innerFunc2, _ := option.Unwrap(result2)
		innerResult2 := innerFunc2("")
		assert.True(t, option.IsNone(innerResult2))

		// Test sequenced with inner None
		innerFunc3 := sequenced("")
		innerResult3 := innerFunc3(10)
		assert.True(t, option.IsNone(innerResult3))
	})

	t.Run("validates None propagation with outer failure", func(t *testing.T) {
		type Config struct {
			Valid bool
		}

		original := func(cfg Config) option.Option[ReaderOption[int, string]] {
			if !cfg.Valid {
				return option.None[ReaderOption[int, string]]()
			}
			return option.Some(func(x int) option.Option[string] {
				return option.Some(fmt.Sprintf("value: %d", x))
			})
		}

		sequenced := Sequence(original)

		invalidCfg := Config{Valid: false}

		// Original with invalid config
		result1 := original(invalidCfg)
		assert.True(t, option.IsNone(result1))

		// Sequenceped returns a Reader, but when called with invalid config, it returns None
		innerFunc2 := sequenced(42)
		innerResult2 := innerFunc2(invalidCfg)
		assert.True(t, option.IsNone(innerResult2))
	})

	t.Run("validates None propagation with inner failure", func(t *testing.T) {
		original := func(threshold int) option.Option[ReaderOption[int, string]] {
			return option.Some(func(value int) option.Option[string] {
				if value < threshold {
					return option.None[string]()
				}
				return option.Some(fmt.Sprintf("valid: %d", value))
			})
		}

		sequenced := Sequence(original)

		// Test with value below threshold
		result1 := original(10)
		innerFunc1, _ := option.Unwrap(result1)
		innerResult1 := innerFunc1(5)
		assert.True(t, option.IsNone(innerResult1))

		// Sequenceped should also return None
		innerFunc2 := sequenced(5)
		innerResult2 := innerFunc2(10)
		assert.True(t, option.IsNone(innerResult2))
	})
}

func TestSequenceReader(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Option[Reader[string, int]]
		original := func(x int) option.Option[Reader[string, int]] {
			return option.Some(func(s string) int {
				return x + len(s)
			})
		}

		// Sequenceped: takes string, returns Reader[int, Option[int]]
		sequenced := SequenceReader(original)

		// Test original
		result1 := original(10)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1("hello")
		assert.Equal(t, 15, value1)

		// Test sequenced
		innerFunc2 := sequenced("hello")
		result2 := innerFunc2(10)
		assert.True(t, option.IsSome(result2))
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, 15, value2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Config, returns Option[Reader[Database, string]]
		query := func(cfg Config) option.Option[Reader[Database, string]] {
			if cfg.Timeout <= 0 {
				return option.None[Reader[Database, string]]()
			}
			return option.Some(func(db Database) string {
				return fmt.Sprintf("Query on %s with timeout %d",
					db.ConnectionString, cfg.Timeout)
			})
		}

		// Sequenceped: takes Database, returns Reader[Config, Option[string]]
		sequenced := SequenceReader(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original with valid inputs
		result1 := query(cfg)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1(db)
		assert.Equal(t, expected, value1)

		// Test sequenced with valid inputs
		innerFunc2 := sequenced(db)
		result2 := innerFunc2(cfg)
		assert.True(t, option.IsSome(result2))
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, expected, value2)
	})

	t.Run("preserves outer None", func(t *testing.T) {
		// Original that returns None at outer level
		original := func(x int) option.Option[Reader[string, int]] {
			if x < 0 {
				return option.None[Reader[string, int]]()
			}
			return option.Some(func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test original with None
		result1 := original(-1)
		assert.True(t, option.IsNone(result1))

		// Test sequenced - the outer None becomes an inner None after flip
		innerFunc2 := sequenced("test")
		result2 := innerFunc2(-1)
		assert.True(t, option.IsNone(result2))
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original: takes bool, returns Option[Reader[int, string]]
		original := func(flag bool) option.Option[Reader[int, string]] {
			if !flag {
				return option.None[Reader[int, string]]()
			}
			return option.Some(func(n int) string {
				if flag {
					return fmt.Sprintf("positive: %d", n)
				}
				return fmt.Sprintf("negative: %d", -n)
			})
		}

		sequenced := SequenceReader(original)

		// Test with true flag
		result1 := original(true)
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1(42)
		assert.Equal(t, "positive: 42", value1)

		innerFunc2 := sequenced(42)
		result2 := innerFunc2(true)
		assert.True(t, option.IsSome(result2))
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, "positive: 42", value2)

		// Test with false flag
		result3 := original(false)
		assert.True(t, option.IsNone(result3))

		innerFunc4 := sequenced(42)
		result4 := innerFunc4(false)
		assert.True(t, option.IsNone(result4))
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) option.Option[Reader[string, int]] {
			return option.Some(func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test with zero values
		result1 := original(0)
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1("")
		assert.Equal(t, 0, value1)

		innerFunc2 := sequenced("")
		result2 := innerFunc2(0)
		value2, _ := option.Unwrap(result2)
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

		// Original: takes User, returns Option[Reader[Session, string]]
		authenticate := func(user User) option.Option[Reader[Session, string]] {
			if user.ID <= 0 {
				return option.None[Reader[Session, string]]()
			}
			return option.Some(func(session Session) string {
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
		assert.True(t, option.IsSome(result1))
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1(session)
		assert.Equal(t, expected, value1)

		// Test sequenced
		innerFunc2 := sequenced(session)
		result2 := innerFunc2(user)
		assert.True(t, option.IsSome(result2))
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, expected, value2)
	})
}

func TestSequenceReaderEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) option.Option[Reader[Empty, int]] {
			return option.Some(func(e Empty) int {
				return x * 2
			})
		}

		sequenced := SequenceReader(original)

		empty := Empty{}

		result1 := original(10)
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1(empty)
		assert.Equal(t, 20, value1)

		innerFunc2 := sequenced(empty)
		result2 := innerFunc2(10)
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, 20, value2)
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) option.Option[Reader[*Data, int]] {
			return option.Some(func(d *Data) int {
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
		innerFunc1, _ := option.Unwrap(result1)
		value1 := innerFunc1(data)
		assert.Equal(t, 142, value1)

		innerFunc2 := sequenced(data)
		result2 := innerFunc2(42)
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, 142, value2)

		// Test with nil pointer
		result3 := original(42)
		innerFunc3, _ := option.Unwrap(result3)
		value3 := innerFunc3(nil)
		assert.Equal(t, 42, value3)

		innerFunc4 := sequenced(nil)
		result4 := innerFunc4(42)
		value4, _ := option.Unwrap(result4)
		assert.Equal(t, 42, value4)
	})

	t.Run("validates None propagation with outer failure", func(t *testing.T) {
		type Config struct {
			Valid bool
		}

		original := func(cfg Config) option.Option[Reader[int, string]] {
			if !cfg.Valid {
				return option.None[Reader[int, string]]()
			}
			return option.Some(func(x int) string {
				return fmt.Sprintf("value: %d", x)
			})
		}

		sequenced := SequenceReader(original)

		invalidCfg := Config{Valid: false}

		// Original with invalid config
		result1 := original(invalidCfg)
		assert.True(t, option.IsNone(result1))

		// Sequenceped returns a Reader, but when called with invalid config, it returns None
		innerFunc2 := sequenced(42)
		result2 := innerFunc2(invalidCfg)
		assert.True(t, option.IsNone(result2))
	})

	t.Run("can be used in composition", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original function
		multiply := func(x int) option.Option[Reader[Config, int]] {
			if x < 0 {
				return option.None[Reader[Config, int]]()
			}
			return option.Some(func(c Config) int {
				return x * c.Multiplier
			})
		}

		// Sequence it
		sequenced := SequenceReader(multiply)

		// Use sequenced version to partially apply config first
		config := Config{Multiplier: 10}
		multiplyBy10 := sequenced(config)

		// Now we have a function int -> Option[int]
		result1 := multiplyBy10(5)
		assert.True(t, option.IsSome(result1))
		value1, _ := option.Unwrap(result1)
		assert.Equal(t, 50, value1)

		result2 := multiplyBy10(10)
		assert.True(t, option.IsSome(result2))
		value2, _ := option.Unwrap(result2)
		assert.Equal(t, 100, value2)

		// Test with negative value
		result3 := multiplyBy10(-5)
		assert.True(t, option.IsNone(result3))
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		// The same inputs should always produce the same outputs
		original := func(x int) option.Option[Reader[string, int]] {
			return option.Some(func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Call multiple times with same inputs
		for range 5 {
			result1 := original(10)
			innerFunc1, _ := option.Unwrap(result1)
			value1 := innerFunc1("hello")
			assert.Equal(t, 15, value1)

			innerFunc2 := sequenced("hello")
			result2 := innerFunc2(10)
			value2, _ := option.Unwrap(result2)
			assert.Equal(t, 15, value2)
		}
	})
}
