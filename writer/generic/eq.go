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
	EQ "github.com/IBM/fp-go/eq"
	T "github.com/IBM/fp-go/tuple"
)

// Constructs an equal predicate for a [Writer]
func Eq[GA ~func() T.Tuple2[A, W], W, A any](w EQ.Eq[W], a EQ.Eq[A]) EQ.Eq[GA] {
	return EQ.FromEquals(func(l, r GA) bool {
		ll := l()
		rr := r()

		return a.Equals(ll.F1, rr.F1) && w.Equals(ll.F2, rr.F2)
	})
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[GA ~func() T.Tuple2[A, W], W, A comparable]() EQ.Eq[GA] {
	return Eq[GA](EQ.FromStrictEquals[W](), EQ.FromStrictEquals[A]())
}
