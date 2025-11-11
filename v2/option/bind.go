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
	f Kleisli[S1, A],
) Kleisli[Option[S1], S2] {
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
) Kleisli[Option[S1], S2] {
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
) Kleisli[Option[S1], S2] {
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
) Kleisli[Option[T], S1] {
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
) Kleisli[Option[S1], S2] {
	return A.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
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
//	result := F.Pipe2(
//	    option.Some(Person{Name: "Alice"}),
//	    option.ApSL(
//	        addressLens,
//	        option.Some(Address{Street: "Main St", City: "NYC"}),
//	    ),
//	)
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa Option[T],
) Kleisli[Option[S], S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an Option that produces the new value.
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
//	// Increment the counter, but return None if it would exceed 100
//	increment := func(v int) option.Option[int] {
//	    if v >= 100 {
//	        return option.None[int]()
//	    }
//	    return option.Some(v + 1)
//	}
//
//	result := F.Pipe1(
//	    option.Some(Counter{Value: 42}),
//	    option.BindL(valueLens, increment),
//	) // Some(Counter{Value: 43})
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Kleisli[Option[S], S] {
	return Bind[S, S, T](lens.Set, function.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in Option).
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
//	result := F.Pipe1(
//	    option.Some(Counter{Value: 21}),
//	    option.LetL(valueLens, double),
//	) // Some(Counter{Value: 42})
func LetL[S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Kleisli[Option[S], S] {
	return Let[S, S, T](lens.Set, function.Flow2(lens.Get, f))
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
//	result := F.Pipe1(
//	    option.Some(Config{Debug: true, Timeout: 30}),
//	    option.LetToL(debugLens, false),
//	) // Some(Config{Debug: false, Timeout: 30})
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Kleisli[Option[S], S] {
	return LetTo[S, S, T](lens.Set, b)
}
