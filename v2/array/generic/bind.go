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

package generic

import (
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
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
//	result := generic.Do[[]State, State](State{})
func Do[GS ~[]S, S any](
	empty S,
) GS {
	return Of[GS](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
// For arrays, this produces the cartesian product where later steps can use values from earlier steps.
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
//	    generic.Do[[]State, State](State{}),
//	    generic.Bind[[]State, []State, []int, State, State, int](
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        func(s State) []int {
//	            return []int{1, 2, 3}
//	        },
//	    ),
//	    generic.Bind[[]State, []State, []int, State, State, int](
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        func(s State) []int {
//	            // This can access s.X from the previous step
//	            return []int{s.X * 10, s.X * 20}
//	        },
//	    ),
//	) // Produces: {1,10}, {1,20}, {2,20}, {2,40}, {3,30}, {3,60}
func Bind[GS1 ~[]S1, GS2 ~[]S2, GT ~[]T, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GT,
) func(GS1) GS2 {
	return C.Bind(
		Chain[GS1, GS2, S1, S2],
		Map[GT, GS2, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~[]S1, GS2 ~[]S2, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GS1) GS2 {
	return F.Let(
		Map[GS1, GS2, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~[]S1, GS2 ~[]S2, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS1, GS2, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~[]S1, GT ~[]T, S1, T any](
	setter func(T) S1,
) func(GT) GS1 {
	return C.BindTo(
		Map[GT, GS1, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel. For arrays, this produces the cartesian product.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y string
//	}
//
//	// These operations are independent and can be combined with ApS
//	xValues := []int{1, 2}
//	yValues := []string{"a", "b"}
//
//	result := F.Pipe2(
//	    generic.Do[[]State, State](State{}),
//	    generic.ApS[[]State, []State, []int, State, State, int](
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        xValues,
//	    ),
//	    generic.ApS[[]State, []State, []string, State, State, string](
//	        func(y string) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        yValues,
//	    ),
//	) // [{1,"a"}, {1,"b"}, {2,"a"}, {2,"b"}]
func ApS[GS1 ~[]S1, GS2 ~[]S2, GT ~[]T, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GT,
) func(GS1) GS2 {
	return A.ApS(
		Ap[GS2, []func(T) S2, GT, S2, T],
		Map[GS1, []func(T) S2, S1, func(T) S2],
		setter,
		fa,
	)
}
