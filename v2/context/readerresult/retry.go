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

package readerresult

import (
	"context"
	"time"

	RD "github.com/IBM/fp-go/v2/reader"
	R "github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

// Retrying retries a ReaderResult computation according to a retry policy with context awareness.
//
// This function implements a retry mechanism for operations that depend on a [context.Context]
// and can fail (Result). It respects context cancellation, meaning that if the context is
// cancelled during retry delays, the operation will stop immediately and return the cancellation error.
//
// The retry loop will continue until one of the following occurs:
//   - The action succeeds and the check function returns false (no retry needed)
//   - The retry policy returns None (retry limit reached)
//   - The check function returns false (indicating success or a non-retryable failure)
//   - The context is cancelled (returns context.Canceled or context.DeadlineExceeded)
//
// Type Parameters:
//   - A: The type of the success value
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
//     on a context.Context and produces a Result[A].
//
//   - check: A predicate function that examines the Result[A] and returns true if the
//     operation should be retried, or false if it should stop. This allows you to
//     distinguish between retryable failures (e.g., network timeouts) and permanent
//     failures (e.g., invalid input).
//
// Returns:
//   - A ReaderResult[A] that, when executed with a context, will perform the retry
//     logic with context cancellation support and return the final result.
//
// Example:
//
//	// Create a retry policy: exponential backoff with a cap, limited to 5 retries
//	policy := M.Concat(
//	    retry.LimitRetries(5),
//	    retry.CapDelay(10*time.Second, retry.ExponentialBackoff(100*time.Millisecond)),
//	)(retry.Monoid)
//
//	// Action that fetches data
//	fetchData := func(status retry.RetryStatus) ReaderResult[string] {
//	    return func(ctx context.Context) Result[string] {
//	        if ctx.Err() != nil {
//	            return result.Left[string](ctx.Err())
//	        }
//	        if status.IterNumber < 3 {
//	            return result.Left[string](fmt.Errorf("temporary error"))
//	        }
//	        return result.Of("success")
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
//	finalResult := retryingFetch(ctx)
//
//go:inline
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check Predicate[Result[A]],
) ReaderResult[A] {

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
	delayWithCancel := func(delay time.Duration) RD.Operator[context.Context, R.RetryStatus, R.RetryStatus] {
		return func(ri Reader[context.Context, R.RetryStatus]) Reader[context.Context, R.RetryStatus] {
			return func(ctx context.Context) R.RetryStatus {
				// Create a timeout context that will be cancelled when either:
				// - The delay duration expires, or
				// - The parent context is cancelled
				timeoutCtx, cancelTimeout := context.WithTimeout(ctx, delay)
				defer cancelTimeout()

				// Wait for either the timeout or parent context cancellation
				<-timeoutCtx.Done()

				// Dispatch to the next action with the original context.
				// WithContextK will handle context cancellation correctly.
				return ri(ctx)
			}
		}
	}

	// get an implementation for the types
	return RG.Retrying(
		RD.Chain[context.Context, Result[A], Trampoline[R.RetryStatus, Result[A]]],
		RD.Map[context.Context, R.RetryStatus, Trampoline[R.RetryStatus, Result[A]]],
		RD.Of[context.Context, Trampoline[R.RetryStatus, Result[A]]],
		RD.Of[context.Context, R.RetryStatus],
		delayWithCancel,

		RD.TailRec,

		policy,
		WithContextK(action),
		check,
	)

}
