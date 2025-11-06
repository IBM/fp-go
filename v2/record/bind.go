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

package record

import (
	Mo "github.com/IBM/fp-go/v2/monoid"
	G "github.com/IBM/fp-go/v2/record/generic"
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
//	result := record.Do[string, State]()
func Do[K comparable, S any]() map[K]S {
	return G.Do[map[K]S, K, S]()
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
// For records, this merges values by key.
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
//	    record.Do[string, State](),
//	    record.Bind(monoid.Record[string, State]())(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        func(s State) map[string]string {
//	            return map[string]string{"a": "Alice", "b": "Bob"}
//	        },
//	    ),
//	    record.Bind(monoid.Record[string, State]())(
//	        func(count int) func(State) State {
//	            return func(s State) State { s.Count = count; return s }
//	        },
//	        func(s State) map[string]int {
//	            // This can access s.Name from the previous step
//	            return map[string]int{"a": len(s.Name), "b": len(s.Name) * 2}
//	        },
//	    ),
//	)
func Bind[S1, T any, K comparable, S2 any](m Mo.Monoid[map[K]S2]) func(setter func(T) func(S1) S2, f func(S1) map[K]T) func(map[K]S1) map[K]S2 {
	return G.Bind[map[K]S1, map[K]S2, map[K]T, K, S1, S2, T](m)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, T any, K comparable, S2 any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(map[K]S1) map[K]S2 {
	return G.Let[map[K]S1, map[K]S2, K, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, T any, K comparable, S2 any](
	setter func(T) func(S1) S2,
	b T,
) func(map[K]S1) map[K]S2 {
	return G.LetTo[map[K]S1, map[K]S2, K, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any, K comparable](setter func(T) S1) func(map[K]T) map[K]S1 {
	return G.BindTo[map[K]S1, map[K]T, K, S1, T](setter)
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
//	    Name  string
//	    Count int
//	}
//
//	// These operations are independent and can be combined with ApS
//	names := map[string]string{"a": "Alice", "b": "Bob"}
//	counts := map[string]int{"a": 10, "b": 20}
//
//	result := F.Pipe2(
//	    record.Do[string, State](),
//	    record.ApS(monoid.Record[string, State]())(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        names,
//	    ),
//	    record.ApS(monoid.Record[string, State]())(
//	        func(count int) func(State) State {
//	            return func(s State) State { s.Count = count; return s }
//	        },
//	        counts,
//	    ),
//	) // map[string]State{"a": {Name: "Alice", Count: 10}, "b": {Name: "Bob", Count: 20}}
func ApS[S1, T any, K comparable, S2 any](m Mo.Monoid[map[K]S2]) func(setter func(T) func(S1) S2, fa map[K]T) func(map[K]S1) map[K]S2 {
	return G.ApS[map[K]S1, map[K]S2, map[K]T, K, S1, S2, T](m)
}
