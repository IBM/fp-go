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

	return L.AssertLaws(t,
		WRT.Eq(eqw, eqa),
		WRT.Eq(eqw, eqb),
		WRT.Eq(eqw, eqc),

		WRT.Of[A](m),
		WRT.Of[B](m),
		WRT.Of[C](m),

		WRT.Of[func(A) A](m),
		WRT.Of[func(A) B](m),
		WRT.Of[func(B) C](m),
		WRT.Of[func(func(A) B) B](m),

		WRT.MonadMap[func(A) A, W, A, A],
		WRT.MonadMap[func(A) B, W, A, B],
		WRT.MonadMap[func(A) C, W, A, C],
		WRT.MonadMap[func(B) C, W, B, C],

		WRT.MonadMap[func(func(B) C) func(func(A) B) func(A) C, W, func(B) C, func(func(A) B) func(A) C],

		WRT.MonadChain[func(A) WRT.Writer[W, A], W, A, A],
		WRT.MonadChain[func(A) WRT.Writer[W, B], W, A, B],
		WRT.MonadChain[func(A) WRT.Writer[W, C], W, A, C],
		WRT.MonadChain[func(B) WRT.Writer[W, C], W, B, C],

		WRT.MonadAp[A, A, W],
		WRT.MonadAp[B, A, W],
		WRT.MonadAp[C, B, W],
		WRT.MonadAp[C, A, W],

		WRT.MonadAp[B, func(A) B, W],
		WRT.MonadAp[func(A) C, func(A) B, W],

		ab,
		bc,
	)

}
