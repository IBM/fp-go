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

	ET "github.com/IBM/fp-go/either"
	EQ "github.com/IBM/fp-go/eq"
	L "github.com/IBM/fp-go/internal/monad/testing"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[E, A, B, C any](t *testing.T,
	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		ET.Eq(eqe, eqa),
		ET.Eq(eqe, eqb),
		ET.Eq(eqe, eqc),

		ET.Of[E, A],
		ET.Of[E, B],
		ET.Of[E, C],

		ET.Of[E, func(A) A],
		ET.Of[E, func(A) B],
		ET.Of[E, func(B) C],
		ET.Of[E, func(func(A) B) B],

		ET.MonadMap[E, A, A],
		ET.MonadMap[E, A, B],
		ET.MonadMap[E, A, C],
		ET.MonadMap[E, B, C],

		ET.MonadMap[E, func(B) C, func(func(A) B) func(A) C],

		ET.MonadChain[E, A, A],
		ET.MonadChain[E, A, B],
		ET.MonadChain[E, A, C],
		ET.MonadChain[E, B, C],

		ET.MonadAp[A, E, A],
		ET.MonadAp[B, E, A],
		ET.MonadAp[C, E, B],
		ET.MonadAp[C, E, A],

		ET.MonadAp[B, E, func(A) B],
		ET.MonadAp[func(A) C, E, func(A) B],

		ab,
		bc,
	)

}
