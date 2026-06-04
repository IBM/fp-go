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

	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
)

// Unfold builds a lazy [SeqResult] from an initial seed by repeatedly applying
// a step function. It is the dual of a fold (an anamorphism).
//
// On each iteration f is called with the current seed B:
//   - [option.None] terminates the sequence normally (no more values).
//   - [option.Some](Pair(nextSeed, value)) emits value (type A) as a [result.Right]
//     and advances the seed to nextSeed for the next iteration.
//   - A [result.Left] error terminates the sequence and propagates that error.
//
// Context cancellation is checked at the start of every iteration. When the
// context is cancelled the sequence emits [result.Left](context.Cause(ctx))
// and stops.
//
// Example — emit integers 0..4:
//
//	counter := Unfold(func(n int) ReaderIOResult[Option[Pair[int, int]]] {
//	    if n >= 5 {
//	        return Of(option.None[Pair[int, int]]())
//	    }
//	    // Head = next seed, Tail = emitted value
//	    return Of(option.Some(pair.MakePair(n+1, n)))
//	})(0)
//	// counter(ctx) yields Right(0), Right(1), Right(2), Right(3), Right(4)
func Unfold[A, B any](
	f func(B) ReaderIOResult[Option[Pair[B, A]]],
) func(B) Reader[context.Context, SeqResult[A]] {
	return func(seed B) Reader[context.Context, SeqResult[A]] {
		return func(ctx context.Context) SeqResult[A] {
			return func(yield func(Result[A]) bool) {
				current := seed
				for {
					if ctx.Err() != nil {
						yield(result.Left[A](context.Cause(ctx)))
						return
					}
					opab, err := result.Unwrap(f(current)(ctx)())
					if err != nil {
						yield(result.Left[A](err))
						return
					}
					pab, ok := option.Unwrap(opab)
					if !ok {
						return
					}
					if !yield(result.Of(pair.Tail(pab))) {
						return
					}
					current = pair.Head(pab)
				}
			}
		}
	}
}
