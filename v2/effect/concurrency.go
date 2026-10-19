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

package effect

import (
	thunk "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
)

// liftTraverse lifts a traversal of Thunk computations to a traversal of effects that
// share the same dependencies.
func liftTraverse[C, A, B any](
	traverse func(thunk.Kleisli[A, B]) thunk.Kleisli[[]A, []B],
) func(Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return func(f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
		return func(as []A) Effect[C, []B] {
			return func(c C) Thunk[[]B] {
				return traverse(function.Flow2(f, reader.Read[Thunk[B]](c)))(as)
			}
		}
	}
}

// TraverseArraySeq applies an effectful function to each element of an array, one after
// the other, and collects the results.
//
// Use it when the effects must not overlap, for example because they depend on a shared
// resource. The first failure stops the traversal.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the input elements
//   - B: The type of the output elements
//
// Parameters:
//   - f: The function applied to each element
//
// Returns:
//   - Kleisli[C, []A, []B]: A function from the input array to an effect of the results
//
// See Also:
//   - TraverseArrayPar: Runs the effects in parallel
//   - TraverseArrayParN: Runs the effects in parallel with bounded concurrency
func TraverseArraySeq[C, A, B any](f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return liftTraverse[C](thunk.TraverseArraySeq[A, B])(f)
}

// TraverseArrayPar applies an effectful function to each element of an array in parallel
// and collects the results in the order of the input.
//
// Every element runs in its own goroutine. The first failure cancels the effects that are
// still running and is returned.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the input elements
//   - B: The type of the output elements
//
// Parameters:
//   - f: The function applied to each element
//
// Returns:
//   - Kleisli[C, []A, []B]: A function from the input array to an effect of the results
//
// See Also:
//   - TraverseArrayParN: Bounds the number of effects that run at the same time
//   - TraverseArraySeq: Runs the effects one after the other
func TraverseArrayPar[C, A, B any](f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return liftTraverse[C](thunk.TraverseArrayPar[A, B])(f)
}

// TraverseArrayParN applies an effectful function to each element of an array with at most
// concurrency effects running at the same time, and collects the results in the order of
// the input.
//
// This corresponds to Effect.forEach with the concurrency option in effect-ts. The first
// failure cancels the effects that are still running, no further elements are started, and
// the error is returned.
//
// Type Parameters:
//   - C: The dependencies required by the effects
//   - A: The type of the input elements
//   - B: The type of the output elements
//
// Parameters:
//   - concurrency: The maximum number of effects that run at the same time. A value less
//     than or equal to zero means no bound, like TraverseArrayPar.
//   - f: The function applied to each element
//
// Returns:
//   - Kleisli[C, []A, []B]: A function from the input array to an effect of the results
//
// See Also:
//   - TraverseArrayPar: Unbounded parallel traversal
//   - WithSemaphore: Bounds concurrency across independent effects
func TraverseArrayParN[C, A, B any](concurrency int, f Kleisli[C, A, B]) Kleisli[C, []A, []B] {
	return liftTraverse[C](function.Bind1st(thunk.TraverseArrayParN[A, B], concurrency))(f)
}

// ApSeq applies the function produced by an effect to the value produced by fa, running
// the function effect first and fa afterwards.
//
// Type Parameters:
//   - B: The type of the result
//   - C: The dependencies required by the effects
//   - A: The type of the value produced by fa
//
// Parameters:
//   - fa: The effect that produces the argument
//
// Returns:
//   - Operator[C, func(A) B, B]: A function from the effect of the function to the result
//
// See Also:
//   - ApPar: Runs both effects in parallel
//   - Ap: Uses the default mode, which is parallel
func ApSeq[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B] {
	return function.Bind2nd(readerreaderioresult.MonadApSeq[C, A, B], fa)
}

// ApPar applies the function produced by an effect to the value produced by fa, running
// both effects in parallel.
//
// If one of the effects fails, the other one is cancelled.
//
// Type Parameters:
//   - B: The type of the result
//   - C: The dependencies required by the effects
//   - A: The type of the value produced by fa
//
// Parameters:
//   - fa: The effect that produces the argument
//
// Returns:
//   - Operator[C, func(A) B, B]: A function from the effect of the function to the result
//
// See Also:
//   - ApSeq: Runs the effects one after the other
func ApPar[B, C, A any](fa Effect[C, A]) Operator[C, func(A) B, B] {
	return function.Bind2nd(readerreaderioresult.MonadApPar[C, A, B], fa)
}
