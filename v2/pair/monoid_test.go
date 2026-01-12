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

package pair

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// TestApplicativeMonoidTail tests the ApplicativeMonoidTail implementation
func TestApplicativeMonoidTail(t *testing.T) {
	t.Run("integer addition and string concatenation", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		p1 := MakePair(5, "hello")
		p2 := MakePair(3, " world")

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, 8, Head(result))
		assert.Equal(t, "hello world", Tail(result))
	})

	t.Run("integer multiplication and addition", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidTail(intMul, intAdd)

		p1 := MakePair(3, 5)
		p2 := MakePair(4, 10)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, 12, Head(result)) // 3 * 4
		assert.Equal(t, 15, Tail(result)) // 5 + 10
	})

	t.Run("boolean AND and OR", func(t *testing.T) {
		boolAnd := M.MakeMonoid(func(a, b bool) bool { return a && b }, true)
		boolOr := M.MakeMonoid(func(a, b bool) bool { return a || b }, false)

		pairMonoid := ApplicativeMonoidTail(boolAnd, boolOr)

		p1 := MakePair(true, false)
		p2 := MakePair(true, true)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, true, Head(result)) // true && true
		assert.Equal(t, true, Tail(result)) // false || true
	})

	t.Run("empty value", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		empty := pairMonoid.Empty()
		assert.Equal(t, 0, Head(empty))
		assert.Equal(t, "", Tail(empty))
	})

	t.Run("left identity law", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		p := MakePair(5, "test")
		result := pairMonoid.Concat(pairMonoid.Empty(), p)

		assert.Equal(t, p, result)
	})

	t.Run("right identity law", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		p := MakePair(5, "test")
		result := pairMonoid.Concat(p, pairMonoid.Empty())

		assert.Equal(t, p, result)
	})

	t.Run("associativity law", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		p1 := MakePair(1, "a")
		p2 := MakePair(2, "b")
		p3 := MakePair(3, "c")

		left := pairMonoid.Concat(pairMonoid.Concat(p1, p2), p3)
		right := pairMonoid.Concat(p1, pairMonoid.Concat(p2, p3))

		assert.Equal(t, left, right)
		assert.Equal(t, 6, Head(left))
		assert.Equal(t, "abc", Tail(left))
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		pairMonoid := ApplicativeMonoidTail(intAdd, intMul)

		pairs := []Pair[int, int]{
			MakePair(1, 2),
			MakePair(3, 4),
			MakePair(5, 6),
		}

		result := pairMonoid.Empty()
		for _, p := range pairs {
			result = pairMonoid.Concat(result, p)
		}

		assert.Equal(t, 9, Head(result))  // 0 + 1 + 3 + 5
		assert.Equal(t, 48, Tail(result)) // 1 * 2 * 4 * 6
	})
}

// TestApplicativeMonoidHead tests the ApplicativeMonoidHead implementation
func TestApplicativeMonoidHead(t *testing.T) {
	t.Run("integer multiplication and addition", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidHead(intMul, intAdd)

		p1 := MakePair(3, 5)
		p2 := MakePair(4, 10)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, 12, Head(result)) // 3 * 4
		assert.Equal(t, 15, Tail(result)) // 5 + 10
	})

	t.Run("string concatenation and boolean OR", func(t *testing.T) {
		strConcat := S.Monoid
		boolOr := M.MakeMonoid(func(a, b bool) bool { return a || b }, false)

		pairMonoid := ApplicativeMonoidHead(strConcat, boolOr)

		p1 := MakePair("hello", false)
		p2 := MakePair(" world", true)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, "hello world", Head(result))
		assert.Equal(t, true, Tail(result))
	})

	t.Run("empty value", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidHead(intMul, intAdd)

		empty := pairMonoid.Empty()
		assert.Equal(t, 1, Head(empty))
		assert.Equal(t, 0, Tail(empty))
	})

	t.Run("left identity law", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidHead(intMul, intAdd)

		p := MakePair(5, 10)
		result := pairMonoid.Concat(pairMonoid.Empty(), p)

		assert.Equal(t, p, result)
	})

	t.Run("right identity law", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidHead(intMul, intAdd)

		p := MakePair(5, 10)
		result := pairMonoid.Concat(p, pairMonoid.Empty())

		assert.Equal(t, p, result)
	})

	t.Run("associativity law", func(t *testing.T) {
		intMul := N.MonoidProduct[int]()
		intAdd := N.MonoidSum[int]()

		pairMonoid := ApplicativeMonoidHead(intMul, intAdd)

		p1 := MakePair(2, 1)
		p2 := MakePair(3, 2)
		p3 := MakePair(4, 3)

		left := pairMonoid.Concat(pairMonoid.Concat(p1, p2), p3)
		right := pairMonoid.Concat(p1, pairMonoid.Concat(p2, p3))

		assert.Equal(t, left, right)
		assert.Equal(t, 24, Head(left)) // 2 * 3 * 4
		assert.Equal(t, 6, Tail(left))  // 1 + 2 + 3
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		pairMonoid := ApplicativeMonoidHead(intAdd, intMul)

		pairs := []Pair[int, int]{
			MakePair(1, 2),
			MakePair(3, 4),
			MakePair(5, 6),
		}

		result := pairMonoid.Empty()
		for _, p := range pairs {
			result = pairMonoid.Concat(result, p)
		}

		assert.Equal(t, 9, Head(result))  // 0 + 1 + 3 + 5
		assert.Equal(t, 48, Tail(result)) // 1 * 2 * 4 * 6
	})
}

// TestApplicativeMonoid tests the ApplicativeMonoid alias
func TestApplicativeMonoid(t *testing.T) {
	t.Run("is alias for ApplicativeMonoidTail", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		monoid1 := ApplicativeMonoid(intAdd, strConcat)
		monoid2 := ApplicativeMonoidTail(intAdd, strConcat)

		p1 := MakePair(5, "hello")
		p2 := MakePair(3, " world")

		result1 := monoid1.Concat(p1, p2)
		result2 := monoid2.Concat(p1, p2)

		assert.Equal(t, result1, result2)
		assert.Equal(t, 8, Head(result1))
		assert.Equal(t, "hello world", Tail(result1))
	})

	t.Run("empty values are identical", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		monoid1 := ApplicativeMonoid(intAdd, strConcat)
		monoid2 := ApplicativeMonoidTail(intAdd, strConcat)

		assert.Equal(t, monoid1.Empty(), monoid2.Empty())
	})
}

// TestMonoidHeadVsTail compares ApplicativeMonoidHead and ApplicativeMonoidTail
func TestMonoidHeadVsTail(t *testing.T) {
	t.Run("same result with commutative operations", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		headMonoid := ApplicativeMonoidHead(intMul, intAdd)
		tailMonoid := ApplicativeMonoidTail(intMul, intAdd)

		p1 := MakePair(2, 3)
		p2 := MakePair(4, 5)

		resultHead := headMonoid.Concat(p1, p2)
		resultTail := tailMonoid.Concat(p1, p2)

		// Both should give same result since operations are commutative
		assert.Equal(t, 8, Head(resultHead)) // 2 * 4
		assert.Equal(t, 8, Tail(resultHead)) // 3 + 5
		assert.Equal(t, 8, Head(resultTail)) // 2 * 4
		assert.Equal(t, 8, Tail(resultTail)) // 3 + 5
	})

	t.Run("different empty values", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		headMonoid := ApplicativeMonoidHead(intMul, intAdd)
		tailMonoid := ApplicativeMonoidTail(intAdd, intMul)

		emptyHead := headMonoid.Empty()
		emptyTail := tailMonoid.Empty()

		assert.Equal(t, 1, Head(emptyHead)) // intMul empty
		assert.Equal(t, 0, Tail(emptyHead)) // intAdd empty
		assert.Equal(t, 0, Head(emptyTail)) // intAdd empty
		assert.Equal(t, 1, Tail(emptyTail)) // intMul empty
	})
}

// TestMonoidLaws verifies monoid laws for all implementations
func TestMonoidLaws(t *testing.T) {
	testCases := []struct {
		name       string
		monoid     M.Monoid[Pair[int, int]]
		p1, p2, p3 Pair[int, int]
	}{
		{
			name:   "ApplicativeMonoidTail",
			monoid: ApplicativeMonoidTail(N.MonoidSum[int](), N.MonoidProduct[int]()),
			p1:     MakePair(1, 2),
			p2:     MakePair(3, 4),
			p3:     MakePair(5, 6),
		},
		{
			name:   "ApplicativeMonoidHead",
			monoid: ApplicativeMonoidHead(N.MonoidProduct[int](), N.MonoidSum[int]()),
			p1:     MakePair(2, 1),
			p2:     MakePair(3, 2),
			p3:     MakePair(4, 3),
		},
		{
			name:   "ApplicativeMonoid",
			monoid: ApplicativeMonoid(N.MonoidSum[int](), N.MonoidSum[int]()),
			p1:     MakePair(1, 2),
			p2:     MakePair(3, 4),
			p3:     MakePair(5, 6),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Run("associativity", func(t *testing.T) {
				left := tc.monoid.Concat(tc.monoid.Concat(tc.p1, tc.p2), tc.p3)
				right := tc.monoid.Concat(tc.p1, tc.monoid.Concat(tc.p2, tc.p3))
				assert.Equal(t, left, right)
			})

			t.Run("left identity", func(t *testing.T) {
				result := tc.monoid.Concat(tc.monoid.Empty(), tc.p1)
				assert.Equal(t, tc.p1, result)
			})

			t.Run("right identity", func(t *testing.T) {
				result := tc.monoid.Concat(tc.p1, tc.monoid.Empty())
				assert.Equal(t, tc.p1, result)
			})
		})
	}
}

// TestMonoidEdgeCases tests edge cases for monoid operations
func TestMonoidEdgeCases(t *testing.T) {
	t.Run("concatenating empty with empty", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(intAdd, strConcat)

		result := pairMonoid.Concat(pairMonoid.Empty(), pairMonoid.Empty())
		assert.Equal(t, pairMonoid.Empty(), result)
	})

	t.Run("chain of operations", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		pairMonoid := ApplicativeMonoidTail(intAdd, intMul)

		result := pairMonoid.Concat(
			pairMonoid.Concat(
				pairMonoid.Concat(MakePair(1, 2), MakePair(2, 3)),
				MakePair(3, 4),
			),
			MakePair(4, 5),
		)

		assert.Equal(t, 10, Head(result))  // 1 + 2 + 3 + 4
		assert.Equal(t, 120, Tail(result)) // 2 * 3 * 4 * 5
	})

	t.Run("zero values", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		pairMonoid := ApplicativeMonoidTail(intAdd, intMul)

		p1 := MakePair(0, 0)
		p2 := MakePair(5, 10)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, 5, Head(result))
		assert.Equal(t, 0, Tail(result)) // 0 * 10 = 0
	})

	t.Run("negative values", func(t *testing.T) {
		intAdd := N.MonoidSum[int]()
		intMul := N.MonoidProduct[int]()

		pairMonoid := ApplicativeMonoidTail(intAdd, intMul)

		p1 := MakePair(-5, -2)
		p2 := MakePair(3, 4)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, -2, Head(result)) // -5 + 3
		assert.Equal(t, -8, Tail(result)) // -2 * 4
	})
}

// TestMonoidWithDifferentTypes tests monoids with various type combinations
func TestMonoidWithDifferentTypes(t *testing.T) {
	t.Run("string and boolean", func(t *testing.T) {
		strConcat := S.Monoid
		boolAnd := M.MakeMonoid(func(a, b bool) bool { return a && b }, true)

		pairMonoid := ApplicativeMonoidTail(strConcat, boolAnd)

		p1 := MakePair("hello", true)
		p2 := MakePair(" world", true)

		result := pairMonoid.Concat(p1, p2)
		// Note: The order depends on the applicative implementation
		assert.Equal(t, " worldhello", Head(result))
		assert.Equal(t, true, Tail(result))
	})

	t.Run("boolean and string", func(t *testing.T) {
		boolOr := M.MakeMonoid(func(a, b bool) bool { return a || b }, false)
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(boolOr, strConcat)

		p1 := MakePair(false, "foo")
		p2 := MakePair(true, "bar")

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, true, Head(result))
		assert.Equal(t, "foobar", Tail(result))
	})

	t.Run("float64 addition and multiplication", func(t *testing.T) {
		floatAdd := N.MonoidSum[float64]()
		floatMul := N.MonoidProduct[float64]()

		pairMonoid := ApplicativeMonoidTail(floatAdd, floatMul)

		p1 := MakePair(1.5, 2.0)
		p2 := MakePair(2.5, 3.0)

		result := pairMonoid.Concat(p1, p2)
		assert.Equal(t, 4.0, Head(result))
		assert.Equal(t, 6.0, Tail(result))
	})
}

// TestMonoidCommutativity tests behavior with non-commutative operations
func TestMonoidCommutativity(t *testing.T) {
	t.Run("string concatenation is not commutative", func(t *testing.T) {
		strConcat := S.Monoid

		pairMonoid := ApplicativeMonoidTail(strConcat, strConcat)

		p1 := MakePair("hello", "foo")
		p2 := MakePair(" world", "bar")

		result1 := pairMonoid.Concat(p1, p2)
		result2 := pairMonoid.Concat(p2, p1)

		// The applicative implementation reverses the order for head values
		assert.Equal(t, " worldhello", Head(result1))
		assert.Equal(t, "foobar", Tail(result1))
		assert.Equal(t, "hello world", Head(result2))
		assert.Equal(t, "barfoo", Tail(result2))
		assert.NotEqual(t, result1, result2)
	})
}
