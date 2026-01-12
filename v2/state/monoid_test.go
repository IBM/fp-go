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

package state

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Counter is a simple state type for testing
type Counter struct {
	count int
}

// Logger is a state type that accumulates log messages
type Logger struct {
	logs []string
}

// TestApplicativeMonoidBasic tests basic monoid operations
func TestApplicativeMonoidBasic(t *testing.T) {
	t.Run("integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](5)
		s2 := Of[Counter](3)

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 10})

		assert.Equal(t, Counter{count: 10}, pair.Head(result))
		assert.Equal(t, 8, pair.Tail(result))
	})

	t.Run("integer multiplication", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		stateMonoid := ApplicativeMonoid[Counter](intMulMonoid)

		s1 := Of[Counter](4)
		s2 := Of[Counter](5)

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})

		assert.Equal(t, Counter{count: 0}, pair.Head(result))
		assert.Equal(t, 20, pair.Tail(result))
	})

	t.Run("string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		stateMonoid := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("Hello")
		s2 := Of[Counter](" World")

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 5})

		assert.Equal(t, Counter{count: 5}, pair.Head(result))
		assert.Equal(t, "Hello World", pair.Tail(result))
	})

	t.Run("boolean AND", func(t *testing.T) {
		boolAndMonoid := M.MakeMonoid(func(a, b bool) bool { return a && b }, true)
		stateMonoid := ApplicativeMonoid[Counter](boolAndMonoid)

		s1 := Of[Counter](true)
		s2 := Of[Counter](true)

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})

		assert.Equal(t, true, pair.Tail(result))

		s3 := Of[Counter](false)
		combined2 := stateMonoid.Concat(s1, s3)
		result2 := combined2(Counter{count: 0})

		assert.Equal(t, false, pair.Tail(result2))
	})
}

// TestApplicativeMonoidEmpty tests the empty element
func TestApplicativeMonoidEmpty(t *testing.T) {
	t.Run("integer addition empty", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		empty := stateMonoid.Empty()
		result := empty(Counter{count: 5})

		assert.Equal(t, Counter{count: 5}, pair.Head(result))
		assert.Equal(t, 0, pair.Tail(result))
	})

	t.Run("integer multiplication empty", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		stateMonoid := ApplicativeMonoid[Counter](intMulMonoid)

		empty := stateMonoid.Empty()
		result := empty(Counter{count: 10})

		assert.Equal(t, Counter{count: 10}, pair.Head(result))
		assert.Equal(t, 1, pair.Tail(result))
	})

	t.Run("string concatenation empty", func(t *testing.T) {
		strMonoid := S.Monoid
		stateMonoid := ApplicativeMonoid[Counter](strMonoid)

		empty := stateMonoid.Empty()
		result := empty(Counter{count: 0})

		assert.Equal(t, Counter{count: 0}, pair.Head(result))
		assert.Equal(t, "", pair.Tail(result))
	})
}

// TestApplicativeMonoidLaws verifies monoid laws
func TestApplicativeMonoidLaws(t *testing.T) {
	t.Run("associativity with integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](1)
		s2 := Of[Counter](2)
		s3 := Of[Counter](3)

		initialState := Counter{count: 0}

		// (s1 • s2) • s3
		left := m.Concat(m.Concat(s1, s2), s3)
		leftResult := left(initialState)

		// s1 • (s2 • s3)
		right := m.Concat(s1, m.Concat(s2, s3))
		rightResult := right(initialState)

		assert.Equal(t, leftResult, rightResult)
		assert.Equal(t, 6, pair.Tail(leftResult))
	})

	t.Run("left identity with integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s := Of[Counter](42)
		initialState := Counter{count: 5}

		// Empty() • s = s
		result := m.Concat(m.Empty(), s)
		expected := s(initialState)
		actual := result(initialState)

		assert.Equal(t, expected, actual)
		assert.Equal(t, 42, pair.Tail(actual))
	})

	t.Run("right identity with integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s := Of[Counter](42)
		initialState := Counter{count: 5}

		// s • Empty() = s
		result := m.Concat(s, m.Empty())
		expected := s(initialState)
		actual := result(initialState)

		assert.Equal(t, expected, actual)
		assert.Equal(t, 42, pair.Tail(actual))
	})

	t.Run("associativity with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("a")
		s2 := Of[Counter]("b")
		s3 := Of[Counter]("c")

		initialState := Counter{count: 0}

		left := m.Concat(m.Concat(s1, s2), s3)
		leftResult := left(initialState)

		right := m.Concat(s1, m.Concat(s2, s3))
		rightResult := right(initialState)

		assert.Equal(t, leftResult, rightResult)
		assert.Equal(t, "abc", pair.Tail(leftResult))
	})

	t.Run("left identity with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s := Of[Counter]("test")
		initialState := Counter{count: 0}

		result := m.Concat(m.Empty(), s)
		expected := s(initialState)
		actual := result(initialState)

		assert.Equal(t, expected, actual)
		assert.Equal(t, "test", pair.Tail(actual))
	})

	t.Run("right identity with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s := Of[Counter]("test")
		initialState := Counter{count: 0}

		result := m.Concat(s, m.Empty())
		expected := s(initialState)
		actual := result(initialState)

		assert.Equal(t, expected, actual)
		assert.Equal(t, "test", pair.Tail(actual))
	})
}

// TestApplicativeMonoidWithStateModification tests monoid with state-modifying computations
func TestApplicativeMonoidWithStateModification(t *testing.T) {
	t.Run("state modification with integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		// Computation that increments counter and returns a value
		incrementAndReturn := func(val int) State[Counter, int] {
			return func(s Counter) Pair[Counter, int] {
				return pair.MakePair(Counter{count: s.count + 1}, val)
			}
		}

		s1 := incrementAndReturn(5)
		s2 := incrementAndReturn(3)

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})

		// State should be incremented twice
		assert.Equal(t, Counter{count: 2}, pair.Head(result))
		// Values should be added
		assert.Equal(t, 8, pair.Tail(result))
	})

	t.Run("state modification with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		stateMonoid := ApplicativeMonoid[Logger](strMonoid)

		// Computation that logs a message and returns it
		logMessage := func(msg string) State[Logger, string] {
			return func(s Logger) Pair[Logger, string] {
				newLogs := append(s.logs, msg)
				return pair.MakePair(Logger{logs: newLogs}, msg)
			}
		}

		s1 := logMessage("Hello")
		s2 := logMessage(" World")

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Logger{logs: []string{}})

		// Both messages should be logged
		assert.Equal(t, []string{"Hello", " World"}, pair.Head(result).logs)
		// Messages should be concatenated
		assert.Equal(t, "Hello World", pair.Tail(result))
	})

	t.Run("complex state transformation", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		// Computation that doubles the counter and returns the old value
		doubleAndReturnOld := func(val int) State[Counter, int] {
			return func(s Counter) Pair[Counter, int] {
				return pair.MakePair(Counter{count: s.count * 2}, val)
			}
		}

		s1 := doubleAndReturnOld(10)
		s2 := doubleAndReturnOld(20)

		combined := stateMonoid.Concat(s1, s2)
		result := combined(Counter{count: 3})

		// State: 3 -> 6 -> 12
		assert.Equal(t, Counter{count: 12}, pair.Head(result))
		// Values: 10 + 20 = 30
		assert.Equal(t, 30, pair.Tail(result))
	})
}

// TestApplicativeMonoidMultipleConcatenations tests chaining multiple concatenations
func TestApplicativeMonoidMultipleConcatenations(t *testing.T) {
	t.Run("chain of integer additions", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		states := []State[Counter, int]{
			Of[Counter](1),
			Of[Counter](2),
			Of[Counter](3),
			Of[Counter](4),
			Of[Counter](5),
		}

		// Fold all states using the monoid
		result := m.Empty()
		for _, s := range states {
			result = m.Concat(result, s)
		}

		finalResult := result(Counter{count: 0})
		assert.Equal(t, Counter{count: 0}, pair.Head(finalResult))
		assert.Equal(t, 15, pair.Tail(finalResult)) // 1+2+3+4+5
	})

	t.Run("chain of string concatenations", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		words := []string{"The", " quick", " brown", " fox"}
		states := make([]State[Counter, string], len(words))
		for i, word := range words {
			states[i] = Of[Counter](word)
		}

		result := m.Empty()
		for _, s := range states {
			result = m.Concat(result, s)
		}

		finalResult := result(Counter{count: 0})
		assert.Equal(t, "The quick brown fox", pair.Tail(finalResult))
	})

	t.Run("nested concatenations", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](1)
		s2 := Of[Counter](2)
		s3 := Of[Counter](3)
		s4 := Of[Counter](4)

		// ((s1 • s2) • s3) • s4
		result := m.Concat(
			m.Concat(
				m.Concat(s1, s2),
				s3,
			),
			s4,
		)

		finalResult := result(Counter{count: 0})
		assert.Equal(t, 10, pair.Tail(finalResult))
	})
}

// TestApplicativeMonoidEdgeCases tests edge cases
func TestApplicativeMonoidEdgeCases(t *testing.T) {
	t.Run("concatenating empty with empty", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		result := m.Concat(m.Empty(), m.Empty())
		finalResult := result(Counter{count: 5})

		assert.Equal(t, Counter{count: 5}, pair.Head(finalResult))
		assert.Equal(t, 0, pair.Tail(finalResult))
	})

	t.Run("zero values with multiplication", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		m := ApplicativeMonoid[Counter](intMulMonoid)

		s1 := Of[Counter](0)
		s2 := Of[Counter](42)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, 0, pair.Tail(finalResult))
	})

	t.Run("empty strings", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("")
		s2 := Of[Counter]("test")

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, "test", pair.Tail(finalResult))
	})

	t.Run("negative numbers with addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](-5)
		s2 := Of[Counter](3)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, -2, pair.Tail(finalResult))
	})
}

// TestApplicativeMonoidWithDifferentTypes tests monoid with various type combinations
func TestApplicativeMonoidWithDifferentTypes(t *testing.T) {
	t.Run("float64 addition", func(t *testing.T) {
		floatAddMonoid := N.MonoidSum[float64]()
		m := ApplicativeMonoid[Counter](floatAddMonoid)

		s1 := Of[Counter](1.5)
		s2 := Of[Counter](2.5)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, 4.0, pair.Tail(finalResult))
	})

	t.Run("float64 multiplication", func(t *testing.T) {
		floatMulMonoid := N.MonoidProduct[float64]()
		m := ApplicativeMonoid[Counter](floatMulMonoid)

		s1 := Of[Counter](2.0)
		s2 := Of[Counter](3.5)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, 7.0, pair.Tail(finalResult))
	})

	t.Run("boolean OR", func(t *testing.T) {
		boolOrMonoid := M.MakeMonoid(func(a, b bool) bool { return a || b }, false)
		m := ApplicativeMonoid[Counter](boolOrMonoid)

		s1 := Of[Counter](false)
		s2 := Of[Counter](true)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})

		assert.Equal(t, true, pair.Tail(finalResult))
	})
}

// TestApplicativeMonoidWithGets tests monoid with Gets operations
func TestApplicativeMonoidWithGets(t *testing.T) {
	t.Run("combining state reads", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		// Read the count from state
		getCount := Gets(func(c Counter) int { return c.count })

		// Combine two reads
		combined := m.Concat(getCount, getCount)
		result := combined(Counter{count: 5})

		// State unchanged
		assert.Equal(t, Counter{count: 5}, pair.Head(result))
		// Value is doubled (5 + 5)
		assert.Equal(t, 10, pair.Tail(result))
	})

	t.Run("combining different state projections", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		getCount := Gets(func(c Counter) int { return c.count })
		getDouble := Gets(func(c Counter) int { return c.count * 2 })

		combined := m.Concat(getCount, getDouble)
		result := combined(Counter{count: 3})

		assert.Equal(t, Counter{count: 3}, pair.Head(result))
		assert.Equal(t, 9, pair.Tail(result)) // 3 + 6
	})
}

// TestApplicativeMonoidCommutativity tests behavior with non-commutative operations
func TestApplicativeMonoidCommutativity(t *testing.T) {
	t.Run("string concatenation order matters", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("Hello")
		s2 := Of[Counter](" World")

		result1 := m.Concat(s1, s2)
		result2 := m.Concat(s2, s1)

		finalResult1 := result1(Counter{count: 0})
		finalResult2 := result2(Counter{count: 0})

		assert.Equal(t, "Hello World", pair.Tail(finalResult1))
		assert.Equal(t, " WorldHello", pair.Tail(finalResult2))
		assert.NotEqual(t, pair.Tail(finalResult1), pair.Tail(finalResult2))
	})

	t.Run("subtraction is not commutative", func(t *testing.T) {
		// Note: Subtraction doesn't form a proper monoid, but we can test it
		subMonoid := M.MakeMonoid(func(a, b int) int { return a - b }, 0)
		m := ApplicativeMonoid[Counter](subMonoid)

		s1 := Of[Counter](10)
		s2 := Of[Counter](3)

		result1 := m.Concat(s1, s2)
		result2 := m.Concat(s2, s1)

		finalResult1 := result1(Counter{count: 0})
		finalResult2 := result2(Counter{count: 0})

		assert.Equal(t, 7, pair.Tail(finalResult1))  // 10 - 3
		assert.Equal(t, -7, pair.Tail(finalResult2)) // 3 - 10
		assert.NotEqual(t, pair.Tail(finalResult1), pair.Tail(finalResult2))
	})
}
