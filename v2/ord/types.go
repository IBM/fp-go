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

package ord

type (
	// Kleisli represents a function that takes a value of type A and returns an Ord[B].
	// This is useful for creating orderings that depend on input values.
	//
	// Type Parameters:
	//   - A: The input type
	//   - B: The type for which ordering is produced
	//
	// Example:
	//
	//	// Create a Kleisli that produces different orderings based on input
	//	var orderingFactory Kleisli[string, int] = func(mode string) Ord[int] {
	//	    if mode == "ascending" {
	//	        return ord.FromStrictCompare[int]()
	//	    }
	//	    return ord.Reverse(ord.FromStrictCompare[int]())
	//	}
	//	ascOrd := orderingFactory("ascending")
	//	descOrd := orderingFactory("descending")
	Kleisli[A, B any] = func(A) Ord[B]

	// Operator represents a function that transforms an Ord[A] into a value of type B.
	// This is commonly used for operations that modify or combine orderings.
	//
	// Type Parameters:
	//   - A: The type for which ordering is defined
	//   - B: The result type of the operation
	//
	// This is equivalent to Kleisli[Ord[A], B] and is used for operations like
	// Contramap, which takes an Ord[A] and produces an Ord[B].
	//
	// Example:
	//
	//	// Contramap is an Operator that transforms Ord[A] to Ord[B]
	//	type Person struct { Age int }
	//	var ageOperator Operator[int, Person] = ord.Contramap(func(p Person) int {
	//	    return p.Age
	//	})
	//	intOrd := ord.FromStrictCompare[int]()
	//	personOrd := ageOperator(intOrd)
	Operator[A, B any] = Kleisli[Ord[A], B]
)
