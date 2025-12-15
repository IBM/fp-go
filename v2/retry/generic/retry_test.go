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

package generic

import (
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/stretchr/testify/assert"
)

// Mock monad for testing - represents a simple state monad that tracks execution
type TestMonad[A any] struct {
	value       A
	delays      []time.Duration
	actionCalls int
}

func newTestMonad[A any](value A) *TestMonad[A] {
	return &TestMonad[A]{
		value:       value,
		delays:      []time.Duration{},
		actionCalls: 0,
	}
}

// Monad operations for TestMonad
func testMonadOf[A any](value A) *TestMonad[A] {
	return newTestMonad(value)
}

func testMonadChain[A, B any](f func(A) *TestMonad[B]) func(*TestMonad[A]) *TestMonad[B] {
	return func(m *TestMonad[A]) *TestMonad[B] {
		result := f(m.value)
		// Preserve delays and action calls
		result.delays = append(m.delays, result.delays...)
		result.actionCalls = m.actionCalls + result.actionCalls
		return result
	}
}

func testMonadChainStatus[A any](f func(R.RetryStatus) *TestMonad[A]) func(*TestMonad[R.RetryStatus]) *TestMonad[A] {
	return func(m *TestMonad[R.RetryStatus]) *TestMonad[A] {
		result := f(m.value)
		result.delays = append(m.delays, result.delays...)
		result.actionCalls = m.actionCalls + result.actionCalls
		return result
	}
}

func testMonadDelay[A any](duration time.Duration) func(*TestMonad[A]) *TestMonad[A] {
	return func(m *TestMonad[A]) *TestMonad[A] {
		m.delays = append(m.delays, duration)
		return m
	}
}

func TestApplyAndDelay(t *testing.T) {
	t.Run("applies policy and delays when policy returns Some", func(t *testing.T) {
		policy := R.ConstantDelay(100 * time.Millisecond)
		status := R.DefaultRetryStatus

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus])
		result := applyDelay(policy, status)

		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 100*time.Millisecond, result.delays[0])

		// Check that status was updated
		assert.Equal(t, uint(1), result.value.IterNumber)
	})

	t.Run("does not delay when policy returns None", func(t *testing.T) {
		policy := R.LimitRetries(0) // Always returns None
		status := R.DefaultRetryStatus

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus])
		result := applyDelay(policy, status)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.delays))

		// Check that status was still updated
		assert.Equal(t, uint(1), result.value.IterNumber)
	})

	t.Run("accumulates delays correctly", func(t *testing.T) {
		policy := R.ExponentialBackoff(100 * time.Millisecond)
		status := R.DefaultRetryStatus

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus])

		// First application
		result := applyDelay(policy, status)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 100*time.Millisecond, result.delays[0])

		// Second application
		newStatus := result.value
		result = applyDelay(policy, newStatus)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 200*time.Millisecond, result.delays[0])
	})
}

func TestRetrying(t *testing.T) {
	t.Run("succeeds on first try when check returns false", func(t *testing.T) {
		policy := R.LimitRetries(3)
		callCount := 0

		action := func(status R.RetryStatus) *TestMonad[string] {
			callCount++
			result := newTestMonad("success")
			result.actionCalls = 1
			return result
		}

		check := func(value string) bool {
			return false // Don't retry
		}

		result := Retrying(
			testMonadChain[string, string],
			testMonadChainStatus[string],
			testMonadOf[string],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, "success", result.value)
		assert.Equal(t, 1, callCount)
		assert.Equal(t, 0, len(result.delays))
	})

	t.Run("retries until success", func(t *testing.T) {
		policy := R.Monoid.Concat(
			R.LimitRetries(5),
			R.ConstantDelay(100*time.Millisecond),
		)
		callCount := 0

		action := func(status R.RetryStatus) *TestMonad[int] {
			callCount++
			result := newTestMonad(callCount)
			result.actionCalls = 1
			return result
		}

		check := func(value int) bool {
			return value < 3 // Retry until we reach 3
		}

		result := Retrying(
			testMonadChain[int, int],
			testMonadChainStatus[int],
			testMonadOf[int],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, 3, result.value)
		assert.Equal(t, 3, callCount)
		assert.Equal(t, 2, len(result.delays)) // 2 retries = 2 delays
		assert.Equal(t, 100*time.Millisecond, result.delays[0])
		assert.Equal(t, 100*time.Millisecond, result.delays[1])
	})

	t.Run("stops when retry limit is reached", func(t *testing.T) {
		policy := R.Monoid.Concat(
			R.LimitRetries(3),
			R.ConstantDelay(50*time.Millisecond),
		)
		callCount := 0

		action := func(status R.RetryStatus) *TestMonad[string] {
			callCount++
			result := newTestMonad("failure")
			result.actionCalls = 1
			return result
		}

		check := func(value string) bool {
			return true // Always retry
		}

		result := Retrying(
			testMonadChain[string, string],
			testMonadChainStatus[string],
			testMonadOf[string],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, "failure", result.value)
		assert.Equal(t, 4, callCount) // Initial + 3 retries
		assert.Equal(t, 3, len(result.delays))
	})

	t.Run("respects exponential backoff", func(t *testing.T) {
		policy := R.Monoid.Concat(
			R.LimitRetries(4),
			R.ExponentialBackoff(100*time.Millisecond),
		)
		callCount := 0

		action := func(status R.RetryStatus) *TestMonad[int] {
			callCount++
			result := newTestMonad(callCount)
			result.actionCalls = 1
			return result
		}

		check := func(value int) bool {
			return value < 4 // Retry until we reach 4
		}

		result := Retrying(
			testMonadChain[int, int],
			testMonadChainStatus[int],
			testMonadOf[int],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, 4, result.value)
		assert.Equal(t, 4, callCount)
		assert.Equal(t, 3, len(result.delays))
		// Verify exponential growth: 100ms, 200ms, 400ms
		assert.Equal(t, 100*time.Millisecond, result.delays[0])
		assert.Equal(t, 200*time.Millisecond, result.delays[1])
		assert.Equal(t, 400*time.Millisecond, result.delays[2])
	})

	t.Run("respects capped delay", func(t *testing.T) {
		policy := R.Monoid.Concat(
			R.LimitRetries(5),
			R.CapDelay(300*time.Millisecond, R.ExponentialBackoff(100*time.Millisecond)),
		)
		callCount := 0

		action := func(status R.RetryStatus) *TestMonad[int] {
			callCount++
			result := newTestMonad(callCount)
			result.actionCalls = 1
			return result
		}

		check := func(value int) bool {
			return value < 5
		}

		result := Retrying(
			testMonadChain[int, int],
			testMonadChainStatus[int],
			testMonadOf[int],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, 5, result.value)
		assert.Equal(t, 5, callCount)
		assert.Equal(t, 4, len(result.delays))
		// Verify capped growth: 100ms, 200ms, 300ms (capped), 300ms (capped)
		assert.Equal(t, 100*time.Millisecond, result.delays[0])
		assert.Equal(t, 200*time.Millisecond, result.delays[1])
		assert.Equal(t, 300*time.Millisecond, result.delays[2])
		assert.Equal(t, 300*time.Millisecond, result.delays[3])
	})

	t.Run("handles immediate success without delay", func(t *testing.T) {
		policy := R.LimitRetries(3)

		action := func(status R.RetryStatus) *TestMonad[int] {
			result := newTestMonad(42)
			result.actionCalls = 1
			return result
		}

		check := func(value int) bool {
			return false // Success on first try
		}

		result := Retrying(
			testMonadChain[int, int],
			testMonadChainStatus[int],
			testMonadOf[int],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.Equal(t, 42, result.value)
		assert.Equal(t, 0, len(result.delays))
	})
}

func TestRetryingWithOption(t *testing.T) {
	t.Run("works with Option monad conceptually", func(t *testing.T) {
		// This test demonstrates how Retrying would work with Option
		// We use a simplified version since we can't easily test with real Option monad

		policy := R.Monoid.Concat(
			R.LimitRetries(3),
			R.ConstantDelay(100*time.Millisecond),
		)

		attempts := 0
		action := func(status R.RetryStatus) *TestMonad[O.Option[string]] {
			attempts++
			// Simulate Option: None on first 2 attempts, Some on 3rd
			if attempts < 3 {
				result := newTestMonad(O.None[string]())
				result.actionCalls = 1
				return result
			}
			result := newTestMonad(O.Some("success"))
			result.actionCalls = 1
			return result
		}

		check := O.IsNone[string] // Retry if None

		result := Retrying(
			testMonadChain[O.Option[string], O.Option[string]],
			testMonadChainStatus[O.Option[string]],
			testMonadOf[O.Option[string]],
			testMonadOf[R.RetryStatus],
			testMonadDelay[R.RetryStatus],
			policy,
			action,
			check,
		)

		assert.True(t, O.IsSome(result.value))
		assert.Equal(t, "success", O.GetOrElse(F.Constant(""))(result.value))
		assert.Equal(t, 3, attempts)
		assert.Equal(t, 2, len(result.delays))
	})
}
