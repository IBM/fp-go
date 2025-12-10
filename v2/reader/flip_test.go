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

func TestSequence(t *testing.T) {
	t.Run("flips parameter order for simple types", func(t *testing.T) {
		// Original: takes int, returns Reader[string, int]
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		// Sequenceped: takes string, returns Reader[int, int]
		sequenced := Sequence(original)

		// Test original
		result1 := original(10)("hello") // 10 + 5 = 15
		assert.Equal(t, 15, result1)

		// Test sequenced
		result2 := sequenced("hello")(10) // 10 + 5 = 15
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

		// Sequenceped: takes Config, returns Reader[Port, string]
		sequenced := Sequence(makeURL)

		config := Config{Host: "localhost"}
		port := Port(8080)

		// Test original
		result1 := makeURL(port)(config)
		assert.Equal(t, "localhost:8080", result1)

		// Test sequenced
		result2 := sequenced(config)(port)
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

		sequenced := Sequence(original)

		// Test with true flag
		assert.Equal(t, "positive: 42", original(true)(42))
		assert.Equal(t, "positive: 42", sequenced(42)(true))

		// Test with false flag
		assert.Equal(t, "negative: -42", original(false)(42))
		assert.Equal(t, "negative: -42", sequenced(42)(false))
	})

	t.Run("works with multiple flips", func(t *testing.T) {
		// Original function
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x * len(s)
			}
		}

		// Sequence once
		sequenced1 := Sequence(original)

		// Sequence twice (should be equivalent to original)
		sequenced2 := Sequence(sequenced1)

		// Test that double flip returns to original behavior
		result1 := original(3)("test")   // 3 * 4 = 12
		result2 := sequenced2(3)("test") // Should also be 12
		assert.Equal(t, result1, result2)
		assert.Equal(t, 12, result2)
	})

	t.Run("works with zero values", func(t *testing.T) {
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		sequenced := Sequence(original)

		// Test with zero values
		result1 := original(0)("")
		result2 := sequenced("")(0)
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

		sequenced := Sequence(executeQuery)

		db := Database{ConnectionString: "localhost:5432"}
		query := Query{SQL: "SELECT * FROM users"}

		expected := "Executing 'SELECT * FROM users' on localhost:5432"

		// Test original
		result1 := executeQuery(query)(db)
		assert.Equal(t, expected, result1)

		// Test sequenced
		result2 := sequenced(db)(query)
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
		sequencedStr := Sequence(strFunc)
		assert.Equal(t, "42", strFunc(42)(true))
		assert.Equal(t, "42", sequencedStr(true)(42))

		// Bool result
		boolFunc := func(s string) Reader[int, bool] {
			return func(n int) bool {
				return len(s) > n
			}
		}
		sequencedBool := Sequence(boolFunc)
		assert.True(t, boolFunc("hello")(3))
		assert.True(t, sequencedBool(3)("hello"))

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
		sequencedSlice := Sequence(sliceFunc)
		expected := []int{5, 6, 7}
		assert.Equal(t, expected, sliceFunc(3)("hello"))
		assert.Equal(t, expected, sequencedSlice("hello")(3))
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

		// Sequence it
		sequenced := Sequence(multiply)

		// Use sequenced version to partially apply config first
		config := Config{Multiplier: 10}
		multiplyBy10 := sequenced(config)

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

		sequenced := Sequence(original)

		data := &Data{Value: 100}

		// Test with non-nil pointer
		assert.Equal(t, 142, original(42)(data))
		assert.Equal(t, 142, sequenced(data)(42))

		// Test with nil pointer
		assert.Equal(t, 42, original(42)(nil))
		assert.Equal(t, 42, sequenced(nil)(42))
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		// The same inputs should always produce the same outputs
		original := func(x int) Reader[string, int] {
			return func(s string) int {
				return x + len(s)
			}
		}

		sequenced := Sequence(original)

		// Call multiple times with same inputs
		for range 5 {
			assert.Equal(t, 15, original(10)("hello"))
			assert.Equal(t, 15, sequenced("hello")(10))
		}
	})
}

func TestSequenceWithChain(t *testing.T) {
	t.Run("sequenced function can be used with Chain", func(t *testing.T) {
		type Config struct {
			BaseValue int
		}

		// A function that takes an int and returns a Reader
		addToBase := func(x int) Reader[Config, int] {
			return func(c Config) int {
				return c.BaseValue + x
			}
		}

		// Sequence it
		sequenced := Sequence(addToBase)

		// Use it in a chain
		config := Config{BaseValue: 100}

		// Create a reader that uses the sequenced function
		reader := sequenced(config)

		// Apply it
		result := reader(42) // 100 + 42
		assert.Equal(t, 142, result)
	})
}

func TestSequenceEdgeCases(t *testing.T) {
	t.Run("works with empty struct", func(t *testing.T) {
		type Empty struct{}

		original := func(x int) Reader[Empty, int] {
			return func(e Empty) int {
				return x * 2
			}
		}

		sequenced := Sequence(original)

		empty := Empty{}
		assert.Equal(t, 20, original(10)(empty))
		assert.Equal(t, 20, sequenced(empty)(10))
	})

	t.Run("works with function types", func(t *testing.T) {
		// Test with function types as parameters
		type Transform func(string) string

		original := func(x int) Reader[Transform, string] {
			return func(t Transform) string {
				return t(fmt.Sprintf("%d", x))
			}
		}

		sequenced := Sequence(original)

		transform := func(s string) string {
			return "value: " + s
		}

		assert.Equal(t, "value: 42", original(42)(transform))
		assert.Equal(t, "value: 42", sequenced(transform)(42))
	})
}

func TestTraverse(t *testing.T) {
	t.Run("basic traverse with two environments", func(t *testing.T) {
		type Database struct {
			UserID int
		}
		type Config struct {
			Prefix string
		}

		// Reader that gets user ID from database
		getUserID := func(db Database) int {
			return db.UserID
		}

		// Kleisli that formats user ID with config
		formatUser := func(id int) Reader[Config, string] {
			return func(c Config) string {
				return fmt.Sprintf("%s%d", c.Prefix, id)
			}
		}

		// Traverse applies formatUser to the result of getUserID
		traversed := Traverse[Database](formatUser)(getUserID)

		// Apply both environments
		config := Config{Prefix: "User-"}
		db := Database{UserID: 42}
		result := traversed(config)(db)

		assert.Equal(t, "User-42", result)
	})

	t.Run("traverse with computation", func(t *testing.T) {
		type Source struct {
			Value int
		}
		type Multiplier struct {
			Factor int
		}

		// Reader that extracts value from source
		getValue := func(s Source) int {
			return s.Value
		}

		// Kleisli that multiplies value with multiplier
		multiply := func(n int) Reader[Multiplier, int] {
			return func(m Multiplier) int {
				return n * m.Factor
			}
		}

		traversed := Traverse[Source](multiply)(getValue)

		source := Source{Value: 10}
		multiplier := Multiplier{Factor: 5}
		result := traversed(multiplier)(source)

		assert.Equal(t, 50, result)
	})

	t.Run("traverse with string transformation", func(t *testing.T) {
		type Input struct {
			Text string
		}
		type Format struct {
			Template string
		}

		// Reader that gets text from input
		getText := func(i Input) string {
			return i.Text
		}

		// Kleisli that formats text with template
		format := func(text string) Reader[Format, string] {
			return func(f Format) string {
				return fmt.Sprintf(f.Template, text)
			}
		}

		traversed := Traverse[Input](format)(getText)

		input := Input{Text: "world"}
		formatCfg := Format{Template: "Hello, %s!"}
		result := traversed(formatCfg)(input)

		assert.Equal(t, "Hello, world!", result)
	})

	t.Run("traverse with boolean logic", func(t *testing.T) {
		type Data struct {
			Value int
		}
		type Threshold struct {
			Limit int
		}

		// Reader that gets value from data
		getValue := func(d Data) int {
			return d.Value
		}

		// Kleisli that checks if value exceeds threshold
		checkThreshold := func(val int) Reader[Threshold, bool] {
			return func(t Threshold) bool {
				return val > t.Limit
			}
		}

		traversed := Traverse[Data](checkThreshold)(getValue)

		data := Data{Value: 100}
		threshold := Threshold{Limit: 50}
		result := traversed(threshold)(data)

		assert.True(t, result)

		threshold2 := Threshold{Limit: 150}
		result2 := traversed(threshold2)(data)

		assert.False(t, result2)
	})

	t.Run("traverse with slice transformation", func(t *testing.T) {
		type Source struct {
			Items []string
		}
		type Config struct {
			Separator string
		}

		// Reader that gets items from source
		getItems := func(s Source) []string {
			return s.Items
		}

		// Kleisli that joins items with separator
		joinItems := func(items []string) Reader[Config, string] {
			return func(c Config) string {
				result := ""
				for i, item := range items {
					if i > 0 {
						result += c.Separator
					}
					result += item
				}
				return result
			}
		}

		traversed := Traverse[Source](joinItems)(getItems)

		source := Source{Items: []string{"a", "b", "c"}}
		config := Config{Separator: ", "}
		result := traversed(config)(source)

		assert.Equal(t, "a, b, c", result)
	})

	t.Run("traverse with struct transformation", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}
		type Database struct {
			TablePrefix string
		}
		type UserRecord struct {
			Table string
			ID    int
			Name  string
		}

		// Reader that gets user
		getUser := func(db Database) User {
			return User{ID: 123, Name: "Alice"}
		}

		// Kleisli that creates user record
		createRecord := func(user User) Reader[Database, UserRecord] {
			return func(db Database) UserRecord {
				return UserRecord{
					Table: db.TablePrefix + "users",
					ID:    user.ID,
					Name:  user.Name,
				}
			}
		}

		traversed := Traverse[Database](createRecord)(getUser)

		db := Database{TablePrefix: "prod_"}
		result := traversed(db)(db)

		assert.Equal(t, UserRecord{
			Table: "prod_users",
			ID:    123,
			Name:  "Alice",
		}, result)
	})

	t.Run("traverse with nil handling", func(t *testing.T) {
		type Source struct {
			Value *int
		}
		type Config struct {
			Default int
		}

		// Reader that gets pointer value
		getValue := func(s Source) *int {
			return s.Value
		}

		// Kleisli that handles nil with default
		handleNil := func(ptr *int) Reader[Config, int] {
			return func(c Config) int {
				if ptr == nil {
					return c.Default
				}
				return *ptr
			}
		}

		traversed := Traverse[Source](handleNil)(getValue)

		config := Config{Default: 999}

		// Test with non-nil value
		val := 42
		source1 := Source{Value: &val}
		result1 := traversed(config)(source1)
		assert.Equal(t, 42, result1)

		// Test with nil value
		source2 := Source{Value: nil}
		result2 := traversed(config)(source2)
		assert.Equal(t, 999, result2)
	})

	t.Run("traverse composition", func(t *testing.T) {
		type Env1 struct {
			Base int
		}
		type Env2 struct {
			Multiplier int
		}

		// Reader that gets base value
		getBase := func(e Env1) int {
			return e.Base
		}

		// Kleisli that multiplies
		multiply := func(n int) Reader[Env2, int] {
			return func(e Env2) int {
				return n * e.Multiplier
			}
		}

		// Another Kleisli that adds
		add := func(n int) Reader[Env2, int] {
			return func(e Env2) int {
				return n + e.Multiplier
			}
		}

		// Traverse with multiply
		traversed1 := Traverse[Env1](multiply)(getBase)
		env1 := Env1{Base: 10}
		env2 := Env2{Multiplier: 5}
		result1 := traversed1(env2)(env1)
		assert.Equal(t, 50, result1)

		// Traverse with add
		traversed2 := Traverse[Env1](add)(getBase)
		result2 := traversed2(env2)(env1)
		assert.Equal(t, 15, result2)
	})

	t.Run("traverse with identity", func(t *testing.T) {
		type Env1 struct {
			Value string
		}
		type Env2 struct {
			Prefix string
		}

		// Reader that gets value
		getValue := func(e Env1) string {
			return e.Value
		}

		// Identity Kleisli (just wraps in Of)
		identity := func(s string) Reader[Env2, string] {
			return Of[Env2](s)
		}

		traversed := Traverse[Env1](identity)(getValue)

		env1 := Env1{Value: "test"}
		env2 := Env2{Prefix: "ignored"}
		result := traversed(env2)(env1)

		assert.Equal(t, "test", result)
	})

	t.Run("traverse with complex computation", func(t *testing.T) {
		type Request struct {
			UserID int
		}
		type Database struct {
			Users map[int]string
		}
		type Response struct {
			UserID   int
			UserName string
			Found    bool
		}

		// Reader that gets user ID from request
		getUserID := func(r Request) int {
			return r.UserID
		}

		// Kleisli that looks up user in database
		lookupUser := func(id int) Reader[Database, Response] {
			return func(db Database) Response {
				name, found := db.Users[id]
				return Response{
					UserID:   id,
					UserName: name,
					Found:    found,
				}
			}
		}

		traversed := Traverse[Request](lookupUser)(getUserID)

		request := Request{UserID: 42}
		db := Database{
			Users: map[int]string{
				42: "Alice",
				99: "Bob",
			},
		}

		result := traversed(db)(request)

		assert.Equal(t, Response{
			UserID:   42,
			UserName: "Alice",
			Found:    true,
		}, result)

		// Test with missing user
		request2 := Request{UserID: 123}
		result2 := traversed(db)(request2)

		assert.Equal(t, Response{
			UserID:   123,
			UserName: "",
			Found:    false,
		}, result2)
	})
}
