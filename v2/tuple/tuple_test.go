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

package tuple

import (
	"encoding/json"
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/ord"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	value := MakeTuple2("Carsten", 1)
	assert.Equal(t, "Tuple2[string, int](Carsten, 1)", value.String())
}

func TestMarshal(t *testing.T) {
	value := MakeTuple3("Carsten", 1, true)

	data, err := json.Marshal(value)
	require.NoError(t, err)

	var unmarshaled Tuple3[string, int, bool]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, value, unmarshaled)
}

func TestMarshalSmallArray(t *testing.T) {
	value := `["Carsten"]`

	var unmarshaled Tuple3[string, int, bool]
	err := json.Unmarshal([]byte(value), &unmarshaled)
	require.NoError(t, err)

	assert.Equal(t, MakeTuple3("Carsten", 0, false), unmarshaled)
}

// Test Of function
func TestOf(t *testing.T) {
	t1 := Of(42)
	assert.Equal(t, Tuple1[int]{F1: 42}, t1)
	assert.Equal(t, 42, t1.F1)
}

// Test First and Second functions
func TestFirstSecond(t *testing.T) {
	t2 := MakeTuple2("hello", 42)
	assert.Equal(t, "hello", First(t2))
	assert.Equal(t, 42, Second(t2))
}

// Test Swap function
func TestSwap(t *testing.T) {
	t2 := MakeTuple2("hello", 42)
	swapped := Swap(t2)
	assert.Equal(t, MakeTuple2(42, "hello"), swapped)
	assert.Equal(t, 42, swapped.F1)
	assert.Equal(t, "hello", swapped.F2)
}

// Test Of2 function
func TestOf2(t *testing.T) {
	pairWith42 := Of2[string](42)
	result := pairWith42("hello")
	assert.Equal(t, MakeTuple2("hello", 42), result)
}

// Test BiMap function
func TestBiMap(t *testing.T) {
	t2 := MakeTuple2(5, "hello")
	mapper := BiMap(
		S.Size,
		func(n int) string { return fmt.Sprintf("%d", n*2) },
	)
	result := mapper(t2)
	assert.Equal(t, MakeTuple2("10", 5), result)
}

// Test Tupled and Untupled functions
func TestTupled2Untupled2(t *testing.T) {
	add := func(a, b int) int { return a + b }

	// Test Tupled2
	tupledAdd := Tupled2(add)
	result := tupledAdd(MakeTuple2(3, 4))
	assert.Equal(t, 7, result)

	// Test Untupled2
	untupledAdd := Untupled2(tupledAdd)
	result2 := untupledAdd(5, 6)
	assert.Equal(t, 11, result2)
}

func TestTupled3Untupled3(t *testing.T) {
	sum3 := func(a, b, c int) int { return a + b + c }

	tupled := Tupled3(sum3)
	result := tupled(MakeTuple3(1, 2, 3))
	assert.Equal(t, 6, result)

	untupled := Untupled3(tupled)
	result2 := untupled(4, 5, 6)
	assert.Equal(t, 15, result2)
}

// Test Map functions
func TestMap1(t *testing.T) {
	t1 := MakeTuple1(5)
	mapper := Map1(func(n int) string { return fmt.Sprintf("%d", n*2) })
	result := mapper(t1)
	assert.Equal(t, MakeTuple1("10"), result)
}

func TestMap2(t *testing.T) {
	t2 := MakeTuple2(5, "hello")
	mapper := Map2(
		func(n int) string { return fmt.Sprintf("%d", n*2) },
		S.Size,
	)
	result := mapper(t2)
	assert.Equal(t, MakeTuple2("10", 5), result)
}

func TestMap3(t *testing.T) {
	t3 := MakeTuple3(1, 2, 3)
	mapper := Map3(
		N.Mul(2),
		N.Mul(3),
		N.Mul(4),
	)
	result := mapper(t3)
	assert.Equal(t, MakeTuple3(2, 6, 12), result)
}

// Test Replicate functions
func TestReplicate1(t *testing.T) {
	result := Replicate1(42)
	assert.Equal(t, MakeTuple1(42), result)
}

func TestReplicate2(t *testing.T) {
	result := Replicate2(42)
	assert.Equal(t, MakeTuple2(42, 42), result)
}

func TestReplicate3(t *testing.T) {
	result := Replicate3(42)
	assert.Equal(t, MakeTuple3(42, 42, 42), result)
}

// Test ToArray and FromArray functions
func TestToArray1FromArray1(t *testing.T) {
	t1 := MakeTuple1(42)
	toArray := ToArray1(func(n int) int { return n })
	arr := toArray(t1)
	assert.Equal(t, []int{42}, arr)

	fromArray := FromArray1(func(n int) int { return n })
	result := fromArray(arr)
	assert.Equal(t, t1, result)
}

func TestToArray2FromArray2(t *testing.T) {
	t2 := MakeTuple2(1, 2)
	toArray := ToArray2(
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t2)
	assert.Equal(t, []int{1, 2}, arr)

	fromArray := FromArray2(
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t2, result)
}

func TestToArray3FromArray3(t *testing.T) {
	t3 := MakeTuple3(1, 2, 3)
	toArray := ToArray3(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t3)
	assert.Equal(t, []int{1, 2, 3}, arr)

	fromArray := FromArray3(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t3, result)
}

// Test Push functions
func TestPush1(t *testing.T) {
	t1 := MakeTuple1(42)
	push := Push1[int]("hello")
	result := push(t1)
	assert.Equal(t, MakeTuple2(42, "hello"), result)
}

func TestPush2(t *testing.T) {
	t2 := MakeTuple2(1, 2)
	push := Push2[int, int](3)
	result := push(t2)
	assert.Equal(t, MakeTuple3(1, 2, 3), result)
}

func TestPush3(t *testing.T) {
	t3 := MakeTuple3(1, 2, 3)
	push := Push3[int, int, int](4)
	result := push(t3)
	assert.Equal(t, MakeTuple4(1, 2, 3, 4), result)
}

// Test Monoid functions
func TestMonoid1(t *testing.T) {
	m := Monoid1(N.MonoidSum[int]())
	t1 := MakeTuple1(5)
	t2 := MakeTuple1(3)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple1(8), result)
	assert.Equal(t, MakeTuple1(0), m.Empty())
}

func TestMonoid2(t *testing.T) {
	m := Monoid2(S.Monoid, N.MonoidSum[int]())
	t1 := MakeTuple2("hello", 5)
	t2 := MakeTuple2(" world", 3)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple2("hello world", 8), result)
	assert.Equal(t, MakeTuple2("", 0), m.Empty())
}

func TestMonoid3(t *testing.T) {
	m := Monoid3(S.Monoid, N.MonoidSum[int](), N.MonoidProduct[int]())
	t1 := MakeTuple3("a", 2, 3)
	t2 := MakeTuple3("b", 4, 5)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple3("ab", 6, 15), result)
}

// Test Ord functions
func TestOrd1(t *testing.T) {
	o := Ord1(O.FromStrictCompare[int]())
	t1 := MakeTuple1(5)
	t2 := MakeTuple1(3)
	t3 := MakeTuple1(5)

	assert.Equal(t, 1, o.Compare(t1, t2))
	assert.Equal(t, -1, o.Compare(t2, t1))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
	assert.False(t, o.Equals(t1, t2))
}

func TestOrd2(t *testing.T) {
	o := Ord2(O.FromStrictCompare[string](), O.FromStrictCompare[int]())
	t1 := MakeTuple2("a", 1)
	t2 := MakeTuple2("b", 2)
	t3 := MakeTuple2("a", 1)
	t4 := MakeTuple2("a", 2)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 1, o.Compare(t2, t1))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.Equal(t, -1, o.Compare(t1, t4))
	assert.True(t, o.Equals(t1, t3))
	assert.False(t, o.Equals(t1, t2))
}

func TestOrd3(t *testing.T) {
	o := Ord3(O.FromStrictCompare[int](), O.FromStrictCompare[int](), O.FromStrictCompare[int]())
	t1 := MakeTuple3(1, 2, 3)
	t2 := MakeTuple3(1, 2, 4)
	t3 := MakeTuple3(1, 2, 3)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

// Test String methods for different tuple sizes
func TestTuple1String(t *testing.T) {
	t1 := MakeTuple1(42)
	assert.Equal(t, "Tuple1[int](42)", t1.String())
}

func TestTuple3String(t *testing.T) {
	t3 := MakeTuple3("test", 42, true)
	assert.Equal(t, "Tuple3[string, int, bool](test, 42, true)", t3.String())
}

func TestTuple4String(t *testing.T) {
	t4 := MakeTuple4(1, 2, 3, 4)
	assert.Equal(t, "Tuple4[int, int, int, int](1, 2, 3, 4)", t4.String())
}

// Test JSON marshaling for different tuple sizes
func TestTuple1JSON(t *testing.T) {
	t1 := MakeTuple1(42)
	data, err := json.Marshal(t1)
	require.NoError(t, err)
	assert.Equal(t, "[42]", string(data))

	var unmarshaled Tuple1[int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t1, unmarshaled)
}

func TestTuple2JSON(t *testing.T) {
	t2 := MakeTuple2("hello", 42)
	data, err := json.Marshal(t2)
	require.NoError(t, err)

	var unmarshaled Tuple2[string, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t2, unmarshaled)
}

func TestTuple4JSON(t *testing.T) {
	t4 := MakeTuple4(1, 2, 3, 4)
	data, err := json.Marshal(t4)
	require.NoError(t, err)

	var unmarshaled Tuple4[int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t4, unmarshaled)
}

func TestTuple5JSON(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	data, err := json.Marshal(t5)
	require.NoError(t, err)

	var unmarshaled Tuple5[int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t5, unmarshaled)
}

// Test JSON unmarshal error cases
func TestUnmarshalInvalidJSON(t *testing.T) {
	var t2 Tuple2[string, int]
	err := json.Unmarshal([]byte("invalid json"), &t2)
	assert.Error(t, err)
}

func TestUnmarshalInvalidType(t *testing.T) {
	var t2 Tuple2[int, int]
	err := json.Unmarshal([]byte(`["string", 42]`), &t2)
	assert.Error(t, err)
}

// Test MakeTuple functions for various sizes
func TestMakeTuple4(t *testing.T) {
	t4 := MakeTuple4(1, "two", 3.0, true)
	assert.Equal(t, 1, t4.F1)
	assert.Equal(t, "two", t4.F2)
	assert.Equal(t, 3.0, t4.F3)
	assert.Equal(t, true, t4.F4)
}

func TestMakeTuple5(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	assert.Equal(t, 1, t5.F1)
	assert.Equal(t, 5, t5.F5)
}

func TestMakeTuple6(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	assert.Equal(t, 1, t6.F1)
	assert.Equal(t, 6, t6.F6)
}

// Test Tupled/Untupled for larger tuples
func TestTupled4Untupled4(t *testing.T) {
	sum4 := func(a, b, c, d int) int { return a + b + c + d }

	tupled := Tupled4(sum4)
	result := tupled(MakeTuple4(1, 2, 3, 4))
	assert.Equal(t, 10, result)

	untupled := Untupled4(tupled)
	result2 := untupled(2, 3, 4, 5)
	assert.Equal(t, 14, result2)
}

func TestTupled5Untupled5(t *testing.T) {
	sum5 := func(a, b, c, d, e int) int { return a + b + c + d + e }

	tupled := Tupled5(sum5)
	result := tupled(MakeTuple5(1, 2, 3, 4, 5))
	assert.Equal(t, 15, result)

	untupled := Untupled5(tupled)
	result2 := untupled(1, 1, 1, 1, 1)
	assert.Equal(t, 5, result2)
}

// Test Map for larger tuples
func TestMap4(t *testing.T) {
	t4 := MakeTuple4(1, 2, 3, 4)
	mapper := Map4(
		N.Mul(2),
		N.Mul(3),
		N.Mul(4),
		N.Mul(5),
	)
	result := mapper(t4)
	assert.Equal(t, MakeTuple4(2, 6, 12, 20), result)
}

func TestMap5(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	mapper := Map5(
		N.Add(1),
		func(n int) int { return n + 2 },
		func(n int) int { return n + 3 },
		func(n int) int { return n + 4 },
		N.Add(5),
	)
	result := mapper(t5)
	assert.Equal(t, MakeTuple5(2, 4, 6, 8, 10), result)
}

// Test Replicate for larger tuples
func TestReplicate4(t *testing.T) {
	result := Replicate4(7)
	assert.Equal(t, MakeTuple4(7, 7, 7, 7), result)
}

func TestReplicate5(t *testing.T) {
	result := Replicate5(9)
	assert.Equal(t, MakeTuple5(9, 9, 9, 9, 9), result)
}

// Test ToArray/FromArray for larger tuples
func TestToArray4FromArray4(t *testing.T) {
	t4 := MakeTuple4(1, 2, 3, 4)
	toArray := ToArray4(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t4)
	assert.Equal(t, []int{1, 2, 3, 4}, arr)

	fromArray := FromArray4(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t4, result)
}

func TestToArray5FromArray5(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	toArray := ToArray5(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t5)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, arr)

	fromArray := FromArray5(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t5, result)
}

// Test Push for larger tuples
func TestPush4(t *testing.T) {
	t4 := MakeTuple4(1, 2, 3, 4)
	push := Push4[int, int, int, int](5)
	result := push(t4)
	assert.Equal(t, MakeTuple5(1, 2, 3, 4, 5), result)
}

func TestPush5(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	push := Push5[int, int, int, int, int](6)
	result := push(t5)
	assert.Equal(t, MakeTuple6(1, 2, 3, 4, 5, 6), result)
}

// Test Monoid for larger tuples
func TestMonoid4(t *testing.T) {
	m := Monoid4(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple4(1, 2, 3, 4)
	t2 := MakeTuple4(5, 6, 7, 8)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple4(6, 8, 10, 12), result)
}

func TestMonoid5(t *testing.T) {
	m := Monoid5(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple5(1, 2, 3, 4, 5)
	t2 := MakeTuple5(1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple5(2, 3, 4, 5, 6), result)
}

// Test Ord for larger tuples
func TestOrd4(t *testing.T) {
	o := Ord4(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple4(1, 2, 3, 4)
	t2 := MakeTuple4(1, 2, 3, 5)
	t3 := MakeTuple4(1, 2, 3, 4)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd5(t *testing.T) {
	o := Ord5(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple5(1, 2, 3, 4, 5)
	t2 := MakeTuple5(1, 2, 3, 4, 6)
	t3 := MakeTuple5(1, 2, 3, 4, 5)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

// Test larger tuple sizes (6-10)
func TestMakeTuple7(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	assert.Equal(t, 1, t7.F1)
	assert.Equal(t, 7, t7.F7)
}

func TestMakeTuple8(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	assert.Equal(t, 1, t8.F1)
	assert.Equal(t, 8, t8.F8)
}

func TestMakeTuple9(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.Equal(t, 1, t9.F1)
	assert.Equal(t, 9, t9.F9)
}

func TestMakeTuple10(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.Equal(t, 1, t10.F1)
	assert.Equal(t, 10, t10.F10)
}

// Test Tupled/Untupled for sizes 6-10
func TestTupled6Untupled6(t *testing.T) {
	sum6 := func(a, b, c, d, e, f int) int { return a + b + c + d + e + f }

	tupled := Tupled6(sum6)
	result := tupled(MakeTuple6(1, 2, 3, 4, 5, 6))
	assert.Equal(t, 21, result)

	untupled := Untupled6(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1)
	assert.Equal(t, 6, result2)
}

func TestTupled7Untupled7(t *testing.T) {
	sum7 := func(a, b, c, d, e, f, g int) int { return a + b + c + d + e + f + g }

	tupled := Tupled7(sum7)
	result := tupled(MakeTuple7(1, 2, 3, 4, 5, 6, 7))
	assert.Equal(t, 28, result)

	untupled := Untupled7(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 7, result2)
}

func TestTupled8Untupled8(t *testing.T) {
	sum8 := func(a, b, c, d, e, f, g, h int) int { return a + b + c + d + e + f + g + h }

	tupled := Tupled8(sum8)
	result := tupled(MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8))
	assert.Equal(t, 36, result)

	untupled := Untupled8(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 8, result2)
}

func TestTupled9Untupled9(t *testing.T) {
	sum9 := func(a, b, c, d, e, f, g, h, i int) int { return a + b + c + d + e + f + g + h + i }

	tupled := Tupled9(sum9)
	result := tupled(MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9))
	assert.Equal(t, 45, result)

	untupled := Untupled9(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 9, result2)
}

func TestTupled10Untupled10(t *testing.T) {
	sum10 := func(a, b, c, d, e, f, g, h, i, j int) int { return a + b + c + d + e + f + g + h + i + j }

	tupled := Tupled10(sum10)
	result := tupled(MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))
	assert.Equal(t, 55, result)

	untupled := Untupled10(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 10, result2)
}

// Test Map for sizes 6-10
func TestMap6(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	mapper := Map6(
		func(n int) int { return n + 1 },
		func(n int) int { return n + 2 },
		func(n int) int { return n + 3 },
		func(n int) int { return n + 4 },
		N.Add(5),
		func(n int) int { return n + 6 },
	)
	result := mapper(t6)
	assert.Equal(t, MakeTuple6(2, 4, 6, 8, 10, 12), result)
}

func TestMap7(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	mapper := Map7(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t7)
	assert.Equal(t, MakeTuple7(2, 4, 6, 8, 10, 12, 14), result)
}

func TestMap8(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	mapper := Map8(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t8)
	assert.Equal(t, MakeTuple8(2, 4, 6, 8, 10, 12, 14, 16), result)
}

func TestMap9(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	mapper := Map9(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t9)
	assert.Equal(t, MakeTuple9(2, 4, 6, 8, 10, 12, 14, 16, 18), result)
}

func TestMap10(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	mapper := Map10(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t10)
	assert.Equal(t, MakeTuple10(2, 4, 6, 8, 10, 12, 14, 16, 18, 20), result)
}

// Test Replicate for sizes 6-10
func TestReplicate6(t *testing.T) {
	result := Replicate6(11)
	assert.Equal(t, MakeTuple6(11, 11, 11, 11, 11, 11), result)
}

func TestReplicate7(t *testing.T) {
	result := Replicate7(13)
	assert.Equal(t, MakeTuple7(13, 13, 13, 13, 13, 13, 13), result)
}

func TestReplicate8(t *testing.T) {
	result := Replicate8(15)
	assert.Equal(t, MakeTuple8(15, 15, 15, 15, 15, 15, 15, 15), result)
}

func TestReplicate9(t *testing.T) {
	result := Replicate9(17)
	assert.Equal(t, MakeTuple9(17, 17, 17, 17, 17, 17, 17, 17, 17), result)
}

func TestReplicate10(t *testing.T) {
	result := Replicate10(19)
	assert.Equal(t, MakeTuple10(19, 19, 19, 19, 19, 19, 19, 19, 19, 19), result)
}

// Test ToArray/FromArray for sizes 6-10
func TestToArray6FromArray6(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	toArray := ToArray6(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t6)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, arr)

	fromArray := FromArray6(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t6, result)
}

func TestToArray7FromArray7(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	toArray := ToArray7(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t7)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, arr)

	fromArray := FromArray7(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t7, result)
}

func TestToArray8FromArray8(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	toArray := ToArray8(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t8)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8}, arr)

	fromArray := FromArray8(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t8, result)
}

func TestToArray9FromArray9(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	toArray := ToArray9(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t9)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, arr)

	fromArray := FromArray9(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t9, result)
}

func TestToArray10FromArray10(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	toArray := ToArray10(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t10)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, arr)

	fromArray := FromArray10(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t10, result)
}

// Test Push for sizes 6-10
func TestPush6(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	push := Push6[int, int, int, int, int, int](7)
	result := push(t6)
	assert.Equal(t, MakeTuple7(1, 2, 3, 4, 5, 6, 7), result)
}

func TestPush7(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	push := Push7[int, int, int, int, int, int, int](8)
	result := push(t7)
	assert.Equal(t, MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8), result)
}

func TestPush8(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	push := Push8[int, int, int, int, int, int, int, int](9)
	result := push(t8)
	assert.Equal(t, MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9), result)
}

func TestPush9(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	push := Push9[int, int, int, int, int, int, int, int, int](10)
	result := push(t9)
	assert.Equal(t, MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), result)
}

// Test String methods for sizes 5-10
func TestTuple5String(t *testing.T) {
	t5 := MakeTuple5(1, 2, 3, 4, 5)
	assert.Equal(t, "Tuple5[int, int, int, int, int](1, 2, 3, 4, 5)", t5.String())
}

func TestTuple6String(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	assert.Equal(t, "Tuple6[int, int, int, int, int, int](1, 2, 3, 4, 5, 6)", t6.String())
}

func TestTuple7String(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	assert.Equal(t, "Tuple7[int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7)", t7.String())
}

func TestTuple8String(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	assert.Equal(t, "Tuple8[int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8)", t8.String())
}

func TestTuple9String(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.Equal(t, "Tuple9[int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9)", t9.String())
}

func TestTuple10String(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.Equal(t, "Tuple10[int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10)", t10.String())
}

// Test JSON for sizes 6-10
func TestTuple6JSON(t *testing.T) {
	t6 := MakeTuple6(1, 2, 3, 4, 5, 6)
	data, err := json.Marshal(t6)
	require.NoError(t, err)

	var unmarshaled Tuple6[int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t6, unmarshaled)
}

func TestTuple7JSON(t *testing.T) {
	t7 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	data, err := json.Marshal(t7)
	require.NoError(t, err)

	var unmarshaled Tuple7[int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t7, unmarshaled)
}

func TestTuple8JSON(t *testing.T) {
	t8 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	data, err := json.Marshal(t8)
	require.NoError(t, err)

	var unmarshaled Tuple8[int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t8, unmarshaled)
}

func TestTuple9JSON(t *testing.T) {
	t9 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	data, err := json.Marshal(t9)
	require.NoError(t, err)

	var unmarshaled Tuple9[int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t9, unmarshaled)
}

func TestTuple10JSON(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	data, err := json.Marshal(t10)
	require.NoError(t, err)

	var unmarshaled Tuple10[int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t10, unmarshaled)
}

// Test Monoid for sizes 6-10
func TestMonoid6(t *testing.T) {
	m := Monoid6(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple6(1, 2, 3, 4, 5, 6)
	t2 := MakeTuple6(1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple6(2, 3, 4, 5, 6, 7), result)
}

func TestMonoid7(t *testing.T) {
	m := Monoid7(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	t2 := MakeTuple7(1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple7(2, 3, 4, 5, 6, 7, 8), result)
}

func TestMonoid8(t *testing.T) {
	m := Monoid8(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	t2 := MakeTuple8(1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple8(2, 3, 4, 5, 6, 7, 8, 9), result)
}

func TestMonoid9(t *testing.T) {
	m := Monoid9(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	t2 := MakeTuple9(1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple9(2, 3, 4, 5, 6, 7, 8, 9, 10), result)
}

func TestMonoid10(t *testing.T) {
	m := Monoid10(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	t2 := MakeTuple10(1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple10(2, 3, 4, 5, 6, 7, 8, 9, 10, 11), result)
}

// Test Ord for sizes 6-10
func TestOrd6(t *testing.T) {
	o := Ord6(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple6(1, 2, 3, 4, 5, 6)
	t2 := MakeTuple6(1, 2, 3, 4, 5, 7)
	t3 := MakeTuple6(1, 2, 3, 4, 5, 6)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd7(t *testing.T) {
	o := Ord7(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)
	t2 := MakeTuple7(1, 2, 3, 4, 5, 6, 8)
	t3 := MakeTuple7(1, 2, 3, 4, 5, 6, 7)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd8(t *testing.T) {
	o := Ord8(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
	t2 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 9)
	t3 := MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd9(t *testing.T) {
	o := Ord9(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
	t2 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 10)
	t3 := MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd10(t *testing.T) {
	o := Ord10(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	t2 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 11)
	t3 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

// Test tuple sizes 11-15
func TestMakeTuple11(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	assert.Equal(t, 1, t11.F1)
	assert.Equal(t, 11, t11.F11)
}

func TestMakeTuple12(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equal(t, 1, t12.F1)
	assert.Equal(t, 12, t12.F12)
}

func TestMakeTuple13(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	assert.Equal(t, 1, t13.F1)
	assert.Equal(t, 13, t13.F13)
}

func TestMakeTuple14(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	assert.Equal(t, 1, t14.F1)
	assert.Equal(t, 14, t14.F14)
}

func TestMakeTuple15(t *testing.T) {
	t15 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	assert.Equal(t, 1, t15.F1)
	assert.Equal(t, 15, t15.F15)
}

// Test Tupled/Untupled for sizes 11-15
func TestTupled11Untupled11(t *testing.T) {
	sum11 := func(a, b, c, d, e, f, g, h, i, j, k int) int {
		return a + b + c + d + e + f + g + h + i + j + k
	}

	tupled := Tupled11(sum11)
	result := tupled(MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11))
	assert.Equal(t, 66, result)

	untupled := Untupled11(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 11, result2)
}

func TestTupled12Untupled12(t *testing.T) {
	sum12 := func(a, b, c, d, e, f, g, h, i, j, k, l int) int {
		return a + b + c + d + e + f + g + h + i + j + k + l
	}

	tupled := Tupled12(sum12)
	result := tupled(MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12))
	assert.Equal(t, 78, result)

	untupled := Untupled12(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 12, result2)
}

func TestTupled13Untupled13(t *testing.T) {
	sum13 := func(a, b, c, d, e, f, g, h, i, j, k, l, m int) int {
		return a + b + c + d + e + f + g + h + i + j + k + l + m
	}

	tupled := Tupled13(sum13)
	result := tupled(MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13))
	assert.Equal(t, 91, result)

	untupled := Untupled13(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 13, result2)
}

func TestTupled14Untupled14(t *testing.T) {
	sum14 := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n int) int {
		return a + b + c + d + e + f + g + h + i + j + k + l + m + n
	}

	tupled := Tupled14(sum14)
	result := tupled(MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14))
	assert.Equal(t, 105, result)

	untupled := Untupled14(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 14, result2)
}

func TestTupled15Untupled15(t *testing.T) {
	sum15 := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o int) int {
		return a + b + c + d + e + f + g + h + i + j + k + l + m + n + o
	}

	tupled := Tupled15(sum15)
	result := tupled(MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15))
	assert.Equal(t, 120, result)

	untupled := Untupled15(tupled)
	result2 := untupled(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	assert.Equal(t, 15, result2)
}

// Test Map for sizes 11-15
func TestMap11(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	mapper := Map11(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t11)
	assert.Equal(t, MakeTuple11(2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22), result)
}

func TestMap12(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	mapper := Map12(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t12)
	assert.Equal(t, MakeTuple12(2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24), result)
}

func TestMap13(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	mapper := Map13(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t13)
	assert.Equal(t, MakeTuple13(2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26), result)
}

func TestMap14(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	mapper := Map14(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t14)
	assert.Equal(t, MakeTuple14(2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28), result)
}

func TestMap15(t *testing.T) {
	t15 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	mapper := Map15(
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
		N.Mul(2),
	)
	result := mapper(t15)
	assert.Equal(t, MakeTuple15(2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30), result)
}

// Test Replicate for sizes 11-15
func TestReplicate11(t *testing.T) {
	result := Replicate11(21)
	assert.Equal(t, MakeTuple11(21, 21, 21, 21, 21, 21, 21, 21, 21, 21, 21), result)
}

func TestReplicate12(t *testing.T) {
	result := Replicate12(23)
	assert.Equal(t, MakeTuple12(23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23, 23), result)
}

func TestReplicate13(t *testing.T) {
	result := Replicate13(25)
	assert.Equal(t, MakeTuple13(25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25, 25), result)
}

func TestReplicate14(t *testing.T) {
	result := Replicate14(27)
	assert.Equal(t, MakeTuple14(27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27, 27), result)
}

func TestReplicate15(t *testing.T) {
	result := Replicate15(29)
	assert.Equal(t, MakeTuple15(29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29, 29), result)
}

// Test ToArray/FromArray for sizes 11-15
func TestToArray11FromArray11(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	toArray := ToArray11(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t11)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}, arr)

	fromArray := FromArray11(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t11, result)
}

func TestToArray12FromArray12(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	toArray := ToArray12(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t12)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}, arr)

	fromArray := FromArray12(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t12, result)
}

func TestToArray13FromArray13(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	toArray := ToArray13(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t13)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, arr)

	fromArray := FromArray13(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t13, result)
}

func TestToArray14FromArray14(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	toArray := ToArray14(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t14)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14}, arr)

	fromArray := FromArray14(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t14, result)
}

func TestToArray15FromArray15(t *testing.T) {
	t15 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	toArray := ToArray15(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	arr := toArray(t15)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15}, arr)

	fromArray := FromArray15(
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
		func(n int) int { return n },
	)
	result := fromArray(arr)
	assert.Equal(t, t15, result)
}

// Test Push for sizes 10-14
func TestPush10(t *testing.T) {
	t10 := MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	push := Push10[int, int, int, int, int, int, int, int, int, int](11)
	result := push(t10)
	assert.Equal(t, MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11), result)
}

func TestPush11(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	push := Push11[int, int, int, int, int, int, int, int, int, int, int](12)
	result := push(t11)
	assert.Equal(t, MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12), result)
}

func TestPush12(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	push := Push12[int, int, int, int, int, int, int, int, int, int, int, int](13)
	result := push(t12)
	assert.Equal(t, MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), result)
}

func TestPush13(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	push := Push13[int, int, int, int, int, int, int, int, int, int, int, int, int](14)
	result := push(t13)
	assert.Equal(t, MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14), result)
}

func TestPush14(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	push := Push14[int, int, int, int, int, int, int, int, int, int, int, int, int, int](15)
	result := push(t14)
	assert.Equal(t, MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15), result)
}

// Test String methods for sizes 11-15
func TestTuple11String(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	assert.Equal(t, "Tuple11[int, int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)", t11.String())
}

func TestTuple12String(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equal(t, "Tuple12[int, int, int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)", t12.String())
}

func TestTuple13String(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	assert.Equal(t, "Tuple13[int, int, int, int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)", t13.String())
}

func TestTuple14String(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	assert.Equal(t, "Tuple14[int, int, int, int, int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)", t14.String())
}

func TestTuple15String(t *testing.T) {
	t15 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	assert.Equal(t, "Tuple15[int, int, int, int, int, int, int, int, int, int, int, int, int, int, int](1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)", t15.String())
}

// Test JSON for sizes 11-15
func TestTuple11JSON(t *testing.T) {
	t11 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	data, err := json.Marshal(t11)
	require.NoError(t, err)

	var unmarshaled Tuple11[int, int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t11, unmarshaled)
}

func TestTuple12JSON(t *testing.T) {
	t12 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	data, err := json.Marshal(t12)
	require.NoError(t, err)

	var unmarshaled Tuple12[int, int, int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t12, unmarshaled)
}

func TestTuple13JSON(t *testing.T) {
	t13 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	data, err := json.Marshal(t13)
	require.NoError(t, err)

	var unmarshaled Tuple13[int, int, int, int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t13, unmarshaled)
}

func TestTuple14JSON(t *testing.T) {
	t14 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	data, err := json.Marshal(t14)
	require.NoError(t, err)

	var unmarshaled Tuple14[int, int, int, int, int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t14, unmarshaled)
}

func TestTuple15JSON(t *testing.T) {
	t15 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	data, err := json.Marshal(t15)
	require.NoError(t, err)

	var unmarshaled Tuple15[int, int, int, int, int, int, int, int, int, int, int, int, int, int, int]
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.Equal(t, t15, unmarshaled)
}

// Test Monoid for sizes 11-15
func TestMonoid11(t *testing.T) {
	m := Monoid11(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	t2 := MakeTuple11(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple11(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12), result)
}

func TestMonoid12(t *testing.T) {
	m := Monoid12(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	t2 := MakeTuple12(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple12(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13), result)
}

func TestMonoid13(t *testing.T) {
	m := Monoid13(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	t2 := MakeTuple13(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple13(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14), result)
}

func TestMonoid14(t *testing.T) {
	m := Monoid14(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	t2 := MakeTuple14(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple14(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15), result)
}

func TestMonoid15(t *testing.T) {
	m := Monoid15(
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
		N.MonoidSum[int](),
	)
	t1 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	t2 := MakeTuple15(1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1)
	result := m.Concat(t1, t2)
	assert.Equal(t, MakeTuple15(2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16), result)
}

// Test Ord for sizes 11-15
func TestOrd11(t *testing.T) {
	o := Ord11(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	t2 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 12)
	t3 := MakeTuple11(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd12(t *testing.T) {
	o := Ord12(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	t2 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 13)
	t3 := MakeTuple12(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd13(t *testing.T) {
	o := Ord13(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	t2 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14)
	t3 := MakeTuple13(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd14(t *testing.T) {
	o := Ord14(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	t2 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 15)
	t3 := MakeTuple14(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}

func TestOrd15(t *testing.T) {
	o := Ord15(
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
		O.FromStrictCompare[int](),
	)
	t1 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	t2 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 16)
	t3 := MakeTuple15(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)

	assert.Equal(t, -1, o.Compare(t1, t2))
	assert.Equal(t, 0, o.Compare(t1, t3))
	assert.True(t, o.Equals(t1, t3))
}
