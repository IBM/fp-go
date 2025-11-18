// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/readerioeither"
)

// Eq implements the equals predicate for values contained in the [StateReaderIOEither] monad
func Eq[
	S, R, E, A any](eqr eq.Eq[ReaderIOEither[R, E, Pair[S, A]]]) func(S) eq.Eq[StateReaderIOEither[S, R, E, A]] {
	return func(s S) eq.Eq[StateReaderIOEither[S, R, E, A]] {
		return eq.FromEquals(func(l, r StateReaderIOEither[S, R, E, A]) bool {
			return eqr.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an [eq.Eq] from the canonical comparison function
func FromStrictEquals[
	S comparable, R any, E, A comparable]() func(R) func(S) eq.Eq[StateReaderIOEither[S, R, E, A]] {
	return function.Flow2(
		readerioeither.FromStrictEquals[R, E, Pair[S, A]](),
		Eq[S, R, E, A],
	)
}
