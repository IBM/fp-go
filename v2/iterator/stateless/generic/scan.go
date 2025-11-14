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
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

func apTuple[A, B any](t Pair[func(A) B, A]) Pair[B, A] {
	return P.MakePair(P.Head(t)(P.Tail(t)), P.Tail(t))
}

func Scan[GV ~func() Option[Pair[GV, V]], GU ~func() Option[Pair[GU, U]], FCT ~func(V, U) V, U, V any](f FCT, initial V) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(GU) func(V) GV

	recurse := func(ma GU, current V) GV {
		return F.Nullary2(
			ma,
			O.Map(F.Flow2(
				P.BiMap(m, F.Bind1st(f, current)),
				apTuple[V, GV],
			)),
		)
	}

	m = F.Curry2(recurse)

	return F.Bind2nd(recurse, initial)
}
