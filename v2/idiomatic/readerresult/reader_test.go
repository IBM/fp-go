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
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	rrErr := FromEither[MyContext](result.Left[int](testError))
	_, err = rrErr(defaultContext)
	assert.Equal(t, testError, err)
}

func TestFromResult(t *testing.T) {
	rr := FromResult[MyContext](42, nil)
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
}

func TestRightReader(t *testing.T) {
	r := func(ctx MyContext) int { return 42 }
	rr := RightReader(r)
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
}

func TestLeftReader(t *testing.T) {
	r := func(ctx MyContext) error { return testError }
	rr := LeftReader[int](r)
	_, err := rr(defaultContext)
	assert.Equal(t, testError, err)
}

func TestLeft(t *testing.T) {
	rr := Left[MyContext, int](testError)
	_, err := rr(defaultContext)
	assert.Equal(t, testError, err)
}

func TestRight(t *testing.T) {
	rr := Right[MyContext](42)
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
}

func TestOf(t *testing.T) {
	rr := Of[MyContext](42)
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
}

func TestFromReader(t *testing.T) {
	r := func(ctx MyContext) string { return string(ctx) }
	rr := FromReader(r)
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, "default", v)
}

func TestMap(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](1),
		Map[MyContext](utils.Double),
	)
	v, err := g(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, v)

	// Test with error
	gErr := F.Pipe1(
		Left[MyContext, int](testError),
		Map[MyContext](utils.Double),
	)
	_, err = gErr(defaultContext)
	assert.Equal(t, testError, err)
}

func TestMonadMap(t *testing.T) {
	rr := Of[MyContext](5)
	doubled := MonadMap(rr, N.Mul(2))
	v, err := doubled(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestChain(t *testing.T) {
	addOne := func(x int) ReaderResult[MyContext, int] {
		return Of[MyContext](x + 1)
	}

	g := F.Pipe1(
		Of[MyContext](5),
		Chain(addOne),
	)
	v, err := g(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 6, v)

	// Test error propagation
	gErr := F.Pipe1(
		Left[MyContext, int](testError),
		Chain(addOne),
	)
	_, err = gErr(defaultContext)
	assert.Equal(t, testError, err)
}

func TestMonadChain(t *testing.T) {
	addOne := func(x int) ReaderResult[MyContext, int] {
		return Of[MyContext](x + 1)
	}

	rr := Of[MyContext](5)
	res := MonadChain(rr, addOne)
	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 6, v)
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](utils.Double),
		Ap[int](Of[MyContext](1)),
	)
	v, err := g(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 2, v)
}

func TestMonadAp(t *testing.T) {
	add := func(x int) func(int) int {
		return func(y int) int { return x + y }
	}
	fabr := Of[MyContext](add(5))
	fa := Of[MyContext](3)
	res := MonadAp(fabr, fa)
	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 8, v)
}

func TestFromPredicate(t *testing.T) {
	isPositive := FromPredicate[MyContext](
		N.MoreThan(0),
		func(x int) error { return fmt.Errorf("%d is not positive", x) },
	)

	v, err := isPositive(5)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 5, v)

	_, err = isPositive(-1)(defaultContext)
	assert.Error(t, err)
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

	v, err := F.Pipe1(Of[MyContext](42), orElse)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	v, err = F.Pipe1(Left[MyContext, int](testError), orElse)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 99, v)
}

func TestOrLeft(t *testing.T) {
	enrichErr := func(err error) reader.Reader[MyContext, error] {
		return func(ctx MyContext) error {
			return fmt.Errorf("enriched: %w", err)
		}
	}

	orLeft := OrLeft[int](enrichErr)

	v, err := F.Pipe1(Of[MyContext](42), orLeft)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	_, err = F.Pipe1(Left[MyContext, int](testError), orLeft)(defaultContext)
	assert.Error(t, err)
}

func TestAsk(t *testing.T) {
	rr := Ask[MyContext]()
	v, err := rr(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, defaultContext, v)
}

func TestAsks(t *testing.T) {
	getLen := func(ctx MyContext) int { return len(string(ctx)) }
	rr := Asks(getLen)
	v, err := rr(defaultContext) // "default" has 7 chars
	assert.NoError(t, err)
	assert.Equal(t, 7, v)
}

func TestChainReaderK(t *testing.T) {
	parseInt := func(s string) (int, error) {
		if s == "42" {
			return 42, nil
		}
		return 0, errors.New("parse error")
	}

	chain := ChainReaderK[MyContext](parseInt)

	v, err := F.Pipe1(Of[MyContext]("42"), chain)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	_, err = F.Pipe1(Of[MyContext]("invalid"), chain)(defaultContext)
	assert.Error(t, err)
}

func TestChainEitherK(t *testing.T) {
	parseInt := func(s string) Result[int] {
		if s == "42" {
			return result.Of(42)
		}
		return result.Left[int](errors.New("parse error"))
	}

	chain := ChainEitherK[MyContext](parseInt)

	v, err := F.Pipe1(Of[MyContext]("42"), chain)(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	_, err = F.Pipe1(Of[MyContext]("invalid"), chain)(defaultContext)
	assert.Error(t, err)
}

func TestChainOptionK(t *testing.T) {
	findEven := func(x int) (int, bool) {
		if x%2 == 0 {
			return x, true
		}
		return 0, false
	}

	notFound := func() error { return errors.New("not even") }
	chain := ChainOptionK[MyContext, int, int](notFound)(findEven)

	res := F.Pipe1(Of[MyContext](4), chain)
	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 4, v)

	res2 := F.Pipe1(Of[MyContext](3), chain)
	_, err = res2(defaultContext)
	assert.Error(t, err)
}

func TestFlatten(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](Of[MyContext]("a")),
		Flatten[MyContext, string],
	)
	v, err := g(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, "a", v)
}

func TestBiMap(t *testing.T) {
	enrichErr := func(e error) error { return fmt.Errorf("enriched: %w", e) }
	double := N.Mul(2)

	res := F.Pipe1(Of[MyContext](5), BiMap[MyContext](enrichErr, double))
	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 10, v)

	res2 := F.Pipe1(Left[MyContext, int](testError), BiMap[MyContext](enrichErr, double))
	_, err = res2(defaultContext)
	assert.Error(t, err)
}

func TestLocal(t *testing.T) {
	type OtherContext int
	toMyContext := func(oc OtherContext) MyContext {
		return MyContext(fmt.Sprintf("ctx-%d", oc))
	}

	rr := Asks(func(ctx MyContext) string { return string(ctx) })
	adapted := Local[string](toMyContext)(rr)

	v, err := adapted(OtherContext(42))
	assert.NoError(t, err)
	assert.Equal(t, "ctx-42", v)
}

func TestRead(t *testing.T) {
	rr := Of[MyContext](42)
	read := Read[int](defaultContext)
	v, err := read(rr)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)
}

func TestFlap(t *testing.T) {
	fabr := Of[MyContext](N.Mul(2))
	flapped := MonadFlap(fabr, 5)
	v, err := flapped(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 10, v)
}

func TestMapLeft(t *testing.T) {
	enrichErr := func(e error) error { return fmt.Errorf("DB error: %w", e) }

	res := F.Pipe1(Of[MyContext](42), MapLeft[MyContext, int](enrichErr))
	v, err := res(defaultContext)
	assert.NoError(t, err)
	assert.Equal(t, 42, v)

	res2 := F.Pipe1(Left[MyContext, int](testError), MapLeft[MyContext, int](enrichErr))
	_, err = res2(defaultContext)
	assert.Error(t, err)
}
