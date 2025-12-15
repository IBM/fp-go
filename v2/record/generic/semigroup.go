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
	S "github.com/IBM/fp-go/v2/semigroup"
)

func UnionSemigroup[N ~map[K]V, K comparable, V any](s S.Semigroup[V]) S.Semigroup[N] {
	return S.MakeSemigroup(func(first N, second N) N {
		return union(S.ToMagma(s), first, second)
	})
}

func UnionLastSemigroup[N ~map[K]V, K comparable, V any]() S.Semigroup[N] {
	return S.MakeSemigroup(unionLast[N])
}

func UnionFirstSemigroup[N ~map[K]V, K comparable, V any]() S.Semigroup[N] {
	return S.MakeSemigroup(func(first N, second N) N {
		return unionLast(second, first)
	})
}
