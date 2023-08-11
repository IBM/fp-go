// Copyright (c) 2023 IBM Corp.
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
	C "github.com/IBM/fp-go/constraints"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/ord"
)

// Constructs an order for [Option]
func Ord[A any](a ord.Ord[A]) ord.Ord[Option[A]] {
	// some convenient shortcuts
	fld := Fold(
		F.Constant(Fold(F.Constant(0), F.Constant1[A](-1))),
		F.Flow2(F.Curry2(a.Compare), F.Bind1st(Fold[A, int], F.Constant(1))),
	)
	// convert to an ordering predicate
	return ord.MakeOrd(F.Uncurry2(fld), Eq(ord.ToEq(a)).Equals)
}

// FromStrictCompare constructs an [Ord] from the canonical comparison function
func FromStrictCompare[A C.Ordered]() ord.Ord[Option[A]] {
	return Ord(ord.FromStrictCompare[A]())
}
