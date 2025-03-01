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

	EQ "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/internal/monad/testing"
	"github.com/IBM/fp-go/v2/io"
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
		io.Eq(eqa),
		io.Eq(eqb),
		io.Eq(eqc),

		io.Of[A],
		io.Of[B],
		io.Of[C],

		io.Of[func(A) A],
		io.Of[func(A) B],
		io.Of[func(B) C],
		io.Of[func(func(A) B) B],

		io.MonadMap[A, A],
		io.MonadMap[A, B],
		io.MonadMap[A, C],
		io.MonadMap[B, C],

		io.MonadMap[func(B) C, func(func(A) B) func(A) C],

		io.MonadChain[A, A],
		io.MonadChain[A, B],
		io.MonadChain[A, C],
		io.MonadChain[B, C],

		io.MonadAp[A, A],
		io.MonadAp[B, A],
		io.MonadAp[C, B],
		io.MonadAp[C, A],

		io.MonadAp[B, func(A) B],
		io.MonadAp[func(A) C, func(A) B],

		ab,
		bc,
	)

}
