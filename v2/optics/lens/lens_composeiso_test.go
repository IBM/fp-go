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
	"strings"
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

// TestComposeIso_IdentityIso verifies that composing with an identity isomorphism
// leaves Get and Set behaviour completely unchanged.
func TestComposeIso_IdentityIso(t *testing.T) {
	idIso := common.MakeIso(
		func(c tcCelsius) tcCelsius { return c },
		func(c tcCelsius) tcCelsius { return c },
	)
	identityLens := F.Pipe1(readingLens, lens.ComposeIso[tcThermometer](idIso))
	therm := tcThermometer{Reading: 42}

	t.Run("Get is unchanged", func(t *testing.T) {
		assert.Equal(t, tcCelsius(42), identityLens.Get(therm))
	})

	t.Run("Set is unchanged", func(t *testing.T) {
		updated := identityLens.Set(tcCelsius(99))(therm)
		assert.Equal(t, tcCelsius(99), updated.Reading)
	})

	t.Run("original is not mutated", func(t *testing.T) {
		_ = identityLens.Set(tcCelsius(0))(therm)
		assert.Equal(t, tcCelsius(42), therm.Reading)
	})
}

// TestComposeIso_StringIso verifies ComposeIso with a non-numeric string isomorphism,
// where equality can be checked exactly rather than with InDelta.
func TestComposeIso_StringIso(t *testing.T) {
	type Tag string
	type Item struct{ Label Tag }

	upperIso := common.MakeIso(
		func(s Tag) string { return strings.ToUpper(string(s)) },
		func(s string) Tag { return Tag(strings.ToLower(s)) },
	)
	labelLens := common.MakeLens(
		func(i Item) Tag { return i.Label },
		func(i Item, l Tag) Item { i.Label = l; return i },
	)
	upperLabelLens := F.Pipe1(labelLens, lens.ComposeIso[Item](upperIso))

	item := Item{Label: "hello"}

	t.Run("Get applies forward iso", func(t *testing.T) {
		assert.Equal(t, "HELLO", upperLabelLens.Get(item))
	})

	t.Run("Set applies reverse iso", func(t *testing.T) {
		updated := upperLabelLens.Set("WORLD")(item)
		assert.Equal(t, Tag("world"), updated.Label)
	})

	t.Run("law GetSet: Set(Get(s))(s) == s", func(t *testing.T) {
		result := upperLabelLens.Set(upperLabelLens.Get(item))(item)
		assert.Equal(t, item, result)
	})

	t.Run("law SetGet: Get(Set(b)(s)) == b", func(t *testing.T) {
		b := "WORLD"
		got := upperLabelLens.Get(upperLabelLens.Set(b)(item))
		assert.Equal(t, b, got)
	})

	t.Run("law SetSet: Set(b2)(Set(b1)(s)) == Set(b2)(s)", func(t *testing.T) {
		b1, b2 := "FIRST", "SECOND"
		setTwice := upperLabelLens.Set(b2)(upperLabelLens.Set(b1)(item))
		setOnce := upperLabelLens.Set(b2)(item)
		assert.Equal(t, setOnce, setTwice)
	})
}

// TestComposeIso_Chained verifies that ComposeIso can be chained after Compose
// to navigate two levels and then apply an isomorphism.
func TestComposeIso_Chained(t *testing.T) {
	type Inner struct{ Value int }
	type Outer struct{ Inner Inner }
	type Wrapped struct{ N int }

	innerLens := common.MakeLens(
		func(o Outer) Inner { return o.Inner },
		func(o Outer, i Inner) Outer { o.Inner = i; return o },
	)
	valueLens := common.MakeLens(
		func(i Inner) int { return i.Value },
		func(i Inner, v int) Inner { i.Value = v; return i },
	)
	wrapIso := common.MakeIso(
		func(n int) Wrapped { return Wrapped{N: n} },
		func(w Wrapped) int { return w.N },
	)

	// Compose two lenses, then apply the iso: Outer → Inner → int → Wrapped
	wrappedLens := F.Pipe1(
		F.Pipe1(innerLens, lens.Compose[Outer](valueLens)),
		lens.ComposeIso[Outer](wrapIso),
	)

	outer := Outer{Inner: Inner{Value: 7}}

	t.Run("Get navigates two levels and wraps", func(t *testing.T) {
		assert.Equal(t, Wrapped{N: 7}, wrappedLens.Get(outer))
	})

	t.Run("Set unwraps and updates two levels deep", func(t *testing.T) {
		updated := wrappedLens.Set(Wrapped{N: 42})(outer)
		assert.Equal(t, 42, updated.Inner.Value)
	})

	t.Run("original is not mutated", func(t *testing.T) {
		_ = wrappedLens.Set(Wrapped{N: 99})(outer)
		assert.Equal(t, 7, outer.Inner.Value)
	})
}

// TestComposeIso_ModifyThroughIso verifies that Modify works correctly through
// a ComposeIso lens, applying a function in the B-space.
func TestComposeIso_ModifyThroughIso(t *testing.T) {
	// Reuse the package-level Celsius→Fahrenheit lens.
	therm := tcThermometer{Reading: 0} // 0 °C = 32 °F

	double := func(f tcFahrenheit) tcFahrenheit { return f * 2 }
	modified := F.Pipe1(therm, lens.Modify[tcThermometer](double)(fahrenheitLens))

	// 0 °C → 32 °F → doubled to 64 °F → stored back as (64-32)*5/9 = 17.78 °C
	assert.InDelta(t, float64(tcCelsius((64-32)*5.0/9.0)), float64(modified.Reading), 0.001)
}

// TestComposeIso_HasNonEmptyName verifies that the composed lens carries a
// non-empty name (used in String() and structured logging).
func TestComposeIso_HasNonEmptyName(t *testing.T) {
	assert.NotEmpty(t, fahrenheitLens.String())
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

// ExampleComposeIso_stringTransform demonstrates ComposeIso with a string
// isomorphism.  The lens focuses on a Name field; the iso uppercases on Get
// and lowercases on Set so the stored value always stays lower-case.
func ExampleComposeIso_stringTransform() {
	type Person struct{ Name string }

	// iso: stored lower-case <-> displayed upper-case
	upperIso := common.MakeIso(
		strings.ToUpper,
		strings.ToLower,
	)
	nameLens := common.MakeLens(
		func(p Person) string { return p.Name },
		func(p Person, n string) Person { p.Name = n; return p },
	)
	displayLens := F.Pipe1(nameLens, lens.ComposeIso[Person](upperIso))

	p := Person{Name: "alice"}
	fmt.Println(displayLens.Get(p))           // stored "alice" → displayed "ALICE"
	updated := displayLens.Set("BOB")(p)      // "BOB" → stored as "bob"
	fmt.Println(updated.Name)
	// Output:
	// ALICE
	// bob
}

// ExampleComposeIso_chainedCompose demonstrates chaining Compose and ComposeIso
// to navigate two levels of nesting and then apply a unit conversion.
func ExampleComposeIso_chainedCompose() {
	type Engine struct{ Horsepower int }
	type Car struct{ Engine Engine }
	type Kilowatts float64

	engineLens := common.MakeLens(
		func(c Car) Engine { return c.Engine },
		func(c Car, e Engine) Car { c.Engine = e; return c },
	)
	hpLens := common.MakeLens(
		func(e Engine) int { return e.Horsepower },
		func(e Engine, hp int) Engine { e.Horsepower = hp; return e },
	)
	// 1 HP ≈ 0.7457 kW
	hpToKw := common.MakeIso(
		func(hp int) Kilowatts { return Kilowatts(float64(hp) * 0.7457) },
		func(kw Kilowatts) int { return int(float64(kw) / 0.7457) },
	)

	// Car → Engine → int (HP) → Kilowatts
	kwLens := F.Pipe1(
		F.Pipe1(engineLens, lens.Compose[Car](hpLens)),
		lens.ComposeIso[Car](hpToKw),
	)

	car := Car{Engine: Engine{Horsepower: 200}}
	fmt.Printf("%.1f kW\n", float64(kwLens.Get(car)))
	updated := kwLens.Set(Kilowatts(100))(car)
	fmt.Printf("%d HP\n", updated.Engine.Horsepower)
	// Output:
	// 149.1 kW
	// 134 HP
}
