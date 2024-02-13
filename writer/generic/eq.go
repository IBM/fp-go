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
	P "github.com/IBM/fp-go/pair"
)

// Constructs an equal predicate for a [Writer]
func Eq[GA ~func() P.Pair[A, W], W, A any](w EQ.Eq[W], a EQ.Eq[A]) EQ.Eq[GA] {
	eqp := P.Eq(a, w)
	return EQ.FromEquals(func(l, r GA) bool {
		return eqp.Equals(l(), r())
	})
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[GA ~func() P.Pair[A, W], W, A comparable]() EQ.Eq[GA] {
	return Eq[GA](EQ.FromStrictEquals[W](), EQ.FromStrictEquals[A]())
}
