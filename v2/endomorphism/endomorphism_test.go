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

package endomorphism

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
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
	// MonadAp composes two endomorphisms (RIGHT-TO-LEFT)
	// MonadAp(double, increment) means: increment first, then double
	composed := MonadAp(double, increment)
	result := composed(5)
	assert.Equal(t, 12, result, "MonadAp should compose right-to-left: (5 + 1) * 2 = 12")

	// Test with different order
	composed2 := MonadAp(increment, double)
	result2 := composed2(5)
	assert.Equal(t, 11, result2, "MonadAp should compose right-to-left: (5 * 2) + 1 = 11")

	// Test with square
	composed3 := MonadAp(square, increment)
	result3 := composed3(5)
	assert.Equal(t, 36, result3, "MonadAp should compose right-to-left: (5 + 1) ^ 2 = 36")
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	// Ap is the curried version of MonadAp
	// Ap(increment) returns a function that composes with increment (RIGHT-TO-LEFT)
	applyIncrement := Ap(increment)

	composed := applyIncrement(double)
	result := composed(5)
	assert.Equal(t, 12, result, "Ap should compose right-to-left: (5 + 1) * 2 = 12")

	// Test with different endomorphism
	composed2 := applyIncrement(square)
	result2 := composed2(5)
	assert.Equal(t, 36, result2, "Ap should compose right-to-left: (5 + 1) ^ 2 = 36")

	// Test with different base endomorphism
	applyDouble := Ap(double)
	composed3 := applyDouble(increment)
	result3 := composed3(5)
	assert.Equal(t, 11, result3, "Ap should compose right-to-left: (5 * 2) + 1 = 11")
}

// TestMonadCompose tests the MonadCompose function
func TestMonadCompose(t *testing.T) {
	// Test basic composition: RIGHT-TO-LEFT execution
	// MonadCompose(double, increment) means: increment first, then double
	composed := MonadCompose(double, increment)
	result := composed(5)
	assert.Equal(t, 12, result, "MonadCompose should execute right-to-left: (5 + 1) * 2 = 12")

	// Test composition order: RIGHT-TO-LEFT execution
	// MonadCompose(increment, double) means: double first, then increment
	composed2 := MonadCompose(increment, double)
	result2 := composed2(5)
	assert.Equal(t, 11, result2, "MonadCompose should execute right-to-left: (5 * 2) + 1 = 11")

	// Test with three compositions: RIGHT-TO-LEFT execution
	// MonadCompose(MonadCompose(double, increment), square) means: square, then increment, then double
	complex := MonadCompose(MonadCompose(double, increment), square)
	result3 := complex(5)
	// 5 -> square -> 25 -> increment -> 26 -> double -> 52
	assert.Equal(t, 52, result3, "MonadCompose should work with nested compositions: square(5)=25, +1=26, *2=52")
}

// TestMonadChain tests the MonadChain function
func TestMonadChain(t *testing.T) {
	// MonadChain executes LEFT-TO-RIGHT (first arg first, second arg second)
	chained := MonadChain(double, increment)
	result := chained(5)
	assert.Equal(t, 11, result, "MonadChain should execute left-to-right: (5 * 2) + 1 = 11")

	chained2 := MonadChain(increment, double)
	result2 := chained2(5)
	assert.Equal(t, 12, result2, "MonadChain should execute left-to-right: (5 + 1) * 2 = 12")

	// Test with negative values
	chained3 := MonadChain(negate, increment)
	result3 := chained3(5)
	assert.Equal(t, -4, result3, "MonadChain should execute left-to-right: -(5) + 1 = -4")
}

// TestChain tests the Chain function
func TestChain(t *testing.T) {
	// Chain(f) returns a function that applies its argument first, then f
	chainWithIncrement := Chain(increment)

	// chainWithIncrement(double) means: double first, then increment
	chained := chainWithIncrement(double)
	result := chained(5)
	assert.Equal(t, 11, result, "Chain should execute left-to-right: (5 * 2) + 1 = 11")

	chainWithDouble := Chain(double)
	// chainWithDouble(increment) means: increment first, then double
	chained2 := chainWithDouble(increment)
	result2 := chained2(5)
	assert.Equal(t, 12, result2, "Chain should execute left-to-right: (5 + 1) * 2 = 12")

	// Test chaining with square
	chainWithSquare := Chain(square)
	// chainWithSquare(double) means: double first, then square
	chained3 := chainWithSquare(double)
	result3 := chained3(3)
	assert.Equal(t, 36, result3, "Chain should execute left-to-right: (3 * 2) ^ 2 = 36")
}

// TestCompose tests the curried Compose function
func TestCompose(t *testing.T) {
	// Compose(g) returns a function that applies g first, then its argument
	composeWithIncrement := Compose(increment)

	// composeWithIncrement(double) means: increment first, then double
	composed := composeWithIncrement(double)
	result := composed(5)
	assert.Equal(t, 12, result, "Compose should execute right-to-left: (5 + 1) * 2 = 12")

	composeWithDouble := Compose(double)
	// composeWithDouble(increment) means: double first, then increment
	composed2 := composeWithDouble(increment)
	result2 := composed2(5)
	assert.Equal(t, 11, result2, "Compose should execute right-to-left: (5 * 2) + 1 = 11")

	// Test composing with square
	composeWithSquare := Compose(square)
	// composeWithSquare(double) means: square first, then double
	composed3 := composeWithSquare(double)
	result3 := composed3(3)
	assert.Equal(t, 18, result3, "Compose should execute right-to-left: (3 ^ 2) * 2 = 18")
}

// TestMonadComposeVsCompose demonstrates the relationship between MonadCompose and Compose
func TestMonadComposeVsCompose(t *testing.T) {
	double := N.Mul(2)
	increment := N.Add(1)

	// MonadCompose takes both functions at once
	monadComposed := MonadCompose(double, increment)
	result1 := monadComposed(5) // (5 + 1) * 2 = 12

	// Compose is the curried version - takes one function, returns a function
	curriedCompose := Compose(increment)
	composed := curriedCompose(double)
	result2 := composed(5) // (5 + 1) * 2 = 12

	assert.Equal(t, result1, result2, "MonadCompose and Compose should produce the same result")
	assert.Equal(t, 12, result1, "Both should execute right-to-left: (5 + 1) * 2 = 12")

	// Demonstrate that Compose(g)(f) is equivalent to MonadCompose(f, g)
	assert.Equal(t, MonadCompose(double, increment)(5), Compose(increment)(double)(5),
		"Compose(g)(f) should equal MonadCompose(f, g)")
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

	// Identity should be neutral for composition (RIGHT-TO-LEFT)
	// Compose(id, double) means: double first, then id
	composed1 := MonadCompose(id, double)
	assert.Equal(t, 10, composed1(5), "Identity should be left neutral: double(5) = 10")

	// Compose(double, id) means: id first, then double
	composed2 := MonadCompose(double, id)
	assert.Equal(t, 10, composed2(5), "Identity should be right neutral: id(5) then double = 10")

	// Test with strings
	idStr := Identity[string]()
	assert.Equal(t, "hello", idStr("hello"), "Identity should work with strings")
}

// TestSemigroup tests the Semigroup function
func TestSemigroup(t *testing.T) {
	sg := Semigroup[int]()

	// Test basic concat (RIGHT-TO-LEFT execution via Compose)
	// Concat(double, increment) means: increment first, then double
	combined := sg.Concat(double, increment)
	result := combined(5)
	assert.Equal(t, 12, result, "Semigroup concat should execute right-to-left: (5 + 1) * 2 = 12")

	// Test associativity: (f . g) . h = f . (g . h)
	f := double
	g := increment
	h := square

	left := sg.Concat(sg.Concat(f, g), h)
	right := sg.Concat(f, sg.Concat(g, h))

	testValue := 3
	assert.Equal(t, left(testValue), right(testValue), "Semigroup should be associative")

	// Test with ConcatAll from semigroup package (RIGHT-TO-LEFT)
	// ConcatAll(double)(increment, square) means: square, then increment, then double
	combined2 := S.ConcatAll(sg)(double)([]Endomorphism[int]{increment, square})
	result2 := combined2(5)
	// 5 -> square -> 25 -> increment -> 26 -> double -> 52
	assert.Equal(t, 52, result2, "Semigroup ConcatAll should execute right-to-left: square(5)=25, +1=26, *2=52")
}

// TestMonoid tests the Monoid function
func TestMonoid(t *testing.T) {
	monoid := Monoid[int]()

	// Test that empty is identity
	empty := monoid.Empty()
	assert.Equal(t, 42, empty(42), "Monoid empty should be identity")

	// Test right identity: x . empty = x (RIGHT-TO-LEFT: empty first, then x)
	// Concat(double, empty) means: empty first, then double
	rightIdentity := monoid.Concat(double, empty)
	assert.Equal(t, 10, rightIdentity(5), "Monoid should satisfy right identity: empty(5) then double = 10")

	// Test left identity: empty . x = x (RIGHT-TO-LEFT: x first, then empty)
	// Concat(empty, double) means: double first, then empty
	leftIdentity := monoid.Concat(empty, double)
	assert.Equal(t, 10, leftIdentity(5), "Monoid should satisfy left identity: double(5) then empty = 10")

	// Test ConcatAll with multiple endomorphisms (RIGHT-TO-LEFT execution)
	combined := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square})
	result := combined(5)
	// RIGHT-TO-LEFT: square(5) = 25, increment(25) = 26, double(26) = 52
	assert.Equal(t, 52, result, "Monoid ConcatAll should execute right-to-left: square(5)=25, +1=26, *2=52")

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
	// Test with strings (RIGHT-TO-LEFT execution)
	addExclamation := func(s string) string {
		return s + "!"
	}
	addPrefix := func(s string) string {
		return "Hello, " + s
	}

	// Compose(addExclamation, addPrefix) means: addPrefix first, then addExclamation
	strComposed := MonadCompose(addExclamation, addPrefix)
	result := strComposed("World")
	assert.Equal(t, "Hello, World!", result, "Compose should execute right-to-left with strings")

	// Test with float64 (RIGHT-TO-LEFT execution)
	doubleFloat := func(x float64) float64 {
		return x * 2.0
	}
	addOne := func(x float64) float64 {
		return x + 1.0
	}

	// Compose(doubleFloat, addOne) means: addOne first, then doubleFloat
	floatComposed := MonadCompose(doubleFloat, addOne)
	resultFloat := floatComposed(5.5)
	// 5.5 + 1.0 = 6.5, 6.5 * 2.0 = 13.0
	assert.Equal(t, 13.0, resultFloat, "Compose should execute right-to-left: (5.5 + 1.0) * 2.0 = 13.0")
}

// TestComplexCompositions tests more complex composition scenarios
func TestComplexCompositions(t *testing.T) {
	// Create a pipeline of transformations (RIGHT-TO-LEFT execution)
	// Innermost Compose is evaluated first in the composition chain
	pipeline := MonadCompose(
		MonadCompose(
			MonadCompose(double, increment),
			square,
		),
		negate,
	)

	// RIGHT-TO-LEFT: negate(5) = -5, square(-5) = 25, increment(25) = 26, double(26) = 52
	result := pipeline(5)
	assert.Equal(t, 52, result, "Complex composition should execute right-to-left")

	// Test using monoid to build the same pipeline (RIGHT-TO-LEFT)
	monoid := Monoid[int]()
	pipelineMonoid := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square, negate})
	resultMonoid := pipelineMonoid(5)
	// RIGHT-TO-LEFT: negate(5) = -5, square(-5) = 25, increment(25) = 26, double(26) = 52
	assert.Equal(t, 52, resultMonoid, "Monoid-based pipeline should match composition (right-to-left)")
}

// TestOperatorType tests the Operator type
func TestOperatorType(t *testing.T) {
	// Create an operator that transforms int endomorphisms
	// This operator takes an endomorphism and returns a new one that applies it twice
	applyTwice := func(f Endomorphism[int]) Endomorphism[int] {
		return func(x int) int {
			return f(f(x))
		}
	}

	// Use the operator
	var op Operator[int] = applyTwice
	doubleDouble := op(double)

	result := doubleDouble(5) // double(double(5)) = double(10) = 20
	assert.Equal(t, 20, result, "Operator should transform endomorphisms correctly")

	// Test with increment
	incrementTwice := op(increment)
	result2 := incrementTwice(5) // increment(increment(5)) = increment(6) = 7
	assert.Equal(t, 7, result2, "Operator should work with different endomorphisms")
}

// BenchmarkCompose benchmarks the Compose function
func BenchmarkCompose(b *testing.B) {
	composed := MonadCompose(double, increment)
	b.ResetTimer()
	for b.Loop() {
		_ = composed(5)
	}
}

// BenchmarkMonoidConcatAll benchmarks ConcatAll with monoid
// TestComposeVsChain demonstrates the key difference between Compose and Chain
func TestComposeVsChain(t *testing.T) {
	double := N.Mul(2)
	increment := N.Add(1)

	// Compose executes RIGHT-TO-LEFT
	// Compose(double, increment) means: increment first, then double
	composed := MonadCompose(double, increment)
	composedResult := composed(5) // (5 + 1) * 2 = 12

	// MonadChain executes LEFT-TO-RIGHT
	// MonadChain(double, increment) means: double first, then increment
	chained := MonadChain(double, increment)
	chainedResult := chained(5) // (5 * 2) + 1 = 11

	assert.Equal(t, 12, composedResult, "Compose should execute right-to-left")
	assert.Equal(t, 11, chainedResult, "MonadChain should execute left-to-right")
	assert.NotEqual(t, composedResult, chainedResult, "Compose and Chain should produce different results with non-commutative operations")

	// To get the same result with Compose, we need to reverse the order
	composedReversed := MonadCompose(increment, double)
	assert.Equal(t, chainedResult, composedReversed(5), "Compose with reversed args should match Chain")

	// Demonstrate with a more complex example
	square := func(x int) int { return x * x }

	// Compose: RIGHT-TO-LEFT
	composed3 := MonadCompose(MonadCompose(square, increment), double)
	// double(5) = 10, increment(10) = 11, square(11) = 121
	result1 := composed3(5)

	// MonadChain: LEFT-TO-RIGHT
	chained3 := MonadChain(MonadChain(double, increment), square)
	// double(5) = 10, increment(10) = 11, square(11) = 121
	result2 := chained3(5)

	assert.Equal(t, 121, result1, "Compose should execute right-to-left")
	assert.Equal(t, 121, result2, "MonadChain should execute left-to-right")
	assert.Equal(t, result1, result2, "Both should produce same result when operations are in correct order")
}

func BenchmarkMonoidConcatAll(b *testing.B) {
	monoid := Monoid[int]()
	combined := M.ConcatAll(monoid)([]Endomorphism[int]{double, increment, square})
	b.ResetTimer()
	for b.Loop() {
		_ = combined(5)
	}
}

// BenchmarkChain benchmarks the Chain function
func BenchmarkChain(b *testing.B) {
	chainWithIncrement := Chain(increment)
	chained := chainWithIncrement(double)
	b.ResetTimer()
	for b.Loop() {
		_ = chained(5)
	}
}

// TestFunctorLaws tests that endomorphisms satisfy the functor laws
func TestFunctorLaws(t *testing.T) {
	// Functor Law 1: Identity
	// map(id) = id
	t.Run("Identity", func(t *testing.T) {
		id := Identity[int]()
		endo := double

		// map(id)(endo) should equal endo
		mapped := MonadMap(id, endo)
		testValue := 5
		assert.Equal(t, endo(testValue), mapped(testValue), "map(id) should equal id")
	})

	// Functor Law 2: Composition
	// map(f . g) = map(f) . map(g)
	t.Run("Composition", func(t *testing.T) {
		f := double
		g := increment
		endo := square

		// Left side: map(f . g)(endo)
		composed := MonadCompose(f, g)
		left := MonadMap(composed, endo)

		// Right side: map(f)(map(g)(endo))
		mappedG := MonadMap(g, endo)
		right := MonadMap(f, mappedG)

		testValue := 3
		assert.Equal(t, left(testValue), right(testValue), "map(f . g) should equal map(f) . map(g)")
	})
}

// TestApplicativeLaws tests that endomorphisms satisfy the applicative functor laws
func TestApplicativeLaws(t *testing.T) {
	// Applicative Law 1: Identity
	// ap(id, v) = v
	t.Run("Identity", func(t *testing.T) {
		id := Identity[int]()
		v := double

		applied := MonadAp(id, v)
		testValue := 5
		assert.Equal(t, v(testValue), applied(testValue), "ap(id, v) should equal v")
	})

	// Applicative Law 2: Composition
	// ap(ap(ap(compose, u), v), w) = ap(u, ap(v, w))
	t.Run("Composition", func(t *testing.T) {
		u := double
		v := increment
		w := square

		// For endomorphisms, ap is just composition
		// Left side: ap(ap(ap(compose, u), v), w) = compose(compose(u, v), w)
		left := MonadCompose(MonadCompose(u, v), w)

		// Right side: ap(u, ap(v, w)) = compose(u, compose(v, w))
		right := MonadCompose(u, MonadCompose(v, w))

		testValue := 3
		assert.Equal(t, left(testValue), right(testValue), "Applicative composition law")
	})

	// Applicative Law 3: Homomorphism
	// ap(pure(f), pure(x)) = pure(f(x))
	t.Run("Homomorphism", func(t *testing.T) {
		// For endomorphisms, "pure" is just the identity function that returns a constant
		// This law is trivially satisfied for endomorphisms
		f := double
		x := 5

		// ap(f, id) applied to x should equal f(x)
		id := Identity[int]()
		applied := MonadAp(f, id)
		assert.Equal(t, f(x), applied(x), "Homomorphism law")
	})
}

// TestMonadLaws tests that endomorphisms satisfy the monad laws
func TestMonadLaws(t *testing.T) {
	// Monad Law 1: Left Identity
	// chain(pure(a), f) = f(a)
	t.Run("LeftIdentity", func(t *testing.T) {
		// For endomorphisms, "pure" is the identity function
		// chain(id, f) = f
		id := Identity[int]()
		f := double

		chained := MonadChain(id, f)
		testValue := 5
		assert.Equal(t, f(testValue), chained(testValue), "chain(id, f) should equal f")
	})

	// Monad Law 2: Right Identity
	// chain(m, pure) = m
	t.Run("RightIdentity", func(t *testing.T) {
		m := double
		id := Identity[int]()

		chained := MonadChain(m, id)
		testValue := 5
		assert.Equal(t, m(testValue), chained(testValue), "chain(m, id) should equal m")
	})

	// Monad Law 3: Associativity
	// chain(chain(m, f), g) = chain(m, x => chain(f(x), g))
	t.Run("Associativity", func(t *testing.T) {
		m := square
		f := double
		g := increment

		// Left side: chain(chain(m, f), g)
		left := MonadChain(MonadChain(m, f), g)

		// Right side: chain(m, chain(f, g))
		// For simple endomorphisms (not Kleisli arrows), this simplifies to:
		right := MonadChain(m, MonadChain(f, g))

		testValue := 3
		assert.Equal(t, left(testValue), right(testValue), "Monad associativity law")
	})
}

// TestMonadComposeVsMonadChain verifies the relationship between Compose and Chain
func TestMonadComposeVsMonadChain(t *testing.T) {
	f := double
	g := increment

	// MonadCompose(f, g) should equal MonadChain(g, f)
	// Because Compose is right-to-left and Chain is left-to-right
	composed := MonadCompose(f, g)
	chained := MonadChain(g, f)

	testValue := 5
	assert.Equal(t, composed(testValue), chained(testValue),
		"MonadCompose(f, g) should equal MonadChain(g, f)")
}

// TestMapEqualsCompose verifies that Map is equivalent to Compose for endomorphisms
func TestMapEqualsCompose(t *testing.T) {
	f := double
	g := increment

	// MonadMap(f, g) should equal MonadCompose(f, g)
	mapped := MonadMap(f, g)
	composed := MonadCompose(f, g)

	testValue := 5
	assert.Equal(t, composed(testValue), mapped(testValue),
		"MonadMap should equal MonadCompose for endomorphisms")

	// Curried versions
	mapF := Map(f)
	composeF := Compose(f)

	mappedG := mapF(g)
	composedG := composeF(g)

	assert.Equal(t, composedG(testValue), mappedG(testValue),
		"Map should equal Compose for endomorphisms (curried)")
}

// TestApEqualsCompose verifies that Ap is equivalent to Compose for endomorphisms
func TestApEqualsCompose(t *testing.T) {
	f := double
	g := increment

	// MonadAp(f, g) should equal MonadCompose(f, g)
	applied := MonadAp(f, g)
	composed := MonadCompose(f, g)

	testValue := 5
	assert.Equal(t, composed(testValue), applied(testValue),
		"MonadAp should equal MonadCompose for endomorphisms")

	// Curried versions
	apG := Ap(g)
	composeG := Compose(g)

	appliedF := apG(f)
	composedF := composeG(f)

	assert.Equal(t, composedF(testValue), appliedF(testValue),
		"Ap should equal Compose for endomorphisms (curried)")
}

// TestChainFirst tests the ChainFirst operation
func TestChainFirst(t *testing.T) {
	double := N.Mul(2)

	// Track side effect
	var sideEffect int
	logEffect := func(x int) int {
		sideEffect = x
		return x + 100 // This result should be discarded
	}

	chained := MonadChainFirst(double, logEffect)
	result := chained(5)

	// Should return double's result (10), not logEffect's result
	assert.Equal(t, 10, result, "ChainFirst should return first result")
	// But side effect should have been executed with double's result
	assert.Equal(t, 10, sideEffect, "ChainFirst should execute second function for effect")
}

// TestBuild tests the Build function
func TestBuild(t *testing.T) {
	t.Run("build with single transformation", func(t *testing.T) {
		// Build applies endomorphism to zero value
		result := Build(double)
		assert.Equal(t, 0, result, "Build(double) on zero value should be 0")
	})

	t.Run("build with composed transformations", func(t *testing.T) {
		// Create a builder that starts from zero and applies transformations
		builder := M.ConcatAll(Monoid[int]())([]Endomorphism[int]{
			N.Add(10),
			N.Mul(2),
			N.Add(5),
		})

		result := Build(builder)
		// RIGHT-TO-LEFT: 0 + 5 = 5, 5 * 2 = 10, 10 + 10 = 20
		assert.Equal(t, 20, result, "Build should apply composed transformations to zero value")
	})

	t.Run("build with identity", func(t *testing.T) {
		result := Build(Identity[int]())
		assert.Equal(t, 0, result, "Build(identity) should return zero value")
	})

	t.Run("build string from empty", func(t *testing.T) {
		builder := M.ConcatAll(Monoid[string]())([]Endomorphism[string]{
			func(s string) string { return s + "Hello" },
			func(s string) string { return s + " " },
			func(s string) string { return s + "World" },
		})

		result := Build(builder)
		// RIGHT-TO-LEFT: "" + "World" = "World", "World" + " " = "World ", "World " + "Hello" = "World Hello"
		assert.Equal(t, "World Hello", result, "Build should work with strings")
	})

	t.Run("build struct with builder pattern", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		withHost := func(host string) Endomorphism[Config] {
			return func(c Config) Config {
				c.Host = host
				return c
			}
		}

		withPort := func(port int) Endomorphism[Config] {
			return func(c Config) Config {
				c.Port = port
				return c
			}
		}

		builder := M.ConcatAll(Monoid[Config]())([]Endomorphism[Config]{
			withHost("localhost"),
			withPort(8080),
		})

		result := Build(builder)
		assert.Equal(t, "localhost", result.Host, "Build should set Host")
		assert.Equal(t, 8080, result.Port, "Build should set Port")
	})

	t.Run("build slice with operations", func(t *testing.T) {
		type IntSlice []int

		appendValue := func(v int) Endomorphism[IntSlice] {
			return func(s IntSlice) IntSlice {
				return append(s, v)
			}
		}

		builder := M.ConcatAll(Monoid[IntSlice]())([]Endomorphism[IntSlice]{
			appendValue(1),
			appendValue(2),
			appendValue(3),
		})

		result := Build(builder)
		// RIGHT-TO-LEFT: append 3, append 2, append 1
		assert.Equal(t, IntSlice{3, 2, 1}, result, "Build should construct slice")
	})
}

// TestBuildAsBuilderPattern demonstrates using Build as a builder pattern
func TestBuildAsBuilderPattern(t *testing.T) {
	type Person struct {
		Name   string
		Age    int
		Email  string
		Active bool
	}

	// Define builder functions
	withName := func(name string) Endomorphism[Person] {
		return func(p Person) Person {
			p.Name = name
			return p
		}
	}

	withAge := func(age int) Endomorphism[Person] {
		return func(p Person) Person {
			p.Age = age
			return p
		}
	}

	withEmail := func(email string) Endomorphism[Person] {
		return func(p Person) Person {
			p.Email = email
			return p
		}
	}

	withActive := func(active bool) Endomorphism[Person] {
		return func(p Person) Person {
			p.Active = active
			return p
		}
	}

	// Build a person using the builder pattern
	personBuilder := M.ConcatAll(Monoid[Person]())([]Endomorphism[Person]{
		withName("Alice"),
		withAge(30),
		withEmail("alice@example.com"),
		withActive(true),
	})

	person := Build(personBuilder)

	assert.Equal(t, "Alice", person.Name)
	assert.Equal(t, 30, person.Age)
	assert.Equal(t, "alice@example.com", person.Email)
	assert.True(t, person.Active)
}

// TestConcatAll tests the ConcatAll function
func TestConcatAll(t *testing.T) {
	t.Run("concat all with multiple endomorphisms", func(t *testing.T) {
		// ConcatAll executes RIGHT-TO-LEFT
		combined := ConcatAll([]Endomorphism[int]{double, increment, square})
		result := combined(5)
		// RIGHT-TO-LEFT: square(5) = 25, increment(25) = 26, double(26) = 52
		assert.Equal(t, 52, result, "ConcatAll should execute right-to-left")
	})

	t.Run("concat all with empty slice", func(t *testing.T) {
		// Empty slice should return identity
		identity := ConcatAll([]Endomorphism[int]{})
		result := identity(42)
		assert.Equal(t, 42, result, "ConcatAll with empty slice should return identity")
	})

	t.Run("concat all with single endomorphism", func(t *testing.T) {
		combined := ConcatAll([]Endomorphism[int]{double})
		result := combined(5)
		assert.Equal(t, 10, result, "ConcatAll with single endomorphism should apply it")
	})

	t.Run("concat all with two endomorphisms", func(t *testing.T) {
		// RIGHT-TO-LEFT: increment first, then double
		combined := ConcatAll([]Endomorphism[int]{double, increment})
		result := combined(5)
		assert.Equal(t, 12, result, "ConcatAll should execute right-to-left: (5 + 1) * 2 = 12")
	})

	t.Run("concat all with strings", func(t *testing.T) {
		appendHello := func(s string) string { return s + "Hello" }
		appendSpace := func(s string) string { return s + " " }
		appendWorld := func(s string) string { return s + "World" }

		// RIGHT-TO-LEFT execution
		combined := ConcatAll([]Endomorphism[string]{appendHello, appendSpace, appendWorld})
		result := combined("")
		// RIGHT-TO-LEFT: "" + "World" = "World", "World" + " " = "World ", "World " + "Hello" = "World Hello"
		assert.Equal(t, "World Hello", result, "ConcatAll should work with strings")
	})

	t.Run("concat all for building structs", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		withHost := func(host string) Endomorphism[Config] {
			return func(c Config) Config {
				c.Host = host
				return c
			}
		}

		withPort := func(port int) Endomorphism[Config] {
			return func(c Config) Config {
				c.Port = port
				return c
			}
		}

		combined := ConcatAll([]Endomorphism[Config]{
			withHost("localhost"),
			withPort(8080),
		})

		result := combined(Config{})
		assert.Equal(t, "localhost", result.Host)
		assert.Equal(t, 8080, result.Port)
	})

	t.Run("concat all is equivalent to monoid ConcatAll", func(t *testing.T) {
		endos := []Endomorphism[int]{double, increment, square}

		result1 := ConcatAll(endos)(5)
		result2 := M.ConcatAll(Monoid[int]())(endos)(5)

		assert.Equal(t, result1, result2, "ConcatAll should be equivalent to M.ConcatAll(Monoid())")
	})
}

// TestReduce tests the Reduce function
func TestReduce(t *testing.T) {
	t.Run("reduce with multiple endomorphisms", func(t *testing.T) {
		// Reduce executes LEFT-TO-RIGHT starting from zero value
		result := Reduce([]Endomorphism[int]{double, increment, square})
		// LEFT-TO-RIGHT: 0 -> double(0) = 0 -> increment(0) = 1 -> square(1) = 1
		assert.Equal(t, 1, result, "Reduce should execute left-to-right from zero value")
	})

	t.Run("reduce with empty slice", func(t *testing.T) {
		// Empty slice should return zero value
		result := Reduce([]Endomorphism[int]{})
		assert.Equal(t, 0, result, "Reduce with empty slice should return zero value")
	})

	t.Run("reduce with single endomorphism", func(t *testing.T) {
		addTen := N.Add(10)
		result := Reduce([]Endomorphism[int]{addTen})
		// 0 + 10 = 10
		assert.Equal(t, 10, result, "Reduce with single endomorphism should apply it to zero")
	})

	t.Run("reduce with sequential transformations", func(t *testing.T) {
		addTen := N.Add(10)
		// LEFT-TO-RIGHT: 0 -> addTen(0) = 10 -> double(10) = 20 -> increment(20) = 21
		result := Reduce([]Endomorphism[int]{addTen, double, increment})
		assert.Equal(t, 21, result, "Reduce should apply transformations left-to-right")
	})

	t.Run("reduce with strings", func(t *testing.T) {
		appendHello := func(s string) string { return s + "Hello" }
		appendSpace := func(s string) string { return s + " " }
		appendWorld := func(s string) string { return s + "World" }

		// LEFT-TO-RIGHT execution
		result := Reduce([]Endomorphism[string]{appendHello, appendSpace, appendWorld})
		// "" -> "Hello" -> "Hello " -> "Hello World"
		assert.Equal(t, "Hello World", result, "Reduce should work with strings left-to-right")
	})

	t.Run("reduce for building structs", func(t *testing.T) {
		type Settings struct {
			Theme    string
			FontSize int
		}

		withTheme := func(theme string) Endomorphism[Settings] {
			return func(s Settings) Settings {
				s.Theme = theme
				return s
			}
		}

		withFontSize := func(size int) Endomorphism[Settings] {
			return func(s Settings) Settings {
				s.FontSize = size
				return s
			}
		}

		// LEFT-TO-RIGHT application
		result := Reduce([]Endomorphism[Settings]{
			withTheme("dark"),
			withFontSize(14),
		})

		assert.Equal(t, "dark", result.Theme)
		assert.Equal(t, 14, result.FontSize)
	})

	t.Run("reduce is equivalent to Build(ConcatAll(reverse))", func(t *testing.T) {
		addTen := N.Add(10)
		endos := []Endomorphism[int]{addTen, double, increment}

		// Reduce applies left-to-right
		result1 := Reduce(endos)

		// Reverse and use ConcatAll (which is right-to-left)
		reversed := []Endomorphism[int]{increment, double, addTen}
		result2 := Build(ConcatAll(reversed))

		assert.Equal(t, result1, result2, "Reduce should be equivalent to Build(ConcatAll(reverse))")
	})
}

// TestConcatAllVsReduce demonstrates the difference between ConcatAll and Reduce
func TestConcatAllVsReduce(t *testing.T) {
	addTen := N.Add(10)

	endos := []Endomorphism[int]{addTen, double, increment}

	// ConcatAll: RIGHT-TO-LEFT composition, returns endomorphism
	concatResult := ConcatAll(endos)(5)
	// 5 -> increment(5) = 6 -> double(6) = 12 -> addTen(12) = 22

	// Reduce: LEFT-TO-RIGHT application, returns value from zero
	reduceResult := Reduce(endos)
	// 0 -> addTen(0) = 10 -> double(10) = 20 -> increment(20) = 21

	assert.NotEqual(t, concatResult, reduceResult, "ConcatAll and Reduce should produce different results")
	assert.Equal(t, 22, concatResult, "ConcatAll should execute right-to-left on input value")
	assert.Equal(t, 21, reduceResult, "Reduce should execute left-to-right from zero value")
}

// TestReduceWithBuild demonstrates using Reduce vs Build with ConcatAll
func TestReduceWithBuild(t *testing.T) {
	addFive := N.Add(5)
	multiplyByThree := N.Mul(3)

	endos := []Endomorphism[int]{addFive, multiplyByThree}

	// Reduce: LEFT-TO-RIGHT from zero
	reduceResult := Reduce(endos)
	// 0 -> addFive(0) = 5 -> multiplyByThree(5) = 15
	assert.Equal(t, 15, reduceResult)

	// Build with ConcatAll: RIGHT-TO-LEFT from zero
	buildResult := Build(ConcatAll(endos))
	// 0 -> multiplyByThree(0) = 0 -> addFive(0) = 5
	assert.Equal(t, 5, buildResult)

	assert.NotEqual(t, reduceResult, buildResult, "Reduce and Build(ConcatAll) produce different results due to execution order")
}
