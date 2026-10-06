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
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	O "github.com/IBM/fp-go/v2/option"
	RE "github.com/IBM/fp-go/v2/readereither"
	RIO "github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	RO "github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/retry"
	"github.com/IBM/fp-go/v2/tailrec"
)

// ExampleOf lifts a pure value into a computation that ignores both environments.
func ExampleOf() {
	fmt.Println(exRun(Of[exCfg, exReq, string](42)))
	// Output: Right[int](42)
}

// ExampleRight creates a successful computation.
func ExampleRight() {
	fmt.Println(exRun(Right[exCfg, exReq, string](42)))
	// Output: Right[int](42)
}

// ExampleLeft creates a failed computation.
func ExampleLeft() {
	fmt.Println(exRun(Left[exCfg, exReq, int]("boom")))
	// Output: Left[string](boom)
}

// ExampleFromEither lifts an Either.
func ExampleFromEither() {
	fmt.Println(exRun(FromEither[exCfg, exReq](exPositive(3))))
	fmt.Println(exRun(FromEither[exCfg, exReq](exPositive(-3))))
	// Output:
	// Right[int](3)
	// Left[string](not positive)
}

// ExampleFromIO lifts a side effect that cannot fail.
func ExampleFromIO() {
	fmt.Println(exRun(FromIO[exCfg, exReq, string](func() int {
		fmt.Println("side effect")
		return 1
	})))
	// Output:
	// side effect
	// Right[int](1)
}

// ExampleRightIO is an alias for FromIO.
func ExampleRightIO() {
	fmt.Println(exRun(RightIO[exCfg, exReq, string](lazy.Of(1))))
	// Output: Right[int](1)
}

// ExampleLeftIO lifts a side effect that produces an error.
func ExampleLeftIO() {
	fmt.Println(exRun(LeftIO[exCfg, exReq, int](lazy.Of("boom"))))
	// Output: Left[string](boom)
}

// ExampleFromIOEither lifts an IOEither.
func ExampleFromIOEither() {
	fmt.Println(exRun(FromIOEither[exCfg, exReq](IOE.Of[string](7))))
	// Output: Right[int](7)
}

// ExampleFromReader lifts a Reader on the outer environment.
func ExampleFromReader() {
	fmt.Println(exRun(FromReader[exReq, string](exFactor)))
	// Output: Right[int](10)
}

// ExampleRightReader is an alias for FromReader.
func ExampleRightReader() {
	fmt.Println(exRun(RightReader[exReq, string](exFactor)))
	// Output: Right[int](10)
}

// ExampleLeftReader derives an error from the outer environment.
func ExampleLeftReader() {
	fmt.Println(exRun(LeftReader[exReq, int](F.Flow2(exFactor, strconv.Itoa))))
	// Output: Left[string](10)
}

// ExampleFromReaderIO lifts a ReaderIO on the outer environment.
func ExampleFromReaderIO() {
	fmt.Println(exRun(FromReaderIO[exReq, string](RIO.Asks(exFactor))))
	// Output: Right[int](10)
}

// ExampleRightReaderIO is an alias for FromReaderIO.
func ExampleRightReaderIO() {
	fmt.Println(exRun(RightReaderIO[exReq, string](RIO.Asks(exFactor))))
	// Output: Right[int](10)
}

// ExampleLeftReaderIO derives an error from the outer environment with a side effect.
func ExampleLeftReaderIO() {
	fmt.Println(exRun(LeftReaderIO[exReq, int](RIO.Of[exCfg]("boom"))))
	// Output: Left[string](boom)
}

// ExampleFromReaderEither lifts a ReaderEither on the outer environment.
func ExampleFromReaderEither() {
	fmt.Println(exRun(FromReaderEither[exCfg, exReq](RE.Of[exCfg, string](5))))
	// Output: Right[int](5)
}

// ExampleFromReaderIOEither lifts a ReaderIOEither on the outer environment.
func ExampleFromReaderIOEither() {
	fmt.Println(exRun(FromReaderIOEither[exReq](RIOE.Of[exCfg, string](5))))
	// Output: Right[int](5)
}

// ExampleFromReaderOption lifts a ReaderOption and maps None to an error.
func ExampleFromReaderOption() {
	fromRO := FromReaderOption[exCfg, exReq, int](lazy.Of("missing"))

	fmt.Println(exRun(fromRO(RO.Of[exCfg](5))))
	fmt.Println(exRun(fromRO(RO.None[exCfg, int]())))
	// Output:
	// Right[int](5)
	// Left[string](missing)
}

// ExampleFromOption lifts an Option and maps None to an error.
func ExampleFromOption() {
	fromOption := FromOption[exCfg, exReq, int](lazy.Of("missing"))

	fmt.Println(exRun(fromOption(O.Some(5))))
	fmt.Println(exRun(fromOption(O.None[int]())))
	// Output:
	// Right[int](5)
	// Left[string](missing)
}

// ExampleFromPredicate succeeds if the predicate holds and fails otherwise.
func ExampleFromPredicate() {
	positive := FromPredicate[exCfg, exReq](N.MoreThan(0), strconv.Itoa)

	fmt.Println(exRun(positive(5)))
	fmt.Println(exRun(positive(-5)))
	// Output:
	// Right[int](5)
	// Left[string](-5)
}

// ExampleAsk returns the outer environment.
func ExampleAsk() {
	fmt.Println(exRun(Ask[exCfg, exReq, string]()))
	// Output: Right[readerreaderioeither.exCfg]({10})
}

// ExampleAsks projects a value from the outer environment.
func ExampleAsks() {
	fmt.Println(exRun(Asks[exReq, string](exFactor)))
	// Output: Right[int](10)
}

// ExampleRead provides the outer environment and returns a ReaderIOEither on the inner one.
func ExampleRead() {
	rioe := F.Pipe1(
		Asks[exReq, string](exFactor),
		Read[exReq, string, int](exCfg{Factor: 3}),
	)

	fmt.Println(rioe(exRequest)())
	// Output: Right[int](3)
}

// ExampleFlatten removes one level of nesting.
func ExampleFlatten() {
	nested := Of[exCfg, exReq, string](Of[exCfg, exReq, string](42))

	fmt.Println(exRun(Flatten(nested)))
	// Output: Right[int](42)
}

// ExampleDelay waits before running the computation.
func ExampleDelay() {
	delayed := F.Pipe1(
		Of[exCfg, exReq, string](42),
		Delay[exCfg, exReq, string, int](time.Millisecond),
	)

	fmt.Println(exRun(delayed))
	// Output: Right[int](42)
}

// ExampleAfter waits until the given time before running the computation.
func ExampleAfter() {
	delayed := F.Pipe1(
		Of[exCfg, exReq, string](42),
		After[exCfg, exReq, string, int](time.Now().Add(time.Millisecond)),
	)

	fmt.Println(exRun(delayed))
	// Output: Right[int](42)
}

// ExampleTailRec runs a stack safe loop; each step reads the decrement from the
// outer environment.
func ExampleTailRec() {
	countdown := TailRec(func(n int) exRRIOE[Trampoline[int, string]] {
		if n <= 0 {
			return Of[exCfg, exReq, string](tailrec.Land[int]("liftoff"))
		}
		return Asks[exReq, string](func(c exCfg) Trampoline[int, string] {
			return tailrec.Bounce[string](n - c.Factor)
		})
	})

	fmt.Println(exRun(countdown(1_000_000)))
	// Output: Right[string](liftoff)
}

// ExampleRetrying retries an action as long as it fails, up to the limit of the policy.
func ExampleRetrying() {
	action := func(status retry.RetryStatus) exRRIOE[int] {
		return FromEither[exCfg, exReq](exPositive(int(status.IterNumber) - 1))
	}

	retried := Retrying(retry.LimitRetries(5), action, E.IsLeft[string, int])

	fmt.Println(exRun(retried))
	// Output: Right[int](1)
}
