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

package generic

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	G "github.com/IBM/fp-go/v2/internal/eq"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[GEA ~func(R) GIOA, GIOA ~func() A, R, A any](e EQ.Eq[A]) func(r R) EQ.Eq[GEA] {
	// comparator for the monad
	eq := G.Eq(
		MonadMap[GEA, func(R) func() func(A) bool, GIOA, func() func(A) bool, R, A, func(A) bool],
		MonadAp[GEA, func(R) func() bool, func(R) func() func(A) bool],
		e,
	)
	// eagerly execute
	return func(ctx R) EQ.Eq[GEA] {
		return EQ.FromEquals(func(l, r GEA) bool {
			return eq(l, r)(ctx)()
		})
	}
}
