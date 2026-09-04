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

package common

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// TestPrismCompose_MethodEquivalentToFreeFunction verifies that the method form
// p.Compose(ab) returns a prism identical in behaviour to the free function
// Compose[S](ab)(p) for every operation.
func TestPrismCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	// outer: Option[int] → int  (extracts from Some)
	outer := PrismFromOption[int]()
	// inner: int → int  (matches only positive numbers)
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })

	method := outer.Compose(inner)
	free := PrismComposePrism[Option[int]](inner)(outer)

	cases := []Option[int]{OptionSome(42), OptionSome(-1), OptionNone[int]()}
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
	outer := PrismFromOption[int]()
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	got := composed.GetOption(OptionSome(42))
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 42, OptionGetOrElse(F.Constant(-1))(got))
}

// TestPrismCompose_Method_GetOption_OuterMiss verifies None is returned when
// the outer prism does not match.
func TestPrismCompose_Method_GetOption_OuterMiss(t *testing.T) {
	outer := PrismFromOption[int]()
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	assert.True(t, OptionIsNone(composed.GetOption(OptionNone[int]())))
}

// TestPrismCompose_Method_GetOption_InnerMiss verifies None is returned when
// the outer prism matches but the inner does not.
func TestPrismCompose_Method_GetOption_InnerMiss(t *testing.T) {
	outer := PrismFromOption[int]()
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	assert.True(t, OptionIsNone(composed.GetOption(OptionSome(-5))))
}

// TestPrismCompose_Method_ReverseGet verifies that ReverseGet threads the value
// back through the inner then outer ReverseGet.
func TestPrismCompose_Method_ReverseGet(t *testing.T) {
	outer := PrismFromOption[int]()
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	// ReverseGet(42): inner.ReverseGet(42)==42, outer.ReverseGet(42)==Some(42)
	assert.Equal(t, OptionSome(42), composed.ReverseGet(42))
}

// TestPrismCompose_Method_PrismLaws verifies the two standard prism laws on the
// method-composed prism.
func TestPrismCompose_Method_PrismLaws(t *testing.T) {
	outer := PrismFromOption[int]()
	inner := PrismFromPredicate(func(n int) bool { return n > 0 })
	composed := outer.Compose(inner)

	t.Run("law 1: GetOption(ReverseGet(a)) == Some(a)", func(t *testing.T) {
		a := 99
		got := composed.GetOption(composed.ReverseGet(a))
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, a, OptionGetOrElse(F.Constant(-1))(got))
	})

	t.Run("law 2: if GetOption(s)==Some(a) then ReverseGet(a) round-trips", func(t *testing.T) {
		s := OptionSome(7)
		got := composed.GetOption(s)
		if OptionIsSome(got) {
			a := OptionGetOrElse(F.Constant(0))(got)
			reconstructed := composed.ReverseGet(a)
			assert.Equal(t, OptionSome(a), composed.GetOption(reconstructed))
		}
	})
}

// TestPrismCompose_Method_Chained verifies that multiple .Compose calls chain
// correctly through three levels of sum type.
func TestPrismCompose_Method_Chained(t *testing.T) {
	// Level 1: Option[Option[int]] → Option[int]  (unwrap outer Some)
	l1 := PrismFromOption[Option[int]]()
	// Level 2: Option[int] → int                  (unwrap inner Some)
	l2 := PrismFromOption[int]()
	// Level 3: int → int                          (only positives)
	l3 := PrismFromPredicate(func(n int) bool { return n > 0 })

	// Chain: Option[Option[int]] → Option[int] → int → int (positive)
	composed := l1.Compose(l2).Compose(l3)

	t.Run("all three levels match", func(t *testing.T) {
		s := OptionSome(OptionSome(5))
		got := composed.GetOption(s)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 5, OptionGetOrElse(F.Constant(-1))(got))
	})

	t.Run("outermost None short-circuits", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(OptionNone[Option[int]]())))
	})

	t.Run("middle None short-circuits", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(OptionSome(OptionNone[int]()))))
	})

	t.Run("innermost predicate miss short-circuits", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(OptionSome(OptionSome(-3)))))
	})

	t.Run("ReverseGet reconstructs all three layers", func(t *testing.T) {
		// inner.ReverseGet(10)==10 → l2.ReverseGet(10)==Some(10) → l1.ReverseGet(Some(10))==Some(Some(10))
		assert.Equal(t, OptionSome(OptionSome(10)), composed.ReverseGet(10))
	})
}

// ---- ComposeLens method tests ----
//
// Fixtures:
//   outer prism: Shape → float64  (matches Circle, extracts Radius via ReverseGet)
//   inner lens:  float64 → int    (truncates radius to whole units)
//
// The prism matches Circle; Rectangle and Triangle cause GetOption to return None.

// circlePrismG127 focuses on the Circle variant of Shape.
// (Defined here rather than reusing the package-level circlePrismPLC from
// prism_lens_compose_test.go to keep the test file self-contained.)
var circlePrismG127 = MakePrism(
	func(s Shape) Option[float64] {
		if c, ok := s.(Circle); ok {
			return OptionSome(c.Radius)
		}
		return OptionNone[float64]()
	},
	func(r float64) Shape { return Circle{Radius: r} },
)

// truncLens focuses on float64 by truncating to int (Get) and storing as float64 (Set).
var truncLens = MakeLens(
	func(r float64) int { return int(r) },
	func(r float64, v int) float64 { return float64(v) },
)

// TestPrismComposeLens_MethodEquivalentToFreeFunction verifies that the method
// form p.ComposeLens(ab) returns an Optional identical in behaviour to the
// free-function PrismComposeLens[S](ab)(p).
func TestPrismComposeLens_MethodEquivalentToFreeFunction(t *testing.T) {
	method := circlePrismG127.ComposeLens(truncLens)
	free := PrismComposeLens[Shape, float64, int](truncLens)(circlePrismG127)

	cases := []Shape{
		Circle{Radius: 3.7},
		Rectangle{Width: 4, Height: 5},
	}
	for _, s := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(s), method.GetOption(s))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(9)(s), method.Set(9)(s))
		})
	}
}

// TestPrismComposeLens_Method_GetOption_Match verifies Some is returned when
// the prism matches and the lens retrieves the truncated radius.
func TestPrismComposeLens_Method_GetOption_Match(t *testing.T) {
	opt := circlePrismG127.ComposeLens(truncLens)

	got := opt.GetOption(Circle{Radius: 7.9})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 7, OptionGetOrElse(F.Constant(-1))(got))
}

// TestPrismComposeLens_Method_GetOption_PrismMiss verifies None is returned
// when the prism does not match the source variant.
func TestPrismComposeLens_Method_GetOption_PrismMiss(t *testing.T) {
	opt := circlePrismG127.ComposeLens(truncLens)

	assert.True(t, OptionIsNone(opt.GetOption(Rectangle{Width: 4, Height: 5})))
	assert.True(t, OptionIsNone(opt.GetOption(Triangle{Base: 3, Height: 4})))
}

// TestPrismComposeLens_Method_Set_Match verifies that Set updates the lens
// focus when the prism matches.
func TestPrismComposeLens_Method_Set_Match(t *testing.T) {
	opt := circlePrismG127.ComposeLens(truncLens)

	updated := opt.Set(10)(Circle{Radius: 1.5})
	assert.Equal(t, Circle{Radius: 10.0}, updated)
}

// TestPrismComposeLens_Method_Set_NoOpOnPrismMiss verifies that Set is a no-op
// when the prism does not match (Optional no-op law / GetSet law).
func TestPrismComposeLens_Method_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := circlePrismG127.ComposeLens(truncLens)

	original := Shape(Rectangle{Width: 4, Height: 5})
	result := opt.Set(99)(original)
	assert.Equal(t, original, result)
}

// TestPrismComposeLens_Method_OptionalLaws verifies the three standard Optional
// laws on the method-composed Optional.
func TestPrismComposeLens_Method_OptionalLaws(t *testing.T) {
	opt := circlePrismG127.ComposeLens(truncLens)

	matching := Shape(Circle{Radius: 3.0})
	nonMatching := Shape(Rectangle{Width: 1, Height: 2})

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := opt.Set(99)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := opt.Set(42)(matching)
		got := opt.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 42, OptionGetOrElse(F.Constant(-1))(got))
	})

	t.Run("law SetSet: Set(c)(Set(b)(s)) == Set(c)(s)", func(t *testing.T) {
		result1 := opt.Set(20)(opt.Set(10)(matching))
		result2 := opt.Set(20)(matching)
		assert.Equal(t, result2, result1)
	})
}

// TestPrismComposeLens_Method_Chained verifies that ComposeLens can be chained
// after Compose (prism ∘ prism ∘ lens) to navigate three levels.
func TestPrismComposeLens_Method_Chained(t *testing.T) {
	// level 1 (Compose): Option[Shape] → Shape  (unwrap Some)
	// level 2 (ComposeLens): Shape → int  (circle radius, truncated)
	outerPrism := PrismFromOption[Shape]()

	opt := outerPrism.Compose(circlePrismG127).ComposeLens(truncLens)

	t.Run("Get navigates through Compose+ComposeLens", func(t *testing.T) {
		got := opt.GetOption(OptionSome[Shape](Circle{Radius: 5.9}))
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 5, OptionGetOrElse(F.Constant(-1))(got))
	})

	t.Run("Set updates through Compose+ComposeLens", func(t *testing.T) {
		s := OptionSome[Shape](Circle{Radius: 1.0})
		updated := opt.Set(7)(s)
		assert.Equal(t, OptionSome[Shape](Circle{Radius: 7.0}), updated)
	})

	t.Run("Set is no-op when outer prism misses", func(t *testing.T) {
		original := OptionNone[Shape]()
		result := opt.Set(99)(original)
		assert.Equal(t, original, result)
	})

	t.Run("Set is no-op when inner prism misses", func(t *testing.T) {
		original := OptionSome[Shape](Rectangle{Width: 3, Height: 4})
		result := opt.Set(99)(original)
		assert.Equal(t, original, result)
	})
}

// TestPrismComposeLens_Method_PipelineUsage verifies that the method form and
// the free-function form used inside F.Pipe1 produce identical results.
func TestPrismComposeLens_Method_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(
		circlePrismG127,
		PrismComposeLens[Shape, float64, int](truncLens),
	)
	methodOpt := circlePrismG127.ComposeLens(truncLens)

	s := Shape(Circle{Radius: 2.5})
	assert.Equal(t, pipelineOpt.GetOption(s), methodOpt.GetOption(s))
	assert.Equal(t, pipelineOpt.Set(8)(s), methodOpt.Set(8)(s))
}

// ---- PrismComposeIso free-function tests ----

// radiusToIntIso converts float64 (radius) to int by truncation and back.
// Get(ReverseGet(n)) == n always; ReverseGet(Get(r)) preserves length, not exact value.
var radiusToIntIso = MakeIso(
	func(r float64) int { return int(r) },
	func(n int) float64 { return float64(n) },
)

// TestPrismComposeIso_GetOption_Match verifies that GetOption maps the iso's
// forward function over the extracted value when the prism matches.
func TestPrismComposeIso_GetOption_Match(t *testing.T) {
	composed := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	got := composed.GetOption(Circle{Radius: 7.9})
	assert.Equal(t, OptionSome(7), got)
}

// TestPrismComposeIso_GetOption_Miss verifies that GetOption returns None when
// the prism does not match, regardless of the iso.
func TestPrismComposeIso_GetOption_Miss(t *testing.T) {
	composed := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	assert.Equal(t, OptionNone[int](), composed.GetOption(Rectangle{Width: 4, Height: 5}))
	assert.Equal(t, OptionNone[int](), composed.GetOption(Triangle{Base: 3, Height: 4}))
}

// TestPrismComposeIso_ReverseGet verifies that ReverseGet threads the value
// back through the iso's ReverseGet and then the prism's ReverseGet.
func TestPrismComposeIso_ReverseGet(t *testing.T) {
	composed := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	// iso.ReverseGet(5) == 5.0, prism.ReverseGet(5.0) == Circle{Radius: 5.0}
	got := composed.ReverseGet(5)
	assert.Equal(t, Circle{Radius: 5.0}, got)
}

// TestPrismComposeIso_PrismLaws verifies both prism laws on the composed prism.
func TestPrismComposeIso_PrismLaws(t *testing.T) {
	composed := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	t.Run("law 1: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		b := 10
		got := composed.GetOption(composed.ReverseGet(b))
		assert.Equal(t, OptionSome(b), got)
	})

	t.Run("law 2: if GetOption(s)==Some(b) then ReverseGet(b) round-trips", func(t *testing.T) {
		s := Shape(Circle{Radius: 3.0})
		got := composed.GetOption(s)
		if OptionIsSome(got) {
			b := OptionGetOrElse(F.Constant(-1))(got)
			assert.Equal(t, OptionSome(b), composed.GetOption(composed.ReverseGet(b)))
		}
	})
}

// TestPrismComposeIso_PipelineUsage verifies the curried form is compatible
// with F.Pipe1.
func TestPrismComposeIso_PipelineUsage(t *testing.T) {
	pipelineResult := F.Pipe1(
		circlePrismG127,
		PrismComposeIso[Shape, float64, int](radiusToIntIso),
	)
	directResult := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	s := Shape(Circle{Radius: 4.7})
	assert.Equal(t, pipelineResult.GetOption(s), directResult.GetOption(s))
	assert.Equal(t, pipelineResult.ReverseGet(4), directResult.ReverseGet(4))
}

// ---- ComposeIso method tests ----

// TestPrismComposeIso_MethodEquivalentToFreeFunction verifies that the method
// form p.ComposeIso(ab) returns a Prism identical in behaviour to the
// free-function PrismComposeIso[S](ab)(p).
func TestPrismComposeIso_MethodEquivalentToFreeFunction(t *testing.T) {
	method := circlePrismG127.ComposeIso(radiusToIntIso)
	free := PrismComposeIso[Shape, float64, int](radiusToIntIso)(circlePrismG127)

	cases := []Shape{
		Circle{Radius: 3.7},
		Rectangle{Width: 4, Height: 5},
	}
	for _, s := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(s), method.GetOption(s))
		})
	}
	t.Run("ReverseGet parity", func(t *testing.T) {
		assert.Equal(t, free.ReverseGet(6), method.ReverseGet(6))
	})
}

// TestPrismComposeIso_Method_GetOption_Match verifies Some is returned when the
// prism matches and the iso maps the extracted value forward.
func TestPrismComposeIso_Method_GetOption_Match(t *testing.T) {
	composed := circlePrismG127.ComposeIso(radiusToIntIso)

	got := composed.GetOption(Circle{Radius: 7.9})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 7, OptionGetOrElse(F.Constant(-1))(got))
}

// TestPrismComposeIso_Method_GetOption_PrismMiss verifies None is returned
// when the prism does not match, regardless of the iso.
func TestPrismComposeIso_Method_GetOption_PrismMiss(t *testing.T) {
	composed := circlePrismG127.ComposeIso(radiusToIntIso)

	assert.True(t, OptionIsNone(composed.GetOption(Rectangle{Width: 4, Height: 5})))
	assert.True(t, OptionIsNone(composed.GetOption(Triangle{Base: 3, Height: 4})))
}

// TestPrismComposeIso_Method_ReverseGet verifies that ReverseGet threads the
// value through iso.ReverseGet then prism.ReverseGet.
func TestPrismComposeIso_Method_ReverseGet(t *testing.T) {
	composed := circlePrismG127.ComposeIso(radiusToIntIso)

	got := composed.ReverseGet(5)
	assert.Equal(t, Circle{Radius: 5.0}, got)
}

// TestPrismComposeIso_Method_PrismLaws verifies both prism laws on the
// method-composed prism.
func TestPrismComposeIso_Method_PrismLaws(t *testing.T) {
	composed := circlePrismG127.ComposeIso(radiusToIntIso)

	t.Run("law 1: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		b := 10
		got := composed.GetOption(composed.ReverseGet(b))
		assert.Equal(t, OptionSome(b), got)
	})

	t.Run("law 2: if GetOption(s)==Some(b) then ReverseGet(b) round-trips", func(t *testing.T) {
		s := Shape(Circle{Radius: 3.0})
		got := composed.GetOption(s)
		if OptionIsSome(got) {
			b := OptionGetOrElse(F.Constant(-1))(got)
			assert.Equal(t, OptionSome(b), composed.GetOption(composed.ReverseGet(b)))
		}
	})
}

// TestPrismComposeIso_Method_Chained verifies that ComposeIso can be chained
// after Compose (prism.Compose(prism).ComposeIso(iso)) to navigate three levels.
func TestPrismComposeIso_Method_Chained(t *testing.T) {
	// level 1 (Compose): Option[Shape] → Shape  (unwrap Some)
	// level 2 (ComposeIso): Shape → int  (circle radius, truncated via iso)
	outerPrism := PrismFromOption[Shape]()

	composed := outerPrism.Compose(circlePrismG127).ComposeIso(radiusToIntIso)

	t.Run("GetOption navigates through Compose+ComposeIso", func(t *testing.T) {
		got := composed.GetOption(OptionSome[Shape](Circle{Radius: 5.9}))
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 5, OptionGetOrElse(F.Constant(-1))(got))
	})

	t.Run("GetOption returns None when outer prism misses", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(OptionNone[Shape]())))
	})

	t.Run("GetOption returns None when inner prism misses", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(OptionSome[Shape](Rectangle{Width: 3, Height: 4}))))
	})

	t.Run("ReverseGet reconstructs all three layers", func(t *testing.T) {
		// iso.ReverseGet(7)==7.0, circlePrism.ReverseGet(7.0)==Circle{7.0}, outerPrism.ReverseGet(Circle{7.0})==Some(Circle{7.0})
		got := composed.ReverseGet(7)
		assert.Equal(t, OptionSome[Shape](Circle{Radius: 7.0}), got)
	})

	t.Run("law 1 holds after chaining: GetOption(ReverseGet(b)) == Some(b)", func(t *testing.T) {
		b := 4
		got := composed.GetOption(composed.ReverseGet(b))
		assert.Equal(t, OptionSome(b), got)
	})
}

// TestPrismComposeIso_Method_PipelineUsage verifies that the method form and the
// free-function form used inside F.Pipe1 produce identical results.
func TestPrismComposeIso_Method_PipelineUsage(t *testing.T) {
	pipelinePrism := F.Pipe1(
		circlePrismG127,
		PrismComposeIso[Shape, float64, int](radiusToIntIso),
	)
	methodPrism := circlePrismG127.ComposeIso(radiusToIntIso)

	s := Shape(Circle{Radius: 2.5})
	assert.Equal(t, pipelinePrism.GetOption(s), methodPrism.GetOption(s))
	assert.Equal(t, pipelinePrism.ReverseGet(8), methodPrism.ReverseGet(8))
}

// ---- ComposeOptional method tests ----
//
// Fixture: circlePrismG127 focuses on the Circle variant of Shape (defined
// above).  positiveRadiusOptG127 is an Optional[float64, float64] that matches
// only strictly-positive radii, giving a genuine two-level partial composition.

// positiveRadiusOptG127 is an Optional[float64, float64] that only matches
// when the radius is strictly positive.
var positiveRadiusOptG127 = MakeOptional(
	func(r float64) Option[float64] {
		if r > 0 {
			return OptionSome(r)
		}
		return OptionNone[float64]()
	},
	func(r float64, v float64) float64 { return v },
)

// TestPrismComposeOptional_MethodEquivalentToFreeFunction verifies that the
// method form p.ComposeOptional(ab) returns an Optional identical in behaviour
// to the free-function PrismComposeOptional[S](ab)(p).
func TestPrismComposeOptional_MethodEquivalentToFreeFunction(t *testing.T) {
	method := circlePrismG127.ComposeOptional(positiveRadiusOptG127)
	free := PrismComposeOptional[Shape, float64, float64](positiveRadiusOptG127)(circlePrismG127)

	cases := []Shape{
		Circle{Radius: 3.7},
		Circle{Radius: 0},
		Rectangle{Width: 4, Height: 5},
	}
	for _, s := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(s), method.GetOption(s))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(9.0)(s), method.Set(9.0)(s))
		})
	}
}

// TestPrismComposeOptional_Method_GetOption_Match verifies Some is returned when
// both the prism and the inner Optional match.
func TestPrismComposeOptional_Method_GetOption_Match(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	got := opt.GetOption(Circle{Radius: 7.9})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 7.9, OptionGetOrElse(F.Constant(-1.0))(got))
}

// TestPrismComposeOptional_Method_GetOption_PrismMiss verifies None is returned
// when the prism does not match the source variant.
func TestPrismComposeOptional_Method_GetOption_PrismMiss(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	assert.True(t, OptionIsNone(opt.GetOption(Rectangle{Width: 4, Height: 5})))
	assert.True(t, OptionIsNone(opt.GetOption(Triangle{Base: 3, Height: 4})))
}

// TestPrismComposeOptional_Method_GetOption_InnerMiss verifies None is returned
// when the prism matches but the inner Optional does not.
func TestPrismComposeOptional_Method_GetOption_InnerMiss(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	// Prism matches (Circle), but inner Optional rejects radius == 0.
	assert.True(t, OptionIsNone(opt.GetOption(Circle{Radius: 0})))
}

// TestPrismComposeOptional_Method_Set_Match verifies that Set updates the
// focused value when the prism matches.
func TestPrismComposeOptional_Method_Set_Match(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	updated := opt.Set(10.0)(Circle{Radius: 1.5})
	assert.Equal(t, Circle{Radius: 10.0}, updated)
}

// TestPrismComposeOptional_Method_Set_NoOpOnPrismMiss verifies that Set is a
// no-op when the prism does not match (Optional GetSet law).
func TestPrismComposeOptional_Method_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	original := Shape(Rectangle{Width: 4, Height: 5})
	result := opt.Set(99.0)(original)
	assert.Equal(t, original, result)
}

// TestPrismComposeOptional_Method_OptionalLaws verifies the three standard
// Optional laws on the method-composed Optional.
func TestPrismComposeOptional_Method_OptionalLaws(t *testing.T) {
	opt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	matching := Shape(Circle{Radius: 3.0})
	nonMatching := Shape(Rectangle{Width: 1, Height: 2})

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := opt.Set(99.0)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := opt.Set(42.0)(matching)
		got := opt.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 42.0, OptionGetOrElse(F.Constant(-1.0))(got))
	})

	t.Run("law SetSet: Set(c)(Set(b)(s)) == Set(c)(s)", func(t *testing.T) {
		result1 := opt.Set(20.0)(opt.Set(10.0)(matching))
		result2 := opt.Set(20.0)(matching)
		assert.Equal(t, result2, result1)
	})
}

// TestPrismComposeOptional_Method_Chained verifies that ComposeOptional can be
// chained after Compose (prism.Compose(prism).ComposeOptional(optional)) to
// navigate three levels.
func TestPrismComposeOptional_Method_Chained(t *testing.T) {
	// level 1 (Compose): Option[Shape] → Shape  (unwrap Some)
	// level 2 (ComposeOptional): Shape → float64  (circle radius, positive only)
	outerPrism := PrismFromOption[Shape]()

	opt := outerPrism.Compose(circlePrismG127).ComposeOptional(positiveRadiusOptG127)

	t.Run("Get navigates through Compose+ComposeOptional when all match", func(t *testing.T) {
		got := opt.GetOption(OptionSome[Shape](Circle{Radius: 5.9}))
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 5.9, OptionGetOrElse(F.Constant(-1.0))(got))
	})

	t.Run("Set updates through Compose+ComposeOptional", func(t *testing.T) {
		s := OptionSome[Shape](Circle{Radius: 1.0})
		updated := opt.Set(7.0)(s)
		assert.Equal(t, OptionSome[Shape](Circle{Radius: 7.0}), updated)
	})

	t.Run("Set is no-op when outer prism misses", func(t *testing.T) {
		original := OptionNone[Shape]()
		result := opt.Set(99.0)(original)
		assert.Equal(t, original, result)
	})

	t.Run("Set is no-op when inner prism misses", func(t *testing.T) {
		original := OptionSome[Shape](Rectangle{Width: 3, Height: 4})
		result := opt.Set(99.0)(original)
		assert.Equal(t, original, result)
	})
}

// TestPrismComposeOptional_Method_PipelineUsage verifies that the method form
// and the free-function form used inside F.Pipe1 produce identical results.
func TestPrismComposeOptional_Method_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(
		circlePrismG127,
		PrismComposeOptional[Shape, float64, float64](positiveRadiusOptG127),
	)
	methodOpt := circlePrismG127.ComposeOptional(positiveRadiusOptG127)

	s := Shape(Circle{Radius: 2.5})
	assert.Equal(t, pipelineOpt.GetOption(s), methodOpt.GetOption(s))
	assert.Equal(t, pipelineOpt.Set(8.0)(s), methodOpt.Set(8.0)(s))
}
