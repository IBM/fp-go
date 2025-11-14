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

// addToMap makes a deep copy of a map and adds a value
func addToMap[A comparable](a A, m map[A]bool) map[A]bool {
	cpy := make(map[A]bool, len(m)+1)
	for k, v := range m {
		cpy[k] = v
	}
	cpy[a] = true
	return cpy
}

func Uniq[AS ~func() Option[Pair[AS, A]], K comparable, A any](f func(A) K) func(as AS) AS {

	var recurse func(as AS, mp map[K]bool) AS

	recurse = func(as AS, mp map[K]bool) AS {
		return F.Nullary2(
			as,
			O.Chain(func(a Pair[AS, A]) Option[Pair[AS, A]] {
				return F.Pipe3(
					P.Tail(a),
					f,
					O.FromPredicate(func(k K) bool {
						_, ok := mp[k]
						return !ok
					}),
					O.Fold(recurse(P.Head(a), mp), func(k K) Option[Pair[AS, A]] {
						return O.Of(P.MakePair(recurse(P.Head(a), addToMap(k, mp)), P.Tail(a)))
					}),
				)
			}),
		)
	}

	return F.Bind2nd(recurse, make(map[K]bool, 0))
}

func StrictUniq[AS ~func() Option[Pair[AS, A]], A comparable](as AS) AS {
	return Uniq[AS](F.Identity[A])(as)
}
