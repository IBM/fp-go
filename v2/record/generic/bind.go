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
	Mo "github.com/IBM/fp-go/v2/monoid"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    Name  string
//	    Count int
//	}
//	result := generic.Do[map[string]State, string, State]()
func Do[GS ~map[K]S, K comparable, S any]() GS {
	return Empty[GS]()
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
// For records, this merges values by key where later steps can use values from earlier steps.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    Name  string
//	    Count int
//	}
//
//	result := F.Pipe2(
//	    generic.Do[map[string]State, string, State](),
//	    generic.Bind[map[string]State, map[string]State, map[string]string, string, State, State, string](
//	        monoid.Record[string, State](),
//	    )(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        func(s State) map[string]string {
//	            return map[string]string{"a": "Alice", "b": "Bob"}
//	        },
//	    ),
//	    generic.Bind[map[string]State, map[string]State, map[string]int, string, State, State, int](
//	        monoid.Record[string, State](),
//	    )(
//	        func(count int) func(State) State {
//	            return func(s State) State { s.Count = count; return s }
//	        },
//	        func(s State) map[string]int {
//	            // This can access s.Name from the previous step
//	            return map[string]int{"a": len(s.Name), "b": len(s.Name) * 2}
//	        },
//	    ),
//	)
func Bind[GS1 ~map[K]S1, GS2 ~map[K]S2, GT ~map[K]T, K comparable, S1, S2, T any](m Mo.Monoid[GS2]) func(setter func(T) func(S1) S2, f func(S1) GT) func(GS1) GS2 {
	c := Chain[GS1](m)
	return func(setter func(T) func(S1) S2, f func(S1) GT) func(GS1) GS2 {
		return C.Bind(
			c,
			Map[GT, GS2, K, T, S2],
			setter,
			f,
		)
	}
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~map[K]S1, GS2 ~map[K]S2, K comparable, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GS1) GS2 {
	return F.Let(
		Map[GS1, GS2, K, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~map[K]S1, GS2 ~map[K]S2, K comparable, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS1, GS2, K, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~map[K]S1, GT ~map[K]T, K comparable, S1, T any](setter func(T) S1) func(GT) GS1 {
	return C.BindTo(
		Map[GT, GS1, K, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel. For records, this merges values by key.
//
// Example:
//
//	type State struct {
//	    Name  string
//	    Score int
//	}
//
//	// These operations are independent and can be combined with ApS
//	names := map[string]string{"player1": "Alice", "player2": "Bob"}
//	scores := map[string]int{"player1": 100, "player2": 200}
//
//	result := F.Pipe2(
//	    generic.Do[map[string]State, string, State](),
//	    generic.ApS[map[string]State, map[string]State, map[string]string, string, State, State, string](
//	        monoid.Record[string, State](),
//	    )(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        names,
//	    ),
//	    generic.ApS[map[string]State, map[string]State, map[string]int, string, State, State, int](
//	        monoid.Record[string, State](),
//	    )(
//	        func(score int) func(State) State {
//	            return func(s State) State { s.Score = score; return s }
//	        },
//	        scores,
//	    ),
//	) // map[string]State{"player1": {Name: "Alice", Score: 100}, "player2": {Name: "Bob", Score: 200}}
func ApS[GS1 ~map[K]S1, GS2 ~map[K]S2, GT ~map[K]T, K comparable, S1, S2, T any](m Mo.Monoid[GS2]) func(setter func(T) func(S1) S2, fa GT) func(GS1) GS2 {
	a := Ap[GS2, map[K]func(T) S2, GT](m)
	return func(setter func(T) func(S1) S2, fa GT) func(GS1) GS2 {
		return A.ApS(
			a,
			Map[GS1, map[K]func(T) S2, K, S1, func(T) S2],
			setter,
			fa,
		)
	}
}
