// Copyright (c) 2023 IBM Corp.
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

package endomorphism

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/stretchr/testify/assert"
)

// Test helper functions
func double(x int) int {
	return x * 2
}

func increment(x int) int {
	return x + 1
}

func square(x int) int {
	return x * x
}

func negate(x int) int {
	return -x
}

// TestCurry2 tests the Curry2 function
func TestCurry2(t *testing.T) {
	add := func(x, y int) int {
		return x + y
	}

	curriedAdd := Curry2(add)
	addFive := curriedAdd(5)

	result := addFive(10)
	assert.Equal(t, 15, result, "Curry2 should curry binary function correctly")

	// Test with different values
	addTen := curriedAdd(10)
	assert.Equal(t, 25, addTen(15), "Curry2 should work with different values")
}

// TestCurry3 tests the Curry3 function
func TestCurry3(t *testing.T) {
	combine := func(x, y, z int) int {
		return x + y + z
	}

	curriedCombine := Curry3(combine)
	addTen := curriedCombine(5)(5)

	result := addTen(20)
	assert.Equal(t, 30, result, "Curry3 should curry ternary function correctly")

	// Test with different values
	addFifteen := curriedCombine(5)(10)
	assert.Equal(t, 35, addFifteen(20), "Curry3 should work with different values")
}

// TestMonadAp tests the MonadAp function
func TestMonadAp(t *testing.T) {
	result := MonadAp(double, 5)
	assert.Equal(t, 10, result, "MonadAp should apply endomorphism to value")

	result2 := MonadAp(increment, 10)
	assert.Equal(t, 11, result2, "MonadAp should work with different endomorphisms")

	result3 := MonadAp(square, 4)
	assert.Equal(t, 16, result3, "MonadAp should work with square function")
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	applyFive := Ap(5)

	result := applyFive(double)
	assert.Equal(t, 10, result, "Ap should apply value to endomorphism")

	result2 := applyFive(increment)
	assert.Equal(t, 6, result2, "Ap should work with different endomorphisms")

	applyTen := Ap(10)
	result3 := applyTen(square)
	assert.Equal(t, 100, result3, "Ap should work with different values")
}

// TestCompose tests the Compose function
func TestCompose(t *testing.T) {
	// Test basic composition: (5 * 2) + 1 = 11
	doubleAndIncrement := Compose(double, increment)
	result := doubleAndIncrement(5)
	assert.Equal(t, 11, result, "Compose should compose endomorphisms correctly")

	// Test composition order: (5 + 1) * 2 = 12
	incrementAndDouble := Compose(increment, double)
	result2 := incrementAndDouble(5)
	assert.Equal(t, 12, result2, "Compose should respect order of composition")

	// Test with three compositions: ((5 * 2) + 1) * ((5 * 2) + 1) = 121
	complex := Compose(Compose(double, increment), square)
	result3 := complex(5)
	assert.Equal(t, 121, result3, "Compose should work with nested compositions")
}

// TestMonadChain tests the MonadChain function
func TestMonadChain(t *testing.T) {
	// MonadChain should behave like Compose
	chained := MonadChain(double, increment)
	result := chained(5)
	assert.Equal(t, 11, result, "MonadChain should chain endomorphisms correctly")

	chained2 := MonadChain(increment, double)
	result2 := chained2(5)
	assert.Equal(t, 12, result2, "MonadChain should respect order")

	// Test with negative values
	chained3 := MonadChain(negate, increment)
	result3 := chained3(5)
	assert.Equal(t, -4, result3, "MonadChain should work with negative values")
}

// TestChain tests the Chain function
func TestChain(t *testing.T) {
	chainWithIncrement := Chain(increment)

	chained := chainWithIncrement(double)
	result := chained(5)
	assert.Equal(t, 11, result, "Chain should create chaining function correctly")

	chainWithDouble := Chain(double)
	chained2 := chainWithDouble(increment)
	result2 := chained2(5)
	assert.Equal(t, 12, result2, "Chain should work with different endomorphisms")

	// Test chaining with square
	chainWithSquare := Chain(square)
	chained3 := chainWithSquare(double)
	result3 := chained3(3)
	assert.Equal(t, 36, result3, "Chain should work with square function")
}

// TestOf tests the Of function
func TestOf(t *testing.T) {
	endo := Of(double)
	result := endo(5)
	assert.Equal(t, 10, result, "Of should convert function to endomorphism")

	endo2 := Of(increment)
	result2 := endo2(10)
	assert.Equal(t, 11, result2, "Of should work with different functions")
}

// TestWrap tests the Wrap function (deprecated)
func TestWrap(t *testing.T) {
	endo := Wrap(double)
	result := endo(5)
	assert.Equal(t, 10, result, "Wrap should convert function to endomorphism")
}

// TestUnwrap tests the Unwrap function (deprecated)
func TestUnwrap(t *testing.T) {
	endo := Of(double)
	unwrapped := Unwrap[func(int) int](endo)
	result := unwrapped(5)
	assert.Equal(t, 10, result, "Unwrap should convert endomorphism to function")
}

// TestIdentity tests the Identity function
func TestIdentity(t *testing.T) {
	id := Identity[int]()

	// Identity should return input unchanged
	assert.Equal(t, 42, id(42), "Identity should return input unchanged")
	assert.Equal(t, 0, id(0), "Identity should work with zero")
	assert.Equal(t, -10, id(-10), "Identity should work with negative values")

	// Identity should be neutral for composition
	composed1 := Compose(id, double)
	assert.Equal(t, 10, composed1(5), "Identity should be right neutral for composition")

	composed2 := Compose(double, id)
	assert.Equal(t, 10, composed2(5), "Identity should be left neutral for composition")

	// Test with strings
	idStr := Identity[string]()
	assert.Equal(t, "hello", idStr("hello"), "Identity should work with strings")
}

// TestSemigroup tests the Semigroup function
func TestSemigroup(t *testing.T) {
	sg := Semigroup[int]()

	// Test basic concat
	combined := sg.Concat(double, increment)
	result := combined(5)
	assert.Equal(t, 11, result, "Semigroup concat should compose endomorphisms")

	// Test associativity: (f . g) . h = f . (g . h)
	f := double
	g := increment
	h := square

	left := sg.Concat(sg.Concat(f, g), h)
	right := sg.Concat(f, sg.Concat(g, h))

	testValue := 3
	assert.Equal(t, left(testValue), right(testValue), "Semigroup should be associative")

	// Test with ConcatAll from semigroup package
	combined2 := S.ConcatAll(sg)(double)([]Endomorphism[int]{increment, square})
	result2 := combined2(5)
	assert.Equal(t, 121, result2, "Semigroup should work with ConcatAll")
}

// TestMonoid tests the Monoid function
func TestMonoid(t *testing.T) {
	monoid := Monoid[int]()

	// Test that empty is identity
	empty := monoid.Empty()
	assert.Equal(t, 42, empty(42), "Monoid empty should be identity")

	// Test right identity: x . empty = x
	rightIdentity := monoid.Concat(double, empty)
	assert.Equal(t, 10, rightIdentity(5), "Monoid should satisfy right identity")

	// Test left identity: empty . x = x
	leftIdentity := monoid.Concat(empty, double)
	assert.Equal(t, 10, leftIdentity(5), "Monoid should satisfy left identity")

	// Test ConcatAll with multiple endomorphisms
	combined := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square})
	result := combined(5)
	// (5 * 2) = 10, (10 + 1) = 11, (11 * 11) = 121
	assert.Equal(t, 121, result, "Monoid should work with ConcatAll")

	// Test ConcatAll with empty list should return identity
	emptyResult := M.ConcatAll(monoid)([]Endomorphism[int]{})
	assert.Equal(t, 42, emptyResult(42), "ConcatAll with no args should return identity")
}

// TestMonoidLaws tests that the Monoid satisfies monoid laws
func TestMonoidLaws(t *testing.T) {
	monoid := Monoid[int]()
	empty := monoid.Empty()

	testCases := []struct {
		name string
		f    Endomorphism[int]
		g    Endomorphism[int]
		h    Endomorphism[int]
	}{
		{"basic", double, increment, square},
		{"with negate", negate, double, increment},
		{"with identity", Identity[int](), double, increment},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testValue := 5

			// Right identity: x . empty = x
			rightId := monoid.Concat(tc.f, empty)
			assert.Equal(t, tc.f(testValue), rightId(testValue), "Right identity law")

			// Left identity: empty . x = x
			leftId := monoid.Concat(empty, tc.f)
			assert.Equal(t, tc.f(testValue), leftId(testValue), "Left identity law")

			// Associativity: (f . g) . h = f . (g . h)
			left := monoid.Concat(monoid.Concat(tc.f, tc.g), tc.h)
			right := monoid.Concat(tc.f, monoid.Concat(tc.g, tc.h))
			assert.Equal(t, left(testValue), right(testValue), "Associativity law")
		})
	}
}

// TestEndomorphismWithDifferentTypes tests endomorphisms with different types
func TestEndomorphismWithDifferentTypes(t *testing.T) {
	// Test with strings
	toUpper := func(s string) string {
		return s + "!"
	}
	addPrefix := func(s string) string {
		return "Hello, " + s
	}

	strComposed := Compose(toUpper, addPrefix)
	result := strComposed("World")
	assert.Equal(t, "Hello, World!", result, "Endomorphism should work with strings")

	// Test with float64
	doubleFloat := func(x float64) float64 {
		return x * 2.0
	}
	addOne := func(x float64) float64 {
		return x + 1.0
	}

	floatComposed := Compose(doubleFloat, addOne)
	resultFloat := floatComposed(5.5)
	assert.Equal(t, 12.0, resultFloat, "Endomorphism should work with float64")
}

// TestComplexCompositions tests more complex composition scenarios
func TestComplexCompositions(t *testing.T) {
	// Create a pipeline of transformations
	pipeline := Compose(
		Compose(
			Compose(double, increment),
			square,
		),
		negate,
	)

	// (5 * 2) = 10, (10 + 1) = 11, (11 * 11) = 121, -(121) = -121
	result := pipeline(5)
	assert.Equal(t, -121, result, "Complex composition should work correctly")

	// Test using monoid to build the same pipeline
	monoid := Monoid[int]()
	pipelineMonoid := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square, negate})
	resultMonoid := pipelineMonoid(5)
	assert.Equal(t, -121, resultMonoid, "Monoid-based pipeline should match composition")
}

// TestOperatorType tests the Operator type
func TestOperatorType(t *testing.T) {
	// Create an operator that lifts an int endomorphism to work on the length of strings
	lengthOperator := func(f Endomorphism[int]) Endomorphism[string] {
		return func(s string) string {
			newLen := f(len(s))
			if newLen > len(s) {
				// Pad with spaces
				for i := len(s); i < newLen; i++ {
					s += " "
				}
			} else if newLen < len(s) {
				// Truncate
				s = s[:newLen]
			}
			return s
		}
	}

	// Use the operator
	var op Operator[int, string] = lengthOperator
	doubleLength := op(double)

	result := doubleLength("hello") // len("hello") = 5, 5 * 2 = 10
	assert.Equal(t, 10, len(result), "Operator should transform endomorphisms correctly")
	assert.Equal(t, "hello     ", result, "Operator should pad string correctly")
}

// BenchmarkCompose benchmarks the Compose function
func BenchmarkCompose(b *testing.B) {
	composed := Compose(double, increment)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = composed(5)
	}
}

// BenchmarkMonoidConcatAll benchmarks ConcatAll with monoid
func BenchmarkMonoidConcatAll(b *testing.B) {
	monoid := Monoid[int]()
	combined := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = combined(5)
	}
}

// BenchmarkChain benchmarks the Chain function
func BenchmarkChain(b *testing.B) {
	chainWithIncrement := Chain(increment)
	chained := chainWithIncrement(double)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chained(5)
	}
}
