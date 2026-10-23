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

package readerreaderioresult

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	RO "github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/readerresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/IBM/fp-go/v2/tailrec"
)

// rxCfg is the outer environment used by most examples in this file set.
var rxCfg = exConfig{Multiplier: 10}

// rxErr is the error used by the failure examples.
var rxErr = errors.New("boom")

// rxMultiplier reads the multiplier from the outer environment.
func rxMultiplier(c exConfig) int { return c.Multiplier }

// rxScale multiplies n with the multiplier of the outer environment.
func rxScale(n int) ReaderReaderIOResult[exConfig, int] {
	return Asks(F.Flow2(rxMultiplier, N.Mul(n)))
}

// rxPositive fails for numbers that are not positive.
func rxPositive(n int) Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

// rxPrintInt prints an int as a side effect.
var rxPrintInt = io.Printf[int]("log: %d\n")

// rxPrintErr prints an error as a side effect.
var rxPrintErr = io.Printf[error]("error: %v\n")

// ExampleOf lifts a pure value.
func ExampleOf() {
	fmt.Println(exRun(rxCfg, Of[exConfig](42)))
	// Output: 42 <nil>
}

// ExampleRight creates a successful computation.
func ExampleRight() {
	fmt.Println(exRun(rxCfg, Right[exConfig](42)))
	// Output: 42 <nil>
}

// ExampleLeft creates a failed computation.
func ExampleLeft() {
	fmt.Println(exRun(rxCfg, Left[exConfig, int](rxErr)))
	// Output: 0 boom
}

// ExampleFromEither lifts an Either with an error on the left.
func ExampleFromEither() {
	fmt.Println(exRun(rxCfg, FromEither[exConfig](E.Right[error](3))))
	// Output: 3 <nil>
}

// ExampleFromResult lifts a Result.
func ExampleFromResult() {
	fmt.Println(exRun(rxCfg, FromResult[exConfig](rxPositive(-3))))
	// Output: 0 not positive
}

// ExampleFromIO lifts a side effect that cannot fail.
func ExampleFromIO() {
	fmt.Println(exRun(rxCfg, FromIO[exConfig](rxPrintInt(1))))
	// Output:
	// log: 1
	// 1 <nil>
}

// ExampleRightIO is an alias for FromIO.
func ExampleRightIO() {
	fmt.Println(exRun(rxCfg, RightIO[exConfig](io.Of(1))))
	// Output: 1 <nil>
}

// ExampleLeftIO lifts a side effect that produces an error.
func ExampleLeftIO() {
	fmt.Println(exRun(rxCfg, LeftIO[exConfig, int](io.Of(rxErr))))
	// Output: 0 boom
}

// ExampleFromIOEither lifts an IOEither with an error on the left.
func ExampleFromIOEither() {
	fmt.Println(exRun(rxCfg, FromIOEither[exConfig](IOE.Of[error](7))))
	// Output: 7 <nil>
}

// ExampleFromIOResult lifts an IOResult.
func ExampleFromIOResult() {
	fmt.Println(exRun(rxCfg, FromIOResult[exConfig](ioresult.Left[int](rxErr))))
	// Output: 0 boom
}

// ExampleFromReader lifts a Reader on the outer environment.
func ExampleFromReader() {
	fmt.Println(exRun(rxCfg, FromReader(rxMultiplier)))
	// Output: 10 <nil>
}

// ExampleRightReader is an alias for FromReader.
func ExampleRightReader() {
	fmt.Println(exRun(rxCfg, RightReader(rxMultiplier)))
	// Output: 10 <nil>
}

// ExampleLeftReader derives an error from the outer environment.
func ExampleLeftReader() {
	fmt.Println(exRun(rxCfg, LeftReader[int](func(c exConfig) error {
		return fmt.Errorf("multiplier %d not supported", c.Multiplier)
	})))
	// Output: 0 multiplier 10 not supported
}

// ExampleFromReaderIO lifts a ReaderIO on the outer environment.
func ExampleFromReaderIO() {
	fmt.Println(exRun(rxCfg, FromReaderIO(RIO.Asks(rxMultiplier))))
	// Output: 10 <nil>
}

// ExampleRightReaderIO is an alias for FromReaderIO.
func ExampleRightReaderIO() {
	fmt.Println(exRun(rxCfg, RightReaderIO(RIO.Asks(rxMultiplier))))
	// Output: 10 <nil>
}

// ExampleLeftReaderIO derives an error from the outer environment with a side effect.
func ExampleLeftReaderIO() {
	fmt.Println(exRun(rxCfg, LeftReaderIO[int](RIO.Of[exConfig](rxErr))))
	// Output: 0 boom
}

// ExampleFromReaderEither lifts a ReaderEither on the outer environment.
func ExampleFromReaderEither() {
	fmt.Println(exRun(rxCfg, FromReaderEither(RE.Of[exConfig, error](5))))
	// Output: 5 <nil>
}

// ExampleFromReaderResult lifts a ReaderResult on the outer environment.
func ExampleFromReaderResult() {
	fmt.Println(exRun(rxCfg, FromReaderResult(readerresult.Asks(rxMultiplier))))
	// Output: 10 <nil>
}

// ExampleFromReaderIOResult lifts a ReaderIOResult on the outer environment.
func ExampleFromReaderIOResult() {
	fmt.Println(exRun(rxCfg, FromReaderIOResult(RIOR.Asks(rxMultiplier))))
	// Output: 10 <nil>
}

// ExampleFromReaderOption lifts a ReaderOption and maps None to an error.
func ExampleFromReaderOption() {
	fromRO := FromReaderOption[exConfig, int](lazy.Of(errors.New("missing")))

	fmt.Println(exRun(rxCfg, fromRO(RO.Of[exConfig](5))))
	fmt.Println(exRun(rxCfg, fromRO(RO.None[exConfig, int]())))
	// Output:
	// 5 <nil>
	// 0 missing
}

// ExampleFromOption lifts an Option and maps None to an error.
func ExampleFromOption() {
	fromOption := FromOption[exConfig, int](lazy.Of(errors.New("missing")))

	fmt.Println(exRun(rxCfg, fromOption(O.Some(5))))
	fmt.Println(exRun(rxCfg, fromOption(O.None[int]())))
	// Output:
	// 5 <nil>
	// 0 missing
}

// ExampleFromPredicate succeeds if the predicate holds and fails otherwise.
func ExampleFromPredicate() {
	positive := FromPredicate[exConfig](N.MoreThan(0), func(n int) error {
		return fmt.Errorf("%d is not positive", n)
	})

	fmt.Println(exRun(rxCfg, positive(5)))
	fmt.Println(exRun(rxCfg, positive(-5)))
	// Output:
	// 5 <nil>
	// 0 -5 is not positive
}

// ExampleAsk returns the outer environment.
func ExampleAsk() {
	fmt.Println(Ask[exConfig]()(rxCfg)(context.Background())())
	// Output: Right[readerreaderioresult.exConfig]({10 map[]})
}

// ExampleAsks projects a value from the outer environment.
func ExampleAsks() {
	fmt.Println(exRun(rxCfg, Asks(rxMultiplier)))
	// Output: 10 <nil>
}

// ExampleRead provides the outer environment and returns a ReaderIOResult on the context.
func ExampleRead() {
	rior := F.Pipe1(rxScale(4), Read[int](exConfig{Multiplier: 3}))

	fmt.Println(rior(context.Background())())
	// Output: Right[int](12)
}

// ExampleReadIO provides the outer environment from an IO.
func ExampleReadIO() {
	rior := F.Pipe1(rxScale(4), ReadIO[int](io.Of(exConfig{Multiplier: 3})))

	fmt.Println(rior(context.Background())())
	// Output: Right[int](12)
}

// ExampleReadIOEither provides the outer environment from an IOEither that may fail.
func ExampleReadIOEither() {
	rior := F.Pipe1(rxScale(4), ReadIOEither[int](IOE.Left[exConfig](rxErr)))

	fmt.Println(rior(context.Background())())
	// Output: Left[*errors.errorString](boom)
}

// ExampleFlatten removes one level of nesting.
func ExampleFlatten() {
	fmt.Println(exRun(rxCfg, Flatten(Of[exConfig](rxScale(4)))))
	// Output: 40 <nil>
}

// ExampleDelay waits before running the computation.
func ExampleDelay() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxScale(4), Delay[exConfig, int](time.Millisecond))))
	// Output: 40 <nil>
}

// ExampleDefer creates the computation anew each time it runs.
func ExampleDefer() {
	calls := 0
	deferred := Defer(func() ReaderReaderIOResult[exConfig, int] {
		calls++
		return Of[exConfig](calls)
	})

	fmt.Println(exRun(rxCfg, deferred))
	fmt.Println(exRun(rxCfg, deferred))
	// Output:
	// 1 <nil>
	// 2 <nil>
}

// ExampleTailRec runs a stack safe loop; each step reads the decrement from the
// outer environment and the context is checked for cancellation before each step.
func ExampleTailRec() {
	countdown := TailRec(func(n int) ReaderReaderIOResult[exConfig, Trampoline[int, string]] {
		if n <= 0 {
			return Of[exConfig](tailrec.Land[int]("liftoff"))
		}
		return Asks(func(c exConfig) Trampoline[int, string] {
			return tailrec.Bounce[string](n - c.Multiplier)
		})
	})

	fmt.Println(exRun(rxCfg, countdown(1_000_000)))
	// Output: liftoff <nil>
}

// ExampleTailRec_cancelled shows that a cancelled context stops the recursion.
func ExampleTailRec_cancelled() {
	endless := TailRec(func(n int) ReaderReaderIOResult[exConfig, Trampoline[int, int]] {
		return Of[exConfig](tailrec.Bounce[int](n + 1))
	})

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	fmt.Println(endless(0)(rxCfg)(ctx)())
	// Output: Left[*errors.errorString](context canceled)
}

// ExampleRetrying retries an action as long as it fails, up to the limit of the policy.
func ExampleRetrying() {
	action := func(status retry.RetryStatus) ReaderReaderIOResult[exConfig, int] {
		return FromResult[exConfig](rxPositive(int(status.IterNumber) - 1))
	}

	fmt.Println(exRun(rxCfg, Retrying(retry.LimitRetries(5), action, result.IsLeft[int])))
	// Output: 1 <nil>
}

// ExampleRetryingI retries an idiomatic action.
func ExampleRetryingI() {
	action := func(status retry.RetryStatus) func(context.Context, exConfig) (int, error) {
		return func(_ context.Context, c exConfig) (int, error) {
			if status.IterNumber < 2 {
				return 0, rxErr
			}
			return int(status.IterNumber) * c.Multiplier, nil
		}
	}

	fmt.Println(exRun(rxCfg, RetryingI(retry.LimitRetries(5), action, result.IsLeft[int])))
	// Output: 20 <nil>
}

// ExampleEitherize lifts a function in idiomatic Go style.
func ExampleEitherize() {
	load := Eitherize(func(c exConfig, _ context.Context) (string, error) {
		return strconv.Itoa(c.Multiplier), nil
	})

	fmt.Println(exRun(rxCfg, load))
	// Output: 10 <nil>
}

// ExampleEitherize1 lifts a function with one argument in idiomatic Go style.
func ExampleEitherize1() {
	scale := Eitherize1(func(c exConfig, _ context.Context, s string) (int, error) {
		n, err := strconv.Atoi(s)
		return n * c.Multiplier, err
	})

	fmt.Println(exRun(rxCfg, scale("4")))
	fmt.Println(exRun(rxCfg, scale("x")))
	// Output:
	// 40 <nil>
	// 0 strconv.Atoi: parsing "x": invalid syntax
}

// ExampleBracket acquires a resource, uses it and releases it, even if using it fails.
func ExampleBracket() {
	acquire := FromIO[exConfig](io.FromImpure(func() { fmt.Println("open") }))
	release := func(_ F.Void, _ Result[int]) ReaderReaderIOResult[exConfig, F.Void] {
		return FromIO[exConfig](io.FromImpure(func() { fmt.Println("close") }))
	}

	fmt.Println(exRun(rxCfg, Bracket(acquire, F.Constant1[F.Void](rxScale(4)), release)))
	fmt.Println(exRun(rxCfg, Bracket(acquire, F.Constant1[F.Void](Left[exConfig, int](rxErr)), release)))
	// Output:
	// open
	// close
	// 40 <nil>
	// open
	// close
	// 0 boom
}

// ExamplePaired turns the computation into a function of the pair of context and environment.
func ExamplePaired() {
	run := Paired(rxScale(3))

	fmt.Println(run(pair.MakePair(context.Background(), exConfig{Multiplier: 5}))())
	// Output: Right[int](15)
}

// ExampleMemoize caches the result per outer environment.
func ExampleMemoize() {
	type key struct{ ID int }
	calls := 0
	cached := Memoize(Asks(func(k key) int {
		calls++
		return k.ID * 10
	}))

	fmt.Println(cached(key{1})(context.Background())())
	fmt.Println(cached(key{1})(context.Background())())
	fmt.Println(cached(key{2})(context.Background())())
	fmt.Println("calls:", calls)
	// Output:
	// Right[int](10)
	// Right[int](10)
	// Right[int](20)
	// calls: 2
}

// ExampleContramapMemoize caches the result by a key derived from the outer environment.
func ExampleContramapMemoize() {
	calls := 0
	cached := F.Pipe1(
		FromIO[exConfig](func() int {
			calls++
			return calls
		}),
		ContramapMemoize[int](rxMultiplier),
	)

	first, _ := exRun(exConfig{Multiplier: 1}, cached)
	second, _ := exRun(exConfig{Multiplier: 1, Table: map[string]string{"a": "b"}}, cached)
	other, _ := exRun(exConfig{Multiplier: 2}, cached)

	fmt.Println(first, second, other)
	// Output: 1 1 2
}

// ExampleFilter filters any container given its filter function.
func ExampleFilter() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		Of[exConfig]([]int{1, 2, 3, 4, 5}),
		Filter[exConfig](A.Filter[int])(N.MoreThan(2)),
	)))
	// Output: [3 4 5] <nil>
}

// ExampleFilterArray keeps the elements of a slice that satisfy the predicate.
func ExampleFilterArray() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		Of[exConfig]([]int{1, 2, 3, 4, 5}),
		FilterArray[exConfig](N.MoreThan(2)),
	)))
	// Output: [3 4 5] <nil>
}

// ExampleFilterIter keeps the elements of a sequence that satisfy the predicate.
func ExampleFilterIter() {
	seq, err := exRun(rxCfg, F.Pipe1(
		Of[exConfig](iter.From(1, 2, 3, 4, 5)),
		FilterIter[exConfig](N.MoreThan(2)),
	))

	fmt.Println(slices.Collect(seq), err)
	// Output: [3 4 5] <nil>
}

// ExampleFilterMap filters and maps any container given its filter-map function.
func ExampleFilterMap() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		Of[exConfig]([]int{-1, 2, -3, 4}),
		FilterMap[exConfig](A.FilterMap[int, string])(F.Flow2(
			O.FromPredicate(N.MoreThan(0)),
			O.Map(S.Format[int]("n=%d")),
		)),
	)))
	// Output: [n=2 n=4] <nil>
}

// ExampleFilterMapArray filters and maps the elements of a slice.
func ExampleFilterMapArray() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		Of[exConfig]([]int{-1, 2, -3, 4}),
		FilterMapArray[exConfig](F.Flow2(
			O.FromPredicate(N.MoreThan(0)),
			O.Map(S.Format[int]("n=%d")),
		)),
	)))
	// Output: [n=2 n=4] <nil>
}

// ExampleFilterMapIter filters and maps the elements of a sequence.
func ExampleFilterMapIter() {
	seq, err := exRun(rxCfg, F.Pipe1(
		Of[exConfig](iter.From(-1, 2, -3, 4)),
		FilterMapIter[exConfig](F.Flow2(
			O.FromPredicate(N.MoreThan(0)),
			O.Map(S.Format[int]("n=%d")),
		)),
	))

	fmt.Println(slices.Collect(seq), err)
	// Output: [n=2 n=4] <nil>
}
