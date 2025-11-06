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

package io

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	INTE "github.com/IBM/fp-go/v2/internal/eq"
)

// Eq implements the equals predicate for values contained in the IO monad.
// It lifts an Eq[A] into an Eq[IO[A]] by executing both IO computations
// and comparing their results.
//
// Example:
//
//	intEq := eq.FromStrictEquals[int]()
//	ioEq := io.Eq(intEq)
//	result := ioEq.Equals(io.Of(42), io.Of(42)) // true
func Eq[A any](e EQ.Eq[A]) EQ.Eq[IO[A]] {
	// comparator for the monad
	eq := INTE.Eq(
		MonadMap[A, func(A) bool],
		MonadAp[A, bool],
		e,
	)
	// eagerly execute
	return EQ.FromEquals(func(l, r IO[A]) bool {
		return eq(l, r)()
	})
}

// FromStrictEquals constructs an Eq[IO[A]] from the canonical comparison function
// for comparable types. This is a convenience function that combines Eq with
// the standard equality operator.
//
// Example:
//
//	ioEq := io.FromStrictEquals[int]()
//	result := ioEq.Equals(io.Of(42), io.Of(42)) // true
func FromStrictEquals[A comparable]() EQ.Eq[IO[A]] {
	return Eq(EQ.FromStrictEquals[A]())
}
