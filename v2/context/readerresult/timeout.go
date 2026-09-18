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

package readerresult

import (
	"context"
	"time"

	"github.com/IBM/fp-go/v2/pair"
)

func pairFromContextCancel(newCtx context.Context, cancelFct context.CancelFunc) pair.Pair[context.CancelFunc, context.Context] {
	return pair.MakePair(cancelFct, newCtx)
}

// WithTimeout adds a timeout to the context for a ReaderResult computation.
//
// This is a convenience wrapper around Local that uses [context.WithTimeout].
// The timeout is relative to when the ReaderResult is executed, not when
// WithTimeout is called. The cancel function is automatically called when
// the computation completes, ensuring proper cleanup. A computation that
// respects its context observes [context.DeadlineExceeded] once the timeout expires.
//
// Type Parameters:
//   - A: The value type of the ReaderResult
//
// Parameters:
//   - timeout: The maximum duration for the computation
//
// Returns:
//   - An Operator that runs the computation with a timeout
//
// Example:
//
//	fetchData := func(ctx context.Context) Result[Data] {
//	    select {
//	    case <-time.After(10 * time.Second):
//	        return result.Of(Data{Value: "slow"})
//	    case <-ctx.Done():
//	        return result.Left[Data](ctx.Err())
//	    }
//	}
//
//	res := F.Pipe1(
//	    fetchData,
//	    readerresult.WithTimeout[Data](5*time.Second),
//	)
//	res(context.Background()) // Left(context.DeadlineExceeded) after 5s
//
// See Also:
//   - WithDeadline: Uses an absolute point in time instead of a duration
//   - Local: The general mechanism for transforming the context
func WithTimeout[A any](timeout time.Duration) Operator[A, A] {
	return Local[A](func(ctx context.Context) pair.Pair[context.CancelFunc, context.Context] {
		return pairFromContextCancel(context.WithTimeout(ctx, timeout))
	})
}

// WithDeadline adds an absolute deadline to the context for a ReaderResult computation.
//
// This is a convenience wrapper around Local that uses [context.WithDeadline].
// Unlike [WithTimeout], the deadline is an absolute time. If the parent context
// already has an earlier deadline, that one takes precedence. The cancel function
// is automatically called when the computation completes.
//
// Type Parameters:
//   - A: The value type of the ReaderResult
//
// Parameters:
//   - deadline: The absolute time by which the computation must complete
//
// Returns:
//   - An Operator that runs the computation with a deadline
//
// Example:
//
//	deadline := time.Now().Add(5 * time.Second)
//	res := F.Pipe1(
//	    fetchData,
//	    readerresult.WithDeadline[Data](deadline),
//	)
//	res(context.Background()) // Left(context.DeadlineExceeded) if the deadline passes
//
// See Also:
//   - WithTimeout: Uses a relative duration instead of an absolute time
//   - Local: The general mechanism for transforming the context
func WithDeadline[A any](deadline time.Time) Operator[A, A] {
	return Local[A](func(ctx context.Context) pair.Pair[context.CancelFunc, context.Context] {
		return pairFromContextCancel(context.WithDeadline(ctx, deadline))
	})
}
