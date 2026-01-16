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

package io

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
)

type (
	ioApplicative[A, B any] struct{}

	// IOApplicative represents the applicative functor type class for IO.
	// It combines the capabilities of Functor (Map) and Pointed (Of) with
	// the ability to apply wrapped functions to wrapped values (Ap).
	//
	// An applicative functor is a functor with two additional operations:
	//   - Of: lifts a pure value into the IO context
	//   - Ap: applies a wrapped function to a wrapped value
	//
	// This allows for function application within the IO context while maintaining
	// the computational structure. The Ap operation uses parallel execution by default
	// for better performance.
	//
	// Type parameters:
	//   - A: the input type
	//   - B: the output type
	IOApplicative[A, B any] = applicative.Applicative[A, B, IO[A], IO[B], IO[func(A) B]]
)

// Of lifts a pure value into the IO context.
// This is the pointed functor operation that wraps a value in an IO computation.
//
// Example:
//
//	app := io.Applicative[int, string]()
//	ioValue := app.Of(42) // IO[int] that returns 42
//	result := ioValue()   // 42
func (o *ioApplicative[A, B]) Of(a A) IO[A] {
	return Of(a)
}

// Map transforms the result of an IO computation by applying a function to it.
// This is the functor operation that allows mapping over wrapped values.
//
// Example:
//
//	app := io.Applicative[int, string]()
//	double := func(x int) int { return x * 2 }
//	ioValue := app.Of(21)
//	doubled := app.Map(double)(ioValue)
//	result := doubled() // 42
func (o *ioApplicative[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

// Ap applies a wrapped function to a wrapped value, both in the IO context.
// This operation uses parallel execution by default, running the function and
// value computations concurrently for better performance.
//
// The Ap operation is useful for applying multi-argument functions in a curried
// fashion within the IO context.
//
// Example:
//
//	app := io.Applicative[int, int]()
//	add := func(a int) func(int) int {
//		return func(b int) int { return a + b }
//	}
//	ioFunc := app.Of(add(10))  // IO[func(int) int]
//	ioValue := app.Of(32)      // IO[int]
//	result := app.Ap(ioValue)(ioFunc)
//	value := result() // 42
func (o *ioApplicative[A, B]) Ap(fa IO[A]) Operator[func(A) B, B] {
	return Ap[B](fa)
}

// Applicative returns an instance of the Applicative type class for IO.
// This provides a structured way to access applicative operations (Of, Map, Ap)
// for IO computations.
//
// The applicative pattern is useful when you need to:
//   - Apply functions with multiple arguments to wrapped values
//   - Combine multiple independent IO computations
//   - Maintain the computational structure while transforming values
//
// Type parameters:
//   - A: the input type for the applicative operations
//   - B: the output type for the applicative operations
//
// Example - Basic usage:
//
//	app := io.Applicative[int, string]()
//	result := app.Map(strconv.Itoa)(app.Of(42))
//	value := result() // "42"
//
// Example - Applying curried functions:
//
//	app := io.Applicative[int, int]()
//	add := func(a int) func(int) int {
//		return func(b int) int { return a + b }
//	}
//	// Create IO computations
//	ioFunc := io.Map(add)(app.Of(10))  // IO[func(int) int]
//	ioValue := app.Of(32)               // IO[int]
//	// Apply the function to the value
//	result := app.Ap(ioValue)(ioFunc)
//	value := result() // 42
//
// Example - Combining multiple IO computations:
//
//	app := io.Applicative[int, int]()
//	multiply := func(a int) func(int) int {
//		return func(b int) int { return a * b }
//	}
//	io1 := app.Of(6)
//	io2 := app.Of(7)
//	ioFunc := io.Map(multiply)(io1)
//	result := app.Ap(io2)(ioFunc)
//	value := result() // 42
func Applicative[A, B any]() IOApplicative[A, B] {
	return &ioApplicative[A, B]{}
}
