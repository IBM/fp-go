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
	"github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/reader"
)

// Compose composes two traversal endomorphisms to create a new traversal that focuses on nested values.
//
// This function is specifically designed for traversal endomorphisms, which work with endomorphism-based
// higher-kinded types. Unlike the standard traversal composition in optics/traversal/generic, this version
// uses applicative operations to efficiently compose traversals that return endomorphisms.
//
// Composition allows you to combine traversals to access deeply nested structures. When you compose
// traversal AB (which focuses on B within A) with traversal SA (which focuses on A within S), you get
// a traversal that focuses on B within S.
//
// The composition is associative, meaning:
//
//	Compose(fmap)(ab)(Compose(fmap)(bc)(cd)) == Compose(fmap)(Compose(fmap)(ab)(bc))(cd)
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
//	    F "github.com/IBM/fp-go/v2/function"
//	    TLG "github.com/IBM/fp-go/v2/optics/traversalendo/lens/generic"
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
//	// Create traversals for nested structures
//	deptTraversal := departmentsTraversal[*Company, *Department]()
//	empTraversal := employeesTraversal[*Department, *Employee]()
//	salaryLens := TLG.FromLens[*Employee, int](thunk.Map)(salaryLensImpl)
//
//	// Compose to get all employees in a company
//	allEmployees := F.Pipe1(
//	    deptTraversal,
//	    Compose[*Company, thunk.Thunk[endomorphism.Endomorphism[*Company]], *Department](thunk.Map)(empTraversal),
//	)
//
//	// Further compose to access all salaries
//	allSalaries := F.Pipe1(
//	    allEmployees,
//	    Compose[*Company, thunk.Thunk[endomorphism.Endomorphism[*Company]], *Employee](thunk.Map)(salaryLens),
//	)
//
// See Also:
//   - optics/traversal/generic.Compose: Standard traversal composition
//   - MakeMonoid: Create a monoid for combining traversal endomorphisms
//   - Concat: Combine two traversal endomorphisms
func Compose[
	S, HKTES, A, HKTEA, B, HKTA, HKTB any](
	fmap functor.MapType[Endomorphism[A], A, HKTEA, HKTA],
) func(Traversal[A, B, HKTEA, HKTB]) func(Traversal[S, A, HKTES, HKTA]) Traversal[S, B, HKTES, HKTB] {
	readA := F.Flow2(
		endomorphism.Read[A],
		fmap,
	)
	return func(ab Traversal[A, B, HKTEA, HKTB]) func(Traversal[S, A, HKTES, HKTA]) Traversal[S, B, HKTES, HKTB] {
		return func(sa Traversal[S, A, HKTES, HKTA]) Traversal[S, B, HKTES, HKTB] {
			return func(f func(B) HKTB) func(S) HKTES {
				return F.Pipe1(
					F.Pipe1(
						readA,
						F.Pipe2(
							f,
							ab,
							reader.Ap[HKTA],
						),
					),
					sa,
				)
			}
		}
	}
}
