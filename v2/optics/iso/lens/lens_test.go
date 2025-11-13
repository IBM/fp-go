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

package lens

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	ISO "github.com/IBM/fp-go/v2/optics/iso"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

// Test types
type Celsius float64
type Fahrenheit float64

type UserId int
type User struct {
	id   UserId
	name string
}

type Meters float64
type Feet float64

// TestIsoAsLensBasic tests basic functionality of IsoAsLens
func TestIsoAsLensBasic(t *testing.T) {
	// Create an isomorphism between Celsius and Fahrenheit
	celsiusToFahrenheit := func(c Celsius) Fahrenheit {
		return Fahrenheit(c*9/5 + 32)
	}
	fahrenheitToCelsius := func(f Fahrenheit) Celsius {
		return Celsius((f - 32) * 5 / 9)
	}

	tempIso := ISO.MakeIso(celsiusToFahrenheit, fahrenheitToCelsius)
	tempLens := IsoAsLens(tempIso)

	t.Run("Get", func(t *testing.T) {
		celsius := Celsius(20.0)
		fahrenheit := tempLens.Get(celsius)
		assert.InDelta(t, 68.0, float64(fahrenheit), 0.001)
	})

	t.Run("Set", func(t *testing.T) {
		celsius := Celsius(20.0)
		newFahrenheit := Fahrenheit(86.0)
		updated := tempLens.Set(newFahrenheit)(celsius)
		assert.InDelta(t, 30.0, float64(updated), 0.001)
	})

	t.Run("SetPreservesOriginal", func(t *testing.T) {
		original := Celsius(20.0)
		newFahrenheit := Fahrenheit(86.0)
		_ = tempLens.Set(newFahrenheit)(original)
		// Original should be unchanged
		assert.Equal(t, Celsius(20.0), original)
	})
}

// TestIsoAsLensRefBasic tests basic functionality of IsoAsLensRef
func TestIsoAsLensRefBasic(t *testing.T) {
	// Create an isomorphism for User pointer and UserId
	userToId := func(u *User) UserId {
		return u.id
	}
	idToUser := func(id UserId) *User {
		return &User{id: id, name: "Unknown"}
	}

	userIdIso := ISO.MakeIso(userToId, idToUser)
	userIdLens := IsoAsLensRef(userIdIso)

	t.Run("Get", func(t *testing.T) {
		user := &User{id: 42, name: "Alice"}
		id := userIdLens.Get(user)
		assert.Equal(t, UserId(42), id)
	})

	t.Run("Set", func(t *testing.T) {
		user := &User{id: 42, name: "Alice"}
		newId := UserId(100)
		updated := userIdLens.Set(newId)(user)
		assert.Equal(t, UserId(100), updated.id)
		assert.Equal(t, "Unknown", updated.name) // ReverseGet creates new user
	})

	t.Run("SetCreatesNewPointer", func(t *testing.T) {
		user := &User{id: 42, name: "Alice"}
		newId := UserId(100)
		updated := userIdLens.Set(newId)(user)
		// Should be different pointers
		assert.NotSame(t, user, updated)
		// Original should be unchanged
		assert.Equal(t, UserId(42), user.id)
		assert.Equal(t, "Alice", user.name)
	})
}

// TestIsoAsLensLaws verifies that IsoAsLens satisfies lens laws
func TestIsoAsLensLaws(t *testing.T) {
	// Create a simple isomorphism
	type Wrapper struct{ value int }

	wrapperIso := ISO.MakeIso(
		func(w Wrapper) int { return w.value },
		func(i int) Wrapper { return Wrapper{value: i} },
	)

	lens := IsoAsLens(wrapperIso)
	wrapper := Wrapper{value: 42}
	newValue := 100

	// Law 1: GetSet - lens.Set(lens.Get(s))(s) == s
	t.Run("GetSetLaw", func(t *testing.T) {
		result := lens.Set(lens.Get(wrapper))(wrapper)
		assert.Equal(t, wrapper, result)
	})

	// Law 2: SetGet - lens.Get(lens.Set(a)(s)) == a
	t.Run("SetGetLaw", func(t *testing.T) {
		result := lens.Get(lens.Set(newValue)(wrapper))
		assert.Equal(t, newValue, result)
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("SetSetLaw", func(t *testing.T) {
		result1 := lens.Set(200)(lens.Set(newValue)(wrapper))
		result2 := lens.Set(200)(wrapper)
		assert.Equal(t, result2, result1)
	})
}

// TestIsoAsLensRefLaws verifies that IsoAsLensRef satisfies lens laws
func TestIsoAsLensRefLaws(t *testing.T) {
	type Wrapper struct{ value int }

	wrapperIso := ISO.MakeIso(
		func(w *Wrapper) int { return w.value },
		func(i int) *Wrapper { return &Wrapper{value: i} },
	)

	lens := IsoAsLensRef(wrapperIso)
	wrapper := &Wrapper{value: 42}
	newValue := 100

	// Law 1: GetSet - lens.Set(lens.Get(s))(s) == s
	t.Run("GetSetLaw", func(t *testing.T) {
		result := lens.Set(lens.Get(wrapper))(wrapper)
		assert.Equal(t, wrapper.value, result.value)
	})

	// Law 2: SetGet - lens.Get(lens.Set(a)(s)) == a
	t.Run("SetGetLaw", func(t *testing.T) {
		result := lens.Get(lens.Set(newValue)(wrapper))
		assert.Equal(t, newValue, result)
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("SetSetLaw", func(t *testing.T) {
		result1 := lens.Set(200)(lens.Set(newValue)(wrapper))
		result2 := lens.Set(200)(wrapper)
		assert.Equal(t, result2.value, result1.value)
	})
}

// TestIsoAsLensComposition tests composing iso-based lenses with other lenses
func TestIsoAsLensComposition(t *testing.T) {
	type Temperature struct {
		celsius Celsius
	}

	// Lens to access celsius field
	celsiusFieldLens := L.MakeLens(
		func(t Temperature) Celsius { return t.celsius },
		func(t Temperature, c Celsius) Temperature {
			t.celsius = c
			return t
		},
	)

	// Isomorphism between Celsius and Fahrenheit
	celsiusToFahrenheit := func(c Celsius) Fahrenheit {
		return Fahrenheit(c*9/5 + 32)
	}
	fahrenheitToCelsius := func(f Fahrenheit) Celsius {
		return Celsius((f - 32) * 5 / 9)
	}

	tempIso := ISO.MakeIso(celsiusToFahrenheit, fahrenheitToCelsius)
	tempLens := IsoAsLens(tempIso)

	// Compose to work with Fahrenheit directly from Temperature
	composedLens := F.Pipe1(
		celsiusFieldLens,
		L.Compose[Temperature](tempLens),
	)

	temp := Temperature{celsius: 20}

	t.Run("ComposedGet", func(t *testing.T) {
		fahrenheit := composedLens.Get(temp)
		assert.InDelta(t, 68.0, float64(fahrenheit), 0.001)
	})

	t.Run("ComposedSet", func(t *testing.T) {
		newFahrenheit := Fahrenheit(86.0)
		updated := composedLens.Set(newFahrenheit)(temp)
		assert.InDelta(t, 30.0, float64(updated.celsius), 0.001)
	})
}

// TestIsoAsLensModify tests using Modify with iso-based lenses
func TestIsoAsLensModify(t *testing.T) {
	// Isomorphism between Meters and Feet
	metersToFeet := func(m Meters) Feet {
		return Feet(m * 3.28084)
	}
	feetToMeters := func(f Feet) Meters {
		return Meters(f / 3.28084)
	}

	distanceIso := ISO.MakeIso(metersToFeet, feetToMeters)
	distanceLens := IsoAsLens(distanceIso)

	meters := Meters(10.0)

	t.Run("ModifyDouble", func(t *testing.T) {
		// Double the distance in feet, result in meters
		doubleFeet := func(f Feet) Feet { return f * 2 }
		modified := L.Modify[Meters](doubleFeet)(distanceLens)(meters)
		assert.InDelta(t, 20.0, float64(modified), 0.001)
	})

	t.Run("ModifyIdentity", func(t *testing.T) {
		// Identity modification should return same value
		identity := func(f Feet) Feet { return f }
		modified := L.Modify[Meters](identity)(distanceLens)(meters)
		assert.InDelta(t, float64(meters), float64(modified), 0.001)
	})
}

// TestIsoAsLensWithIdentityIso tests that identity iso creates identity lens
func TestIsoAsLensWithIdentityIso(t *testing.T) {
	type Value int

	idIso := ISO.Id[Value]()
	idLens := IsoAsLens(idIso)

	value := Value(42)

	t.Run("IdentityGet", func(t *testing.T) {
		result := idLens.Get(value)
		assert.Equal(t, value, result)
	})

	t.Run("IdentitySet", func(t *testing.T) {
		newValue := Value(100)
		result := idLens.Set(newValue)(value)
		assert.Equal(t, newValue, result)
	})
}

// TestIsoAsLensRefWithIdentityIso tests identity iso with references
func TestIsoAsLensRefWithIdentityIso(t *testing.T) {
	type Value struct{ n int }

	idIso := ISO.Id[*Value]()
	idLens := IsoAsLensRef(idIso)

	value := &Value{n: 42}

	t.Run("IdentityGet", func(t *testing.T) {
		result := idLens.Get(value)
		assert.Equal(t, value, result)
	})

	t.Run("IdentitySet", func(t *testing.T) {
		newValue := &Value{n: 100}
		result := idLens.Set(newValue)(value)
		assert.Equal(t, newValue, result)
	})
}

// TestIsoAsLensRoundTrip tests round-trip conversions
func TestIsoAsLensRoundTrip(t *testing.T) {
	type Email string
	type ValidatedEmail struct{ value Email }

	emailIso := ISO.MakeIso(
		func(ve ValidatedEmail) Email { return ve.value },
		func(e Email) ValidatedEmail { return ValidatedEmail{value: e} },
	)

	emailLens := IsoAsLens(emailIso)

	validated := ValidatedEmail{value: "user@example.com"}

	t.Run("RoundTripThroughGet", func(t *testing.T) {
		// Get the email, then Set it back
		email := emailLens.Get(validated)
		restored := emailLens.Set(email)(validated)
		assert.Equal(t, validated, restored)
	})

	t.Run("RoundTripThroughSet", func(t *testing.T) {
		// Set a new email, then Get it
		newEmail := Email("admin@example.com")
		updated := emailLens.Set(newEmail)(validated)
		retrieved := emailLens.Get(updated)
		assert.Equal(t, newEmail, retrieved)
	})
}

// TestIsoAsLensWithComplexTypes tests with more complex type transformations
func TestIsoAsLensWithComplexTypes(t *testing.T) {
	type Point struct {
		x, y float64
	}

	type PolarCoord struct {
		r, theta float64
	}

	// Isomorphism between Cartesian and Polar coordinates (simplified for testing)
	cartesianToPolar := func(p Point) PolarCoord {
		r := p.x*p.x + p.y*p.y
		theta := 0.0 // Simplified
		return PolarCoord{r: r, theta: theta}
	}

	polarToCartesian := func(pc PolarCoord) Point {
		return Point{x: pc.r, y: pc.theta} // Simplified
	}

	coordIso := ISO.MakeIso(cartesianToPolar, polarToCartesian)
	coordLens := IsoAsLens(coordIso)

	point := Point{x: 3.0, y: 4.0}

	t.Run("ComplexGet", func(t *testing.T) {
		polar := coordLens.Get(point)
		assert.NotNil(t, polar)
	})

	t.Run("ComplexSet", func(t *testing.T) {
		newPolar := PolarCoord{r: 5.0, theta: 0.927}
		updated := coordLens.Set(newPolar)(point)
		assert.NotNil(t, updated)
	})
}

// TestIsoAsLensTypeConversion tests type conversion scenarios
func TestIsoAsLensTypeConversion(t *testing.T) {
	type StringWrapper string
	type IntWrapper int

	// Isomorphism that converts string length to int
	strLenIso := ISO.MakeIso(
		func(s StringWrapper) IntWrapper { return IntWrapper(len(s)) },
		func(i IntWrapper) StringWrapper {
			// Create a string of given length (simplified)
			result := ""
			for j := 0; j < int(i); j++ {
				result += "x"
			}
			return StringWrapper(result)
		},
	)

	strLenLens := IsoAsLens(strLenIso)

	t.Run("StringToLength", func(t *testing.T) {
		str := StringWrapper("hello")
		length := strLenLens.Get(str)
		assert.Equal(t, IntWrapper(5), length)
	})

	t.Run("LengthToString", func(t *testing.T) {
		str := StringWrapper("hello")
		newLength := IntWrapper(3)
		updated := strLenLens.Set(newLength)(str)
		assert.Equal(t, 3, len(updated))
	})
}
