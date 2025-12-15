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
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromStrictEquals(t *testing.T) {
	t.Run("int equality", func(t *testing.T) {
		intEq := FromStrictEquals[int]()

		assert.True(t, intEq.Equals(42, 42))
		assert.True(t, intEq.Equals(0, 0))
		assert.True(t, intEq.Equals(-1, -1))

		assert.False(t, intEq.Equals(42, 43))
		assert.False(t, intEq.Equals(0, 1))
		assert.False(t, intEq.Equals(-1, 1))
	})

	t.Run("string equality", func(t *testing.T) {
		stringEq := FromStrictEquals[string]()

		assert.True(t, stringEq.Equals("hello", "hello"))
		assert.True(t, stringEq.Equals("", ""))
		assert.True(t, stringEq.Equals("test", "test"))

		assert.False(t, stringEq.Equals("hello", "Hello"))
		assert.False(t, stringEq.Equals("hello", "world"))
		assert.False(t, stringEq.Equals("", "test"))
	})

	t.Run("float equality", func(t *testing.T) {
		floatEq := FromStrictEquals[float64]()

		assert.True(t, floatEq.Equals(1.0, 1.0))
		assert.True(t, floatEq.Equals(0.0, 0.0))
		assert.True(t, floatEq.Equals(-1.5, -1.5))

		assert.False(t, floatEq.Equals(1.0, 1.1))
		assert.False(t, floatEq.Equals(0.0, 0.1))
	})

	t.Run("bool equality", func(t *testing.T) {
		boolEq := FromStrictEquals[bool]()

		assert.True(t, boolEq.Equals(true, true))
		assert.True(t, boolEq.Equals(false, false))

		assert.False(t, boolEq.Equals(true, false))
		assert.False(t, boolEq.Equals(false, true))
	})
}

func TestFromEquals(t *testing.T) {
	t.Run("case-insensitive string equality", func(t *testing.T) {
		caseInsensitiveEq := FromEquals(strings.EqualFold)

		assert.True(t, caseInsensitiveEq.Equals("hello", "HELLO"))
		assert.True(t, caseInsensitiveEq.Equals("Hello", "hello"))
		assert.True(t, caseInsensitiveEq.Equals("TeSt", "test"))

		assert.False(t, caseInsensitiveEq.Equals("hello", "world"))
	})

	t.Run("approximate float equality", func(t *testing.T) {
		epsilon := 0.0001
		approxEq := FromEquals(func(a, b float64) bool {
			return math.Abs(a-b) < epsilon
		})

		assert.True(t, approxEq.Equals(1.0, 1.0))
		assert.True(t, approxEq.Equals(1.0, 1.00009))
		assert.True(t, approxEq.Equals(1.00009, 1.0))

		assert.False(t, approxEq.Equals(1.0, 1.001))
		assert.False(t, approxEq.Equals(1.0, 2.0))
	})

	t.Run("custom struct equality", func(t *testing.T) {
		type Point struct {
			X, Y int
		}

		pointEq := FromEquals(func(a, b Point) bool {
			return a.X == b.X && a.Y == b.Y
		})

		assert.True(t, pointEq.Equals(Point{1, 2}, Point{1, 2}))
		assert.True(t, pointEq.Equals(Point{0, 0}, Point{0, 0}))

		assert.False(t, pointEq.Equals(Point{1, 2}, Point{2, 1}))
		assert.False(t, pointEq.Equals(Point{1, 2}, Point{1, 3}))
	})

	t.Run("slice equality", func(t *testing.T) {
		sliceEq := FromEquals(func(a, b []int) bool {
			if len(a) != len(b) {
				return false
			}
			for i := range a {
				if a[i] != b[i] {
					return false
				}
			}
			return true
		})

		assert.True(t, sliceEq.Equals([]int{1, 2, 3}, []int{1, 2, 3}))
		assert.True(t, sliceEq.Equals([]int{}, []int{}))
		assert.True(t, sliceEq.Equals(nil, nil))

		assert.False(t, sliceEq.Equals([]int{1, 2, 3}, []int{1, 2, 4}))
		assert.False(t, sliceEq.Equals([]int{1, 2}, []int{1, 2, 3}))
		assert.False(t, sliceEq.Equals([]int{1, 2, 3}, []int{}))
	})
}

func TestEmpty(t *testing.T) {
	t.Run("always returns true", func(t *testing.T) {
		emptyInt := Empty[int]()

		assert.True(t, emptyInt.Equals(1, 2))
		assert.True(t, emptyInt.Equals(42, 43))
		assert.True(t, emptyInt.Equals(0, 100))
		assert.True(t, emptyInt.Equals(-1, 1))
	})

	t.Run("works with any type", func(t *testing.T) {
		emptyString := Empty[string]()
		assert.True(t, emptyString.Equals("hello", "world"))
		assert.True(t, emptyString.Equals("", "test"))

		emptyBool := Empty[bool]()
		assert.True(t, emptyBool.Equals(true, false))

		type Custom struct{ Value int }
		emptyCustom := Empty[Custom]()
		assert.True(t, emptyCustom.Equals(Custom{1}, Custom{2}))
	})
}

func TestEquals(t *testing.T) {
	t.Run("curried equality for int", func(t *testing.T) {
		intEq := FromStrictEquals[int]()
		equals42 := Equals(intEq)(42)

		assert.True(t, equals42(42))
		assert.False(t, equals42(43))
		assert.False(t, equals42(0))
		assert.False(t, equals42(-42))
	})

	t.Run("curried equality for string", func(t *testing.T) {
		stringEq := FromStrictEquals[string]()
		equalsHello := Equals(stringEq)("hello")

		assert.True(t, equalsHello("hello"))
		assert.False(t, equalsHello("Hello"))
		assert.False(t, equalsHello("world"))
		assert.False(t, equalsHello(""))
	})

	t.Run("partial application", func(t *testing.T) {
		intEq := FromStrictEquals[int]()
		equalsFunc := Equals(intEq)

		equals10 := equalsFunc(10)
		equals20 := equalsFunc(20)

		assert.True(t, equals10(10))
		assert.False(t, equals10(20))

		assert.True(t, equals20(20))
		assert.False(t, equals20(10))
	})
}

func TestContramap(t *testing.T) {
	type Person struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("compare by ID", func(t *testing.T) {
		personEqByID := Contramap(func(p Person) int {
			return p.ID
		})(FromStrictEquals[int]())

		p1 := Person{ID: 1, Name: "Alice", Age: 30}
		p2 := Person{ID: 1, Name: "Bob", Age: 25}
		p3 := Person{ID: 2, Name: "Alice", Age: 30}

		assert.True(t, personEqByID.Equals(p1, p2))  // Same ID
		assert.False(t, personEqByID.Equals(p1, p3)) // Different ID
	})

	t.Run("compare by name", func(t *testing.T) {
		personEqByName := Contramap(func(p Person) string {
			return p.Name
		})(FromStrictEquals[string]())

		p1 := Person{ID: 1, Name: "Alice", Age: 30}
		p2 := Person{ID: 2, Name: "Alice", Age: 25}
		p3 := Person{ID: 1, Name: "Bob", Age: 30}

		assert.True(t, personEqByName.Equals(p1, p2))  // Same name
		assert.False(t, personEqByName.Equals(p1, p3)) // Different name
	})

	t.Run("compare by age", func(t *testing.T) {
		personEqByAge := Contramap(func(p Person) int {
			return p.Age
		})(FromStrictEquals[int]())

		p1 := Person{ID: 1, Name: "Alice", Age: 30}
		p2 := Person{ID: 2, Name: "Bob", Age: 30}
		p3 := Person{ID: 1, Name: "Alice", Age: 25}

		assert.True(t, personEqByAge.Equals(p1, p2))  // Same age
		assert.False(t, personEqByAge.Equals(p1, p3)) // Different age
	})

	t.Run("case-insensitive name comparison", func(t *testing.T) {
		caseInsensitiveEq := FromEquals(strings.EqualFold)

		personEqByNameCI := Contramap(func(p Person) string {
			return p.Name
		})(caseInsensitiveEq)

		p1 := Person{ID: 1, Name: "Alice", Age: 30}
		p2 := Person{ID: 2, Name: "ALICE", Age: 25}
		p3 := Person{ID: 1, Name: "Bob", Age: 30}

		assert.True(t, personEqByNameCI.Equals(p1, p2))  // Same name (case-insensitive)
		assert.False(t, personEqByNameCI.Equals(p1, p3)) // Different name
	})

	t.Run("nested contramap", func(t *testing.T) {
		type Address struct {
			Street string
			City   string
		}

		type User struct {
			Name    string
			Address Address
		}

		// Compare users by city
		userEqByCity := Contramap(func(u User) string {
			return u.Address.City
		})(FromStrictEquals[string]())

		u1 := User{Name: "Alice", Address: Address{Street: "Main St", City: "NYC"}}
		u2 := User{Name: "Bob", Address: Address{Street: "Oak Ave", City: "NYC"}}
		u3 := User{Name: "Charlie", Address: Address{Street: "Main St", City: "LA"}}

		assert.True(t, userEqByCity.Equals(u1, u2))  // Same city
		assert.False(t, userEqByCity.Equals(u1, u3)) // Different city
	})
}

func TestSemigroup(t *testing.T) {
	type User struct {
		Username string
		Email    string
	}

	t.Run("combine two equality predicates", func(t *testing.T) {
		usernameEq := Contramap(func(u User) string {
			return u.Username
		})(FromStrictEquals[string]())

		emailEq := Contramap(func(u User) string {
			return u.Email
		})(FromStrictEquals[string]())

		// Both username AND email must match
		userEq := Semigroup[User]().Concat(usernameEq, emailEq)

		u1 := User{Username: "alice", Email: "alice@example.com"}
		u2 := User{Username: "alice", Email: "alice@example.com"}
		u3 := User{Username: "alice", Email: "different@example.com"}
		u4 := User{Username: "bob", Email: "alice@example.com"}

		assert.True(t, userEq.Equals(u1, u2))  // Both match
		assert.False(t, userEq.Equals(u1, u3)) // Email differs
		assert.False(t, userEq.Equals(u1, u4)) // Username differs
	})

	t.Run("combine multiple equality predicates", func(t *testing.T) {
		type Product struct {
			ID    int
			Name  string
			Price float64
		}

		idEq := Contramap(func(p Product) int {
			return p.ID
		})(FromStrictEquals[int]())

		nameEq := Contramap(func(p Product) string {
			return p.Name
		})(FromStrictEquals[string]())

		priceEq := Contramap(func(p Product) float64 {
			return p.Price
		})(FromStrictEquals[float64]())

		sg := Semigroup[Product]()
		productEq := sg.Concat(sg.Concat(idEq, nameEq), priceEq)

		p1 := Product{ID: 1, Name: "Widget", Price: 9.99}
		p2 := Product{ID: 1, Name: "Widget", Price: 9.99}
		p3 := Product{ID: 1, Name: "Widget", Price: 10.99}

		assert.True(t, productEq.Equals(p1, p2))  // All match
		assert.False(t, productEq.Equals(p1, p3)) // Price differs
	})

	t.Run("associativity", func(t *testing.T) {
		eq1 := FromStrictEquals[int]()
		eq2 := FromStrictEquals[int]()
		eq3 := FromStrictEquals[int]()

		sg := Semigroup[int]()

		// (eq1 <> eq2) <> eq3
		left := sg.Concat(sg.Concat(eq1, eq2), eq3)

		// eq1 <> (eq2 <> eq3)
		right := sg.Concat(eq1, sg.Concat(eq2, eq3))

		// Both should behave the same
		assert.True(t, left.Equals(42, 42))
		assert.True(t, right.Equals(42, 42))
		assert.False(t, left.Equals(42, 43))
		assert.False(t, right.Equals(42, 43))
	})
}

func TestMonoid(t *testing.T) {
	t.Run("empty is identity", func(t *testing.T) {
		intEq := FromStrictEquals[int]()
		monoid := Monoid[int]()
		empty := monoid.Empty()

		// empty <> eq = eq
		leftIdentity := monoid.Concat(empty, intEq)
		assert.True(t, leftIdentity.Equals(42, 42))
		assert.False(t, leftIdentity.Equals(42, 43))

		// eq <> empty = eq
		rightIdentity := monoid.Concat(intEq, empty)
		assert.True(t, rightIdentity.Equals(42, 42))
		assert.False(t, rightIdentity.Equals(42, 43))
	})

	t.Run("empty always returns true", func(t *testing.T) {
		monoid := Monoid[string]()
		empty := monoid.Empty()

		assert.True(t, empty.Equals("hello", "world"))
		assert.True(t, empty.Equals("", "test"))
		assert.True(t, empty.Equals("same", "same"))
	})

	t.Run("concat with empty", func(t *testing.T) {
		stringEq := FromStrictEquals[string]()
		monoid := Monoid[string]()

		// Combining with empty should preserve the original behavior
		combined := monoid.Concat(stringEq, monoid.Empty())

		assert.True(t, combined.Equals("hello", "hello"))
		assert.False(t, combined.Equals("hello", "world"))
	})
}

// Test type class laws
func TestEqLaws(t *testing.T) {
	intEq := FromStrictEquals[int]()

	t.Run("reflexivity", func(t *testing.T) {
		// For all x, Equals(x, x) = true
		values := []int{0, 1, -1, 42, 100, -100}
		for _, x := range values {
			assert.True(t, intEq.Equals(x, x), "reflexivity failed for %d", x)
		}
	})

	t.Run("symmetry", func(t *testing.T) {
		// For all x, y, Equals(x, y) = Equals(y, x)
		pairs := [][2]int{
			{1, 1}, {1, 2}, {42, 42}, {0, 1}, {-1, 1},
		}
		for _, pair := range pairs {
			x, y := pair[0], pair[1]
			assert.Equal(t, intEq.Equals(x, y), intEq.Equals(y, x),
				"symmetry failed for (%d, %d)", x, y)
		}
	})

	t.Run("transitivity", func(t *testing.T) {
		// If Equals(x, y) and Equals(y, z), then Equals(x, z)
		triples := [][3]int{
			{1, 1, 1},
			{42, 42, 42},
			{0, 0, 0},
		}
		for _, triple := range triples {
			x, y, z := triple[0], triple[1], triple[2]
			if intEq.Equals(x, y) && intEq.Equals(y, z) {
				assert.True(t, intEq.Equals(x, z),
					"transitivity failed for (%d, %d, %d)", x, y, z)
			}
		}
	})
}
