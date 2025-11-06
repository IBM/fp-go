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
	C "github.com/IBM/fp-go/v2/constraints"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ord"
)

// Ord constructs an ordering for Option[A] given an ordering for A.
// The ordering follows these rules:
//   - None is considered less than any Some value
//   - Two None values are equal
//   - Two Some values are compared using the provided Ord for A
//
// Example:
//
//	intOrd := ord.FromStrictCompare[int]()
//	optOrd := Ord(intOrd)
//	optOrd.Compare(None[int](), Some(5)) // -1 (None < Some)
//	optOrd.Compare(Some(3), Some(5)) // -1 (3 < 5)
//	optOrd.Compare(Some(5), Some(3)) // 1 (5 > 3)
//	optOrd.Compare(None[int](), None[int]()) // 0 (equal)
func Ord[A any](a ord.Ord[A]) ord.Ord[Option[A]] {
	// some convenient shortcuts
	fld := Fold(
		F.Constant(Fold(F.Constant(0), F.Constant1[A](-1))),
		F.Flow2(F.Curry2(a.Compare), F.Bind1st(Fold[A, int], F.Constant(1))),
	)
	// convert to an ordering predicate
	return ord.MakeOrd(F.Uncurry2(fld), Eq(ord.ToEq(a)).Equals)
}

// FromStrictCompare constructs an Ord for Option[A] using Go's built-in comparison operators for type A.
// This is a convenience function for ordered types (types that support <, >, ==).
//
// Example:
//
//	optOrd := FromStrictCompare[int]()
//	optOrd.Compare(Some(5), Some(10)) // -1
//	optOrd.Compare(None[int](), Some(5)) // -1
func FromStrictCompare[A C.Ordered]() ord.Ord[Option[A]] {
	return Ord(ord.FromStrictCompare[A]())
}
