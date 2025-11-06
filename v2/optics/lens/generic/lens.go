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
	L "github.com/IBM/fp-go/v2/optics/lens"
)

// AsTraversal converts a lens to a traversal
func AsTraversal[R ~func(func(A) HKTA) func(S) HKTS, S, A, HKTS, HKTA any](
	fmap func(HKTA, func(A) S) HKTS,
) func(L.Lens[S, A]) R {
	return func(sa L.Lens[S, A]) R {
		return func(f func(a A) HKTA) func(S) HKTS {
			return func(s S) HKTS {
				return fmap(f(sa.Get(s)), func(a A) S {
					return sa.Set(a)(s)
				})
			}
		}
	}
}
