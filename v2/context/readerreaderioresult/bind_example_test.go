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
	"fmt"
	"strconv"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	RD "github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
	S "github.com/IBM/fp-go/v2/string"
)

// rxState is the state that the do-notation examples build up.
type rxState struct {
	X    int
	Name string
}

// rxSetX sets the X field of the state.
func rxSetX(x int) func(rxState) rxState {
	return func(s rxState) rxState {
		s.X = x
		return s
	}
}

// rxSetName sets the Name field of the state.
func rxSetName(name string) func(rxState) rxState {
	return func(s rxState) rxState {
		s.Name = name
		return s
	}
}

// rxGetX reads the X field of the state.
func rxGetX(s rxState) int { return s.X }

// rxXLens focuses on the X field of the state.
var rxXLens = L.MakeLens(rxGetX, func(s rxState, x int) rxState {
	s.X = x
	return s
})

// rxDo starts the do-notation with an empty state.
var rxDo = Do[exConfig](rxState{})

// ExampleDo starts a do-notation pipeline with an initial state.
func ExampleDo() {
	fmt.Println(exRun(rxCfg, Do[exConfig](rxState{X: 1})))
	// Output: {1 } <nil>
}

// ExampleBind adds the result of a computation that depends on the state.
func ExampleBind() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		Bind(rxSetX, F.Constant1[rxState](rxScale(4))),
		Bind(rxSetName, F.Flow3(rxGetX, S.Format[int]("x=%d"), Of[exConfig, string])),
	)))
	// Output: {40 x=40} <nil>
}

// ExampleLet adds the result of a pure computation on the state.
func ExampleLet() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetTo[exConfig](rxSetX, 2),
		Let[exConfig](rxSetName, F.Flow2(rxGetX, strconv.Itoa)),
	)))
	// Output: {2 2} <nil>
}

// ExampleLetTo adds a constant to the state.
func ExampleLetTo() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, LetTo[exConfig](rxSetName, "bob"))))
	// Output: {0 bob} <nil>
}

// ExampleBindTo starts the do-notation from an existing computation.
func ExampleBindTo() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxScale(1),
		BindTo[exConfig](func(x int) rxState { return rxState{X: x} }),
	)))
	// Output: {10 } <nil>
}

// ExampleApS adds the result of an independent computation to the state.
func ExampleApS() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		ApS(rxSetX, rxScale(3)),
		ApS(rxSetName, Of[exConfig]("bob")),
	)))
	// Output: {30 bob} <nil>
}

// ExampleApSL adds the result of an independent computation through a lens.
func ExampleApSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApSL(rxXLens, rxScale(3)))))
	// Output: {30 } <nil>
}

// ExampleBindL updates a field through a lens with a computation.
func ExampleBindL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 4),
		BindL(rxXLens, rxScale),
	)))
	// Output: {40 } <nil>
}

// ExampleLetL updates a field through a lens with a pure function.
func ExampleLetL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 4),
		LetL[exConfig](rxXLens, N.Mul(3)),
	)))
	// Output: {12 } <nil>
}

// ExampleLetToL sets a field through a lens.
func ExampleLetToL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, LetToL[exConfig](rxXLens, 4))))
	// Output: {4 } <nil>
}

// ExampleBindIOEitherK adds the result of an IOEither to the state.
func ExampleBindIOEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxDo,
		BindIOEitherK[exConfig](rxSetX, F.Constant1[rxState](IOE.Left[int](rxErr))),
	)))
	// Output: {0 } boom
}

// ExampleBindIOResultK adds the result of an IOResult to the state.
func ExampleBindIOResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxDo,
		BindIOResultK[exConfig](rxSetX, F.Constant1[rxState](ioresult.Of(5))),
	)))
	// Output: {5 } <nil>
}

// ExampleBindIOK adds the result of an IO to the state.
func ExampleBindIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxDo,
		BindIOK[exConfig](rxSetX, F.Constant1[rxState](io.Of(5))),
	)))
	// Output: {5 } <nil>
}

// ExampleBindReaderK adds the result of a Reader on the outer environment to the state.
func ExampleBindReaderK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxDo,
		BindReaderK(rxSetX, F.Constant1[rxState](RD.Reader[exConfig, int](rxMultiplier))),
	)))
	// Output: {10 } <nil>
}

// ExampleBindReaderIOK adds the result of a ReaderIO on the outer environment to the state.
func ExampleBindReaderIOK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(
		rxDo,
		BindReaderIOK(rxSetX, F.Constant1[rxState](RIO.Asks(rxMultiplier))),
	)))
	// Output: {10 } <nil>
}

// ExampleBindEitherK adds the result of an Either to the state.
func ExampleBindEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, -1),
		BindEitherK[exConfig](rxSetX, F.Flow2(rxGetX, rxPositive)),
	)))
	// Output: {0 } not positive
}

// ExampleBindIOEitherKL updates a field through a lens with an IOEither.
func ExampleBindIOEitherKL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 3),
		BindIOEitherKL[exConfig](rxXLens, F.Flow2(rxPositive, IOE.FromEither[error, int])),
	)))
	// Output: {3 } <nil>
}

// ExampleBindIOKL updates a field through a lens with an IO.
func ExampleBindIOKL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 3),
		BindIOKL[exConfig](rxXLens, F.Flow2(N.Mul(2), io.Of[int])),
	)))
	// Output: {6 } <nil>
}

// ExampleBindReaderKL updates a field through a lens with a Reader on the outer environment.
func ExampleBindReaderKL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 3),
		BindReaderKL(rxXLens, rxScaleR),
	)))
	// Output: {30 } <nil>
}

// ExampleBindReaderIOKL updates a field through a lens with a ReaderIO on the outer environment.
func ExampleBindReaderIOKL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 3),
		BindReaderIOKL(rxXLens, rxScaleRIO),
	)))
	// Output: {30 } <nil>
}

// ExampleBindI adds the result of an idiomatic function to the state.
func ExampleBindI() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, 3),
		BindI(rxSetX, F.Flow2(rxGetX, rxScaleI)),
	)))
	// Output: {30 } <nil>
}

// ExampleBindIL updates a field through a lens with an idiomatic function.
func ExampleBindIL() {
	fmt.Println(exRun(rxCfg, F.Pipe2(
		rxDo,
		LetToL[exConfig](rxXLens, -3),
		BindIL(rxXLens, rxPositiveI),
	)))
	// Output: {0 } not positive
}

// ExampleApIOEitherS adds the result of an independent IOEither to the state.
func ExampleApIOEitherS() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApIOEitherS[exConfig](rxSetX, IOE.Of[error](7)))))
	// Output: {7 } <nil>
}

// ExampleApIOS adds the result of an independent IO to the state.
func ExampleApIOS() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApIOS[exConfig](rxSetName, io.Of("bob")))))
	// Output: {0 bob} <nil>
}

// ExampleApReaderS adds the result of an independent Reader on the outer environment to the state.
func ExampleApReaderS() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApReaderS(rxSetX, rxMultiplier))))
	// Output: {10 } <nil>
}

// ExampleApReaderIOS adds the result of an independent ReaderIO on the outer environment to the state.
func ExampleApReaderIOS() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApReaderIOS(rxSetX, RIO.Asks(rxMultiplier)))))
	// Output: {10 } <nil>
}

// ExampleApEitherS adds an independent Either to the state.
func ExampleApEitherS() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApEitherS[exConfig](rxSetX, E.Right[error](7)))))
	// Output: {7 } <nil>
}

// ExampleApIOEitherSL sets a field through a lens from an independent IOEither.
func ExampleApIOEitherSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApIOEitherSL[exConfig](rxXLens, IOE.Of[error](7)))))
	// Output: {7 } <nil>
}

// ExampleApIOSL sets a field through a lens from an independent IO.
func ExampleApIOSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApIOSL[exConfig](rxXLens, io.Of(7)))))
	// Output: {7 } <nil>
}

// ExampleApReaderSL sets a field through a lens from an independent Reader.
func ExampleApReaderSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApReaderSL(rxXLens, rxMultiplier))))
	// Output: {10 } <nil>
}

// ExampleApReaderIOSL sets a field through a lens from an independent ReaderIO.
func ExampleApReaderIOSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApReaderIOSL(rxXLens, RIO.Asks(rxMultiplier)))))
	// Output: {10 } <nil>
}

// ExampleApEitherSL sets a field through a lens from an independent Either.
func ExampleApEitherSL() {
	fmt.Println(exRun(rxCfg, F.Pipe1(rxDo, ApEitherSL[exConfig](rxXLens, E.Left[int](rxErr)))))
	// Output: {0 } boom
}
