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
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RO "github.com/IBM/fp-go/v2/readeroption"
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

func TestOfLazy(t *testing.T) {
	t.Run("evaluates lazy computation ignoring environment", func(t *testing.T) {
		lazyValue := func() int { return 42 }
		rr := OfLazy[MyContext](lazyValue)
		res := rr(defaultContext)
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("defers computation until ReaderResult is executed", func(t *testing.T) {
		executed := false
		lazyComputation := func() string {
			executed = true
			return "computed"
		}
		rr := OfLazy[MyContext](lazyComputation)

		// Computation should not be executed yet
		assert.False(t, executed, "lazy computation should not be executed during ReaderResult creation")

		// Execute the ReaderResult
		res := rr(defaultContext)

		// Now computation should be executed
		assert.True(t, executed, "lazy computation should be executed when ReaderResult runs")
		assert.Equal(t, result.Of("computed"), res)
	})

	t.Run("evaluates lazy computation each time ReaderResult is called", func(t *testing.T) {
		counter := 0
		lazyCounter := func() int {
			counter++
			return counter
		}
		rr := OfLazy[MyContext](lazyCounter)

		// First execution
		res1 := rr(defaultContext)
		assert.Equal(t, result.Of(1), res1)

		// Second execution
		res2 := rr(defaultContext)
		assert.Equal(t, result.Of(2), res2)

		// Third execution
		res3 := rr(defaultContext)
		assert.Equal(t, result.Of(3), res3)
	})

	t.Run("works with different types", func(t *testing.T) {
		lazyString := func() string { return "hello" }
		rr1 := OfLazy[MyContext](lazyString)
		assert.Equal(t, result.Of("hello"), rr1(defaultContext))

		lazySlice := func() []int { return []int{1, 2, 3} }
		rr2 := OfLazy[MyContext](lazySlice)
		assert.Equal(t, result.Of([]int{1, 2, 3}), rr2(defaultContext))

		lazyStruct := func() MyContext { return "test" }
		rr3 := OfLazy[string](lazyStruct)
		assert.Equal(t, result.Of(MyContext("test")), rr3("ignored"))
	})

	t.Run("can be composed with other ReaderResult operations", func(t *testing.T) {
		lazyValue := func() int { return 10 }
		rr := F.Pipe1(
			OfLazy[MyContext](lazyValue),
			Map[MyContext](N.Mul(2)),
		)
		res := rr(defaultContext)
		assert.Equal(t, result.Of(20), res)
	})

	t.Run("ignores environment completely", func(t *testing.T) {
		lazyValue := func() string { return "constant" }
		rr := OfLazy[MyContext](lazyValue)

		// Different environments should produce same result
		ctx1 := MyContext("context1")
		ctx2 := MyContext("context2")

		assert.Equal(t, result.Of("constant"), rr(ctx1))
		assert.Equal(t, result.Of("constant"), rr(ctx2))
	})

	t.Run("always wraps result in success", func(t *testing.T) {
		lazyValue := func() int { return 42 }
		rr := OfLazy[MyContext](lazyValue)
		res := rr(defaultContext)

		// Verify it's a successful Result
		assert.True(t, result.IsRight(res))
		assert.Equal(t, result.Of(42), res)
	})
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

func TestFromOption(t *testing.T) {
	onNone := lazy.Of(errors.New("not found"))

	t.Run("Some yields Right", func(t *testing.T) {
		res := FromOption[MyContext, int](onNone)(option.Some(42))(defaultContext)
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("None yields Left with onNone error", func(t *testing.T) {
		res := FromOption[MyContext, int](onNone)(option.None[int]())(defaultContext)
		assert.Equal(t, result.Left[int](errors.New("not found")), res)
	})

	t.Run("environment is ignored for Some", func(t *testing.T) {
		const ctx1 MyContext = "ctx1"
		const ctx2 MyContext = "ctx2"
		lift := FromOption[MyContext, int](onNone)(option.Some(7))
		assert.Equal(t, lift(ctx1), lift(ctx2))
	})

	t.Run("environment is ignored for None", func(t *testing.T) {
		const ctx1 MyContext = "ctx1"
		const ctx2 MyContext = "ctx2"
		lift := FromOption[MyContext, int](onNone)(option.None[int]())
		assert.Equal(t, lift(ctx1), lift(ctx2))
	})

	t.Run("composition with Map", func(t *testing.T) {
		res := F.Pipe1(
			FromOption[MyContext, int](onNone)(option.Some(3)),
			Map[MyContext](utils.Double),
		)(defaultContext)
		assert.Equal(t, result.Of(6), res)
	})
}

func TestFromReaderOption(t *testing.T) {
	onNone := lazy.Of(errors.New("not found"))

	t.Run("Some-yielding ReaderOption gives Right", func(t *testing.T) {
		ro := RO.Of[MyContext](42)
		res := FromReaderOption[MyContext, int](onNone)(ro)(defaultContext)
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("None-yielding ReaderOption gives Left", func(t *testing.T) {
		ro := RO.None[MyContext, int]()
		res := FromReaderOption[MyContext, int](onNone)(ro)(defaultContext)
		assert.Equal(t, result.Left[int](errors.New("not found")), res)
	})

	t.Run("environment is forwarded to the ReaderOption", func(t *testing.T) {
		type Config struct{ port int }
		ro := RO.Asks(func(cfg Config) int { return cfg.port })
		res := FromReaderOption[Config, int](onNone)(ro)(Config{port: 8080})
		assert.Equal(t, result.Of(8080), res)
	})

	t.Run("None with environment-dependent ReaderOption gives Left", func(t *testing.T) {
		type Config struct{ port int }
		ro := func(cfg Config) option.Option[int] { return option.None[int]() }
		res := FromReaderOption[Config, int](onNone)(ro)(Config{port: 9})
		assert.Equal(t, result.Left[int](errors.New("not found")), res)
	})

	t.Run("composition with Map", func(t *testing.T) {
		ro := RO.Of[MyContext](5)
		res := F.Pipe1(
			FromReaderOption[MyContext, int](onNone)(ro),
			Map[MyContext](utils.Double),
		)(defaultContext)
		assert.Equal(t, result.Of(10), res)
	})
}

func TestMonadChainFirst(t *testing.T) {
	t.Run("success: preserves original value and executes side effect", func(t *testing.T) {
		var sideEffect int
		record := func(x int) ReaderResult[MyContext, string] {
			sideEffect = x
			return Of[MyContext]("logged")
		}

		rr := Of[MyContext](42)
		res := MonadChainFirst(rr, record)(defaultContext)

		assert.Equal(t, result.Of(42), res)
		assert.Equal(t, 42, sideEffect)
	})

	t.Run("failure: propagates error without invoking side effect", func(t *testing.T) {
		var called bool
		record := func(x int) ReaderResult[MyContext, string] {
			called = true
			return Of[MyContext]("logged")
		}

		rr := Left[MyContext, int](testError)
		res := MonadChainFirst(rr, record)(defaultContext)

		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called, "side-effect function should not be called on failure")
	})

	t.Run("success: discards return value of side-effect computation", func(t *testing.T) {
		sideEffect := func(x int) ReaderResult[MyContext, string] {
			return Of[MyContext]("ignored")
		}

		rr := Of[MyContext](99)
		res := MonadChainFirst(rr, sideEffect)(defaultContext)

		assert.Equal(t, result.Of(99), res)
	})

	t.Run("side-effect failure propagates as error", func(t *testing.T) {
		sideError := errors.New("side effect failed")
		failingSideEffect := func(x int) ReaderResult[MyContext, string] {
			return Left[MyContext, string](sideError)
		}

		rr := Of[MyContext](10)
		res := MonadChainFirst(rr, failingSideEffect)(defaultContext)

		assert.Equal(t, result.Left[int](sideError), res)
	})

	t.Run("side-effect receives the environment", func(t *testing.T) {
		var capturedCtx MyContext
		record := func(_ int) ReaderResult[MyContext, string] {
			return func(ctx MyContext) result.Result[string] {
				capturedCtx = ctx
				return result.Of("ok")
			}
		}

		rr := Of[MyContext](7)
		_ = MonadChainFirst(rr, record)(defaultContext)

		assert.Equal(t, defaultContext, capturedCtx)
	})
}

func TestChainFirst(t *testing.T) {
	t.Run("success: preserves original value and executes side effect", func(t *testing.T) {
		var sideEffect int
		record := func(x int) ReaderResult[MyContext, string] {
			sideEffect = x
			return Of[MyContext]("logged")
		}

		res := F.Pipe1(Of[MyContext](42), ChainFirst[MyContext, int, string](record))(defaultContext)

		assert.Equal(t, result.Of(42), res)
		assert.Equal(t, 42, sideEffect)
	})

	t.Run("failure: propagates error without invoking side effect", func(t *testing.T) {
		var called bool
		record := func(x int) ReaderResult[MyContext, string] {
			called = true
			return Of[MyContext]("logged")
		}

		res := F.Pipe1(Left[MyContext, int](testError), ChainFirst[MyContext, int, string](record))(defaultContext)

		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called, "side-effect function should not be called on failure")
	})

	t.Run("composes with other operators via Pipe", func(t *testing.T) {
		var log []int
		record := func(x int) ReaderResult[MyContext, string] {
			log = append(log, x)
			return Of[MyContext]("ok")
		}

		res := F.Pipe1(
			Of[MyContext](5),
			ChainFirst[MyContext, int, string](record),
		)(defaultContext)

		assert.Equal(t, result.Of(5), res)
		assert.Equal(t, []int{5}, log)
	})

	t.Run("side-effect failure propagates as error", func(t *testing.T) {
		sideError := errors.New("side effect failed")
		failingSideEffect := func(x int) ReaderResult[MyContext, string] {
			return Left[MyContext, string](sideError)
		}

		res := F.Pipe1(Of[MyContext](10), ChainFirst[MyContext, int, string](failingSideEffect))(defaultContext)

		assert.Equal(t, result.Left[int](sideError), res)
	})

	t.Run("chained multiple times keeps first value", func(t *testing.T) {
		var log []string
		appendLog := func(label string) func(int) ReaderResult[MyContext, string] {
			return func(x int) ReaderResult[MyContext, string] {
				log = append(log, fmt.Sprintf("%s:%d", label, x))
				return Of[MyContext]("ok")
			}
		}

		res := F.Pipe1(
			F.Pipe1(
				Of[MyContext](3),
				ChainFirst[MyContext, int, string](appendLog("first")),
			),
			ChainFirst[MyContext, int, string](appendLog("second")),
		)(defaultContext)

		assert.Equal(t, result.Of(3), res)
		assert.Equal(t, []string{"first:3", "second:3"}, log)
	})
}

func TestMonadChainReaderK(t *testing.T) {
	t.Run("sequences with a Reader-returning Kleisli and threads environment", func(t *testing.T) {
		// f reads from the environment using the value from ma
		f := func(x int) reader.Reader[MyContext, string] {
			return func(ctx MyContext) string { return fmt.Sprintf("%s:%d", ctx, x) }
		}
		rr := Of[MyContext](42)
		res := MonadChainReaderK(rr, f)(defaultContext)
		assert.Equal(t, result.Of("default:42"), res)
	})

	t.Run("error in ma is propagated; f is not called", func(t *testing.T) {
		var called bool
		f := func(x int) reader.Reader[MyContext, string] {
			called = true
			return func(ctx MyContext) string { return "ok" }
		}
		res := MonadChainReaderK(Left[MyContext, int](testError), f)(defaultContext)
		assert.Equal(t, result.Left[string](testError), res)
		assert.False(t, called)
	})
}

func TestChainReaderK(t *testing.T) {
	t.Run("curried form works in Pipe", func(t *testing.T) {
		f := func(x int) reader.Reader[MyContext, string] {
			return func(ctx MyContext) string { return fmt.Sprintf("%s:%d", ctx, x) }
		}
		res := F.Pipe1(Of[MyContext](7), ChainReaderK[MyContext, int, string](f))(defaultContext)
		assert.Equal(t, result.Of("default:7"), res)
	})

	t.Run("error propagates unchanged", func(t *testing.T) {
		f := func(x int) reader.Reader[MyContext, string] {
			return func(ctx MyContext) string { return "ok" }
		}
		res := F.Pipe1(Left[MyContext, int](testError), ChainReaderK[MyContext, int, string](f))(defaultContext)
		assert.Equal(t, result.Left[string](testError), res)
	})
}

func TestMonadApReader(t *testing.T) {
	t.Run("applies a Reader value to a wrapped function", func(t *testing.T) {
		add5 := func(x int) int { return x + 5 }
		fab := Of[MyContext](add5)
		fa := func(_ MyContext) int { return 10 }
		res := MonadApReader[int](fab, fa)(defaultContext)
		assert.Equal(t, result.Of(15), res)
	})

	t.Run("Reader reads from the environment", func(t *testing.T) {
		prefix := func(s string) string { return "Hello " + s }
		fab := Of[MyContext](prefix)
		fa := func(ctx MyContext) string { return string(ctx) }
		res := MonadApReader[string](fab, fa)(defaultContext)
		assert.Equal(t, result.Of("Hello default"), res)
	})

	t.Run("error in fab propagates; fa is not applied", func(t *testing.T) {
		fab := Left[MyContext, func(int) int](testError)
		fa := func(_ MyContext) int { return 99 }
		res := MonadApReader[int](fab, fa)(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
	})
}

func TestApReader(t *testing.T) {
	t.Run("curried form works in Pipe", func(t *testing.T) {
		double := func(x int) int { return x * 2 }
		fab := Of[MyContext](double)
		fa := func(_ MyContext) int { return 6 }
		res := F.Pipe1(fab, ApReader[int, MyContext, int](fa))(defaultContext)
		assert.Equal(t, result.Of(12), res)
	})

	t.Run("error in wrapped function propagates", func(t *testing.T) {
		fab := Left[MyContext, func(int) int](testError)
		fa := func(_ MyContext) int { return 6 }
		res := F.Pipe1(fab, ApReader[int, MyContext, int](fa))(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
	})
}

func TestMonadChainResultK(t *testing.T) {
	t.Run("chains with a Result-returning function on success", func(t *testing.T) {
		double := func(x int) result.Result[int] { return result.Of(x * 2) }
		res := MonadChainResultK(Of[MyContext](5), double)(defaultContext)
		assert.Equal(t, result.Of(10), res)
	})

	t.Run("failure in chained function propagates", func(t *testing.T) {
		chainErr := errors.New("chain error")
		fail := func(x int) result.Result[int] { return result.Left[int](chainErr) }
		res := MonadChainResultK(Of[MyContext](5), fail)(defaultContext)
		assert.Equal(t, result.Left[int](chainErr), res)
	})

	t.Run("error in ma propagates; f is not called", func(t *testing.T) {
		var called bool
		f := func(x int) result.Result[int] { called = true; return result.Of(x) }
		res := MonadChainResultK(Left[MyContext, int](testError), f)(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called)
	})
}

func TestChainResultK(t *testing.T) {
	t.Run("curried form works in Pipe", func(t *testing.T) {
		double := func(x int) result.Result[int] { return result.Of(x * 2) }
		res := F.Pipe1(Of[MyContext](4), ChainResultK[MyContext](double))(defaultContext)
		assert.Equal(t, result.Of(8), res)
	})

	t.Run("error propagates unchanged", func(t *testing.T) {
		double := func(x int) result.Result[int] { return result.Of(x * 2) }
		res := F.Pipe1(Left[MyContext, int](testError), ChainResultK[MyContext](double))(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
	})
}

func TestMonadAlt(t *testing.T) {
	t.Run("returns first when it succeeds", func(t *testing.T) {
		first := Of[MyContext](42)
		second := func() ReaderResult[MyContext, int] { return Of[MyContext](99) }
		res := MonadAlt(first, second)(defaultContext)
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("falls back to second when first fails", func(t *testing.T) {
		first := Left[MyContext, int](testError)
		second := func() ReaderResult[MyContext, int] { return Of[MyContext](99) }
		res := MonadAlt(first, second)(defaultContext)
		assert.Equal(t, result.Of(99), res)
	})

	t.Run("second is lazy: not evaluated when first succeeds", func(t *testing.T) {
		var evaluated bool
		first := Of[MyContext](1)
		second := func() ReaderResult[MyContext, int] {
			evaluated = true
			return Of[MyContext](99)
		}
		_ = MonadAlt(first, second)(defaultContext)
		assert.False(t, evaluated, "second should not be evaluated when first succeeds")
	})

	t.Run("second failure is propagated when first also fails", func(t *testing.T) {
		secondErr := errors.New("second failed")
		first := Left[MyContext, int](testError)
		second := func() ReaderResult[MyContext, int] { return Left[MyContext, int](secondErr) }
		res := MonadAlt(first, second)(defaultContext)
		assert.Equal(t, result.Left[int](secondErr), res)
	})
}

func TestAlt(t *testing.T) {
	t.Run("returns original value when it succeeds", func(t *testing.T) {
		second := func() ReaderResult[MyContext, int] { return Of[MyContext](99) }
		res := F.Pipe1(Of[MyContext](7), Alt(second))(defaultContext)
		assert.Equal(t, result.Of(7), res)
	})

	t.Run("falls back to second when original fails", func(t *testing.T) {
		second := func() ReaderResult[MyContext, int] { return Of[MyContext](99) }
		res := F.Pipe1(Left[MyContext, int](testError), Alt(second))(defaultContext)
		assert.Equal(t, result.Of(99), res)
	})

	t.Run("second is lazy: not evaluated on success", func(t *testing.T) {
		var evaluated bool
		second := func() ReaderResult[MyContext, int] {
			evaluated = true
			return Of[MyContext](99)
		}
		_ = F.Pipe1(Of[MyContext](1), Alt(second))(defaultContext)
		assert.False(t, evaluated)
	})
}

func TestMonadChainFirstI(t *testing.T) {
	t.Run("success: preserves original value and executes idiomatic side effect", func(t *testing.T) {
		var capturedVal int
		f := func(x int) func(MyContext) (string, error) {
			return func(ctx MyContext) (string, error) {
				capturedVal = x
				return "logged", nil
			}
		}
		res := MonadChainFirstI(Of[MyContext](42), f)(defaultContext)
		assert.Equal(t, result.Of(42), res)
		assert.Equal(t, 42, capturedVal)
	})

	t.Run("failure: propagates error without calling f", func(t *testing.T) {
		var called bool
		f := func(x int) func(MyContext) (string, error) {
			return func(ctx MyContext) (string, error) {
				called = true
				return "", nil
			}
		}
		res := MonadChainFirstI(Left[MyContext, int](testError), f)(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called)
	})

	t.Run("side-effect error propagates", func(t *testing.T) {
		sideError := errors.New("side failed")
		f := func(x int) func(MyContext) (string, error) {
			return func(ctx MyContext) (string, error) { return "", sideError }
		}
		res := MonadChainFirstI(Of[MyContext](10), f)(defaultContext)
		assert.Equal(t, result.Left[int](sideError), res)
	})
}

func TestChainFirstI(t *testing.T) {
	t.Run("curried form preserves original value", func(t *testing.T) {
		var capturedVal int
		f := func(x int) func(MyContext) (string, error) {
			return func(ctx MyContext) (string, error) {
				capturedVal = x
				return "logged", nil
			}
		}
		res := F.Pipe1(Of[MyContext](5), ChainFirstI[MyContext, int, string](f))(defaultContext)
		assert.Equal(t, result.Of(5), res)
		assert.Equal(t, 5, capturedVal)
	})

	t.Run("error propagates without calling f", func(t *testing.T) {
		var called bool
		f := func(x int) func(MyContext) (string, error) {
			return func(ctx MyContext) (string, error) { called = true; return "", nil }
		}
		res := F.Pipe1(Left[MyContext, int](testError), ChainFirstI[MyContext, int, string](f))(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called)
	})
}

func TestMonadChainFirstResultK(t *testing.T) {
	t.Run("success: preserves original value after Result side effect", func(t *testing.T) {
		validate := func(x int) result.Result[string] {
			if x > 0 {
				return result.Of("valid")
			}
			return result.Left[string](errors.New("not positive"))
		}
		res := MonadChainFirstResultK(Of[MyContext](42), validate)(defaultContext)
		assert.Equal(t, result.Of(42), res)
	})

	t.Run("side-effect failure propagates", func(t *testing.T) {
		sideError := errors.New("validation failed")
		fail := func(x int) result.Result[string] { return result.Left[string](sideError) }
		res := MonadChainFirstResultK(Of[MyContext](42), fail)(defaultContext)
		assert.Equal(t, result.Left[int](sideError), res)
	})

	t.Run("error in ma propagates without calling f", func(t *testing.T) {
		var called bool
		f := func(x int) result.Result[string] { called = true; return result.Of("ok") }
		res := MonadChainFirstResultK(Left[MyContext, int](testError), f)(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called)
	})
}

func TestChainFirstResultK(t *testing.T) {
	t.Run("curried form preserves original value", func(t *testing.T) {
		validate := func(x int) result.Result[string] { return result.Of("ok") }
		res := F.Pipe1(Of[MyContext](3), ChainFirstResultK[MyContext](validate))(defaultContext)
		assert.Equal(t, result.Of(3), res)
	})

	t.Run("side-effect failure propagates in pipeline", func(t *testing.T) {
		sideError := errors.New("fail")
		fail := func(x int) result.Result[string] { return result.Left[string](sideError) }
		res := F.Pipe1(Of[MyContext](3), ChainFirstResultK[MyContext](fail))(defaultContext)
		assert.Equal(t, result.Left[int](sideError), res)
	})
}

func TestMonadChainFirstResultIK(t *testing.T) {
	t.Run("success: preserves original value after idiomatic side effect", func(t *testing.T) {
		log := func(x int) (string, error) { return fmt.Sprintf("logged %d", x), nil }
		res := MonadChainFirstResultIK(Of[MyContext](7), log)(defaultContext)
		assert.Equal(t, result.Of(7), res)
	})

	t.Run("idiomatic side-effect error propagates", func(t *testing.T) {
		sideError := errors.New("log failed")
		fail := func(x int) (string, error) { return "", sideError }
		res := MonadChainFirstResultIK(Of[MyContext](7), fail)(defaultContext)
		assert.Equal(t, result.Left[int](sideError), res)
	})

	t.Run("error in ma propagates without calling f", func(t *testing.T) {
		var called bool
		f := func(x int) (string, error) { called = true; return "", nil }
		res := MonadChainFirstResultIK(Left[MyContext, int](testError), f)(defaultContext)
		assert.Equal(t, result.Left[int](testError), res)
		assert.False(t, called)
	})
}

func TestChainFirstResultIK(t *testing.T) {
	t.Run("curried form preserves original value", func(t *testing.T) {
		log := func(x int) (string, error) { return "ok", nil }
		res := F.Pipe1(Of[MyContext](9), ChainFirstResultIK[MyContext](log))(defaultContext)
		assert.Equal(t, result.Of(9), res)
	})

	t.Run("idiomatic side-effect error propagates in pipeline", func(t *testing.T) {
		sideError := errors.New("fail")
		fail := func(x int) (string, error) { return "", sideError }
		res := F.Pipe1(Of[MyContext](9), ChainFirstResultIK[MyContext](fail))(defaultContext)
		assert.Equal(t, result.Left[int](sideError), res)
	})
}
