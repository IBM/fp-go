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

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus], policy)
		result := applyDelay(status)

		assert.NotNil(t, result)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 100*time.Millisecond, result.delays[0])

		// Check that status was updated
		assert.Equal(t, uint(1), result.value.IterNumber)
	})

	t.Run("does not delay when policy returns None", func(t *testing.T) {
		policy := R.LimitRetries(0) // Always returns None
		status := R.DefaultRetryStatus

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus], policy)
		result := applyDelay(status)

		assert.NotNil(t, result)
		assert.Equal(t, 0, len(result.delays))

		// Check that status was still updated
		assert.Equal(t, uint(1), result.value.IterNumber)
	})

	t.Run("accumulates delays correctly", func(t *testing.T) {
		policy := R.ExponentialBackoff(100 * time.Millisecond)
		status := R.DefaultRetryStatus

		applyDelay := applyAndDelay(testMonadOf[R.RetryStatus], testMonadDelay[R.RetryStatus], policy)

		// First application
		result := applyDelay(status)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 100*time.Millisecond, result.delays[0])

		// Second application
		newStatus := result.value
		result = applyDelay(newStatus)
		assert.Equal(t, 1, len(result.delays))
		assert.Equal(t, 200*time.Millisecond, result.delays[0])
	})
}
