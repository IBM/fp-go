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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	TG "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// Compose composes two traversal endomorphisms to create a new traversal that focuses on nested values.
//
// This function is specifically designed for traversal endomorphisms, which work with endomorphism-based
// higher-kinded types. It is derived from the standard traversal composition in optics/traversal/generic:
// the inner traversal endomorphism AB is first converted into a regular traversal via ToTraversal, and the
// result is then composed with the outer traversal SA.
//
// Composition allows you to combine traversals to access deeply nested structures. When you compose
// traversal AB (which focuses on B within A) with traversal SA (which focuses on A within S), you get
// a traversal that focuses on B within S.
//
// The composition is associative. For traversal endomorphisms sa (A within S), ab (B within A)
// and bc (C within B), using the matching fmap at each level:
//
//	Compose(fmap)(bc)(Compose(fmap)(ab)(sa)) == Compose(fmap)(Compose(fmap)(bc)(ab))(sa)
//
// Type Parameters:
//   - S: The outer source type
//   - HKTES: Higher-kinded type for Endomorphism[S]
//   - A: The intermediate type
//   - HKTEA: Higher-kinded type for Endomorphism[A]
//   - B: The final focus type
//   - HKTA: Higher-kinded type for A
//   - HKTB: Higher-kinded type for B
//
// Parameters:
//   - fmap: Functor map operation to transform Endomorphism[A] to A within the effect context
//
// Returns:
//   - A function that takes a traversal AB and returns a function that takes a traversal SA
//     and returns the composed traversal SB
//
// Example:
//
//	import (
//	    thunk "github.com/IBM/fp-go/v2/context/readerioresult"
//	    "github.com/IBM/fp-go/v2/endomorphism"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	type Company struct {
//	    Departments []*Department
//	}
//	type Department struct {
//	    Employees []*Employee
//	}
//	type Employee struct {
//	    Salary int
//	}
//
//	// Traversal endomorphisms for each level, e.g. built with FromArrayLens and FromLens
//	var (
//	    deptTrav   Traversal[*Company, *Department, thunk.ReaderIOResult[endomorphism.Endomorphism[*Company]], thunk.ReaderIOResult[*Department]]
//	    empTrav    Traversal[*Department, *Employee, thunk.ReaderIOResult[endomorphism.Endomorphism[*Department]], thunk.ReaderIOResult[*Employee]]
//	    salaryTrav Traversal[*Employee, int, thunk.ReaderIOResult[endomorphism.Endomorphism[*Employee]], thunk.ReaderIOResult[int]]
//	)
//
//	// Compose to focus on all employees of a company
//	allEmployees := F.Pipe1(
//	    deptTrav,
//	    Compose[*Company, *Employee, thunk.ReaderIOResult[endomorphism.Endomorphism[*Company]], thunk.ReaderIOResult[*Employee]](
//	        thunk.Map[endomorphism.Endomorphism[*Department], *Department],
//	    )(empTrav),
//	)
//
//	// Further compose to focus on all salaries
//	allSalaries := F.Pipe1(
//	    allEmployees,
//	    Compose[*Company, int, thunk.ReaderIOResult[endomorphism.Endomorphism[*Company]], thunk.ReaderIOResult[int]](
//	        thunk.Map[endomorphism.Endomorphism[*Employee], *Employee],
//	    )(salaryTrav),
//	)
//
// See Also:
//   - optics/traversal/generic.Compose: Standard traversal composition
//   - MakeMonoid: Create a monoid for combining traversal endomorphisms
//   - Concat: Combine two traversal endomorphisms
func Compose[
	S, B, HKTES, HKTB, A, HKTEA, HKTA any](
	fmap functor.MapType[Endomorphism[A], A, HKTEA, HKTA],
) func(Traversal[A, B, HKTEA, HKTB]) func(Traversal[S, A, HKTES, HKTA]) Traversal[S, B, HKTES, HKTB] {
	return F.Flow2(
		ToTraversal[B, HKTB](fmap),
		TG.Compose[
			Traversal[A, B, HKTA, HKTB],
			Traversal[S, A, HKTES, HKTA],
			Traversal[S, B, HKTES, HKTB],
		],
	)
}
