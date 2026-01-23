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

package function

// Identity returns its argument unchanged.
//
// This is the identity function from category theory, which satisfies:
//   - Identity(x) = x for all x
//
// It's useful as a default transformation or when you need a function that
// does nothing but is required by an API.
//
// Example:
//
//	result := Identity(42)        // 42
//	result := Identity("hello")   // "hello"
//
//	// Useful in higher-order functions
//	values := []int{1, 2, 3}
//	mapped := Map(Identity[int])(values)  // [1, 2, 3]
func Identity[A any](a A) A {
	return a
}

// Constant creates a nullary function that always returns the same value.
//
// This creates a function with no parameters that returns the constant value 'a'.
// Useful for lazy evaluation or when you need a function that produces a fixed value.
//
// Parameters:
//   - a: The constant value to return
//
// Returns:
//   - A function that takes no arguments and returns 'a'
//
// Example:
//
//	getFortyTwo := Constant(42)
//	result := getFortyTwo()  // 42
//
//	getMessage := Constant("Hello")
//	msg := getMessage()  // "Hello"
//
//go:inline
func Constant[A any](a A) func() A {
	//go:inline
	return func() A {
		return a
	}
}

// Constant1 creates a unary function that always returns the same value, ignoring its input.
//
// This creates a function that takes one parameter but ignores it and always returns
// the constant value 'a'. Useful for providing default values or placeholder functions.
//
// Type Parameters:
//   - B: The type of the ignored input parameter
//   - A: The type of the constant return value
//
// Parameters:
//   - a: The constant value to return
//
// Returns:
//   - A function that takes a B and returns 'a'
//
// Example:
//
//	alwaysZero := Constant1[string, int](0)
//	result := alwaysZero("anything")  // 0
//
//	defaultName := Constant1[int, string]("Unknown")
//	name := defaultName(42)  // "Unknown"
//
//go:inline
func Constant1[B, A any](a A) func(B) A {
	//go:inline
	return func(_ B) A {
		return a
	}
}

// Constant2 creates a binary function that always returns the same value, ignoring its inputs.
//
// This creates a function that takes two parameters but ignores both and always returns
// the constant value 'a'.
//
// Type Parameters:
//   - B: The type of the first ignored input parameter
//   - C: The type of the second ignored input parameter
//   - A: The type of the constant return value
//
// Parameters:
//   - a: The constant value to return
//
// Returns:
//   - A function that takes a B and C and returns 'a'
//
// Example:
//
//	alwaysTrue := Constant2[int, string, bool](true)
//	result := alwaysTrue(42, "test")  // true
//
//go:inline
func Constant2[B, C, A any](a A) func(B, C) A {
	//go:inline
	return func(_ B, _ C) A {
		return a
	}
}

// IsNil checks if a pointer is nil.
//
// Parameters:
//   - a: A pointer to check
//
// Returns:
//   - true if the pointer is nil, false otherwise
//
// Example:
//
//	var ptr *int
//	IsNil(ptr)  // true
//
//	value := 42
//	IsNil(&value)  // false
//
//go:inline
func IsNil[A any](a *A) bool {
	return a == nil
}

// IsNonNil checks if a pointer is not nil.
//
// This is the logical negation of IsNil.
//
// Parameters:
//   - a: A pointer to check
//
// Returns:
//   - true if the pointer is not nil, false otherwise
//
// Example:
//
//	var ptr *int
//	IsNonNil(ptr)  // false
//
//	value := 42
//	IsNonNil(&value)  // true
//
//go:inline
func IsNonNil[A any](a *A) bool {
	return a != nil
}

// Swap returns a new binary function with the parameter order reversed.
//
// Given a function f(a, b), Swap returns a function g(b, a) where g(b, a) = f(a, b).
// This is useful when you have a function but need to call it with arguments in
// a different order.
//
// Type Parameters:
//   - T1: The type of the first parameter (becomes second)
//   - T2: The type of the second parameter (becomes first)
//   - R: The return type
//
// Parameters:
//   - f: The function to swap
//
// Returns:
//   - A new function with swapped parameters
//
// Example:
//
//	divide := func(a, b float64) float64 { return a / b }
//	divideSwapped := Swap(divide)
//
//	result1 := divide(10, 2)         // 5.0 (10 / 2)
//	result2 := divideSwapped(10, 2)  // 0.2 (2 / 10)
//
//	subtract := func(a, b int) int { return a - b }
//	subtractSwapped := Swap(subtract)
//	result := subtractSwapped(5, 10)  // 5 (10 - 5)
func Swap[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(t2 T2, t1 T1) R {
		return f(t1, t2)
	}
}

// First returns the first of two input values, ignoring the second.
//
// This is a projection function that selects the first element of a pair.
// Also known as the K combinator in combinatory logic.
//
// Type Parameters:
//   - T1: The type of the first value (returned)
//   - T2: The type of the second value (ignored)
//
// Parameters:
//   - t1: The first value
//   - t2: The second value (ignored)
//
// Returns:
//   - The first value
//
// Example:
//
//	result := First(42, "hello")  // 42
//	result := First(true, 100)    // true
//
//go:inline
func First[T1, T2 any](t1 T1, _ T2) T1 {
	return t1
}

// Second returns the second of two input values, ignoring the first.
//
// This is a projection function that selects the second element of a pair.
// Identical to SK combinator in combinatory logic.
//
// Type Parameters:
//   - T1: The type of the first value (ignored)
//   - T2: The type of the second value (returned)
//
// Parameters:
//   - t1: The first value (ignored)
//   - t2: The second value
//
// Returns:
//   - The second value
//
// Example:
//
//	result := Second(42, "hello")  // "hello"
//	result := Second(true, 100)    // 100
//
//go:inline
func Second[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}

// Zero returns the zero value of the given type.
func Zero[A any]() A {
	var zero A
	return zero
}
