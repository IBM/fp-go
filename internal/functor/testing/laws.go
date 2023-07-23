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

package testing

import (
	"testing"

	E "github.com/IBM/fp-go/eq"
	F "github.com/IBM/fp-go/function"
	"github.com/stretchr/testify/assert"
)

// Functor identity law
//
// F.map(fa, a => a) <-> fa
func AssertIdentity[HKTA, A any](t *testing.T, eq E.Eq[HKTA], fmap func(HKTA, func(A) A) HKTA) func(fa HKTA) bool {
	return func(fa HKTA) bool {
		return assert.True(t, eq.Equals(fa, fmap(fa, F.Identity[A])), "Functor identity law")
	}
}

// Functor composition law
//
// F.map(fa, a => bc(ab(a))) <-> F.map(F.map(fa, ab), bc)
func AssertComposition[HKTA, HKTB, HKTC, A, B, C any](
	t *testing.T,

	eq E.Eq[HKTC],

	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,
	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	return func(fa HKTA) bool {
		return assert.True(t, eq.Equals(fac(fa, F.Flow2(ab, bc)), fbc(fab(fa, ab), bc)), "Functor composition law")
	}
}

// AssertLaws asserts the functor laws `identity` and `composition`
func AssertLaws[HKTA, HKTB, HKTC, A, B, C any](t *testing.T,
	eqa E.Eq[HKTA],
	eqc E.Eq[HKTC],

	faa func(HKTA, func(A) A) HKTA,
	fab func(HKTA, func(A) B) HKTB,
	fac func(HKTA, func(A) C) HKTC,
	fbc func(HKTB, func(B) C) HKTC,
	ab func(A) B,
	bc func(B) C,
) func(fa HKTA) bool {
	identity := AssertIdentity(t, eqa, faa)
	composition := AssertComposition(t, eqc, fab, fac, fbc, ab, bc)

	return func(fa HKTA) bool {
		return identity(fa) && composition(fa)
	}
}
