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

import (
	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// Semigroup returns a Semigroup instance for Eq[A].
// A Semigroup provides a way to combine two values of the same type.
// For Eq, the combination uses logical AND - two values are equal only if
// they are equal according to BOTH equality predicates.
//
// Type Parameters:
//   - A: The type for which equality predicates are being combined
//
// Returns:
//   - A Semigroup[Eq[A]] that combines equality predicates with logical AND
//
// The Concat operation satisfies:
//   - Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
// Example - Combine multiple equality checks:
//
//	type User struct {
//	    Username string
//	    Email    string
//	}
//
//	usernameEq := eq.Contramap(func(u User) string {
//	    return u.Username
//	})(eq.FromStrictEquals[string]())
//
//	emailEq := eq.Contramap(func(u User) string {
//	    return u.Email
//	})(eq.FromStrictEquals[string]())
//
//	// Users are equal only if BOTH username AND email match
//	userEq := eq.Semigroup[User]().Concat(usernameEq, emailEq)
//
//	u1 := User{Username: "alice", Email: "alice@example.com"}
//	u2 := User{Username: "alice", Email: "alice@example.com"}
//	u3 := User{Username: "alice", Email: "different@example.com"}
//
//	assert.True(t, userEq.Equals(u1, u2))   // Both match
//	assert.False(t, userEq.Equals(u1, u3))  // Email differs
//
// Example - Combine multiple field checks:
//
//	type Product struct {
//	    ID    int
//	    Name  string
//	    Price float64
//	}
//
//	idEq := eq.Contramap(func(p Product) int { return p.ID })(eq.FromStrictEquals[int]())
//	nameEq := eq.Contramap(func(p Product) string { return p.Name })(eq.FromStrictEquals[string]())
//	priceEq := eq.Contramap(func(p Product) float64 { return p.Price })(eq.FromStrictEquals[float64]())
//
//	sg := eq.Semigroup[Product]()
//	// All three fields must match
//	productEq := sg.Concat(sg.Concat(idEq, nameEq), priceEq)
//
// Use cases:
//   - Combining multiple field comparisons for struct equality
//   - Building complex equality predicates from simpler ones
//   - Ensuring all conditions are met (logical AND of predicates)
func Semigroup[A any]() S.Semigroup[Eq[A]] {
	return S.MakeSemigroup(func(x, y Eq[A]) Eq[A] {
		return FromEquals(func(a, b A) bool {
			return x.Equals(a, b) && y.Equals(a, b)
		})
	})
}

// Monoid returns a Monoid instance for Eq[A].
// A Monoid extends Semigroup with an identity element (Empty).
// For Eq, the identity is an equality predicate that always returns true.
//
// Type Parameters:
//   - A: The type for which the equality monoid is defined
//
// Returns:
//   - A Monoid[Eq[A]] with:
//   - Concat: Combines equality predicates with logical AND (from Semigroup)
//   - Empty: An equality predicate that always returns true (identity element)
//
// Monoid Laws:
//  1. Left Identity: Concat(Empty(), x) = x
//  2. Right Identity: Concat(x, Empty()) = x
//  3. Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
//
// Example - Using the identity element:
//
//	monoid := eq.Monoid[int]()
//	intEq := eq.FromStrictEquals[int]()
//
//	// Empty is the identity - combining with it doesn't change behavior
//	leftIdentity := monoid.Concat(monoid.Empty(), intEq)
//	rightIdentity := monoid.Concat(intEq, monoid.Empty())
//
//	assert.True(t, leftIdentity.Equals(42, 42))
//	assert.False(t, leftIdentity.Equals(42, 43))
//	assert.True(t, rightIdentity.Equals(42, 42))
//	assert.False(t, rightIdentity.Equals(42, 43))
//
// Example - Empty always returns true:
//
//	monoid := eq.Monoid[string]()
//	alwaysTrue := monoid.Empty()
//
//	assert.True(t, alwaysTrue.Equals("hello", "world"))
//	assert.True(t, alwaysTrue.Equals("same", "same"))
//	assert.True(t, alwaysTrue.Equals("", "anything"))
//
// Example - Building complex equality with fold:
//
//	type Person struct {
//	    FirstName string
//	    LastName  string
//	    Age       int
//	}
//
//	firstNameEq := eq.Contramap(func(p Person) string { return p.FirstName })(eq.FromStrictEquals[string]())
//	lastNameEq := eq.Contramap(func(p Person) string { return p.LastName })(eq.FromStrictEquals[string]())
//	ageEq := eq.Contramap(func(p Person) int { return p.Age })(eq.FromStrictEquals[int]())
//
//	monoid := eq.Monoid[Person]()
//	// Combine all predicates - all fields must match
//	personEq := monoid.Concat(monoid.Concat(firstNameEq, lastNameEq), ageEq)
//
// Use cases:
//   - Providing a neutral element for equality combinations
//   - Generic algorithms that require a Monoid instance
//   - Folding multiple equality predicates into one
//   - Default "accept everything" equality predicate
func Monoid[A any]() M.Monoid[Eq[A]] {
	return M.MakeMonoid(Semigroup[A]().Concat, Empty[A]())
}
