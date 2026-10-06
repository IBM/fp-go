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

package readerreaderioresult

import (
	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
)

// TailRec implements stack-safe tail recursion for the ReaderReaderIOResult monad.
//
// The step function f maps the current state to a computation that produces a
// Trampoline:
//   - Bounce(A): continue the recursion with the new state A
//   - Land(B): stop the recursion with the final result B
//
// An error produced by any step stops the recursion and is returned as is.
//
// The recursion runs in a loop, so the stack does not grow with the number of
// steps. Before each step the context.Context is checked: if it has been
// cancelled, the recursion stops with the cancellation cause and the step
// function is not called again.
//
// Type Parameters:
//   - R: The outer environment type
//   - A: The state type that changes during the recursion
//   - B: The final result type
//
// Parameters:
//   - f: A Kleisli arrow that computes the next step from the current state
//
// Returns:
//   - Kleisli[R, A, B]: A Kleisli arrow that runs the recursion from an initial state
//
// See Also:
//   - readerreaderioeither.TailRec: Tail recursion without the cancellation check
//   - readerioresult.TailRec: Context-aware tail recursion without the outer environment
//
//go:inline
func TailRec[R, A, B any](f Kleisli[R, A, Trampoline[A, B]]) Kleisli[R, A, B] {
	return RRIOE.TailRec(function.Flow2(
		f,
		reader.Map[R](RIOE.WithContext[Trampoline[A, B]]),
	))
}
