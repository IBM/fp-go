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

package array

import (
	"testing"

	A "github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/eq"
	LT "github.com/IBM/fp-go/optics/lens/testing"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	"github.com/stretchr/testify/assert"
)

var (
	sEq = eq.FromEquals(S.Eq)
)

func TestLaws(t *testing.T) {
	headLaws := LT.AssertLaws(t, O.Eq(sEq), A.Eq(sEq))(AtHead[string]())

	assert.True(t, headLaws(A.Empty[string](), O.None[string]()))
	assert.True(t, headLaws(A.Empty[string](), O.Of("a")))
	assert.True(t, headLaws(A.From("a", "b"), O.None[string]()))
	assert.True(t, headLaws(A.From("a", "b"), O.Of("c")))
}
