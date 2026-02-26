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

// Package constant provides the Const functor, a phantom type that ignores its second type parameter.
//
// The Const functor is a fundamental building block in functional programming that wraps a value
// of type E while having a phantom type parameter A. This makes it useful for:
//   - Accumulating values during traversals (e.g., collecting metadata)
//   - Implementing optics (lenses, prisms) where you need to track information
//   - Building applicative functors that combine values using a semigroup
//
// # The Const Functor
//
// Const[E, A] wraps a value of type E and has a phantom type parameter A that doesn't affect
// the runtime value. This allows it to participate in functor and applicative operations while
// maintaining the wrapped value unchanged.
//
// # Key Properties
//
//   - Map operations ignore the function and preserve the wrapped value
//   - Ap operations combine wrapped values using a semigroup
//   - The phantom type A allows type-safe composition with other functors
//
// # Example Usage
//
//	// Accumulate string values
//	c1 := Make[string, int]("hello")
//	c2 := Make[string, int]("world")
//
//	// Map doesn't change the wrapped value
//	mapped := Map[string, int, string](strconv.Itoa)(c1)  // Still contains "hello"
//
//	// Ap combines values using a semigroup
//	combined := Ap[string, int, int](S.Monoid)(c1)(c2)  // Contains "helloworld"
package constant

import (
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Const is a functor that wraps a value of type E with a phantom type parameter A.
//
// The Const functor is useful for accumulating values during traversals or implementing
// optics. The type parameter A is phantom - it doesn't affect the runtime value but allows
// the type to participate in functor and applicative operations.
//
// Type Parameters:
//   - E: The type of the wrapped value (the actual data)
//   - A: The phantom type parameter (not stored, only used for type-level operations)
//
// Example:
//
//	// Create a Const that wraps a string
//	c := Make[string, int]("metadata")
//
//	// The int type parameter is phantom - no int value is stored
//	value := Unwrap(c)  // "metadata"
type Const[E, A any] struct {
	value E
}

// Make creates a Const value wrapping the given value.
//
// This is the primary constructor for Const values. The second type parameter A
// is phantom and must be specified explicitly when needed for type inference.
//
// Type Parameters:
//   - E: The type of the value to wrap
//   - A: The phantom type parameter
//
// Parameters:
//   - e: The value to wrap
//
// Returns:
//   - A Const[E, A] wrapping the value
//
// Example:
//
//	c := Make[string, int]("hello")
//	value := Unwrap(c)  // "hello"
func Make[E, A any](e E) Const[E, A] {
	return Const[E, A]{value: e}
}

// Unwrap extracts the wrapped value from a Const.
//
// This is the inverse of Make, retrieving the actual value stored in the Const.
//
// Type Parameters:
//   - E: The type of the wrapped value
//   - A: The phantom type parameter
//
// Parameters:
//   - c: The Const to unwrap
//
// Returns:
//   - The wrapped value of type E
//
// Example:
//
//	c := Make[string, int]("world")
//	value := Unwrap(c)  // "world"
func Unwrap[E, A any](c Const[E, A]) E {
	return c.value
}

// Of creates a Const containing the monoid's empty value, ignoring the input.
//
// This implements the Applicative's "pure" operation for Const. It creates a Const
// wrapping the monoid's identity element, regardless of the input value.
//
// Type Parameters:
//   - E: The type of the wrapped value (must have a monoid)
//   - A: The input type (ignored)
//
// Parameters:
//   - m: The monoid providing the empty value
//
// Returns:
//   - A function that ignores its input and returns Const[E, A] with the empty value
//
// Example:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	of := Of[string, int](S.Monoid)
//	c := of(42)  // Const[string, int] containing ""
//	value := Unwrap(c)  // ""
func Of[E, A any](m M.Monoid[E]) func(A) Const[E, A] {
	return F.Constant1[A](Make[E, A](m.Empty()))
}

// MonadMap applies a function to the phantom type parameter without changing the wrapped value.
//
// This implements the Functor's map operation for Const. Since the type parameter A is phantom,
// the function is never actually called - the wrapped value E remains unchanged.
//
// Type Parameters:
//   - E: The type of the wrapped value
//   - A: The input phantom type
//   - B: The output phantom type
//
// Parameters:
//   - fa: The Const to map over
//   - _: The function to apply (ignored)
//
// Returns:
//   - A Const[E, B] with the same wrapped value
//
// Example:
//
//	c := Make[string, int]("hello")
//	mapped := MonadMap(c, func(i int) string { return strconv.Itoa(i) })
//	// mapped still contains "hello", function was never called
func MonadMap[E, A, B any](fa Const[E, A], _ func(A) B) Const[E, B] {
	return Make[E, B](fa.value)
}

// MonadAp combines two Const values using a semigroup.
//
// This implements the Applicative's ap operation for Const. It combines the wrapped
// values from both Const instances using the provided semigroup, ignoring the function
// type in the first argument.
//
// Type Parameters:
//   - E: The type of the wrapped values (must have a semigroup)
//   - A: The input phantom type
//   - B: The output phantom type
//
// Parameters:
//   - s: The semigroup for combining wrapped values
//
// Returns:
//   - A function that takes two Const values and combines their wrapped values
//
// Example:
//
//	import S "github.com/IBM/fp-go/v2/string"
//
//	ap := MonadAp[string, int, int](S.Monoid)
//	c1 := Make[string, func(int) int]("hello")
//	c2 := Make[string, int]("world")
//	result := ap(c1, c2)  // Const containing "helloworld"
func MonadAp[E, A, B any](s S.Semigroup[E]) func(fab Const[E, func(A) B], fa Const[E, A]) Const[E, B] {
	return func(fab Const[E, func(A) B], fa Const[E, A]) Const[E, B] {
		return Make[E, B](s.Concat(fab.value, fa.value))
	}
}

// Map applies a function to the phantom type parameter without changing the wrapped value.
//
// This is the curried version of MonadMap, providing a more functional programming style.
// The function is never actually called since A is a phantom type.
//
// Type Parameters:
//   - E: The type of the wrapped value
//   - A: The input phantom type
//   - B: The output phantom type
//
// Parameters:
//   - f: The function to apply (ignored)
//
// Returns:
//   - A function that transforms Const[E, A] to Const[E, B]
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	c := Make[string, int]("data")
//	mapped := F.Pipe1(c, Map[string, int, string](strconv.Itoa))
//	// mapped still contains "data"
func Map[E, A, B any](f func(A) B) func(fa Const[E, A]) Const[E, B] {
	return F.Bind2nd(MonadMap[E, A, B], f)
}

// Ap combines Const values using a semigroup in a curried style.
//
// This is the curried version of MonadAp, providing data-last style for better composition.
// It combines the wrapped values from both Const instances using the provided semigroup.
//
// Type Parameters:
//   - E: The type of the wrapped values (must have a semigroup)
//   - A: The input phantom type
//   - B: The output phantom type
//
// Parameters:
//   - s: The semigroup for combining wrapped values
//
// Returns:
//   - A curried function for combining Const values
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    S "github.com/IBM/fp-go/v2/string"
//	)
//
//	c1 := Make[string, int]("hello")
//	c2 := Make[string, func(int) int]("world")
//	result := F.Pipe1(c1, Ap[string, int, int](S.Monoid)(c2))
//	// result contains "helloworld"
func Ap[E, A, B any](s S.Semigroup[E]) func(fa Const[E, A]) func(fab Const[E, func(A) B]) Const[E, B] {
	monadap := MonadAp[E, A, B](s)
	return func(fa Const[E, A]) func(fab Const[E, func(A) B]) Const[E, B] {
		return F.Bind2nd(monadap, fa)
	}
}
