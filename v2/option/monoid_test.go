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

package option

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/semigroup"
	"github.com/stretchr/testify/assert"
)

// TestSemigroupAssociativity tests the associativity law for Semigroup
func TestSemigroupAssociativity(t *testing.T) {
	intSemigroup := S.MakeSemigroup(func(a, b int) int { return a + b })
	optSemigroup := Semigroup[int]()(intSemigroup)

	a := Some(1)
	b := Some(2)
	c := Some(3)

	// Test that (a • b) • c = a • (b • c)
	left := optSemigroup.Concat(optSemigroup.Concat(a, b), c)
	right := optSemigroup.Concat(a, optSemigroup.Concat(b, c))

	assert.Equal(t, left, right)
	assert.Equal(t, Some(6), left)
}

// TestSemigroupWithNone tests Semigroup behavior with None values
func TestSemigroupWithNone(t *testing.T) {
	intSemigroup := S.MakeSemigroup(func(a, b int) int { return a + b })
	optSemigroup := Semigroup[int]()(intSemigroup)

	t.Run("None with None", func(t *testing.T) {
		result := optSemigroup.Concat(None[int](), None[int]())
		assert.Equal(t, None[int](), result)
	})

	t.Run("associativity with None", func(t *testing.T) {
		a := None[int]()
		b := Some(2)
		c := Some(3)

		left := optSemigroup.Concat(optSemigroup.Concat(a, b), c)
		right := optSemigroup.Concat(a, optSemigroup.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Some(5), left)
	})
}

// TestMonoidIdentityLaws tests the identity laws for Monoid
func TestMonoidIdentityLaws(t *testing.T) {
	intSemigroup := S.MakeSemigroup(func(a, b int) int { return a + b })
	optMonoid := Monoid[int]()(intSemigroup)

	t.Run("left identity with Some", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity with Some", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("left identity with None", func(t *testing.T) {
		x := None[int]()
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity with None", func(t *testing.T) {
		x := None[int]()
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})
}

// TestAlternativeMonoidIdentityLaws tests identity laws for AlternativeMonoid
func TestAlternativeMonoidIdentityLaws(t *testing.T) {
	intMonoid := N.MonoidSum[int]()
	optMonoid := AlternativeMonoid(intMonoid)

	t.Run("left identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("empty is Some(0)", func(t *testing.T) {
		empty := optMonoid.Empty()
		assert.Equal(t, Some(0), empty)
	})
}

// TestAltMonoidIdentityLaws tests identity laws for AltMonoid
func TestAltMonoidIdentityLaws(t *testing.T) {
	optMonoid := AltMonoid[int]()

	t.Run("left identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("associativity", func(t *testing.T) {
		a := Some(1)
		b := None[int]()
		c := Some(3)

		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Some(1), left)
	})
}

// TestFirstMonoid tests the FirstMonoid implementation
func TestFirstMonoid(t *testing.T) {
	optMonoid := FirstMonoid[int]()

	t.Run("both Some values - returns first", func(t *testing.T) {
		result := optMonoid.Concat(Some(2), Some(3))
		assert.Equal(t, Some(2), result)
	})

	t.Run("left Some, right None", func(t *testing.T) {
		result := optMonoid.Concat(Some(2), None[int]())
		assert.Equal(t, Some(2), result)
	})

	t.Run("left None, right Some", func(t *testing.T) {
		result := optMonoid.Concat(None[int](), Some(3))
		assert.Equal(t, Some(3), result)
	})

	t.Run("both None", func(t *testing.T) {
		result := optMonoid.Concat(None[int](), None[int]())
		assert.Equal(t, None[int](), result)
	})

	t.Run("empty value", func(t *testing.T) {
		empty := optMonoid.Empty()
		assert.Equal(t, None[int](), empty)
	})

	t.Run("left identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("associativity", func(t *testing.T) {
		a := Some(1)
		b := Some(2)
		c := Some(3)

		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Some(1), left)
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		// Should return the first Some value encountered
		result := optMonoid.Concat(
			optMonoid.Concat(None[int](), Some(1)),
			optMonoid.Concat(Some(2), Some(3)),
		)
		assert.Equal(t, Some(1), result)
	})

	t.Run("with strings", func(t *testing.T) {
		strMonoid := FirstMonoid[string]()

		result := strMonoid.Concat(Some("first"), Some("second"))
		assert.Equal(t, Some("first"), result)

		result = strMonoid.Concat(None[string](), Some("second"))
		assert.Equal(t, Some("second"), result)
	})
}

// TestLastMonoid tests the LastMonoid implementation
func TestLastMonoid(t *testing.T) {
	optMonoid := LastMonoid[int]()

	t.Run("both Some values - returns last", func(t *testing.T) {
		result := optMonoid.Concat(Some(2), Some(3))
		assert.Equal(t, Some(3), result)
	})

	t.Run("left Some, right None", func(t *testing.T) {
		result := optMonoid.Concat(Some(2), None[int]())
		assert.Equal(t, Some(2), result)
	})

	t.Run("left None, right Some", func(t *testing.T) {
		result := optMonoid.Concat(None[int](), Some(3))
		assert.Equal(t, Some(3), result)
	})

	t.Run("both None", func(t *testing.T) {
		result := optMonoid.Concat(None[int](), None[int]())
		assert.Equal(t, None[int](), result)
	})

	t.Run("empty value", func(t *testing.T) {
		empty := optMonoid.Empty()
		assert.Equal(t, None[int](), empty)
	})

	t.Run("left identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(optMonoid.Empty(), x)
		assert.Equal(t, x, result)
	})

	t.Run("right identity", func(t *testing.T) {
		x := Some(5)
		result := optMonoid.Concat(x, optMonoid.Empty())
		assert.Equal(t, x, result)
	})

	t.Run("associativity", func(t *testing.T) {
		a := Some(1)
		b := Some(2)
		c := Some(3)

		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))

		assert.Equal(t, left, right)
		assert.Equal(t, Some(3), left)
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		// Should return the last Some value encountered
		result := optMonoid.Concat(
			optMonoid.Concat(Some(1), Some(2)),
			optMonoid.Concat(Some(3), None[int]()),
		)
		assert.Equal(t, Some(3), result)
	})

	t.Run("with strings", func(t *testing.T) {
		strMonoid := LastMonoid[string]()

		result := strMonoid.Concat(Some("first"), Some("second"))
		assert.Equal(t, Some("second"), result)

		result = strMonoid.Concat(Some("first"), None[string]())
		assert.Equal(t, Some("first"), result)
	})
}

// TestFirstMonoidVsAltMonoid verifies FirstMonoid and AltMonoid have the same behavior
func TestFirstMonoidVsAltMonoid(t *testing.T) {
	firstMonoid := FirstMonoid[int]()
	altMonoid := AltMonoid[int]()

	testCases := []struct {
		name  string
		left  Option[int]
		right Option[int]
	}{
		{"both Some", Some(1), Some(2)},
		{"left Some, right None", Some(1), None[int]()},
		{"left None, right Some", None[int](), Some(2)},
		{"both None", None[int](), None[int]()},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			firstResult := firstMonoid.Concat(tc.left, tc.right)
			altResult := altMonoid.Concat(tc.left, tc.right)
			assert.Equal(t, firstResult, altResult, "FirstMonoid and AltMonoid should behave the same")
		})
	}
}

// TestFirstMonoidVsLastMonoid verifies the difference between FirstMonoid and LastMonoid
func TestFirstMonoidVsLastMonoid(t *testing.T) {
	firstMonoid := FirstMonoid[int]()
	lastMonoid := LastMonoid[int]()

	t.Run("both Some - different results", func(t *testing.T) {
		firstResult := firstMonoid.Concat(Some(1), Some(2))
		lastResult := lastMonoid.Concat(Some(1), Some(2))

		assert.Equal(t, Some(1), firstResult)
		assert.Equal(t, Some(2), lastResult)
		assert.NotEqual(t, firstResult, lastResult)
	})

	t.Run("with None - same results", func(t *testing.T) {
		testCases := []struct {
			name     string
			left     Option[int]
			right    Option[int]
			expected Option[int]
		}{
			{"left Some, right None", Some(1), None[int](), Some(1)},
			{"left None, right Some", None[int](), Some(2), Some(2)},
			{"both None", None[int](), None[int](), None[int]()},
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

// TestMonoidComparison compares different monoid implementations
func TestMonoidComparison(t *testing.T) {
	t.Run("Monoid vs AlternativeMonoid with addition", func(t *testing.T) {
		intSemigroup := S.MakeSemigroup(func(a, b int) int { return a + b })
		regularMonoid := Monoid[int]()(intSemigroup)

		intMonoid := M.MakeMonoid(func(a, b int) int { return a + b }, 0)
		altMonoid := AlternativeMonoid(intMonoid)

		// Both should combine Some values the same way
		assert.Equal(t,
			regularMonoid.Concat(Some(2), Some(3)),
			altMonoid.Concat(Some(2), Some(3)),
		)

		// But empty values differ
		assert.Equal(t, None[int](), regularMonoid.Empty())
		assert.Equal(t, Some(0), altMonoid.Empty())
	})
}

// TestMonoidLaws verifies monoid laws for all monoid implementations
func TestMonoidLaws(t *testing.T) {
	t.Run("Monoid with addition", func(t *testing.T) {
		intSemigroup := N.SemigroupSum[int]()
		optMonoid := Monoid[int]()(intSemigroup)

		a := Some(1)
		b := Some(2)
		c := Some(3)

		// Associativity: (a • b) • c = a • (b • c)
		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))
		assert.Equal(t, left, right)

		// Left identity: Empty() • a = a
		leftId := optMonoid.Concat(optMonoid.Empty(), a)
		assert.Equal(t, a, leftId)

		// Right identity: a • Empty() = a
		rightId := optMonoid.Concat(a, optMonoid.Empty())
		assert.Equal(t, a, rightId)
	})

	t.Run("FirstMonoid laws", func(t *testing.T) {
		optMonoid := FirstMonoid[int]()

		a := Some(1)
		b := Some(2)
		c := Some(3)

		// Associativity
		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))
		assert.Equal(t, left, right)

		// Left identity
		leftId := optMonoid.Concat(optMonoid.Empty(), a)
		assert.Equal(t, a, leftId)

		// Right identity
		rightId := optMonoid.Concat(a, optMonoid.Empty())
		assert.Equal(t, a, rightId)
	})

	t.Run("LastMonoid laws", func(t *testing.T) {
		optMonoid := LastMonoid[int]()

		a := Some(1)
		b := Some(2)
		c := Some(3)

		// Associativity
		left := optMonoid.Concat(optMonoid.Concat(a, b), c)
		right := optMonoid.Concat(a, optMonoid.Concat(b, c))
		assert.Equal(t, left, right)

		// Left identity
		leftId := optMonoid.Concat(optMonoid.Empty(), a)
		assert.Equal(t, a, leftId)

		// Right identity
		rightId := optMonoid.Concat(a, optMonoid.Empty())
		assert.Equal(t, a, rightId)
	})
}
