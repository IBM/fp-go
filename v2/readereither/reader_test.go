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
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, error](1),
		Map[MyContext, error](utils.Double),
	)

	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext, error](utils.Double),
		Ap[int](Of[MyContext, error](1)),
	)
	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestFlatten(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, string](Of[MyContext, string]("a")),
		Flatten[MyContext, string, string],
	)

	assert.Equal(t, ET.Of[string]("a"), g(defaultContext))
}

func TestChainLeftFunc(t *testing.T) {
	type Config struct {
		errorCode int
	}

	// Test with Right - should pass through unchanged
	t.Run("Right passes through", func(t *testing.T) {
		g := F.Pipe1(
			Right[Config, string](42),
			ChainLeft(func(err string) ReaderEither[Config, int, int] {
				return Left[Config, int](999)
			}),
		)
		result := g(Config{errorCode: 500})
		assert.Equal(t, ET.Right[int](42), result)
	})

	// Test with Left - error transformation with config
	t.Run("Left transforms error with config", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("error"),
			ChainLeft(func(err string) ReaderEither[Config, int, int] {
				return func(cfg Config) Either[int, int] {
					return ET.Left[int](cfg.errorCode)
				}
			}),
		)
		result := g(Config{errorCode: 500})
		assert.Equal(t, ET.Left[int](500), result)
	})

	// Test with Left - successful recovery
	t.Run("Left recovers successfully", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("recoverable"),
			ChainLeft(func(err string) ReaderEither[Config, int, int] {
				if err == "recoverable" {
					return Right[Config, int](999)
				}
				return Left[Config, int](0)
			}),
		)
		result := g(Config{errorCode: 500})
		assert.Equal(t, ET.Right[int](999), result)
	})
}

func TestChainFirstLeftFunc(t *testing.T) {
	type Config struct {
		logEnabled bool
	}

	logged := false

	// Test with Right - should not call function
	t.Run("Right does not call function", func(t *testing.T) {
		logged = false
		g := F.Pipe1(
			Right[Config, string](42),
			ChainFirstLeft[int](func(err string) ReaderEither[Config, int, string] {
				logged = true
				return Right[Config, int]("logged")
			}),
		)
		result := g(Config{logEnabled: true})
		assert.Equal(t, ET.Right[string](42), result)
		assert.False(t, logged)
	})

	// Test with Left - calls function but preserves original error
	t.Run("Left calls function but preserves error", func(t *testing.T) {
		logged = false
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderEither[Config, int, string] {
				return func(cfg Config) Either[int, string] {
					if cfg.logEnabled {
						logged = true
					}
					return ET.Right[int]("side effect done")
				}
			}),
		)
		result := g(Config{logEnabled: true})
		assert.Equal(t, ET.Left[int]("original error"), result)
		assert.True(t, logged)
	})

	// Test with Left - preserves original error even if side effect fails
	t.Run("Left preserves error even if side effect fails", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderEither[Config, int, string] {
				return Left[Config, string](999) // Side effect fails
			}),
		)
		result := g(Config{logEnabled: true})
		assert.Equal(t, ET.Left[int]("original error"), result)
	})
}

func TestTapLeftFunc(t *testing.T) {
	// TapLeft is an alias for ChainFirstLeft, so just a basic sanity test
	type Config struct{}

	sideEffectRan := false

	g := F.Pipe1(
		Left[Config, int]("error"),
		TapLeft[int](func(err string) ReaderEither[Config, string, int] {
			sideEffectRan = true
			return Right[Config, string](0)
		}),
	)

	result := g(Config{})
	assert.Equal(t, ET.Left[int]("error"), result)
	assert.True(t, sideEffectRan)
}

func TestOrElse(t *testing.T) {
	type Config struct {
		fallbackValue int
	}

	// Test OrElse with Right - should pass through unchanged
	rightValue := Of[Config, string](42)
	recover := OrElse(func(err string) ReaderEither[Config, string, int] {
		return Left[Config, int]("should not be called")
	})
	result := recover(rightValue)(Config{fallbackValue: 0})
	assert.Equal(t, ET.Right[string](42), result)

	// Test OrElse with Left - should recover with fallback
	leftValue := Left[Config, int]("not found")
	recoverWithFallback := OrElse(func(err string) ReaderEither[Config, string, int] {
		if err == "not found" {
			return func(cfg Config) ET.Either[string, int] {
				return ET.Right[string](cfg.fallbackValue)
			}
		}
		return Left[Config, int](err)
	})
	result = recoverWithFallback(leftValue)(Config{fallbackValue: 99})
	assert.Equal(t, ET.Right[string](99), result)

	// Test OrElse with Left - should propagate other errors
	leftValue = Left[Config, int]("fatal error")
	result = recoverWithFallback(leftValue)(Config{fallbackValue: 99})
	assert.Equal(t, ET.Left[int]("fatal error"), result)

	// Test error type widening
	type ValidationError struct{ field string }
	type AppError struct{ code int }

	validationErr := Left[Config, int](ValidationError{field: "username"})
	wideningRecover := OrElse(func(ve ValidationError) ReaderEither[Config, AppError, int] {
		if ve.field == "username" {
			return Right[Config, AppError](100)
		}
		return Left[Config, int](AppError{code: 400})
	})
	appResult := wideningRecover(validationErr)(Config{})
	assert.Equal(t, ET.Right[AppError](100), appResult)
}

func TestReadEither(t *testing.T) {
	type Config struct {
		apiKey string
		host   string
	}

	// Test with Right context - should execute the ReaderEither
	t.Run("Right context executes computation", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "secret", host: "localhost"})

		computation := func(cfg Config) Either[string, int] {
			if cfg.apiKey == "secret" {
				return ET.Right[string](42)
			}
			return ET.Left[int]("invalid key")
		}

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Right[string](42), result)
	})

	// Test with Right context but computation fails
	t.Run("Right context with failing computation", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "wrong", host: "localhost"})

		computation := func(cfg Config) Either[string, int] {
			if cfg.apiKey == "secret" {
				return ET.Right[string](42)
			}
			return ET.Left[int]("invalid key")
		}

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Left[int]("invalid key"), result)
	})

	// Test with Left context - should short-circuit without executing
	t.Run("Left context short-circuits", func(t *testing.T) {
		invalidConfig := ET.Left[Config]("config not found")

		executed := false
		computation := func(cfg Config) Either[string, int] {
			executed = true
			return ET.Right[string](42)
		}

		result := ReadEither[string, int](invalidConfig)(computation)
		assert.Equal(t, ET.Left[int]("config not found"), result)
		assert.False(t, executed, "computation should not be executed with Left context")
	})

	// Test with complex ReaderEither computation
	t.Run("Complex ReaderEither computation", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "secret", host: "api.example.com"})

		// A more complex computation using the config
		computation := F.Pipe2(
			Ask[Config, string](),
			Map[Config, string](func(cfg Config) string {
				return cfg.host + "/data"
			}),
			Chain(func(url string) ReaderEither[Config, string, int] {
				return func(cfg Config) Either[string, int] {
					if cfg.apiKey != "" {
						return ET.Right[string](len(url))
					}
					return ET.Left[int]("no API key")
				}
			}),
		)

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Right[string](20), result) // len("api.example.com/data") = 20
	})

	// Test error type consistency
	t.Run("Error type consistency", func(t *testing.T) {
		type AppError struct {
			code    int
			message string
		}

		configError := AppError{code: 404, message: "config not found"}
		invalidConfig := ET.Left[Config](configError)

		computation := func(cfg Config) Either[AppError, string] {
			return ET.Right[AppError]("success")
		}

		result := ReadEither[AppError, string](invalidConfig)(computation)
		assert.Equal(t, ET.Left[string](configError), result)
	})

	// Test with chained operations
	t.Run("Chained operations with ReadEither", func(t *testing.T) {
		config1 := ET.Right[string](Config{apiKey: "key1", host: "host1"})
		config2 := ET.Right[string](Config{apiKey: "key2", host: "host2"})

		computation := func(cfg Config) Either[string, string] {
			return ET.Right[string](cfg.host)
		}

		// Apply first config
		result1 := ReadEither[string, string](config1)(computation)
		assert.Equal(t, ET.Right[string]("host1"), result1)

		// Apply second config
		result2 := ReadEither[string, string](config2)(computation)
		assert.Equal(t, ET.Right[string]("host2"), result2)
	})

	// Test with FromReader
	t.Run("ReadEither with FromReader", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "secret", host: "localhost"})

		// Create a ReaderEither from a Reader
		readerComputation := func(cfg Config) int {
			return len(cfg.apiKey)
		}

		computation := FromReader[string](readerComputation)

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Right[string](6), result) // len("secret") = 6
	})

	// Test with Of (pure value)
	t.Run("ReadEither with pure value", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "secret", host: "localhost"})
		computation := Of[Config, string](100)

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Right[string](100), result)
	})

	// Test with Left computation
	t.Run("ReadEither with Left computation", func(t *testing.T) {
		validConfig := ET.Right[string](Config{apiKey: "secret", host: "localhost"})
		computation := Left[Config, int]("computation error")

		result := ReadEither[string, int](validConfig)(computation)
		assert.Equal(t, ET.Left[int]("computation error"), result)
	})

	// Test composition with Read
	t.Run("ReadEither vs Read comparison", func(t *testing.T) {
		config := Config{apiKey: "secret", host: "localhost"}
		computation := func(cfg Config) Either[string, int] {
			return ET.Right[string](len(cfg.apiKey))
		}

		// Using Read directly
		resultRead := Read[string, int](config)(computation)

		// Using ReadEither with Right
		resultReadEither := ReadEither[string, int](ET.Right[string](config))(computation)

		assert.Equal(t, resultRead, resultReadEither)
	})
}
