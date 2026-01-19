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

package statereaderioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOR "github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	RES "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type testState struct {
	counter int
}

func TestOf(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()
	result := Of[testState](42)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Fold(
		func(err error) bool {
			t.Fatalf("Expected Success but got Error: %v", err)
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
	ctx := t.Context()
	result := Right[testState](100)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 100, P.Tail(p))
		assert.Equal(t, 5, P.Head(p).counter)
		return p
	})(res)
}

func TestLeft(t *testing.T) {
	state := testState{counter: 10}
	ctx := t.Context()
	testErr := errors.New("test error")
	result := Left[testState, int](testErr)
	res := result(state)(ctx)()

	assert.True(t, RES.IsLeft(res))
}

func TestMonadMap(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	result := MonadMap(
		Of[testState](21),
		N.Mul(2),
	)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestMap(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	result := F.Pipe1(
		Of[testState](21),
		Map[testState](N.Mul(2)),
	)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestMonadChain(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	result := MonadChain(
		Of[testState](5),
		func(x int) StateReaderIOResult[testState, string] {
			return Of[testState](fmt.Sprintf("value: %d", x))
		},
	)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "value: 5", P.Tail(p))
		return p
	})(res)
}

func TestChain(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	result := F.Pipe1(
		Of[testState](5),
		Chain(func(x int) StateReaderIOResult[testState, string] {
			return Of[testState](fmt.Sprintf("value: %d", x))
		}),
	)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "value: 5", P.Tail(p))
		return p
	})(res)
}

func TestMonadAp(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	fab := Of[testState](N.Mul(2))
	fa := Of[testState](21)
	result := MonadAp(fab, fa)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestAp(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	fa := Of[testState](21)
	result := F.Pipe1(
		Of[testState](N.Mul(2)),
		Ap[int](fa),
	)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(res)
}

func TestFromIOResult(t *testing.T) {
	state := testState{counter: 3}
	ctx := t.Context()

	ior := IOR.Of(55)
	result := FromIOResult[testState](ior)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 55, P.Tail(p))
		assert.Equal(t, 3, P.Head(p).counter)
		return p
	})(res)
}

func TestFromState(t *testing.T) {
	initialState := testState{counter: 10}
	ctx := t.Context()

	// State computation that increments counter and returns it
	stateComp := func(s testState) P.Pair[testState, int] {
		newState := testState{counter: s.counter + 1}
		return P.MakePair(newState, newState.counter)
	}

	result := FromState(stateComp)
	res := result(initialState)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 11, P.Tail(p))         // Incremented value
		assert.Equal(t, 11, P.Head(p).counter) // State updated
		return p
	})(res)
}

func TestFromIO(t *testing.T) {
	state := testState{counter: 8}
	ctx := t.Context()

	ioVal := func() int { return 99 }
	result := FromIO[testState](ioVal)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 99, P.Tail(p))
		assert.Equal(t, 8, P.Head(p).counter)
		return p
	})(res)
}

func TestFromResult(t *testing.T) {
	state := testState{counter: 12}
	ctx := t.Context()

	// Test Success case
	resultSuccess := FromResult[testState](RES.Of(42))
	resSuccess := resultSuccess(state)(ctx)()
	assert.True(t, RES.IsRight(resSuccess))

	// Test Error case
	resultError := FromResult[testState](RES.Left[int](errors.New("error")))
	resError := resultError(state)(ctx)()
	assert.True(t, RES.IsLeft(resError))
}

func TestLocal(t *testing.T) {
	state := testState{counter: 0}
	ctx := context.WithValue(t.Context(), "key", "value1")

	// Create a computation that uses the context
	comp := Asks(func(c context.Context) StateReaderIOResult[testState, string] {
		val := c.Value("key").(string)
		return Of[testState](val)
	})

	// Modify context before running computation
	result := Local[testState, string](
		func(c context.Context) context.Context {
			return context.WithValue(c, "key", "value2")
		},
	)(comp)

	res := result(state)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "value2", P.Tail(p))
		return p
	})(res)
}

func TestAsks(t *testing.T) {
	state := testState{counter: 0}
	ctx := context.WithValue(t.Context(), "multiplier", 7)

	result := Asks(func(c context.Context) StateReaderIOResult[testState, int] {
		mult := c.Value("multiplier").(int)
		return Of[testState](mult * 5)
	})

	res := result(state)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 35, P.Tail(p))
		return p
	})(res)
}

func TestFromResultK(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	validate := func(x int) RES.Result[int] {
		if x > 0 {
			return RES.Of(x * 2)
		}
		return RES.Left[int](errors.New("negative"))
	}

	kleisli := FromResultK[testState](validate)

	// Test with valid input
	resultValid := kleisli(5)
	resValid := resultValid(state)(ctx)()
	assert.True(t, RES.IsRight(resValid))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 10, P.Tail(p))
		return p
	})(resValid)

	// Test with invalid input
	resultInvalid := kleisli(-5)
	resInvalid := resultInvalid(state)(ctx)()
	assert.True(t, RES.IsLeft(resInvalid))
}

func TestFromIOK(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	ioFunc := func(x int) io.IO[int] {
		return func() int { return x * 3 }
	}

	kleisli := FromIOK[testState](ioFunc)
	result := kleisli(7)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 21, P.Tail(p))
		return p
	})(res)
}

func TestFromIOResultK(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	iorFunc := func(x int) IOR.IOResult[int] {
		if x > 0 {
			return IOR.Of(x * 4)
		}
		return IOR.Left[int](errors.New("invalid"))
	}

	kleisli := FromIOResultK[testState](iorFunc)
	result := kleisli(3)
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 12, P.Tail(p))
		return p
	})(res)
}

func TestChainResultK(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	validate := func(x int) RES.Result[string] {
		if x > 0 {
			return RES.Of(fmt.Sprintf("valid: %d", x))
		}
		return RES.Left[string](errors.New("invalid"))
	}

	result := F.Pipe1(
		Of[testState](42),
		ChainResultK[testState](validate),
	)

	res := result(state)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "valid: 42", P.Tail(p))
		return p
	})(res)
}

func TestChainIOResultK(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	iorFunc := func(x int) IOR.IOResult[string] {
		return IOR.Of(fmt.Sprintf("result: %d", x))
	}

	result := F.Pipe1(
		Of[testState](100),
		ChainIOResultK[testState](iorFunc),
	)

	res := result(state)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "result: 100", P.Tail(p))
		return p
	})(res)
}

func TestDo(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	type Result struct {
		value int
	}

	result := Do[testState](Result{})
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, Result]) P.Pair[testState, Result] {
		assert.Equal(t, 0, P.Tail(p).value)
		return p
	})(res)
}

func TestBindTo(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	type Result struct {
		value int
	}

	result := F.Pipe1(
		Of[testState](42),
		BindTo[testState](func(v int) Result {
			return Result{value: v}
		}),
	)

	res := result(state)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, Result]) P.Pair[testState, Result] {
		assert.Equal(t, 42, P.Tail(p).value)
		return p
	})(res)
}

func TestStatefulComputation(t *testing.T) {
	initialState := testState{counter: 0}
	ctx := t.Context()

	// Create a computation that modifies state
	incrementAndGet := func(s testState) P.Pair[testState, int] {
		newState := testState{counter: s.counter + 1}
		return P.MakePair(newState, newState.counter)
	}

	// Chain multiple stateful operations
	result := F.Pipe2(
		FromState(incrementAndGet),
		Chain(func(v1 int) StateReaderIOResult[testState, int] {
			return FromState(incrementAndGet)
		}),
		Chain(func(v2 int) StateReaderIOResult[testState, int] {
			return FromState(incrementAndGet)
		}),
	)

	res := result(initialState)(ctx)()
	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, int]) P.Pair[testState, int] {
		assert.Equal(t, 3, P.Tail(p))         // Last incremented value
		assert.Equal(t, 3, P.Head(p).counter) // State updated three times
		return p
	})(res)
}

func TestErrorPropagation(t *testing.T) {
	state := testState{counter: 0}
	ctx := t.Context()

	testErr := errors.New("test error")

	// Chain operations where the second one fails
	result := F.Pipe1(
		Of[testState](42),
		Chain(func(x int) StateReaderIOResult[testState, int] {
			return Left[testState, int](testErr)
		}),
	)

	res := result(state)(ctx)()
	assert.True(t, RES.IsLeft(res))
}

func TestPointed(t *testing.T) {
	p := Pointed[testState, int]()
	assert.NotNil(t, p)

	result := p.Of(42)
	state := testState{counter: 0}
	ctx := t.Context()
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
}

func TestFunctor(t *testing.T) {
	f := Functor[testState, int, string]()
	assert.NotNil(t, f)

	mapper := f.Map(func(x int) string { return fmt.Sprintf("%d", x) })
	result := mapper(Of[testState](42))

	state := testState{counter: 0}
	ctx := t.Context()
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}

func TestApplicative(t *testing.T) {
	a := Applicative[testState, int, string]()
	assert.NotNil(t, a)

	fab := Of[testState](func(x int) string { return fmt.Sprintf("%d", x) })
	fa := Of[testState](42)
	result := a.Ap(fa)(fab)

	state := testState{counter: 0}
	ctx := t.Context()
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}

func TestMonad(t *testing.T) {
	m := Monad[testState, int, string]()
	assert.NotNil(t, m)

	fa := m.Of(42)
	result := m.Chain(func(x int) StateReaderIOResult[testState, string] {
		return Of[testState](fmt.Sprintf("%d", x))
	})(fa)

	state := testState{counter: 0}
	ctx := t.Context()
	res := result(state)(ctx)()

	assert.True(t, RES.IsRight(res))
	RES.Map(func(p P.Pair[testState, string]) P.Pair[testState, string] {
		assert.Equal(t, "42", P.Tail(p))
		return p
	})(res)
}
