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

package iooption

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    Name  string
//	    Age   int
//	}
//	result := iooption.Do(State{})
func Do[S any](
	empty S,
) IOOption[S] {
	return Of(empty)
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
//	    Name  string
//	    Age   int
//	}
//
//	result := F.Pipe2(
//	    iooption.Do(State{}),
//	    iooption.Bind(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        func(s State) iooption.IOOption[string] {
//	            return iooption.FromIO(io.Of("Alice"))
//	        },
//	    ),
//	    iooption.Bind(
//	        func(age int) func(State) State {
//	            return func(s State) State { s.Age = age; return s }
//	        },
//	        func(s State) iooption.IOOption[int] {
//	            // This can access s.Name from the previous step
//	            return iooption.FromIO(io.Of(len(s.Name) * 10))
//	        },
//	    ),
//	)
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return chain.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return functor.Let(
		Map[S1, S2],
		setter,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return functor.LetTo(
		Map[S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return chain.BindTo(
		Map[T, S1],
		setter,
	)
}

//go:inline
func BindToP[S1, T any](
	setter Prism[S1, T],
) Operator[T, S1] {
	return BindTo(setter.ReverseGet)
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
//	    Age   int
//	}
//
//	// These operations are independent and can be combined with ApS
//	getName := iooption.Some("Alice")
//	getAge := iooption.Some(30)
//
//	result := F.Pipe2(
//	    iooption.Do(State{}),
//	    iooption.ApS(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        getName,
//	    ),
//	    iooption.ApS(
//	        func(age int) func(State) State {
//	            return func(s State) State { s.Age = age; return s }
//	        },
//	        getAge,
//	    ),
//	)
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOOption[T],
) Operator[S1, S2] {
	return apply.ApS(
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
//	type State struct {
//	    Name  string
//	    Age   int
//	}
//
//	ageLens := lens.MakeLens(
//	    func(s State) int { return s.Age },
//	    func(s State, a int) State { s.Age = a; return s },
//	)
//
//	result := F.Pipe2(
//	    iooption.Of(State{Name: "Alice"}),
//	    iooption.ApSL(ageLens, iooption.Some(30)),
//	)
func ApSL[S, T any](
	lens Lens[S, T],
	fa IOOption[T],
) Operator[S, S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an IOOption that produces the new value.
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
//	increment := func(v int) iooption.IOOption[int] {
//	    return iooption.FromIO(io.Of(v + 1))
//	}
//
//	result := F.Pipe1(
//	    iooption.Of(Counter{Value: 42}),
//	    iooption.BindL(valueLens, increment),
//	) // IOOption[Counter{Value: 43}]
func BindL[S, T any](
	lens Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in IOOption).
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
//	    iooption.Of(Counter{Value: 21}),
//	    iooption.LetL(valueLens, double),
//	) // IOOption[Counter{Value: 42}]
func LetL[S, T any](
	lens Lens[S, T],
	f func(T) T,
) Operator[S, S] {
	return Let(lens.Set, F.Flow2(lens.Get, f))
}

// LetToL attaches a constant value to a context using a lens-based setter.
// This is a convenience function that combines LetTo with a lens, allowing you to use
// optics to set nested fields to specific values.
//
// The lens parameter provides the setter for a field within the structure S.
// Unlike LetL which transforms the current value, LetToL simply replaces it with
// the provided constant value b.
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
//	    iooption.Of(Config{Debug: true, Timeout: 30}),
//	    iooption.LetToL(debugLens, false),
//	) // IOOption[Config{Debug: false, Timeout: 30}]
func LetToL[S, T any](
	lens Lens[S, T],
	b T,
) Operator[S, S] {
	return LetTo(lens.Set, b)
}
