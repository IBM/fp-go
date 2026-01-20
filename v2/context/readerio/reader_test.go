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
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

func TestMonadMap(t *testing.T) {
	rio := Of(5)
	doubled := MonadMap(rio, N.Mul(2))

	result := doubled(t.Context())()
	assert.Equal(t, 10, result)
}

func TestMap(t *testing.T) {
	g := F.Pipe1(
		Of(1),
		Map(utils.Double),
	)

	assert.Equal(t, 2, g(t.Context())())
}

func TestMonadMapTo(t *testing.T) {
	rio := Of(42)
	replaced := MonadMapTo(rio, "constant")

	result := replaced(t.Context())()
	assert.Equal(t, "constant", result)
}

func TestMapTo(t *testing.T) {
	result := F.Pipe1(
		Of(42),
		MapTo[int]("constant"),
	)

	assert.Equal(t, "constant", result(t.Context())())
}

func TestMonadChain(t *testing.T) {
	rio1 := Of(5)
	result := MonadChain(rio1, func(n int) ReaderIO[int] {
		return Of(n * 3)
	})

	assert.Equal(t, 15, result(t.Context())())
}

func TestChain(t *testing.T) {
	result := F.Pipe1(
		Of(5),
		Chain(func(n int) ReaderIO[int] {
			return Of(n * 3)
		}),
	)

	assert.Equal(t, 15, result(t.Context())())
}

func TestMonadChainFirst(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadChainFirst(rio, func(n int) ReaderIO[string] {
		sideEffect = n
		return Of("side effect")
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirst(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		ChainFirst(func(n int) ReaderIO[string] {
			sideEffect = n
			return Of("side effect")
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTap(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadTap(rio, func(n int) ReaderIO[func()] {
		sideEffect = n
		return Of(func() {})
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTap(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		Tap(func(n int) ReaderIO[func()] {
			sideEffect = n
			return Of(func() {})
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestOf(t *testing.T) {
	rio := Of(100)
	result := rio(t.Context())()

	assert.Equal(t, 100, result)
}

func TestMonadAp(t *testing.T) {
	fabIO := Of(N.Mul(2))
	faIO := Of(5)
	result := MonadAp(fabIO, faIO)

	assert.Equal(t, 10, result(t.Context())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of(utils.Double),
		Ap[int](Of(1)),
	)

	assert.Equal(t, 2, g(t.Context())())
}

func TestMonadApSeq(t *testing.T) {
	fabIO := Of(N.Add(10))
	faIO := Of(5)
	result := MonadApSeq(fabIO, faIO)

	assert.Equal(t, 15, result(t.Context())())
}

func TestApSeq(t *testing.T) {
	g := F.Pipe1(
		Of(N.Add(10)),
		ApSeq[int](Of(5)),
	)

	assert.Equal(t, 15, g(t.Context())())
}

func TestMonadApPar(t *testing.T) {
	fabIO := Of(N.Add(10))
	faIO := Of(5)
	result := MonadApPar(fabIO, faIO)

	assert.Equal(t, 15, result(t.Context())())
}

func TestApPar(t *testing.T) {
	g := F.Pipe1(
		Of(N.Add(10)),
		ApPar[int](Of(5)),
	)

	assert.Equal(t, 15, g(t.Context())())
}

func TestAsk(t *testing.T) {
	rio := Ask()
	ctx := context.WithValue(t.Context(), "key", "value")
	result := rio(ctx)()

	assert.Equal(t, ctx, result)
}

func TestFromIO(t *testing.T) {
	ioAction := G.Of(42)
	rio := FromIO(ioAction)

	result := rio(t.Context())()
	assert.Equal(t, 42, result)
}

func TestFromReader(t *testing.T) {
	rdr := func(ctx context.Context) int {
		return 42
	}

	rio := FromReader(rdr)
	result := rio(t.Context())()

	assert.Equal(t, 42, result)
}

func TestFromLazy(t *testing.T) {
	lazy := func() int { return 42 }
	rio := FromLazy(lazy)

	result := rio(t.Context())()
	assert.Equal(t, 42, result)
}

func TestMonadChainIOK(t *testing.T) {
	rio := Of(5)
	result := MonadChainIOK(rio, func(n int) G.IO[int] {
		return G.Of(n * 4)
	})

	assert.Equal(t, 20, result(t.Context())())
}

func TestChainIOK(t *testing.T) {
	result := F.Pipe1(
		Of(5),
		ChainIOK(func(n int) G.IO[int] {
			return G.Of(n * 4)
		}),
	)

	assert.Equal(t, 20, result(t.Context())())
}

func TestMonadChainFirstIOK(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadChainFirstIOK(rio, func(n int) G.IO[string] {
		sideEffect = n
		return G.Of("side effect")
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirstIOK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		ChainFirstIOK(func(n int) G.IO[string] {
			sideEffect = n
			return G.Of("side effect")
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTapIOK(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadTapIOK(rio, func(n int) G.IO[func()] {
		sideEffect = n
		return G.Of(func() {})
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTapIOK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		TapIOK(func(n int) G.IO[func()] {
			sideEffect = n
			return G.Of(func() {})
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestDefer(t *testing.T) {
	counter := 0
	rio := Defer(func() ReaderIO[int] {
		counter++
		return Of(counter)
	})

	result1 := rio(t.Context())()
	result2 := rio(t.Context())()

	assert.Equal(t, 1, result1)
	assert.Equal(t, 2, result2)
}

func TestMemoize(t *testing.T) {
	counter := 0
	rio := Of(0)
	memoized := Memoize(MonadMap(rio, func(int) int {
		counter++
		return counter
	}))

	result1 := memoized(t.Context())()
	result2 := memoized(t.Context())()

	assert.Equal(t, 1, result1)
	assert.Equal(t, 1, result2) // Same value, memoized
}

func TestFlatten(t *testing.T) {
	nested := Of(Of(42))
	flattened := Flatten(nested)

	result := flattened(t.Context())()
	assert.Equal(t, 42, result)
}

func TestMonadFlap(t *testing.T) {
	fabIO := Of(N.Mul(3))
	result := MonadFlap(fabIO, 7)

	assert.Equal(t, 21, result(t.Context())())
}

func TestFlap(t *testing.T) {
	result := F.Pipe1(
		Of(N.Mul(3)),
		Flap[int](7),
	)

	assert.Equal(t, 21, result(t.Context())())
}

func TestMonadChainReaderK(t *testing.T) {
	rio := Of(5)
	result := MonadChainReaderK(rio, func(n int) reader.Reader[context.Context, int] {
		return func(ctx context.Context) int { return n * 2 }
	})

	assert.Equal(t, 10, result(t.Context())())
}

func TestChainReaderK(t *testing.T) {
	result := F.Pipe1(
		Of(5),
		ChainReaderK(func(n int) reader.Reader[context.Context, int] {
			return func(ctx context.Context) int { return n * 2 }
		}),
	)

	assert.Equal(t, 10, result(t.Context())())
}

func TestMonadChainFirstReaderK(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadChainFirstReaderK(rio, func(n int) reader.Reader[context.Context, string] {
		return func(ctx context.Context) string {
			sideEffect = n
			return "side effect"
		}
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestChainFirstReaderK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		ChainFirstReaderK(func(n int) reader.Reader[context.Context, string] {
			return func(ctx context.Context) string {
				sideEffect = n
				return "side effect"
			}
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestMonadTapReaderK(t *testing.T) {
	sideEffect := 0
	rio := Of(42)
	result := MonadTapReaderK(rio, func(n int) reader.Reader[context.Context, func()] {
		return func(ctx context.Context) func() {
			sideEffect = n
			return func() {}
		}
	})

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestTapReaderK(t *testing.T) {
	sideEffect := 0
	result := F.Pipe1(
		Of(42),
		TapReaderK(func(n int) reader.Reader[context.Context, func()] {
			return func(ctx context.Context) func() {
				sideEffect = n
				return func() {}
			}
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 42, value)
	assert.Equal(t, 42, sideEffect)
}

func TestRead(t *testing.T) {
	rio := Of(42)
	ctx := t.Context()
	ioAction := Read[int](ctx)(rio)
	result := ioAction()

	assert.Equal(t, 42, result)
}

func TestComplexPipeline(t *testing.T) {
	// Test a complex pipeline combining multiple operations
	result := F.Pipe3(
		Ask(),
		Map(func(ctx context.Context) int { return 5 }),
		Chain(func(n int) ReaderIO[int] {
			return Of(n * 2)
		}),
		Map(N.Add(10)),
	)

	assert.Equal(t, 20, result(t.Context())()) // (5 * 2) + 10 = 20
}

func TestFromIOWithChain(t *testing.T) {
	ioAction := G.Of(10)

	result := F.Pipe1(
		FromIO(ioAction),
		Chain(func(n int) ReaderIO[int] {
			return Of(n + 5)
		}),
	)

	assert.Equal(t, 15, result(t.Context())())
}

func TestTapWithLogging(t *testing.T) {
	// Simulate logging scenario
	logged := []int{}

	result := F.Pipe3(
		Of(42),
		Tap(func(n int) ReaderIO[func()] {
			logged = append(logged, n)
			return Of(func() {})
		}),
		Map(N.Mul(2)),
		Tap(func(n int) ReaderIO[func()] {
			logged = append(logged, n)
			return Of(func() {})
		}),
	)

	value := result(t.Context())()
	assert.Equal(t, 84, value)
	assert.Equal(t, []int{42, 84}, logged)
}

func TestReadIO(t *testing.T) {
	// Test basic ReadIO functionality
	contextIO := G.Of(context.WithValue(t.Context(), "testKey", "testValue"))
	rio := FromReader(func(ctx context.Context) string {
		if val := ctx.Value("testKey"); val != nil {
			return val.(string)
		}
		return "default"
	})

	ioAction := ReadIO[string](contextIO)(rio)
	result := ioAction()

	assert.Equal(t, "testValue", result)
}

func TestReadIOWithBackground(t *testing.T) {
	// Test ReadIO with plain background context
	contextIO := G.Of(t.Context())
	rio := Of(42)

	ioAction := ReadIO[int](contextIO)(rio)
	result := ioAction()

	assert.Equal(t, 42, result)
}

func TestReadIOWithChain(t *testing.T) {
	// Test ReadIO with chained operations
	contextIO := G.Of(context.WithValue(t.Context(), "multiplier", 3))

	result := F.Pipe1(
		FromReader(func(ctx context.Context) int {
			if val := ctx.Value("multiplier"); val != nil {
				return val.(int)
			}
			return 1
		}),
		Chain(func(n int) ReaderIO[int] {
			return Of(n * 10)
		}),
	)

	ioAction := ReadIO[int](contextIO)(result)
	value := ioAction()

	assert.Equal(t, 30, value) // 3 * 10
}

func TestReadIOWithMap(t *testing.T) {
	// Test ReadIO with Map operations
	contextIO := G.Of(t.Context())

	result := F.Pipe2(
		Of(5),
		Map(N.Mul(2)),
		Map(N.Add(10)),
	)

	ioAction := ReadIO[int](contextIO)(result)
	value := ioAction()

	assert.Equal(t, 20, value) // (5 * 2) + 10
}

func TestReadIOWithSideEffects(t *testing.T) {
	// Test ReadIO with side effects in context creation
	counter := 0
	contextIO := func() context.Context {
		counter++
		return context.WithValue(t.Context(), "counter", counter)
	}

	rio := FromReader(func(ctx context.Context) int {
		if val := ctx.Value("counter"); val != nil {
			return val.(int)
		}
		return 0
	})

	ioAction := ReadIO[int](contextIO)(rio)
	result := ioAction()

	assert.Equal(t, 1, result)
	assert.Equal(t, 1, counter)
}

func TestReadIOMultipleExecutions(t *testing.T) {
	// Test that ReadIO creates fresh effects on each execution
	counter := 0
	contextIO := func() context.Context {
		counter++
		return t.Context()
	}

	rio := Of(42)
	ioAction := ReadIO[int](contextIO)(rio)

	result1 := ioAction()
	result2 := ioAction()

	assert.Equal(t, 42, result1)
	assert.Equal(t, 42, result2)
	assert.Equal(t, 2, counter) // Context IO executed twice
}

func TestReadIOComparisonWithRead(t *testing.T) {
	// Compare ReadIO with Read to show the difference
	ctx := context.WithValue(t.Context(), "key", "value")

	rio := FromReader(func(ctx context.Context) string {
		if val := ctx.Value("key"); val != nil {
			return val.(string)
		}
		return "default"
	})

	// Using Read (direct context)
	ioAction1 := Read[string](ctx)(rio)
	result1 := ioAction1()

	// Using ReadIO (context wrapped in IO)
	contextIO := G.Of(ctx)
	ioAction2 := ReadIO[string](contextIO)(rio)
	result2 := ioAction2()

	assert.Equal(t, result1, result2)
	assert.Equal(t, "value", result1)
	assert.Equal(t, "value", result2)
}

func TestReadIOWithComplexContext(t *testing.T) {
	// Test ReadIO with complex context manipulation
	type contextKey string
	const (
		userKey  contextKey = "user"
		tokenKey contextKey = "token"
	)

	contextIO := G.Of(
		context.WithValue(
			context.WithValue(t.Context(), userKey, "Alice"),
			tokenKey,
			"secret123",
		),
	)

	rio := FromReader(func(ctx context.Context) map[string]string {
		result := make(map[string]string)
		if user := ctx.Value(userKey); user != nil {
			result["user"] = user.(string)
		}
		if token := ctx.Value(tokenKey); token != nil {
			result["token"] = token.(string)
		}
		return result
	})

	ioAction := ReadIO[map[string]string](contextIO)(rio)
	result := ioAction()

	assert.Equal(t, "Alice", result["user"])
	assert.Equal(t, "secret123", result["token"])
}

func TestReadIOWithAsk(t *testing.T) {
	// Test ReadIO combined with Ask
	contextIO := G.Of(context.WithValue(t.Context(), "data", 100))

	result := F.Pipe1(
		Ask(),
		Map(func(ctx context.Context) int {
			if val := ctx.Value("data"); val != nil {
				return val.(int)
			}
			return 0
		}),
	)

	ioAction := ReadIO[int](contextIO)(result)
	value := ioAction()

	assert.Equal(t, 100, value)
}
