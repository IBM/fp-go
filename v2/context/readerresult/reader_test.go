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

package readerresult

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestMapTo(t *testing.T) {
	t.Run("executes original reader and returns constant value on success", func(t *testing.T) {
		executed := false
		originalReader := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		// Apply MapTo operator
		toDone := MapTo[int]("done")
		resultReader := toDone(originalReader)

		// Execute the resulting reader
		result := resultReader(t.Context())

		// Verify the constant value is returned
		assert.Equal(t, E.Of[error]("done"), result)
		// Verify the original reader WAS executed (side effect occurred)
		assert.True(t, executed, "original reader should be executed to allow side effects")
	})

	t.Run("executes reader in functional pipeline", func(t *testing.T) {
		executed := false
		step1 := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](100)
		}

		pipeline := F.Pipe1(
			step1,
			MapTo[int]("complete"),
		)

		result := pipeline(t.Context())

		assert.Equal(t, E.Of[error]("complete"), result)
		assert.True(t, executed, "original reader should be executed in pipeline")
	})

	t.Run("executes reader with side effects", func(t *testing.T) {
		sideEffectOccurred := false
		readerWithSideEffect := func(ctx context.Context) E.Either[error, int] {
			sideEffectOccurred = true
			return E.Of[error](42)
		}

		resultReader := MapTo[int](true)(readerWithSideEffect)
		result := resultReader(t.Context())

		assert.Equal(t, E.Of[error](true), result)
		assert.True(t, sideEffectOccurred, "side effect should occur")
	})

	t.Run("preserves errors from original reader", func(t *testing.T) {
		executed := false
		testErr := assert.AnError
		failingReader := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Left[int](testErr)
		}

		resultReader := MapTo[int]("done")(failingReader)
		result := resultReader(t.Context())

		assert.Equal(t, E.Left[string](testErr), result)
		assert.True(t, executed, "failing reader should still be executed")
	})
}

func TestMonadMapTo(t *testing.T) {
	t.Run("executes original reader and returns constant value on success", func(t *testing.T) {
		executed := false
		originalReader := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		// Apply MonadMapTo
		resultReader := MonadMapTo(originalReader, "done")

		// Execute the resulting reader
		result := resultReader(t.Context())

		// Verify the constant value is returned
		assert.Equal(t, E.Of[error]("done"), result)
		// Verify the original reader WAS executed (side effect occurred)
		assert.True(t, executed, "original reader should be executed to allow side effects")
	})

	t.Run("executes complex computation with side effects", func(t *testing.T) {
		computationExecuted := false
		complexReader := func(ctx context.Context) E.Either[error, string] {
			computationExecuted = true
			return E.Of[error]("complex result")
		}

		resultReader := MonadMapTo(complexReader, 42)
		result := resultReader(t.Context())

		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, computationExecuted, "complex computation should be executed")
	})

	t.Run("preserves errors from original reader", func(t *testing.T) {
		executed := false
		testErr := assert.AnError
		failingReader := func(ctx context.Context) E.Either[error, []string] {
			executed = true
			return E.Left[[]string](testErr)
		}

		resultReader := MonadMapTo(failingReader, 99)
		result := resultReader(t.Context())

		assert.Equal(t, E.Left[int](testErr), result)
		assert.True(t, executed, "failing reader should still be executed")
	})
}

func TestChainTo(t *testing.T) {
	t.Run("executes first reader then second reader on success", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		firstReader := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Of[error](42)
		}

		secondReader := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("result")
		}

		// Apply ChainTo operator
		thenSecond := ChainTo[int](secondReader)
		resultReader := thenSecond(firstReader)

		// Execute the resulting reader
		result := resultReader(t.Context())

		// Verify the second reader's result is returned
		assert.Equal(t, E.Of[error]("result"), result)
		// Verify both readers were executed
		assert.True(t, firstExecuted, "first reader should be executed")
		assert.True(t, secondExecuted, "second reader should be executed")
	})

	t.Run("executes both readers in functional pipeline", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		step1 := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Of[error](100)
		}

		step2 := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("complete")
		}

		pipeline := F.Pipe1(
			step1,
			ChainTo[int](step2),
		)

		result := pipeline(t.Context())

		assert.Equal(t, E.Of[error]("complete"), result)
		assert.True(t, firstExecuted, "first reader should be executed in pipeline")
		assert.True(t, secondExecuted, "second reader should be executed in pipeline")
	})

	t.Run("executes first reader with side effects", func(t *testing.T) {
		sideEffectOccurred := false
		readerWithSideEffect := func(ctx context.Context) E.Either[error, int] {
			sideEffectOccurred = true
			return E.Of[error](42)
		}

		secondReader := func(ctx context.Context) E.Either[error, bool] {
			return E.Of[error](true)
		}

		resultReader := ChainTo[int](secondReader)(readerWithSideEffect)
		result := resultReader(t.Context())

		assert.Equal(t, E.Of[error](true), result)
		assert.True(t, sideEffectOccurred, "side effect should occur in first reader")
	})

	t.Run("preserves error from first reader without executing second", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false
		testErr := assert.AnError

		failingReader := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Left[int](testErr)
		}

		secondReader := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("result")
		}

		resultReader := ChainTo[int](secondReader)(failingReader)
		result := resultReader(t.Context())

		assert.Equal(t, E.Left[string](testErr), result)
		assert.True(t, firstExecuted, "first reader should be executed")
		assert.False(t, secondExecuted, "second reader should not be executed on error")
	})
}

func TestMonadChainTo(t *testing.T) {
	t.Run("executes first reader then second reader on success", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		firstReader := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Of[error](42)
		}

		secondReader := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("result")
		}

		// Apply MonadChainTo
		resultReader := MonadChainTo(firstReader, secondReader)

		// Execute the resulting reader
		result := resultReader(t.Context())

		// Verify the second reader's result is returned
		assert.Equal(t, E.Of[error]("result"), result)
		// Verify both readers were executed
		assert.True(t, firstExecuted, "first reader should be executed")
		assert.True(t, secondExecuted, "second reader should be executed")
	})

	t.Run("executes complex first computation with side effects", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		complexFirstReader := func(ctx context.Context) E.Either[error, []int] {
			firstExecuted = true
			return E.Of[error]([]int{1, 2, 3})
		}

		secondReader := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("done")
		}

		resultReader := MonadChainTo(complexFirstReader, secondReader)
		result := resultReader(t.Context())

		assert.Equal(t, E.Of[error]("done"), result)
		assert.True(t, firstExecuted, "complex first computation should be executed")
		assert.True(t, secondExecuted, "second reader should be executed")
	})

	t.Run("preserves error from first reader without executing second", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false
		testErr := assert.AnError

		failingReader := func(ctx context.Context) E.Either[error, map[string]int] {
			firstExecuted = true
			return E.Left[map[string]int](testErr)
		}

		secondReader := func(ctx context.Context) E.Either[error, float64] {
			secondExecuted = true
			return E.Of[error](3.14)
		}

		resultReader := MonadChainTo(failingReader, secondReader)
		result := resultReader(t.Context())

		assert.Equal(t, E.Left[float64](testErr), result)
		assert.True(t, firstExecuted, "first reader should be executed")
		assert.False(t, secondExecuted, "second reader should not be executed on error")
	})
}

func TestOrElse(t *testing.T) {
	ctx := t.Context()

	// Test OrElse with Right - should pass through unchanged
	t.Run("Right value unchanged", func(t *testing.T) {
		rightValue := Of(42)
		recover := OrElse(func(err error) ReaderResult[int] {
			return Left[int](errors.New("should not be called"))
		})
		res := recover(rightValue)(ctx)
		assert.Equal(t, E.Of[error](42), res)
	})

	// Test OrElse with Left - should recover with fallback
	t.Run("Left value recovered", func(t *testing.T) {
		leftValue := Left[int](errors.New("not found"))
		recoverWithFallback := OrElse(func(err error) ReaderResult[int] {
			if err.Error() == "not found" {
				return func(ctx context.Context) E.Either[error, int] {
					return E.Of[error](99)
				}
			}
			return Left[int](err)
		})
		res := recoverWithFallback(leftValue)(ctx)
		assert.Equal(t, E.Of[error](99), res)
	})

	// Test OrElse with Left - should propagate other errors
	t.Run("Left value propagated", func(t *testing.T) {
		leftValue := Left[int](errors.New("fatal error"))
		recoverWithFallback := OrElse(func(err error) ReaderResult[int] {
			if err.Error() == "not found" {
				return Of(99)
			}
			return Left[int](err)
		})
		res := recoverWithFallback(leftValue)(ctx)
		assert.True(t, E.IsLeft(res))
		val, err := E.UnwrapError(res)
		assert.Equal(t, 0, val)
		assert.Equal(t, "fatal error", err.Error())
	})

	// Test OrElse with context-aware recovery
	t.Run("Context-aware recovery", func(t *testing.T) {
		type ctxKey string
		ctxWithValue := context.WithValue(ctx, ctxKey("fallback"), 123)

		leftValue := Left[int](errors.New("use fallback"))
		ctxRecover := OrElse(func(err error) ReaderResult[int] {
			if err.Error() == "use fallback" {
				return func(ctx context.Context) E.Either[error, int] {
					if val := ctx.Value(ctxKey("fallback")); val != nil {
						return E.Of[error](val.(int))
					}
					return E.Left[int](errors.New("no fallback"))
				}
			}
			return Left[int](err)
		})
		res := ctxRecover(leftValue)(ctxWithValue)
		assert.Equal(t, E.Of[error](123), res)
	})
}

// TestFromIO tests the FromIO function
func TestFromIO(t *testing.T) {
	t.Run("lifts IO computation into ReaderResult", func(t *testing.T) {
		ioOp := func() int { return 42 }
		rr := FromIO(ioOp)
		result := rr(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("executes IO side effects", func(t *testing.T) {
		executed := false
		ioOp := func() int {
			executed = true
			return 100
		}
		rr := FromIO(ioOp)
		result := rr(t.Context())
		assert.True(t, executed, "IO operation should be executed")
		assert.Equal(t, E.Of[error](100), result)
	})
}

// TestFromIOResult tests the FromIOResult function
func TestFromIOResult(t *testing.T) {
	t.Run("lifts IOResult into ReaderResult on success", func(t *testing.T) {
		ioResult := func() E.Either[error, int] {
			return E.Of[error](42)
		}
		rr := FromIOResult(ioResult)
		result := rr(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("lifts IOResult into ReaderResult on error", func(t *testing.T) {
		testErr := errors.New("io error")
		ioResult := func() E.Either[error, int] {
			return E.Left[int](testErr)
		}
		rr := FromIOResult(ioResult)
		result := rr(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
	})
}

// TestFromReader tests the FromReader function
func TestFromReader(t *testing.T) {
	t.Run("lifts Reader into ReaderResult", func(t *testing.T) {
		reader := func(ctx context.Context) int {
			return 42
		}
		rr := FromReader(reader)
		result := rr(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("Reader can access context", func(t *testing.T) {
		type ctxKey string
		ctx := context.WithValue(t.Context(), ctxKey("key"), "value")
		reader := func(ctx context.Context) string {
			return ctx.Value(ctxKey("key")).(string)
		}
		rr := FromReader(reader)
		result := rr(ctx)
		assert.Equal(t, E.Of[error]("value"), result)
	})
}

// TestFromEither tests the FromEither function
func TestFromEither(t *testing.T) {
	t.Run("lifts Right Either into ReaderResult", func(t *testing.T) {
		either := E.Of[error](42)
		rr := FromEither(either)
		result := rr(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("lifts Left Either into ReaderResult", func(t *testing.T) {
		testErr := errors.New("test error")
		either := E.Left[int](testErr)
		rr := FromEither(either)
		result := rr(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
	})
}

// TestLeftRight tests the Left and Right functions
func TestLeftRight(t *testing.T) {
	t.Run("Left creates error ReaderResult", func(t *testing.T) {
		testErr := errors.New("test error")
		rr := Left[int](testErr)
		result := rr(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("Right creates success ReaderResult", func(t *testing.T) {
		rr := Right(42)
		result := rr(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})
}

// TestMonadMapAndMap tests MonadMap and Map functions
func TestMonadMapAndMap(t *testing.T) {
	t.Run("MonadMap transforms success value", func(t *testing.T) {
		rr := Of(42)
		mapped := MonadMap(rr, func(x int) string {
			return F.Pipe1(x, func(n int) string { return "value: " + F.Pipe1(n, func(i int) string { return string(rune(i + 48)) }) })
		})
		result := mapped(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("Map creates operator that transforms value", func(t *testing.T) {
		toString := Map(func(x int) string {
			return "value"
		})
		rr := Of(42)
		result := toString(rr)(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("Map preserves errors", func(t *testing.T) {
		testErr := errors.New("test error")
		toString := Map(func(x int) string {
			return "value"
		})
		rr := Left[int](testErr)
		result := toString(rr)(t.Context())
		assert.Equal(t, E.Left[string](testErr), result)
	})
}

// TestMonadChainAndChain tests MonadChain and Chain functions
func TestMonadChainAndChain(t *testing.T) {
	t.Run("MonadChain sequences computations", func(t *testing.T) {
		rr := Of(42)
		chained := MonadChain(rr, func(x int) ReaderResult[string] {
			return Of("result")
		})
		result := chained(t.Context())
		assert.Equal(t, E.Of[error]("result"), result)
	})

	t.Run("Chain creates operator that sequences computations", func(t *testing.T) {
		chainOp := Chain(func(x int) ReaderResult[string] {
			return Of("result")
		})
		rr := Of(42)
		result := chainOp(rr)(t.Context())
		assert.Equal(t, E.Of[error]("result"), result)
	})

	t.Run("Chain short-circuits on error", func(t *testing.T) {
		executed := false
		testErr := errors.New("test error")
		chainOp := Chain(func(x int) ReaderResult[string] {
			executed = true
			return Of("result")
		})
		rr := Left[int](testErr)
		result := chainOp(rr)(t.Context())
		assert.False(t, executed, "Chain should not execute on error")
		assert.Equal(t, E.Left[string](testErr), result)
	})
}

// TestAsk tests the Ask function
func TestAsk(t *testing.T) {
	t.Run("Ask returns the context", func(t *testing.T) {
		ctx := t.Context()
		rr := Ask()
		result := rr(ctx)
		assert.True(t, E.IsRight(result))
	})

	t.Run("Ask can be used in chain to access context", func(t *testing.T) {
		type ctxKey string
		ctx := context.WithValue(t.Context(), ctxKey("key"), "value")
		pipeline := F.Pipe1(
			Ask(),
			Chain(func(c context.Context) ReaderResult[string] {
				val := c.Value(ctxKey("key"))
				if val != nil {
					return Of(val.(string))
				}
				return Left[string](errors.New("key not found"))
			}),
		)
		result := pipeline(ctx)
		assert.Equal(t, E.Of[error]("value"), result)
	})
}

// TestMonadChainEitherK tests MonadChainEitherK and ChainEitherK
func TestMonadChainEitherK(t *testing.T) {
	t.Run("MonadChainEitherK sequences with Either function", func(t *testing.T) {
		rr := Of(42)
		chained := MonadChainEitherK(rr, func(x int) E.Either[error, string] {
			if x > 0 {
				return E.Of[error]("positive")
			}
			return E.Left[string](errors.New("not positive"))
		})
		result := chained(t.Context())
		assert.Equal(t, E.Of[error]("positive"), result)
	})

	t.Run("ChainEitherK creates operator", func(t *testing.T) {
		validate := ChainEitherK(func(x int) E.Either[error, int] {
			if x > 0 {
				return E.Of[error](x)
			}
			return E.Left[int](errors.New("must be positive"))
		})
		result := validate(Of(42))(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})
}

// TestMonadFlap tests MonadFlap and Flap
func TestMonadFlap(t *testing.T) {
	t.Run("MonadFlap applies value to function", func(t *testing.T) {
		fabRR := Of(func(x int) string {
			return "value"
		})
		result := MonadFlap(fabRR, 42)(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("Flap creates operator", func(t *testing.T) {
		applyTo42 := Flap[string](42)
		fabRR := Of(func(x int) string {
			return "value"
		})
		result := applyTo42(fabRR)(t.Context())
		assert.True(t, E.IsRight(result))
	})
}

// TestRead functions
func TestReadFunctions(t *testing.T) {
	t.Run("Read executes ReaderResult with context", func(t *testing.T) {
		rr := Of(42)
		ctx := t.Context()
		runWithCtx := Read[int](ctx)
		result := runWithCtx(rr)
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("ReadEither executes with Result context on success", func(t *testing.T) {
		rr := Of(42)
		ctxResult := E.Of[error](t.Context())
		runWithCtxResult := ReadEither[int](ctxResult)
		result := runWithCtxResult(rr)
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("ReadEither returns error when context Result is error", func(t *testing.T) {
		rr := Of(42)
		testErr := errors.New("context error")
		ctxResult := E.Left[context.Context](testErr)
		runWithCtxResult := ReadEither[int](ctxResult)
		result := runWithCtxResult(rr)
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("ReadResult is alias for ReadEither", func(t *testing.T) {
		rr := Of(42)
		ctxResult := E.Of[error](t.Context())
		runWithCtxResult := ReadResult[int](ctxResult)
		result := runWithCtxResult(rr)
		assert.Equal(t, E.Of[error](42), result)
	})
}

// TestMonadChainFirst tests MonadChainFirst and ChainFirst
func TestMonadChainFirst(t *testing.T) {
	t.Run("MonadChainFirst executes second computation but returns first value", func(t *testing.T) {
		secondExecuted := false
		rr := Of(42)
		withSideEffect := MonadChainFirst(rr, func(x int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				secondExecuted = true
				return E.Of[error]("logged")
			}
		})
		result := withSideEffect(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, secondExecuted, "second computation should execute")
	})

	t.Run("ChainFirst creates operator", func(t *testing.T) {
		secondExecuted := false
		logValue := ChainFirst(func(x int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				secondExecuted = true
				return E.Of[error]("logged")
			}
		})
		result := logValue(Of(42))(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, secondExecuted)
	})
}

// TestChainIOK tests ChainIOK and MonadChainIOK
func TestChainIOK(t *testing.T) {
	t.Run("MonadChainIOK sequences with IO computation", func(t *testing.T) {
		ioExecuted := false
		rr := Of(42)
		withIO := MonadChainIOK(rr, func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "done"
			}
		})
		result := withIO(t.Context())
		assert.Equal(t, E.Of[error]("done"), result)
		assert.True(t, ioExecuted)
	})

	t.Run("ChainIOK creates operator", func(t *testing.T) {
		ioExecuted := false
		logIO := ChainIOK(func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		result := logIO(Of(42))(t.Context())
		assert.Equal(t, E.Of[error]("logged"), result)
		assert.True(t, ioExecuted)
	})
}

// TestChainFirstIOK tests ChainFirstIOK, MonadChainFirstIOK, and TapIOK
func TestChainFirstIOK(t *testing.T) {
	t.Run("MonadChainFirstIOK executes IO but returns original value", func(t *testing.T) {
		ioExecuted := false
		rr := Of(42)
		withLog := MonadChainFirstIOK(rr, func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		result := withLog(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, ioExecuted)
	})

	t.Run("ChainFirstIOK creates operator", func(t *testing.T) {
		ioExecuted := false
		logIO := ChainFirstIOK(func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		result := logIO(Of(42))(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, ioExecuted)
	})

	t.Run("TapIOK is alias for ChainFirstIOK", func(t *testing.T) {
		ioExecuted := false
		tapLog := TapIOK(func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		result := tapLog(Of(42))(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, ioExecuted)
	})

	t.Run("MonadTapIOK is alias for MonadChainFirstIOK", func(t *testing.T) {
		ioExecuted := false
		rr := Of(42)
		withLog := MonadTapIOK(rr, func(x int) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		result := withLog(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.True(t, ioExecuted)
	})
}

// TestChainIOEitherK tests ChainIOEitherK and ChainIOResultK
func TestChainIOEitherK(t *testing.T) {
	t.Run("ChainIOEitherK sequences with IOResult on success", func(t *testing.T) {
		ioResultOp := ChainIOEitherK(func(x int) func() E.Either[error, string] {
			return func() E.Either[error, string] {
				if x > 0 {
					return E.Of[error]("positive")
				}
				return E.Left[string](errors.New("not positive"))
			}
		})
		result := ioResultOp(Of(42))(t.Context())
		assert.Equal(t, E.Of[error]("positive"), result)
	})

	t.Run("ChainIOEitherK propagates IOResult error", func(t *testing.T) {
		testErr := errors.New("io error")
		ioResultOp := ChainIOEitherK(func(x int) func() E.Either[error, string] {
			return func() E.Either[error, string] {
				return E.Left[string](testErr)
			}
		})
		result := ioResultOp(Of(42))(t.Context())
		assert.Equal(t, E.Left[string](testErr), result)
	})

	t.Run("ChainIOResultK is alias for ChainIOEitherK", func(t *testing.T) {
		ioResultOp := ChainIOResultK(func(x int) func() E.Either[error, string] {
			return func() E.Either[error, string] {
				return E.Of[error]("value")
			}
		})
		result := ioResultOp(Of(42))(t.Context())
		assert.Equal(t, E.Of[error]("value"), result)
	})
}

// TestReadIO tests ReadIO, ReadIOEither, and ReadIOResult
func TestReadIO(t *testing.T) {
	t.Run("ReadIO executes with IO context", func(t *testing.T) {
		getCtx := func() context.Context { return t.Context() }
		rr := Of(42)
		runWithIO := ReadIO[int](getCtx)
		ioResult := runWithIO(rr)
		result := ioResult()
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("ReadIOEither executes with IOResult context on success", func(t *testing.T) {
		getCtx := func() E.Either[error, context.Context] {
			return E.Of[error](t.Context())
		}
		rr := Of(42)
		runWithIOResult := ReadIOEither[int](getCtx)
		ioResult := runWithIOResult(rr)
		result := ioResult()
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("ReadIOEither returns error when IOResult context is error", func(t *testing.T) {
		testErr := errors.New("context error")
		getCtx := func() E.Either[error, context.Context] {
			return E.Left[context.Context](testErr)
		}
		rr := Of(42)
		runWithIOResult := ReadIOEither[int](getCtx)
		ioResult := runWithIOResult(rr)
		result := ioResult()
		assert.Equal(t, E.Left[int](testErr), result)
	})

	t.Run("ReadIOResult is alias for ReadIOEither", func(t *testing.T) {
		getCtx := func() E.Either[error, context.Context] {
			return E.Of[error](t.Context())
		}
		rr := Of(42)
		runWithIOResult := ReadIOResult[int](getCtx)
		ioResult := runWithIOResult(rr)
		result := ioResult()
		assert.Equal(t, E.Of[error](42), result)
	})
}

// TestChainFirstLeft tests ChainFirstLeft, ChainFirstLeftIOK, and TapLeftIOK
func TestChainFirstLeft(t *testing.T) {
	t.Run("ChainFirstLeft executes on error but preserves it", func(t *testing.T) {
		errorHandled := false
		testErr := errors.New("test error")
		logError := ChainFirstLeft[int](func(err error) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				errorHandled = true
				return E.Of[error]("logged")
			}
		})
		rr := Left[int](testErr)
		result := logError(rr)(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
		assert.True(t, errorHandled, "error handler should execute")
	})

	t.Run("ChainFirstLeft does not execute on success", func(t *testing.T) {
		errorHandled := false
		logError := ChainFirstLeft[int](func(err error) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				errorHandled = true
				return E.Of[error]("logged")
			}
		})
		rr := Of(42)
		result := logError(rr)(t.Context())
		assert.Equal(t, E.Of[error](42), result)
		assert.False(t, errorHandled, "error handler should not execute on success")
	})

	t.Run("ChainFirstLeftIOK executes IO on error", func(t *testing.T) {
		ioExecuted := false
		testErr := errors.New("test error")
		logErrorIO := ChainFirstLeftIOK[int](func(err error) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		rr := Left[int](testErr)
		result := logErrorIO(rr)(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
		assert.True(t, ioExecuted)
	})

	t.Run("TapLeftIOK is alias for ChainFirstLeftIOK", func(t *testing.T) {
		ioExecuted := false
		testErr := errors.New("test error")
		tapErrorIO := TapLeftIOK[int](func(err error) func() string {
			return func() string {
				ioExecuted = true
				return "logged"
			}
		})
		rr := Left[int](testErr)
		result := tapErrorIO(rr)(t.Context())
		assert.Equal(t, E.Left[int](testErr), result)
		assert.True(t, ioExecuted)
	})
}

// TestFromPredicate tests the FromPredicate function
func TestFromPredicate(t *testing.T) {
	t.Run("FromPredicate returns Right when predicate is true", func(t *testing.T) {
		isPositive := FromPredicate(
			func(x int) bool { return x > 0 },
			func(x int) error { return errors.New("not positive") },
		)
		result := isPositive(42)(t.Context())
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("FromPredicate returns Left when predicate is false", func(t *testing.T) {
		isPositive := FromPredicate(
			func(x int) bool { return x > 0 },
			func(x int) error { return errors.New("not positive") },
		)
		result := isPositive(-1)(t.Context())
		assert.True(t, E.IsLeft(result))
	})
}

// TestMonadAp tests MonadAp and Ap
func TestMonadAp(t *testing.T) {
	t.Run("MonadAp applies function to value", func(t *testing.T) {
		fabRR := Of(func(x int) string {
			return "value"
		})
		faRR := Of(42)
		result := MonadAp(fabRR, faRR)(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("Ap creates function that applies", func(t *testing.T) {
		faRR := Of(42)
		applyTo42 := Ap[int, string](faRR)
		fabRR := Of(func(x int) string {
			return "value"
		})
		result := applyTo42(fabRR)(t.Context())
		assert.True(t, E.IsRight(result))
	})
}

// TestChainOptionK tests the ChainOptionK function
func TestChainOptionK(t *testing.T) {
	t.Run("ChainOptionK returns Right when Option is Some", func(t *testing.T) {
		chainOpt := ChainOptionK[int, string](func() error {
			return errors.New("value not found")
		})
		optKleisli := func(x int) option.Option[string] {
			if x > 0 {
				return option.Some("value")
			}
			return option.None[string]()
		}
		operator := chainOpt(optKleisli)
		result := operator(Of(42))(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("ChainOptionK returns Left when Option is None", func(t *testing.T) {
		chainOpt := ChainOptionK[int, string](func() error {
			return errors.New("value not found")
		})
		optKleisli := func(x int) option.Option[string] {
			return option.None[string]()
		}
		operator := chainOpt(optKleisli)
		result := operator(Of(42))(t.Context())
		assert.True(t, E.IsLeft(result))
	})
}
