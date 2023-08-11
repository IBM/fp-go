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
	G "github.com/IBM/fp-go/array/generic"
	O "github.com/IBM/fp-go/ord"
)

// Sort implements a stable sort on the array given the provided ordering
func Sort[T any](ord O.Ord[T]) func(ma []T) []T {
	return G.Sort[[]T](ord)
}

// SortByKey implements a stable sort on the array given the provided ordering on an extracted key
func SortByKey[K, T any](ord O.Ord[K], f func(T) K) func(ma []T) []T {
	return G.SortByKey[[]T](ord, f)
}

// SortBy implements a stable sort on the array given the provided ordering
func SortBy[T any](ord []O.Ord[T]) func(ma []T) []T {
	return G.SortBy[[]T, []O.Ord[T]](ord)
}
