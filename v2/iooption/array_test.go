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

package iooption

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray_Success(t *testing.T) {
	f := func(n int) IOOption[int] {
		return Of(n * 2)
	}

	input := []int{1, 2, 3, 4, 5}
	result := TraverseArray(f)(input)()

	assert.Equal(t, O.Some([]int{2, 4, 6, 8, 10}), result)
}

func TestTraverseArray_WithNone(t *testing.T) {
	f := func(n int) IOOption[int] {
		if n > 0 {
			return Of(n * 2)
		}
		return None[int]()
	}

	input := []int{1, 2, -3, 4}
	result := TraverseArray(f)(input)()

	assert.Equal(t, O.None[[]int](), result)
}

func TestTraverseArray_EmptyArray(t *testing.T) {
	f := func(n int) IOOption[int] {
		return Of(n * 2)
	}

	input := []int{}
	result := TraverseArray(f)(input)()

	assert.Equal(t, O.Some([]int{}), result)
}

func TestTraverseArrayWithIndex_Success(t *testing.T) {
	f := func(idx, n int) IOOption[int] {
		return Of(n + idx)
	}

	input := []int{10, 20, 30}
	result := TraverseArrayWithIndex(f)(input)()

	assert.Equal(t, O.Some([]int{10, 21, 32}), result)
}

func TestTraverseArrayWithIndex_WithNone(t *testing.T) {
	f := func(idx, n int) IOOption[int] {
		if idx < 2 {
			return Of(n + idx)
		}
		return None[int]()
	}

	input := []int{10, 20, 30}
	result := TraverseArrayWithIndex(f)(input)()

	assert.Equal(t, O.None[[]int](), result)
}

func TestTraverseArrayWithIndex_EmptyArray(t *testing.T) {
	f := func(idx, n int) IOOption[int] {
		return Of(n + idx)
	}

	input := []int{}
	result := TraverseArrayWithIndex(f)(input)()

	assert.Equal(t, O.Some([]int{}), result)
}

func TestSequenceArray_AllSome(t *testing.T) {
	input := []IOOption[int]{
		Of(1),
		Of(2),
		Of(3),
	}

	result := SequenceArray(input)()

	assert.Equal(t, O.Some([]int{1, 2, 3}), result)
}

func TestSequenceArray_WithNone(t *testing.T) {
	input := []IOOption[int]{
		Of(1),
		None[int](),
		Of(3),
	}

	result := SequenceArray(input)()

	assert.Equal(t, O.None[[]int](), result)
}

func TestSequenceArray_Empty(t *testing.T) {
	input := []IOOption[int]{}

	result := SequenceArray(input)()

	assert.Equal(t, O.Some([]int{}), result)
}

func TestSequenceArray_AllNone(t *testing.T) {
	input := []IOOption[int]{
		None[int](),
		None[int](),
		None[int](),
	}

	result := SequenceArray(input)()

	assert.Equal(t, O.None[[]int](), result)
}

func TestTraverseArray_Composition(t *testing.T) {
	// Test composing traverse with other operations
	f := func(n int) IOOption[int] {
		if n%2 == 0 {
			return Of(n / 2)
		}
		return None[int]()
	}

	input := []int{2, 4, 6, 8}
	result := F.Pipe1(
		input,
		TraverseArray(f),
	)()

	assert.Equal(t, O.Some([]int{1, 2, 3, 4}), result)
}

func TestTraverseArray_WithMap(t *testing.T) {
	// Test traverse followed by map
	f := func(n int) IOOption[int] {
		return Of(n * 2)
	}

	input := []int{1, 2, 3}
	result := F.Pipe2(
		input,
		TraverseArray(f),
		Map(func(arr []int) int {
			sum := 0
			for _, v := range arr {
				sum += v
			}
			return sum
		}),
	)()

	assert.Equal(t, O.Some(12), result) // (1*2 + 2*2 + 3*2) = 12
}

func TestTraverseArrayWithIndex_UseIndex(t *testing.T) {
	// Test that index is properly used
	f := func(idx, n int) IOOption[string] {
		return Of(fmt.Sprintf("%d", idx*n*2))
	}

	input := []int{1, 2, 3}
	result := TraverseArrayWithIndex(f)(input)()

	assert.Equal(t, O.Some([]string{"0", "4", "12"}), result)
}


