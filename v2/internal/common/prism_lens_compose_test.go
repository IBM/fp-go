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
// prism_test.go (same package).  The Prism focuses on the Circle variant
// and extracts its Radius; the Lens then focuses on that radius value.
// ---------------------------------------------------------------------------

// circlePrismPLC focuses on the Circle variant of Shape and extracts its Radius.
var circlePrismPLC = MakePrism(
	func(s Shape) Option[float64] {
		if c, ok := s.(Circle); ok {
			return OptionSome(c.Radius)
		}
		return OptionNone[float64]()
	},
	func(r float64) Shape { return Circle{Radius: r} },
)

// radiusLens focuses on the float64 value directly (identity lens on float64),
// giving PrismComposeLens an A→B lens to compose with.
var radiusLens = MakeLens(
	func(r float64) float64 { return r },
	func(r float64, v float64) float64 { return v },
)

// ---------------------------------------------------------------------------
// TestPrismComposeLens_EquivalentToManualComposition verifies that
// PrismComposeLens[S, A, B](l)(p) produces the same Optional as manually
// composing PrismAsOptional with OptionalComposeLens.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_EquivalentToManualComposition(t *testing.T) {
	composed := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)
	manual := F.Pipe1(
		PrismAsOptional(circlePrismPLC),
		OptionalComposeLens[Shape](radiusLens),
	)

	cases := []Shape{
		Circle{Radius: 3.0},
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
// TestPrismComposeLens_GetOption tests the GetOption direction.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_GetOption(t *testing.T) {
	opt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)

	t.Run("returns Some when prism matches", func(t *testing.T) {
		got := opt.GetOption(Circle{Radius: 7.5})
		assert.Equal(t, OptionSome(7.5), got)
	})

	t.Run("returns None when prism does not match", func(t *testing.T) {
		got := opt.GetOption(Rectangle{Width: 4, Height: 5})
		assert.Equal(t, OptionNone[float64](), got)
	})

	t.Run("returns None for a different non-matching variant", func(t *testing.T) {
		got := opt.GetOption(Triangle{Base: 3, Height: 4})
		assert.Equal(t, OptionNone[float64](), got)
	})
}

// ---------------------------------------------------------------------------
// TestPrismComposeLens_Set tests the Set direction.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_Set(t *testing.T) {
	opt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)

	t.Run("updates the field when prism matches", func(t *testing.T) {
		updated := opt.Set(99.0)(Circle{Radius: 1.0})
		assert.Equal(t, Circle{Radius: 99.0}, updated)
	})

	t.Run("is a no-op when prism does not match (GetSet law)", func(t *testing.T) {
		original := Shape(Rectangle{Width: 4, Height: 5})
		result := opt.Set(99.0)(original)
		assert.Equal(t, original, result)
	})

	t.Run("is a no-op for every non-matching variant", func(t *testing.T) {
		original := Shape(Triangle{Base: 3, Height: 4})
		result := opt.Set(99.0)(original)
		assert.Equal(t, original, result)
	})
}

// ---------------------------------------------------------------------------
// TestPrismComposeLens_OptionalLaws verifies all three standard Optional laws.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_OptionalLaws(t *testing.T) {
	opt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)

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
// TestPrismComposeLens_Immutability verifies that Set does not mutate the
// original value (pass-by-value semantics for the Shape interface).
// ---------------------------------------------------------------------------

func TestPrismComposeLens_Immutability(t *testing.T) {
	opt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)

	original := Shape(Circle{Radius: 5.0})
	_ = opt.Set(99.0)(original)
	// original is an interface value wrapping a by-value Circle; it must be unchanged.
	assert.Equal(t, Circle{Radius: 5.0}, original)
}

// ---------------------------------------------------------------------------
// TestPrismComposeLens_NonIdentityLens verifies composition with a Lens that
// performs a real transformation (scaling the radius by a factor).
// ---------------------------------------------------------------------------

func TestPrismComposeLens_NonIdentityLens(t *testing.T) {
	// Lens: float64 (raw radius) → int (rounded radius in whole units).
	// Get rounds down; Set stores the int as float64.
	roundLens := MakeLens(
		func(r float64) int { return int(r) },
		func(r float64, v int) float64 { return float64(v) },
	)

	opt := PrismComposeLens[Shape, float64, int](roundLens)(circlePrismPLC)

	t.Run("Get applies the lens transformation", func(t *testing.T) {
		got := opt.GetOption(Circle{Radius: 7.9})
		assert.Equal(t, OptionSome(7), got)
	})

	t.Run("Set applies the lens transformation in reverse", func(t *testing.T) {
		updated := opt.Set(10)(Circle{Radius: 1.5})
		assert.Equal(t, Circle{Radius: 10.0}, updated)
	})

	t.Run("Set is still no-op when prism misses", func(t *testing.T) {
		original := Shape(Rectangle{Width: 3, Height: 4})
		result := opt.Set(10)(original)
		assert.Equal(t, original, result)
	})
}

// ---------------------------------------------------------------------------
// TestPrismComposeLens_PipelineUsage verifies that the curried form is
// compatible with F.Pipe1.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(
		circlePrismPLC,
		PrismComposeLens[Shape, float64, float64](radiusLens),
	)
	directOpt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)

	s := Shape(Circle{Radius: 2.5})
	assert.Equal(t, pipelineOpt.GetOption(s), directOpt.GetOption(s))
	assert.Equal(t, pipelineOpt.Set(8.0)(s), directOpt.Set(8.0)(s))
}

// ---------------------------------------------------------------------------
// TestPrismComposeLens_ModifyOption verifies that OptionalModifyOption correctly
// returns None when the prism misses, and Some with the modified value otherwise.
// ---------------------------------------------------------------------------

func TestPrismComposeLens_ModifyOption(t *testing.T) {
	opt := PrismComposeLens[Shape, float64, float64](radiusLens)(circlePrismPLC)
	double := func(r float64) float64 { return r * 2 }

	t.Run("returns Some with modified value when prism matches", func(t *testing.T) {
		result := OptionalModifyOption[Shape](double)(opt)(Circle{Radius: 3.0})
		assert.Equal(t, OptionSome[Shape](Circle{Radius: 6.0}), result)
	})

	t.Run("returns None when prism does not match", func(t *testing.T) {
		result := OptionalModifyOption[Shape](double)(opt)(Rectangle{Width: 4, Height: 5})
		assert.Equal(t, OptionNone[Shape](), result)
	})
}
