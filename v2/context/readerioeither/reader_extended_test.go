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
	"time"

	E "github.com/IBM/fp-go/v2/either"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestFromEither(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	rightVal := E.Right[error](42)
	result := FromEither(rightVal)(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left
	err := errors.New("test error")
	leftVal := E.Left[int](err)
	result = FromEither(leftVal)(ctx)()
	assert.Equal(t, E.Left[int](err), result)
}

func TestLeftRight(t *testing.T) {
	ctx := t.Context()

	// Test Left
	err := errors.New("test error")
	result := Left[int](err)(ctx)()
	assert.True(t, E.IsLeft(result))

	// Test Right
	result = Right(42)(ctx)()
	assert.True(t, E.IsRight(result))
	val, _ := E.Unwrap(result)
	assert.Equal(t, 42, val)
}

func TestOf(t *testing.T) {
	ctx := t.Context()
	result := Of(42)(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestMonadMap(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadMap(Right(42), func(x int) int { return x * 2 })(ctx)()
	assert.Equal(t, E.Right[error](84), result)

	// Test with Left
	err := errors.New("test error")
	result = MonadMap(Left[int](err), func(x int) int { return x * 2 })(ctx)()
	assert.Equal(t, E.Left[int](err), result)
}

func TestMonadMapTo(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadMapTo(Right(42), "hello")(ctx)()
	assert.Equal(t, E.Right[error]("hello"), result)

	// Test with Left
	err := errors.New("test error")
	result = MonadMapTo(Left[int](err), "hello")(ctx)()
	assert.Equal(t, E.Left[string](err), result)
}

func TestMonadChain(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChain(Right(42), func(x int) ReaderIOEither[int] {
		return Right(x * 2)
	})(ctx)()
	assert.Equal(t, E.Right[error](84), result)

	// Test with Left
	err := errors.New("test error")
	result = MonadChain(Left[int](err), func(x int) ReaderIOEither[int] {
		return Right(x * 2)
	})(ctx)()
	assert.Equal(t, E.Left[int](err), result)

	// Test where function returns Left
	result = MonadChain(Right(42), func(x int) ReaderIOEither[int] {
		return Left[int](errors.New("chain error"))
	})(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestMonadChainFirst(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChainFirst(Right(42), func(x int) ReaderIOEither[string] {
		return Right("ignored")
	})(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left in first
	err := errors.New("test error")
	result = MonadChainFirst(Left[int](err), func(x int) ReaderIOEither[string] {
		return Right("ignored")
	})(ctx)()
	assert.Equal(t, E.Left[int](err), result)

	// Test with Left in second
	result = MonadChainFirst(Right(42), func(x int) ReaderIOEither[string] {
		return Left[string](errors.New("chain error"))
	})(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestMonadApSeq(t *testing.T) {
	ctx := t.Context()

	// Test with both Right
	fct := Right(func(x int) int { return x * 2 })
	val := Right(42)
	result := MonadApSeq(fct, val)(ctx)()
	assert.Equal(t, E.Right[error](84), result)

	// Test with Left function
	err := errors.New("function error")
	fct = Left[func(int) int](err)
	result = MonadApSeq(fct, val)(ctx)()
	assert.Equal(t, E.Left[int](err), result)

	// Test with Left value
	fct = Right(func(x int) int { return x * 2 })
	err = errors.New("value error")
	val = Left[int](err)
	result = MonadApSeq(fct, val)(ctx)()
	assert.Equal(t, E.Left[int](err), result)
}

func TestMonadApPar(t *testing.T) {
	ctx := t.Context()

	// Test with both Right
	fct := Right(func(x int) int { return x * 2 })
	val := Right(42)
	result := MonadApPar(fct, val)(ctx)()
	assert.Equal(t, E.Right[error](84), result)
}

func TestFromPredicate(t *testing.T) {
	ctx := t.Context()

	pred := func(x int) bool { return x > 0 }
	onFalse := func(x int) error { return fmt.Errorf("value %d is not positive", x) }

	// Test with predicate true
	result := FromPredicate(pred, onFalse)(42)(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with predicate false
	result = FromPredicate(pred, onFalse)(-1)(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestAsk(t *testing.T) {
	ctx := context.WithValue(t.Context(), "key", "value")
	result := Ask()(ctx)()
	assert.True(t, E.IsRight(result))
	retrievedCtx, _ := E.Unwrap(result)
	assert.Equal(t, "value", retrievedCtx.Value("key"))
}

func TestMonadChainEitherK(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChainEitherK(Right(42), func(x int) E.Either[error, int] {
		return E.Right[error](x * 2)
	})(ctx)()
	assert.Equal(t, E.Right[error](84), result)

	// Test with Left in Either
	result = MonadChainEitherK(Right(42), func(x int) E.Either[error, int] {
		return E.Left[int](errors.New("either error"))
	})(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestMonadChainFirstEitherK(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChainFirstEitherK(Right(42), func(x int) E.Either[error, string] {
		return E.Right[error]("ignored")
	})(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left in Either
	result = MonadChainFirstEitherK(Right(42), func(x int) E.Either[error, string] {
		return E.Left[string](errors.New("either error"))
	})(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestChainOptionKFunc(t *testing.T) {
	ctx := t.Context()

	onNone := func() error { return errors.New("none error") }

	// Test with Some
	chainFunc := ChainOptionK[int, int](onNone)
	result := chainFunc(func(x int) O.Option[int] {
		return O.Some(x * 2)
	})(Right(42))(ctx)()
	assert.Equal(t, E.Right[error](84), result)

	// Test with None
	result = chainFunc(func(x int) O.Option[int] {
		return O.None[int]()
	})(Right(42))(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestFromIOEither(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	ioe := func() E.Either[error, int] {
		return E.Right[error](42)
	}
	result := FromIOEither(ioe)(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left
	err := errors.New("test error")
	ioe = func() E.Either[error, int] {
		return E.Left[int](err)
	}
	result = FromIOEither(ioe)(ctx)()
	assert.Equal(t, E.Left[int](err), result)
}

func TestFromIO(t *testing.T) {
	ctx := t.Context()

	io := func() int { return 42 }
	result := FromIO(io)(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestFromLazy(t *testing.T) {
	ctx := t.Context()

	lazy := func() int { return 42 }
	result := FromLazy(lazy)(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestNeverWithCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())

	// Start Never in a goroutine
	done := make(chan E.Either[error, int])
	go func() {
		done <- Never[int]()(ctx)()
	}()

	// Cancel the context
	cancel()

	// Should receive cancellation error
	result := <-done
	assert.True(t, E.IsLeft(result))
}

func TestMonadChainIOK(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChainIOK(Right(42), func(x int) func() int {
		return func() int { return x * 2 }
	})(ctx)()
	assert.Equal(t, E.Right[error](84), result)
}

func TestMonadChainFirstIOK(t *testing.T) {
	ctx := t.Context()

	// Test with Right
	result := MonadChainFirstIOK(Right(42), func(x int) func() string {
		return func() string { return "ignored" }
	})(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestDelayFunc(t *testing.T) {
	ctx := t.Context()
	delay := 100 * time.Millisecond

	start := time.Now()
	delayFunc := Delay[int](delay)
	result := delayFunc(Right(42))(ctx)()
	elapsed := time.Since(start)

	assert.True(t, E.IsRight(result))
	assert.GreaterOrEqual(t, elapsed, delay)
}

func TestDefer(t *testing.T) {
	ctx := t.Context()
	count := 0

	gen := func() ReaderIOEither[int] {
		count++
		return Right(count)
	}

	deferred := Defer(gen)

	// First call
	result1 := deferred(ctx)()
	assert.Equal(t, E.Right[error](1), result1)

	// Second call should generate new value
	result2 := deferred(ctx)()
	assert.Equal(t, E.Right[error](2), result2)
}

func TestTryCatch(t *testing.T) {
	ctx := t.Context()

	// Test success
	result := TryCatch(func(ctx context.Context) func() (int, error) {
		return func() (int, error) {
			return 42, nil
		}
	})(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test error
	err := errors.New("test error")
	result = TryCatch(func(ctx context.Context) func() (int, error) {
		return func() (int, error) {
			return 0, err
		}
	})(ctx)()
	assert.Equal(t, E.Left[int](err), result)
}

func TestMonadAlt(t *testing.T) {
	ctx := t.Context()

	// Test with Right (alternative not called)
	result := MonadAlt(Right(42), func() ReaderIOEither[int] {
		return Right(99)
	})(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left (alternative called)
	err := errors.New("test error")
	result = MonadAlt(Left[int](err), func() ReaderIOEither[int] {
		return Right(99)
	})(ctx)()
	assert.Equal(t, E.Right[error](99), result)
}

func TestMemoize(t *testing.T) {
	ctx := t.Context()
	count := 0

	rdr := Memoize(FromLazy(func() int {
		count++
		return count
	}))

	// First call
	result1 := rdr(ctx)()
	assert.Equal(t, E.Right[error](1), result1)

	// Second call should return memoized value
	result2 := rdr(ctx)()
	assert.Equal(t, E.Right[error](1), result2)
}

func TestFlatten(t *testing.T) {
	ctx := t.Context()

	nested := Right(Right(42))
	result := Flatten(nested)(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestMonadFlap(t *testing.T) {
	ctx := t.Context()
	fab := Right(func(x int) int { return x * 2 })
	result := MonadFlap(fab, 42)(ctx)()
	assert.Equal(t, E.Right[error](84), result)
}

func TestWithContext(t *testing.T) {
	// Test with non-canceled context
	ctx := t.Context()
	result := WithContext(Right(42))(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with canceled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	result = WithContext(Right(42))(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestMonadAp(t *testing.T) {
	ctx := t.Context()

	// Test with both Right
	fct := Right(func(x int) int { return x * 2 })
	val := Right(42)
	result := MonadAp(fct, val)(ctx)()
	assert.Equal(t, E.Right[error](84), result)
}

// Test traverse functions
func TestSequenceArray(t *testing.T) {
	ctx := t.Context()

	// Test with all Right
	arr := []ReaderIOEither[int]{Right(1), Right(2), Right(3)}
	result := SequenceArray(arr)(ctx)()
	assert.True(t, E.IsRight(result))
	vals, _ := E.Unwrap(result)
	assert.Equal(t, []int{1, 2, 3}, vals)

	// Test with one Left
	err := errors.New("test error")
	arr = []ReaderIOEither[int]{Right(1), Left[int](err), Right(3)}
	result = SequenceArray(arr)(ctx)()
	assert.True(t, E.IsLeft(result))
}

func TestTraverseArray(t *testing.T) {
	ctx := t.Context()

	// Test transformation
	arr := []int{1, 2, 3}
	result := TraverseArray(func(x int) ReaderIOEither[int] {
		return Right(x * 2)
	})(arr)(ctx)()
	assert.True(t, E.IsRight(result))
	vals, _ := E.Unwrap(result)
	assert.Equal(t, []int{2, 4, 6}, vals)
}

func TestSequenceRecord(t *testing.T) {
	ctx := t.Context()

	// Test with all Right
	rec := map[string]ReaderIOEither[int]{
		"a": Right(1),
		"b": Right(2),
	}
	result := SequenceRecord(rec)(ctx)()
	assert.True(t, E.IsRight(result))
	vals, _ := E.Unwrap(result)
	assert.Equal(t, 1, vals["a"])
	assert.Equal(t, 2, vals["b"])
}

func TestTraverseRecord(t *testing.T) {
	ctx := t.Context()

	// Test transformation
	rec := map[string]int{"a": 1, "b": 2}
	result := TraverseRecord[string](func(x int) ReaderIOEither[int] {
		return Right(x * 2)
	})(rec)(ctx)()
	assert.True(t, E.IsRight(result))
	vals, _ := E.Unwrap(result)
	assert.Equal(t, 2, vals["a"])
	assert.Equal(t, 4, vals["b"])
}

// Test monoid functions
func TestAltSemigroup(t *testing.T) {
	ctx := t.Context()

	sg := AltSemigroup[int]()

	// Test with Right (first succeeds)
	result := sg.Concat(Right(42), Right(99))(ctx)()
	assert.Equal(t, E.Right[error](42), result)

	// Test with Left then Right (fallback)
	err := errors.New("test error")
	result = sg.Concat(Left[int](err), Right(99))(ctx)()
	assert.Equal(t, E.Right[error](99), result)
}

// Test Do notation
func TestDo(t *testing.T) {
	ctx := t.Context()

	type State struct {
		Value int
	}

	result := Do(State{Value: 42})(ctx)()
	assert.True(t, E.IsRight(result))
	state, _ := E.Unwrap(result)
	assert.Equal(t, 42, state.Value)
}
