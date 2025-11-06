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

package array

import (
	G "github.com/IBM/fp-go/v2/array/generic"
)

// Do creates an empty context of type S to be used with the Bind operation.
// This is the starting point for monadic do-notation style computations.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//	result := array.Do(State{})
//
//go:inline
func Do[S any](
	empty S,
) []S {
	return G.Do[[]S, S](empty)
}

// Bind attaches the result of a computation to a context S1 to produce a context S2.
// The setter function defines how to update the context with the computation result.
// This enables monadic composition where each step can produce multiple results.
//
// Example:
//
//	result := F.Pipe2(
//	    array.Do(struct{ X, Y int }{}),
//	    array.Bind(
//	        func(x int) func(s struct{}) struct{ X int } {
//	            return func(s struct{}) struct{ X int } { return struct{ X int }{x} }
//	        },
//	        func(s struct{}) []int { return []int{1, 2} },
//	    ),
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) []T,
) func([]S1) []S2 {
	return G.Bind[[]S1, []S2, []T, S1, S2, T](setter, f)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Unlike Bind, the computation function returns a plain value T rather than []T.
//
// Example:
//
//	result := array.Let(
//	    func(sum int) func(s struct{ X int }) struct{ X, Sum int } {
//	        return func(s struct{ X int }) struct{ X, Sum int } {
//	            return struct{ X, Sum int }{s.X, sum}
//	        }
//	    },
//	    func(s struct{ X int }) int { return s.X * 2 },
//	)
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func([]S1) []S2 {
	return G.Let[[]S1, []S2, S1, S2, T](setter, f)
}

// LetTo attaches a constant value to a context S1 to produce a context S2.
// This is useful for adding constant values to the context.
//
// Example:
//
//	result := array.LetTo(
//	    func(name string) func(s struct{ X int }) struct{ X int; Name string } {
//	        return func(s struct{ X int }) struct{ X int; Name string } {
//	            return struct{ X int; Name string }{s.X, name}
//	        }
//	    },
//	    "constant",
//	)
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func([]S1) []S2 {
	return G.LetTo[[]S1, []S2, S1, S2, T](setter, b)
}

// BindTo initializes a new state S1 from a value T.
// This is typically the first operation after Do to start building the context.
//
// Example:
//
//	result := F.Pipe2(
//	    []int{1, 2, 3},
//	    array.BindTo(func(x int) struct{ X int } {
//	        return struct{ X int }{x}
//	    }),
//	)
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) func([]T) []S1 {
	return G.BindTo[[]S1, []T, S1, T](setter)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering
// the context and the value concurrently (using applicative semantics).
// This produces all combinations of context values and array values.
//
// Example:
//
//	result := array.ApS(
//	    func(y int) func(s struct{ X int }) struct{ X, Y int } {
//	        return func(s struct{ X int }) struct{ X, Y int } {
//	            return struct{ X, Y int }{s.X, y}
//	        }
//	    },
//	    []int{10, 20},
//	)
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa []T,
) func([]S1) []S2 {
	return G.ApS[[]S1, []S2, []T, S1, S2, T](setter, fa)
}
