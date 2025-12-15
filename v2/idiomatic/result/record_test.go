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

// TestTraverseRecordG_Success tests successful traversal of a map
func TestTraverseRecordG_Success(t *testing.T) {
	parse := strconv.Atoi

	input := map[string]string{"a": "1", "b": "2", "c": "3"}
	result, err := TraverseRecordG[map[string]string, map[string]int](parse)(input)

	require.NoError(t, err)
	assert.Equal(t, 1, result["a"])
	assert.Equal(t, 2, result["b"])
	assert.Equal(t, 3, result["c"])
}

// TestTraverseRecordG_Error tests that traversal short-circuits on error
func TestTraverseRecordG_Error(t *testing.T) {
	parse := strconv.Atoi

	input := map[string]string{"a": "1", "b": "bad", "c": "3"}
	result, err := TraverseRecordG[map[string]string, map[string]int](parse)(input)

	require.Error(t, err)
	assert.Nil(t, result)
}

// TestTraverseRecordG_EmptyMap tests traversal of an empty map
func TestTraverseRecordG_EmptyMap(t *testing.T) {
	parse := strconv.Atoi

	input := map[string]string{}
	result, err := TraverseRecordG[map[string]string, map[string]int](parse)(input)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result) // Should be an empty map, not nil
}

// TestTraverseRecordG_CustomMapType tests with custom map types
func TestTraverseRecordG_CustomMapType(t *testing.T) {
	type StringMap map[string]string
	type IntMap map[string]int

	parse := strconv.Atoi

	input := StringMap{"x": "10", "y": "20"}
	result, err := TraverseRecordG[StringMap, IntMap](parse)(input)

	require.NoError(t, err)
	assert.Equal(t, IntMap{"x": 10, "y": 20}, result)
}

// TestTraverseRecord_Success tests successful traversal
func TestTraverseRecord_Success(t *testing.T) {
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

	input := map[string]string{"a": "1", "b": "2"}
	result, err := TraverseRecord[string](validate)(input)

	require.NoError(t, err)
	assert.Equal(t, 2, result["a"])
	assert.Equal(t, 4, result["b"])
}

// TestTraverseRecord_ValidationError tests validation failure
func TestTraverseRecord_ValidationError(t *testing.T) {
	validate := func(s string) (int, error) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return 0, err
		}
		if n < 0 {
			return 0, errors.New("negative number")
		}
		return n, nil
	}

	input := map[string]string{"a": "1", "b": "-5"}
	result, err := TraverseRecord[string](validate)(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "negative")
	assert.Nil(t, result)
}

// TestTraverseRecordWithIndexG_Success tests successful indexed traversal
func TestTraverseRecordWithIndexG_Success(t *testing.T) {
	annotate := func(k string, v string) (string, error) {
		if S.IsEmpty(v) {
			return "", fmt.Errorf("empty value for key %s", k)
		}
		return fmt.Sprintf("%s=%s", k, v), nil
	}

	input := map[string]string{"a": "1", "b": "2"}
	result, err := TraverseRecordWithIndexG[map[string]string, map[string]string](annotate)(input)

	require.NoError(t, err)
	assert.Equal(t, "a=1", result["a"])
	assert.Equal(t, "b=2", result["b"])
}

// TestTraverseRecordWithIndexG_Error tests error handling with key
func TestTraverseRecordWithIndexG_Error(t *testing.T) {
	annotate := func(k string, v string) (string, error) {
		if S.IsEmpty(v) {
			return "", fmt.Errorf("empty value for key %s", k)
		}
		return v, nil
	}

	input := map[string]string{"a": "1", "b": ""}
	result, err := TraverseRecordWithIndexG[map[string]string, map[string]string](annotate)(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "key b")
	assert.Nil(t, result)
}

// TestTraverseRecordWithIndexG_EmptyMap tests empty map
func TestTraverseRecordWithIndexG_EmptyMap(t *testing.T) {
	annotate := func(k string, v string) (string, error) {
		return k + ":" + v, nil
	}

	input := map[string]string{}
	result, err := TraverseRecordWithIndexG[map[string]string, map[string]string](annotate)(input)

	require.NoError(t, err)
	assert.Empty(t, result)
	assert.NotNil(t, result)
}

// TestTraverseRecordWithIndex_Success tests successful indexed traversal
func TestTraverseRecordWithIndex_Success(t *testing.T) {
	check := func(k string, v int) (string, error) {
		if v < 0 {
			return "", fmt.Errorf("negative value for key %s", k)
		}
		return fmt.Sprintf("%s:%d", k, v*2), nil
	}

	input := map[string]int{"a": 1, "b": 2}
	result, err := TraverseRecordWithIndex(check)(input)

	require.NoError(t, err)
	assert.Equal(t, "a:2", result["a"])
	assert.Equal(t, "b:4", result["b"])
}

// TestTraverseRecordWithIndex_Error tests error with key info
func TestTraverseRecordWithIndex_Error(t *testing.T) {
	check := func(k string, v int) (int, error) {
		if v < 0 {
			return 0, fmt.Errorf("negative value for key %s", k)
		}
		return v, nil
	}

	input := map[string]int{"ok": 1, "bad": -5}
	result, err := TraverseRecordWithIndex(check)(input)

	require.Error(t, err)
	assert.Contains(t, err.Error(), "key bad")
	assert.Nil(t, result)
}

// TestTraverseRecordWithIndex_TypeTransformation tests transforming types with key
func TestTraverseRecordWithIndex_TypeTransformation(t *testing.T) {
	prefixKey := func(k string, v string) (string, error) {
		return k + "_" + v, nil
	}

	input := map[string]string{"prefix": "value", "another": "test"}
	result, err := TraverseRecordWithIndex(prefixKey)(input)

	require.NoError(t, err)
	assert.Equal(t, "prefix_value", result["prefix"])
	assert.Equal(t, "another_test", result["another"])
}

// TestTraverseRecord_IntKeys tests with integer keys
func TestTraverseRecord_IntKeys(t *testing.T) {
	double := func(n int) (int, error) {
		return n * 2, nil
	}

	input := map[int]int{1: 10, 2: 20, 3: 30}
	result, err := TraverseRecord[int](double)(input)

	require.NoError(t, err)
	assert.Equal(t, 20, result[1])
	assert.Equal(t, 40, result[2])
	assert.Equal(t, 60, result[3])
}

// TestTraverseRecordG_PreservesKeys tests that keys are preserved
func TestTraverseRecordG_PreservesKeys(t *testing.T) {
	identity := func(s string) (string, error) {
		return s, nil
	}

	input := map[string]string{"key1": "val1", "key2": "val2"}
	result, err := TraverseRecordG[map[string]string, map[string]string](identity)(input)

	require.NoError(t, err)
	assert.Contains(t, result, "key1")
	assert.Contains(t, result, "key2")
	assert.Equal(t, "val1", result["key1"])
	assert.Equal(t, "val2", result["key2"])
}
