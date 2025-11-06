// Copyright (c) 2025 IBM Corp.
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

package predicate

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// Test predicates for reuse
var (
	isPositive      = func(n int) bool { return n > 0 }
	isEven          = func(n int) bool { return n%2 == 0 }
	isNegative      = func(n int) bool { return n < 0 }
	isGreaterThan10 = func(n int) bool { return n > 10 }
)

// TestNot tests the Not function
func TestNot(t *testing.T) {
	t.Run("negates true to false", func(t *testing.T) {
		notPositive := Not(isPositive)
		assert.False(t, notPositive(5))
		assert.False(t, notPositive(1))
	})

	t.Run("negates false to true", func(t *testing.T) {
		notPositive := Not(isPositive)
		assert.True(t, notPositive(-5))
		assert.True(t, notPositive(0))
	})

	t.Run("double negation returns original", func(t *testing.T) {
		doubleNegated := Not(Not(isPositive))
		assert.True(t, doubleNegated(5))
		assert.False(t, doubleNegated(-5))
	})
}

// TestAnd tests the And function
func TestAnd(t *testing.T) {
	t.Run("returns true when both predicates are true", func(t *testing.T) {
		isPositiveAndEven := F.Pipe1(isPositive, And(isEven))
		assert.True(t, isPositiveAndEven(2))
		assert.True(t, isPositiveAndEven(4))
		assert.True(t, isPositiveAndEven(100))
	})

	t.Run("returns false when first predicate is false", func(t *testing.T) {
		isPositiveAndEven := F.Pipe1(isPositive, And(isEven))
		assert.False(t, isPositiveAndEven(-2))
		assert.False(t, isPositiveAndEven(-4))
	})

	t.Run("returns false when second predicate is false", func(t *testing.T) {
		isPositiveAndEven := F.Pipe1(isPositive, And(isEven))
		assert.False(t, isPositiveAndEven(1))
		assert.False(t, isPositiveAndEven(3))
		assert.False(t, isPositiveAndEven(5))
	})

	t.Run("returns false when both predicates are false", func(t *testing.T) {
		isPositiveAndEven := F.Pipe1(isPositive, And(isEven))
		assert.False(t, isPositiveAndEven(-1))
		assert.False(t, isPositiveAndEven(-3))
	})

	t.Run("chains multiple And operations", func(t *testing.T) {
		isPositiveEvenAndGreaterThan10 := F.Pipe2(
			isPositive,
			And(isEven),
			And(isGreaterThan10),
		)
		assert.True(t, isPositiveEvenAndGreaterThan10(12))
		assert.False(t, isPositiveEvenAndGreaterThan10(8))
		assert.False(t, isPositiveEvenAndGreaterThan10(11))
	})
}

// TestOr tests the Or function
func TestOr(t *testing.T) {
	t.Run("returns true when first predicate is true", func(t *testing.T) {
		isPositiveOrEven := F.Pipe1(isPositive, Or(isEven))
		assert.True(t, isPositiveOrEven(1))
		assert.True(t, isPositiveOrEven(3))
		assert.True(t, isPositiveOrEven(5))
	})

	t.Run("returns true when second predicate is true", func(t *testing.T) {
		isPositiveOrEven := F.Pipe1(isPositive, Or(isEven))
		assert.True(t, isPositiveOrEven(-2))
		assert.True(t, isPositiveOrEven(-4))
		assert.True(t, isPositiveOrEven(0))
	})

	t.Run("returns true when both predicates are true", func(t *testing.T) {
		isPositiveOrEven := F.Pipe1(isPositive, Or(isEven))
		assert.True(t, isPositiveOrEven(2))
		assert.True(t, isPositiveOrEven(4))
		assert.True(t, isPositiveOrEven(100))
	})

	t.Run("returns false when both predicates are false", func(t *testing.T) {
		isPositiveOrEven := F.Pipe1(isPositive, Or(isEven))
		assert.False(t, isPositiveOrEven(-1))
		assert.False(t, isPositiveOrEven(-3))
		assert.False(t, isPositiveOrEven(-5))
	})

	t.Run("chains multiple Or operations", func(t *testing.T) {
		isPositiveOrEvenOrNegative := F.Pipe2(
			isPositive,
			Or(isEven),
			Or(isNegative),
		)
		assert.True(t, isPositiveOrEvenOrNegative(5))  // positive
		assert.True(t, isPositiveOrEvenOrNegative(2))  // even
		assert.True(t, isPositiveOrEvenOrNegative(-3)) // negative
		assert.True(t, isPositiveOrEvenOrNegative(0))  // even
	})
}

// TestContraMap tests the ContraMap function
func TestContraMap(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	t.Run("transforms predicate to work with different type", func(t *testing.T) {
		isAdult := func(age int) bool { return age >= 18 }
		getAge := func(p Person) int { return p.Age }
		isPersonAdult := F.Pipe1(isAdult, ContraMap(getAge))

		assert.True(t, isPersonAdult(Person{Name: "Alice", Age: 25}))
		assert.True(t, isPersonAdult(Person{Name: "Bob", Age: 18}))
		assert.False(t, isPersonAdult(Person{Name: "Charlie", Age: 15}))
	})

	t.Run("works with string length", func(t *testing.T) {
		isLongEnough := func(n int) bool { return n >= 5 }
		getLength := func(s string) int { return len(s) }
		isStringLongEnough := F.Pipe1(isLongEnough, ContraMap(getLength))

		assert.True(t, isStringLongEnough("hello"))
		assert.True(t, isStringLongEnough("world!"))
		assert.False(t, isStringLongEnough("hi"))
		assert.False(t, isStringLongEnough(""))
	})

	t.Run("composes with other operations", func(t *testing.T) {
		type Product struct {
			Name  string
			Price int
		}

		isExpensive := func(price int) bool { return price > 100 }
		isCheap := func(price int) bool { return price < 50 }
		getPrice := func(p Product) int { return p.Price }

		isExpensiveProduct := F.Pipe1(isExpensive, ContraMap(getPrice))
		isCheapProduct := F.Pipe1(isCheap, ContraMap(getPrice))
		isExtremePrice := F.Pipe1(isExpensiveProduct, Or(isCheapProduct))

		assert.True(t, isExtremePrice(Product{Name: "Luxury", Price: 200}))
		assert.True(t, isExtremePrice(Product{Name: "Budget", Price: 30}))
		assert.False(t, isExtremePrice(Product{Name: "Mid-range", Price: 75}))
	})
}

// TestSemigroupAny tests the SemigroupAny function
func TestSemigroupAny(t *testing.T) {
	s := SemigroupAny[int]()

	t.Run("combines predicates with OR logic", func(t *testing.T) {
		combined := s.Concat(isPositive, isEven)
		assert.True(t, combined(4))   // both true
		assert.True(t, combined(3))   // first true
		assert.True(t, combined(-2))  // second true
		assert.False(t, combined(-3)) // both false
	})

	t.Run("is associative", func(t *testing.T) {
		// (a OR b) OR c == a OR (b OR c)
		left := s.Concat(s.Concat(isPositive, isEven), isNegative)
		right := s.Concat(isPositive, s.Concat(isEven, isNegative))

		testValues := []int{-5, -2, 0, 1, 2, 5}
		for _, v := range testValues {
			assert.Equal(t, left(v), right(v), "associativity failed for value %d", v)
		}
	})

	t.Run("combines multiple predicates", func(t *testing.T) {
		combined := s.Concat(s.Concat(isPositive, isEven), isGreaterThan10)
		assert.True(t, combined(15))  // positive and > 10
		assert.True(t, combined(2))   // even
		assert.True(t, combined(1))   // positive
		assert.False(t, combined(-3)) // none
	})
}

// TestSemigroupAll tests the SemigroupAll function
func TestSemigroupAll(t *testing.T) {
	s := SemigroupAll[int]()

	t.Run("combines predicates with AND logic", func(t *testing.T) {
		combined := s.Concat(isPositive, isEven)
		assert.True(t, combined(4))   // both true
		assert.False(t, combined(3))  // first true only
		assert.False(t, combined(-2)) // second true only
		assert.False(t, combined(-3)) // both false
	})

	t.Run("is associative", func(t *testing.T) {
		// (a AND b) AND c == a AND (b AND c)
		isLessThan100 := func(n int) bool { return n < 100 }
		left := s.Concat(s.Concat(isPositive, isEven), isLessThan100)
		right := s.Concat(isPositive, s.Concat(isEven, isLessThan100))

		testValues := []int{-5, -2, 0, 1, 2, 50, 150}
		for _, v := range testValues {
			assert.Equal(t, left(v), right(v), "associativity failed for value %d", v)
		}
	})

	t.Run("combines multiple predicates", func(t *testing.T) {
		combined := s.Concat(s.Concat(isPositive, isEven), isGreaterThan10)
		assert.True(t, combined(12))  // all true
		assert.False(t, combined(8))  // not > 10
		assert.False(t, combined(11)) // not even
		assert.False(t, combined(-2)) // not positive
	})
}

// TestMonoidAny tests the MonoidAny function
func TestMonoidAny(t *testing.T) {
	m := MonoidAny[int]()

	t.Run("has identity element that returns false", func(t *testing.T) {
		empty := m.Empty()
		assert.False(t, empty(0))
		assert.False(t, empty(5))
		assert.False(t, empty(-5))
	})

	t.Run("identity is left identity", func(t *testing.T) {
		// empty OR p == p
		combined := m.Concat(m.Empty(), isPositive)
		assert.True(t, combined(5))
		assert.False(t, combined(-5))
	})

	t.Run("identity is right identity", func(t *testing.T) {
		// p OR empty == p
		combined := m.Concat(isPositive, m.Empty())
		assert.True(t, combined(5))
		assert.False(t, combined(-5))
	})

	t.Run("reduces empty list to identity", func(t *testing.T) {
		predicates := []Predicate[int]{}
		result := m.Empty()
		for _, p := range predicates {
			result = m.Concat(result, p)
		}
		assert.False(t, result(5))
	})

	t.Run("reduces list of predicates", func(t *testing.T) {
		predicates := []Predicate[int]{isPositive, isEven, isGreaterThan10}
		result := m.Empty()
		for _, p := range predicates {
			result = m.Concat(result, p)
		}
		assert.True(t, result(15))  // positive
		assert.True(t, result(2))   // even
		assert.True(t, result(11))  // > 10
		assert.False(t, result(-3)) // none
	})
}

// TestMonoidAll tests the MonoidAll function
func TestMonoidAll(t *testing.T) {
	m := MonoidAll[int]()

	t.Run("has identity element that returns true", func(t *testing.T) {
		empty := m.Empty()
		assert.True(t, empty(0))
		assert.True(t, empty(5))
		assert.True(t, empty(-5))
	})

	t.Run("identity is left identity", func(t *testing.T) {
		// empty AND p == p
		combined := m.Concat(m.Empty(), isPositive)
		assert.True(t, combined(5))
		assert.False(t, combined(-5))
	})

	t.Run("identity is right identity", func(t *testing.T) {
		// p AND empty == p
		combined := m.Concat(isPositive, m.Empty())
		assert.True(t, combined(5))
		assert.False(t, combined(-5))
	})

	t.Run("reduces empty list to identity", func(t *testing.T) {
		predicates := []Predicate[int]{}
		result := m.Empty()
		for _, p := range predicates {
			result = m.Concat(result, p)
		}
		assert.True(t, result(5))
	})

	t.Run("reduces list of predicates", func(t *testing.T) {
		isLessThan100 := func(n int) bool { return n < 100 }
		predicates := []Predicate[int]{isPositive, isEven, isLessThan100}
		result := m.Empty()
		for _, p := range predicates {
			result = m.Concat(result, p)
		}
		assert.True(t, result(50))   // all true
		assert.False(t, result(51))  // not even
		assert.False(t, result(-2))  // not positive
		assert.False(t, result(150)) // not < 100
	})
}

// TestComplexScenarios tests complex combinations of predicates
func TestComplexScenarios(t *testing.T) {
	t.Run("complex boolean logic", func(t *testing.T) {
		// (positive AND even) OR (negative AND odd)
		positiveAndEven := F.Pipe1(isPositive, And(isEven))
		isOdd := Not(isEven)
		negativeAndOdd := F.Pipe1(isNegative, And(isOdd))
		complex := F.Pipe1(positiveAndEven, Or(negativeAndOdd))

		assert.True(t, complex(2))   // positive and even
		assert.True(t, complex(4))   // positive and even
		assert.True(t, complex(-1))  // negative and odd
		assert.True(t, complex(-3))  // negative and odd
		assert.False(t, complex(1))  // positive but odd
		assert.False(t, complex(-2)) // negative but even
		assert.False(t, complex(0))  // neither
	})

	t.Run("contramap with complex predicates", func(t *testing.T) {
		type User struct {
			Name  string
			Age   int
			Score int
		}

		isAdultAge := func(age int) bool { return age >= 18 }
		hasHighScore := func(score int) bool { return score >= 80 }

		getAge := func(u User) int { return u.Age }
		getScore := func(u User) int { return u.Score }

		isAdult := F.Pipe1(isAdultAge, ContraMap(getAge))
		hasGoodScore := F.Pipe1(hasHighScore, ContraMap(getScore))
		isQualified := F.Pipe1(isAdult, And(hasGoodScore))

		assert.True(t, isQualified(User{Name: "Alice", Age: 25, Score: 90}))
		assert.False(t, isQualified(User{Name: "Bob", Age: 16, Score: 90}))
		assert.False(t, isQualified(User{Name: "Charlie", Age: 25, Score: 70}))
		assert.False(t, isQualified(User{Name: "Dave", Age: 16, Score: 70}))
	})

	t.Run("monoid with contramap", func(t *testing.T) {
		type Item struct {
			Price int
			Stock int
		}

		m := MonoidAll[Item]()

		isAffordable := func(price int) bool { return price < 100 }
		isInStock := func(stock int) bool { return stock > 0 }

		getPrice := func(i Item) int { return i.Price }
		getStock := func(i Item) int { return i.Stock }

		isAffordableItem := F.Pipe1(isAffordable, ContraMap(getPrice))
		isInStockItem := F.Pipe1(isInStock, ContraMap(getStock))

		canBuy := m.Concat(isAffordableItem, isInStockItem)

		assert.True(t, canBuy(Item{Price: 50, Stock: 10}))
		assert.False(t, canBuy(Item{Price: 150, Stock: 10}))
		assert.False(t, canBuy(Item{Price: 50, Stock: 0}))
		assert.False(t, canBuy(Item{Price: 150, Stock: 0}))
	})
}
