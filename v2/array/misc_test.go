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

package array

import (
	"testing"

	OR "github.com/IBM/fp-go/v2/ord"
	"github.com/stretchr/testify/assert"
)

func TestAnyWithIndex(t *testing.T) {
	src := []int{1, 2, 3, 4, 5}
	checker := AnyWithIndex(func(i, x int) bool {
		return i == 2 && x == 3
	})
	assert.True(t, checker(src))

	checker2 := AnyWithIndex(func(i, x int) bool {
		return i == 10
	})
	assert.False(t, checker2(src))
}

func TestSemigroup(t *testing.T) {
	sg := Semigroup[int]()
	result := sg.Concat([]int{1, 2}, []int{3, 4})
	assert.Equal(t, []int{1, 2, 3, 4}, result)
}

func TestArrayConcatAll(t *testing.T) {
	result := ArrayConcatAll(
		[]int{1, 2},
		[]int{3, 4},
		[]int{5, 6},
	)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, result)

	// Test with empty arrays
	result2 := ArrayConcatAll(
		[]int{},
		[]int{1},
		[]int{},
	)
	assert.Equal(t, []int{1}, result2)
}

func TestMonad(t *testing.T) {
	m := Monad[int, string]()

	// Test Map
	mapFn := m.Map(func(x int) string {
		return string(rune('a' + x - 1))
	})
	mapped := mapFn([]int{1, 2, 3})
	assert.Equal(t, []string{"a", "b", "c"}, mapped)

	// Test Chain
	chainFn := m.Chain(func(x int) []string {
		return []string{string(rune('a' + x - 1))}
	})
	chained := chainFn([]int{1, 2})
	assert.Equal(t, []string{"a", "b"}, chained)

	// Test Of
	ofResult := m.Of(42)
	assert.Equal(t, []int{42}, ofResult)
}

func TestSortByKey(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Charlie", 35},
	}

	sorter := SortByKey(OR.FromStrictCompare[int](), func(p Person) int {
		return p.Age
	})
	result := sorter(people)

	assert.Equal(t, "Bob", result[0].Name)
	assert.Equal(t, "Alice", result[1].Name)
	assert.Equal(t, "Charlie", result[2].Name)
}

func TestUniqByKey(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	people := []Person{
		{"Alice", 30},
		{"Bob", 25},
		{"Alice", 35},
		{"Charlie", 30},
	}

	uniquer := Uniq(func(p Person) string {
		return p.Name
	})
	result := uniquer(people)

	assert.Equal(t, 3, len(result))
	assert.Equal(t, "Alice", result[0].Name)
	assert.Equal(t, "Bob", result[1].Name)
	assert.Equal(t, "Charlie", result[2].Name)
}
