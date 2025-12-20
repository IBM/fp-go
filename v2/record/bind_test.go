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

package record

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

type TestState struct {
	Name    string
	Count   int
	Version int
}

func TestDo(t *testing.T) {
	result := Do[string, TestState]()
	assert.NotNil(t, result)
	assert.Empty(t, result)
	assert.Equal(t, map[string]TestState{}, result)
}

func TestBindTo(t *testing.T) {
	input := map[string]string{"a": "Alice", "b": "Bob"}
	result := F.Pipe1(
		input,
		BindTo[TestState, string, string](func(name string) TestState {
			return TestState{Name: name}
		}),
	)
	expected := map[string]TestState{
		"a": {Name: "Alice"},
		"b": {Name: "Bob"},
	}
	assert.Equal(t, expected, result)
}

func TestLet(t *testing.T) {
	input := map[string]TestState{
		"a": {Name: "Alice"},
		"b": {Name: "Bob"},
	}
	result := F.Pipe1(
		input,
		Let[TestState, int, string](
			func(length int) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Count = length
					return s
				}
			},
			func(s TestState) int {
				return len(s.Name)
			},
		),
	)
	expected := map[string]TestState{
		"a": {Name: "Alice", Count: 5},
		"b": {Name: "Bob", Count: 3},
	}
	assert.Equal(t, expected, result)
}

func TestLetTo(t *testing.T) {
	input := map[string]TestState{
		"a": {Name: "Alice"},
		"b": {Name: "Bob"},
	}
	result := F.Pipe1(
		input,
		LetTo[TestState, int, string](
			func(version int) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Version = version
					return s
				}
			},
			2,
		),
	)
	expected := map[string]TestState{
		"a": {Name: "Alice", Version: 2},
		"b": {Name: "Bob", Version: 2},
	}
	assert.Equal(t, expected, result)
}

func TestBind(t *testing.T) {
	monoid := MergeMonoid[string, TestState]()

	// Bind chains computations where each step can depend on previous results
	result := F.Pipe1(
		map[string]string{"x": "test"},
		Bind[string, int](monoid)(
			func(length int) func(string) TestState {
				return func(s string) TestState {
					return TestState{Name: s, Count: length}
				}
			},
			func(s string) map[string]int {
				return map[string]int{"x": len(s)}
			},
		),
	)

	expected := map[string]TestState{
		"x": {Name: "test", Count: 4},
	}
	assert.Equal(t, expected, result)
}

func TestApS(t *testing.T) {
	monoid := MergeMonoid[string, TestState]()

	// ApS applies independent computations
	names := map[string]string{"x": "Alice"}
	counts := map[string]int{"x": 10}

	result := F.Pipe2(
		map[string]TestState{"x": {}},
		ApS[TestState, string](monoid)(
			func(name string) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Name = name
					return s
				}
			},
			names,
		),
		ApS[TestState, int](monoid)(
			func(count int) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Count = count
					return s
				}
			},
			counts,
		),
	)

	expected := map[string]TestState{
		"x": {Name: "Alice", Count: 10},
	}
	assert.Equal(t, expected, result)
}

func TestBindChain(t *testing.T) {
	// Test a complete do-notation chain with BindTo, Let, and LetTo
	result := F.Pipe3(
		map[string]string{"x": "Alice", "y": "Bob"},
		BindTo[TestState, string, string](func(name string) TestState {
			return TestState{Name: name}
		}),
		Let[TestState, int, string](
			func(count int) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Count = count
					return s
				}
			},
			func(s TestState) int {
				return len(s.Name)
			},
		),
		LetTo[TestState, int, string](
			func(version int) func(TestState) TestState {
				return func(s TestState) TestState {
					s.Version = version
					return s
				}
			},
			1,
		),
	)

	expected := map[string]TestState{
		"x": {Name: "Alice", Count: 5, Version: 1},
		"y": {Name: "Bob", Count: 3, Version: 1},
	}
	assert.Equal(t, expected, result)
}

func TestBindWithDependentComputation(t *testing.T) {
	// Test Bind where the computation creates new keys based on input
	monoid := MergeMonoid[string, TestState]()

	result := F.Pipe1(
		map[string]int{"x": 5},
		Bind[int, string](monoid)(
			func(str string) func(int) TestState {
				return func(n int) TestState {
					return TestState{Name: str, Count: n}
				}
			},
			func(n int) map[string]string {
				// Create a string based on the number
				result := ""
				for i := 0; i < n; i++ {
					result += "a"
				}
				return map[string]string{"x": result}
			},
		),
	)

	expected := map[string]TestState{
		"x": {Name: "aaaaa", Count: 5},
	}
	assert.Equal(t, expected, result)
}
