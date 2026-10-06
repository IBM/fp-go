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
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	RD "github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RO "github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// rxScaleR multiplies n with the multiplier as a Reader on the outer environment.
func rxScaleR(n int) RD.Reader[exConfig, int] {
	return F.Flow2(rxMultiplier, N.Mul(n))
}

// rxScaleRIO multiplies n with the multiplier as a ReaderIO on the outer environment.
func rxScaleRIO(n int) RIO.ReaderIO[exConfig, int] {
	return RIO.Asks(rxScaleR(n))
}

// rxAboveMultiplier fails if n does not exceed the multiplier of the outer environment.
func rxAboveMultiplier(n int) RE.ReaderEither[exConfig, error, int] {
	return func(c exConfig) Result[int] {
		return rxPositive(n - c.Multiplier)
	}
}

// rxAboveMultiplierRO is None if n does not exceed the multiplier of the outer environment.
func rxAboveMultiplierRO(n int) RO.ReaderOption[exConfig, int] {
	return func(c exConfig) Option[int] {
		return O.FromPredicate(N.MoreThan(c.Multiplier))(n)
	}
}

// rxScaleI multiplies n with the multiplier in idiomatic Go style.
func rxScaleI(n int) func(context.Context, exConfig) (int, error) {
	return func(_ context.Context, c exConfig) (int, error) {
		return n * c.Multiplier, nil
	}
}

// rxPositiveI fails for numbers that are not positive, in idiomatic Go style.
func rxPositiveI(n int) func(context.Context, exConfig) (int, error) {
	return func(context.Context, exConfig) (int, error) {
		return result.Unwrap(rxPositive(n))
	}
}

// rxLogErr logs an error and succeeds.
func rxLogErr(e error) ReaderReaderIOResult[exConfig, error] {
	return FromIO[exConfig](rxPrintErr(e))
}

// ExampleMap transforms the success value.
func ExampleMap() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](21), Map[exConfig](N.Mul(2)))))
	// Output: 42 <nil>
}

// ExampleMonadMap is the non curried version of Map.
func ExampleMonadMap() {
	fmt.Println(exRun(rxCfg, MonadMap(Of[exConfig](21), N.Mul(2))))
	// Output: 42 <nil>
}

// ExampleMapTo replaces the success value.
func ExampleMapTo() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](21), MapTo[exConfig, int]("done"))))
	// Output: done <nil>
}

// ExampleMonadMapTo is the non curried version of MapTo.
func ExampleMonadMapTo() {
	fmt.Println(exRun(rxCfg, MonadMapTo(Of[exConfig](21), "done")))
	// Output: done <nil>
}

// ExampleChain sequences a computation that depends on the previous result.
func ExampleChain() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), Chain(rxScale))))
	// Output: 40 <nil>
}

// ExampleMonadChain is the non curried version of Chain.
func ExampleMonadChain() {
	fmt.Println(exRun(rxCfg, MonadChain(Of[exConfig](4), rxScale)))
	// Output: 40 <nil>
}

// ExampleChainFirst runs a computation for its effect and keeps the original value.
func ExampleChainFirst() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainFirst(rxScale))))
	// Output: 4 <nil>
}

// ExampleMonadChainFirst is the non curried version of ChainFirst.
func ExampleMonadChainFirst() {
	fmt.Println(exRun(rxCfg, MonadChainFirst(Of[exConfig](4), rxScale)))
	// Output: 4 <nil>
}

// ExampleTap is an alias for ChainFirst.
func ExampleTap() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), Tap(rxScale))))
	// Output: 4 <nil>
}

// ExampleMonadTap is the non curried version of Tap.
func ExampleMonadTap() {
	fmt.Println(exRun(rxCfg, MonadTap(Left[exConfig, int](rxErr), rxScale)))
	// Output: 0 boom
}

// ExampleChainEitherK chains a function that returns an Either.
func ExampleChainEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](-3), ChainEitherK[exConfig](rxPositive))))
	// Output: 0 not positive
}

// ExampleMonadChainEitherK is the non curried version of ChainEitherK.
func ExampleMonadChainEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainEitherK(Of[exConfig](3), rxPositive)))
	// Output: 3 <nil>
}

// ExampleChainResultK chains a function that returns a Result.
func ExampleChainResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](3), ChainResultK[exConfig](rxPositive))))
	// Output: 3 <nil>
}

// ExampleChainFirstEitherK validates the value with an Either and keeps it.
func ExampleChainFirstEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](3), ChainFirstEitherK[exConfig](rxPositive))))
	// Output: 3 <nil>
}

// ExampleMonadChainFirstEitherK is the non curried version of ChainFirstEitherK.
func ExampleMonadChainFirstEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstEitherK(Of[exConfig](0), rxPositive)))
	// Output: 0 not positive
}

// ExampleTapEitherK is an alias for ChainFirstEitherK.
func ExampleTapEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](0), TapEitherK[exConfig](rxPositive))))
	// Output: 0 not positive
}

// ExampleMonadTapEitherK is the non curried version of TapEitherK.
func ExampleMonadTapEitherK() {
	fmt.Println(exRun(rxCfg, MonadTapEitherK(Of[exConfig](3), rxPositive)))
	// Output: 3 <nil>
}

// ExampleChainReaderK chains a Reader on the outer environment.
func ExampleChainReaderK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainReaderK(rxScaleR))))
	// Output: 40 <nil>
}

// ExampleMonadChainReaderK is the non curried version of ChainReaderK.
func ExampleMonadChainReaderK() {
	fmt.Println(exRun(rxCfg, MonadChainReaderK(Of[exConfig](4), rxScaleR)))
	// Output: 40 <nil>
}

// ExampleChainFirstReaderK runs a Reader and keeps the original value.
func ExampleChainFirstReaderK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainFirstReaderK(rxScaleR))))
	// Output: 4 <nil>
}

// ExampleMonadChainFirstReaderK is the non curried version of ChainFirstReaderK.
func ExampleMonadChainFirstReaderK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstReaderK(Of[exConfig](4), rxScaleR)))
	// Output: 4 <nil>
}

// ExampleTapReaderK is an alias for ChainFirstReaderK.
func ExampleTapReaderK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), TapReaderK(rxScaleR))))
	// Output: 4 <nil>
}

// ExampleMonadTapReaderK is the non curried version of TapReaderK.
func ExampleMonadTapReaderK() {
	fmt.Println(exRun(rxCfg, MonadTapReaderK(Of[exConfig](4), rxScaleR)))
	// Output: 4 <nil>
}

// ExampleChainReaderIOK chains a ReaderIO on the outer environment.
func ExampleChainReaderIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainReaderIOK(rxScaleRIO))))
	// Output: 40 <nil>
}

// ExampleMonadChainReaderIOK is the non curried version of ChainReaderIOK.
func ExampleMonadChainReaderIOK() {
	fmt.Println(exRun(rxCfg, MonadChainReaderIOK(Of[exConfig](4), rxScaleRIO)))
	// Output: 40 <nil>
}

// ExampleChainFirstReaderIOK runs a ReaderIO and keeps the original value.
func ExampleChainFirstReaderIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainFirstReaderIOK(rxScaleRIO))))
	// Output: 4 <nil>
}

// ExampleMonadChainFirstReaderIOK is the non curried version of ChainFirstReaderIOK.
func ExampleMonadChainFirstReaderIOK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstReaderIOK(Of[exConfig](4), rxScaleRIO)))
	// Output: 4 <nil>
}

// ExampleTapReaderIOK is an alias for ChainFirstReaderIOK.
func ExampleTapReaderIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), TapReaderIOK(rxScaleRIO))))
	// Output: 4 <nil>
}

// ExampleMonadTapReaderIOK is the non curried version of TapReaderIOK.
func ExampleMonadTapReaderIOK() {
	fmt.Println(exRun(rxCfg, MonadTapReaderIOK(Of[exConfig](4), rxScaleRIO)))
	// Output: 4 <nil>
}

// ExampleChainReaderEitherK chains a ReaderEither on the outer environment.
func ExampleChainReaderEitherK() {
	aboveMultiplier := ChainReaderEitherK(rxAboveMultiplier)

	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](15))))
	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](5))))
	// Output:
	// 5 <nil>
	// 0 not positive
}

// ExampleMonadChainReaderEitherK is the non curried version of ChainReaderEitherK.
func ExampleMonadChainReaderEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainReaderEitherK(Of[exConfig](15), rxAboveMultiplier)))
	// Output: 5 <nil>
}

// ExampleChainFirstReaderEitherK validates with a ReaderEither and keeps the value.
func ExampleChainFirstReaderEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](15), ChainFirstReaderEitherK(rxAboveMultiplier))))
	// Output: 15 <nil>
}

// ExampleMonadChainFirstReaderEitherK is the non curried version of ChainFirstReaderEitherK.
func ExampleMonadChainFirstReaderEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstReaderEitherK(Of[exConfig](5), rxAboveMultiplier)))
	// Output: 0 not positive
}

// ExampleTapReaderEitherK is an alias for ChainFirstReaderEitherK.
func ExampleTapReaderEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](15), TapReaderEitherK(rxAboveMultiplier))))
	// Output: 15 <nil>
}

// ExampleMonadTapReaderEitherK is the non curried version of TapReaderEitherK.
func ExampleMonadTapReaderEitherK() {
	fmt.Println(exRun(rxCfg, MonadTapReaderEitherK(Of[exConfig](15), rxAboveMultiplier)))
	// Output: 15 <nil>
}

// ExampleChainReaderOptionK chains a ReaderOption and maps None to an error.
func ExampleChainReaderOptionK() {
	aboveMultiplier := ChainReaderOptionK[exConfig, int, int](lazy.Of(errors.New("too small")))(rxAboveMultiplierRO)

	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](15))))
	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](5))))
	// Output:
	// 15 <nil>
	// 0 too small
}

// ExampleChainFirstReaderOptionK validates with a ReaderOption and keeps the value.
func ExampleChainFirstReaderOptionK() {
	aboveMultiplier := ChainFirstReaderOptionK[exConfig, int, int](lazy.Of(errors.New("too small")))(rxAboveMultiplierRO)

	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](15))))
	// Output: 15 <nil>
}

// ExampleTapReaderOptionK is an alias for ChainFirstReaderOptionK.
func ExampleTapReaderOptionK() {
	aboveMultiplier := TapReaderOptionK[exConfig, int, int](lazy.Of(errors.New("too small")))(rxAboveMultiplierRO)

	fmt.Println(exRun(rxCfg, aboveMultiplier(Of[exConfig](5))))
	// Output: 0 too small
}

// ExampleChainIOEitherK chains a function that returns an IOEither.
func ExampleChainIOEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		Of[exConfig](-1),
		ChainIOEitherK[exConfig](F.Flow2(rxPositive, IOE.FromEither[error, int])),
	)))
	// Output: 0 not positive
}

// ExampleMonadChainIOEitherK is the non curried version of ChainIOEitherK.
func ExampleMonadChainIOEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainIOEitherK(Of[exConfig](1), F.Flow2(rxPositive, IOE.FromEither[error, int]))))
	// Output: 1 <nil>
}

// ExampleChainIOK chains a side effect that cannot fail.
func ExampleChainIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainIOK[exConfig](F.Flow2(strconv.Itoa, io.Of[string])))))
	// Output: 7 <nil>
}

// ExampleMonadChainIOK is the non curried version of ChainIOK.
func ExampleMonadChainIOK() {
	fmt.Println(exRun(rxCfg, MonadChainIOK(Of[exConfig](7), rxPrintInt)))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleChainFirstIOK runs a side effect and keeps the original value.
func ExampleChainFirstIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainFirstIOK[exConfig](rxPrintInt))))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleMonadChainFirstIOK is the non curried version of ChainFirstIOK.
func ExampleMonadChainFirstIOK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstIOK(Of[exConfig](7), rxPrintInt)))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleTapIOK is an alias for ChainFirstIOK; the side effect is skipped on failure.
func ExampleTapIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), TapIOK[exConfig](rxPrintInt))))
	// Output: 0 boom
}

// ExampleMonadTapIOK is the non curried version of TapIOK.
func ExampleMonadTapIOK() {
	fmt.Println(exRun(rxCfg, MonadTapIOK(Of[exConfig](7), rxPrintInt)))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleChainOptionK chains a function that returns an Option and maps None to an error.
func ExampleChainOptionK() {
	positive := ChainOptionK[exConfig, int, int](lazy.Of(errors.New("not positive")))(O.FromPredicate(N.MoreThan(0)))

	fmt.Println(exRun(rxCfg, positive(Of[exConfig](3))))
	fmt.Println(exRun(rxCfg, positive(Of[exConfig](-3))))
	// Output:
	// 3 <nil>
	// 0 not positive
}

// ExampleAp applies a wrapped function to a wrapped value.
func ExampleAp() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](N.Mul(2)), Ap[int](Of[exConfig](21)))))
	// Output: 42 <nil>
}

// ExampleMonadAp is the non curried version of Ap.
func ExampleMonadAp() {
	fmt.Println(exRun(rxCfg, MonadAp(Of[exConfig](N.Mul(2)), rxScale(2))))
	// Output: 40 <nil>
}

// ExampleMonadApSeq evaluates the function before the value.
func ExampleMonadApSeq() {
	fmt.Println(exRun(rxCfg, MonadApSeq(Of[exConfig](S.Format[int]("#%d")), rxScale(2))))
	// Output: #20 <nil>
}

// ExampleMonadApPar evaluates the function and the value concurrently.
func ExampleMonadApPar() {
	fmt.Println(exRun(rxCfg, MonadApPar(Of[exConfig](N.Mul(2)), Left[exConfig, int](rxErr))))
	// Output: 0 boom
}

// ExampleFlap applies a wrapped function to a plain value.
func ExampleFlap() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](N.Mul(2)), Flap[exConfig, int](21))))
	// Output: 42 <nil>
}

// ExampleMonadFlap is the non curried version of Flap.
func ExampleMonadFlap() {
	fmt.Println(exRun(rxCfg, MonadFlap(Of[exConfig](N.Mul(2)), 21)))
	// Output: 42 <nil>
}

// ExampleAlt falls back to another computation on failure.
func ExampleAlt() {
	fallback := Alt(lazy.Of(Of[exConfig](0)))

	fmt.Println(exRun(rxCfg, fallback(Of[exConfig](42))))
	fmt.Println(exRun(rxCfg, fallback(Left[exConfig, int](rxErr))))
	// Output:
	// 42 <nil>
	// 0 <nil>
}

// ExampleMonadAlt is the non curried version of Alt.
func ExampleMonadAlt() {
	fmt.Println(exRun(rxCfg, MonadAlt(Left[exConfig, int](rxErr), lazy.Of(rxScale(1)))))
	// Output: 10 <nil>
}

// ExampleMapLeft transforms the error.
func ExampleMapLeft() {
	wrap := func(e error) error { return fmt.Errorf("wrapped: %w", e) }

	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), MapLeft[exConfig, int](wrap))))
	// Output: 0 wrapped: boom
}

// ExampleMonadMapLeft is the non curried version of MapLeft.
func ExampleMonadMapLeft() {
	wrap := func(e error) error { return fmt.Errorf("wrapped: %w", e) }

	fmt.Println(exRun(rxCfg, MonadMapLeft(Left[exConfig, int](rxErr), wrap)))
	// Output: 0 wrapped: boom
}

// ExampleChainLeft recovers from an error with a new computation.
func ExampleChainLeft() {
	recoverOne := ChainLeft(func(_ error) ReaderReaderIOResult[exConfig, int] {
		return rxScale(1)
	})

	fmt.Println(exRun(rxCfg, recoverOne(Left[exConfig, int](rxErr))))
	fmt.Println(exRun(rxCfg, recoverOne(Of[exConfig](1))))
	// Output:
	// 10 <nil>
	// 1 <nil>
}

// ExampleMonadChainLeft is the non curried version of ChainLeft.
func ExampleMonadChainLeft() {
	fmt.Println(exRun(rxCfg, MonadChainLeft(Left[exConfig, int](rxErr), F.Constant1[error](Of[exConfig](0)))))
	// Output: 0 <nil>
}

// ExampleChainFirstLeft runs a computation on the error and keeps the original error.
func ExampleChainFirstLeft() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), ChainFirstLeft[int](rxLogErr))))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleMonadChainFirstLeft is the non curried version of ChainFirstLeft.
func ExampleMonadChainFirstLeft() {
	fmt.Println(exRun(rxCfg, MonadChainFirstLeft(Of[exConfig](1), rxLogErr)))
	// Output: 1 <nil>
}

// ExampleTapLeft is an alias for ChainFirstLeft.
func ExampleTapLeft() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), TapLeft[int](rxLogErr))))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleMonadTapLeft is the non curried version of TapLeft.
func ExampleMonadTapLeft() {
	fmt.Println(exRun(rxCfg, MonadTapLeft(Left[exConfig, int](rxErr), rxLogErr)))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleChainFirstLeftIOK runs a side effect on the error and keeps it.
func ExampleChainFirstLeftIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), ChainFirstLeftIOK[int, exConfig](rxPrintErr))))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleMonadChainFirstLeftIOK is the non curried version of ChainFirstLeftIOK.
func ExampleMonadChainFirstLeftIOK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstLeftIOK(Left[exConfig, int](rxErr), rxPrintErr)))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleTapLeftIOK is an alias for ChainFirstLeftIOK; nothing happens on success.
func ExampleTapLeftIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](1), TapLeftIOK[int, exConfig](rxPrintErr))))
	// Output: 1 <nil>
}

// ExampleMonadTapLeftIOK is the non curried version of TapLeftIOK.
func ExampleMonadTapLeftIOK() {
	fmt.Println(exRun(rxCfg, MonadTapLeftIOK(Left[exConfig, int](rxErr), rxPrintErr)))
	// Output:
	// error: boom
	// 0 boom
}

// ExampleFromIdiomatic lifts an idiomatic Kleisli arrow.
func ExampleFromIdiomatic() {
	fmt.Println(exRun(rxCfg, FromIdiomatic(rxScaleI)(4)))
	// Output: 40 <nil>
}

// ExampleChainI chains an idiomatic Kleisli arrow.
func ExampleChainI() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainI(rxScaleI))))
	// Output: 40 <nil>
}

// ExampleMonadChainI is the non curried version of ChainI.
func ExampleMonadChainI() {
	fmt.Println(exRun(rxCfg, MonadChainI(Of[exConfig](-4), rxPositiveI)))
	// Output: 0 not positive
}

// ExampleChainFirstI runs an idiomatic Kleisli arrow and keeps the original value.
func ExampleChainFirstI() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), ChainFirstI(rxScaleI))))
	// Output: 4 <nil>
}

// ExampleMonadChainFirstI is the non curried version of ChainFirstI.
func ExampleMonadChainFirstI() {
	fmt.Println(exRun(rxCfg, MonadChainFirstI(Of[exConfig](-4), rxPositiveI)))
	// Output: 0 not positive
}

// ExampleTapI is an alias for ChainFirstI.
func ExampleTapI() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](4), TapI(rxPositiveI))))
	// Output: 4 <nil>
}

// ExampleMonadTapI is the non curried version of TapI.
func ExampleMonadTapI() {
	fmt.Println(exRun(rxCfg, MonadTapI(Of[exConfig](4), rxScaleI)))
	// Output: 4 <nil>
}

// ExampleChainLeftI recovers from an error with an idiomatic function.
func ExampleChainLeftI() {
	recoverI := ChainLeftI(func(_ error) func(context.Context, exConfig) (int, error) {
		return func(_ context.Context, c exConfig) (int, error) {
			return c.Multiplier, nil
		}
	})

	fmt.Println(exRun(rxCfg, recoverI(Left[exConfig, int](rxErr))))
	// Output: 10 <nil>
}

// ExampleMonadChainLeftI is the non curried version of ChainLeftI.
func ExampleMonadChainLeftI() {
	rethrow := func(e error) func(context.Context, exConfig) (int, error) {
		return func(context.Context, exConfig) (int, error) {
			return 0, fmt.Errorf("rethrown: %w", e)
		}
	}

	fmt.Println(exRun(rxCfg, MonadChainLeftI(Left[exConfig, int](rxErr), rethrow)))
	// Output: 0 rethrown: boom
}

// ExampleTraverseArrayI applies an idiomatic function to each element and collects the results.
func ExampleTraverseArrayI() {
	fmt.Println(exRun(rxCfg, TraverseArrayI(rxScaleI)([]int{1, 2, 3})))
	fmt.Println(exRun(rxCfg, TraverseArrayI(rxPositiveI)([]int{1, -2, 3})))
	// Output:
	// [10 20 30] <nil>
	// [] not positive
}
