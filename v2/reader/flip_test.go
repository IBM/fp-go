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

package reader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlip(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Reader[string, int]
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		// Flipped: takes string, returns Reader[int, int]
		flipped := Flip(original)

		// Test original
		result1 := original(10)("hello") // 10 + 5 = 15
		assert.Equal(t, 15, result1)

		// Test flipped
		result2 := flipped("hello")(10) // 10 + 5 = 15
		assert.Equal(t, 15, result2)
	})

	t.Run("flips parameter order for struct types", func(t *testing.T) {
		type Config struct {
			Host string
		}
		type Port int

		// Original: takes Port, returns Reader[Config, string]
		makeURL := func(port Port) Reader[Config, string] {
			return func(c Config) string {
				return fmt.Sprintf("%s:%d", c.Host, port)
			}
		}

		// Flipped: takes Config, returns Reader[Port, string]
		flipped := Flip(makeURL)

		config := Config{Host: "localhost"}
		port := Port(8080)

		// Test original
		result1 := makeURL(port)(config)
		assert.Equal(t, "localhost:8080", result1)

		// Test flipped
		result2 := flipped(config)(port)
		assert.Equal(t, "localhost:8080", result2)
	})

	t.Run("preserves computation logic", func(t *testing.T) {
		// Original: takes bool, returns Reader[int, string]
		original := func(flag bool) Reader[int, string] {
			return func(n int) string {
				if flag {
					return fmt.Sprintf("positive: %d", n)
				}
				return fmt.Sprintf("negative: %d", -n)
			}
		}

		flipped := Flip(original)

		// Test with true flag
		assert.Equal(t, "positive: 42", original(true)(42))
		assert.Equal(t, "positive: 42", flipped(42)(true))

		// Test with false flag
		assert.Equal(t, "negative: -42", original(false)(42))
		assert.Equal(t, "negative: -42", flipped(42)(false))
	})

	t.Run("works with multiple flips", func(t *testing.T) {
		// Original function
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x * len(s)
			}
		}

		// Flip once
		flipped1 := Flip(original)

		// Flip twice (should be equivalent to original)
		flipped2 := Flip(flipped1)

		// Test that double flip returns to original behavior
		result1 := original(3)("test") // 3 * 4 = 12
		result2 := flipped2(3)("test") // Should also be 12
		assert.Equal(t, result1, result2)
		assert.Equal(t, 12, result2)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		flipped := Flip(original)

		// Test with zero values
		result1 := original(0)("")
		result2 := flipped("")(0)
		assert.Equal(t, 0, result1)
		assert.Equal(t, 0, result2)
	})

	t.Run("works with complex computations", func(t *testing.T) {
		type Database struct {
			ConnectionString string
		}
		type Query struct {
			SQL string
		}

		// Original: takes Query, returns Reader[Database, string]
		executeQuery := func(query Query) Reader[Database, string] {
			return func(db Database) string {
				return fmt.Sprintf("Executing '%s' on %s", query.SQL, db.ConnectionString)
			}
		}

		flipped := Flip(executeQuery)

		db := Database{ConnectionString: "localhost:5432"}
		query := Query{SQL: "SELECT * FROM users"}

		expected := "Executing 'SELECT * FROM users' on localhost:5432"

		// Test original
		result1 := executeQuery(query)(db)
		assert.Equal(t, expected, result1)

		// Test flipped
		result2 := flipped(db)(query)
		assert.Equal(t, expected, result2)
	})

	t.Run("works with different result types", func(t *testing.T) {
		// Test with various result types

		// String result
		strFunc := func(x int) Reader[bool, string] {
			return func(b bool) string {
				if b {
					return fmt.Sprintf("%d", x)
				}
				return ""
			}
		}
		flippedStr := Flip(strFunc)
		assert.Equal(t, "42", strFunc(42)(true))
		assert.Equal(t, "42", flippedStr(true)(42))

		// Bool result
		boolFunc := func(s string) Reader[int, bool] {
			return func(n int) bool {
				return len(s) > n
			}
		}
		flippedBool := Flip(boolFunc)
		assert.True(t, boolFunc("hello")(3))
		assert.True(t, flippedBool(3)("hello"))

		// Slice result
		sliceFunc := func(n int) Reader[string, []int] {
			return func(s string) []int {
				result := make([]int, n)
				for i := 0; i < n; i++ {
					result[i] = len(s) + i
				}
				return result
			}
		}
		flippedSlice := Flip(sliceFunc)
		expected := []int{5, 6, 7}
		assert.Equal(t, expected, sliceFunc(3)("hello"))
		assert.Equal(t, expected, flippedSlice("hello")(3))
	})

	t.Run("can be used in composition", func(t *testing.T) {
		type Config struct {
			Multiplier int
		}

		// Original function
		multiply := func(x int) Reader[Config, int] {
			return func(c Config) int {
				return x * c.Multiplier
			}
		}

		// Flip it
		flipped := Flip(multiply)

		// Use flipped version to partially apply config first
		config := Config{Multiplier: 10}
		multiplyBy10 := flipped(config)

		// Now we have a simple function int -> int
		assert.Equal(t, 50, multiplyBy10(5))
		assert.Equal(t, 100, multiplyBy10(10))
		assert.Equal(t, 150, multiplyBy10(15))
	})

	t.Run("works with pointer types", func(t *testing.T) {
		type Data struct {
			Value int
		}

		original := func(x int) Reader[*Data, int] {
			return func(d *Data) int {
				if d == nil {
					return x
				}
				return x + d.Value
			}
		}

		flipped := Flip(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		assert.Equal(t, 142, original(42)(data))
		assert.Equal(t, 142, flipped(data)(42))

		// Test with nil pointer
		assert.Equal(t, 42, original(42)(nil))
		assert.Equal(t, 42, flipped(nil)(42))
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		// The same inputs should always produce the same outputs
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		flipped := Flip(original)

		// Call multiple times with same inputs
		for range 5 {
			assert.Equal(t, 15, original(10)("hello"))
			assert.Equal(t, 15, flipped("hello")(10))
		}
	})
}

func TestFlipWithChain(t *testing.T) {
	t.Run("flipped function can be used with Chain", func(t *testing.T) {
		type Config struct {
			BaseValue int
		}

		// A function that takes an int and returns a Reader
		addToBase := func(x int) Reader[Config, int] {
			return func(c Config) int {
				return c.BaseValue + x
			}
		}

		// Flip it
		flipped := Flip(addToBase)

		// Use it in a chain
		config := Config{BaseValue: 100}

		// Create a reader that uses the flipped function
		reader := flipped(config)

		// Apply it
		result := reader(42) // 100 + 42
		assert.Equal(t, 142, result)
	})
}

func TestFlipEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) Reader[Empty, int] {
			return func(e Empty) int {
				return x * 2
			}
		}

		flipped := Flip(original)

		empty := Empty{}
		assert.Equal(t, 20, original(10)(empty))
		assert.Equal(t, 20, flipped(empty)(10))
	})

	t.Run("works with function types", func(t *testing.T) {
		// Test with function types as parameters
		type Transform func(string) string

		original := func(x int) Reader[Transform, string] {
			return func(t Transform) string {
				return t(fmt.Sprintf("%d", x))
			}
		}

		flipped := Flip(original)

		transform := func(s string) string {
			return "value: " + s
		}

		assert.Equal(t, "value: 42", original(42)(transform))
		assert.Equal(t, "value: 42", flipped(transform)(42))
	})
}
