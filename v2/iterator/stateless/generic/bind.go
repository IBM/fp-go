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
	"github.com/IBM/fp-go/v2/internal/apply"
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
//	result := generic.Do[Iterator[State]](State{})
func Do[GS ~func() Option[Pair[GS, S]], S any](
	empty S,
) GS {
	return Of[GS](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
// For iterators, this produces the cartesian product where later steps can use values from earlier steps.
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
//	    generic.Do[Iterator[State]](State{}),
//	    generic.Bind[Iterator[State], Iterator[State], Iterator[int], State, State, int](
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        func(s State) Iterator[int] {
//	            return generic.Of[Iterator[int]](1, 2, 3)
//	        },
//	    ),
//	    generic.Bind[Iterator[State], Iterator[State], Iterator[int], State, State, int](
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        func(s State) Iterator[int] {
//	            // This can access s.X from the previous step
//	            return generic.Of[Iterator[int]](s.X * 10, s.X * 20)
//	        },
//	    ),
//	) // Produces: {1,10}, {1,20}, {2,20}, {2,40}, {3,30}, {3,60}
func Bind[GS1 ~func() Option[Pair[GS1, S1]], GS2 ~func() Option[Pair[GS2, S2]], GA ~func() Option[Pair[GA, A]], S1, S2, A any](
	setter func(A) func(S1) S2,
	f func(S1) GA,
) func(GS1) GS2 {

	return C.Bind(
		Chain[GS2, GS1, S1, S2],
		Map[GS2, GA, func(A) S2, A, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~func() Option[Pair[GS1, S1]], GS2 ~func() Option[Pair[GS2, S2]], S1, S2, A any](
	key func(A) func(S1) S2,
	f func(S1) A,
) func(GS1) GS2 {
	return F.Let(
		Map[GS2, GS1, func(S1) S2, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~func() Option[Pair[GS1, S1]], GS2 ~func() Option[Pair[GS2, S2]], S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS2, GS1, func(S1) S2, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~func() Option[Pair[GS1, S1]], GA ~func() Option[Pair[GA, A]], S1, A any](
	setter func(A) S1,
) func(GA) GS1 {
	return C.BindTo(
		Map[GS1, GA, func(A) S1, A, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel. For iterators, this produces the cartesian product.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y string
//	}
//
//	// These operations are independent and can be combined with ApS
//	xIter := generic.Of[Iterator[int]](1, 2, 3)
//	yIter := generic.Of[Iterator[string]]("a", "b")
//
//	result := F.Pipe2(
//	    generic.Do[Iterator[State]](State{}),
//	    generic.ApS[Iterator[func(int) State], Iterator[State], Iterator[State], Iterator[int], State, State, int](
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        xIter,
//	    ),
//	    generic.ApS[Iterator[func(string) State], Iterator[State], Iterator[State], Iterator[string], State, State, string](
//	        func(y string) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        yIter,
//	    ),
//	) // Produces: {1,"a"}, {1,"b"}, {2,"a"}, {2,"b"}, {3,"a"}, {3,"b"}
func ApS[GAS2 ~func() Option[Pair[GAS2, func(A) S2]], GS1 ~func() Option[Pair[GS1, S1]], GS2 ~func() Option[Pair[GS2, S2]], GA ~func() Option[Pair[GA, A]], S1, S2, A any](
	setter func(A) func(S1) S2,
	fa GA,
) func(GS1) GS2 {
	return apply.ApS(
		Ap[GAS2, GS2, GA, A, S2],
		Map[GAS2, GS1, func(S1) func(A) S2, S1, func(A) S2],
		setter,
		fa,
	)
}
