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
)

// UnionSemigroup creates a semigroup for maps that combines two maps using the provided
// semigroup for resolving conflicts when the same key exists in both maps.
//
// When concatenating two maps:
//   - Keys that exist in only one map are included in the result
//   - Keys that exist in both maps have their values combined using the provided semigroup
//
// This is useful when you want custom conflict resolution logic beyond simple "first wins"
// or "last wins" semantics.
//
// Example:
//
//	// Create a semigroup that sums values for duplicate keys
//	sumSemigroup := number.SemigroupSum[int]()
//	mapSemigroup := UnionSemigroup[string, int](sumSemigroup)
//
//	map1 := map[string]int{"a": 1, "b": 2}
//	map2 := map[string]int{"b": 3, "c": 4}
//	result := mapSemigroup.Concat(map1, map2)
//	// result: {"a": 1, "b": 5, "c": 4}  // b values are summed: 2 + 3 = 5
//
// Example with string concatenation:
//
//	stringSemigroup := string.Semigroup
//	mapSemigroup := UnionSemigroup[string, string](stringSemigroup)
//
//	map1 := map[string]string{"a": "Hello", "b": "World"}
//	map2 := map[string]string{"b": "!", "c": "Goodbye"}
//	result := mapSemigroup.Concat(map1, map2)
//	// result: {"a": "Hello", "b": "World!", "c": "Goodbye"}
//
//go:inline
func UnionSemigroup[K comparable, V any](s Semigroup[V]) Semigroup[Record[K, V]] {
	return G.UnionSemigroup[Record[K, V]](s)
}

// UnionLastSemigroup creates a semigroup for maps where the last (right) value wins
// when the same key exists in both maps being concatenated.
//
// This is the most common conflict resolution strategy and is equivalent to using
// the standard map merge operation where right-side values take precedence.
//
// When concatenating two maps:
//   - Keys that exist in only one map are included in the result
//   - Keys that exist in both maps take the value from the second (right) map
//
// Example:
//
//	semigroup := UnionLastSemigroup[string, int]()
//
//	map1 := map[string]int{"a": 1, "b": 2}
//	map2 := map[string]int{"b": 3, "c": 4}
//	result := semigroup.Concat(map1, map2)
//	// result: {"a": 1, "b": 3, "c": 4}  // b takes value from map2 (last wins)
//
// This is useful for:
//   - Configuration overrides (later configs override earlier ones)
//   - Applying updates to a base map
//   - Merging user preferences where newer values should win
//
//go:inline
func UnionLastSemigroup[K comparable, V any]() Semigroup[Record[K, V]] {
	return G.UnionLastSemigroup[Record[K, V]]()
}

// UnionFirstSemigroup creates a semigroup for maps where the first (left) value wins
// when the same key exists in both maps being concatenated.
//
// This is useful when you want to preserve original values and ignore updates for
// keys that already exist.
//
// When concatenating two maps:
//   - Keys that exist in only one map are included in the result
//   - Keys that exist in both maps keep the value from the first (left) map
//
// Example:
//
//	semigroup := UnionFirstSemigroup[string, int]()
//
//	map1 := map[string]int{"a": 1, "b": 2}
//	map2 := map[string]int{"b": 3, "c": 4}
//	result := semigroup.Concat(map1, map2)
//	// result: {"a": 1, "b": 2, "c": 4}  // b keeps value from map1 (first wins)
//
// This is useful for:
//   - Default values (defaults are set first, user values don't override)
//   - Caching (first cached value is kept, subsequent updates ignored)
//   - Immutable registries (first registration wins, duplicates are ignored)
//
//go:inline
func UnionFirstSemigroup[K comparable, V any]() Semigroup[Record[K, V]] {
	return G.UnionFirstSemigroup[Record[K, V]]()
}
