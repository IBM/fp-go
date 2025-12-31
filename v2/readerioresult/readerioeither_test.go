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
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	value int
}

func TestMonadMap(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadMap(Of[testContext](5), N.Mul(2))
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadMapTo(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadMapTo(Of[testContext](5), 42)
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestMapTo(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(Of[testContext](5), MapTo[testContext, int](42))
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestMonadChainFirst(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainFirst(
		Of[testContext](5),
		func(x int) ReaderIOResult[testContext, string] {
			return Of[testContext](fmt.Sprintf("%d", x))
		},
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestChainFirst(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		ChainFirst(func(x int) ReaderIOResult[testContext, string] {
			return Of[testContext](fmt.Sprintf("%d", x))
		}),
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestMonadChainEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainEitherK(
		Of[testContext](5),
		func(x int) E.Either[error, int] {
			return result.Of(x * 2)
		},
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadChainFirstEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainFirstEitherK(
		Of[testContext](5),
		func(x int) E.Either[error, string] {
			return result.Of(fmt.Sprintf("%d", x))
		},
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestChainFirstEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		ChainFirstEitherK[testContext](func(x int) E.Either[error, string] {
			return result.Of(fmt.Sprintf("%d", x))
		}),
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestMonadChainReaderK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainReaderK(
		Of[testContext](5),
		func(x int) R.Reader[testContext, int] {
			return func(c testContext) int { return x + c.value }
		},
	)
	assert.Equal(t, result.Of(15), res(ctx)())
}

func TestMonadChainIOEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainIOEitherK(
		Of[testContext](5),
		func(x int) IOResult[int] {
			return ioresult.Of(x * 2)
		},
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestChainIOEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		ChainIOEitherK[testContext](func(x int) IOResult[int] {
			return ioresult.Of(x * 2)
		}),
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadChainIOK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainIOK(
		Of[testContext](5),
		func(x int) IO[int] {
			return func() int { return x * 2 }
		},
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestChainIOK(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		ChainIOK[testContext](func(x int) IO[int] {
			return func() int { return x * 2 }
		}),
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadChainFirstIOK(t *testing.T) {
	ctx := testContext{value: 10}
	res := MonadChainFirstIOK(
		Of[testContext](5),
		func(x int) IO[string] {
			return func() string { return fmt.Sprintf("%d", x) }
		},
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestChainFirstIOK(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		ChainFirstIOK[testContext](func(x int) IO[string] {
			return func() string { return fmt.Sprintf("%d", x) }
		}),
	)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestChainOptionK(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Some
	resultSome := F.Pipe1(
		Of[testContext](5),
		ChainOptionK[testContext, int, int](func() error {
			return errors.New("none")
		})(func(x int) Option[int] {
			return O.Some(x * 2)
		}),
	)
	assert.Equal(t, result.Of(10), resultSome(ctx)())

	// Test with None
	resultNone := F.Pipe1(
		Of[testContext](5),
		ChainOptionK[testContext, int, int](func() error {
			return errors.New("none")
		})(func(x int) Option[int] {
			return O.None[int]()
		}),
	)
	assert.True(t, E.IsLeft(resultNone(ctx)()))
}

func TestMonadApSeq(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext](N.Mul(2))
	fa := Of[testContext](5)
	res := MonadApSeq(fab, fa)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadApPar(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext](N.Mul(2))
	fa := Of[testContext](5)
	res := MonadApPar(fab, fa)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestChain(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		Chain(func(x int) ReaderIOResult[testContext, int] {
			return Of[testContext](x * 2)
		}),
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestThrowError(t *testing.T) {
	ctx := testContext{value: 10}
	result := ThrowError[testContext, int](errors.New("test error"))
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestFlatten(t *testing.T) {
	ctx := testContext{value: 10}
	nested := Of[testContext](Of[testContext](5))
	res := Flatten(nested)
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestFromEither(t *testing.T) {
	ctx := testContext{value: 10}
	res := FromEither[testContext](result.Of(5))
	assert.Equal(t, result.Of(5), res(ctx)())
}

func TestRightReader(t *testing.T) {
	ctx := testContext{value: 10}
	rdr := func(c testContext) int { return c.value }
	res := RightReader(rdr)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestLeftReader(t *testing.T) {
	ctx := testContext{value: 10}
	reader := func(c testContext) error { return errors.New("test") }
	res := LeftReader[int](reader)
	assert.True(t, E.IsLeft(res(ctx)()))
}

func TestRightIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() int { return 42 }
	res := RightIO[testContext](ioVal)
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestLeftIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() error { return errors.New("test") }
	res := LeftIO[testContext, int](ioVal)
	assert.True(t, E.IsLeft(res(ctx)()))
}

func TestFromIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() int { return 42 }
	res := FromIO[testContext](ioVal)
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestFromIOEither(t *testing.T) {
	ctx := testContext{value: 10}
	ioe := ioresult.Of(42)
	res := FromIOEither[testContext](ioe)
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestFromReaderEither(t *testing.T) {
	ctx := testContext{value: 10}
	re := RE.Of[testContext, error](42)
	res := FromReaderEither(re)
	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestAsk(t *testing.T) {
	ctx := testContext{value: 10}
	res := Ask[testContext]()
	assert.Equal(t, result.Of(ctx), res(ctx)())
}

func TestAsks(t *testing.T) {
	ctx := testContext{value: 10}
	res := Asks(func(c testContext) int { return c.value })
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestFromOption(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Some
	resultSome := FromOption[testContext, int](func() error {
		return errors.New("none")
	})(O.Some(42))
	assert.Equal(t, result.Of(42), resultSome(ctx)())

	// Test with None
	resultNone := FromOption[testContext, int](func() error {
		return errors.New("none")
	})(O.None[int]())
	assert.True(t, E.IsLeft(resultNone(ctx)()))
}

func TestFromPredicate(t *testing.T) {
	ctx := testContext{value: 10}

	// Test predicate true
	resultTrue := FromPredicate[testContext](
		N.MoreThan(0),
		func(x int) error { return errors.New("negative") },
	)(5)
	assert.Equal(t, result.Of(5), resultTrue(ctx)())

	// Test predicate false
	resultFalse := FromPredicate[testContext](
		N.MoreThan(0),
		func(x int) error { return errors.New("negative") },
	)(-5)
	assert.True(t, E.IsLeft(resultFalse(ctx)()))
}

func TestFold(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := Fold(
		func(e error) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext]("error: " + e.Error())
		},
		func(x int) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext](fmt.Sprintf("value: %d", x))
		},
	)(Of[testContext](42))
	assert.Equal(t, "value: 42", resultRight(ctx)())

	// Test Left case
	resultLeft := Fold(
		func(e error) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext]("error: " + e.Error())
		},
		func(x int) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext](fmt.Sprintf("value: %d", x))
		},
	)(Left[testContext, int](errors.New("test")))
	assert.Equal(t, "error: test", resultLeft(ctx)())
}

func TestGetOrElse(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := GetOrElse(func(e error) RIO.ReaderIO[testContext, int] {
		return RIO.Of[testContext](0)
	})(Of[testContext](42))
	assert.Equal(t, 42, resultRight(ctx)())

	// Test Left case
	resultLeft := GetOrElse(func(e error) RIO.ReaderIO[testContext, int] {
		return RIO.Of[testContext](0)
	})(Left[testContext, int](errors.New("test")))
	assert.Equal(t, 0, resultLeft(ctx)())
}

func TestMonadBiMap(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := MonadBiMap(
		Of[testContext](5),
		error.Error,
		strconv.Itoa,
	)
	assert.Equal(t, E.Of[string]("5"), resultRight(ctx)())

	// Test Left case
	resultLeft := MonadBiMap(
		Left[testContext, int](errors.New("test")),
		error.Error,
		strconv.Itoa,
	)
	assert.Equal(t, E.Left[string]("test"), resultLeft(ctx)())
}

func TestBiMap(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](5),
		BiMap[testContext](
			error.Error,
			strconv.Itoa,
		),
	)
	assert.Equal(t, E.Of[string]("5"), res(ctx)())
}

func TestSwap(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right becomes Left
	resultRight := Swap(Of[testContext](5))
	res := resultRight(ctx)()
	assert.True(t, E.IsLeft(res))

	// Test Left becomes Right
	resultLeft := Swap(Left[testContext, int](errors.New("test")))
	assert.True(t, E.IsRight(resultLeft(ctx)()))
}

func TestDefer(t *testing.T) {
	ctx := testContext{value: 10}
	callCount := 0
	res := Defer(func() ReaderIOResult[testContext, int] {
		callCount++
		return Of[testContext](42)
	})

	// First call
	assert.Equal(t, result.Of(42), res(ctx)())
	assert.Equal(t, 1, callCount)

	// Second call
	assert.Equal(t, result.Of(42), res(ctx)())
	assert.Equal(t, 2, callCount)
}

func TestMonadAlt(t *testing.T) {
	ctx := testContext{value: 10}

	// Test first succeeds
	resultFirst := MonadAlt(
		Of[testContext](42),
		func() ReaderIOResult[testContext, int] {
			return Of[testContext](99)
		},
	)
	assert.Equal(t, result.Of(42), resultFirst(ctx)())

	// Test first fails, second succeeds
	resultSecond := MonadAlt(
		Left[testContext, int](errors.New("first")),
		func() ReaderIOResult[testContext, int] {
			return Of[testContext](99)
		},
	)
	assert.Equal(t, result.Of(99), resultSecond(ctx)())
}

func TestAlt(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Left[testContext, int](errors.New("first")),
		Alt(func() ReaderIOResult[testContext, int] {
			return Of[testContext](99)
		}),
	)
	assert.Equal(t, result.Of(99), res(ctx)())
}

func TestMemoize(t *testing.T) {
	ctx := testContext{value: 10}
	callCount := 0
	res := Memoize(func(c testContext) IOResult[int] {
		return func() E.Either[error, int] {
			callCount++
			return result.Of(c.value * 2)
		}
	})

	// First call
	assert.Equal(t, result.Of(20), res(ctx)())
	assert.Equal(t, 1, callCount)

	// Second call should use memoized value
	assert.Equal(t, result.Of(20), res(ctx)())
	assert.Equal(t, 1, callCount)
}

func TestMonadFlap(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext](N.Mul(2))
	res := MonadFlap(fab, 5)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestFlap(t *testing.T) {
	ctx := testContext{value: 10}
	res := F.Pipe1(
		Of[testContext](N.Mul(2)),
		Flap[testContext, int](5),
	)
	assert.Equal(t, result.Of(10), res(ctx)())
}

func TestMonadMapLeft(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadMapLeft(
		Left[testContext, int](errors.New("test")),
		func(e error) string { return e.Error() + "!" },
	)
	res := result(ctx)()
	assert.True(t, E.IsLeft(res))
	// Verify the error was transformed
	E.Fold(
		func(s string) int {
			assert.Equal(t, "test!", s)
			return 0
		},
		func(i int) int { return i },
	)(res)
}

func TestMapLeft(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Left[testContext, int](errors.New("test")),
		MapLeft[testContext, int](func(e error) string { return e.Error() + "!" }),
	)
	res := result(ctx)()
	assert.True(t, E.IsLeft(res))
	// Verify the error was transformed
	E.Fold(
		func(s string) int {
			assert.Equal(t, "test!", s)
			return 0
		},
		func(i int) int { return i },
	)(res)
}

func TestLocal(t *testing.T) {
	ctx2 := testContext{value: 20}

	rdr := Asks(func(c testContext) int { return c.value })
	res := Local[int](func(c testContext) testContext {
		return testContext{value: c.value * 2}
	})(rdr)

	assert.Equal(t, result.Of(40), res(ctx2)())
}

func TestRightReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	rio := func(c testContext) IO[int] {
		return func() int { return c.value * 2 }
	}
	res := RightReaderIO(rio)
	assert.Equal(t, result.Of(20), res(ctx)())
}

func TestLeftReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	rio := func(c testContext) IO[error] {
		return func() error { return errors.New("test") }
	}
	res := LeftReaderIO[int](rio)
	assert.True(t, E.IsLeft(res(ctx)()))
}

func TestLet(t *testing.T) {
	type State struct {
		a int
		b string
	}

	ctx := context.Background()
	res := F.Pipe2(
		Do[context.Context](State{}),
		Let[context.Context](func(b string) func(State) State {
			return func(s State) State { return State{a: s.a, b: b} }
		}, func(s State) string { return "test" }),
		Map[context.Context](func(s State) string { return s.b }),
	)

	assert.Equal(t, result.Of("test"), res(ctx)())
}

func TestLetTo(t *testing.T) {
	type State struct {
		a int
		b string
	}

	ctx := context.Background()
	res := F.Pipe2(
		Do[context.Context](State{}),
		LetTo[context.Context](func(b string) func(State) State {
			return func(s State) State { return State{a: s.a, b: b} }
		}, "constant"),
		Map[context.Context](func(s State) string { return s.b }),
	)

	assert.Equal(t, result.Of("constant"), res(ctx)())
}

func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	ctx := context.Background()
	res := F.Pipe2(
		Of[context.Context](42),
		BindTo[context.Context](func(v int) State { return State{value: v} }),
		Map[context.Context](func(s State) int { return s.value }),
	)

	assert.Equal(t, result.Of(42), res(ctx)())
}

func TestBracket(t *testing.T) {
	ctx := testContext{value: 10}
	released := false

	res := Bracket(
		Of[testContext](42),
		func(x int) ReaderIOResult[testContext, string] {
			return Of[testContext](fmt.Sprintf("%d", x))
		},
		func(x int, result E.Either[error, string]) ReaderIOResult[testContext, int] {
			released = true
			return Of[testContext](0)
		},
	)

	assert.Equal(t, result.Of("42"), res(ctx)())
	assert.True(t, released)
}

func TestWithResource(t *testing.T) {
	ctx := testContext{value: 10}
	released := false

	res := WithResource[string](
		Of[testContext](42),
		func(x int) ReaderIOResult[testContext, int] {
			released = true
			return Of[testContext](0)
		},
	)(func(x int) ReaderIOResult[testContext, string] {
		return Of[testContext](fmt.Sprintf("%d", x))
	})

	assert.Equal(t, result.Of("42"), res(ctx)())
	assert.True(t, released)
}

func TestMonad(t *testing.T) {
	m := Monad[testContext, int, string]()
	assert.NotNil(t, m)
}

func TestTraverseArrayWithIndex(t *testing.T) {
	ctx := testContext{value: 10}
	res := TraverseArrayWithIndex(func(i int, x int) ReaderIOResult[testContext, int] {
		return Of[testContext](x + i)
	})([]int{1, 2, 3})

	assert.Equal(t, result.Of([]int{1, 3, 5}), res(ctx)())
}

func TestTraverseRecord(t *testing.T) {
	ctx := testContext{value: 10}
	res := TraverseRecord[string](func(x int) ReaderIOResult[testContext, int] {
		return Of[testContext](x * 2)
	})(map[string]int{"a": 1, "b": 2})

	expected := map[string]int{"a": 2, "b": 4}
	assert.Equal(t, result.Of(expected), res(ctx)())
}

func TestTraverseRecordWithIndex(t *testing.T) {
	ctx := testContext{value: 10}
	res := TraverseRecordWithIndex(func(k string, x int) ReaderIOResult[testContext, string] {
		return Of[testContext](fmt.Sprintf("%s:%d", k, x))
	})(map[string]int{"a": 1, "b": 2})

	assert.True(t, E.IsRight(res(ctx)()))
}

func TestSequenceRecord(t *testing.T) {
	ctx := testContext{value: 10}
	res := SequenceRecord(map[string]ReaderIOResult[testContext, int]{
		"a": Of[testContext](1),
		"b": Of[testContext](2),
	})

	expected := map[string]int{"a": 1, "b": 2}
	assert.Equal(t, result.Of(expected), res(ctx)())
}

func TestSequenceT1(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT1(Of[testContext](42))
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestSequenceT3(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT3(
		Of[testContext](1),
		Of[testContext]("a"),
		Of[testContext](true),
	)
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestSequenceT4(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT4(
		Of[testContext](1),
		Of[testContext]("a"),
		Of[testContext](true),
		Of[testContext](3.14),
	)
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestWithLock(t *testing.T) {
	ctx := testContext{value: 10}
	unlocked := false

	res := F.Pipe1(
		Of[testContext](42),
		WithLock[testContext, int](func() context.CancelFunc {
			return func() { unlocked = true }
		}),
	)

	assert.Equal(t, result.Of(42), res(ctx)())
	assert.True(t, unlocked)
}

func TestOrElse(t *testing.T) {
	type Config struct {
		fallbackValue int
	}
	ctx := Config{fallbackValue: 99}

	// Test OrElse with Right - should pass through unchanged
	rightValue := Of[Config](42)
	recover := OrElse(func(err error) ReaderIOResult[Config, int] {
		return Left[Config, int](errors.New("should not be called"))
	})
	res := recover(rightValue)(ctx)()
	assert.Equal(t, result.Of(42), res)

	// Test OrElse with Left - should recover with fallback
	leftValue := Left[Config, int](errors.New("not found"))
	recoverWithFallback := OrElse(func(err error) ReaderIOResult[Config, int] {
		if err.Error() == "not found" {
			return func(cfg Config) IOResult[int] {
				return func() result.Result[int] {
					return result.Of(cfg.fallbackValue)
				}
			}
		}
		return Left[Config, int](err)
	})
	res = recoverWithFallback(leftValue)(ctx)()
	assert.Equal(t, result.Of(99), res)

	// Test OrElse with Left - should propagate other errors
	leftValue = Left[Config, int](errors.New("fatal error"))
	res = recoverWithFallback(leftValue)(ctx)()
	assert.True(t, result.IsLeft(res))
	val, err := result.UnwrapError(res)
	assert.Equal(t, 0, val)
	assert.Equal(t, "fatal error", err.Error())
}
