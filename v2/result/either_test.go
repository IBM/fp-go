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
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError := Left[string](err)

	assert.True(t, IsLeft(withError))
	assert.False(t, IsRight(withError))
}

func TestIsRight(t *testing.T) {
	noError := Right("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {
	e := errors.New("s")
	assert.Equal(t, F.Pipe1(Right("abc"), Map(utils.StringLen)), Right(3))

	val2 := F.Pipe1(Left[string](e), Map(utils.StringLen))
	exp2 := Left[int](e)

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

	s := S.Semigroup

	assert.Equal(t, "foobar", F.Pipe1(Right("bar"), Reduce(s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string](errors.New("bar")), Reduce(s.Concat, "foo")))

}
func TestAp(t *testing.T) {
	f := S.Size

	maError := errors.New("maError")
	mabError := errors.New("mabError")

	assert.Equal(t, Right(3), F.Pipe1(Right(f), Ap[int](Right("abc"))))
	assert.Equal(t, Left[int](maError), F.Pipe1(Right(f), Ap[int](Left[string](maError))))
	assert.Equal(t, Left[int](mabError), F.Pipe1(Left[func(string) int](mabError), Ap[int](Left[string](maError))))
}

func TestAlt(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Left[int](a)))))
	assert.Equal(t, Right(2), F.Pipe1(Left[int](b), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Left[int](b), F.Pipe1(Left[int](a), Alt(F.Constant(Left[int](b)))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[int])
	maError := errors.New("maError")

	assert.Equal(t, Right("abc"), F.Pipe1(Right("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string](maError), F.Pipe1(Left[string](maError), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")

	f := ChainOptionK[int, int](F.Constant(a))(func(n int) Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right(1), f(Right(1)))
	assert.Equal(t, Left[int](a), f(Right(-1)))
	assert.Equal(t, Left[int](b), f(Left[int](b)))
}

func TestFromOption(t *testing.T) {
	none := errors.New("none")

	assert.Equal(t, Left[int](none), FromOption[int](F.Constant(none))(O.None[int]()))
	assert.Equal(t, Right(1), FromOption[int](F.Constant(none))(O.Some(1)))
}

func TestStringer(t *testing.T) {
	e := Of("foo")
	exp := "Right[string](foo)"

	assert.Equal(t, exp, e.String())

	var s fmt.Stringer = &e
	assert.Equal(t, exp, s.String())
}

func TestFromIO(t *testing.T) {
	f := IO.Of("abc")
	e := FromIO(f)

	assert.Equal(t, Right("abc"), e)
}

// TestOrElse tests recovery from error
func TestOrElse(t *testing.T) {
	// Test basic recovery from Left
	recover := OrElse(func(err error) Result[int] {
		return Right(0) // default value
	})

	leftResult := Left[int](errors.New("fail"))
	assert.Equal(t, Right(0), recover(leftResult))

	// Test that Right values pass through unchanged
	rightResult := Right(42)
	assert.Equal(t, Right(42), recover(rightResult))

	// Test conditional recovery
	recoverSpecific := OrElse(func(err error) Result[int] {
		if err.Error() == "not found" {
			return Right(0) // default for not found
		}
		return Left[int](err) // propagate other errors
	})

	notFoundErr := errors.New("not found")
	assert.Equal(t, Right(0), recoverSpecific(Left[int](notFoundErr)))

	otherErr := errors.New("other error")
	assert.Equal(t, Left[int](otherErr), recoverSpecific(Left[int](otherErr)))
}

// TestZeroEqualsDefaultInitialization tests that Zero returns the same value as default initialization
func TestZeroEqualsDefaultInitialization(t *testing.T) {
	// Default initialization of Result
	var defaultInit Result[int]

	// Zero function
	zero := Zero[int]()

	// They should be equal
	assert.Equal(t, defaultInit, zero, "Zero should equal default initialization")
	assert.Equal(t, IsRight(defaultInit), IsRight(zero), "Both should be Right")
	assert.Equal(t, IsLeft(defaultInit), IsLeft(zero), "Both should not be Left")
}
