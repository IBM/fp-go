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

// addToMap makes a deep copy of a map and adds a value
func addToMap[A comparable](a A, m map[A]bool) map[A]bool {
	cpy := make(map[A]bool, len(m)+1)
	for k, v := range m {
		cpy[k] = v
	}
	cpy[a] = true
	return cpy
}

func Uniq[AS ~func() O.Option[T.Tuple2[AS, A]], K comparable, A any](f func(A) K) func(as AS) AS {

	var recurse func(as AS, mp map[K]bool) AS

	recurse = func(as AS, mp map[K]bool) AS {
		return F.Nullary2(
			as,
			O.Chain(func(a T.Tuple2[AS, A]) O.Option[T.Tuple2[AS, A]] {
				return F.Pipe3(
					a.F2,
					f,
					O.FromPredicate(func(k K) bool {
						_, ok := mp[k]
						return !ok
					}),
					O.Fold(recurse(a.F1, mp), func(k K) O.Option[T.Tuple2[AS, A]] {
						return O.Of(T.MakeTuple2(recurse(a.F1, addToMap(k, mp)), a.F2))
					}),
				)
			}),
		)
	}

	return F.Bind2nd(recurse, make(map[K]bool, 0))
}

func StrictUniq[AS ~func() O.Option[T.Tuple2[AS, A]], A comparable](as AS) AS {
	return Uniq[AS](F.Identity[A])(as)
}
