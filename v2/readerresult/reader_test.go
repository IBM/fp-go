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

package readerresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

var (
	testError = errors.New("test error")
)

func TestFromEither(t *testing.T) {
	rr := FromEither[MyContext](result.Of(42))
	assert.Equal(t, result.Of(42), rr(defaultContext))

	rrErr := FromEither[MyContext](result.Left[int](testError))
	assert.Equal(t, result.Left[int](testError), rrErr(defaultContext))
}

func TestFromResult(t *testing.T) {
	rr := FromResult[MyContext](result.Of(42))
	assert.Equal(t, result.Of(42), rr(defaultContext))
}

func TestRightReader(t *testing.T) {
	r := func(ctx MyContext) int { return 42 }
	rr := RightReader(r)
	assert.Equal(t, result.Of(42), rr(defaultContext))
}

func TestLeftReader(t *testing.T) {
	r := func(ctx MyContext) error { return testError }
	rr := LeftReader[int](r)
	assert.Equal(t, result.Left[int](testError), rr(defaultContext))
}

func TestLeft(t *testing.T) {
	rr := Left[MyContext, int](testError)
	assert.Equal(t, result.Left[int](testError), rr(defaultContext))
}

func TestRight(t *testing.T) {
	rr := Right[MyContext](42)
	assert.Equal(t, result.Of(42), rr(defaultContext))
}

func TestOf(t *testing.T) {
	rr := Of[MyContext](42)
	assert.Equal(t, result.Of(42), rr(defaultContext))
}

func TestFromReader(t *testing.T) {
	r := func(ctx MyContext) string { return string(ctx) }
	rr := FromReader(r)
	assert.Equal(t, result.Of("default"), rr(defaultContext))
}

func TestMap(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](1),
		Map[MyContext](utils.Double),
	)
	assert.Equal(t, result.Of(2), g(defaultContext))

	// Test with error
	gErr := F.Pipe1(
		Left[MyContext, int](testError),
		Map[MyContext](utils.Double),
	)
	assert.Equal(t, result.Left[int](testError), gErr(defaultContext))
}

func TestMonadMap(t *testing.T) {
	rr := Of[MyContext](5)
	doubled := MonadMap(rr, N.Mul(2))
	assert.Equal(t, result.Of(10), doubled(defaultContext))
}

func TestChain(t *testing.T) {
	addOne := func(x int) ReaderResult[MyContext, int] {
		return Of[MyContext](x + 1)
	}

	g := F.Pipe1(
		Of[MyContext](5),
		Chain(addOne),
	)
	assert.Equal(t, result.Of(6), g(defaultContext))

	// Test error propagation
	gErr := F.Pipe1(
		Left[MyContext, int](testError),
		Chain(addOne),
	)
	assert.Equal(t, result.Left[int](testError), gErr(defaultContext))
}

func TestMonadChain(t *testing.T) {
	addOne := func(x int) ReaderResult[MyContext, int] {
		return Of[MyContext](x + 1)
	}

	rr := Of[MyContext](5)
	res := MonadChain(rr, addOne)
	assert.Equal(t, result.Of(6), res(defaultContext))
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](utils.Double),
		Ap[int](Of[MyContext](1)),
	)
	assert.Equal(t, result.Of(2), g(defaultContext))
}

func TestMonadAp(t *testing.T) {
	add := func(x int) func(int) int {
		return func(y int) int { return x + y }
	}
	fabr := Of[MyContext](add(5))
	fa := Of[MyContext](3)
	res := MonadAp(fabr, fa)
	assert.Equal(t, result.Of(8), res(defaultContext))
}

func TestFromPredicate(t *testing.T) {
	isPositive := FromPredicate[MyContext](
		N.MoreThan(0),
		func(x int) error { return fmt.Errorf("%d is not positive", x) },
	)

	assert.Equal(t, result.Of(5), isPositive(5)(defaultContext))
	res := isPositive(-1)(defaultContext)
	assert.True(t, result.IsLeft(res))
}

func TestFold(t *testing.T) {
	handleError := func(err error) reader.Reader[MyContext, string] {
		return func(ctx MyContext) string { return "Error: " + err.Error() }
	}
	handleSuccess := func(x int) reader.Reader[MyContext, string] {
		return func(ctx MyContext) string { return fmt.Sprintf("Success: %d", x) }
	}

	fold := Fold(handleError, handleSuccess)

	res1 := fold(Of[MyContext](42))(defaultContext)
	assert.Equal(t, "Success: 42", res1)

	res2 := fold(Left[MyContext, int](testError))(defaultContext)
	assert.Equal(t, "Error: "+testError.Error(), res2)
}

func TestGetOrElse(t *testing.T) {
	defaultVal := func(err error) reader.Reader[MyContext, int] {
		return func(ctx MyContext) int { return 0 }
	}

	getOrElse := GetOrElse(defaultVal)

	res1 := getOrElse(Of[MyContext](42))(defaultContext)
	assert.Equal(t, 42, res1)

	res2 := getOrElse(Left[MyContext, int](testError))(defaultContext)
	assert.Equal(t, 0, res2)
}

func TestOrElse(t *testing.T) {
	fallback := func(err error) ReaderResult[MyContext, int] {
		return Of[MyContext](99)
	}

	orElse := OrElse(fallback)

	res1 := F.Pipe1(Of[MyContext](42), orElse)(defaultContext)
	assert.Equal(t, result.Of(42), res1)

	res2 := F.Pipe1(Left[MyContext, int](testError), orElse)(defaultContext)
	assert.Equal(t, result.Of(99), res2)
}

func TestOrLeft(t *testing.T) {
	enrichErr := func(err error) reader.Reader[MyContext, error] {
		return func(ctx MyContext) error {
			return fmt.Errorf("enriched: %w", err)
		}
	}

	orLeft := OrLeft[MyContext, int](enrichErr)

	res1 := F.Pipe1(Of[MyContext](42), orLeft)(defaultContext)
	assert.Equal(t, result.Of(42), res1)

	res2 := F.Pipe1(Left[MyContext, int](testError), orLeft)(defaultContext)
	assert.True(t, result.IsLeft(res2))
}

func TestAsk(t *testing.T) {
	rr := Ask[MyContext]()
	assert.Equal(t, result.Of(defaultContext), rr(defaultContext))
}

func TestAsks(t *testing.T) {
	getLen := func(ctx MyContext) int { return len(string(ctx)) }
	rr := Asks(getLen)
	assert.Equal(t, result.Of(7), rr(defaultContext)) // "default" has 7 chars
}

func TestChainEitherK(t *testing.T) {
	parseInt := func(s string) result.Result[int] {
		if s == "42" {
			return result.Of(42)
		}
		return result.Left[int](errors.New("parse error"))
	}

	chain := ChainEitherK[MyContext](parseInt)

	res1 := F.Pipe1(Of[MyContext]("42"), chain)(defaultContext)
	assert.Equal(t, result.Of(42), res1)

	res2 := F.Pipe1(Of[MyContext]("invalid"), chain)(defaultContext)
	assert.True(t, result.IsLeft(res2))
}

func TestChainOptionK(t *testing.T) {
	findEven := func(x int) option.Option[int] {
		if x%2 == 0 {
			return option.Some(x)
		}
		return option.None[int]()
	}

	notFound := func() error { return errors.New("not even") }
	chain := ChainOptionK[MyContext, int, int](notFound)(findEven)

	res1 := F.Pipe1(Of[MyContext](4), chain)(defaultContext)
	assert.Equal(t, result.Of(4), res1)

	res2 := F.Pipe1(Of[MyContext](3), chain)(defaultContext)
	assert.True(t, result.IsLeft(res2))
}

func TestFlatten(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](Of[MyContext]("a")),
		Flatten[MyContext, string],
	)
	assert.Equal(t, result.Of("a"), g(defaultContext))
}

func TestBiMap(t *testing.T) {
	enrichErr := func(e error) error { return fmt.Errorf("enriched: %w", e) }
	double := N.Mul(2)

	res1 := F.Pipe1(Of[MyContext](5), BiMap[MyContext](enrichErr, double))(defaultContext)
	assert.Equal(t, result.Of(10), res1)

	res2 := F.Pipe1(Left[MyContext, int](testError), BiMap[MyContext](enrichErr, double))(defaultContext)
	assert.True(t, result.IsLeft(res2))
}

func TestLocal(t *testing.T) {
	type OtherContext int
	toMyContext := func(oc OtherContext) MyContext {
		return MyContext(fmt.Sprintf("ctx-%d", oc))
	}

	rr := Asks(func(ctx MyContext) string { return string(ctx) })
	adapted := Local[string](toMyContext)(rr)

	res := adapted(OtherContext(42))
	assert.Equal(t, result.Of("ctx-42"), res)
}

func TestRead(t *testing.T) {
	rr := Of[MyContext](42)
	read := Read[int](defaultContext)
	res := read(rr)
	assert.Equal(t, result.Of(42), res)
}

func TestFlap(t *testing.T) {
	fabr := Of[MyContext](N.Mul(2))
	flapped := MonadFlap(fabr, 5)
	assert.Equal(t, result.Of(10), flapped(defaultContext))
}

func TestMapLeft(t *testing.T) {
	enrichErr := func(e error) error { return fmt.Errorf("DB error: %w", e) }

	res1 := F.Pipe1(Of[MyContext](42), MapLeft[MyContext, int](enrichErr))(defaultContext)
	assert.Equal(t, result.Of(42), res1)

	res2 := F.Pipe1(Left[MyContext, int](testError), MapLeft[MyContext, int](enrichErr))(defaultContext)
	assert.True(t, result.IsLeft(res2))
}
