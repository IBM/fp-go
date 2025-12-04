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
		ChainIOK[ReaderTestConfig, int, int](func(n int) G.IO[int] {
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
		Flap[ReaderTestConfig, int, int](7),
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
		ChainFirstIOK[ReaderTestConfig, int, string](func(n int) G.IO[string] {
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
		TapIOK[ReaderTestConfig, int, func()](func(n int) G.IO[func()] {
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
