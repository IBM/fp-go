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
	E "github.com/IBM/fp-go/v2/eq"
)

func equals[M ~map[K]V, K comparable, V any](left, right M, eq func(V, V) bool) bool {
	if len(left) != len(right) {
		return false
	}
	for k, v1 := range left {
		if v2, ok := right[k]; !ok || !eq(v1, v2) {
			return false
		}
	}
	return true
}

func Eq[M ~map[K]V, K comparable, V any](e E.Eq[V]) E.Eq[M] {
	eq := e.Equals
	return E.FromEquals(func(left, right M) bool {
		return equals(left, right, eq)
	})
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[M ~map[K]V, K, V comparable]() E.Eq[M] {
	return Eq[M](E.FromStrictEquals[V]())
}
