// Copyright (c) 2024 IBM Corp.
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
	P "github.com/IBM/fp-go/v2/pair"

	M "github.com/IBM/fp-go/v2/monoid"
)

// AssertLaws asserts the apply monad laws for the [P.Pair] monad
func assertLawsHead[E, A, B, C any](t *testing.T,
	m M.Monoid[E],

	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	fofc := P.Pointed[C](m)
	fofaa := P.Pointed[func(A) A](m)
	fofbc := P.Pointed[func(B) C](m)
	fofabb := P.Pointed[func(func(A) B) B](m)

	fmap := P.Functor[func(B) C, E, func(func(A) B) func(A) C]()

	fapabb := P.Applicative[func(A) B, E, B](m)
	fapabac := P.Applicative[func(A) B, E, func(A) C](m)

	maa := P.Monad[A, E, A](m)
	mab := P.Monad[A, E, B](m)
	mac := P.Monad[A, E, C](m)
	mbc := P.Monad[B, E, C](m)

	return L.MonadAssertLaws(t,
		P.Eq(eqa, eqe),
		P.Eq(eqb, eqe),
		P.Eq(eqc, eqe),

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

// AssertLaws asserts the apply monad laws for the [P.Pair] monad
func assertLawsTail[E, A, B, C any](t *testing.T,
	m M.Monoid[E],

	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	fofc := P.PointedTail[C](m)
	fofaa := P.PointedTail[func(A) A](m)
	fofbc := P.PointedTail[func(B) C](m)
	fofabb := P.PointedTail[func(func(A) B) B](m)

	fmap := P.FunctorTail[func(B) C, E, func(func(A) B) func(A) C]()

	fapabb := P.ApplicativeTail[func(A) B, E, B](m)
	fapabac := P.ApplicativeTail[func(A) B, E, func(A) C](m)

	maa := P.MonadTail[A, E, A](m)
	mab := P.MonadTail[A, E, B](m)
	mac := P.MonadTail[A, E, C](m)
	mbc := P.MonadTail[B, E, C](m)

	return L.MonadAssertLaws(t,
		P.Eq(eqe, eqa),
		P.Eq(eqe, eqb),
		P.Eq(eqe, eqc),

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

// AssertLaws asserts the apply monad laws for the [P.Pair] monad
func AssertLaws[E, A, B, C any](t *testing.T,
	m M.Monoid[E],

	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(A) bool {

	head := assertLawsHead(t, m, eqe, eqa, eqb, eqc, ab, bc)
	tail := assertLawsHead(t, m, eqe, eqa, eqb, eqc, ab, bc)

	return func(a A) bool {
		return head(a) && tail(a)
	}
}
