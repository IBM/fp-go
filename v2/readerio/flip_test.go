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

package readerio

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

func TestSequenceReader(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns IO[Reader[string, int]]
		original := func(x int) io.IO[reader.Reader[string, int]] {
			return io.Of(func(s string) int {
				return x + len(s)
			})
		}

		// Sequenceped: takes string, returns Reader[int, IO[int]]
		sequenced := SequenceReader(original)

		// Test original
		result1 := original(10)()("hello") // 10 + 5 = 15
		assert.Equal(t, 15, result1)

		// Test sequenced - note different call pattern
		result2 := sequenced("hello")(10)() // 10 + 5 = 15
		assert.Equal(t, 15, result2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Database, returns IO[Reader[Config, string]]
		query := func(db Database) io.IO[reader.Reader[Config, string]] {
			return io.Of(func(cfg Config) string {
				return fmt.Sprintf("Query on %s with timeout %d", db.ConnectionString, cfg.Timeout)
			})
		}

		// Sequenceped: takes Config, returns Reader[Database, IO[string]]
		sequenced := SequenceReader(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original
		result1 := query(db)()(cfg)
		assert.Equal(t, expected, result1)

		// Test sequenced - note different call pattern
		result2 := sequenced(cfg)(db)()
		assert.Equal(t, expected, result2)
	})

	t.Run("preserves IO effects", func(t *testing.T) {
		callCount := 0

		// Original with side effect
		original := func(x int) io.IO[reader.Reader[string, int]] {
			return func() reader.Reader[string, int] {
				callCount++
				return func(s string) int {
					return x + len(s)
				}
			}
		}

		sequenced := SequenceReader(original)

		// Execute original
		callCount = 0
		_ = original(10)()("test")
		assert.Equal(t, 1, callCount)

		// Execute sequenced - note different call pattern
		callCount = 0
		_ = sequenced("test")(10)()
		assert.Equal(t, 1, callCount)
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original function
		original := func(x int) io.IO[reader.Reader[string, int]] {
			return io.Of(func(s string) int {
				return x * len(s)
			})
		}

		// Sequence once
		sequenced := SequenceReader(original)

		// Test that flip produces correct results
		result1 := original(3)()("test")  // 3 * 4 = 12
		result2 := sequenced("test")(3)() // Should also be 12
		assert.Equal(t, result1, result2)
		assert.Equal(t, 12, result2)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) io.IO[reader.Reader[string, int]] {
			return io.Of(func(s string) int {
				return x + len(s)
			})
		}

		sequenced := SequenceReader(original)

		// Test with zero values
		result1 := original(0)()("")
		result2 := sequenced("")(0)()
		assert.Equal(t, 0, result1)
		assert.Equal(t, 0, result2)
	})

	t.Run("works with complex computations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		type Session struct {
			Token string
		}

		// Original: takes User, returns IO[Reader[Session, string]]
		authenticate := func(user User) io.IO[reader.Reader[Session, string]] {
			return io.Of(func(session Session) string {
				return fmt.Sprintf("User %s (ID: %d) authenticated with token %s",
					user.Name, user.ID, session.Token)
			})
		}

		sequenced := SequenceReader(authenticate)

		user := User{ID: 42, Name: "Alice"}
		session := Session{Token: "abc123"}

		expected := "User Alice (ID: 42) authenticated with token abc123"

		// Test original
		result1 := authenticate(user)()(session)
		assert.Equal(t, expected, result1)

		// Test sequenced - note different call pattern
		result2 := sequenced(session)(user)()
		assert.Equal(t, expected, result2)
	})
}

func TestSequence(t *testing.T) {
	t.Run("flips nested ReaderIO parameter order", func(t *testing.T) {
		// Original: takes int, returns IO[ReaderIO[string, int]]
		original := func(x int) io.IO[ReaderIO[string, int]] {
			return io.Of(func(s string) io.IO[int] {
				return io.Of(x + len(s))
			})
		}

		// Sequenceped: takes string, returns IO[ReaderIO[int, int]]
		sequenced := Sequence(original)

		// Test original
		result1 := original(10)()("hello")() // 10 + 5 = 15
		assert.Equal(t, 15, result1)

		// Test sequenced - returns func(string) func(int) IO[int]
		innerFunc := sequenced("hello")
		result2 := innerFunc(10)() // 10 + 5 = 15
		assert.Equal(t, 15, result2)
	})

	t.Run("flips nested ReaderIO for struct types", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Config struct {
			Timeout int
		}

		// Original: takes Database, returns IO[ReaderIO[Config, string]]
		query := func(db Database) io.IO[ReaderIO[Config, string]] {
			return io.Of(func(cfg Config) io.IO[string] {
				return io.Of(fmt.Sprintf("Query on %s with timeout %d", db.ConnectionString, cfg.Timeout))
			})
		}

		// Sequenceped: takes Config, returns IO[ReaderIO[Database, string]]
		sequenced := Sequence(query)

		db := Database{ConnectionString: "localhost:5432"}
		cfg := Config{Timeout: 30}

		expected := "Query on localhost:5432 with timeout 30"

		// Test original
		result1 := query(db)()(cfg)()
		assert.Equal(t, expected, result1)

		// Test sequenced - returns func(Config) func(Database) IO[string]
		innerFunc := sequenced(cfg)
		result2 := innerFunc(db)()
		assert.Equal(t, expected, result2)
	})

	t.Run("preserves nested IO effects", func(t *testing.T) {
		outerCallCount := 0
		innerCallCount := 0

		// Original with side effects at both levels
		original := func(x int) io.IO[ReaderIO[string, int]] {
			return func() ReaderIO[string, int] {
				outerCallCount++
				return func(s string) io.IO[int] {
					return func() int {
						innerCallCount++
						return x + len(s)
					}
				}
			}
		}

		sequenced := Sequence(original)

		// Execute original
		outerCallCount = 0
		innerCallCount = 0
		_ = original(10)()("test")()
		assert.Equal(t, 1, outerCallCount)
		assert.Equal(t, 1, innerCallCount)

		// Execute sequenced
		outerCallCount = 0
		innerCallCount = 0
		innerFunc := sequenced("test")
		_ = innerFunc(10)()
		assert.Equal(t, 1, outerCallCount)
		assert.Equal(t, 1, innerCallCount)
	})

	t.Run("works with multiple flips", func(t *testing.T) {
		// Original function
		original := func(x int) io.IO[ReaderIO[string, int]] {
			return io.Of(func(s string) io.IO[int] {
				return io.Of(x * len(s))
			})
		}

		// Sequence once
		sequenced1 := Sequence(original)

		// Test that flip produces correct results
		result1 := original(3)()("test")() // 3 * 4 = 12
		innerFunc := sequenced1("test")
		result2 := innerFunc(3)() // Should also be 12
		assert.Equal(t, result1, result2)
		assert.Equal(t, 12, result2)
	})
}

func TestSequenceReaderEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) io.IO[reader.Reader[Empty, int]] {
			return io.Of(func(e Empty) int {
				return x * 2
			})
		}

		sequenced := SequenceReader(original)

		empty := Empty{}
		assert.Equal(t, 20, original(10)()(empty))
		assert.Equal(t, 20, sequenced(empty)(10)())
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) io.IO[reader.Reader[*Data, int]] {
			return io.Of(func(d *Data) int {
				if d == nil {
					return x
				}
				return x + d.Value
			})
		}

		sequenced := SequenceReader(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		assert.Equal(t, 142, original(42)()(data))
		assert.Equal(t, 142, sequenced(data)(42)())

		// Test with nil pointer
		assert.Equal(t, 42, original(42)()(nil))
		assert.Equal(t, 42, sequenced(nil)(42)())
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) io.IO[ReaderIO[Empty, int]] {
			return io.Of(func(e Empty) io.IO[int] {
				return io.Of(x * 2)
			})
		}

		sequenced := Sequence(original)

		empty := Empty{}
		assert.Equal(t, 20, original(10)()(empty)())
		innerFunc := sequenced(empty)
		assert.Equal(t, 20, innerFunc(10)())
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) io.IO[ReaderIO[*Data, int]] {
			return io.Of(func(d *Data) io.IO[int] {
				return io.Of(func() int {
					if d == nil {
						return x
					}
					return x + d.Value
				}())
			})
		}

		sequenced := Sequence(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		assert.Equal(t, 142, original(42)()(data)())
		innerFunc1 := sequenced(data)
		assert.Equal(t, 142, innerFunc1(42)())

		// Test with nil pointer
		assert.Equal(t, 42, original(42)()(nil)())
		innerFunc2 := sequenced(nil)
		assert.Equal(t, 42, innerFunc2(42)())
	})
}
