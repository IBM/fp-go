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
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestUnionSemigroup(t *testing.T) {
	// Test with sum semigroup - values should be added for duplicate keys
	sumSemigroup := N.SemigroupSum[int]()
	mapSemigroup := UnionSemigroup[string](sumSemigroup)

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 3, "c": 4}
	result := mapSemigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 5, "c": 4}
	assert.Equal(t, expected, result)
}

func TestUnionSemigroupString(t *testing.T) {
	// Test with string semigroup - strings should be concatenated
	stringSemigroup := S.Semigroup
	mapSemigroup := UnionSemigroup[string](stringSemigroup)

	map1 := map[string]string{"a": "Hello", "b": "World"}
	map2 := map[string]string{"b": "!", "c": "Goodbye"}
	result := mapSemigroup.Concat(map1, map2)

	expected := map[string]string{"a": "Hello", "b": "World!", "c": "Goodbye"}
	assert.Equal(t, expected, result)
}

func TestUnionSemigroupProduct(t *testing.T) {
	// Test with product semigroup - values should be multiplied
	prodSemigroup := N.SemigroupProduct[int]()
	mapSemigroup := UnionSemigroup[string](prodSemigroup)

	map1 := map[string]int{"a": 2, "b": 3}
	map2 := map[string]int{"b": 4, "c": 5}
	result := mapSemigroup.Concat(map1, map2)

	expected := map[string]int{"a": 2, "b": 12, "c": 5}
	assert.Equal(t, expected, result)
}

func TestUnionSemigroupEmpty(t *testing.T) {
	// Test with empty maps
	sumSemigroup := N.SemigroupSum[int]()
	mapSemigroup := UnionSemigroup[string](sumSemigroup)

	map1 := map[string]int{"a": 1}
	empty := map[string]int{}

	result1 := mapSemigroup.Concat(map1, empty)
	assert.Equal(t, map1, result1)

	result2 := mapSemigroup.Concat(empty, map1)
	assert.Equal(t, map1, result2)
}

func TestUnionLastSemigroup(t *testing.T) {
	// Test that last (right) value wins for duplicate keys
	semigroup := UnionLastSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 3, "c": 4}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 3, "c": 4}
	assert.Equal(t, expected, result)
}

func TestUnionLastSemigroupNoOverlap(t *testing.T) {
	// Test with no overlapping keys
	semigroup := UnionLastSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"c": 3, "d": 4}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	assert.Equal(t, expected, result)
}

func TestUnionLastSemigroupAllOverlap(t *testing.T) {
	// Test with all keys overlapping
	semigroup := UnionLastSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 10, "b": 20}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 10, "b": 20}
	assert.Equal(t, expected, result)
}

func TestUnionLastSemigroupEmpty(t *testing.T) {
	// Test with empty maps
	semigroup := UnionLastSemigroup[string, int]()

	map1 := map[string]int{"a": 1}
	empty := map[string]int{}

	result1 := semigroup.Concat(map1, empty)
	assert.Equal(t, map1, result1)

	result2 := semigroup.Concat(empty, map1)
	assert.Equal(t, map1, result2)
}

func TestUnionFirstSemigroup(t *testing.T) {
	// Test that first (left) value wins for duplicate keys
	semigroup := UnionFirstSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"b": 3, "c": 4}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 2, "c": 4}
	assert.Equal(t, expected, result)
}

func TestUnionFirstSemigroupNoOverlap(t *testing.T) {
	// Test with no overlapping keys
	semigroup := UnionFirstSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"c": 3, "d": 4}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}
	assert.Equal(t, expected, result)
}

func TestUnionFirstSemigroupAllOverlap(t *testing.T) {
	// Test with all keys overlapping
	semigroup := UnionFirstSemigroup[string, int]()

	map1 := map[string]int{"a": 1, "b": 2}
	map2 := map[string]int{"a": 10, "b": 20}
	result := semigroup.Concat(map1, map2)

	expected := map[string]int{"a": 1, "b": 2}
	assert.Equal(t, expected, result)
}

func TestUnionFirstSemigroupEmpty(t *testing.T) {
	// Test with empty maps
	semigroup := UnionFirstSemigroup[string, int]()

	map1 := map[string]int{"a": 1}
	empty := map[string]int{}

	result1 := semigroup.Concat(map1, empty)
	assert.Equal(t, map1, result1)

	result2 := semigroup.Concat(empty, map1)
	assert.Equal(t, map1, result2)
}

// Test associativity law for UnionSemigroup
func TestUnionSemigroupAssociativity(t *testing.T) {
	sumSemigroup := N.SemigroupSum[int]()
	mapSemigroup := UnionSemigroup[string](sumSemigroup)

	map1 := map[string]int{"a": 1}
	map2 := map[string]int{"a": 2, "b": 3}
	map3 := map[string]int{"b": 4, "c": 5}

	// (map1 + map2) + map3
	left := mapSemigroup.Concat(mapSemigroup.Concat(map1, map2), map3)
	// map1 + (map2 + map3)
	right := mapSemigroup.Concat(map1, mapSemigroup.Concat(map2, map3))

	assert.Equal(t, left, right)
}

// Test associativity law for UnionLastSemigroup
func TestUnionLastSemigroupAssociativity(t *testing.T) {
	semigroup := UnionLastSemigroup[string, int]()

	map1 := map[string]int{"a": 1}
	map2 := map[string]int{"a": 2, "b": 3}
	map3 := map[string]int{"b": 4, "c": 5}

	// (map1 + map2) + map3
	left := semigroup.Concat(semigroup.Concat(map1, map2), map3)
	// map1 + (map2 + map3)
	right := semigroup.Concat(map1, semigroup.Concat(map2, map3))

	assert.Equal(t, left, right)
}

// Test associativity law for UnionFirstSemigroup
func TestUnionFirstSemigroupAssociativity(t *testing.T) {
	semigroup := UnionFirstSemigroup[string, int]()

	map1 := map[string]int{"a": 1}
	map2 := map[string]int{"a": 2, "b": 3}
	map3 := map[string]int{"b": 4, "c": 5}

	// (map1 + map2) + map3
	left := semigroup.Concat(semigroup.Concat(map1, map2), map3)
	// map1 + (map2 + map3)
	right := semigroup.Concat(map1, semigroup.Concat(map2, map3))

	assert.Equal(t, left, right)
}
