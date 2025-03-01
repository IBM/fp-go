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

package prism

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// AsTraversal converts a prism to a traversal
func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
	fof func(S) HKTS,
	fmap func(HKTA, func(A) S) HKTS,
) func(Prism[S, A]) R {
	return func(sa Prism[S, A]) R {
		return func(f func(a A) HKTA) func(S) HKTS {
			return func(s S) HKTS {
				return F.Pipe2(
					s,
					sa.GetOption,
					O.Fold(
						F.Nullary2(F.Constant(s), fof),
						func(a A) HKTS {
							return fmap(f(a), func(a A) S {
								return prismModify(F.Constant1[A](a), sa, s)
							})
						},
					),
				)
			}
		}
	}
}
