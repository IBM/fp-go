// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, result.Of(2), F.Pipe1(
		Of(1),
		Map(utils.Double),
	)())

}

func TestChainEitherK(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	f := ChainEitherK(func(n int) Result[int] {
		if n > 0 {
			return result.Of(n)
		}
		return result.Left[int](a)

	})
	assert.Equal(t, result.Right(1), f(Right(1))())
	assert.Equal(t, result.Left[int](a), f(Right(-1))())
	assert.Equal(t, result.Left[int](b), f(Left[int](b))())
}

func TestChainOptionK(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	f := ChainOptionK[int, int](F.Constant(a))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})

	assert.Equal(t, result.Right(1), f(Right(1))())
	assert.Equal(t, result.Left[int](a), f(Right(-1))())
	assert.Equal(t, result.Left[int](b), f(Left[int](b))())
}

func TestFromOption(t *testing.T) {

	a := errors.New("a")

	f := FromOption[int](F.Constant(a))
	assert.Equal(t, result.Right(1), f(O.Some(1))())
	assert.Equal(t, result.Left[int](a), f(O.None[int]())())
}

func TestChainIOK(t *testing.T) {
	b := errors.New("b")

	f := ChainIOK(func(n int) io.IO[string] {
		return func() string {
			return fmt.Sprintf("%d", n)
		}
	})

	assert.Equal(t, result.Right("1"), f(Right(1))())
	assert.Equal(t, result.Left[string](b), f(Left[int](b))())
}

func TestChainWithIO(t *testing.T) {

	r := F.Pipe1(
		Of("test"),
		// sad, we need the generics version ...
		io.Map(result.IsRight[string]),
	)

	assert.True(t, r())
}

func TestChainFirst(t *testing.T) {

	foo := errors.New("foo")

	f := func(a string) IOResult[int] {
		if len(a) > 2 {
			return Of(len(a))
		}
		return Left[int](foo)
	}
	good := Of("foo")
	bad := Of("a")
	ch := ChainFirst(f)

	assert.Equal(t, result.Of("foo"), F.Pipe1(good, ch)())
	assert.Equal(t, result.Left[string](foo), F.Pipe1(bad, ch)())
}

func TestChainFirstIOK(t *testing.T) {
	f := func(a string) io.IO[int] {
		return io.Of(len(a))
	}
	good := Of("foo")
	ch := ChainFirstIOK(f)

	assert.Equal(t, result.Of("foo"), F.Pipe1(good, ch)())
}

func TestApFirst(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApFirst[string](Of("b")),
	)

	assert.Equal(t, result.Of("a"), x())
}

func TestApSecond(t *testing.T) {

	x := F.Pipe1(
		Of("a"),
		ApSecond[string](Of("b")),
	)

	assert.Equal(t, result.Of("b"), x())
}

func TestOrElse(t *testing.T) {
	// Test basic recovery from Left
	recover := OrElse(func(e error) IOResult[int] {
		return Right(0)
	})

	res := recover(Left[int](fmt.Errorf("error")))()
	assert.Equal(t, result.Of(0), res)

	// Test Right value passes through unchanged
	res = recover(Right(42))()
	assert.Equal(t, result.Of(42), res)

	// Test selective recovery - recover some errors, propagate others
	selectiveRecover := OrElse(func(err error) IOResult[int] {
		if err.Error() == "not found" {
			return Right(0) // default value for "not found"
		}
		return Left[int](err) // propagate other errors
	})
	notFoundResult := selectiveRecover(Left[int](fmt.Errorf("not found")))()
	assert.Equal(t, result.Of(0), notFoundResult)

	permissionErr := fmt.Errorf("permission denied")
	permissionResult := selectiveRecover(Left[int](permissionErr))()
	assert.Equal(t, result.Left[int](permissionErr), permissionResult)

	// Test chaining multiple OrElse operations
	firstRecover := OrElse(func(err error) IOResult[int] {
		if err.Error() == "error1" {
			return Right(1)
		}
		return Left[int](err)
	})
	secondRecover := OrElse(func(err error) IOResult[int] {
		if err.Error() == "error2" {
			return Right(2)
		}
		return Left[int](err)
	})

	result1 := F.Pipe1(Left[int](fmt.Errorf("error1")), firstRecover)()
	assert.Equal(t, result.Of(1), result1)

	result2 := F.Pipe1(Left[int](fmt.Errorf("error2")), F.Flow2(firstRecover, secondRecover))()
	assert.Equal(t, result.Of(2), result2)
}
