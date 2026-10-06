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

package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
)

// TailRec implements stack-safe tail recursion for the ReaderReaderIOEither monad.
//
// The step function f maps the current state to a computation that produces a
// Trampoline:
//   - Bounce(A): continue the recursion with the new state A
//   - Land(B): stop the recursion with the final result B
//
// A Left value produced by any step stops the recursion and is returned as is.
//
// The recursion runs in a loop, so the stack does not grow with the number of
// steps. The step function is only called when the resulting IO runs, and both
// environments R and C are passed unchanged to every step.
//
// Type Parameters:
//   - R: The outer environment type
//   - C: The inner environment type
//   - E: The error type
//   - A: The state type that changes during the recursion
//   - B: The final result type
//
// Parameters:
//   - f: A Kleisli arrow that computes the next step from the current state
//
// Returns:
//   - Kleisli[R, C, E, A, B]: A Kleisli arrow that runs the recursion from an initial state
//
// See Also:
//   - readerioeither.TailRec: Tail recursion with a single environment
//   - Chain: Sequences computations without loop semantics
func TailRec[R, C, E, A, B any](f Kleisli[R, C, E, A, Trampoline[A, B]]) Kleisli[R, C, E, A, B] {
	return func(a A) ReaderReaderIOEither[R, C, E, B] {
		return func(r R) ReaderIOEither[C, E, B] {
			return func(c C) IOEither[E, B] {
				return func() Either[E, B] {
					current := a
					for {
						next := f(current)(r)(c)()
						rec, e := either.Unwrap(next)
						if either.IsLeft(next) {
							return either.Left[B](e)
						}
						if rec.Landed {
							return either.Right[E](rec.Land)
						}
						current = rec.Bounce
					}
				}
			}
		}
	}
}
