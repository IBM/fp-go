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

package readerresult

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

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

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			return "success", nil
		}
	}

	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	result, err := retrying(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
}

// TestRetrying_SuccessAfterRetries tests that Retrying eventually succeeds
// after a few failed attempts.
func TestRetrying_SuccessAfterRetries(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			// Fail on first 3 attempts, succeed on 4th
			if status.IterNumber < 3 {
				return "", fmt.Errorf("attempt %d failed", status.IterNumber)
			}
			return fmt.Sprintf("success on attempt %d", status.IterNumber), nil
		}
	}

	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	result, err := retrying(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "success on attempt 3", result)
}

// TestRetrying_ExhaustsRetries tests that Retrying stops after the retry limit
// is reached and returns the last error.
func TestRetrying_ExhaustsRetries(t *testing.T) {
	policy := R.LimitRetries(3)

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			return "", fmt.Errorf("attempt %d failed", status.IterNumber)
		}
	}

	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	result, err := retrying(ctx)

	assert.Error(t, err)
	assert.Equal(t, "", result)
	assert.Equal(t, "attempt 3 failed", err.Error())
}

// TestRetrying_ActionChecksContextCancellation tests that actions can check
// the context and return early if it's cancelled.
func TestRetrying_ActionChecksContextCancellation(t *testing.T) {
	policy := R.LimitRetries(10)

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			attemptCount++

			// Check context at the start of the action
			if ctx.Err() != nil {
				return "", ctx.Err()
			}

			// Simulate work that might take time
			time.Sleep(10 * time.Millisecond)

			// Check context again after work
			if ctx.Err() != nil {
				return "", ctx.Err()
			}

			// Always fail to trigger retries
			return "", fmt.Errorf("attempt %d failed", status.IterNumber)
		}
	}

	check := func(val string, err error) bool {
		// Don't retry on context errors
		if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			return false
		}
		return err != nil
	}

	retrying := Retrying(policy, action, check)

	// Create a context that we'll cancel after a short time
	ctx, cancel := context.WithCancel(context.Background())

	// Start the retry operation in a goroutine
	type result struct {
		val string
		err error
	}
	resultChan := make(chan result, 1)
	go func() {
		val, err := retrying(ctx)
		resultChan <- result{val, err}
	}()

	// Cancel the context after allowing a couple attempts
	time.Sleep(50 * time.Millisecond)
	cancel()

	// Wait for the result
	res := <-resultChan

	// Should have stopped due to context cancellation
	assert.Error(t, res.err)

	// Should have stopped early (not all 10 attempts)
	assert.Less(t, attemptCount, 10, "Should stop retrying when action detects context cancellation")
}

// TestRetrying_ContextCancelledBeforeStart tests that if the context is already
// cancelled before starting, the operation fails immediately.
func TestRetrying_ContextCancelledBeforeStart(t *testing.T) {
	policy := testRetryPolicy()

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			attemptCount++
			// Check context before doing work
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "", fmt.Errorf("attempt %d failed", status.IterNumber)
		}
	}

	check := func(val string, err error) bool {
		// Don't retry on context errors
		if err != nil && errors.Is(err, context.Canceled) {
			return false
		}
		return err != nil
	}

	retrying := Retrying(policy, action, check)

	// Create an already-cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := retrying(ctx)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled))

	// Should have attempted at most once
	assert.LessOrEqual(t, attemptCount, 1)
}

// TestRetrying_ContextTimeoutInAction tests that actions respect context deadlines.
func TestRetrying_ContextTimeoutInAction(t *testing.T) {
	policy := R.LimitRetries(10)

	attemptCount := 0

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			attemptCount++

			// Check context before doing work
			if ctx.Err() != nil {
				return "", ctx.Err()
			}

			// Simulate some work
			time.Sleep(50 * time.Millisecond)

			// Check context after work
			if ctx.Err() != nil {
				return "", ctx.Err()
			}

			// Always fail to trigger retries
			return "", fmt.Errorf("attempt %d failed", status.IterNumber)
		}
	}

	check := func(val string, err error) bool {
		// Don't retry on context errors
		if err != nil && (errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)) {
			return false
		}
		return err != nil
	}

	retrying := Retrying(policy, action, check)

	// Create a context with a short timeout
	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	startTime := time.Now()
	_, err := retrying(ctx)
	elapsed := time.Since(startTime)

	assert.Error(t, err)

	// Should have stopped before completing all 10 retries
	assert.Less(t, attemptCount, 10, "Should stop retrying when action detects context timeout")

	// Should have stopped around the timeout duration
	assert.Less(t, elapsed, 500*time.Millisecond, "Should stop soon after timeout")
}

// TestRetrying_CheckFunctionStopsRetry tests that the check function can
// stop retrying even when errors occur.
func TestRetrying_CheckFunctionStopsRetry(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			if status.IterNumber == 0 {
				return "", fmt.Errorf("retryable error")
			}
			return "", fmt.Errorf("permanent error")
		}
	}

	// Only retry on "retryable error"
	check := func(val string, err error) bool {
		return err != nil && err.Error() == "retryable error"
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	_, err := retrying(ctx)

	assert.Error(t, err)
	assert.Equal(t, "permanent error", err.Error())
}

// TestRetrying_ExponentialBackoff tests that exponential backoff is applied.
func TestRetrying_ExponentialBackoff(t *testing.T) {
	// Use a policy with measurable delays
	policy := R.Monoid.Concat(
		R.LimitRetries(3),
		R.ExponentialBackoff(50*time.Millisecond),
	)

	startTime := time.Now()

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			if status.IterNumber < 2 {
				return "", fmt.Errorf("retry")
			}
			return "success", nil
		}
	}

	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	result, err := retrying(ctx)
	elapsed := time.Since(startTime)

	assert.NoError(t, err)
	assert.Equal(t, "success", result)
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

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			// Extract value from context
			requestID, ok := ctx.Value(requestIDKey).(string)
			if !ok {
				return "", fmt.Errorf("missing request ID")
			}

			if status.IterNumber < 1 {
				return "", fmt.Errorf("retry needed")
			}

			return fmt.Sprintf("processed request %s", requestID), nil
		}
	}

	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)

	// Create context with a value
	ctx := context.WithValue(context.Background(), requestIDKey, "12345")

	result, err := retrying(ctx)

	assert.NoError(t, err)
	assert.Equal(t, "processed request 12345", result)
}

// TestRetrying_RetryStatusProgression tests that the RetryStatus is properly
// updated on each iteration.
func TestRetrying_RetryStatusProgression(t *testing.T) {
	policy := testRetryPolicy()

	var iterations []uint

	action := func(status R.RetryStatus) ReaderResult[int] {
		return func(ctx context.Context) (int, error) {
			iterations = append(iterations, status.IterNumber)
			if status.IterNumber < 3 {
				return 0, fmt.Errorf("retry")
			}
			return int(status.IterNumber), nil
		}
	}

	check := func(val int, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	result, err := retrying(ctx)

	assert.NoError(t, err)
	assert.Equal(t, 3, result)
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

	action := func(status R.RetryStatus) ReaderResult[string] {
		return func(ctx context.Context) (string, error) {
			attemptCount++
			// Always fail to trigger retries
			return "", fmt.Errorf("attempt %d failed", status.IterNumber)
		}
	}

	// Always retry on errors (don't check for context cancellation in check function)
	check := func(val string, err error) bool {
		return err != nil
	}

	retrying := Retrying(policy, action, check)

	// Create a context that we'll cancel during the retry delay
	ctx, cancel := context.WithCancel(context.Background())

	// Start the retry operation in a goroutine
	type result struct {
		val string
		err error
	}
	resultChan := make(chan result, 1)
	startTime := time.Now()
	go func() {
		val, err := retrying(ctx)
		resultChan <- result{val, err}
	}()

	// Wait for the first attempt to complete and the delay to start
	time.Sleep(50 * time.Millisecond)

	// Cancel the context during the retry delay
	cancel()

	// Wait for the result
	res := <-resultChan
	elapsed := time.Since(startTime)

	// Should have stopped due to context cancellation
	assert.Error(t, res.err)

	// Should have attempted only once or twice (not all 10 attempts)
	// because the context was cancelled during the delay
	assert.LessOrEqual(t, attemptCount, 2, "Should stop retrying when context is cancelled during delay")

	// Should have stopped quickly after cancellation, not waiting for all delays
	// With 10 retries and 200ms delays, it would take ~2 seconds without cancellation
	// With cancellation during first delay, it should complete in well under 500ms
	assert.Less(t, elapsed, 500*time.Millisecond, "Should stop immediately when context is cancelled during delay")

	// When context is cancelled during the delay, the retry mechanism
	// detects the cancellation and returns a context error
	assert.True(t, errors.Is(res.err, context.Canceled), "Should return context.Canceled when cancelled during delay")
}
