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

package ioeither

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestChainConsumer(t *testing.T) {
	t.Run("executes consumer with Right value", func(t *testing.T) {
		var captured int
		consumer := func(x int) {
			captured = x
		}

		result := F.Pipe1(
			Right[error](42),
			ChainConsumer[error](consumer),
		)

		// Execute the IOEither
		result()

		// Verify the consumer was called with the correct value
		assert.Equal(t, 42, captured)
	})

	t.Run("does not execute consumer with Left value", func(t *testing.T) {
		called := false
		consumer := func(x int) {
			called = true
		}

		result := F.Pipe1(
			Left[int](errors.New("error")),
			ChainConsumer[error](consumer),
		)

		// Execute the IOEither
		result()

		// Verify the consumer was NOT called
		assert.False(t, called)
	})

	t.Run("returns Right with empty struct for Right input", func(t *testing.T) {
		consumer := func(x int) {
			// no-op consumer
		}

		result := F.Pipe1(
			Right[error](100),
			ChainConsumer[error](consumer),
		)

		// Execute and verify return type
		output := result()
		assert.Equal(t, E.Of[error](struct{}{}), output)
	})

	t.Run("returns Left unchanged for Left input", func(t *testing.T) {
		consumer := func(x int) {
			// no-op consumer
		}

		err := errors.New("test error")
		result := F.Pipe1(
			Left[int](err),
			ChainConsumer[error](consumer),
		)

		// Execute and verify error is preserved
		output := result()
		assert.True(t, E.IsLeft(output))
		_, leftErr := E.Unwrap(output)
		assert.Equal(t, err, leftErr)
	})

	t.Run("can be chained with other IOEither operations", func(t *testing.T) {
		var captured int
		consumer := func(x int) {
			captured = x
		}

		result := F.Pipe2(
			Right[error](21),
			ChainConsumer[error](consumer),
			Map[error](func(struct{}) int { return captured * 2 }),
		)

		output := result()
		assert.Equal(t, E.Right[error](42), output)
		assert.Equal(t, 21, captured)
	})

	t.Run("works with string values", func(t *testing.T) {
		var captured string
		consumer := func(s string) {
			captured = s
		}

		result := F.Pipe1(
			Right[error]("hello"),
			ChainConsumer[error](consumer),
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
			Right[error](user),
			ChainConsumer[error](consumer),
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
			Right[error](42),
			ChainConsumer[error](consumer1),
			ChainConsumer[error](consumer2),
		)

		result()
		assert.Equal(t, []int{42, 999}, values)
	})

	t.Run("consumer with side effects on Right values only", func(t *testing.T) {
		counter := 0
		consumer := func(x int) {
			counter += x
		}

		op := ChainConsumer[error](consumer)

		// Execute with Right values
		io1 := op(Right[error](10))
		io2 := op(Right[error](20))
		io3 := op(Left[int](errors.New("error")))
		io4 := op(Right[error](30))

		io1()
		assert.Equal(t, 10, counter)

		io2()
		assert.Equal(t, 30, counter)

		io3() // Should not increment counter
		assert.Equal(t, 30, counter)

		io4()
		assert.Equal(t, 60, counter)
	})

	t.Run("consumer in a pipeline with Map and Chain", func(t *testing.T) {
		var log []string
		logger := func(s string) {
			log = append(log, s)
		}

		result := F.Pipe3(
			Right[error]("start"),
			ChainConsumer[error](logger),
			Map[error](func(struct{}) string { return "middle" }),
			Chain(func(s string) IOEither[error, string] {
				logger(s)
				return Right[error]("end")
			}),
		)

		output := result()
		assert.Equal(t, E.Right[error]("end"), output)
		assert.Equal(t, []string{"start", "middle"}, log)
	})

	t.Run("error propagation through consumer chain", func(t *testing.T) {
		var captured []int
		consumer := func(x int) {
			captured = append(captured, x)
		}

		err := errors.New("early error")
		result := F.Pipe3(
			Left[int](err),
			ChainConsumer[error](consumer),
			Map[error](func(struct{}) int { return 100 }),
			ChainConsumer[error](consumer),
		)

		output := result()
		assert.True(t, E.IsLeft(output))
		_, leftErr := E.Unwrap(output)
		assert.Equal(t, err, leftErr)
		assert.Empty(t, captured) // Consumer never called
	})

	t.Run("consumer with pointer types", func(t *testing.T) {
		var captured *int
		consumer := func(p *int) {
			captured = p
		}

		value := 42
		result := F.Pipe1(
			Right[error](&value),
			ChainConsumer[error](consumer),
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

		op := ChainConsumer[error](consumer)

		// Create multiple IOEithers and execute them
		for i := 1; i <= 5; i++ {
			io := op(Right[error](i))
			io()
		}

		assert.Equal(t, []int{1, 2, 3, 4, 5}, accumulated)
	})

	t.Run("consumer with map accumulation", func(t *testing.T) {
		counts := make(map[string]int)
		consumer := func(s string) {
			counts[s]++
		}

		op := ChainConsumer[error](consumer)

		words := []string{"hello", "world", "hello", "test", "world", "hello"}
		for _, word := range words {
			io := op(Right[error](word))
			io()
		}

		assert.Equal(t, 3, counts["hello"])
		assert.Equal(t, 2, counts["world"])
		assert.Equal(t, 1, counts["test"])
	})

	t.Run("lazy evaluation - consumer not called until IOEither executed", func(t *testing.T) {
		called := false
		consumer := func(x int) {
			called = true
		}

		// Create the IOEither but don't execute it
		io := F.Pipe1(
			Right[error](42),
			ChainConsumer[error](consumer),
		)

		// Consumer should not be called yet
		assert.False(t, called)

		// Now execute
		io()

		// Consumer should be called now
		assert.True(t, called)
	})

	t.Run("consumer with different error types", func(t *testing.T) {
		var captured int
		consumer := func(x int) {
			captured = x
		}

		// Test with string error type
		result1 := F.Pipe1(
			Right[string](100),
			ChainConsumer[string](consumer),
		)
		result1()
		assert.Equal(t, 100, captured)

		// Test with custom error type
		type CustomError struct {
			Code    int
			Message string
		}
		result2 := F.Pipe1(
			Right[CustomError](200),
			ChainConsumer[CustomError](consumer),
		)
		result2()
		assert.Equal(t, 200, captured)
	})

	t.Run("consumer in error recovery scenario", func(t *testing.T) {
		var successLog []int
		successConsumer := func(x int) {
			successLog = append(successLog, x)
		}

		result := F.Pipe2(
			Left[int](errors.New("initial error")),
			ChainLeft(func(e error) IOEither[error, int] {
				// Recover from error
				return Right[error](42)
			}),
			ChainConsumer[error](successConsumer),
		)

		output := result()
		assert.True(t, E.IsRight(output))
		assert.Equal(t, []int{42}, successLog)
	})

	t.Run("consumer composition with ChainFirst", func(t *testing.T) {
		var log []string
		logger := func(s string) {
			log = append(log, "Logged: "+s)
		}

		result := F.Pipe2(
			Right[error]("test"),
			ChainConsumer[error](logger),
			ChainFirst(func(_ struct{}) IOEither[error, int] {
				return Right[error](42)
			}),
		)

		output := result()
		assert.True(t, E.IsRight(output))
		assert.Equal(t, []string{"Logged: test"}, log)
	})
}
