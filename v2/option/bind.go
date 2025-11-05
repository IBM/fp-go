// Copyright (c) 2025 IBM Corp.
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

package option

import (
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type S to be used with the Bind operation.
// This is the starting point for building up a context using do-notation style.
//
// Example:
//
//	type Result struct {
//	    x int
//	    y string
//	}
//	result := Do(Result{})
func Do[S any](
	empty S,
) Option[S] {
	return Of(empty)
}

// Bind attaches the result of a computation to a context S1 to produce a context S2.
// This is used in do-notation style to sequentially build up a context.
//
// Example:
//
//	type State struct { x int; y int }
//	result := F.Pipe2(
//	    Do(State{}),
//	    Bind(func(x int) func(State) State {
//	        return func(s State) State { s.x = x; return s }
//	    }, func(s State) Option[int] { return Some(42) }),
//	)
func Bind[S1, S2, A any](
	setter func(A) func(S1) S2,
	f func(S1) Option[A],
) func(Option[S1]) Option[S2] {
	return C.Bind(
		Chain[S1, S2],
		Map[A, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Unlike Bind, the computation function returns a plain value, not an Option.
//
// Example:
//
//	type State struct { x int; computed int }
//	result := F.Pipe2(
//	    Do(State{x: 5}),
//	    Let(func(c int) func(State) State {
//	        return func(s State) State { s.computed = c; return s }
//	    }, func(s State) int { return s.x * 2 }),
//	)
func Let[S1, S2, B any](
	key func(B) func(S1) S2,
	f func(S1) B,
) func(Option[S1]) Option[S2] {
	return F.Let(
		Map[S1, S2],
		key,
		f,
	)
}

// LetTo attaches a constant value to a context S1 to produce a context S2.
//
// Example:
//
//	type State struct { x int; name string }
//	result := F.Pipe2(
//	    Do(State{x: 5}),
//	    LetTo(func(n string) func(State) State {
//	        return func(s State) State { s.name = n; return s }
//	    }, "example"),
//	)
func LetTo[S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(Option[S1]) Option[S2] {
	return F.LetTo(
		Map[S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state S1 from a value T.
// This is typically used as the first operation after creating an Option value.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe1(
//	    Some(42),
//	    BindTo(func(x int) State { return State{value: x} }),
//	)
func BindTo[S1, T any](
	setter func(T) S1,
) func(Option[T]) Option[S1] {
	return C.BindTo(
		Map[T, S1],
		setter,
	)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering the context and the value concurrently.
// This uses the applicative functor pattern, allowing parallel composition.
//
// Example:
//
//	type State struct { x int; y int }
//	result := F.Pipe2(
//	    Do(State{}),
//	    ApS(func(x int) func(State) State {
//	        return func(s State) State { s.x = x; return s }
//	    }, Some(42)),
//	)
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Option[T],
) func(Option[S1]) Option[S2] {
	return A.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}
