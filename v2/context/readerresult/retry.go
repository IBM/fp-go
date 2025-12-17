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

//go:inline
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check func(Result[A]) bool,
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
		RD.Chain[context.Context, Result[A], Result[A]],
		RD.Chain[context.Context, R.RetryStatus, Result[A]],
		RD.Of[context.Context, Result[A]],
		RD.Of[context.Context, R.RetryStatus],
		delayWithCancel,

		policy,
		WithContextK(action),
		check,
	)

}
