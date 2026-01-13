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

package either

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFirstMonoid tests the FirstMonoid implementation
func TestFirstMonoid(t *testing.T) {
	zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
	m := FirstMonoid(zero)

	t.Run("both Right values - returns first", func(t *testing.T) {
		result := m.Concat(Right[error](2), Right[error](3))
		assert.Equal(t, Right[error](2), result)
	})

	t.Run("left Right, right Left", func(t *testing.T) {
		result := m.Concat(Right[error](2), Left[int](errors.New("err")))
		assert.Equal(t, Right[error](2), result)
	})

	t.Run("left Left, right Right", func(t *testing.T) {
		result := m.Concat(Left[int](errors.New("err")), Right[error](3))
		assert.Equal(t, Right[error](3), result)
	})

	t.Run("both Left", func(t *testing.T) {
		err1 := errors.New("err1")
		err2 := errors.New("err2")
		result := m.Concat(Left[int](err1), Left[int](err2))
		// Should return the second Left
		assert.True(t, IsLeft(result))
		_, leftErr := Unwrap(result)
		assert.Equal(t, err2, leftErr)
	})

	t.Run("empty value", func(t *testing.T) {
		empty := m.Empty()
		assert.True(t, IsLeft(empty))
		_, leftErr := Unwrap(empty)
		assert.Equal(t, "empty", leftErr.Error())
	})

	t.Run("left identity", func(t *testing.T) {
		x := Right[error](5)
		result := m.Concat(m.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Right[error](5)
		result := m.Concat(x, m.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("associativity", func(t *testing.T) {
		a := Right[error](1)
		b := Right[error](2)
		c := Right[error](3)

		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Right[error](1), left)
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		// Should return the first Right value encountered
		result := m.Concat(
			m.Concat(Left[int](errors.New("err1")), Right[error](1)),
			m.Concat(Right[error](2), Right[error](3)),
		)
		assert.Equal(t, Right[error](1), result)
	})

	t.Run("with strings", func(t *testing.T) {
		zeroStr := func() Either[error, string] { return Left[string](errors.New("empty")) }
		strMonoid := FirstMonoid(zeroStr)

		result := strMonoid.Concat(Right[error]("first"), Right[error]("second"))
		assert.Equal(t, Right[error]("first"), result)

		result = strMonoid.Concat(Left[string](errors.New("err")), Right[error]("second"))
		assert.Equal(t, Right[error]("second"), result)
	})
}

// TestLastMonoid tests the LastMonoid implementation
func TestLastMonoid(t *testing.T) {
	zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
	m := LastMonoid(zero)

	t.Run("both Right values - returns last", func(t *testing.T) {
		result := m.Concat(Right[error](2), Right[error](3))
		assert.Equal(t, Right[error](3), result)
	})

	t.Run("left Right, right Left", func(t *testing.T) {
		result := m.Concat(Right[error](2), Left[int](errors.New("err")))
		assert.Equal(t, Right[error](2), result)
	})

	t.Run("left Left, right Right", func(t *testing.T) {
		result := m.Concat(Left[int](errors.New("err")), Right[error](3))
		assert.Equal(t, Right[error](3), result)
	})

	t.Run("both Left", func(t *testing.T) {
		err1 := errors.New("err1")
		err2 := errors.New("err2")
		result := m.Concat(Left[int](err1), Left[int](err2))
		// Should return the first Left
		assert.True(t, IsLeft(result))
		_, leftErr := Unwrap(result)
		assert.Equal(t, err1, leftErr)
	})

	t.Run("empty value", func(t *testing.T) {
		empty := m.Empty()
		assert.True(t, IsLeft(empty))
		_, leftErr := Unwrap(empty)
		assert.Equal(t, "empty", leftErr.Error())
	})

	t.Run("left identity", func(t *testing.T) {
		x := Right[error](5)
		result := m.Concat(m.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Right[error](5)
		result := m.Concat(x, m.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("associativity", func(t *testing.T) {
		a := Right[error](1)
		b := Right[error](2)
		c := Right[error](3)

		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Right[error](3), left)
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		// Should return the last Right value encountered
		result := m.Concat(
			m.Concat(Right[error](1), Right[error](2)),
			m.Concat(Right[error](3), Left[int](errors.New("err"))),
		)
		assert.Equal(t, Right[error](3), result)
	})

	t.Run("with strings", func(t *testing.T) {
		zeroStr := func() Either[error, string] { return Left[string](errors.New("empty")) }
		strMonoid := LastMonoid(zeroStr)

		result := strMonoid.Concat(Right[error]("first"), Right[error]("second"))
		assert.Equal(t, Right[error]("second"), result)

		result = strMonoid.Concat(Right[error]("first"), Left[string](errors.New("err")))
		assert.Equal(t, Right[error]("first"), result)
	})
}

// TestFirstMonoidVsAltMonoid verifies FirstMonoid and AltMonoid have the same behavior
func TestFirstMonoidVsAltMonoid(t *testing.T) {
	zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
	firstMonoid := FirstMonoid(zero)
	altMonoid := AltMonoid(zero)

	testCases := []struct {
		name  string
		left  Either[error, int]
		right Either[error, int]
	}{
		{"both Right", Right[error](1), Right[error](2)},
		{"left Right, right Left", Right[error](1), Left[int](errors.New("err"))},
		{"left Left, right Right", Left[int](errors.New("err")), Right[error](2)},
		{"both Left", Left[int](errors.New("err1")), Left[int](errors.New("err2"))},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			firstResult := firstMonoid.Concat(tc.left, tc.right)
			altResult := altMonoid.Concat(tc.left, tc.right)

			// Both should have the same Right/Left status
			assert.Equal(t, IsRight(firstResult), IsRight(altResult), "FirstMonoid and AltMonoid should have same Right/Left status")

			if IsRight(firstResult) {
				rightVal1, _ := Unwrap(firstResult)
				rightVal2, _ := Unwrap(altResult)
				assert.Equal(t, rightVal1, rightVal2, "FirstMonoid and AltMonoid should have same Right value")
			}
		})
	}
}

// TestFirstMonoidVsLastMonoid verifies the difference between FirstMonoid and LastMonoid
func TestFirstMonoidVsLastMonoid(t *testing.T) {
	zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
	firstMonoid := FirstMonoid(zero)
	lastMonoid := LastMonoid(zero)

	t.Run("both Right - different results", func(t *testing.T) {
		firstResult := firstMonoid.Concat(Right[error](1), Right[error](2))
		lastResult := lastMonoid.Concat(Right[error](1), Right[error](2))

		assert.Equal(t, Right[error](1), firstResult)
		assert.Equal(t, Right[error](2), lastResult)
		assert.NotEqual(t, firstResult, lastResult)
	})

	t.Run("with Left values - different behavior", func(t *testing.T) {
		err1 := errors.New("err1")
		err2 := errors.New("err2")

		// Both Left: FirstMonoid returns second, LastMonoid returns first
		firstResult := firstMonoid.Concat(Left[int](err1), Left[int](err2))
		lastResult := lastMonoid.Concat(Left[int](err1), Left[int](err2))

		assert.True(t, IsLeft(firstResult))
		assert.True(t, IsLeft(lastResult))
		_, leftErr1 := Unwrap(firstResult)
		_, leftErr2 := Unwrap(lastResult)
		assert.Equal(t, err2, leftErr1)
		assert.Equal(t, err1, leftErr2)
	})

	t.Run("mixed values - same results", func(t *testing.T) {
		testCases := []struct {
			name     string
			left     Either[error, int]
			right    Either[error, int]
			expected Either[error, int]
		}{
			{"left Right, right Left", Right[error](1), Left[int](errors.New("err")), Right[error](1)},
			{"left Left, right Right", Left[int](errors.New("err")), Right[error](2), Right[error](2)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				firstResult := firstMonoid.Concat(tc.left, tc.right)
				lastResult := lastMonoid.Concat(tc.left, tc.right)

				assert.Equal(t, tc.expected, firstResult)
				assert.Equal(t, tc.expected, lastResult)
				assert.Equal(t, firstResult, lastResult)
			})
		}
	})
}

// TestMonoidLaws verifies monoid laws for FirstMonoid and LastMonoid
func TestMonoidLaws(t *testing.T) {
	t.Run("FirstMonoid laws", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := FirstMonoid(zero)

		a := Right[error](1)
		b := Right[error](2)
		c := Right[error](3)

		// Associativity: (a • b) • c = a • (b • c)
		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))
		assert.Equal(t, left, right)

		// Left identity: Empty() • a = a
		leftId := m.Concat(m.Empty(), a)
		assert.Equal(t, a, leftId)

		// Right identity: a • Empty() = a
		rightId := m.Concat(a, m.Empty())
		assert.Equal(t, a, rightId)
	})

	t.Run("LastMonoid laws", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := LastMonoid(zero)

		a := Right[error](1)
		b := Right[error](2)
		c := Right[error](3)

		// Associativity: (a • b) • c = a • (b • c)
		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))
		assert.Equal(t, left, right)

		// Left identity: Empty() • a = a
		leftId := m.Concat(m.Empty(), a)
		assert.Equal(t, a, leftId)

		// Right identity: a • Empty() = a
		rightId := m.Concat(a, m.Empty())
		assert.Equal(t, a, rightId)
	})

	t.Run("FirstMonoid laws with Left values", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := FirstMonoid(zero)

		a := Left[int](errors.New("err1"))
		b := Left[int](errors.New("err2"))
		c := Left[int](errors.New("err3"))

		// Associativity with Left values
		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))
		assert.Equal(t, left, right)
	})

	t.Run("LastMonoid laws with Left values", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := LastMonoid(zero)

		a := Left[int](errors.New("err1"))
		b := Left[int](errors.New("err2"))
		c := Left[int](errors.New("err3"))

		// Associativity with Left values
		left := m.Concat(m.Concat(a, b), c)
		right := m.Concat(a, m.Concat(b, c))
		assert.Equal(t, left, right)
	})
}

// TestMonoidEdgeCases tests edge cases for monoid operations
func TestMonoidEdgeCases(t *testing.T) {
	t.Run("FirstMonoid with empty concatenations", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := FirstMonoid(zero)

		// Empty with empty
		result := m.Concat(m.Empty(), m.Empty())
		assert.True(t, IsLeft(result))
	})

	t.Run("LastMonoid with empty concatenations", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := LastMonoid(zero)

		// Empty with empty
		result := m.Concat(m.Empty(), m.Empty())
		assert.True(t, IsLeft(result))
	})

	t.Run("FirstMonoid chain of operations", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := FirstMonoid(zero)

		// Chain multiple operations
		result := m.Concat(
			m.Concat(
				m.Concat(Left[int](errors.New("err1")), Left[int](errors.New("err2"))),
				Right[error](1),
			),
			m.Concat(Right[error](2), Right[error](3)),
		)
		assert.Equal(t, Right[error](1), result)
	})

	t.Run("LastMonoid chain of operations", func(t *testing.T) {
		zero := func() Either[error, int] { return Left[int](errors.New("empty")) }
		m := LastMonoid(zero)

		// Chain multiple operations
		result := m.Concat(
			m.Concat(Right[error](1), Right[error](2)),
			m.Concat(
				Right[error](3),
				m.Concat(Right[error](4), Left[int](errors.New("err"))),
			),
		)
		assert.Equal(t, Right[error](4), result)
	})
}
