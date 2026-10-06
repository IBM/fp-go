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
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	RD "github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
)

// exMakeCfg builds the outer environment from a plain factor.
func exMakeCfg(factor int) exCfg { return exCfg{Factor: factor} }

// exRunWith runs a computation whose outer environment is a plain factor.
func exRunWith[A any](factor int, m ReaderReaderIOEither[int, exReq, string, A]) Either[string, A] {
	return m(factor)(exRequest)()
}

// ExampleLocal changes the outer environment with a pure function.
func ExampleLocal() {
	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		Local[exReq, string, int](exMakeCfg),
	)))
	// Output: Right[int](12)
}

// ExampleLocalIOK changes the outer environment with an IO.
func ExampleLocalIOK() {
	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		LocalIOK[exReq, string, int](F.Flow2(exMakeCfg, io.Of[exCfg])),
	)))
	// Output: Right[int](12)
}

// ExampleLocalIOEitherK changes the outer environment with an IOEither that may fail.
func ExampleLocalIOEitherK() {
	toCfg := F.Flow3(exPositive, E.Map[string](exMakeCfg), IOE.FromEither[string, exCfg])

	fmt.Println(exRunWith(-3, F.Pipe1(
		exScale(4),
		LocalIOEitherK[exReq, int](toCfg),
	)))
	// Output: Left[string](not positive)
}

// ExampleLocalEitherK changes the outer environment with an Either that may fail.
func ExampleLocalEitherK() {
	toCfg := F.Flow2(exPositive, E.Map[string](exMakeCfg))

	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		LocalEitherK[exReq, int](toCfg),
	)))
	// Output: Right[int](12)
}

// ExampleLocalReaderIOEitherK derives the outer environment from the old one and the inner environment.
func ExampleLocalReaderIOEitherK() {
	toCfg := func(n int) RIOE.ReaderIOEither[exReq, string, exCfg] {
		return RIOE.Asks[string](func(r exReq) exCfg { return exCfg{Factor: n * len(r.User)} })
	}

	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		LocalReaderIOEitherK[int](toCfg),
	)))
	// Output: Right[int](60)
}

// ExampleLocalReaderK derives the outer environment from the old one and the inner environment.
func ExampleLocalReaderK() {
	toCfg := func(n int) RD.Reader[exReq, exCfg] {
		return func(r exReq) exCfg { return exCfg{Factor: n * len(r.User)} }
	}

	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		LocalReaderK[string, int](toCfg),
	)))
	// Output: Right[int](60)
}

// ExampleLocalReaderReaderIOEitherK derives the outer environment with a full computation.
func ExampleLocalReaderReaderIOEitherK() {
	toCfg := func(n int) ReaderReaderIOEither[int, exReq, string, exCfg] {
		return Asks[exReq, string](func(outer int) exCfg { return exCfg{Factor: n + outer} })
	}

	fmt.Println(exRunWith(3, F.Pipe1(
		exScale(4),
		LocalReaderReaderIOEitherK[int](toCfg),
	)))
	// Output: Right[int](24)
}

// ExampleSequence swaps the order of the outer environments of a nested computation.
func ExampleSequence() {
	nested := Of[string, exReq, string](exScale(4))

	fmt.Println(Sequence(nested)(exConfig)("unused")(exRequest)())
	// Output: Right[int](40)
}

// ExampleSequenceReader turns a computation that produces a Reader into one that takes its environment first.
func ExampleSequenceReader() {
	nested := Of[string, exReq, string](RD.Reader[exCfg, int](exFactor))

	fmt.Println(SequenceReader(nested)(exConfig)("unused")(exRequest)())
	// Output: Right[int](10)
}

// ExampleSequenceReaderIO turns a computation that produces a ReaderIO into one that takes its environment first.
func ExampleSequenceReaderIO() {
	nested := Of[string, exReq, string](RIO.Asks(exFactor))

	fmt.Println(SequenceReaderIO(nested)(exConfig)("unused")(exRequest)())
	// Output: Right[int](10)
}

// ExampleTraverse maps with a computation on another outer environment and moves that environment first.
func ExampleTraverse() {
	traversed := Traverse[string](exScale)(Of[string, exReq, string](4))

	fmt.Println(traversed(exConfig)("unused")(exRequest)())
	// Output: Right[int](40)
}

// ExampleTraverseReader maps with a Reader and moves its environment first.
func ExampleTraverseReader() {
	traversed := TraverseReader[string, exCfg, exReq, string](exScaleR)(Of[string, exReq, string](4))

	fmt.Println(traversed(exConfig)("unused")(exRequest)())
	// Output: Right[int](40)
}

// ExampleTraverseArray applies a computation to each element and collects the results.
func ExampleTraverseArray() {
	fmt.Println(exRun(TraverseArray(exScale)([]int{1, 2, 3})))
	fmt.Println(exRun(TraverseArray(F.Flow2(exPositive, FromEither[exCfg, exReq, string, int]))([]int{1, -2, 3})))
	// Output:
	// Right[[]int]([10 20 30])
	// Left[string](not positive)
}

// ExampleApplicativeMonoid combines the results of two computations with a monoid.
func ExampleApplicativeMonoid() {
	m := ApplicativeMonoid[exCfg, exReq, string](N.MonoidSum[int]())

	fmt.Println(exRun(m.Concat(exScale(1), exScale(2))))
	fmt.Println(exRun(m.Empty()))
	// Output:
	// Right[int](30)
	// Right[int](0)
}

// ExampleApplicativeMonoidSeq combines the results of two computations sequentially.
func ExampleApplicativeMonoidSeq() {
	m := ApplicativeMonoidSeq[exCfg, exReq, string](N.MonoidSum[int]())

	fmt.Println(exRun(m.Concat(exScale(1), Left[exCfg, exReq, int]("boom"))))
	// Output: Left[string](boom)
}

// ExampleApplicativeMonoidPar combines the results of two computations concurrently.
func ExampleApplicativeMonoidPar() {
	m := ApplicativeMonoidPar[exCfg, exReq, string](N.MonoidSum[int]())

	fmt.Println(exRun(m.Concat(exScale(1), exScale(2))))
	// Output: Right[int](30)
}

// ExampleAlternativeMonoid combines successes and skips failures.
func ExampleAlternativeMonoid() {
	m := AlternativeMonoid[exCfg, exReq, string](N.MonoidSum[int]())

	fmt.Println(exRun(m.Concat(exScale(1), exScale(2))))
	fmt.Println(exRun(m.Concat(Left[exCfg, exReq, int]("boom"), exScale(2))))
	// Output:
	// Right[int](30)
	// Right[int](20)
}

// ExampleAltMonoid picks the first successful computation.
func ExampleAltMonoid() {
	m := AltMonoid(lazy.Of(Left[exCfg, exReq, string]("none")))

	fmt.Println(exRun(m.Concat(Left[exCfg, exReq, string]("boom"), exGreet(1))))
	fmt.Println(exRun(m.Concat(m.Empty(), F.Pipe1(exScale(1), Map[exCfg, exReq, string](strconv.Itoa)))))
	// Output:
	// Right[string](hi alice )
	// Right[string](10)
}
