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

// Package retry provides functional retry policies and combinators for implementing
// retry logic with configurable backoff strategies.
//
// This package offers a composable approach to retrying operations, allowing you to:
//   - Define retry policies that determine when and how long to wait between retries
//   - Combine multiple policies using monoid operations
//   - Implement various backoff strategies (constant, exponential, etc.)
//   - Limit the number of retries or cap the maximum delay
//
// # Basic Usage
//
// Create a simple retry policy that retries up to 3 times with exponential backoff:
//
//	policy := M.Concat(
//		LimitRetries(3),
//		ExponentialBackoff(100 * time.Millisecond),
//	)(Monoid)
//
// # Retry Policies
//
// A RetryPolicy is a function that takes a RetryStatus and returns an optional delay.
// If the policy returns None, retrying stops. If it returns Some(delay), the operation
// will be retried after the specified delay.
//
// # Combining Policies
//
// Policies can be combined using the Monoid instance. When combining policies:
//   - If either policy returns None, the combined policy returns None
//   - If both return a delay, the larger delay is used
//
// Example combining a retry limit with exponential backoff:
//
//	policy := M.Concat(
//		LimitRetries(5),
//		CapDelay(5*time.Second, ExponentialBackoff(100*time.Millisecond)),
//	)(Monoid)
package retry

import (
	"math"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/ord"
)

// RetryStatus tracks the current state of a retry operation.
// It contains information about the iteration number, cumulative delay,
// and the delay from the previous attempt.
type RetryStatus struct {
	// IterNumber is the iteration number, where 0 is the first try.
	// This increments by 1 for each retry attempt.
	IterNumber uint

	// CumulativeDelay is the total delay incurred so far from all retries.
	// This is the sum of all previous delays.
	CumulativeDelay time.Duration

	// PreviousDelay is the delay from the latest attempt.
	// This will always be None on the first run (IterNumber == 0).
	PreviousDelay Option[time.Duration]
}

// RetryPolicy is a function that takes a RetryStatus and possibly returns
// a delay duration. Iteration numbers start at zero and increase by one on
// each retry.
//
// A None return value from the policy indicates that the retry limit has been
// reached and no further retries should be attempted. A Some(duration) return
// value indicates that the operation should be retried after waiting for the
// specified duration.
//
// Example creating a custom policy:
//
//	// Retry up to 3 times with a fixed 1 second delay
//	customPolicy := func(status RetryStatus) Option[time.Duration] {
//		if status.IterNumber < 3 {
//			return O.Some(1 * time.Second)
//		}
//		return O.None[time.Duration]()
//	}
type RetryPolicy = func(RetryStatus) Option[time.Duration]

const emptyDuration = time.Duration(0)

var ordDuration = ord.FromStrictCompare[time.Duration]()

var IterNumberLens = L.MakeLensWithName(
	func(rs RetryStatus) uint { return rs.IterNumber },
	func(rs RetryStatus, iter uint) RetryStatus { rs.IterNumber = iter; return rs },
	"RetryStatus.IterNumber",
)

var CumulativeDelayLens = L.MakeLensWithName(
	func(rs RetryStatus) time.Duration { return rs.CumulativeDelay },
	func(rs RetryStatus, delay time.Duration) RetryStatus { rs.CumulativeDelay = delay; return rs },
	"RetryStatus.CumulativeDelay",
)

var PreviousDelayLens = L.MakeLensWithName(
	func(rs RetryStatus) Option[time.Duration] { return rs.PreviousDelay },
	func(rs RetryStatus, delay Option[time.Duration]) RetryStatus { rs.PreviousDelay = delay; return rs },
	"RetryStatus.PreviousDelay",
)

// IterNumber is an accessor function that extracts the iteration number
// from a RetryStatus. This is useful for functional composition.
//
// Example:
//
//	status := RetryStatus{IterNumber: 3, CumulativeDelay: 0, PreviousDelay: O.None[time.Duration]()}
//	iter := IterNumber(status) // returns 3
func IterNumber(rs RetryStatus) uint {
	return rs.IterNumber
}

// Monoid is the Monoid instance for RetryPolicy. You can collapse multiple
// retry strategies into one using the monoid's Concat operation.
//
// The semantics of this combination are:
//
//  1. If either policy returns None, the combined policy returns None.
//     This allows you to limit retries by combining with LimitRetries.
//
//  2. If both policies return a delay, the larger delay will be used.
//     This is natural when combining multiple policies to achieve a
//     certain effect, such as exponential backoff with a cap.
//
// Example combining policies:
//
//	// Retry up to 5 times with exponential backoff, capped at 10 seconds
//	policy := M.Concat(
//		M.Concat(
//			LimitRetries(5),
//			ExponentialBackoff(100*time.Millisecond),
//		)(Monoid),
//		CapDelay(10*time.Second, ConstantDelay(0)),
//	)(Monoid)
var Monoid = M.FunctionMonoid[RetryStatus](O.ApplicativeMonoid(M.MakeMonoid(
	ord.MaxSemigroup(ordDuration).Concat, emptyDuration)))

// LimitRetries creates a retry policy that retries immediately (with zero delay),
// but only up to i times. After i retries, the policy returns None, stopping
// further retry attempts.
//
// The iteration count starts at 0, so LimitRetries(3) will allow the initial
// attempt plus 3 retries (4 total attempts).
//
// Example:
//
//	// Allow up to 3 retries (4 total attempts)
//	policy := LimitRetries(3)
//
//	// Combine with a delay strategy
//	policyWithDelay := M.Concat(
//		LimitRetries(3),
//		ConstantDelay(1*time.Second),
//	)(Monoid)
func LimitRetries(i uint) RetryPolicy {
	return F.Flow3(
		IterNumber,
		O.FromPredicate(N.LessThan(i)),
		O.Map(F.Constant1[uint](emptyDuration)),
	)
}

// ConstantDelay creates a retry policy that always returns the same delay
// duration, allowing unlimited retries. This policy never returns None,
// so it should typically be combined with LimitRetries to prevent infinite retries.
//
// Example:
//
//	// Retry with a constant 500ms delay, up to 5 times
//	policy := M.Concat(
//		LimitRetries(5),
//		ConstantDelay(500*time.Millisecond),
//	)(Monoid)
func ConstantDelay(delay time.Duration) RetryPolicy {
	return F.Constant1[RetryStatus](O.Of(delay))
}

// CapDelay sets an upper bound on the delay returned by a retry policy.
// Any delay greater than maxDelay will be capped to maxDelay.
//
// This function does not terminate retrying. For example, the policy
// CapDelay(maxDelay, ExponentialBackoff(n)) will never stop retrying;
// it will reach a state where it retries forever with a delay of maxDelay
// between each attempt. To get termination, you need to combine this with
// LimitRetries or another limiting policy.
//
// Example:
//
//	// Exponential backoff starting at 100ms, capped at 5 seconds, up to 10 retries
//	policy := M.Concat(
//		LimitRetries(10),
//		CapDelay(5*time.Second, ExponentialBackoff(100*time.Millisecond)),
//	)(Monoid)
func CapDelay(maxDelay time.Duration, policy RetryPolicy) RetryPolicy {
	return F.Flow2(
		policy,
		O.Map(F.Bind1st(ord.Min(ordDuration), maxDelay)),
	)
}

// ExponentialBackoff creates a retry policy where the delay grows exponentially
// with each iteration. Each delay increases by a factor of two.
//
// The delay for iteration n is: delay * 2^n
//
// For example, with an initial delay of 100ms:
//   - Iteration 0: 100ms
//   - Iteration 1: 200ms
//   - Iteration 2: 400ms
//   - Iteration 3: 800ms
//   - etc.
//
// This policy never returns None, so it should be combined with LimitRetries
// and/or CapDelay to prevent unbounded delays.
//
// Example:
//
//	// Exponential backoff starting at 100ms, capped at 10s, up to 5 retries
//	policy := M.Concat(
//		LimitRetries(5),
//		CapDelay(10*time.Second, ExponentialBackoff(100*time.Millisecond)),
//	)(Monoid)
func ExponentialBackoff(delay time.Duration) RetryPolicy {
	return func(status RetryStatus) Option[time.Duration] {
		return O.Some(delay * time.Duration(math.Pow(2, float64(status.IterNumber))))
	}
}

// DefaultRetryStatus is the initial retry status used when starting a retry operation.
// It represents the state before any retries have been attempted:
//   - IterNumber: 0 (first attempt)
//   - CumulativeDelay: 0 (no delays yet)
//   - PreviousDelay: None (no previous attempt)
//
// This is exported primarily to allow user code to test their retry handlers
// and retry policies.
//
// Example:
//
//	policy := LimitRetries(3)
//	result := policy(DefaultRetryStatus) // Returns Some(0) for immediate retry
var DefaultRetryStatus = RetryStatus{
	IterNumber:      0,
	CumulativeDelay: 0,
	PreviousDelay:   O.None[time.Duration](),
}

var getOrElseDelay = O.GetOrElse(F.Constant(emptyDuration))

// ApplyPolicy applies a retry policy to the current status and returns the
// updated status for the next iteration. This function:
//   - Calls the policy with the current status to get the next delay
//   - Increments the iteration number
//   - Adds the delay to the cumulative delay
//   - Stores the delay as the previous delay for the next iteration
//
// This is useful for testing policies or implementing custom retry logic.
//
// Example:
//
//	policy := ExponentialBackoff(100 * time.Millisecond)
//	status := DefaultRetryStatus
//
//	// First retry
//	status = ApplyPolicy(policy, status)
//	// status.IterNumber == 1, status.PreviousDelay == Some(100ms)
//
//	// Second retry
//	status = ApplyPolicy(policy, status)
//	// status.IterNumber == 2, status.PreviousDelay == Some(200ms)
func ApplyPolicy(policy RetryPolicy, status RetryStatus) RetryStatus {
	previousDelay := policy(status)
	return RetryStatus{
		IterNumber:      status.IterNumber + 1,
		CumulativeDelay: status.CumulativeDelay + getOrElseDelay(previousDelay),
		PreviousDelay:   previousDelay,
	}
}

// Always creates a constant function that always returns the same value,
// ignoring the RetryStatus parameter. This is particularly useful as the
// check callback in Retrying functions to retry an operation unconditionally,
// independent of the operation's status or result.
//
// When used as a check callback in Retrying, Always(true) will cause the
// operation to retry on every iteration until the retry policy terminates
// (e.g., via LimitRetries or context cancellation). This is useful when you
// want to retry an operation a fixed number of times regardless of whether
// it succeeds or fails.
//
// Parameters:
//   - a: The constant value to return
//
// Returns:
//
//	A function that takes a RetryStatus and always returns the provided value a.
//
// Example with Retrying:
//
//	// Retry exactly 3 times with exponential backoff, regardless of success/failure
//	policy := M.Concat(
//		LimitRetries(3),
//		ExponentialBackoff(100*time.Millisecond),
//	)(Monoid)
//
//	action := func(status RetryStatus) ReaderResult[string] {
//		return func(ctx context.Context) (string, error) {
//			// This will be called 4 times total (initial + 3 retries)
//			return fetchData(ctx)
//		}
//	}
//
//	// Always retry, regardless of the result
//	retrying := Retrying(policy, action, Always(true))
//
// Example with custom logic:
//
//	// Create a function that always returns false
//	neverRetry := Always(false)
//	shouldRetry := neverRetry(status) // always returns false
//
//go:inline
func Always[A any](a A) func(RetryStatus) A {
	return F.Constant1[RetryStatus](a)
}
