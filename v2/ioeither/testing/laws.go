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

	"github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
	"github.com/IBM/fp-go/v2/ioeither"
)

// AssertLaws asserts the apply monad laws for the `IOEither` monad
func AssertLaws[E, A, B, C any](t *testing.T,
	eqe EQ.Eq[E],
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		ioeither.Eq(either.Eq(eqe, eqa)),
		ioeither.Eq(either.Eq(eqe, eqb)),
		ioeither.Eq(either.Eq(eqe, eqc)),

		ioeither.Of[E, A],
		ioeither.Of[E, B],
		ioeither.Of[E, C],

		ioeither.Of[E, func(A) A],
		ioeither.Of[E, func(A) B],
		ioeither.Of[E, func(B) C],
		ioeither.Of[E, func(func(A) B) B],

		ioeither.MonadMap[E, A, A],
		ioeither.MonadMap[E, A, B],
		ioeither.MonadMap[E, A, C],
		ioeither.MonadMap[E, B, C],

		ioeither.MonadMap[E, func(B) C, func(func(A) B) func(A) C],

		ioeither.MonadChain[E, A, A],
		ioeither.MonadChain[E, A, B],
		ioeither.MonadChain[E, A, C],
		ioeither.MonadChain[E, B, C],

		ioeither.MonadAp[A, E, A],
		ioeither.MonadAp[B, E, A],
		ioeither.MonadAp[C, E, B],
		ioeither.MonadAp[C, E, A],

		ioeither.MonadAp[B, E, func(A) B],
		ioeither.MonadAp[func(A) C, E, func(A) B],

		ab,
		bc,
	)

}
