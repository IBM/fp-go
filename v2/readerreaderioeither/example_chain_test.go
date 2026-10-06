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

package readerreaderioeither

import (
	"fmt"
	"strconv"
	"strings"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	RD "github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	RO "github.com/IBM/fp-go/v2/readeroption"
	S "github.com/IBM/fp-go/v2/string"
)

// exScale multiplies n with the factor of the outer environment.
func exScale(n int) exRRIOE[int] {
	return Asks[exReq, string](F.Flow2(exFactor, N.Mul(n)))
}

// exGreet greets the user of the inner environment n times.
func exGreet(n int) exRRIOE[string] {
	return func(_ exCfg) ReaderIOEither[exReq, string, string] {
		return func(r exReq) IOEither[string, string] {
			return IOE.Of[string](strings.Repeat("hi "+r.User+" ", n))
		}
	}
}

// exScaleR multiplies n with the factor as a Reader on the outer environment.
func exScaleR(n int) RD.Reader[exCfg, int] {
	return F.Flow2(exFactor, N.Mul(n))
}

// exScaleRIO multiplies n with the factor as a ReaderIO on the outer environment.
func exScaleRIO(n int) RIO.ReaderIO[exCfg, int] {
	return RIO.Asks(exScaleR(n))
}

// exAboveFactor fails if n does not exceed the factor of the outer environment.
func exAboveFactor(n int) RE.ReaderEither[exCfg, string, int] {
	return func(c exCfg) Either[string, int] {
		return exPositive(n - c.Factor)
	}
}

// exAboveFactorRO is None if n does not exceed the factor of the outer environment.
func exAboveFactorRO(n int) RO.ReaderOption[exCfg, int] {
	return func(c exCfg) Option[int] {
		return O.FromPredicate(N.MoreThan(c.Factor))(n)
	}
}

// exPrintInt prints an int as a side effect.
var exPrintInt = io.Printf[int]("log: %d\n")

// exPrintErr prints an error as a side effect.
var exPrintErr = io.Printf[string]("error: %s\n")

// ExampleMap transforms the success value.
func ExampleMap() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](21),
		Map[exCfg, exReq, string](N.Mul(2)),
	)))
	// Output: Right[int](42)
}

// ExampleMonadMap is the non curried version of Map.
func ExampleMonadMap() {
	fmt.Println(exRun(MonadMap(Of[exCfg, exReq, string](21), N.Mul(2))))
	// Output: Right[int](42)
}

// ExampleMapTo replaces the success value.
func ExampleMapTo() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](21),
		MapTo[exCfg, exReq, string, int]("done"),
	)))
	// Output: Right[string](done)
}

// ExampleMonadMapTo is the non curried version of MapTo.
func ExampleMonadMapTo() {
	fmt.Println(exRun(MonadMapTo(Of[exCfg, exReq, string](21), "done")))
	// Output: Right[string](done)
}

// ExampleChain sequences a computation that depends on the previous result.
func ExampleChain() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](2),
		Chain(exGreet),
	)))
	// Output: Right[string](hi alice hi alice )
}

// ExampleMonadChain is the non curried version of Chain.
func ExampleMonadChain() {
	fmt.Println(exRun(MonadChain(Of[exCfg, exReq, string](4), exScale)))
	// Output: Right[int](40)
}

// ExampleChainFirst runs a computation for its effect and keeps the original value.
func ExampleChainFirst() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainFirst(exScale),
	)))
	// Output: Right[int](4)
}

// ExampleMonadChainFirst is the non curried version of ChainFirst.
func ExampleMonadChainFirst() {
	fmt.Println(exRun(MonadChainFirst(Of[exCfg, exReq, string](4), exGreet)))
	// Output: Right[int](4)
}

// ExampleTap is an alias for ChainFirst.
func ExampleTap() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		Tap(exScale),
	)))
	// Output: Right[int](4)
}

// ExampleMonadTap is the non curried version of Tap.
func ExampleMonadTap() {
	fmt.Println(exRun(MonadTap(Left[exCfg, exReq, int]("boom"), exScale)))
	// Output: Left[string](boom)
}

// ExampleChainEitherK chains a function that returns an Either.
func ExampleChainEitherK() {
	positive := ChainEitherK[exCfg, exReq](exPositive)

	fmt.Println(exRun(positive(Of[exCfg, exReq, string](3))))
	fmt.Println(exRun(positive(Of[exCfg, exReq, string](-3))))
	// Output:
	// Right[int](3)
	// Left[string](not positive)
}

// ExampleMonadChainEitherK is the non curried version of ChainEitherK.
func ExampleMonadChainEitherK() {
	fmt.Println(exRun(MonadChainEitherK(Of[exCfg, exReq, string](-3), exPositive)))
	// Output: Left[string](not positive)
}

// ExampleChainFirstEitherK validates the value with an Either and keeps it.
func ExampleChainFirstEitherK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](3),
		ChainFirstEitherK[exCfg, exReq](exPositive),
	)))
	// Output: Right[int](3)
}

// ExampleMonadChainFirstEitherK is the non curried version of ChainFirstEitherK.
func ExampleMonadChainFirstEitherK() {
	fmt.Println(exRun(MonadChainFirstEitherK(Of[exCfg, exReq, string](0), exPositive)))
	// Output: Left[string](not positive)
}

// ExampleTapEitherK is an alias for ChainFirstEitherK.
func ExampleTapEitherK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](0),
		TapEitherK[exCfg, exReq](exPositive),
	)))
	// Output: Left[string](not positive)
}

// ExampleMonadTapEitherK is the non curried version of TapEitherK.
func ExampleMonadTapEitherK() {
	fmt.Println(exRun(MonadTapEitherK(Of[exCfg, exReq, string](3), exPositive)))
	// Output: Right[int](3)
}

// ExampleChainReaderK chains a Reader on the outer environment.
func ExampleChainReaderK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainReaderK[exReq, string](exScaleR),
	)))
	// Output: Right[int](40)
}

// ExampleMonadChainReaderK is the non curried version of ChainReaderK.
func ExampleMonadChainReaderK() {
	fmt.Println(exRun(MonadChainReaderK(Of[exCfg, exReq, string](4), exScaleR)))
	// Output: Right[int](40)
}

// ExampleChainFirstReaderK runs a Reader and keeps the original value.
func ExampleChainFirstReaderK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainFirstReaderK[exReq, string](exScaleR),
	)))
	// Output: Right[int](4)
}

// ExampleMonadChainFirstReaderK is the non curried version of ChainFirstReaderK.
func ExampleMonadChainFirstReaderK() {
	fmt.Println(exRun(MonadChainFirstReaderK(Of[exCfg, exReq, string](4), exScaleR)))
	// Output: Right[int](4)
}

// ExampleTapReaderK is an alias for ChainFirstReaderK.
func ExampleTapReaderK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		TapReaderK[exReq, string](exScaleR),
	)))
	// Output: Right[int](4)
}

// ExampleMonadTapReaderK is the non curried version of TapReaderK.
func ExampleMonadTapReaderK() {
	fmt.Println(exRun(MonadTapReaderK(Of[exCfg, exReq, string](4), exScaleR)))
	// Output: Right[int](4)
}

// ExampleChainReaderIOK chains a ReaderIO on the outer environment.
func ExampleChainReaderIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainReaderIOK[exReq, string](exScaleRIO),
	)))
	// Output: Right[int](40)
}

// ExampleMonadChainReaderIOK is the non curried version of ChainReaderIOK.
func ExampleMonadChainReaderIOK() {
	fmt.Println(exRun(MonadChainReaderIOK(Of[exCfg, exReq, string](4), exScaleRIO)))
	// Output: Right[int](40)
}

// ExampleChainFirstReaderIOK runs a ReaderIO and keeps the original value.
func ExampleChainFirstReaderIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainFirstReaderIOK[exReq, string](exScaleRIO),
	)))
	// Output: Right[int](4)
}

// ExampleMonadChainFirstReaderIOK is the non curried version of ChainFirstReaderIOK.
func ExampleMonadChainFirstReaderIOK() {
	fmt.Println(exRun(MonadChainFirstReaderIOK(Of[exCfg, exReq, string](4), exScaleRIO)))
	// Output: Right[int](4)
}

// ExampleTapReaderIOK is an alias for ChainFirstReaderIOK.
func ExampleTapReaderIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		TapReaderIOK[exReq, string](exScaleRIO),
	)))
	// Output: Right[int](4)
}

// ExampleMonadTapReaderIOK is the non curried version of TapReaderIOK.
func ExampleMonadTapReaderIOK() {
	fmt.Println(exRun(MonadTapReaderIOK(Of[exCfg, exReq, string](4), exScaleRIO)))
	// Output: Right[int](4)
}

// ExampleChainReaderEitherK chains a ReaderEither on the outer environment.
func ExampleChainReaderEitherK() {
	aboveFactor := ChainReaderEitherK[exReq](exAboveFactor)

	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](15))))
	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](5))))
	// Output:
	// Right[int](5)
	// Left[string](not positive)
}

// ExampleMonadChainReaderEitherK is the non curried version of ChainReaderEitherK.
func ExampleMonadChainReaderEitherK() {
	fmt.Println(exRun(MonadChainReaderEitherK(Of[exCfg, exReq, string](15), exAboveFactor)))
	// Output: Right[int](5)
}

// ExampleChainFirstReaderEitherK validates with a ReaderEither and keeps the value.
func ExampleChainFirstReaderEitherK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](15),
		ChainFirstReaderEitherK[exReq](exAboveFactor),
	)))
	// Output: Right[int](15)
}

// ExampleMonadChainFirstReaderEitherK is the non curried version of ChainFirstReaderEitherK.
func ExampleMonadChainFirstReaderEitherK() {
	fmt.Println(exRun(MonadChainFirstReaderEitherK(Of[exCfg, exReq, string](5), exAboveFactor)))
	// Output: Left[string](not positive)
}

// ExampleTapReaderEitherK is an alias for ChainFirstReaderEitherK.
func ExampleTapReaderEitherK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](15),
		TapReaderEitherK[exReq](exAboveFactor),
	)))
	// Output: Right[int](15)
}

// ExampleMonadTapReaderEitherK is the non curried version of TapReaderEitherK.
func ExampleMonadTapReaderEitherK() {
	fmt.Println(exRun(MonadTapReaderEitherK(Of[exCfg, exReq, string](15), exAboveFactor)))
	// Output: Right[int](15)
}

// ExampleChainReaderIOEitherK chains a ReaderIOEither on the outer environment.
func ExampleChainReaderIOEitherK() {
	scale := func(n int) RIOE.ReaderIOEither[exCfg, string, int] {
		return RIOE.Asks[string](exScaleR(n))
	}

	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](4),
		ChainReaderIOEitherK[exReq](scale),
	)))
	// Output: Right[int](40)
}

// ExampleChainReaderOptionK chains a ReaderOption and maps None to an error.
func ExampleChainReaderOptionK() {
	aboveFactor := ChainReaderOptionK[exCfg, exReq, int, int](lazy.Of("too small"))(exAboveFactorRO)

	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](15))))
	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](5))))
	// Output:
	// Right[int](15)
	// Left[string](too small)
}

// ExampleChainFirstReaderOptionK validates with a ReaderOption and keeps the value.
func ExampleChainFirstReaderOptionK() {
	aboveFactor := ChainFirstReaderOptionK[exCfg, exReq, int, int](lazy.Of("too small"))(exAboveFactorRO)

	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](15))))
	// Output: Right[int](15)
}

// ExampleTapReaderOptionK is an alias for ChainFirstReaderOptionK.
func ExampleTapReaderOptionK() {
	aboveFactor := TapReaderOptionK[exCfg, exReq, int, int](lazy.Of("too small"))(exAboveFactorRO)

	fmt.Println(exRun(aboveFactor(Of[exCfg, exReq, string](5))))
	// Output: Left[string](too small)
}

// ExampleChainIOEitherK chains a function that returns an IOEither.
func ExampleChainIOEitherK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](-1),
		ChainIOEitherK[exCfg, exReq](F.Flow2(exPositive, IOE.FromEither[string, int])),
	)))
	// Output: Left[string](not positive)
}

// ExampleMonadChainIOEitherK is the non curried version of ChainIOEitherK.
func ExampleMonadChainIOEitherK() {
	fmt.Println(exRun(MonadChainIOEitherK(
		Of[exCfg, exReq, string](1),
		F.Flow2(exPositive, IOE.FromEither[string, int]),
	)))
	// Output: Right[int](1)
}

// ExampleChainIOK chains a side effect that cannot fail.
func ExampleChainIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](7),
		ChainIOK[exCfg, exReq, string](F.Flow2(strconv.Itoa, io.Of[string])),
	)))
	// Output: Right[string](7)
}

// ExampleMonadChainIOK is the non curried version of ChainIOK.
func ExampleMonadChainIOK() {
	fmt.Println(exRun(MonadChainIOK(Of[exCfg, exReq, string](7), exPrintInt)))
	// Output:
	// log: 7
	// Right[int](7)
}

// ExampleChainFirstIOK runs a side effect and keeps the original value.
func ExampleChainFirstIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](7),
		ChainFirstIOK[exCfg, exReq, string](exPrintInt),
	)))
	// Output:
	// log: 7
	// Right[int](7)
}

// ExampleMonadChainFirstIOK is the non curried version of ChainFirstIOK.
func ExampleMonadChainFirstIOK() {
	fmt.Println(exRun(MonadChainFirstIOK(Of[exCfg, exReq, string](7), exPrintInt)))
	// Output:
	// log: 7
	// Right[int](7)
}

// ExampleTapIOK is an alias for ChainFirstIOK; the side effect is skipped on failure.
func ExampleTapIOK() {
	fmt.Println(exRun(F.Pipe1(
		Left[exCfg, exReq, int]("boom"),
		TapIOK[exCfg, exReq, string](exPrintInt),
	)))
	// Output: Left[string](boom)
}

// ExampleMonadTapIOK is the non curried version of TapIOK.
func ExampleMonadTapIOK() {
	fmt.Println(exRun(MonadTapIOK(Of[exCfg, exReq, string](7), exPrintInt)))
	// Output:
	// log: 7
	// Right[int](7)
}

// ExampleChainOptionK chains a function that returns an Option and maps None to an error.
func ExampleChainOptionK() {
	positive := ChainOptionK[exCfg, exReq, int, int](lazy.Of("not positive"))(O.FromPredicate(N.MoreThan(0)))

	fmt.Println(exRun(positive(Of[exCfg, exReq, string](3))))
	fmt.Println(exRun(positive(Of[exCfg, exReq, string](-3))))
	// Output:
	// Right[int](3)
	// Left[string](not positive)
}

// ExampleAp applies a wrapped function to a wrapped value.
func ExampleAp() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](N.Mul(2)),
		Ap[int](Of[exCfg, exReq, string](21)),
	)))
	// Output: Right[int](42)
}

// ExampleMonadAp is the non curried version of Ap.
func ExampleMonadAp() {
	fmt.Println(exRun(MonadAp(Of[exCfg, exReq, string](N.Mul(2)), exScale(2))))
	// Output: Right[int](40)
}

// ExampleMonadApSeq evaluates the function before the value.
func ExampleMonadApSeq() {
	fmt.Println(exRun(MonadApSeq(Of[exCfg, exReq, string](S.Format[int]("#%d")), exScale(2))))
	// Output: Right[string](#20)
}

// ExampleMonadApPar evaluates the function and the value concurrently.
func ExampleMonadApPar() {
	fmt.Println(exRun(MonadApPar(Of[exCfg, exReq, string](N.Mul(2)), Left[exCfg, exReq, int]("boom"))))
	// Output: Left[string](boom)
}

// ExampleFlap applies a wrapped function to a plain value.
func ExampleFlap() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](N.Mul(2)),
		Flap[exCfg, exReq, string, int](21),
	)))
	// Output: Right[int](42)
}

// ExampleMonadFlap is the non curried version of Flap.
func ExampleMonadFlap() {
	fmt.Println(exRun(MonadFlap(Of[exCfg, exReq, string](N.Mul(2)), 21)))
	// Output: Right[int](42)
}

// ExampleAlt falls back to another computation on failure.
func ExampleAlt() {
	fallback := Alt(lazy.Of(Of[exCfg, exReq, string](0)))

	fmt.Println(exRun(fallback(Of[exCfg, exReq, string](42))))
	fmt.Println(exRun(fallback(Left[exCfg, exReq, int]("boom"))))
	// Output:
	// Right[int](42)
	// Right[int](0)
}

// ExampleMonadAlt is the non curried version of Alt.
func ExampleMonadAlt() {
	fmt.Println(exRun(MonadAlt(Left[exCfg, exReq, int]("boom"), lazy.Of(exScale(1)))))
	// Output: Right[int](10)
}

// ExampleMapLeft transforms the error.
func ExampleMapLeft() {
	fmt.Println(exRun(F.Pipe1(
		Left[exCfg, exReq, int]("boom"),
		MapLeft[exCfg, exReq, int](strings.ToUpper),
	)))
	// Output: Left[string](BOOM)
}

// ExampleMonadMapLeft is the non curried version of MapLeft; it can change the error type.
func ExampleMonadMapLeft() {
	res := MonadMapLeft(Left[exCfg, exReq, int]("boom"), S.Size)(exConfig)(exRequest)()

	fmt.Println(res)
	// Output: Left[int](4)
}

// ExampleChainLeft recovers from an error with a new computation.
func ExampleChainLeft() {
	recoverLen := ChainLeft(func(e string) exRRIOE[int] {
		return exScale(len(e))
	})

	fmt.Println(exRun(recoverLen(Left[exCfg, exReq, int]("boom"))))
	fmt.Println(exRun(recoverLen(Of[exCfg, exReq, string](1))))
	// Output:
	// Right[int](40)
	// Right[int](1)
}

// ExampleMonadChainLeft is the non curried version of ChainLeft.
func ExampleMonadChainLeft() {
	res := MonadChainLeft(
		Left[exCfg, exReq, int]("boom"),
		F.Flow2(strings.ToUpper, Left[exCfg, exReq, int, string]),
	)

	fmt.Println(exRun(res))
	// Output: Left[string](BOOM)
}

// ExampleChainFirstLeft runs a computation on the error and keeps the original error.
func ExampleChainFirstLeft() {
	logErr := func(e string) exRRIOE[F.Void] {
		return FromIO[exCfg, exReq, string](io.FromImpure(func() { fmt.Println("error:", e) }))
	}

	fmt.Println(exRun(F.Pipe1(
		Left[exCfg, exReq, int]("boom"),
		ChainFirstLeft[int](logErr),
	)))
	// Output:
	// error: boom
	// Left[string](boom)
}

// ExampleMonadChainFirstLeft is the non curried version of ChainFirstLeft.
func ExampleMonadChainFirstLeft() {
	fmt.Println(exRun(MonadChainFirstLeft(
		Of[exCfg, exReq, string](1),
		F.Flow2(strings.ToUpper, Left[exCfg, exReq, int, string]),
	)))
	// Output: Right[int](1)
}

// ExampleTapLeft is an alias for ChainFirstLeft.
func ExampleTapLeft() {
	fmt.Println(exRun(F.Pipe1(
		Left[exCfg, exReq, int]("boom"),
		TapLeft[int](F.Flow2(exPrintErr, FromIO[exCfg, exReq, string, string])),
	)))
	// Output:
	// error: boom
	// Left[string](boom)
}

// ExampleMonadTapLeft is the non curried version of TapLeft.
func ExampleMonadTapLeft() {
	fmt.Println(exRun(MonadTapLeft(
		Left[exCfg, exReq, int]("boom"),
		F.Flow2(exPrintErr, FromIO[exCfg, exReq, string, string]),
	)))
	// Output:
	// error: boom
	// Left[string](boom)
}

// ExampleChainFirstLeftIOK runs a side effect on the error and keeps it.
func ExampleChainFirstLeftIOK() {
	fmt.Println(exRun(F.Pipe1(
		Left[exCfg, exReq, int]("boom"),
		ChainFirstLeftIOK[int, exCfg, exReq](exPrintErr),
	)))
	// Output:
	// error: boom
	// Left[string](boom)
}

// ExampleMonadChainFirstLeftIOK is the non curried version of ChainFirstLeftIOK.
func ExampleMonadChainFirstLeftIOK() {
	fmt.Println(exRun(MonadChainFirstLeftIOK(Left[exCfg, exReq, int]("boom"), exPrintErr)))
	// Output:
	// error: boom
	// Left[string](boom)
}

// ExampleTapLeftIOK is an alias for ChainFirstLeftIOK; nothing happens on success.
func ExampleTapLeftIOK() {
	fmt.Println(exRun(F.Pipe1(
		Of[exCfg, exReq, string](1),
		TapLeftIOK[int, exCfg, exReq](exPrintErr),
	)))
	// Output: Right[int](1)
}

// ExampleMonadTapLeftIOK is the non curried version of TapLeftIOK.
func ExampleMonadTapLeftIOK() {
	fmt.Println(exRun(MonadTapLeftIOK(Left[exCfg, exReq, int]("boom"), exPrintErr)))
	// Output:
	// error: boom
	// Left[string](boom)
}
