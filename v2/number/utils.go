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

// Package number provides utility functions for numeric operations with a functional programming approach.
// It includes curried arithmetic operations, comparison functions, and min/max utilities that work with
// generic numeric types.
package number

import (
	C "github.com/IBM/fp-go/v2/constraints"
)

// Number is a constraint that represents all numeric types including integers, floats, and complex numbers.
// It combines the Integer, Float, and Complex constraints from the constraints package.
type Number interface {
	C.Integer | C.Float | C.Complex
}

// Add is a curried function that adds two numbers.
// It takes a right operand and returns a function that takes a left operand,
// returning their sum (left + right).
//
// This curried form is useful for partial application and function composition.
//
// Example:
//
//	add5 := Add(5)
//	result := add5(10) // returns 15
func Add[T Number](right T) func(T) T {
	return func(left T) T {
		return left + right
	}
}

// Sub is a curried function that subtracts two numbers.
// It takes a right operand and returns a function that takes a left operand,
// returning their difference (left - right).
//
// This curried form is useful for partial application and function composition.
//
// Example:
//
//	sub5 := Sub(5)
//	result := sub5(10) // returns 5
func Sub[T Number](right T) func(T) T {
	return func(left T) T {
		return left - right
	}
}

// Mul is a curried function that multiplies two numbers.
// It takes a right operand and returns a function that takes a left operand,
// returning their product (left * right).
//
// This curried form is useful for partial application and function composition.
//
// Example:
//
//	double := Mul(2)
//	result := double(10) // returns 20
func Mul[T Number](right T) func(T) T {
	return func(left T) T {
		return left * right
	}
}

// Div is a curried function that divides two numbers.
// It takes a right operand (divisor) and returns a function that takes a left operand (dividend),
// returning their quotient (left / right).
//
// Note: Division by zero will cause a panic for integer types or return infinity for float types.
//
// This curried form is useful for partial application and function composition.
//
// Example:
//
//	halve := Div(2)
//	result := halve(10) // returns 5
func Div[T Number](right T) func(T) T {
	return func(left T) T {
		return left / right
	}
}

// Inc increments a number by 1.
// It works with any numeric type that satisfies the Number constraint.
//
// Example:
//
//	result := Inc(5) // returns 6
func Inc[T Number](value T) T {
	return value + 1
}

// Min returns the minimum of two ordered values.
// If the values are considered equal, the first argument is returned.
//
// This function works with any type that satisfies the Ordered constraint,
// including all numeric types and strings.
//
// Example:
//
//	result := Min(5, 10) // returns 5
//	result := Min(3.14, 2.71) // returns 2.71
func Min[A C.Ordered](a, b A) A {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two ordered values.
// If the values are considered equal, the first argument is returned.
//
// This function works with any type that satisfies the Ordered constraint,
// including all numeric types and strings.
//
// Example:
//
//	result := Max(5, 10) // returns 10
//	result := Max(3.14, 2.71) // returns 3.14
func Max[A C.Ordered](a, b A) A {
	if a > b {
		return a
	}
	return b
}

// MoreThan is a curried comparison function that checks if a value is more than (greater than) another.
// It takes a threshold value 'a' and returns a predicate function that checks if 'a' is less than its argument,
// meaning the argument is more than 'a'.
//
// This curried form is useful for creating reusable predicates and function composition.
//
// Example:
//
//	moreThan10 := MoreThan(10)
//	result := moreThan10(15) // returns true (15 is more than 10)
//	result := moreThan10(10) // returns false (10 is not more than 10)
//	result := moreThan10(5)  // returns false (5 is not more than 10)
func MoreThan[A C.Ordered](a A) func(A) bool {
	return func(b A) bool {
		return a < b
	}
}

// LessThan is a curried comparison function that checks if a value is less than another.
// It takes a threshold value 'a' and returns a predicate function that checks if 'a' is greater than its argument,
// meaning the argument is less than 'a'.
//
// This curried form is useful for creating reusable predicates and function composition.
//
// Example:
//
//	lessThan10 := LessThan(10)
//	result := lessThan10(5)  // returns true (5 is less than 10)
//	result := lessThan10(10) // returns false (10 is not less than 10)
//	result := lessThan10(15) // returns false (15 is not less than 10)
func LessThan[A C.Ordered](a A) func(A) bool {
	return func(b A) bool {
		return a > b
	}
}
