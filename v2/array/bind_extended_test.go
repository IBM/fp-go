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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

// TestLet tests the Let function
func TestLet(t *testing.T) {
	type State1 struct {
		X int
	}
	type State2 struct {
		X      int
		Double int
	}

	// Test Let with pure computation
	result := F.Pipe1(
		[]State1{{X: 5}, {X: 10}},
		Let(
			func(double int) func(s State1) State2 {
				return func(s State1) State2 {
					return State2{X: s.X, Double: double}
				}
			},
			func(s State1) int { return s.X * 2 },
		),
	)

	expected := []State2{{X: 5, Double: 10}, {X: 10, Double: 20}}
	assert.Equal(t, expected, result)

	// Test Let with empty array
	empty := []State1{}
	result2 := F.Pipe1(
		empty,
		Let(
			func(double int) func(s State1) State2 {
				return func(s State1) State2 {
					return State2{X: s.X, Double: double}
				}
			},
			func(s State1) int { return s.X * 2 },
		),
	)
	assert.Equal(t, []State2{}, result2)
}

// TestLetTo tests the LetTo function
func TestLetTo(t *testing.T) {
	type State1 struct {
		X int
	}
	type State2 struct {
		X    int
		Name string
	}

	// Test LetTo with constant value
	result := F.Pipe1(
		[]State1{{X: 1}, {X: 2}},
		LetTo(
			func(name string) func(s State1) State2 {
				return func(s State1) State2 {
					return State2{X: s.X, Name: name}
				}
			},
			"constant",
		),
	)

	expected := []State2{{X: 1, Name: "constant"}, {X: 2, Name: "constant"}}
	assert.Equal(t, expected, result)

	// Test LetTo with different constant
	result2 := F.Pipe1(
		[]State1{{X: 10}},
		LetTo(
			func(name string) func(s State1) State2 {
				return func(s State1) State2 {
					return State2{X: s.X, Name: name}
				}
			},
			"test",
		),
	)

	expected2 := []State2{{X: 10, Name: "test"}}
	assert.Equal(t, expected2, result2)
}

// TestBindTo tests the BindTo function
func TestBindTo(t *testing.T) {
	type State struct {
		X int
	}

	// Test BindTo with integers
	result := F.Pipe1(
		[]int{1, 2, 3},
		BindTo(func(x int) State {
			return State{X: x}
		}),
	)

	expected := []State{{X: 1}, {X: 2}, {X: 3}}
	assert.Equal(t, expected, result)

	// Test BindTo with strings
	type StringState struct {
		Value string
	}

	result2 := F.Pipe1(
		[]string{"hello", "world"},
		BindTo(func(s string) StringState {
			return StringState{Value: s}
		}),
	)

	expected2 := []StringState{{Value: "hello"}, {Value: "world"}}
	assert.Equal(t, expected2, result2)

	// Test BindTo with empty array
	empty := []int{}
	result3 := F.Pipe1(
		empty,
		BindTo(func(x int) State {
			return State{X: x}
		}),
	)
	assert.Equal(t, []State{}, result3)
}

// TestDoWithLetAndBindTo tests combining Do, Let, LetTo, and BindTo
func TestDoWithLetAndBindTo(t *testing.T) {
	type State1 struct {
		X int
	}
	type State2 struct {
		X      int
		Double int
	}
	type State3 struct {
		X      int
		Double int
		Name   string
	}

	// Test complex pipeline
	result := F.Pipe3(
		[]int{5, 10},
		BindTo(func(x int) State1 {
			return State1{X: x}
		}),
		Let(
			func(double int) func(s State1) State2 {
				return func(s State1) State2 {
					return State2{X: s.X, Double: double}
				}
			},
			func(s State1) int { return s.X * 2 },
		),
		LetTo(
			func(name string) func(s State2) State3 {
				return func(s State2) State3 {
					return State3{X: s.X, Double: s.Double, Name: name}
				}
			},
			"result",
		),
	)

	expected := []State3{
		{X: 5, Double: 10, Name: "result"},
		{X: 10, Double: 20, Name: "result"},
	}
	assert.Equal(t, expected, result)
}
