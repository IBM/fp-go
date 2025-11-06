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

	RA "github.com/IBM/fp-go/v2/array"
	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
)

// AssertLaws asserts the apply monad laws for the array monad
func AssertLaws[A, B, C any](t *testing.T,
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		RA.Eq(eqa),
		RA.Eq(eqb),
		RA.Eq(eqc),

		RA.Of[A],
		RA.Of[B],
		RA.Of[C],

		RA.Of[func(A) A],
		RA.Of[func(A) B],
		RA.Of[func(B) C],
		RA.Of[func(func(A) B) B],

		RA.MonadMap[A, A],
		RA.MonadMap[A, B],
		RA.MonadMap[A, C],
		RA.MonadMap[B, C],

		RA.MonadMap[func(B) C, func(func(A) B) func(A) C],

		RA.MonadChain[A, A],
		RA.MonadChain[A, B],
		RA.MonadChain[A, C],
		RA.MonadChain[B, C],

		RA.MonadAp[A, A],
		RA.MonadAp[B, A],
		RA.MonadAp[C, B],
		RA.MonadAp[C, A],

		RA.MonadAp[B, func(A) B],
		RA.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
