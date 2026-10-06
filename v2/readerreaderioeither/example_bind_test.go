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

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	RD "github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
	S "github.com/IBM/fp-go/v2/string"
)

// exState is the state that the do-notation examples build up.
type exState struct {
	X    int
	Name string
}

// exSetX sets the X field of the state.
func exSetX(x int) func(exState) exState {
	return func(s exState) exState {
		s.X = x
		return s
	}
}

// exSetName sets the Name field of the state.
func exSetName(name string) func(exState) exState {
	return func(s exState) exState {
		s.Name = name
		return s
	}
}

// exGetX reads the X field of the state.
func exGetX(s exState) int { return s.X }

// exXLens focuses on the X field of the state.
var exXLens = L.MakeLens(exGetX, func(s exState, x int) exState {
	s.X = x
	return s
})

// exDo starts the do-notation with an empty state.
var exDo = Do[exCfg, exReq, string](exState{})

// ExampleDo starts a do-notation pipeline with an initial state.
func ExampleDo() {
	fmt.Println(exRun(Do[exCfg, exReq, string](exState{X: 1})))
	// Output: Right[readerreaderioeither.exState]({1 })
}

// ExampleBind adds the result of a computation that depends on the state.
func ExampleBind() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		Bind(exSetX, F.Constant1[exState](exScale(4))),
		Bind(exSetName, F.Flow3(exGetX, S.Format[int]("x=%d"), Of[exCfg, exReq, string, string])),
	)))
	// Output: Right[readerreaderioeither.exState]({40 x=40})
}

// ExampleLet adds the result of a pure computation on the state.
func ExampleLet() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetTo[exCfg, exReq, string](exSetX, 2),
		Let[exCfg, exReq, string](exSetName, F.Flow2(exGetX, strconv.Itoa)),
	)))
	// Output: Right[readerreaderioeither.exState]({2 2})
}

// ExampleLetTo adds a constant to the state.
func ExampleLetTo() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		LetTo[exCfg, exReq, string](exSetName, "bob"),
	)))
	// Output: Right[readerreaderioeither.exState]({0 bob})
}

// ExampleBindTo starts the do-notation from an existing computation.
func ExampleBindTo() {
	fmt.Println(exRun(F.Pipe1(
		exScale(1),
		BindTo[exCfg, exReq, string](func(x int) exState { return exState{X: x} }),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleApS adds the result of an independent computation to the state.
func ExampleApS() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		ApS(exSetX, exScale(3)),
		ApS(exSetName, exGreet(1)),
	)))
	// Output: Right[readerreaderioeither.exState]({30 hi alice })
}

// ExampleApSL adds the result of an independent computation through a lens.
func ExampleApSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApSL(exXLens, exScale(3)),
	)))
	// Output: Right[readerreaderioeither.exState]({30 })
}

// ExampleBindL updates a field through a lens with a computation.
func ExampleBindL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 4),
		BindL(exXLens, exScale),
	)))
	// Output: Right[readerreaderioeither.exState]({40 })
}

// ExampleLetL updates a field through a lens with a pure function.
func ExampleLetL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 4),
		LetL[exCfg, exReq, string](exXLens, N.Mul(3)),
	)))
	// Output: Right[readerreaderioeither.exState]({12 })
}

// ExampleLetToL sets a field through a lens.
func ExampleLetToL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 4),
	)))
	// Output: Right[readerreaderioeither.exState]({4 })
}

// ExampleBindIOEitherK adds the result of an IOEither to the state.
func ExampleBindIOEitherK() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		BindIOEitherK[exCfg, exReq](exSetX, func(_ exState) IOEither[string, int] {
			return IOE.Left[int]("failed")
		}),
	)))
	// Output: Left[string](failed)
}

// ExampleBindIOK adds the result of an IO to the state.
func ExampleBindIOK() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		BindIOK[exCfg, exReq, string](exSetX, F.Constant1[exState](io.Of(5))),
	)))
	// Output: Right[readerreaderioeither.exState]({5 })
}

// ExampleBindReaderK adds the result of a Reader on the outer environment to the state.
func ExampleBindReaderK() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		BindReaderK[exReq, string](exSetX, F.Constant1[exState](RD.Reader[exCfg, int](exFactor))),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleBindReaderIOK adds the result of a ReaderIO on the outer environment to the state.
func ExampleBindReaderIOK() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		BindReaderIOK[exReq, string](exSetX, F.Constant1[exState](RIO.Asks(exFactor))),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleBindEitherK adds the result of an Either to the state.
func ExampleBindEitherK() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, -1),
		BindEitherK[exCfg, exReq](exSetX, F.Flow2(exGetX, exPositive)),
	)))
	// Output: Left[string](not positive)
}

// ExampleBindIOEitherKL updates a field through a lens with an IOEither.
func ExampleBindIOEitherKL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 3),
		BindIOEitherKL[exCfg, exReq](exXLens, F.Flow2(exPositive, IOE.FromEither[string, int])),
	)))
	// Output: Right[readerreaderioeither.exState]({3 })
}

// ExampleBindIOKL updates a field through a lens with an IO.
func ExampleBindIOKL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 3),
		BindIOKL[exCfg, exReq, string](exXLens, F.Flow2(N.Mul(2), io.Of[int])),
	)))
	// Output: Right[readerreaderioeither.exState]({6 })
}

// ExampleBindReaderKL updates a field through a lens with a Reader on the outer environment.
func ExampleBindReaderKL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 3),
		BindReaderKL[exReq, string](exXLens, exScaleR),
	)))
	// Output: Right[readerreaderioeither.exState]({30 })
}

// ExampleBindReaderIOKL updates a field through a lens with a ReaderIO on the outer environment.
func ExampleBindReaderIOKL() {
	fmt.Println(exRun(F.Pipe2(
		exDo,
		LetToL[exCfg, exReq, string](exXLens, 3),
		BindReaderIOKL[exReq, string](exXLens, exScaleRIO),
	)))
	// Output: Right[readerreaderioeither.exState]({30 })
}

// ExampleApIOEitherS adds the result of an independent IOEither to the state.
func ExampleApIOEitherS() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApIOEitherS[exCfg, exReq](exSetX, IOE.Of[string](7)),
	)))
	// Output: Right[readerreaderioeither.exState]({7 })
}

// ExampleApIOS adds the result of an independent IO to the state.
func ExampleApIOS() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApIOS[exCfg, exReq, string](exSetName, io.Of("bob")),
	)))
	// Output: Right[readerreaderioeither.exState]({0 bob})
}

// ExampleApReaderS adds the result of an independent Reader on the outer environment to the state.
func ExampleApReaderS() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApReaderS[exReq, string](exSetX, exFactor),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleApReaderIOS adds the result of an independent ReaderIO on the outer environment to the state.
func ExampleApReaderIOS() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApReaderIOS[exReq, string](exSetX, RIO.Asks(exFactor)),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleApEitherS adds an independent Either to the state.
func ExampleApEitherS() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApEitherS[exCfg, exReq](exSetX, E.Right[string](7)),
	)))
	// Output: Right[readerreaderioeither.exState]({7 })
}

// ExampleApIOEitherSL sets a field through a lens from an independent IOEither.
func ExampleApIOEitherSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApIOEitherSL[exCfg, exReq](exXLens, IOE.Of[string](7)),
	)))
	// Output: Right[readerreaderioeither.exState]({7 })
}

// ExampleApIOSL sets a field through a lens from an independent IO.
func ExampleApIOSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApIOSL[exCfg, exReq, string](exXLens, io.Of(7)),
	)))
	// Output: Right[readerreaderioeither.exState]({7 })
}

// ExampleApReaderSL sets a field through a lens from an independent Reader.
func ExampleApReaderSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApReaderSL[exReq, string](exXLens, exFactor),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleApReaderIOSL sets a field through a lens from an independent ReaderIO.
func ExampleApReaderIOSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApReaderIOSL[exReq, string](exXLens, RIO.Asks(exFactor)),
	)))
	// Output: Right[readerreaderioeither.exState]({10 })
}

// ExampleApEitherSL sets a field through a lens from an independent Either.
func ExampleApEitherSL() {
	fmt.Println(exRun(F.Pipe1(
		exDo,
		ApEitherSL[exCfg, exReq](exXLens, E.Left[int]("invalid")),
	)))
	// Output: Left[string](invalid)
}
