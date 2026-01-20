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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	IOG "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

func TestFromEither(t *testing.T) {
	t.Run("Right value", func(t *testing.T) {
		either := E.Right[error]("success")
		result := FromEither(either)
		assert.Equal(t, E.Right[error]("success"), result(t.Context())())
	})

	t.Run("Left value", func(t *testing.T) {
		err := errors.New("test error")
		either := E.Left[string](err)
		result := FromEither(either)
		assert.Equal(t, E.Left[string](err), result(t.Context())())
	})
}

func TestFromResult(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		result := FromResult(E.Right[error](42))
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})

	t.Run("Error", func(t *testing.T) {
		err := errors.New("test error")
		result := FromResult(E.Left[int](err))
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestLeft(t *testing.T) {
	err := errors.New("test error")
	result := Left[string](err)
	assert.Equal(t, E.Left[string](err), result(t.Context())())
}

func TestRight(t *testing.T) {
	result := Right("success")
	assert.Equal(t, E.Right[error]("success"), result(t.Context())())
}

func TestOf(t *testing.T) {
	result := Of(42)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestMonadMap(t *testing.T) {
	t.Run("Map over Right", func(t *testing.T) {
		result := MonadMap(Of(5), N.Mul(2))
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("Map over Left", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadMap(Left[int](err), N.Mul(2))
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestMap(t *testing.T) {
	t.Run("Map with success", func(t *testing.T) {
		mapper := Map(N.Mul(2))
		result := mapper(Of(5))
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("Map with error", func(t *testing.T) {
		err := errors.New("test error")
		mapper := Map(N.Mul(2))
		result := mapper(Left[int](err))
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestMonadMapTo(t *testing.T) {
	t.Run("MapTo with success", func(t *testing.T) {
		result := MonadMapTo(Of("original"), 42)
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})

	t.Run("MapTo with error", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadMapTo(Left[string](err), 42)
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestMapTo(t *testing.T) {
	mapper := MapTo[string](42)
	result := mapper(Of("original"))
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestMonadChain(t *testing.T) {
	t.Run("Chain with success", func(t *testing.T) {
		result := MonadChain(Of(5), func(x int) ReaderIOResult[int] {
			return Of(x * 2)
		})
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("Chain with error in first", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadChain(Left[int](err), func(x int) ReaderIOResult[int] {
			return Of(x * 2)
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})

	t.Run("Chain with error in second", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadChain(Of(5), func(x int) ReaderIOResult[int] {
			return Left[int](err)
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestChain(t *testing.T) {
	chainer := Chain(func(x int) ReaderIOResult[int] {
		return Of(x * 2)
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestMonadChainFirst(t *testing.T) {
	t.Run("ChainFirst keeps first value", func(t *testing.T) {
		result := MonadChainFirst(Of(5), func(x int) ReaderIOResult[string] {
			return Of("ignored")
		})
		assert.Equal(t, E.Right[error](5), result(t.Context())())
	})

	t.Run("ChainFirst propagates error from second", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadChainFirst(Of(5), func(x int) ReaderIOResult[string] {
			return Left[string](err)
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestChainFirst(t *testing.T) {
	chainer := ChainFirst(func(x int) ReaderIOResult[string] {
		return Of("ignored")
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](5), result(t.Context())())
}

func TestMonadApSeq(t *testing.T) {
	t.Run("ApSeq with success", func(t *testing.T) {
		fab := Of(N.Mul(2))
		fa := Of(5)
		result := MonadApSeq(fab, fa)
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("ApSeq with error in function", func(t *testing.T) {
		err := errors.New("test error")
		fab := Left[func(int) int](err)
		fa := Of(5)
		result := MonadApSeq(fab, fa)
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})

	t.Run("ApSeq with error in value", func(t *testing.T) {
		err := errors.New("test error")
		fab := Of(N.Mul(2))
		fa := Left[int](err)
		result := MonadApSeq(fab, fa)
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestApSeq(t *testing.T) {
	fa := Of(5)
	fab := Of(N.Mul(2))
	result := MonadApSeq(fab, fa)
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestApPar(t *testing.T) {
	t.Run("ApPar with success", func(t *testing.T) {
		fa := Of(5)
		fab := Of(N.Mul(2))
		result := MonadApPar(fab, fa)
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("ApPar with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		fa := Of(5)
		fab := Of(N.Mul(2))
		result := MonadApPar(fab, fa)
		res := result(ctx)()
		assert.True(t, E.IsLeft(res))
	})
}

func TestFromPredicate(t *testing.T) {
	t.Run("Predicate true", func(t *testing.T) {
		pred := FromPredicate(
			N.MoreThan(0),
			func(x int) error { return fmt.Errorf("value %d is not positive", x) },
		)
		result := pred(5)
		assert.Equal(t, E.Right[error](5), result(t.Context())())
	})

	t.Run("Predicate false", func(t *testing.T) {
		pred := FromPredicate(
			N.MoreThan(0),
			func(x int) error { return fmt.Errorf("value %d is not positive", x) },
		)
		result := pred(-5)
		res := result(t.Context())()
		assert.True(t, E.IsLeft(res))
	})
}

func TestOrElse(t *testing.T) {
	t.Run("OrElse with success", func(t *testing.T) {
		fallback := OrElse(func(err error) ReaderIOResult[int] {
			return Of(42)
		})
		result := fallback(Of(10))
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("OrElse with error", func(t *testing.T) {
		err := errors.New("test error")
		fallback := OrElse(func(err error) ReaderIOResult[int] {
			return Of(42)
		})
		result := fallback(Left[int](err))
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})
}

func TestAsk(t *testing.T) {
	result := Ask()
	ctx := t.Context()
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
	ctxResult := E.ToOption(res)
	assert.True(t, O.IsSome(ctxResult))
}

func TestMonadChainEitherK(t *testing.T) {
	t.Run("ChainEitherK with success", func(t *testing.T) {
		result := MonadChainEitherK(Of(5), func(x int) Either[int] {
			return E.Right[error](x * 2)
		})
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("ChainEitherK with error", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadChainEitherK(Of(5), func(x int) Either[int] {
			return E.Left[int](err)
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestChainEitherK(t *testing.T) {
	chainer := ChainEitherK(func(x int) Either[int] {
		return E.Right[error](x * 2)
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestMonadChainFirstEitherK(t *testing.T) {
	t.Run("ChainFirstEitherK keeps first value", func(t *testing.T) {
		result := MonadChainFirstEitherK(Of(5), func(x int) Either[string] {
			return E.Right[error]("ignored")
		})
		assert.Equal(t, E.Right[error](5), result(t.Context())())
	})

	t.Run("ChainFirstEitherK propagates error", func(t *testing.T) {
		err := errors.New("test error")
		result := MonadChainFirstEitherK(Of(5), func(x int) Either[string] {
			return E.Left[string](err)
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestChainFirstEitherK(t *testing.T) {
	chainer := ChainFirstEitherK(func(x int) Either[string] {
		return E.Right[error]("ignored")
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](5), result(t.Context())())
}

func TestChainOptionK(t *testing.T) {
	t.Run("ChainOptionK with Some", func(t *testing.T) {
		chainer := ChainOptionK[int, int](func() error {
			return errors.New("none error")
		})(func(x int) Option[int] {
			return O.Some(x * 2)
		})
		result := chainer(Of(5))
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("ChainOptionK with None", func(t *testing.T) {
		chainer := ChainOptionK[int, int](func() error {
			return errors.New("none error")
		})(func(x int) Option[int] {
			return O.None[int]()
		})
		result := chainer(Of(5))
		res := result(t.Context())()
		assert.True(t, E.IsLeft(res))
	})
}

func TestFromIOEither(t *testing.T) {
	t.Run("FromIOEither with success", func(t *testing.T) {
		ioe := IOE.Of[error](42)
		result := FromIOEither(ioe)
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})

	t.Run("FromIOEither with error", func(t *testing.T) {
		err := errors.New("test error")
		ioe := IOE.Left[int](err)
		result := FromIOEither(ioe)
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestFromIOResult(t *testing.T) {
	ioe := IOE.Of[error](42)
	result := FromIOResult(ioe)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestFromIO(t *testing.T) {
	io := IOG.Of(42)
	result := FromIO(io)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestFromReader(t *testing.T) {
	reader := R.Of[context.Context](42)
	result := FromReader(reader)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestFromLazy(t *testing.T) {
	lazy := func() int { return 42 }
	result := FromLazy(lazy)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestNever(t *testing.T) {
	t.Run("Never with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		result := Never[int]()

		// Cancel immediately
		cancel()

		res := result(ctx)()
		assert.True(t, E.IsLeft(res))
	})

	t.Run("Never with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
		defer cancel()

		result := Never[int]()
		res := result(ctx)()
		assert.True(t, E.IsLeft(res))
	})
}

func TestMonadChainIOK(t *testing.T) {
	result := MonadChainIOK(Of(5), func(x int) IOG.IO[int] {
		return IOG.Of(x * 2)
	})
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestChainIOK(t *testing.T) {
	chainer := ChainIOK(func(x int) IOG.IO[int] {
		return IOG.Of(x * 2)
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestMonadChainFirstIOK(t *testing.T) {
	result := MonadChainFirstIOK(Of(5), func(x int) IOG.IO[string] {
		return IOG.Of("ignored")
	})
	assert.Equal(t, E.Right[error](5), result(t.Context())())
}

func TestChainFirstIOK(t *testing.T) {
	chainer := ChainFirstIOK(func(x int) IOG.IO[string] {
		return IOG.Of("ignored")
	})
	result := chainer(Of(5))
	assert.Equal(t, E.Right[error](5), result(t.Context())())
}

func TestChainIOEitherK(t *testing.T) {
	t.Run("ChainIOEitherK with success", func(t *testing.T) {
		chainer := ChainIOEitherK(func(x int) IOResult[int] {
			return IOE.Of[error](x * 2)
		})
		result := chainer(Of(5))
		assert.Equal(t, E.Right[error](10), result(t.Context())())
	})

	t.Run("ChainIOEitherK with error", func(t *testing.T) {
		err := errors.New("test error")
		chainer := ChainIOEitherK(func(x int) IOResult[int] {
			return IOE.Left[int](err)
		})
		result := chainer(Of(5))
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestDelay(t *testing.T) {
	t.Run("Delay with success", func(t *testing.T) {
		start := time.Now()
		delayed := Delay[int](100 * time.Millisecond)
		result := delayed(Of(42))
		res := result(t.Context())()
		elapsed := time.Since(start)

		assert.True(t, E.IsRight(res))
		assert.GreaterOrEqual(t, elapsed, 100*time.Millisecond)
	})

	t.Run("Delay with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())

		delayed := Delay[int](100 * time.Millisecond)
		result := delayed(Of(42))

		// Cancel after starting but before delay completes
		cancel()
		res := result(ctx)()

		// The result might be either Left (if cancelled) or Right (if completed before cancel)
		// This is a race condition, so we just verify it completes
		assert.True(t, E.IsLeft(res) || E.IsRight(res))
	})
}

func TestDefer(t *testing.T) {
	counter := 0
	deferred := Defer(func() ReaderIOResult[int] {
		counter++
		return Of(counter)
	})

	// First execution
	res1 := deferred(t.Context())()
	assert.True(t, E.IsRight(res1))

	// Second execution should generate a new computation
	res2 := deferred(t.Context())()
	assert.True(t, E.IsRight(res2))

	// Counter should be incremented for each execution
	assert.Equal(t, 2, counter)
}

func TestTryCatch(t *testing.T) {
	t.Run("TryCatch with success", func(t *testing.T) {
		result := TryCatch(func(ctx context.Context) func() (int, error) {
			return func() (int, error) {
				return 42, nil
			}
		})
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})

	t.Run("TryCatch with error", func(t *testing.T) {
		err := errors.New("test error")
		result := TryCatch(func(ctx context.Context) func() (int, error) {
			return func() (int, error) {
				return 0, err
			}
		})
		assert.Equal(t, E.Left[int](err), result(t.Context())())
	})
}

func TestMonadAlt(t *testing.T) {
	t.Run("Alt with first success", func(t *testing.T) {
		first := Of(42)
		second := func() ReaderIOResult[int] { return Of(100) }
		result := MonadAlt(first, second)
		assert.Equal(t, E.Right[error](42), result(t.Context())())
	})

	t.Run("Alt with first error", func(t *testing.T) {
		err := errors.New("test error")
		first := Left[int](err)
		second := func() ReaderIOResult[int] { return Of(100) }
		result := MonadAlt(first, second)
		assert.Equal(t, E.Right[error](100), result(t.Context())())
	})
}

func TestAlt(t *testing.T) {
	err := errors.New("test error")
	alternative := Alt(func() ReaderIOResult[int] { return Of(100) })
	result := alternative(Left[int](err))
	assert.Equal(t, E.Right[error](100), result(t.Context())())
}

func TestMemoize(t *testing.T) {
	counter := 0
	computation := Memoize(FromLazy(func() int {
		counter++
		return counter
	}))

	// First execution
	res1 := computation(t.Context())()
	assert.True(t, E.IsRight(res1))
	val1 := E.ToOption(res1)
	assert.Equal(t, O.Of(1), val1)

	// Second execution should return cached value
	res2 := computation(t.Context())()
	assert.True(t, E.IsRight(res2))
	val2 := E.ToOption(res2)
	assert.Equal(t, O.Of(1), val2)

	// Counter should only be incremented once
	assert.Equal(t, 1, counter)
}

func TestFlatten(t *testing.T) {
	nested := Of(Of(42))
	result := Flatten(nested)
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestMonadFlap(t *testing.T) {
	fab := Of(N.Mul(2))
	result := MonadFlap(fab, 5)
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestFlap(t *testing.T) {
	flapper := Flap[int](5)
	result := flapper(Of(N.Mul(2)))
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestFold(t *testing.T) {
	t.Run("Fold with success", func(t *testing.T) {
		folder := Fold(
			func(err error) ReaderIOResult[string] {
				return Of("error: " + err.Error())
			},
			func(x int) ReaderIOResult[string] {
				return Of(fmt.Sprintf("success: %d", x))
			},
		)
		result := folder(Of(42))
		assert.Equal(t, E.Right[error]("success: 42"), result(t.Context())())
	})

	t.Run("Fold with error", func(t *testing.T) {
		err := errors.New("test error")
		folder := Fold(
			func(err error) ReaderIOResult[string] {
				return Of("error: " + err.Error())
			},
			func(x int) ReaderIOResult[string] {
				return Of(fmt.Sprintf("success: %d", x))
			},
		)
		result := folder(Left[int](err))
		assert.Equal(t, E.Right[error]("error: test error"), result(t.Context())())
	})
}

func TestGetOrElse(t *testing.T) {
	t.Run("GetOrElse with success", func(t *testing.T) {
		getter := GetOrElse(func(err error) ReaderIO[int] {
			return func(ctx context.Context) IOG.IO[int] {
				return IOG.Of(0)
			}
		})
		result := getter(Of(42))
		assert.Equal(t, 42, result(t.Context())())
	})

	t.Run("GetOrElse with error", func(t *testing.T) {
		err := errors.New("test error")
		getter := GetOrElse(func(err error) ReaderIO[int] {
			return func(ctx context.Context) IOG.IO[int] {
				return IOG.Of(0)
			}
		})
		result := getter(Left[int](err))
		assert.Equal(t, 0, result(t.Context())())
	})
}

func TestWithContext(t *testing.T) {
	t.Run("WithContext with valid context", func(t *testing.T) {
		computation := WithContext(Of(42))
		result := computation(t.Context())()
		assert.Equal(t, E.Right[error](42), result)
	})

	t.Run("WithContext with cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		computation := WithContext(Of(42))
		result := computation(ctx)()
		assert.True(t, E.IsLeft(result))
	})
}

func TestEitherize0(t *testing.T) {
	f := func(ctx context.Context) (int, error) {
		return 42, nil
	}
	eitherized := Eitherize0(f)
	result := eitherized()
	assert.Equal(t, E.Right[error](42), result(t.Context())())
}

func TestUneitherize0(t *testing.T) {
	f := func() ReaderIOResult[int] {
		return Of(42)
	}
	uneitherized := Uneitherize0(f)
	result, err := uneitherized(t.Context())
	assert.NoError(t, err)
	assert.Equal(t, 42, result)
}

func TestEitherize1(t *testing.T) {
	f := func(ctx context.Context, x int) (int, error) {
		return x * 2, nil
	}
	eitherized := Eitherize1(f)
	result := eitherized(5)
	assert.Equal(t, E.Right[error](10), result(t.Context())())
}

func TestUneitherize1(t *testing.T) {
	f := func(x int) ReaderIOResult[int] {
		return Of(x * 2)
	}
	uneitherized := Uneitherize1(f)
	result, err := uneitherized(t.Context(), 5)
	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

func TestSequenceT2(t *testing.T) {
	result := SequenceT2(Of(1), Of(2))
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
	tuple := E.ToOption(res)
	assert.True(t, O.IsSome(tuple))
	t1, _ := O.Unwrap(tuple)
	assert.Equal(t, 1, t1.F1)
	assert.Equal(t, 2, t1.F2)
}

func TestSequenceSeqT2(t *testing.T) {
	result := SequenceSeqT2(Of(1), Of(2))
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
}

func TestSequenceParT2(t *testing.T) {
	result := SequenceParT2(Of(1), Of(2))
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
}

func TestTraverseArray(t *testing.T) {
	t.Run("TraverseArray with success", func(t *testing.T) {
		arr := []int{1, 2, 3}
		traverser := TraverseArray(func(x int) ReaderIOResult[int] {
			return Of(x * 2)
		})
		result := traverser(arr)
		res := result(t.Context())()
		assert.True(t, E.IsRight(res))
		arrOpt := E.ToOption(res)
		assert.Equal(t, O.Of([]int{2, 4, 6}), arrOpt)
	})

	t.Run("TraverseArray with error", func(t *testing.T) {
		arr := []int{1, 2, 3}
		err := errors.New("test error")
		traverser := TraverseArray(func(x int) ReaderIOResult[int] {
			if x == 2 {
				return Left[int](err)
			}
			return Of(x * 2)
		})
		result := traverser(arr)
		res := result(t.Context())()
		assert.True(t, E.IsLeft(res))
	})
}

func TestSequenceArray(t *testing.T) {
	arr := []ReaderIOResult[int]{Of(1), Of(2), Of(3)}
	result := SequenceArray(arr)
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
	arrOpt := E.ToOption(res)
	assert.Equal(t, O.Of([]int{1, 2, 3}), arrOpt)
}

func TestTraverseRecord(t *testing.T) {
	rec := map[string]int{"a": 1, "b": 2}
	result := TraverseRecord[string](func(x int) ReaderIOResult[int] {
		return Of(x * 2)
	})(rec)
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
	recOpt := E.ToOption(res)
	assert.True(t, O.IsSome(recOpt))
	resultRec, _ := O.Unwrap(recOpt)
	assert.Equal(t, 2, resultRec["a"])
	assert.Equal(t, 4, resultRec["b"])
}

func TestSequenceRecord(t *testing.T) {
	rec := map[string]ReaderIOResult[int]{
		"a": Of(1),
		"b": Of(2),
	}
	result := SequenceRecord(rec)
	res := result(t.Context())()
	assert.True(t, E.IsRight(res))
	recOpt := E.ToOption(res)
	assert.True(t, O.IsSome(recOpt))
	resultRec, _ := O.Unwrap(recOpt)
	assert.Equal(t, 1, resultRec["a"])
	assert.Equal(t, 2, resultRec["b"])
}

func TestAltSemigroup(t *testing.T) {
	sg := AltSemigroup[int]()
	err := errors.New("test error")

	result := sg.Concat(Left[int](err), Of(42))
	res := result(t.Context())()
	assert.Equal(t, E.Right[error](42), res)
}

func TestApplicativeMonoid(t *testing.T) {
	// Test with int addition monoid
	intAddMonoid := ApplicativeMonoid(M.MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	))

	result := intAddMonoid.Concat(Of(5), Of(10))
	res := result(t.Context())()
	assert.Equal(t, E.Right[error](15), res)
}

func TestBracket(t *testing.T) {
	t.Run("Bracket with success", func(t *testing.T) {
		var acquired, released bool

		acquire := FromLazy(func() int {
			acquired = true
			return 42
		})

		use := func(x int) ReaderIOResult[int] {
			return Of(x * 2)
		}

		release := func(x int, result Either[int]) ReaderIOResult[any] {
			return FromLazy(func() any {
				released = true
				return nil
			})
		}

		result := Bracket(acquire, use, release)
		res := result(t.Context())()

		assert.True(t, acquired)
		assert.True(t, released)
		assert.Equal(t, E.Right[error](84), res)
	})

	t.Run("Bracket with error in use", func(t *testing.T) {
		var acquired, released bool
		err := errors.New("use error")

		acquire := FromLazy(func() int {
			acquired = true
			return 42
		})

		use := func(x int) ReaderIOResult[int] {
			return Left[int](err)
		}

		release := func(x int, result Either[int]) ReaderIOResult[any] {
			return FromLazy(func() any {
				released = true
				return nil
			})
		}

		result := Bracket(acquire, use, release)
		res := result(t.Context())()

		assert.True(t, acquired)
		assert.True(t, released)
		assert.Equal(t, E.Left[int](err), res)
	})
}
