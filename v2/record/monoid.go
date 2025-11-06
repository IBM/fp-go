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

package record

import (
	M "github.com/IBM/fp-go/v2/monoid"
	G "github.com/IBM/fp-go/v2/record/generic"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// UnionMonoid computes the union of two maps of the same type
func UnionMonoid[K comparable, V any](s S.Semigroup[V]) M.Monoid[map[K]V] {
	return G.UnionMonoid[map[K]V](s)
}

// UnionLastMonoid computes the union of two maps of the same type giving the last map precedence
func UnionLastMonoid[K comparable, V any]() M.Monoid[map[K]V] {
	return G.UnionLastMonoid[map[K]V]()
}

// UnionFirstMonoid computes the union of two maps of the same type giving the first map precedence
func UnionFirstMonoid[K comparable, V any]() M.Monoid[map[K]V] {
	return G.UnionFirstMonoid[map[K]V]()
}

// MergeMonoid computes the union of two maps of the same type giving the last map precedence
func MergeMonoid[K comparable, V any]() M.Monoid[map[K]V] {
	return G.UnionLastMonoid[map[K]V]()
}
