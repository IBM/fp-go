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

package effect

import (
	"fmt"

	ctxthunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
)

// exErrBoom is the failure used throughout the error-handling examples.
func exErrBoom() error { return fmt.Errorf("primary failed") }

// exErrRecover turns any error into a usable fallback value; as a Kleisli arrow
// it is what ChainLeft expects.
func exErrRecover(err error) Effect[exCoreConfig, int] {
	return Of[exCoreConfig](-1)
}

// exErrLog prints the error and yields nothing of interest — the shape every
// *Left tap wants.
func exErrLog(err error) Effect[exCoreConfig, string] {
	return FromIO[exCoreConfig](func() string {
		fmt.Println("handled:", err)
		return ""
	})
}

// exErrLogIO is the same observation expressed as a plain IO Kleisli arrow.
func exErrLogIO(err error) io.IO[string] {
	return func() string {
		fmt.Println("handled:", err)
		return ""
	}
}

// ---------------------------------------------------------------------------
// ChainLeft — recovery
// ---------------------------------------------------------------------------

// ExampleChainLeft demonstrates recovering from a failure by supplying a
// replacement effect.  A successful effect passes through untouched.
func ExampleChainLeft() {
	recovered := F.Pipe1(Fail[exCoreConfig, int](exErrBoom()), ChainLeft(exErrRecover))
	fmt.Println(exCoreRun(exCoreConfig{}, recovered))

	untouched := F.Pipe1(Of[exCoreConfig](42), ChainLeft(exErrRecover))
	fmt.Println(exCoreRun(exCoreConfig{}, untouched))
	// Output:
	// -1 <nil>
	// 42 <nil>
}

// ExampleChainLeft_rethrow demonstrates that the handler may itself fail, which
// is how you translate one error into another.
func ExampleChainLeft_rethrow() {
	translate := func(err error) Effect[exCoreConfig, int] {
		return Fail[exCoreConfig, int](fmt.Errorf("wrapped: %w", err))
	}

	eff := F.Pipe1(Fail[exCoreConfig, int](exErrBoom()), ChainLeft(translate))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 0 wrapped: primary failed
}

// ExampleMonadChainLeft demonstrates the uncurried form, which takes the effect
// and the handler as two arguments instead of returning an operator.
func ExampleMonadChainLeft() {
	eff := MonadChainLeft(Fail[exCoreConfig, int](exErrBoom()), exErrRecover)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// -1 <nil>
}

// ---------------------------------------------------------------------------
// Alt — fallback
// ---------------------------------------------------------------------------

// ExampleAlt demonstrates falling back to a second effect when the first fails.
// Unlike ChainLeft the fallback does not see the error.
func ExampleAlt() {
	eff := F.Pipe1(
		Fail[exCoreConfig, int](exErrBoom()),
		Alt(lazy.Of(Of[exCoreConfig](99))),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 99 <nil>
}

// ExampleAlt_lazy demonstrates that the fallback is only built when it is
// actually needed.
func ExampleAlt_lazy() {
	built := 0
	fallback := func() Effect[exCoreConfig, int] {
		built++
		return Of[exCoreConfig](99)
	}

	// The primary succeeds, so the fallback is never constructed.
	ok := F.Pipe1(Of[exCoreConfig](1), Alt(fallback))
	exCoreRun(exCoreConfig{}, ok)
	fmt.Println("after success:", built)

	// The primary fails, so the fallback is constructed and used.
	bad := F.Pipe1(Fail[exCoreConfig, int](exErrBoom()), Alt(fallback))
	exCoreRun(exCoreConfig{}, bad)
	fmt.Println("after failure:", built)
	// Output:
	// after success: 0
	// after failure: 1
}

// ExampleAlt_chained demonstrates trying a series of candidates in order.
func ExampleAlt_chained() {
	eff := F.Pipe2(
		Fail[exCoreConfig, int](fmt.Errorf("first unavailable")),
		Alt(lazy.Of(Fail[exCoreConfig, int](fmt.Errorf("second unavailable")))),
		Alt(lazy.Of(Of[exCoreConfig](3))),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 3 <nil>
}

// ExampleMonadAlt demonstrates the uncurried form of Alt.
func ExampleMonadAlt() {
	eff := MonadAlt(Fail[exCoreConfig, int](exErrBoom()), lazy.Of(Of[exCoreConfig](99)))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 99 <nil>
}

// ---------------------------------------------------------------------------
// *Left taps — observe the error without handling it
// ---------------------------------------------------------------------------

// ExampleChainFirstLeft demonstrates observing a failure without recovering
// from it: the original error still propagates.
func ExampleChainFirstLeft() {
	eff := F.Pipe1(
		Fail[exCoreConfig, int](exErrBoom()),
		ChainFirstLeft[int](exErrLog),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleChainFirstLeft_successIsUntouched demonstrates that the observer never
// runs on the success path.
func ExampleChainFirstLeft_successIsUntouched() {
	eff := F.Pipe2(
		Of[exCoreConfig](21),
		ChainFirstLeft[int](exErrLog),
		Map[exCoreConfig](N.Mul(2)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 42 <nil>
}

// ExampleMonadChainFirstLeft demonstrates the uncurried form.
func ExampleMonadChainFirstLeft() {
	eff := MonadChainFirstLeft(Fail[exCoreConfig, int](exErrBoom()), exErrLog)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleTapLeft demonstrates that TapLeft is ChainFirstLeft under an
// observation-flavoured name.
func ExampleTapLeft() {
	eff := F.Pipe1(Fail[exCoreConfig, int](exErrBoom()), TapLeft[int](exErrLog))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleMonadTapLeft demonstrates the uncurried form of TapLeft.
func ExampleMonadTapLeft() {
	eff := MonadTapLeft(Fail[exCoreConfig, int](exErrBoom()), exErrLog)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleChainFirstLeftIOK demonstrates observing a failure with a plain IO
// action, without wrapping it in an Effect.
func ExampleChainFirstLeftIOK() {
	eff := F.Pipe1(
		Fail[exCoreConfig, int](exErrBoom()),
		ChainFirstLeftIOK[int, exCoreConfig](exErrLogIO),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleMonadChainFirstLeftIOK demonstrates the uncurried form.
func ExampleMonadChainFirstLeftIOK() {
	eff := MonadChainFirstLeftIOK(Fail[exCoreConfig, int](exErrBoom()), exErrLogIO)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleTapLeftIOK demonstrates the observation-flavoured name for
// ChainFirstLeftIOK.
func ExampleTapLeftIOK() {
	eff := F.Pipe1(
		Fail[exCoreConfig, int](exErrBoom()),
		TapLeftIOK[int, exCoreConfig](exErrLogIO),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleMonadTapLeftIOK demonstrates the uncurried form of TapLeftIOK.
func ExampleMonadTapLeftIOK() {
	eff := MonadTapLeftIOK(Fail[exCoreConfig, int](exErrBoom()), exErrLogIO)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}

// ExampleChainFirstLeftThunkK demonstrates observing a failure with a Thunk, so
// the observer can itself perform IO that fails and can read the runtime
// context.Context.
func ExampleChainFirstLeftThunkK() {
	logThunk := func(err error) ctxthunk.ReaderIOResult[string] {
		return ctxthunk.FromIO[string](func() string {
			fmt.Println("handled:", err)
			return ""
		})
	}

	eff := F.Pipe1(
		Fail[exCoreConfig, int](exErrBoom()),
		ChainFirstLeftThunkK[exCoreConfig, int](logThunk),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// handled: primary failed
	// 0 primary failed
}
