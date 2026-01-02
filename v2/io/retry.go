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

// Package io provides the IO monad for managing side effects in a functional way.
// This file contains retry combinators for the IO monad.
package io

import (
	R "github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

type (
	// RetryStatus is an IO computation that returns retry status information.
	// It wraps the retry.RetryStatus type in the IO monad, allowing retry
	// status to be computed lazily as part of an IO effect.
	RetryStatus = IO[R.RetryStatus]
)

// Retrying retries an IO action according to a retry policy until it succeeds or the policy gives up.
//
// This function implements retry logic for IO computations that don't raise exceptions but
// signal failure through their result value. The retry behavior is controlled by three parameters:
//
// Parameters:
//   - policy: A RetryPolicy that determines the delay between retries and when to stop.
//     Policies can be combined using the retry.Monoid to create complex retry strategies.
//   - action: A Kleisli arrow (function) that takes the current RetryStatus and returns an IO[A].
//     The action is executed on each attempt, receiving updated status information including
//     the iteration number, cumulative delay, and previous delay.
//   - check: A predicate function that examines the result of the action and returns true if
//     the operation should be retried, or false if it succeeded. This allows you to define
//     custom success criteria based on the result value.
//
// The function will:
//  1. Execute the action with the current retry status
//  2. Apply the check predicate to the result
//  3. If check returns false (success), return the result
//  4. If check returns true (should retry), apply the policy to get the next delay
//  5. If the policy returns None, stop retrying and return the last result
//  6. If the policy returns Some(delay), wait for that duration and retry from step 1
//
// The action receives RetryStatus information on each attempt, which includes:
//   - IterNumber: The current attempt number (0-indexed, so 0 is the first attempt)
//   - CumulativeDelay: The total time spent waiting between retries so far
//   - PreviousDelay: The delay from the last retry (None on the first attempt)
//
// This information can be used for logging, implementing custom backoff strategies,
// or making decisions within the action itself.
//
// Example - Retry HTTP request with exponential backoff:
//
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(5),
//	    retry.ExponentialBackoff(100 * time.Millisecond),
//	)
//
//	result := io.Retrying(
//	    policy,
//	    func(status retry.RetryStatus) io.IO[*http.Response] {
//	        log.Printf("Attempt %d (cumulative delay: %v)", status.IterNumber, status.CumulativeDelay)
//	        return io.Of(http.Get("https://api.example.com/data"))
//	    },
//	    func(resp *http.Response) bool {
//	        // Retry on server errors (5xx status codes)
//	        return resp.StatusCode >= 500
//	    },
//	)
//
// Example - Retry until a condition is met:
//
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(10),
//	    retry.ConstantDelay(500 * time.Millisecond),
//	)
//
//	result := io.Retrying(
//	    policy,
//	    func(status retry.RetryStatus) io.IO[string] {
//	        return fetchStatus()
//	    },
//	    func(status string) bool {
//	        // Retry until status is "ready"
//	        return status != "ready"
//	    },
//	)
//
//go:inline
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check Predicate[A],
) IO[A] {
	// Delegate to the generic retry implementation, providing the IO monad operations
	return RG.Retrying(
		Chain[A, Trampoline[R.RetryStatus, A]],
		Map[R.RetryStatus, Trampoline[R.RetryStatus, A]],
		Of[Trampoline[R.RetryStatus, A]],
		Of[R.RetryStatus],    // Pure/return for the status type
		Delay[R.RetryStatus], // Delay operation for the status type

		TailRec,

		policy,
		action,
		check,
	)
}
