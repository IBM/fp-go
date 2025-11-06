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
	"github.com/IBM/fp-go/v2/lazy"
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
		lazy.Eq(eqa),
		lazy.Eq(eqb),
		lazy.Eq(eqc),

		lazy.Of[A],
		lazy.Of[B],
		lazy.Of[C],

		lazy.Of[func(A) A],
		lazy.Of[func(A) B],
		lazy.Of[func(B) C],
		lazy.Of[func(func(A) B) B],

		lazy.MonadMap[A, A],
		lazy.MonadMap[A, B],
		lazy.MonadMap[A, C],
		lazy.MonadMap[B, C],

		lazy.MonadMap[func(B) C, func(func(A) B) func(A) C],

		lazy.MonadChain[A, A],
		lazy.MonadChain[A, B],
		lazy.MonadChain[A, C],
		lazy.MonadChain[B, C],

		lazy.MonadAp[A, A],
		lazy.MonadAp[B, A],
		lazy.MonadAp[C, B],
		lazy.MonadAp[C, A],

		lazy.MonadAp[B, func(A) B],
		lazy.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
