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

package ioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, E.Of[error](2), F.Pipe1(
		Of[error](1),
		Map[error](utils.Double),
	)())

}

func TestChainEitherK(t *testing.T) {
	f := ChainEitherK(func(n int) E.Either[string, int] {
		if n > 0 {
			return E.Of[string](n)
		}
		return E.Left[int]("a")

	})
	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("a"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})

	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestFromOption(t *testing.T) {
	f := FromOption[int](F.Constant("a"))
	assert.Equal(t, E.Right[string](1), f(O.Some(1))())
	assert.Equal(t, E.Left[int]("a"), f(O.None[int]())())
}

func TestChainIOK(t *testing.T) {
	f := ChainIOK[string](func(n int) IO[string] {
		return func() string {
			return fmt.Sprintf("%d", n)
		}
	})

	assert.Equal(t, E.Right[string]("1"), f(Right[string](1))())
	assert.Equal(t, E.Left[string]("b"), f(Left[int]("b"))())
}

func TestChainWithIO(t *testing.T) {

	r := F.Pipe1(
		Of[error]("test"),
		// sad, we need the generics version ...
		io.Map(E.IsRight[error, string]),
	)

	assert.True(t, r())
}

func TestChainFirst(t *testing.T) {
	f := func(a string) IOEither[string, int] {
		if len(a) > 2 {
			return Of[string](len(a))
		}
		return Left[int]("foo")
	}
	good := Of[string]("foo")
	bad := Of[string]("a")
	ch := ChainFirst(f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
	assert.Equal(t, E.Left[string]("foo"), F.Pipe1(bad, ch)())
}

func TestChainFirstIOK(t *testing.T) {
	f := func(a string) IO[int] {
		return io.Of(len(a))
	}
	good := Of[string]("foo")
	ch := ChainFirstIOK[string](f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApFirst[string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("a"), x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApSecond[string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("b"), x())
}

func TestMonadChainLeft(t *testing.T) {
	// Test with Left value - should apply the function
	t.Run("Left value applies function", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int]("error1"),
			func(e string) IOEither[string, int] {
				return Left[int]("transformed: " + e)
			},
		)
		assert.Equal(t, E.Left[int]("transformed: error1"), result())
	})

	// Test with Left value - function returns Right (error recovery)
	t.Run("Left value recovers to Right", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int]("recoverable"),
			func(e string) IOEither[string, int] {
				if e == "recoverable" {
					return Right[string](42)
				}
				return Left[int](e)
			},
		)
		assert.Equal(t, E.Right[string](42), result())
	})

	// Test with Right value - should pass through unchanged
	t.Run("Right value passes through", func(t *testing.T) {
		result := MonadChainLeft(
			Right[string](100),
			func(e string) IOEither[string, int] {
				return Left[int]("should not be called")
			},
		)
		assert.Equal(t, E.Right[string](100), result())
	})

	// Test error type transformation
	t.Run("Error type transformation", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int]("404"),
			func(e string) IOEither[int, int] {
				return Left[int](404)
			},
		)
		assert.Equal(t, E.Left[int](404), result())
	})
}

func TestChainLeft(t *testing.T) {
	// Test with Left value - should apply the function
	t.Run("Left value applies function", func(t *testing.T) {
		chainFn := ChainLeft(func(e string) IOEither[string, int] {
			return Left[int]("chained: " + e)
		})
		result := F.Pipe1(
			Left[int]("original"),
			chainFn,
		)
		assert.Equal(t, E.Left[int]("chained: original"), result())
	})

	// Test with Left value - function returns Right (error recovery)
	t.Run("Left value recovers to Right", func(t *testing.T) {
		chainFn := ChainLeft(func(e string) IOEither[string, int] {
			if e == "network error" {
				return Right[string](0) // default value
			}
			return Left[int](e)
		})
		result := F.Pipe1(
			Left[int]("network error"),
			chainFn,
		)
		assert.Equal(t, E.Right[string](0), result())
	})

	// Test with Right value - should pass through unchanged
	t.Run("Right value passes through", func(t *testing.T) {
		chainFn := ChainLeft(func(e string) IOEither[string, int] {
			return Left[int]("should not be called")
		})
		result := F.Pipe1(
			Right[string](42),
			chainFn,
		)
		assert.Equal(t, E.Right[string](42), result())
	})

	// Test composition with other operations
	t.Run("Composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			Left[int]("error"),
			ChainLeft(func(e string) IOEither[string, int] {
				return Left[int]("handled: " + e)
			}),
			Map[string](utils.Double),
		)
		assert.Equal(t, E.Left[int]("handled: error"), result())
	})
}

func TestMonadChainFirstLeft(t *testing.T) {
	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves original error", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Left[int]("original error"),
			func(e string) IOEither[string, int] {
				sideEffectCalled = true
				return Left[int]("new error") // This error is ignored, original is returned
			},
		)
		actualResult := result()
		assert.True(t, sideEffectCalled)
		assert.Equal(t, E.Left[int]("original error"), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var capturedError string
		result := MonadChainFirstLeft(
			Left[int]("validation failed"),
			func(e string) IOEither[string, int] {
				capturedError = e
				return Right[string](999) // This Right value is ignored, original Left is returned
			},
		)
		actualResult := result()
		assert.Equal(t, "validation failed", capturedError)
		assert.Equal(t, E.Left[int]("validation failed"), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Right[string](42),
			func(e string) IOEither[string, int] {
				sideEffectCalled = true
				return Left[int]("should not be called")
			},
		)
		assert.False(t, sideEffectCalled)
		assert.Equal(t, E.Right[string](42), result())
	})

	// Test that side effects are executed but original error is always preserved
	t.Run("Side effects executed but original error preserved", func(t *testing.T) {
		effectCount := 0
		result := MonadChainFirstLeft(
			Left[int]("original error"),
			func(e string) IOEither[string, int] {
				effectCount++
				// Try to return Right, but original Left should still be returned
				return Right[string](999)
			},
		)
		actualResult := result()
		assert.Equal(t, 1, effectCount)
		assert.Equal(t, E.Left[int]("original error"), actualResult)
	})
}

func TestChainFirstLeft(t *testing.T) {
	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves error", func(t *testing.T) {
		var captured string
		chainFn := ChainFirstLeft[int](func(e string) IOEither[string, int] {
			captured = e
			return Left[int]("ignored error")
		})
		result := F.Pipe1(
			Left[int]("test error"),
			chainFn,
		)
		actualResult := result()
		assert.Equal(t, "test error", captured)
		assert.Equal(t, E.Left[int]("test error"), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var captured string
		chainFn := ChainFirstLeft[int](func(e string) IOEither[string, int] {
			captured = e
			return Right[string](42) // This Right is ignored, original Left is returned
		})
		result := F.Pipe1(
			Left[int]("test error"),
			chainFn,
		)
		actualResult := result()
		assert.Equal(t, "test error", captured)
		assert.Equal(t, E.Left[int]("test error"), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		called := false
		chainFn := ChainFirstLeft[int](func(e string) IOEither[string, int] {
			called = true
			return Right[string](0)
		})
		result := F.Pipe1(
			Right[string](100),
			chainFn,
		)
		assert.False(t, called)
		assert.Equal(t, E.Right[string](100), result())
	})

	// Test that original error is always preserved regardless of what f returns
	t.Run("Original error always preserved", func(t *testing.T) {
		chainFn := ChainFirstLeft[int](func(e string) IOEither[string, int] {
			// Try to return Right, but original Left should still be returned
			return Right[string](999)
		})

		result := F.Pipe1(
			Left[int]("original"),
			chainFn,
		)
		assert.Equal(t, E.Left[int]("original"), result())
	})

	// Test with IO side effects - original Left is always preserved
	t.Run("IO side effects with Left preservation", func(t *testing.T) {
		effectCount := 0
		chainFn := ChainFirstLeft[int](func(e string) IOEither[string, int] {
			return FromIO[string](func() int {
				effectCount++
				return 0
			})
		})

		// Even though FromIO wraps in Right, the original Left is preserved
		result := F.Pipe1(
			Left[int]("error"),
			chainFn,
		)

		assert.Equal(t, E.Left[int]("error"), result())
		assert.Equal(t, 1, effectCount)
	})

	// Test logging with Left preservation
	t.Run("Logging with Left preservation", func(t *testing.T) {
		errorLog := []string{}
		logError := ChainFirstLeft[string](func(e string) IOEither[string, string] {
			errorLog = append(errorLog, "Logged: "+e)
			return Left[string]("log entry") // This is ignored, original is preserved
		})

		result := F.Pipe2(
			Left[string]("step1"),
			logError,
			ChainLeft(func(e string) IOEither[string, string] {
				return Left[string]("step2")
			}),
		)

		actualResult := result()
		assert.Equal(t, []string{"Logged: step1"}, errorLog)
		assert.Equal(t, E.Left[string]("step2"), actualResult)
	})
}

func TestOrElse(t *testing.T) {
	// Test that OrElse recovers from a Left
	t.Run("OrElse recovers from Left", func(t *testing.T) {
		recover := OrElse(func(err string) IOEither[string, int] {
			return Right[string](42)
		})

		result := F.Pipe1(
			Left[int]("error"),
			recover,
		)
		assert.Equal(t, E.Right[string](42), result())
	})

	// When input is Right, should pass through unchanged
	t.Run("OrElse passes through Right unchanged", func(t *testing.T) {
		recover := OrElse(func(err string) IOEither[string, int] {
			return Right[string](42)
		})

		result := F.Pipe1(
			Right[string](100),
			recover,
		)
		assert.Equal(t, E.Right[string](100), result())
	})

	// Test that OrElse can also return a Left (propagate different error)
	t.Run("OrElse can propagate different error", func(t *testing.T) {
		recoverOrFail := OrElse(func(err string) IOEither[string, int] {
			if err == "recoverable" {
				return Right[string](0)
			}
			return Left[int]("unrecoverable: " + err)
		})

		recoverable := F.Pipe1(
			Left[int]("recoverable"),
			recoverOrFail,
		)
		assert.Equal(t, E.Right[string](0), recoverable())

		unrecoverable := F.Pipe1(
			Left[int]("fatal"),
			recoverOrFail,
		)
		assert.Equal(t, E.Left[int]("unrecoverable: fatal"), unrecoverable())
	})

	// Test composition with other operations
	t.Run("OrElse composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			Left[int]("error"),
			OrElse(func(err string) IOEither[string, int] {
				return Right[string](21)
			}),
			Map[string](utils.Double),
		)
		assert.Equal(t, E.Right[string](42), result())
	})
}
