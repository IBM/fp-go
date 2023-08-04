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
	P "github.com/IBM/fp-go/predicate"
	T "github.com/IBM/fp-go/tuple"
)

// DropWhile creates an [Iterator] that drops elements from the [Iterator] as long as the predicate is true; afterwards, returns every element.
// Note, the [Iterator] does not produce any output until the predicate first becomes false
func DropWhile[GU ~func() O.Option[T.Tuple2[GU, U]], U any](pred func(U) bool) func(GU) GU {
	// avoid cyclic references
	var m func(O.Option[T.Tuple2[GU, U]]) O.Option[T.Tuple2[GU, U]]

	fromPred := O.FromPredicate(P.Not(P.ContraMap(T.Second[GU, U])(pred)))

	recurse := func(mu GU) GU {
		return F.Nullary2(
			mu,
			m,
		)
	}

	m = O.Chain(func(t T.Tuple2[GU, U]) O.Option[T.Tuple2[GU, U]] {
		return F.Pipe2(
			t,
			fromPred,
			O.Fold(recurse(t.F1), O.Of[T.Tuple2[GU, U]]),
		)
	})

	return recurse
}
