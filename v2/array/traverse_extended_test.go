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
	"strconv"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestMonadTraverse tests the MonadTraverse function
func TestMonadTraverse(t *testing.T) {
	// Test converting integers to strings via Option
	numbers := []int{1, 2, 3}

	result := MonadTraverse(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		numbers,
		func(n int) O.Option[string] {
			return O.Some(strconv.Itoa(n))
		},
	)

	assert.True(t, O.IsSome(result))
	assert.Equal(t, []string{"1", "2", "3"}, O.GetOrElse(func() []string { return []string{} })(result))

	// Test with a function that can return None
	result2 := MonadTraverse(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		numbers,
		func(n int) O.Option[string] {
			if n == 2 {
				return O.None[string]()
			}
			return O.Some(strconv.Itoa(n))
		},
	)

	assert.True(t, O.IsNone(result2))

	// Test with empty array
	empty := []int{}
	result3 := MonadTraverse(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		empty,
		func(n int) O.Option[string] {
			return O.Some(strconv.Itoa(n))
		},
	)

	assert.True(t, O.IsSome(result3))
	assert.Equal(t, []string{}, O.GetOrElse(func() []string { return nil })(result3))
}

// TestTraverseWithIndex tests the TraverseWithIndex function
func TestTraverseWithIndex(t *testing.T) {
	// Test with index-aware transformation
	words := []string{"a", "b", "c"}

	traverser := TraverseWithIndex(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		func(idx int, s string) O.Option[string] {
			return O.Some(s + strconv.Itoa(idx))
		},
	)

	result := traverser(words)
	assert.True(t, O.IsSome(result))
	assert.Equal(t, []string{"a0", "b1", "c2"}, O.GetOrElse(func() []string { return []string{} })(result))

	// Test with conditional None based on index
	traverser2 := TraverseWithIndex(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		func(idx int, s string) O.Option[string] {
			if idx == 1 {
				return O.None[string]()
			}
			return O.Some(s)
		},
	)

	result2 := traverser2(words)
	assert.True(t, O.IsNone(result2))
}

// TestMonadTraverseWithIndex tests the MonadTraverseWithIndex function
func TestMonadTraverseWithIndex(t *testing.T) {
	// Test with index-aware transformation
	numbers := []int{10, 20, 30}

	result := MonadTraverseWithIndex(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		numbers,
		func(idx, n int) O.Option[string] {
			return O.Some(strconv.Itoa(n * idx))
		},
	)

	assert.True(t, O.IsSome(result))
	// Expected: [10*0, 20*1, 30*2] = ["0", "20", "60"]
	assert.Equal(t, []string{"0", "20", "60"}, O.GetOrElse(func() []string { return []string{} })(result))

	// Test with None at specific index
	result2 := MonadTraverseWithIndex(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
		numbers,
		func(idx, n int) O.Option[string] {
			if idx == 2 {
				return O.None[string]()
			}
			return O.Some(strconv.Itoa(n))
		},
	)

	assert.True(t, O.IsNone(result2))
}

// TestMakeTraverseType tests the MakeTraverseType function
func TestMakeTraverseType(t *testing.T) {
	// Create a traverse type for Option
	traverseType := MakeTraverseType[int, string, O.Option[string], O.Option[[]string], O.Option[func(string) []string]]()

	// Use it to traverse an array
	numbers := []int{1, 2, 3}
	result := traverseType(
		O.Of[[]string],
		O.Map[[]string, func(string) []string],
		O.Ap[[]string, string],
	)(func(n int) O.Option[string] {
		return O.Some(strconv.Itoa(n * 2))
	})(numbers)

	assert.True(t, O.IsSome(result))
	assert.Equal(t, []string{"2", "4", "6"}, O.GetOrElse(func() []string { return []string{} })(result))
}
