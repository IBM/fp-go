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

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
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
	result := MonadMap(Of[testContext, error](5), func(x int) int { return x * 2 })
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
		ChainFirst[testContext, error](func(x int) ReaderIOEither[testContext, error, string] {
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
		ChainFirstEitherK[testContext, error](func(x int) E.Either[error, string] {
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
		ChainIOEitherK[testContext, error](func(x int) IOE.IOEither[error, int] {
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
		ChainOptionK[testContext, int, int, error](func() error {
			return errors.New("none")
		})(func(x int) O.Option[int] {
			return O.Some(x * 2)
		}),
	)
	assert.Equal(t, E.Right[error](10), resultSome(ctx)())

	// Test with None
	resultNone := F.Pipe1(
		Of[testContext, error](5),
		ChainOptionK[testContext, int, int, error](func() error {
			return errors.New("none")
		})(func(x int) O.Option[int] {
			return O.None[int]()
		}),
	)
	assert.True(t, E.IsLeft(resultNone(ctx)()))
}

func TestMonadApSeq(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext, error](func(x int) int { return x * 2 })
	fa := Of[testContext, error](5)
	result := MonadApSeq(fab, fa)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestMonadApPar(t *testing.T) {
	ctx := testContext{value: 10}
	fab := Of[testContext, error](func(x int) int { return x * 2 })
	fa := Of[testContext, error](5)
	result := MonadApPar(fab, fa)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestChain(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		Chain[testContext, error](func(x int) ReaderIOEither[testContext, error, int] {
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
	resultTrue := FromPredicate[testContext, error](
		func(x int) bool { return x > 0 },
		func(x int) error { return errors.New("negative") },
	)(5)
	assert.Equal(t, E.Right[error](5), resultTrue(ctx)())

	// Test predicate false
	resultFalse := FromPredicate[testContext, error](
		func(x int) bool { return x > 0 },
		func(x int) error { return errors.New("negative") },
	)(-5)
	assert.True(t, E.IsLeft(resultFalse(ctx)()))
}

func TestFold(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := Fold[testContext, error, int, string](
		func(e error) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext]("error: " + e.Error())
		},
		func(x int) RIO.ReaderIO[testContext, string] {
			return RIO.Of[testContext](fmt.Sprintf("value: %d", x))
		},
	)(Of[testContext, error](42))
	assert.Equal(t, "value: 42", resultRight(ctx)())

	// Test Left case
	resultLeft := Fold[testContext, error, int, string](
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
	resultRight := GetOrElse[testContext, error](func(e error) RIO.ReaderIO[testContext, int] {
		return RIO.Of[testContext](0)
	})(Of[testContext, error](42))
	assert.Equal(t, 42, resultRight(ctx)())

	// Test Left case
	resultLeft := GetOrElse[testContext, error](func(e error) RIO.ReaderIO[testContext, int] {
		return RIO.Of[testContext](0)
	})(Left[testContext, int](errors.New("test")))
	assert.Equal(t, 0, resultLeft(ctx)())
}

func TestOrElse(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := OrElse[testContext, error, int, string](func(e error) ReaderIOEither[testContext, string, int] {
		return Left[testContext, int]("alternative")
	})(Of[testContext, error](42))
	assert.Equal(t, E.Right[string](42), resultRight(ctx)())

	// Test Left case
	resultLeft := OrElse[testContext, error, int, string](func(e error) ReaderIOEither[testContext, string, int] {
		return Of[testContext, string](99)
	})(Left[testContext, int](errors.New("test")))
	assert.Equal(t, E.Right[string](99), resultLeft(ctx)())
}

func TestMonadBiMap(t *testing.T) {
	ctx := testContext{value: 10}

	// Test Right case
	resultRight := MonadBiMap(
		Of[testContext, error](5),
		func(e error) string { return e.Error() },
		func(x int) string { return fmt.Sprintf("%d", x) },
	)
	assert.Equal(t, E.Right[string]("5"), resultRight(ctx)())

	// Test Left case
	resultLeft := MonadBiMap(
		Left[testContext, int](errors.New("test")),
		func(e error) string { return e.Error() },
		func(x int) string { return fmt.Sprintf("%d", x) },
	)
	assert.Equal(t, E.Left[string]("test"), resultLeft(ctx)())
}

func TestBiMap(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](5),
		BiMap[testContext, error, string](
			func(e error) string { return e.Error() },
			func(x int) string { return fmt.Sprintf("%d", x) },
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
		Alt[testContext, error](func() ReaderIOEither[testContext, error, int] {
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
	fab := Of[testContext, error](func(x int) int { return x * 2 })
	result := MonadFlap(fab, 5)
	assert.Equal(t, E.Right[error](10), result(ctx)())
}

func TestFlap(t *testing.T) {
	ctx := testContext{value: 10}
	result := F.Pipe1(
		Of[testContext, error](func(x int) int { return x * 2 }),
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

func TestMonadFromReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	result := MonadFromReaderIO[testContext, error](
		5,
		func(x int) RIO.ReaderIO[testContext, int] {
			return func(c testContext) io.IO[int] {
				return func() int { return x + c.value }
			}
		},
	)
	assert.Equal(t, E.Right[error](15), result(ctx)())
}

func TestFromReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	result := FromReaderIO[testContext, error](func(x int) RIO.ReaderIO[testContext, int] {
		return func(c testContext) io.IO[int] {
			return func() int { return x + c.value }
		}
	})(5)
	assert.Equal(t, E.Right[error](15), result(ctx)())
}

func TestRightReaderIO(t *testing.T) {
	ctx := testContext{value: 10}
	rio := func(c testContext) io.IO[int] {
		return func() int { return c.value * 2 }
	}
	result := RightReaderIO[testContext, error](rio)
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

	result := WithResource[string, testContext, error, int, int](
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

func TestTraverseArrayWithIndex(t *testing.T) {
	ctx := testContext{value: 10}
	result := TraverseArrayWithIndex[testContext, error](func(i int, x int) ReaderIOEither[testContext, error, int] {
		return Of[testContext, error](x + i)
	})([]int{1, 2, 3})

	assert.Equal(t, E.Right[error]([]int{1, 3, 5}), result(ctx)())
}

func TestTraverseRecord(t *testing.T) {
	ctx := testContext{value: 10}
	result := TraverseRecord[string, testContext, error](func(x int) ReaderIOEither[testContext, error, int] {
		return Of[testContext, error](x * 2)
	})(map[string]int{"a": 1, "b": 2})

	expected := map[string]int{"a": 2, "b": 4}
	assert.Equal(t, E.Right[error](expected), result(ctx)())
}

func TestTraverseRecordWithIndex(t *testing.T) {
	ctx := testContext{value: 10}
	result := TraverseRecordWithIndex[string, testContext, error](func(k string, x int) ReaderIOEither[testContext, error, string] {
		return Of[testContext, error](fmt.Sprintf("%s:%d", k, x))
	})(map[string]int{"a": 1, "b": 2})

	res := result(ctx)()
	assert.True(t, E.IsRight(res))
}

func TestSequenceRecord(t *testing.T) {
	ctx := testContext{value: 10}
	result := SequenceRecord[string, testContext, error](map[string]ReaderIOEither[testContext, error, int]{
		"a": Of[testContext, error](1),
		"b": Of[testContext, error](2),
	})

	expected := map[string]int{"a": 1, "b": 2}
	assert.Equal(t, E.Right[error](expected), result(ctx)())
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
