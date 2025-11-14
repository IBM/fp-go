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
	PR "github.com/IBM/fp-go/v2/predicate"
)

// DropWhile creates an [Iterator] that drops elements from the [Iterator] as long as the predicate is true; afterwards, returns every element.
// Note, the [Iterator] does not produce any output until the predicate first becomes false
func DropWhile[GU ~func() Option[Pair[GU, U]], U any](pred Predicate[U]) func(GU) GU {
	// avoid cyclic references
	var m func(Option[Pair[GU, U]]) Option[Pair[GU, U]]

	fromPred := O.FromPredicate(PR.Not(PR.ContraMap(P.Tail[GU, U])(pred)))

	recurse := func(mu GU) GU {
		return F.Nullary2(
			mu,
			m,
		)
	}

	m = O.Chain(func(t Pair[GU, U]) Option[Pair[GU, U]] {
		return F.Pipe2(
			t,
			fromPred,
			O.Fold(recurse(Next(t)), O.Of[Pair[GU, U]]),
		)
	})

	return recurse
}
