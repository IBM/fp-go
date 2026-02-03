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
	"context"
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		Map[context.Context, error](utils.Double),
	)

	assert.Equal(t, E.Of[error](2), g(context.Background())())
}

func TestOrLeft(t *testing.T) {
	f := OrLeft[int](func(s string) readerio.ReaderIO[context.Context, string] {
		return readerio.Of[context.Context](s + "!")
	})

	g1 := F.Pipe1(
		Right[context.Context, string](1),
		f,
	)

	g2 := F.Pipe1(
		Left[context.Context, int]("a"),
		f,
	)

	assert.Equal(t, E.Of[string](1), g1(context.Background())())
	assert.Equal(t, E.Left[int]("a!"), g2(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context, error](utils.Double),
		Ap[int](Right[context.Context, error](1)),
	)

	assert.Equal(t, E.Right[error](2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		ChainReaderK[error](func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, E.Right[error]("1"), g(context.Background())())
}

func TestOrElse(t *testing.T) {
	type Config struct {
		retryLimit int
		fallback   int
	}

	// Test basic recovery from Left
	t.Run("Left value recovered", func(t *testing.T) {
		rioe := Left[Config, int](errors.New("error"))
		recover := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			return Of[Config, error](0) // default value
		})
		result := recover(rioe)(Config{})()
		assert.Equal(t, E.Right[error](0), result)
	})

	// Test Right value passes through unchanged
	t.Run("Right value unchanged", func(t *testing.T) {
		rioe := Right[Config, error](42)
		recover := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			return Of[Config, error](0)
		})
		result := recover(rioe)(Config{})()
		assert.Equal(t, E.Right[error](42), result)
	})

	// Test conditional recovery using config
	t.Run("Conditional recovery with config", func(t *testing.T) {
		rioe := Left[Config, int](errors.New("retryable"))
		recoverWithConfig := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			return func(cfg Config) IOEither[error, int] {
				if err.Error() == "retryable" && cfg.retryLimit > 0 {
					return IOE.Right[error](cfg.fallback)
				}
				return IOE.Left[int](err)
			}
		})
		result := recoverWithConfig(rioe)(Config{retryLimit: 3, fallback: 99})()
		assert.Equal(t, E.Right[error](99), result)
	})

	// Test error propagation
	t.Run("Error propagation", func(t *testing.T) {
		otherErr := errors.New("other error")
		rioe := Left[Config, int](otherErr)
		recoverSpecific := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			if err.Error() == "retryable" {
				return Of[Config, error](0)
			}
			return Left[Config, int](err) // propagate other errors
		})
		result := recoverSpecific(rioe)(Config{})()
		assert.Equal(t, E.Left[int](otherErr), result)
	})

	// Test chaining multiple OrElse operations
	t.Run("Chaining OrElse operations", func(t *testing.T) {
		firstRecover := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			if err.Error() == "error1" {
				return Of[Config, error](1)
			}
			return Left[Config, int](err)
		})
		secondRecover := OrElse(func(err error) ReaderIOEither[Config, error, int] {
			if err.Error() == "error2" {
				return Of[Config, error](2)
			}
			return Left[Config, int](err)
		})

		result1 := F.Pipe1(Left[Config, int](errors.New("error1")), firstRecover)(Config{})()
		assert.Equal(t, E.Right[error](1), result1)

		result2 := F.Pipe1(Left[Config, int](errors.New("error2")), F.Flow2(firstRecover, secondRecover))(Config{})()
		assert.Equal(t, E.Right[error](2), result2)
	})
}

func TestOrElseWFunc(t *testing.T) {
	type Config struct {
		retryEnabled bool
	}

	// Test with Right - should pass through unchanged
	t.Run("Right passes through", func(t *testing.T) {
		rioe := Right[Config, string](42)
		handler := OrElse(func(err string) ReaderIOEither[Config, int, int] {
			return Left[Config, int](999)
		})
		result := handler(rioe)(Config{retryEnabled: true})()
		assert.Equal(t, E.Right[int](42), result)
	})

	// Test with Left - error type widening
	t.Run("Left with error type widening", func(t *testing.T) {
		rioe := Left[Config, int]("network error")
		handler := OrElse(func(err string) ReaderIOEither[Config, int, int] {
			return func(cfg Config) IOEither[int, int] {
				if cfg.retryEnabled {
					return IOE.Right[int](100)
				}
				return IOE.Left[int](404)
			}
		})
		result := handler(rioe)(Config{retryEnabled: true})()
		assert.Equal(t, E.Right[int](100), result)
	})
}

func TestChainLeftFunc(t *testing.T) {
	type Config struct {
		errorCode int
	}

	// Test with Right - should pass through unchanged
	t.Run("Right passes through", func(t *testing.T) {
		g := F.Pipe1(
			Right[Config, string](42),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				return Left[Config, int](999)
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Right[int](42), result)
	})

	// Test with Left - error transformation with config
	t.Run("Left transforms error with config", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("error"),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				return func(cfg Config) IOEither[int, int] {
					return IOE.Left[int](cfg.errorCode)
				}
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Left[int](500), result)
	})

	// Test with Left - successful recovery
	t.Run("Left recovers successfully", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("recoverable"),
			ChainLeft(func(err string) ReaderIOEither[Config, int, int] {
				if err == "recoverable" {
					return Right[Config, int](999)
				}
				return Left[Config, int](0)
			}),
		)
		result := g(Config{errorCode: 500})()
		assert.Equal(t, E.Right[int](999), result)
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
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				logged = true
				return Right[Config, int]("logged")
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Right[string](42), result)
		assert.False(t, logged)
	})

	// Test with Left - calls function but preserves original error
	t.Run("Left calls function but preserves error", func(t *testing.T) {
		logged = false
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				return func(cfg Config) IOEither[int, string] {
					if cfg.logEnabled {
						logged = true
					}
					return IOE.Right[int]("side effect done")
				}
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Left[int]("original error"), result)
		assert.True(t, logged)
	})

	// Test with Left - preserves original error even if side effect fails
	t.Run("Left preserves error even if side effect fails", func(t *testing.T) {
		g := F.Pipe1(
			Left[Config, int]("original error"),
			ChainFirstLeft[int](func(err string) ReaderIOEither[Config, int, string] {
				return Left[Config, string](999) // Side effect fails
			}),
		)
		result := g(Config{logEnabled: true})()
		assert.Equal(t, E.Left[int]("original error"), result)
	})
}

func TestTapLeft(t *testing.T) {
	// TapLeft is an alias for ChainFirstLeft, so just a basic sanity test
	type Config struct{}

	sideEffectRan := false

	g := F.Pipe1(
		Left[Config, int]("error"),
		TapLeft[int](func(err string) ReaderIOEither[Config, string, int] {
			sideEffectRan = true
			return Right[Config, string](0)
		}),
	)

	result := g(Config{})()
	assert.Equal(t, E.Left[int]("error"), result)
	assert.True(t, sideEffectRan)
}

func TestReadIOEither(t *testing.T) {
	type Config struct {
		baseURL string
		timeout int
	}

	// Test with Right IOEither - should execute ReaderIOEither with the environment
	t.Run("Right IOEither provides environment", func(t *testing.T) {
		// IOEither that successfully produces a config
		configIO := IOE.Right[error](Config{baseURL: "https://api.example.com", timeout: 30})

		// ReaderIOEither that uses the config
		computation := func(cfg Config) IOEither[error, string] {
			return IOE.Right[error](cfg.baseURL + "/users")
		}

		// Execute using ReadIOEither
		result := ReadIOEither[string](configIO)(computation)()
		assert.Equal(t, E.Right[error]("https://api.example.com/users"), result)
	})

	// Test with Left IOEither - should propagate error without executing ReaderIOEither
	t.Run("Left IOEither propagates error", func(t *testing.T) {
		configError := errors.New("failed to load config")
		configIO := IOE.Left[Config](configError)

		executed := false
		computation := func(cfg Config) IOEither[error, string] {
			executed = true
			return IOE.Right[error]("should not execute")
		}

		result := ReadIOEither[string](configIO)(computation)()
		assert.Equal(t, E.Left[string](configError), result)
		assert.False(t, executed, "ReaderIOEither should not execute when IOEither is Left")
	})

	// Test with Right IOEither but ReaderIOEither fails
	t.Run("Right IOEither but ReaderIOEither fails", func(t *testing.T) {
		configIO := IOE.Right[error](Config{baseURL: "https://api.example.com", timeout: 30})

		computationError := errors.New("computation failed")
		computation := func(cfg Config) IOEither[error, string] {
			// Use the config but fail
			if cfg.timeout < 60 {
				return IOE.Left[string](computationError)
			}
			return IOE.Right[error]("success")
		}

		result := ReadIOEither[string](configIO)(computation)()
		assert.Equal(t, E.Left[string](computationError), result)
	})

	// Test chaining with ReadIOEither
	t.Run("Chaining with ReadIOEither", func(t *testing.T) {
		// First get the config
		configIO := IOE.Right[error](Config{baseURL: "https://api.example.com", timeout: 30})

		// Chain multiple operations
		result := F.Pipe2(
			Of[Config, error](10),
			Map[Config, error](func(x int) int { return x * 2 }),
			ReadIOEither[int](configIO),
		)()

		assert.Equal(t, E.Right[error](20), result)
	})

	// Test with complex error type
	t.Run("Complex error type", func(t *testing.T) {
		type AppError struct {
			Code    int
			Message string
		}

		configIO := IOE.Left[Config](AppError{Code: 500, Message: "Internal error"})

		computation := func(cfg Config) IOEither[AppError, string] {
			return IOE.Right[AppError]("success")
		}

		result := ReadIOEither[string](configIO)(computation)()
		assert.Equal(t, E.Left[string](AppError{Code: 500, Message: "Internal error"}), result)
	})
}

func TestReadIO(t *testing.T) {
	type Config struct {
		baseURL string
		version string
	}

	// Test basic execution - IO provides environment
	t.Run("IO provides environment successfully", func(t *testing.T) {
		// IO that produces a config (cannot fail)
		configIO := func() Config {
			return Config{baseURL: "https://api.example.com", version: "v1"}
		}

		// ReaderIOEither that uses the config
		computation := func(cfg Config) IOEither[error, string] {
			return IOE.Right[error](cfg.baseURL + "/" + cfg.version)
		}

		result := ReadIO[error, string](configIO)(computation)()
		assert.Equal(t, E.Right[error]("https://api.example.com/v1"), result)
	})

	// Test when ReaderIOEither fails
	t.Run("ReaderIOEither fails after IO succeeds", func(t *testing.T) {
		configIO := func() Config {
			return Config{baseURL: "https://api.example.com", version: "v1"}
		}

		computationError := errors.New("validation failed")
		computation := func(cfg Config) IOEither[error, string] {
			// Validate config
			if cfg.version != "v2" {
				return IOE.Left[string](computationError)
			}
			return IOE.Right[error]("success")
		}

		result := ReadIO[error, string](configIO)(computation)()
		assert.Equal(t, E.Left[string](computationError), result)
	})

	// Test with side effects in IO
	t.Run("IO with side effects", func(t *testing.T) {
		counter := 0
		configIO := func() Config {
			counter++
			return Config{baseURL: fmt.Sprintf("https://api%d.example.com", counter), version: "v1"}
		}

		computation := func(cfg Config) IOEither[error, string] {
			return IOE.Right[error](cfg.baseURL)
		}

		result := ReadIO[error, string](configIO)(computation)()
		assert.Equal(t, E.Right[error]("https://api1.example.com"), result)
		assert.Equal(t, 1, counter, "IO should execute exactly once")
	})

	// Test chaining with ReadIO
	t.Run("Chaining with ReadIO", func(t *testing.T) {
		configIO := func() Config {
			return Config{baseURL: "https://api.example.com", version: "v1"}
		}

		result := F.Pipe2(
			Of[Config, error](42),
			Map[Config, error](func(x int) string { return fmt.Sprintf("value-%d", x) }),
			ReadIO[error, string](configIO),
		)()

		assert.Equal(t, E.Right[error]("value-42"), result)
	})

	// Test with different error types
	t.Run("Different error types", func(t *testing.T) {
		configIO := func() int {
			return 100
		}

		computation := func(cfg int) IOEither[string, int] {
			if cfg < 200 {
				return IOE.Left[int]("value too low")
			}
			return IOE.Right[string](cfg)
		}

		result := ReadIO[string, int](configIO)(computation)()
		assert.Equal(t, E.Left[int]("value too low"), result)
	})

	// Test ReadIO vs ReadIOEither - ReadIO cannot fail during environment loading
	t.Run("ReadIO always provides environment", func(t *testing.T) {
		// This demonstrates that ReadIO's IO always succeeds
		configIO := func() Config {
			// Even if we wanted to fail here, we can't - IO cannot fail
			return Config{baseURL: "fallback", version: "v0"}
		}

		executed := false
		computation := func(cfg Config) IOEither[error, string] {
			executed = true
			return IOE.Right[error](cfg.baseURL)
		}

		result := ReadIO[error, string](configIO)(computation)()
		assert.Equal(t, E.Right[error]("fallback"), result)
		assert.True(t, executed, "ReaderIOEither should always execute with ReadIO")
	})

	// Test with complex computation
	t.Run("Complex computation with environment", func(t *testing.T) {
		type Env struct {
			multiplier int
			offset     int
		}

		envIO := func() Env {
			return Env{multiplier: 3, offset: 10}
		}

		computation := func(env Env) IOEither[error, int] {
			return func() Either[error, int] {
				// Simulate some computation using the environment
				result := env.multiplier*5 + env.offset
				if result > 20 {
					return E.Right[error](result)
				}
				return E.Left[int](errors.New("result too small"))
			}
		}

		result := ReadIO[error, int](envIO)(computation)()
		assert.Equal(t, E.Right[error](25), result)
	})
}

// TestChainLeftIdenticalToOrElse proves that ChainLeft and OrElse are identical functions.
// This test verifies that both functions produce the same results for all scenarios with reader context:
// - Left values with error recovery using reader context
// - Left values with error transformation
// - Right values passing through unchanged
// - Error type widening
func TestChainLeftIdenticalToOrElse(t *testing.T) {
	type Config struct {
		fallbackValue int
		retryEnabled  bool
	}

	// Test 1: Left value with error recovery using reader context
	t.Run("Left value recovery with reader - ChainLeft equals OrElse", func(t *testing.T) {
		recoveryFn := func(e string) ReaderIOEither[Config, string, int] {
			if e == "recoverable" {
				return func(cfg Config) IOE.IOEither[string, int] {
					return IOE.Right[string](cfg.fallbackValue)
				}
			}
			return Left[Config, int](e)
		}

		input := Left[Config, int]("recoverable")
		cfg := Config{fallbackValue: 42, retryEnabled: true}

		// Using ChainLeft
		resultChainLeft := ChainLeft(recoveryFn)(input)(cfg)()

		// Using OrElse
		resultOrElse := OrElse(recoveryFn)(input)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Right[string](42), resultChainLeft)
	})

	// Test 2: Left value with error transformation
	t.Run("Left value transformation - ChainLeft equals OrElse", func(t *testing.T) {
		transformFn := func(e string) ReaderIOEither[Config, string, int] {
			return Left[Config, int]("transformed: " + e)
		}

		input := Left[Config, int]("original error")
		cfg := Config{fallbackValue: 0, retryEnabled: false}

		// Using ChainLeft
		resultChainLeft := ChainLeft(transformFn)(input)(cfg)()

		// Using OrElse
		resultOrElse := OrElse(transformFn)(input)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Left[int]("transformed: original error"), resultChainLeft)
	})

	// Test 3: Right value - both should pass through unchanged
	t.Run("Right value passthrough - ChainLeft equals OrElse", func(t *testing.T) {
		handlerFn := func(e string) ReaderIOEither[Config, string, int] {
			return Left[Config, int]("should not be called")
		}

		input := Right[Config, string](100)
		cfg := Config{fallbackValue: 0, retryEnabled: false}

		// Using ChainLeft
		resultChainLeft := ChainLeft(handlerFn)(input)(cfg)()

		// Using OrElse
		resultOrElse := OrElse(handlerFn)(input)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Right[string](100), resultChainLeft)
	})

	// Test 4: Error type widening
	t.Run("Error type widening - ChainLeft equals OrElse", func(t *testing.T) {
		widenFn := func(e string) ReaderIOEither[Config, int, int] {
			return Left[Config, int](404)
		}

		input := Left[Config, int]("not found")
		cfg := Config{fallbackValue: 0, retryEnabled: false}

		// Using ChainLeft
		resultChainLeft := ChainLeft(widenFn)(input)(cfg)()

		// Using OrElse
		resultOrElse := OrElse(widenFn)(input)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Left[int](404), resultChainLeft)
	})

	// Test 5: Composition in pipeline with reader context
	t.Run("Pipeline composition with reader - ChainLeft equals OrElse", func(t *testing.T) {
		recoveryFn := func(e string) ReaderIOEither[Config, string, int] {
			if e == "network error" {
				return func(cfg Config) IOE.IOEither[string, int] {
					if cfg.retryEnabled {
						return IOE.Right[string](cfg.fallbackValue)
					}
					return IOE.Left[int]("retry disabled")
				}
			}
			return Left[Config, int](e)
		}

		input := Left[Config, int]("network error")
		cfg := Config{fallbackValue: 99, retryEnabled: true}

		// Using ChainLeft in pipeline
		resultChainLeft := F.Pipe1(input, ChainLeft(recoveryFn))(cfg)()

		// Using OrElse in pipeline
		resultOrElse := F.Pipe1(input, OrElse(recoveryFn))(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Right[string](99), resultChainLeft)
	})

	// Test 6: Multiple chained operations with reader context
	t.Run("Multiple operations with reader - ChainLeft equals OrElse", func(t *testing.T) {
		handler1 := func(e string) ReaderIOEither[Config, string, int] {
			if e == "error1" {
				return Right[Config, string](1)
			}
			return Left[Config, int](e)
		}

		handler2 := func(e string) ReaderIOEither[Config, string, int] {
			if e == "error2" {
				return func(cfg Config) IOE.IOEither[string, int] {
					return IOE.Right[string](cfg.fallbackValue)
				}
			}
			return Left[Config, int](e)
		}

		input := Left[Config, int]("error2")
		cfg := Config{fallbackValue: 2, retryEnabled: false}

		// Using ChainLeft
		resultChainLeft := F.Pipe2(
			input,
			ChainLeft(handler1),
			ChainLeft(handler2),
		)(cfg)()

		// Using OrElse
		resultOrElse := F.Pipe2(
			input,
			OrElse(handler1),
			OrElse(handler2),
		)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Right[string](2), resultChainLeft)
	})

	// Test 7: Reader context is properly threaded through both functions
	t.Run("Reader context threading - ChainLeft equals OrElse", func(t *testing.T) {
		var chainLeftCfg, orElseCfg *Config

		recoveryFn := func(e string) ReaderIOEither[Config, string, int] {
			return func(cfg Config) IOE.IOEither[string, int] {
				// Capture the config to verify it's passed correctly
				if chainLeftCfg == nil {
					chainLeftCfg = &cfg
				} else {
					orElseCfg = &cfg
				}
				return IOE.Right[string](cfg.fallbackValue)
			}
		}

		input := Left[Config, int]("error")
		cfg := Config{fallbackValue: 123, retryEnabled: true}

		// Using ChainLeft
		resultChainLeft := ChainLeft(recoveryFn)(input)(cfg)()

		// Using OrElse
		resultOrElse := OrElse(recoveryFn)(input)(cfg)()

		// Both should produce identical results
		assert.Equal(t, resultOrElse, resultChainLeft)
		assert.Equal(t, E.Right[string](123), resultChainLeft)

		// Verify both received the same config
		assert.NotNil(t, chainLeftCfg)
		assert.NotNil(t, orElseCfg)
		assert.Equal(t, *chainLeftCfg, *orElseCfg)
		assert.Equal(t, cfg, *chainLeftCfg)
	})
}
