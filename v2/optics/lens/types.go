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

// Package lens provides functional optics for zooming into and updating parts of immutable data structures.
//
// A lens is a composable, first-class reference to a subpart of a data structure that enables
// getting and setting values in a purely functional way without mutation.
package lens

import (
	"github.com/IBM/fp-go/v2/endomorphism"
)

type (
	// Endomorphism is a function from a type to itself (A → A).
	// It represents transformations that preserve the type.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lens is a functional reference to a subpart of a data structure.
	//
	// A Lens[S, A] provides a composable way to focus on a field of type A within
	// a structure of type S. It consists of two operations:
	//   - Get: Extracts the focused value from the structure (S → A)
	//   - Set: Updates the focused value in the structure, returning a new structure (A → S → S)
	//
	// Lenses maintain immutability by always returning new copies of the structure
	// when setting values, never modifying the original.
	//
	// Type Parameters:
	//   - S: The source/structure type (the whole)
	//   - A: The focus/field type (the part)
	//
	// Lens Laws:
	//
	// A well-behaved lens must satisfy three laws:
	//
	// 1. GetSet (You get what you set):
	//    lens.Set(lens.Get(s))(s) == s
	//
	// 2. SetGet (You set what you get):
	//    lens.Get(lens.Set(a)(s)) == a
	//
	// 3. SetSet (Setting twice is the same as setting once):
	//    lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	//
	// Example Usage:
	//
	//	type Person struct {
	//	    Name string
	//	    Age  int
	//	}
	//
	//	// Create a lens focusing on the Name field
	//	nameLens := lens.MakeLens(
	//	    func(p Person) string { return p.Name },
	//	    func(name string) func(Person) Person {
	//	        return func(p Person) Person {
	//	            return Person{Name: name, Age: p.Age}
	//	        }
	//	    },
	//	)
	//
	//	person := Person{Name: "Alice", Age: 30}
	//	name := nameLens.Get(person)           // Returns: "Alice"
	//	updated := nameLens.Set("Bob")(person) // Returns: Person{Name: "Bob", Age: 30}
	//	// Original person remains unchanged (immutability preserved)
	Lens[S, A any] struct {
		// Get extracts the focused value of type A from structure S.
		Get func(s S) A

		// Set returns a function that updates the focused value in structure S.
		// The returned function takes a structure S and returns a new structure S
		// with the focused value updated to a. The original structure is never modified.
		Set func(a A) Endomorphism[S]

		// name is an end user facing identifier for the lens
		name string
	}

	// Kleisli represents a function that takes a value of type A and returns a Lens[S, B].
	// This is useful for composing lenses in a monadic style, allowing for dynamic lens creation
	// based on input values.
	//
	// Type Parameters:
	//   - S: The source/structure type
	//   - A: The input type
	//   - B: The focus type of the resulting lens
	Kleisli[S, A, B any] = func(A) Lens[S, B]

	// Operator is a specialized Kleisli that takes a Lens[S, A] and returns a Lens[S, B].
	// This enables lens transformations and compositions where one lens is used to derive another.
	//
	// Type Parameters:
	//   - S: The source/structure type
	//   - A: The focus type of the input lens
	//   - B: The focus type of the resulting lens
	Operator[S, A, B any] = Kleisli[S, Lens[S, A], B]
)
