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
	AR "github.com/IBM/fp-go/v2/internal/array"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// FromArray returns a traversal from an array
func FromArray[GA ~[]A, GB ~[]B, A, B, HKTB, HKTAB, HKTRB any](
	fof func(GB) HKTRB,
	fmap func(func(GB) func(B) GB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
) G.Traversal[GA, A, HKTRB, HKTB] {
	return func(f func(A) HKTB) func(s GA) HKTRB {
		return func(s GA) HKTRB {
			return AR.MonadTraverse(fof, fmap, fap, s, f)
		}
	}
}
