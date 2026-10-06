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
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	RD "github.com/IBM/fp-go/v2/reader"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/result"
)

// rxMakeCfg builds the outer environment from a plain multiplier.
func rxMakeCfg(m int) exConfig { return exConfig{Multiplier: m} }

// rxRunWith runs a computation whose outer environment is a plain multiplier.
func rxRunWith[A any](m int, eff ReaderReaderIOResult[int, A]) (A, error) {
	return result.Unwrap(eff(m)(context.Background())())
}

// ExampleLocal changes the outer environment with a pure function.
func ExampleLocal() {
	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), Local[int](rxMakeCfg))))
	// Output: 12 <nil>
}

// ExampleLocalIOK changes the outer environment with an IO.
func ExampleLocalIOK() {
	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalIOK[int](F.Flow2(rxMakeCfg, io.Of[exConfig])))))
	// Output: 12 <nil>
}

// ExampleLocalIOEitherK changes the outer environment with an IOEither that may fail.
func ExampleLocalIOEitherK() {
	toCfg := F.Flow3(rxPositive, result.Map(rxMakeCfg), ioresult.FromResult[exConfig])

	fmt.Println(rxRunWith(-3, F.Pipe1(rxScale(4), LocalIOEitherK[int](toCfg))))
	// Output: 0 not positive
}

// ExampleLocalIOResultK changes the outer environment with an IOResult that may fail.
func ExampleLocalIOResultK() {
	toCfg := F.Flow3(rxPositive, result.Map(rxMakeCfg), ioresult.FromResult[exConfig])

	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalIOResultK[int](toCfg))))
	// Output: 12 <nil>
}

// ExampleLocalResultK changes the outer environment with a Result that may fail.
func ExampleLocalResultK() {
	toCfg := F.Flow2(rxPositive, result.Map(rxMakeCfg))

	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalResultK[int](toCfg))))
	// Output: 12 <nil>
}

// ExampleLocalReaderIOEitherK derives the outer environment with a ReaderIOResult on the context.
func ExampleLocalReaderIOEitherK() {
	toCfg := func(m int) RIOR.ReaderIOResult[context.Context, exConfig] {
		return RIOR.Of[context.Context](rxMakeCfg(m * 2))
	}

	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalReaderIOEitherK[int](toCfg))))
	// Output: 24 <nil>
}

// ExampleLocalReaderIOResultK derives the outer environment with a ReaderIOResult on the context.
func ExampleLocalReaderIOResultK() {
	toCfg := func(_ int) RIOR.ReaderIOResult[context.Context, exConfig] {
		return RIOR.Left[context.Context, exConfig](errors.New("no config"))
	}

	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalReaderIOResultK[int](toCfg))))
	// Output: 0 no config
}

// ExampleLocalReaderK derives the outer environment with a Reader on the context.
func ExampleLocalReaderK() {
	type key struct{}
	toCfg := func(m int) RD.Reader[context.Context, exConfig] {
		return func(ctx context.Context) exConfig {
			return rxMakeCfg(m * ctx.Value(key{}).(int))
		}
	}

	eff := F.Pipe1(rxScale(4), LocalReaderK[int](toCfg))
	ctx := context.WithValue(context.Background(), key{}, 5)

	fmt.Println(eff(3)(ctx)())
	// Output: Right[int](60)
}

// ExampleLocalReaderReaderIOEitherK derives the outer environment with a full computation.
func ExampleLocalReaderReaderIOEitherK() {
	toCfg := func(m int) ReaderReaderIOResult[int, exConfig] {
		return Asks(func(outer int) exConfig { return rxMakeCfg(m + outer) })
	}

	fmt.Println(rxRunWith(3, F.Pipe1(rxScale(4), LocalReaderReaderIOEitherK[int](toCfg))))
	// Output: 24 <nil>
}

// ExampleSequence swaps the order of the outer environments of a nested computation.
func ExampleSequence() {
	nested := Of[string](rxScale(4))

	fmt.Println(Sequence(nested)(rxCfg)("unused")(context.Background())())
	// Output: Right[int](40)
}

// ExampleSequenceReader turns a computation that produces a Reader into one that takes its environment first.
func ExampleSequenceReader() {
	nested := Of[string](RD.Reader[exConfig, int](rxMultiplier))

	fmt.Println(SequenceReader(nested)(rxCfg)("unused")(context.Background())())
	// Output: Right[int](10)
}

// ExampleSequenceReaderIO turns a computation that produces a ReaderIO into one that takes its environment first.
func ExampleSequenceReaderIO() {
	nested := Of[string](RIO.Asks(rxMultiplier))

	fmt.Println(SequenceReaderIO(nested)(rxCfg)("unused")(context.Background())())
	// Output: Right[int](10)
}

// ExampleTraverse maps with a computation on another outer environment and moves that environment first.
func ExampleTraverse() {
	traversed := Traverse[string](rxScale)(Of[string](4))

	fmt.Println(traversed(rxCfg)("unused")(context.Background())())
	// Output: Right[int](40)
}

// ExampleTraverseReader maps with a Reader and moves its environment first.
func ExampleTraverseReader() {
	traversed := TraverseReader[string](rxScaleR)(Of[string](4))

	fmt.Println(traversed(rxCfg)("unused")(context.Background())())
	// Output: Right[int](40)
}

// ExampleTraverseArray applies a computation to each element and collects the results.
func ExampleTraverseArray() {
	fmt.Println(exRun(rxCfg, TraverseArray(rxScale)([]int{1, 2, 3})))
	fmt.Println(exRun(rxCfg, TraverseArray(F.Flow2(rxPositive, FromResult[exConfig, int]))([]int{1, -2, 3})))
	// Output:
	// [10 20 30] <nil>
	// [] not positive
}

// ExampleApplicativeMonoid combines the results of two computations with a monoid.
func ExampleApplicativeMonoid() {
	m := ApplicativeMonoid[exConfig](N.MonoidSum[int]())

	fmt.Println(exRun(rxCfg, m.Concat(rxScale(1), rxScale(2))))
	fmt.Println(exRun(rxCfg, m.Empty()))
	// Output:
	// 30 <nil>
	// 0 <nil>
}

// ExampleApplicativeMonoidSeq combines the results of two computations sequentially.
func ExampleApplicativeMonoidSeq() {
	m := ApplicativeMonoidSeq[exConfig](N.MonoidSum[int]())

	fmt.Println(exRun(rxCfg, m.Concat(rxScale(1), Left[exConfig, int](rxErr))))
	// Output: 0 boom
}

// ExampleApplicativeMonoidPar combines the results of two computations concurrently.
func ExampleApplicativeMonoidPar() {
	m := ApplicativeMonoidPar[exConfig](N.MonoidSum[int]())

	fmt.Println(exRun(rxCfg, m.Concat(rxScale(1), rxScale(2))))
	// Output: 30 <nil>
}

// ExampleAlternativeMonoid combines successes and skips failures.
func ExampleAlternativeMonoid() {
	m := AlternativeMonoid[exConfig](N.MonoidSum[int]())

	fmt.Println(exRun(rxCfg, m.Concat(rxScale(1), rxScale(2))))
	fmt.Println(exRun(rxCfg, m.Concat(Left[exConfig, int](rxErr), rxScale(2))))
	// Output:
	// 30 <nil>
	// 20 <nil>
}

// ExampleAltMonoid picks the first successful computation.
func ExampleAltMonoid() {
	m := AltMonoid(lazy.Of(Left[exConfig, string](errors.New("none"))))

	fmt.Println(exRun(rxCfg, m.Concat(Left[exConfig, string](rxErr), Of[exConfig]("first"))))
	fmt.Println(exRun(rxCfg, m.Concat(m.Empty(), F.Pipe1(rxScale(1), Map[exConfig](strconv.Itoa)))))
	// Output:
	// first <nil>
	// 10 <nil>
}
