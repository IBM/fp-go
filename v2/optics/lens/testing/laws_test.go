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
	L "github.com/IBM/fp-go/v2/optics/common"
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
	valueLens  = L.MakeLensRef((*Inner).GetValue, (*Inner).SetValue)

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

	laws := AssertLaws(
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

	laws := AssertLaws(
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

	streetName := L.LensComposeLens[*Address](streetLens)(addrLens)

	laws := AssertLaws(
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
