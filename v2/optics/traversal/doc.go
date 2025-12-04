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

/*
Package traversal provides traversals - optics for focusing on multiple values simultaneously.

# Overview

A Traversal is an optic that focuses on zero or more values within a data structure,
allowing you to view, modify, or fold over multiple elements at once. Unlike lenses
which focus on a single field, or prisms which focus on one variant, traversals can
target collections, multiple fields, or any number of values.

Traversals are the most general optic and sit at the bottom of the optics hierarchy.
They are essential for:
  - Working with collections (arrays, slices, maps)
  - Batch operations on multiple fields
  - Filtering and transforming multiple values
  - Aggregating data from multiple sources
  - Applying the same operation to all matching elements

# Mathematical Foundation

A Traversal[S, A] is defined using higher-kinded types and applicative functors.
In practical terms, it provides operations to:
  - Modify: Apply a function to all focused values
  - Set: Replace all focused values with a constant
  - FoldMap: Map each value to a monoid and combine results
  - GetAll: Collect all focused values into a list

Traversals must satisfy the traversal laws:
 1. Identity: traverse(Identity, id) == Identity
 2. Composition: traverse(Compose(F, G), f) == Compose(traverse(F, traverse(G, f)))

These laws ensure that traversals compose properly and behave consistently.

# Basic Usage

Creating a traversal for array elements:

	import (
		A "github.com/IBM/fp-go/v2/array"
		T "github.com/IBM/fp-go/v2/optics/traversal"
		TA "github.com/IBM/fp-go/v2/optics/traversal/array"
	)

	numbers := []int{1, 2, 3, 4, 5}

	// Get all elements
	all := T.GetAll(numbers)(TA.Traversal[int]())
	// Result: [1, 2, 3, 4, 5]

	// Modify all elements
	doubled := F.Pipe2(
		numbers,
		TA.Traversal[int](),
		T.Modify[[]int, int](N.Mul(2)),
	)
	// Result: [2, 4, 6, 8, 10]

	// Set all elements to a constant
	allTens := F.Pipe2(
		numbers,
		TA.Traversal[int](),
		T.Set[[]int, int](10),
	)
	// Result: [10, 10, 10, 10, 10]

# Identity Traversal

The identity traversal focuses on the entire structure:

	idTrav := T.Id[int, int]()

	value := 42
	result := T.Modify[int, int](N.Mul(2))(idTrav)(value)
	// Result: 84

# Folding with Traversals

Aggregate values using monoids:

	import (
		M "github.com/IBM/fp-go/v2/monoid"
		N "github.com/IBM/fp-go/v2/number"
	)

	numbers := []int{1, 2, 3, 4, 5}

	// Sum all elements
	sum := F.Pipe2(
		numbers,
		TA.Traversal[int](),
		T.FoldMap[int, []int, int](F.Identity[int]),
	)(N.MonoidSum[int]())
	// Result: 15

	// Product of all elements
	product := F.Pipe2(
		numbers,
		TA.Traversal[int](),
		T.FoldMap[int, []int, int](F.Identity[int]),
	)(N.MonoidProduct[int]())
	// Result: 120

# Composing Traversals

Traversals can be composed to focus on nested collections:

	type Person struct {
		Name    string
		Friends []string
	}

	people := []Person{
		{Name: "Alice", Friends: []string{"Bob", "Charlie"}},
		{Name: "Bob", Friends: []string{"Alice", "David"}},
	}

	// Traversal for people array
	peopleTrav := TA.Traversal[Person]()

	// Traversal for friends array within a person
	friendsTrav := T.MakeTraversal(func(p Person) []string {
		return p.Friends
	})

	// Compose to access all friends of all people
	allFriendsTrav := F.Pipe1(
		peopleTrav,
		T.Compose[[]Person, Person, string, ...](friendsTrav),
	)

	// Get all friends
	allFriends := T.GetAll(people)(allFriendsTrav)
	// Result: ["Bob", "Charlie", "Alice", "David"]

# Working with Records (Maps)

Traverse over map values:

	import TR "github.com/IBM/fp-go/v2/optics/traversal/record"

	scores := map[string]int{
		"Alice": 85,
		"Bob":   92,
		"Charlie": 78,
	}

	// Get all scores
	allScores := F.Pipe2(
		scores,
		TR.Traversal[string, int](),
		T.GetAll[map[string]int, int],
	)
	// Result: [85, 92, 78] (order may vary)

	// Increase all scores by 5
	boosted := F.Pipe2(
		scores,
		TR.Traversal[string, int](),
		T.Modify[map[string]int, int](func(score int) int {
			return score + 5
		}),
	)
	// Result: {"Alice": 90, "Bob": 97, "Charlie": 83}

# Working with Either Types

Traverse over the Right values:

	import (
		E "github.com/IBM/fp-go/v2/either"
		TE "github.com/IBM/fp-go/v2/optics/traversal/either"
	)

	results := []E.Either[string, int]{
		E.Right[string](10),
		E.Left[int]("error"),
		E.Right[string](20),
	}

	// Traversal for array of Either
	arrayTrav := TA.Traversal[E.Either[string, int]]()

	// Traversal for Right values
	rightTrav := TE.Traversal[string, int]()

	// Compose to access all Right values
	allRightsTrav := F.Pipe1(
		arrayTrav,
		T.Compose[[]E.Either[string, int], E.Either[string, int], int, ...](rightTrav),
	)

	// Get all Right values
	rights := T.GetAll(results)(allRightsTrav)
	// Result: [10, 20]

	// Double all Right values
	doubled := F.Pipe2(
		results,
		allRightsTrav,
		T.Modify[[]E.Either[string, int], int](N.Mul(2)),
	)
	// Result: [Right(20), Left("error"), Right(40)]

# Working with Option Types

Traverse over Some values:

	import (
		O "github.com/IBM/fp-go/v2/option"
		TO "github.com/IBM/fp-go/v2/optics/traversal/option"
	)

	values := []O.Option[int]{
		O.Some(1),
		O.None[int](),
		O.Some(2),
		O.None[int](),
		O.Some(3),
	}

	// Compose array and option traversals
	allSomesTrav := F.Pipe1(
		TA.Traversal[O.Option[int]](),
		T.Compose[[]O.Option[int], O.Option[int], int, ...](TO.Traversal[int]()),
	)

	// Get all Some values
	somes := T.GetAll(values)(allSomesTrav)
	// Result: [1, 2, 3]

	// Increment all Some values
	incremented := F.Pipe2(
		values,
		allSomesTrav,
		T.Modify[[]O.Option[int], int](N.Add(1)),
	)
	// Result: [Some(2), None, Some(3), None, Some(4)]

# Real-World Example: Nested Data Structures

	type Department struct {
		Name      string
		Employees []Employee
	}

	type Employee struct {
		Name   string
		Salary int
	}

	company := []Department{
		{
			Name: "Engineering",
			Employees: []Employee{
				{Name: "Alice", Salary: 100000},
				{Name: "Bob", Salary: 95000},
			},
		},
		{
			Name: "Sales",
			Employees: []Employee{
				{Name: "Charlie", Salary: 80000},
				{Name: "David", Salary: 85000},
			},
		},
	}

	// Traversal for departments
	deptTrav := TA.Traversal[Department]()

	// Traversal for employees within a department
	empTrav := T.MakeTraversal(func(d Department) []Employee {
		return d.Employees
	})

	// Traversal for employee array
	empArrayTrav := TA.Traversal[Employee]()

	// Compose to access all employees
	allEmpTrav := F.Pipe2(
		deptTrav,
		T.Compose[[]Department, Department, []Employee, ...](empTrav),
		T.Compose[[]Department, []Employee, Employee, ...](empArrayTrav),
	)

	// Get all employee names
	names := F.Pipe2(
		company,
		allEmpTrav,
		T.FoldMap[[]string, []Department, Employee](func(e Employee) []string {
			return []string{e.Name}
		}),
	)(A.Monoid[string]())
	// Result: ["Alice", "Bob", "Charlie", "David"]

	// Give everyone a 10% raise
	withRaises := F.Pipe2(
		company,
		allEmpTrav,
		T.Modify[[]Department, Employee](func(e Employee) Employee {
			e.Salary = int(float64(e.Salary) * 1.1)
			return e
		}),
	)

# Real-World Example: Filtering with Traversals

	type Product struct {
		Name  string
		Price float64
		InStock bool
	}

	products := []Product{
		{Name: "Laptop", Price: 999.99, InStock: true},
		{Name: "Mouse", Price: 29.99, InStock: false},
		{Name: "Keyboard", Price: 79.99, InStock: true},
	}

	// Create a traversal that only focuses on in-stock products
	inStockTrav := T.MakeTraversal(func(ps []Product) []Product {
		return A.Filter(func(p Product) bool {
			return p.InStock
		})(ps)
	})

	// Apply discount to in-stock items
	discounted := F.Pipe2(
		products,
		inStockTrav,
		T.Modify[[]Product, Product](func(p Product) Product {
			p.Price = p.Price * 0.9
			return p
		}),
	)
	// Only Laptop and Keyboard prices are reduced

# Real-World Example: Data Aggregation

	type Order struct {
		ID     string
		Items  []OrderItem
		Status string
	}

	type OrderItem struct {
		Product  string
		Quantity int
		Price    float64
	}

	orders := []Order{
		{
			ID: "001",
			Items: []OrderItem{
				{Product: "Widget", Quantity: 2, Price: 10.0},
				{Product: "Gadget", Quantity: 1, Price: 25.0},
			},
			Status: "completed",
		},
		{
			ID: "002",
			Items: []OrderItem{
				{Product: "Widget", Quantity: 5, Price: 10.0},
			},
			Status: "completed",
		},
	}

	// Traversal for orders
	orderTrav := TA.Traversal[Order]()

	// Traversal for items within an order
	itemsTrav := T.MakeTraversal(func(o Order) []OrderItem {
		return o.Items
	})

	// Traversal for item array
	itemArrayTrav := TA.Traversal[OrderItem]()

	// Compose to access all items
	allItemsTrav := F.Pipe2(
		orderTrav,
		T.Compose[[]Order, Order, []OrderItem, ...](itemsTrav),
		T.Compose[[]Order, []OrderItem, OrderItem, ...](itemArrayTrav),
	)

	// Calculate total revenue
	totalRevenue := F.Pipe2(
		orders,
		allItemsTrav,
		T.FoldMap[float64, []Order, OrderItem](func(item OrderItem) float64 {
			return float64(item.Quantity) * item.Price
		}),
	)(N.MonoidSum[float64]())
	// Result: 95.0 (2*10 + 1*25 + 5*10)

# Traversals in the Optics Hierarchy

Traversals are the most general optic:

	Iso[S, A]
	    ↓
	Lens[S, A]
	    ↓
	Optional[S, A]
	    ↓
	Traversal[S, A]

	Prism[S, A]
	    ↓
	Optional[S, A]
	    ↓
	Traversal[S, A]

This means:
  - Every Iso, Lens, Prism, and Optional can be converted to a Traversal
  - Traversals are the most flexible but least specific optic
  - Use more specific optics when possible for better type safety

# Performance Considerations

Traversals can be efficient but consider:
  - Each traversal operation may iterate over all elements
  - Composition creates nested iterations
  - FoldMap is often more efficient than GetAll followed by reduction
  - Modify creates new copies (immutability)

For best performance:
  - Use specialized traversals (array, record, etc.) when available
  - Avoid unnecessary composition
  - Consider batch operations
  - Cache composed traversals

# Type Safety

Traversals are fully type-safe:
  - Compile-time type checking
  - Generic type parameters ensure correctness
  - Composition maintains type relationships
  - No runtime type assertions

# Function Reference

Core Functions:
  - Id: Create an identity traversal
  - Modify: Apply a function to all focused values
  - Set: Replace all focused values with a constant
  - Compose: Compose two traversals

Aggregation:
  - FoldMap: Map each value to a monoid and combine
  - Fold: Fold over all values using a monoid
  - GetAll: Collect all focused values into a list

# Specialized Traversals

The package includes specialized sub-packages for common patterns:
  - array: Traversals for arrays and slices
  - record: Traversals for maps
  - either: Traversals for Either types
  - option: Traversals for Option types

Each specialized package provides optimized implementations for its data type.

# Related Packages

  - github.com/IBM/fp-go/v2/optics/lens: Lenses for single fields
  - github.com/IBM/fp-go/v2/optics/prism: Prisms for sum types
  - github.com/IBM/fp-go/v2/optics/optional: Optionals for maybe values
  - github.com/IBM/fp-go/v2/optics/traversal/array: Array traversals
  - github.com/IBM/fp-go/v2/optics/traversal/record: Record/map traversals
  - github.com/IBM/fp-go/v2/optics/traversal/either: Either traversals
  - github.com/IBM/fp-go/v2/optics/traversal/option: Option traversals
  - github.com/IBM/fp-go/v2/array: Array utilities
  - github.com/IBM/fp-go/v2/monoid: Monoid type class
*/
package traversal
