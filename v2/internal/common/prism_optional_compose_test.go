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

package common

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// ---------------------------------------------------------------------------
// Test fixtures
//
// We reuse the Shape / Circle / Rectangle / Triangle sum type defined in
// prism_test.go (same package).
//
// circlePrismPOC focuses on the Circle variant and exposes its Radius.
// positiveRadiusOpt is an Optional[float64, float64] that matches only when
// the radius is strictly positive, giving a genuine two-level partial optic.
// ---------------------------------------------------------------------------

// circlePrismPOC focuses on the Circle variant of Shape and extracts its Radius.
var circlePrismPOC = MakePrism(
	func(s Shape) Option[float64] {
		if c, ok := s.(Circle); ok {
			return OptionSome(c.Radius)
		}
		return OptionNone[float64]()
	},
	func(r float64) Shape { return Circle{Radius: r} },
)

// positiveRadiusOpt is an Optional[float64, float64] that matches only when
// the radius is strictly positive.  Set is unconditional (it always stores the
// value), while GetOption rejects non-positive values.
var positiveRadiusOpt = MakeOptional(
	func(r float64) Option[float64] {
		if r > 0 {
			return OptionSome(r)
		}
		return OptionNone[float64]()
	},
	func(r float64, v float64) float64 { return v },
)

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_EquivalentToManualComposition verifies that
// PrismComposeOptional[S](opt)(p) produces the same result as manually
// composing PrismAsOptional with OptionalComposeOptional.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_EquivalentToManualComposition(t *testing.T) {
	composed := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)
	manual := F.Pipe1(
		PrismAsOptional(circlePrismPOC),
		OptionalComposeOptional[Shape](positiveRadiusOpt),
	)

	cases := []Shape{
		Circle{Radius: 3.0},
		Circle{Radius: 0},
		Rectangle{Width: 4, Height: 5},
	}
	for _, s := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, manual.GetOption(s), composed.GetOption(s))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, manual.Set(9.0)(s), composed.Set(9.0)(s))
		})
	}
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_GetOption tests the GetOption direction.
// ---------------------------------------------------------------------------

// TestPrismComposeOptional_GetOption_Match verifies Some is returned when both
// the prism and the inner Optional match.
func TestPrismComposeOptional_GetOption_Match(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	got := opt.GetOption(Circle{Radius: 7.5})
	assert.Equal(t, OptionSome(7.5), got)
}

// TestPrismComposeOptional_GetOption_PrismMiss verifies None is returned when
// the prism does not match, regardless of the inner Optional.
func TestPrismComposeOptional_GetOption_PrismMiss(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	assert.Equal(t, OptionNone[float64](), opt.GetOption(Rectangle{Width: 4, Height: 5}))
	assert.Equal(t, OptionNone[float64](), opt.GetOption(Triangle{Base: 3, Height: 4}))
}

// TestPrismComposeOptional_GetOption_InnerMiss verifies None is returned when
// the prism matches but the inner Optional does not.
func TestPrismComposeOptional_GetOption_InnerMiss(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	// Prism matches (it is a Circle), but inner Optional rejects radius == 0.
	assert.Equal(t, OptionNone[float64](), opt.GetOption(Circle{Radius: 0}))
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_Set tests the Set direction.
// ---------------------------------------------------------------------------

// TestPrismComposeOptional_Set_Match verifies that Set updates the value when
// the prism matches.
func TestPrismComposeOptional_Set_Match(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	updated := opt.Set(99.0)(Circle{Radius: 1.0})
	assert.Equal(t, Circle{Radius: 99.0}, updated)
}

// TestPrismComposeOptional_Set_NoOpOnPrismMiss verifies that Set is a no-op when
// the prism does not match.
func TestPrismComposeOptional_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	original := Shape(Rectangle{Width: 4, Height: 5})
	result := opt.Set(99.0)(original)
	assert.Equal(t, original, result)
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_OptionalLaws verifies all three Optional laws.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_OptionalLaws(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	matching := Shape(Circle{Radius: 3.0})
	nonMatching := Shape(Rectangle{Width: 1, Height: 2})

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := opt.Set(99.0)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := opt.Set(42.0)(matching)
		got := opt.GetOption(updated)
		assert.Equal(t, OptionSome(42.0), got)
	})

	t.Run("law SetSet: Set(c)(Set(b)(s)) == Set(c)(s)", func(t *testing.T) {
		result1 := opt.Set(20.0)(opt.Set(10.0)(matching))
		result2 := opt.Set(20.0)(matching)
		assert.Equal(t, result2, result1)
	})
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_TotalInnerOptional verifies that when the inner
// Optional is total (always matches), PrismComposeOptional behaves identically
// to PrismComposeLens with an equivalent Lens.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_TotalInnerOptional(t *testing.T) {
	// A total Optional[float64, float64] that always succeeds.
	totalOpt := MakeOptional(
		func(r float64) Option[float64] { return OptionSome(r) },
		func(r float64, v float64) float64 { return v },
	)
	// Equivalent Lens[float64, float64] (identity).
	idLens := MakeLens(
		func(r float64) float64 { return r },
		func(r float64, v float64) float64 { return v },
	)

	viaOptional := PrismComposeOptional[Shape, float64, float64](totalOpt)(circlePrismPOC)
	viaLens := PrismComposeLens[Shape, float64, float64](idLens)(circlePrismPOC)

	cases := []Shape{
		Circle{Radius: 5.0},
		Rectangle{Width: 3, Height: 4},
	}
	for _, s := range cases {
		assert.Equal(t, viaLens.GetOption(s), viaOptional.GetOption(s), "GetOption parity for %v", s)
		assert.Equal(t, viaLens.Set(8.0)(s), viaOptional.Set(8.0)(s), "Set parity for %v", s)
	}
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_NonIdentityOptional verifies composition with an
// inner Optional that changes the focus type.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_NonIdentityOptional(t *testing.T) {
	// Optional[float64, int] that truncates the radius and rejects non-positives.
	truncPositiveOpt := MakeOptional(
		func(r float64) Option[int] {
			if r > 0 {
				return OptionSome(int(r))
			}
			return OptionNone[int]()
		},
		func(r float64, v int) float64 { return float64(v) },
	)
	opt := PrismComposeOptional[Shape, float64, int](truncPositiveOpt)(circlePrismPOC)

	t.Run("Get truncates positive radius to int", func(t *testing.T) {
		got := opt.GetOption(Circle{Radius: 7.9})
		assert.Equal(t, OptionSome(7), got)
	})

	t.Run("Get returns None for zero radius (inner miss)", func(t *testing.T) {
		assert.Equal(t, OptionNone[int](), opt.GetOption(Circle{Radius: 0}))
	})

	t.Run("Get returns None for non-Circle variant (outer miss)", func(t *testing.T) {
		assert.Equal(t, OptionNone[int](), opt.GetOption(Rectangle{Width: 4, Height: 5}))
	})

	t.Run("Set applies the inner Optional's setter", func(t *testing.T) {
		updated := opt.Set(10)(Circle{Radius: 1.5})
		assert.Equal(t, Circle{Radius: 10.0}, updated)
	})

	t.Run("Set is a no-op when prism misses", func(t *testing.T) {
		original := Shape(Rectangle{Width: 3, Height: 4})
		result := opt.Set(10)(original)
		assert.Equal(t, original, result)
	})
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_PipelineUsage verifies that the curried form is
// compatible with F.Pipe1.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(
		circlePrismPOC,
		PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt),
	)
	directOpt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)

	s := Shape(Circle{Radius: 2.5})
	assert.Equal(t, pipelineOpt.GetOption(s), directOpt.GetOption(s))
	assert.Equal(t, pipelineOpt.Set(8.0)(s), directOpt.Set(8.0)(s))
}

// ---------------------------------------------------------------------------
// TestPrismComposeOptional_ModifyOption verifies that OptionalModifyOption
// returns None when the prism misses and Some with the modified value otherwise.
// ---------------------------------------------------------------------------

func TestPrismComposeOptional_ModifyOption(t *testing.T) {
	opt := PrismComposeOptional[Shape, float64, float64](positiveRadiusOpt)(circlePrismPOC)
	double := func(r float64) float64 { return r * 2 }

	t.Run("returns Some with modified value when both optics match", func(t *testing.T) {
		result := OptionalModifyOption[Shape](double)(opt)(Circle{Radius: 3.0})
		assert.Equal(t, OptionSome[Shape](Circle{Radius: 6.0}), result)
	})

	t.Run("returns None when prism does not match", func(t *testing.T) {
		result := OptionalModifyOption[Shape](double)(opt)(Rectangle{Width: 4, Height: 5})
		assert.Equal(t, OptionNone[Shape](), result)
	})

	t.Run("returns None when inner Optional does not match", func(t *testing.T) {
		result := OptionalModifyOption[Shape](double)(opt)(Circle{Radius: 0})
		assert.Equal(t, OptionNone[Shape](), result)
	})
}
