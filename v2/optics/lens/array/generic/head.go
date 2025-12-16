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
	AA "github.com/IBM/fp-go/v2/array/generic"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
)

// AtHead focusses on the head of an array. The setter works as follows
// - if the new value is none, the result will be an empty array
// - if the new value is some and the array is empty, it creates a new array with one element
// - if the new value is some and the array is not empty, it replaces the head
func AtHead[AS []A, A any]() L.Lens[AS, O.Option[A]] {
	return L.MakeLensWithName(AA.Head[AS, A], func(as AS, a O.Option[A]) AS {
		return O.MonadFold(a, AA.Empty[AS], func(v A) AS {
			if AA.IsEmpty(as) {
				return AA.Of[AS](v)
			}
			cpy := AA.Copy(as)
			cpy[0] = v
			return cpy
		})
	}, "Head")
}
