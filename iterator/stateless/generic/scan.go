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
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

func apTuple[A, B any](t T.Tuple2[func(A) B, A]) T.Tuple2[B, A] {
	return T.MakeTuple2(t.F1(t.F2), t.F2)
}

func Scan[GV ~func() O.Option[T.Tuple2[GV, V]], GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(V, U) V, U, V any](f FCT, initial V) func(ma GU) GV {
	// pre-declare to avoid cyclic reference
	var m func(GU) func(V) GV

	recurse := func(ma GU, current V) GV {
		return F.Nullary2(
			ma,
			O.Map(F.Flow2(
				T.Map2(m, F.Bind1st(f, current)),
				apTuple[V, GV],
			)),
		)
	}

	m = F.Curry2(recurse)

	return F.Bind2nd(recurse, initial)
}
