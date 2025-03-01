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
	EQ "github.com/IBM/fp-go/v2/eq"
	P "github.com/IBM/fp-go/v2/pair"
)

// Constructs an equal predicate for a [State]
func Eq[GA ~func(S) P.Pair[A, S], S, A any](w EQ.Eq[S], a EQ.Eq[A]) func(S) EQ.Eq[GA] {
	eqp := P.Eq(a, w)
	return func(s S) EQ.Eq[GA] {
		return EQ.FromEquals(func(l, r GA) bool {
			return eqp.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[GA ~func(S) P.Pair[A, S], S, A comparable]() func(S) EQ.Eq[GA] {
	return Eq[GA](EQ.FromStrictEquals[S](), EQ.FromStrictEquals[A]())
}
