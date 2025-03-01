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

package either

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid/testing"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestApplySemigroup(t *testing.T) {

	sg := ApplySemigroup[string](N.SemigroupSum[int]())

	la := Left[int]("a")
	lb := Left[int]("b")
	r1 := Right[string](1)
	r2 := Right[string](2)
	r3 := Right[string](3)

	assert.Equal(t, la, sg.Concat(la, lb))
	assert.Equal(t, lb, sg.Concat(r1, lb))
	assert.Equal(t, la, sg.Concat(la, r2))
	assert.Equal(t, lb, sg.Concat(r1, lb))
	assert.Equal(t, r3, sg.Concat(r1, r2))
}

func TestApplicativeMonoid(t *testing.T) {

	m := ApplicativeMonoid[string](N.MonoidSum[int]())

	la := Left[int]("a")
	lb := Left[int]("b")
	r1 := Right[string](1)
	r2 := Right[string](2)
	r3 := Right[string](3)

	assert.Equal(t, la, m.Concat(la, m.Empty()))
	assert.Equal(t, lb, m.Concat(m.Empty(), lb))
	assert.Equal(t, r1, m.Concat(r1, m.Empty()))
	assert.Equal(t, r2, m.Concat(m.Empty(), r2))
	assert.Equal(t, r3, m.Concat(r1, r2))
}

func TestApplicativeMonoidLaws(t *testing.T) {
	m := ApplicativeMonoid[string](N.MonoidSum[int]())
	M.AssertLaws(t, m)([]Either[string, int]{Left[int]("a"), Right[string](1)})
}
