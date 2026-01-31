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

package decode

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/optics/lens"
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
func Do[I, S any](
	empty S,
) Decode[I, S] {
	return Of[I](empty)
}

// Bind attaches the result of a computation to a context S1 to produce a context S2.
// This is used in do-notation style to sequentially build up a context.
//
// Example:
//
//	type State struct { x int; y int }
//	decoder := F.Pipe2(
//	    Do[string](State{}),
//	    Bind(func(x int) func(State) State {
//	        return func(s State) State { s.x = x; return s }
//	    }, func(s State) Decode[string, int] {
//	        return Of[string](42)
//	    }),
//	)
//	result := decoder("input") // Returns validation.Success(State{x: 42})
func Bind[I, S1, S2, A any](
	setter func(A) func(S1) S2,
	f Kleisli[I, S1, A],
) Operator[I, S1, S2] {
	return C.Bind(
		Chain[I, S1, S2],
		Map[I, A, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Unlike Bind, the computation function returns a plain value, not wrapped in Decode.
//
// Example:
//
//	type State struct { x int; computed int }
//	decoder := F.Pipe2(
//	    Do[string](State{x: 5}),
//	    Let[string](func(c int) func(State) State {
//	        return func(s State) State { s.computed = c; return s }
//	    }, func(s State) int { return s.x * 2 }),
//	)
//	result := decoder("input") // Returns validation.Success(State{x: 5, computed: 10})
func Let[I, S1, S2, B any](
	key func(B) func(S1) S2,
	f func(S1) B,
) Operator[I, S1, S2] {
	return F.Let(
		Map[I, S1, S2],
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
func LetTo[I, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) Operator[I, S1, S2] {
	return F.LetTo(
		Map[I, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state S1 from a value T.
// This is typically used as the first operation after creating a Decode value.
//
// Example:
//
//	type State struct { value int }
//	decoder := F.Pipe1(
//	    Of[string](42),
//	    BindTo[string](func(x int) State { return State{value: x} }),
//	)
//	result := decoder("input") // Returns validation.Success(State{value: 42})
func BindTo[I, S1, T any](
	setter func(T) S1,
) Operator[I, T, S1] {
	return C.BindTo(
		Map[I, T, S1],
		setter,
	)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering the context and the value concurrently.
// This uses the applicative functor pattern, allowing parallel composition.
//
// IMPORTANT: Unlike Bind which fails fast, ApS aggregates ALL validation errors from both the context
// and the value. If both validations fail, all errors are collected and returned together.
// This is useful for validating multiple independent fields and reporting all errors at once.
//
// Example:
//
//	type State struct { x int; y int }
//	decoder := F.Pipe2(
//	    Do[string](State{}),
//	    ApS(func(x int) func(State) State {
//	        return func(s State) State { s.x = x; return s }
//	    }, Of[string](42)),
//	)
//	result := decoder("input") // Returns validation.Success(State{x: 42})
//
// Error aggregation example:
//
//	// Both decoders fail - errors are aggregated
//	decoder1 := func(input string) Validation[State] {
//	    return validation.Failures[State](/* errors */)
//	}
//	decoder2 := func(input string) Validation[int] {
//	    return validation.Failures[int](/* errors */)
//	}
//	combined := ApS(setter, decoder2)(decoder1)
//	result := combined("input") // Contains BOTH sets of errors
func ApS[I, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Decode[I, T],
) Operator[I, S1, S2] {
	return A.ApS(
		Ap[S2, I, T],
		Map[I, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// IMPORTANT: Like ApS, this function aggregates ALL validation errors. If both the context
// and the value fail validation, all errors are collected and returned together.
// This enables comprehensive error reporting for complex nested structures.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// This eliminates the need to manually write setter functions.
//
// Example:
//
//	type Address struct {
//	    Street string
//	    City   string
//	}
//
//	type Person struct {
//	    Name    string
//	    Address Address
//	}
//
//	// Create a lens for the Address field
//	addressLens := lens.MakeLens(
//	    func(p Person) Address { return p.Address },
//	    func(p Person, a Address) Person { p.Address = a; return p },
//	)
//
//	// Use ApSL to update the address
//	decoder := F.Pipe2(
//	    Of[string](Person{Name: "Alice"}),
//	    ApSL(
//	        addressLens,
//	        Of[string](Address{Street: "Main St", City: "NYC"}),
//	    ),
//	)
//	result := decoder("input") // Returns validation.Success(Person{...})
func ApSL[I, S, T any](
	lens L.Lens[S, T],
	fa Decode[I, T],
) Operator[I, S, S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// a Validation that produces the new value.
//
// Unlike ApSL, BindL uses monadic sequencing, meaning the computation f can depend on
// the current value of the focused field.
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
//	increment := func(v int) Decode[string, int] {
//	    return func(input string) Validation[int] {
//	        if v >= 100 {
//	            return validation.Failures[int](/* errors */)
//	        }
//	        return validation.Success(v + 1)
//	    }
//	}
//
//	decoder := F.Pipe1(
//	    Of[string](Counter{Value: 42}),
//	    BindL(valueLens, increment),
//	)
//	result := decoder("input") // Returns validation.Success(Counter{Value: 43})
func BindL[I, S, T any](
	lens L.Lens[S, T],
	f Kleisli[I, T, T],
) Operator[I, S, S] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in Validation).
//
// This is useful for pure transformations that cannot fail, such as mathematical operations,
// string manipulations, or other deterministic updates.
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
//	decoder := F.Pipe1(
//	    Of[string](Counter{Value: 21}),
//	    LetL(valueLens, double),
//	)
//	result := decoder("input") // Returns validation.Success(Counter{Value: 42})
func LetL[I, S, T any](
	lens L.Lens[S, T],
	f Endomorphism[T],
) Operator[I, S, S] {
	return Let[I](lens.Set, function.Flow2(lens.Get, f))
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
//	decoder := F.Pipe1(
//	    Of[string](Config{Debug: true, Timeout: 30}),
//	    LetToL(debugLens, false),
//	)
//	result := decoder("input") // Returns validation.Success(Config{Debug: false, Timeout: 30})
func LetToL[I, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[I, S, S] {
	return LetTo[I](lens.Set, b)
}
