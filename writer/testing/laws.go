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

	EQ "github.com/IBM/fp-go/eq"
	L "github.com/IBM/fp-go/internal/monad/testing"
	M "github.com/IBM/fp-go/monoid"
	WRT "github.com/IBM/fp-go/writer"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[W, A, B, C any](t *testing.T,
	m M.Monoid[W],

	eqw EQ.Eq[W],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	fofc := WRT.Pointed[W, C](m)
	fofaa := WRT.Pointed[W, func(A) A](m)
	fofbc := WRT.Pointed[W, func(B) C](m)
	fofabb := WRT.Pointed[W, func(func(A) B) B](m)

	fmap := WRT.Functor[W, func(B) C, func(func(A) B) func(A) C]()

	fapabb := WRT.Applicative[W, func(A) B, B](m)
	fapabac := WRT.Applicative[W, func(A) B, func(A) C](m)

	maa := WRT.Monad[W, A, A](m)
	mab := WRT.Monad[W, A, B](m)
	mac := WRT.Monad[W, A, C](m)
	mbc := WRT.Monad[W, B, C](m)

	return L.MonadAssertLaws(t,
		WRT.Eq(eqw, eqa),
		WRT.Eq(eqw, eqb),
		WRT.Eq(eqw, eqc),

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
