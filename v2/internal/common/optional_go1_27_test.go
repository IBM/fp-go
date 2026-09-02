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

// optionals used across the method-compose tests.
// outer: Person → Name  (None when Name is empty)
// inner: string → first rune  (None when string is empty)

func makeNameOptional() Optional[Person, string] {
	return MakeOptional(
		func(p Person) Option[string] {
			if p.Name != "" {
				return OptionSome(p.Name)
			}
			return OptionNone[string]()
		},
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)
}

func makeFirstCharOptional() Optional[string, rune] {
	return MakeOptional(
		func(s string) Option[rune] {
			if len(s) > 0 {
				return OptionSome(rune(s[0]))
			}
			return OptionNone[rune]()
		},
		func(s string, r rune) string {
			if len(s) > 0 {
				return string(r) + s[1:]
			}
			return string(r)
		},
	)
}

// TestOptionalCompose_MethodEquivalentToFreeFunction verifies that the method
// form o.Compose(ab) returns an optional identical in behaviour to the
// free-function Compose[S](ab)(o).
func TestOptionalCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	outer := makeNameOptional()
	inner := makeFirstCharOptional()

	method := outer.Compose(inner)
	free := OptionalComposeOptional[Person](inner)(outer)

	persons := []Person{
		{Name: "Alice", Age: 30},
		{Name: "", Age: 20},
	}
	for _, p := range persons {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(p), method.GetOption(p))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set('Z')(p), method.Set('Z')(p))
		})
	}
}

// TestOptionalCompose_Method_GetOption_Match verifies Some is returned when
// both optionals produce a value.
func TestOptionalCompose_Method_GetOption_Match(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	got := composed.GetOption(Person{Name: "Alice"})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 'A', OptionGetOrElse(F.Constant(rune(0)))(got))
}

// TestOptionalCompose_Method_GetOption_OuterNone verifies None is returned when
// the outer optional produces None.
func TestOptionalCompose_Method_GetOption_OuterNone(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	assert.True(t, OptionIsNone(composed.GetOption(Person{Name: ""})))
}

// TestOptionalCompose_Method_GetOption_InnerNone verifies None is returned when
// the outer optional matches but the inner optional produces None.
func TestOptionalCompose_Method_GetOption_InnerNone(t *testing.T) {
	// inner optional that always returns None (simulates empty focus)
	alwaysNone := MakeOptional(
		func(_ string) Option[rune] { return OptionNone[rune]() },
		func(s string, _ rune) string { return s },
	)
	composed := makeNameOptional().Compose(alwaysNone)

	assert.True(t, OptionIsNone(composed.GetOption(Person{Name: "Alice"})))
}

// TestOptionalCompose_Method_Set_Match verifies that Set updates the focus when
// both optionals match.
func TestOptionalCompose_Method_Set_Match(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	updated := composed.Set('B')(Person{Name: "Alice", Age: 30})
	assert.Equal(t, "Blice", updated.Name)
	assert.Equal(t, 30, updated.Age)
}

// TestOptionalCompose_Method_Set_NoOpOnOuterNone verifies that Set is a no-op
// when the outer optional does not match (optional no-op law).
func TestOptionalCompose_Method_Set_NoOpOnOuterNone(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	original := Person{Name: "", Age: 42}
	result := composed.Set('Z')(original)
	assert.Equal(t, original, result)
}

// TestOptionalCompose_Method_OptionalLaws verifies the three optional laws on
// the method-composed optional.
func TestOptionalCompose_Method_OptionalLaws(t *testing.T) {
	composed := makeNameOptional().Compose(makeFirstCharOptional())

	matching := Person{Name: "Alice", Age: 30}
	nonMatching := Person{Name: "", Age: 30}

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := composed.Set('Z')(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := composed.Set('B')(matching)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 'B', OptionGetOrElse(F.Constant(rune(0)))(got))
	})

	t.Run("law SetSet: Set(b)(Set(a)(s)) == Set(b)(s)", func(t *testing.T) {
		result1 := composed.Set('X')(composed.Set('Y')(matching))
		result2 := composed.Set('X')(matching)
		assert.Equal(t, result2, result1)
	})
}

// TestOptionalCompose_Method_Chained verifies that three .Compose calls chain
// correctly through three levels of optional structure.
func TestOptionalCompose_Method_Chained(t *testing.T) {
	// level 1: Response → *Info  (None when info is nil)
	// level 2: *Info → *Employment  (None when employment is nil)
	// level 3: *Employment → *Phone  (None when phone is nil)
	l1 := OptionalFromPredicateRef[Response](F.IsNonNil[Info])((*Response).GetInfo, (*Response).SetInfo)
	l2 := MakeOptionalRef(
		func(i *Info) Option[*Employment] {
			if i.employment != nil {
				return OptionSome(i.employment)
			}
			return OptionNone[*Employment]()
		},
		func(i *Info, e *Employment) *Info { i.employment = e; return i },
	)
	l3 := MakeOptionalRef(
		func(e *Employment) Option[*Phone] {
			if e.phone != nil {
				return OptionSome(e.phone)
			}
			return OptionNone[*Phone]()
		},
		func(e *Employment, p *Phone) *Employment { e.phone = p; return e },
	)

	composed := l1.Compose(l2).Compose(l3)

	phone := &Phone{number: "555-1234"}
	emp := &Employment{phone: phone}
	info := &Info{employment: emp}
	resp := &Response{info: info}

	t.Run("all three levels present", func(t *testing.T) {
		got := composed.GetOption(resp)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, phone, OptionGetOrElse(F.Constant[*Phone](nil))(got))
	})

	t.Run("middle level nil short-circuits", func(t *testing.T) {
		resp2 := &Response{info: &Info{employment: nil}}
		assert.True(t, OptionIsNone(composed.GetOption(resp2)))
	})

	t.Run("outermost level nil short-circuits", func(t *testing.T) {
		resp3 := &Response{info: nil}
		assert.True(t, OptionIsNone(composed.GetOption(resp3)))
	})

	t.Run("Set is no-op when chain breaks", func(t *testing.T) {
		resp4 := &Response{info: nil}
		result := composed.Set(&Phone{number: "999"})(resp4)
		assert.Nil(t, result.info)
	})
}

// ---- ComposeLens method tests ----
//
// Fixtures:
//   outer: Wrapper → *Config  (None when config is nil)
//   inner: Lens[*Config, int]  (always focuses on Timeout)

type Wrapper struct {
	cfg *Config
}

func makeConfigOptional() Optional[Wrapper, *Config] {
	return MakeOptional(
		func(w Wrapper) Option[*Config] {
			if w.cfg != nil {
				return OptionSome(w.cfg)
			}
			return OptionNone[*Config]()
		},
		func(w Wrapper, c *Config) Wrapper {
			w.cfg = c
			return w
		},
	)
}

func makeTimeoutLens() Lens[*Config, int] {
	return MakeLens(
		func(c *Config) int { return c.Timeout },
		func(c *Config, t int) *Config { c2 := *c; c2.Timeout = t; return &c2 },
	)
}

// TestOptionalComposeLens_MethodEquivalentToFreeFunction verifies that the
// method form o.ComposeLens(ab) returns an optional identical in behaviour to
// the free-function OptionalComposeLens[S](ab)(o).
func TestOptionalComposeLens_MethodEquivalentToFreeFunction(t *testing.T) {
	outer := makeConfigOptional()
	inner := makeTimeoutLens()

	method := outer.ComposeLens(inner)
	free := OptionalComposeLens[Wrapper](inner)(outer)

	wrappers := []Wrapper{
		{cfg: &Config{Timeout: 5, Retries: 3}},
		{cfg: nil},
	}
	for _, w := range wrappers {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(w), method.GetOption(w))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(w), method.Set(99)(w))
		})
	}
}

// TestOptionalComposeLens_Method_GetOption_Match verifies that Some is returned
// when the outer optional matches (cfg is non-nil).
func TestOptionalComposeLens_Method_GetOption_Match(t *testing.T) {
	composed := makeConfigOptional().ComposeLens(makeTimeoutLens())

	w := Wrapper{cfg: &Config{Timeout: 42, Retries: 1}}
	got := composed.GetOption(w)
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 42, OptionGetOrElse(F.Constant(0))(got))
}

// TestOptionalComposeLens_Method_GetOption_OuterNone verifies that None is
// returned when the outer optional produces None (cfg is nil).
func TestOptionalComposeLens_Method_GetOption_OuterNone(t *testing.T) {
	composed := makeConfigOptional().ComposeLens(makeTimeoutLens())

	assert.True(t, OptionIsNone(composed.GetOption(Wrapper{cfg: nil})))
}

// TestOptionalComposeLens_Method_Set_Match verifies that Set updates the lens
// focus when the outer optional matches.
func TestOptionalComposeLens_Method_Set_Match(t *testing.T) {
	composed := makeConfigOptional().ComposeLens(makeTimeoutLens())

	w := Wrapper{cfg: &Config{Timeout: 10, Retries: 2}}
	updated := composed.Set(99)(w)
	assert.Equal(t, 99, updated.cfg.Timeout)
	assert.Equal(t, 2, updated.cfg.Retries) // unrelated field preserved
}

// TestOptionalComposeLens_Method_Set_NoOpOnOuterNone verifies that Set is a
// no-op when the outer optional does not match (Optional no-op law).
func TestOptionalComposeLens_Method_Set_NoOpOnOuterNone(t *testing.T) {
	composed := makeConfigOptional().ComposeLens(makeTimeoutLens())

	original := Wrapper{cfg: nil}
	result := composed.Set(99)(original)
	assert.Equal(t, original, result)
}

// TestOptionalComposeLens_Method_OptionalLaws verifies the three optional laws
// on the method-composed optional.
func TestOptionalComposeLens_Method_OptionalLaws(t *testing.T) {
	composed := makeConfigOptional().ComposeLens(makeTimeoutLens())

	matching := Wrapper{cfg: &Config{Timeout: 7, Retries: 3}}
	nonMatching := Wrapper{cfg: nil}

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := composed.Set(99)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := composed.Set(55)(matching)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 55, OptionGetOrElse(F.Constant(0))(got))
	})

	t.Run("law SetSet: Set(b)(Set(a)(s)) == Set(b)(s)", func(t *testing.T) {
		result1 := composed.Set(20)(composed.Set(10)(matching))
		result2 := composed.Set(20)(matching)
		assert.Equal(t, result2.cfg.Timeout, result1.cfg.Timeout)
		assert.Equal(t, result2.cfg.Retries, result1.cfg.Retries)
	})
}

// TestOptionalComposeLens_Method_Chained verifies that ComposeLens can be
// chained after Compose to navigate through mixed optional/lens structure.
func TestOptionalComposeLens_Method_Chained(t *testing.T) {
	// level 1: Response → *Info  (None when info is nil) — reuse responseOptional
	// level 2 (ComposeLens): *Info → *Employment  (lens, always present after info is non-nil)
	//
	// We add an employment field via a lens to keep the test self-contained.
	empLens := MakeLens(
		func(i *Info) *Employment { return i.employment },
		func(i *Info, e *Employment) *Info { i2 := *i; i2.employment = e; return &i2 },
	)

	composed := responseOptional.ComposeLens(empLens)

	emp := &Employment{phone: &Phone{number: "555-9876"}}
	resp := &Response{info: &Info{employment: emp}}

	t.Run("Get navigates through optional then lens", func(t *testing.T) {
		got := composed.GetOption(resp)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, emp, OptionGetOrElse(F.Constant[*Employment](nil))(got))
	})

	t.Run("Get returns None when outer optional misses", func(t *testing.T) {
		resp2 := &Response{info: nil}
		assert.True(t, OptionIsNone(composed.GetOption(resp2)))
	})

	t.Run("Set updates through optional then lens", func(t *testing.T) {
		newEmp := &Employment{phone: &Phone{number: "999-0000"}}
		updated := composed.Set(newEmp)(resp)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, newEmp, OptionGetOrElse(F.Constant[*Employment](nil))(got))
	})

	t.Run("Set is no-op when outer optional misses", func(t *testing.T) {
		resp2 := &Response{info: nil}
		result := composed.Set(&Employment{})(resp2)
		assert.Nil(t, result.info)
	})
}

// ---- ComposePrism method tests ----
//
// Fixtures:
//   outer: Wrapper → *Config  (None when config is nil)       — reuses makeConfigOptional()
//   inner: Prism[*Config, int]  (matches only when Timeout > 0)

// makePositiveTimeoutPrism returns a Prism[*Config, int] that extracts Timeout
// only when it is strictly positive, and reconstructs a *Config with that
// Timeout via ReverseGet.
func makePositiveTimeoutPrism() Prism[*Config, int] {
	return MakePrism(
		func(c *Config) Option[int] {
			if c.Timeout > 0 {
				return OptionSome(c.Timeout)
			}
			return OptionNone[int]()
		},
		func(t int) *Config { return &Config{Timeout: t} },
	)
}

// TestOptionalComposePrism_MethodEquivalentToFreeFunction verifies that the
// method form o.ComposePrism(ab) returns an optional identical in behaviour to
// the free-function OptionalComposePrism[S](ab)(o).
func TestOptionalComposePrism_MethodEquivalentToFreeFunction(t *testing.T) {
	outer := makeConfigOptional()
	inner := makePositiveTimeoutPrism()

	method := outer.ComposePrism(inner)
	free := OptionalComposePrism[Wrapper](inner)(outer)

	wrappers := []Wrapper{
		{cfg: &Config{Timeout: 5, Retries: 3}},
		{cfg: &Config{Timeout: 0, Retries: 3}}, // prism misses
		{cfg: nil},                              // outer misses
	}
	for _, w := range wrappers {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(w), method.GetOption(w))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(99)(w), method.Set(99)(w))
		})
	}
}

// TestOptionalComposePrism_Method_GetOption_Match verifies that Some is returned
// when the outer optional matches and the prism also matches.
func TestOptionalComposePrism_Method_GetOption_Match(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	w := Wrapper{cfg: &Config{Timeout: 42, Retries: 1}}
	got := composed.GetOption(w)
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 42, OptionGetOrElse(F.Constant(0))(got))
}

// TestOptionalComposePrism_Method_GetOption_OuterNone verifies that None is
// returned when the outer optional produces None (cfg is nil).
func TestOptionalComposePrism_Method_GetOption_OuterNone(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	assert.True(t, OptionIsNone(composed.GetOption(Wrapper{cfg: nil})))
}

// TestOptionalComposePrism_Method_GetOption_PrismMiss verifies that None is
// returned when the outer optional matches but the prism does not (Timeout ≤ 0).
func TestOptionalComposePrism_Method_GetOption_PrismMiss(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	w := Wrapper{cfg: &Config{Timeout: 0, Retries: 2}}
	assert.True(t, OptionIsNone(composed.GetOption(w)))
}

// TestOptionalComposePrism_Method_Set_Match verifies that Set updates through
// the prism when both the outer optional and the prism match.
func TestOptionalComposePrism_Method_Set_Match(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	w := Wrapper{cfg: &Config{Timeout: 10, Retries: 2}}
	updated := composed.Set(99)(w)
	got := composed.GetOption(updated)
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 99, OptionGetOrElse(F.Constant(0))(got))
}

// TestOptionalComposePrism_Method_Set_NoOpOnOuterNone verifies that Set is a
// no-op when the outer optional does not match (cfg is nil).
func TestOptionalComposePrism_Method_Set_NoOpOnOuterNone(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	original := Wrapper{cfg: nil}
	result := composed.Set(99)(original)
	assert.Equal(t, original, result)
}

// TestOptionalComposePrism_Method_Set_NoOpOnPrismMiss verifies that Set is a
// no-op when the outer optional matches but the prism does not.
func TestOptionalComposePrism_Method_Set_NoOpOnPrismMiss(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	original := Wrapper{cfg: &Config{Timeout: 0, Retries: 3}}
	result := composed.Set(99)(original)
	assert.Equal(t, original.cfg.Timeout, result.cfg.Timeout)
	assert.Equal(t, original.cfg.Retries, result.cfg.Retries)
}

// TestOptionalComposePrism_Method_OptionalLaws verifies the three optional laws
// on the method-composed optional.
func TestOptionalComposePrism_Method_OptionalLaws(t *testing.T) {
	composed := makeConfigOptional().ComposePrism(makePositiveTimeoutPrism())

	matching := Wrapper{cfg: &Config{Timeout: 7, Retries: 3}}
	nonMatchingOuter := Wrapper{cfg: nil}
	nonMatchingPrism := Wrapper{cfg: &Config{Timeout: 0, Retries: 3}}

	t.Run("law GetSet: outer None => Set is no-op", func(t *testing.T) {
		result := composed.Set(99)(nonMatchingOuter)
		assert.Equal(t, nonMatchingOuter, result)
	})

	t.Run("law GetSet: prism miss => Set is no-op", func(t *testing.T) {
		result := composed.Set(99)(nonMatchingPrism)
		assert.Equal(t, nonMatchingPrism.cfg.Timeout, result.cfg.Timeout)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := composed.Set(55)(matching)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 55, OptionGetOrElse(F.Constant(0))(got))
	})

	t.Run("law SetSet: Set(b)(Set(a)(s)) == Set(b)(s)", func(t *testing.T) {
		result1 := composed.Set(20)(composed.Set(10)(matching))
		result2 := composed.Set(20)(matching)
		got1 := composed.GetOption(result1)
		got2 := composed.GetOption(result2)
		assert.Equal(t, got2, got1)
	})
}

// TestOptionalComposePrism_Method_Chained verifies that ComposePrism can be
// chained after ComposeLens to navigate a three-level structure:
//
//	Frame   → *Canvas  (optional, None when nil)
//	*Canvas → Shape    (lens, always present)
//	Shape   → float64  (prism, None for non-Circle shapes)
func TestOptionalComposePrism_Method_Chained(t *testing.T) {
	// Use a self-contained structure to avoid coupling to Info internals.
	// Three levels:
	//   Optional[Frame, *Canvas]  — None when canvas is nil
	//   Lens[*Canvas, Shape]      — always succeeds
	//   Prism[Shape, float64]     — matches Circle only
	type Canvas struct {
		shape Shape
	}
	type Frame struct {
		canvas *Canvas
	}

	frameOptional := MakeOptional(
		func(f Frame) Option[*Canvas] {
			if f.canvas != nil {
				return OptionSome(f.canvas)
			}
			return OptionNone[*Canvas]()
		},
		func(f Frame, c *Canvas) Frame { f.canvas = c; return f },
	)

	canvasShapeLens := MakeLens(
		func(c *Canvas) Shape { return c.shape },
		func(c *Canvas, s Shape) *Canvas { c2 := *c; c2.shape = s; return &c2 },
	)

	circlePrism := MakePrism(
		func(s Shape) Option[float64] {
			if c, ok := s.(Circle); ok {
				return OptionSome(c.Radius)
			}
			return OptionNone[float64]()
		},
		func(r float64) Shape { return Circle{Radius: r} },
	)

	// chain: Optional[Frame, *Canvas] ∘ Lens[*Canvas, Shape] ∘ Prism[Shape, float64]
	composed := frameOptional.ComposeLens(canvasShapeLens).ComposePrism(circlePrism)

	circleFrame := Frame{canvas: &Canvas{shape: Circle{Radius: 5.0}}}
	rectFrame := Frame{canvas: &Canvas{shape: Rectangle{Width: 4, Height: 3}}}
	nilFrame := Frame{canvas: nil}

	t.Run("all levels match — returns Some(radius)", func(t *testing.T) {
		got := composed.GetOption(circleFrame)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 5.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("prism misses (Rectangle) — returns None", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(rectFrame)))
	})

	t.Run("outer optional misses (nil canvas) — returns None", func(t *testing.T) {
		assert.True(t, OptionIsNone(composed.GetOption(nilFrame)))
	})

	t.Run("Set updates radius when all levels match", func(t *testing.T) {
		updated := composed.Set(9.0)(circleFrame)
		got := composed.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 9.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("Set is no-op when prism misses", func(t *testing.T) {
		result := composed.Set(9.0)(rectFrame)
		_, isRect := result.canvas.shape.(Rectangle)
		assert.True(t, isRect, "shape must still be Rectangle")
	})

	t.Run("Set is no-op when outer optional misses", func(t *testing.T) {
		result := composed.Set(9.0)(nilFrame)
		assert.Nil(t, result.canvas)
	})
}
