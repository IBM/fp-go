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

package ioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	I "github.com/IBM/fp-go/io"
	IG "github.com/IBM/fp-go/io/generic"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, E.Of[error](2), F.Pipe1(
		Of[error](1),
		Map[error](utils.Double),
	)())

}

func TestChainEitherK(t *testing.T) {
	f := ChainEitherK(func(n int) E.Either[string, int] {
		if n > 0 {
			return E.Of[string](n)
		}
		return E.Left[int]("a")

	})
	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("a"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})

	assert.Equal(t, E.Right[string](1), f(Right[string](1))())
	assert.Equal(t, E.Left[int]("a"), f(Right[string](-1))())
	assert.Equal(t, E.Left[int]("b"), f(Left[int]("b"))())
}

func TestFromOption(t *testing.T) {
	f := FromOption[int](F.Constant("a"))
	assert.Equal(t, E.Right[string](1), f(O.Some(1))())
	assert.Equal(t, E.Left[int]("a"), f(O.None[int]())())
}

func TestChainIOK(t *testing.T) {
	f := ChainIOK[string](func(n int) I.IO[string] {
		return I.MakeIO(func() string {
			return fmt.Sprintf("%d", n)
		})
	})

	assert.Equal(t, E.Right[string]("1"), f(Right[string](1))())
	assert.Equal(t, E.Left[string, string]("b"), f(Left[int]("b"))())
}

func TestChainWithIO(t *testing.T) {

	r := F.Pipe1(
		Of[error]("test"),
		// sad, we need the generics version ...
		IG.Map[IOEither[error, string], I.IO[bool]](E.IsRight[error, string]),
	)

	assert.True(t, r())
}

func TestChainFirst(t *testing.T) {
	f := func(a string) IOEither[string, int] {
		if len(a) > 2 {
			return Of[string](len(a))
		}
		return Left[int]("foo")
	}
	good := Of[string]("foo")
	bad := Of[string]("a")
	ch := ChainFirst(f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
	assert.Equal(t, E.Left[string, string]("foo"), F.Pipe1(bad, ch)())
}

func TestChainFirstIOK(t *testing.T) {
	f := func(a string) I.IO[int] {
		return I.Of(len(a))
	}
	good := Of[string]("foo")
	ch := ChainFirstIOK[string](f)

	assert.Equal(t, E.Of[string]("foo"), F.Pipe1(good, ch)())
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApFirst[string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("a"), x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of[error]("a"),
		ApSecond[string](Of[error]("b")),
	)

	assert.Equal(t, E.Of[error]("b"), x())
}

func TestOrElse(t *testing.T) {
	// Test that OrElse recovers from a Left
	recover := OrElse(func(err string) IOEither[string, int] {
		return Right[string](42)
	})

	// When input is Left, should recover
	leftResult := F.Pipe1(
		Left[int]("error"),
		recover,
	)
	assert.Equal(t, E.Right[string](42), leftResult())

	// When input is Right, should pass through unchanged
	rightResult := F.Pipe1(
		Right[string](100),
		recover,
	)
	assert.Equal(t, E.Right[string](100), rightResult())

	// Test that OrElse can also return a Left (propagate different error)
	recoverOrFail := OrElse(func(err string) IOEither[string, int] {
		if err == "recoverable" {
			return Right[string](0)
		}
		return Left[int]("unrecoverable: " + err)
	})

	recoverable := F.Pipe1(
		Left[int]("recoverable"),
		recoverOrFail,
	)
	assert.Equal(t, E.Right[string](0), recoverable())

	unrecoverable := F.Pipe1(
		Left[int]("fatal"),
		recoverOrFail,
	)
	assert.Equal(t, E.Left[int]("unrecoverable: fatal"), unrecoverable())
}
