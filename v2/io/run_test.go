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

package io

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// TestRun_BasicValue tests that Run executes a simple IO computation
func TestRun_BasicValue(t *testing.T) {
	io := Of(42)
	result := Run(io)
	assert.Equal(t, 42, result)
}

// TestRun_String tests Run with string values
func TestRun_String(t *testing.T) {
	io := Of("Hello, World!")
	result := Run(io)
	assert.Equal(t, "Hello, World!", result)
}

// TestRun_WithMap tests Run with a mapped computation
func TestRun_WithMap(t *testing.T) {
	io := F.Pipe1(
		Of(5),
		Map(N.Mul(2)),
	)
	result := Run(io)
	assert.Equal(t, 10, result)
}

// TestRun_WithChain tests Run with chained computations
func TestRun_WithChain(t *testing.T) {
	io := F.Pipe1(
		Of(3),
		Chain(func(x int) IO[int] {
			return Of(x * x)
		}),
	)
	result := Run(io)
	assert.Equal(t, 9, result)
}

// TestRun_ComposedOperations tests Run with multiple composed operations
func TestRun_ComposedOperations(t *testing.T) {
	io := F.Pipe3(
		Of(5),
		Map(N.Mul(2)), // 10
		Map(N.Add(3)), // 13
		Map(N.Sub(1)), // 12
	)
	result := Run(io)
	assert.Equal(t, 12, result)
}

// TestRun_WithSideEffect tests that Run executes side effects
func TestRun_WithSideEffect(t *testing.T) {
	counter := 0
	io := func() int {
		counter++
		return counter
	}

	// First execution
	result1 := Run(io)
	assert.Equal(t, 1, result1)
	assert.Equal(t, 1, counter)

	// Second execution (side effect happens again)
	result2 := Run(io)
	assert.Equal(t, 2, result2)
	assert.Equal(t, 2, counter)
}

// TestRun_LazyEvaluation tests that IO is lazy until Run is called
func TestRun_LazyEvaluation(t *testing.T) {
	executed := false
	io := func() bool {
		executed = true
		return true
	}

	// IO created but not executed
	assert.False(t, executed, "IO should not execute until Run is called")

	// Now execute
	result := Run(io)
	assert.True(t, executed, "IO should execute when Run is called")
	assert.True(t, result)
}

// TestRun_WithFlatten tests Run with nested IO
func TestRun_WithFlatten(t *testing.T) {
	nested := Of(Of(42))
	flattened := Flatten(nested)
	result := Run(flattened)
	assert.Equal(t, 42, result)
}

// TestRun_WithAp tests Run with applicative operations
func TestRun_WithAp(t *testing.T) {
	double := N.Mul(2)
	io := F.Pipe1(
		Of(double),
		Ap[int](Of(21)),
	)
	result := Run(io)
	assert.Equal(t, 42, result)
}

// TestRun_DifferentTypes tests Run with various types
func TestRun_DifferentTypes(t *testing.T) {
	// Test with bool
	boolIO := Of(true)
	assert.True(t, Run(boolIO))

	// Test with float
	floatIO := Of(3.14)
	assert.Equal(t, 3.14, Run(floatIO))

	// Test with slice
	sliceIO := Of([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2, 3}, Run(sliceIO))

	// Test with struct
	type Person struct {
		Name string
		Age  int
	}
	personIO := Of(Person{Name: "Alice", Age: 30})
	assert.Equal(t, Person{Name: "Alice", Age: 30}, Run(personIO))
}

// TestRun_WithApFirst tests Run with ApFirst combinator
func TestRun_WithApFirst(t *testing.T) {
	io := F.Pipe1(
		Of("first"),
		ApFirst[string](Of("second")),
	)
	result := Run(io)
	assert.Equal(t, "first", result)
}

// TestRun_WithApSecond tests Run with ApSecond combinator
func TestRun_WithApSecond(t *testing.T) {
	io := F.Pipe1(
		Of("first"),
		ApSecond[string](Of("second")),
	)
	result := Run(io)
	assert.Equal(t, "second", result)
}

// TestRun_MultipleExecutions tests that Run can be called multiple times
func TestRun_MultipleExecutions(t *testing.T) {
	io := Of(100)

	// Execute multiple times
	result1 := Run(io)
	result2 := Run(io)
	result3 := Run(io)

	assert.Equal(t, 100, result1)
	assert.Equal(t, 100, result2)
	assert.Equal(t, 100, result3)
}

// TestRun_WithChainedSideEffects tests Run with multiple side effects
func TestRun_WithChainedSideEffects(t *testing.T) {
	log := []string{}

	io := F.Pipe2(
		func() string {
			log = append(log, "step1")
			return "a"
		},
		Chain(func(s string) IO[string] {
			return func() string {
				log = append(log, "step2")
				return s + "b"
			}
		}),
		Chain(func(s string) IO[string] {
			return func() string {
				log = append(log, "step3")
				return s + "c"
			}
		}),
	)

	result := Run(io)
	assert.Equal(t, "abc", result)
	assert.Equal(t, []string{"step1", "step2", "step3"}, log)
}

// TestRun_ZeroValue tests Run with zero values
func TestRun_ZeroValue(t *testing.T) {
	// Test with zero int
	intIO := Of(0)
	assert.Equal(t, 0, Run(intIO))

	// Test with empty string
	strIO := Of("")
	assert.Equal(t, "", Run(strIO))

	// Test with nil slice
	var nilSlice []int
	sliceIO := Of(nilSlice)
	assert.Nil(t, Run(sliceIO))
}
