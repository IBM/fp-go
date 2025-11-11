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

package either

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type S to be used with the Bind operation.
// This is the starting point for do-notation style computations.
//
// Example:
//
//	type State struct { x, y int }
//	result := either.Do[error](State{})
//
//go:inline
func Do[E, S any](
	empty S,
) Either[E, S] {
	return Of[E](empty)
}

// Bind attaches the result of a computation to a context S1 to produce a context S2.
// This enables building up complex computations in a pipeline.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Do[error](State{}),
//	    either.Bind(
//	        func(v int) func(State) State {
//	            return func(s State) State { return State{value: v} }
//	        },
//	        func(s State) either.Either[error, int] {
//	            return either.Right[error](42)
//	        },
//	    ),
//	)
//
//go:inline
func Bind[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[E, S1, T],
) Operator[E, S1, S2] {
	return C.Bind(
		Chain[E, S1, S2],
		Map[E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Similar to Bind but for pure (non-Either) computations.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Right[error](State{value: 10}),
//	    either.Let(
//	        func(v int) func(State) State {
//	            return func(s State) State { return State{value: s.value + v} }
//	        },
//	        func(s State) int { return 32 },
//	    ),
//	) // Right(State{value: 42})
//
//go:inline
func Let[E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[E, S1, S2] {
	return F.Let(
		Map[E, S1, S2],
		key,
		f,
	)
}

// LetTo attaches a constant value to a context S1 to produce a context S2.
//
// Example:
//
//	type State struct { name string }
//	result := F.Pipe2(
//	    either.Right[error](State{}),
//	    either.LetTo(
//	        func(n string) func(State) State {
//	            return func(s State) State { return State{name: n} }
//	        },
//	        "Alice",
//	    ),
//	) // Right(State{name: "Alice"})
//
//go:inline
func LetTo[E, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[E, S1, S2] {
	return F.LetTo(
		Map[E, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state S1 from a value T.
// This is typically used to start a bind chain.
//
// Example:
//
//	type State struct { value int }
//	result := F.Pipe2(
//	    either.Right[error](42),
//	    either.BindTo(func(v int) State { return State{value: v} }),
//	) // Right(State{value: 42})
//
//go:inline
func BindTo[E, S1, T any](
	setter func(T) S1,
) Operator[E, T, S1] {
	return C.BindTo(
		Map[E, T, S1],
		setter,
	)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering the context and the value concurrently.
// Uses applicative semantics rather than monadic sequencing.
//
// Example:
//
//	type State struct { x, y int }
//	result := F.Pipe2(
//	    either.Right[error](State{x: 10}),
//	    either.ApS(
//	        func(y int) func(State) State {
//	            return func(s State) State { return State{x: s.x, y: y} }
//	        },
//	        either.Right[error](32),
//	    ),
//	) // Right(State{x: 10, y: 32})
//
//go:inline
func ApS[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[E, T],
) Operator[E, S1, S2] {
	return A.ApS(
		Ap[S2, E, T],
		Map[E, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// This eliminates the need to manually write setter functions and enables working with
// nested fields in a type-safe manner.
//
// Unlike BindL, ApSL uses applicative semantics, meaning the computation fa is independent
// of the current state and can be evaluated concurrently.
//
// Type Parameters:
//   - E: Error type for the Either
//   - S: Structure type containing the field to update
//   - T: Type of the field being updated
//
// Parameters:
//   - lens: A Lens[S, T] that focuses on a field of type T within structure S
//   - fa: An Either[E, T] computation that produces the value to set
//
// Returns:
//   - An endomorphism that updates the focused field in the Either context
//
// Example:
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	ageLens := lens.MakeLens(
//	    func(p Person) int { return p.Age },
//	    func(p Person, a int) Person { p.Age = a; return p },
//	)
//
//	result := F.Pipe2(
//	    either.Right[error](Person{Name: "Alice", Age: 25}),
//	    either.ApSL(ageLens, either.Right[error](30)),
//	) // Right(Person{Name: "Alice", Age: 30})
//
//go:inline
func ApSL[E, S, T any](
	lens Lens[S, T],
	fa Either[E, T],
) Endomorphism[Either[E, S]] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an Either that produces the new value.
//
// Unlike ApSL, BindL uses monadic sequencing, meaning the computation f can depend on
// the current value of the focused field.
//
// Type Parameters:
//   - E: Error type for the Either
//   - S: Structure type containing the field to update
//   - T: Type of the field being updated
//
// Parameters:
//   - lens: A Lens[S, T] that focuses on a field of type T within structure S
//   - f: A function that takes the current field value and returns an Either[E, T]
//
// Returns:
//   - An endomorphism that updates the focused field based on its current value
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//
//	// Increment the counter, but fail if it would exceed 100
//	increment := func(v int) either.Either[error, int] {
//	    if v >= 100 {
//	        return either.Left[int](errors.New("counter overflow"))
//	    }
//	    return either.Right[error](v + 1)
//	}
//
//	result := F.Pipe1(
//	    either.Right[error](Counter{Value: 42}),
//	    either.BindL(valueLens, increment),
//	) // Right(Counter{Value: 43})
//
//go:inline
func BindL[E, S, T any](
	lens Lens[S, T],
	f Kleisli[E, T, T],
) Endomorphism[Either[E, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in Either).
//
// This is useful for pure transformations that cannot fail, such as mathematical operations,
// string manipulations, or other deterministic updates.
//
// Type Parameters:
//   - E: Error type for the Either
//   - S: Structure type containing the field to update
//   - T: Type of the field being updated
//
// Parameters:
//   - lens: A Lens[S, T] that focuses on a field of type T within structure S
//   - f: An endomorphism (T → T) that transforms the current field value
//
// Returns:
//   - An endomorphism that updates the focused field with the transformed value
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//
//	// Double the counter value
//	double := func(v int) int { return v * 2 }
//
//	result := F.Pipe1(
//	    either.Right[error](Counter{Value: 21}),
//	    either.LetL(valueLens, double),
//	) // Right(Counter{Value: 42})
//
//go:inline
func LetL[E, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[Either[E, S]] {
	return Let[E](lens.Set, function.Flow2(lens.Get, f))
}

// LetToL attaches a constant value to a context using a lens-based setter.
// This is a convenience function that combines LetTo with a lens, allowing you to use
// optics to set nested fields to specific values.
//
// The lens parameter provides the setter for a field within the structure S.
// Unlike LetL which transforms the current value, LetToL simply replaces it with
// the provided constant value b.
//
// This is useful for resetting fields, initializing values, or setting fields to
// predetermined constants.
//
// Type Parameters:
//   - E: Error type for the Either
//   - S: Structure type containing the field to update
//   - T: Type of the field being updated
//
// Parameters:
//   - lens: A Lens[S, T] that focuses on a field of type T within structure S
//   - b: The constant value to set the field to
//
// Returns:
//   - An endomorphism that sets the focused field to the constant value
//
// Example:
//
//	type Config struct {
//	    Debug   bool
//	    Timeout int
//	}
//
//	debugLens := lens.MakeLens(
//	    func(c Config) bool { return c.Debug },
//	    func(c Config, d bool) Config { c.Debug = d; return c },
//	)
//
//	result := F.Pipe1(
//	    either.Right[error](Config{Debug: true, Timeout: 30}),
//	    either.LetToL(debugLens, false),
//	) // Right(Config{Debug: false, Timeout: 30})
//
//go:inline
func LetToL[E, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[Either[E, S]] {
	return LetTo[E](lens.Set, b)
}
