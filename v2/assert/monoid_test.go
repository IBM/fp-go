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

package assert

import (
	"testing"
)

// TestApplicativeMonoid_Empty tests that Empty returns an assertion that always passes
func TestApplicativeMonoid_Empty(t *testing.T) {
	m := ApplicativeMonoid()
	empty := m.Empty()

	result := empty(t)
	if !result {
		t.Error("Expected Empty() to return an assertion that always passes")
	}
}

// TestApplicativeMonoid_Concat_BothPass tests that Concat returns true when both assertions pass
func TestApplicativeMonoid_Concat_BothPass(t *testing.T) {
	m := ApplicativeMonoid()

	assertion1 := Equal(42)(42)
	assertion2 := Equal("hello")("hello")

	combined := m.Concat(assertion1, assertion2)
	result := combined(t)

	if !result {
		t.Error("Expected Concat to pass when both assertions pass")
	}
}

// TestApplicativeMonoid_Concat_FirstFails tests that Concat returns false when first assertion fails
func TestApplicativeMonoid_Concat_FirstFails(t *testing.T) {
	mockT := &testing.T{}

	m := ApplicativeMonoid()

	assertion1 := Equal(42)(43) // This will fail
	assertion2 := Equal("hello")("hello")

	combined := m.Concat(assertion1, assertion2)
	result := combined(mockT)

	if result {
		t.Error("Expected Concat to fail when first assertion fails")
	}
}

// TestApplicativeMonoid_Concat_SecondFails tests that Concat returns false when second assertion fails
func TestApplicativeMonoid_Concat_SecondFails(t *testing.T) {
	mockT := &testing.T{}

	m := ApplicativeMonoid()

	assertion1 := Equal(42)(42)
	assertion2 := Equal("hello")("world") // This will fail

	combined := m.Concat(assertion1, assertion2)
	result := combined(mockT)

	if result {
		t.Error("Expected Concat to fail when second assertion fails")
	}
}

// TestApplicativeMonoid_Concat_BothFail tests that Concat returns false when both assertions fail
func TestApplicativeMonoid_Concat_BothFail(t *testing.T) {
	mockT := &testing.T{}

	m := ApplicativeMonoid()

	assertion1 := Equal(42)(43)           // This will fail
	assertion2 := Equal("hello")("world") // This will fail

	combined := m.Concat(assertion1, assertion2)
	result := combined(mockT)

	if result {
		t.Error("Expected Concat to fail when both assertions fail")
	}
}

// TestApplicativeMonoid_LeftIdentity tests the left identity law: Concat(Empty(), a) = a
func TestApplicativeMonoid_LeftIdentity(t *testing.T) {
	m := ApplicativeMonoid()

	assertion := Equal(42)(42)

	// Concat(Empty(), assertion) should behave the same as assertion
	combined := m.Concat(m.Empty(), assertion)

	result1 := assertion(t)
	result2 := combined(t)

	if result1 != result2 {
		t.Error("Left identity law violated: Concat(Empty(), a) should equal a")
	}
}

// TestApplicativeMonoid_RightIdentity tests the right identity law: Concat(a, Empty()) = a
func TestApplicativeMonoid_RightIdentity(t *testing.T) {
	m := ApplicativeMonoid()

	assertion := Equal(42)(42)

	// Concat(assertion, Empty()) should behave the same as assertion
	combined := m.Concat(assertion, m.Empty())

	result1 := assertion(t)
	result2 := combined(t)

	if result1 != result2 {
		t.Error("Right identity law violated: Concat(a, Empty()) should equal a")
	}
}

// TestApplicativeMonoid_Associativity tests the associativity law: Concat(Concat(a, b), c) = Concat(a, Concat(b, c))
func TestApplicativeMonoid_Associativity(t *testing.T) {
	m := ApplicativeMonoid()

	a1 := Equal(1)(1)
	a2 := Equal(2)(2)
	a3 := Equal(3)(3)

	// Concat(Concat(a1, a2), a3)
	left := m.Concat(m.Concat(a1, a2), a3)

	// Concat(a1, Concat(a2, a3))
	right := m.Concat(a1, m.Concat(a2, a3))

	result1 := left(t)
	result2 := right(t)

	if result1 != result2 {
		t.Error("Associativity law violated: Concat(Concat(a, b), c) should equal Concat(a, Concat(b, c))")
	}
}

// TestApplicativeMonoid_AssociativityWithFailure tests associativity when assertions fail
func TestApplicativeMonoid_AssociativityWithFailure(t *testing.T) {
	mockT := &testing.T{}
	m := ApplicativeMonoid()

	a1 := Equal(1)(1)
	a2 := Equal(2)(3) // This will fail
	a3 := Equal(3)(3)

	// Concat(Concat(a1, a2), a3)
	left := m.Concat(m.Concat(a1, a2), a3)

	// Concat(a1, Concat(a2, a3))
	right := m.Concat(a1, m.Concat(a2, a3))

	result1 := left(mockT)
	result2 := right(mockT)

	if result1 != result2 {
		t.Error("Associativity law violated even with failures")
	}

	if result1 || result2 {
		t.Error("Expected both to fail when one assertion fails")
	}
}

// TestApplicativeMonoid_ComplexAssertions tests combining complex assertions
func TestApplicativeMonoid_ComplexAssertions(t *testing.T) {
	m := ApplicativeMonoid()

	arr := []int{1, 2, 3, 4, 5}
	mp := map[string]int{"a": 1, "b": 2}

	arrayAssertions := m.Concat(
		ArrayNotEmpty(arr),
		m.Concat(
			ArrayLength[int](5)(arr),
			ArrayContains(3)(arr),
		),
	)

	mapAssertions := m.Concat(
		RecordNotEmpty(mp),
		RecordLength[string, int](2)(mp),
	)

	combined := m.Concat(arrayAssertions, mapAssertions)

	result := combined(t)
	if !result {
		t.Error("Expected complex combined assertions to pass")
	}
}

// TestApplicativeMonoid_ComplexAssertionsWithFailure tests complex assertions when one fails
func TestApplicativeMonoid_ComplexAssertionsWithFailure(t *testing.T) {
	mockT := &testing.T{}
	m := ApplicativeMonoid()

	arr := []int{1, 2, 3}
	mp := map[string]int{"a": 1, "b": 2}

	arrayAssertions := m.Concat(
		ArrayNotEmpty(arr),
		m.Concat(
			ArrayLength[int](5)(arr), // This will fail - array has 3 elements, not 5
			ArrayContains(3)(arr),
		),
	)

	mapAssertions := m.Concat(
		RecordNotEmpty(mp),
		RecordLength[string, int](2)(mp),
	)

	combined := m.Concat(arrayAssertions, mapAssertions)

	result := combined(mockT)
	if result {
		t.Error("Expected complex combined assertions to fail when one assertion fails")
	}
}

// TestApplicativeMonoid_MultipleConcat tests chaining multiple Concat operations
func TestApplicativeMonoid_MultipleConcat(t *testing.T) {
	m := ApplicativeMonoid()

	a1 := Equal(1)(1)
	a2 := Equal(2)(2)
	a3 := Equal(3)(3)
	a4 := Equal(4)(4)

	combined := m.Concat(
		m.Concat(a1, a2),
		m.Concat(a3, a4),
	)

	result := combined(t)
	if !result {
		t.Error("Expected multiple Concat operations to pass when all assertions pass")
	}
}

// TestApplicativeMonoid_WithStringAssertions tests combining string assertions
func TestApplicativeMonoid_WithStringAssertions(t *testing.T) {
	m := ApplicativeMonoid()

	str := "hello world"

	combined := m.Concat(
		StringNotEmpty(str),
		StringLength[any, any](11)(str),
	)

	result := combined(t)
	if !result {
		t.Error("Expected string assertions to pass")
	}
}

// TestApplicativeMonoid_WithBooleanAssertions tests combining boolean assertions
func TestApplicativeMonoid_WithBooleanAssertions(t *testing.T) {
	m := ApplicativeMonoid()

	combined := m.Concat(
		Equal(true)(true),
		m.Concat(
			Equal(false)(false),
			Equal(true)(true),
		),
	)

	result := combined(t)
	if !result {
		t.Error("Expected boolean assertions to pass")
	}
}

// TestApplicativeMonoid_WithErrorAssertions tests combining error assertions
func TestApplicativeMonoid_WithErrorAssertions(t *testing.T) {
	m := ApplicativeMonoid()

	combined := m.Concat(
		NoError(nil),
		Equal("test")("test"),
	)

	result := combined(t)
	if !result {
		t.Error("Expected error assertions to pass")
	}
}

// TestApplicativeMonoid_EmptyWithMultipleConcat tests Empty with multiple Concat operations
func TestApplicativeMonoid_EmptyWithMultipleConcat(t *testing.T) {
	m := ApplicativeMonoid()

	assertion := Equal(42)(42)

	// Multiple Empty values should still act as identity
	combined := m.Concat(
		m.Empty(),
		m.Concat(
			assertion,
			m.Empty(),
		),
	)

	result1 := assertion(t)
	result2 := combined(t)

	if result1 != result2 {
		t.Error("Multiple Empty values should still act as identity")
	}
}

// TestApplicativeMonoid_OnlyEmpty tests using only Empty values
func TestApplicativeMonoid_OnlyEmpty(t *testing.T) {
	m := ApplicativeMonoid()

	// Concat of Empty values should still be Empty (identity)
	combined := m.Concat(m.Empty(), m.Empty())

	result := combined(t)
	if !result {
		t.Error("Expected Concat of Empty values to pass")
	}
}

// TestApplicativeMonoid_RealWorldExample tests a realistic use case
func TestApplicativeMonoid_RealWorldExample(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	m := ApplicativeMonoid()

	validateUser := func(u User) Reader {
		return m.Concat(
			StringNotEmpty(u.Name),
			m.Concat(
				That(func(age int) bool { return age > 0 })(u.Age),
				m.Concat(
					That(func(age int) bool { return age < 150 })(u.Age),
					That(func(email string) bool {
						for _, ch := range email {
							if ch == '@' {
								return true
							}
						}
						return false
					})(u.Email),
				),
			),
		)
	}

	validUser := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
	result := validateUser(validUser)(t)

	if !result {
		t.Error("Expected valid user to pass all validations")
	}
}

// TestApplicativeMonoid_RealWorldExampleWithFailure tests a realistic use case with failure
func TestApplicativeMonoid_RealWorldExampleWithFailure(t *testing.T) {
	mockT := &testing.T{}

	type User struct {
		Name  string
		Age   int
		Email string
	}

	m := ApplicativeMonoid()

	validateUser := func(u User) Reader {
		return m.Concat(
			StringNotEmpty(u.Name),
			m.Concat(
				That(func(age int) bool { return age > 0 })(u.Age),
				m.Concat(
					That(func(age int) bool { return age < 150 })(u.Age),
					That(func(email string) bool {
						for _, ch := range email {
							if ch == '@' {
								return true
							}
						}
						return false
					})(u.Email),
				),
			),
		)
	}

	invalidUser := User{Name: "Bob", Age: 200, Email: "bob@test.com"} // Age > 150
	result := validateUser(invalidUser)(mockT)

	if result {
		t.Error("Expected invalid user to fail validation")
	}
}

// TestApplicativeMonoid_IntegrationWithAllOf demonstrates relationship with AllOf
func TestApplicativeMonoid_IntegrationWithAllOf(t *testing.T) {
	m := ApplicativeMonoid()
	arr := []int{1, 2, 3, 4, 5}

	// Using ApplicativeMonoid directly
	manualCombination := m.Concat(
		ArrayNotEmpty(arr),
		m.Concat(
			ArrayLength[int](5)(arr),
			ArrayContains(3)(arr),
		),
	)

	// Using AllOf (which uses ApplicativeMonoid internally)
	allOfCombination := AllOf([]Reader{
		ArrayNotEmpty(arr),
		ArrayLength[int](5)(arr),
		ArrayContains(3)(arr),
	})

	result1 := manualCombination(t)
	result2 := allOfCombination(t)

	if result1 != result2 {
		t.Error("Expected manual combination and AllOf to produce same result")
	}

	if !result1 || !result2 {
		t.Error("Expected both combinations to pass")
	}
}
