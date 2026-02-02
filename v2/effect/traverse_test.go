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

package effect

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {
	t.Run("traverses empty array", func(t *testing.T) {
		input := []int{}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("traverses array with single element", func(t *testing.T) {
		input := []int{42}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"42"}, result)
	})

	t.Run("traverses array with multiple elements", func(t *testing.T) {
		input := []int{1, 2, 3, 4, 5}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"1", "2", "3", "4", "5"}, result)
	})

	t.Run("transforms to different type", func(t *testing.T) {
		input := []string{"hello", "world", "test"}
		kleisli := TraverseArray(func(s string) Effect[TestContext, int] {
			return Of[TestContext](len(s))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []int{5, 5, 4}, result)
	})

	t.Run("stops on first error", func(t *testing.T) {
		expectedErr := errors.New("traverse error")
		input := []int{1, 2, 3, 4, 5}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			if x == 3 {
				return Fail[TestContext, string](expectedErr)
			}
			return Of[TestContext](strconv.Itoa(x))
		})

		_, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("handles complex transformations", func(t *testing.T) {
		type User struct {
			ID   int
			Name string
		}

		input := []int{1, 2, 3}
		kleisli := TraverseArray(func(id int) Effect[TestContext, User] {
			return Of[TestContext](User{
				ID:   id,
				Name: fmt.Sprintf("User%d", id),
			})
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Len(t, result, 3)
		assert.Equal(t, 1, result[0].ID)
		assert.Equal(t, "User1", result[0].Name)
		assert.Equal(t, 2, result[1].ID)
		assert.Equal(t, "User2", result[1].Name)
		assert.Equal(t, 3, result[2].ID)
		assert.Equal(t, "User3", result[2].Name)
	})

	t.Run("chains with other operations", func(t *testing.T) {
		input := []int{1, 2, 3}

		eff := Chain(func(strings []string) Effect[TestContext, int] {
			total := 0
			for _, s := range strings {
				val, _ := strconv.Atoi(s)
				total += val
			}
			return Of[TestContext](total)
		})(TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x * 2))
		})(input))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 12, result) // (1*2) + (2*2) + (3*2) = 2 + 4 + 6 = 12
	})

	t.Run("uses context in transformation", func(t *testing.T) {
		input := []int{1, 2, 3}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			return Chain(func(ctx TestContext) Effect[TestContext, string] {
				return Of[TestContext](fmt.Sprintf("%s-%d", ctx.Value, x))
			})(Of[TestContext](TestContext{Value: "prefix"}))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []string{"prefix-1", "prefix-2", "prefix-3"}, result)
	})

	t.Run("preserves order", func(t *testing.T) {
		input := []int{5, 3, 8, 1, 9, 2}
		kleisli := TraverseArray(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 10)
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []int{50, 30, 80, 10, 90, 20}, result)
	})

	t.Run("handles large arrays", func(t *testing.T) {
		size := 1000
		input := make([]int, size)
		for i := 0; i < size; i++ {
			input[i] = i
		}

		kleisli := TraverseArray(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Len(t, result, size)
		assert.Equal(t, 0, result[0])
		assert.Equal(t, 1998, result[999])
	})

	t.Run("composes multiple traversals", func(t *testing.T) {
		input := []int{1, 2, 3}

		// First traversal: int -> string
		kleisli1 := TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})

		// Second traversal: string -> int (length)
		kleisli2 := TraverseArray(func(s string) Effect[TestContext, int] {
			return Of[TestContext](len(s))
		})

		eff := Chain(kleisli2)(kleisli1(input))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, []int{1, 1, 1}, result) // All single-digit numbers have length 1
	})

	t.Run("handles nil array", func(t *testing.T) {
		var input []int
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})

		result, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Empty(t, result) // TraverseArray returns empty slice for nil input
	})

	t.Run("works with Map for post-processing", func(t *testing.T) {
		input := []int{1, 2, 3}

		eff := Map[TestContext](func(strings []string) string {
			result := ""
			for _, s := range strings {
				result += s + ","
			}
			return result
		})(TraverseArray(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x))
		})(input))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "1,2,3,", result)
	})

	t.Run("error in middle of array", func(t *testing.T) {
		expectedErr := errors.New("middle error")
		input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			if x == 5 {
				return Fail[TestContext, string](expectedErr)
			}
			return Of[TestContext](strconv.Itoa(x))
		})

		_, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("error at end of array", func(t *testing.T) {
		expectedErr := errors.New("end error")
		input := []int{1, 2, 3, 4, 5}
		kleisli := TraverseArray(func(x int) Effect[TestContext, string] {
			if x == 5 {
				return Fail[TestContext, string](expectedErr)
			}
			return Of[TestContext](strconv.Itoa(x))
		})

		_, err := runEffect(kleisli(input), TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}
