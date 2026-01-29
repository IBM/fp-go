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

	"github.com/stretchr/testify/assert"
)

// TestMonadReduceWithIndex tests the MonadReduceWithIndex function
func TestMonadReduceWithIndex(t *testing.T) {
	// Test with integers - sum with index multiplication
	numbers := []int{1, 2, 3, 4, 5}
	result := MonadReduceWithIndex(numbers, func(idx, acc, val int) int {
		return acc + (val * idx)
	}, 0)
	// Expected: 0*1 + 1*2 + 2*3 + 3*4 + 4*5 = 0 + 2 + 6 + 12 + 20 = 40
	assert.Equal(t, 40, result)

	// Test with empty array
	empty := []int{}
	result2 := MonadReduceWithIndex(empty, func(idx, acc, val int) int {
		return acc + val
	}, 10)
	assert.Equal(t, 10, result2)

	// Test with strings - concatenate with index
	words := []string{"a", "b", "c"}
	result3 := MonadReduceWithIndex(words, func(idx int, acc, val string) string {
		return acc + val + string(rune('0'+idx))
	}, "")
	assert.Equal(t, "a0b1c2", result3)
}

// TestAppend tests the Append function
func TestAppend(t *testing.T) {
	// Test appending to non-empty array
	arr := []int{1, 2, 3}
	result := Append(arr, 4)
	assert.Equal(t, []int{1, 2, 3, 4}, result)
	// Verify original array is unchanged
	assert.Equal(t, []int{1, 2, 3}, arr)

	// Test appending to empty array
	empty := []int{}
	result2 := Append(empty, 1)
	assert.Equal(t, []int{1}, result2)

	// Test appending strings
	words := []string{"hello", "world"}
	result3 := Append(words, "!")
	assert.Equal(t, []string{"hello", "world", "!"}, result3)

	// Test appending to nil array
	var nilArr []int
	result4 := Append(nilArr, 42)
	assert.Equal(t, []int{42}, result4)
}

// TestStrictEquals tests the StrictEquals function
func TestStrictEquals(t *testing.T) {
	eq := StrictEquals[int]()

	// Test equal arrays
	arr1 := []int{1, 2, 3}
	arr2 := []int{1, 2, 3}
	assert.True(t, eq.Equals(arr1, arr2))

	// Test different arrays
	arr3 := []int{1, 2, 4}
	assert.False(t, eq.Equals(arr1, arr3))

	// Test different lengths
	arr4 := []int{1, 2}
	assert.False(t, eq.Equals(arr1, arr4))

	// Test empty arrays
	empty1 := []int{}
	empty2 := []int{}
	assert.True(t, eq.Equals(empty1, empty2))

	// Test with strings
	strEq := StrictEquals[string]()
	words1 := []string{"hello", "world"}
	words2 := []string{"hello", "world"}
	words3 := []string{"hello", "there"}
	assert.True(t, strEq.Equals(words1, words2))
	assert.False(t, strEq.Equals(words1, words3))
}
