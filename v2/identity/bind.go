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

package identity

import (
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) S {
	return empty
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(S1) S2 {
	return C.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(S1) S2 {
	return F.Let(
		Map[S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(S1) S2 {
	return F.LetTo(
		Map[S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(T) S1 {
	return C.BindTo(
		Map[T, S1],
		setter,
	)
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
//	result := F.Pipe2(
//	    identity.Do(State{}),
//	    identity.ApS(
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        42,
//	    ),
//	    identity.ApS(
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        100,
//	    ),
//	) // State{X: 42, Y: 100}
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa T,
) func(S1) S2 {
	return A.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}
