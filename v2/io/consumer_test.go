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

package io

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestChainConsumer(t *testing.T) {
	t.Run("executes consumer with IO value", func(t *testing.T) {
		var captured int
		consumer := func(x int) {
			captured = x
		}

		result := F.Pipe1(
			Of(42),
			ChainConsumer(consumer),
		)

		// Execute the IO
		result()

		// Verify the consumer was called with the correct value
		assert.Equal(t, 42, captured)
	})

	t.Run("returns empty struct", func(t *testing.T) {
		consumer := func(x int) {
			// no-op consumer
		}

		result := F.Pipe1(
			Of(100),
			ChainConsumer(consumer),
		)

		// Execute and verify return type
		output := result()
		assert.Equal(t, struct{}{}, output)
	})

	t.Run("can be chained with other IO operations", func(t *testing.T) {
		var captured int
		consumer := func(x int) {
			captured = x
		}

		result := F.Pipe2(
			Of(21),
			ChainConsumer(consumer),
			Map(func(struct{}) int { return captured * 2 }),
		)

		output := result()
		assert.Equal(t, 42, output)
		assert.Equal(t, 21, captured)
	})

	t.Run("works with string values", func(t *testing.T) {
		var captured string
		consumer := func(s string) {
			captured = s
		}

		result := F.Pipe1(
			Of("hello"),
			ChainConsumer(consumer),
		)

		result()
		assert.Equal(t, "hello", captured)
	})

	t.Run("works with complex types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		var captured User
		consumer := func(u User) {
			captured = u
		}

		user := User{Name: "Alice", Age: 30}
		result := F.Pipe1(
			Of(user),
			ChainConsumer(consumer),
		)

		result()
		assert.Equal(t, user, captured)
	})

	t.Run("multiple consumers in sequence", func(t *testing.T) {
		var values []int
		consumer1 := func(x int) {
			values = append(values, x)
		}
		consumer2 := func(_ struct{}) {
			values = append(values, 999)
		}

		result := F.Pipe2(
			Of(42),
			ChainConsumer(consumer1),
			ChainConsumer(consumer2),
		)

		result()
		assert.Equal(t, []int{42, 999}, values)
	})

	t.Run("consumer with side effects", func(t *testing.T) {
		counter := 0
		consumer := func(x int) {
			counter += x
		}

		// Execute multiple times
		op := ChainConsumer(consumer)
		io1 := op(Of(10))
		io2 := op(Of(20))
		io3 := op(Of(30))

		io1()
		assert.Equal(t, 10, counter)

		io2()
		assert.Equal(t, 30, counter)

		io3()
		assert.Equal(t, 60, counter)
	})

	t.Run("consumer in a pipeline with Map", func(t *testing.T) {
		var log []string
		logger := func(s string) {
			log = append(log, s)
		}

		result := F.Pipe3(
			Of("start"),
			ChainConsumer(logger),
			Map(func(struct{}) string { return "middle" }),
			Chain(func(s string) IO[string] {
				logger(s)
				return Of("end")
			}),
		)

		output := result()
		assert.Equal(t, "end", output)
		assert.Equal(t, []string{"start", "middle"}, log)
	})

	t.Run("consumer does not affect IO chain on panic recovery", func(t *testing.T) {
		var captured int
		safeConsumer := func(x int) {
			captured = x
		}

		result := F.Pipe2(
			Of(42),
			ChainConsumer(safeConsumer),
			Map(func(struct{}) string { return "success" }),
		)

		output := result()
		assert.Equal(t, "success", output)
		assert.Equal(t, 42, captured)
	})

	t.Run("consumer with pointer types", func(t *testing.T) {
		var captured *int
		consumer := func(p *int) {
			captured = p
		}

		value := 42
		result := F.Pipe1(
			Of(&value),
			ChainConsumer(consumer),
		)

		result()
		assert.Equal(t, &value, captured)
		assert.Equal(t, 42, *captured)
	})

	t.Run("consumer with slice accumulation", func(t *testing.T) {
		var accumulated []int
		consumer := func(x int) {
			accumulated = append(accumulated, x)
		}

		op := ChainConsumer(consumer)

		// Create multiple IOs and execute them
		for i := 1; i <= 5; i++ {
			io := op(Of(i))
			io()
		}

		assert.Equal(t, []int{1, 2, 3, 4, 5}, accumulated)
	})

	t.Run("consumer with map accumulation", func(t *testing.T) {
		counts := make(map[string]int)
		consumer := func(s string) {
			counts[s]++
		}

		op := ChainConsumer(consumer)

		words := []string{"hello", "world", "hello", "test", "world", "hello"}
		for _, word := range words {
			io := op(Of(word))
			io()
		}

		assert.Equal(t, 3, counts["hello"])
		assert.Equal(t, 2, counts["world"])
		assert.Equal(t, 1, counts["test"])
	})

	t.Run("lazy evaluation - consumer not called until IO executed", func(t *testing.T) {
		called := false
		consumer := func(x int) {
			called = true
		}

		// Create the IO but don't execute it
		io := F.Pipe1(
			Of(42),
			ChainConsumer(consumer),
		)

		// Consumer should not be called yet
		assert.False(t, called)

		// Now execute
		io()

		// Consumer should be called now
		assert.True(t, called)
	})
}
