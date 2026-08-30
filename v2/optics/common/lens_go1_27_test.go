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
