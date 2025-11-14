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

package stateless

import (
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//	result := stateless.Do(State{})
func Do[S any](
	empty S,
) Iterator[S] {
	return G.Do[Iterator[S]](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
// For iterators, this produces the cartesian product of all values.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//
//	result := F.Pipe2(
//	    stateless.Do(State{}),
//	    stateless.Bind(
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        func(s State) stateless.Iterator[int] {
//	            return stateless.Of(1, 2, 3)
//	        },
//	    ),
//	    stateless.Bind(
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        func(s State) stateless.Iterator[int] {
//	            // This can access s.X from the previous step
//	            return stateless.Of(s.X * 10, s.X * 20)
//	        },
//	    ),
//	) // Produces: {1,10}, {1,20}, {2,20}, {2,40}, {3,30}, {3,60}
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return G.Bind[Iterator[S1], Iterator[S2]](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return G.Let[Iterator[S1], Iterator[S2]](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return G.LetTo[Iterator[S1], Iterator[S2]](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return G.BindTo[Iterator[S1], Iterator[T]](setter)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//
//	// These operations are independent and can be combined with ApS
//	xValues := stateless.Of(1, 2, 3)
//	yValues := stateless.Of(10, 20)
//
//	result := F.Pipe2(
//	    stateless.Do(State{}),
//	    stateless.ApS(
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        xValues,
//	    ),
//	    stateless.ApS(
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        yValues,
//	    ),
//	) // Produces all combinations: {1,10}, {1,20}, {2,10}, {2,20}, {3,10}, {3,20}
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Iterator[T],
) Operator[S1, S2] {
	return G.ApS[Iterator[func(T) S2], Iterator[S1], Iterator[S2]](setter, fa)
}
