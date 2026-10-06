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

import "github.com/IBM/fp-go/v2/context/readerreaderioresult"

// TailRec implements stack-safe tail recursion for Effect.
//
// The step function f maps the current state to an effect that produces a
// Trampoline:
//   - Bounce(A): continue the recursion with the new state A
//   - Land(B): stop the recursion with the final result B
//
// An error produced by any step stops the recursion and is returned as is.
//
// The recursion runs in a loop, so the stack does not grow with the number of
// steps; use it instead of a recursive Chain for loops with many iterations.
// Before each step the context.Context is checked: if it has been cancelled,
// the recursion stops with the cancellation cause.
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The state type that changes during the recursion
//   - B: The final result type
//
// Parameters:
//   - f: A Kleisli arrow that computes the next step from the current state
//
// Returns:
//   - Kleisli[C, A, B]: A Kleisli arrow that runs the recursion from an initial state
//
// See Also:
//   - Chain: Sequences effects without loop semantics
//   - Retrying: Repeats an effect according to a retry policy
//
//go:inline
func TailRec[C, A, B any](f Kleisli[C, A, Trampoline[A, B]]) Kleisli[C, A, B] {
	return readerreaderioresult.TailRec(f)
}
