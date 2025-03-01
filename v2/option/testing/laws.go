// Copyright (c) 2025 IBM Corp.
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
	O "github.com/IBM/fp-go/v2/option"
)

// AssertLaws asserts the apply monad laws for the [Option] monad
func AssertLaws[A, B, C any](t *testing.T,
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		O.Eq(eqa),
		O.Eq(eqb),
		O.Eq(eqc),

		O.Of[A],
		O.Of[B],
		O.Of[C],

		O.Of[func(A) A],
		O.Of[func(A) B],
		O.Of[func(B) C],
		O.Of[func(func(A) B) B],

		O.MonadMap[A, A],
		O.MonadMap[A, B],
		O.MonadMap[A, C],
		O.MonadMap[B, C],

		O.MonadMap[func(B) C, func(func(A) B) func(A) C],

		O.MonadChain[A, A],
		O.MonadChain[A, B],
		O.MonadChain[A, C],
		O.MonadChain[B, C],

		O.MonadAp[A, A],
		O.MonadAp[B, A],
		O.MonadAp[C, B],
		O.MonadAp[C, A],

		O.MonadAp[B, func(A) B],
		O.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
