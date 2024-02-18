// Copyright (c) 2024 IBM Corp.
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
	EQ "github.com/IBM/fp-go/eq"
	P "github.com/IBM/fp-go/pair"
	RIOE "github.com/IBM/fp-go/readerioeither"
	G "github.com/IBM/fp-go/statereaderioeither/generic"
)

// Eq implements the equals predicate for values contained in the [StateReaderIOEither] monad
func Eq[
	S, R, E, A any](eqr EQ.Eq[RIOE.ReaderIOEither[R, E, P.Pair[A, S]]]) func(S) EQ.Eq[StateReaderIOEither[S, R, E, A]] {
	return G.Eq[StateReaderIOEither[S, R, E, A]](eqr)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[
	S, R any, E, A comparable]() func(R) func(S) EQ.Eq[StateReaderIOEither[S, R, E, A]] {
	return G.FromStrictEquals[StateReaderIOEither[S, R, E, A]]()
}
