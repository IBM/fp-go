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

package effect

import (
	"errors"
	"testing"
	"time"

	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

func TestRetrying(t *testing.T) {
	t.Run("succeeds on first attempt", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(3)

		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				return Of[TestContext]("success")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res) // retry on error
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 1, attemptCount)
	})

	t.Run("retries on failure and eventually succeeds", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(5)

		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				if attemptCount < 3 {
					return Fail[TestContext, string](errors.New("temporary error"))
				}
				return Of[TestContext]("success after retries")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res) // retry on error
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "success after retries", result)
		assert.Equal(t, 3, attemptCount)
	})

	t.Run("exhausts retry limit", func(t *testing.T) {
		attemptCount := 0
		maxRetries := uint(3)
		policy := retry.LimitRetries(maxRetries)

		eff := Retrying(
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				return Fail[TestContext, string](errors.New("persistent error"))
			},
			func(res Result[string]) bool {
				return result.IsLeft(res) // retry on error
			},
		)

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, int(maxRetries+1), attemptCount) // initial attempt + retries
	})

	t.Run("does not retry on success", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(5)

		eff := Retrying[TestContext, int](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, int] {
				attemptCount++
				return Of[TestContext](42)
			},
			func(res Result[int]) bool {
				return result.IsLeft(res) // retry on error
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
		assert.Equal(t, 1, attemptCount)
	})

	t.Run("uses custom retry predicate", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(5)

		eff := Retrying[TestContext, int](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, int] {
				attemptCount++
				return Of[TestContext](attemptCount * 10)
			},
			func(res Result[int]) bool {
				// Retry if value is less than 30
				if result.IsRight(res) {
					val, _ := result.Unwrap(res)
					return val < 30
				}
				return true
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 30, result)
		assert.Equal(t, 3, attemptCount)
	})

	t.Run("tracks retry status", func(t *testing.T) {
		var statuses []retry.RetryStatus
		policy := retry.LimitRetries(3)

		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				statuses = append(statuses, status)
				if len(statuses) < 3 {
					return Fail[TestContext, string](errors.New("retry"))
				}
				return Of[TestContext]("done")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res)
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "done", result)
		assert.Len(t, statuses, 3)
		// First attempt has iteration 0
		assert.Equal(t, uint(0), statuses[0].IterNumber)
		assert.Equal(t, uint(1), statuses[1].IterNumber)
		assert.Equal(t, uint(2), statuses[2].IterNumber)
	})

	t.Run("works with exponential backoff", func(t *testing.T) {
		attemptCount := 0
		policy := retry.Monoid.Concat(
			retry.LimitRetries(3),
			retry.ExponentialBackoff(10*time.Millisecond),
		)

		startTime := time.Now()
		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				if attemptCount < 3 {
					return Fail[TestContext, string](errors.New("retry"))
				}
				return Of[TestContext]("success")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res)
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})
		elapsed := time.Since(startTime)

		assert.NoError(t, err)
		assert.Equal(t, "success", result)
		assert.Equal(t, 3, attemptCount)
		// Should have some delay due to backoff
		assert.Greater(t, elapsed, 10*time.Millisecond)
	})

	t.Run("combines with other effect operations", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(3)

		eff := Map[TestContext](func(s string) string {
			return "mapped: " + s
		})(Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				if attemptCount < 2 {
					return Fail[TestContext, string](errors.New("retry"))
				}
				return Of[TestContext]("success")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res)
			},
		))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "mapped: success", result)
		assert.Equal(t, 2, attemptCount)
	})

	t.Run("retries with different error types", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(5)
		errors := []error{
			errors.New("error 1"),
			errors.New("error 2"),
			errors.New("error 3"),
		}

		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				if attemptCount < len(errors) {
					err := errors[attemptCount]
					attemptCount++
					return Fail[TestContext, string](err)
				}
				attemptCount++
				return Of[TestContext]("finally succeeded")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res)
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "finally succeeded", result)
		assert.Equal(t, 4, attemptCount)
	})

	t.Run("no retry when predicate returns false", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(5)

		eff := Retrying(
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				return Fail[TestContext, string](errors.New("error"))
			},
			func(res Result[string]) bool {
				return false // never retry
			},
		)

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, 1, attemptCount) // only initial attempt
	})

	t.Run("retries with context access", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(3)
		ctx := TestContext{Value: "retry-context"}

		eff := Retrying[TestContext, string](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, string] {
				attemptCount++
				if attemptCount < 2 {
					return Fail[TestContext, string](errors.New("retry"))
				}
				return Of[TestContext]("success with context")
			},
			func(res Result[string]) bool {
				return result.IsLeft(res)
			},
		)

		result, err := runEffect(eff, ctx)

		assert.NoError(t, err)
		assert.Equal(t, "success with context", result)
		assert.Equal(t, 2, attemptCount)
	})
}

func TestRetryingWithComplexScenarios(t *testing.T) {
	t.Run("retry with state accumulation", func(t *testing.T) {
		type State struct {
			Attempts []int
			Value    string
		}

		policy := retry.LimitRetries(4)

		eff := Retrying[TestContext, State](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, State] {
				state := State{
					Attempts: make([]int, status.IterNumber+1),
					Value:    "attempt",
				}
				for i := uint(0); i <= status.IterNumber; i++ {
					state.Attempts[i] = int(i)
				}

				if status.IterNumber < 2 {
					return Fail[TestContext, State](errors.New("retry"))
				}
				return Of[TestContext](state)
			},
			func(res Result[State]) bool {
				return result.IsLeft(res)
			},
		)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "attempt", result.Value)
		assert.Equal(t, []int{0, 1, 2}, result.Attempts)
	})

	t.Run("retry with chain operations", func(t *testing.T) {
		attemptCount := 0
		policy := retry.LimitRetries(3)

		eff := Chain(func(x int) Effect[TestContext, string] {
			return Of[TestContext]("final: " + string(rune('0'+x)))
		})(Retrying[TestContext, int](
			policy,
			func(status retry.RetryStatus) Effect[TestContext, int] {
				attemptCount++
				if attemptCount < 2 {
					return Fail[TestContext, int](errors.New("retry"))
				}
				return Of[TestContext](attemptCount)
			},
			func(res Result[int]) bool {
				return result.IsLeft(res)
			},
		))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Contains(t, result, "final:")
	})
}
