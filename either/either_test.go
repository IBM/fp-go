// Copyright (c) 2023 IBM Corp.
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

	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	O "github.com/IBM/fp-go/option"
	S "github.com/IBM/fp-go/string"
	"github.com/stretchr/testify/assert"
)

func TestDefault(t *testing.T) {
	var e Either[error, string]

	assert.Equal(t, Of[error](""), e)
}

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError := Left[string](err)

	assert.True(t, IsLeft(withError))
	assert.False(t, IsRight(withError))
}

func TestIsRight(t *testing.T) {
	noError := Right[error]("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {

	assert.Equal(t, F.Pipe1(Right[error]("abc"), Map[error](utils.StringLen)), Right[error](3))

	val2 := F.Pipe1(Left[string, string]("s"), Map[string](utils.StringLen))
	exp2 := Left[int]("s")

	assert.Equal(t, val2, exp2)
}

func TestUnwrapError(t *testing.T) {
	a := ""
	err := errors.New("Some error")
	withError := Left[string](err)

	res, extracted := UnwrapError(withError)
	assert.Equal(t, a, res)
	assert.Equal(t, extracted, err)

}

func TestReduce(t *testing.T) {

	s := S.Semigroup()

	assert.Equal(t, "foobar", F.Pipe1(Right[string]("bar"), Reduce[string](s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string, string]("bar"), Reduce[string](s.Concat, "foo")))

}

func TestAp(t *testing.T) {
	f := S.Size

	assert.Equal(t, Right[string](3), F.Pipe1(Right[string](f), Ap[int, string, string](Right[string]("abc"))))
	assert.Equal(t, Left[int]("maError"), F.Pipe1(Right[string](f), Ap[int, string, string](Left[string, string]("maError"))))
	assert.Equal(t, Left[int]("mabError"), F.Pipe1(Left[func(string) int]("mabError"), Ap[int, string, string](Left[string, string]("maError"))))
}

func TestAlt(t *testing.T) {
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Right[string](1), F.Pipe1(Right[string](1), Alt(F.Constant(Left[int]("a")))))
	assert.Equal(t, Right[string](2), F.Pipe1(Left[int]("b"), Alt(F.Constant(Right[string](2)))))
	assert.Equal(t, Left[int]("b"), F.Pipe1(Left[int]("a"), Alt(F.Constant(Left[int]("b")))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[string, int])

	assert.Equal(t, Right[string]("abc"), F.Pipe1(Right[string]("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string, string]("maError"), F.Pipe1(Left[string, string]("maError"), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("a"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right[string](1), f(Right[string](1)))
	assert.Equal(t, Left[int]("a"), f(Right[string](-1)))
	assert.Equal(t, Left[int]("b"), f(Left[int]("b")))
}

func TestFromOption(t *testing.T) {
	assert.Equal(t, Left[int]("none"), FromOption[int](F.Constant("none"))(O.None[int]()))
	assert.Equal(t, Right[string](1), FromOption[int](F.Constant("none"))(O.Some(1)))
}
