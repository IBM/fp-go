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
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

type testContext struct {
	value int
}

func TestMonadMap(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadMap(Of[testContext, error](5), N.Mul(2))
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadMapTo(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadMapTo(Of[testContext, error](5), 42)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestMapTo(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(Of[testContext, error](5), MapTo[testContext, error, int](42))
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestMonadChainFirst(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainFirst(
		Of[testContext, error](5),
		func(x int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%d", x))
		},
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestChainFirst(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		ChainFirst(func(x int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%d", x))
		}),
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestMonadChainEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainEitherK(
		Of[testContext, error](5),
		func(x int) E.Either[error, int] {
			return E.Right[error](x * 2)
		},
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadChainFirstEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainFirstEitherK(
		Of[testContext, error](5),
		func(x int) E.Either[error, string] {
			return E.Right[error](fmt.Sprintf("%d", x))
		},
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestChainFirstEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		ChainFirstEitherK[testContext](func(x int) E.Either[error, string] {
			return E.Right[error](fmt.Sprintf("%d", x))
		}),
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestMonadChainReaderK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainReaderK(
		Of[testContext, error](5),
		func(x int) R.Reader[testContext, int] {
			return func(c testContext) int { return x + c.value }
		},
	)
	assert.Equal(t, E.Right[error](15), result(ctx)())
}

func TestMonadChainIOEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainIOEitherK(
		Of[testContext, error](5),
		func(x int) IOE.IOEither[error, int] {
			return IOE.Right[error](x * 2)
		},
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestChainIOEitherK(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		ChainIOEitherK[testContext](func(x int) IOE.IOEither[error, int] {
			return IOE.Right[error](x * 2)
		}),
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadChainIOK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainIOK(
		Of[testContext, error](5),
		func(x int) io.IO[int] {
			return func() int { return x * 2 }
		},
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestChainIOK(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		ChainIOK[testContext, error](func(x int) io.IO[int] {
			return func() int { return x * 2 }
		}),
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadChainFirstIOK(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadChainFirstIOK(
		Of[testContext, error](5),
		func(x int) io.IO[string] {
			return func() string { return fmt.Sprintf("%d", x) }
		},
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestChainFirstIOK(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		ChainFirstIOK[testContext, error](func(x int) io.IO[string] {
			return func() string { return fmt.Sprintf("%d", x) }
		}),
	)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestChainOptionK(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Some
	resultSome := F.Pipe1(
		Of[testContext, error](5),
		ChainOptionK[testContext, int, int](func() error {
			return errors.New("none")
		})(func(x int) O.Option[int] {
			return O.Some(x * 2)
		}),
	)
	assert.Equal(t, E.Right[error](10), resultSome(ctx)())

	// Test with None
	resultNone := F.Pipe1(
		Of[testContext, error](5),
		ChainOptionK[testContext, int, int](func() error {
			return errors.New("none")
		})(func(x int) O.Option[int] {
			return O.None[int]()
		}),
	)
	assert.True(t, E.IsLeft(resultNone(ctx)()))
}

func TestMonadApSeq(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext, error](N.Mul(2))
	fa := Of[testContext, error](5)
	result := MonadApSeq(fab, fa)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadApPar(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext, error](N.Mul(2))
	fa := Of[testContext, error](5)
	result := MonadApPar(fab, fa)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestChain(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		Chain(func(x int) ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](x * 2)
		}),
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestThrowError(t *testing.T) {
	ctx := testContext{value: 10}
	result := ThrowError[testContext, int](errors.New("test error"))
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestFlatten(t *testing.T) {
	ctx := testContext{value: 10}
	nested := Of[testContext, error](Of[testContext, error](5))
	result := Flatten(nested)
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestFromEither(t *testing.T) {
	ctx := testContext{value: 10}
	result := FromEither[testContext](E.Right[error](5))
	assert.Equal(t, E.Right[error](5), result(ctx)())
}

func TestRightReader(t *testing.T) {
	ctx := testContext{value: 10}
	reader := func(c testContext) int { return c.value }
	result := RightReader[error](reader)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestLeftReader(t *testing.T) {
	ctx := testContext{value: 10}
	reader := func(c testContext) error { return errors.New("test") }
	result := LeftReader[int](reader)
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestRightIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() int { return 42 }
	result := RightIO[testContext, error](ioVal)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestLeftIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() error { return errors.New("test") }
	result := LeftIO[testContext, int](ioVal)
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestFromIO(t *testing.T) {
	ctx := testContext{value: 10}
	ioVal := func() int { return 42 }
	result := FromIO[testContext, error](ioVal)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestFromIOEither(t *testing.T) {
	ctx := testContext{value: 10}
	ioe := IOE.Right[error](42)
	result := FromIOEither[testContext](ioe)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestFromReaderEither(t *testing.T) {
	ctx := testContext{value: 10}
	re := RE.Of[testContext, error](42)
	result := FromReaderEither(re)
	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestAsk(t *testing.T) {
	ctx := testContext{value: 10}
	result := Ask[testContext, error]()
	assert.Equal(t, E.Right[error](ctx), result(ctx)())
}

func TestAsks(t *testing.T) {
	ctx := testContext{value: 10}
	result := Asks[error](func(c testContext) int { return c.value })
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestFromOption(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Some
	resultSome := FromOption[testContext, int](func() error {
		return errors.New("none")
	})(O.Some(42))
	assert.Equal(t, E.Right[error](42), resultSome(ctx)())

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
	assert.Equal(t, E.Right[error](5), resultTrue(ctx)())

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
	)(Of[testContext, error](42))
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
	})(Of[testContext, error](42))
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
		Of[testContext, error](5),
		error.Error,
		strconv.Itoa,
	)
	assert.Equal(t, E.Right[string]("5"), resultRight(ctx)())

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
	result := F.Pipe1(
		Of[testContext, error](5),
		BiMap[testContext](
			error.Error,
			strconv.Itoa,
		),
	)
	assert.Equal(t, E.Right[string]("5"), result(ctx)())
}

func TestSwap(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right becomes Left
	resultRight := Swap(Of[testContext, error](5))
	res := resultRight(ctx)()
	assert.True(t, E.IsLeft(res))

	// Test Left becomes Right
	resultLeft := Swap(Left[testContext, int](errors.New("test")))
	assert.True(t, E.IsRight(resultLeft(ctx)()))
}

func TestDefer(t *testing.T) {
	ctx := testContext{value: 10}
	callCount := 0
	result := Defer(func() ReaderIOEither[testContext, error, int] {
		callCount++
		return Of[testContext, error](42)
	})

	// First call
	assert.Equal(t, E.Right[error](42), result(ctx)())
	assert.Equal(t, 1, callCount)

	// Second call
	assert.Equal(t, E.Right[error](42), result(ctx)())
	assert.Equal(t, 2, callCount)
}

func TestTryCatch(t *testing.T) {
	ctx := testContext{value: 10}

	// Test success
	resultSuccess := TryCatch(
		func(c testContext) func() (int, error) {
			return func() (int, error) { return c.value * 2, nil }
		},
		func(err error) error { return err },
	)
	assert.Equal(t, E.Right[error](20), resultSuccess(ctx)())

	// Test error
	resultError := TryCatch(
		func(c testContext) func() (int, error) {
			return func() (int, error) { return 0, errors.New("test error") }
		},
		func(err error) error { return err },
	)
	assert.True(t, E.IsLeft(resultError(ctx)()))
}

func TestMonadAlt(t *testing.T) {
	ctx := testContext{value: 10}

	// Test first succeeds
	resultFirst := MonadAlt(
		Of[testContext, error](42),
		func() ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](99)
		},
	)
	assert.Equal(t, E.Right[error](42), resultFirst(ctx)())

	// Test first fails, second succeeds
	resultSecond := MonadAlt(
		Left[testContext, int](errors.New("first")),
		func() ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](99)
		},
	)
	assert.Equal(t, E.Right[error](99), resultSecond(ctx)())
}

func TestAlt(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Left[testContext, int](errors.New("first")),
		Alt(func() ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](99)
		}),
	)
	assert.Equal(t, E.Right[error](99), result(ctx)())
}

func TestMemoize(t *testing.T) {
	ctx := testContext{value: 10}
	callCount := 0
	result := Memoize(func(c testContext) IOE.IOEither[error, int] {
		return func() E.Either[error, int] {
			callCount++
			return E.Right[error](c.value * 2)
		}
	})

	// First call
	assert.Equal(t, E.Right[error](20), result(ctx)())
	assert.Equal(t, 1, callCount)

	// Second call should use memoized value
	assert.Equal(t, E.Right[error](20), result(ctx)())
	assert.Equal(t, 1, callCount)
}

func TestMonadFlap(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext, error](N.Mul(2))
	result := MonadFlap(fab, 5)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestFlap(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](N.Mul(2)),
		Flap[testContext, error, int](5),
	)
	assert.Equal(t, E.Right[error](10), result(ctx)())
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

	rdr := Asks[error](func(c testContext) int { return c.value })
	result := Local[error, int](func(c testContext) testContext {
		return testContext{value: c.value * 2}
	})(rdr)

	assert.Equal(t, E.Right[error](40), result(ctx2)())
}

func TestRightReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	rio := func(c testContext) io.IO[int] {
		return func() int { return c.value * 2 }
	}
	result := RightReaderIO[error](rio)
	assert.Equal(t, E.Right[error](20), result(ctx)())
}

func TestLeftReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	rio := func(c testContext) io.IO[error] {
		return func() error { return errors.New("test") }
	}
	result := LeftReaderIO[int](rio)
	assert.True(t, E.IsLeft(result(ctx)()))
}

func TestLet(t *testing.T) {
	type State struct {
		a int
		b string
	}

	ctx := context.Background()
	result := F.Pipe2(
		Do[context.Context, error](State{}),
		Let[context.Context, error](func(b string) func(State) State {
			return func(s State) State { return State{a: s.a, b: b} }
		}, func(s State) string { return "test" }),
		Map[context.Context, error](func(s State) string { return s.b }),
	)

	assert.Equal(t, E.Right[error]("test"), result(ctx)())
}

func TestLetTo(t *testing.T) {
	type State struct {
		a int
		b string
	}

	ctx := context.Background()
	result := F.Pipe2(
		Do[context.Context, error](State{}),
		LetTo[context.Context, error](func(b string) func(State) State {
			return func(s State) State { return State{a: s.a, b: b} }
		}, "constant"),
		Map[context.Context, error](func(s State) string { return s.b }),
	)

	assert.Equal(t, E.Right[error]("constant"), result(ctx)())
}

func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	ctx := context.Background()
	result := F.Pipe2(
		Of[context.Context, error](42),
		BindTo[context.Context, error](func(v int) State { return State{value: v} }),
		Map[context.Context, error](func(s State) int { return s.value }),
	)

	assert.Equal(t, E.Right[error](42), result(ctx)())
}

func TestBracket(t *testing.T) {
	ctx := testContext{value: 10}
	released := false

	result := Bracket(
		Of[testContext, error](42),
		func(x int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%d", x))
		},
		func(x int, result E.Either[error, string]) ReaderIOEither[testContext, error, int] {
			released = true
			return Of[testContext, error](0)
		},
	)

	assert.Equal(t, E.Right[error]("42"), result(ctx)())
	assert.True(t, released)
}

func TestWithResource(t *testing.T) {
	ctx := testContext{value: 10}
	released := false

	result := WithResource[string](
		Of[testContext, error](42),
		func(x int) ReaderIOEither[testContext, error, int] {
			released = true
			return Of[testContext, error](0)
		},
	)(func(x int) ReaderIOEither[testContext, error, string] {
		return Of[testContext, error](fmt.Sprintf("%d", x))
	})

	assert.Equal(t, E.Right[error]("42"), result(ctx)())
	assert.True(t, released)
}

func TestMonad(t *testing.T) {
	m := Monad[testContext, error, int, string]()
	assert.NotNil(t, m)
}

func TestTraverseArrayDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("empty array", func(t *testing.T) {
		f := TraverseArray(func(a int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](strconv.Itoa(a))
		})
		result := f([]int{})
		assert.Equal(t, E.Right[error]([]string{}), result(ctx)())
	})

	t.Run("successful transformation", func(t *testing.T) {
		f := TraverseArray(func(a int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](strconv.Itoa(a * 2))
		})
		result := f([]int{1, 2, 3})
		assert.Equal(t, E.Right[error]([]string{"2", "4", "6"}), result(ctx)())
	})

	t.Run("first element fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("first error")
		f := TraverseArray(func(a int) ReaderIOEither[testContext, error, string] {
			if a == 1 {
				return Left[testContext, string](expectedErr)
			}
			return Of[testContext, error](strconv.Itoa(a))
		})
		result := f([]int{1, 2, 3})
		assert.Equal(t, E.Left[[]string](expectedErr), result(ctx)())
	})

	t.Run("middle element fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("middle error")
		f := TraverseArray(func(a int) ReaderIOEither[testContext, error, string] {
			if a == 2 {
				return Left[testContext, string](expectedErr)
			}
			return Of[testContext, error](strconv.Itoa(a))
		})
		result := f([]int{1, 2, 3})
		assert.Equal(t, E.Left[[]string](expectedErr), result(ctx)())
	})
}

func TestTraverseArrayWithIndexDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("basic functionality", func(t *testing.T) {
		result := TraverseArrayWithIndex(func(i int, x int) ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](x + i)
		})([]int{1, 2, 3})

		assert.Equal(t, E.Right[error]([]int{1, 3, 5}), result(ctx)())
	})

	t.Run("empty array", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a string) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{})
		assert.Equal(t, E.Right[error]([]string{}), result(ctx)())
	})

	t.Run("fails at specific index", func(t *testing.T) {
		expectedErr := fmt.Errorf("error at index 1")
		f := TraverseArrayWithIndex(func(i int, a string) ReaderIOEither[testContext, error, string] {
			if i == 1 {
				return Left[testContext, string](expectedErr)
			}
			return Of[testContext, error](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{"a", "b", "c"})
		assert.Equal(t, E.Left[[]string](expectedErr), result(ctx)())
	})
}

func TestTraverseRecordDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("basic functionality", func(t *testing.T) {
		result := TraverseRecord[string](func(x int) ReaderIOEither[testContext, error, int] {
			return Of[testContext, error](x * 2)
		})(map[string]int{"a": 1, "b": 2})

		expected := map[string]int{"a": 2, "b": 4}
		assert.Equal(t, E.Right[error](expected), result(ctx)())
	})

	t.Run("empty record", func(t *testing.T) {
		f := TraverseRecord[string](func(a int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](strconv.Itoa(a))
		})
		result := f(map[string]int{})
		assert.Equal(t, E.Right[error](map[string]string{}), result(ctx)())
	})

	t.Run("one value fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("transformation error")
		f := TraverseRecord[string](func(a int) ReaderIOEither[testContext, error, string] {
			if a == 2 {
				return Left[testContext, string](expectedErr)
			}
			return Of[testContext, error](strconv.Itoa(a))
		})
		input := map[string]int{"a": 1, "b": 2, "c": 3}
		result := f(input)
		assert.Equal(t, E.Left[map[string]string](expectedErr), result(ctx)())
	})
}

func TestTraverseRecordWithIndexDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("basic functionality", func(t *testing.T) {
		result := TraverseRecordWithIndex(func(k string, x int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%s:%d", k, x))
		})(map[string]int{"a": 1, "b": 2})

		res := result(ctx)()
		assert.True(t, E.IsRight(res))
	})

	t.Run("empty record", func(t *testing.T) {
		f := TraverseRecordWithIndex(func(k string, a int) ReaderIOEither[testContext, error, string] {
			return Of[testContext, error](fmt.Sprintf("%s:%d", k, a))
		})
		result := f(map[string]int{})
		assert.Equal(t, E.Right[error](map[string]string{}), result(ctx)())
	})

	t.Run("fails for specific key", func(t *testing.T) {
		expectedErr := fmt.Errorf("error for key y")
		f := TraverseRecordWithIndex(func(k string, a int) ReaderIOEither[testContext, error, string] {
			if k == "y" {
				return Left[testContext, string](expectedErr)
			}
			return Of[testContext, error](fmt.Sprintf("%s:%d", k, a))
		})
		input := map[string]int{"x": 1, "y": 2}
		result := f(input)
		assert.Equal(t, E.Left[map[string]string](expectedErr), result(ctx)())
	})
}

func TestSequenceArrayDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("empty array", func(t *testing.T) {
		computations := []ReaderIOEither[testContext, error, int]{}
		result := SequenceArray(computations)
		assert.Equal(t, E.Right[error]([]int{}), result(ctx)())
	})

	t.Run("all successful", func(t *testing.T) {
		computations := []ReaderIOEither[testContext, error, int]{
			Of[testContext, error](1),
			Of[testContext, error](2),
			Of[testContext, error](3),
		}
		result := SequenceArray(computations)
		assert.Equal(t, E.Right[error]([]int{1, 2, 3}), result(ctx)())
	})

	t.Run("first computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("first computation error")
		computations := []ReaderIOEither[testContext, error, int]{
			Left[testContext, int](expectedErr),
			Of[testContext, error](2),
			Of[testContext, error](3),
		}
		result := SequenceArray(computations)
		assert.Equal(t, E.Left[[]int](expectedErr), result(ctx)())
	})

	t.Run("middle computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("middle computation error")
		computations := []ReaderIOEither[testContext, error, int]{
			Of[testContext, error](1),
			Left[testContext, int](expectedErr),
			Of[testContext, error](3),
		}
		result := SequenceArray(computations)
		assert.Equal(t, E.Left[[]int](expectedErr), result(ctx)())
	})
}

func TestSequenceRecordDetailed(t *testing.T) {
	ctx := testContext{value: 10}

	t.Run("basic functionality", func(t *testing.T) {
		result := SequenceRecord(map[string]ReaderIOEither[testContext, error, int]{
			"a": Of[testContext, error](1),
			"b": Of[testContext, error](2),
		})

		expected := map[string]int{"a": 1, "b": 2}
		assert.Equal(t, E.Right[error](expected), result(ctx)())
	})

	t.Run("empty record", func(t *testing.T) {
		computations := map[string]ReaderIOEither[testContext, error, int]{}
		result := SequenceRecord(computations)
		assert.Equal(t, E.Right[error](map[string]int{}), result(ctx)())
	})

	t.Run("one computation fails", func(t *testing.T) {
		expectedErr := fmt.Errorf("posts computation error")
		computations := map[string]ReaderIOEither[testContext, error, int]{
			"users":    Of[testContext, error](100),
			"posts":    Left[testContext, int](expectedErr),
			"comments": Of[testContext, error](200),
		}
		result := SequenceRecord(computations)
		assert.Equal(t, E.Left[map[string]int](expectedErr), result(ctx)())
	})
}

func TestSequenceT1(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT1(Of[testContext, error](42))
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestSequenceT3(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT3(
		Of[testContext, error](1),
		Of[testContext, error]("a"),
		Of[testContext, error](true),
	)
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestSequenceT4(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceT4(
		Of[testContext, error](1),
		Of[testContext, error]("a"),
		Of[testContext, error](true),
		Of[testContext, error](3.14),
	)
	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestWithLock(t *testing.T) {
	ctx := testContext{value: 10}
	unlocked := false

	result := F.Pipe1(
		Of[testContext, error](42),
		WithLock[testContext, error, int](func() context.CancelFunc {
			return func() { unlocked = true }
		}),
	)

	assert.Equal(t, E.Right[error](42), result(ctx)())
	assert.True(t, unlocked)
}

func TestMonadChainFirstLeft(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves original error", func(t *testing.T) {
		sideEffectCalled := false
		originalErr := errors.New("original error")
		result := MonadChainFirstLeft(
			Left[testContext, int](originalErr),
			func(e error) ReaderIOEither[testContext, error, int] {
				sideEffectCalled = true
				return Left[testContext, int](errors.New("new error")) // This error is ignored
			},
		)
		actualResult := result(ctx)()
		assert.True(t, sideEffectCalled)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var capturedError error
		originalErr := errors.New("validation failed")
		result := MonadChainFirstLeft(
			Left[testContext, int](originalErr),
			func(e error) ReaderIOEither[testContext, error, int] {
				capturedError = e
				return Right[testContext, error](999) // This Right value is ignored
			},
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, capturedError)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Right[testContext, error](42),
			func(e error) ReaderIOEither[testContext, error, int] {
				sideEffectCalled = true
				return Left[testContext, int](errors.New("should not be called"))
			},
		)
		assert.False(t, sideEffectCalled)
		assert.Equal(t, E.Right[error](42), result(ctx)())
	})

	// Test that side effects are executed but original error is always preserved
	t.Run("Side effects executed but original error preserved", func(t *testing.T) {
		effectCount := 0
		originalErr := errors.New("original error")
		result := MonadChainFirstLeft(
			Left[testContext, int](originalErr),
			func(e error) ReaderIOEither[testContext, error, int] {
				effectCount++
				// Try to return Right, but original Left should still be returned
				return Right[testContext, error](999)
			},
		)
		actualResult := result(ctx)()
		assert.Equal(t, 1, effectCount)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})
}

func TestChainFirstLeft(t *testing.T) {
	ctx := testContext{value: 10}

	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves error", func(t *testing.T) {
		var captured error
		originalErr := errors.New("test error")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOEither[testContext, error, int] {
			captured = e
			return Left[testContext, int](errors.New("ignored error"))
		})
		result := F.Pipe1(
			Left[testContext, int](originalErr),
			chainFn,
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, captured)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var captured error
		originalErr := errors.New("test error")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOEither[testContext, error, int] {
			captured = e
			return Right[testContext, error](42) // This Right is ignored
		})
		result := F.Pipe1(
			Left[testContext, int](originalErr),
			chainFn,
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, captured)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		called := false
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOEither[testContext, error, int] {
			called = true
			return Right[testContext, error](0)
		})
		result := F.Pipe1(
			Right[testContext, error](100),
			chainFn,
		)
		assert.False(t, called)
		assert.Equal(t, E.Right[error](100), result(ctx)())
	})

	// Test that original error is always preserved regardless of what f returns
	t.Run("Original error always preserved", func(t *testing.T) {
		originalErr := errors.New("original")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOEither[testContext, error, int] {
			// Try to return Right, but original Left should still be returned
			return Right[testContext, error](999)
		})

		result := F.Pipe1(
			Left[testContext, int](originalErr),
			chainFn,
		)
		assert.Equal(t, E.Left[int](originalErr), result(ctx)())
	})

	// Test logging with Left preservation
	t.Run("Logging with Left preservation", func(t *testing.T) {
		errorLog := []string{}
		originalErr := errors.New("step1")
		logError := ChainFirstLeft[string](func(e error) ReaderIOEither[testContext, error, string] {
			errorLog = append(errorLog, "Logged: "+e.Error())
			return Left[testContext, string](errors.New("log entry")) // This is ignored
		})

		result := F.Pipe2(
			Left[testContext, string](originalErr),
			logError,
			ChainLeft(func(e error) ReaderIOEither[testContext, error, string] {
				return Left[testContext, string](errors.New("step2"))
			}),
		)

		actualResult := result(ctx)()
		assert.Equal(t, []string{"Logged: step1"}, errorLog)
		assert.Equal(t, E.Left[string](errors.New("step2")), actualResult)
	})
}
