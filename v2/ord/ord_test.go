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

// Test MakeOrd
func TestMakeOrd(t *testing.T) {
	intOrd := MakeOrd(
		func(a, b int) int {
			if a < b {
				return -1
			} else if a > b {
				return 1
			}
			return 0
		},
		func(a, b int) bool {
			return a == b
		},
	)

	assert.Equal(t, -1, intOrd.Compare(3, 5))
	assert.Equal(t, 1, intOrd.Compare(5, 3))
	assert.Equal(t, 0, intOrd.Compare(5, 5))

	assert.True(t, intOrd.Equals(5, 5))
	assert.False(t, intOrd.Equals(5, 3))
}

// Test FromCompare
func TestFromCompare(t *testing.T) {
	intOrd := FromCompare(func(a, b int) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}
		return 0
	})

	// Test compare
	assert.Equal(t, -1, intOrd.Compare(3, 5))
	assert.Equal(t, 1, intOrd.Compare(5, 3))
	assert.Equal(t, 0, intOrd.Compare(5, 5))

	// Test equals (derived from compare)
	assert.True(t, intOrd.Equals(5, 5))
	assert.False(t, intOrd.Equals(5, 3))
}

// Test FromStrictCompare
func TestFromStrictCompare(t *testing.T) {
	intOrd := FromStrictCompare[int]()

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"less than", 3, 5, -1},
		{"greater than", 5, 3, 1},
		{"equal", 5, 5, 0},
		{"negative numbers", -5, -3, -1},
		{"mixed signs", -3, 5, -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := intOrd.Compare(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}

	// Test equals
	assert.True(t, intOrd.Equals(5, 5))
	assert.False(t, intOrd.Equals(5, 3))
}

// Test FromStrictCompare with strings
func TestFromStrictCompare_String(t *testing.T) {
	stringOrd := FromStrictCompare[string]()

	assert.Equal(t, -1, stringOrd.Compare("apple", "banana"))
	assert.Equal(t, 1, stringOrd.Compare("banana", "apple"))
	assert.Equal(t, 0, stringOrd.Compare("apple", "apple"))

	assert.True(t, stringOrd.Equals("apple", "apple"))
	assert.False(t, stringOrd.Equals("apple", "banana"))
}

// Test FromStrictCompare with floats
func TestFromStrictCompare_Float(t *testing.T) {
	floatOrd := FromStrictCompare[float64]()

	assert.Equal(t, -1, floatOrd.Compare(3.14, 3.15))
	assert.Equal(t, 1, floatOrd.Compare(3.15, 3.14))
	assert.Equal(t, 0, floatOrd.Compare(3.14, 3.14))
}

// Test Reverse
func TestReverse(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	reversedOrd := Reverse(intOrd)

	// Original order
	assert.Equal(t, -1, intOrd.Compare(3, 5))
	assert.Equal(t, 1, intOrd.Compare(5, 3))

	// Reversed order
	assert.Equal(t, 1, reversedOrd.Compare(3, 5))
	assert.Equal(t, -1, reversedOrd.Compare(5, 3))
	assert.Equal(t, 0, reversedOrd.Compare(5, 5))

	// Equals should be the same
	assert.True(t, reversedOrd.Equals(5, 5))
	assert.False(t, reversedOrd.Equals(5, 3))
}

// Test Contramap
func TestContramap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	intOrd := FromStrictCompare[int]()

	// Order persons by age
	personOrd := Contramap(func(p Person) int {
		return p.Age
	})(intOrd)

	p1 := Person{Name: "Alice", Age: 30}
	p2 := Person{Name: "Bob", Age: 25}
	p3 := Person{Name: "Charlie", Age: 30}

	assert.Equal(t, 1, personOrd.Compare(p1, p2))  // 30 > 25
	assert.Equal(t, -1, personOrd.Compare(p2, p1)) // 25 < 30
	assert.Equal(t, 0, personOrd.Compare(p1, p3))  // 30 == 30

	assert.True(t, personOrd.Equals(p1, p3))
	assert.False(t, personOrd.Equals(p1, p2))
}

// Test ToEq
func TestToEq(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	intEq := ToEq(intOrd)

	assert.True(t, intEq.Equals(5, 5))
	assert.False(t, intEq.Equals(5, 3))
}

// Test Min
func TestMin(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	min := Min(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a smaller", 3, 5, 3},
		{"b smaller", 5, 3, 3},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -5},
		{"mixed signs", -5, 3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := min(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Max
func TestMax(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	max := Max(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a larger", 5, 3, 5},
		{"b larger", 3, 5, 5},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -3},
		{"mixed signs", -5, 3, 3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := max(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Clamp
func TestClamp(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	clamp := Clamp(intOrd)(0, 100)

	tests := []struct {
		name     string
		input    int
		expected int
	}{
		{"below minimum", -10, 0},
		{"at minimum", 0, 0},
		{"within range", 50, 50},
		{"at maximum", 100, 100},
		{"above maximum", 150, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clamp(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Lt (less than)
func TestLt(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	isLessThan5 := Lt(intOrd)(5)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"less than", 3, true},
		{"equal", 5, false},
		{"greater than", 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isLessThan5(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Leq (less than or equal)
func TestLeq(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	isAtMost5 := Leq(intOrd)(5)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"less than", 3, true},
		{"equal", 5, true},
		{"greater than", 7, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAtMost5(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Gt (greater than)
func TestGt(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	isGreaterThan5 := Gt(intOrd)(5)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"less than", 3, false},
		{"equal", 5, false},
		{"greater than", 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isGreaterThan5(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Geq (greater than or equal)
func TestGeq(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	isAtLeast5 := Geq(intOrd)(5)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"less than", 3, false},
		{"equal", 5, true},
		{"greater than", 7, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isAtLeast5(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Between
func TestBetween(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	isBetween3And7 := Between(intOrd)(3, 7)

	tests := []struct {
		name     string
		input    int
		expected bool
	}{
		{"below range", 2, false},
		{"at lower bound", 3, true},
		{"within range", 5, true},
		{"at upper bound", 7, false},
		{"above range", 8, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isBetween3And7(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Semigroup
func TestSemigroup(t *testing.T) {
	type Person struct {
		LastName  string
		FirstName string
	}

	stringOrd := FromStrictCompare[string]()

	// Order by last name
	byLastName := Contramap(func(p Person) string {
		return p.LastName
	})(stringOrd)

	// Order by first name
	byFirstName := Contramap(func(p Person) string {
		return p.FirstName
	})(stringOrd)

	// Combine: order by last name, then first name
	sg := Semigroup[Person]()
	personOrd := sg.Concat(byLastName, byFirstName)

	p1 := Person{LastName: "Smith", FirstName: "Alice"}
	p2 := Person{LastName: "Smith", FirstName: "Bob"}
	p3 := Person{LastName: "Jones", FirstName: "Charlie"}

	// Same last name, different first name
	assert.Equal(t, -1, personOrd.Compare(p1, p2)) // Alice < Bob

	// Different last names
	assert.Equal(t, 1, personOrd.Compare(p1, p3)) // Smith > Jones
}

// Test Monoid
func TestMonoid(t *testing.T) {
	m := Monoid[int]()

	// Empty ordering considers everything equal
	emptyOrd := m.Empty()
	assert.Equal(t, 0, emptyOrd.Compare(5, 3))
	assert.Equal(t, 0, emptyOrd.Compare(3, 5))
	assert.True(t, emptyOrd.Equals(5, 3))

	// Concat with empty returns the original
	intOrd := FromStrictCompare[int]()
	combined := m.Concat(intOrd, emptyOrd)

	assert.Equal(t, -1, combined.Compare(3, 5))
	assert.Equal(t, 1, combined.Compare(5, 3))
	assert.Equal(t, 0, combined.Compare(5, 5))
}

// Test MaxSemigroup
func TestMaxSemigroup(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	maxSg := MaxSemigroup(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a larger", 5, 3, 5},
		{"b larger", 3, 5, 5},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -3},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maxSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test MinSemigroup
func TestMinSemigroup(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	minSg := MinSemigroup(intOrd)

	tests := []struct {
		name     string
		a        int
		b        int
		expected int
	}{
		{"a smaller", 3, 5, 3},
		{"b smaller", 5, 3, 3},
		{"equal", 5, 5, 5},
		{"negative", -5, -3, -5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := minSg.Concat(tt.a, tt.b)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Ord laws
func TestOrdLaws(t *testing.T) {
	intOrd := FromStrictCompare[int]()
	testValues := []int{1, 2, 3, 5, 10}

	for _, x := range testValues {
		// Reflexivity: Compare(x, x) = 0
		assert.Equal(t, 0, intOrd.Compare(x, x), "Reflexivity failed for %d", x)

		for _, y := range testValues {
			cxy := intOrd.Compare(x, y)
			cyx := intOrd.Compare(y, x)

			// Antisymmetry: if Compare(x, y) <= 0 and Compare(y, x) <= 0 then x = y
			if cxy <= 0 && cyx <= 0 {
				assert.Equal(t, x, y, "Antisymmetry failed for %d and %d", x, y)
			}

			// Totality: Compare(x, y) <= 0 or Compare(y, x) <= 0
			assert.True(t, cxy <= 0 || cyx <= 0, "Totality failed for %d and %d", x, y)

			for _, z := range testValues {
				cyz := intOrd.Compare(y, z)
				cxz := intOrd.Compare(x, z)

				// Transitivity: if Compare(x, y) <= 0 and Compare(y, z) <= 0 then Compare(x, z) <= 0
				if cxy <= 0 && cyz <= 0 {
					assert.True(t, cxz <= 0, "Transitivity failed for %d, %d, %d", x, y, z)
				}
			}
		}
	}
}

// Test complex example with multi-level sorting
func TestComplexSorting(t *testing.T) {
	type Employee struct {
		Department string
		Salary     int
		Name       string
	}

	stringOrd := FromStrictCompare[string]()
	intOrd := FromStrictCompare[int]()

	// Order by department
	byDept := Contramap(func(e Employee) string {
		return e.Department
	})(stringOrd)

	// Order by salary (descending)
	bySalary := Reverse(Contramap(func(e Employee) int {
		return e.Salary
	})(intOrd))

	// Order by name
	byName := Contramap(func(e Employee) string {
		return e.Name
	})(stringOrd)

	// Combine: dept, then salary (desc), then name
	sg := Semigroup[Employee]()
	employeeOrd := sg.Concat(sg.Concat(byDept, bySalary), byName)

	e1 := Employee{Department: "IT", Salary: 100000, Name: "Alice"}
	e2 := Employee{Department: "IT", Salary: 100000, Name: "Bob"}
	e3 := Employee{Department: "IT", Salary: 90000, Name: "Charlie"}
	e4 := Employee{Department: "HR", Salary: 80000, Name: "David"}

	// Same dept, same salary, different name
	assert.Equal(t, -1, employeeOrd.Compare(e1, e2)) // Alice < Bob

	// Same dept, different salary
	assert.Equal(t, -1, employeeOrd.Compare(e1, e3)) // 100000 > 90000 (reversed)

	// Different dept
	assert.Equal(t, 1, employeeOrd.Compare(e1, e4)) // IT > HR
}

// Benchmark tests
func BenchmarkCompare(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = intOrd.Compare(i, i+1)
	}
}

func BenchmarkMin(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	min := Min(intOrd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = min(i, i+1)
	}
}

func BenchmarkMax(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	max := Max(intOrd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = max(i, i+1)
	}
}

func BenchmarkClamp(b *testing.B) {
	intOrd := FromStrictCompare[int]()
	clamp := Clamp(intOrd)(0, 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = clamp(i % 150)
	}
}
