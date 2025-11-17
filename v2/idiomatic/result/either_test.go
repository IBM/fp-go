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

package result

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/option"
	"github.com/IBM/fp-go/v2/internal/utils"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError, e1 := Left[string](err)

	assert.True(t, IsLeft(withError, e1))
	assert.False(t, IsRight(withError, e1))
}

func TestIsRight(t *testing.T) {
	noError, err := Right("Carsten")

	assert.True(t, IsRight(noError, err))
	assert.False(t, IsLeft(noError, err))
}

func TestMapEither(t *testing.T) {

	AssertEq(Pipe2("abc", Of, Map(utils.StringLen)))(Right(3))(t)

	e := errors.New("s")

	AssertEq(Left[int](e))(Pipe2(e, Left[string], Map(utils.StringLen)))(t)
}

func TestAp(t *testing.T) {
	f := S.Size

	maError := errors.New("maError")
	mabError := errors.New("mabError")
	AssertEq(Right(3))(Pipe2(f, Right, Ap[int](Right("abc"))))(t)
	AssertEq(Left[int](maError))(Pipe2(f, Right, Ap[int](Left[string](maError))))(t)
	AssertEq(Left[int](mabError))(Pipe2(mabError, Left[func(string) int], Ap[int](Left[string](maError))))(t)
	AssertEq(Left[int](mabError))(Pipe2(mabError, Left[func(string) int], Ap[int](Right("abc"))))(t)
}

func TestAlt(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	AssertEq(Right(1))(Pipe2(1, Right, Alt(func() (int, error) { return Right(2) })))(t)
	AssertEq(Right(1))(Pipe2(1, Right, Alt(func() (int, error) { return Left[int](a) })))(t)
	AssertEq(Right(2))(Pipe2(b, Left[int], Alt(func() (int, error) { return Right(2) })))(t)
	AssertEq(Left[int](b))(Pipe2(a, Left[int], Alt(func() (int, error) { return Left[int](b) })))(t)
}

func TestChainFirst(t *testing.T) {
	f := func(s string) (int, error) {
		return Of(S.Size((s)))
	}

	maError := errors.New("maError")

	AssertEq(Right("abc"))(Pipe2("abc", Right, ChainFirst(f)))(t)
	AssertEq(Left[string](maError))(Pipe2(maError, Left[string], ChainFirst(f)))(t)
}

func TestChainOptionK(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")
	f := ChainOptionK[int, int](F.Constant(a))(func(n int) (int, bool) {
		if n > 0 {
			return option.Some(n)
		}
		return option.None[int]()
	})
	AssertEq(Right(1))(f(Right(1)))(t)
	AssertEq(Left[int](a))(f(Right(-1)))(t)
	AssertEq(Left[int](b))(f(Left[int](b)))(t)
}

func TestFromOption(t *testing.T) {
	none := errors.New("none")
	AssertEq(Left[int](none))(FromOption[int](F.Constant(none))(option.None[int]()))(t)
	AssertEq(Right(1))(FromOption[int](F.Constant(none))(option.Some(1)))(t)
}

func TestStringer(t *testing.T) {
	e := ToString(Of("foo"))
	exp := "Right[string](foo)"

	assert.Equal(t, exp, e)

}
