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

package option

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type (
	Street struct {
		name string
	}

	Address struct {
		street *Street
	}

	Inner struct {
		Value int
		Foo   string
	}

	InnerOpt struct {
		Value *int
		Foo   *string
	}

	Outer struct {
		inner *Inner
	}

	OuterOpt struct {
		inner *InnerOpt
	}
)

func (outer Outer) GetInner() *Inner {
	return outer.inner
}

func (outer Outer) SetInner(inner *Inner) Outer {
	outer.inner = inner
	return outer
}

func (outer OuterOpt) GetInnerOpt() *InnerOpt {
	return outer.inner
}

func (outer OuterOpt) SetInnerOpt(inner *InnerOpt) OuterOpt {
	outer.inner = inner
	return outer
}

func (inner *Inner) GetValue() int {
	return inner.Value
}

func (inner *Inner) SetValue(value int) *Inner {
	inner.Value = value
	return inner
}

func (inner *InnerOpt) GetValue() *int {
	return inner.Value
}

func (inner *InnerOpt) SetValue(value *int) *InnerOpt {
	inner.Value = value
	return inner
}

func (street *Street) GetName() string {
	return street.name
}

func (street *Street) SetName(name string) *Street {
	street.name = name
	return street
}

func (addr *Address) GetStreet() *Street {
	return addr.street
}

func (addr *Address) SetStreet(s *Street) *Address {
	addr.street = s
	return addr
}

var (
	streetLens = L.MakeLensRef((*Street).GetName, (*Street).SetName)
	addrLens   = L.MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

	sampleStreet  = Street{name: "Schönaicherstr"}
	sampleAddress = Address{street: &sampleStreet}
)

func TestComposeOption(t *testing.T) {
	// default inner object
	defaultInner := &Inner{
		Value: 0,
		Foo:   "foo",
	}
	// access to the value
	value := L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	// access to inner
	inner := FromNillable(L.MakeLens(Outer.GetInner, Outer.SetInner))
	// compose lenses
	lens := F.Pipe1(
		inner,
		ComposeOption[Outer, int](defaultInner)(value),
	)
	outer1 := Outer{inner: &Inner{Value: 1, Foo: "a"}}
	// the checks
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(O.Some(1))(Outer{}))
	assert.Equal(t, O.None[int](), lens.Get(Outer{}))
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(O.Some(1))(Outer{inner: &Inner{Value: 2, Foo: "foo"}}))
	assert.Equal(t, O.Some(1), lens.Get(Outer{inner: &Inner{Value: 1, Foo: "foo"}}))
	assert.Equal(t, outer1, L.Modify[Outer](F.Identity[Option[int]])(lens)(outer1))
}

func TestComposeOptions(t *testing.T) {
	// default inner object
	defaultValue1 := 1
	defaultFoo1 := "foo1"
	defaultInner := &InnerOpt{
		Value: &defaultValue1,
		Foo:   &defaultFoo1,
	}
	// access to the value
	value := FromNillable(L.MakeLensRef((*InnerOpt).GetValue, (*InnerOpt).SetValue))
	// access to inner
	inner := FromNillable(L.MakeLens(OuterOpt.GetInnerOpt, OuterOpt.SetInnerOpt))
	// compose lenses
	lens := F.Pipe1(
		inner,
		Compose[OuterOpt, *int](defaultInner)(value),
	)
	// additional settings
	defaultValue2 := 2
	defaultFoo2 := "foo2"
	outer1 := OuterOpt{inner: &InnerOpt{Value: &defaultValue2, Foo: &defaultFoo2}}
	// the checks
	assert.Equal(t, OuterOpt{inner: &InnerOpt{Value: &defaultValue1, Foo: &defaultFoo1}}, lens.Set(O.Some(&defaultValue1))(OuterOpt{}))
	assert.Equal(t, O.None[*int](), lens.Get(OuterOpt{}))
	assert.Equal(t, OuterOpt{inner: &InnerOpt{Value: &defaultValue1, Foo: &defaultFoo2}}, lens.Set(O.Some(&defaultValue1))(OuterOpt{inner: &InnerOpt{Value: &defaultValue2, Foo: &defaultFoo2}}))
	assert.Equal(t, O.Some(&defaultValue1), lens.Get(OuterOpt{inner: &InnerOpt{Value: &defaultValue1, Foo: &defaultFoo1}}))
	assert.Equal(t, outer1, L.Modify[OuterOpt](F.Identity[Option[*int]])(lens)(outer1))
}

func TestFromNullableProp(t *testing.T) {
	// default inner object
	defaultInner := &Inner{
		Value: 0,
		Foo:   "foo",
	}
	// access to the value
	value := L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	// access to inner
	inner := FromNullableProp[Outer](O.FromNillable[Inner], defaultInner)(L.MakeLens(Outer.GetInner, Outer.SetInner))
	// compose
	lens := F.Pipe1(
		inner,
		L.Compose[Outer](value),
	)
	outer1 := Outer{inner: &Inner{Value: 1, Foo: "a"}}
	// the checks
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(1)(Outer{}))
	assert.Equal(t, 0, lens.Get(Outer{}))
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(1)(Outer{inner: &Inner{Value: 2, Foo: "foo"}}))
	assert.Equal(t, 1, lens.Get(Outer{inner: &Inner{Value: 1, Foo: "foo"}}))
	assert.Equal(t, outer1, L.Modify[Outer](F.Identity[int])(lens)(outer1))
}

func TestFromPredicateRef(t *testing.T) {
	type Person struct {
		age int
	}

	ageLens := L.MakeLensRef(
		func(p *Person) int { return p.age },
		func(p *Person, age int) *Person {
			p.age = age
			return p
		},
	)

	adultLens := FromPredicateRef[Person](func(age int) bool { return age >= 18 }, 0)(ageLens)

	adult := &Person{age: 25}
	assert.Equal(t, O.Some(25), adultLens.Get(adult))

	minor := &Person{age: 15}
	assert.Equal(t, O.None[int](), adultLens.Get(minor))
}

func TestFromNillableRef(t *testing.T) {
	type Config struct {
		timeout *int
	}

	timeoutLens := L.MakeLensRef(
		func(c *Config) *int { return c.timeout },
		func(c *Config, t *int) *Config {
			c.timeout = t
			return c
		},
	)

	optLens := FromNillableRef(timeoutLens)

	config := &Config{timeout: nil}
	assert.Equal(t, O.None[*int](), optLens.Get(config))

	timeout := 30
	configWithTimeout := &Config{timeout: &timeout}
	assert.True(t, O.IsSome(optLens.Get(configWithTimeout)))
}

func TestFromNullablePropRef(t *testing.T) {
	type Config struct {
		timeout *int
	}

	timeoutLens := L.MakeLensRef(
		func(c *Config) *int { return c.timeout },
		func(c *Config, t *int) *Config {
			c.timeout = t
			return c
		},
	)

	defaultTimeout := 30
	safeLens := FromNullablePropRef[Config](O.FromNillable[int], &defaultTimeout)(timeoutLens)

	config := &Config{timeout: nil}
	assert.Equal(t, &defaultTimeout, safeLens.Get(config))
}

func TestFromOptionRef(t *testing.T) {
	type Settings struct {
		retries Option[int]
	}

	retriesLens := L.MakeLensRef(
		func(s *Settings) Option[int] { return s.retries },
		func(s *Settings, r Option[int]) *Settings {
			s.retries = r
			return s
		},
	)

	safeLens := FromOptionRef[Settings](3)(retriesLens)

	settings := &Settings{retries: O.None[int]()}
	assert.Equal(t, 3, safeLens.Get(settings))

	settingsWithRetries := &Settings{retries: O.Some(5)}
	assert.Equal(t, 5, safeLens.Get(settingsWithRetries))
}

func TestFromOption(t *testing.T) {
	type Config struct {
		retries Option[int]
	}

	retriesLens := L.MakeLens(
		func(c Config) Option[int] { return c.retries },
		func(c Config, r Option[int]) Config { c.retries = r; return c },
	)

	defaultRetries := 3
	safeLens := FromOption[Config](defaultRetries)(retriesLens)

	// Test with None - should return default
	config := Config{retries: O.None[int]()}
	assert.Equal(t, defaultRetries, safeLens.Get(config))

	// Test with Some - should return the value
	configWithRetries := Config{retries: O.Some(5)}
	assert.Equal(t, 5, safeLens.Get(configWithRetries))

	// Test setter - should always set Some
	updated := safeLens.Set(10)(config)
	assert.Equal(t, O.Some(10), updated.retries)

	// Test setter on existing Some - should replace
	updated2 := safeLens.Set(7)(configWithRetries)
	assert.Equal(t, O.Some(7), updated2.retries)
}

func TestAsTraversal(t *testing.T) {
	type Data struct {
		value int
	}

	valueLens := L.MakeLens(
		func(d Data) int { return d.value },
		func(d Data, v int) Data { d.value = v; return d },
	)

	// Convert lens to traversal
	traversal := AsTraversal[Data, int]()(valueLens)

	// Test that traversal is created (basic smoke test)
	assert.NotNil(t, traversal)

	// The traversal should work with the data
	data := Data{value: 42}

	// Verify the traversal can be used (it's a function that takes a functor)
	// This is a basic smoke test to ensure the conversion works
	assert.NotNil(t, data)
	assert.Equal(t, 42, valueLens.Get(data))
}

func TestComposeOptionsEdgeCases(t *testing.T) {
	// Test setting None when inner doesn't exist
	defaultValue1 := 1
	defaultFoo1 := "foo1"
	defaultInner := &InnerOpt{
		Value: &defaultValue1,
		Foo:   &defaultFoo1,
	}

	value := FromNillable(L.MakeLensRef((*InnerOpt).GetValue, (*InnerOpt).SetValue))
	inner := FromNillable(L.MakeLens(OuterOpt.GetInnerOpt, OuterOpt.SetInnerOpt))
	lens := F.Pipe1(
		inner,
		Compose[OuterOpt, *int](defaultInner)(value),
	)

	// Setting None when inner doesn't exist should be a no-op
	emptyOuter := OuterOpt{}
	result := lens.Set(O.None[*int]())(emptyOuter)
	assert.Equal(t, O.None[*InnerOpt](), inner.Get(result))

	// Setting None when inner exists should unset the value
	defaultValue2 := 2
	defaultFoo2 := "foo2"
	outerWithInner := OuterOpt{inner: &InnerOpt{Value: &defaultValue2, Foo: &defaultFoo2}}
	result2 := lens.Set(O.None[*int]())(outerWithInner)
	assert.NotNil(t, result2.inner)
	assert.Nil(t, result2.inner.Value)
	assert.Equal(t, &defaultFoo2, result2.inner.Foo)
}

func TestComposeOptionEdgeCases(t *testing.T) {
	defaultInner := &Inner{
		Value: 0,
		Foo:   "foo",
	}

	value := L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	inner := FromNillable(L.MakeLens(Outer.GetInner, Outer.SetInner))
	lens := F.Pipe1(
		inner,
		ComposeOption[Outer, int](defaultInner)(value),
	)

	// Setting None should remove the inner entirely
	outerWithInner := Outer{inner: &Inner{Value: 42, Foo: "bar"}}
	result := lens.Set(O.None[int]())(outerWithInner)
	assert.Nil(t, result.inner)

	// Getting from empty should return None
	emptyOuter := Outer{}
	assert.Equal(t, O.None[int](), lens.Get(emptyOuter))
}

func TestFromPredicateEdgeCases(t *testing.T) {
	type Score struct {
		points int
	}

	pointsLens := L.MakeLens(
		func(s Score) int { return s.points },
		func(s Score, p int) Score { s.points = p; return s },
	)

	// Only positive scores are valid
	validLens := FromPredicate[Score](func(p int) bool { return p > 0 }, 0)(pointsLens)

	// Test with valid score
	validScore := Score{points: 100}
	assert.Equal(t, O.Some(100), validLens.Get(validScore))

	// Test with invalid score (zero)
	zeroScore := Score{points: 0}
	assert.Equal(t, O.None[int](), validLens.Get(zeroScore))

	// Test with invalid score (negative)
	negativeScore := Score{points: -10}
	assert.Equal(t, O.None[int](), validLens.Get(negativeScore))

	// Test setting None sets the nil value
	result := validLens.Set(O.None[int]())(validScore)
	assert.Equal(t, 0, result.points)

	// Test setting Some sets the value
	result2 := validLens.Set(O.Some(50))(zeroScore)
	assert.Equal(t, 50, result2.points)
}

func TestFromNullablePropEdgeCases(t *testing.T) {
	type Container struct {
		item *string
	}

	itemLens := L.MakeLens(
		func(c Container) *string { return c.item },
		func(c Container, i *string) Container { c.item = i; return c },
	)

	defaultItem := "default"
	safeLens := FromNullableProp[Container](O.FromNillable[string], &defaultItem)(itemLens)

	// Test with nil - should return default
	emptyContainer := Container{item: nil}
	assert.Equal(t, &defaultItem, safeLens.Get(emptyContainer))

	// Test with value - should return the value
	value := "actual"
	containerWithItem := Container{item: &value}
	assert.Equal(t, &value, safeLens.Get(containerWithItem))

	// Test setter
	newValue := "new"
	updated := safeLens.Set(&newValue)(emptyContainer)
	assert.Equal(t, &newValue, updated.item)
}

// Lens Law Tests for LensO types

func TestFromNillableLensLaws(t *testing.T) {
	type Config struct {
		timeout *int
	}

	timeoutLens := L.MakeLens(
		func(c Config) *int { return c.timeout },
		func(c Config, t *int) Config { c.timeout = t; return c },
	)

	optLens := FromNillable(timeoutLens)

	// Equality predicates
	eqInt := EQT.Eq[*int]()
	eqOptInt := O.Eq(eqInt)
	eqConfig := func(a, b Config) bool {
		if a.timeout == nil && b.timeout == nil {
			return true
		}
		if a.timeout == nil || b.timeout == nil {
			return false
		}
		return *a.timeout == *b.timeout
	}

	// Test structures
	timeout30 := 30
	timeout60 := 60
	configNil := Config{timeout: nil}
	config30 := Config{timeout: &timeout30}

	// Law 1: get(set(a)(s)) = a
	t.Run("GetSet", func(t *testing.T) {
		// Setting Some and getting back
		result := optLens.Get(optLens.Set(O.Some(&timeout60))(config30))
		assert.True(t, eqOptInt.Equals(result, O.Some(&timeout60)))

		// Setting None and getting back
		result2 := optLens.Get(optLens.Set(O.None[*int]())(config30))
		assert.True(t, eqOptInt.Equals(result2, O.None[*int]()))
	})

	// Law 2: set(get(s))(s) = s
	t.Run("SetGet", func(t *testing.T) {
		// With Some value
		result := optLens.Set(optLens.Get(config30))(config30)
		assert.True(t, eqConfig(result, config30))

		// With None value
		result2 := optLens.Set(optLens.Get(configNil))(configNil)
		assert.True(t, eqConfig(result2, configNil))
	})

	// Law 3: set(a)(set(a)(s)) = set(a)(s)
	t.Run("SetSet", func(t *testing.T) {
		// Setting Some twice
		once := optLens.Set(O.Some(&timeout60))(config30)
		twice := optLens.Set(O.Some(&timeout60))(once)
		assert.True(t, eqConfig(once, twice))

		// Setting None twice
		once2 := optLens.Set(O.None[*int]())(config30)
		twice2 := optLens.Set(O.None[*int]())(once2)
		assert.True(t, eqConfig(once2, twice2))
	})
}

func TestFromNillableRefLensLaws(t *testing.T) {
	type Settings struct {
		maxRetries *int
	}

	retriesLens := L.MakeLensRef(
		func(s *Settings) *int { return s.maxRetries },
		func(s *Settings, r *int) *Settings { s.maxRetries = r; return s },
	)

	optLens := FromNillableRef(retriesLens)

	// Equality predicates
	eqInt := EQT.Eq[*int]()
	eqOptInt := O.Eq(eqInt)
	eqSettings := func(a, b *Settings) bool {
		if a == nil && b == nil {
			return true
		}
		if a == nil || b == nil {
			return false
		}
		if a.maxRetries == nil && b.maxRetries == nil {
			return true
		}
		if a.maxRetries == nil || b.maxRetries == nil {
			return false
		}
		return *a.maxRetries == *b.maxRetries
	}

	// Test structures
	retries3 := 3
	retries5 := 5
	settingsNil := &Settings{maxRetries: nil}
	settings3 := &Settings{maxRetries: &retries3}

	// Law 1: get(set(a)(s)) = a
	t.Run("GetSet", func(t *testing.T) {
		result := optLens.Get(optLens.Set(O.Some(&retries5))(settings3))
		assert.True(t, eqOptInt.Equals(result, O.Some(&retries5)))

		result2 := optLens.Get(optLens.Set(O.None[*int]())(settings3))
		assert.True(t, eqOptInt.Equals(result2, O.None[*int]()))
	})

	// Law 2: set(get(s))(s) = s
	t.Run("SetGet", func(t *testing.T) {
		result := optLens.Set(optLens.Get(settings3))(settings3)
		assert.True(t, eqSettings(result, settings3))

		result2 := optLens.Set(optLens.Get(settingsNil))(settingsNil)
		assert.True(t, eqSettings(result2, settingsNil))
	})

	// Law 3: set(a)(set(a)(s)) = set(a)(s)
	t.Run("SetSet", func(t *testing.T) {
		once := optLens.Set(O.Some(&retries5))(settings3)
		twice := optLens.Set(O.Some(&retries5))(once)
		assert.True(t, eqSettings(once, twice))

		once2 := optLens.Set(O.None[*int]())(settings3)
		twice2 := optLens.Set(O.None[*int]())(once2)
		assert.True(t, eqSettings(once2, twice2))
	})
}

func TestComposeOptionLensLaws(t *testing.T) {
	defaultInner := &Inner{Value: 0, Foo: "default"}

	value := L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	inner := FromNillable(L.MakeLens(Outer.GetInner, Outer.SetInner))
	lens := F.Pipe1(inner, ComposeOption[Outer, int](defaultInner)(value))

	// Equality predicates
	eqInt := EQT.Eq[int]()
	eqOptInt := O.Eq(eqInt)
	eqOuter := func(a, b Outer) bool {
		if a.inner == nil && b.inner == nil {
			return true
		}
		if a.inner == nil || b.inner == nil {
			return false
		}
		return a.inner.Value == b.inner.Value && a.inner.Foo == b.inner.Foo
	}

	// Test structures
	outerNil := Outer{inner: nil}
	outer42 := Outer{inner: &Inner{Value: 42, Foo: "test"}}

	// Law 1: get(set(a)(s)) = a
	t.Run("GetSet", func(t *testing.T) {
		result := lens.Get(lens.Set(O.Some(100))(outer42))
		assert.True(t, eqOptInt.Equals(result, O.Some(100)))

		result2 := lens.Get(lens.Set(O.None[int]())(outer42))
		assert.True(t, eqOptInt.Equals(result2, O.None[int]()))
	})

	// Law 2: set(get(s))(s) = s
	t.Run("SetGet", func(t *testing.T) {
		result := lens.Set(lens.Get(outer42))(outer42)
		assert.True(t, eqOuter(result, outer42))

		result2 := lens.Set(lens.Get(outerNil))(outerNil)
		assert.True(t, eqOuter(result2, outerNil))
	})

	// Law 3: set(a)(set(a)(s)) = set(a)(s)
	t.Run("SetSet", func(t *testing.T) {
		once := lens.Set(O.Some(100))(outer42)
		twice := lens.Set(O.Some(100))(once)
		assert.True(t, eqOuter(once, twice))

		once2 := lens.Set(O.None[int]())(outer42)
		twice2 := lens.Set(O.None[int]())(once2)
		assert.True(t, eqOuter(once2, twice2))
	})
}

func TestComposeOptionsLensLaws(t *testing.T) {
	defaultValue := 1
	defaultFoo := "default"
	defaultInner := &InnerOpt{Value: &defaultValue, Foo: &defaultFoo}

	value := FromNillable(L.MakeLensRef((*InnerOpt).GetValue, (*InnerOpt).SetValue))
	inner := FromNillable(L.MakeLens(OuterOpt.GetInnerOpt, OuterOpt.SetInnerOpt))
	lens := F.Pipe1(inner, Compose[OuterOpt, *int](defaultInner)(value))

	// Equality predicates
	eqIntPtr := EQT.Eq[*int]()
	eqOptIntPtr := O.Eq(eqIntPtr)
	eqOuterOpt := func(a, b OuterOpt) bool {
		if a.inner == nil && b.inner == nil {
			return true
		}
		if a.inner == nil || b.inner == nil {
			return false
		}
		aVal := a.inner.Value
		bVal := b.inner.Value
		if aVal == nil && bVal == nil {
			return true
		}
		if aVal == nil || bVal == nil {
			return false
		}
		return *aVal == *bVal
	}

	// Test structures
	val42 := 42
	val100 := 100
	outerNil := OuterOpt{inner: nil}
	outer42 := OuterOpt{inner: &InnerOpt{Value: &val42, Foo: &defaultFoo}}

	// Law 1: get(set(a)(s)) = a
	t.Run("GetSet", func(t *testing.T) {
		result := lens.Get(lens.Set(O.Some(&val100))(outer42))
		assert.True(t, eqOptIntPtr.Equals(result, O.Some(&val100)))

		result2 := lens.Get(lens.Set(O.None[*int]())(outer42))
		assert.True(t, eqOptIntPtr.Equals(result2, O.None[*int]()))
	})

	// Law 2: set(get(s))(s) = s
	t.Run("SetGet", func(t *testing.T) {
		result := lens.Set(lens.Get(outer42))(outer42)
		assert.True(t, eqOuterOpt(result, outer42))

		result2 := lens.Set(lens.Get(outerNil))(outerNil)
		assert.True(t, eqOuterOpt(result2, outerNil))
	})

	// Law 3: set(a)(set(a)(s)) = set(a)(s)
	t.Run("SetSet", func(t *testing.T) {
		once := lens.Set(O.Some(&val100))(outer42)
		twice := lens.Set(O.Some(&val100))(once)
		assert.True(t, eqOuterOpt(once, twice))

		once2 := lens.Set(O.None[*int]())(outer42)
		twice2 := lens.Set(O.None[*int]())(once2)
		assert.True(t, eqOuterOpt(once2, twice2))
	})
}

func TestFromPredicateLensLaws(t *testing.T) {
	type Score struct {
		points int
	}

	pointsLens := L.MakeLens(
		func(s Score) int { return s.points },
		func(s Score, p int) Score { s.points = p; return s },
	)

	// Only positive scores are valid
	validLens := FromPredicate[Score](func(p int) bool { return p > 0 }, 0)(pointsLens)

	// Equality predicates
	eqInt := EQT.Eq[int]()
	eqOptInt := O.Eq(eqInt)
	eqScore := func(a, b Score) bool { return a.points == b.points }

	// Test structures
	scoreZero := Score{points: 0}
	score100 := Score{points: 100}

	// Law 1: get(set(a)(s)) = a
	t.Run("GetSet", func(t *testing.T) {
		result := validLens.Get(validLens.Set(O.Some(50))(score100))
		assert.True(t, eqOptInt.Equals(result, O.Some(50)))

		result2 := validLens.Get(validLens.Set(O.None[int]())(score100))
		assert.True(t, eqOptInt.Equals(result2, O.None[int]()))
	})

	// Law 2: set(get(s))(s) = s
	t.Run("SetGet", func(t *testing.T) {
		result := validLens.Set(validLens.Get(score100))(score100)
		assert.True(t, eqScore(result, score100))

		result2 := validLens.Set(validLens.Get(scoreZero))(scoreZero)
		assert.True(t, eqScore(result2, scoreZero))
	})

	// Law 3: set(a)(set(a)(s)) = set(a)(s)
	t.Run("SetSet", func(t *testing.T) {
		once := validLens.Set(O.Some(75))(score100)
		twice := validLens.Set(O.Some(75))(once)
		assert.True(t, eqScore(once, twice))

		once2 := validLens.Set(O.None[int]())(score100)
		twice2 := validLens.Set(O.None[int]())(once2)
		assert.True(t, eqScore(once2, twice2))
	})
}
