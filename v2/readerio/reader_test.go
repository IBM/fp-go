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
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	G "github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

type ReaderTestConfig struct {
	Value int
	Name  string
}

func TestFromIO(t *testing.T) {
	ioAction := G.Of(42)
	rio := FromIO[ReaderTestConfig](ioAction)

	config := ReaderTestConfig{Value: 10, Name: "test"}
	result := rio(config)()

	assert.Equal(t, 42, result)
}

func TestFromReader(t *testing.T) {
	reader := func(config ReaderTestConfig) int {
		return config.Value * 2
	}

	rio := FromReader(reader)
	config := ReaderTestConfig{Value: 5, Name: "test"}
	result := rio(config)()

	assert.Equal(t, 10, result)
}

func TestOf(t *testing.T) {
	rio := Of[ReaderTestConfig](100)
	config := ReaderTestConfig{Value: 1, Name: "test"}
	result := rio(config)()

	assert.Equal(t, 100, result)
}

func TestMonadMap(t *testing.T) {
	rio := Of[ReaderTestConfig](5)
	doubled := MonadMap(rio, N.Mul(2))

	config := ReaderTestConfig{Value: 1, Name: "test"}
	result := doubled(config)()

	assert.Equal(t, 10, result)
}

func TestMap(t *testing.T) {
	g := F.Pipe1(
		Of[context.Context](1),
		Map[context.Context](utils.Double),
	)

	assert.Equal(t, 2, g(context.Background())())
}

func TestMonadChain(t *testing.T) {
	rio1 := Of[ReaderTestConfig](5)
	result := MonadChain(rio1, func(n int) ReaderIO[ReaderTestConfig, int] {
		return Of[ReaderTestConfig](n * 3)
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestChain(t *testing.T) {
	result := F.Pipe1(
		Of[ReaderTestConfig](5),
		Chain(func(n int) ReaderIO[ReaderTestConfig, int] {
			return Of[ReaderTestConfig](n * 3)
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestMonadAp(t *testing.T) {
	fabIO := Of[ReaderTestConfig](N.Mul(2))
	faIO := Of[ReaderTestConfig](5)
	result := MonadAp(fabIO, faIO)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 10, result(config)())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[context.Context](utils.Double),
		Ap[int](Of[context.Context](1)),
	)

	assert.Equal(t, 2, g(context.Background())())
}

func TestMonadApSeq(t *testing.T) {
	fabIO := Of[ReaderTestConfig](N.Add(10))
	faIO := Of[ReaderTestConfig](5)
	result := MonadApSeq(fabIO, faIO)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestMonadApPar(t *testing.T) {
	fabIO := Of[ReaderTestConfig](N.Add(10))
	faIO := Of[ReaderTestConfig](5)
	result := MonadApPar(fabIO, faIO)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestAsk(t *testing.T) {
	rio := Ask[ReaderTestConfig]()
	config := ReaderTestConfig{Value: 42, Name: "test"}
	result := rio(config)()

	assert.Equal(t, config, result)
	assert.Equal(t, 42, result.Value)
	assert.Equal(t, "test", result.Name)
}

func TestAsks(t *testing.T) {
	rio := Asks(func(c ReaderTestConfig) int {
		return c.Value * 2
	})

	config := ReaderTestConfig{Value: 21, Name: "test"}
	result := rio(config)()

	assert.Equal(t, 42, result)
}

func TestMonadChainIOK(t *testing.T) {
	rio := Of[ReaderTestConfig](5)
	result := MonadChainIOK(rio, func(n int) G.IO[int] {
		return G.Of(n * 4)
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 20, result(config)())
}

func TestChainIOK(t *testing.T) {
	result := F.Pipe1(
		Of[ReaderTestConfig](5),
		ChainIOK[ReaderTestConfig](func(n int) G.IO[int] {
			return G.Of(n * 4)
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 20, result(config)())
}

func TestDefer(t *testing.T) {
	counter := 0
	rio := Defer(func() ReaderIO[ReaderTestConfig, int] {
		counter++
		return Of[ReaderTestConfig](counter)
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	result1 := rio(config)()
	result2 := rio(config)()

	assert.Equal(t, 1, result1)
	assert.Equal(t, 2, result2)
}

func TestMemoize(t *testing.T) {
	counter := 0
	rio := Of[ReaderTestConfig](0)
	memoized := Memoize(MonadMap(rio, func(int) int {
		counter++
		return counter
	}))

	config := ReaderTestConfig{Value: 1, Name: "test"}
	result1 := memoized(config)()
	result2 := memoized(config)()

	assert.Equal(t, 1, result1)
	assert.Equal(t, 1, result2) // Same value, memoized
}

func TestMemoizeWithDifferentContexts(t *testing.T) {
	rio := Ask[ReaderTestConfig]()
	memoized := Memoize(MonadMap(rio, func(c ReaderTestConfig) int {
		return c.Value
	}))

	config1 := ReaderTestConfig{Value: 10, Name: "first"}
	config2 := ReaderTestConfig{Value: 20, Name: "second"}

	result1 := memoized(config1)()
	result2 := memoized(config2)() // Should still return 10 (memoized from first call)

	assert.Equal(t, 10, result1)
	assert.Equal(t, 10, result2) // Memoized value from first context
}

func TestFlatten(t *testing.T) {
	nested := Of[ReaderTestConfig](Of[ReaderTestConfig](42))
	flattened := Flatten(nested)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	result := flattened(config)()

	assert.Equal(t, 42, result)
}

func TestMonadFlap(t *testing.T) {
	fabIO := Of[ReaderTestConfig](N.Mul(3))
	result := MonadFlap(fabIO, 7)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 21, result(config)())
}

func TestFlap(t *testing.T) {
	result := F.Pipe1(
		Of[ReaderTestConfig](N.Mul(3)),
		Flap[ReaderTestConfig, int](7),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	assert.Equal(t, 21, result(config)())
}

func TestComplexPipeline(t *testing.T) {
	// Test a complex pipeline combining multiple operations
	result := F.Pipe3(
		Ask[ReaderTestConfig](),
		Map[ReaderTestConfig](func(c ReaderTestConfig) int { return c.Value }),
		Chain(func(n int) ReaderIO[ReaderTestConfig, int] {
			return Of[ReaderTestConfig](n * 2)
		}),
		Map[ReaderTestConfig](N.Add(10)),
	)

	config := ReaderTestConfig{Value: 5, Name: "test"}
	assert.Equal(t, 20, result(config)()) // (5 * 2) + 10 = 20
}

func TestFromIOWithChain(t *testing.T) {
	ioAction := G.Of(10)

	result := F.Pipe1(
		FromIO[ReaderTestConfig](ioAction),
		Chain(func(n int) ReaderIO[ReaderTestConfig, int] {
			return MonadMap(Ask[ReaderTestConfig](), func(c ReaderTestConfig) int {
				return n + c.Value
			})
		}),
	)

	config := ReaderTestConfig{Value: 5, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestFromReaderWithMap(t *testing.T) {
	reader := func(c ReaderTestConfig) string {
		return c.Name
	}

	result := F.Pipe1(
		FromReader(reader),
		Map[ReaderTestConfig](func(s string) string {
			return s + " modified"
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "original"}
	assert.Equal(t, "original modified", result(config)())
}

func TestMonadMapTo(t *testing.T) {
	rio := Of[ReaderTestConfig](42)
	replaced := MonadMapTo(rio, "constant")

	config := ReaderTestConfig{Value: 10, Name: "test"}
	result := replaced(config)()

	assert.Equal(t, "constant", result)
}

func TestMapTo(t *testing.T) {
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		MapTo[ReaderTestConfig, int]("constant"),
	)

	config := ReaderTestConfig{Value: 10, Name: "test"}
	assert.Equal(t, "constant", result(config)())
}

func TestMapToExecutesSideEffects(t *testing.T) {
	t.Run("executes original ReaderIO and returns constant value", func(t *testing.T) {
		executed := false
		originalReaderIO := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				executed = true
				return 42
			}
		}

		// Apply MapTo operator
		toDone := MapTo[ReaderTestConfig, int]("done")
		resultReaderIO := toDone(originalReaderIO)

		// Execute the resulting ReaderIO
		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := resultReaderIO(config)()

		// Verify the constant value is returned
		assert.Equal(t, "done", result)
		// Verify the original ReaderIO WAS executed (side effect occurred)
		assert.True(t, executed, "original ReaderIO should be executed to allow side effects")
	})

	t.Run("executes ReaderIO in functional pipeline", func(t *testing.T) {
		executed := false
		step1 := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				executed = true
				return 100
			}
		}

		pipeline := F.Pipe1(
			step1,
			MapTo[ReaderTestConfig, int]("complete"),
		)

		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := pipeline(config)()

		assert.Equal(t, "complete", result)
		assert.True(t, executed, "original ReaderIO should be executed in pipeline")
	})

	t.Run("executes ReaderIO with side effects", func(t *testing.T) {
		sideEffectOccurred := false
		readerIOWithSideEffect := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				sideEffectOccurred = true
				return 42
			}
		}

		resultReaderIO := MapTo[ReaderTestConfig, int](true)(readerIOWithSideEffect)
		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := resultReaderIO(config)()

		assert.Equal(t, true, result)
		assert.True(t, sideEffectOccurred, "side effect should occur")
	})

	t.Run("executes complex computation with side effects", func(t *testing.T) {
		computationExecuted := false
		complexReaderIO := func(c ReaderTestConfig) G.IO[string] {
			return func() string {
				computationExecuted = true
				return "complex result"
			}
		}

		resultReaderIO := MapTo[ReaderTestConfig, string](99)(complexReaderIO)
		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := resultReaderIO(config)()

		assert.Equal(t, 99, result)
		assert.True(t, computationExecuted, "complex computation should be executed")
	})
}

func TestMonadMapToExecutesSideEffects(t *testing.T) {
	t.Run("executes original ReaderIO and returns constant value", func(t *testing.T) {
		executed := false
		originalReaderIO := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				executed = true
				return 42
			}
		}

		// Apply MonadMapTo
		resultReaderIO := MonadMapTo(originalReaderIO, "done")

		// Execute the resulting ReaderIO
		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := resultReaderIO(config)()

		// Verify the constant value is returned
		assert.Equal(t, "done", result)
		// Verify the original ReaderIO WAS executed (side effect occurred)
		assert.True(t, executed, "original ReaderIO should be executed to allow side effects")
	})

	t.Run("executes complex computation with side effects", func(t *testing.T) {
		computationExecuted := false
		complexReaderIO := func(c ReaderTestConfig) G.IO[string] {
			return func() string {
				computationExecuted = true
				return "complex result"
			}
		}

		resultReaderIO := MonadMapTo(complexReaderIO, 42)
		config := ReaderTestConfig{Value: 10, Name: "test"}
		result := resultReaderIO(config)()

		assert.Equal(t, 42, result)
		assert.True(t, computationExecuted, "complex computation should be executed")
	})

	t.Run("executes ReaderIO with logging side effect", func(t *testing.T) {
		logged := []string{}
		loggingReaderIO := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				logged = append(logged, "computation executed")
				return c.Value * 2
			}
		}

		resultReaderIO := MonadMapTo(loggingReaderIO, "result")
		config := ReaderTestConfig{Value: 5, Name: "test"}
		result := resultReaderIO(config)()

		assert.Equal(t, "result", result)
		assert.Equal(t, []string{"computation executed"}, logged)
	})

	t.Run("executes ReaderIO accessing environment", func(t *testing.T) {
		accessedEnv := false
		envReaderIO := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				accessedEnv = true
				return c.Value + 10
			}
		}

		resultReaderIO := MonadMapTo(envReaderIO, []int{1, 2, 3})
		config := ReaderTestConfig{Value: 20, Name: "test"}
		result := resultReaderIO(config)()

		assert.Equal(t, []int{1, 2, 3}, result)
		assert.True(t, accessedEnv, "ReaderIO should access environment during execution")
	})
}

func TestMonadChainFirst(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadChainFirst(rio, func(n int) ReaderIO[ReaderTestConfig, string] {
		sideEffect = n
		return Of[ReaderTestConfig]("side effect")
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirst(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		ChainFirst(func(n int) ReaderIO[ReaderTestConfig, string] {
			sideEffect = n
			return Of[ReaderTestConfig]("side effect")
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTap(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadTap(rio, func(n int) ReaderIO[ReaderTestConfig, func()] {
		sideEffect = n
		return Of[ReaderTestConfig](func() {})
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTap(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		Tap(func(n int) ReaderIO[ReaderTestConfig, func()] {
			sideEffect = n
			return Of[ReaderTestConfig](func() {})
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadChainFirstIOK(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadChainFirstIOK(rio, func(n int) G.IO[string] {
		sideEffect = n
		return G.Of("side effect")
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirstIOK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		ChainFirstIOK[ReaderTestConfig](func(n int) G.IO[string] {
			sideEffect = n
			return G.Of("side effect")
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTapIOK(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadTapIOK(rio, func(n int) G.IO[func()] {
		sideEffect = n
		return G.Of(func() {})
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTapIOK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		TapIOK[ReaderTestConfig](func(n int) G.IO[func()] {
			sideEffect = n
			return G.Of(func() {})
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadChainReaderK(t *testing.T) {
	rio := Of[ReaderTestConfig](5)
	result := MonadChainReaderK(rio, func(n int) func(ReaderTestConfig) int {
		return func(c ReaderTestConfig) int { return n + c.Value }
	})

	config := ReaderTestConfig{Value: 10, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestChainReaderK(t *testing.T) {
	result := F.Pipe1(
		Of[ReaderTestConfig](5),
		ChainReaderK(func(n int) func(ReaderTestConfig) int {
			return func(c ReaderTestConfig) int { return n + c.Value }
		}),
	)

	config := ReaderTestConfig{Value: 10, Name: "test"}
	assert.Equal(t, 15, result(config)())
}

func TestMonadChainFirstReaderK(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadChainFirstReaderK(rio, func(n int) func(ReaderTestConfig) string {
		return func(c ReaderTestConfig) string {
			sideEffect = n
			return "side effect"
		}
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirstReaderK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		ChainFirstReaderK(func(n int) func(ReaderTestConfig) string {
			return func(c ReaderTestConfig) string {
				sideEffect = n
				return "side effect"
			}
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTapReaderK(t *testing.T) {
	sideEffect := 0
	rio := Of[ReaderTestConfig](42)
	result := MonadTapReaderK(rio, func(n int) func(ReaderTestConfig) func() {
		return func(c ReaderTestConfig) func() {
			sideEffect = n
			return func() {}
		}
	})

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTapReaderK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of[ReaderTestConfig](42),
		TapReaderK(func(n int) func(ReaderTestConfig) func() {
			return func(c ReaderTestConfig) func() {
				sideEffect = n
				return func() {}
			}
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestRead(t *testing.T) {
	rio := Of[ReaderTestConfig](42)
	config := ReaderTestConfig{Value: 10, Name: "test"}
	ioAction := Read[int](config)(rio)
	result := ioAction()

	assert.Equal(t, 42, result)
}

func TestReadIO(t *testing.T) {
	t.Run("basic usage with IO environment", func(t *testing.T) {
		// Create a ReaderIO that uses the config
		rio := Of[ReaderTestConfig](42)

		// Create an IO that produces the config
		configIO := G.Of(ReaderTestConfig{Value: 21, Name: "test"})

		// Use ReadIO to execute the ReaderIO with the IO environment
		result := ReadIO[int](configIO)(rio)()

		assert.Equal(t, 42, result)
	})

	t.Run("chains IO effects correctly", func(t *testing.T) {
		// Track execution order
		executionOrder := []string{}

		// Create an IO that produces the config with a side effect
		configIO := func() ReaderTestConfig {
			executionOrder = append(executionOrder, "load config")
			return ReaderTestConfig{Value: 10, Name: "test"}
		}

		// Create a ReaderIO that uses the config with a side effect
		rio := func(c ReaderTestConfig) G.IO[int] {
			return func() int {
				executionOrder = append(executionOrder, "use config")
				return c.Value * 3
			}
		}

		// Execute the composed computation
		result := ReadIO[int](configIO)(rio)()

		assert.Equal(t, 30, result)
		assert.Equal(t, []string{"load config", "use config"}, executionOrder)
	})

	t.Run("works with complex environment loading", func(t *testing.T) {
		// Simulate loading config from a file or database
		loadConfigFromDB := func() ReaderTestConfig {
			// Simulate side effect
			return ReaderTestConfig{Value: 100, Name: "production"}
		}

		// A computation that depends on the loaded config
		getConnectionString := func(c ReaderTestConfig) G.IO[string] {
			return G.Of(c.Name + ":" + S.Format[int]("%d")(c.Value))
		}

		result := ReadIO[string](loadConfigFromDB)(getConnectionString)()

		assert.Equal(t, "production:100", result)
	})

	t.Run("composes with other ReaderIO operations", func(t *testing.T) {
		configIO := G.Of(ReaderTestConfig{Value: 5, Name: "test"})

		// Build a pipeline using ReaderIO operations
		pipeline := F.Pipe2(
			Ask[ReaderTestConfig](),
			Map[ReaderTestConfig](func(c ReaderTestConfig) int { return c.Value }),
			Chain(func(n int) ReaderIO[ReaderTestConfig, int] {
				return Of[ReaderTestConfig](n * 4)
			}),
		)

		result := ReadIO[int](configIO)(pipeline)()

		assert.Equal(t, 20, result)
	})

	t.Run("handles environment with multiple fields", func(t *testing.T) {
		configIO := G.Of(ReaderTestConfig{Value: 42, Name: "answer"})

		// Access both fields from the environment
		rio := func(c ReaderTestConfig) G.IO[string] {
			return G.Of(c.Name + "=" + S.Format[int]("%d")(c.Value))
		}

		result := ReadIO[string](configIO)(rio)()

		assert.Equal(t, "answer=42", result)
	})

	t.Run("lazy evaluation - IO not executed until called", func(t *testing.T) {
		executed := false

		configIO := func() ReaderTestConfig {
			executed = true
			return ReaderTestConfig{Value: 1, Name: "test"}
		}

		rio := Of[ReaderTestConfig](42)

		// Create the composed IO but don't execute it yet
		composedIO := ReadIO[int](configIO)(rio)

		// Config IO should not be executed yet
		assert.False(t, executed)

		// Now execute it
		result := composedIO()

		// Now it should be executed
		assert.True(t, executed)
		assert.Equal(t, 42, result)
	})

	t.Run("works with ChainIOK", func(t *testing.T) {
		configIO := G.Of(ReaderTestConfig{Value: 10, Name: "test"})

		pipeline := F.Pipe1(
			Of[ReaderTestConfig](5),
			ChainIOK[ReaderTestConfig](func(n int) G.IO[int] {
				return G.Of(n * 2)
			}),
		)

		result := ReadIO[int](configIO)(pipeline)()

		assert.Equal(t, 10, result)
	})

	t.Run("comparison with Read - different input types", func(t *testing.T) {
		rio := func(c ReaderTestConfig) G.IO[int] {
			return G.Of(c.Value + 10)
		}

		config := ReaderTestConfig{Value: 5, Name: "test"}

		// Using Read with a pure value
		resultRead := Read[int](config)(rio)()

		// Using ReadIO with an IO value
		resultReadIO := ReadIO[int](G.Of(config))(rio)()

		// Both should produce the same result
		assert.Equal(t, 15, resultRead)
		assert.Equal(t, 15, resultReadIO)
	})
}

func TestTapWithLogging(t *testing.T) {
	// Simulate logging scenario
	logged := []int{}

	result := F.Pipe3(
		Of[ReaderTestConfig](42),
		Tap(func(n int) ReaderIO[ReaderTestConfig, func()] {
			logged = append(logged, n)
			return Of[ReaderTestConfig](func() {})
		}),
		Map[ReaderTestConfig](N.Mul(2)),
		Tap(func(n int) ReaderIO[ReaderTestConfig, func()] {
			logged = append(logged, n)
			return Of[ReaderTestConfig](func() {})
		}),
	)

	config := ReaderTestConfig{Value: 1, Name: "test"}
	value := result(config)()
	assert.Equal(t, 84, value)
	assert.Equal(t, []int{42, 84}, logged)
}
