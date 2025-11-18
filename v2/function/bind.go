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

// Bind1st performs partial application by fixing the first argument of a binary function.
//
// Given a function f(a, b) and a value for 'a', this returns a new function g(b)
// where g(b) = f(a, b). This is useful for creating specialized functions from
// more general ones.
//
// Type Parameters:
//   - T1: The type of the first parameter (to be bound)
//   - T2: The type of the second parameter (remains free)
//   - R: The return type
//
// Parameters:
//   - f: The binary function to partially apply
//   - t1: The value to bind to the first parameter
//
// Returns:
//   - A unary function with the first parameter fixed
//
// Example:
//
//	multiply := func(a, b int) int { return a * b }
//	double := Bind1st(multiply, 2)
//	result := double(5)  // 10 (2 * 5)
//
//	divide := func(a, b float64) float64 { return a / b }
//	divideBy10 := Bind1st(divide, 10.0)
//	result := divideBy10(2.0)  // 5.0 (10 / 2)
//
//go:inline
func Bind1st[T1, T2, R any](f func(T1, T2) R, t1 T1) func(T2) R {
	//go:inline
	return func(t2 T2) R {
		return f(t1, t2)
	}
}

// Bind2nd performs partial application by fixing the second argument of a binary function.
//
// Given a function f(a, b) and a value for 'b', this returns a new function g(a)
// where g(a) = f(a, b). This is useful for creating specialized functions from
// more general ones.
//
// Type Parameters:
//   - T1: The type of the first parameter (remains free)
//   - T2: The type of the second parameter (to be bound)
//   - R: The return type
//
// Parameters:
//   - f: The binary function to partially apply
//   - t2: The value to bind to the second parameter
//
// Returns:
//   - A unary function with the second parameter fixed
//
// Example:
//
//	multiply := func(a, b int) int { return a * b }
//	triple := Bind2nd(multiply, 3)
//	result := triple(5)  // 15 (5 * 3)
//
//	divide := func(a, b float64) float64 { return a / b }
//	halve := Bind2nd(divide, 2.0)
//	result := halve(10.0)  // 5.0 (10 / 2)
//
//go:inline
func Bind2nd[T1, T2, R any](f func(T1, T2) R, t2 T2) func(T1) R {
	//go:inline
	return func(t1 T1) R {
		return f(t1, t2)
	}
}

// SK is the SK combinator from SKI combinator calculus.
//
// This function takes two arguments and returns the second, ignoring the first.
// It's identical to the Second function and represents the K combinator applied
// to the second argument.
//
// In combinatory logic: SK = λx.λy.y
//
// Type Parameters:
//   - T1: The type of the first parameter (ignored)
//   - T2: The type of the second parameter (returned)
//
// Parameters:
//   - _: The first value (ignored)
//   - t2: The second value
//
// Returns:
//   - The second value
//
// Example:
//
//	result := SK(42, "hello")  // "hello"
//	result := SK(true, 100)    // 100
//
//go:inline
func SK[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}
