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

package stateio

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do starts a do-notation chain for building computations in a fluent style.
// This is typically used with Bind, Let, and other combinators to compose
// stateful computations with side effects.
//
// Example:
//
//	type Result struct {
//	    name string
//	    age  int
//	}
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    Bind(...),
//	    Let(...),
//	)
//
//go:inline
func Do[ST, A any](
	empty A,
) StateIO[ST, A] {
	return Of[ST](empty)
}

// Bind executes a computation and binds its result to a field in the accumulator state.
// This is used in do-notation to sequence dependent computations.
//
// The setter function takes the computed value and returns a function that updates
// the accumulator state. The computation function (f) receives the current accumulator
// state and returns a StateIO computation.
//
// Example:
//
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    Bind(
//	        func(name string) func(Result) Result {
//	            return func(r Result) Result { return Result{name: name, age: r.age} }
//	        },
//	        func(r Result) StateIO[AppState, string] {
//	            return Of[AppState]("John")
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
// The key function takes the computed value and returns a function that updates
// the accumulator state. The computation function (f) receives the current accumulator
// state and returns a pure value.
//
// Example:
//
//	result := function.Pipe2(
//	    Do[AppState](Result{age: 25}),
//	    Let(
//	        func(isAdult bool) func(Result) Result {
//	            return func(r Result) Result { return Result{age: r.age, isAdult: isAdult} }
//	        },
//	        func(r Result) bool { return r.age >= 18 },
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
// This is useful for setting fixed values in the accumulator during do-notation.
//
// Example:
//
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    LetTo(
//	        func(status string) func(Result) Result {
//	            return func(r Result) Result { return Result{status: status} }
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
// after getting an initial value. This transforms a StateIO[S, T] into StateIO[S, S1]
// by applying a constructor function.
//
// Example:
//
//	result := function.Pipe2(
//	    Of[AppState](42),
//	    BindTo[AppState](func(x int) Result { return Result{value: x} }),
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
// This is the applicative version of Bind, useful for parallel-style composition
// where computations don't depend on each other's results.
//
// Example:
//
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    ApS(
//	        func(count int) func(Result) Result {
//	            return func(r Result) Result { return Result{count: count} }
//	        },
//	        Of[AppState](42),
//	    ),
//	)
//
//go:inline
func ApS[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa StateIO[ST, T],
) Operator[ST, S1, S2] {
	return A.ApS(
		Ap[S2, ST, T],
		Map[ST, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a lens-based variant of ApS for working with nested structures.
// It uses a lens to focus on a specific field in the accumulator state,
// making it easier to update nested fields without manual destructuring.
//
// Example:
//
//	nameLens := lens.Prop[Result, string]("name")
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    ApSL(nameLens, Of[AppState]("John")),
//	)
//
//go:inline
func ApSL[ST, S, T any](
	lens Lens[S, T],
	fa StateIO[ST, T],
) Endomorphism[StateIO[ST, S]] {
	return ApS(lens.Set, fa)
}

// BindL is a lens-based variant of Bind for working with nested structures.
// It uses a lens to focus on a specific field in the accumulator state,
// allowing you to update that field based on a computation that depends on its current value.
//
// Example:
//
//	counterLens := lens.Prop[Result, int]("counter")
//	result := function.Pipe2(
//	    Do[AppState](Result{counter: 0}),
//	    BindL(counterLens, func(n int) StateIO[AppState, int] {
//	        return Of[AppState](n + 1)
//	    }),
//	)
//
//go:inline
func BindL[ST, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, T, T],
) Endomorphism[StateIO[ST, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL is a lens-based variant of Let for working with nested structures.
// It uses a lens to focus on a specific field in the accumulator state,
// allowing you to update that field using a pure function.
//
// Example:
//
//	counterLens := lens.Prop[Result, int]("counter")
//	result := function.Pipe2(
//	    Do[AppState](Result{counter: 5}),
//	    LetL(counterLens, N.Mul(2)),
//	)
//
//go:inline
func LetL[ST, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[StateIO[ST, S]] {
	return Let[ST](lens.Set, function.Flow2(lens.Get, f))
}

// LetToL is a lens-based variant of LetTo for working with nested structures.
// It uses a lens to focus on a specific field in the accumulator state,
// allowing you to set that field to a constant value.
//
// Example:
//
//	statusLens := lens.Prop[Result, string]("status")
//	result := function.Pipe2(
//	    Do[AppState](Result{}),
//	    LetToL(statusLens, "active"),
//	)
//
//go:inline
func LetToL[ST, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[StateIO[ST, S]] {
	return LetTo[ST](lens.Set, b)
}
