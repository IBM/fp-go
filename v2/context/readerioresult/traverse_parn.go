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

package readerioresult

import (
	"context"
	"sync"

	"github.com/IBM/fp-go/v2/result"
)

// TraverseArrayParN transforms an array by applying f to each element, running at most
// concurrency computations at the same time.
//
// TraverseArrayPar starts one goroutine per element, which is a problem for large
// inputs or for downstream services that limit the number of concurrent requests.
// TraverseArrayParN bounds the parallelism with a fixed number of workers. The order
// of the results matches the order of the input.
//
// The first failure cancels the context passed to the computations that are still
// running, no further elements are started, and the error is returned. Cancellation
// of the outer context stops the traversal in the same way.
//
// Parameters:
//   - concurrency: The maximum number of computations that run at the same time. A
//     value less than or equal to zero means no bound, like TraverseArrayPar.
//   - f: The function applied to each element
//
// Returns:
//   - Kleisli[[]A, []B]: A function from the input array to a ReaderIOResult of the results
//
// See Also:
//   - TraverseArrayPar: Unbounded parallel traversal
//   - TraverseArraySeq: Sequential traversal
func TraverseArrayParN[A, B any](concurrency int, f Kleisli[A, B]) Kleisli[[]A, []B] {
	return func(as []A) ReaderIOResult[[]B] {
		return func(ctx context.Context) IOResult[[]B] {
			return func() Result[[]B] {
				// quick check for cancellation
				if ctx.Err() != nil {
					return result.Left[[]B](context.Cause(ctx))
				}

				count := len(as)
				workers := concurrency
				if workers <= 0 || workers > count {
					workers = count
				}

				// the sub-context lets the first failure cancel all other computations
				ctxSub, cancelSub := context.WithCancelCause(ctx)
				defer cancelSub(nil)

				results := make([]B, count)
				indexes := make(chan int)

				var (
					wg       sync.WaitGroup
					errOnce  sync.Once
					firstErr error
				)

				wg.Add(workers)
				for range workers {
					go func() {
						defer wg.Done()
						for i := range indexes {
							b, err := result.Unwrap(WithContext(f(as[i]))(ctxSub)())
							if err != nil {
								errOnce.Do(func() {
									firstErr = err
									cancelSub(err)
								})
								continue
							}
							results[i] = b
						}
					}()
				}

				aborted := false
			feed:
				for i := range count {
					select {
					case indexes <- i:
					case <-ctxSub.Done():
						aborted = true
						break feed
					}
				}
				close(indexes)
				wg.Wait()

				if firstErr != nil {
					return result.Left[[]B](firstErr)
				}
				if aborted {
					return result.Left[[]B](context.Cause(ctxSub))
				}
				return result.Of(results)
			}
		}
	}
}
