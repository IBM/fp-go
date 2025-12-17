// Copyright (c) 2025 IBM Corp.
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

	"github.com/IBM/fp-go/v2/result"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a test retry policy
func testRetryPolicy() R.RetryPolicy {
	return R.Monoid.Concat(
		R.LimitRetries(5),
		R.CapDelay(1*time.Second, R.ExponentialBackoff(10*time.Millisecond)),
	)
}

// TestRetrying_SuccessOnFirstAttempt tests that Retrying succeeds immediately
// when the action succeeds on the first attempt.
func TestRetrying_SuccessOnFirstAttempt(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				return result.Of("success")
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()

	assert.Equal(t, result.Of("success"), res)
}

// TestRetrying_SuccessAfterRetries tests that Retrying eventually succeeds
// after a few failed attempts.
func TestRetrying_SuccessAfterRetries(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				// Fail on first 3 attempts, succeed on 4th
				if status.IterNumber < 3 {
					return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
				}
				return result.Of(fmt.Sprintf("success on attempt %d", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()

	assert.Equal(t, result.Of("success on attempt 3"), res)
}

// TestRetrying_ExhaustsRetries tests that Retrying stops after the retry limit
// is reached and returns the last error.
func TestRetrying_ExhaustsRetries(t *testing.T) {
	policy := R.LimitRetries(3)

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()

	assert.True(t, result.IsLeft(res))
	assert.Equal(t, result.Left[string](fmt.Errorf("attempt 3 failed")), res)
}

// TestRetrying_ActionChecksContextCancellation tests that actions can check
// the context and return early if it's cancelled.
func TestRetrying_ActionChecksContextCancellation(t *testing.T) {
	policy := R.LimitRetries(10)

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				attemptCount++

				// Check context at the start of the action
				if ctx.Err() != nil {
					return result.Left[string](ctx.Err())
				}

				// Simulate work that might take time
				time.Sleep(10 * time.Millisecond)

				// Check context again after work
				if ctx.Err() != nil {
					return result.Left[string](ctx.Err())
				}

				// Always fail to trigger retries
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		// Don't retry on context errors
		val, err := result.Unwrap(r)
		_ = val
		if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			return false
		}
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)

	// Create a context that we'll cancel after a short time
	ctx, cancel := context.WithCancel(t.Context())

	// Start the retry operation in a goroutine
	resultChan := make(chan Result[string], 1)
	go func() {
		res := retrying(ctx)()
		resultChan <- res
	}()

	// Cancel the context after allowing a couple attempts
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the result
	res := <-resultChan

	// Should have stopped due to context cancellation
	assert.True(t, result.IsLeft(res))

	// Should have stopped early (not all 10 attempts)
	assert.Less(t, attemptCount, 10, "Should stop retrying when action detects context cancellation")

	// The error should be related to context cancellation or an early attempt
	val, err := result.Unwrap(res)
	_ = val
	assert.Error(t, err)
}

// TestRetrying_ContextCancelledBeforeStart tests that if the context is already
// cancelled before starting, the operation fails immediately.
func TestRetrying_ContextCancelledBeforeStart(t *testing.T) {
	policy := testRetryPolicy()

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				attemptCount++
				// Check context before doing work
				if ctx.Err() != nil {
					return result.Left[string](ctx.Err())
				}
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		// Don't retry on context errors
		val, err := result.Unwrap(r)
		_ = val
		if err != nil && errors.Is(err, context.Canceled) {
			return false
		}
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	res := retrying(ctx)()

	assert.True(t, result.IsLeft(res))
	val, err := result.Unwrap(res)
	_ = val
	assert.True(t, errors.Is(err, context.Canceled))

	// Should have attempted at most once
	assert.LessOrEqual(t, attemptCount, 1)
}

// TestRetrying_ContextTimeoutInAction tests that actions respect context deadlines.
func TestRetrying_ContextTimeoutInAction(t *testing.T) {
	policy := R.LimitRetries(10)

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				attemptCount++

				// Check context before doing work
				if ctx.Err() != nil {
					return result.Left[string](ctx.Err())
				}

				// Simulate some work
				time.Sleep(50 * time.Millisecond)

				// Check context after work
				if ctx.Err() != nil {
					return result.Left[string](ctx.Err())
				}

				// Always fail to trigger retries
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		// Don't retry on context errors
		val, err := result.Unwrap(r)
		_ = val
		if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			return false
		}
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(t.Context(), 150*time.Millisecond)
	defer cancel()

	startTime := time.Now()
	res := retrying(ctx)()
	elapsed := time.Since(startTime)

	assert.True(t, result.IsLeft(res))

	// Should have stopped before completing all 10 retries
	assert.Less(t, attemptCount, 10, "Should stop retrying when action detects context timeout")

	// Should have stopped around the timeout duration
	assert.Less(t, elapsed, 500*time.Millisecond, "Should stop soon after timeout")
}

// TestRetrying_CheckFunctionStopsRetry tests that the check function can
// stop retrying even when errors occur.
func TestRetrying_CheckFunctionStopsRetry(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				if status.IterNumber == 0 {
					return result.Left[string](fmt.Errorf("retryable error"))
				}
				return result.Left[string](fmt.Errorf("permanent error"))
			}
		}
	}

	// Only retry on "retryable error"
	check := func(r Result[string]) bool {
		return result.IsLeft(r) && result.Fold(
			func(err error) bool { return err.Error() == "retryable error" },
			func(string) bool { return false },
		)(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()

	assert.Equal(t, result.Left[string](fmt.Errorf("permanent error")), res)
}

// TestRetrying_ExponentialBackoff tests that exponential backoff is applied.
func TestRetrying_ExponentialBackoff(t *testing.T) {
	// Use a policy with measurable delays
	policy := R.Monoid.Concat(
		R.LimitRetries(3),
		R.ExponentialBackoff(50*time.Millisecond),
	)

	startTime := time.Now()

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				if status.IterNumber < 2 {
					return result.Left[string](fmt.Errorf("retry"))
				}
				return result.Of("success")
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()
	elapsed := time.Since(startTime)

	assert.Equal(t, result.Of("success"), res)
	// With exponential backoff starting at 50ms:
	// Iteration 0: no delay
	// Iteration 1: 50ms delay
	// Iteration 2: 100ms delay
	// Total should be at least 150ms
	assert.GreaterOrEqual(t, elapsed, 150*time.Millisecond)
}

// TestRetrying_ContextValuePropagation tests that context values are properly
// propagated through the retry mechanism.
func TestRetrying_ContextValuePropagation(t *testing.T) {
	policy := R.LimitRetries(2)

	type contextKey string
	const requestIDKey contextKey = "requestID"

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				// Extract value from context
				requestID, ok := ctx.Value(requestIDKey).(string)
				if !ok {
					return result.Left[string](fmt.Errorf("missing request ID"))
				}

				if status.IterNumber < 1 {
					return result.Left[string](fmt.Errorf("retry needed"))
				}

				return result.Of(fmt.Sprintf("processed request %s", requestID))
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)

	// Create context with a value
	ctx := context.WithValue(t.Context(), requestIDKey, "12345")

	res := retrying(ctx)()

	assert.Equal(t, result.Of("processed request 12345"), res)
}

// TestRetrying_RetryStatusProgression tests that the RetryStatus is properly
// updated on each iteration.
func TestRetrying_RetryStatusProgression(t *testing.T) {
	policy := testRetryPolicy()

	var iterations []uint

	action := func(status R.RetryStatus) ReaderIOResult[int] {
		return func(ctx context.Context) IOResult[int] {
			return func() Result[int] {
				iterations = append(iterations, status.IterNumber)
				if status.IterNumber < 3 {
					return result.Left[int](fmt.Errorf("retry"))
				}
				return result.Of(int(status.IterNumber))
			}
		}
	}

	check := func(r Result[int]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := t.Context()

	res := retrying(ctx)()

	assert.Equal(t, result.Of(3), res)
	// Should have attempted iterations 0, 1, 2, 3
	assert.Equal(t, []uint{0, 1, 2, 3}, iterations)
}

// TestRetrying_ContextCancelledDuringDelay tests that the retry operation
// stops immediately when the context is cancelled during a retry delay,
// even if there are still retries remaining according to the policy.
func TestRetrying_ContextCancelledDuringDelay(t *testing.T) {
	// Use a policy with significant delays to ensure we can cancel during the delay
	policy := R.Monoid.Concat(
		R.LimitRetries(10),
		R.ConstantDelay(200*time.Millisecond),
	)

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderIOResult[string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				attemptCount++
				// Always fail to trigger retries
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	// Always retry on errors (don't check for context cancellation in check function)
	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)

	// Create a context that we'll cancel during the retry delay
	ctx, cancel := context.WithCancel(t.Context())

	// Start the retry operation in a goroutine
	resultChan := make(chan Result[string], 1)
	startTime := time.Now()
	go func() {
		res := retrying(ctx)()
		resultChan <- res
	}()

	// Wait for the first attempt to complete and the delay to start
	time.Sleep(50 * time.Millisecond)

	// Cancel the context during the retry delay
	cancel()

	// Wait for the result
	res := <-resultChan
	elapsed := time.Since(startTime)

	// Should have stopped due to context cancellation
	assert.True(t, result.IsLeft(res))

	// Should have attempted only once or twice (not all 10 attempts)
	// because the context was cancelled during the delay
	assert.LessOrEqual(t, attemptCount, 2, "Should stop retrying when context is cancelled during delay")

	// Should have stopped quickly after cancellation, not waiting for all delays
	// With 10 retries and 200ms delays, it would take ~2 seconds without cancellation
	// With cancellation during first delay, it should complete in well under 500ms
	assert.Less(t, elapsed, 500*time.Millisecond, "Should stop immediately when context is cancelled during delay")

	// When context is cancelled during the delay, the retry mechanism
	// detects the cancellation and returns a context error
	val, err := result.Unwrap(res)
	_ = val
	assert.Error(t, err)
	// The error should be a context cancellation error since cancellation
	// happened during the delay between retries
	assert.True(t, errors.Is(err, context.Canceled), "Should return context.Canceled when cancelled during delay")
}
