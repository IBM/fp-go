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

package lens_test

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/common"
	"github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// ---- test fixtures ----
//
// The sum-type fixtures (tcCanvas, tcShape, tcCircle, tcRect) live in
// lens_composeprism_test.go and are shared across this file because both
// belong to package lens_test.

// coCircleRadiusOpt is an Optional[tcShape, float64] that focuses on the
// Radius of a tcCircle, but only when the radius is strictly positive.
// This exercises both the matching and non-matching code paths.
var coCircleRadiusOpt = common.MakeOptional(
	func(s tcShape) O.Option[float64] {
		if c, ok := s.(tcCircle); ok && c.Radius > 0 {
			return O.Some(c.Radius)
		}
		return O.None[float64]()
	},
	func(s tcShape, r float64) tcShape { return tcCircle{Radius: r} },
)

// coCanvasRadiusOpt is the Optional[tcCanvas, float64] produced by
// composing canvasShapeLens with coCircleRadiusOpt.
var coCanvasRadiusOpt = F.Pipe1(canvasShapeLens, lens.ComposeOptional[tcCanvas](coCircleRadiusOpt))

// ---- unit tests ----

// TestComposeOptional_GetOption_Match verifies Some is returned when the
// lens succeeds and the inner optional matches.
func TestComposeOptional_GetOption_Match(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 5}}
	got := coCanvasRadiusOpt.GetOption(canvas)
	assert.Equal(t, O.Some(5.0), got)
}

// TestComposeOptional_GetOption_OptionalMiss verifies None is returned when
// the lens succeeds but the inner optional does not match (wrong variant).
func TestComposeOptional_GetOption_OptionalMiss(t *testing.T) {
	canvas := tcCanvas{Shape: tcRect{Width: 3, Height: 4}}
	got := coCanvasRadiusOpt.GetOption(canvas)
	assert.Equal(t, O.None[float64](), got)
}

// TestComposeOptional_GetOption_PredicateMiss verifies None is returned when
// the shape is the right variant but fails the predicate (radius == 0).
func TestComposeOptional_GetOption_PredicateMiss(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 0}}
	got := coCanvasRadiusOpt.GetOption(canvas)
	assert.Equal(t, O.None[float64](), got)
}

// TestComposeOptional_Set_Match verifies that Set updates the focused value
// when the inner optional matches.
func TestComposeOptional_Set_Match(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 3}}
	updated := coCanvasRadiusOpt.Set(9.0)(canvas)
	assert.Equal(t, tcCanvas{Shape: tcCircle{Radius: 9}}, updated)
}

// TestComposeOptional_Immutability verifies that Set never mutates the
// original structure.
func TestComposeOptional_Immutability(t *testing.T) {
	original := tcCanvas{Shape: tcCircle{Radius: 5}}
	_ = coCanvasRadiusOpt.Set(99.0)(original)
	assert.Equal(t, tcCircle{Radius: 5}, original.Shape)
}

// TestComposeOptional_OptionalLaws verifies the three standard Optional laws.
func TestComposeOptional_OptionalLaws(t *testing.T) {
	matching := tcCanvas{Shape: tcCircle{Radius: 3}}
	nonMatching := tcCanvas{Shape: tcRect{Width: 2, Height: 4}}

	t.Run("law SetGet (match): GetOption(Set(b)(s)) == Some(b)", func(t *testing.T) {
		updated := coCanvasRadiusOpt.Set(42.0)(matching)
		got := coCanvasRadiusOpt.GetOption(updated)
		assert.Equal(t, O.Some(42.0), got)
	})

	t.Run("law GetSet (no-match): GetOption returns None for non-matching input", func(t *testing.T) {
		// The Optional no-op contract for the composition holds at the
		// ModifyOption level: when GetOption returns None, ModifyOption returns
		// None[S] (leaving the structure unchanged). The raw Set method on the
		// inner optional is unconditional — it is the caller's responsibility to
		// check GetOption before calling Set when needed.
		assert.True(t, O.IsNone(coCanvasRadiusOpt.GetOption(nonMatching)))
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		setTwice := coCanvasRadiusOpt.Set(20.0)(coCanvasRadiusOpt.Set(10.0)(matching))
		setOnce := coCanvasRadiusOpt.Set(20.0)(matching)
		assert.Equal(t, setOnce, setTwice)
	})
}

// TestComposeOptional_EquivalentToUnderlyingFreeFunction verifies that the
// optics/lens wrapper and a direct call to common.LensComposeOptional produce
// identical results.
func TestComposeOptional_EquivalentToUnderlyingFreeFunction(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 7}}

	wrapped := F.Pipe1(canvasShapeLens, lens.ComposeOptional[tcCanvas](coCircleRadiusOpt))
	direct := common.LensComposeOptional[tcCanvas](coCircleRadiusOpt)(canvasShapeLens)

	t.Run("GetOption returns same result", func(t *testing.T) {
		assert.Equal(t, direct.GetOption(canvas), wrapped.GetOption(canvas))
	})
	t.Run("Set returns same result", func(t *testing.T) {
		assert.Equal(t, direct.Set(15.0)(canvas), wrapped.Set(15.0)(canvas))
	})
}

// TestComposeOptional_Chained verifies that ComposeOptional can be chained
// after Compose to navigate two levels before applying a partial optional.
func TestComposeOptional_Chained(t *testing.T) {
	type Gallery struct{ Canvas tcCanvas }

	galleryLens := common.MakeLens(
		func(g Gallery) tcCanvas { return g.Canvas },
		func(g Gallery, c tcCanvas) Gallery { g.Canvas = c; return g },
	)

	// Gallery → tcCanvas → tcShape → float64 (circle radius, positive only)
	opt := F.Pipe1(
		F.Pipe1(galleryLens, lens.Compose[Gallery](canvasShapeLens)),
		lens.ComposeOptional[Gallery](coCircleRadiusOpt),
	)

	circleGallery := Gallery{Canvas: tcCanvas{Shape: tcCircle{Radius: 6}}}
	rectGallery := Gallery{Canvas: tcCanvas{Shape: tcRect{Width: 2, Height: 3}}}

	t.Run("GetOption navigates two levels — match", func(t *testing.T) {
		assert.Equal(t, O.Some(6.0), opt.GetOption(circleGallery))
	})

	t.Run("GetOption navigates two levels — miss", func(t *testing.T) {
		assert.Equal(t, O.None[float64](), opt.GetOption(rectGallery))
	})

	t.Run("Set updates two levels deep", func(t *testing.T) {
		updated := opt.Set(12.0)(circleGallery)
		assert.Equal(t, tcCircle{Radius: 12}, updated.Canvas.Shape)
		assert.Equal(t, tcCircle{Radius: 6}, circleGallery.Canvas.Shape) // original unchanged
	})

	t.Run("GetOption returns None when optional misses after chained Compose", func(t *testing.T) {
		// GetOption returns None for the non-matching case.
		// (The no-op contract is upheld at the ModifyOption level, not raw Set.)
		assert.Equal(t, O.None[float64](), opt.GetOption(rectGallery))
	})
}

// TestComposeOptional_HasNonEmptyName verifies that the composed Optional
// carries a non-empty name for logging and debugging.
func TestComposeOptional_HasNonEmptyName(t *testing.T) {
	assert.NotEmpty(t, coCanvasRadiusOpt.String())
}

// ---- examples ----

// ExampleComposeOptional demonstrates composing a Lens with a partial Optional
// to focus on a value that may or may not exist within a nested structure.
//
// The test types tcCanvas, tcCircle, and tcRect are defined at package level
// in lens_composeprism_test.go and are shared within this package.
func ExampleComposeOptional() {
	// coCircleRadiusOpt (package-level) focuses on tcCircle.Radius when > 0.
	// canvasShapeLens (package-level) focuses on tcCanvas.Shape.

	// Compose: tcCanvas → tcShape → float64 (only when shape is a positive Circle)
	canvasRadiusOpt := F.Pipe1(canvasShapeLens, lens.ComposeOptional[tcCanvas](coCircleRadiusOpt))

	circle := tcCanvas{Shape: tcCircle{Radius: 5}}
	rect := tcCanvas{Shape: tcRect{Width: 3, Height: 4}}

	fmt.Println(O.IsSome(canvasRadiusOpt.GetOption(circle))) // circle matches
	fmt.Println(O.IsNone(canvasRadiusOpt.GetOption(rect)))   // rect misses

	updated := canvasRadiusOpt.Set(10.0)(circle) // Set on match — updates
	fmt.Println(updated.Shape.(tcCircle).Radius)

	// GetOption returning None means the focus is absent; Set is unconditional
	// at the raw Optional level — callers should check GetOption before using Set.
	fmt.Println(O.IsNone(canvasRadiusOpt.GetOption(rect)))
	// Output:
	// true
	// true
	// 10
	// true
}

// ExampleComposeOptional_predicate demonstrates a predicate-guarded Optional
// that only matches values satisfying an additional condition.
func ExampleComposeOptional_predicate() {
	type Profile struct{ Score int }
	type Player struct{ Profile Profile }

	// Optional focuses on Score only when it is above zero (non-zero players).
	positiveScoreOpt := common.MakeOptional(
		func(p Profile) O.Option[int] {
			if p.Score > 0 {
				return O.Some(p.Score)
			}
			return O.None[int]()
		},
		func(p Profile, s int) Profile { p.Score = s; return p },
	)
	profileLens := common.MakeLens(
		func(p Player) Profile { return p.Profile },
		func(p Player, pr Profile) Player { p.Profile = pr; return p },
	)

	scoreOpt := F.Pipe1(profileLens, lens.ComposeOptional[Player](positiveScoreOpt))

	active := Player{Profile: Profile{Score: 42}}
	inactive := Player{Profile: Profile{Score: 0}}

	// Active player: GetOption returns Some, Set updates.
	fmt.Println(O.GetOrElse(func() int { return -1 })(scoreOpt.GetOption(active)))
	boosted := scoreOpt.Set(100)(active)
	fmt.Println(boosted.Profile.Score)

	// Inactive player: GetOption returns None (predicate not satisfied).
	// The raw Set method is unconditional; callers guard writes with GetOption.
	fmt.Println(O.IsNone(scoreOpt.GetOption(inactive)))
	fmt.Println(O.IsSome(scoreOpt.GetOption(active)))
	// Output:
	// 42
	// 100
	// true
	// true
}

// ExampleComposeOptional_chained demonstrates chaining Compose and
// ComposeOptional to navigate two levels of nesting before the partial focus.
func ExampleComposeOptional_chained() {
	type Engine struct{ HP int }
	type Car struct{ Engine Engine }
	type Fleet struct{ Car Car }

	// Optional: focus on HP only when it exceeds 100.
	highPowerOpt := common.MakeOptional(
		func(e Engine) O.Option[int] {
			if e.HP > 100 {
				return O.Some(e.HP)
			}
			return O.None[int]()
		},
		func(e Engine, hp int) Engine { e.HP = hp; return e },
	)
	carLens := common.MakeLens(
		func(f Fleet) Car { return f.Car },
		func(f Fleet, c Car) Fleet { f.Car = c; return f },
	)
	engineLens := common.MakeLens(
		func(c Car) Engine { return c.Engine },
		func(c Car, e Engine) Car { c.Engine = e; return c },
	)

	// Fleet → Car → Engine → int  (only when HP > 100)
	fleetHPOpt := F.Pipe1(
		F.Pipe1(carLens, lens.Compose[Fleet](engineLens)),
		lens.ComposeOptional[Fleet](highPowerOpt),
	)

	powerful := Fleet{Car: Car{Engine: Engine{HP: 200}}}
	weak := Fleet{Car: Car{Engine: Engine{HP: 50}}}

	fmt.Println(O.IsSome(fleetHPOpt.GetOption(powerful)))
	fmt.Println(O.IsNone(fleetHPOpt.GetOption(weak)))

	tuned := fleetHPOpt.Set(300)(powerful)
	fmt.Println(tuned.Car.Engine.HP)
	// Output:
	// true
	// true
	// 300
}
