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
//go:build either_any

package either

import (
	"fmt"
)

type (
	either struct {
		value   any
		isRight bool
	}

	// Either defines a data structure that logically holds either an E or an A. The flag discriminates the cases
	Either[E, A any] either
)

// String prints some debug info for the object
//
//go:noinline
func eitherString(s *either) string {
	if s.isRight {
		return fmt.Sprintf("Right[%T](%v)", s.value, s.value)
	}
	return fmt.Sprintf("Left[%T](%v)", s.value, s.value)
}

// Format prints some debug info for the object
//
//go:noinline
func eitherFormat(e *either, f fmt.State, c rune) {
	switch c {
	case 's':
		fmt.Fprint(f, eitherString(e))
	default:
		fmt.Fprint(f, eitherString(e))
	}
}

// String prints some debug info for the object
func (s Either[E, A]) String() string {
	return eitherString((*either)(&s))
}

// Format prints some debug info for the object
func (s Either[E, A]) Format(f fmt.State, c rune) {
	eitherFormat((*either)(&s), f, c)
}

// IsLeft tests if the Either is a Left value.
// Rather use [Fold] or [MonadFold] if you need to access the values.
// Inverse is [IsRight].
//
// Example:
//
//	either.IsLeft(either.Left[int](errors.New("err"))) // true
//	either.IsLeft(either.Right[error](42)) // false
//
//go:inline
func IsLeft[E, A any](val Either[E, A]) bool {
	return !val.isRight
}

// IsRight tests if the Either is a Right value.
// Rather use [Fold] or [MonadFold] if you need to access the values.
// Inverse is [IsLeft].
//
// Example:
//
//	either.IsRight(either.Right[error](42)) // true
//	either.IsRight(either.Left[int](errors.New("err"))) // false
//
//go:inline
func IsRight[E, A any](val Either[E, A]) bool {
	return val.isRight
}

// Left creates a new Either representing a Left (error/failure) value.
// By convention, Left represents the error case.
//
// Example:
//
//	result := either.Left[int](errors.New("something went wrong"))
//
//go:inline
func Left[A, E any](value E) Either[E, A] {
	return Either[E, A]{value, false}
}

// Right creates a new Either representing a Right (success) value.
// By convention, Right represents the success case.
//
// Example:
//
//	result := either.Right[error](42)
//
//go:inline
func Right[E, A any](value A) Either[E, A] {
	return Either[E, A]{value, true}
}

// MonadFold extracts the value from an Either by providing handlers for both cases.
// This is the fundamental pattern matching operation for Either.
//
// Example:
//
//	result := either.MonadFold(
//	    either.Right[error](42),
//	    func(err error) string { return "Error: " + err.Error() },
//	    func(n int) string { return fmt.Sprintf("Value: %d", n) },
//	) // "Value: 42"
//
//go:inline
func MonadFold[E, A, B any](ma Either[E, A], onLeft func(e E) B, onRight func(a A) B) B {
	if ma.isRight {
		return onRight(ma.value.(A))
	}
	return onLeft(ma.value.(E))
}

// Unwrap converts an Either into the idiomatic Go tuple (value, error).
// For Right values, returns (value, zero-error).
// For Left values, returns (zero-value, error).
//
// Example:
//
//	val, err := either.Unwrap(either.Right[error](42)) // 42, nil
//	val, err := either.Unwrap(either.Left[int](errors.New("fail"))) // 0, error
//
//go:inline
func Unwrap[E, A any](ma Either[E, A]) (A, E) {
	if ma.isRight {
		var e E
		return ma.value.(A), e
	} else {
		var a A
		return a, ma.value.(E)
	}
}
