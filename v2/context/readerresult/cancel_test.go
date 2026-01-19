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
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

// TestWithContext tests the WithContext function
func TestWithContext(t *testing.T) {
	t.Run("executes wrapped ReaderResult when context is not cancelled", func(t *testing.T) {
		executed := false
		computation := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		wrapped := WithContext(computation)
		result := wrapped(t.Context())

		assert.True(t, executed, "computation should be executed")
		assert.Equal(t, E.Of[error](42), result)
	})

	t.Run("returns cancellation error when context is cancelled", func(t *testing.T) {
		executed := false
		computation := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		wrapped := WithContext(computation)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := wrapped(ctx)

		assert.False(t, executed, "computation should not be executed when context is cancelled")
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("returns deadline exceeded error when context times out", func(t *testing.T) {
		executed := false
		computation := func(ctx context.Context) E.Either[error, int] {
			executed = true
			time.Sleep(100 * time.Millisecond)
			return E.Of[error](42)
		}

		wrapped := WithContext(computation)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()

		time.Sleep(20 * time.Millisecond) // Wait for timeout

		result := wrapped(ctx)

		assert.False(t, executed, "computation should not be executed when context has timed out")
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("preserves errors from wrapped computation", func(t *testing.T) {
		testErr := errors.New("computation error")
		computation := func(ctx context.Context) E.Either[error, int] {
			return E.Left[int](testErr)
		}

		wrapped := WithContext(computation)
		result := wrapped(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("prevents expensive computation when context is already cancelled", func(t *testing.T) {
		expensiveExecuted := false
		expensiveComputation := func(ctx context.Context) E.Either[error, int] {
			expensiveExecuted = true
			// Simulate expensive operation
			time.Sleep(1 * time.Second)
			return E.Of[error](42)
		}

		wrapped := WithContext(expensiveComputation)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		start := time.Now()
		result := wrapped(ctx)
		duration := time.Since(start)

		assert.False(t, expensiveExecuted, "expensive computation should not execute")
		assert.True(t, E.IsLeft(result))
		assert.Less(t, duration, 100*time.Millisecond, "should return immediately")
	})

	t.Run("works with context.WithCancelCause", func(t *testing.T) {
		executed := false
		computation := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		wrapped := WithContext(computation)

		customErr := errors.New("custom cancellation reason")
		ctx, cancel := context.WithCancelCause(t.Context())
		cancel(customErr)

		result := wrapped(ctx)

		assert.False(t, executed, "computation should not be executed")
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, customErr, err)
	})

	t.Run("can be nested for multiple cancellation checks", func(t *testing.T) {
		executed := false
		computation := func(ctx context.Context) E.Either[error, int] {
			executed = true
			return E.Of[error](42)
		}

		doubleWrapped := WithContext(WithContext(computation))

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := doubleWrapped(ctx)

		assert.False(t, executed, "computation should not be executed")
		assert.True(t, E.IsLeft(result))
	})
}

// TestWithContextK tests the WithContextK function
func TestWithContextK(t *testing.T) {
	t.Run("wraps Kleisli arrow with cancellation checking", func(t *testing.T) {
		executed := false
		processUser := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				executed = true
				return E.Of[error]("user-" + string(rune(id+48)))
			}
		}

		safeProcessUser := WithContextK(processUser)

		result := safeProcessUser(123)(t.Context())

		assert.True(t, executed, "Kleisli should be executed")
		assert.True(t, E.IsRight(result))
	})

	t.Run("prevents Kleisli execution when context is cancelled", func(t *testing.T) {
		executed := false
		processUser := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				executed = true
				return E.Of[error]("user")
			}
		}

		safeProcessUser := WithContextK(processUser)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := safeProcessUser(123)(ctx)

		assert.False(t, executed, "Kleisli should not be executed when context is cancelled")
		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, context.Canceled, err)
	})

	t.Run("works in Chain pipeline", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		getUser := WithContextK(func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				firstExecuted = true
				return E.Of[error]("Alice")
			}
		})

		getOrders := WithContextK(func(name string) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				secondExecuted = true
				return E.Of[error](5)
			}
		})

		pipeline := F.Pipe2(
			Of(123),
			Chain(getUser),
			Chain(getOrders),
		)

		result := pipeline(t.Context())

		assert.True(t, firstExecuted, "first step should execute")
		assert.True(t, secondExecuted, "second step should execute")
		assert.Equal(t, E.Of[error](5), result)
	})

	t.Run("stops pipeline on cancellation", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		getUser := WithContextK(func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				firstExecuted = true
				return E.Of[error]("Alice")
			}
		})

		getOrders := WithContextK(func(name string) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				secondExecuted = true
				return E.Of[error](5)
			}
		})

		pipeline := F.Pipe2(
			Of(123),
			Chain(getUser),
			Chain(getOrders),
		)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := pipeline(ctx)

		assert.False(t, firstExecuted, "first step should not execute")
		assert.False(t, secondExecuted, "second step should not execute")
		assert.True(t, E.IsLeft(result))
	})

	t.Run("respects timeout in multi-step pipeline", func(t *testing.T) {
		step1Executed := false
		step2Executed := false

		step1 := WithContextK(func(x int) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				step1Executed = true
				time.Sleep(50 * time.Millisecond)
				return E.Of[error](x * 2)
			}
		})

		step2 := WithContextK(func(x int) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				step2Executed = true
				return E.Of[error](x + 10)
			}
		})

		pipeline := F.Pipe2(
			Of(5),
			Chain(step1),
			Chain(step2),
		)

		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()

		time.Sleep(20 * time.Millisecond) // Wait for timeout

		result := pipeline(ctx)

		assert.False(t, step1Executed, "step1 should not execute after timeout")
		assert.False(t, step2Executed, "step2 should not execute after timeout")
		assert.True(t, E.IsLeft(result))
	})

	t.Run("preserves errors from Kleisli computation", func(t *testing.T) {
		testErr := errors.New("kleisli error")
		failingKleisli := func(id int) ReaderResult[string] {
			return func(ctx context.Context) E.Either[error, string] {
				return E.Left[string](testErr)
			}
		}

		safeKleisli := WithContextK(failingKleisli)
		result := safeKleisli(123)(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})
}

// TestWithContextIntegration tests integration scenarios
func TestWithContextIntegration(t *testing.T) {
	t.Run("WithContext in complex pipeline with multiple operations", func(t *testing.T) {
		step1Executed := false
		step2Executed := false
		step3Executed := false

		step1 := WithContext(func(ctx context.Context) E.Either[error, int] {
			step1Executed = true
			return E.Of[error](10)
		})

		step2 := WithContextK(func(x int) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				step2Executed = true
				return E.Of[error](x * 2)
			}
		})

		step3 := WithContext(func(ctx context.Context) E.Either[error, string] {
			step3Executed = true
			return E.Of[error]("done")
		})

		pipeline := F.Pipe2(
			step1,
			Chain(step2),
			ChainTo[int](step3),
		)

		result := pipeline(t.Context())

		assert.True(t, step1Executed)
		assert.True(t, step2Executed)
		assert.True(t, step3Executed)
		assert.Equal(t, E.Of[error]("done"), result)
	})

	t.Run("early cancellation prevents all subsequent operations", func(t *testing.T) {
		step1Executed := false
		step2Executed := false
		step3Executed := false

		step1 := WithContext(func(ctx context.Context) E.Either[error, int] {
			step1Executed = true
			return E.Of[error](10)
		})

		step2 := WithContextK(func(x int) ReaderResult[int] {
			return func(ctx context.Context) E.Either[error, int] {
				step2Executed = true
				return E.Of[error](x * 2)
			}
		})

		step3 := WithContext(func(ctx context.Context) E.Either[error, string] {
			step3Executed = true
			return E.Of[error]("done")
		})

		pipeline := F.Pipe2(
			step1,
			Chain(step2),
			ChainTo[int](step3),
		)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := pipeline(ctx)

		assert.False(t, step1Executed, "no steps should execute")
		assert.False(t, step2Executed, "no steps should execute")
		assert.False(t, step3Executed, "no steps should execute")
		assert.True(t, E.IsLeft(result))
	})

	t.Run("WithContext with Map and Chain", func(t *testing.T) {
		computation := WithContext(func(ctx context.Context) E.Either[error, int] {
			return E.Of[error](42)
		})

		pipeline := F.Pipe2(
			computation,
			Map(N.Mul(2)),
			Map(reader.Of[int]("result")),
		)

		result := pipeline(t.Context())
		assert.Equal(t, E.Of[error]("result"), result)
	})
}
