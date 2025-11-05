// Copyright (c) 2023 IBM Corp.
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
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type (
	Street struct {
		num  int
		name string
	}

	Address struct {
		city   string
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
	streetLens = MakeLensRef((*Street).GetName, (*Street).SetName)
	addrLens   = MakeLensRef((*Address).GetStreet, (*Address).SetStreet)

	sampleStreet  = Street{num: 220, name: "Schönaicherstr"}
	sampleAddress = Address{city: "Böblingen", street: &sampleStreet}
)

func TestLens(t *testing.T) {
	// read the value
	assert.Equal(t, sampleStreet.name, streetLens.Get(&sampleStreet))
	// new street
	newName := "Böblingerstr"
	// update
	old := sampleStreet
	updated := streetLens.Set(newName)(&sampleStreet)
	assert.Equal(t, old, sampleStreet)
	// validate the new name
	assert.Equal(t, newName, streetLens.Get(updated))
}

func TestAddressCompose(t *testing.T) {
	// compose
	streetName := Compose[*Address](streetLens)(addrLens)
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))
	// new street
	newName := "Böblingerstr"
	updated := streetName.Set(newName)(&sampleAddress)
	// check that we have not modified the original
	assert.Equal(t, sampleStreet.name, streetName.Get(&sampleAddress))
	assert.Equal(t, newName, streetName.Get(updated))
}

func TestIMap(t *testing.T) {

	type S struct {
		a int
	}

	sa := F.Pipe1(
		Id[S](),
		IMap[S](
			func(s S) int { return s.a },
			func(a int) S { return S{a} },
		),
	)

	assert.Equal(t, 1, sa.Get(S{1}))
	assert.Equal(t, S{2}, sa.Set(2)(S{1}))
}

func TestPassByValue(t *testing.T) {

	testLens := MakeLens(func(s Street) string { return s.name }, func(s Street, value string) Street {
		s.name = value
		return s
	})

	s1 := Street{1, "value1"}
	s2 := testLens.Set("value2")(s1)

	assert.Equal(t, "value1", s1.name)
	assert.Equal(t, "value2", s2.name)
}

func TestFromNullableProp(t *testing.T) {
	// default inner object
	defaultInner := &Inner{
		Value: 0,
		Foo:   "foo",
	}
	// access to the value
	value := MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	// access to inner
	inner := FromNullableProp[Outer](O.FromNillable[Inner], defaultInner)(MakeLens(Outer.GetInner, Outer.SetInner))
	// compose
	lens := F.Pipe1(
		inner,
		Compose[Outer](value),
	)
	outer1 := Outer{inner: &Inner{Value: 1, Foo: "a"}}
	// the checks
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(1)(Outer{}))
	assert.Equal(t, 0, lens.Get(Outer{}))
	assert.Equal(t, Outer{inner: &Inner{Value: 1, Foo: "foo"}}, lens.Set(1)(Outer{inner: &Inner{Value: 2, Foo: "foo"}}))
	assert.Equal(t, 1, lens.Get(Outer{inner: &Inner{Value: 1, Foo: "foo"}}))
	assert.Equal(t, outer1, Modify[Outer](F.Identity[int])(lens)(outer1))
}

func TestComposeOption(t *testing.T) {
	// default inner object
	defaultInner := &Inner{
		Value: 0,
		Foo:   "foo",
	}
	// access to the value
	value := MakeLensRef((*Inner).GetValue, (*Inner).SetValue)
	// access to inner
	inner := FromNillable(MakeLens(Outer.GetInner, Outer.SetInner))
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
	assert.Equal(t, outer1, Modify[Outer](F.Identity[O.Option[int]])(lens)(outer1))
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
	value := FromNillable(MakeLensRef((*InnerOpt).GetValue, (*InnerOpt).SetValue))
	// access to inner
	inner := FromNillable(MakeLens(OuterOpt.GetInnerOpt, OuterOpt.SetInnerOpt))
	// compose lenses
	lens := F.Pipe1(
		inner,
		ComposeOptions[OuterOpt, *int](defaultInner)(value),
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
	assert.Equal(t, outer1, Modify[OuterOpt](F.Identity[O.Option[*int]])(lens)(outer1))
}

func TestIdRef(t *testing.T) {
	idLens := IdRef[Street]()
	street := &Street{num: 1, name: "Main"}

	assert.Equal(t, street, idLens.Get(street))

	newStreet := &Street{num: 2, name: "Oak"}
	result := idLens.Set(newStreet)(street)
	assert.Equal(t, newStreet, result)
	assert.Equal(t, 1, street.num) // Original unchanged
}

func TestComposeRef(t *testing.T) {
	composedLens := ComposeRef[Address, *Street](streetLens)(addrLens)

	assert.Equal(t, sampleStreet.name, composedLens.Get(&sampleAddress))

	newName := "NewStreet"
	updated := composedLens.Set(newName)(&sampleAddress)
	assert.Equal(t, newName, composedLens.Get(updated))
	assert.Equal(t, sampleStreet.name, sampleAddress.street.name) // Original unchanged
}

func TestFromPredicateRef(t *testing.T) {
	type Person struct {
		age int
	}

	ageLens := MakeLensRef(
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

	timeoutLens := MakeLensRef(
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

	timeoutLens := MakeLensRef(
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
		retries O.Option[int]
	}

	retriesLens := MakeLensRef(
		func(s *Settings) O.Option[int] { return s.retries },
		func(s *Settings, r O.Option[int]) *Settings {
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

func TestMakeLensCurried(t *testing.T) {
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
	assert.Equal(t, "Main", nameLens.Get(street))

	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name)
}

func TestMakeLensRefCurried(t *testing.T) {
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
	assert.Equal(t, "Main", nameLens.Get(street))

	updated := nameLens.Set("Oak")(street)
	assert.Equal(t, "Oak", updated.name)
	assert.Equal(t, "Main", street.name)
}
