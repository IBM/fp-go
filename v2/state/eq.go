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

package state

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/pair"
)

// Constructs an equal predicate for a [State]
func Eq[S, A any](w eq.Eq[S], a eq.Eq[A]) func(S) eq.Eq[State[S, A]] {
	eqp := pair.Eq(w, a)
	return func(s S) eq.Eq[State[S, A]] {
		return eq.FromEquals(func(l, r State[S, A]) bool {
			return eqp.Equals(l(s), r(s))
		})
	}
}

// FromStrictEquals constructs an [eq.Eq] from the canonical comparison function
func FromStrictEquals[S, A comparable]() func(S) eq.Eq[State[S, A]] {
	return Eq(eq.FromStrictEquals[S](), eq.FromStrictEquals[A]())
}
