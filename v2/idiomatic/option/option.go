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

// Package option implements the Option monad using idiomatic Go data types.
//
// Unlike the standard option package which uses wrapper structs, this package represents
// Options as tuples (value, bool) where the boolean indicates presence (true) or absence (false).
// This approach is more idiomatic in Go and has better performance characteristics.
//
// Example:
//
//	// Creating Options
//	some := Some(42)           // (42, true)
//	none := None[int]()        // (0, false)
//
//	// Using Options
//	result, ok := some         // ok == true, result == 42
//	result, ok := none         // ok == false, result == 0
//
//	// Transforming Options
//	doubled := Map(N.Mul(2))(some)  // (84, true)
package option

import (
	"github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/predicate"
)

// FromPredicate returns a function that creates an Option based on a predicate.
// The returned function will wrap a value in Some if the predicate is satisfied, otherwise None.
//
// Parameters:
//   - pred: A predicate function that determines if a value should be wrapped in Some
//
// Example:
//
//	isPositive := FromPredicate(N.MoreThan(0))
//	result := isPositive(5)  // Some(5)
//	result := isPositive(-1) // None
func FromPredicate[A any](pred func(A) bool) Kleisli[A, A] {
	return func(a A) (A, bool) {
		return a, pred(a)
	}
}

// FromZero returns a function that creates an Option based on whether a value is the zero value.
// Returns Some if the value is the zero value, None otherwise.
//
// Example:
//
//	checkZero := FromZero[int]()
//	result := checkZero(0)  // Some(0)
//	result := checkZero(5)  // None
//
//go:inline
func FromZero[A comparable]() Kleisli[A, A] {
	return FromPredicate(P.IsZero[A]())
}

// FromNonZero returns a function that creates an Option based on whether a value is non-zero.
// Returns Some if the value is non-zero, None otherwise.
//
// Example:
//
//	checkNonZero := FromNonZero[int]()
//	result := checkNonZero(5)  // Some(5)
//	result := checkNonZero(0)  // None
//
//go:inline
func FromNonZero[A comparable]() Kleisli[A, A] {
	return FromPredicate(P.IsNonZero[A]())
}

// FromEq returns a function that creates an Option based on equality with a given value.
// The returned function takes a value to compare against and returns a Kleisli function.
//
// Parameters:
//   - pred: An equality predicate
//
// Example:
//
//	import "github.com/IBM/fp-go/v2/eq"
//	equals42 := FromEq(eq.FromStrictEquals[int]())(42)
//	result := equals42(42)  // Some(42)
//	result := equals42(10)  // None
//
//go:inline
func FromEq[A any](pred eq.Eq[A]) func(A) Kleisli[A, A] {
	return F.Flow2(P.IsEqual(pred), FromPredicate[A])
}

// FromNillable converts a pointer to an Option.
// Returns Some if the pointer is non-nil, None otherwise.
//
// Parameters:
//   - a: A pointer that may be nil
//
// Example:
//
//	var ptr *int = nil
//	result := FromNillable(ptr) // None
//	val := 42
//	result := FromNillable(&val) // Some(&val)
func FromNillable[A any](a *A) (*A, bool) {
	return a, F.IsNonNil(a)
}

// Ap is the curried applicative functor for Option.
// Returns a function that applies an Option-wrapped function to the given Option value.
//
// Parameters:
//   - fa: The value of the Option
//   - faok: Whether the Option contains a value (true for Some, false for None)
//
// Example:
//
//	fa := Some(5)
//	applyTo5 := Ap[int](fa)
//	fab := Some(N.Mul(2))
//	result := applyTo5(fab) // Some(10)
func Ap[B, A any](fa A, faok bool) Operator[func(A) B, B] {
	if faok {
		return func(fab func(A) B, fabok bool) (b B, bok bool) {
			if fabok {
				return fab(fa), true
			}
			return
		}
	}
	return func(_ func(A) B, _ bool) (b B, bok bool) {
		return
	}
}

// Map returns a function that applies a transformation to the value inside an Option.
// If the Option is None, returns None.
//
// Parameters:
//   - f: A transformation function to apply to the Option value
//
// Example:
//
//	double := Map(N.Mul(2))
//	result := double(Some(5)) // Some(10)
//	result := double(None[int]()) // None
func Map[A, B any](f func(a A) B) Operator[A, B] {
	return func(fa A, faok bool) (b B, bok bool) {
		if faok {
			return f(fa), true
		}
		return
	}
}

// MapTo returns a function that replaces the value inside an Option with a constant.
//
// Parameters:
//   - b: The constant value to replace with
//
// Example:
//
//	replaceWith42 := MapTo[string, int](42)
//	result := replaceWith42(Some("hello")) // Some(42)
func MapTo[A, B any](b B) Operator[A, B] {
	return func(_ A, faok bool) (B, bool) {
		return b, faok
	}
}

// Fold provides a way to handle both Some and None cases of an Option.
// Returns a function that applies onNone if the Option is None, or onSome if it's Some.
//
// Parameters:
//   - onNone: Function to call when the Option is None
//   - onSome: Function to call when the Option is Some, receives the wrapped value
//
// Example:
//
//	handler := Fold(
//	    func() string { return "no value" },
//	    func(x int) string { return fmt.Sprintf("value: %d", x) },
//	)
//	result := handler(Some(42)) // "value: 42"
//	result := handler(None[int]()) // "no value"
func Fold[A, B any](onNone func() B, onSome func(A) B) func(A, bool) B {
	return func(a A, aok bool) B {
		if aok {
			return onSome(a)
		}
		return onNone()
	}
}

// GetOrElse returns a function that extracts the value from an Option or returns a default.
//
// Parameters:
//   - onNone: Function that provides the default value when the Option is None
//
// Example:
//
//	getOrZero := GetOrElse(func() int { return 0 })
//	result := getOrZero(Some(42)) // 42
//	result := getOrZero(None[int]()) // 0
func GetOrElse[A any](onNone func() A) func(A, bool) A {
	return func(a A, aok bool) A {
		if aok {
			return a
		}
		return onNone()
	}
}

// Chain returns a function that applies an Option-returning function to an Option value.
// This is the curried form of the monadic bind operation.
//
// Parameters:
//   - f: A function that takes a value and returns an Option
//
// Example:
//
//	validate := Chain(func(x int) (int, bool) {
//	    if x > 0 { return x * 2, true }
//	    return 0, false
//	})
//	result := validate(Some(5)) // Some(10)
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return func(a A, aok bool) (b B, bok bool) {
		if aok {
			return f(a)
		}
		return
	}
}

// ChainTo returns a function that ignores its input Option and returns a fixed Option.
//
// Parameters:
//   - b: The value of the replacement Option
//   - bok: Whether the replacement Option contains a value
//
// Example:
//
//	replaceWith := ChainTo(Some("hello"))
//	result := replaceWith(Some(42)) // Some("hello")
func ChainTo[A, B any](b B, bok bool) Operator[A, B] {
	return func(_ A, aok bool) (B, bool) {
		return b, bok && aok
	}
}

// ChainFirst returns a function that applies an Option-returning function but keeps the original value.
//
// Parameters:
//   - f: A function that takes a value and returns an Option (result is used only for success/failure)
//
// Example:
//
//	logAndKeep := ChainFirst(func(x int) (string, bool) {
//	    fmt.Println(x)
//	    return "logged", true
//	})
//	result := logAndKeep(Some(5)) // Some(5)
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return func(a A, aok bool) (A, bool) {
		if aok {
			_, bok := f(a)
			return a, bok
		}
		return a, false
	}
}

// Alt returns a function that provides an alternative Option if the input is None.
//
// Parameters:
//   - that: A function that provides an alternative Option
//
// Example:
//
//	withDefault := Alt(func() (int, bool) { return 0, true })
//	result := withDefault(Some(5)) // Some(5)
//	result := withDefault(None[int]()) // Some(0)
func Alt[A any](that func() (A, bool)) Operator[A, A] {
	return func(a A, aok bool) (A, bool) {
		if aok {
			return a, aok
		}
		return that()
	}
}

// Reduce folds an Option into a single value using a reducer function.
// If the Option is None, returns the initial value.
//
// Parameters:
//   - f: A reducer function that combines the accumulator with the Option value
//   - initial: The initial/default value to use
//
// Example:
//
//	sum := Reduce(func(acc, val int) int { return acc + val }, 0)
//	result := sum(Some(5)) // 5
//	result := sum(None[int]()) // 0
func Reduce[A, B any](f func(B, A) B, initial B) func(A, bool) B {
	return func(a A, aok bool) B {
		if aok {
			return f(initial, a)
		}
		return initial
	}
}

// Filter keeps the Option if it's Some and the predicate is satisfied, otherwise returns None.
//
// Parameters:
//   - pred: A predicate function to test the Option value
//
// Example:
//
//	isPositive := Filter(N.MoreThan(0))
//	result := isPositive(Some(5)) // Some(5)
//	result := isPositive(Some(-1)) // None
//	result := isPositive(None[int]()) // None
func Filter[A any](pred func(A) bool) Operator[A, A] {
	return func(a A, aok bool) (A, bool) {
		return a, aok && pred(a)
	}
}

// Flap returns a function that applies a value to an Option-wrapped function.
//
// Parameters:
//   - a: The value to apply to the function
//
// Example:
//
//	applyFive := Flap[int](5)
//	fab := Some(N.Mul(2))
//	result := applyFive(fab) // Some(10)
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return func(f func(A) B, fabok bool) (b B, bok bool) {
		if fabok {
			return f(a), true
		}
		return
	}
}
