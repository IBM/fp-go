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
	"testing"

	"github.com/IBM/fp-go/v2/io"
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
		stateIOMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](5)
		s2 := Of[Counter](3)

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 10})()

		assert.Equal(t, Counter{count: 10}, pair.Head(result))
		assert.Equal(t, 8, pair.Tail(result))
	})

	t.Run("integer multiplication", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		stateIOMonoid := ApplicativeMonoid[Counter](intMulMonoid)

		s1 := Of[Counter](4)
		s2 := Of[Counter](5)

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})()

		assert.Equal(t, Counter{count: 0}, pair.Head(result))
		assert.Equal(t, 20, pair.Tail(result))
	})

	t.Run("string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		stateIOMonoid := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("Hello")
		s2 := Of[Counter](" World")

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 5})()

		assert.Equal(t, Counter{count: 5}, pair.Head(result))
		assert.Equal(t, "Hello World", pair.Tail(result))
	})

	t.Run("boolean AND", func(t *testing.T) {
		boolAndMonoid := M.MakeMonoid(func(a, b bool) bool { return a && b }, true)
		stateIOMonoid := ApplicativeMonoid[Counter](boolAndMonoid)

		s1 := Of[Counter](true)
		s2 := Of[Counter](true)

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})()

		assert.Equal(t, true, pair.Tail(result))

		s3 := Of[Counter](false)
		combined2 := stateIOMonoid.Concat(s1, s3)
		result2 := combined2(Counter{count: 0})()

		assert.Equal(t, false, pair.Tail(result2))
	})
}

// TestApplicativeMonoidEmpty tests the empty element
func TestApplicativeMonoidEmpty(t *testing.T) {
	t.Run("integer addition empty", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateIOMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		empty := stateIOMonoid.Empty()
		result := empty(Counter{count: 5})()

		assert.Equal(t, Counter{count: 5}, pair.Head(result))
		assert.Equal(t, 0, pair.Tail(result))
	})

	t.Run("integer multiplication empty", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		stateIOMonoid := ApplicativeMonoid[Counter](intMulMonoid)

		empty := stateIOMonoid.Empty()
		result := empty(Counter{count: 10})()

		assert.Equal(t, Counter{count: 10}, pair.Head(result))
		assert.Equal(t, 1, pair.Tail(result))
	})

	t.Run("string concatenation empty", func(t *testing.T) {
		strMonoid := S.Monoid
		stateIOMonoid := ApplicativeMonoid[Counter](strMonoid)

		empty := stateIOMonoid.Empty()
		result := empty(Counter{count: 0})()

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
		leftResult := left(initialState)()

		// s1 • (s2 • s3)
		right := m.Concat(s1, m.Concat(s2, s3))
		rightResult := right(initialState)()

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
		expected := s(initialState)()
		actual := result(initialState)()

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
		expected := s(initialState)()
		actual := result(initialState)()

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
		leftResult := left(initialState)()

		right := m.Concat(s1, m.Concat(s2, s3))
		rightResult := right(initialState)()

		assert.Equal(t, leftResult, rightResult)
		assert.Equal(t, "abc", pair.Tail(leftResult))
	})

	t.Run("left identity with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s := Of[Counter]("test")
		initialState := Counter{count: 0}

		result := m.Concat(m.Empty(), s)
		expected := s(initialState)()
		actual := result(initialState)()

		assert.Equal(t, expected, actual)
		assert.Equal(t, "test", pair.Tail(actual))
	})

	t.Run("right identity with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s := Of[Counter]("test")
		initialState := Counter{count: 0}

		result := m.Concat(s, m.Empty())
		expected := s(initialState)()
		actual := result(initialState)()

		assert.Equal(t, expected, actual)
		assert.Equal(t, "test", pair.Tail(actual))
	})
}

// TestApplicativeMonoidWithStateModification tests monoid with state-modifying computations
func TestApplicativeMonoidWithStateModification(t *testing.T) {
	t.Run("state modification with integer addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateIOMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		// Computation that increments counter and returns a value
		incrementAndReturn := func(val int) StateIO[Counter, int] {
			return func(s Counter) IO[Pair[Counter, int]] {
				return func() Pair[Counter, int] {
					return pair.MakePair(Counter{count: s.count + 1}, val)
				}
			}
		}

		s1 := incrementAndReturn(5)
		s2 := incrementAndReturn(3)

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 0})()

		// State should be incremented twice
		assert.Equal(t, Counter{count: 2}, pair.Head(result))
		// Values should be added
		assert.Equal(t, 8, pair.Tail(result))
	})

	t.Run("state modification with string concatenation", func(t *testing.T) {
		strMonoid := S.Monoid
		stateIOMonoid := ApplicativeMonoid[Logger](strMonoid)

		// Computation that logs a message and returns it
		logMessage := func(msg string) StateIO[Logger, string] {
			return func(s Logger) IO[Pair[Logger, string]] {
				return func() Pair[Logger, string] {
					newLogs := append(s.logs, msg)
					return pair.MakePair(Logger{logs: newLogs}, msg)
				}
			}
		}

		s1 := logMessage("Hello")
		s2 := logMessage(" World")

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Logger{logs: []string{}})()

		// Both messages should be logged
		assert.Equal(t, []string{"Hello", " World"}, pair.Head(result).logs)
		// Messages should be concatenated
		assert.Equal(t, "Hello World", pair.Tail(result))
	})

	t.Run("complex state transformation", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		stateIOMonoid := ApplicativeMonoid[Counter](intAddMonoid)

		// Computation that doubles the counter and returns the old value
		doubleAndReturnOld := func(val int) StateIO[Counter, int] {
			return func(s Counter) IO[Pair[Counter, int]] {
				return func() Pair[Counter, int] {
					return pair.MakePair(Counter{count: s.count * 2}, val)
				}
			}
		}

		s1 := doubleAndReturnOld(10)
		s2 := doubleAndReturnOld(20)

		combined := stateIOMonoid.Concat(s1, s2)
		result := combined(Counter{count: 3})()

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

		states := []StateIO[Counter, int]{
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

		finalResult := result(Counter{count: 0})()
		assert.Equal(t, Counter{count: 0}, pair.Head(finalResult))
		assert.Equal(t, 15, pair.Tail(finalResult)) // 1+2+3+4+5
	})

	t.Run("chain of string concatenations", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		words := []string{"The", " quick", " brown", " fox"}
		states := make([]StateIO[Counter, string], len(words))
		for i, word := range words {
			states[i] = Of[Counter](word)
		}

		result := m.Empty()
		for _, s := range states {
			result = m.Concat(result, s)
		}

		finalResult := result(Counter{count: 0})()
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

		finalResult := result(Counter{count: 0})()
		assert.Equal(t, 10, pair.Tail(finalResult))
	})
}

// TestApplicativeMonoidEdgeCases tests edge cases
func TestApplicativeMonoidEdgeCases(t *testing.T) {
	t.Run("concatenating empty with empty", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		result := m.Concat(m.Empty(), m.Empty())
		finalResult := result(Counter{count: 5})()

		assert.Equal(t, Counter{count: 5}, pair.Head(finalResult))
		assert.Equal(t, 0, pair.Tail(finalResult))
	})

	t.Run("zero values with multiplication", func(t *testing.T) {
		intMulMonoid := N.MonoidProduct[int]()
		m := ApplicativeMonoid[Counter](intMulMonoid)

		s1 := Of[Counter](0)
		s2 := Of[Counter](42)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})()

		assert.Equal(t, 0, pair.Tail(finalResult))
	})

	t.Run("empty strings", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Counter](strMonoid)

		s1 := Of[Counter]("")
		s2 := Of[Counter]("test")

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})()

		assert.Equal(t, "test", pair.Tail(finalResult))
	})

	t.Run("negative numbers with addition", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		s1 := Of[Counter](-5)
		s2 := Of[Counter](3)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})()

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
		finalResult := result(Counter{count: 0})()

		assert.Equal(t, 4.0, pair.Tail(finalResult))
	})

	t.Run("float64 multiplication", func(t *testing.T) {
		floatMulMonoid := N.MonoidProduct[float64]()
		m := ApplicativeMonoid[Counter](floatMulMonoid)

		s1 := Of[Counter](2.0)
		s2 := Of[Counter](3.5)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})()

		assert.Equal(t, 7.0, pair.Tail(finalResult))
	})

	t.Run("boolean OR", func(t *testing.T) {
		boolOrMonoid := M.MakeMonoid(func(a, b bool) bool { return a || b }, false)
		m := ApplicativeMonoid[Counter](boolOrMonoid)

		s1 := Of[Counter](false)
		s2 := Of[Counter](true)

		result := m.Concat(s1, s2)
		finalResult := result(Counter{count: 0})()

		assert.Equal(t, true, pair.Tail(finalResult))
	})
}

// TestApplicativeMonoidWithIO tests monoid with actual IO effects
func TestApplicativeMonoidWithIO(t *testing.T) {
	t.Run("combining IO effects", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		// StateIO that performs IO and modifies state
		effectfulComputation := func(val int, increment int) StateIO[Counter, int] {
			return func(s Counter) IO[Pair[Counter, int]] {
				return func() Pair[Counter, int] {
					// Simulate IO effect
					newCount := s.count + increment
					return pair.MakePair(Counter{count: newCount}, val)
				}
			}
		}

		s1 := effectfulComputation(10, 1)
		s2 := effectfulComputation(20, 2)

		combined := m.Concat(s1, s2)
		result := combined(Counter{count: 0})()

		// State should be incremented by 1 then by 2
		assert.Equal(t, Counter{count: 3}, pair.Head(result))
		// Values should be added
		assert.Equal(t, 30, pair.Tail(result))
	})

	t.Run("IO effects with logging", func(t *testing.T) {
		strMonoid := S.Monoid
		m := ApplicativeMonoid[Logger](strMonoid)

		// StateIO that logs and returns a message
		logAndReturn := func(msg string) StateIO[Logger, string] {
			return func(s Logger) IO[Pair[Logger, string]] {
				return func() Pair[Logger, string] {
					// Simulate IO logging effect
					newLogs := append(append([]string{}, s.logs...), msg)
					return pair.MakePair(Logger{logs: newLogs}, msg)
				}
			}
		}

		s1 := logAndReturn("First")
		s2 := logAndReturn("Second")
		s3 := logAndReturn("Third")

		combined := m.Concat(m.Concat(s1, s2), s3)
		result := combined(Logger{logs: []string{}})()

		assert.Equal(t, []string{"First", "Second", "Third"}, pair.Head(result).logs)
		assert.Equal(t, "FirstSecondThird", pair.Tail(result))
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

		finalResult1 := result1(Counter{count: 0})()
		finalResult2 := result2(Counter{count: 0})()

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

		finalResult1 := result1(Counter{count: 0})()
		finalResult2 := result2(Counter{count: 0})()

		assert.Equal(t, 7, pair.Tail(finalResult1))  // 10 - 3
		assert.Equal(t, -7, pair.Tail(finalResult2)) // 3 - 10
		assert.NotEqual(t, pair.Tail(finalResult1), pair.Tail(finalResult2))
	})

	t.Run("state modification order matters", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		// Computation that uses current state value
		useState := func(multiplier int) StateIO[Counter, int] {
			return func(s Counter) IO[Pair[Counter, int]] {
				return func() Pair[Counter, int] {
					result := s.count * multiplier
					newState := Counter{count: s.count + 1}
					return pair.MakePair(newState, result)
				}
			}
		}

		s1 := useState(2)
		s2 := useState(3)

		// Order matters because state is threaded through
		result1 := m.Concat(s1, s2)
		result2 := m.Concat(s2, s1)

		finalResult1 := result1(Counter{count: 5})()
		finalResult2 := result2(Counter{count: 5})()

		// s1 then s2: (5*2) + ((5+1)*3) = 10 + 18 = 28
		assert.Equal(t, 28, pair.Tail(finalResult1))
		// s2 then s1: (5*3) + ((5+1)*2) = 15 + 12 = 27
		assert.Equal(t, 27, pair.Tail(finalResult2))
		assert.NotEqual(t, pair.Tail(finalResult1), pair.Tail(finalResult2))
	})
}

// TestApplicativeMonoidWithFromIO tests monoid with FromIO
func TestApplicativeMonoidWithFromIO(t *testing.T) {
	t.Run("combining FromIO computations", func(t *testing.T) {
		intAddMonoid := N.MonoidSum[int]()
		m := ApplicativeMonoid[Counter](intAddMonoid)

		// Create StateIO from IO
		io1 := io.Of(5)
		io2 := io.Of(3)

		s1 := FromIO[Counter](io1)
		s2 := FromIO[Counter](io2)

		combined := m.Concat(s1, s2)
		result := combined(Counter{count: 10})()

		assert.Equal(t, Counter{count: 10}, pair.Head(result))
		assert.Equal(t, 8, pair.Tail(result))
	})
}
