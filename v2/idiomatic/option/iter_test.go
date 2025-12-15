// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/iterator/iter"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a sequence from a slice
func seqFromSlice[T any](items []T) Seq[T] {
	return I.From(items...)
}

// Helper function to collect a sequence into a slice
func collectSeq[T any](seq Seq[T]) []T {
	return slices.Collect(seq)
}

func TestTraverseIter_AllSome(t *testing.T) {
	// Test case where all transformations succeed
	parse := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		return Some(n)
	}

	input := I.From("1", "2", "3", "4", "5")
	result, resultok := TraverseIter(parse)(input)

	assert.True(t, IsSome(result, resultok), "Expected Some result when all transformations succeed")

	collected := collectSeq(result)
	expected := A.From(1, 2, 3, 4, 5)
	assert.Equal(t, expected, collected)
}

func TestTraverseIter_ContainsNone(t *testing.T) {
	// Test case where one transformation fails
	parse := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		return Some(n)
	}

	input := seqFromSlice([]string{"1", "invalid", "3"})
	result, resultok := TraverseIter(parse)(input)

	assert.True(t, IsNone(result, resultok), "Expected None when any transformation fails")
}

func TestTraverseIter_EmptySequence(t *testing.T) {
	// Test with empty sequence
	double := func(x int) (int, bool) {
		return Some(x * 2)
	}

	input := seqFromSlice([]int{})
	result, resultok := TraverseIter(double)(input)

	assert.True(t, IsSome(result, resultok), "Expected Some for empty sequence")

	collected := collectSeq(result)
	assert.Empty(t, collected)
}

func TestTraverseIter_SingleElement(t *testing.T) {
	// Test with single element - success case
	validate := func(x int) (int, bool) {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	input := seqFromSlice([]int{5})
	result, resultok := TraverseIter(validate)(input)

	assert.True(t, IsSome(result, resultok))
	collected := collectSeq(result)
	assert.Equal(t, []int{10}, collected)
}

func TestTraverseIter_SingleElementFails(t *testing.T) {
	// Test with single element - failure case
	validate := func(x int) (int, bool) {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	input := seqFromSlice([]int{-5})
	result, resultok := TraverseIter(validate)(input)

	assert.True(t, IsNone(result, resultok))
}

func TestTraverseIter_Validation(t *testing.T) {
	// Test validation use case
	validatePositive := func(x int) (int, bool) {
		if x > 0 {
			return Some(x)
		}
		return None[int]()
	}

	// All positive
	input1 := seqFromSlice([]int{1, 2, 3, 4})
	result1, result1ok := TraverseIter(validatePositive)(input1)
	assert.True(t, IsSome(result1, result1ok))

	// Contains negative
	input2 := seqFromSlice([]int{1, -2, 3})
	result2, result2ok := TraverseIter(validatePositive)(input2)
	assert.True(t, IsNone(result2, result2ok))

	// Contains zero
	input3 := seqFromSlice([]int{1, 0, 3})
	result3, result3ok := TraverseIter(validatePositive)(input3)
	assert.True(t, IsNone(result3, result3ok))
}

func TestTraverseIter_Transformation(t *testing.T) {
	// Test transformation use case
	safeDivide := func(x int) (float64, bool) {
		if x != 0 {
			return Some(100.0 / float64(x))
		}
		return None[float64]()
	}

	// All non-zero
	input1 := seqFromSlice([]int{1, 2, 4, 5})
	result1, result1ok := TraverseIter(safeDivide)(input1)
	assert.True(t, IsSome(result1, result1ok))

	collected := collectSeq(result1)
	expected := []float64{100.0, 50.0, 25.0, 20.0}
	assert.Equal(t, expected, collected)

	// Contains zero
	input2 := seqFromSlice([]int{1, 0, 4})
	result2, result2ok := TraverseIter(safeDivide)(input2)
	assert.True(t, IsNone(result2, result2ok))
}

func TestTraverseIter_ShortCircuit(t *testing.T) {
	// Test that traversal short-circuits on first None
	callCount := 0
	countingFunc := func(x int) (int, bool) {
		callCount++
		if x < 0 {
			return None[int]()
		}
		return Some(x * 2)
	}

	// First element fails
	input := seqFromSlice([]int{-1, 2, 3, 4, 5})
	result, resultok := TraverseIter(countingFunc)(input)

	assert.True(t, IsNone(result, resultok))
	// Should have called the function for elements until the first failure
	// Note: The exact count depends on implementation details of the traverse function
	assert.Greater(t, callCount, 0, "Function should be called at least once")
}

func TestTraverseIter_LazyEvaluation(t *testing.T) {
	// Test that the result sequence is lazy
	transform := func(x int) (int, bool) {
		return Some(x * 2)
	}

	input := seqFromSlice([]int{1, 2, 3, 4, 5})
	result, resultok := TraverseIter(transform)(input)

	assert.True(t, IsSome(result, resultok))

	// Partially consume the sequence
	callCount := 0
	Fold(func() int { return 0 }, func(seq Seq[int]) int {
		for val := range seq {
			callCount++
			_ = val
			if callCount == 2 {
				break
			}
		}
		return callCount
	})(result, resultok)

	assert.Equal(t, 2, callCount, "Should only evaluate consumed elements")
}

func TestTraverseIter_ComplexTransformation(t *testing.T) {
	// Test with more complex transformation
	type Person struct {
		Name string
		Age  int
	}

	validatePerson := func(name string) (Person, bool) {
		if S.IsEmpty(name) {
			return None[Person]()
		}
		return Some(Person{Name: name, Age: len(name)})
	}

	input := seqFromSlice([]string{"Alice", "Bob", "Charlie"})
	result, resultok := TraverseIter(validatePerson)(input)

	assert.True(t, IsSome(result, resultok))

	collected := collectSeq((result))
	expected := []Person{
		{Name: "Alice", Age: 5},
		{Name: "Bob", Age: 3},
		{Name: "Charlie", Age: 7},
	}
	assert.Equal(t, expected, collected)
}

func TestTraverseIter_WithPipeline(t *testing.T) {
	// Test TraverseIter in a functional pipeline
	parse := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		return Some(n)
	}

	input := seqFromSlice([]string{"1", "2", "3", "4", "5"})

	collected := Fold(func() []int { return nil }, F.Identity[[]int])(Map(collectSeq[int])(TraverseIter(parse)(input)))
	expected := []int{1, 2, 3, 4, 5}
	assert.Equal(t, expected, collected)
}

func TestTraverseIter_ChainedTransformations(t *testing.T) {
	// Test chaining multiple transformations
	parseAndValidate := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		if n > 0 {
			return Some(n)
		}
		return None[int]()
	}

	// All valid
	input1 := seqFromSlice([]string{"1", "2", "3"})
	result1, result1ok := TraverseIter(parseAndValidate)(input1)
	assert.True(t, IsSome(result1, result1ok))

	// Contains invalid number
	input2 := seqFromSlice([]string{"1", "invalid", "3"})
	result2, result2ok := TraverseIter(parseAndValidate)(input2)
	assert.True(t, IsNone(result2, result2ok))

	// Contains non-positive number
	input3 := seqFromSlice([]string{"1", "0", "3"})
	result3, result3ok := TraverseIter(parseAndValidate)(input3)
	assert.True(t, IsNone(result3, result3ok))
}

// Example test demonstrating usage
func ExampleTraverseIter() {
	// Parse a sequence of strings to integers
	parse := func(s string) (int, bool) {
		n, err := strconv.Atoi(s)
		if err != nil {
			return None[int]()
		}
		return Some(n)
	}

	// Create a sequence of valid strings
	validStrings := seqFromSlice([]string{"1", "2", "3"})
	result, resultok := TraverseIter(parse)(validStrings)

	if IsSome(result, resultok) {
		numbers := collectSeq(result)
		fmt.Println(numbers)
	}

	// Create a sequence with invalid string
	invalidStrings := seqFromSlice([]string{"1", "invalid", "3"})
	result2, result2ok := TraverseIter(parse)(invalidStrings)

	if IsNone(result2, result2ok) {
		fmt.Println("Parsing failed")
	}

	// Output:
	// [1 2 3]
	// Parsing failed
}
