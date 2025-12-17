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
	"fmt"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/result"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// Test configuration type
type testConfig struct {
	maxRetries int
	baseURL    string
}

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

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
			return func() Result[string] {
				return result.Of("success")
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	cfg := testConfig{maxRetries: 5, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Of("success"), res)
}

// TestRetrying_SuccessAfterRetries tests that Retrying eventually succeeds
// after a few failed attempts.
func TestRetrying_SuccessAfterRetries(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
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
	cfg := testConfig{maxRetries: 5, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Of("success on attempt 3"), res)
}

// TestRetrying_ExhaustsRetries tests that Retrying stops after the retry limit
// is reached and returns the last error.
func TestRetrying_ExhaustsRetries(t *testing.T) {
	policy := R.LimitRetries(3)

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
			return func() Result[string] {
				return result.Left[string](fmt.Errorf("attempt %d failed", status.IterNumber))
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	cfg := testConfig{maxRetries: 3, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.True(t, result.IsLeft(res))
	assert.Equal(t, result.Left[string](fmt.Errorf("attempt 3 failed")), res)
}

// TestRetrying_CheckFunctionStopsRetry tests that the check function can
// stop retrying even when errors occur.
func TestRetrying_CheckFunctionStopsRetry(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
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
	cfg := testConfig{maxRetries: 5, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Left[string](fmt.Errorf("permanent error")), res)
}

// TestRetrying_UsesContext tests that the Reader context is properly passed
// through to the action.
func TestRetrying_UsesContext(t *testing.T) {
	policy := testRetryPolicy()

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
			return func() Result[string] {
				// Use the config from context
				if status.IterNumber < 2 {
					return result.Left[string](fmt.Errorf("retry needed"))
				}
				return result.Of(fmt.Sprintf("success with baseURL: %s", cfg.baseURL))
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	cfg := testConfig{maxRetries: 5, baseURL: "https://test.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Of("success with baseURL: https://test.example.com"), res)
}

// TestRetrying_RetryStatusProgression tests that the RetryStatus is properly
// updated on each iteration.
func TestRetrying_RetryStatusProgression(t *testing.T) {
	policy := testRetryPolicy()

	var iterations []uint

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, int] {
		return func(cfg testConfig) IOResult[int] {
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
	cfg := testConfig{maxRetries: 5, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Of(3), res)
	// Should have attempted iterations 0, 1, 2, 3
	assert.Equal(t, []uint{0, 1, 2, 3}, iterations)
}

// TestRetrying_WithContextContext tests using context.Context as the Reader type.
func TestRetrying_WithContextContext(t *testing.T) {
	policy := R.LimitRetries(2)

	action := func(status R.RetryStatus) ReaderIOResult[context.Context, string] {
		return func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				if status.IterNumber == 0 {
					return result.Left[string](fmt.Errorf("first attempt failed"))
				}
				return result.Of("success")
			}
		}
	}

	check := func(r Result[string]) bool {
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	ctx := context.Background()

	res := retrying(ctx)()

	assert.Equal(t, result.Of("success"), res)
}

// TestRetrying_NoRetryOnSuccess tests that successful results are not retried
// even if the check function would return true.
func TestRetrying_NoRetryOnSuccess(t *testing.T) {
	policy := testRetryPolicy()

	callCount := 0

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
			return func() Result[string] {
				callCount++
				return result.Of("success")
			}
		}
	}

	// This check would normally trigger a retry, but since the result is Right,
	// it should not be called or should not cause a retry
	check := func(r Result[string]) bool {
		// Only retry on Left (error) results
		return result.IsLeft(r)
	}

	retrying := Retrying(policy, action, check)
	cfg := testConfig{maxRetries: 5, baseURL: "https://api.example.com"}

	res := retrying(cfg)()

	assert.Equal(t, result.Of("success"), res)
	// Should only be called once since it succeeded immediately
	assert.Equal(t, 1, callCount)
}

// TestRetrying_ExponentialBackoff tests that exponential backoff is applied.
func TestRetrying_ExponentialBackoff(t *testing.T) {
	// Use a policy with measurable delays
	policy := R.Monoid.Concat(
		R.LimitRetries(3),
		R.ExponentialBackoff(50*time.Millisecond),
	)

	startTime := time.Now()

	action := func(status R.RetryStatus) ReaderIOResult[testConfig, string] {
		return func(cfg testConfig) IOResult[string] {
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
	cfg := testConfig{maxRetries: 3, baseURL: "https://api.example.com"}

	res := retrying(cfg)()
	elapsed := time.Since(startTime)

	assert.Equal(t, result.Of("success"), res)
	// With exponential backoff starting at 50ms:
	// Iteration 0: no delay
	// Iteration 1: 50ms delay
	// Iteration 2: 100ms delay
	// Total should be at least 150ms
	assert.GreaterOrEqual(t, elapsed, 150*time.Millisecond)
}
