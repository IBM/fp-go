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

	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
	ST "github.com/IBM/fp-go/v2/state"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[S, A, B, C any](t *testing.T,
	eqw EQ.Eq[S],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,

	s S,
) func(a A) bool {

	fofc := ST.Pointed[S, C]()
	fofaa := ST.Pointed[S, func(A) A]()
	fofbc := ST.Pointed[S, func(B) C]()
	fofabb := ST.Pointed[S, func(func(A) B) B]()

	fmap := ST.Functor[S, func(B) C, func(func(A) B) func(A) C]()

	fapabb := ST.Applicative[S, func(A) B, B]()
	fapabac := ST.Applicative[S, func(A) B, func(A) C]()

	maa := ST.Monad[S, A, A]()
	mab := ST.Monad[S, A, B]()
	mac := ST.Monad[S, A, C]()
	mbc := ST.Monad[S, B, C]()

	return L.MonadAssertLaws(t,
		ST.Eq(eqw, eqa)(s),
		ST.Eq(eqw, eqb)(s),
		ST.Eq(eqw, eqc)(s),

		fofc,
		fofaa,
		fofbc,
		fofabb,

		fmap,

		fapabb,
		fapabac,

		maa,
		mab,
		mac,
		mbc,

		ab,
		bc,
	)

}
