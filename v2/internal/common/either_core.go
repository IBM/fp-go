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

package common

type (
	// Either defines a data structure that logically holds either an E or an A. The flag discriminates the cases
	Either[E, A any] struct {
		a A
		e E
		l bool
	}

	// EitherKleisli represents a Kleisli arrow for the Either monad.
	// It's a function from A to Either[E, B], used for composing operations that may fail.
	EitherKleisli[E, A, B any] = func(A) Either[E, B]

	// EitherOperator represents a function that transforms one Either into another.
	// It takes an Either[E, A] and produces an Either[E, B].
	EitherOperator[E, A, B any] = EitherKleisli[E, Either[E, A], B]

	EitherTraversable[E, A, B, GA, GB any] = func(EitherKleisli[E, A, B]) EitherKleisli[E, GA, GB]
)

// EitherIsLeft tests if the Either is a Left value.
// Rather use [EitherFold] or [EitherMonadFold] if you need to access the values.
// Inverse is [EitherIsRight].
//
// Example:
//
//	either.EitherIsLeft(either.Left[int](errors.New("err"))) // true
//	either.EitherIsLeft(either.Right[error](42)) // false
//
//go:inline
func EitherIsLeft[E, A any](val Either[E, A]) bool {
	return val.l
}

// EitherIsRight tests if the Either is a Right value.
// Rather use [EitherFold] or [EitherMonadFold] if you need to access the values.
// Inverse is [EitherIsLeft].
//
// Example:
//
//	either.EitherIsRight(either.Right[error](42)) // true
//	either.EitherIsRight(either.Left[int](errors.New("err"))) // false
//
//go:inline
func EitherIsRight[E, A any](val Either[E, A]) bool {
	return !val.l
}

// EitherLeft creates a new Either representing a EitherLeft (error/failure) value.
// By convention, EitherLeft represents the error case.
//
// Example:
//
//	result := either.EitherLeft[int](errors.New("something went wrong"))
//
//go:inline
func EitherLeft[A, E any](value E) Either[E, A] {
	return Either[E, A]{l: true, e: value}
}

// EitherRight creates a new Either representing a EitherRight (success) value.
// By convention, EitherRight represents the success case.
//
// Example:
//
//	result := either.EitherRight[error](42)
//
//go:inline
func EitherRight[E, A any](value A) Either[E, A] {
	return Either[E, A]{a: value}
}

// EitherMonadFold extracts the value from an Either by providing handlers for both cases.
// This is the fundamental pattern matching operation for Either.
//
// Example:
//
//	result := either.EitherMonadFold(
//	    either.Right[error](42),
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Value: %d", n) },
//	) // "Value: 42"
//
//go:inline
func EitherMonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if !ma.l {
		return onRight(ma.a)
	}
	return onLeft(ma.e)
}

// EitherUnwrap converts an Either into the idiomatic Go tuple (value, error).
// For Right values, returns (value, zero-error).
// For Left values, returns (zero-value, error).
//
// Example:
//
//	val, err := either.EitherUnwrap(either.Right[error](42)) // 42, nil
//	val, err := either.EitherUnwrap(either.Left[int](errors.New("fail"))) // 0, error
//
//go:inline
func EitherUnwrap[E, A any](ma Either[E, A]) (A, E) {
	return ma.a, ma.e
}
