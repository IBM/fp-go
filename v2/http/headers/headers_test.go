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

package headers

import (
	"net/http"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/eq"
	LT "github.com/IBM/fp-go/v2/optics/lens/testing"
	O "github.com/IBM/fp-go/v2/option"
	RG "github.com/IBM/fp-go/v2/record/generic"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

var (
	sEq      = eq.FromEquals(S.Eq)
	valuesEq = RG.Eq[http.Header](A.Eq(sEq))
)

func TestLaws(t *testing.T) {
	name := "Content-Type"
	fieldLaws := LT.AssertLaws[http.Header, O.Option[string]](t, O.Eq(sEq), valuesEq)(AtValue(name))

	n := O.None[string]()
	s1 := O.Some("s1")

	def := make(http.Header)

	v1 := make(http.Header)
	v1.Set(name, "v1")

	v2 := make(http.Header)
	v2.Set("Other-Header", "v2")

	assert.True(t, fieldLaws(def, n))
	assert.True(t, fieldLaws(v1, n))
	assert.True(t, fieldLaws(v2, n))

	assert.True(t, fieldLaws(def, s1))
	assert.True(t, fieldLaws(v1, s1))
	assert.True(t, fieldLaws(v2, s1))
}
