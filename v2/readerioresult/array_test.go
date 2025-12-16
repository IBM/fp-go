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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	TST "github.com/IBM/fp-go/v2/internal/testing"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {

	e := errors.New("e")

	f := TraverseArray(func(a string) ReaderIOResult[context.Context, string] {
		if S.IsNonEmpty(a) {
			return Right[context.Context](a + a)
		}
		return Left[context.Context, string](e)
	})
	ctx := context.Background()
	assert.Equal(t, result.Of(A.Empty[string]()), F.Pipe1(A.Empty[string](), f)(ctx)())
	assert.Equal(t, result.Of([]string{"aa", "bb"}), F.Pipe1([]string{"a", "b"}, f)(ctx)())
	assert.Equal(t, result.Left[[]string](e), F.Pipe1([]string{"a", ""}, f)(ctx)())
}

func TestSequenceArray(t *testing.T) {

	s := TST.SequenceArrayTest(
		FromStrictEquals[context.Context, bool]()(context.Background()),
		Pointed[context.Context, string](),
		Pointed[context.Context, bool](),
		Functor[context.Context, []string, bool](),
		SequenceArray[context.Context, string],
	)

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("TestSequenceArray %d", i), s(i))
	}
}

func TestSequenceArrayError(t *testing.T) {

	s := TST.SequenceArrayErrorTest(
		FromStrictEquals[context.Context, bool]()(context.Background()),
		Left[context.Context],
		Left[context.Context, bool],
		Pointed[context.Context, string](),
		Pointed[context.Context, bool](),
		Functor[context.Context, []string, bool](),
		SequenceArray[context.Context, string],
	)
	// run across four bits
	s(4)(t)
}

func TestMonadReduceArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	readers := []ReaderIOResult[Config, int]{
		Of[Config](11),
		Of[Config](12),
		Of[Config](13),
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadReduceArray(readers, sum, 0)
	res := r(config)()

	assert.Equal(t, result.Of(36), res) // 11 + 12 + 13
}

func TestMonadReduceArrayWithError(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	testErr := errors.New("test error")
	readers := []ReaderIOResult[Config, int]{
		Of[Config](11),
		Left[Config, int](testErr),
		Of[Config](13),
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadReduceArray(readers, sum, 0)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	val, err := result.Unwrap(res)
	assert.Equal(t, 0, val)
	assert.Equal(t, testErr, err)
}

func TestReduceArray(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 5}

	product := func(acc, val int) int { return acc * val }
	reducer := ReduceArray[Config](product, 1)

	readers := []ReaderIOResult[Config, int]{
		Of[Config](10),
		Of[Config](15),
	}

	r := reducer(readers)
	res := r(config)()

	assert.Equal(t, result.Of(150), res) // 10 * 15
}

func TestReduceArrayWithError(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 5}

	testErr := errors.New("multiplication error")
	product := func(acc, val int) int { return acc * val }
	reducer := ReduceArray[Config](product, 1)

	readers := []ReaderIOResult[Config, int]{
		Of[Config](10),
		Left[Config, int](testErr),
	}

	r := reducer(readers)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestMonadReduceArrayM(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	readers := []ReaderIOResult[Config, int]{
		Of[Config](5),
		Of[Config](10),
		Of[Config](15),
	}

	intAddMonoid := N.MonoidSum[int]()
	r := MonadReduceArrayM(readers, intAddMonoid)
	res := r(config)()

	assert.Equal(t, result.Of(30), res) // 5 + 10 + 15
}

func TestMonadReduceArrayMWithError(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	testErr := errors.New("monoid error")
	readers := []ReaderIOResult[Config, int]{
		Of[Config](5),
		Left[Config, int](testErr),
		Of[Config](15),
	}

	intAddMonoid := N.MonoidSum[int]()
	r := MonadReduceArrayM(readers, intAddMonoid)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestReduceArrayM(t *testing.T) {
	type Config struct{ Scale int }
	config := Config{Scale: 3}

	intMultMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)
	reducer := ReduceArrayM[Config](intMultMonoid)

	readers := []ReaderIOResult[Config, int]{
		Of[Config](3),
		Of[Config](6),
	}

	r := reducer(readers)
	res := r(config)()

	assert.Equal(t, result.Of(18), res) // 3 * 6
}

func TestReduceArrayMWithError(t *testing.T) {
	type Config struct{ Scale int }
	config := Config{Scale: 3}

	testErr := errors.New("scale error")
	intMultMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)
	reducer := ReduceArrayM[Config](intMultMonoid)

	readers := []ReaderIOResult[Config, int]{
		Of[Config](3),
		Left[Config, int](testErr),
	}

	r := reducer(readers)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestMonadTraverseReduceArray(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 10}

	numbers := []int{1, 2, 3, 4}
	multiply := func(n int) ReaderIOResult[Config, int] {
		return Of[Config](n * 10)
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadTraverseReduceArray(numbers, multiply, sum, 0)
	res := r(config)()

	assert.Equal(t, result.Of(100), res) // 10 + 20 + 30 + 40
}

func TestMonadTraverseReduceArrayWithError(t *testing.T) {
	type Config struct{ Multiplier int }
	config := Config{Multiplier: 10}

	testErr := errors.New("transform error")
	numbers := []int{1, 2, 3, 4}
	multiply := func(n int) ReaderIOResult[Config, int] {
		if n == 3 {
			return Left[Config, int](testErr)
		}
		return Of[Config](n * 10)
	}

	sum := func(acc, val int) int { return acc + val }
	r := MonadTraverseReduceArray(numbers, multiply, sum, 0)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestTraverseReduceArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	addBase := func(n int) ReaderIOResult[Config, int] {
		return Of[Config](n + 10)
	}

	product := func(acc, val int) int { return acc * val }
	transformer := TraverseReduceArray(addBase, product, 1)

	r := transformer([]int{2, 3, 4})
	res := r(config)()

	assert.Equal(t, result.Of(2184), res) // 12 * 13 * 14
}

func TestTraverseReduceArrayWithError(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	testErr := errors.New("addition error")
	addBase := func(n int) ReaderIOResult[Config, int] {
		if n == 3 {
			return Left[Config, int](testErr)
		}
		return Of[Config](n + 10)
	}

	product := func(acc, val int) int { return acc * val }
	transformer := TraverseReduceArray(addBase, product, 1)

	r := transformer([]int{2, 3, 4})
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestMonadTraverseReduceArrayM(t *testing.T) {
	type Config struct{ Offset int }
	config := Config{Offset: 100}

	numbers := []int{1, 2, 3}
	addOffset := func(n int) ReaderIOResult[Config, int] {
		return Of[Config](n + 100)
	}

	intSumMonoid := N.MonoidSum[int]()
	r := MonadTraverseReduceArrayM(numbers, addOffset, intSumMonoid)
	res := r(config)()

	assert.Equal(t, result.Of(306), res) // 101 + 102 + 103
}

func TestMonadTraverseReduceArrayMWithError(t *testing.T) {
	type Config struct{ Offset int }
	config := Config{Offset: 100}

	testErr := errors.New("offset error")
	numbers := []int{1, 2, 3}
	addOffset := func(n int) ReaderIOResult[Config, int] {
		if n == 2 {
			return Left[Config, int](testErr)
		}
		return Of[Config](n + 100)
	}

	intSumMonoid := N.MonoidSum[int]()
	r := MonadTraverseReduceArrayM(numbers, addOffset, intSumMonoid)
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestTraverseReduceArrayM(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	scale := func(n int) ReaderIOResult[Config, int] {
		return Of[Config](n * 5)
	}

	intProdMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)
	transformer := TraverseReduceArrayM(scale, intProdMonoid)
	r := transformer([]int{2, 3, 4})
	res := r(config)()

	assert.Equal(t, result.Of(3000), res) // 10 * 15 * 20
}

func TestTraverseReduceArrayMWithError(t *testing.T) {
	type Config struct{ Factor int }
	config := Config{Factor: 5}

	testErr := errors.New("scaling error")
	scale := func(n int) ReaderIOResult[Config, int] {
		if n == 3 {
			return Left[Config, int](testErr)
		}
		return Of[Config](n * 5)
	}

	intProdMonoid := M.MakeMonoid(func(a, b int) int { return a * b }, 1)
	transformer := TraverseReduceArrayM(scale, intProdMonoid)
	r := transformer([]int{2, 3, 4})
	res := r(config)()

	assert.True(t, result.IsLeft(res))
	_, err := result.Unwrap(res)
	assert.Equal(t, testErr, err)
}

func TestReduceArrayEmptyArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	sum := func(acc, val int) int { return acc + val }
	reducer := ReduceArray[Config](sum, 100)

	readers := []ReaderIOResult[Config, int]{}
	r := reducer(readers)
	res := r(config)()

	assert.Equal(t, result.Of(100), res) // Should return initial value
}

func TestTraverseReduceArrayEmptyArray(t *testing.T) {
	type Config struct{ Base int }
	config := Config{Base: 10}

	addBase := func(n int) ReaderIOResult[Config, int] {
		return Of[Config](n + 10)
	}

	sum := func(acc, val int) int { return acc + val }
	transformer := TraverseReduceArray(addBase, sum, 50)

	r := transformer([]int{})
	res := r(config)()

	assert.Equal(t, result.Of(50), res) // Should return initial value
}
