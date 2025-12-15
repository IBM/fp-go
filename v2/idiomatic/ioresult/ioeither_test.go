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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	result, err := F.Pipe1(
		Of(1),
		Map(utils.Double),
	)()
	assert.NoError(t, err)
	assert.Equal(t, 2, result)

}

func TestChainEitherK(t *testing.T) {
	f := ChainResultK(func(n int) (int, error) {
		if n > 0 {
			return n, nil
		}
		return 0, errors.New("a")
	})
	result1, err1 := f(Right(1))()
	assert.NoError(t, err1)
	assert.Equal(t, 1, result1)

	_, err2 := f(Right(-1))()
	assert.Error(t, err2)
	assert.Equal(t, "a", err2.Error())

	_, err3 := f(Left[int](fmt.Errorf("b")))()
	assert.Error(t, err3)
	assert.Equal(t, "b", err3.Error())
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](func() error { return fmt.Errorf("a") })(func(n int) (int, bool) {
		if n > 0 {
			return n, true
		}
		return 0, false
	})

	result1, err1 := f(Right(1))()
	assert.NoError(t, err1)
	assert.Equal(t, 1, result1)

	_, err2 := f(Right(-1))()
	assert.Error(t, err2)
	assert.Equal(t, "a", err2.Error())

	_, err3 := f(Left[int](fmt.Errorf("b")))()
	assert.Error(t, err3)
	assert.Equal(t, "b", err3.Error())
}

func TestFromOption(t *testing.T) {
	f := FromOption[int](func() error { return errors.New("a") })

	result1, err1 := f(1, true)()
	assert.NoError(t, err1)
	assert.Equal(t, 1, result1)

	_, err2 := f(0, false)()
	assert.Error(t, err2)
	assert.Equal(t, "a", err2.Error())
}

func TestChainIOK(t *testing.T) {
	f := ChainIOK(func(n int) IO[string] {
		return func() string {
			return fmt.Sprintf("%d", n)
		}
	})

	result1, err1 := f(Right(1))()
	assert.NoError(t, err1)
	assert.Equal(t, "1", result1)

	_, err2 := f(Left[int](errors.New("b")))()
	assert.Error(t, err2)
	assert.Equal(t, "b", err2.Error())
}

func TestChainWithIO(t *testing.T) {

	r := F.Pipe1(
		Of("test"),
		ChainIOK(func(s string) IO[bool] {
			return func() bool {
				return S.IsNonEmpty(s)
			}
		}),
	)

	result, err := r()
	assert.NoError(t, err)
	assert.True(t, result)
}

func TestChainFirst(t *testing.T) {
	f := func(a string) IOResult[int] {
		if len(a) > 2 {
			return Of(len(a))
		}
		return Left[int](errors.New("foo"))
	}
	good := Of("foo")
	bad := Of("a")
	ch := ChainFirst(f)

	result1, err1 := F.Pipe1(good, ch)()
	assert.NoError(t, err1)
	assert.Equal(t, "foo", result1)

	_, err2 := F.Pipe1(bad, ch)()
	assert.Error(t, err2)
	assert.Equal(t, "foo", err2.Error())
}

func TestChainFirstIOK(t *testing.T) {
	f := func(a string) IO[int] {
		return io.Of(len(a))
	}
	good := Of("foo")
	ch := ChainFirstIOK(f)

	result, err := F.Pipe1(good, ch)()
	assert.NoError(t, err)
	assert.Equal(t, "foo", result)
}

func TestMonadChainLeft(t *testing.T) {
	// Test with Left value - should apply the function
	t.Run("Left value applies function", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int](errors.New("error1")),
			func(e error) IOResult[int] {
				return Left[int](errors.New("transformed: " + e.Error()))
			},
		)
		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "transformed: error1", err.Error())
	})

	// Test with Left value - function returns Right (error recovery)
	t.Run("Left value recovers to Right", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int](errors.New("recoverable")),
			func(e error) IOResult[int] {
				if e.Error() == "recoverable" {
					return Right(42)
				}
				return Left[int](e)
			},
		)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	// Test with Right value - should pass through unchanged
	t.Run("Right value passes through", func(t *testing.T) {
		result := MonadChainLeft(
			Right(100),
			func(e error) IOResult[int] {
				return Left[int](errors.New("should not be called"))
			},
		)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 100, val)
	})

	// Test error type transformation
	t.Run("Error type transformation", func(t *testing.T) {
		result := MonadChainLeft(
			Left[int](errors.New("404")),
			func(e error) IOResult[int] {
				return Left[int](errors.New("404"))
			},
		)
		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "404", err.Error())
	})
}

func TestChainLeft(t *testing.T) {
	// Test with Left value - should apply the function
	t.Run("Left value applies function", func(t *testing.T) {
		chainFn := ChainLeft(func(e error) IOResult[int] {
			return Left[int](errors.New("chained: " + e.Error()))
		})
		result := F.Pipe1(
			Left[int](errors.New("original")),
			chainFn,
		)
		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "chained: original", err.Error())
	})

	// Test with Left value - function returns Right (error recovery)
	t.Run("Left value recovers to Right", func(t *testing.T) {
		chainFn := ChainLeft(func(e error) IOResult[int] {
			if e.Error() == "network error" {
				return Right(0) // default value
			}
			return Left[int](e)
		})
		result := F.Pipe1(
			Left[int](errors.New("network error")),
			chainFn,
		)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 0, val)
	})

	// Test with Right value - should pass through unchanged
	t.Run("Right value passes through", func(t *testing.T) {
		chainFn := ChainLeft(func(e error) IOResult[int] {
			return Left[int](errors.New("should not be called"))
		})
		result := F.Pipe1(
			Right(42),
			chainFn,
		)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	// Test composition with other operations
	t.Run("Composition with Map", func(t *testing.T) {
		result := F.Pipe2(
			Left[int](errors.New("error")),
			ChainLeft(func(e error) IOResult[int] {
				return Left[int](errors.New("handled: " + e.Error()))
			}),
			Map(utils.Double),
		)
		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "handled: error", err.Error())
	})
}

func TestMonadChainFirstLeft(t *testing.T) {
	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves original error", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Left[int](errors.New("original error")),
			func(e error) IOResult[int] {
				sideEffectCalled = true
				return Left[int](errors.New("new error")) // This error is ignored, original is returned
			},
		)
		_, err := result()
		assert.True(t, sideEffectCalled)
		assert.Error(t, err)
		assert.Equal(t, "original error", err.Error())
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var capturedError string
		result := MonadChainFirstLeft(
			Left[int](errors.New("validation failed")),
			func(e error) IOResult[int] {
				capturedError = e.Error()
				return Right(999) // This Right value is ignored, original Left is returned
			},
		)
		_, err := result()
		assert.Equal(t, "validation failed", capturedError)
		assert.Error(t, err)
		assert.Equal(t, "validation failed", err.Error())
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Right(42),
			func(e error) IOResult[int] {
				sideEffectCalled = true
				return Left[int](errors.New("should not be called"))
			},
		)
		assert.False(t, sideEffectCalled)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	// Test that side effects are executed but original error is always preserved
	t.Run("Side effects executed but original error preserved", func(t *testing.T) {
		effectCount := 0
		result := MonadChainFirstLeft(
			Left[int](errors.New("original error")),
			func(e error) IOResult[int] {
				effectCount++
				// Try to return Right, but original Left should still be returned
				return Right(999)
			},
		)
		_, err := result()
		assert.Equal(t, 1, effectCount)
		assert.Error(t, err)
		assert.Equal(t, "original error", err.Error())
	})
}

func TestChainFirstLeft(t *testing.T) {
	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves error", func(t *testing.T) {
		var captured string
		chainFn := ChainFirstLeft[int](func(e error) IOResult[int] {
			captured = e.Error()
			return Left[int](errors.New("ignored error"))
		})
		result := F.Pipe1(
			Left[int](errors.New("test error")),
			chainFn,
		)
		_, err := result()
		assert.Equal(t, "test error", captured)
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var captured string
		chainFn := ChainFirstLeft[int](func(e error) IOResult[int] {
			captured = e.Error()
			return Right(42) // This Right is ignored, original Left is returned
		})
		result := F.Pipe1(
			Left[int](errors.New("test error")),
			chainFn,
		)
		_, err := result()
		assert.Equal(t, "test error", captured)
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		called := false
		chainFn := ChainFirstLeft[int](func(e error) IOResult[int] {
			called = true
			return Right(0)
		})
		result := F.Pipe1(
			Right(100),
			chainFn,
		)
		assert.False(t, called)
		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 100, val)
	})

	// Test that original error is always preserved regardless of what f returns
	t.Run("Original error always preserved", func(t *testing.T) {
		chainFn := ChainFirstLeft[int](func(e error) IOResult[int] {
			// Try to return Right, but original Left should still be returned
			return Right(999)
		})

		result := F.Pipe1(
			Left[int](errors.New("original")),
			chainFn,
		)
		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "original", err.Error())
	})

	// Test with IO side effects - original Left is always preserved
	t.Run("IO side effects with Left preservation", func(t *testing.T) {
		effectCount := 0
		chainFn := ChainFirstLeft[int](func(e error) IOResult[int] {
			return FromIO(func() int {
				effectCount++
				return 0
			})
		})

		// Even though FromIO wraps in Right, the original Left is preserved
		result := F.Pipe1(
			Left[int](errors.New("error")),
			chainFn,
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "error", err.Error())
		assert.Equal(t, 1, effectCount)
	})

	// Test logging with Left preservation
	t.Run("Logging with Left preservation", func(t *testing.T) {
		errorLog := []string{}
		logError := ChainFirstLeft[string](func(e error) IOResult[string] {
			errorLog = append(errorLog, "Logged: "+e.Error())
			return Left[string](errors.New("log entry")) // This is ignored, original is preserved
		})

		result := F.Pipe2(
			Left[string](errors.New("step1")),
			logError,
			ChainLeft(func(e error) IOResult[string] {
				return Left[string](errors.New("step2"))
			}),
		)

		_, err := result()
		assert.Equal(t, []string{"Logged: step1"}, errorLog)
		assert.Error(t, err)
		assert.Equal(t, "step2", err.Error())
	})
}
