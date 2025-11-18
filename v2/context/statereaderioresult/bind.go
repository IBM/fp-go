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

package statereaderioresult

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
//	    statereaderioresult.Do[AppState](State{}),
//	    statereaderioresult.Bind(...),
//	    statereaderioresult.Let(...),
//	)
//
//go:inline
func Do[ST, A any](
	empty A,
) StateReaderIOResult[ST, A] {
	return Of[ST](empty)
}

// Bind executes a computation and binds its result to a field in the accumulator state.
// This is used in do-notation to sequence dependent computations.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioresult.Do[AppState](State{}),
//	    statereaderioresult.Bind(
//	        func(name string) func(State) State {
//	            return func(s State) State { return State{name: name, age: s.age} }
//	        },
//	        func(s State) statereaderioresult.StateReaderIOResult[AppState, string] {
//	            return statereaderioresult.Of[AppState]("John")
//	        },
//	    ),
//	)
//
//go:inline
func Bind[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[ST, S1, T],
) Operator[ST, S1, S2] {
	return C.Bind(
		Chain[ST, S1, S2],
		Map[ST, T, S2],
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
//	    statereaderioresult.Do[AppState](State{age: 25}),
//	    statereaderioresult.Let(
//	        func(isAdult bool) func(State) State {
//	            return func(s State) State { return State{age: s.age, isAdult: isAdult} }
//	        },
//	        func(s State) bool { return s.age >= 18 },
//	    ),
//	)
//
//go:inline
func Let[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[ST, S1, S2] {
	return F.Let(
		Map[ST, S1, S2],
		key,
		f,
	)
}

// LetTo binds a constant value to a field in the accumulator state.
//
// Example:
//
//	result := function.Pipe2(
//	    statereaderioresult.Do[AppState](State{}),
//	    statereaderioresult.LetTo(
//	        func(status string) func(State) State {
//	            return func(s State) State { return State{...s, status: status} }
//	        },
//	        "active",
//	    ),
//	)
//
//go:inline
func LetTo[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[ST, S1, S2] {
	return F.LetTo(
		Map[ST, S1, S2],
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
//	    statereaderioresult.Of[AppState](42),
//	    statereaderioresult.BindTo[AppState](func(x int) State { return State{value: x} }),
//	)
//
//go:inline
func BindTo[ST, S1, T any](
	setter func(T) S1,
) Operator[ST, T, S1] {
	return C.BindTo(
		Map[ST, T, S1],
		setter,
	)
}

// ApS applies a computation in sequence and binds the result to a field.
// This is the applicative version of Bind.
//
//go:inline
func ApS[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa StateReaderIOResult[ST, T],
) Operator[ST, S1, S2] {
	return A.ApS(
		Ap[S2, ST, T],
		Map[ST, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a lens-based variant of ApS for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func ApSL[ST, S, T any](
	lens Lens[S, T],
	fa StateReaderIOResult[ST, T],
) Endomorphism[StateReaderIOResult[ST, S]] {
	return ApS(lens.Set, fa)
}

// BindL is a lens-based variant of Bind for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func BindL[ST, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, T, T],
) Endomorphism[StateReaderIOResult[ST, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL is a lens-based variant of Let for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func LetL[ST, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[StateReaderIOResult[ST, S]] {
	return Let[ST](lens.Set, function.Flow2(lens.Get, f))
}

// LetToL is a lens-based variant of LetTo for working with nested structures.
// It uses a lens to focus on a specific field in the state.
//
//go:inline
func LetToL[ST, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[StateReaderIOResult[ST, S]] {
	return LetTo[ST](lens.Set, b)
}
