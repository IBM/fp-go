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

// Package generic provides generic implementations of traversal optics.
//
// This package contains the core generic implementations of traversal operations
// that work with higher-kinded types. These implementations are used by the
// concrete traversal package to provide type-safe traversal operations.
//
// A Traversal is an optic that focuses on zero or more values of type A within
// a structure S. Unlike lenses (which focus on exactly one value) or prisms
// (which focus on zero or one value), traversals can focus on any number of values.
//
// The generic package uses higher-kinded types (HKTS, HKTA) to abstract over
// different effect types like Identity, Const, Option, etc. This allows the same
// traversal to be used for different operations:
//   - With Identity: Modify all focused values
//   - With Const: Extract or fold all focused values
//   - With Option: Optionally modify values
//
// See Also:
//   - optics/traversal: Concrete traversal operations
//   - optics/lens: For focusing on a single field
//   - optics/prism: For focusing on a variant
package generic

import (
	AR "github.com/IBM/fp-go/v2/array/generic"
	C "github.com/IBM/fp-go/v2/constant"
	"github.com/IBM/fp-go/v2/endomorphism"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/pointed"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/predicate"
)

type (
	// Traversal is a generic optic that focuses on zero or more values of type A within a structure S.
	//
	// A Traversal is represented as a function that takes a transformation function and returns
	// a function that applies this transformation within a higher-kinded type context. This
	// abstraction allows traversals to work with various effect types to support different
	// operations like modification, folding, and filtering.
	//
	// Type Parameters:
	//   - S: The source type (the structure being traversed)
	//   - A: The focus type (the values being accessed or modified)
	//   - HKTS: Higher-kinded type for S (the effect context for the structure)
	//   - HKTA: Higher-kinded type for A (the effect context for the values)
	//
	// The function signature func(func(A) HKTA) func(S) HKTS means:
	//   - Given a way to transform A into HKTA (the transformation in context)
	//   - Return a way to transform S into HKTS (the structure in context)
	//
	// Common instantiations:
	//   - Traversal[S, A, S, A]: For modification using Identity
	//   - Traversal[S, A, Const[M, S], Const[M, A]]: For folding using Const
	//   - Traversal[S, A, Option[S], Option[A]]: For optional modification
	//
	// Example:
	//
	//	// A traversal for array elements
	//	arrayTraversal := func(f func(int) Identity[int]) func([]int) Identity[[]int] {
	//	    return func(arr []int) Identity[[]int] {
	//	        result := make([]int, len(arr))
	//	        for i, v := range arr {
	//	            result[i] = f(v).Value
	//	        }
	//	        return Identity{Value: result}
	//	    }
	//	}
	//
	// See Also:
	//   - Compose: Compose two traversals
	//   - Filter: Filter traversal targets by predicate
	//   - FoldMap: Fold traversal targets using a monoid
	//   - GetAll: Extract all traversal targets
	Traversal[S, A, HKTS, HKTA any] = func(func(A) HKTA) func(S) HKTS
)

// Compose composes two traversals to create a new traversal that focuses on nested values.
//
// Composition allows you to combine traversals to access deeply nested structures.
// When you compose traversal AB (which focuses on B within A) with traversal SA
// (which focuses on A within S), you get a traversal that focuses on B within S.
//
// The composition is associative, meaning:
//
//	Compose(ab)(Compose(bc)(cd)) == Compose(Compose(ab)(bc))(cd)
//
// Type Parameters:
//   - TAB: Traversal from A to B
//   - TSA: Traversal from S to A
//   - TSB: Resulting traversal from S to B
//   - S: The outer source type
//   - A: The intermediate type
//   - B: The final focus type
//   - HKTS, HKTA, HKTB: Higher-kinded types for the respective types
//
// Parameters:
//   - ab: A traversal that focuses on B within A
//
// Returns:
//   - A function that takes a traversal SA and returns a composed traversal SB
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	)
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
//	// Traversal for departments in a company
//	deptTraversal := departmentsTraversal[Company]()
//
//	// Traversal for employees in a department
//	empTraversal := employeesTraversal[Department]()
//
//	// Compose to get all employees in a company
//	allEmployees := F.Pipe1(
//	    deptTraversal,
//	    Compose[Department, Company](empTraversal),
//	)
//
//	// Further compose to access all salaries
//	allSalaries := F.Pipe1(
//	    allEmployees,
//	    Compose[Employee, Company](salaryLens),
//	)
//
// See Also:
//   - Filter: Filter traversal targets by predicate
//   - optics/traversal.Compose: Concrete version with Identity
func Compose[
	TAB ~func(func(B) HKTB) func(A) HKTA,
	TSA ~func(func(A) HKTA) func(S) HKTS,
	TSB ~func(func(B) HKTB) func(S) HKTS,
	S, A, B, HKTS, HKTA, HKTB any](ab TAB) func(TSA) TSB {
	return func(sa TSA) TSB {
		return F.Flow2(ab, sa)
	}
}

// FromTraversable creates a traversal from a traversable data structure.
//
// This function converts a traverse function (which operates on a traversable
// data structure like arrays, trees, or other containers) into a traversal optic.
// The traverse function must satisfy the traversable laws.
//
// A traversable data structure is one that can be traversed from left to right,
// performing an action on each element and collecting the results. Common examples
// include arrays, slices, trees, and other container types.
//
// Type Parameters:
//   - TAB: The resulting traversal type
//   - A: The element type within the traversable structure
//   - HKTTA: Higher-kinded type for the traversable structure
//   - HKTFA: Higher-kinded type for the transformed elements
//   - HKTAA: Higher-kinded type for the result
//
// Parameters:
//   - traverseF: A traverse function that applies an effectful function to each element
//
// Returns:
//   - A traversal that can be used with other optics operations
//
// Example:
//
//	import (
//	    AR "github.com/IBM/fp-go/v2/array"
//	    "github.com/IBM/fp-go/v2/identity"
//	)
//
//	// Create a traversal from array's traverse function
//	arrayTraversal := FromTraversable(AR.Traverse[int, Identity[int], Identity[[]int]])
//
//	// Use the traversal to modify all elements
//	numbers := []int{1, 2, 3, 4}
//	doubled := arrayTraversal(func(n int) Identity[int] {
//	    return identity.Of(n * 2)
//	})(numbers)
//
// See Also:
//   - array.Traverse: Traverse function for arrays
//   - optics/traversal/array: Pre-built array traversals
func FromTraversable[
	TAB ~func(func(A) HKTFA) func(HKTTA) HKTAA,
	A,
	HKTTA,
	HKTFA,
	HKTAA any](
	traverseF func(HKTTA, func(A) HKTFA) HKTAA,
) TAB {
	return F.Bind1st(F.Bind2nd[HKTTA, func(A) HKTFA, HKTAA], traverseF)
}

// FoldMap maps each target to a monoid and combines the results.
//
// This function allows you to extract and combine all values focused by a traversal
// using a monoid operation. It maps each focused value through a function f that
// produces a monoid value, then combines all these values using the monoid's
// associative operation.
//
// FoldMap is useful for aggregating data from multiple locations in a structure,
// such as summing numbers, concatenating strings, or collecting values into a list.
//
// Type Parameters:
//   - S: The source type
//   - M: The monoid type (must have an associative binary operation and identity)
//   - A: The focus type
//
// Parameters:
//   - f: A function that maps each focused value to a monoid value
//
// Returns:
//   - A function that takes a traversal and returns a function that folds the structure
//
// Example:
//
//	import (
//	    C "github.com/IBM/fp-go/v2/constant"
//	    F "github.com/IBM/fp-go/v2/function"
//	    N "github.com/IBM/fp-go/v2/number"
//	)
//
//	type Person struct {
//	    Name string
//	    Age  int
//	}
//
//	people := []Person{
//	    {Name: "Alice", Age: 30},
//	    {Name: "Bob", Age: 25},
//	    {Name: "Carol", Age: 35},
//	}
//
//	// Sum all ages using FoldMap
//	totalAge := FoldMap[[]Person, int, Person](func(p Person) int {
//	    return p.Age
//	})(peopleTraversal)(people)
//	// Result: 90
//
//	// Collect all names
//	names := FoldMap[[]Person, []string, Person](func(p Person) []string {
//	    return []string{p.Name}
//	})(peopleTraversal)(people)
//	// Result: []string{"Alice", "Bob", "Carol"}
//
// See Also:
//   - Fold: Fold without mapping (identity function)
//   - GetAll: Extract all values into a slice
//   - monoid: Monoid operations
func FoldMap[S, M, A any](f func(A) M) func(sa Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
	return func(sa Traversal[S, A, C.Const[M, S], C.Const[M, A]]) func(S) M {
		return F.Flow2(
			F.Pipe1(
				F.Flow2(f, C.Make[M, A]),
				sa,
			),
			C.Unwrap[M, S],
		)
	}
}

// Fold combines all targets of a traversal using their monoid operation.
//
// This is a specialized version of FoldMap that uses the identity function,
// meaning it directly combines the focused values without any transformation.
// The focused type A must form a monoid (have an associative binary operation
// and an identity element).
//
// Type Parameters:
//   - S: The source type
//   - A: The focus type (must be a monoid)
//
// Parameters:
//   - sa: A traversal that focuses on monoid values
//
// Returns:
//   - A function that folds the structure into a single monoid value
//
// Example:
//
//	import (
//	    C "github.com/IBM/fp-go/v2/constant"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Sum all numbers in a nested structure
//	numbers := [][]int{{1, 2}, {3, 4}, {5}}
//	sum := Fold[[][]int, int](nestedArrayTraversal)(numbers)
//	// Result: 15
//
//	// Concatenate all strings
//	words := [][]string{{"Hello", " "}, {"World", "!"}}
//	text := Fold[[][]string, string](nestedArrayTraversal)(words)
//	// Result: "Hello World!"
//
// See Also:
//   - FoldMap: Fold with a mapping function
//   - GetAll: Extract all values into a slice
func Fold[S, A any](sa Traversal[S, A, C.Const[A, S], C.Const[A, A]]) func(S) A {
	return FoldMap[S](F.Identity[A])(sa)
}

// GetAll extracts all targets of a traversal into a slice.
//
// This function collects all values focused by a traversal into a slice,
// preserving their order. It's useful when you need to extract multiple
// values from a structure for inspection or further processing.
//
// GetAll is implemented using FoldMap with the array monoid, which means
// it efficiently collects values in a single pass through the structure.
//
// Type Parameters:
//   - GA: The slice type (must be []A)
//   - S: The source type
//   - A: The focus type (element type)
//
// Parameters:
//   - s: The source value to extract from
//
// Returns:
//   - A function that takes a traversal and returns all focused values as a slice
//
// Example:
//
//	import (
//	    C "github.com/IBM/fp-go/v2/constant"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	type Company struct {
//	    Departments []Department
//	}
//	type Department struct {
//	    Name      string
//	    Employees []Employee
//	}
//	type Employee struct {
//	    Name   string
//	    Salary int
//	}
//
//	company := Company{
//	    Departments: []Department{
//	        {Name: "Engineering", Employees: []Employee{
//	            {Name: "Alice", Salary: 100000},
//	            {Name: "Bob", Salary: 90000},
//	        }},
//	        {Name: "Sales", Employees: []Employee{
//	            {Name: "Carol", Salary: 80000},
//	        }},
//	    },
//	}
//
//	// Get all employee salaries
//	salaries := GetAll[[]int](company)(allSalariesTraversal)
//	// Result: []int{100000, 90000, 80000}
//
//	// Get all employee names
//	names := GetAll[[]string](company)(allNamesTraversal)
//	// Result: []string{"Alice", "Bob", "Carol"}
//
// See Also:
//   - FoldMap: Fold with a custom mapping function
//   - Fold: Fold without collecting into a slice
//   - optics/traversal.GetAll: Concrete version with Identity
func GetAll[GA ~[]A, S, A any](s S) func(sa Traversal[S, A, C.Const[GA, S], C.Const[GA, A]]) GA {
	fmap := FoldMap[S](AR.Of[GA, A])
	return func(sa Traversal[S, A, C.Const[GA, S], C.Const[GA, A]]) GA {
		return fmap(sa)(s)
	}
}

// Filter creates a function that filters the targets of a traversal based on a predicate.
//
// This function allows you to refine a traversal to only focus on values that satisfy
// a given predicate. It works by converting the predicate into a prism, then converting
// that prism into a traversal, and finally composing it with the original traversal.
//
// The filtering is selective: when modifying values through the filtered traversal,
// only values that satisfy the predicate will be transformed. Values that don't
// satisfy the predicate remain unchanged.
//
// Type Parameters:
//   - S: The source type
//   - A: The focus type (the values being filtered)
//   - HKTS: Higher-kinded type for S (functor/applicative context)
//   - HKTA: Higher-kinded type for A (functor/applicative context)
//
// Parameters:
//   - fof: Function to lift A into the higher-kinded type HKTA (pure/of operation)
//   - fmap: Function to map over HKTA (functor map operation)
//
// Returns:
//   - A function that takes a predicate and returns an endomorphism on traversals
//
// Example:
//
//	import (
//	    AR "github.com/IBM/fp-go/v2/array"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	    N "github.com/IBM/fp-go/v2/number"
//	    AI "github.com/IBM/fp-go/v2/optics/traversal/array/identity"
//	)
//
//	// Create a traversal for array elements
//	arrayTraversal := AI.FromArray[int]()
//	baseTraversal := F.Pipe1(
//	    Id[[]int, []int](),
//	    Compose[[]int, []int, []int, int](arrayTraversal),
//	)
//
//	// Filter to only positive numbers
//	isPositive := N.MoreThan(0)
//	filteredTraversal := F.Pipe1(
//	    baseTraversal,
//	    Filter[[]int, int](identity.Of[int], identity.Map[int, int])(isPositive),
//	)
//
//	// Double only positive numbers
//	numbers := []int{-2, -1, 0, 1, 2, 3}
//	result := filteredTraversal(func(n int) int { return n * 2 })(numbers)
//	// result: [-2, -1, 0, 2, 4, 6]
//
// See Also:
//   - prism.FromPredicate: Creates a prism from a predicate
//   - prism.AsTraversal: Converts a prism to a traversal
//   - Compose: Composes two traversals
func Filter[
	S, HKTS, A, HKTA any](
	fof pointed.OfType[A, HKTA],
	fmap functor.MapType[A, A, HKTA, HKTA],
) func(predicate.Predicate[A]) endomorphism.Endomorphism[Traversal[S, A, HKTS, HKTA]] {
	return F.Flow3(
		prism.FromPredicate,
		prism.AsTraversal[Traversal[A, A, HKTA, HKTA]](fof, fmap),
		Compose[
			Traversal[A, A, HKTA, HKTA],
			Traversal[S, A, HKTS, HKTA],
			Traversal[S, A, HKTS, HKTA]],
	)
}

// Empty creates an empty traversal that focuses on no values.
//
// An empty traversal is the identity element of the traversal monoid. When composed
// with another traversal using Concat, it has no effect. This is useful as a base
// case when building traversals dynamically or when you need a traversal that
// performs no operations.
//
// The empty traversal simply lifts the source value S into the higher-kinded type
// context HKTS without applying any transformation to focused values.
//
// Type Parameters:
//   - A: The focus type (not used, as no values are focused)
//   - HKTA: Higher-kinded type for A (not used)
//   - S: The source type
//   - HKTS: Higher-kinded type for S
//
// Parameters:
//   - fof: Function to lift S into the higher-kinded type HKTS (pure/of operation)
//
// Returns:
//   - A traversal that focuses on no values
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	)
//
//	// Create an empty traversal
//	emptyTrav := Empty[int, identity.Identity[int], []int, identity.Identity[[]int]](
//	    identity.Of[[]int],
//	)
//
//	// Applying it returns the original value unchanged
//	numbers := []int{1, 2, 3}
//	result := emptyTrav(func(n int) identity.Identity[int] {
//	    return identity.Of(n * 2)
//	})(numbers)
//	// result.Value == []int{1, 2, 3} (unchanged)
//
//	// Empty is the identity for Concat
//	trav := someTraversal
//	F.Pipe1(trav, Concat(fchain)(Empty(fof))) // same as trav
//	F.Pipe1(Empty(fof), Concat(fchain)(trav)) // same as trav
//
// See Also:
//   - Concat: Combine two traversals
//   - Monoid: Create a monoid instance for traversals
func Empty[A, HKTA, S, HKTS any](
	fof pointed.OfType[S, HKTS],
) Traversal[S, A, HKTS, HKTA] {
	return F.Constant1[func(A) HKTA](fof)
}

// Concat combines two traversals into a single traversal that focuses on all values from both.
//
// This function creates a new traversal that applies both input traversals in sequence,
// combining their effects using the chain operation. The resulting traversal will focus
// on all values that either the left or right traversal focuses on.
//
// Concat is associative, meaning:
//
//	Concat(a)(Concat(b)(c)) == Concat(Concat(a)(b))(c)
//
// Together with Empty, Concat forms a monoid for traversals, allowing you to combine
// multiple traversals using standard monoid operations like fold or reduce.
//
// Type Parameters:
//   - A: The focus type
//   - HKTA: Higher-kinded type for A
//   - S: The source type
//   - HKTS: Higher-kinded type for S
//
// Parameters:
//   - fchain: Function to chain/sequence effects in the higher-kinded type context
//
// Returns:
//   - A function that takes two traversals and returns their concatenation
//
// Example:
//
//	import (
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	)
//
//	type Person struct {
//	    FirstName string
//	    LastName  string
//	}
//
//	// Traversal for first name
//	firstNameTrav := /* ... */
//
//	// Traversal for last name
//	lastNameTrav := /* ... */
//
//	// Combine to traverse both names
//	bothNamesTrav := Concat[string, identity.Identity[string]](
//	    identity.Chain[Person, Person],
//	)(firstNameTrav, lastNameTrav)
//
//	person := Person{FirstName: "John", LastName: "Doe"}
//
//	// Modify both names
//	result := bothNamesTrav(func(name string) identity.Identity[string] {
//	    return identity.Of(strings.ToUpper(name))
//	})(person)
//	// result.Value == Person{FirstName: "JOHN", LastName: "DOE"}
//
// See Also:
//   - Empty: Create an empty traversal (identity element)
//   - Monoid: Create a monoid instance for traversals
//   - Compose: Compose traversals for nested access
func Concat[A, HKTA, S, HKTS any](
	fchain chain.ChainType[S, HKTS, HKTS],
) func(l, r Traversal[S, A, HKTS, HKTA]) Traversal[S, A, HKTS, HKTA] {
	return func(l, r Traversal[S, A, HKTS, HKTA]) Traversal[S, A, HKTS, HKTA] {
		return func(f func(A) HKTA) func(S) HKTS {
			return F.Flow2(
				r(f),
				fchain(l(f)),
			)
		}
	}
}

// Monoid creates a monoid instance for traversals.
//
// A monoid provides two operations: an identity element (Empty) and an associative
// binary operation (Concat). This allows traversals to be combined using standard
// monoid operations, making it easy to build complex traversals from simpler ones.
//
// The monoid laws are satisfied:
//  1. Identity: Concat(Empty, x) == x and Concat(x, Empty) == x
//  2. Associativity: Concat(Concat(a, b), c) == Concat(a, Concat(b, c))
//
// This monoid instance enables you to use traversals with generic monoid operations
// like fold, reduce, or mconcat to combine multiple traversals into one.
//
// Note: This monoid works sequentially, processing traversals one after another using
// the chain operation. For better performance with parallel composition, consider using
// optics/traversalendo/generic.MakeMonoid instead, which uses applicative operations
// to combine traversals more efficiently.
//
// Type Parameters:
//   - S: The source type
//   - A: The focus type
//   - HKTS: Higher-kinded type for S
//   - HKTA: Higher-kinded type for A
//
// Parameters:
//   - fof: Function to lift S into the higher-kinded type HKTS (for Empty)
//   - fchain: Function to chain/sequence effects (for Concat)
//
// Returns:
//   - A monoid instance for traversals
//
// Example:
//
//	import (
//	    AR "github.com/IBM/fp-go/v2/array"
//	    F "github.com/IBM/fp-go/v2/function"
//	    "github.com/IBM/fp-go/v2/identity"
//	    M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	type Record struct {
//	    Field1 int
//	    Field2 int
//	    Field3 int
//	}
//
//	// Create individual field traversals
//	field1Trav := /* ... */
//	field2Trav := /* ... */
//	field3Trav := /* ... */
//
//	// Create monoid for traversals
//	travMonoid := Monoid[Record, int, identity.Identity[Record], identity.Identity[int]](
//	    identity.Of[Record],
//	    identity.Chain[Record, Record],
//	)
//
//	// Combine all field traversals using monoid
//	allFieldsTrav := M.Fold(travMonoid)([]Traversal[Record, int, ...]{
//	    field1Trav,
//	    field2Trav,
//	    field3Trav,
//	})
//
//	// Now modify all fields at once
//	record := Record{Field1: 1, Field2: 2, Field3: 3}
//	result := allFieldsTrav(func(n int) identity.Identity[int] {
//	    return identity.Of(n * 10)
//	})(record)
//	// result.Value == Record{Field1: 10, Field2: 20, Field3: 30}
//
// See Also:
//   - Empty: The identity element of the monoid
//   - Concat: The binary operation of the monoid
//   - monoid.Fold: Combine multiple monoid values
//   - monoid.MakeMonoid: Create a monoid instance
func Monoid[S, A, HKTS, HKTA any](
	fof pointed.OfType[S, HKTS],
	fchain chain.ChainType[S, HKTS, HKTS],
) M.Monoid[Traversal[S, A, HKTS, HKTA]] {
	return M.MakeMonoid(
		Concat[A, HKTA](fchain),
		Empty[A, HKTA](fof),
	)
}
