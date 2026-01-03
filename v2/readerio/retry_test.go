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

package readerio

import (
	"fmt"
	"strings"
	"testing"
	"time"

	G "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// Test environment types for retry tests
type RetryTestEnv struct {
	Logger     func(string)
	MaxRetries int
	Config     string
}

// Test helpers and common policies

var expLogBackoff = R.ExponentialBackoff(10)

// our retry policy with a 2s cap
var testLogPolicy = R.CapDelay(
	2*time.Second,
	R.Monoid.Concat(expLogBackoff, R.LimitRetries(20)),
)

// TestRetrying_BasicSuccess tests that Retrying succeeds when the check predicate returns false
func TestRetrying_BasicSuccess(t *testing.T) {
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, string] {
		return func(e RetryTestEnv) G.IO[string] {
			return G.Of(fmt.Sprintf("Retrying %d", status.IterNumber))
		}
	}
	check := func(value string) bool {
		return !strings.Contains(value, "5")
	}

	r := Retrying(testLogPolicy, action, check)

	assert.Equal(t, "Retrying 5", r(env)())
}

// TestRetrying_ImmediateSuccess tests that no retries occur when the first attempt succeeds
func TestRetrying_ImmediateSuccess(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				return 42
			}())
		}
	}
	check := func(value int) bool {
		return false // Never retry
	}

	policy := R.LimitRetries(5)
	result := Retrying(policy, action, check)

	assert.Equal(t, 42, result(env)())
	assert.Equal(t, 1, attempts, "Should only execute once when immediately successful")
}

// TestRetrying_MaxRetriesReached tests that retrying stops when the policy limit is reached
func TestRetrying_MaxRetriesReached(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, string] {
		return func(e RetryTestEnv) G.IO[string] {
			return G.Of(func() string {
				attempts++
				return fmt.Sprintf("attempt_%d", attempts)
			}())
		}
	}
	check := func(value string) bool {
		return true // Always retry
	}

	policy := R.LimitRetries(3)
	result := Retrying(policy, action, check)

	finalResult := result(env)()
	assert.Equal(t, "attempt_4", finalResult, "Should execute initial attempt + 3 retries")
	assert.Equal(t, 4, attempts, "Should execute 4 times total (1 initial + 3 retries)")
}

// TestRetrying_StatusTracking tests that RetryStatus is correctly updated across retries
func TestRetrying_StatusTracking(t *testing.T) {
	var statuses []R.RetryStatus
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				statuses = append(statuses, status)
				return len(statuses)
			}())
		}
	}
	check := func(value int) bool {
		return value < 3 // Retry until we've done 3 attempts
	}

	policy := R.Monoid.Concat(
		R.LimitRetries(5),
		R.ConstantDelay(10*time.Millisecond),
	)
	result := Retrying(policy, action, check)

	result(env)()

	assert.Equal(t, 3, len(statuses), "Should have 3 status records")

	// Check first attempt
	assert.Equal(t, uint(0), statuses[0].IterNumber, "First attempt should be iteration 0")
	assert.True(t, O.IsNone(statuses[0].PreviousDelay), "First attempt should have no previous delay")

	// Check second attempt
	assert.Equal(t, uint(1), statuses[1].IterNumber, "Second attempt should be iteration 1")
	assert.True(t, O.IsSome(statuses[1].PreviousDelay), "Second attempt should have previous delay")

	// Check third attempt
	assert.Equal(t, uint(2), statuses[2].IterNumber, "Third attempt should be iteration 2")
	assert.True(t, O.IsSome(statuses[2].PreviousDelay), "Third attempt should have previous delay")
}

// TestRetrying_ConstantDelay tests retry with constant delay between attempts
func TestRetrying_ConstantDelay(t *testing.T) {
	attempts := 0
	delay := 50 * time.Millisecond
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, bool] {
		return func(e RetryTestEnv) G.IO[bool] {
			return G.Of(func() bool {
				attempts++
				return attempts >= 3
			}())
		}
	}
	check := func(value bool) bool {
		return !value // Retry while false
	}

	policy := R.Monoid.Concat(
		R.LimitRetries(5),
		R.ConstantDelay(delay),
	)

	start := time.Now()
	result := Retrying(policy, action, check)
	result(env)()
	elapsed := time.Since(start)

	assert.Equal(t, 3, attempts)
	// Should have 2 delays (between attempt 1-2 and 2-3)
	expectedMinDelay := 2 * delay
	assert.GreaterOrEqual(t, elapsed, expectedMinDelay, "Should wait at least 2 delays")
}

// TestRetrying_ExponentialBackoff tests retry with exponential backoff
func TestRetrying_ExponentialBackoff(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				return attempts
			}())
		}
	}
	check := func(value int) bool {
		return value < 4 // Retry until 4th attempt
	}

	baseDelay := 10 * time.Millisecond
	policy := R.Monoid.Concat(
		R.LimitRetries(5),
		R.ExponentialBackoff(baseDelay),
	)

	start := time.Now()
	result := Retrying(policy, action, check)
	result(env)()
	elapsed := time.Since(start)

	assert.Equal(t, 4, attempts)
	// Exponential backoff: 10ms, 20ms, 40ms = 70ms minimum
	expectedMinDelay := 70 * time.Millisecond
	assert.GreaterOrEqual(t, elapsed, expectedMinDelay, "Should wait with exponential backoff")
}

// TestRetrying_CapDelay tests that delay capping works correctly
func TestRetrying_CapDelay(t *testing.T) {
	var delays []time.Duration
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				if O.IsSome(status.PreviousDelay) {
					delay, _ := O.Unwrap(status.PreviousDelay)
					delays = append(delays, delay)
				}
				return len(delays)
			}())
		}
	}
	check := func(value int) bool {
		return value < 5 // Do 5 retries
	}

	maxDelay := 50 * time.Millisecond
	policy := R.Monoid.Concat(
		R.LimitRetries(10),
		R.CapDelay(maxDelay, R.ExponentialBackoff(10*time.Millisecond)),
	)

	result := Retrying(policy, action, check)
	result(env)()

	// All delays should be capped at maxDelay
	for i, delay := range delays {
		assert.LessOrEqual(t, delay, maxDelay,
			"Delay %d should be capped at %v, got %v", i, maxDelay, delay)
	}
}

// TestRetrying_EnvironmentAccess tests that the action can access the environment
func TestRetrying_EnvironmentAccess(t *testing.T) {
	env := RetryTestEnv{
		Config:     "production",
		MaxRetries: 3,
	}

	var configValues []string
	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				configValues = append(configValues, e.Config)
				return int(status.IterNumber)
			}())
		}
	}
	check := func(value int) bool {
		return value < 2 // Retry twice
	}

	policy := R.LimitRetries(5)
	result := Retrying(policy, action, check)
	result(env)()

	assert.Equal(t, 3, len(configValues), "Should have accessed environment 3 times")
	for i, config := range configValues {
		assert.Equal(t, "production", config, "Attempt %d should have correct config", i)
	}
}

// TestRetrying_WithLogging tests that retry status can be used for logging via environment
func TestRetrying_WithLogging(t *testing.T) {
	logs := []string{}
	env := RetryTestEnv{
		Logger: func(msg string) {
			logs = append(logs, msg)
		},
		Config: "test",
	}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				e.Logger(fmt.Sprintf("Attempt %d, cumulative delay: %v",
					status.IterNumber, status.CumulativeDelay))
				return int(status.IterNumber)
			}())
		}
	}
	check := func(value int) bool {
		return value < 3
	}

	policy := R.Monoid.Concat(
		R.LimitRetries(5),
		R.ConstantDelay(10*time.Millisecond),
	)

	result := Retrying(policy, action, check)
	result(env)()

	assert.Equal(t, 4, len(logs), "Should have 4 log entries")
	assert.Contains(t, logs[0], "Attempt 0")
	assert.Contains(t, logs[1], "Attempt 1")
	assert.Contains(t, logs[2], "Attempt 2")
	assert.Contains(t, logs[3], "Attempt 3")
}

// TestRetrying_PredicateBasedRetry tests retry based on result value
func TestRetrying_PredicateBasedRetry(t *testing.T) {
	values := []int{1, 2, 3, 4, 5}
	index := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				val := values[index]
				index++
				return val
			}())
		}
	}
	check := func(value int) bool {
		return value < 5 // Retry until we get 5
	}

	policy := R.LimitRetries(10)
	result := Retrying(policy, action, check)

	assert.Equal(t, 5, result(env)())
	assert.Equal(t, 5, index, "Should have tried 5 times")
}

// TestRetrying_NoRetryOnSuccess tests that successful operations don't retry
func TestRetrying_NoRetryOnSuccess(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, string] {
		return func(e RetryTestEnv) G.IO[string] {
			return G.Of(func() string {
				attempts++
				return "success"
			}())
		}
	}
	check := func(value string) bool {
		return value != "success" // Don't retry on success
	}

	policy := R.LimitRetries(5)
	result := Retrying(policy, action, check)

	assert.Equal(t, "success", result(env)())
	assert.Equal(t, 1, attempts, "Should only execute once on immediate success")
}

// TestRetrying_ZeroRetries tests behavior with zero retries allowed
func TestRetrying_ZeroRetries(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				return attempts
			}())
		}
	}
	check := func(value int) bool {
		return true // Always want to retry
	}

	policy := R.LimitRetries(0)
	result := Retrying(policy, action, check)

	assert.Equal(t, 1, result(env)())
	assert.Equal(t, 1, attempts, "Should execute once even with 0 retries")
}

// TestRetrying_CumulativeDelay tests that cumulative delay is tracked correctly
func TestRetrying_CumulativeDelay(t *testing.T) {
	var cumulativeDelays []time.Duration
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				cumulativeDelays = append(cumulativeDelays, status.CumulativeDelay)
				return len(cumulativeDelays)
			}())
		}
	}
	check := func(value int) bool {
		return value < 4 // Do 4 attempts
	}

	delay := 20 * time.Millisecond
	policy := R.Monoid.Concat(
		R.LimitRetries(5),
		R.ConstantDelay(delay),
	)

	result := Retrying(policy, action, check)
	result(env)()

	assert.Equal(t, 4, len(cumulativeDelays))
	assert.Equal(t, time.Duration(0), cumulativeDelays[0], "First attempt should have 0 cumulative delay")

	// Each subsequent attempt should have increasing cumulative delay
	for i := 1; i < len(cumulativeDelays); i++ {
		assert.Greater(t, cumulativeDelays[i], cumulativeDelays[i-1],
			"Cumulative delay should increase with each retry")
	}
}

// TestRetrying_ComplexPredicate tests retry with complex success criteria
func TestRetrying_ComplexPredicate(t *testing.T) {
	type Result struct {
		StatusCode int
		Body       string
	}

	results := []Result{
		{StatusCode: 500, Body: "error"},
		{StatusCode: 503, Body: "unavailable"},
		{StatusCode: 200, Body: "success"},
	}
	index := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, Result] {
		return func(e RetryTestEnv) G.IO[Result] {
			return G.Of(func() Result {
				result := results[index]
				index++
				return result
			}())
		}
	}
	check := func(r Result) bool {
		// Retry on server errors (5xx)
		return r.StatusCode >= 500
	}

	policy := R.LimitRetries(5)
	result := Retrying(policy, action, check)

	finalResult := result(env)()
	assert.Equal(t, 200, finalResult.StatusCode)
	assert.Equal(t, "success", finalResult.Body)
	assert.Equal(t, 3, index, "Should have tried 3 times")
}

// TestRetrying_PolicyReturnsNone tests behavior when policy returns None immediately
func TestRetrying_PolicyReturnsNone(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				return attempts
			}())
		}
	}
	check := func(value int) bool {
		return true // Always want to retry
	}

	// Policy that never allows retries
	policy := func(status R.RetryStatus) O.Option[time.Duration] {
		return O.None[time.Duration]()
	}

	result := Retrying(policy, action, check)
	result(env)()

	assert.Equal(t, 1, attempts, "Should only execute once when policy returns None")
}

// TestRetrying_CombinedPolicies tests combining multiple retry policies
func TestRetrying_CombinedPolicies(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, string] {
		return func(e RetryTestEnv) G.IO[string] {
			return G.Of(func() string {
				attempts++
				return fmt.Sprintf("attempt_%d", attempts)
			}())
		}
	}
	check := func(value string) bool {
		return true // Always retry
	}

	// Combine limit and exponential backoff with cap
	policy := R.Monoid.Concat(
		R.LimitRetries(3),
		R.CapDelay(100*time.Millisecond, R.ExponentialBackoff(20*time.Millisecond)),
	)

	result := Retrying(policy, action, check)
	result(env)()

	assert.Equal(t, 4, attempts, "Should respect the retry limit")
}

// TestRetrying_DifferentEnvironments tests that different environments can be used
func TestRetrying_DifferentEnvironments(t *testing.T) {
	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, string] {
		return func(e RetryTestEnv) G.IO[string] {
			return G.Of(func() string {
				return e.Config
			}())
		}
	}
	check := func(value string) bool {
		return false // Never retry
	}

	policy := R.LimitRetries(5)
	result := Retrying(policy, action, check)

	env1 := RetryTestEnv{Config: "dev"}
	env2 := RetryTestEnv{Config: "prod"}

	assert.Equal(t, "dev", result(env1)())
	assert.Equal(t, "prod", result(env2)())
}

// TestRetrying_EnvironmentBasedRetryLogic tests retry logic that depends on environment
func TestRetrying_EnvironmentBasedRetryLogic(t *testing.T) {
	attempts := 0
	env := RetryTestEnv{
		Config:     "test",
		MaxRetries: 3,
	}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				// Use environment to determine behavior
				if attempts >= e.MaxRetries {
					return 100 // Success value
				}
				return attempts
			}())
		}
	}
	check := func(value int) bool {
		return value < 100 // Retry until we get 100
	}

	policy := R.LimitRetries(10)
	result := Retrying(policy, action, check)

	finalResult := result(env)()
	assert.Equal(t, 100, finalResult)
	assert.Equal(t, 3, attempts, "Should retry based on environment MaxRetries")
}

// TestRetrying_LongRunningRetry tests retry over multiple attempts with realistic delays
func TestRetrying_LongRunningRetry(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping long-running test in short mode")
	}

	attempts := 0
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, bool] {
		return func(e RetryTestEnv) G.IO[bool] {
			return G.Of(func() bool {
				attempts++
				// Succeed on 5th attempt
				return attempts >= 5
			}())
		}
	}
	check := func(value bool) bool {
		return !value
	}

	policy := R.Monoid.Concat(
		R.LimitRetries(10),
		R.ConstantDelay(100*time.Millisecond),
	)

	start := time.Now()
	result := Retrying(policy, action, check)
	result(env)()
	elapsed := time.Since(start)

	assert.Equal(t, 5, attempts)
	// Should have 4 delays of 100ms each = 400ms minimum
	expectedMinDelay := 400 * time.Millisecond
	assert.GreaterOrEqual(t, elapsed, expectedMinDelay)
}

// TestRetrying_StackSafety tests that deep recursion doesn't cause stack overflow
func TestRetrying_StackSafety(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stack safety test in short mode")
	}

	attempts := 0
	maxAttempts := 1000
	env := RetryTestEnv{Config: "test"}

	action := func(status R.RetryStatus) ReaderIO[RetryTestEnv, int] {
		return func(e RetryTestEnv) G.IO[int] {
			return G.Of(func() int {
				attempts++
				return attempts
			}())
		}
	}
	check := func(value int) bool {
		return value < maxAttempts
	}

	policy := R.LimitRetries(uint(maxAttempts))
	result := Retrying(policy, action, check)

	finalResult := result(env)()
	assert.Equal(t, maxAttempts, finalResult)
	assert.Equal(t, maxAttempts, attempts, "Should handle many retries without stack overflow")
}
