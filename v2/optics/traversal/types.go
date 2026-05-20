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

package traversal

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Traversal is an optic that focuses on zero or more values of type A within a structure S.
	//
	// A Traversal allows you to view, modify, or traverse multiple values simultaneously within
	// a data structure. Unlike a Lens which focuses on exactly one value, or a Prism which
	// focuses on zero or one value, a Traversal can focus on any number of values (0, 1, or many).
	//
	// The Traversal type is defined as a function that takes a transformation function and
	// returns a function that applies this transformation within a higher-kinded type context.
	// This abstraction allows traversals to work with various effect types (Identity, Const,
	// Option, etc.) to support different operations like modification, folding, and filtering.
	//
	// Type Parameters:
	//   - S: The source type (the structure being traversed)
	//   - A: The focus type (the values being accessed or modified)
	//   - HKTS: Higher-kinded type for S (the effect context for the structure)
	//   - HKTA: Higher-kinded type for A (the effect context for the values)
	//
	// Common Operations:
	//
	// Modify: Transform all focused values
	//
	//	numbers := []int{1, 2, 3, 4}
	//	doubled := Modify[[]int, int](func(n int) int { return n * 2 })(arrayTraversal)(numbers)
	//	// Result: []int{2, 4, 6, 8}
	//
	// GetAll: Extract all focused values
	//
	//	numbers := []int{1, 2, 3, 4}
	//	values := GetAll[int](numbers)(arrayTraversal)
	//	// Result: []int{1, 2, 3, 4}
	//
	// FoldMap: Combine all focused values using a monoid
	//
	//	numbers := []int{1, 2, 3, 4}
	//	sum := FoldMap[[]int, int, int](func(n int) int { return n })(arrayTraversal)(numbers)
	//	// Result: 10
	//
	// Composition:
	//
	// Traversals compose naturally, allowing you to traverse nested structures:
	//
	//	type Company struct {
	//	    Departments []Department
	//	}
	//	type Department struct {
	//	    Employees []Employee
	//	}
	//	type Employee struct {
	//	    Salary int
	//	}
	//
	//	// Compose traversals to access all employee salaries
	//	allSalaries := F.Pipe1(
	//	    departmentsTraversal,
	//	    Compose[Company](employeesTraversal),
	//	    Compose[Company](salaryLens),
	//	)
	//
	//	// Give everyone a 10% raise
	//	updated := Modify[Company, int](func(s int) int {
	//	    return s * 110 / 100
	//	})(allSalaries)(company)
	//
	// Filtering:
	//
	// Traversals can be filtered to focus only on values matching a predicate:
	//
	//	import (
	//	    F "github.com/IBM/fp-go/v2/function"
	//	    "github.com/IBM/fp-go/v2/identity"
	//	    N "github.com/IBM/fp-go/v2/number"
	//	)
	//
	//	numbers := []int{-2, -1, 0, 1, 2, 3}
	//	isPositive := N.MoreThan(0)
	//
	//	// Create a filtered traversal
	//	positiveTraversal := F.Pipe1(
	//	    arrayTraversal,
	//	    Filter[[]int, int](identity.Of[int], identity.Map[int, int])(isPositive),
	//	)
	//
	//	// Double only positive numbers
	//	result := Modify[[]int, int](func(n int) int {
	//	    return n * 2
	//	})(positiveTraversal)(numbers)
	//	// Result: []int{-2, -1, 0, 2, 4, 6}
	//
	// Optics Hierarchy:
	//
	// Traversal is the most general optic in the hierarchy. Other optics can be
	// converted to traversals:
	//
	//	Iso[S, A]    →  Lens[S, A]    →  Optional[S, A]  →  Traversal[S, A]
	//	                                       ↑
	//	                                  Prism[S, A]
	//
	// Use Cases:
	//   - Modifying all elements in a collection
	//   - Extracting values from nested structures
	//   - Applying transformations to multiple fields
	//   - Filtering and transforming data simultaneously
	//   - Batch operations on structured data
	//
	// Laws:
	//
	// Traversals must satisfy the traversal laws:
	//  1. Identity: traverse(Identity) = Identity
	//  2. Composition: traverse(Compose(f, g)) = Compose(traverse(f), traverse(g))
	//  3. Naturality: t . traverse(f) = traverse(t . f) for natural transformation t
	//
	// Performance:
	//
	// Traversals are efficient for batch operations:
	//   - Single pass over the structure
	//   - No intermediate allocations for composed traversals
	//   - Type-safe with no reflection overhead
	//
	// See Also:
	//   - Modify: Transform all focused values
	//   - GetAll: Extract all focused values
	//   - FoldMap: Combine values using a monoid
	//   - Filter: Create a filtered traversal
	//   - Compose: Compose two traversals
	//   - optics/lens: For focusing on a single field
	//   - optics/prism: For focusing on a variant
	//   - optics/optional: For focusing on an optional value
	Traversal[S, A, HKTS, HKTA any] = G.Traversal[S, A, HKTS, HKTA]

	Predicate[A any] = predicate.Predicate[A]
)
