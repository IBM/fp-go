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
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/function"
)

// Flatten removes one level of nesting from an Effect that produces another Effect.
//
// The outer effect runs first; on success the inner effect runs with the same
// context. A failure of either effect fails the result.
//
// Type Parameters:
//   - C: The context type required by both effects
//   - A: The type of the success value of the inner effect
//
// Parameters:
//   - mma: An Effect whose success value is itself an Effect
//
// Returns:
//   - Effect[C, A]: The flattened effect
//
// See Also:
//   - Chain: Flatten is Chain with the identity function
//   - FlattenThunk: Flattens an Effect that produces a Thunk
//
//go:inline
func Flatten[C, A any](mma Effect[C, Effect[C, A]]) Effect[C, A] {
	return readerreaderioresult.Flatten(mma)
}

// FlattenIO runs the IO action produced by an Effect and succeeds with its value.
//
// This is useful in the middle of a pipeline, when a Map step produces an IO
// action that still has to be executed.
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The type of the value produced by the IO action
//
// Parameters:
//   - mma: An Effect whose success value is an IO[A]
//
// Returns:
//   - Effect[C, A]: An effect that runs the produced IO action
//
// See Also:
//   - ChainIOK: FlattenIO is ChainIOK with the identity function
//   - FlattenIOResult: Similar but the IO action may fail
//
//go:inline
func FlattenIO[C, A any](mma Effect[C, IO[A]]) Effect[C, A] {
	return readerreaderioresult.FlattenIO(mma)
}

// FlattenIOResult runs the IOResult produced by an Effect.
//
// A successful IOResult yields its value, a failed one fails the effect.
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The type of the success value of the IOResult
//
// Parameters:
//   - mma: An Effect whose success value is an IOResult[A]
//
// Returns:
//   - Effect[C, A]: An effect that runs the produced IOResult
//
// See Also:
//   - ChainIOResultK: FlattenIOResult is ChainIOResultK with the identity function
//   - FlattenIO: Similar but the IO action cannot fail
//   - FlattenThunk: Similar but the action also sees context.Context
//
//go:inline
func FlattenIOResult[C, A any](mma Effect[C, IOResult[A]]) Effect[C, A] {
	return readerreaderioresult.FlattenIOResult(mma)
}

// FlattenThunk runs the Thunk produced by an Effect with the runtime context.Context.
//
// A successful Thunk yields its value, a failed one fails the effect.
//
// Type Parameters:
//   - C: The context type required by the effect
//   - A: The type of the success value of the Thunk
//
// Parameters:
//   - mma: An Effect whose success value is a Thunk[A]
//
// Returns:
//   - Effect[C, A]: An effect that runs the produced Thunk
//
// See Also:
//   - ChainThunkK: FlattenThunk is ChainThunkK with the identity function
//   - FlattenIOResult: Similar but the action does not see context.Context
//
//go:inline
func FlattenThunk[C, A any](mma Effect[C, Thunk[A]]) Effect[C, A] {
	return ChainThunkK[C](function.Identity[Thunk[A]])(mma)
}
