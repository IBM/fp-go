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

	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

// Test Logger function
func TestLogger(t *testing.T) {
	logger := Logger[int]()
	logFunc := logger("test")

	// Test with Some
	result := logFunc(Some(42))
	assert.Equal(t, Some(42), result)

	// Test with None
	result = logFunc(None[int]())
	assert.Equal(t, None[int](), result)
}

// Test TraverseArrayG with custom slice types
func TestTraverseArrayG(t *testing.T) {
	type MySlice []int
	type MyResultSlice []string

	f := func(x int) Option[string] {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	result := TraverseArrayG[MySlice, MyResultSlice](f)(MySlice{1, 2, 3})
	expected := Some(MyResultSlice{"1", "2", "3"})
	assert.Equal(t, expected, result)

	// Test with failure
	result = TraverseArrayG[MySlice, MyResultSlice](f)(MySlice{1, -1, 3})
	assert.Equal(t, None[MyResultSlice](), result)
}

// Test SequenceArrayG with custom slice types
func TestSequenceArrayG(t *testing.T) {
	type MySlice []int

	input := []Option[int]{Some(1), Some(2), Some(3)}
	result := SequenceArrayG[MySlice](input)
	expected := Some(MySlice{1, 2, 3})
	assert.Equal(t, expected, result)

	// Test with None
	input = []Option[int]{Some(1), None[int](), Some(3)}
	result = SequenceArrayG[MySlice](input)
	assert.Equal(t, None[MySlice](), result)
}

// Test CompactArrayG with custom slice types
func TestCompactArrayG(t *testing.T) {
	type MySlice []int

	input := []Option[int]{Some(1), None[int](), Some(3), Some(5)}
	result := CompactArrayG[[]Option[int], MySlice](input)
	expected := MySlice{1, 3, 5}
	assert.Equal(t, expected, result)
}

// Test TraverseRecordG with custom map types
func TestTraverseRecordG(t *testing.T) {
	type MyMap map[string]int
	type MyResultMap map[string]string

	f := func(x int) Option[string] {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	input := MyMap{"a": 1, "b": 2}
	result := TraverseRecordG[MyMap, MyResultMap](f)(input)

	assert.True(t, IsSome(result))
	unwrapped, _ := Unwrap(result)
	assert.Equal(t, "1", unwrapped["a"])
	assert.Equal(t, "2", unwrapped["b"])
}

// Test SequenceRecordG with custom map types
func TestSequenceRecordG(t *testing.T) {
	type MyMap map[string]int

	input := map[string]Option[int]{"a": Some(1), "b": Some(2)}
	result := SequenceRecordG[MyMap](input)

	assert.True(t, IsSome(result))
	unwrapped, _ := Unwrap(result)
	assert.Equal(t, 1, unwrapped["a"])
	assert.Equal(t, 2, unwrapped["b"])
}

// Test CompactRecordG with custom map types
func TestCompactRecordG(t *testing.T) {
	type MyMap map[string]int

	input := map[string]Option[int]{"a": Some(1), "b": None[int](), "c": Some(3)}
	result := CompactRecordG[map[string]Option[int], MyMap](input)

	expected := MyMap{"a": 1, "c": 3}
	assert.Equal(t, expected, result)
}

// Test Optionize3 through Optionize10
func TestOptionize3(t *testing.T) {
	f := func(a, b, c int) (int, bool) {
		if a > 0 && b > 0 && c > 0 {
			return a + b + c, true
		}
		return 0, false
	}

	optF := Optionize3(f)
	assert.Equal(t, Some(6), optF(1, 2, 3))
	assert.Equal(t, None[int](), optF(-1, 2, 3))
}

func TestOptionize4(t *testing.T) {
	f := func(a, b, c, d int) (int, bool) {
		sum := a + b + c + d
		return sum, sum > 0
	}

	optF := Optionize4(f)
	assert.Equal(t, Some(10), optF(1, 2, 3, 4))
	assert.Equal(t, None[int](), optF(-5, 1, 1, 1))
}

func TestOptionize5(t *testing.T) {
	f := func(a, b, c, d, e int) (int, bool) {
		sum := a + b + c + d + e
		return sum, sum > 0
	}

	optF := Optionize5(f)
	assert.Equal(t, Some(15), optF(1, 2, 3, 4, 5))
}

// Test Unoptionize3 through Unoptionize10
func TestUnoptionize3(t *testing.T) {
	f := func(a, b, c int) Option[int] {
		if a > 0 && b > 0 && c > 0 {
			return Some(a + b + c)
		}
		return None[int]()
	}

	unoptF := Unoptionize3(f)
	val, ok := unoptF(1, 2, 3)
	assert.True(t, ok)
	assert.Equal(t, 6, val)

	_, ok = unoptF(-1, 2, 3)
	assert.False(t, ok)
}

func TestUnoptionize4(t *testing.T) {
	f := func(a, b, c, d int) Option[int] {
		return Some(a + b + c + d)
	}

	unoptF := Unoptionize4(f)
	val, ok := unoptF(1, 2, 3, 4)
	assert.True(t, ok)
	assert.Equal(t, 10, val)
}

// Test SequenceT5 through SequenceT10
func TestSequenceT5(t *testing.T) {
	result := SequenceT5(Some(1), Some(2), Some(3), Some(4), Some(5))
	expected := Some(T.MakeTuple5(1, 2, 3, 4, 5))
	assert.Equal(t, expected, result)

	// Test with None
	result = SequenceT5(Some(1), None[int](), Some(3), Some(4), Some(5))
	assert.Equal(t, None[T.Tuple5[int, int, int, int, int]](), result)
}

func TestSequenceT6(t *testing.T) {
	result := SequenceT6(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6))
	expected := Some(T.MakeTuple6(1, 2, 3, 4, 5, 6))
	assert.Equal(t, expected, result)
}

func TestSequenceT7(t *testing.T) {
	result := SequenceT7(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7))
	expected := Some(T.MakeTuple7(1, 2, 3, 4, 5, 6, 7))
	assert.Equal(t, expected, result)
}

func TestSequenceT8(t *testing.T) {
	result := SequenceT8(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8))
	expected := Some(T.MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, expected, result)
}

func TestSequenceT9(t *testing.T) {
	result := SequenceT9(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8), Some(9))
	expected := Some(T.MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9))
	assert.Equal(t, expected, result)
}

func TestSequenceT10(t *testing.T) {
	result := SequenceT10(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8), Some(9), Some(10))
	expected := Some(T.MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple4 through SequenceTuple10
func TestSequenceTuple4(t *testing.T) {
	tuple := T.MakeTuple4(Some(1), Some(2), Some(3), Some(4))
	result := SequenceTuple4(tuple)
	expected := Some(T.MakeTuple4(1, 2, 3, 4))
	assert.Equal(t, expected, result)
}

func TestSequenceTuple5(t *testing.T) {
	tuple := T.MakeTuple5(Some(1), Some(2), Some(3), Some(4), Some(5))
	result := SequenceTuple5(tuple)
	expected := Some(T.MakeTuple5(1, 2, 3, 4, 5))
	assert.Equal(t, expected, result)
}

func TestSequenceTuple6(t *testing.T) {
	tuple := T.MakeTuple6(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6))
	result := SequenceTuple6(tuple)
	expected := Some(T.MakeTuple6(1, 2, 3, 4, 5, 6))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple3 through TraverseTuple10
func TestTraverseTuple3(t *testing.T) {
	f1 := func(x int) Option[int] { return Some(x * 2) }
	f2 := func(s string) Option[string] { return Some(s + "!") }
	f3 := func(b bool) Option[bool] { return Some(!b) }

	traverse := TraverseTuple3(f1, f2, f3)
	tuple := T.MakeTuple3(5, "hello", true)
	result := traverse(tuple)

	expected := Some(T.MakeTuple3(10, "hello!", false))
	assert.Equal(t, expected, result)
}

func TestTraverseTuple4(t *testing.T) {
	f1 := func(x int) Option[int] { return Some(x * 2) }
	f2 := func(x int) Option[int] { return Some(x + 1) }
	f3 := func(x int) Option[int] { return Some(x - 1) }
	f4 := func(x int) Option[int] { return Some(x * 3) }

	traverse := TraverseTuple4(f1, f2, f3, f4)
	tuple := T.MakeTuple4(1, 2, 3, 4)
	result := traverse(tuple)

	expected := Some(T.MakeTuple4(2, 3, 2, 12))
	assert.Equal(t, expected, result)
}

// Test JSON marshaling edge cases
func TestJSONMarshalNone(t *testing.T) {
	none := None[int]()
	data, err := none.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("null"), data)
}

func TestJSONMarshalSome(t *testing.T) {
	some := Some(42)
	data, err := some.MarshalJSON()
	assert.NoError(t, err)
	assert.Equal(t, []byte("42"), data)
}

// Test Format method
func TestFormat(t *testing.T) {
	some := Some(42)
	formatted := fmt.Sprintf("%s", some)
	assert.Contains(t, formatted, "Some")
	assert.Contains(t, formatted, "42")

	none := None[int]()
	formatted = fmt.Sprintf("%s", none)
	assert.Contains(t, formatted, "None")
}

// Test edge cases for MonadFold
func TestMonadFoldEdgeCases(t *testing.T) {
	// Test with complex types
	type ComplexType struct {
		value int
		name  string
	}

	some := Some(ComplexType{value: 42, name: "test"})
	result := MonadFold(some,
		func() string { return "none" },
		func(ct ComplexType) string { return ct.name },
	)
	assert.Equal(t, "test", result)

	none := None[ComplexType]()
	result = MonadFold(none,
		func() string { return "none" },
		func(ct ComplexType) string { return ct.name },
	)
	assert.Equal(t, "none", result)
}

// Test TraverseArrayWithIndexG
func TestTraverseArrayWithIndexG(t *testing.T) {
	type MySlice []int
	type MyResultSlice []string

	f := func(i int, x int) Option[string] {
		return Some(fmt.Sprintf("%d:%d", i, x))
	}

	result := TraverseArrayWithIndexG[MySlice, MyResultSlice](f)(MySlice{10, 20, 30})
	expected := Some(MyResultSlice{"0:10", "1:20", "2:30"})
	assert.Equal(t, expected, result)
}

// Test TraverseRecordWithIndexG
func TestTraverseRecordWithIndexG(t *testing.T) {
	type MyMap map[string]int
	type MyResultMap map[string]string

	f := func(k string, v int) Option[string] {
		return Some(fmt.Sprintf("%s=%d", k, v))
	}

	input := MyMap{"a": 1, "b": 2}
	result := TraverseRecordWithIndexG[MyMap, MyResultMap](f)(input)

	assert.True(t, IsSome(result))
}

// Test SequenceTuple1
func TestSequenceTuple1(t *testing.T) {
	tuple := T.MakeTuple1(Some(42))
	result := SequenceTuple1(tuple)
	expected := Some(T.MakeTuple1(42))
	assert.Equal(t, expected, result)

	// Test with None
	tuple = T.MakeTuple1(None[int]())
	result = SequenceTuple1(tuple)
	assert.Equal(t, None[T.Tuple1[int]](), result)
}

// Test TraverseTuple1
func TestTraverseTuple1(t *testing.T) {
	f := func(x int) Option[int] { return Some(x * 2) }

	traverse := TraverseTuple1(f)
	tuple := T.MakeTuple1(5)
	result := traverse(tuple)

	expected := Some(T.MakeTuple1(10))
	assert.Equal(t, expected, result)
}

// Test Unoptionize2
func TestUnoptionize2(t *testing.T) {
	f := func(a, b int) Option[int] {
		return Some(a + b)
	}

	unoptF := Unoptionize2(f)
	val, ok := unoptF(2, 3)
	assert.True(t, ok)
	assert.Equal(t, 5, val)
}

// Test Unoptionize5
func TestUnoptionize5(t *testing.T) {
	f := func(a, b, c, d, e int) Option[int] {
		return Some(a + b + c + d + e)
	}

	unoptF := Unoptionize5(f)
	val, ok := unoptF(1, 2, 3, 4, 5)
	assert.True(t, ok)
	assert.Equal(t, 15, val)
}

// Test Unoptionize6
func TestUnoptionize6(t *testing.T) {
	f := func(a, b, c, d, e, f int) Option[int] {
		return Some(a + b + c + d + e + f)
	}

	unoptF := Unoptionize6(f)
	val, ok := unoptF(1, 2, 3, 4, 5, 6)
	assert.True(t, ok)
	assert.Equal(t, 21, val)
}

// Test Optionize6
func TestOptionize6(t *testing.T) {
	f := func(a, b, c, d, e, f int) (int, bool) {
		sum := a + b + c + d + e + f
		return sum, sum > 0
	}

	optF := Optionize6(f)
	assert.Equal(t, Some(21), optF(1, 2, 3, 4, 5, 6))
}

// Test Optionize7
func TestOptionize7(t *testing.T) {
	f := func(a, b, c, d, e, f, g int) (int, bool) {
		sum := a + b + c + d + e + f + g
		return sum, sum > 0
	}

	optF := Optionize7(f)
	assert.Equal(t, Some(28), optF(1, 2, 3, 4, 5, 6, 7))
}

// Test Unoptionize7
func TestUnoptionize7(t *testing.T) {
	f := func(a, b, c, d, e, f, g int) Option[int] {
		return Some(a + b + c + d + e + f + g)
	}

	unoptF := Unoptionize7(f)
	val, ok := unoptF(1, 2, 3, 4, 5, 6, 7)
	assert.True(t, ok)
	assert.Equal(t, 28, val)
}

// Test Optionize8
func TestOptionize8(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h int) (int, bool) {
		sum := a + b + c + d + e + f + g + h
		return sum, sum > 0
	}

	optF := Optionize8(f)
	assert.Equal(t, Some(36), optF(1, 2, 3, 4, 5, 6, 7, 8))
}

// Test Unoptionize8
func TestUnoptionize8(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h int) Option[int] {
		return Some(a + b + c + d + e + f + g + h)
	}

	unoptF := Unoptionize8(f)
	val, ok := unoptF(1, 2, 3, 4, 5, 6, 7, 8)
	assert.True(t, ok)
	assert.Equal(t, 36, val)
}

// Test Optionize9
func TestOptionize9(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i int) (int, bool) {
		sum := a + b + c + d + e + f + g + h + i
		return sum, sum > 0
	}

	optF := Optionize9(f)
	assert.Equal(t, Some(45), optF(1, 2, 3, 4, 5, 6, 7, 8, 9))
}

// Test Unoptionize9
func TestUnoptionize9(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i int) Option[int] {
		return Some(a + b + c + d + e + f + g + h + i)
	}

	unoptF := Unoptionize9(f)
	val, ok := unoptF(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.True(t, ok)
	assert.Equal(t, 45, val)
}

// Test Optionize10
func TestOptionize10(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j int) (int, bool) {
		sum := a + b + c + d + e + f + g + h + i + j
		return sum, sum > 0
	}

	optF := Optionize10(f)
	assert.Equal(t, Some(55), optF(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
}

// Test Unoptionize10
func TestUnoptionize10(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j int) Option[int] {
		return Some(a + b + c + d + e + f + g + h + i + j)
	}

	unoptF := Unoptionize10(f)
	val, ok := unoptF(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.True(t, ok)
	assert.Equal(t, 55, val)
}

// Test SequenceTuple7
func TestSequenceTuple7(t *testing.T) {
	tuple := T.MakeTuple7(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7))
	result := SequenceTuple7(tuple)
	expected := Some(T.MakeTuple7(1, 2, 3, 4, 5, 6, 7))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple8
func TestSequenceTuple8(t *testing.T) {
	tuple := T.MakeTuple8(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8))
	result := SequenceTuple8(tuple)
	expected := Some(T.MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple9
func TestSequenceTuple9(t *testing.T) {
	tuple := T.MakeTuple9(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8), Some(9))
	result := SequenceTuple9(tuple)
	expected := Some(T.MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple10
func TestSequenceTuple10(t *testing.T) {
	tuple := T.MakeTuple10(Some(1), Some(2), Some(3), Some(4), Some(5), Some(6), Some(7), Some(8), Some(9), Some(10))
	result := SequenceTuple10(tuple)
	expected := Some(T.MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	assert.Equal(t, expected, result)
}

// Test TryCatch with success case
func TestTryCatchSuccess(t *testing.T) {
	result := TryCatch(func() (int, error) {
		return 42, nil
	})
	assert.Equal(t, Some(42), result)
}
