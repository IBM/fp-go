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

package array

import (
	"github.com/IBM/fp-go/v2/internal/array"
)

func Traverse[A, B, HKTB, HKTAB, HKTRB any](
	fof func([]B) HKTRB,
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	f func(A) HKTB) func([]A) HKTRB {
	return array.Traverse[[]A](fof, fmap, fap, f)
}

func MonadTraverse[A, B, HKTB, HKTAB, HKTRB any](
	fof func([]B) HKTRB,
	fmap func(func([]B) func(B) []B) func(HKTRB) HKTAB,
	fap func(HKTB) func(HKTAB) HKTRB,

	ta []A,
	f func(A) HKTB) HKTRB {

	return array.MonadTraverse(fof, fmap, fap, ta, f)
}
