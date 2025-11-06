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
	"sort"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/ord"
)

// Sort implements a stable sort on the array given the provided ordering
func Sort[GA ~[]T, T any](ord O.Ord[T]) func(ma GA) GA {
	return SortByKey[GA](ord, F.Identity[T])
}

// SortByKey implements a stable sort on the array given the provided ordering on an extracted key
func SortByKey[GA ~[]T, K, T any](ord O.Ord[K], f func(T) K) func(ma GA) GA {

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
			return ord.Compare(f(cpy[i]), f(cpy[j])) < 0
		})
		return cpy
	}
}

// SortBy implements a stable sort on the array given the provided ordering
func SortBy[GA ~[]T, GO ~[]O.Ord[T], T any](ord GO) func(ma GA) GA {
	return F.Pipe2(
		ord,
		Fold[GO](O.Monoid[T]()),
		Sort[GA, T],
	)
}
