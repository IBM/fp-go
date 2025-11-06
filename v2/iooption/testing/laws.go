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
	IOO "github.com/IBM/fp-go/v2/iooption"
)

// AssertLaws asserts the apply monad laws for the `Either` monad
func AssertLaws[A, B, C any](t *testing.T,
	eqa EQ.Eq[A],
	eqb EQ.Eq[B],
	eqc EQ.Eq[C],

	ab func(A) B,
	bc func(B) C,
) func(a A) bool {

	return L.AssertLaws(t,
		IOO.Eq(eqa),
		IOO.Eq(eqb),
		IOO.Eq(eqc),

		IOO.Of[A],
		IOO.Of[B],
		IOO.Of[C],

		IOO.Of[func(A) A],
		IOO.Of[func(A) B],
		IOO.Of[func(B) C],
		IOO.Of[func(func(A) B) B],

		IOO.MonadMap[A, A],
		IOO.MonadMap[A, B],
		IOO.MonadMap[A, C],
		IOO.MonadMap[B, C],

		IOO.MonadMap[func(B) C, func(func(A) B) func(A) C],

		IOO.MonadChain[A, A],
		IOO.MonadChain[A, B],
		IOO.MonadChain[A, C],
		IOO.MonadChain[B, C],

		IOO.MonadAp[A, A],
		IOO.MonadAp[B, A],
		IOO.MonadAp[C, B],
		IOO.MonadAp[C, A],

		IOO.MonadAp[B, func(A) B],
		IOO.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
