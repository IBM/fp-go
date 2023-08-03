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
	N "github.com/IBM/fp-go/number/integer"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

func Take[GU ~func() O.Option[T.Tuple2[GU, U]], U any](n int) func(ma GU) GU {
	// pre-declare to avoid cyclic reference
	var recurse func(ma GU, idx int) GU

	fromPred := O.FromPredicate(N.Between(0, n))

	recurse = func(ma GU, idx int) GU {
		return func() O.Option[T.Tuple2[GU, U]] {
			return F.Pipe2(
				idx,
				fromPred,
				O.Chain(F.Ignore1of1[int](F.Nullary2(
					ma,
					O.Map(func(t T.Tuple2[GU, U]) T.Tuple2[GU, U] {
						return T.MakeTuple2(recurse(t.F1, idx+1), t.F2)
					}),
				))),
			)
		}
	}

	return F.Bind2nd(recurse, 0)
}
