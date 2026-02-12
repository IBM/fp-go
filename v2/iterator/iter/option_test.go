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

package iter

import (
	"fmt"
	"slices"
	"strconv"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

// TestMonadChainOptionK_AllSome tests MonadChainOptionK when all values produce Some
func TestMonadChainOptionK_AllSome(t *testing.T) {
	// Function that always returns Some
	double := func(x int) O.Option[int] {
		return O.Some(x * 2)
	}

	seq := From(1, 2, 3, 4, 5)
	result := MonadChainOptionK(seq, double)
	values := slices.Collect(result)

	expected := A.From(2, 4, 6, 8, 10)
	assert.Equal(t, expected, values)
}

// TestMonadChainOptionK_AllNone tests MonadChainOptionK when all values produce None
func TestMonadChainOptionK_AllNone(t *testing.T) {
	// Function that always returns None
	alwaysNone := func(x int) O.Option[int] {
		return O.None[int]()
	}

	seq := From(1, 2, 3, 4, 5)
	result := MonadChainOptionK(seq, alwaysNone)
	values := slices.Collect(result)

	assert.Empty(t, values)
}

// TestMonadChainOptionK_MixedSomeNone tests MonadChainOptionK with mixed Some and None
func TestMonadChainOptionK_MixedSomeNone(t *testing.T) {
	// Function that returns Some for even numbers, None for odd
	evenOnly := func(x int) O.Option[int] {
		if x%2 == 0 {
			return O.Some(x)
		}
		return O.None[int]()
	}

	seq := From(1, 2, 3, 4, 5, 6)
	result := MonadChainOptionK(seq, evenOnly)
	values := slices.Collect(result)

	expected := A.From(2, 4, 6)
	assert.Equal(t, expected, values)
}

// TestMonadChainOptionK_ParseStrings tests parsing strings to integers
func TestMonadChainOptionK_ParseStrings(t *testing.T) {
	// Parse strings to integers, returning None for invalid strings
	parseNum := func(s string) O.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return O.Some(n)
		}
		return O.None[int]()
	}

	seq := From("1", "invalid", "2", "3", "bad", "4")
	result := MonadChainOptionK(seq, parseNum)
	values := slices.Collect(result)

	expected := A.From(1, 2, 3, 4)
	assert.Equal(t, expected, values)
}

// TestMonadChainOptionK_EmptySequence tests MonadChainOptionK with empty sequence
func TestMonadChainOptionK_EmptySequence(t *testing.T) {
	double := func(x int) O.Option[int] {
		return O.Some(x * 2)
	}

	seq := From[int]()
	result := MonadChainOptionK(seq, double)
	values := slices.Collect(result)

	assert.Empty(t, values)
}

// TestMonadChainOptionK_TypeTransformation tests transforming types
func TestMonadChainOptionK_TypeTransformation(t *testing.T) {
	// Convert integers to strings, only for positive numbers
	positiveToString := func(x int) O.Option[string] {
		if x > 0 {
			return O.Some(fmt.Sprintf("num_%d", x))
		}
		return O.None[string]()
	}

	seq := From(-2, -1, 0, 1, 2, 3)
	result := MonadChainOptionK(seq, positiveToString)
	values := slices.Collect(result)

	expected := A.From("num_1", "num_2", "num_3")
	assert.Equal(t, expected, values)
}

// TestMonadChainOptionK_ComplexType tests with complex types
func TestMonadChainOptionK_ComplexType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Extract age only for adults
	getAdultAge := func(p Person) O.Option[int] {
		if p.Age >= 18 {
			return O.Some(p.Age)
		}
		return O.None[int]()
	}

	seq := From(
		Person{"Alice", 25},
		Person{"Bob", 15},
		Person{"Charlie", 30},
		Person{"David", 12},
	)
	result := MonadChainOptionK(seq, getAdultAge)
	values := slices.Collect(result)

	expected := A.From(25, 30)
	assert.Equal(t, expected, values)
}

// TestChainOptionK_BasicUsage tests ChainOptionK basic functionality
func TestChainOptionK_BasicUsage(t *testing.T) {
	// Create a reusable operator
	parsePositive := ChainOptionK(func(x int) O.Option[int] {
		if x > 0 {
			return O.Some(x)
		}
		return O.None[int]()
	})

	seq := From(-1, 2, -3, 4, 5, -6)
	result := parsePositive(seq)
	values := slices.Collect(result)

	expected := A.From(2, 4, 5)
	assert.Equal(t, expected, values)
}

// TestChainOptionK_WithPipe tests ChainOptionK in a pipeline
func TestChainOptionK_WithPipe(t *testing.T) {
	// Validate and transform in a pipeline
	validateRange := ChainOptionK(func(x int) O.Option[int] {
		if x >= 0 && x <= 100 {
			return O.Some(x)
		}
		return O.None[int]()
	})

	result := F.Pipe2(
		From(-10, 20, 150, 50, 200, 75),
		validateRange,
		Map(func(x int) int { return x * 2 }),
	)
	values := slices.Collect(result)

	expected := A.From(40, 100, 150)
	assert.Equal(t, expected, values)
}

// TestChainOptionK_Composition tests composing multiple ChainOptionK operations
func TestChainOptionK_Composition(t *testing.T) {
	// First filter: only positive
	onlyPositive := ChainOptionK(func(x int) O.Option[int] {
		if x > 0 {
			return O.Some(x)
		}
		return O.None[int]()
	})

	// Second filter: only even
	onlyEven := ChainOptionK(func(x int) O.Option[int] {
		if x%2 == 0 {
			return O.Some(x)
		}
		return O.None[int]()
	})

	result := F.Pipe2(
		From(-2, -1, 0, 1, 2, 3, 4, 5, 6),
		onlyPositive,
		onlyEven,
	)
	values := slices.Collect(result)

	expected := A.From(2, 4, 6)
	assert.Equal(t, expected, values)
}

// TestChainOptionK_StringParsing tests parsing with ChainOptionK
func TestChainOptionK_StringParsing(t *testing.T) {
	// Create a reusable string parser
	parseInt := ChainOptionK(func(s string) O.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return O.Some(n)
		}
		return O.None[int]()
	})

	result := F.Pipe1(
		From("10", "abc", "20", "xyz", "30"),
		parseInt,
	)
	values := slices.Collect(result)

	expected := A.From(10, 20, 30)
	assert.Equal(t, expected, values)
}

// TestFlatMapOptionK_Equivalence tests that FlatMapOptionK is equivalent to ChainOptionK
func TestFlatMapOptionK_Equivalence(t *testing.T) {
	validate := func(x int) O.Option[int] {
		if x >= 0 && x <= 10 {
			return O.Some(x)
		}
		return O.None[int]()
	}

	seq := From(-5, 0, 5, 10, 15)

	// Using ChainOptionK
	result1 := ChainOptionK(validate)(seq)
	values1 := slices.Collect(result1)

	// Using FlatMapOptionK
	result2 := FlatMapOptionK(validate)(seq)
	values2 := slices.Collect(result2)

	// Both should produce the same result
	assert.Equal(t, values1, values2)
	assert.Equal(t, A.From(0, 5, 10), values1)
}

// TestFlatMapOptionK_WithMap tests FlatMapOptionK combined with Map
func TestFlatMapOptionK_WithMap(t *testing.T) {
	// Validate age and convert to category
	validateAge := FlatMapOptionK(func(age int) O.Option[string] {
		if age >= 18 && age <= 120 {
			return O.Some(fmt.Sprintf("Valid age: %d", age))
		}
		return O.None[string]()
	})

	result := F.Pipe1(
		From(15, 25, 150, 30, 200),
		validateAge,
	)
	values := slices.Collect(result)

	expected := A.From("Valid age: 25", "Valid age: 30")
	assert.Equal(t, expected, values)
}

// TestChainOptionK_LookupOperation tests using ChainOptionK for lookup operations
func TestChainOptionK_LookupOperation(t *testing.T) {
	// Simulate a lookup table
	lookup := map[string]int{
		"one":   1,
		"two":   2,
		"three": 3,
	}

	lookupValue := ChainOptionK(func(key string) O.Option[int] {
		if val, ok := lookup[key]; ok {
			return O.Some(val)
		}
		return O.None[int]()
	})

	result := F.Pipe1(
		From("one", "invalid", "two", "missing", "three"),
		lookupValue,
	)
	values := slices.Collect(result)

	expected := A.From(1, 2, 3)
	assert.Equal(t, expected, values)
}

// TestMonadChainOptionK_EarlyTermination tests that iteration stops when yield returns false
func TestMonadChainOptionK_EarlyTermination(t *testing.T) {
	callCount := 0
	countCalls := func(x int) O.Option[int] {
		callCount++
		return O.Some(x)
	}

	seq := From(1, 2, 3, 4, 5)
	result := MonadChainOptionK(seq, countCalls)

	// Collect only first 3 elements
	collected := make([]int, 0)
	for v := range result {
		collected = append(collected, v)
		if len(collected) >= 3 {
			break
		}
	}

	// Should have called the function only 3 times due to early termination
	assert.Equal(t, 3, callCount)
	assert.Equal(t, A.From(1, 2, 3), collected)
}

// TestChainOptionK_WithReduce tests ChainOptionK with reduction
func TestChainOptionK_WithReduce(t *testing.T) {
	// Parse and sum valid numbers
	parseInt := ChainOptionK(func(s string) O.Option[int] {
		if n, err := strconv.Atoi(s); err == nil {
			return O.Some(n)
		}
		return O.None[int]()
	})

	result := F.Pipe1(
		From("10", "invalid", "20", "bad", "30"),
		parseInt,
	)

	sum := MonadReduce(result, func(acc, x int) int {
		return acc + x
	}, 0)

	assert.Equal(t, 60, sum)
}

// TestFlatMapOptionK_NestedOptions tests FlatMapOptionK with nested option handling
func TestFlatMapOptionK_NestedOptions(t *testing.T) {
	type Result struct {
		Value int
		Valid bool
	}

	// Extract value only if valid
	extractValid := FlatMapOptionK(func(r Result) O.Option[int] {
		if r.Valid {
			return O.Some(r.Value)
		}
		return O.None[int]()
	})

	seq := From(
		Result{10, true},
		Result{20, false},
		Result{30, true},
		Result{40, false},
		Result{50, true},
	)

	result := F.Pipe1(seq, extractValid)
	values := slices.Collect(result)

	expected := A.From(10, 30, 50)
	assert.Equal(t, expected, values)
}
