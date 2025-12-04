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

package retry

import (
	"testing"
	"time"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestIterNumber(t *testing.T) {
	t.Run("returns iteration number from status", func(t *testing.T) {
		status := RetryStatus{
			IterNumber:      5,
			CumulativeDelay: 1000 * time.Millisecond,
			PreviousDelay:   O.Some(500 * time.Millisecond),
		}
		assert.Equal(t, uint(5), IterNumber(status))
	})

	t.Run("returns zero for default status", func(t *testing.T) {
		assert.Equal(t, uint(0), IterNumber(DefaultRetryStatus))
	})
}

func TestLimitRetries(t *testing.T) {
	t.Run("allows retries up to limit", func(t *testing.T) {
		policy := LimitRetries(3)

		// First attempt (iter 0)
		result := policy(RetryStatus{IterNumber: 0, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		assert.Equal(t, time.Duration(0), O.GetOrElse(func() time.Duration { return -1 })(result))

		// Second attempt (iter 1)
		result = policy(RetryStatus{IterNumber: 1, CumulativeDelay: 0, PreviousDelay: O.Some(time.Duration(0))})
		assert.True(t, O.IsSome(result))

		// Third attempt (iter 2)
		result = policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.Some(time.Duration(0))})
		assert.True(t, O.IsSome(result))

		// Fourth attempt (iter 3) - should fail
		result = policy(RetryStatus{IterNumber: 3, CumulativeDelay: 0, PreviousDelay: O.Some(time.Duration(0))})
		assert.True(t, O.IsNone(result))
	})

	t.Run("returns None when limit is reached", func(t *testing.T) {
		policy := LimitRetries(2)
		result := policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsNone(result))
	})

	t.Run("zero limit prevents all retries", func(t *testing.T) {
		policy := LimitRetries(0)
		result := policy(DefaultRetryStatus)
		assert.True(t, O.IsNone(result))
	})
}

func TestConstantDelay(t *testing.T) {
	t.Run("returns constant delay for any status", func(t *testing.T) {
		delay := 500 * time.Millisecond
		policy := ConstantDelay(delay)

		// First attempt
		result := policy(DefaultRetryStatus)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, delay, O.GetOrElse(func() time.Duration { return 0 })(result))

		// Later attempts
		result = policy(RetryStatus{IterNumber: 10, CumulativeDelay: 5 * time.Second, PreviousDelay: O.Some(delay)})
		assert.True(t, O.IsSome(result))
		assert.Equal(t, delay, O.GetOrElse(func() time.Duration { return 0 })(result))
	})

	t.Run("works with zero delay", func(t *testing.T) {
		policy := ConstantDelay(0)
		result := policy(DefaultRetryStatus)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, time.Duration(0), O.GetOrElse(func() time.Duration { return -1 })(result))
	})
}

func TestExponentialBackoff(t *testing.T) {
	t.Run("grows delay exponentially", func(t *testing.T) {
		baseDelay := 100 * time.Millisecond
		policy := ExponentialBackoff(baseDelay)

		testCases := []struct {
			iterNumber    uint
			expectedDelay time.Duration
		}{
			{0, 100 * time.Millisecond},  // 100 * 2^0
			{1, 200 * time.Millisecond},  // 100 * 2^1
			{2, 400 * time.Millisecond},  // 100 * 2^2
			{3, 800 * time.Millisecond},  // 100 * 2^3
			{4, 1600 * time.Millisecond}, // 100 * 2^4
		}

		for _, tc := range testCases {
			status := RetryStatus{
				IterNumber:      tc.iterNumber,
				CumulativeDelay: 0,
				PreviousDelay:   O.None[time.Duration](),
			}
			result := policy(status)
			assert.True(t, O.IsSome(result))
			assert.Equal(t, tc.expectedDelay, O.GetOrElse(func() time.Duration { return 0 })(result))
		}
	})

	t.Run("never returns None", func(t *testing.T) {
		policy := ExponentialBackoff(100 * time.Millisecond)
		// Even with very high iteration numbers
		result := policy(RetryStatus{IterNumber: 100, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
	})
}

func TestCapDelay(t *testing.T) {
	t.Run("caps delay at maximum", func(t *testing.T) {
		maxDelay := 1 * time.Second
		basePolicy := ExponentialBackoff(100 * time.Millisecond)
		policy := CapDelay(maxDelay, basePolicy)

		// Small iteration - should not be capped
		result := policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		delay := O.GetOrElse(func() time.Duration { return 0 })(result)
		assert.Equal(t, 400*time.Millisecond, delay)

		// Large iteration - should be capped
		result = policy(RetryStatus{IterNumber: 10, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		delay = O.GetOrElse(func() time.Duration { return 0 })(result)
		assert.Equal(t, maxDelay, delay)
	})

	t.Run("preserves None from underlying policy", func(t *testing.T) {
		maxDelay := 1 * time.Second
		basePolicy := LimitRetries(2)
		policy := CapDelay(maxDelay, basePolicy)

		// Within limit
		result := policy(RetryStatus{IterNumber: 1, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))

		// Beyond limit
		result = policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsNone(result))
	})
}

func TestApplyPolicy(t *testing.T) {
	t.Run("increments iteration number", func(t *testing.T) {
		policy := ConstantDelay(100 * time.Millisecond)
		status := DefaultRetryStatus

		newStatus := ApplyPolicy(policy, status)
		assert.Equal(t, uint(1), newStatus.IterNumber)

		newStatus = ApplyPolicy(policy, newStatus)
		assert.Equal(t, uint(2), newStatus.IterNumber)
	})

	t.Run("accumulates delay", func(t *testing.T) {
		policy := ConstantDelay(100 * time.Millisecond)
		status := DefaultRetryStatus

		newStatus := ApplyPolicy(policy, status)
		assert.Equal(t, 100*time.Millisecond, newStatus.CumulativeDelay)

		newStatus = ApplyPolicy(policy, newStatus)
		assert.Equal(t, 200*time.Millisecond, newStatus.CumulativeDelay)

		newStatus = ApplyPolicy(policy, newStatus)
		assert.Equal(t, 300*time.Millisecond, newStatus.CumulativeDelay)
	})

	t.Run("stores previous delay", func(t *testing.T) {
		policy := ExponentialBackoff(100 * time.Millisecond)
		status := DefaultRetryStatus

		newStatus := ApplyPolicy(policy, status)
		assert.True(t, O.IsSome(newStatus.PreviousDelay))
		assert.Equal(t, 100*time.Millisecond, O.GetOrElse(func() time.Duration { return 0 })(newStatus.PreviousDelay))

		newStatus = ApplyPolicy(policy, newStatus)
		assert.True(t, O.IsSome(newStatus.PreviousDelay))
		assert.Equal(t, 200*time.Millisecond, O.GetOrElse(func() time.Duration { return 0 })(newStatus.PreviousDelay))
	})

	t.Run("handles None from policy", func(t *testing.T) {
		policy := LimitRetries(1)
		status := RetryStatus{IterNumber: 1, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()}

		newStatus := ApplyPolicy(policy, status)
		assert.Equal(t, uint(2), newStatus.IterNumber)
		assert.True(t, O.IsNone(newStatus.PreviousDelay))
		assert.Equal(t, time.Duration(0), newStatus.CumulativeDelay)
	})
}

func TestMonoidCombination(t *testing.T) {
	t.Run("combines limit with constant delay", func(t *testing.T) {
		policy := Monoid.Concat(
			LimitRetries(3),
			ConstantDelay(500*time.Millisecond),
		)

		// Within limit - should return delay
		result := policy(RetryStatus{IterNumber: 0, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 500*time.Millisecond, O.GetOrElse(func() time.Duration { return 0 })(result))

		result = policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))

		// Beyond limit - should return None
		result = policy(RetryStatus{IterNumber: 3, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsNone(result))
	})

	t.Run("takes maximum delay when both return Some", func(t *testing.T) {
		policy := Monoid.Concat(
			ConstantDelay(100*time.Millisecond),
			ConstantDelay(500*time.Millisecond),
		)

		result := policy(DefaultRetryStatus)
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 500*time.Millisecond, O.GetOrElse(func() time.Duration { return 0 })(result))
	})

	t.Run("returns None if either policy returns None", func(t *testing.T) {
		policy := Monoid.Concat(
			LimitRetries(2),
			ConstantDelay(500*time.Millisecond),
		)

		// Beyond limit
		result := policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsNone(result))
	})

	t.Run("complex policy combination", func(t *testing.T) {
		// Exponential backoff, capped at 5 seconds, up to 10 retries
		policy := Monoid.Concat(
			LimitRetries(10),
			CapDelay(5*time.Second, ExponentialBackoff(100*time.Millisecond)),
		)

		// Early iterations - exponential growth
		result := policy(RetryStatus{IterNumber: 2, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 400*time.Millisecond, O.GetOrElse(func() time.Duration { return 0 })(result))

		// Later iterations - capped
		result = policy(RetryStatus{IterNumber: 8, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsSome(result))
		assert.Equal(t, 5*time.Second, O.GetOrElse(func() time.Duration { return 0 })(result))

		// Beyond limit
		result = policy(RetryStatus{IterNumber: 10, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()})
		assert.True(t, O.IsNone(result))
	})
}

func TestDefaultRetryStatus(t *testing.T) {
	t.Run("has correct initial values", func(t *testing.T) {
		assert.Equal(t, uint(0), DefaultRetryStatus.IterNumber)
		assert.Equal(t, time.Duration(0), DefaultRetryStatus.CumulativeDelay)
		assert.True(t, O.IsNone(DefaultRetryStatus.PreviousDelay))
	})
}

func TestRetryStatusProgression(t *testing.T) {
	t.Run("simulates full retry sequence", func(t *testing.T) {
		policy := Monoid.Concat(
			LimitRetries(3),
			ExponentialBackoff(100*time.Millisecond),
		)

		status := DefaultRetryStatus

		// First retry
		status = ApplyPolicy(policy, status)
		assert.Equal(t, uint(1), status.IterNumber)
		assert.Equal(t, 100*time.Millisecond, status.CumulativeDelay)
		assert.True(t, O.IsSome(status.PreviousDelay))

		// Second retry
		status = ApplyPolicy(policy, status)
		assert.Equal(t, uint(2), status.IterNumber)
		assert.Equal(t, 300*time.Millisecond, status.CumulativeDelay) // 100 + 200

		// Third retry
		status = ApplyPolicy(policy, status)
		assert.Equal(t, uint(3), status.IterNumber)
		assert.Equal(t, 700*time.Millisecond, status.CumulativeDelay) // 100 + 200 + 400

		// Fourth attempt - should stop
		status = ApplyPolicy(policy, status)
		assert.Equal(t, uint(4), status.IterNumber)
		assert.True(t, O.IsNone(status.PreviousDelay))
		assert.Equal(t, 700*time.Millisecond, status.CumulativeDelay) // No additional delay
	})
}
