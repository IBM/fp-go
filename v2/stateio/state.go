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
	"github.com/IBM/fp-go/v2/internal/statet"
	"github.com/IBM/fp-go/v2/io"
)

// Of creates a StateIO that wraps a pure value.
// The value is wrapped and the state is passed through unchanged.
//
// This is the Pointed/Applicative 'of' operation that lifts a pure value
// into the StateIO context.
//
// Example:
//
//	result := Of[AppState](42)
//	// Returns a computation containing 42 that passes state through unchanged
func Of[S, A any](a A) StateIO[S, A] {
	return statet.Of[StateIO[S, A]](io.Of[Pair[S, A]], a)
}

// MonadMap transforms the value of a StateIO using the provided function.
// The state is threaded through the computation unchanged.
// This is the functor map operation.
//
// Example:
//
//	result := MonadMap(
//	    Of[AppState](21),
//	    func(x int) int { return x * 2 },
//	) // Result contains 42
func MonadMap[S, A, B any](fa StateIO[S, A], f func(A) B) StateIO[S, B] {
	return statet.MonadMap[StateIO[S, A], StateIO[S, B]](
		io.MonadMap[Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Map is the curried version of [MonadMap].
// Returns a function that transforms a StateIO.
//
// Example:
//
//	double := Map[AppState](func(x int) int { return x * 2 })
//	result := function.Pipe1(Of[AppState](21), double)
func Map[S, A, B any](f func(A) B) Operator[S, A, B] {
	return statet.Map[StateIO[S, A], StateIO[S, B]](
		io.Map[Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadChain sequences two computations, passing the result of the first to a function
// that produces the second computation. This is the monadic bind operation.
// The state is threaded through both computations sequentially.
//
// Example:
//
//	result := MonadChain(
//	    Of[AppState](5),
//	    func(x int) StateIO[AppState, string] {
//	        return Of[AppState](fmt.Sprintf("value: %d", x))
//	    },
//	)
func MonadChain[S, A, B any](fa StateIO[S, A], f Kleisli[S, A, B]) StateIO[S, B] {
	return statet.MonadChain(
		io.MonadChain[Pair[S, A], Pair[S, B]],
		fa,
		f,
	)
}

// Chain is the curried version of [MonadChain].
// Returns a function that sequences computations.
//
// Example:
//
//	stringify := Chain(func(x int) StateIO[AppState, string] {
//	    return Of[AppState](fmt.Sprintf("%d", x))
//	})
//	result := function.Pipe1(Of[AppState](42), stringify)
func Chain[S, A, B any](f Kleisli[S, A, B]) Operator[S, A, B] {
	return statet.Chain[StateIO[S, A]](
		io.Chain[Pair[S, A], Pair[S, B]],
		f,
	)
}

// MonadAp applies a function wrapped in a StateIO to a value wrapped in a StateIO.
// The state is threaded through both computations sequentially.
// This is the applicative apply operation.
//
// Example:
//
//	fab := Of[AppState](func(x int) int { return x * 2 })
//	fa := Of[AppState](21)
//	result := MonadAp(fab, fa) // Result contains 42
func MonadAp[B, S, A any](fab StateIO[S, func(A) B], fa StateIO[S, A]) StateIO[S, B] {
	return statet.MonadAp[StateIO[S, A], StateIO[S, B]](
		io.MonadMap[Pair[S, A], Pair[S, B]],
		io.MonadChain[Pair[S, func(A) B], Pair[S, B]],
		fab,
		fa,
	)
}

// Ap is the curried version of [MonadAp].
// Returns a function that applies a wrapped function to the given wrapped value.
func Ap[B, S, A any](fa StateIO[S, A]) Operator[S, func(A) B, B] {
	return statet.Ap[StateIO[S, A], StateIO[S, B], StateIO[S, func(A) B]](
		io.Map[Pair[S, A], Pair[S, B]],
		io.Chain[Pair[S, func(A) B], Pair[S, B]],
		fa,
	)
}

// FromIO lifts an IO computation into StateIO.
// The IO computation is executed and its result is wrapped in StateIO.
// The state is passed through unchanged.
//
// Example:
//
//	ioAction := io.Of(42)
//	stateIOAction := FromIO[AppState](ioAction)
func FromIO[S, A any](fa IO[A]) StateIO[S, A] {
	return statet.FromF[StateIO[S, A]](
		io.MonadMap[A],
		fa,
	)
}

// Combinators

// FromIOK lifts an IO-returning function into a Kleisli arrow for StateIO.
// This is useful for composing functions that return IO actions with StateIO computations.
//
// Example:
//
//	readFile := func(path string) IO[string] { ... }
//	kleisli := FromIOK[AppState](readFile)
//	// kleisli can now be used with Chain
func FromIOK[S, A, B any](f func(A) IO[B]) Kleisli[S, A, B] {
	return function.Flow2(
		f,
		FromIO[S, B],
	)
}
