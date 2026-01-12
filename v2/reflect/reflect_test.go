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

package reflect

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestReduceWithIndex_IntSum tests reducing integers with index awareness
func TestReduceWithIndex_IntSum(t *testing.T) {
	input := []int{10, 20, 30}
	reflectVal := reflect.ValueOf(input)

	// Sum values plus their indices: (0+10) + (1+20) + (2+30) = 63
	reducer := ReduceWithIndex(func(i int, acc int, v reflect.Value) int {
		return acc + i + int(v.Int())
	}, 0)

	result := reducer(reflectVal)
	assert.Equal(t, 63, result)
}

// TestReduceWithIndex_StringConcat tests concatenating strings with indices
func TestReduceWithIndex_StringConcat(t *testing.T) {
	input := []string{"a", "b", "c"}
	reflectVal := reflect.ValueOf(input)

	// Concatenate with indices: "0:a,1:b,2:c"
	reducer := ReduceWithIndex(func(i int, acc string, v reflect.Value) string {
		if acc == "" {
			return string(rune('0'+i)) + ":" + v.String()
		}
		return acc + "," + string(rune('0'+i)) + ":" + v.String()
	}, "")

	result := reducer(reflectVal)
	assert.Equal(t, "0:a,1:b,2:c", result)
}

// TestReduceWithIndex_EmptySlice tests reducing an empty slice
func TestReduceWithIndex_EmptySlice(t *testing.T) {
	input := []int{}
	reflectVal := reflect.ValueOf(input)

	reducer := ReduceWithIndex(func(i int, acc int, v reflect.Value) int {
		return acc + int(v.Int())
	}, 42)

	result := reducer(reflectVal)
	assert.Equal(t, 42, result, "Should return initial value for empty slice")
}

// TestReduceWithIndex_SingleElement tests reducing a single-element slice
func TestReduceWithIndex_SingleElement(t *testing.T) {
	input := []int{100}
	reflectVal := reflect.ValueOf(input)

	reducer := ReduceWithIndex(func(i int, acc int, v reflect.Value) int {
		return acc + i + int(v.Int())
	}, 0)

	result := reducer(reflectVal)
	assert.Equal(t, 100, result, "Should process single element correctly")
}

// TestReduceWithIndex_BuildStruct tests building a complex structure
func TestReduceWithIndex_BuildStruct(t *testing.T) {
	type Result struct {
		Sum   int
		Count int
	}

	input := []int{5, 10, 15}
	reflectVal := reflect.ValueOf(input)

	reducer := ReduceWithIndex(func(i int, acc Result, v reflect.Value) Result {
		return Result{
			Sum:   acc.Sum + int(v.Int()),
			Count: acc.Count + 1,
		}
	}, Result{Sum: 0, Count: 0})

	result := reducer(reflectVal)
	assert.Equal(t, 30, result.Sum)
	assert.Equal(t, 3, result.Count)
}

// TestReduce_IntSum tests basic integer summation
func TestReduce_IntSum(t *testing.T) {
	input := []int{10, 20, 30}
	reflectVal := reflect.ValueOf(input)

	reducer := Reduce(func(acc int, v reflect.Value) int {
		return acc + int(v.Int())
	}, 0)

	result := reducer(reflectVal)
	assert.Equal(t, 60, result)
}

// TestReduce_IntProduct tests integer multiplication
func TestReduce_IntProduct(t *testing.T) {
	input := []int{2, 3, 4}
	reflectVal := reflect.ValueOf(input)

	reducer := Reduce(func(acc int, v reflect.Value) int {
		return acc * int(v.Int())
	}, 1)

	result := reducer(reflectVal)
	assert.Equal(t, 24, result)
}

// TestReduce_StringConcat tests string concatenation
func TestReduce_StringConcat(t *testing.T) {
	input := []string{"Hello", " ", "World"}
	reflectVal := reflect.ValueOf(input)

	reducer := Reduce(func(acc string, v reflect.Value) string {
		return acc + v.String()
	}, "")

	result := reducer(reflectVal)
	assert.Equal(t, "Hello World", result)
}

// TestReduce_EmptySlice tests reducing an empty slice
func TestReduce_EmptySlice(t *testing.T) {
	input := []int{}
	reflectVal := reflect.ValueOf(input)

	reducer := Reduce(func(acc int, v reflect.Value) int {
		return acc + int(v.Int())
	}, 100)

	result := reducer(reflectVal)
	assert.Equal(t, 100, result, "Should return initial value for empty slice")
}

// TestReduce_FindMax tests finding maximum value
func TestReduce_FindMax(t *testing.T) {
	input := []int{3, 7, 2, 9, 1, 5}
	reflectVal := reflect.ValueOf(input)

	reducer := Reduce(func(acc int, v reflect.Value) int {
		val := int(v.Int())
		if val > acc {
			return val
		}
		return acc
	}, input[0])

	result := reducer(reflectVal)
	assert.Equal(t, 9, result)
}

// TestReduce_CountElements tests counting elements matching a condition
func TestReduce_CountElements(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	reflectVal := reflect.ValueOf(input)

	// Count even numbers
	reducer := Reduce(func(acc int, v reflect.Value) int {
		if int(v.Int())%2 == 0 {
			return acc + 1
		}
		return acc
	}, 0)

	result := reducer(reflectVal)
	assert.Equal(t, 5, result, "Should count 5 even numbers")
}

// TestMap_IntToString tests mapping integers to strings
func TestMap_IntToString(t *testing.T) {
	input := []int{1, 2, 3}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) string {
		return "num:" + string(rune('0'+int(v.Int())))
	})

	result := mapper(reflectVal)
	expected := []string{"num:1", "num:2", "num:3"}
	assert.Equal(t, expected, result)
}

// TestMap_DoubleInts tests doubling integer values
func TestMap_DoubleInts(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) int {
		return int(v.Int()) * 2
	})

	result := mapper(reflectVal)
	expected := []int{2, 4, 6, 8, 10}
	assert.Equal(t, expected, result)
}

// TestMap_ExtractField tests extracting a field from structs
func TestMap_ExtractField(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	input := []Person{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) string {
		return v.FieldByName("Name").String()
	})

	result := mapper(reflectVal)
	expected := []string{"Alice", "Bob", "Charlie"}
	assert.Equal(t, expected, result)
}

// TestMap_EmptySlice tests mapping an empty slice
func TestMap_EmptySlice(t *testing.T) {
	input := []int{}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) int {
		return int(v.Int()) * 2
	})

	result := mapper(reflectVal)
	assert.Empty(t, result, "Should return empty slice")
	assert.NotNil(t, result, "Should not return nil")
}

// TestMap_SingleElement tests mapping a single-element slice
func TestMap_SingleElement(t *testing.T) {
	input := []int{42}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) int {
		return int(v.Int()) * 2
	})

	result := mapper(reflectVal)
	expected := []int{84}
	assert.Equal(t, expected, result)
}

// TestMap_BoolToInt tests mapping booleans to integers
func TestMap_BoolToInt(t *testing.T) {
	input := []bool{true, false, true, true, false}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) int {
		if v.Bool() {
			return 1
		}
		return 0
	})

	result := mapper(reflectVal)
	expected := []int{1, 0, 1, 1, 0}
	assert.Equal(t, expected, result)
}

// TestMap_ComplexTransformation tests a complex transformation
func TestMap_ComplexTransformation(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	reflectVal := reflect.ValueOf(input)

	type Result struct {
		Original int
		Squared  int
		IsEven   bool
	}

	mapper := Map(func(v reflect.Value) Result {
		val := int(v.Int())
		return Result{
			Original: val,
			Squared:  val * val,
			IsEven:   val%2 == 0,
		}
	})

	result := mapper(reflectVal)
	assert.Len(t, result, 5)
	assert.Equal(t, 1, result[0].Original)
	assert.Equal(t, 1, result[0].Squared)
	assert.False(t, result[0].IsEven)
	assert.Equal(t, 4, result[3].Original)
	assert.Equal(t, 16, result[3].Squared)
	assert.True(t, result[3].IsEven)
}

// TestMap_StringLength tests mapping strings to their lengths
func TestMap_StringLength(t *testing.T) {
	input := []string{"a", "ab", "abc", "abcd"}
	reflectVal := reflect.ValueOf(input)

	mapper := Map(func(v reflect.Value) int {
		return len(v.String())
	})

	result := mapper(reflectVal)
	expected := []int{1, 2, 3, 4}
	assert.Equal(t, expected, result)
}

// TestIntegration_MapThenReduce tests combining Map and Reduce operations
func TestIntegration_MapThenReduce(t *testing.T) {
	input := []int{1, 2, 3, 4, 5}
	reflectVal := reflect.ValueOf(input)

	// First map: square each number
	mapper := Map(func(v reflect.Value) int {
		val := int(v.Int())
		return val * val
	})
	squared := mapper(reflectVal)

	// Then reduce: sum all squared values
	squaredReflect := reflect.ValueOf(squared)
	reducer := Reduce(func(acc int, v reflect.Value) int {
		return acc + int(v.Int())
	}, 0)
	result := reducer(squaredReflect)

	// 1^2 + 2^2 + 3^2 + 4^2 + 5^2 = 1 + 4 + 9 + 16 + 25 = 55
	assert.Equal(t, 55, result)
}

// TestIntegration_ReduceWithIndexToMap tests using ReduceWithIndex to build a map
func TestIntegration_ReduceWithIndexToMap(t *testing.T) {
	input := []string{"apple", "banana", "cherry"}
	reflectVal := reflect.ValueOf(input)

	// Build a map with index as key
	reducer := ReduceWithIndex(func(i int, acc map[int]string, v reflect.Value) map[int]string {
		acc[i] = v.String()
		return acc
	}, make(map[int]string))

	result := reducer(reflectVal)
	expected := map[int]string{
		0: "apple",
		1: "banana",
		2: "cherry",
	}
	assert.Equal(t, expected, result)
}
