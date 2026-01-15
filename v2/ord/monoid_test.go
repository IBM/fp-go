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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Semigroup laws
func TestSemigroup_Associativity(t *testing.T) {
	type Person struct {
		LastName   string
		FirstName  string
		MiddleName string
	}

	stringOrd := FromStrictCompare[string]()

	byLastName := Contramap(func(p Person) string { return p.LastName })(stringOrd)
	byFirstName := Contramap(func(p Person) string { return p.FirstName })(stringOrd)
	byMiddleName := Contramap(func(p Person) string { return p.MiddleName })(stringOrd)

	sg := Semigroup[Person]()

	// Test associativity: (a <> b) <> c == a <> (b <> c)
	left := sg.Concat(sg.Concat(byLastName, byFirstName), byMiddleName)
	right := sg.Concat(byLastName, sg.Concat(byFirstName, byMiddleName))

	p1 := Person{LastName: "Smith", FirstName: "John", MiddleName: "A"}
	p2 := Person{LastName: "Smith", FirstName: "John", MiddleName: "B"}

	assert.Equal(t, left.Compare(p1, p2), right.Compare(p1, p2), "Associativity should hold")
}

// Test Semigroup with three levels
func TestSemigroup_ThreeLevels(t *testing.T) {
	type Employee struct {
		Department string
		Level      int
		Name       string
	}

	stringOrd := FromStrictCompare[string]()
	intOrd := FromStrictCompare[int]()

	byDept := Contramap(func(e Employee) string { return e.Department })(stringOrd)
	byLevel := Contramap(func(e Employee) int { return e.Level })(intOrd)
	byName := Contramap(func(e Employee) string { return e.Name })(stringOrd)

	sg := Semigroup[Employee]()
	employeeOrd := sg.Concat(sg.Concat(byDept, byLevel), byName)

	e1 := Employee{Department: "IT", Level: 3, Name: "Alice"}
	e2 := Employee{Department: "IT", Level: 3, Name: "Bob"}
	e3 := Employee{Department: "IT", Level: 2, Name: "Charlie"}
	e4 := Employee{Department: "HR", Level: 3, Name: "David"}

	// Same dept, same level, different name
	assert.Equal(t, -1, employeeOrd.Compare(e1, e2), "Alice < Bob")

	// Same dept, different level
	assert.Equal(t, 1, employeeOrd.Compare(e1, e3), "Level 3 > Level 2")

	// Different dept
	assert.Equal(t, -1, employeeOrd.Compare(e4, e1), "HR < IT")
}

// Test Monoid identity laws
func TestMonoid_IdentityLaws(t *testing.T) {
	m := Monoid[int]()
	intOrd := FromStrictCompare[int]()
	emptyOrd := m.Empty()

	// Left identity: empty <> x == x
	leftIdentity := m.Concat(emptyOrd, intOrd)
	assert.Equal(t, -1, leftIdentity.Compare(3, 5), "Left identity: 3 < 5")
	assert.Equal(t, 1, leftIdentity.Compare(5, 3), "Left identity: 5 > 3")

	// Right identity: x <> empty == x
	rightIdentity := m.Concat(intOrd, emptyOrd)
	assert.Equal(t, -1, rightIdentity.Compare(3, 5), "Right identity: 3 < 5")
	assert.Equal(t, 1, rightIdentity.Compare(5, 3), "Right identity: 5 > 3")
}

// Test Monoid with multiple empty concatenations
func TestMonoid_MultipleEmpty(t *testing.T) {
	m := Monoid[int]()
	emptyOrd := m.Empty()

	// Concatenating multiple empty orderings should still be empty
	combined := m.Concat(m.Concat(emptyOrd, emptyOrd), emptyOrd)

	assert.Equal(t, 0, combined.Compare(5, 3), "Multiple empties: always equal")
	assert.Equal(t, 0, combined.Compare(3, 5), "Multiple empties: always equal")
	assert.True(t, combined.Equals(5, 3), "Multiple empties: always equal")
}

// Test MaxSemigroup with edge cases
func TestMaxSemigroup_EdgeCases(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	maxSg := MaxSemigroup(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"both positive", 5, 3, 5},
		{"both negative", -5, -3, -3},
		{"mixed signs", -5, 3, 3},
		{"zero and positive", 0, 5, 5},
		{"zero and negative", 0, -5, 0},
		{"both zero", 0, 0, 0},
		{"equal positive", 5, 5, 5},
		{"equal negative", -5, -5, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MinSemigroup with edge cases
func TestMinSemigroup_EdgeCases(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	minSg := MinSemigroup(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"both positive", 5, 3, 3},
		{"both negative", -5, -3, -5},
		{"mixed signs", -5, 3, -5},
		{"zero and positive", 0, 5, 0},
		{"zero and negative", 0, -5, -5},
		{"both zero", 0, 0, 0},
		{"equal positive", 5, 5, 5},
		{"equal negative", -5, -5, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MaxSemigroup with strings
func TestMaxSemigroup_Strings(t *testing.T) {
	stringOrd := FromStrictCompare[string]()
	maxSg := MaxSemigroup(stringOrd)

	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"alphabetical", "apple", "banana", "banana"},
		{"same string", "apple", "apple", "apple"},
		{"empty and non-empty", "", "apple", "apple"},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MinSemigroup with strings
func TestMinSemigroup_Strings(t *testing.T) {
	stringOrd := FromStrictCompare[string]()
	minSg := MinSemigroup(stringOrd)

	tests := []struct {
		name     string
		a        string
		b        string
		expected string
	}{
		{"alphabetical", "apple", "banana", "apple"},
		{"same string", "apple", "apple", "apple"},
		{"empty and non-empty", "", "apple", ""},
		{"both empty", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MaxSemigroup associativity
func TestMaxSemigroup_Associativity(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	maxSg := MaxSemigroup(intOrd)

	// (a <> b) <> c == a <> (b <> c)
	a, b, c := 5, 3, 7

	left := maxSg.Concat(maxSg.Concat(a, b), c)
	right := maxSg.Concat(a, maxSg.Concat(b, c))

	assert.Equal(t, left, right, "MaxSemigroup should be associative")
	assert.Equal(t, 7, left, "Should return maximum value")
}

// Test MinSemigroup associativity
func TestMinSemigroup_Associativity(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	minSg := MinSemigroup(intOrd)

	// (a <> b) <> c == a <> (b <> c)
	a, b, c := 5, 3, 7

	left := minSg.Concat(minSg.Concat(a, b), c)
	right := minSg.Concat(a, minSg.Concat(b, c))

	assert.Equal(t, left, right, "MinSemigroup should be associative")
	assert.Equal(t, 3, left, "Should return minimum value")
}

// Test Semigroup with reversed ordering
func TestSemigroup_WithReverse(t *testing.T) {
	type Person struct {
		Age  int
		Name string
	}

	intOrd := FromStrictCompare[int]()
	stringOrd := FromStrictCompare[string]()

	// Order by age descending, then by name ascending
	byAge := Contramap(func(p Person) int { return p.Age })(Reverse(intOrd))
	byName := Contramap(func(p Person) string { return p.Name })(stringOrd)

	sg := Semigroup[Person]()
	personOrd := sg.Concat(byAge, byName)

	p1 := Person{Age: 30, Name: "Alice"}
	p2 := Person{Age: 30, Name: "Bob"}
	p3 := Person{Age: 25, Name: "Charlie"}

	// Same age, different name
	assert.Equal(t, -1, personOrd.Compare(p1, p2), "Alice < Bob (same age)")

	// Different age (descending)
	assert.Equal(t, -1, personOrd.Compare(p1, p3), "30 > 25 (descending)")
}

// Benchmark MaxSemigroup
func BenchmarkMaxSemigroup(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	maxSg := MaxSemigroup(intOrd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = maxSg.Concat(i, i+1)
	}
}

// Benchmark MinSemigroup
func BenchmarkMinSemigroup(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	minSg := MinSemigroup(intOrd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = minSg.Concat(i, i+1)
	}
}

// Benchmark Semigroup concatenation
func BenchmarkSemigroup_Concat(b *testing.B) {
	type Person struct {
		LastName  string
		FirstName string
	}

	stringOrd := FromStrictCompare[string]()
	byLastName := Contramap(func(p Person) string { return p.LastName })(stringOrd)
	byFirstName := Contramap(func(p Person) string { return p.FirstName })(stringOrd)

	sg := Semigroup[Person]()
	personOrd := sg.Concat(byLastName, byFirstName)

	p1 := Person{LastName: "Smith", FirstName: "Alice"}
	p2 := Person{LastName: "Smith", FirstName: "Bob"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = personOrd.Compare(p1, p2)
	}
}
