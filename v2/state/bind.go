// Copyright (c) 2024 - 2025 IBM Corp.
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

package state

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do initializes a do-notation computation with an empty value.
// This is the entry point for building complex stateful computations using
// the do-notation pattern, which allows for imperative-style sequencing of
// monadic operations.
//
// The do-notation pattern is useful for building pipelines where you need to:
//   - Bind intermediate results to names
//   - Sequence multiple stateful operations
//   - Build up complex state transformations step by step
//
// Example:
//
//	type MyState struct {
//	    x int
//	    y int
//	}
//
//	type Result struct {
//	    sum int
//	    product int
//	}
//
//	computation := function.Pipe3(
//	    Do[MyState](Result{}),
//	    Bind(func(r Result) func(int) Result {
//	        return func(x int) Result { r.sum = x; return r }
//	    }, Gets(func(s MyState) int { return s.x })),
//	    Bind(func(r Result) func(int) Result {
//	        return func(y int) Result { r.product = r.sum * y; return r }
//	    }, Gets(func(s MyState) int { return s.y })),
//	    Map[MyState](func(r Result) int { return r.product }),
//	)
//
//go:inline
func Do[ST, A any](
	empty A,
) State[ST, A] {
	return Of[ST](empty)
}

// Bind sequences a stateful computation and binds its result to a field in an
// accumulator structure. This is a key building block for do-notation, allowing
// you to extract values from State computations and incorporate them into a
// growing result structure.
//
// The setter function takes the computed value T and returns a function that
// updates the accumulator from S1 to S2 by setting the field to T.
//
// Parameters:
//   - setter: A function that takes a value T and returns a function to update
//     the accumulator structure from S1 to S2
//   - f: A Kleisli arrow that takes the current accumulator S1 and produces a
//     State computation yielding T
//
// Example:
//
//	type Accumulator struct {
//	    value int
//	    doubled int
//	}
//
//	// Bind the result of a computation to the 'doubled' field
//	computation := Bind(
//	    func(d int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.doubled = d
//	            return acc
//	        }
//	    },
//	    func(acc Accumulator) State[MyState, int] {
//	        return Of[MyState](acc.value * 2)
//	    },
//	)
//
//go:inline
func Bind[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[ST, S1, T],
) Operator[ST, S1, S2] {
	return C.Bind(
		Chain[ST, Kleisli[ST, S1, S2], S1, S2],
		Map[ST, func(T) S2, T, S2],
		setter,
		f,
	)
}

// Let computes a pure value from the current accumulator and binds it to a field.
// Unlike Bind, this doesn't execute a stateful computation - it simply applies a
// pure function to the accumulator and stores the result.
//
// This is useful in do-notation when you need to compute derived values without
// performing stateful operations.
//
// Parameters:
//   - key: A function that takes the computed value T and returns a function to
//     update the accumulator from S1 to S2
//   - f: A pure function that extracts or computes T from the current accumulator S1
//
// Example:
//
//	type Accumulator struct {
//	    x int
//	    y int
//	    sum int
//	}
//
//	// Compute sum from x and y without state operations
//	computation := Let(
//	    func(s int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.sum = s
//	            return acc
//	        }
//	    },
//	    func(acc Accumulator) int {
//	        return acc.x + acc.y
//	    },
//	)
//
//go:inline
func Let[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) Operator[ST, S1, S2] {
	return F.Let(
		Map[ST, func(S1) S2, S1, S2],
		key,
		f,
	)
}

// LetTo binds a constant value to a field in the accumulator.
// This is a specialized version of Let where the value is already known
// and doesn't need to be computed from the accumulator.
//
// This is useful for initializing fields with constant values or for
// setting default values in do-notation pipelines.
//
// Parameters:
//   - key: A function that takes the constant value T and returns a function to
//     update the accumulator from S1 to S2
//   - b: The constant value to bind
//
// Example:
//
//	type Accumulator struct {
//	    status string
//	    value int
//	}
//
//	// Set a constant status
//	computation := LetTo(
//	    func(s string) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.status = s
//	            return acc
//	        }
//	    },
//	    "initialized",
//	)
//
//go:inline
func LetTo[ST, S1, S2, T any](
	key func(T) func(S1) S2,
	b T,
) Operator[ST, S1, S2] {
	return F.LetTo(
		Map[ST, func(S1) S2, S1, S2],
		key,
		b,
	)
}

// BindTo creates an initial accumulator structure from a value.
// This is typically the first operation in a do-notation pipeline,
// converting a simple value into a structure that can accumulate
// additional fields.
//
// Parameters:
//   - setter: A function that takes a value T and creates an accumulator structure S1
//
// Example:
//
//	type Accumulator struct {
//	    initial int
//	    doubled int
//	}
//
//	// Start a pipeline by binding the initial value
//	computation := function.Pipe2(
//	    Of[MyState](42),
//	    BindTo(func(x int) Accumulator {
//	        return Accumulator{initial: x}
//	    }),
//	    Bind(func(d int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.doubled = d
//	            return acc
//	        }
//	    }, func(acc Accumulator) State[MyState, int] {
//	        return Of[MyState](acc.initial * 2)
//	    }),
//	)
//
//go:inline
func BindTo[ST, S1, T any](
	setter func(T) S1,
) Operator[ST, T, S1] {
	return C.BindTo(
		Map[ST, func(T) S1, T, S1],
		setter,
	)
}

// ApS applies a State computation in an applicative style and binds the result
// to a field in the accumulator. Unlike Bind, which uses monadic sequencing,
// ApS uses applicative composition, which can be more efficient when the
// computation doesn't depend on the accumulator value.
//
// This is useful when you have independent State computations that can be
// composed without depending on each other's results.
//
// Parameters:
//   - setter: A function that takes the computed value T and returns a function to
//     update the accumulator from S1 to S2
//   - fa: A State computation that produces a value of type T
//
// Example:
//
//	type Accumulator struct {
//	    counter int
//	    timestamp int64
//	}
//
//	// Apply an independent state computation
//	getTimestamp := Gets(func(s MyState) int64 { return s.timestamp })
//
//	computation := ApS(
//	    func(ts int64) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.timestamp = ts
//	            return acc
//	        }
//	    },
//	    getTimestamp,
//	)
//
//go:inline
func ApS[ST, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa State[ST, T],
) Operator[ST, S1, S2] {
	return A.ApS(
		Ap[S2, ST, T],
		Map[ST, func(S1) func(T) S2, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a lens-based version of ApS that uses a lens to focus on a specific
// field in the accumulator structure. This provides a more convenient and
// type-safe way to update nested fields.
//
// A lens provides both a getter and setter for a field, making it easier to
// work with complex data structures without manually writing setter functions.
//
// Parameters:
//   - lens: A lens focusing on field T within structure S
//   - fa: A State computation that produces a value of type T
//
// Example:
//
//	type MyState struct {
//	    counter int
//	}
//
//	type Accumulator struct {
//	    value int
//	    doubled int
//	}
//
//	// Create a lens for the 'doubled' field
//	doubledLens := MakeLens(
//	    func(acc Accumulator) int { return acc.doubled },
//	    func(d int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.doubled = d
//	            return acc
//	        }
//	    },
//	)
//
//	computation := ApSL(doubledLens, Of[MyState](42))
//
//go:inline
func ApSL[ST, S, T any](
	lens Lens[S, T],
	fa State[ST, T],
) Endomorphism[State[ST, S]] {
	return ApS(lens.Set, fa)
}

// BindL is a lens-based version of Bind that focuses on a specific field,
// extracts its value, applies a stateful computation, and updates the field
// with the result. This is particularly useful for updating nested fields
// based on their current values.
//
// The computation receives the current value of the focused field and produces
// a new value through a State computation.
//
// Parameters:
//   - lens: A lens focusing on field T within structure S
//   - f: A Kleisli arrow that takes the current field value and produces a
//     State computation yielding the new value
//
// Example:
//
//	type MyState struct {
//	    multiplier int
//	}
//
//	type Accumulator struct {
//	    value int
//	}
//
//	valueLens := MakeLens(
//	    func(acc Accumulator) int { return acc.value },
//	    func(v int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.value = v
//	            return acc
//	        }
//	    },
//	)
//
//	// Double the value using state
//	computation := BindL(
//	    valueLens,
//	    func(v int) State[MyState, int] {
//	        return Gets(func(s MyState) int {
//	            return v * s.multiplier
//	        })
//	    },
//	)
//
//go:inline
func BindL[ST, S, T any](
	lens Lens[S, T],
	f Kleisli[ST, T, T],
) Endomorphism[State[ST, S]] {
	return Bind(lens.Set, function.Flow2(lens.Get, f))
}

// LetL is a lens-based version of Let that focuses on a specific field,
// extracts its value, applies a pure transformation, and updates the field
// with the result. This is useful for pure transformations of nested fields.
//
// Unlike BindL, this doesn't perform stateful computations - it only applies
// a pure function to the field value.
//
// Parameters:
//   - lens: A lens focusing on field T within structure S
//   - f: An endomorphism (pure function from T to T) that transforms the field value
//
// Example:
//
//	type Accumulator struct {
//	    counter int
//	    message string
//	}
//
//	counterLens := MakeLens(
//	    func(acc Accumulator) int { return acc.counter },
//	    func(c int) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.counter = c
//	            return acc
//	        }
//	    },
//	)
//
//	// Increment the counter
//	computation := LetL(
//	    counterLens,
//	    func(c int) int { return c + 1 },
//	)
//
//go:inline
func LetL[ST, S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Endomorphism[State[ST, S]] {
	return Let[ST](lens.Set, function.Flow2(lens.Get, f))
}

// LetToL is a lens-based version of LetTo that sets a specific field to a
// constant value. This provides a convenient way to update nested fields
// with known values.
//
// Parameters:
//   - lens: A lens focusing on field T within structure S
//   - b: The constant value to set
//
// Example:
//
//	type Accumulator struct {
//	    status string
//	    value int
//	}
//
//	statusLens := MakeLens(
//	    func(acc Accumulator) string { return acc.status },
//	    func(s string) func(Accumulator) Accumulator {
//	        return func(acc Accumulator) Accumulator {
//	            acc.status = s
//	            return acc
//	        }
//	    },
//	)
//
//	// Set status to "completed"
//	computation := LetToL(statusLens, "completed")
//
//go:inline
func LetToL[ST, S, T any](
	lens Lens[S, T],
	b T,
) Endomorphism[State[ST, S]] {
	return LetTo[ST](lens.Set, b)
}
