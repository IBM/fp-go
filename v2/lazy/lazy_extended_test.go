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

package lazy

import (
	"testing"
	"time"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	result := Of(42)
	assert.Equal(t, 42, result())
}

func TestFromLazy(t *testing.T) {
	original := func() int { return 42 }
	wrapped := FromLazy(original)
	assert.Equal(t, 42, wrapped())
}

func TestFromImpure(t *testing.T) {
	counter := 0
	impure := func() {
		counter++
	}
	lazy := FromImpure(impure)
	lazy()
	assert.Equal(t, 1, counter)
}

func TestMonadOf(t *testing.T) {
	result := MonadOf(42)
	assert.Equal(t, 42, result())
}

func TestMonadMap(t *testing.T) {
	result := MonadMap(Of(5), N.Mul(2))
	assert.Equal(t, 10, result())
}

func TestMonadMapTo(t *testing.T) {
	result := MonadMapTo(Of("ignored"), 42)
	assert.Equal(t, 42, result())
}

func TestMapTo(t *testing.T) {
	mapper := MapTo[string](42)
	result := mapper(Of("ignored"))
	assert.Equal(t, 42, result())
}

func TestMonadChain(t *testing.T) {
	result := MonadChain(Of(5), func(x int) Lazy[int] {
		return Of(x * 2)
	})
	assert.Equal(t, 10, result())
}

func TestMonadChainFirst(t *testing.T) {
	result := MonadChainFirst(Of(5), func(x int) Lazy[string] {
		return Of("ignored")
	})
	assert.Equal(t, 5, result())
}

func TestChainFirst(t *testing.T) {
	chainer := ChainFirst(func(x int) Lazy[string] {
		return Of("ignored")
	})
	result := chainer(Of(5))
	assert.Equal(t, 5, result())
}

func TestMonadChainTo(t *testing.T) {
	result := MonadChainTo(Of(5), Of(10))
	assert.Equal(t, 10, result())
}

func TestChainTo(t *testing.T) {
	chainer := ChainTo[int](Of(10))
	result := chainer(Of(5))
	assert.Equal(t, 10, result())
}

func TestMonadAp(t *testing.T) {
	lazyFunc := Of(N.Mul(2))
	lazyValue := Of(5)
	result := MonadAp(lazyFunc, lazyValue)
	assert.Equal(t, 10, result())
}

func TestMonadApFirst(t *testing.T) {
	result := MonadApFirst(Of(5), Of(10))
	assert.Equal(t, 5, result())
}

func TestMonadApSecond(t *testing.T) {
	result := MonadApSecond(Of(5), Of(10))
	assert.Equal(t, 10, result())
}

func TestNow(t *testing.T) {
	before := time.Now()
	result := Now()
	after := time.Now()

	assert.True(t, result.After(before) || result.Equal(before))
	assert.True(t, result.Before(after) || result.Equal(after))
}

func TestDefer(t *testing.T) {
	counter := 0
	deferred := Defer(func() Lazy[int] {
		counter++
		return Of(counter)
	})

	// First execution
	result1 := deferred()
	assert.Equal(t, 1, result1)

	// Second execution should generate a new computation
	result2 := deferred()
	assert.Equal(t, 2, result2)
}

func TestDo(t *testing.T) {
	type State struct {
		Value int
	}
	result := Do(State{Value: 42})
	assert.Equal(t, State{Value: 42}, result())
}

func TestLet(t *testing.T) {
	type State struct {
		Value int
	}

	result := F.Pipe2(
		Do(State{}),
		Let(
			func(v int) func(State) State {
				return func(s State) State { s.Value = v; return s }
			},
			func(s State) int { return 42 },
		),
		Map(func(s State) int { return s.Value }),
	)

	assert.Equal(t, 42, result())
}

func TestLetTo(t *testing.T) {
	type State struct {
		Value int
	}

	result := F.Pipe2(
		Do(State{}),
		LetTo(
			func(v int) func(State) State {
				return func(s State) State { s.Value = v; return s }
			},
			42,
		),
		Map(func(s State) int { return s.Value }),
	)

	assert.Equal(t, 42, result())
}

func TestBindTo(t *testing.T) {
	type State struct {
		Value int
	}

	result := F.Pipe2(
		Of(42),
		BindTo(func(v int) State { return State{Value: v} }),
		Map(func(s State) int { return s.Value }),
	)

	assert.Equal(t, 42, result())
}

func TestBindL(t *testing.T) {
	type Config struct {
		Port int
	}
	type State struct {
		Config Config
	}

	// Create a lens manually
	configLens := L.MakeLens(
		func(s State) Config { return s.Config },
		func(s State, cfg Config) State { s.Config = cfg; return s },
	)

	result := F.Pipe2(
		Do(State{Config: Config{Port: 8080}}),
		BindL(configLens, func(cfg Config) Lazy[Config] {
			return Of(Config{Port: cfg.Port + 1})
		}),
		Map(func(s State) int { return s.Config.Port }),
	)

	assert.Equal(t, 8081, result())
}

func TestLetL(t *testing.T) {
	type Config struct {
		Port int
	}
	type State struct {
		Config Config
	}

	// Create a lens manually
	configLens := L.MakeLens(
		func(s State) Config { return s.Config },
		func(s State, cfg Config) State { s.Config = cfg; return s },
	)

	result := F.Pipe2(
		Do(State{Config: Config{Port: 8080}}),
		LetL(configLens, func(cfg Config) Config {
			return Config{Port: cfg.Port + 1}
		}),
		Map(func(s State) int { return s.Config.Port }),
	)

	assert.Equal(t, 8081, result())
}

func TestLetToL(t *testing.T) {
	type Config struct {
		Port int
	}
	type State struct {
		Config Config
	}

	// Create a lens manually
	configLens := L.MakeLens(
		func(s State) Config { return s.Config },
		func(s State, cfg Config) State { s.Config = cfg; return s },
	)

	result := F.Pipe2(
		Do(State{}),
		LetToL(configLens, Config{Port: 8080}),
		Map(func(s State) int { return s.Config.Port }),
	)

	assert.Equal(t, 8080, result())
}

func TestApSL(t *testing.T) {
	type Config struct {
		Port int
	}
	type State struct {
		Config Config
	}

	// Create a lens manually
	configLens := L.MakeLens(
		func(s State) Config { return s.Config },
		func(s State, cfg Config) State { s.Config = cfg; return s },
	)

	result := F.Pipe2(
		Do(State{}),
		ApSL(configLens, Of(Config{Port: 8080})),
		Map(func(s State) int { return s.Config.Port }),
	)

	assert.Equal(t, 8080, result())
}

func TestSequenceT1(t *testing.T) {
	result := SequenceT1(Of(42))
	tuple := result()
	assert.Equal(t, 42, tuple.F1)
}

func TestSequenceT2(t *testing.T) {
	result := SequenceT2(Of(42), Of("hello"))
	tuple := result()
	assert.Equal(t, 42, tuple.F1)
	assert.Equal(t, "hello", tuple.F2)
}

func TestSequenceT3(t *testing.T) {
	result := SequenceT3(Of(42), Of("hello"), Of(true))
	tuple := result()
	assert.Equal(t, 42, tuple.F1)
	assert.Equal(t, "hello", tuple.F2)
	assert.Equal(t, true, tuple.F3)
}

func TestSequenceT4(t *testing.T) {
	result := SequenceT4(Of(42), Of("hello"), Of(true), Of(3.14))
	tuple := result()
	assert.Equal(t, 42, tuple.F1)
	assert.Equal(t, "hello", tuple.F2)
	assert.Equal(t, true, tuple.F3)
	assert.Equal(t, 3.14, tuple.F4)
}

func TestTraverseArray(t *testing.T) {
	numbers := []int{1, 2, 3}
	result := F.Pipe1(
		numbers,
		TraverseArray(func(x int) Lazy[int] {
			return Of(x * 2)
		}),
	)
	assert.Equal(t, []int{2, 4, 6}, result())
}

func TestTraverseArrayWithIndex(t *testing.T) {
	numbers := []int{10, 20, 30}
	result := F.Pipe1(
		numbers,
		TraverseArrayWithIndex(func(i int, x int) Lazy[int] {
			return Of(x + i)
		}),
	)
	assert.Equal(t, []int{10, 21, 32}, result())
}

func TestSequenceArray(t *testing.T) {
	lazies := []Lazy[int]{Of(1), Of(2), Of(3)}
	result := SequenceArray(lazies)
	assert.Equal(t, []int{1, 2, 3}, result())
}

func TestMonadTraverseArray(t *testing.T) {
	numbers := []int{1, 2, 3}
	result := MonadTraverseArray(numbers, func(x int) Lazy[int] {
		return Of(x * 2)
	})
	assert.Equal(t, []int{2, 4, 6}, result())
}

func TestTraverseRecord(t *testing.T) {
	record := map[string]int{"a": 1, "b": 2}
	result := F.Pipe1(
		record,
		TraverseRecord[string](func(x int) Lazy[int] {
			return Of(x * 2)
		}),
	)
	resultMap := result()
	assert.Equal(t, 2, resultMap["a"])
	assert.Equal(t, 4, resultMap["b"])
}

func TestTraverseRecordWithIndex(t *testing.T) {
	record := map[string]int{"a": 10, "b": 20}
	result := F.Pipe1(
		record,
		TraverseRecordWithIndex(func(k string, x int) Lazy[int] {
			if k == "a" {
				return Of(x + 1)
			}
			return Of(x + 2)
		}),
	)
	resultMap := result()
	assert.Equal(t, 11, resultMap["a"])
	assert.Equal(t, 22, resultMap["b"])
}

func TestSequenceRecord(t *testing.T) {
	record := map[string]Lazy[int]{
		"a": Of(1),
		"b": Of(2),
	}
	result := SequenceRecord(record)
	resultMap := result()
	assert.Equal(t, 1, resultMap["a"])
	assert.Equal(t, 2, resultMap["b"])
}

func TestMonadTraverseRecord(t *testing.T) {
	record := map[string]int{"a": 1, "b": 2}
	result := MonadTraverseRecord(record, func(x int) Lazy[int] {
		return Of(x * 2)
	})
	resultMap := result()
	assert.Equal(t, 2, resultMap["a"])
	assert.Equal(t, 4, resultMap["b"])
}

func TestApplySemigroup(t *testing.T) {
	sg := ApplySemigroup(M.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	))

	result := sg.Concat(Of(5), Of(10))
	assert.Equal(t, 15, result())
}

func TestApplicativeMonoid(t *testing.T) {
	mon := ApplicativeMonoid(M.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	))

	// Test Empty
	empty := mon.Empty()
	assert.Equal(t, 0, empty())

	// Test Concat
	result := mon.Concat(Of(5), Of(10))
	assert.Equal(t, 15, result())

	// Test identity laws
	left := mon.Concat(mon.Empty(), Of(5))
	assert.Equal(t, 5, left())

	right := mon.Concat(Of(5), mon.Empty())
	assert.Equal(t, 5, right())
}

func TestEq(t *testing.T) {
	eq := Eq(EQ.FromEquals(func(a, b int) bool { return a == b }))

	assert.True(t, eq.Equals(Of(42), Of(42)))
	assert.False(t, eq.Equals(Of(42), Of(43)))
}

func TestComplexDoNotation(t *testing.T) {
	// Test a more complex do-notation scenario
	result := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, func(s utils.Initial) Lazy[string] {
			return Of("Doe")
		}),
		Bind(utils.SetGivenName, func(s utils.WithLastName) Lazy[string] {
			return Of("John")
		}),
		Map(utils.GetFullName),
	)

	assert.Equal(t, "John Doe", result())
}

func TestChainComposition(t *testing.T) {
	// Test chaining multiple operations
	double := func(x int) Lazy[int] {
		return Of(x * 2)
	}

	addTen := func(x int) Lazy[int] {
		return Of(x + 10)
	}

	result := F.Pipe2(
		Of(5),
		Chain(double),
		Chain(addTen),
	)

	assert.Equal(t, 20, result())
}

func TestMapComposition(t *testing.T) {
	// Test mapping multiple transformations
	result := F.Pipe3(
		Of(5),
		Map(N.Mul(2)),
		Map(N.Add(10)),
		Map(reader.Ask[int]()),
	)

	assert.Equal(t, 20, result())
}
