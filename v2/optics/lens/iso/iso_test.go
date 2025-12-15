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

package iso

import (
	"testing"

	EQT "github.com/IBM/fp-go/v2/eq/testing"
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type (
	Inner struct {
		Value *int
		Foo   string
	}

	Outer struct {
		inner Inner
	}
)

func (outer Outer) GetInner() Inner {
	return outer.inner
}

func (outer Outer) SetInner(inner Inner) Outer {
	outer.inner = inner
	return outer
}

func (inner Inner) GetValue() *int {
	return inner.Value
}

func (inner Inner) SetValue(value *int) Inner {
	inner.Value = value
	return inner
}

func TestIso(t *testing.T) {

	eqOptInt := O.Eq(EQT.Eq[int]())
	eqOuter := EQT.Eq[Outer]()

	emptyOuter := Outer{}

	// iso
	intIso := FromNillable[int]()

	innerFromOuter := L.MakeLens(Outer.GetInner, Outer.SetInner)
	valueFromInner := L.MakeLens(Inner.GetValue, Inner.SetValue)

	optValueFromInner := F.Pipe1(
		valueFromInner,
		Compose[Inner](intIso),
	)

	optValueFromOuter := F.Pipe1(
		innerFromOuter,
		L.Compose[Outer](optValueFromInner),
	)

	// try some access
	require.True(t, eqOptInt.Equals(optValueFromOuter.Get(emptyOuter), O.None[int]()))

	updatedOuter := optValueFromOuter.Set(O.Some(1))(emptyOuter)

	require.True(t, eqOptInt.Equals(optValueFromOuter.Get(updatedOuter), O.Some(1)))
	secondOuter := optValueFromOuter.Set(O.None[int]())(updatedOuter)
	require.True(t, eqOptInt.Equals(optValueFromOuter.Get(secondOuter), O.None[int]()))

	// check if this obeys laws
	laws := LT.AssertLaws(
		t,
		eqOptInt,
		eqOuter,
	)(optValueFromOuter)

	assert.True(t, laws(emptyOuter, O.Some(2)))
}
