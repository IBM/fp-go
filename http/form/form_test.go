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

package form

import (
	"net/url"
	"testing"

	A "github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	LT "github.com/IBM/fp-go/optics/lens/testing"
	O "github.com/IBM/fp-go/option"
	RG "github.com/IBM/fp-go/record/generic"
	S "github.com/IBM/fp-go/string"
	"github.com/stretchr/testify/assert"
)

var (
	sEq      = eq.FromEquals(S.Eq)
	valuesEq = RG.Eq[url.Values](A.Eq(sEq))
)

func TestLaws(t *testing.T) {
	name := "Content-Type"
	fieldLaws := LT.AssertLaws[url.Values, O.Option[string]](t, O.Eq(sEq), valuesEq)(AtValue(name))

	n := O.None[string]()
	s1 := O.Some("s1")

	v1 := F.Pipe1(
		Default,
		WithValue(name)("v1"),
	)

	v2 := F.Pipe1(
		Default,
		WithValue("Other-Header")("v2"),
	)

	assert.True(t, fieldLaws(Default, n))
	assert.True(t, fieldLaws(v1, n))
	assert.True(t, fieldLaws(v2, n))

	assert.True(t, fieldLaws(Default, s1))
	assert.True(t, fieldLaws(v1, s1))
	assert.True(t, fieldLaws(v2, s1))
}

func TestFormField(t *testing.T) {

	v1 := F.Pipe1(
		Default,
		WithValue("h1")("v1"),
	)

	v2 := F.Pipe1(
		v1,
		WithValue("h2")("v2"),
	)

	// make sure the code does not change structures
	assert.False(t, valuesEq.Equals(Default, v1))
	assert.False(t, valuesEq.Equals(Default, v2))
	assert.False(t, valuesEq.Equals(v1, v2))

	// check for existence of values
	assert.Equal(t, "v1", v1.Get("h1"))
	assert.Equal(t, "v1", v2.Get("h1"))
	assert.Equal(t, "v2", v2.Get("h2"))

	// check getter on lens

	l1 := AtValue("h1")
	l2 := AtValue("h2")

	assert.Equal(t, O.Of("v1"), l1.Get(v1))
	assert.Equal(t, O.Of("v1"), l1.Get(v2))
	assert.Equal(t, O.Of("v2"), l2.Get(v2))
}
