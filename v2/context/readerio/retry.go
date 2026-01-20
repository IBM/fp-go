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
	"github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

// Retrying retries a ReaderIO computation according to a retry policy.
//
// This function implements a retry mechanism for operations that depend on a [context.Context]
// and perform side effects (IO). The retry loop continues until one of the following occurs:
//   - The action succeeds and the check function returns false (no retry needed)
//   - The retry policy returns None (retry limit reached)
//   - The check function returns false (indicating success or a non-retryable condition)
//
// Type Parameters:
//   - A: The type of the value produced by the action
//
// Parameters:
//
//   - policy: A RetryPolicy that determines when and how long to wait between retries.
//     The policy receives a RetryStatus on each iteration and returns an optional delay.
//     If it returns None, retrying stops. Common policies include LimitRetries,
//     ExponentialBackoff, and CapDelay from the retry package.
//
//   - action: A Kleisli arrow that takes a RetryStatus and returns a ReaderIO[A].
//     This function is called on each retry attempt and receives information about the
//     current retry state (iteration number, cumulative delay, etc.).
//
//   - check: A predicate function that examines the result A and returns true if the
//     operation should be retried, or false if it should stop. This allows you to
//     distinguish between retryable conditions and successful/permanent results.
//
// Returns:
//   - A ReaderIO[A] that, when executed with a context, will perform the retry logic
//     and return the final result.
//
// Example:
//
//	// Create a retry policy: exponential backoff with a cap, limited to 5 retries
//	policy := M.Concat(
//	    retry.LimitRetries(5),
//	    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
//	)(retry.Monoid)
//
//	// Action that fetches data, with retry status information
//	fetchData := func(status retry.RetryStatus) ReaderIO[string] {
//	    return func(ctx context.Context) IO[string] {
//	        return func() string {
//	            // Simulate an operation that might fail
//	            if status.IterNumber < 3 {
//	                return ""  // Empty result indicates failure
//	            }
//	            return "success"
//	        }
//	    }
//	}
//
//	// Check function: retry if result is empty
//	shouldRetry := func(s string) bool {
//	    return s == ""
//	}
//
//	// Create the retrying computation
//	retryingFetch := Retrying(policy, fetchData, shouldRetry)
//
//	// Execute
//	ctx := t.Context()
//	result := retryingFetch(ctx)() // Returns "success" after 3 attempts
//
//go:inline
func Retrying[A any](
	policy retry.RetryPolicy,
	action Kleisli[retry.RetryStatus, A],
	check Predicate[A],
) ReaderIO[A] {
	// get an implementation for the types
	return RG.Retrying(
		Chain[A, Trampoline[retry.RetryStatus, A]],
		Map[retry.RetryStatus, Trampoline[retry.RetryStatus, A]],
		Of[Trampoline[retry.RetryStatus, A]],
		Of[retry.RetryStatus],
		Delay[retry.RetryStatus],

		TailRec,

		policy,
		action,
		check,
	)
}
