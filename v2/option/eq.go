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

package option

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
)

// Eq constructs an equality predicate for Option[A] given an equality predicate for A.
// Two Options are equal if:
//   - Both are None, or
//   - Both are Some and their contained values are equal according to the provided Eq
//
// Example:
//
//	intEq := eq.FromStrictEquals[int]()
//	optEq := Eq(intEq)
//	optEq.Equals(Some(42), Some(42)) // true
//	optEq.Equals(Some(42), Some(43)) // false
//	optEq.Equals(None[int](), None[int]()) // true
//	optEq.Equals(Some(42), None[int]()) // false
func Eq[A any](a EQ.Eq[A]) EQ.Eq[Option[A]] {
	// some convenient shortcuts
	fld := Fold(
		F.Constant(Fold(F.ConstTrue, F.Constant1[A](false))),
		F.Flow2(F.Curry2(a.Equals), F.Bind1st(Fold[A, bool], F.ConstFalse)),
	)
	// convert to an equals predicate
	return EQ.FromEquals(F.Uncurry2(fld))
}

// FromStrictEquals constructs an Eq for Option[A] using Go's built-in equality (==) for type A.
// This is a convenience function for comparable types.
//
// Example:
//
//	optEq := FromStrictEquals[int]()
//	optEq.Equals(Some(42), Some(42)) // true
//	optEq.Equals(None[int](), None[int]()) // true
func FromStrictEquals[A comparable]() EQ.Eq[Option[A]] {
	return Eq(EQ.FromStrictEquals[A]())
}
