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

package result

import (
	"github.com/IBM/fp-go/v2/either"
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
func Do[S any](
	empty S,
) Result[S] {
	return either.Do[error](empty)
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
//	        func(s State) either.Result[int] {
//	            return either.Right[error](42)
//	        },
//	    ),
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return either.Bind(setter, f)
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
func Let[S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return either.Let[error](key, f)
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
func LetTo[S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return either.LetTo[error](key, b)
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
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return either.BindTo[error](setter)
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
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[S1, S2] {
	return either.ApS(setter, fa)
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
//   - fa: An Result[T] computation that produces the value to set
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
func ApSL[S, T any](
	lens Lens[S, T],
	fa Result[T],
) Operator[S, S] {
	return either.ApSL(lens, fa)
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
//   - f: A function that takes the current field value and returns an Result[T]
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
//	increment := func(v int) either.Result[int] {
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
func BindL[S, T any](
	lens Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return either.BindL(lens, f)
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
func LetL[S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Operator[S, S] {
	return either.LetL[error](lens, f)
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
func LetToL[S, T any](
	lens Lens[S, T],
	b T,
) Operator[S, S] {
	return either.LetToL[error](lens, b)
}
