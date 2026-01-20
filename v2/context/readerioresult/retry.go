// Copyright (c) 2023 - 2025 IBM Corp.
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

package readerioresult

import (
	"context"
	"time"

	RIO "github.com/IBM/fp-go/v2/context/readerio"
	R "github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

// Retrying retries a ReaderIOResult computation according to a retry policy with context awareness.
//
// This function implements a retry mechanism for operations that depend on a [context.Context],
// perform side effects (IO), and can fail (Result). It respects context cancellation, meaning
// that if the context is cancelled during retry delays, the operation will stop immediately
// and return the cancellation error.
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
//   - action: A Kleisli arrow that takes a RetryStatus and returns a ReaderIOResult[A].
//     This function is called on each retry attempt and receives information about the
//     current retry state (iteration number, cumulative delay, etc.). The action depends
//     on a context.Context and produces a Result[A]. The context passed to the action
//     will be the same context used for retry delays, so cancellation is properly propagated.
//
//   - check: A predicate function that examines the Result[A] and returns true if the
//     operation should be retried, or false if it should stop. This allows you to
//     distinguish between retryable failures (e.g., network timeouts) and permanent
//     failures (e.g., invalid input). Note that context cancellation errors will
//     automatically stop retrying regardless of this function's return value.
//
// Returns:
//
//	A ReaderIOResult[A] that, when executed with a context, will perform the retry
//	logic with context cancellation support and return the final result.
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
// Example:
//
//	// Create a retry policy: exponential backoff with a cap, limited to 5 retries
//	policy := M.Concat(
//	    retry.LimitRetries(5),
//	    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
//	)(retry.Monoid)
//
//	// Action that fetches data, with retry status information
//	fetchData := func(status retry.RetryStatus) ReaderIOResult[string] {
//	    return func(ctx context.Context) IOResult[string] {
//	        return func() Result[string] {
//	            // Check if context is cancelled
//	            if ctx.Err() != nil {
//	                return result.Left[string](ctx.Err())
//	            }
//	            // Simulate an HTTP request that might fail
//	            if status.IterNumber < 3 {
//	                return result.Left[string](fmt.Errorf("temporary error"))
//	            }
//	            return result.Of("success")
//	        }
//	    }
//	}
//
//	// Check function: retry on any error except context cancellation
//	shouldRetry := func(r Result[string]) bool {
//	    return result.IsLeft(r) && !errors.Is(result.GetLeft(r), context.Canceled)
//	}
//
//	// Create the retrying computation
//	retryingFetch := Retrying(policy, fetchData, shouldRetry)
//
//	// Execute with a cancellable context
//	ctx, cancel := context.WithTimeout(t.Context(), 5*time.Second)
//	defer cancel()
//	ioResult := retryingFetch(ctx)
//	finalResult := ioResult()
//
// See also:
//   - retry.RetryPolicy for available retry policies
//   - retry.RetryStatus for information passed to the action
//   - context.Context for context cancellation semantics
//
//go:inline
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check Predicate[Result[A]],
) ReaderIOResult[A] {

	// delayWithCancel implements a context-aware delay mechanism for retry operations.
	// It creates a timeout context that will be cancelled when either:
	//   1. The delay duration expires (normal case), or
	//   2. The parent context is cancelled (early termination)
	//
	// The function waits on timeoutCtx.Done(), which will be signaled in either case:
	//   - If the delay expires, timeoutCtx is cancelled by the timeout
	//   - If the parent ctx is cancelled, timeoutCtx inherits the cancellation
	//
	// After the wait completes, we dispatch to the next action by calling ri(ctx)().
	// This works correctly because the action is wrapped in WithContextK, which handles
	// context cancellation by checking ctx.Err() and returning an appropriate error
	// (context.Canceled or context.DeadlineExceeded) when the context is cancelled.
	//
	// This design ensures that:
	//   - Retry delays respect context cancellation and terminate immediately
	//   - The cancellation error propagates correctly through the retry chain
	//   - No unnecessary delays occur when the context is already cancelled
	delayWithCancel := func(delay time.Duration) RIO.Operator[R.RetryStatus, R.RetryStatus] {
		return func(ri ReaderIO[R.RetryStatus]) ReaderIO[R.RetryStatus] {
			return func(ctx context.Context) IO[R.RetryStatus] {
				return func() R.RetryStatus {
					// Create a timeout context that will be cancelled when either:
					// - The delay duration expires, or
					// - The parent context is cancelled
					timeoutCtx, cancelTimeout := context.WithTimeout(ctx, delay)
					defer cancelTimeout()

					// Wait for either the timeout or parent context cancellation
					<-timeoutCtx.Done()

					// Dispatch to the next action with the original context.
					// WithContextK will handle context cancellation correctly.
					return ri(ctx)()
				}
			}
		}
	}

	// get an implementation for the types
	return RG.Retrying(
		RIO.Chain[Result[A], Trampoline[R.RetryStatus, Result[A]]],
		RIO.Map[R.RetryStatus, Trampoline[R.RetryStatus, Result[A]]],
		RIO.Of[Trampoline[R.RetryStatus, Result[A]]],
		RIO.Of[R.RetryStatus],
		delayWithCancel,

		RIO.TailRec,

		policy,
		WithContextK(action),
		check,
	)

}
