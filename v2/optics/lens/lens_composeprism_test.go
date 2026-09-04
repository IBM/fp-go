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

// ---- test types ----

// tcShape is a sum type (tagged union) representing geometric shapes.
type tcShape interface{ isShape() }

// tcCircle is the Circle variant of tcShape.
type tcCircle struct{ Radius float64 }

// tcRect is the Rectangle variant of tcShape.
type tcRect struct{ Width, Height float64 }

func (tcCircle) isShape() {}
func (tcRect) isShape()   {}

// tcCanvas holds the shape that is drawn on it.
type tcCanvas struct{ Shape tcShape }

// ---- helpers ----

// canvasShapeLens focuses on the Shape field of tcCanvas.
var canvasShapeLens = common.MakeLens(
	func(c tcCanvas) tcShape { return c.Shape },
	func(c tcCanvas, s tcShape) tcCanvas { c.Shape = s; return c },
)

// circlePrism focuses on the tcCircle variant of tcShape.
var circlePrism = common.MakePrism(
	func(s tcShape) O.Option[tcCircle] {
		if c, ok := s.(tcCircle); ok {
			return O.Some(c)
		}
		return O.None[tcCircle]()
	},
	func(c tcCircle) tcShape { return c },
)

// rectPrism focuses on the tcRect variant of tcShape.
var rectPrism = common.MakePrism(
	func(s tcShape) O.Option[tcRect] {
		if r, ok := s.(tcRect); ok {
			return O.Some(r)
		}
		return O.None[tcRect]()
	},
	func(r tcRect) tcShape { return r },
)

// canvasCircleOpt is the Optional[tcCanvas, tcCircle] produced by ComposePrism.
var canvasCircleOpt = F.Pipe1(canvasShapeLens, lens.ComposePrism[tcCanvas](circlePrism))

// ---- tests ----

// TestComposePrism_GetOption_Match verifies Some is returned when the Prism matches.
func TestComposePrism_GetOption_Match(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 5}}
	got := canvasCircleOpt.GetOption(canvas)
	assert.Equal(t, O.Some(tcCircle{Radius: 5}), got)
}

// TestComposePrism_GetOption_NoMatch verifies None is returned when the Prism does not match.
func TestComposePrism_GetOption_NoMatch(t *testing.T) {
	canvas := tcCanvas{Shape: tcRect{Width: 3, Height: 4}}
	got := canvasCircleOpt.GetOption(canvas)
	assert.Equal(t, O.None[tcCircle](), got)
}

// TestComposePrism_Set_Match verifies that Set updates the value when the Prism matches.
func TestComposePrism_Set_Match(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 5}}
	updated := canvasCircleOpt.Set(tcCircle{Radius: 10})(canvas)
	assert.Equal(t, tcCanvas{Shape: tcCircle{Radius: 10}}, updated)
}

// TestComposePrism_Set_NoOpOnMiss verifies that Set is a no-op when the Prism does not match.
func TestComposePrism_Set_NoOpOnMiss(t *testing.T) {
	canvas := tcCanvas{Shape: tcRect{Width: 3, Height: 4}}
	unchanged := canvasCircleOpt.Set(tcCircle{Radius: 10})(canvas)
	assert.Equal(t, canvas, unchanged)
}

// TestComposePrism_Immutability verifies that Set does not mutate the original structure.
func TestComposePrism_Immutability(t *testing.T) {
	original := tcCanvas{Shape: tcCircle{Radius: 5}}
	_ = canvasCircleOpt.Set(tcCircle{Radius: 99})(original)
	assert.Equal(t, tcCircle{Radius: 5}, original.Shape)
}

// TestComposePrism_OptionalLaws verifies the two Optional laws.
func TestComposePrism_OptionalLaws(t *testing.T) {
	circle := tcCanvas{Shape: tcCircle{Radius: 7}}
	rect := tcCanvas{Shape: tcRect{Width: 2, Height: 3}}

	t.Run("law SetGet (match): GetOption(Set(b)(s)) == Some(b)", func(t *testing.T) {
		newCircle := tcCircle{Radius: 42}
		updated := canvasCircleOpt.Set(newCircle)(circle)
		result := canvasCircleOpt.GetOption(updated)
		assert.Equal(t, O.Some(newCircle), result)
	})

	t.Run("law GetSet (no-match): Set(b)(s) == s when GetOption(s) == None", func(t *testing.T) {
		assert.Equal(t, O.None[tcCircle](), canvasCircleOpt.GetOption(rect))
		unchanged := canvasCircleOpt.Set(tcCircle{Radius: 99})(rect)
		assert.Equal(t, rect, unchanged)
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		c1, c2 := tcCircle{Radius: 1}, tcCircle{Radius: 2}
		setTwice := canvasCircleOpt.Set(c2)(canvasCircleOpt.Set(c1)(circle))
		setOnce := canvasCircleOpt.Set(c2)(circle)
		assert.Equal(t, setOnce, setTwice)
	})
}

// TestComposePrism_EquivalentToUnderlyingFreeFunction verifies that the lens-package
// wrapper produces the same result as calling common.LensComposePrism directly.
func TestComposePrism_EquivalentToUnderlyingFreeFunction(t *testing.T) {
	canvas := tcCanvas{Shape: tcCircle{Radius: 3}}
	newCircle := tcCircle{Radius: 9}

	wrapped := F.Pipe1(canvasShapeLens, lens.ComposePrism[tcCanvas](circlePrism))
	direct := common.LensComposePrism[tcCanvas](circlePrism)(canvasShapeLens)

	t.Run("GetOption returns same result", func(t *testing.T) {
		assert.Equal(t, direct.GetOption(canvas), wrapped.GetOption(canvas))
	})
	t.Run("Set returns same result", func(t *testing.T) {
		assert.Equal(t, direct.Set(newCircle)(canvas), wrapped.Set(newCircle)(canvas))
	})
}

// TestComposePrism_MultipleVariants verifies independent Optionals over the same Lens.
func TestComposePrism_MultipleVariants(t *testing.T) {
	circleOpt := F.Pipe1(canvasShapeLens, lens.ComposePrism[tcCanvas](circlePrism))
	rectOpt := F.Pipe1(canvasShapeLens, lens.ComposePrism[tcCanvas](rectPrism))

	circleCanvas := tcCanvas{Shape: tcCircle{Radius: 5}}
	assert.True(t, O.IsSome(circleOpt.GetOption(circleCanvas)))
	assert.True(t, O.IsNone(rectOpt.GetOption(circleCanvas)))

	rectCanvas := tcCanvas{Shape: tcRect{Width: 2, Height: 3}}
	assert.True(t, O.IsNone(circleOpt.GetOption(rectCanvas)))
	assert.True(t, O.IsSome(rectOpt.GetOption(rectCanvas)))
}

// ---- examples ----

// ExampleComposePrism demonstrates focusing on a specific variant of a sum type
// inside a struct via a Lens + Prism composition.
//
// The test types tcCanvas, tcCircle, and tcRect are defined at package level
// in this file.
func ExampleComposePrism() {
	shapeLens := common.MakeLens(
		func(c tcCanvas) tcShape { return c.Shape },
		func(c tcCanvas, s tcShape) tcCanvas { c.Shape = s; return c },
	)
	circPrism := common.MakePrism(
		func(s tcShape) O.Option[tcCircle] {
			if c, ok := s.(tcCircle); ok {
				return O.Some(c)
			}
			return O.None[tcCircle]()
		},
		func(c tcCircle) tcShape { return c },
	)

	circleOpt := F.Pipe1(shapeLens, lens.ComposePrism[tcCanvas](circPrism))

	// GetOption returns Some when the shape is a Circle.
	c := tcCanvas{Shape: tcCircle{Radius: 5}}
	fmt.Println(O.IsSome(circleOpt.GetOption(c)))

	// GetOption returns None when the shape is not a Circle.
	r := tcCanvas{Shape: tcRect{Width: 3, Height: 4}}
	fmt.Println(O.IsNone(circleOpt.GetOption(r)))

	// Set updates the Circle; is a no-op for Rect.
	updated := circleOpt.Set(tcCircle{Radius: 10})(c)
	fmt.Println(updated.Shape.(tcCircle).Radius)
	unchanged := circleOpt.Set(tcCircle{Radius: 10})(r)
	fmt.Println(unchanged.Shape.(tcRect).Width)
	// Output:
	// true
	// true
	// 10
	// 3
}
