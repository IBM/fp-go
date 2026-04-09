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

/*
Package constraints defines a set of useful type constraints for generic programming in Go.

# Overview

This package provides type constraints that can be used with Go generics to restrict
type parameters to specific categories of types. These constraints are similar to those
in Go's standard constraints package but are defined here for consistency within the
fp-go project.

# Type Constraints

Ordered - Types that support comparison operators:

	type Ordered interface {
	    Integer | Float | ~string
	}

Used for types that can be compared using <, <=, >, >= operators.

Integer - All integer types (signed and unsigned):

	type Integer interface {
	    Signed | Unsigned
	}

Signed - Signed integer types:

	type Signed interface {
	    ~int | ~int8 | ~int16 | ~int32 | ~int64
	}

Unsigned - Unsigned integer types:

	type Unsigned interface {
	    ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}

Float - Floating-point types:

	type Float interface {
	    ~float32 | ~float64
	}

Complex - Complex number types:

	type Complex interface {
	    ~complex64 | ~complex128
	}

# Usage Examples

Using Ordered constraint for comparison:

	import C "github.com/IBM/fp-go/v2/constraints"

	func Min[T C.Ordered](a, b T) T {
	    if a < b {
	        return a
	    }
	    return b
	}

	result := Min(5, 3)        // 3
	result := Min(3.14, 2.71)  // 2.71
	result := Min("apple", "banana")  // "apple"

Using Integer constraint:

	func Abs[T C.Integer](n T) T {
	    if n < 0 {
	        return -n
	    }
	    return n
	}

	result := Abs(-42)  // 42
	result := Abs(uint(10))  // 10

Using Float constraint:

	func Average[T C.Float](a, b T) T {
	    return (a + b) / 2
	}

	result := Average(3.14, 2.86)  // 3.0

Using Complex constraint:

	func Magnitude[T C.Complex](c T) float64 {
	    r, i := real(c), imag(c)
	    return math.Sqrt(r*r + i*i)
	}

	c := complex(3, 4)
	result := Magnitude(c)  // 5.0

# Combining Constraints

Constraints can be combined to create more specific type restrictions:

	type Number interface {
	    C.Integer | C.Float | C.Complex
	}

	func Add[T Number](a, b T) T {
	    return a + b
	}

# Tilde Operator

The ~ operator in type constraints means "underlying type". For example, ~int
matches not only int but also any type whose underlying type is int:

	type MyInt int

	func Double[T C.Integer](n T) T {
	    return n * 2
	}

	var x MyInt = 5
	result := Double(x)  // Works because MyInt's underlying type is int

# Related Packages

  - number: Provides algebraic structures and utilities for numeric types
  - ord: Provides ordering operations using these constraints
  - eq: Provides equality operations for comparable types
*/
package constraints

// Ordered is a constraint that permits any ordered type: any type that supports
// the operators < <= >= >. Ordered types include integers, floats, and strings.
//
// This constraint is commonly used for comparison operations, sorting, and
// finding minimum/maximum values.
//
// Example:
//
//	func Max[T Ordered](a, b T) T {
//	    if a > b {
//	        return a
//	    }
//	    return b
//	}
type Ordered interface {
	Integer | Float | ~string
}

// Signed is a constraint that permits any signed integer type.
// This includes int, int8, int16, int32, and int64, as well as any
// types whose underlying type is one of these.
//
// Example:
//
//	func Negate[T Signed](n T) T {
//	    return -n
//	}
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned is a constraint that permits any unsigned integer type.
// This includes uint, uint8, uint16, uint32, uint64, and uintptr, as well
// as any types whose underlying type is one of these.
//
// Example:
//
//	func IsEven[T Unsigned](n T) bool {
//	    return n%2 == 0
//	}
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer is a constraint that permits any integer type, both signed and unsigned.
// This is a union of the Signed and Unsigned constraints.
//
// Example:
//
//	func Abs[T Integer](n T) T {
//	    if n < 0 {
//	        return -n
//	    }
//	    return n
//	}
type Integer interface {
	Signed | Unsigned
}

// Float is a constraint that permits any floating-point type.
// This includes float32 and float64, as well as any types whose
// underlying type is one of these.
//
// Example:
//
//	func Round[T Float](f T) T {
//	    return T(math.Round(float64(f)))
//	}
type Float interface {
	~float32 | ~float64
}

// Complex is a constraint that permits any complex numeric type.
// This includes complex64 and complex128, as well as any types whose
// underlying type is one of these.
//
// Example:
//
//	func Conjugate[T Complex](c T) T {
//	    return complex(real(c), -imag(c))
//	}
type Complex interface {
	~complex64 | ~complex128
}
