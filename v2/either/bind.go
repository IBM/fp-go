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
	f func(S1) Either[E, T],
) func(Either[E, S1]) Either[E, S2] {
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
) func(Either[E, S1]) Either[E, S2] {
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
) func(Either[E, S1]) Either[E, S2] {
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
) func(Either[E, T]) Either[E, S1] {
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
) func(Either[E, S1]) Either[E, S2] {
	return A.ApS(
		Ap[S2, E, T],
		Map[E, S1, func(T) S2],
		setter,
		fa,
	)
}
