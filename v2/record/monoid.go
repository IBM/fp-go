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
	G "github.com/IBM/fp-go/v2/record/generic"
	S "github.com/IBM/fp-go/v2/semigroup"
)

// UnionMonoid computes the union of two maps of the same type
//
//go:inline
func UnionMonoid[K comparable, V any](s S.Semigroup[V]) Monoid[Record[K, V]] {
	return G.UnionMonoid[Record[K, V]](s)
}

// UnionLastMonoid computes the union of two maps of the same type giving the last map precedence
//
//go:inline
func UnionLastMonoid[K comparable, V any]() Monoid[Record[K, V]] {
	return G.UnionLastMonoid[Record[K, V]]()
}

// UnionFirstMonoid computes the union of two maps of the same type giving the first map precedence
//
//go:inline
func UnionFirstMonoid[K comparable, V any]() Monoid[Record[K, V]] {
	return G.UnionFirstMonoid[Record[K, V]]()
}

// MergeMonoid computes the union of two maps of the same type giving the last map precedence
//
//go:inline
func MergeMonoid[K comparable, V any]() Monoid[Record[K, V]] {
	return G.UnionLastMonoid[Record[K, V]]()
}

// IntersectionMonoid computes the intersection of two maps of the same type,
// combining values for keys that exist in both maps using the provided semigroup.
//
// The empty element is an empty map. When concatenating two maps:
//   - Only keys that exist in both maps are kept in the result
//   - Values for shared keys are combined using the provided semigroup
//
// Note that this is not a true monoid in the mathematical sense: the empty map
// acts as an absorbing element (annihilator) rather than an identity, since
// intersecting with the empty map always yields the empty map.
//
//go:inline
func IntersectionMonoid[K comparable, V any](s S.Semigroup[V]) Monoid[Record[K, V]] {
	return G.IntersectionMonoid[Record[K, V]](s)
}

// IntersectionLastMonoid computes the intersection of two maps of the same type,
// keeping only keys present in both maps and giving the last (right) map precedence
// for the values of those shared keys.
//
// The empty element is an empty map. When concatenating two maps:
//   - Only keys that exist in both maps are kept in the result
//   - For shared keys the value from the second (right) map is used
//
// Note that this is not a true monoid in the mathematical sense: the empty map
// acts as an absorbing element (annihilator) rather than an identity, since
// intersecting with the empty map always yields the empty map.
//
//go:inline
func IntersectionLastMonoid[K comparable, V any]() Monoid[Record[K, V]] {
	return G.IntersectionLastMonoid[Record[K, V]]()
}

// IntersectionFirstMonoid computes the intersection of two maps of the same type,
// keeping only keys present in both maps and giving the first (left) map precedence
// for the values of those shared keys.
//
// The empty element is an empty map. When concatenating two maps:
//   - Only keys that exist in both maps are kept in the result
//   - For shared keys the value from the first (left) map is used
//
// Note that this is not a true monoid in the mathematical sense: the empty map
// acts as an absorbing element (annihilator) rather than an identity, since
// intersecting with the empty map always yields the empty map.
//
//go:inline
func IntersectionFirstMonoid[K comparable, V any]() Monoid[Record[K, V]] {
	return G.IntersectionFirstMonoid[Record[K, V]]()
}
