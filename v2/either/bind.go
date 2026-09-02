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

package either

import (
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/internal/common"
	FN "github.com/IBM/fp-go/v2/function"
)

// Do creates an empty context of type S to be used with the Bind operation.
// This is the starting point for do-notation style computations.
//
// Example:
//
//	type State struct { x, y int }
//	result := either.Do[error](State{})
//
//go:inline
func Do[E, S any](
	empty S,
) Either[E, S] {
	return Of[E](empty)
}

// Bind attaches the result of a computation to a context S1 to produce a context S2.
// This enables building up complex computations in a pipeline.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Do[error](State{}),
//	    either.Bind(
//	        func(v int) func(State) State {
//	            return func(s State) State { return State{value: v} }
//	        },
//	        func(s State) either.Either[error, int] {
//	            return either.Right[error](42)
//	        },
//	    ),
//	)
//
//go:inline
func Bind[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[E, S1, T],
) Operator[E, S1, S2] {
	return C.Bind(
		Chain[E, S1, S2],
		Map[E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Similar to Bind but for pure (non-Either) computations.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Right[error](State{value: 10}),
//	    either.Let(
//	        func(v int) func(State) State {
//	            return func(s State) State { return State{value: s.value + v} }
//	        },
//	        func(s State) int { return 32 },
//	    ),
//	) // Right(State{value: 42})
//
//go:inline
func Let[E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[E, S1, S2] {
	return F.Let(
		Map[E, S1, S2],
		key,
		f,
	)
}

// LetTo attaches a constant value to a context S1 to produce a context S2.
//
// Example:
//
//	type State struct { name string }
//	result := F.Pipe2(
//	    either.Right[error](State{}),
//	    either.LetTo(
//	        func(n string) func(State) State {
//	            return func(s State) State { return State{name: n} }
//	        },
//	        "Alice",
//	    ),
//	) // Right(State{name: "Alice"})
//
//go:inline
func LetTo[E, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[E, S1, S2] {
	return F.LetTo(
		Map[E, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state S1 from a value T.
// This is typically used to start a bind chain.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Right[error](42),
//	    either.BindTo(func(v int) State { return State{value: v} }),
//	) // Right(State{value: 42})
//
//go:inline
func BindTo[E, S1, T any](
	setter func(T) S1,
) Operator[E, T, S1] {
	return C.BindTo(
		Map[E, T, S1],
		setter,
	)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering the context and the value concurrently.
// Uses applicative semantics rather than monadic sequencing.
//
// Example:
//
//	type State struct { x, y int }
//	result := F.Pipe2(
//	    either.Right[error](State{x: 10}),
//	    either.ApS(
//	        func(y int) func(State) State {
//	            return func(s State) State { return State{x: s.x, y: y} }
//	        },
//	        either.Right[error](32),
//	    ),
//	) // Right(State{x: 10, y: 32})
//
//go:inline
func ApS[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[E, T],
) Operator[E, S1, S2] {
	return A.ApS(
		Ap[S2, E, T],
		Map[E, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a Lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// Unlike ApS, which requires a manually-written curried setter function, ApSL derives the
// setter directly from the provided Lens.
//
// Type Parameters:
//   - E: The error type (Left side of Either)
//   - S: The context structure type (unchanged — both input and output are S)
//   - T: The type of the value being attached
//
// Parameters:
//   - lens: A Lens[S, T] providing both the getter and the setter for the focused field
//   - fa: An Either[E, T] supplying the value to attach
//
// Returns:
//   - An Operator[E, S, S] that updates the S context via the lens
//
// See Also:
//   - ApS: The setter-function variant
//   - BindL: The monadic lens-based variant
//
//go:inline
func ApSL[E, S, T any](
	lens L.Lens[S, T],
	fa Either[E, T],
) Operator[E, S, S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a Lens.
//
// The lens provides both the getter (to read the current field value and pass it to f)
// and the setter (to write the new value back into the structure S).
//
// Type Parameters:
//   - E: The error type (Left side of Either)
//   - S: The context structure type (unchanged — both input and output are S)
//   - T: The type of the value being focused and transformed
//
// Parameters:
//   - lens: A Lens[S, T] providing both getter and setter for the focused field
//   - f: A Kleisli[E, T, T] — a function from the current T to Either[E, T]
//
// Returns:
//   - An Operator[E, S, S] that reads T via lens.Get, runs f, and writes the result via lens.Set
//
// See Also:
//   - Bind: The setter-function variant
//   - ApSL: The applicative lens-based variant
//
//go:inline
func BindL[E, S, T any](
	lens L.Lens[S, T],
	f Kleisli[E, T, T],
) Operator[E, S, S] {
	return Bind(lens.Set, FN.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure transformation to a context using a lens-based setter.
// This is a convenience function that combines Let with a Lens.
//
// The lens provides both the getter (to read the current field value and pass it to f)
// and the setter (to write the transformed value back into the structure S).
//
// Type Parameters:
//   - E: The error type (Left side of Either)
//   - S: The context structure type (unchanged — both input and output are S)
//   - T: The type of the value being focused and transformed
//
// Parameters:
//   - lens: A Lens[S, T] providing both getter and setter for the focused field
//   - f: A pure function T → T applied to the current field value
//
// Returns:
//   - An Operator[E, S, S] that reads T via lens.Get, applies f, and writes the result via lens.Set
//
// See Also:
//   - Let: The setter-function variant
//   - BindL: The monadic lens-based variant
//
//go:inline
func LetL[E, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[E, S, S] {
	return Let[E](lens.Set, FN.Flow2(lens.Get, f))
}

// LetToL attaches a constant value to a context using a lens-based setter.
// This is a convenience function that combines LetTo with a Lens.
//
// Unlike LetL, which transforms the current field value, LetToL replaces it with
// the provided constant b, regardless of the current value.
//
// Type Parameters:
//   - E: The error type (Left side of Either)
//   - S: The context structure type (unchanged — both input and output are S)
//   - T: The type of the value being set
//
// Parameters:
//   - lens: A Lens[S, T] providing the setter for the focused field
//   - b: The constant value of type T to set in every S
//
// Returns:
//   - An Operator[E, S, S] that sets the focused field to b via the lens
//
// See Also:
//   - LetTo: The setter-function variant
//   - LetL: The lens-based transformation variant
//
//go:inline
func LetToL[E, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[E, S, S] {
	return LetTo[E](lens.Set, b)
}
