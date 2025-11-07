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
