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
	"sort"

	O "github.com/IBM/fp-go/ord"
)

// Sort implements a stable sort on the array given the provided ordering
func Sort[GA ~[]T, T any](ord O.Ord[T]) func(ma GA) GA {

	return func(ma GA) GA {
		// nothing to sort
		l := len(ma)
		if l < 2 {
			return ma
		}
		// copy
		cpy := make(GA, l)
		copy(cpy, ma)
		sort.Slice(cpy, func(i, j int) bool {
			return ord.Compare(cpy[i], cpy[j]) < 0
		})
		return cpy
	}
}
