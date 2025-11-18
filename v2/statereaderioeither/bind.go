// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do starts a do-notation chain for building computations in a fluent style.
// This is typically used with Bind, Let, and other combinators to compose
// stateful, context-dependent computations that can fail.
//
// Example:
//
//	type State struct {
//	    name string
//	    age  int
//	}
//	result := function.Pipe2(
//	    statereaderioeither.Do[AppState, Config, error](State{}),
//	    statereaderioeither.Bind(...),
//	    statereaderioeither.Let(...),
//	)
//
//go:inline
func Do[ST, R, E, A any](
	empty A,
) StateReaderIOEither[ST, R, E, A] {
	return Of[ST, R, E](empty)
}

// Bind executes a computation and binds its result to a field in the accumulator state.
// This is used in do-notation to sequence dependent computations.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioeither.Do[AppState, Config, error](State{}),
//	    statereaderioeither.Bind(
//	        func(name string) func(State) State {
//	            return func(s State) State { return State{name: name, age: s.age} }
//	        },
//	        func(s State) statereaderioeither.StateReaderIOEither[AppState, Config, error, string] {
//	            return statereaderioeither.Of[AppState, Config, error]("John")
//	        },
//	    ),
//	)
//
//go:inline
func Bind[ST, R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[ST, R, E, S1, T],
) Operator[ST, R, E, S1, S2] {
	return C.Bind(
		Chain[ST, R, E, S1, S2],
		Map[ST, R, E, T, S2],
		setter,
		f,
	)
}

// Let computes a derived value and binds it to a field in the accumulator state.
// Unlike Bind, this does not execute a monadic computation, just a pure function.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioeither.Do[AppState, Config, error](State{age: 25}),
//	    statereaderioeither.Let(
//	        func(isAdult bool) func(State) State {
//	            return func(s State) State { return State{age: s.age, isAdult: isAdult} }
//	        },
//	        func(s State) bool { return s.age >= 18 },
//	    ),
//	)
//
//go:inline
func Let[ST, R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[ST, R, E, S1, S2] {
	return F.Let(
		Map[ST, R, E, S1, S2],
		key,
		f,
	)
}

// LetTo binds a constant value to a field in the accumulator state.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioeither.Do[AppState, Config, error](State{}),
//	    statereaderioeither.LetTo(
//	        func(status string) func(State) State {
//	            return func(s State) State { return State{...s, status: status} }
//	        },
//	        "active",
//	    ),
//	)
//
//go:inline
func LetTo[ST, R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[ST, R, E, S1, S2] {
	return F.LetTo(
		Map[ST, R, E, S1, S2],
		key,
		b,
	)
}

// BindTo wraps a value in a simple constructor, typically used to start a do-notation chain
// after getting an initial value.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioeither.Of[AppState, Config, error](42),
//	    statereaderioeither.BindTo[AppState, Config, error](func(x int) State { return State{value: x} }),
//	)
//
//go:inline
func BindTo[ST, R, E, S1, T any](
	setter func(T) S1,
) Operator[ST, R, E, T, S1] {
	return C.BindTo(
		Map[ST, R, E, T, S1],
		setter,
	)
}

// ApS applies a computation in sequence and binds the result to a field.
// This is the applicative version of Bind.
//
//go:inline
func ApS[ST, R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa StateReaderIOEither[ST, R, E, T],
) Operator[ST, R, E, S1, S2] {
	return A.ApS(
		Ap[S2, ST, R, E, T],
		Map[ST, R, E, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a lens-based variant of ApS for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func ApSL[ST, R, E, S, T any](
	lens Lens[S, T],
	fa StateReaderIOEither[ST, R, E, T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return ApS(lens.Set, fa)
}

// BindL is a lens-based variant of Bind for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func BindL[ST, R, E, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, R, E, T, T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL is a lens-based variant of Let for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func LetL[ST, R, E, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return Let[ST, R, E](lens.Set, function.Flow2(lens.Get, f))
}

// LetToL is a lens-based variant of LetTo for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func LetToL[ST, R, E, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[StateReaderIOEither[ST, R, E, S]] {
	return LetTo[ST, R, E](lens.Set, b)
}
