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

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//	result := identity.Do(State{})
func Do[S any](
	empty S,
) S {
	return empty
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
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
//	    identity.Do(State{}),
//	    identity.Bind(
//	        func(x int) func(State) State {
//	            return func(s State) State { s.X = x; return s }
//	        },
//	        func(s State) int {
//	            return 42
//	        },
//	    ),
//	    identity.Bind(
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        func(s State) int {
//	            // This can access s.X from the previous step
//	            return s.X * 2
//	        },
//	    ),
//	) // State{X: 42, Y: 84}
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

// Let attaches the result of a computation to a context [S1] to produce a context [S2].
// Similar to Bind, but uses the Functor's Map operation instead of the Monad's Chain.
// This is useful when you want to add a computed value to the context without needing
// the full power of monadic composition.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	    Sum int
//	}
//
//	result := F.Pipe2(
//	    identity.Do(State{X: 10, Y: 20}),
//	    identity.Let(
//	        func(sum int) func(State) State {
//	            return func(s State) State { s.Sum = sum; return s }
//	        },
//	        func(s State) int {
//	            return s.X + s.Y
//	        },
//	    ),
//	) // State{X: 10, Y: 20, Sum: 30}
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

// LetTo attaches a constant value to a context [S1] to produce a context [S2].
// This is a specialized version of Let that doesn't require a computation function,
// useful when you want to add a known value to the context.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	    Constant string
//	}
//
//	result := F.Pipe2(
//	    identity.Do(State{X: 10, Y: 20}),
//	    identity.LetTo(
//	        func(c string) func(State) State {
//	            return func(s State) State { s.Constant = c; return s }
//	        },
//	        "fixed value",
//	    ),
//	) // State{X: 10, Y: 20, Constant: "fixed value"}
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

// BindTo initializes a new state [S1] from a value [T].
// This is typically used as the first operation in a do-notation chain to convert
// a plain value into a context that can be used with subsequent Bind operations.
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//
//	result := F.Pipe2(
//	    42,
//	    identity.BindTo(func(x int) State {
//	        return State{X: x}
//	    }),
//	    identity.Bind(
//	        func(y int) func(State) State {
//	            return func(s State) State { s.Y = y; return s }
//	        },
//	        func(s State) int {
//	            return s.X * 2
//	        },
//	    ),
//	) // State{X: 42, Y: 84}
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
