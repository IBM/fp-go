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

package semigroup

import (
	F "github.com/IBM/fp-go/function"
)

/*
*
HKTA = HKT<A>
HKTFA = HKT<func(A)A>
*/
func ApplySemigroup[A, HKTA, HKTFA any](
	fmap func(HKTA, func(A) func(A) A) HKTFA,
	fap func(HKTFA, HKTA) HKTA,

	s Semigroup[A],
) Semigroup[HKTA] {

	cb := F.Curry2(s.Concat)
	return MakeSemigroup(func(first HKTA, second HKTA) HKTA {
		return fap(fmap(first, cb), second)
	})
}
