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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
	P "github.com/IBM/fp-go/v2/pair"
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// Eq implements the equals predicate for values contained in the [StateReaderIOEither] monad
func Eq[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any](eqr EQ.Eq[RIOEA]) func(S) EQ.Eq[SRIOEA] {
	return func(s S) EQ.Eq[SRIOEA] {
		return EQ.FromEquals(func(l, r SRIOEA) bool {
			return eqr.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R any, E, A comparable]() func(R) func(S) EQ.Eq[SRIOEA] {
	return F.Flow2(
		G.FromStrictEquals[RIOEA](),
		Eq[SRIOEA, RIOEA, IOEA, S, R, E, A],
	)
}
