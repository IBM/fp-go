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

package monoid

import (
	"math"
	"testing"

	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/stretchr/testify/assert"
)

// assertMonoidLaws checks monoid laws for a given value
func assertMonoidLaws[A any](t *testing.T, m Monoid[A], a A) bool {
	e := m.Empty()
	return assert.Equal(t, a, m.Concat(a, e), "Monoid right identity") &&
		assert.Equal(t, a, m.Concat(e, a), "Monoid left identity")
}

// assertMonoidLawsForAll checks monoid laws for all values in a slice
func assertMonoidLawsForAll[A any](t *testing.T, m Monoid[A], data []A) bool {
	result := true
	for _, value := range data {
		result = result && assertMonoidLaws(t, m, value)
	}
	return result
}

// Test MakeMonoid creates a valid monoid
func TestMakeMonoid(t *testing.T) {
	// Integer addition monoid
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	assert.Equal(t, 8, addMonoid.Concat(5, 3))
	assert.Equal(t, 0, addMonoid.Empty())

	// Integer multiplication monoid
	mulMonoid := MakeMonoid(
		func(a, b int) int { return a * b },
		1,
	)

	assert.Equal(t, 15, mulMonoid.Concat(5, 3))
	assert.Equal(t, 1, mulMonoid.Empty())

	// String concatenation monoid
	stringMonoid := MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)

	assert.Equal(t, "HelloWorld", stringMonoid.Concat("Hello", "World"))
	assert.Equal(t, "", stringMonoid.Empty())
}

// Test monoid laws for integer addition
func TestMonoidLaws_IntAddition(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	testData := []int{0, 1, -1, 5, 10, -10, 100}

	// Test using helper function
	assert.True(t, assertMonoidLawsForAll(t, addMonoid, testData))

	// Explicit law tests
	for _, x := range testData {
		// Left identity: Empty() + x = x
		assert.Equal(t, x, addMonoid.Concat(addMonoid.Empty(), x),
			"Left identity failed for %d", x)

		// Right identity: x + Empty() = x
		assert.Equal(t, x, addMonoid.Concat(x, addMonoid.Empty()),
			"Right identity failed for %d", x)
	}

	// Associativity (from Semigroup)
	assert.Equal(t,
		addMonoid.Concat(addMonoid.Concat(1, 2), 3),
		addMonoid.Concat(1, addMonoid.Concat(2, 3)),
	)
}

// Test monoid laws for integer multiplication
func TestMonoidLaws_IntMultiplication(t *testing.T) {
	mulMonoid := MakeMonoid(
		func(a, b int) int { return a * b },
		1,
	)

	testData := []int{1, 2, 3, 5, 10}

	assert.True(t, assertMonoidLawsForAll(t, mulMonoid, testData))

	// Explicit tests
	for _, x := range testData {
		assert.Equal(t, x, mulMonoid.Concat(mulMonoid.Empty(), x))
		assert.Equal(t, x, mulMonoid.Concat(x, mulMonoid.Empty()))
	}
}

// Test monoid laws for string concatenation
func TestMonoidLaws_String(t *testing.T) {
	stringMonoid := MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)

	testData := []string{"", "a", "hello", "world", "test"}

	assert.True(t, assertMonoidLawsForAll(t, stringMonoid, testData))
}

// Test monoid laws for boolean AND
func TestMonoidLaws_BooleanAnd(t *testing.T) {
	andMonoid := MakeMonoid(
		func(a, b bool) bool { return a && b },
		true,
	)

	testData := []bool{true, false}

	assert.True(t, assertMonoidLawsForAll(t, andMonoid, testData))

	// Specific tests
	assert.Equal(t, true, andMonoid.Concat(true, true))
	assert.Equal(t, false, andMonoid.Concat(true, false))
	assert.Equal(t, false, andMonoid.Concat(false, true))
	assert.Equal(t, false, andMonoid.Concat(false, false))
	assert.Equal(t, true, andMonoid.Empty())
}

// Test monoid laws for boolean OR
func TestMonoidLaws_BooleanOr(t *testing.T) {
	orMonoid := MakeMonoid(
		func(a, b bool) bool { return a || b },
		false,
	)

	testData := []bool{true, false}

	assert.True(t, assertMonoidLawsForAll(t, orMonoid, testData))

	// Specific tests
	assert.Equal(t, true, orMonoid.Concat(true, true))
	assert.Equal(t, true, orMonoid.Concat(true, false))
	assert.Equal(t, true, orMonoid.Concat(false, true))
	assert.Equal(t, false, orMonoid.Concat(false, false))
	assert.Equal(t, false, orMonoid.Empty())
}

// Test Reverse swaps argument order
func TestReverse(t *testing.T) {
	// Subtraction is not commutative, so reverse will give different results
	subMonoid := MakeMonoid(
		func(a, b int) int { return a - b },
		0,
	)

	reversedMonoid := Reverse(subMonoid)

	// Original: 10 - 3 = 7
	assert.Equal(t, 7, subMonoid.Concat(10, 3))

	// Reversed: 3 - 10 = -7
	assert.Equal(t, -7, reversedMonoid.Concat(10, 3))

	// Empty should be the same
	assert.Equal(t, subMonoid.Empty(), reversedMonoid.Empty())

	// Test with string concatenation
	stringMonoid := MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)

	reversedString := Reverse(stringMonoid)

	assert.Equal(t, "HelloWorld", stringMonoid.Concat("Hello", "World"))
	assert.Equal(t, "WorldHello", reversedString.Concat("Hello", "World"))
}

// Test Reverse preserves monoid laws
func TestReverse_PreservesLaws(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	reversedMonoid := Reverse(addMonoid)

	testData := []int{0, 1, 5, 10}

	// Reversed monoid should still satisfy monoid laws
	assert.True(t, assertMonoidLawsForAll(t, reversedMonoid, testData))
}

// Test ToSemigroup conversion
func TestToSemigroup(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	sg := ToSemigroup(addMonoid)

	// Should work as a semigroup
	assert.Equal(t, 8, sg.Concat(5, 3))
	assert.Equal(t, 15, sg.Concat(10, 5))

	// Verify it's the same underlying operation
	assert.Equal(t, addMonoid.Concat(5, 3), sg.Concat(5, 3))
}

// Test ConcatAll with various inputs
func TestConcatAll(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	concatAll := ConcatAll(addMonoid)

	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"empty slice", []int{}, 0},
		{"single element", []int{5}, 5},
		{"multiple elements", []int{1, 2, 3, 4, 5}, 15},
		{"with zeros", []int{0, 5, 0, 3}, 8},
		{"negative numbers", []int{-1, -2, -3}, -6},
		{"mixed signs", []int{10, -5, 3, -2}, 6},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := concatAll(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test ConcatAll with multiplication
func TestConcatAll_Multiplication(t *testing.T) {
	mulMonoid := MakeMonoid(
		func(a, b int) int { return a * b },
		1,
	)

	concatAll := ConcatAll(mulMonoid)

	tests := []struct {
		name     string
		input    []int
		expected int
	}{
		{"empty slice", []int{}, 1},
		{"single element", []int{5}, 5},
		{"multiple elements", []int{2, 3, 4}, 24},
		{"with one", []int{1, 5, 1, 3}, 15},
		{"with zero", []int{2, 0, 3}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := concatAll(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test ConcatAll with strings
func TestConcatAll_String(t *testing.T) {
	stringMonoid := MakeMonoid(
		func(a, b string) string { return a + b },
		"",
	)

	concatAll := ConcatAll(stringMonoid)

	tests := []struct {
		name     string
		input    []string
		expected string
	}{
		{"empty slice", []string{}, ""},
		{"single element", []string{"hello"}, "hello"},
		{"multiple elements", []string{"Hello", " ", "World", "!"}, "Hello World!"},
		{"with empty strings", []string{"a", "", "b", "", "c"}, "abc"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := concatAll(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test Fold (alias for ConcatAll)
func TestFold(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	fold := Fold(addMonoid)
	concatAll := ConcatAll(addMonoid)

	input := []int{1, 2, 3, 4, 5}

	// Fold and ConcatAll should produce the same result
	assert.Equal(t, concatAll(input), fold(input))
	assert.Equal(t, 15, fold(input))
}

// Test GenericConcatAll with custom slice type
func TestGenericConcatAll(t *testing.T) {
	type IntSlice []int

	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	concatAll := GenericConcatAll[IntSlice](addMonoid)

	tests := []struct {
		name     string
		input    IntSlice
		expected int
	}{
		{"empty slice", IntSlice{}, 0},
		{"single element", IntSlice{5}, 5},
		{"multiple elements", IntSlice{1, 2, 3, 4, 5}, 15},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := concatAll(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Test FunctionMonoid
func TestFunctionMonoid(t *testing.T) {
	// Create a monoid for functions that return integers
	intAddMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	funcMonoid := FunctionMonoid[string](intAddMonoid)

	// Create some functions
	f1 := func(s string) int { return len(s) }
	f2 := func(s string) int { return len(s) * 2 }
	f3 := func(s string) int { return 10 }

	// Combine functions
	combined := funcMonoid.Concat(f1, f2)

	// Test combined function
	assert.Equal(t, 15, combined("hello")) // len("hello") + len("hello")*2 = 5 + 10 = 15
	assert.Equal(t, 9, combined("abc"))    // len("abc") + len("abc")*2 = 3 + 6 = 9

	// Test with three functions
	combined3 := funcMonoid.Concat(funcMonoid.Concat(f1, f2), f3)
	assert.Equal(t, 25, combined3("hello")) // 5 + 10 + 10 = 25

	// Test empty function
	emptyFunc := funcMonoid.Empty()
	assert.Equal(t, 0, emptyFunc("anything"))
	assert.Equal(t, 0, emptyFunc(""))

	// Test identity laws
	input := "test"
	assert.Equal(t, f1(input), funcMonoid.Concat(funcMonoid.Empty(), f1)(input))
	assert.Equal(t, f1(input), funcMonoid.Concat(f1, funcMonoid.Empty())(input))
}

// Test FunctionMonoid with different types
func TestFunctionMonoid_DifferentTypes(t *testing.T) {
	// Monoid for functions that return booleans (OR)
	boolOrMonoid := MakeMonoid(
		func(a, b bool) bool { return a || b },
		false,
	)

	funcMonoid := FunctionMonoid[int](boolOrMonoid)

	isEven := func(n int) bool { return n%2 == 0 }
	isPositive := func(n int) bool { return n > 0 }

	// Combine: returns true if number is even OR positive
	combined := funcMonoid.Concat(isEven, isPositive)

	assert.True(t, combined(4))   // even and positive
	assert.True(t, combined(3))   // odd but positive
	assert.True(t, combined(-2))  // even but negative
	assert.False(t, combined(-3)) // odd and negative
}

// Test ApplicativeMonoid with a simple applicative
func TestApplicativeMonoid(t *testing.T) {
	// Simple Option-like type for testing
	type Option[A any] struct {
		value *A
	}

	some := func(a int) Option[int] {
		return Option[int]{value: &a}
	}

	none := func() Option[int] {
		return Option[int]{value: nil}
	}

	fmap := func(opt Option[int], f func(int) func(int) int) Option[func(int) int] {
		if opt.value == nil {
			return Option[func(int) int]{value: nil}
		}
		fn := f(*opt.value)
		return Option[func(int) int]{value: &fn}
	}

	fap := func(optF Option[func(int) int], opt Option[int]) Option[int] {
		if optF.value == nil || opt.value == nil {
			return none()
		}
		result := (*optF.value)(*opt.value)
		return some(result)
	}

	intAddMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	optMonoid := ApplicativeMonoid(some, fmap, fap, intAddMonoid)

	// Test concat of Some values
	opt1 := some(5)
	opt2 := some(3)
	result := optMonoid.Concat(opt1, opt2)
	assert.NotNil(t, result.value)
	assert.Equal(t, 8, *result.value)

	// Test empty
	empty := optMonoid.Empty()
	assert.NotNil(t, empty.value)
	assert.Equal(t, 0, *empty.value)

	// Test identity laws
	assert.Equal(t, opt1, optMonoid.Concat(optMonoid.Empty(), opt1))
	assert.Equal(t, opt1, optMonoid.Concat(opt1, optMonoid.Empty()))
}

// Test AlternativeMonoid
func TestAlternativeMonoid(t *testing.T) {
	// Simple Option-like type for testing
	type Option[A any] struct {
		value *A
	}

	some := func(a int) Option[int] {
		return Option[int]{value: &a}
	}

	none := func() Option[int] {
		return Option[int]{value: nil}
	}

	fmap := func(opt Option[int], f func(int) func(int) int) Option[func(int) int] {
		if opt.value == nil {
			return Option[func(int) int]{value: nil}
		}
		fn := f(*opt.value)
		return Option[func(int) int]{value: &fn}
	}

	fap := func(optF Option[func(int) int], opt Option[int]) Option[int] {
		if optF.value == nil || opt.value == nil {
			return none()
		}
		result := (*optF.value)(*opt.value)
		return some(result)
	}

	falt := func(first Option[int], second func() Option[int]) Option[int] {
		if first.value != nil {
			return first
		}
		return second()
	}

	intAddMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	optMonoid := AlternativeMonoid(some, fmap, fap, falt, intAddMonoid)

	// Test concat of Some values - should add them
	opt1 := some(5)
	opt2 := some(3)
	result := optMonoid.Concat(opt1, opt2)
	assert.NotNil(t, result.value)
	assert.Equal(t, 8, *result.value)

	// Test concat with None - should fall back
	result2 := optMonoid.Concat(none(), opt2)
	assert.NotNil(t, result2.value)
	assert.Equal(t, 3, *result2.value)

	// Test concat: first Some wins in alternative
	result3 := optMonoid.Concat(opt1, opt2)
	assert.NotNil(t, result3.value)
	// The result should be 8 (5 + 3) from applicative behavior

	// Test empty
	empty := optMonoid.Empty()
	assert.NotNil(t, empty.value)
	assert.Equal(t, 0, *empty.value)

	// Test identity laws
	assert.Equal(t, opt1, optMonoid.Concat(optMonoid.Empty(), opt1))
	assert.Equal(t, opt1, optMonoid.Concat(opt1, optMonoid.Empty()))
}

// Test AltMonoid
func TestAltMonoid(t *testing.T) {
	// Simple Option-like type
	type Option[A any] struct {
		value *A
	}

	some := func(a int) Option[int] {
		return Option[int]{value: &a}
	}

	none := func() Option[int] {
		return Option[int]{value: nil}
	}

	falt := func(first Option[int], second func() Option[int]) Option[int] {
		if first.value != nil {
			return first
		}
		return second()
	}

	optMonoid := AltMonoid(none, falt)

	// Test concat: first Some wins
	opt1 := some(5)
	opt2 := some(3)
	result := optMonoid.Concat(opt1, opt2)
	assert.NotNil(t, result.value)
	assert.Equal(t, 5, *result.value)

	// Test concat: None falls back to second
	result2 := optMonoid.Concat(none(), opt2)
	assert.NotNil(t, result2.value)
	assert.Equal(t, 3, *result2.value)

	// Test concat: both None
	result3 := optMonoid.Concat(none(), none())
	assert.Nil(t, result3.value)

	// Test empty
	empty := optMonoid.Empty()
	assert.Nil(t, empty.value)
}

// Test monoid with max operation
func TestMonoid_Max(t *testing.T) {
	maxMonoid := MakeMonoid(
		func(a, b int) int {
			if a > b {
				return a
			}
			return b
		},
		math.MinInt,
	)

	numbers := []int{3, 7, 2, 9, 1, 5}
	maximum := ConcatAll(maxMonoid)(numbers)
	assert.Equal(t, 9, maximum)

	// Empty slice returns identity (MinInt)
	emptyMax := ConcatAll(maxMonoid)([]int{})
	assert.Equal(t, math.MinInt, emptyMax)
}

// Test monoid with min operation
func TestMonoid_Min(t *testing.T) {
	minMonoid := MakeMonoid(
		func(a, b int) int {
			if a < b {
				return a
			}
			return b
		},
		math.MaxInt,
	)

	numbers := []int{3, 7, 2, 9, 1, 5}
	minimum := ConcatAll(minMonoid)(numbers)
	assert.Equal(t, 1, minimum)

	// Empty slice returns identity (MaxInt)
	emptyMin := ConcatAll(minMonoid)([]int{})
	assert.Equal(t, math.MaxInt, emptyMin)
}

// Test monoid with list concatenation
func TestMonoid_ListConcat(t *testing.T) {
	listMonoid := MakeMonoid(
		func(a, b []int) []int {
			result := make([]int, len(a)+len(b))
			copy(result, a)
			copy(result[len(a):], b)
			return result
		},
		[]int{},
	)

	lists := [][]int{{1, 2}, {3, 4}, {5}}
	flattened := ConcatAll(listMonoid)(lists)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, flattened)

	// Empty slice returns empty list
	emptyList := ConcatAll(listMonoid)([][]int{})
	assert.Equal(t, []int{}, emptyList)

	// Test identity laws
	list := []int{1, 2, 3}
	assert.Equal(t, list, listMonoid.Concat(listMonoid.Empty(), list))
	assert.Equal(t, list, listMonoid.Concat(list, listMonoid.Empty()))
}

// Test that monoid interface is properly implemented
func TestMonoidInterface(t *testing.T) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	// Verify it implements Monoid interface
	var _ Monoid[int] = addMonoid

	// Verify it also implements Semigroup interface (through embedding)
	var _ S.Semigroup[int] = addMonoid
}

// Benchmark ConcatAll
func BenchmarkConcatAll(b *testing.B) {
	addMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	concatAll := ConcatAll(addMonoid)
	numbers := make([]int, 1000)
	for i := range numbers {
		numbers[i] = i
	}

	b.ResetTimer()
	for b.Loop() {
		_ = concatAll(numbers)
	}
}

// Benchmark FunctionMonoid
func BenchmarkFunctionMonoid(b *testing.B) {
	intAddMonoid := MakeMonoid(
		func(a, b int) int { return a + b },
		0,
	)

	funcMonoid := FunctionMonoid[string](intAddMonoid)

	f1 := func(s string) int { return len(s) }
	f2 := func(s string) int { return len(s) * 2 }

	combined := funcMonoid.Concat(f1, f2)

	b.ResetTimer()
	for b.Loop() {
		_ = combined("benchmark")
	}
}
