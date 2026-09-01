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
	"github.com/stretchr/testify/assert"
)

// ---- test types ----

type tcCelsius float64
type tcFahrenheit float64

type tcThermometer struct {
	Reading tcCelsius
}

// celsiusToFahrenheitIso is an isomorphism between Celsius and Fahrenheit.
var celsiusToFahrenheitIso = common.MakeIso(
	func(c tcCelsius) tcFahrenheit { return tcFahrenheit(c*9/5 + 32) },
	func(f tcFahrenheit) tcCelsius { return tcCelsius((f - 32) * 5 / 9) },
)

// readingLens focuses on the Reading field of Thermometer.
var readingLens = common.MakeLens(
	func(t tcThermometer) tcCelsius { return t.Reading },
	func(t tcThermometer, c tcCelsius) tcThermometer { t.Reading = c; return t },
)

// fahrenheitLens reads/writes the thermometer's reading in Fahrenheit.
var fahrenheitLens = F.Pipe1(readingLens, lens.ComposeIso[tcThermometer](celsiusToFahrenheitIso))

// ---- tests ----

// TestComposeIso_Get verifies that Get applies the forward direction of the isomorphism.
func TestComposeIso_Get(t *testing.T) {
	therm := tcThermometer{Reading: 100}
	got := fahrenheitLens.Get(therm)
	assert.InDelta(t, 212.0, float64(got), 0.001)
}

// TestComposeIso_Set verifies that Set applies the reverse direction of the isomorphism.
func TestComposeIso_Set(t *testing.T) {
	therm := tcThermometer{Reading: 100}
	updated := fahrenheitLens.Set(tcFahrenheit(32))(therm)
	assert.InDelta(t, 0.0, float64(updated.Reading), 0.001)
}

// TestComposeIso_SetDoesNotMutateOriginal verifies immutability.
func TestComposeIso_SetDoesNotMutateOriginal(t *testing.T) {
	therm := tcThermometer{Reading: 100}
	_ = fahrenheitLens.Set(tcFahrenheit(32))(therm)
	assert.InDelta(t, 100.0, float64(therm.Reading), 0.001)
}

// TestComposeIso_LensLaws verifies all three lens laws for a ComposeIso result.
func TestComposeIso_LensLaws(t *testing.T) {
	therm := tcThermometer{Reading: 20}

	t.Run("law GetSet: Set(Get(s))(s) == s", func(t *testing.T) {
		result := fahrenheitLens.Set(fahrenheitLens.Get(therm))(therm)
		assert.InDelta(t, float64(therm.Reading), float64(result.Reading), 0.001)
	})

	t.Run("law SetGet: Get(Set(b)(s)) == b", func(t *testing.T) {
		b := tcFahrenheit(86)
		got := fahrenheitLens.Get(fahrenheitLens.Set(b)(therm))
		assert.InDelta(t, float64(b), float64(got), 0.001)
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		b1, b2 := tcFahrenheit(68), tcFahrenheit(77)
		result1 := fahrenheitLens.Set(b2)(fahrenheitLens.Set(b1)(therm))
		result2 := fahrenheitLens.Set(b2)(therm)
		assert.InDelta(t, float64(result2.Reading), float64(result1.Reading), 0.001)
	})
}

// TestComposeIso_EquivalentToFreeFunction verifies that ComposeIso in the lens package
// produces the same result as the underlying LensComposeIso free function.
func TestComposeIso_EquivalentToFreeFunction(t *testing.T) {
	therm := tcThermometer{Reading: 25}

	free := common.LensComposeIso[tcThermometer](celsiusToFahrenheitIso)(readingLens)
	method := fahrenheitLens

	t.Run("Get returns same value", func(t *testing.T) {
		assert.InDelta(t, float64(free.Get(therm)), float64(method.Get(therm)), 0.001)
	})

	t.Run("Set returns same result", func(t *testing.T) {
		b := tcFahrenheit(77)
		assert.InDelta(t,
			float64(free.Set(b)(therm).Reading),
			float64(method.Set(b)(therm).Reading),
			0.001,
		)
	})
}

// ---- examples ----

// ExampleComposeIso demonstrates reading and writing a Thermometer's temperature
// in Fahrenheit via a Celsius↔Fahrenheit isomorphism.
func ExampleComposeIso() {
	type Celsius float64
	type Fahrenheit float64
	type Thermometer struct{ Reading Celsius }

	celsiusIso := common.MakeIso(
		func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
		func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
	)
	rLens := common.MakeLens(
		func(t Thermometer) Celsius { return t.Reading },
		func(t Thermometer, c Celsius) Thermometer { t.Reading = c; return t },
	)

	fLens := F.Pipe1(rLens, lens.ComposeIso[Thermometer](celsiusIso))

	t := Thermometer{Reading: 100}
	fmt.Printf("%.0f°F\n", float64(fLens.Get(t)))
	updated := fLens.Set(Fahrenheit(32))(t)
	fmt.Printf("%.0f°C\n", float64(updated.Reading))
	// Output:
	// 212°F
	// 0°C
}

// ExampleComposeIso_newtypeWrapper demonstrates using ComposeIso to unwrap a newtype.
func ExampleComposeIso_newtypeWrapper() {
	type Miles float64
	type Km float64
	type Route struct{ Distance Miles }

	milesToKm := common.MakeIso(
		func(m Miles) Km { return Km(m * 1.60934) },
		func(k Km) Miles { return Miles(k / 1.60934) },
	)
	distLens := common.MakeLens(
		func(r Route) Miles { return r.Distance },
		func(r Route, m Miles) Route { r.Distance = m; return r },
	)

	kmLens := F.Pipe1(distLens, lens.ComposeIso[Route](milesToKm))

	route := Route{Distance: 10}
	fmt.Printf("%.2f km\n", float64(kmLens.Get(route)))
	updated := kmLens.Set(Km(16.0934))(route)
	fmt.Printf("%.2f miles\n", float64(updated.Distance))
	// Output:
	// 16.09 km
	// 10.00 miles
}
