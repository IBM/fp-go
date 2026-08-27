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
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestAt_Success verifies that At returns Some for valid in-bounds indices.
func TestAt_Success(t *testing.T) {
	data := []int{10, 20, 30, 40, 50}

	t.Run("first element", func(t *testing.T) {
		assert.Equal(t, O.Some(10), At[int](0)(data))
	})

	t.Run("middle element", func(t *testing.T) {
		assert.Equal(t, O.Some(30), At[int](2)(data))
	})

	t.Run("last element", func(t *testing.T) {
		assert.Equal(t, O.Some(50), At[int](4)(data))
	})

	t.Run("single-element slice", func(t *testing.T) {
		assert.Equal(t, O.Some(42), At[int](0)([]int{42}))
	})
}

// TestAt_Failure verifies that At returns None for out-of-bounds and negative indices.
func TestAt_Failure(t *testing.T) {
	data := []int{10, 20, 30}

	t.Run("index equals length", func(t *testing.T) {
		assert.Equal(t, O.None[int](), At[int](3)(data))
	})

	t.Run("index beyond length", func(t *testing.T) {
		assert.Equal(t, O.None[int](), At[int](100)(data))
	})

	t.Run("negative index", func(t *testing.T) {
		assert.Equal(t, O.None[int](), At[int](-1)(data))
	})

	t.Run("empty slice", func(t *testing.T) {
		assert.Equal(t, O.None[int](), At[int](0)([]int{}))
	})

	t.Run("nil slice", func(t *testing.T) {
		assert.Equal(t, O.None[int](), At[int](0)(nil))
	})
}

// TestAt_EdgeCases verifies boundary conditions for At.
func TestAt_EdgeCases(t *testing.T) {
	data := []string{"a", "b", "c"}

	t.Run("index zero on non-empty slice", func(t *testing.T) {
		assert.Equal(t, O.Some("a"), At[string](0)(data))
	})

	t.Run("index len-1 (last valid)", func(t *testing.T) {
		assert.Equal(t, O.Some("c"), At[string](2)(data))
	})

	t.Run("index len (first invalid)", func(t *testing.T) {
		assert.Equal(t, O.None[string](), At[string](3)(data))
	})
}

// TestAt_Integration verifies At composes correctly with other combinators.
func TestAt_Integration(t *testing.T) {
	data := []int{1, 4, 9, 16, 25}

	t.Run("compose with option.Map", func(t *testing.T) {
		// Retrieve element and double it if present.
		result := F.Pipe1(At[int](2)(data), O.Map(func(x int) int { return x * 2 }))
		assert.Equal(t, O.Some(18), result)
	})

	t.Run("compose with option.GetOrElse", func(t *testing.T) {
		present := F.Pipe1(At[int](1)(data), O.GetOrElse(func() int { return 0 }))
		absent := F.Pipe1(At[int](99)(data), O.GetOrElse(func() int { return 0 }))
		assert.Equal(t, 4, present)
		assert.Equal(t, 0, absent)
	})

	t.Run("use as kleisli in chain", func(t *testing.T) {
		// Use At as a FilterMap callback to extract the second element of each sub-slice.
		nested := [][]int{{1, 2, 3}, {4, 5}, {6}}
		result := FilterMap[[]int, int](At[int](1))(nested)
		assert.Equal(t, []int{2, 5}, result)
	})
}

// ExampleAt_found demonstrates retrieving an element that exists.
func ExampleAt_found() {
	result := At[string](1)([]string{"foo", "bar", "baz"})
	fmt.Println(result)
	// Output:
	// Some[string](bar)
}

// ExampleAt_notFound demonstrates that At returns None for an out-of-bounds index.
func ExampleAt_notFound() {
	result := At[string](5)([]string{"foo", "bar", "baz"})
	fmt.Println(O.IsNone(result))
	// Output:
	// true
}

// ExampleAt_negative demonstrates that a negative index yields None.
func ExampleAt_negative() {
	result := At[int](-1)([]int{10, 20, 30})
	fmt.Println(O.IsNone(result))
	// Output:
	// true
}
