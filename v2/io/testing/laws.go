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

	return L.MonadAssertLaws(t,
		io.Eq(eqa),
		io.Eq(eqb),
		io.Eq(eqc),

		io.Pointed[C](),
		io.Pointed[func(A) A](),
		io.Pointed[func(B) C](),
		io.Pointed[func(func(A) B) B](),

		io.Functor[func(B) C, func(func(A) B) func(A) C](),

		io.Applicative[func(A) B, B](),
		io.Applicative[func(A) B, func(A) C](),

		io.Monad[A, A](),
		io.Monad[A, B](),
		io.Monad[A, C](),
		io.Monad[B, C](),

		ab,
		bc,
	)

}
