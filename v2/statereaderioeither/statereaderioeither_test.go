// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/stretchr/testify/assert"
)

type testState struct {
	counter int
}

type testContext struct {
	multiplier int
}

func TestOf(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}
	result := Of[testState, testContext, error](42)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Fold(
		func(err error) bool {
			t.Fatalf("Expected Right but got Left: %v", err)
			return false
		},
		func(p P.Pair[testState, int]) bool {
			assert.Equal(t, 42, P.Tail(p))
			assert.Equal(t, 0, P.Head(p).counter) // State unchanged
			return true
		},
	)(res)
}

func TestRight(t *testing.T) {
	state := testState{counter: 5}
	ctx := testContext{multiplier: 3}
	result := Right[testState, testContext, error](100)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 100, P.Tail(p))
		assert.Equal(t, 5, P.Head(p).counter)
		return p
	})(res)
}

func TestLeft(t *testing.T) {
	state := testState{counter: 10}
	ctx := testContext{multiplier: 2}
	testErr := errors.New("test error")
	result := Left[testState, testContext, int](testErr)
	res := result(state)(ctx)()

	assert.True(t, E.IsLeft(res))
}

func TestMonadMap(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	result := MonadMap(
		Of[testState, testContext, error](21),
		N.Mul(2),
	)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestMap(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	result := F.Pipe1(
		Of[testState, testContext, error](21),
		Map[testState, testContext, error](N.Mul(2)),
	)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestMonadChain(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	result := MonadChain(
		Of[testState, testContext, error](5),
		func(x int) StateReaderIOEither[testState, testContext, error, string] {
			return Of[testState, testContext, error](fmt.Sprintf("value: %d", x))
		},
	)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "value: 5", P.Tail(p))
		return p
	})(res)
}

func TestChain(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	result := F.Pipe1(
		Of[testState, testContext, error](5),
		Chain(func(x int) StateReaderIOEither[testState, testContext, error, string] {
			return Of[testState, testContext, error](fmt.Sprintf("value: %d", x))
		}),
	)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "value: 5", P.Tail(p))
		return p
	})(res)
}

func TestMonadAp(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	fab := Of[testState, testContext, error](N.Mul(2))
	fa := Of[testState, testContext, error](21)
	result := MonadAp(fab, fa)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestAp(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	fa := Of[testState, testContext, error](21)
	result := F.Pipe1(
		Of[testState, testContext, error](N.Mul(2)),
		Ap[int](fa),
	)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestFromReaderIOEither(t *testing.T) {
	state := testState{counter: 5}
	ctx := testContext{multiplier: 2}

	rioe := RIOE.Of[testContext, error](42)
	result := FromReaderIOEither[testState](rioe)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		assert.Equal(t, 5, P.Head(p).counter) // State unchanged
		return p
	})(res)
}

func TestFromReaderEither(t *testing.T) {
	state := testState{counter: 7}
	ctx := testContext{multiplier: 3}

	re := RE.Of[testContext, error](100)
	result := FromReaderEither[testState](re)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 100, P.Tail(p))
		assert.Equal(t, 7, P.Head(p).counter)
		return p
	})(res)
}

func TestFromIOEither(t *testing.T) {
	state := testState{counter: 3}
	ctx := testContext{multiplier: 4}

	ioe := IOE.Right[error](55)
	result := FromIOEither[testState, testContext](ioe)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 55, P.Tail(p))
		assert.Equal(t, 3, P.Head(p).counter)
		return p
	})(res)
}

func TestFromState(t *testing.T) {
	initialState := testState{counter: 10}
	ctx := testContext{multiplier: 2}

	// State computation that increments counter and returns it
	stateComp := func(s testState) P.Pair[testState, int] {
		newState := testState{counter: s.counter + 1}
		return P.MakePair(newState, newState.counter)
	}

	result := FromState[testContext, error](stateComp)
	res := result(initialState)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 11, P.Tail(p))         // Incremented value
		assert.Equal(t, 11, P.Head(p).counter) // State updated
		return p
	})(res)
}

func TestFromIO(t *testing.T) {
	state := testState{counter: 8}
	ctx := testContext{multiplier: 2}

	ioVal := func() int { return 99 }
	result := FromIO[testState, testContext, error](ioVal)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 99, P.Tail(p))
		assert.Equal(t, 8, P.Head(p).counter)
		return p
	})(res)
}

func TestFromReader(t *testing.T) {
	state := testState{counter: 6}
	ctx := testContext{multiplier: 5}

	reader := func(c testContext) int { return c.multiplier * 10 }
	result := FromReader[testState, error](reader)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 50, P.Tail(p))
		assert.Equal(t, 6, P.Head(p).counter)
		return p
	})(res)
}

func TestFromEither(t *testing.T) {
	state := testState{counter: 12}
	ctx := testContext{multiplier: 3}

	// Test Right case
	resultRight := FromEither[testState, testContext](E.Right[error](42))
	resRight := resultRight(state)(ctx)()
	assert.True(t, E.IsRight(resRight))

	// Test Left case
	resultLeft := FromEither[testState, testContext](E.Left[int](errors.New("error")))
	resLeft := resultLeft(state)(ctx)()
	assert.True(t, E.IsLeft(resLeft))
}

func TestLocal(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	// Create a computation that uses the context
	comp := Asks(func(c testContext) StateReaderIOEither[testState, testContext, error, int] {
		return Of[testState, testContext, error](c.multiplier * 10)
	})

	// Modify context before running computation
	result := Local[testState, error, int, int](
		func(c testContext) testContext {
			return testContext{multiplier: c.multiplier * 2}
		},
	)(comp)

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 40, P.Tail(p)) // (2 * 2) * 10
		return p
	})(res)
}

func TestAsks(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 7}

	result := Asks(func(c testContext) StateReaderIOEither[testState, testContext, error, int] {
		return Of[testState, testContext, error](c.multiplier * 5)
	})

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 35, P.Tail(p))
		return p
	})(res)
}

func TestFromEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	validate := func(x int) E.Either[error, int] {
		if x > 0 {
			return E.Right[error](x * 2)
		}
		return E.Left[int](errors.New("negative"))
	}

	kleisli := FromEitherK[testState, testContext](validate)

	// Test with valid input
	resultValid := kleisli(5)
	resValid := resultValid(state)(ctx)()
	assert.True(t, E.IsRight(resValid))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 10, P.Tail(p))
		return p
	})(resValid)

	// Test with invalid input
	resultInvalid := kleisli(-5)
	resInvalid := resultInvalid(state)(ctx)()
	assert.True(t, E.IsLeft(resInvalid))
}

func TestFromIOK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	ioFunc := func(x int) io.IO[int] {
		return func() int { return x * 3 }
	}

	kleisli := FromIOK[testState, testContext, error](ioFunc)
	result := kleisli(7)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 21, P.Tail(p))
		return p
	})(res)
}

func TestFromIOEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	ioeFunc := func(x int) IOE.IOEither[error, int] {
		if x > 0 {
			return IOE.Right[error](x * 4)
		}
		return IOE.Left[int](errors.New("invalid"))
	}

	kleisli := FromIOEitherK[testState, testContext](ioeFunc)
	result := kleisli(3)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 12, P.Tail(p))
		return p
	})(res)
}

func TestFromReaderIOEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	rioeFunc := func(x int) RIOE.ReaderIOEither[testContext, error, int] {
		return func(c testContext) IOE.IOEither[error, int] {
			return IOE.Right[error](x * c.multiplier)
		}
	}

	kleisli := FromReaderIOEitherK[testState](rioeFunc)
	result := kleisli(5)
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 10, P.Tail(p))
		return p
	})(res)
}

func TestChainEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	validate := func(x int) E.Either[error, string] {
		if x > 0 {
			return E.Right[error](fmt.Sprintf("valid: %d", x))
		}
		return E.Left[string](errors.New("invalid"))
	}

	result := F.Pipe1(
		Of[testState, testContext, error](42),
		ChainEitherK[testState, testContext](validate),
	)

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "valid: 42", P.Tail(p))
		return p
	})(res)
}

func TestChainIOEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	ioeFunc := func(x int) IOE.IOEither[error, string] {
		return IOE.Right[error](fmt.Sprintf("result: %d", x))
	}

	result := F.Pipe1(
		Of[testState, testContext, error](100),
		ChainIOEitherK[testState, testContext](ioeFunc),
	)

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "result: 100", P.Tail(p))
		return p
	})(res)
}

func TestChainReaderIOEitherK(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 3}

	rioeFunc := func(x int) RIOE.ReaderIOEither[testContext, error, int] {
		return func(c testContext) IOE.IOEither[error, int] {
			return IOE.Right[error](x * c.multiplier)
		}
	}

	result := F.Pipe1(
		Of[testState, testContext, error](5),
		ChainReaderIOEitherK[testState](rioeFunc),
	)

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 15, P.Tail(p))
		return p
	})(res)
}

func TestDo(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	type Result struct {
		value int
	}

	result := Do[testState, testContext, error](Result{})
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, Result]) P.Pair[testState, Result] {
		assert.Equal(t, 0, P.Tail(p).value)
		return p
	})(res)
}

func TestBindTo(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	type Result struct {
		value int
	}

	result := F.Pipe1(
		Of[testState, testContext, error](42),
		BindTo[testState, testContext, error](func(v int) Result {
			return Result{value: v}
		}),
	)

	res := result(state)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, Result]) P.Pair[testState, Result] {
		assert.Equal(t, 42, P.Tail(p).value)
		return p
	})(res)
}

func TestStatefulComputation(t *testing.T) {
	initialState := testState{counter: 0}
	ctx := testContext{multiplier: 10}

	// Create a computation that modifies state
	incrementAndGet := func(s testState) P.Pair[testState, int] {
		newState := testState{counter: s.counter + 1}
		return P.MakePair(newState, newState.counter)
	}

	// Chain multiple stateful operations
	result := F.Pipe2(
		FromState[testContext, error](incrementAndGet),
		Chain(func(v1 int) StateReaderIOEither[testState, testContext, error, int] {
			return FromState[testContext, error](incrementAndGet)
		}),
		Chain(func(v2 int) StateReaderIOEither[testState, testContext, error, int] {
			return FromState[testContext, error](incrementAndGet)
		}),
	)

	res := result(initialState)(ctx)()
	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 3, P.Tail(p))         // Last incremented value
		assert.Equal(t, 3, P.Head(p).counter) // State updated three times
		return p
	})(res)
}

func TestErrorPropagation(t *testing.T) {
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}

	testErr := errors.New("test error")

	// Chain operations where the second one fails
	result := F.Pipe1(
		Of[testState, testContext, error](42),
		Chain(func(x int) StateReaderIOEither[testState, testContext, error, int] {
			return Left[testState, testContext, int](testErr)
		}),
	)

	res := result(state)(ctx)()
	assert.True(t, E.IsLeft(res))
}

func TestPointed(t *testing.T) {
	p := Pointed[testState, testContext, error, int]()
	assert.NotNil(t, p)

	result := p.Of(42)
	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
}

func TestFunctor(t *testing.T) {
	f := Functor[testState, testContext, error, int, string]()
	assert.NotNil(t, f)

	mapper := f.Map(func(x int) string { return fmt.Sprintf("%d", x) })
	result := mapper(Of[testState, testContext, error](42))

	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}

func TestApplicative(t *testing.T) {
	a := Applicative[testState, testContext, error, int, string]()
	assert.NotNil(t, a)

	fab := Of[testState, testContext, error](func(x int) string { return fmt.Sprintf("%d", x) })
	fa := Of[testState, testContext, error](42)
	result := a.Ap(fa)(fab)

	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}

func TestMonad(t *testing.T) {
	m := Monad[testState, testContext, error, int, string]()
	assert.NotNil(t, m)

	fa := m.Of(42)
	result := m.Chain(func(x int) StateReaderIOEither[testState, testContext, error, string] {
		return Of[testState, testContext, error](fmt.Sprintf("%d", x))
	})(fa)

	state := testState{counter: 0}
	ctx := testContext{multiplier: 2}
	res := result(state)(ctx)()

	assert.True(t, E.IsRight(res))
	E.Map[error](func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}
