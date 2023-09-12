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

package monoid

import (
	S "github.com/IBM/fp-go/semigroup"
)

func AlternativeMonoid[A, HKTA, HKTFA any](
	fof func(A) HKTA,

	fmap func(HKTA, func(A) func(A) A) HKTFA,
	fap func(HKTFA, HKTA) HKTA,

	falt func(HKTA, func() HKTA) HKTA,

	m Monoid[A],

) Monoid[HKTA] {

	sg := ApplicativeMonoid(fof, fmap, fap, m)

	return MakeMonoid(
		func(first, second HKTA) HKTA {
			snd := func() HKTA { return second }

			return falt(sg.Concat(first, second), func() HKTA {
				return falt(first, snd)
			})
		},
		sg.Empty(),
	)
}

func AltMonoid[HKTA any](
	fzero func() HKTA,
	falt func(HKTA, func() HKTA) HKTA,

) Monoid[HKTA] {

	return MakeMonoid(
		S.AltSemigroup(falt).Concat,
		fzero(),
	)
}
