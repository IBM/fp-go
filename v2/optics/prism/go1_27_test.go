//go:build go1.27

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

package prism

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestPrismCompose_MethodEquivalentToFreeFunction verifies that the method form
// p.Compose(ab) returns a prism identical in behaviour to the free function
// Compose[S](ab)(p) for every operation.
func TestPrismCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	// outer: Option[int] → int  (extracts from Some)
	outer := FromOption[int]()
	// inner: int → int  (matches only positive numbers)
	inner := FromPredicate(func(n int) bool { return n > 0 })

	method := outer.Compose(inner)
	free := Compose[Option[int]](inner)(outer)

	cases := []Option[int]{O.Some(42), O.Some(-1), O.None[int]()}
	for _, c := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(c), method.GetOption(c))
		})
	}

	t.Run("ReverseGet parity", func(t *testing.T) {
		assert.Equal(t, free.ReverseGet(7), method.ReverseGet(7))
	})
}

// TestPrismCompose_Method_GetOption_Match verifies Some is returned when both
// constituent prisms match.
func TestPrismCompose_Method_GetOption_Match(t *testing.T) {
	outer := FromOption[int]()
	inner := FromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	got := composed.GetOption(O.Some(42))
	assert.True(t, O.IsSome(got))
	assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(got))
}

// TestPrismCompose_Method_GetOption_OuterMiss verifies None is returned when
// the outer prism does not match.
func TestPrismCompose_Method_GetOption_OuterMiss(t *testing.T) {
	outer := FromOption[int]()
	inner := FromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	assert.True(t, O.IsNone(composed.GetOption(O.None[int]())))
}

// TestPrismCompose_Method_GetOption_InnerMiss verifies None is returned when
// the outer prism matches but the inner does not.
func TestPrismCompose_Method_GetOption_InnerMiss(t *testing.T) {
	outer := FromOption[int]()
	inner := FromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	assert.True(t, O.IsNone(composed.GetOption(O.Some(-5))))
}

// TestPrismCompose_Method_ReverseGet verifies that ReverseGet threads the value
// back through the inner then outer ReverseGet.
func TestPrismCompose_Method_ReverseGet(t *testing.T) {
	outer := FromOption[int]()
	inner := FromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	// ReverseGet(42): inner.ReverseGet(42)==42, outer.ReverseGet(42)==Some(42)
	assert.Equal(t, O.Some(42), composed.ReverseGet(42))
}

// TestPrismCompose_Method_PrismLaws verifies the two standard prism laws on the
// method-composed prism.
func TestPrismCompose_Method_PrismLaws(t *testing.T) {
	outer := FromOption[int]()
	inner := FromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		a := 99
		got := composed.GetOption(composed.ReverseGet(a))
		assert.True(t, O.IsSome(got))
		assert.Equal(t, a, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("law 2: if GetOption(s)==Some(a) then ReverseGet(a) round-trips", func(t *testing.T) {
		s := O.Some(7)
		got := composed.GetOption(s)
		if O.IsSome(got) {
			a := O.GetOrElse(F.Constant(0))(got)
			reconstructed := composed.ReverseGet(a)
			assert.Equal(t, O.Some(a), composed.GetOption(reconstructed))
		}
	})
}

// TestPrismCompose_Method_Chained verifies that multiple .Compose calls chain
// correctly through three levels of sum type.
func TestPrismCompose_Method_Chained(t *testing.T) {
	// Level 1: Option[Option[int]] → Option[int]  (unwrap outer Some)
	l1 := FromOption[Option[int]]()
	// Level 2: Option[int] → int                  (unwrap inner Some)
	l2 := FromOption[int]()
	// Level 3: int → int                          (only positives)
	l3 := FromPredicate(func(n int) bool { return n > 0 })

	// Chain: Option[Option[int]] → Option[int] → int → int (positive)
	composed := l1.Compose(l2).Compose(l3)

	t.Run("all three levels match", func(t *testing.T) {
		s := O.Some(O.Some(5))
		got := composed.GetOption(s)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 5, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("outermost None short-circuits", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(O.None[Option[int]]())))
	})

	t.Run("middle None short-circuits", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(O.Some(O.None[int]()))))
	})

	t.Run("innermost predicate miss short-circuits", func(t *testing.T) {
		assert.True(t, O.IsNone(composed.GetOption(O.Some(O.Some(-3)))))
	})

	t.Run("ReverseGet reconstructs all three layers", func(t *testing.T) {
		// inner.ReverseGet(10)==10 → l2.ReverseGet(10)==Some(10) → l1.ReverseGet(Some(10))==Some(Some(10))
		assert.Equal(t, O.Some(O.Some(10)), composed.ReverseGet(10))
	})
}
