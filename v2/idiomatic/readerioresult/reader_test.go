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

package readerioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

type TestConfig struct {
	Multiplier int
	Prefix     string
}

func TestFromIOResult(t *testing.T) {
	t.Run("lifts successful IOResult", func(t *testing.T) {
		ioResult := ioresult.Of(42)

		readerIOResult := FromIOResult[TestConfig](ioResult)
		cfg := TestConfig{Multiplier: 5}

		result, err := readerIOResult(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("lifts failing IOResult", func(t *testing.T) {
		expectedError := errors.New("io error")
		ioResult := ioresult.Left[int](expectedError)

		readerIOResult := FromIOResult[TestConfig](ioResult)
		cfg := TestConfig{Multiplier: 5}

		_, err := readerIOResult(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("ignores environment", func(t *testing.T) {
		ioResult := ioresult.Of("constant")

		readerIOResult := FromIOResult[TestConfig](ioResult)

		// Different configs should produce same result
		result1, _ := readerIOResult(TestConfig{Multiplier: 1})()
		result2, _ := readerIOResult(TestConfig{Multiplier: 100})()

		assert.Equal(t, result1, result2)
		assert.Equal(t, "constant", result1)
	})
}

func TestRightIO(t *testing.T) {
	t.Run("lifts IO as success", func(t *testing.T) {
		counter := 0
		io := func() int {
			counter++
			return counter
		}

		readerIOResult := RightIO[TestConfig](io)
		cfg := TestConfig{}

		result, err := readerIOResult(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 1, result)
		assert.Equal(t, 1, counter)
	})

	t.Run("always succeeds", func(t *testing.T) {
		io := io.Of("success")

		readerIOResult := RightIO[TestConfig](io)
		cfg := TestConfig{}

		result, err := readerIOResult(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "success", result)
	})
}

func TestLeftIO(t *testing.T) {
	t.Run("lifts IO error as failure", func(t *testing.T) {
		expectedError := errors.New("io error")
		io := io.Of(expectedError)

		readerIOResult := LeftIO[TestConfig, int](io)
		cfg := TestConfig{}

		_, err := readerIOResult(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("always fails", func(t *testing.T) {
		io := io.Of(errors.New("always fails"))

		readerIOResult := LeftIO[TestConfig, string](io)
		cfg := TestConfig{}

		_, err := readerIOResult(cfg)()
		assert.Error(t, err)
	})
}

func TestFromReaderIO(t *testing.T) {
	t.Run("lifts ReaderIO as success", func(t *testing.T) {
		readerIO := func(cfg TestConfig) func() int {
			return func() int {
				return cfg.Multiplier * 10
			}
		}

		readerIOResult := FromReaderIO(readerIO)
		cfg := TestConfig{Multiplier: 5}

		result, err := readerIOResult(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 50, result)
	})

	t.Run("uses environment", func(t *testing.T) {
		readerIO := func(cfg TestConfig) func() string {
			return func() string {
				return fmt.Sprintf("%s:%d", cfg.Prefix, cfg.Multiplier)
			}
		}

		readerIOResult := FromReaderIO(readerIO)

		result1, _ := readerIOResult(TestConfig{Prefix: "A", Multiplier: 1})()
		result2, _ := readerIOResult(TestConfig{Prefix: "B", Multiplier: 2})()

		assert.Equal(t, "A:1", result1)
		assert.Equal(t, "B:2", result2)
	})
}

func TestMonadMap(t *testing.T) {
	t.Run("transforms success value", func(t *testing.T) {
		getValue := Right[TestConfig](10)
		double := N.Mul(2)

		result := MonadMap(getValue, double)
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedError := errors.New("error")
		getValue := Left[TestConfig, int](expectedError)
		double := N.Mul(2)

		result := MonadMap(getValue, double)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		getValue := Right[TestConfig](5)

		result := F.Pipe3(
			getValue,
			Map[TestConfig](N.Mul(2)),
			Map[TestConfig](N.Add(3)),
			Map[TestConfig](S.Format[int]("result:%d")),
		)

		cfg := TestConfig{}
		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "result:13", value)
	})
}

func TestMap(t *testing.T) {
	t.Run("curried version works in pipeline", func(t *testing.T) {
		double := Map[TestConfig](N.Mul(2))
		getValue := Right[TestConfig](10)

		result := F.Pipe1(getValue, double)
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})
}

func TestMonadMapTo(t *testing.T) {
	t.Run("replaces value with constant", func(t *testing.T) {
		getValue := Right[TestConfig](10)

		result := MonadMapTo(getValue, "constant")
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "constant", value)
	})

	t.Run("propagates error", func(t *testing.T) {
		expectedError := errors.New("error")
		getValue := Left[TestConfig, int](expectedError)

		result := MonadMapTo(getValue, "constant")
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestMonadChain(t *testing.T) {
	t.Run("sequences dependent computations", func(t *testing.T) {
		getUser := Right[TestConfig](User{ID: 1, Name: "Alice"})
		getUserPosts := func(user User) ReaderIOResult[TestConfig, []string] {
			return func(cfg TestConfig) IOResult[[]string] {
				return func() ([]string, error) {
					return []string{
						fmt.Sprintf("Post 1 by %s", user.Name),
						fmt.Sprintf("Post 2 by %s", user.Name),
					}, nil
				}
			}
		}

		result := MonadChain(getUser, getUserPosts)
		cfg := TestConfig{}

		posts, err := result(cfg)()
		assert.NoError(t, err)
		assert.Len(t, posts, 2)
		assert.Contains(t, posts[0], "Alice")
	})

	t.Run("propagates first error", func(t *testing.T) {
		expectedError := errors.New("first error")
		getUser := Left[TestConfig, User](expectedError)
		getUserPosts := func(user User) ReaderIOResult[TestConfig, []string] {
			return Right[TestConfig]([]string{"should not be called"})
		}

		result := MonadChain(getUser, getUserPosts)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("propagates second error", func(t *testing.T) {
		expectedError := errors.New("second error")
		getUser := Right[TestConfig](User{ID: 1, Name: "Alice"})
		getUserPosts := func(user User) ReaderIOResult[TestConfig, []string] {
			return Left[TestConfig, []string](expectedError)
		}

		result := MonadChain(getUser, getUserPosts)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("shares environment", func(t *testing.T) {
		getValue := Ask[TestConfig]()
		transform := func(cfg TestConfig) ReaderIOResult[TestConfig, string] {
			return func(cfg2 TestConfig) IOResult[string] {
				return func() (string, error) {
					// Both should see the same config
					assert.Equal(t, cfg.Multiplier, cfg2.Multiplier)
					return fmt.Sprintf("%s:%d", cfg.Prefix, cfg.Multiplier), nil
				}
			}
		}

		result := MonadChain(getValue, transform)
		cfg := TestConfig{Prefix: "test", Multiplier: 42}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "test:42", value)
	})
}

func TestChain(t *testing.T) {
	t.Run("curried version works in pipeline", func(t *testing.T) {
		double := func(x int) ReaderIOResult[TestConfig, int] {
			return Right[TestConfig](x * 2)
		}

		result := F.Pipe1(
			Right[TestConfig](10),
			Chain(double),
		)

		cfg := TestConfig{}
		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 20, value)
	})
}

func TestMonadChainFirst(t *testing.T) {
	t.Run("executes side effect but returns first value", func(t *testing.T) {
		sideEffectCalled := false
		getUser := Right[TestConfig](User{ID: 1, Name: "Alice"})
		logUser := func(user User) ReaderIOResult[TestConfig, string] {
			return func(cfg TestConfig) IOResult[string] {
				return func() (string, error) {
					sideEffectCalled = true
					return "logged", nil
				}
			}
		}

		result := MonadChainFirst(getUser, logUser)
		cfg := TestConfig{}

		user, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "Alice", user.Name)
		assert.True(t, sideEffectCalled)
	})

	t.Run("propagates first error", func(t *testing.T) {
		expectedError := errors.New("first error")
		getUser := Left[TestConfig, User](expectedError)
		logUser := func(user User) ReaderIOResult[TestConfig, string] {
			return Right[TestConfig]("should not be called")
		}

		result := MonadChainFirst(getUser, logUser)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("propagates second error", func(t *testing.T) {
		expectedError := errors.New("second error")
		getUser := Right[TestConfig](User{ID: 1, Name: "Alice"})
		logUser := func(user User) ReaderIOResult[TestConfig, string] {
			return Left[TestConfig, string](expectedError)
		}

		result := MonadChainFirst(getUser, logUser)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestMonadAp(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		fab := Right[TestConfig](N.Mul(2))
		fa := Right[TestConfig](21)

		result := MonadAp(fab, fa)
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("propagates function error", func(t *testing.T) {
		expectedError := errors.New("function error")
		fab := Left[TestConfig, func(int) int](expectedError)
		fa := Right[TestConfig](21)

		result := MonadAp(fab, fa)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("propagates value error", func(t *testing.T) {
		expectedError := errors.New("value error")
		fab := Right[TestConfig](N.Mul(2))
		fa := Left[TestConfig, int](expectedError)

		result := MonadAp(fab, fa)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestRightAndLeft(t *testing.T) {
	t.Run("Right creates successful value", func(t *testing.T) {
		result := Right[TestConfig](42)
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("Left creates error", func(t *testing.T) {
		expectedError := errors.New("error")
		result := Left[TestConfig, int](expectedError)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("Of is alias for Right", func(t *testing.T) {
		result1 := Right[TestConfig](42)
		result2 := Of[TestConfig](42)
		cfg := TestConfig{}

		value1, _ := result1(cfg)()
		value2, _ := result2(cfg)()

		assert.Equal(t, value1, value2)
	})
}

func TestFlatten(t *testing.T) {
	t.Run("removes one level of nesting", func(t *testing.T) {
		inner := Right[TestConfig](42)
		outer := Right[TestConfig](inner)

		result := Flatten(outer)
		cfg := TestConfig{}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("propagates outer error", func(t *testing.T) {
		expectedError := errors.New("outer error")
		outer := Left[TestConfig, ReaderIOResult[TestConfig, int]](expectedError)

		result := Flatten(outer)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("propagates inner error", func(t *testing.T) {
		expectedError := errors.New("inner error")
		inner := Left[TestConfig, int](expectedError)
		outer := Right[TestConfig](inner)

		result := Flatten(outer)
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})
}

func TestAsk(t *testing.T) {
	t.Run("retrieves environment", func(t *testing.T) {
		result := Ask[TestConfig]()
		cfg := TestConfig{Multiplier: 42, Prefix: "test"}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, cfg, value)
	})

	t.Run("always succeeds", func(t *testing.T) {
		result := Ask[TestConfig]()
		cfg := TestConfig{}

		_, err := result(cfg)()
		assert.NoError(t, err)
	})
}

func TestAsks(t *testing.T) {
	t.Run("extracts value from environment", func(t *testing.T) {
		getMultiplier := func(cfg TestConfig) int {
			return cfg.Multiplier
		}

		result := Asks(getMultiplier)
		cfg := TestConfig{Multiplier: 42}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("works with different extractors", func(t *testing.T) {
		getPrefix := func(cfg TestConfig) string {
			return cfg.Prefix
		}

		result := Asks(getPrefix)
		cfg := TestConfig{Prefix: "test"}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "test", value)
	})
}

func TestLocal(t *testing.T) {
	t.Run("transforms environment", func(t *testing.T) {
		// Computation that uses TestConfig
		computation := func(cfg TestConfig) IOResult[string] {
			return func() (string, error) {
				return fmt.Sprintf("%s:%d", cfg.Prefix, cfg.Multiplier), nil
			}
		}

		// Transform function that modifies the config
		transform := func(cfg TestConfig) TestConfig {
			return TestConfig{
				Prefix:     "modified-" + cfg.Prefix,
				Multiplier: cfg.Multiplier * 2,
			}
		}

		result := Local[string](transform)(computation)
		cfg := TestConfig{Prefix: "test", Multiplier: 5}

		value, err := result(cfg)()
		assert.NoError(t, err)
		assert.Equal(t, "modified-test:10", value)
	})
}

func TestRead(t *testing.T) {
	t.Run("provides environment to computation", func(t *testing.T) {
		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return cfg.Multiplier * 10, nil
			}
		}

		cfg := TestConfig{Multiplier: 5}
		result := Read[int](cfg)(computation)

		value, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 50, value)
	})
}

// Helper type for tests
type User struct {
	ID   int
	Name string
}

func TestReadIO(t *testing.T) {
	t.Run("executes computation with IO environment", func(t *testing.T) {
		// IO that produces the config
		loadConfig := func() TestConfig {
			return TestConfig{Multiplier: 7, Prefix: "loaded"}
		}

		// Computation that uses the config
		computation := func(cfg TestConfig) IOResult[string] {
			return func() (string, error) {
				return fmt.Sprintf("%s:%d", cfg.Prefix, cfg.Multiplier), nil
			}
		}

		result := ReadIO[string](loadConfig)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, "loaded:7", value)
	})

	t.Run("executes IO before computation", func(t *testing.T) {
		executionOrder := []string{}

		// IO that tracks execution
		loadConfig := func() TestConfig {
			executionOrder = append(executionOrder, "load-config")
			return TestConfig{Multiplier: 5}
		}

		// Computation that tracks execution
		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				executionOrder = append(executionOrder, "compute")
				return cfg.Multiplier * 10, nil
			}
		}

		result := ReadIO[int](loadConfig)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, 50, value)
		assert.Equal(t, []string{"load-config", "compute"}, executionOrder)
	})

	t.Run("propagates computation error", func(t *testing.T) {
		expectedError := errors.New("computation failed")

		loadConfig := func() TestConfig {
			return TestConfig{Multiplier: 5}
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return 0, expectedError
			}
		}

		result := ReadIO[int](loadConfig)(computation)
		_, err := result()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("works with different environment types", func(t *testing.T) {
		// Using a simple string as environment
		loadEnv := func() string {
			return "test-env"
		}

		computation := func(env string) IOResult[string] {
			return func() (string, error) {
				return "env:" + env, nil
			}
		}

		result := ReadIO[string](loadEnv)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, "env:test-env", value)
	})

	t.Run("IO is executed on each call", func(t *testing.T) {
		counter := 0
		loadConfig := func() TestConfig {
			counter++
			return TestConfig{Multiplier: counter}
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return cfg.Multiplier, nil
			}
		}

		result := ReadIO[int](loadConfig)(computation)

		// First execution
		value1, _ := result()
		assert.Equal(t, 1, value1)

		// Second execution - IO runs again
		value2, _ := result()
		assert.Equal(t, 2, value2)
	})
}

func TestReadIOResult(t *testing.T) {
	t.Run("executes computation with successful IOResult environment", func(t *testing.T) {
		// IOResult that successfully produces config
		loadConfig := func() (TestConfig, error) {
			return TestConfig{Multiplier: 8, Prefix: "success"}, nil
		}

		// Computation that uses the config
		computation := func(cfg TestConfig) IOResult[string] {
			return func() (string, error) {
				return fmt.Sprintf("%s:%d", cfg.Prefix, cfg.Multiplier), nil
			}
		}

		result := ReadIOResult[string](loadConfig)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, "success:8", value)
	})

	t.Run("propagates environment loading error", func(t *testing.T) {
		expectedError := errors.New("failed to load config")

		// IOResult that fails to produce config
		loadConfig := func() (TestConfig, error) {
			return TestConfig{}, expectedError
		}

		// Computation should not be executed
		computationCalled := false
		computation := func(cfg TestConfig) IOResult[string] {
			return func() (string, error) {
				computationCalled = true
				return "should not reach here", nil
			}
		}

		result := ReadIOResult[string](loadConfig)(computation)
		_, err := result()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.False(t, computationCalled, "computation should not be called when environment loading fails")
	})

	t.Run("propagates computation error after successful environment load", func(t *testing.T) {
		expectedError := errors.New("computation failed")

		loadConfig := func() (TestConfig, error) {
			return TestConfig{Multiplier: 5}, nil
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return 0, expectedError
			}
		}

		result := ReadIOResult[int](loadConfig)(computation)
		_, err := result()

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("chains environment loading and computation", func(t *testing.T) {
		executionOrder := []string{}

		loadConfig := func() (TestConfig, error) {
			executionOrder = append(executionOrder, "load-config")
			return TestConfig{Multiplier: 3}, nil
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				executionOrder = append(executionOrder, "compute")
				return cfg.Multiplier * 10, nil
			}
		}

		result := ReadIOResult[int](loadConfig)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, 30, value)
		assert.Equal(t, []string{"load-config", "compute"}, executionOrder)
	})

	t.Run("works with validation in environment loading", func(t *testing.T) {
		// IOResult that validates config
		loadConfig := func() (TestConfig, error) {
			cfg := TestConfig{Multiplier: -1}
			if cfg.Multiplier < 0 {
				return TestConfig{}, errors.New("invalid multiplier: must be positive")
			}
			return cfg, nil
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return cfg.Multiplier * 10, nil
			}
		}

		result := ReadIOResult[int](loadConfig)(computation)
		_, err := result()

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid multiplier")
	})

	t.Run("IOResult is executed on each call", func(t *testing.T) {
		counter := 0
		loadConfig := func() (TestConfig, error) {
			counter++
			if counter == 1 {
				return TestConfig{}, errors.New("first attempt fails")
			}
			return TestConfig{Multiplier: counter}, nil
		}

		computation := func(cfg TestConfig) IOResult[int] {
			return func() (int, error) {
				return cfg.Multiplier, nil
			}
		}

		result := ReadIOResult[int](loadConfig)(computation)

		// First execution - fails
		_, err1 := result()
		assert.Error(t, err1)

		// Second execution - succeeds
		value2, err2 := result()
		assert.NoError(t, err2)
		assert.Equal(t, 2, value2)
	})

	t.Run("works with complex environment types", func(t *testing.T) {
		type DatabaseConfig struct {
			Host     string
			Port     int
			Username string
		}

		loadDBConfig := func() (DatabaseConfig, error) {
			// Simulate loading from environment variables
			return DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				Username: "admin",
			}, nil
		}

		computation := func(cfg DatabaseConfig) IOResult[string] {
			return func() (string, error) {
				return fmt.Sprintf("%s@%s:%d", cfg.Username, cfg.Host, cfg.Port), nil
			}
		}

		result := ReadIOResult[string](loadDBConfig)(computation)
		value, err := result()

		assert.NoError(t, err)
		assert.Equal(t, "admin@localhost:5432", value)
	})
}
