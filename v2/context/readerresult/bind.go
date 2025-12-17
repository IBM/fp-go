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

package readerresult

import (
	F "github.com/IBM/fp-go/v2/function"
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//	result := readereither.Do(State{})
//
//go:inline
func Do[S any](
	empty S,
) ReaderResult[S] {
	return G.Do[ReaderResult[S]](empty)
}

// Bind attaches the result of an EFFECTFUL computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access the context.Context from the environment.
//
// IMPORTANT: Bind is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The function parameter takes state and returns a ReaderResult[T], which is effectful because
// it depends on context.Context (can be cancelled, has deadlines, carries values).
//
// For PURE FUNCTIONS (side-effect free), use:
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//   - Let: For pure functions without errors (State -> Value)
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//
//	result := F.Pipe2(
//	    readereither.Do(State{}),
//	    readereither.Bind(
//	        func(uid string) func(State) State {
//	            return func(s State) State { s.UserID = uid; return s }
//	        },
//	        func(s State) readereither.ReaderResult[string] {
//	            return func(ctx context.Context) either.Either[error, string] {
//	                if uid, ok := ctx.Value("userID").(string); ok {
//	                    return either.Right[error](uid)
//	                }
//	                return either.Left[string](errors.New("no userID"))
//	            }
//	        },
//	    ),
//	    readereither.Bind(
//	        func(tid string) func(State) State {
//	            return func(s State) State { s.TenantID = tid; return s }
//	        },
//	        func(s State) readereither.ReaderResult[string] {
//	            // This can access s.UserID from the previous step
//	            return func(ctx context.Context) either.Either[error, string] {
//	                return either.Right[error]("tenant-" + s.UserID)
//	            }
//	        },
//	    ),
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Kleisli[ReaderResult[S1], S2] {
	return G.Bind[ReaderResult[S1], ReaderResult[S2]](setter, F.Flow2(f, WithContext))
}

// Let attaches the result of a PURE computation to a context [S1] to produce a context [S2].
//
// IMPORTANT: Let is for PURE FUNCTIONS (side-effect free) that don't depend on context.Context.
// The function parameter takes state and returns a value directly, with no errors or effects.
//
// For EFFECTFUL FUNCTIONS (that need context.Context), use:
//   - Bind: For effectful ReaderResult computations (State -> ReaderResult[Value])
//
// For PURE FUNCTIONS with error handling, use:
//   - BindResultK: For pure functions with errors (State -> (Value, error))
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Kleisli[ReaderResult[S1], S2] {
	return G.Let[ReaderResult[S1], ReaderResult[S2]](setter, f)
}

// LetTo attaches a constant value to a context [S1] to produce a context [S2].
// This is a PURE operation (side-effect free) that simply sets a field to a constant value.
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Kleisli[ReaderResult[S1], S2] {
	return G.LetTo[ReaderResult[S1], ReaderResult[S2]](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return G.BindTo[ReaderResult[S1], ReaderResult[T]](setter)
}

//go:inline
func BindToP[S1, T any](
	setter Prism[S1, T],
) Operator[T, S1] {
	return BindTo(setter.ReverseGet)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent EFFECTFUL computations to be combined without one depending on the result of the other.
//
// IMPORTANT: ApS is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The ReaderResult parameter is effectful because it depends on context.Context.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    UserID   string
//	    TenantID string
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUserID := func(ctx context.Context) either.Either[error, string] {
//	    return either.Right[error](ctx.Value("userID").(string))
//	}
//	getTenantID := func(ctx context.Context) either.Either[error, string] {
//	    return either.Right[error](ctx.Value("tenantID").(string))
//	}
//
//	result := F.Pipe2(
//	    readereither.Do(State{}),
//	    readereither.ApS(
//	        func(uid string) func(State) State {
//	            return func(s State) State { s.UserID = uid; return s }
//	        },
//	        getUserID,
//	    ),
//	    readereither.ApS(
//	        func(tid string) func(State) State {
//	            return func(s State) State { s.TenantID = tid; return s }
//	        },
//	        getTenantID,
//	    ),
//	)
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderResult[T],
) Kleisli[ReaderResult[S1], S2] {
	return G.ApS[ReaderResult[S1], ReaderResult[S2]](setter, fa)
}

// ApSL is a variant of ApS that uses a lens to focus on a specific field in the state.
// Instead of providing a setter function, you provide a lens that knows how to get and set
// the field. This is more convenient when working with nested structures.
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - fa: A ReaderResult computation that produces a value of type T
//
// Returns:
//   - A function that transforms ReaderResult[S] to ReaderResult[S] by setting the focused field
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
//	getAge := func(ctx context.Context) either.Either[error, int] {
//	    return either.Right[error](30)
//	}
//
//	result := F.Pipe1(
//	    readereither.Do(Person{Name: "Alice", Age: 25}),
//	    readereither.ApSL(ageLens, getAge),
//	)
//
//go:inline
func ApSL[S, T any](
	lens Lens[S, T],
	fa ReaderResult[T],
) Kleisli[ReaderResult[S], S] {
	return ApS(lens.Set, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific field in the state.
// It combines the lens-based field access with monadic composition for EFFECTFUL computations.
//
// IMPORTANT: BindL is for EFFECTFUL FUNCTIONS that depend on context.Context.
// The function parameter returns a ReaderResult, which is effectful.
//
// It allows you to:
// 1. Extract a field value using the lens
// 2. Use that value in an effectful computation that may fail
// 3. Update the field with the result
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - f: A function that takes the current field value and returns a ReaderResult computation
//
// Returns:
//   - A function that transforms ReaderResult[S] to ReaderResult[S]
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
//	increment := func(v int) readereither.ReaderResult[int] {
//	    return func(ctx context.Context) either.Either[error, int] {
//	        if v >= 100 {
//	            return either.Left[int](errors.New("value too large"))
//	        }
//	        return either.Right[error](v + 1)
//	    }
//	}
//
//	result := F.Pipe1(
//	    readereither.Of[error](Counter{Value: 42}),
//	    readereither.BindL(valueLens, increment),
//	)
//
//go:inline
func BindL[S, T any](
	lens Lens[S, T],
	f Kleisli[T, T],
) Kleisli[ReaderResult[S], S] {
	return Bind(lens.Set, F.Flow2(lens.Get, F.Flow2(f, WithContext)))
}

// LetL is a variant of Let that uses a lens to focus on a specific field in the state.
// It applies a PURE transformation to the focused field without any effects.
//
// IMPORTANT: LetL is for PURE FUNCTIONS (side-effect free) that don't depend on context.Context.
// The function parameter is a pure endomorphism (T -> T) with no errors or effects.
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - f: A pure function that transforms the field value
//
// Returns:
//   - A function that transforms ReaderResult[S] to ReaderResult[S]
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
//	double := func(v int) int { return v * 2 }
//
//	result := F.Pipe1(
//	    readereither.Of[error](Counter{Value: 21}),
//	    readereither.LetL(valueLens, double),
//	)
//	// result when executed will be Right(Counter{Value: 42})
//
//go:inline
func LetL[S, T any](
	lens Lens[S, T],
	f Endomorphism[T],
) Kleisli[ReaderResult[S], S] {
	return Let(lens.Set, F.Flow2(lens.Get, f))
}

// LetToL is a variant of LetTo that uses a lens to focus on a specific field in the state.
// It sets the focused field to a constant value. This is a PURE operation (side-effect free).
//
// Parameters:
//   - lens: A lens that focuses on a field of type T within state S
//   - b: The constant value to set
//
// Returns:
//   - A function that transforms ReaderResult[S] to ReaderResult[S]
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
//	    readereither.Of[error](Config{Debug: true, Timeout: 30}),
//	    readereither.LetToL(debugLens, false),
//	)
//	// result when executed will be Right(Config{Debug: false, Timeout: 30})
//
//go:inline
func LetToL[S, T any](
	lens Lens[S, T],
	b T,
) Kleisli[ReaderResult[S], S] {
	return LetTo(lens.Set, b)
}
