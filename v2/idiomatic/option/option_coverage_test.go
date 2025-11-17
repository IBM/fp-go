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
	"testing"

	"github.com/stretchr/testify/assert"
)

// Test Logger function
func TestLogger(t *testing.T) {
	logger := Logger[int]()
	logFunc := logger("test")

	// Test with Some
	result, resultok := logFunc(Some(42))
	AssertEq(Some(42))(result, resultok)(t)

	// Test with None
	result, resultok = logFunc(None[int]())
	AssertEq(None[int]())(result, resultok)(t)
}

// Test TraverseArrayG with custom slice types
func TestTraverseArrayG(t *testing.T) {
	type MySlice []int
	type MyResultSlice []string

	f := func(x int) (string, bool) {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	result, resultok := TraverseArrayG[MySlice, MyResultSlice](f)(MySlice{1, 2, 3})
	AssertEq(Some(MyResultSlice{"1", "2", "3"}))(result, resultok)(t)

	// Test with failure
	result, resultok = TraverseArrayG[MySlice, MyResultSlice](f)(MySlice{1, -1, 3})
	AssertEq(None[MyResultSlice]())(result, resultok)(t)
}

// Test TraverseRecordG with custom map types
func TestTraverseRecordG(t *testing.T) {
	type MyMap map[string]int
	type MyResultMap map[string]string

	f := func(x int) (string, bool) {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	input := MyMap{"a": 1, "b": 2}
	result, resultok := TraverseRecordG[MyMap, MyResultMap](f)(input)

	assert.True(t, IsSome(result, resultok))
	assert.Equal(t, "1", result["a"])
	assert.Equal(t, "2", result["b"])
}

// Test TraverseTuple3 through TraverseTuple10
func TestTraverseTuple3(t *testing.T) {
	f1 := func(x int) (int, bool) { return Some(x * 2) }
	f2 := func(s string) (string, bool) { return Some(s + "!") }
	f3 := func(b bool) (bool, bool) { return Some(!b) }

	traverse := TraverseTuple3(f1, f2, f3)
	r1, r2, r3, resultok := traverse(5, "hello", true)
	assert.True(t, resultok)
	assert.Equal(t, r1, 10)
	assert.Equal(t, r2, "hello!")
	assert.Equal(t, r3, false)
}

func TestTraverseTuple4(t *testing.T) {
	f1 := func(x int) (int, bool) { return Some(x * 2) }
	f2 := func(x int) (int, bool) { return Some(x + 1) }
	f3 := func(x int) (int, bool) { return Some(x - 1) }
	f4 := func(x int) (int, bool) { return Some(x * 3) }

	traverse := TraverseTuple4(f1, f2, f3, f4)
	r1, r2, r3, r4, resultok := traverse(1, 2, 3, 4)
	assert.True(t, resultok)
	assert.Equal(t, r1, 2)
	assert.Equal(t, r2, 3)
	assert.Equal(t, r3, 2)
	assert.Equal(t, r4, 12)
}

// Test edge cases for MonadFold
func TestMonadFoldEdgeCases(t *testing.T) {
	// Test with complex types
	type ComplexType struct {
		value int
		name  string
	}

	result := Fold(
		func() string { return "none" },
		func(ct ComplexType) string { return ct.name },
	)(Some(ComplexType{value: 42, name: "test"}))

	assert.Equal(t, "test", result)

	result = Fold(func() string { return "none" },
		func(ct ComplexType) string { return ct.name },
	)(None[ComplexType]())

	assert.Equal(t, "none", result)
}

// Test TraverseArrayWithIndexG
func TestTraverseArrayWithIndexG(t *testing.T) {
	type MySlice []int
	type MyResultSlice []string

	f := func(i int, x int) (string, bool) {
		return Some(fmt.Sprintf("%d:%d", i, x))
	}

	result, resultok := TraverseArrayWithIndexG[MySlice, MyResultSlice](f)(MySlice{10, 20, 30})
	AssertEq(Some(MyResultSlice{"0:10", "1:20", "2:30"}))(result, resultok)(t)
}

// Test TraverseRecordWithIndexG
func TestTraverseRecordWithIndexG(t *testing.T) {
	type MyMap map[string]int
	type MyResultMap map[string]string

	f := func(k string, v int) (string, bool) {
		return Some(fmt.Sprintf("%s=%d", k, v))
	}

	input := MyMap{"a": 1, "b": 2}
	result, resultok := TraverseRecordWithIndexG[MyMap, MyResultMap](f)(input)

	assert.True(t, IsSome(result, resultok))
}

// Test TraverseTuple1
func TestTraverseTuple1(t *testing.T) {
	f := func(x int) (int, bool) { return Some(x * 2) }

	traverse := TraverseTuple1(f)
	result, resultok := traverse(5)

	assert.True(t, resultok)
	assert.Equal(t, 10, result)
}
