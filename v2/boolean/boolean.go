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

package boolean

import (
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/ord"
)

var (
	// MonoidAny is the boolean [monoid.Monoid] under disjunction
	MonoidAny = monoid.MakeMonoid(
		func(l, r bool) bool {
			return l || r
		},
		false,
	)

	// MonoidAll is the boolean [monoid.Monoid] under conjuction
	MonoidAll = monoid.MakeMonoid(
		func(l, r bool) bool {
			return l && r
		},
		true,
	)

	// Eq is the equals predicate for boolean
	Eq = eq.FromStrictEquals[bool]()

	// Ord is the strict ordering for boolean
	Ord = ord.MakeOrd(func(l, r bool) int {
		if l {
			if r {
				return 0
			}
			return +1
		}
		if r {
			return -1
		}
		return 0
	}, func(l, r bool) bool {
		return l == r
	})
)
