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

package generic

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	G "github.com/IBM/fp-go/v2/internal/eq"
)

// Eq implements the equals predicate for values contained in the IO monad
func Eq[GA ~func() A, A any](e EQ.Eq[A]) EQ.Eq[GA] {
	// comparator for the monad
	eq := G.Eq(
		MonadMap[GA, func() func(A) bool, A, func(A) bool],
		MonadAp[GA, func() bool, func() func(A) bool, A, bool],
		e,
	)
	// eagerly execute
	return EQ.FromEquals(func(l, r GA) bool {
		return eq(l, r)()
	})
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[GA ~func() A, A comparable]() EQ.Eq[GA] {
	return Eq[GA](EQ.FromStrictEquals[A]())
}
