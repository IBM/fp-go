// Copyright (c) 2024 - 2025 IBM Corp.
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

package stateio

import (
	"fmt"
	"testing"

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

type MonadTestState struct {
	Value int
}

// Test Left Identity Law: Of(a) >>= f ≡ f(a)
func TestMonadLeftIdentity(t *testing.T) {
	initial := MonadTestState{Value: 0}
	a := 42

	f := func(x int) StateIO[MonadTestState, string] {
		return func(s MonadTestState) IO[Pair[MonadTestState, string]] {
			return func() Pair[MonadTestState, string] {
				newState := MonadTestState{Value: s.Value + x}
				return pair.MakePair(newState, fmt.Sprintf("%d", x))
			}
		}
	}

	// Left side: Of(a) >>= f
	left := MonadChain(Of[MonadTestState](a), f)
	leftResult := left(initial)()

	// Right side: f(a)
	right := f(a)
	rightResult := right(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test Right Identity Law: m >>= Of ≡ m
func TestMonadRightIdentity(t *testing.T) {
	initial := MonadTestState{Value: 10}

	m := func(s MonadTestState) IO[Pair[MonadTestState, int]] {
		return func() Pair[MonadTestState, int] {
			newState := MonadTestState{Value: s.Value * 2}
			return pair.MakePair(newState, newState.Value)
		}
	}

	// Left side: m >>= Of
	left := MonadChain(m, func(x int) StateIO[MonadTestState, int] {
		return Of[MonadTestState](x)
	})
	leftResult := left(initial)()

	// Right side: m
	rightResult := m(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test Associativity Law: (m >>= f) >>= g ≡ m >>= (x => f(x) >>= g)
func TestMonadAssociativity(t *testing.T) {
	initial := MonadTestState{Value: 5}

	m := Of[MonadTestState](10)

	f := func(x int) StateIO[MonadTestState, int] {
		return func(s MonadTestState) IO[Pair[MonadTestState, int]] {
			return func() Pair[MonadTestState, int] {
				newState := MonadTestState{Value: s.Value + x}
				return pair.MakePair(newState, x*2)
			}
		}
	}

	g := func(y int) StateIO[MonadTestState, string] {
		return func(s MonadTestState) IO[Pair[MonadTestState, string]] {
			return func() Pair[MonadTestState, string] {
				newState := MonadTestState{Value: s.Value + y}
				return pair.MakePair(newState, fmt.Sprintf("%d", y))
			}
		}
	}

	// Left side: (m >>= f) >>= g
	left := MonadChain(MonadChain(m, f), g)
	leftResult := left(initial)()

	// Right side: m >>= (x => f(x) >>= g)
	right := MonadChain(m, func(x int) StateIO[MonadTestState, string] {
		return MonadChain(f(x), g)
	})
	rightResult := right(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test Functor Identity Law: Map(id) ≡ id
func TestFunctorIdentity(t *testing.T) {
	initial := MonadTestState{Value: 7}

	m := Of[MonadTestState](42)

	// Map with identity function
	mapped := MonadMap(m, F.Identity[int])
	mappedResult := mapped(initial)()

	// Original computation
	originalResult := m(initial)()

	assert.Equal(t, pair.Tail(originalResult), pair.Tail(mappedResult))
	assert.Equal(t, pair.Head(originalResult).Value, pair.Head(mappedResult).Value)
}

// Test Functor Composition Law: Map(f . g) ≡ Map(f) . Map(g)
func TestFunctorComposition(t *testing.T) {
	initial := MonadTestState{Value: 3}

	m := Of[MonadTestState](10)

	f := func(x int) int { return x * 2 }
	g := func(x int) int { return x + 5 }

	// Left side: Map(f . g)
	left := MonadMap(m, F.Flow2(g, f))
	leftResult := left(initial)()

	// Right side: Map(f) . Map(g)
	right := F.Pipe1(m, F.Flow2(Map[MonadTestState](g), Map[MonadTestState](f)))
	rightResult := right(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test Applicative Identity Law: Ap(Of(id), v) ≡ v
func TestApplicativeIdentity(t *testing.T) {
	initial := MonadTestState{Value: 1}

	v := Of[MonadTestState](42)

	// Ap(Of(id), v)
	applied := MonadAp(Of[MonadTestState](F.Identity[int]), v)
	appliedResult := applied(initial)()

	// v
	originalResult := v(initial)()

	assert.Equal(t, pair.Tail(originalResult), pair.Tail(appliedResult))
	assert.Equal(t, pair.Head(originalResult).Value, pair.Head(appliedResult).Value)
}

// Test Applicative Homomorphism Law: Ap(Of(f), Of(x)) ≡ Of(f(x))
func TestApplicativeHomomorphism(t *testing.T) {
	initial := MonadTestState{Value: 2}

	f := func(x int) int { return x * 3 }
	x := 7

	// Left side: Ap(Of(f), Of(x))
	left := MonadAp(Of[MonadTestState](f), Of[MonadTestState](x))
	leftResult := left(initial)()

	// Right side: Of(f(x))
	right := Of[MonadTestState](f(x))
	rightResult := right(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test Applicative Interchange Law: Ap(u, Of(y)) ≡ Ap(Of(f => f(y)), u)
func TestApplicativeInterchange(t *testing.T) {
	initial := MonadTestState{Value: 4}

	u := Of[MonadTestState](func(x int) int { return x + 10 })
	y := 5

	// Left side: Ap(u, Of(y))
	left := MonadAp(u, Of[MonadTestState](y))
	leftResult := left(initial)()

	// Right side: Ap(Of(f => f(y)), u)
	right := MonadAp(
		Of[MonadTestState](func(f func(int) int) int { return f(y) }),
		u,
	)
	rightResult := right(initial)()

	assert.Equal(t, pair.Tail(rightResult), pair.Tail(leftResult))
	assert.Equal(t, pair.Head(rightResult).Value, pair.Head(leftResult).Value)
}

// Test that StateIO implements Pointed interface
func TestPointed(t *testing.T) {
	pointed := Pointed[MonadTestState, int]()
	computation := pointed.Of(42)

	initial := MonadTestState{Value: 0}
	result := computation(initial)()

	assert.Equal(t, 42, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

// Test that StateIO implements Functor interface
func TestFunctor(t *testing.T) {
	functor := Functor[MonadTestState, int, string]()

	computation := Of[MonadTestState](42)
	mapped := functor.Map(func(x int) string { return fmt.Sprintf("%d", x) })(computation)

	initial := MonadTestState{Value: 0}
	result := mapped(initial)()

	assert.Equal(t, "42", pair.Tail(result))
}

// Test that StateIO implements Applicative interface
func TestApplicative(t *testing.T) {
	applicative := Applicative[MonadTestState, int, string]()

	fab := Of[MonadTestState](func(x int) string { return fmt.Sprintf("%d", x) })
	fa := Of[MonadTestState](42)
	result := applicative.Ap(fa)(fab)

	initial := MonadTestState{Value: 0}
	output := result(initial)()

	assert.Equal(t, "42", pair.Tail(output))
}

// Test that StateIO implements Monad interface
func TestMonad(t *testing.T) {
	monad := Monad[MonadTestState, int, string]()

	computation := monad.Of(42)
	chained := monad.Chain(func(x int) StateIO[MonadTestState, string] {
		return Of[MonadTestState](fmt.Sprintf("%d", x))
	})(computation)

	initial := MonadTestState{Value: 0}
	result := chained(initial)()

	assert.Equal(t, "42", pair.Tail(result))
}

// Test Eq functionality
func TestEq(t *testing.T) {
	initial := MonadTestState{Value: 0}

	comp1 := Of[MonadTestState](42)
	comp2 := Of[MonadTestState](42)
	comp3 := Of[MonadTestState](43)

	// Create equality predicate for IO[Pair[MonadTestState, int]]
	eqIO := EQ.FromEquals(func(l, r IO[Pair[MonadTestState, int]]) bool {
		lResult := l()
		rResult := r()
		return pair.Tail(lResult) == pair.Tail(rResult) &&
			pair.Head(lResult).Value == pair.Head(rResult).Value
	})

	eq := Eq(eqIO)(initial)

	assert.True(t, eq.Equals(comp1, comp2))
	assert.False(t, eq.Equals(comp1, comp3))
}

// Test FromStrictEquals
func TestFromStrictEquals(t *testing.T) {
	initial := MonadTestState{Value: 0}

	comp1 := Of[MonadTestState](42)
	comp2 := Of[MonadTestState](42)
	comp3 := Of[MonadTestState](43)

	eq := FromStrictEquals[MonadTestState, int]()(initial)

	assert.True(t, eq.Equals(comp1, comp2))
	assert.False(t, eq.Equals(comp1, comp3))
}
