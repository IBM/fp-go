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

// TestLensCompose_MethodEquivalentToFreeFunction verifies that the method form
// l.Compose(ab) returns a lens identical in behaviour to the free function
// Compose[S](ab)(l).
func TestLensCompose_MethodEquivalentToFreeFunction(t *testing.T) {
	// Use the package-level test fixtures (Street/Address/streetLens/addrLens).
	sampleStreet := Street{num: 1, name: "Main St"}
	sampleAddress := Address{city: "Springfield", street: &sampleStreet}

	method := addrLens.Compose(streetLens)
	free := LensComposeLens[*Address](streetLens)(addrLens)

	t.Run("Get returns same value", func(t *testing.T) {
		assert.Equal(t, free.Get(&sampleAddress), method.Get(&sampleAddress))
	})

	t.Run("Set returns same result", func(t *testing.T) {
		newName := "Oak Ave"
		assert.Equal(t,
			free.Set(newName)(&sampleAddress).street.name,
			method.Set(newName)(&sampleAddress).street.name,
		)
	})
}

// TestLensCompose_Method_Get verifies that the composed lens retrieves the
// correct deeply-nested value.
func TestLensCompose_Method_Get(t *testing.T) {
	sampleStreet := Street{num: 42, name: "Elm St"}
	sampleAddress := Address{city: "Shelbyville", street: &sampleStreet}

	composed := addrLens.Compose(streetLens)

	assert.Equal(t, "Elm St", composed.Get(&sampleAddress))
}

// TestLensCompose_Method_Immutability verifies that Set does not mutate the
// original pointer or any nested pointer.
func TestLensCompose_Method_Immutability(t *testing.T) {
	sampleStreet := Street{num: 10, name: "Original"}
	sampleAddress := Address{city: "Shelbyville", street: &sampleStreet}
	originalStreetPtr := sampleAddress.street

	composed := addrLens.Compose(streetLens)
	updated := composed.Set("Updated")(&sampleAddress)

	t.Run("original name is unchanged", func(t *testing.T) {
		assert.Equal(t, "Original", sampleAddress.street.name)
	})

	t.Run("updated name is set", func(t *testing.T) {
		assert.Equal(t, "Updated", updated.street.name)
	})

	t.Run("nested pointer is a new allocation", func(t *testing.T) {
		assert.NotSame(t, originalStreetPtr, updated.street)
	})
}

// TestLensCompose_Method_LensLaws verifies that the composed method lens
// satisfies all three lens laws.
func TestLensCompose_Method_LensLaws(t *testing.T) {
	sampleStreet := Street{num: 5, name: "Schönaicherstr"}
	sampleAddress := Address{city: "Böblingen", street: &sampleStreet}
	newName := "Böblingerstr"

	l := addrLens.Compose(streetLens)

	t.Run("law GetSet: Set(Get(s))(s) == s", func(t *testing.T) {
		result := l.Set(l.Get(&sampleAddress))(&sampleAddress)
		assert.Equal(t, sampleAddress.street.name, result.street.name)
		assert.Equal(t, sampleAddress.street.num, result.street.num)
	})

	t.Run("law SetGet: Get(Set(b)(s)) == b", func(t *testing.T) {
		assert.Equal(t, newName, l.Get(l.Set(newName)(&sampleAddress)))
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		b1, b2 := "First St", "Second St"
		result1 := l.Set(b2)(l.Set(b1)(&sampleAddress))
		result2 := l.Set(b2)(&sampleAddress)
		assert.Equal(t, result2.street.name, result1.street.name)
	})
}

// ---- ComposeIso method tests ----

// streetNameLens focuses on the name field of a Street value.
// addrLens (*Address → *Street) + streetNameLens (*Street → string) are composed
// in the tests below using .Compose first, giving a Lens[*Address, string] which is
// then used with .ComposeIso.
var streetNameLens = MakeLensRef((*Street).GetName, (*Street).SetName)

// addrStreetNameLens is a pre-composed Lens[*Address, string] used as the base
// lens for all ComposeIso tests in this file.
var addrStreetNameLens = addrLens.Compose(streetNameLens)

// TestLensComposeIso_MethodEquivalentToFreeFunction verifies that the method form
// l.ComposeIso(ab) returns a lens identical in behaviour to the free function
// LensComposeIso[S](ab)(l).
func TestLensComposeIso_MethodEquivalentToFreeFunction(t *testing.T) {
	// Isomorphism: street name (string) <-> upper-case string
	toUpper := func(s string) string {
		result := make([]byte, len(s))
		for i := range len(s) {
			c := s[i]
			if c >= 'a' && c <= 'z' {
				c -= 32
			}
			result[i] = c
		}
		return string(result)
	}
	toLower := func(s string) string {
		result := make([]byte, len(s))
		for i := range len(s) {
			c := s[i]
			if c >= 'A' && c <= 'Z' {
				c += 32
			}
			result[i] = c
		}
		return string(result)
	}
	upperIso := MakeIso(toUpper, toLower)

	sampleStreetLocal := Street{num: 1, name: "Main St"}
	sampleAddressLocal := Address{city: "Springfield", street: &sampleStreetLocal}

	method := addrStreetNameLens.ComposeIso(upperIso)
	free := LensComposeIso[*Address](upperIso)(addrStreetNameLens)

	t.Run("Get returns same value", func(t *testing.T) {
		assert.Equal(t, free.Get(&sampleAddressLocal), method.Get(&sampleAddressLocal))
	})

	t.Run("Set returns same result", func(t *testing.T) {
		assert.Equal(t,
			free.Set("OAK AVE")(&sampleAddressLocal).street.name,
			method.Set("OAK AVE")(&sampleAddressLocal).street.name,
		)
	})
}

// TestLensComposeIso_Method_Get verifies that ComposeIso applies the forward iso on Get.
func TestLensComposeIso_Method_Get(t *testing.T) {
	// Isomorphism that wraps the street name in brackets.
	bracketIso := MakeIso(
		func(s string) string { return "[" + s + "]" },
		func(s string) string {
			if len(s) >= 2 {
				return s[1 : len(s)-1]
			}
			return s
		},
	)

	sampleStreetLocal := Street{num: 7, name: "Elm St"}
	sampleAddressLocal := Address{city: "Shelbyville", street: &sampleStreetLocal}

	composed := addrStreetNameLens.ComposeIso(bracketIso)
	assert.Equal(t, "[Elm St]", composed.Get(&sampleAddressLocal))
}

// TestLensComposeIso_Method_Set verifies that Set applies the reverse iso direction.
func TestLensComposeIso_Method_Set(t *testing.T) {
	bracketIso := MakeIso(
		func(s string) string { return "[" + s + "]" },
		func(s string) string {
			if len(s) >= 2 {
				return s[1 : len(s)-1]
			}
			return s
		},
	)

	sampleStreetLocal := Street{num: 7, name: "Elm St"}
	sampleAddressLocal := Address{city: "Shelbyville", street: &sampleStreetLocal}

	composed := addrStreetNameLens.ComposeIso(bracketIso)
	// Setting "[Oak Ave]" should store "Oak Ave" via the reverse iso.
	updated := composed.Set("[Oak Ave]")(&sampleAddressLocal)
	assert.Equal(t, "Oak Ave", updated.street.name)
}

// TestLensComposeIso_Method_Immutability verifies that Set does not mutate the original.
func TestLensComposeIso_Method_Immutability(t *testing.T) {
	identityIso := MakeIso(
		func(s string) string { return s },
		func(s string) string { return s },
	)

	sampleStreetLocal := Street{num: 10, name: "Original"}
	sampleAddressLocal := Address{city: "Shelbyville", street: &sampleStreetLocal}
	originalStreetPtr := sampleAddressLocal.street

	composed := addrStreetNameLens.ComposeIso(identityIso)
	updated := composed.Set("Updated")(&sampleAddressLocal)

	t.Run("original name is unchanged", func(t *testing.T) {
		assert.Equal(t, "Original", sampleAddressLocal.street.name)
	})

	t.Run("updated name is set", func(t *testing.T) {
		assert.Equal(t, "Updated", updated.street.name)
	})

	t.Run("nested pointer is a new allocation", func(t *testing.T) {
		assert.NotSame(t, originalStreetPtr, updated.street)
	})
}

// TestLensComposeIso_Method_LensLaws verifies all three lens laws for a ComposeIso result.
func TestLensComposeIso_Method_LensLaws(t *testing.T) {
	identityIso := MakeIso(
		func(s string) string { return s },
		func(s string) string { return s },
	)

	sampleStreetLocal := Street{num: 5, name: "Schönaicherstr"}
	sampleAddressLocal := Address{city: "Böblingen", street: &sampleStreetLocal}
	newName := "Böblingerstr"

	l := addrStreetNameLens.ComposeIso(identityIso)

	t.Run("law GetSet: Set(Get(s))(s) == s", func(t *testing.T) {
		result := l.Set(l.Get(&sampleAddressLocal))(&sampleAddressLocal)
		assert.Equal(t, sampleAddressLocal.street.name, result.street.name)
		assert.Equal(t, sampleAddressLocal.street.num, result.street.num)
	})

	t.Run("law SetGet: Get(Set(b)(s)) == b", func(t *testing.T) {
		assert.Equal(t, newName, l.Get(l.Set(newName)(&sampleAddressLocal)))
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		b1, b2 := "First St", "Second St"
		result1 := l.Set(b2)(l.Set(b1)(&sampleAddressLocal))
		result2 := l.Set(b2)(&sampleAddressLocal)
		assert.Equal(t, result2.street.name, result1.street.name)
	})
}

// TestLensComposeIso_Method_Chained verifies that ComposeIso can be chained with Compose.
func TestLensComposeIso_Method_Chained(t *testing.T) {
	type StreetNum int
	type NumWrapper struct{ val StreetNum }

	numIso := MakeIso(
		func(n StreetNum) NumWrapper { return NumWrapper{val: n} },
		func(w NumWrapper) StreetNum { return w.val },
	)

	// Build a *Street → StreetNum lens and chain: *Address → *Street → StreetNum → NumWrapper.
	streetNumLens := MakeLens(
		func(s *Street) StreetNum { return StreetNum(s.num) },
		func(s *Street, n StreetNum) *Street { cp := *s; cp.num = int(n); return &cp },
	)
	wrapperLens := addrLens.Compose(streetNumLens).ComposeIso(numIso)

	sampleStreetLocal := Street{num: 42, name: "Broadway"}
	sampleAddressLocal := Address{city: "NYC", street: &sampleStreetLocal}

	t.Run("Get navigates through Compose+ComposeIso", func(t *testing.T) {
		w := wrapperLens.Get(&sampleAddressLocal)
		assert.Equal(t, NumWrapper{val: 42}, w)
	})

	t.Run("Set updates through Compose+ComposeIso", func(t *testing.T) {
		updated := wrapperLens.Set(NumWrapper{val: 99})(&sampleAddressLocal)
		assert.Equal(t, 99, updated.street.num)
		assert.Equal(t, 42, sampleAddressLocal.street.num) // original unchanged
	})
}

// TestLensComposeIso_Method_PipelineUsage verifies usage in an F.Pipe1 pipeline, which
// mirrors how the underlying free function is used at the package level.
func TestLensComposeIso_Method_PipelineUsage(t *testing.T) {
	upperIso := MakeIso(
		func(s string) string {
			b := []byte(s)
			for i, c := range b {
				if c >= 'a' && c <= 'z' {
					b[i] = c - 32
				}
			}
			return string(b)
		},
		func(s string) string {
			b := []byte(s)
			for i, c := range b {
				if c >= 'A' && c <= 'Z' {
					b[i] = c + 32
				}
			}
			return string(b)
		},
	)

	sampleStreetLocal := Street{num: 1, name: "main street"}
	sampleAddressLocal := Address{city: "Springfield", street: &sampleStreetLocal}

	// Use F.Pipe1 with LensComposeIso — same semantics as method form
	pipelineLens := F.Pipe1(addrStreetNameLens, LensComposeIso[*Address](upperIso))
	methodLens := addrStreetNameLens.ComposeIso(upperIso)

	assert.Equal(t, pipelineLens.Get(&sampleAddressLocal), methodLens.Get(&sampleAddressLocal))
	assert.Equal(t, "MAIN STREET", methodLens.Get(&sampleAddressLocal))
}

// TestLensCompose_Method_Chained verifies that multiple .Compose calls can be
// chained to navigate three or more levels deep.
func TestLensCompose_Method_Chained(t *testing.T) {
	type City struct{ name string }
	type District struct{ city City }
	type Region struct{ district District }

	cityNameLens := MakeLens(
		func(c City) string { return c.name },
		func(c City, n string) City { c.name = n; return c },
	)
	districtCityLens := MakeLens(
		func(d District) City { return d.city },
		func(d District, c City) District { d.city = c; return d },
	)
	regionDistrictLens := MakeLens(
		func(r Region) District { return r.district },
		func(r Region, d District) Region { r.district = d; return r },
	)

	// Chained method composition: Region → District → City → string
	regionCityNameLens := regionDistrictLens.Compose(districtCityLens).Compose(cityNameLens)

	region := Region{District{City{"Böblingen"}}}

	t.Run("Get navigates three levels", func(t *testing.T) {
		assert.Equal(t, "Böblingen", regionCityNameLens.Get(region))
	})

	t.Run("Set updates three levels deep", func(t *testing.T) {
		updated := regionCityNameLens.Set("Stuttgart")(region)
		assert.Equal(t, "Stuttgart", regionCityNameLens.Get(updated))
		assert.Equal(t, "Böblingen", regionCityNameLens.Get(region)) // original unchanged
	})
}

// ---- ComposePrism method tests ----

// shapeLens is a Lens[*Address, Shape] used for the ComposePrism tests.
// It is defined here simply for test purposes — a Street's name is re-interpreted
// as a Shape stored in a thin wrapper via a local test struct.

// addrShapeHolder is a small struct that pairs a *Address with a Shape field,
// used by the ComposePrism tests below.
type addrShapeHolder struct {
	addr  *Address
	shape Shape
}

// shapeHolderLens focuses on the shape field of an addrShapeHolder.
var shapeHolderLens = MakeLens(
	func(h addrShapeHolder) Shape { return h.shape },
	func(h addrShapeHolder, s Shape) addrShapeHolder { h.shape = s; return h },
)

// circlePrismForLensTest focuses on the Circle variant of Shape.
var circlePrismForLensTest = MakePrism(
	func(s Shape) Option[float64] {
		if c, ok := s.(Circle); ok {
			return OptionSome(c.Radius)
		}
		return OptionNone[float64]()
	},
	func(r float64) Shape { return Circle{Radius: r} },
)

// TestLensComposePrism_MethodEquivalentToFreeFunction verifies that the method
// form l.ComposePrism(ab) returns an Optional identical in behaviour to the
// free-function LensComposePrism[S](ab)(l).
func TestLensComposePrism_MethodEquivalentToFreeFunction(t *testing.T) {
	method := shapeHolderLens.ComposePrism(circlePrismForLensTest)
	free := LensComposePrism[addrShapeHolder](circlePrismForLensTest)(shapeHolderLens)

	cases := []addrShapeHolder{
		{shape: Circle{Radius: 3.0}},
		{shape: Rectangle{Width: 4, Height: 5}},
	}
	for _, c := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(c), method.GetOption(c))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(7.0)(c), method.Set(7.0)(c))
		})
	}
}

// TestLensComposePrism_Method_GetOption_Match verifies Some is returned when
// the lens focuses correctly and the prism matches the variant.
func TestLensComposePrism_Method_GetOption_Match(t *testing.T) {
	opt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	got := opt.GetOption(addrShapeHolder{shape: Circle{Radius: 5.0}})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 5.0, OptionGetOrElse(F.Constant(0.0))(got))
}

// TestLensComposePrism_Method_GetOption_PrismMiss verifies None is returned
// when the prism does not match the current variant.
func TestLensComposePrism_Method_GetOption_PrismMiss(t *testing.T) {
	opt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	got := opt.GetOption(addrShapeHolder{shape: Rectangle{Width: 4, Height: 5}})
	assert.True(t, OptionIsNone(got))
}

// TestLensComposePrism_Method_Set_Match verifies that Set updates the shape
// when the prism matches the current variant.
func TestLensComposePrism_Method_Set_Match(t *testing.T) {
	opt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	h := addrShapeHolder{shape: Circle{Radius: 1.0}}
	updated := opt.Set(9.0)(h)
	assert.Equal(t, Circle{Radius: 9.0}, updated.shape)
}

// TestLensComposePrism_Method_Set_NoOpOnPrismMiss verifies that Set is a no-op
// when the prism does not match (Optional no-op law).
func TestLensComposePrism_Method_Set_NoOpOnPrismMiss(t *testing.T) {
	opt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	h := addrShapeHolder{shape: Rectangle{Width: 4, Height: 5}}
	result := opt.Set(9.0)(h)
	assert.Equal(t, h, result)
}

// TestLensComposePrism_Method_OptionalLaws verifies the three standard Optional
// laws on the method-composed Optional.
func TestLensComposePrism_Method_OptionalLaws(t *testing.T) {
	opt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	matching := addrShapeHolder{shape: Circle{Radius: 3.0}}
	nonMatching := addrShapeHolder{shape: Rectangle{Width: 1, Height: 2}}

	t.Run("law GetSet: GetOption==None => Set is no-op", func(t *testing.T) {
		result := opt.Set(99.0)(nonMatching)
		assert.Equal(t, nonMatching, result)
	})

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := opt.Set(42.0)(matching)
		got := opt.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 42.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		result1 := opt.Set(20.0)(opt.Set(10.0)(matching))
		result2 := opt.Set(20.0)(matching)
		assert.Equal(t, result2, result1)
	})
}

// TestLensComposePrism_Method_Chained verifies that ComposePrism can be
// chained after Compose to navigate through a nested structure and then
// discriminate on a sum type.
func TestLensComposePrism_Method_Chained(t *testing.T) {
	type Outer struct{ holder addrShapeHolder }

	outerLens := MakeLens(
		func(o Outer) addrShapeHolder { return o.holder },
		func(o Outer, h addrShapeHolder) Outer { o.holder = h; return o },
	)

	// Compose lens twice, then apply ComposePrism.
	opt := outerLens.Compose(shapeHolderLens).ComposePrism(circlePrismForLensTest)

	o := Outer{holder: addrShapeHolder{shape: Circle{Radius: 7.0}}}

	t.Run("Get navigates through Compose+ComposePrism", func(t *testing.T) {
		got := opt.GetOption(o)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 7.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("Set updates through Compose+ComposePrism", func(t *testing.T) {
		updated := opt.Set(50.0)(o)
		assert.Equal(t, Circle{Radius: 50.0}, updated.holder.shape)
		assert.Equal(t, Circle{Radius: 7.0}, o.holder.shape) // original unchanged
	})

	t.Run("Set is no-op when prism misses after chained Compose", func(t *testing.T) {
		o2 := Outer{holder: addrShapeHolder{shape: Rectangle{Width: 3, Height: 4}}}
		result := opt.Set(50.0)(o2)
		assert.Equal(t, o2, result)
	})
}

// TestLensComposePrism_Method_PipelineUsage verifies that the method form and
// the free-function form used inside F.Pipe1 produce identical results.
func TestLensComposePrism_Method_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(shapeHolderLens, LensComposePrism[addrShapeHolder](circlePrismForLensTest))
	methodOpt := shapeHolderLens.ComposePrism(circlePrismForLensTest)

	h := addrShapeHolder{shape: Circle{Radius: 2.5}}
	assert.Equal(t, pipelineOpt.GetOption(h), methodOpt.GetOption(h))
	assert.Equal(t, pipelineOpt.Set(8.0)(h), methodOpt.Set(8.0)(h))
}

// ---- ComposeOptional method tests ----

// circleRadiusOptional focuses on the Radius of the Circle variant in a Shape,
// but only when the radius is positive. This makes it a genuinely partial optional
// (not a total function), demonstrating both the Some and None code paths.
var circleRadiusOptional = MakeOptional(
	func(s Shape) Option[float64] {
		if c, ok := s.(Circle); ok && c.Radius > 0 {
			return OptionSome(c.Radius)
		}
		return OptionNone[float64]()
	},
	func(s Shape, r float64) Shape { return Circle{Radius: r} },
)

// TestLensComposeOptional_MethodEquivalentToFreeFunction verifies that the method
// form l.ComposeOptional(ab) returns an Optional identical in behaviour to the
// free-function LensComposeOptional[S](ab)(l).
func TestLensComposeOptional_MethodEquivalentToFreeFunction(t *testing.T) {
	method := shapeHolderLens.ComposeOptional(circleRadiusOptional)
	free := LensComposeOptional[addrShapeHolder](circleRadiusOptional)(shapeHolderLens)

	cases := []addrShapeHolder{
		{shape: Circle{Radius: 3.0}},
		{shape: Rectangle{Width: 4, Height: 5}},
	}
	for _, c := range cases {
		t.Run("GetOption parity", func(t *testing.T) {
			assert.Equal(t, free.GetOption(c), method.GetOption(c))
		})
		t.Run("Set parity", func(t *testing.T) {
			assert.Equal(t, free.Set(7.0)(c), method.Set(7.0)(c))
		})
	}
}

// TestLensComposeOptional_Method_GetOption_Match verifies Some is returned when
// the lens focuses correctly and the optional matches the value.
func TestLensComposeOptional_Method_GetOption_Match(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	got := opt.GetOption(addrShapeHolder{shape: Circle{Radius: 5.0}})
	assert.True(t, OptionIsSome(got))
	assert.Equal(t, 5.0, OptionGetOrElse(F.Constant(0.0))(got))
}

// TestLensComposeOptional_Method_GetOption_OptionalMiss verifies None is returned
// when the optional predicate does not match the focused value.
func TestLensComposeOptional_Method_GetOption_OptionalMiss(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	// Rectangle does not match the Circle-with-positive-radius predicate.
	got := opt.GetOption(addrShapeHolder{shape: Rectangle{Width: 4, Height: 5}})
	assert.True(t, OptionIsNone(got))
}

// TestLensComposeOptional_Method_GetOption_ZeroRadius verifies None is returned
// when the shape is a Circle but its radius is not positive (predicate miss).
func TestLensComposeOptional_Method_GetOption_ZeroRadius(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	got := opt.GetOption(addrShapeHolder{shape: Circle{Radius: 0}})
	assert.True(t, OptionIsNone(got))
}

// TestLensComposeOptional_Method_Set_Match verifies that Set updates the shape
// when the optional matches the focused value.
func TestLensComposeOptional_Method_Set_Match(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	h := addrShapeHolder{shape: Circle{Radius: 1.0}}
	updated := opt.Set(9.0)(h)
	assert.Equal(t, Circle{Radius: 9.0}, updated.shape)
}

// TestLensComposeOptional_Method_OptionalLaws verifies the three standard Optional
// laws on the method-composed Optional.
func TestLensComposeOptional_Method_OptionalLaws(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	matching := addrShapeHolder{shape: Circle{Radius: 3.0}}
	nonMatching := addrShapeHolder{shape: Rectangle{Width: 1, Height: 2}}

	t.Run("law SetGet: GetOption==Some => GetOption(Set(b)(s))==Some(b)", func(t *testing.T) {
		updated := opt.Set(42.0)(matching)
		got := opt.GetOption(updated)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 42.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		result1 := opt.Set(20.0)(opt.Set(10.0)(matching))
		result2 := opt.Set(20.0)(matching)
		assert.Equal(t, result2, result1)
	})

	t.Run("law ModifyOption is no-op when GetOption returns None", func(t *testing.T) {
		// ModifyOption returns None[S] when GetOption returns None, preserving
		// the Optional no-op contract at the ModifyOption level.
		result := OptionalModifyOption[addrShapeHolder](func(r float64) float64 { return r + 1 })(opt)(nonMatching)
		assert.Equal(t, OptionNone[addrShapeHolder](), result)
	})
}

// TestLensComposeOptional_Method_Immutability verifies that Set does not mutate
// the original structure.
func TestLensComposeOptional_Method_Immutability(t *testing.T) {
	opt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	original := addrShapeHolder{shape: Circle{Radius: 3.0}}
	_ = opt.Set(99.0)(original)
	assert.Equal(t, Circle{Radius: 3.0}, original.shape)
}

// TestLensComposeOptional_Method_Chained verifies that ComposeOptional can be
// chained after Compose to navigate through a nested structure and then apply
// a partial optional focus.
func TestLensComposeOptional_Method_Chained(t *testing.T) {
	type Outer struct{ holder addrShapeHolder }

	outerLens := MakeLens(
		func(o Outer) addrShapeHolder { return o.holder },
		func(o Outer, h addrShapeHolder) Outer { o.holder = h; return o },
	)

	opt := outerLens.Compose(shapeHolderLens).ComposeOptional(circleRadiusOptional)

	o := Outer{holder: addrShapeHolder{shape: Circle{Radius: 7.0}}}

	t.Run("Get navigates through Compose+ComposeOptional", func(t *testing.T) {
		got := opt.GetOption(o)
		assert.True(t, OptionIsSome(got))
		assert.Equal(t, 7.0, OptionGetOrElse(F.Constant(0.0))(got))
	})

	t.Run("Set updates through Compose+ComposeOptional", func(t *testing.T) {
		updated := opt.Set(50.0)(o)
		assert.Equal(t, Circle{Radius: 50.0}, updated.holder.shape)
		assert.Equal(t, Circle{Radius: 7.0}, o.holder.shape) // original unchanged
	})

	t.Run("ModifyOption is no-op when optional misses after chained Compose", func(t *testing.T) {
		o2 := Outer{holder: addrShapeHolder{shape: Rectangle{Width: 3, Height: 4}}}
		result := OptionalModifyOption[Outer](func(r float64) float64 { return r * 2 })(opt)(o2)
		assert.Equal(t, OptionNone[Outer](), result)
	})
}

// TestLensComposeOptional_Method_PipelineUsage verifies that the method form and
// the free-function form used inside F.Pipe1 produce identical results.
func TestLensComposeOptional_Method_PipelineUsage(t *testing.T) {
	pipelineOpt := F.Pipe1(shapeHolderLens, LensComposeOptional[addrShapeHolder](circleRadiusOptional))
	methodOpt := shapeHolderLens.ComposeOptional(circleRadiusOptional)

	h := addrShapeHolder{shape: Circle{Radius: 2.5}}
	assert.Equal(t, pipelineOpt.GetOption(h), methodOpt.GetOption(h))
	assert.Equal(t, pipelineOpt.Set(8.0)(h), methodOpt.Set(8.0)(h))
}
