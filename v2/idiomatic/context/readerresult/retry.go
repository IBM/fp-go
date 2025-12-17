// Copyright (c) 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package readerresult

import (
	RS "github.com/IBM/fp-go/v2/context/readerresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	R "github.com/IBM/fp-go/v2/retry"
)

// Retrying retries a ReaderResult computation according to a retry policy with context awareness.
//
// This is the idiomatic wrapper around the functional [github.com/IBM/fp-go/v2/context/readerresult.Retrying]
// function. It provides a more Go-friendly API by working with (value, error) tuples instead of Result types.
//
// The function implements a retry mechanism for operations that depend on a [context.Context] and can fail.
// It respects context cancellation, meaning that if the context is cancelled during retry delays, the
// operation will stop immediately and return the cancellation error.
//
// The retry loop will continue until one of the following occurs:
//   - The action succeeds and the check function returns false (no retry needed)
//   - The retry policy returns None (retry limit reached)
//   - The check function returns false (indicating success or a non-retryable failure)
//   - The context is cancelled (returns context.Canceled or context.DeadlineExceeded)
//
// Parameters:
//
//   - policy: A RetryPolicy that determines when and how long to wait between retries.
//     The policy receives a RetryStatus on each iteration and returns an optional delay.
//     If it returns None, retrying stops. Common policies include LimitRetries,
//     ExponentialBackoff, and CapDelay from the retry package.
//
//   - action: A Kleisli arrow that takes a RetryStatus and returns a ReaderResult[A].
//     This function is called on each retry attempt and receives information about the
//     current retry state (iteration number, cumulative delay, etc.). The action depends
//     on a context.Context and produces (A, error). The context passed to the action
//     will be the same context used for retry delays, so cancellation is properly propagated.
//
//   - check: A predicate function that examines the result value and error, returning true
//     if the operation should be retried, or false if it should stop. This allows you to
//     distinguish between retryable failures (e.g., network timeouts) and permanent
//     failures (e.g., invalid input). The function receives both the value and error from
//     the action's result. Note that context cancellation errors will automatically stop
//     retrying regardless of this function's return value.
//
// Returns:
//
//	A ReaderResult[A] that, when executed with a context, will perform the retry
//	logic with context cancellation support and return the final (value, error) tuple.
//
// Type Parameters:
//   - A: The type of the success value
//
// Context Cancellation:
//
// The retry mechanism respects context cancellation in two ways:
//  1. During retry delays: If the context is cancelled while waiting between retries,
//     the operation stops immediately and returns the context error.
//  2. During action execution: If the action itself checks the context and returns
//     an error due to cancellation, the retry loop will stop (assuming the check
//     function doesn't force a retry on context errors).
//
// Implementation Details:
//
// This function wraps the functional [github.com/IBM/fp-go/v2/context/readerresult.Retrying]
// by converting between the idiomatic (value, error) tuple representation and the functional
// Result[A] representation. The conversion is handled by ToReaderResult and FromReaderResult,
// ensuring seamless integration with the underlying retry mechanism that uses delayWithCancel
// to properly handle context cancellation during delays.
//
// Example:
//
//	// Create a retry policy: exponential backoff with a cap, limited to 5 retries
//	policy := retry.Monoid.Concat(
//	    retry.LimitRetries(5),
//	    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
//	)
//
//	// Action that fetches data, with retry status information
//	fetchData := func(status retry.RetryStatus) ReaderResult[string] {
//	    return func(ctx context.Context) (string, error) {
//	        // Check if context is cancelled
//	        if ctx.Err() != nil {
//	            return "", ctx.Err()
//	        }
//	        // Simulate an HTTP request that might fail
//	        if status.IterNumber < 3 {
//	            return "", fmt.Errorf("temporary error")
//	        }
//	        return "success", nil
//	    }
//	}
//
//	// Check function: retry on any error except context cancellation
//	shouldRetry := func(val string, err error) bool {
//	    return err != nil && !errors.Is(err, context.Canceled)
//	}
//
//	// Create the retrying computation
//	retryingFetch := Retrying(policy, fetchData, shouldRetry)
//
//	// Execute with a cancellable context
//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//	defer cancel()
//	result, err := retryingFetch(ctx)
//
// See also:
//   - retry.RetryPolicy for available retry policies
//   - retry.RetryStatus for information passed to the action
//   - context.Context for context cancellation semantics
//   - github.com/IBM/fp-go/v2/context/readerresult.Retrying for the underlying functional implementation
//
//go:inline
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check func(A, error) bool,
) ReaderResult[A] {
	return F.Pipe1(
		RS.Retrying(
			policy,
			F.Flow2(
				action,
				ToReaderResult,
			),
			func(a Result[A]) bool {
				return check(result.Unwrap(a))
			},
		),
		FromReaderResult,
	)
}
