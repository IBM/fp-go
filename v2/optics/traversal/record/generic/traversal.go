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
	R "github.com/IBM/fp-go/v2/internal/record"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
)

// FromRecord returns a traversal from a record
func FromRecord[MA ~map[K]A, MB ~map[K]B, K comparable, A, B, HKTB, HKTAB, HKTRB any](
	fof func(MB) HKTRB,
	fmap func(func(MB) func(B) MB) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,
) G.Traversal[MA, A, HKTRB, HKTB] {
	return func(f func(A) HKTB) func(s MA) HKTRB {
		return func(s MA) HKTRB {
			return R.MonadTraverse(fof, fmap, fap, s, f)
		}
	}
}
