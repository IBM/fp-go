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
	"fmt"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Domain types shared across all composition tests
// ---------------------------------------------------------------------------

// compShape is a sum type (interface) used in composition tests.
type compShape interface{ isCompShape() }

type compCircle struct{ Radius float64 }
type compRect struct{ Width, Height float64 }

func (compCircle) isCompShape() {}
func (compRect) isCompShape()   {}

// compCirclePrism selects the compCircle variant from a compShape.
var compCirclePrism = MakePrism(
	func(s compShape) O.Option[compCircle] {
		if c, ok := s.(compCircle); ok {
			return O.Some(c)
		}
		return O.None[compCircle]()
	},
	func(c compCircle) compShape { return c },
)

// compRadiusLens focuses on the Radius field of a compCircle.
var compRadiusLens = common.MakeLens(
	func(c compCircle) float64 { return c.Radius },
	func(c compCircle, r float64) compCircle { c.Radius = r; return c },
)

// positiveFloatPrism admits only positive float64 values.
var positiveFloatPrism = MakePrism(
	func(r float64) O.Option[float64] {
		if r > 0 {
			return O.Some(r)
		}
		return O.None[float64]()
	},
	func(r float64) float64 { return r },
)

// positiveRadiusOptional is an Optional[compCircle, float64] that only admits
// positive radii.
var positiveRadiusOptional = common.MakeOptional(
	func(c compCircle) O.Option[float64] {
		if c.Radius > 0 {
			return O.Some(c.Radius)
		}
		return O.None[float64]()
	},
	func(c compCircle, r float64) compCircle { c.Radius = r; return c },
)

// ---------------------------------------------------------------------------
// Compose (Prism ∘ Prism → Prism)
// ---------------------------------------------------------------------------

// TestCompose_GetOption_Match verifies that composing two prisms returns Some
// when both individual prisms match.
//
// Composition: Prism[compShape, compCircle] ∘ Prism[compCircle, float64]
// where the outer prism selects the Circle variant and the inner prism
// extracts a positive radius from it.
func TestCompose_GetOption_Match(t *testing.T) {
	// Build a prism that extracts a positive radius from a compCircle.
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	got := composed.GetOption(compCircle{Radius: 5})
	assert.True(t, O.IsSome(got))
	assert.Equal(t, 5.0, O.GetOrElse(F.Constant(0.0))(got))
}

// TestCompose_GetOption_OuterMiss verifies None when the outer prism misses.
func TestCompose_GetOption_OuterMiss(t *testing.T) {
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	got := composed.GetOption(compRect{Width: 3, Height: 4})
	assert.True(t, O.IsNone(got))
}

// TestCompose_GetOption_InnerMiss verifies None when the inner prism misses.
func TestCompose_GetOption_InnerMiss(t *testing.T) {
	// Inner prism admits only positive radii.
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] {
			if c.Radius > 0 {
				return O.Some(c.Radius)
			}
			return O.None[float64]()
		},
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	// Circle with zero radius: outer prism matches, inner misses.
	got := composed.GetOption(compCircle{Radius: 0})
	assert.True(t, O.IsNone(got))
}

// TestCompose_ReverseGet verifies that ReverseGet threads the value through
// both prisms in reverse order.
func TestCompose_ReverseGet(t *testing.T) {
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	s := composed.ReverseGet(7.0)
	c, ok := s.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 7.0, c.Radius)
}

// TestCompose_PrismLaw1 verifies GetOption(ReverseGet(b)) == Some(b).
func TestCompose_PrismLaw1(t *testing.T) {
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	for _, r := range []float64{1, 2.5, 100} {
		got := composed.GetOption(composed.ReverseGet(r))
		assert.True(t, O.IsSome(got), "expected Some for radius %v", r)
		assert.Equal(t, r, O.GetOrElse(F.Constant(0.0))(got))
	}
}

// TestCompose_PrismLaw2 verifies that if GetOption(s) == Some(a) then
// GetOption(ReverseGet(a)) == Some(a).
func TestCompose_PrismLaw2(t *testing.T) {
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](innerPrism)(compCirclePrism)

	s := compCircle{Radius: 3.0}
	opt := composed.GetOption(s)
	if O.IsSome(opt) {
		a := O.GetOrElse(F.Constant(0.0))(opt)
		second := composed.GetOption(composed.ReverseGet(a))
		assert.True(t, O.IsSome(second))
		assert.Equal(t, a, O.GetOrElse(F.Constant(0.0))(second))
	}
}

// TestCompose_WithParseInt demonstrates composing a head-of-slice prism with
// a string-to-int parsing prism, mirroring the pattern used in at_head_test.go.
func TestCompose_WithParseInt(t *testing.T) {
	// Prism[[]string, string] ∘ Prism[string, int]
	composed := Compose[[]string](ParseInt())(Head[string]())

	t.Run("first element parses to int", func(t *testing.T) {
		got := composed.GetOption([]string{"42", "not-a-number"})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("first element does not parse", func(t *testing.T) {
		got := composed.GetOption([]string{"abc"})
		assert.True(t, O.IsNone(got))
	})

	t.Run("empty slice gives None", func(t *testing.T) {
		got := composed.GetOption([]string{})
		assert.True(t, O.IsNone(got))
	})

	t.Run("ReverseGet reconstructs slice", func(t *testing.T) {
		s := composed.ReverseGet(7)
		assert.Equal(t, []string{"7"}, s)
	})
}

// ---------------------------------------------------------------------------
// ComposeIso (Prism ∘ Iso → Prism)
// ---------------------------------------------------------------------------

// circleToRadiusIso maps compCircle ↔ float64 via the Radius field.
// This is not a true isomorphism (a float64 cannot round-trip back to the same
// compCircle if other fields existed), but it satisfies the iso laws for this
// single-field struct.
var circleToRadiusIso = common.MakeIso(
	func(c compCircle) float64 { return c.Radius },
	func(r float64) compCircle { return compCircle{Radius: r} },
)

// radiusToIntIso maps float64 → int by truncation.
var radiusToIntIso = common.MakeIso(
	func(f float64) int { return int(f) },
	func(i int) float64 { return float64(i) },
)

// TestComposeIso_GetOption_Match verifies Some is returned and the iso's
// forward function is applied to the focused value.
func TestComposeIso_GetOption_Match(t *testing.T) {
	// Prism[compShape, compCircle] ∘ Iso[compCircle, float64]
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	got := composed.GetOption(compCircle{Radius: 7.9})
	assert.True(t, O.IsSome(got))
	assert.Equal(t, 7.9, O.GetOrElse(F.Constant(-1.0))(got))
}

// TestComposeIso_GetOption_Miss verifies None when the prism does not match.
func TestComposeIso_GetOption_Miss(t *testing.T) {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	got := composed.GetOption(compRect{Width: 1, Height: 2})
	assert.True(t, O.IsNone(got))
}

// TestComposeIso_ReverseGet verifies that ReverseGet threads through the iso's
// ReverseGet and then the prism's ReverseGet.
func TestComposeIso_ReverseGet(t *testing.T) {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	s := composed.ReverseGet(5.0)
	c, ok := s.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 5.0, c.Radius)
}

// TestComposeIso_PrismLaw1 verifies GetOption(ReverseGet(b)) == Some(b).
func TestComposeIso_PrismLaw1(t *testing.T) {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	for _, r := range []float64{0, 1, 42.5, 100} {
		got := composed.GetOption(composed.ReverseGet(r))
		assert.True(t, O.IsSome(got))
		assert.Equal(t, r, O.GetOrElse(F.Constant(-1.0))(got))
	}
}

// TestComposeIso_PrismLaw2 verifies round-trip consistency.
func TestComposeIso_PrismLaw2(t *testing.T) {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	s := compCircle{Radius: 3.0}
	opt := composed.GetOption(s)
	if O.IsSome(opt) {
		a := O.GetOrElse(F.Constant(-1.0))(opt)
		second := composed.GetOption(composed.ReverseGet(a))
		assert.True(t, O.IsSome(second))
		assert.Equal(t, a, O.GetOrElse(F.Constant(-1.0))(second))
	}
}

// TestComposeIso_IdentityIso verifies that composing with the identity iso
// leaves GetOption results unchanged.
func TestComposeIso_IdentityIso(t *testing.T) {
	idIso := common.MakeIso(
		func(c compCircle) compCircle { return c },
		func(c compCircle) compCircle { return c },
	)
	composed := ComposeIso[compShape](idIso)(compCirclePrism)

	original := compCircle{Radius: 4.0}
	origOpt := compCirclePrism.GetOption(original)
	compOpt := composed.GetOption(original)
	assert.Equal(t, origOpt, compOpt)
}

// TestComposeIso_IntegerTruncation demonstrates composing with a non-identity
// iso that changes the focus type (compCircle → float64 → int).
func TestComposeIso_IntegerTruncation(t *testing.T) {
	// Compose two isos: compCircle → float64 → int
	step1 := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism) // Prism[compShape, float64]
	step2 := ComposeIso[compShape](radiusToIntIso)(step1)              // Prism[compShape, int]

	t.Run("GetOption truncates to int", func(t *testing.T) {
		got := step2.GetOption(compCircle{Radius: 7.9})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 7, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("ReverseGet reconstructs compShape from int", func(t *testing.T) {
		s := step2.ReverseGet(5)
		c, ok := s.(compCircle)
		assert.True(t, ok)
		assert.Equal(t, 5.0, c.Radius)
	})

	t.Run("prism law 1: GetOption(ReverseGet(n)) == Some(n)", func(t *testing.T) {
		for _, n := range []int{0, 1, 42, 100} {
			got := step2.GetOption(step2.ReverseGet(n))
			assert.True(t, O.IsSome(got))
			assert.Equal(t, n, O.GetOrElse(F.Constant(-1))(got))
		}
	})
}

// TestComposeIso_Chaining verifies that ComposeIso can be chained: the result
// of one ComposeIso can be the input of another.
func TestComposeIso_Chaining(t *testing.T) {
	// compCircle → float64 → string
	strIso := common.MakeIso(
		func(f float64) string { return strconv.FormatFloat(f, 'f', -1, 64) },
		func(s string) float64 {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		},
	)

	step1 := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism) // Prism[compShape, float64]
	step2 := ComposeIso[compShape](strIso)(step1)                      // Prism[compShape, string]

	t.Run("GetOption extracts and transforms through both isos", func(t *testing.T) {
		got := step2.GetOption(compCircle{Radius: 9.5})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, "9.5", O.GetOrElse(F.Constant(""))(got))
	})

	t.Run("ReverseGet threads through both isos in reverse", func(t *testing.T) {
		s := step2.ReverseGet("3.0")
		c, ok := s.(compCircle)
		assert.True(t, ok)
		assert.Equal(t, 3.0, c.Radius)
	})
}

// ---------------------------------------------------------------------------
// ComposeLens (Prism ∘ Lens → Optional)
// ---------------------------------------------------------------------------

// TestComposeLens_GetOption_Match verifies Some is returned and the lens
// focuses on the correct field.
func TestComposeLens_GetOption_Match(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	got := opt.GetOption(compCircle{Radius: 4.0})
	assert.True(t, O.IsSome(got))
	assert.Equal(t, 4.0, O.GetOrElse(F.Constant(0.0))(got))
}

// TestComposeLens_GetOption_Miss verifies None when the prism does not match.
func TestComposeLens_GetOption_Miss(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	got := opt.GetOption(compRect{Width: 2, Height: 3})
	assert.True(t, O.IsNone(got))
}

// TestComposeLens_Set_Match verifies that Set updates the focused field when
// the prism matches.
func TestComposeLens_Set_Match(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	result := opt.Set(99.0)(compCircle{Radius: 1.0})
	c, ok := result.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 99.0, c.Radius)
}

// TestComposeLens_Set_NoOpOnPrismMiss verifies that Set is a no-op when the
// prism does not match (Optional law 1).
func TestComposeLens_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	original := compRect{Width: 2, Height: 3}
	result := opt.Set(99.0)(original)
	assert.Equal(t, original, result)
}

// TestComposeLens_OptionalLaws verifies all three standard Optional laws.
func TestComposeLens_OptionalLaws(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	matchingShape := compCircle{Radius: 5.0}
	missingShape := compRect{Width: 1, Height: 2}

	t.Run("law 1 – GetSet: no-op on None", func(t *testing.T) {
		assert.Equal(t, missingShape, opt.Set(99.0)(missingShape))
	})

	t.Run("law 2 – SetGet: get what you set", func(t *testing.T) {
		updated := opt.Set(42.0)(matchingShape)
		got := opt.GetOption(updated)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42.0, O.GetOrElse(F.Constant(0.0))(got))
	})

	t.Run("law 3 – SetSet: last set wins", func(t *testing.T) {
		first := opt.Set(1.0)(matchingShape)
		both := opt.Set(2.0)(first)
		once := opt.Set(2.0)(matchingShape)
		assert.Equal(t, once, both)
	})
}

// TestComposeLens_Immutability verifies that Set does not mutate the original
// value.
func TestComposeLens_Immutability(t *testing.T) {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	original := compCircle{Radius: 3.0}
	_ = opt.Set(99.0)(original)
	assert.Equal(t, 3.0, original.Radius)
}

// TestComposeLens_NonIdentityLens verifies composition with a Lens that
// transforms the field type.
func TestComposeLens_NonIdentityLens(t *testing.T) {
	// Lens that exposes the truncated integer radius.
	truncLens := common.MakeLens(
		func(c compCircle) int { return int(c.Radius) },
		func(c compCircle, n int) compCircle { c.Radius = float64(n); return c },
	)
	opt := ComposeLens[compShape](truncLens)(compCirclePrism)

	t.Run("GetOption truncates radius to int", func(t *testing.T) {
		got := opt.GetOption(compCircle{Radius: 3.9})
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 3, O.GetOrElse(F.Constant(-1))(got))
	})

	t.Run("Set updates radius via lens", func(t *testing.T) {
		result := opt.Set(7)(compCircle{Radius: 1.5})
		c, ok := result.(compCircle)
		assert.True(t, ok)
		assert.Equal(t, 7.0, c.Radius)
	})
}

// ---------------------------------------------------------------------------
// ComposeOptional (Prism ∘ Optional → Optional)
// ---------------------------------------------------------------------------

// TestComposeOptional_GetOption_Match verifies Some is returned when both the
// prism and the inner Optional match.
func TestComposeOptional_GetOption_Match(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	got := opt.GetOption(compCircle{Radius: 5.0})
	assert.True(t, O.IsSome(got))
	assert.Equal(t, 5.0, O.GetOrElse(F.Constant(0.0))(got))
}

// TestComposeOptional_GetOption_PrismMiss verifies None when the prism misses.
func TestComposeOptional_GetOption_PrismMiss(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	got := opt.GetOption(compRect{Width: 1, Height: 2})
	assert.True(t, O.IsNone(got))
}

// TestComposeOptional_GetOption_InnerMiss verifies None when the prism matches
// but the inner Optional misses.
func TestComposeOptional_GetOption_InnerMiss(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	// Circle with zero radius: outer prism matches, inner Optional misses.
	got := opt.GetOption(compCircle{Radius: 0})
	assert.True(t, O.IsNone(got))
}

// TestComposeOptional_Set_Match verifies that Set updates the focused value
// when both optics match.
func TestComposeOptional_Set_Match(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	result := opt.Set(10.0)(compCircle{Radius: 3.0})
	c, ok := result.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 10.0, c.Radius)
}

// TestComposeOptional_Set_NoOpOnPrismMiss verifies that Set is a no-op when
// the prism does not match (Optional law 1 – outer miss).
func TestComposeOptional_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	original := compRect{Width: 2, Height: 3}
	result := opt.Set(10.0)(original)
	assert.Equal(t, original, result)
}

// TestComposeOptional_GetOption_InnerMissDoesNotAffectSet verifies the
// distinction between GetOption (returns None for zero-radius) and Set (which
// in the composed Optional still performs the write because Set in the inner
// Optional is unconditional).
//
// The Optional law only guarantees: if GetOption(s) = None then Set(b)(s) = s.
// When the prism matches but the inner Optional's GetOption returns None, Set
// behaviour depends on the inner Optional's Set implementation. In our
// positiveRadiusOptional the Set function is unconditional (it doesn't re-check
// the predicate), so Set still updates the field even when GetOption returns None.
func TestComposeOptional_GetOption_InnerMissDoesNotAffectSet(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	// GetOption returns None for zero radius.
	assert.True(t, O.IsNone(opt.GetOption(compCircle{Radius: 0})))
	// But Set still updates because the prism matched.
	result := opt.Set(10.0)(compCircle{Radius: 0})
	c, ok := result.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 10.0, c.Radius)
}

// TestComposeOptional_OptionalLaws verifies all three standard Optional laws.
func TestComposeOptional_OptionalLaws(t *testing.T) {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	matchingShape := compCircle{Radius: 5.0}
	missingShape := compRect{Width: 1, Height: 2}

	t.Run("law 1 – GetSet: no-op on None", func(t *testing.T) {
		assert.Equal(t, missingShape, opt.Set(99.0)(missingShape))
	})

	t.Run("law 2 – SetGet: get what you set", func(t *testing.T) {
		updated := opt.Set(42.0)(matchingShape)
		got := opt.GetOption(updated)
		assert.True(t, O.IsSome(got))
		assert.Equal(t, 42.0, O.GetOrElse(F.Constant(0.0))(got))
	})

	t.Run("law 3 – SetSet: last set wins", func(t *testing.T) {
		first := opt.Set(1.0)(matchingShape)
		both := opt.Set(2.0)(first)
		once := opt.Set(2.0)(matchingShape)
		assert.Equal(t, once, both)
	})
}

// TestComposeOptional_EquivalentToComposeLens verifies that ComposeOptional
// gives the same results as ComposeLens when the inner Optional is total
// (i.e., derived from the same Lens).
func TestComposeOptional_EquivalentToComposeLens(t *testing.T) {
	// Build an Optional that wraps the same getter/setter as compRadiusLens.
	totalOpt := common.MakeOptional(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(c compCircle, r float64) compCircle { c.Radius = r; return c },
	)

	lensResult := ComposeLens[compShape](compRadiusLens)(compCirclePrism)
	optResult := ComposeOptional[compShape](totalOpt)(compCirclePrism)

	s := compCircle{Radius: 6.0}
	assert.Equal(t, lensResult.GetOption(s), optResult.GetOption(s))
	assert.Equal(t, lensResult.Set(12.0)(s), optResult.Set(12.0)(s))
}

// ---------------------------------------------------------------------------
// Cross-function integration tests
// ---------------------------------------------------------------------------

// TestCompose_ThenComposeIso verifies that Compose followed by ComposeIso
// correctly drills into a nested structure.
func TestCompose_ThenComposeIso(t *testing.T) {
	strIso := common.MakeIso(
		func(r float64) string { return strconv.FormatFloat(r, 'f', 1, 64) },
		func(s string) float64 {
			f, _ := strconv.ParseFloat(s, 64)
			return f
		},
	)

	// Prism[compShape, compCircle] ∘ Iso[compCircle, compCircle] is trivial;
	// instead compose the circle-focused prism with an iso over float64.
	// Build a prism that pulls out a positive float64 from a compShape via compCircle.
	innerPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	byRadius := Compose[compShape](innerPrism)(compCirclePrism) // Prism[compShape, float64]
	withStr := ComposeIso[compShape](strIso)(byRadius)          // Prism[compShape, string]

	got := withStr.GetOption(compCircle{Radius: 3.5})
	assert.True(t, O.IsSome(got))
	assert.Equal(t, "3.5", O.GetOrElse(F.Constant(""))(got))

	s := withStr.ReverseGet("7.0")
	c, ok := s.(compCircle)
	assert.True(t, ok)
	assert.Equal(t, 7.0, c.Radius)
}

// ---------------------------------------------------------------------------
// Example functions — Compose
// ---------------------------------------------------------------------------

// ExampleCompose_match demonstrates that GetOption returns Some when both the
// outer and inner prisms match their respective variants.
func ExampleCompose_match() {
	// Prism[compShape, compCircle] ∘ Prism[compCircle, float64]
	radiusPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](radiusPrism)(compCirclePrism)

	fmt.Println(composed.GetOption(compCircle{Radius: 5}))
	// Output:
	// Some[float64](5)
}

// ExampleCompose_outerMiss demonstrates that GetOption returns None when the
// outer prism does not match.
func ExampleCompose_outerMiss() {
	radiusPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](radiusPrism)(compCirclePrism)

	fmt.Println(composed.GetOption(compRect{Width: 3, Height: 4}))
	// Output:
	// None[float64]
}

// ExampleCompose_innerMiss demonstrates that GetOption returns None when the
// outer prism matches but the inner prism does not.
func ExampleCompose_innerMiss() {
	// Inner prism only admits positive radii.
	positivePrism := MakePrism(
		func(c compCircle) O.Option[float64] {
			if c.Radius > 0 {
				return O.Some(c.Radius)
			}
			return O.None[float64]()
		},
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](positivePrism)(compCirclePrism)

	fmt.Println(composed.GetOption(compCircle{Radius: 0})) // outer matches, inner misses
	// Output:
	// None[float64]
}

// ExampleCompose_reverseGet demonstrates that ReverseGet threads the value
// through the inner prism's ReverseGet and then the outer prism's ReverseGet.
func ExampleCompose_reverseGet() {
	radiusPrism := MakePrism(
		func(c compCircle) O.Option[float64] { return O.Some(c.Radius) },
		func(r float64) compCircle { return compCircle{Radius: r} },
	)
	composed := Compose[compShape](radiusPrism)(compCirclePrism)

	s := composed.ReverseGet(7.0)
	fmt.Println(s)
	// Output:
	// {7}
}

// ExampleCompose_pipeline demonstrates the curried form used inside F.Pipe1
// to compose a head-of-slice prism with a string-to-int parsing prism.
func ExampleCompose_pipeline() {
	// Prism[[]string, string] ∘ Prism[string, int] → Prism[[]string, int]
	composed := F.Pipe1(Head[string](), Compose[[]string](ParseInt()))

	fmt.Println(composed.GetOption([]string{"42", "ignored"}))
	fmt.Println(composed.GetOption([]string{"abc"}))
	fmt.Println(composed.GetOption([]string{}))
	// Output:
	// Some[int](42)
	// None[int]
	// None[int]
}

// ExampleCompose_eitherChain demonstrates composing Either-based prisms to
// extract a value from a nested Either.
func ExampleCompose_eitherChain() {
	// Outer: Either[string, Either[string, int]] → Either[string, int]
	outerPrism := FromEither[string, Either[string, int]]()
	// Inner: Either[string, int] → int
	innerPrism := FromEither[string, int]()

	composed := F.Pipe1(outerPrism, Compose[Either[string, Either[string, int]]](innerPrism))

	nested := E.Right[string](E.Right[string](42))
	fmt.Println(composed.GetOption(nested))
	fmt.Println(E.IsRight(composed.ReverseGet(99)))
	// Output:
	// Some[int](42)
	// true
}

// ---------------------------------------------------------------------------
// Example functions — ComposeIso
// ---------------------------------------------------------------------------

// ExampleComposeIso_match demonstrates that GetOption returns Some and the
// iso's forward function is applied to the focused value.
func ExampleComposeIso_match() {
	// Iso: compCircle ↔ float64 via the Radius field.
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	fmt.Println(composed.GetOption(compCircle{Radius: 7.5}))
	// Output:
	// Some[float64](7.5)
}

// ExampleComposeIso_miss demonstrates that GetOption returns None when the
// prism does not match, regardless of the iso.
func ExampleComposeIso_miss() {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	fmt.Println(composed.GetOption(compRect{Width: 2, Height: 3}))
	// Output:
	// None[float64]
}

// ExampleComposeIso_reverseGet demonstrates that ReverseGet applies the iso's
// ReverseGet then the prism's ReverseGet to reconstruct the source.
func ExampleComposeIso_reverseGet() {
	composed := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism)

	s := composed.ReverseGet(5.0)
	fmt.Println(s)
	// Output:
	// {5}
}

// ExampleComposeIso_pipeline demonstrates the curried form used inside F.Pipe1.
func ExampleComposeIso_pipeline() {
	// Prism[compShape, compCircle] composed with Iso[compCircle, float64].
	composed := F.Pipe1(compCirclePrism, ComposeIso[compShape](circleToRadiusIso))

	fmt.Println(composed.GetOption(compCircle{Radius: 3.5}))
	fmt.Println(composed.GetOption(compRect{Width: 1, Height: 2}))
	// Output:
	// Some[float64](3.5)
	// None[float64]
}

// ExampleComposeIso_chained demonstrates chaining two ComposeIso calls to
// convert through multiple types: compShape → float64 → int.
func ExampleComposeIso_chained() {
	step1 := ComposeIso[compShape](circleToRadiusIso)(compCirclePrism) // Prism[compShape, float64]
	step2 := ComposeIso[compShape](radiusToIntIso)(step1)              // Prism[compShape, int]

	fmt.Println(step2.GetOption(compCircle{Radius: 9.9}))
	fmt.Println(step2.ReverseGet(3))
	// Output:
	// Some[int](9)
	// {3}
}

// ---------------------------------------------------------------------------
// Example functions — ComposeLens
// ---------------------------------------------------------------------------

// ExampleComposeLens_getOption demonstrates that GetOption returns the focused
// field value when the prism matches.
func ExampleComposeLens_getOption() {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	fmt.Println(opt.GetOption(compCircle{Radius: 4}))
	// Output:
	// Some[float64](4)
}

// ExampleComposeLens_miss demonstrates that GetOption returns None when the
// prism does not match.
func ExampleComposeLens_miss() {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	fmt.Println(opt.GetOption(compRect{Width: 2, Height: 3}))
	// Output:
	// None[float64]
}

// ExampleComposeLens_set demonstrates that Set updates the focused field when
// the prism matches.
func ExampleComposeLens_set() {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	result := opt.Set(10.0)(compCircle{Radius: 1.0})
	fmt.Println(result)
	// Output:
	// {10}
}

// ExampleComposeLens_setNoOp demonstrates that Set is a no-op when the prism
// does not match (Optional law 1 — GetSet).
func ExampleComposeLens_setNoOp() {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	original := compRect{Width: 2, Height: 3}
	result := opt.Set(99.0)(original)
	fmt.Println(result == original)
	// Output:
	// true
}

// ExampleComposeLens_pipeline demonstrates the curried form used inside F.Pipe1.
func ExampleComposeLens_pipeline() {
	opt := F.Pipe1(compCirclePrism, ComposeLens[compShape](compRadiusLens))

	fmt.Println(opt.GetOption(compCircle{Radius: 6}))
	// Output:
	// Some[float64](6)
}

// ExampleComposeLens_getOrElse demonstrates extracting the focused value with
// a fallback using O.GetOrElse and lazy.Of.
func ExampleComposeLens_getOrElse() {
	opt := ComposeLens[compShape](compRadiusLens)(compCirclePrism)

	radius := O.GetOrElse(lazy.Of(0.0))(opt.GetOption(compCircle{Radius: 8}))
	missing := O.GetOrElse(lazy.Of(0.0))(opt.GetOption(compRect{Width: 1, Height: 2}))
	fmt.Println(radius)
	fmt.Println(missing)
	// Output:
	// 8
	// 0
}

// ---------------------------------------------------------------------------
// Example functions — ComposeOptional
// ---------------------------------------------------------------------------

// ExampleComposeOptional_match demonstrates that GetOption returns Some when
// both the prism and the inner Optional match.
func ExampleComposeOptional_match() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	fmt.Println(opt.GetOption(compCircle{Radius: 5}))
	// Output:
	// Some[float64](5)
}

// ExampleComposeOptional_prismMiss demonstrates that GetOption returns None
// when the outer prism does not match.
func ExampleComposeOptional_prismMiss() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	fmt.Println(opt.GetOption(compRect{Width: 2, Height: 3}))
	// Output:
	// None[float64]
}

// ExampleComposeOptional_innerMiss demonstrates that GetOption returns None
// when the outer prism matches but the inner Optional does not.
func ExampleComposeOptional_innerMiss() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	fmt.Println(opt.GetOption(compCircle{Radius: 0})) // inner optional rejects zero
	// Output:
	// None[float64]
}

// ExampleComposeOptional_set demonstrates that Set updates the focused value
// when both the prism and inner Optional match.
func ExampleComposeOptional_set() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	result := opt.Set(10.0)(compCircle{Radius: 3.0})
	fmt.Println(result)
	// Output:
	// {10}
}

// ExampleComposeOptional_setNoOp demonstrates that Set is a no-op when the
// prism does not match (Optional law 1 — GetSet).
func ExampleComposeOptional_setNoOp() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	original := compRect{Width: 2, Height: 3}
	result := opt.Set(99.0)(original)
	fmt.Println(result == original)
	// Output:
	// true
}

// ExampleComposeOptional_pipeline demonstrates the curried form used inside
// F.Pipe1.
func ExampleComposeOptional_pipeline() {
	opt := F.Pipe1(compCirclePrism, ComposeOptional[compShape](positiveRadiusOptional))

	fmt.Println(opt.GetOption(compCircle{Radius: 8}))
	fmt.Println(opt.GetOption(compCircle{Radius: 0}))
	// Output:
	// Some[float64](8)
	// None[float64]
}

// ExampleComposeOptional_getOrElse demonstrates extracting the focused value
// with a fallback using O.GetOrElse and lazy.Of.
func ExampleComposeOptional_getOrElse() {
	opt := ComposeOptional[compShape](positiveRadiusOptional)(compCirclePrism)

	radius := O.GetOrElse(lazy.Of(0.0))(opt.GetOption(compCircle{Radius: 4}))
	missing := O.GetOrElse(lazy.Of(0.0))(opt.GetOption(compCircle{Radius: 0}))
	fmt.Println(radius)
	fmt.Println(missing)
	// Output:
	// 4
	// 0
}
