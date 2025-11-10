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

package either

import (
	"errors"
	"testing"

	O "github.com/IBM/fp-go/v2/option"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

// Test MonadChainOptionK
func TestMonadChainOptionK(t *testing.T) {
	f := func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n * 2)
		}
		return O.None[int]()
	}

	onNone := func() error { return errors.New("none") }

	result := MonadChainOptionK(onNone, Right[error](5), f)
	assert.Equal(t, Right[error](10), result)

	result = MonadChainOptionK(onNone, Right[error](-1), f)
	assert.Equal(t, Left[int](errors.New("none")), result)

	result = MonadChainOptionK(onNone, Left[int](errors.New("error")), f)
	assert.Equal(t, Left[int](errors.New("error")), result)
}

// Test Uneitherize1
func TestUneitherize1(t *testing.T) {
	f := func(x int) Either[error, string] {
		if x > 0 {
			return Right[error]("positive")
		}
		return Left[string](errors.New("negative"))
	}

	uneitherized := Uneitherize1(f)
	result, err := uneitherized(5)
	assert.NoError(t, err)
	assert.Equal(t, "positive", result)

	_, err = uneitherized(-1)
	assert.Error(t, err)
}

// Test Uneitherize2
func TestUneitherize2(t *testing.T) {
	f := func(x, y int) Either[error, int] {
		if x > 0 && y > 0 {
			return Right[error](x + y)
		}
		return Left[int](errors.New("invalid"))
	}

	uneitherized := Uneitherize2(f)
	result, err := uneitherized(5, 3)
	assert.NoError(t, err)
	assert.Equal(t, 8, result)

	result, err = uneitherized(-1, 3)
	assert.Error(t, err)
}

// Test Uneitherize3
func TestUneitherize3(t *testing.T) {
	f := func(x, y, z int) Either[error, int] {
		return Right[error](x + y + z)
	}

	uneitherized := Uneitherize3(f)
	result, err := uneitherized(1, 2, 3)
	assert.NoError(t, err)
	assert.Equal(t, 6, result)
}

// Test Uneitherize4
func TestUneitherize4(t *testing.T) {
	f := func(a, b, c, d int) Either[error, int] {
		return Right[error](a + b + c + d)
	}

	uneitherized := Uneitherize4(f)
	result, err := uneitherized(1, 2, 3, 4)
	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

// Test SequenceT1
func TestSequenceT1(t *testing.T) {
	result := SequenceT1(
		Right[error](1),
	)
	expected := Right[error](T.MakeTuple1(1))
	assert.Equal(t, expected, result)

	result = SequenceT1(
		Left[int](errors.New("error")),
	)
	assert.True(t, IsLeft(result))
}

// Test SequenceT2
func TestSequenceT2(t *testing.T) {
	result := SequenceT2(
		Right[error](1),
		Right[error]("hello"),
	)
	expected := Right[error](T.MakeTuple2(1, "hello"))
	assert.Equal(t, expected, result)

	result = SequenceT2(
		Left[int](errors.New("error")),
		Right[error]("hello"),
	)
	assert.True(t, IsLeft(result))
}

// Test SequenceT3
func TestSequenceT3(t *testing.T) {
	result := SequenceT3(
		Right[error](1),
		Right[error]("hello"),
		Right[error](true),
	)
	expected := Right[error](T.MakeTuple3(1, "hello", true))
	assert.Equal(t, expected, result)
}

// Test SequenceT4
func TestSequenceT4(t *testing.T) {
	result := SequenceT4(
		Right[error](1),
		Right[error](2),
		Right[error](3),
		Right[error](4),
	)
	expected := Right[error](T.MakeTuple4(1, 2, 3, 4))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple1
func TestSequenceTuple1(t *testing.T) {
	tuple := T.MakeTuple1(Right[error](42))
	result := SequenceTuple1(tuple)
	expected := Right[error](T.MakeTuple1(42))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple2
func TestSequenceTuple2(t *testing.T) {
	tuple := T.MakeTuple2(Right[error](1), Right[error]("hello"))
	result := SequenceTuple2(tuple)
	expected := Right[error](T.MakeTuple2(1, "hello"))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple3
func TestSequenceTuple3(t *testing.T) {
	tuple := T.MakeTuple3(Right[error](1), Right[error](2), Right[error](3))
	result := SequenceTuple3(tuple)
	expected := Right[error](T.MakeTuple3(1, 2, 3))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple4
func TestSequenceTuple4(t *testing.T) {
	tuple := T.MakeTuple4(Right[error](1), Right[error](2), Right[error](3), Right[error](4))
	result := SequenceTuple4(tuple)
	expected := Right[error](T.MakeTuple4(1, 2, 3, 4))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple1
func TestTraverseTuple1(t *testing.T) {
	f := func(x int) Either[error, string] {
		if x > 0 {
			return Right[error]("positive")
		}
		return Left[string](errors.New("negative"))
	}

	tuple := T.MakeTuple1(5)
	result := TraverseTuple1(f)(tuple)
	expected := Right[error](T.MakeTuple1("positive"))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple2
func TestTraverseTuple2(t *testing.T) {
	f1 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f2 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}

	tuple := T.MakeTuple2(1, 2)
	result := TraverseTuple2(f1, f2)(tuple)
	expected := Right[error](T.MakeTuple2(2, 4))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple3
func TestTraverseTuple3(t *testing.T) {
	f1 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f2 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f3 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}

	tuple := T.MakeTuple3(1, 2, 3)
	result := TraverseTuple3(f1, f2, f3)(tuple)
	expected := Right[error](T.MakeTuple3(2, 4, 6))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple4
func TestTraverseTuple4(t *testing.T) {
	f1 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f2 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f3 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f4 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}

	tuple := T.MakeTuple4(1, 2, 3, 4)
	result := TraverseTuple4(f1, f2, f3, f4)(tuple)
	expected := Right[error](T.MakeTuple4(2, 4, 6, 8))
	assert.Equal(t, expected, result)
}

// Test Eitherize5
func TestEitherize5(t *testing.T) {
	f := func(a, b, c, d, e int) (int, error) {
		return a + b + c + d + e, nil
	}

	eitherized := Eitherize5(f)
	result := eitherized(1, 2, 3, 4, 5)
	assert.Equal(t, Right[error](15), result)
}

// Test Uneitherize5
func TestUneitherize5(t *testing.T) {
	f := func(a, b, c, d, e int) Either[error, int] {
		return Right[error](a + b + c + d + e)
	}

	uneitherized := Uneitherize5(f)
	result, err := uneitherized(1, 2, 3, 4, 5)
	assert.NoError(t, err)
	assert.Equal(t, 15, result)
}

// Test SequenceT5
func TestSequenceT5(t *testing.T) {
	result := SequenceT5(
		Right[error](1),
		Right[error](2),
		Right[error](3),
		Right[error](4),
		Right[error](5),
	)
	expected := Right[error](T.MakeTuple5(1, 2, 3, 4, 5))
	assert.Equal(t, expected, result)
}

// Test SequenceTuple5
func TestSequenceTuple5(t *testing.T) {
	tuple := T.MakeTuple5(
		Right[error](1),
		Right[error](2),
		Right[error](3),
		Right[error](4),
		Right[error](5),
	)
	result := SequenceTuple5(tuple)
	expected := Right[error](T.MakeTuple5(1, 2, 3, 4, 5))
	assert.Equal(t, expected, result)
}

// Test TraverseTuple5
func TestTraverseTuple5(t *testing.T) {
	f1 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f2 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f3 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f4 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}
	f5 := func(x int) Either[error, int] {
		return Right[error](x * 2)
	}

	tuple := T.MakeTuple5(1, 2, 3, 4, 5)
	result := TraverseTuple5(f1, f2, f3, f4, f5)(tuple)
	expected := Right[error](T.MakeTuple5(2, 4, 6, 8, 10))
	assert.Equal(t, expected, result)
}

// Test higher arity functions (6-10) - sample tests
func TestEitherize6(t *testing.T) {
	f := func(a, b, c, d, e, f int) (int, error) {
		return a + b + c + d + e + f, nil
	}

	eitherized := Eitherize6(f)
	result := eitherized(1, 2, 3, 4, 5, 6)
	assert.Equal(t, Right[error](21), result)
}

func TestSequenceT6(t *testing.T) {
	result := SequenceT6(
		Right[error](1),
		Right[error](2),
		Right[error](3),
		Right[error](4),
		Right[error](5),
		Right[error](6),
	)
	assert.True(t, IsRight(result))
}

func TestEitherize7(t *testing.T) {
	f := func(a, b, c, d, e, f, g int) (int, error) {
		return a + b + c + d + e + f + g, nil
	}

	eitherized := Eitherize7(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7)
	assert.Equal(t, Right[error](28), result)
}

func TestEitherize8(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h int) (int, error) {
		return a + b + c + d + e + f + g + h, nil
	}

	eitherized := Eitherize8(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8)
	assert.Equal(t, Right[error](36), result)
}

func TestEitherize9(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i int) (int, error) {
		return a + b + c + d + e + f + g + h + i, nil
	}

	eitherized := Eitherize9(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.Equal(t, Right[error](45), result)
}

func TestEitherize10(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j, nil
	}

	eitherized := Eitherize10(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.Equal(t, Right[error](55), result)
}

func TestEitherize11(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j + k, nil
	}

	eitherized := Eitherize11(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	assert.Equal(t, Right[error](66), result)
}

func TestEitherize12(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j + k + l, nil
	}

	eitherized := Eitherize12(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.Equal(t, Right[error](78), result)
}

func TestEitherize13(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j + k + l + m, nil
	}

	eitherized := Eitherize13(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	assert.Equal(t, Right[error](91), result)
}

func TestEitherize14(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j + k + l + m + n, nil
	}

	eitherized := Eitherize14(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	assert.Equal(t, Right[error](105), result)
}

func TestEitherize15(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o int) (int, error) {
		return a + b + c + d + e + f + g + h + i + j + k + l + m + n + o, nil
	}

	eitherized := Eitherize15(f)
	result := eitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	assert.Equal(t, Right[error](120), result)
}

// Test Uneitherize functions for higher arities
func TestUneitherize6(t *testing.T) {
	f := func(a, b, c, d, e, f int) Either[error, int] {
		return Right[error](a + b + c + d + e + f)
	}

	uneitherized := Uneitherize6(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6)
	assert.NoError(t, err)
	assert.Equal(t, 21, result)
}

func TestUneitherize7(t *testing.T) {
	f := func(a, b, c, d, e, f, g int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g)
	}

	uneitherized := Uneitherize7(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7)
	assert.NoError(t, err)
	assert.Equal(t, 28, result)
}

func TestUneitherize8(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h)
	}

	uneitherized := Uneitherize8(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8)
	assert.NoError(t, err)
	assert.Equal(t, 36, result)
}

func TestUneitherize9(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i)
	}

	uneitherized := Uneitherize9(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.NoError(t, err)
	assert.Equal(t, 45, result)
}

func TestUneitherize10(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j)
	}

	uneitherized := Uneitherize10(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.NoError(t, err)
	assert.Equal(t, 55, result)
}

func TestUneitherize11(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j + k)
	}

	uneitherized := Uneitherize11(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11)
	assert.NoError(t, err)
	assert.Equal(t, 66, result)
}

func TestUneitherize12(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j + k + l)
	}

	uneitherized := Uneitherize12(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12)
	assert.NoError(t, err)
	assert.Equal(t, 78, result)
}

func TestUneitherize13(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j + k + l + m)
	}

	uneitherized := Uneitherize13(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13)
	assert.NoError(t, err)
	assert.Equal(t, 91, result)
}

func TestUneitherize14(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j + k + l + m + n)
	}

	uneitherized := Uneitherize14(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14)
	assert.NoError(t, err)
	assert.Equal(t, 105, result)
}

func TestUneitherize15(t *testing.T) {
	f := func(a, b, c, d, e, f, g, h, i, j, k, l, m, n, o int) Either[error, int] {
		return Right[error](a + b + c + d + e + f + g + h + i + j + k + l + m + n + o)
	}

	uneitherized := Uneitherize15(f)
	result, err := uneitherized(1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15)
	assert.NoError(t, err)
	assert.Equal(t, 120, result)
}
