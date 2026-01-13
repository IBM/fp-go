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

package eq

// Contramap creates an Eq[B] from an Eq[A] by providing a function that maps B to A.
// This is a contravariant functor operation that allows you to transform equality predicates
// by mapping the input type. It's particularly useful for comparing complex types by
// extracting comparable fields.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// The name "contramap" comes from category theory, where it represents a contravariant
// functor. Unlike regular map (covariant), which transforms the output, contramap
// transforms the input in the opposite direction.
//
// Type Parameters:
//   - A: The type that has an existing Eq instance
//   - B: The type for which we want to create an Eq instance
//
// Parameters:
//   - f: A function that extracts or converts a value of type B to type A
//
// Returns:
//   - A function that takes an Eq[A] and returns an Eq[B]
//
// The resulting Eq[B] compares two B values by:
//  1. Applying f to both values to get A values
//  2. Using the original Eq[A] to compare those A values
//
// Example - Compare structs by a single field:
//
//	type Person struct {
//	    ID   int
//	    Name string
//	    Age  int
//	}
//
//	// Compare persons by ID only
//	personEqByID := eq.Contramap(func(p Person) int {
//	    return p.ID
//	})(eq.FromStrictEquals[int]())
//
//	p1 := Person{ID: 1, Name: "Alice", Age: 30}
//	p2 := Person{ID: 1, Name: "Bob", Age: 25}
//	assert.True(t, personEqByID.Equals(p1, p2))  // Same ID, different names
//
// Example - Case-insensitive string comparison:
//
//	type User struct {
//	    Username string
//	    Email    string
//	}
//
//	caseInsensitiveEq := eq.FromEquals(func(a, b string) bool {
//	    return strings.EqualFold(a, b)
//	})
//
//	userEqByUsername := eq.Contramap(func(u User) string {
//	    return u.Username
//	})(caseInsensitiveEq)
//
//	u1 := User{Username: "Alice", Email: "alice@example.com"}
//	u2 := User{Username: "ALICE", Email: "different@example.com"}
//	assert.True(t, userEqByUsername.Equals(u1, u2))  // Case-insensitive match
//
// Example - Nested field access:
//
//	type Address struct {
//	    City string
//	}
//
//	type Person struct {
//	    Name    string
//	    Address Address
//	}
//
//	// Compare persons by city
//	personEqByCity := eq.Contramap(func(p Person) string {
//	    return p.Address.City
//	})(eq.FromStrictEquals[string]())
//
// Contramap Law:
// Contramap must satisfy: Contramap(f)(Contramap(g)(eq)) = Contramap(g ∘ f)(eq)
// This means contramapping twice is the same as contramapping with the composed function.
func Contramap[A, B any](f func(b B) A) func(Eq[A]) Eq[B] {
	return func(fa Eq[A]) Eq[B] {
		equals := fa.Equals
		return FromEquals(func(x, y B) bool {
			return equals(f(x), f(y))
		})
	}
}
