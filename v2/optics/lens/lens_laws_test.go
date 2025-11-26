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

	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// TestModify tests the Modify function
func TestModify(t *testing.T) {
	type Counter struct {
		Value int
	}

	valueLens := MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter {
			c.Value = v
			return c
		},
	)

	counter := Counter{Value: 5}

	// Test increment
	increment := func(v int) int { return v + 1 }
	modifyIncrement := Modify[Counter](increment)(valueLens)
	incremented := modifyIncrement(counter)
	assert.Equal(t, 6, incremented.Value)
	assert.Equal(t, 5, counter.Value) // Original unchanged

	// Test double
	double := func(v int) int { return v * 2 }
	modifyDouble := Modify[Counter](double)(valueLens)
	doubled := modifyDouble(counter)
	assert.Equal(t, 10, doubled.Value)
	assert.Equal(t, 5, counter.Value) // Original unchanged

	// Test identity (no change)
	identity := func(v int) int { return v }
	modifyIdentity := Modify[Counter](identity)(valueLens)
	unchanged := modifyIdentity(counter)
	assert.Equal(t, counter, unchanged)
}

func TestModifyRef(t *testing.T) {
	valueLens := MakeLensRef(
		func(s *Street) int { return s.num },
		func(s *Street, num int) *Street {
			s.num = num
			return s
		},
	)

	street := &Street{num: 10, name: "Main"}

	// Test increment
	increment := func(v int) int { return v + 1 }
	modifyIncrement := Modify[*Street](increment)(valueLens)
	incremented := modifyIncrement(street)
	assert.Equal(t, 11, incremented.num)
	assert.Equal(t, 10, street.num) // Original unchanged
}

// Lens Laws Tests

func TestMakeLensLaws(t *testing.T) {
	nameLens := MakeLens(
		func(s Street) string { return s.name },
		func(s Street, name string) Street {
			s.name = name
			return s
		},
	)

	street := Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet - lens.Set(lens.Get(s))(s) == s
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street, result)
	})

	// Law 2: SetGet - lens.Get(lens.Set(a)(s)) == a
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2, result1)
	})
}

func TestMakeLensRefLaws(t *testing.T) {
	nameLens := MakeLensRef(
		(*Street).GetName,
		(*Street).SetName,
	)

	street := &Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet - lens.Set(lens.Get(s))(s) == s
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street.name, result.name)
		assert.Equal(t, street.num, result.num)
	})

	// Law 2: SetGet - lens.Get(lens.Set(a)(s)) == a
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet - lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2.name, result1.name)
		assert.Equal(t, result2.num, result1.num)
	})
}

func TestMakeLensCurriedLaws(t *testing.T) {
	nameLens := MakeLensCurried(
		func(s Street) string { return s.name },
		func(name string) func(Street) Street {
			return func(s Street) Street {
				s.name = name
				return s
			}
		},
	)

	street := Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street, result)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2, result1)
	})
}

func TestMakeLensRefCurriedLaws(t *testing.T) {
	nameLens := MakeLensRefCurried(
		func(s *Street) string { return s.name },
		func(name string) func(*Street) *Street {
			return func(s *Street) *Street {
				s.name = name
				return s
			}
		},
	)

	street := &Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street.name, result.name)
		assert.Equal(t, street.num, result.num)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2.name, result1.name)
		assert.Equal(t, result2.num, result1.num)
	})
}

func TestMakeLensWithEqLaws(t *testing.T) {
	nameLens := MakeLensWithEq(
		EQ.FromStrictEquals[string](),
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street.name, result.name)
		assert.Equal(t, street.num, result.num)
		// With Eq optimization, should return same pointer
		assert.Same(t, street, result)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2.name, result1.name)
		assert.Equal(t, result2.num, result1.num)
	})
}

func TestMakeLensStrictLaws(t *testing.T) {
	nameLens := MakeLensStrict(
		func(s *Street) string { return s.name },
		func(s *Street, name string) *Street {
			s.name = name
			return s
		},
	)

	street := &Street{num: 1, name: "Main"}
	newName := "Oak"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := nameLens.Set(nameLens.Get(street))(street)
		assert.Equal(t, street.name, result.name)
		assert.Equal(t, street.num, result.num)
		// With strict equality optimization, should return same pointer
		assert.Same(t, street, result)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := nameLens.Get(nameLens.Set(newName)(street))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := nameLens.Set("Elm")(nameLens.Set(newName)(street))
		result2 := nameLens.Set("Elm")(street)
		assert.Equal(t, result2.name, result1.name)
		assert.Equal(t, result2.num, result1.num)
	})
}

func TestIdLaws(t *testing.T) {
	idLens := Id[Street]()
	street := Street{num: 1, name: "Main"}
	newStreet := Street{num: 2, name: "Oak"}

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := idLens.Set(idLens.Get(street))(street)
		assert.Equal(t, street, result)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := idLens.Get(idLens.Set(newStreet)(street))
		assert.Equal(t, newStreet, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		anotherStreet := Street{num: 3, name: "Elm"}
		result1 := idLens.Set(anotherStreet)(idLens.Set(newStreet)(street))
		result2 := idLens.Set(anotherStreet)(street)
		assert.Equal(t, result2, result1)
	})
}

func TestIdRefLaws(t *testing.T) {
	idLens := IdRef[Street]()
	street := &Street{num: 1, name: "Main"}
	newStreet := &Street{num: 2, name: "Oak"}

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := idLens.Set(idLens.Get(street))(street)
		assert.Equal(t, street.name, result.name)
		assert.Equal(t, street.num, result.num)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := idLens.Get(idLens.Set(newStreet)(street))
		assert.Equal(t, newStreet, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		anotherStreet := &Street{num: 3, name: "Elm"}
		result1 := idLens.Set(anotherStreet)(idLens.Set(newStreet)(street))
		result2 := idLens.Set(anotherStreet)(street)
		assert.Equal(t, result2, result1)
	})
}

func TestComposeLaws(t *testing.T) {
	streetLens := MakeLensRef((*Street).GetName, (*Street).SetName)
	addrLens := MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

	// Compose to get street name from address
	streetNameLens := Compose[*Address](streetLens)(addrLens)

	sampleStreet := Street{num: 220, name: "Schönaicherstr"}
	sampleAddress := Address{city: "Böblingen", street: &sampleStreet}
	newName := "Böblingerstr"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := streetNameLens.Set(streetNameLens.Get(&sampleAddress))(&sampleAddress)
		assert.Equal(t, sampleAddress.street.name, result.street.name)
		assert.Equal(t, sampleAddress.street.num, result.street.num)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := streetNameLens.Get(streetNameLens.Set(newName)(&sampleAddress))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := streetNameLens.Set("Elm St")(streetNameLens.Set(newName)(&sampleAddress))
		result2 := streetNameLens.Set("Elm St")(&sampleAddress)
		assert.Equal(t, result2.street.name, result1.street.name)
	})
}

func TestComposeRefLaws(t *testing.T) {
	streetLens := MakeLensRef((*Street).GetName, (*Street).SetName)
	addrLens := MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

	// Compose using ComposeRef
	streetNameLens := ComposeRef[Address](streetLens)(addrLens)

	sampleStreet := Street{num: 220, name: "Schönaicherstr"}
	sampleAddress := Address{city: "Böblingen", street: &sampleStreet}
	newName := "Böblingerstr"

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := streetNameLens.Set(streetNameLens.Get(&sampleAddress))(&sampleAddress)
		assert.Equal(t, sampleAddress.street.name, result.street.name)
		assert.Equal(t, sampleAddress.street.num, result.street.num)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := streetNameLens.Get(streetNameLens.Set(newName)(&sampleAddress))
		assert.Equal(t, newName, result)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		result1 := streetNameLens.Set("Elm St")(streetNameLens.Set(newName)(&sampleAddress))
		result2 := streetNameLens.Set("Elm St")(&sampleAddress)
		assert.Equal(t, result2.street.name, result1.street.name)
	})
}

func TestIMapLaws(t *testing.T) {
	type Celsius float64
	type Fahrenheit float64

	celsiusToFahrenheit := func(c Celsius) Fahrenheit {
		return Fahrenheit(c*9/5 + 32)
	}

	fahrenheitToCelsius := func(f Fahrenheit) Celsius {
		return Celsius((f - 32) * 5 / 9)
	}

	type Weather struct {
		Temperature Celsius
	}

	tempCelsiusLens := MakeLens(
		func(w Weather) Celsius { return w.Temperature },
		func(w Weather, t Celsius) Weather {
			w.Temperature = t
			return w
		},
	)

	// Create a lens that works with Fahrenheit
	tempFahrenheitLens := F.Pipe1(
		tempCelsiusLens,
		IMap[Weather](celsiusToFahrenheit, fahrenheitToCelsius),
	)

	weather := Weather{Temperature: 20} // 20°C
	newTempF := Fahrenheit(86)          // 86°F (30°C)

	// Law 1: GetSet
	t.Run("GetSet", func(t *testing.T) {
		result := tempFahrenheitLens.Set(tempFahrenheitLens.Get(weather))(weather)
		// Allow small floating point differences
		assert.InDelta(t, float64(weather.Temperature), float64(result.Temperature), 0.0001)
	})

	// Law 2: SetGet
	t.Run("SetGet", func(t *testing.T) {
		result := tempFahrenheitLens.Get(tempFahrenheitLens.Set(newTempF)(weather))
		assert.InDelta(t, float64(newTempF), float64(result), 0.0001)
	})

	// Law 3: SetSet
	t.Run("SetSet", func(t *testing.T) {
		anotherTempF := Fahrenheit(95) // 95°F (35°C)
		result1 := tempFahrenheitLens.Set(anotherTempF)(tempFahrenheitLens.Set(newTempF)(weather))
		result2 := tempFahrenheitLens.Set(anotherTempF)(weather)
		assert.InDelta(t, float64(result2.Temperature), float64(result1.Temperature), 0.0001)
	})
}

func TestIMapIdentity(t *testing.T) {
	// IMap with identity functions should behave like the original lens
	type S struct {
		a int
	}

	originalLens := MakeLens(
		func(s S) int { return s.a },
		func(s S, a int) S {
			s.a = a
			return s
		},
	)

	// Apply IMap with identity functions
	identityMappedLens := F.Pipe1(
		originalLens,
		IMap[S](F.Identity[int], F.Identity[int]),
	)

	s := S{a: 42}

	// Both lenses should behave identically
	assert.Equal(t, originalLens.Get(s), identityMappedLens.Get(s))
	assert.Equal(t, originalLens.Set(100)(s), identityMappedLens.Set(100)(s))
}

func TestIMapComposition(t *testing.T) {
	// IMap(f, g) ∘ IMap(h, k) = IMap(f ∘ h, k ∘ g)
	type S struct {
		value int
	}

	baseLens := MakeLens(
		func(s S) int { return s.value },
		func(s S, v int) S {
			s.value = v
			return s
		},
	)

	// First transformation: int -> float64
	intToFloat := func(i int) float64 { return float64(i) }
	floatToInt := func(f float64) int { return int(f) }

	// Second transformation: float64 -> string
	floatToString := func(f float64) string { return F.Pipe1(f, func(x float64) string { return "value" }) }
	stringToFloat := func(s string) float64 { return 42.0 }

	// Compose IMap twice
	lens1 := F.Pipe1(baseLens, IMap[S](intToFloat, floatToInt))
	lens2 := F.Pipe1(lens1, IMap[S](floatToString, stringToFloat))

	// Direct composition
	lens3 := F.Pipe1(
		baseLens,
		IMap[S](
			F.Flow2(intToFloat, floatToString),
			F.Flow2(stringToFloat, floatToInt),
		),
	)

	s := S{value: 10}

	// Both should produce the same results
	assert.Equal(t, lens2.Get(s), lens3.Get(s))
	assert.Equal(t, lens2.Set("test")(s), lens3.Set("test")(s))
}

func TestModifyLaws(t *testing.T) {
	// Modify should satisfy: Modify(id) = id
	// and: Modify(f ∘ g) = Modify(f) ∘ Modify(g)

	type S struct {
		value int
	}

	lens := MakeLens(
		func(s S) int { return s.value },
		func(s S, v int) S {
			s.value = v
			return s
		},
	)

	s := S{value: 10}

	// Modify with identity should return the same value
	t.Run("ModifyIdentity", func(t *testing.T) {
		modifyIdentity := Modify[S](F.Identity[int])(lens)
		result := modifyIdentity(s)
		assert.Equal(t, s, result)
	})

	// Modify composition: Modify(f ∘ g) = Modify(f) ∘ Modify(g)
	t.Run("ModifyComposition", func(t *testing.T) {
		f := N.Mul(2)
		g := N.Add(3)

		// Modify(f ∘ g)
		composed := F.Flow2(g, f)
		modifyComposed := Modify[S](composed)(lens)
		result1 := modifyComposed(s)

		// Modify(f) ∘ Modify(g)
		modifyG := Modify[S](g)(lens)
		intermediate := modifyG(s)
		modifyF := Modify[S](f)(lens)
		result2 := modifyF(intermediate)

		assert.Equal(t, result1, result2)
	})
}

func TestComposeAssociativity(t *testing.T) {
	// Test that lens composition is associative:
	// (l1 ∘ l2) ∘ l3 = l1 ∘ (l2 ∘ l3)

	type Level3 struct {
		value string
	}

	type Level2 struct {
		level3 Level3
	}

	type Level1 struct {
		level2 Level2
	}

	lens12 := MakeLens(
		func(l1 Level1) Level2 { return l1.level2 },
		func(l1 Level1, l2 Level2) Level1 {
			l1.level2 = l2
			return l1
		},
	)

	lens23 := MakeLens(
		func(l2 Level2) Level3 { return l2.level3 },
		func(l2 Level2, l3 Level3) Level2 {
			l2.level3 = l3
			return l2
		},
	)

	lens3Value := MakeLens(
		func(l3 Level3) string { return l3.value },
		func(l3 Level3, v string) Level3 {
			l3.value = v
			return l3
		},
	)

	// (lens12 ∘ lens23) ∘ lens3Value
	composed1 := F.Pipe2(
		lens12,
		Compose[Level1](lens23),
		Compose[Level1](lens3Value),
	)

	// lens12 ∘ (lens23 ∘ lens3Value)
	composed2 := F.Pipe1(
		lens12,
		Compose[Level1](F.Pipe1(lens23, Compose[Level2](lens3Value))),
	)

	l1 := Level1{
		level2: Level2{
			level3: Level3{value: "test"},
		},
	}

	// Both compositions should behave identically
	assert.Equal(t, composed1.Get(l1), composed2.Get(l1))
	assert.Equal(t, composed1.Set("new")(l1), composed2.Set("new")(l1))
}
