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

package result

import (
	"errors"
	"fmt"
	"strconv"
	"testing"

	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTraverseArrayG_Success tests successful traversal of an array with all valid elements
func TestTraverseArrayG_Success(t *testing.T) {
	parse := strconv.Atoi

	result, err := TraverseArrayG[[]string, []int](parse)([]string{"1", "2", "3"})

	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

// TestTraverseArrayG_Error tests that traversal short-circuits on first error
func TestTraverseArrayG_Error(t *testing.T) {
	parse := strconv.Atoi

	result, err := TraverseArrayG[[]string, []int](parse)([]string{"1", "bad", "3"})

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestTraverseArrayG_EmptyArray tests traversal of an empty array
func TestTraverseArrayG_EmptyArray(t *testing.T) {
	parse := strconv.Atoi

	result, err := TraverseArrayG[[]string, []int](parse)([]string{})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result) // Should be an empty slice, not nil
}

// TestTraverseArrayG_SingleElement tests traversal with a single element
func TestTraverseArrayG_SingleElement(t *testing.T) {
	parse := strconv.Atoi

	result, err := TraverseArrayG[[]string, []int](parse)([]string{"42"})

	require.NoError(t, err)
	assert.Equal(t, []int{42}, result)
}

// TestTraverseArrayG_ShortCircuit tests that processing stops at first error
func TestTraverseArrayG_ShortCircuit(t *testing.T) {
	callCount := 0
	parse := func(s string) (int, error) {
		callCount++
		if s == "error" {
			return 0, errors.New("parse error")
		}
		return len(s), nil
	}

	_, err := TraverseArrayG[[]string, []int](parse)([]string{"ok", "error", "should-not-process"})

	require.Error(t, err)
	assert.Equal(t, 2, callCount, "should stop after encountering error")
}

// TestTraverseArrayG_CustomSliceType tests with custom slice types
func TestTraverseArrayG_CustomSliceType(t *testing.T) {
	type StringSlice []string
	type IntSlice []int

	parse := strconv.Atoi

	input := StringSlice{"1", "2", "3"}
	result, err := TraverseArrayG[StringSlice, IntSlice](parse)(input)

	require.NoError(t, err)
	assert.Equal(t, IntSlice{1, 2, 3}, result)
}

// TestTraverseArray_Success tests successful traversal
func TestTraverseArray_Success(t *testing.T) {
	validate := func(s string) (int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if n < 0 {
			return 0, errors.New("negative number")
		}
		return n * 2, nil
	}

	result, err := TraverseArray(validate)([]string{"1", "2", "3"})

	require.NoError(t, err)
	assert.Equal(t, []int{2, 4, 6}, result)
}

// TestTraverseArray_ValidationError tests validation failure
func TestTraverseArray_ValidationError(t *testing.T) {
	validate := func(s string) (int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if n < 0 {
			return 0, errors.New("negative number")
		}
		return n * 2, nil
	}

	result, err := TraverseArray(validate)([]string{"1", "-5", "3"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "negative number")
	assert.Nil(t, result)
}

// TestTraverseArray_ParseError tests parse failure
func TestTraverseArray_ParseError(t *testing.T) {
	validate := func(s string) (int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		return n, nil
	}

	result, err := TraverseArray(validate)([]string{"1", "not-a-number", "3"})

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestTraverseArray_EmptyArray tests empty array
func TestTraverseArray_EmptyArray(t *testing.T) {
	identity := func(n int) (int, error) {
		return n, nil
	}

	result, err := TraverseArray(identity)([]int{})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// TestTraverseArray_DifferentTypes tests transformation between different types
func TestTraverseArray_DifferentTypes(t *testing.T) {
	toLength := func(s string) (int, error) {
		if S.IsEmpty(s) {
			return 0, errors.New("empty string")
		}
		return len(s), nil
	}

	result, err := TraverseArray(toLength)([]string{"a", "bb", "ccc"})

	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3}, result)
}

// TestTraverseArrayWithIndexG_Success tests successful indexed traversal
func TestTraverseArrayWithIndexG_Success(t *testing.T) {
	annotate := func(i int, s string) (string, error) {
		if S.IsEmpty(s) {
			return "", fmt.Errorf("empty string at index %d", i)
		}
		return fmt.Sprintf("[%d]=%s", i, s), nil
	}

	result, err := TraverseArrayWithIndexG[[]string, []string](annotate)([]string{"a", "b", "c"})

	require.NoError(t, err)
	assert.Equal(t, []string{"[0]=a", "[1]=b", "[2]=c"}, result)
}

// TestTraverseArrayWithIndexG_Error tests error handling with index
func TestTraverseArrayWithIndexG_Error(t *testing.T) {
	annotate := func(i int, s string) (string, error) {
		if S.IsEmpty(s) {
			return "", fmt.Errorf("empty string at index %d", i)
		}
		return fmt.Sprintf("[%d]=%s", i, s), nil
	}

	result, err := TraverseArrayWithIndexG[[]string, []string](annotate)([]string{"a", "", "c"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "index 1")
	assert.Nil(t, result)
}

// TestTraverseArrayWithIndexG_EmptyArray tests empty array
func TestTraverseArrayWithIndexG_EmptyArray(t *testing.T) {
	annotate := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d:%s", i, s), nil
	}

	result, err := TraverseArrayWithIndexG[[]string, []string](annotate)([]string{})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// TestTraverseArrayWithIndexG_IndexValidation tests that indices are correct
func TestTraverseArrayWithIndexG_IndexValidation(t *testing.T) {
	indices := []int{}
	collect := func(i int, s string) (string, error) {
		indices = append(indices, i)
		return s, nil
	}

	_, err := TraverseArrayWithIndexG[[]string, []string](collect)([]string{"a", "b", "c", "d"})

	require.NoError(t, err)
	assert.Equal(t, []int{0, 1, 2, 3}, indices)
}

// TestTraverseArrayWithIndexG_ShortCircuit tests short-circuit behavior
func TestTraverseArrayWithIndexG_ShortCircuit(t *testing.T) {
	maxIndex := -1
	process := func(i int, s string) (string, error) {
		maxIndex = i
		if i == 2 {
			return "", errors.New("stop at index 2")
		}
		return s, nil
	}

	_, err := TraverseArrayWithIndexG[[]string, []string](process)([]string{"a", "b", "c", "d", "e"})

	require.Error(t, err)
	assert.Equal(t, 2, maxIndex, "should stop at index 2")
}

// TestTraverseArrayWithIndexG_CustomSliceType tests with custom slice types
func TestTraverseArrayWithIndexG_CustomSliceType(t *testing.T) {
	type StringSlice []string
	type ResultSlice []string

	annotate := func(i int, s string) (string, error) {
		return fmt.Sprintf("%d:%s", i, s), nil
	}

	input := StringSlice{"x", "y", "z"}
	result, err := TraverseArrayWithIndexG[StringSlice, ResultSlice](annotate)(input)

	require.NoError(t, err)
	assert.Equal(t, ResultSlice{"0:x", "1:y", "2:z"}, result)
}

// TestTraverseArrayWithIndex_Success tests successful indexed traversal
func TestTraverseArrayWithIndex_Success(t *testing.T) {
	check := func(i int, s string) (string, error) {
		if S.IsEmpty(s) {
			return "", fmt.Errorf("empty value at position %d", i)
		}
		return fmt.Sprintf("%d_%s", i, s), nil
	}

	result, err := TraverseArrayWithIndex(check)([]string{"a", "b", "c"})

	require.NoError(t, err)
	assert.Equal(t, []string{"0_a", "1_b", "2_c"}, result)
}

// TestTraverseArrayWithIndex_Error tests error with position info
func TestTraverseArrayWithIndex_Error(t *testing.T) {
	check := func(i int, s string) (string, error) {
		if S.IsEmpty(s) {
			return "", fmt.Errorf("empty value at position %d", i)
		}
		return s, nil
	}

	result, err := TraverseArrayWithIndex(check)([]string{"ok", "ok", "", "ok"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "position 2")
	assert.Nil(t, result)
}

// TestTraverseArrayWithIndex_TypeTransformation tests transforming types with index
func TestTraverseArrayWithIndex_TypeTransformation(t *testing.T) {
	multiply := func(i int, n int) (int, error) {
		return n * (i + 1), nil
	}

	result, err := TraverseArrayWithIndex(multiply)([]int{10, 20, 30})

	require.NoError(t, err)
	assert.Equal(t, []int{10, 40, 90}, result) // [10*(0+1), 20*(1+1), 30*(2+1)]
}

// TestTraverseArrayWithIndex_EmptyArray tests empty array
func TestTraverseArrayWithIndex_EmptyArray(t *testing.T) {
	process := func(i int, s string) (int, error) {
		return i, nil
	}

	result, err := TraverseArrayWithIndex(process)([]string{})

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// TestTraverseArrayWithIndex_SingleElement tests single element processing
func TestTraverseArrayWithIndex_SingleElement(t *testing.T) {
	annotate := func(i int, s string) (string, error) {
		return fmt.Sprintf("item_%d:%s", i, s), nil
	}

	result, err := TraverseArrayWithIndex(annotate)([]string{"solo"})

	require.NoError(t, err)
	assert.Equal(t, []string{"item_0:solo"}, result)
}

// TestTraverseArrayWithIndex_ConditionalLogic tests using index for conditional logic
func TestTraverseArrayWithIndex_ConditionalLogic(t *testing.T) {
	// Only accept even indices
	process := func(i int, s string) (string, error) {
		if i%2 != 0 {
			return "", fmt.Errorf("odd index %d not allowed", i)
		}
		return s, nil
	}

	result, err := TraverseArrayWithIndex(process)([]string{"ok", "fail"})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "odd index 1")
	assert.Nil(t, result)
}

// TestTraverseArray_LargeArray tests traversal with a larger array
func TestTraverseArray_LargeArray(t *testing.T) {
	// Create array with 1000 elements
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	double := func(n int) (int, error) {
		return n * 2, nil
	}

	result, err := TraverseArray(double)(input)

	require.NoError(t, err)
	assert.Len(t, result, 1000)
	assert.Equal(t, 0, result[0])
	assert.Equal(t, 1998, result[999])
}

// TestTraverseArrayG_PreservesOrder tests that order is preserved
func TestTraverseArrayG_PreservesOrder(t *testing.T) {
	reverse := func(s string) (string, error) {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes), nil
	}

	result, err := TraverseArrayG[[]string, []string](reverse)([]string{"abc", "def", "ghi"})

	require.NoError(t, err)
	assert.Equal(t, []string{"cba", "fed", "ihg"}, result)
}

// TestTraverseArrayWithIndex_BoundaryCheck tests boundary conditions with index
func TestTraverseArrayWithIndex_BoundaryCheck(t *testing.T) {
	// Reject if index exceeds a threshold
	limitedProcess := func(i int, s string) (string, error) {
		if i >= 100 {
			return "", errors.New("index too large")
		}
		return s, nil
	}

	// Should succeed with index < 100
	result, err := TraverseArrayWithIndex(limitedProcess)([]string{"a", "b", "c"})
	require.NoError(t, err)
	assert.Equal(t, []string{"a", "b", "c"}, result)
}
