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

package iter

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates a sequence containing a single element, typically used to start a do-notation chain.
// This is the entry point for monadic composition using do-notation style.
//
// Type Parameters:
//   - S: The type of the state/structure being built
//
// Parameters:
//   - empty: The initial value to wrap in a sequence
//
// Returns:
//   - A sequence containing the single element
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	// Start a do-notation chain
//	result := Do(User{})
//	// yields: User{Name: "", Age: 0}
//
//go:inline
func Do[S any](
	empty S,
) Seq[S] {
	return Of(empty)
}

// Bind performs a monadic bind operation in do-notation style, chaining a computation
// that produces a sequence and updating the state with the result.
//
// This function is the core of do-notation for sequences. It takes a Kleisli arrow
// (a function that returns a sequence) and a setter function that updates the state
// with the result. The setter is curried to allow partial application.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the Kleisli arrow
//
// Parameters:
//   - setter: A curried function that takes a value T and returns a function that updates S1 to S2
//   - f: A Kleisli arrow that takes S1 and produces a sequence of T
//
// Returns:
//   - An Operator that transforms Seq[S1] to Seq[S2]
//
// Example:
//
//	type State struct {
//	    Value int
//	    Double int
//	}
//
//	setValue := func(v int) func(State) State {
//	    return func(s State) State { s.Value = v; return s }
//	}
//
//	getValues := func(s State) Seq[int] {
//	    return From(1, 2, 3)
//	}
//
//	result := F.Pipe2(
//	    Do(State{}),
//	    Bind(setValue, getValues),
//	)
//	// yields: State{Value: 1}, State{Value: 2}, State{Value: 3}
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return C.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let performs a pure computation in do-notation style, updating the state with a computed value.
//
// Unlike Bind, Let doesn't perform a monadic operation - it simply computes a value from
// the current state and updates the state with that value. This is useful for intermediate
// calculations that don't require sequencing.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of the computed value
//
// Parameters:
//   - key: A curried function that takes a value T and returns a function that updates S1 to S2
//   - f: A function that computes T from S1
//
// Returns:
//   - An Operator that transforms Seq[S1] to Seq[S2]
//
// Example:
//
//	type State struct {
//	    Value int
//	    Double int
//	}
//
//	setDouble := func(d int) func(State) State {
//	    return func(s State) State { s.Double = d; return s }
//	}
//
//	computeDouble := func(s State) int {
//	    return s.Value * 2
//	}
//
//	result := F.Pipe2(
//	    Do(State{Value: 5}),
//	    Let(setDouble, computeDouble),
//	)
//	// yields: State{Value: 5, Double: 10}
//
//go:inline
func Let[S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return F.Let(
		Map[S1, S2],
		key,
		f,
	)
}

// LetTo sets a field in the state to a constant value in do-notation style.
//
// This is a specialized version of Let that doesn't compute the value from the state,
// but instead uses a fixed value. It's useful for setting constants or default values.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of the value to set
//
// Parameters:
//   - key: A curried function that takes a value T and returns a function that updates S1 to S2
//   - b: The constant value to set
//
// Returns:
//   - An Operator that transforms Seq[S1] to Seq[S2]
//
// Example:
//
//	type State struct {
//	    Name string
//	    Status string
//	}
//
//	setStatus := func(s string) func(State) State {
//	    return func(st State) State { st.Status = s; return st }
//	}
//
//	result := F.Pipe2(
//	    Do(State{Name: "Alice"}),
//	    LetTo(setStatus, "active"),
//	)
//	// yields: State{Name: "Alice", Status: "active"}
//
//go:inline
func LetTo[S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return F.LetTo(
		Map[S1, S2],
		key,
		b,
	)
}

// BindTo wraps a value into a structure using a setter function.
//
// This is typically used at the beginning of a do-notation chain to convert a simple
// value into a structured state. It's the inverse of extracting a value from a structure.
//
// Type Parameters:
//   - S1: The structure type to create
//   - T: The value type to wrap
//
// Parameters:
//   - setter: A function that takes a value T and creates a structure S1
//
// Returns:
//   - An Operator that transforms Seq[T] to Seq[S1]
//
// Example:
//
//	type State struct {
//	    Value int
//	}
//
//	createState := func(v int) State {
//	    return State{Value: v}
//	}
//
//	result := F.Pipe2(
//	    From(1, 2, 3),
//	    BindTo(createState),
//	)
//	// yields: State{Value: 1}, State{Value: 2}, State{Value: 3}
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return C.BindTo(
		Map[T, S1],
		setter,
	)
}

// BindToP wraps a value into a structure using a Prism's ReverseGet function.
//
// This is a specialized version of BindTo that works with Prisms (optics that focus
// on a case of a sum type). It uses the Prism's ReverseGet to construct the structure.
//
// Type Parameters:
//   - S1: The structure type to create
//   - T: The value type to wrap
//
// Parameters:
//   - setter: A Prism that can construct S1 from T
//
// Returns:
//   - An Operator that transforms Seq[T] to Seq[S1]
//
// Example:
//
//	// Assuming a Prism for wrapping int into a Result type
//	result := F.Pipe2(
//	    From(1, 2, 3),
//	    BindToP(successPrism),
//	)
//	// yields: Success(1), Success(2), Success(3)
//
//go:inline
func BindToP[S1, T any](
	setter Prism[S1, T],
) Operator[T, S1] {
	return BindTo(setter.ReverseGet)
}

// ApS applies a sequence of values to update a state using applicative style.
//
// This function combines applicative application with state updates. It takes a sequence
// of values and a setter function, and produces an operator that applies each value
// to update the state. This is useful for parallel composition of independent computations.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of values in the sequence
//
// Parameters:
//   - setter: A curried function that takes a value T and returns a function that updates S1 to S2
//   - fa: A sequence of values to apply
//
// Returns:
//   - An Operator that transforms Seq[S1] to Seq[S2]
//
// Example:
//
//	type State struct {
//	    X int
//	    Y int
//	}
//
//	setY := func(y int) func(State) State {
//	    return func(s State) State { s.Y = y; return s }
//	}
//
//	yValues := From(10, 20, 30)
//
//	result := F.Pipe2(
//	    Do(State{X: 5}),
//	    ApS(setY, yValues),
//	)
//	// yields: State{X: 5, Y: 10}, State{X: 5, Y: 20}, State{X: 5, Y: 30}
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Seq[T],
) Operator[S1, S2] {
	return A.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL applies a sequence of values to update a state field using a Lens.
//
// This is a specialized version of ApS that works with Lenses (optics that focus on
// a field of a structure). It uses the Lens's Set function to update the field.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field being updated
//
// Parameters:
//   - lens: A Lens focusing on the field to update
//   - fa: A sequence of values to set
//
// Returns:
//   - An Endomorphism on Seq[S] (transforms Seq[S] to Seq[S])
//
// Example:
//
//	type State struct {
//	    Name string
//	    Age  int
//	}
//
//	ageLens := lens.Prop[State, int]("Age")
//	ages := From(25, 30, 35)
//
//	result := F.Pipe2(
//	    Do(State{Name: "Alice"}),
//	    ApSL(ageLens, ages),
//	)
//	// yields: State{Name: "Alice", Age: 25}, State{Name: "Alice", Age: 30}, State{Name: "Alice", Age: 35}
//
//go:inline
func ApSL[S, T any](
	lens Lens[S, T],
	fa Seq[T],
) Endomorphism[Seq[S]] {
	return ApS(lens.Set, fa)
}

// BindL performs a monadic bind on a field of a structure using a Lens.
//
// This function combines Lens-based field access with monadic binding. It extracts
// a field value using the Lens's Get, applies a Kleisli arrow to produce a sequence,
// and updates the field with each result using the Lens's Set.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field being accessed and updated
//
// Parameters:
//   - lens: A Lens focusing on the field to bind
//   - f: A Kleisli arrow that takes the field value and produces a sequence
//
// Returns:
//   - An Endomorphism on Seq[S]
//
// Example:
//
//	type State struct {
//	    Value int
//	}
//
//	valueLens := lens.Prop[State, int]("Value")
//
//	multiplyValues := func(v int) Seq[int] {
//	    return From(v, v*2, v*3)
//	}
//
//	result := F.Pipe2(
//	    Do(State{Value: 5}),
//	    BindL(valueLens, multiplyValues),
//	)
//	// yields: State{Value: 5}, State{Value: 10}, State{Value: 15}
//
//go:inline
func BindL[S, T any](
	lens Lens[S, T],
	f Kleisli[T, T],
) Endomorphism[Seq[S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL performs a pure computation on a field of a structure using a Lens.
//
// This function extracts a field value using the Lens's Get, applies a pure function
// to compute a new value, and updates the field using the Lens's Set. It's useful
// for transforming fields without monadic effects.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field being transformed
//
// Parameters:
//   - lens: A Lens focusing on the field to transform
//   - f: An Endomorphism that transforms the field value
//
// Returns:
//   - An Endomorphism on Seq[S]
//
// Example:
//
//	type State struct {
//	    Count int
//	}
//
//	countLens := lens.Prop[State, int]("Count")
//
//	increment := func(n int) int { return n + 1 }
//
//	result := F.Pipe2(
//	    Do(State{Count: 5}),
//	    LetL(countLens, increment),
//	)
//	// yields: State{Count: 6}
//
//go:inline
func LetL[S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[Seq[S]] {
	return Let(lens.Set, function.Flow2(lens.Get, f))
}

// LetToL sets a field of a structure to a constant value using a Lens.
//
// This is a specialized version of LetL that sets a field to a fixed value rather
// than computing it from the current value. It's useful for setting defaults or
// resetting fields.
//
// Type Parameters:
//   - S: The state type
//   - T: The type of the field being set
//
// Parameters:
//   - lens: A Lens focusing on the field to set
//   - b: The constant value to set
//
// Returns:
//   - An Endomorphism on Seq[S]
//
// Example:
//
//	type State struct {
//	    Status string
//	}
//
//	statusLens := lens.Prop[State, string]("Status")
//
//	result := F.Pipe2(
//	    Do(State{Status: "pending"}),
//	    LetToL(statusLens, "active"),
//	)
//	// yields: State{Status: "active"}
//
//go:inline
func LetToL[S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[Seq[S]] {
	return LetTo(lens.Set, b)
}
