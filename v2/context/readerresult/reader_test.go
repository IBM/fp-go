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
		result := resultReader(context.Background())

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

		result := pipeline(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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

		result := pipeline(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

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
		result := resultReader(context.Background())

		assert.Equal(t, E.Left[float64](testErr), result)
		assert.True(t, firstExecuted, "first reader should be executed")
		assert.False(t, secondExecuted, "second reader should not be executed on error")
	})
}

func TestOrElse(t *testing.T) {
	ctx := context.Background()

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
