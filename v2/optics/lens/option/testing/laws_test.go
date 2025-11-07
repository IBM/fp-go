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

package testing

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	L "github.com/IBM/fp-go/v2/optics/lens"
	LI "github.com/IBM/fp-go/v2/optics/lens/iso"
	LO "github.com/IBM/fp-go/v2/optics/lens/option"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
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

func (outer *OuterOpt) GetInner() *InnerOpt {
	return outer.inner
}

func (outer *OuterOpt) SetInner(inner *InnerOpt) *OuterOpt {
	outer.inner = inner
	return outer
}

func (inner *InnerOpt) GetValue() *int {
	return inner.Value
}

func (inner *InnerOpt) SetValue(value *int) *InnerOpt {
	inner.Value = value
	return inner
}

func (outer *Outer) GetInner() *Inner {
	return outer.inner
}

func (outer *Outer) SetInner(inner *Inner) *Outer {
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
	outerLens  = LO.FromNillableRef(L.MakeLensRef((*Outer).GetInner, (*Outer).SetInner))
	valueLens  = L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)

	outerOptLens = LO.FromNillableRef(L.MakeLensRef((*OuterOpt).GetInner, (*OuterOpt).SetInner))
	valueOptLens = L.MakeLensRef((*InnerOpt).GetValue, (*InnerOpt).SetValue)

	sampleStreet  = Street{num: 220, name: "Schönaicherstr"}
	sampleAddress = Address{city: "Böblingen", street: &sampleStreet}
	sampleStreet2 = Street{num: 220, name: "Neue Str"}

	defaultInner = Inner{
		Value: -1,
		Foo:   "foo",
	}

	emptyOuter = Outer{}

	defaultInnerOpt = InnerOpt{
		Value: &defaultInner.Value,
		Foo:   &defaultInner.Foo,
	}

	emptyOuterOpt = OuterOpt{}
)

func TestStreetLensLaws(t *testing.T) {
	// some comparison
	eqs := EQT.Eq[*Street]()
	eqa := EQT.Eq[string]()

	laws := LT.AssertLaws(
		t,
		eqa,
		eqs,
	)(streetLens)

	cpy := sampleStreet
	assert.True(t, laws(&sampleStreet, "Neue Str."))
	assert.Equal(t, cpy, sampleStreet)
}

func TestAddrLensLaws(t *testing.T) {
	// some comparison
	eqs := EQT.Eq[*Address]()
	eqa := EQT.Eq[*Street]()

	laws := LT.AssertLaws(
		t,
		eqa,
		eqs,
	)(addrLens)

	cpyAddr := sampleAddress
	cpyStreet := sampleStreet2
	assert.True(t, laws(&sampleAddress, &sampleStreet2))
	assert.Equal(t, cpyAddr, sampleAddress)
	assert.Equal(t, cpyStreet, sampleStreet2)
}

func TestCompose(t *testing.T) {
	// some comparison
	eqs := EQT.Eq[*Address]()
	eqa := EQT.Eq[string]()

	streetName := L.Compose[*Address](streetLens)(addrLens)

	laws := LT.AssertLaws(
		t,
		eqa,
		eqs,
	)(streetName)

	cpyAddr := sampleAddress
	cpyStreet := sampleStreet
	assert.True(t, laws(&sampleAddress, "Neue Str."))
	assert.Equal(t, cpyAddr, sampleAddress)
	assert.Equal(t, cpyStreet, sampleStreet)
}

func TestOuterLensLaws(t *testing.T) {
	// some equal predicates
	eqValue := EQT.Eq[int]()
	eqOptValue := O.Eq(eqValue)
	// lens to access a value from outer
	valueFromOuter := LO.ComposeOption[*Outer, int](&defaultInner)(valueLens)(outerLens)
	// try to access the value, this should get an option
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(&emptyOuter), O.None[int]()))
	// update the object
	withValue := valueFromOuter.Set(O.Some(1))(&emptyOuter)
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(&emptyOuter), O.None[int]()))
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(withValue), O.Some(1)))
	// updating with none should remove the inner
	nextValue := valueFromOuter.Set(O.None[int]())(withValue)
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(nextValue), O.None[int]()))
	// check if this meets the laws

	eqOuter := EQT.Eq[*Outer]()

	laws := LT.AssertLaws(
		t,
		eqOptValue,
		eqOuter,
	)(valueFromOuter)

	assert.True(t, laws(&emptyOuter, O.Some(2)))
	assert.True(t, laws(&emptyOuter, O.None[int]()))

	assert.True(t, laws(withValue, O.Some(2)))
	assert.True(t, laws(withValue, O.None[int]()))
}

func TestOuterOptLensLaws(t *testing.T) {
	// some equal predicates
	eqValue := EQT.Eq[int]()
	eqOptValue := O.Eq(eqValue)
	intIso := LI.FromNillable[int]()
	// lens to access a value from outer
	valueFromOuter := F.Pipe3(
		valueOptLens,
		LI.Compose[*InnerOpt](intIso),
		LO.Compose[*OuterOpt, int](&defaultInnerOpt),
		I.Ap[L.Lens[*OuterOpt, O.Option[int]]](outerOptLens),
	)

	// try to access the value, this should get an option
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(&emptyOuterOpt), O.None[int]()))
	// update the object
	withValue := valueFromOuter.Set(O.Some(1))(&emptyOuterOpt)
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(&emptyOuterOpt), O.None[int]()))
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(withValue), O.Some(1)))
	// updating with none should remove the inner
	nextValue := valueFromOuter.Set(O.None[int]())(withValue)
	assert.True(t, eqOptValue.Equals(valueFromOuter.Get(nextValue), O.None[int]()))
	// check if this meets the laws

	eqOuter := EQT.Eq[*OuterOpt]()

	laws := LT.AssertLaws(
		t,
		eqOptValue,
		eqOuter,
	)(valueFromOuter)

	assert.True(t, laws(&emptyOuterOpt, O.Some(2)))
	assert.True(t, laws(&emptyOuterOpt, O.None[int]()))

	assert.True(t, laws(withValue, O.Some(2)))
	assert.True(t, laws(withValue, O.None[int]()))
}
