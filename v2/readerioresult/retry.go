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
	"github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/retry"
)

// Retrying retries a ReaderIOResult computation according to a retry policy.
//
// This function implements a retry mechanism for operations that depend on a context (Reader),
// perform side effects (IO), and can fail (Result). It will repeatedly execute the action
// according to the retry policy until either:
//   - The action succeeds and the check function returns false (no retry needed)
//   - The retry policy returns None (retry limit reached)
//   - The check function returns false (indicating success or a non-retryable failure)
//
// Parameters:
//
//   - policy: A RetryPolicy that determines when and how long to wait between retries.
//     The policy receives a RetryStatus on each iteration and returns an optional delay.
//     If it returns None, retrying stops. Common policies include LimitRetries,
//     ExponentialBackoff, and CapDelay from the retry package.
//
//   - action: A Kleisli arrow that takes a RetryStatus and returns a ReaderIOResult[R, A].
//     This function is called on each retry attempt and receives information about the
//     current retry state (iteration number, cumulative delay, etc.). The action depends
//     on a context of type R and produces a Result[A].
//
//   - check: A predicate function that examines the Result[A] and returns true if the
//     operation should be retried, or false if it should stop. This allows you to
//     distinguish between retryable failures (e.g., network timeouts) and permanent
//     failures (e.g., invalid input).
//
// Returns:
//
//	A ReaderIOResult[R, A] that, when executed with a context, will perform the retry
//	logic and return the final result.
//
// Type Parameters:
//   - R: The type of the context/environment required by the action
//   - A: The type of the success value
//
// Example:
//
//	type Config struct {
//	    MaxRetries int
//	    BaseURL    string
//	}
//
//	// Create a retry policy: exponential backoff with a cap, limited to 5 retries
//	policy := M.Concat(
//	    retry.LimitRetries(5),
//	    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
//	)(retry.Monoid)
//
//	// Action that fetches data, with retry status information
//	fetchData := func(status retry.RetryStatus) ReaderIOResult[Config, string] {
//	    return func(cfg Config) IOResult[string] {
//	        return func() Result[string] {
//	            // Simulate an HTTP request that might fail
//	            if status.IterNumber < 3 {
//	                return result.Left[string](fmt.Errorf("temporary error"))
//	            }
//	            return result.Right[error]("success")
//	        }
//	    }
//	}
//
//	// Check function: retry on any error
//	shouldRetry := func(r Result[string]) bool {
//	    return result.IsLeft(r)
//	}
//
//	// Create the retrying computation
//	retryingFetch := Retrying(policy, fetchData, shouldRetry)
//
//	// Execute with a config
//	cfg := Config{MaxRetries: 5, BaseURL: "https://api.example.com"}
//	ioResult := retryingFetch(cfg)
//	finalResult := ioResult()
//
// See also:
//   - retry.RetryPolicy for available retry policies
//   - retry.RetryStatus for information passed to the action
//   - readerioeither.Retrying for the underlying implementation
//
//go:inline
func Retrying[R, A any](
	policy retry.RetryPolicy,
	action Kleisli[R, retry.RetryStatus, A],
	check func(Result[A]) bool,
) ReaderIOResult[R, A] {
	// get an implementation for the types
	return readerioeither.Retrying(policy, action, check)
}
