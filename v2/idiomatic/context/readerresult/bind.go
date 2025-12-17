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

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	"github.com/IBM/fp-go/v2/idiomatic/result"
	AP "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/reader"
	RES "github.com/IBM/fp-go/v2/result"
)

// Do initializes a do-notation context with an empty state.
//
// This is the starting point for do-notation style composition, which allows
// imperative-style sequencing of ReaderResult computations while maintaining
// functional purity.
//
// Type Parameters:
//   - S: The state type
//
// Parameters:
//   - empty: The initial empty state
//
// Returns:
//   - A ReaderResult[S] containing the initial state
//
//go:inline
func Do[S any](
	empty S,
) ReaderResult[S] {
	return RR.Do[context.Context](empty)
}

// Bind sequences an EFFECTFUL ReaderResult computation and updates the state with its result.
//
// IMPORTANT: Bind is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The Kleisli parameter (State -> ReaderResult[T]) is effectful because ReaderResult
// depends on context.Context (can be cancelled, has deadlines, carries values).
//
// For PURE FUNCTIONS (side-effect free), use:
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//   - Let: For pure functions without errors (State -> Value)
//
// This is the core operation for do-notation, allowing you to chain computations
// where each step can depend on the accumulated state and update it with new values.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the computation
//
// Parameters:
//   - setter: A function that takes the computation result and returns a state updater
//   - f: A Kleisli arrow that produces the next effectful computation based on current state
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return C.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		WithContextK(f),
	)
}

// Let attaches the result of a PURE computation to a state.
//
// IMPORTANT: Let is for PURE FUNCTIONS (side-effect free) that don't depend on context.Context.
// The function parameter (State -> Value) is pure - it only reads from state with no effects.
//
// For EFFECTFUL FUNCTIONS (that need context.Context), use:
//   - Bind: For effectful ReaderResult computations (State -> ReaderResult[Value])
//
// For PURE FUNCTIONS with error handling, use:
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//
// Unlike Bind, Let works with pure functions (not ReaderResult computations).
// This is useful for deriving values from the current state without performing
// any effects.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value computed
//
// Parameters:
//   - setter: A function that takes the computed value and returns a state updater
//   - f: A pure function that computes a value from the current state
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return RR.Let[context.Context](setter, f)
}

// LetTo attaches a constant value to a state.
// This is a PURE operation (side-effect free).
//
// This is a simplified version of Let for when you want to add a constant
// value to the state without computing it.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of the constant value
//
// Parameters:
//   - setter: A function that takes the constant and returns a state updater
//   - b: The constant value to attach
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return RR.LetTo[context.Context](setter, b)
}

// BindTo initializes do-notation by binding a value to a state.
//
// This is typically used as the first operation after a computation to
// start building up a state structure.
//
// Type Parameters:
//   - S1: The state type to create
//   - T: The type of the initial value
//
// Parameters:
//   - setter: A function that creates the initial state from a value
//
// Returns:
//   - An Operator that transforms ReaderResult[T] to ReaderResult[S1]
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return RR.BindTo[context.Context](setter)
}

// BindToP initializes do-notation by binding a value to a state using a Prism.
//
// This is a variant of BindTo that uses a prism instead of a setter function.
// Prisms are useful for working with sum types and optional values.
//
// Type Parameters:
//   - S1: The state type to create
//   - T: The type of the initial value
//
// Parameters:
//   - setter: A prism that can construct the state from a value
//
// Returns:
//   - An Operator that transforms ReaderResult[T] to ReaderResult[S1]
//
//go:inline
func BindToP[S1, T any](
	setter Prism[S1, T],
) Operator[T, S1] {
	return BindTo(setter.ReverseGet)
}

// ApS attaches a value to a context using applicative style.
//
// IMPORTANT: ApS is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The ReaderResult parameter is effectful because it depends on context.Context.
//
// Unlike Bind (which sequences operations), ApS can be used when operations are
// independent and can conceptually run in parallel.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the computation
//
// Parameters:
//   - setter: A function that takes the computation result and returns a state updater
//   - fa: An effectful ReaderResult computation
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderResult[T],
) Operator[S1, S2] {
	return AP.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a variant of ApS that uses a lens to focus on a specific field in the state.
//
// IMPORTANT: ApSL is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The ReaderResult parameter is effectful because it depends on context.Context.
//
// Instead of providing a setter function, you provide a lens that knows how to get and set
// the field. This is more convenient when working with nested structures.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field to update
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - fa: An effectful ReaderResult computation that produces a value of type T
//
// Returns:
//   - An Operator that transforms ReaderResult[S] to ReaderResult[S]
//
//go:inline
func ApSL[S, T any](
	lens Lens[S, T],
	fa ReaderResult[T],
) Operator[S, S] {
	return ApS(lens.Set, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific field in the state.
//
// IMPORTANT: BindL is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The Kleisli parameter returns a ReaderResult, which is effectful.
//
// It combines lens-based field access with monadic composition, allowing you to:
// 1. Extract a field value using the lens
// 2. Use that value in an effectful computation that may fail
// 3. Update the field with the result
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field to update
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - f: An effectful Kleisli arrow that transforms the field value
//
// Returns:
//   - An Operator that transforms ReaderResult[S] to ReaderResult[S]
//
//go:inline
func BindL[S, T any](
	lens Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return RR.BindL(lens, WithContextK(f))
}

// LetL is a variant of Let that uses a lens to focus on a specific field in the state.
//
// IMPORTANT: LetL is for PURE FUNCTIONS (side-effect free) that don't depend on context.Context.
// The endomorphism parameter is a pure function (T -> T) with no errors or effects.
//
// It applies a pure transformation to the focused field without any effects.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field to update
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - f: A pure endomorphism that transforms the field value
//
// Returns:
//   - An Operator that transforms ReaderResult[S] to ReaderResult[S]
//
//go:inline
func LetL[S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Operator[S, S] {
	return RR.LetL[context.Context](lens, f)
}

// LetToL is a variant of LetTo that uses a lens to focus on a specific field in the state.
//
// IMPORTANT: LetToL is for setting constant values. This is a PURE operation (side-effect free).
//
// It sets the focused field to a constant value.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field to update
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - b: The constant value to set
//
// Returns:
//   - An Operator that transforms ReaderResult[S] to ReaderResult[S]
//
//go:inline
func LetToL[S, T any](
	lens Lens[S, T],
	b T,
) Operator[S, S] {
	return RR.LetToL[context.Context](lens, b)
}

// BindReaderK binds a Reader computation (context-dependent but error-free) into the do-notation chain.
//
// IMPORTANT: This is for functions that depend on context.Context but don't return errors.
// The Reader[Context, T] is effectful because it depends on context.Context.
// Use this when you need context values but the operation cannot fail.
//
//go:inline
func BindReaderK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[context.Context, S1, T],
) Operator[S1, S2] {
	return RR.BindReaderK(setter, f)
}

// BindEitherK binds a Result (Either) computation into the do-notation chain.
//
// IMPORTANT: This is for PURE FUNCTIONS (side-effect free) that return Result[T].
// The function (State -> Result[T]) is pure - it only depends on state, not context.
// Use this for pure error-handling logic that doesn't need context.
//
//go:inline
func BindEitherK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f RES.Kleisli[S1, T],
) Operator[S1, S2] {
	return RR.BindEitherK[context.Context](setter, f)
}

// BindResultK binds an idiomatic Go function (returning value and error) into the do-notation chain.
//
// IMPORTANT: This is for PURE FUNCTIONS (side-effect free) that return (Value, error).
// The function (State -> (Value, error)) is pure - it only depends on state, not context.
// Use this for pure computations with error handling that don't need context.
//
// For EFFECTFUL FUNCTIONS (that need context.Context), use Bind instead.
//
//go:inline
func BindResultK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[S1, S2] {
	return RR.BindResultK[context.Context](setter, f)
}

// BindToReader converts a Reader computation into a ReaderResult and binds it to create an initial state.
//
// IMPORTANT: Reader[Context, T] is EFFECTFUL because it depends on context.Context.
// Use this when you have a context-dependent computation that cannot fail.
//
//go:inline
func BindToReader[
	S1, T any](
	setter func(T) S1,
) func(Reader[context.Context, T]) ReaderResult[S1] {
	return RR.BindToReader[context.Context](setter)
}

// BindToEither converts a Result (Either) into a ReaderResult and binds it to create an initial state.
//
// IMPORTANT: Result[T] is PURE (side-effect free) - it doesn't depend on context.
// Use this to lift pure error-handling values into the ReaderResult context.
//
//go:inline
func BindToEither[
	S1, T any](
	setter func(T) S1,
) func(Result[T]) ReaderResult[S1] {
	return RR.BindToEither[context.Context](setter)
}

// BindToResult converts an idiomatic Go tuple (value, error) into a ReaderResult and binds it to create an initial state.
//
// IMPORTANT: The (Value, error) tuple is PURE (side-effect free) - it doesn't depend on context.
// Use this to lift pure Go error-handling results into the ReaderResult context.
//
//go:inline
func BindToResult[
	S1, T any](
	setter func(T) S1,
) func(T, error) ReaderResult[S1] {
	return RR.BindToResult[context.Context](setter)
}

// ApReaderS applies a Reader computation in applicative style, combining it with the current state.
//
// IMPORTANT: Reader[Context, T] is EFFECTFUL because it depends on context.Context.
// Use this for context-dependent operations that cannot fail.
//
//go:inline
func ApReaderS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[context.Context, T],
) Operator[S1, S2] {
	return RR.ApReaderS(setter, fa)
}

// ApResultS applies an idiomatic Go tuple (value, error) in applicative style.
//
// IMPORTANT: The (Value, error) tuple is PURE (side-effect free) - it doesn't depend on context.
// Use this for pure Go error-handling results.
//
//go:inline
func ApResultS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
) func(T, error) Operator[S1, S2] {
	return RR.ApResultS[context.Context](setter)
}

// ApEitherS applies a Result (Either) in applicative style, combining it with the current state.
//
// IMPORTANT: Result[T] is PURE (side-effect free) - it doesn't depend on context.
// Use this for pure error-handling values.
//
//go:inline
func ApEitherS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[S1, S2] {
	return RR.ApEitherS[context.Context](setter, fa)
}
