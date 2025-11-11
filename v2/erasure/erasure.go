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

// Package erasure provides utilities for type erasure and type-safe conversion between
// generic types and the any type. This is useful when working with heterogeneous collections
// or when interfacing with APIs that require type erasure.
//
// The package provides functions to:
//   - Erase typed values to any (via pointers)
//   - Unerase any values back to their original types
//   - Safely unerase with error handling
//   - Convert type-safe functions to erased functions
//
// Example usage:
//
//	// Basic erasure and unerasure
//	erased := erasure.Erase(42)
//	value := erasure.Unerase[int](erased) // value == 42
//
//	// Safe unerasure with error handling
//	result := erasure.SafeUnerase[int](erased)
//	// result is Either[error, int]
//
//	// Function erasure
//	typedFunc := strconv.Itoa
//	erasedFunc := erasure.Erase1(typedFunc)
//	result := erasedFunc(erasure.Erase(42)) // returns erased "42"
package erasure

import (
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
)

// Erase converts a variable of type T to an any by returning a pointer to that variable.
// This allows type-safe storage and retrieval of values in heterogeneous collections.
//
// The value is stored as a pointer, which means the type information is preserved
// and can be recovered using Unerase or SafeUnerase.
//
// Example:
//
//	erased := Erase(42)
//	// erased is any, but internally holds *int
func Erase[T any](t T) any {
	return &t
}

// Unerase converts an erased variable back to its original value.
// This function panics if the type assertion fails, so use SafeUnerase
// for error handling.
//
// Example:
//
//	erased := Erase(42)
//	value := Unerase[int](erased) // value == 42
//
// Panics if the erased value is not of type *T.
func Unerase[T any](t any) T {
	return *t.(*T)
}

// SafeUnerase converts an erased variable back to its original value with error handling.
// Returns Either[error, T] where Left contains an error if the type assertion fails,
// and Right contains the unerased value if successful.
//
// This is the safe alternative to Unerase that doesn't panic on type mismatch.
//
// Example:
//
//	erased := Erase(42)
//	result := SafeUnerase[int](erased)
//	// result is Right(42)
//
//	wrongType := SafeUnerase[string](erased)
//	// wrongType is Left(error) with message about type mismatch
func SafeUnerase[T any](t any) E.Either[error, T] {
	return F.Pipe2(
		t,
		E.ToType[*T](errors.OnSome[any]("Value of type [%T] is not erased")),
		E.Map[error](F.Deref[T]),
	)
}

// Erase0 converts a type-safe nullary function into an erased function.
// The resulting function returns an erased value.
//
// Example:
//
//	typedFunc := func() int { return 42 }
//	erasedFunc := Erase0(typedFunc)
//	result := erasedFunc() // returns erased 42
func Erase0[T1 any](f func() T1) func() any {
	return F.Nullary2(f, Erase[T1])
}

// Erase1 converts a type-safe unary function into an erased function.
// The resulting function takes an erased argument and returns an erased value.
//
// Example:
//
//	typedFunc := strconv.Itoa
//	erasedFunc := Erase1(typedFunc)
//	result := erasedFunc(Erase(42)) // returns erased "42"
func Erase1[T1, T2 any](f func(T1) T2) func(any) any {
	return F.Flow3(
		Unerase[T1],
		f,
		Erase[T2],
	)
}

// Erase2 converts a type-safe binary function into an erased function.
// The resulting function takes two erased arguments and returns an erased value.
//
// Example:
//
//	typedFunc := func(x, y int) int { return x + y }
//	erasedFunc := Erase2(typedFunc)
//	result := erasedFunc(Erase(10), Erase(32)) // returns erased 42
func Erase2[T1, T2, T3 any](f func(T1, T2) T3) func(any, any) any {
	return func(t1, t2 any) any {
		return Erase(f(Unerase[T1](t1), Unerase[T2](t2)))
	}
}
