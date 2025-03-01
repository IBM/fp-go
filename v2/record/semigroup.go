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

package record

import (
	G "github.com/IBM/fp-go/v2/record/generic"
	S "github.com/IBM/fp-go/v2/semigroup"
)

func UnionSemigroup[K comparable, V any](s S.Semigroup[V]) S.Semigroup[map[K]V] {
	return G.UnionSemigroup[map[K]V](s)
}

func UnionLastSemigroup[K comparable, V any]() S.Semigroup[map[K]V] {
	return G.UnionLastSemigroup[map[K]V]()
}

func UnionFirstSemigroup[K comparable, V any]() S.Semigroup[map[K]V] {
	return G.UnionFirstSemigroup[map[K]V]()
}
