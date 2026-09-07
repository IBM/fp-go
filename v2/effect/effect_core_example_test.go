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
	"context"
	"fmt"

	ctxthunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// ---------------------------------------------------------------------------
// Shared fixtures
// ---------------------------------------------------------------------------

// exCoreConfig is the dependency used by the core examples.
type exCoreConfig struct {
	Multiplier int
	Prefix     string
}

// exCoreMultiplier and exCorePrefix are named getters, so the pipelines below
// stay point-free.
func exCoreMultiplier(c exCoreConfig) int { return c.Multiplier }
func exCorePrefix(c exCoreConfig) string  { return c.Prefix }
func exCoreLabel(n int) string            { return fmt.Sprintf("n=%d", n) }
func exCoreErr(msg string) error          { return fmt.Errorf("%s", msg) }
func exCoreHalve(n int) Result[int]       { return R.Of(n / 2) }
func exCoreRequirePositive(n int) Result[int] {
	if n <= 0 {
		return R.Left[int](exCoreErr("not positive"))
	}
	return R.Of(n)
}

// exCoreRun provides the dependency and runs the effect, returning (value, error).
func exCoreRun[A any](cfg exCoreConfig, eff Effect[exCoreConfig, A]) (A, error) {
	return RunSync(Provide[A](cfg)(eff))(context.Background())
}

// ---------------------------------------------------------------------------
// Constructors
// ---------------------------------------------------------------------------

// ExampleOf demonstrates lifting a plain value into an Effect that always
// succeeds and ignores its dependency.
func ExampleOf() {
	fmt.Println(exCoreRun(exCoreConfig{}, Of[exCoreConfig](42)))
	// Output:
	// 42 <nil>
}

// ExampleSucceed demonstrates that Succeed is a synonym of Of, named for
// readability at the head of a pipeline.
func ExampleSucceed() {
	fmt.Println(exCoreRun(exCoreConfig{}, Succeed[exCoreConfig]("ready")))
	// Output:
	// ready <nil>
}

// ExampleFail demonstrates constructing an Effect that always fails.
func ExampleFail() {
	fmt.Println(exCoreRun(exCoreConfig{}, Fail[exCoreConfig, int](exCoreErr("boom"))))
	// Output:
	// 0 boom
}

// ExampleFromResult demonstrates lifting an already-computed Result, turning a
// pure success-or-error value into an Effect.
func ExampleFromResult() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromResult[exCoreConfig](R.Of(7))))
	fmt.Println(exCoreRun(exCoreConfig{}, FromResult[exCoreConfig](R.Left[int](exCoreErr("bad input")))))
	// Output:
	// 7 <nil>
	// 0 bad input
}

// ExampleFromIO demonstrates lifting a side effect that cannot fail.  The IO is
// not executed until the effect is run.
func ExampleFromIO() {
	executed := false
	action := func() string {
		executed = true
		return "side effect ran"
	}

	eff := FromIO[exCoreConfig](action)
	fmt.Println("before run:", executed)

	value, _ := exCoreRun(exCoreConfig{}, eff)
	fmt.Println(value)
	// Output:
	// before run: false
	// side effect ran
}

// ExampleFromThunk demonstrates lifting a Thunk — an IO that may fail and sees
// the runtime context.Context — into an Effect that ignores its dependency.
func ExampleFromThunk() {
	eff := FromThunk[exCoreConfig](ctxthunk.Of("from thunk"))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// from thunk <nil>
}

// ExampleAsks demonstrates reading a projection of the dependency, which is the
// point-free way to depend on a single field.
func ExampleAsks() {
	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 6}, Asks(exCoreMultiplier)))
	fmt.Println(exCoreRun(exCoreConfig{Prefix: "LOG"}, Asks(exCorePrefix)))
	// Output:
	// 6 <nil>
	// LOG <nil>
}

// ExampleSuspend demonstrates deferring the *construction* of an effect until
// it is run, which is how you break recursion or delay expensive setup.
func ExampleSuspend() {
	built := 0

	eff := Suspend(func() Effect[exCoreConfig, int] {
		built++
		return Of[exCoreConfig](42)
	})

	fmt.Println("after Suspend:", built)

	exCoreRun(exCoreConfig{}, eff)
	exCoreRun(exCoreConfig{}, eff)
	fmt.Println("after two runs:", built)
	// Output:
	// after Suspend: 0
	// after two runs: 2
}

// ---------------------------------------------------------------------------
// Map / Chain
// ---------------------------------------------------------------------------

// ExampleMap demonstrates transforming the success value with a pure function.
func ExampleMap() {
	eff := F.Pipe2(
		Of[exCoreConfig](21),
		Map[exCoreConfig](N.Mul(2)),
		Map[exCoreConfig](exCoreLabel),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// n=42 <nil>
}

// ExampleMap_errorsPropagate demonstrates that Map is skipped on failure.
func ExampleMap_errorsPropagate() {
	eff := F.Pipe1(
		Fail[exCoreConfig, int](exCoreErr("upstream failed")),
		Map[exCoreConfig](N.Mul(2)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 0 upstream failed
}

// ExampleChain demonstrates sequencing an effect that itself depends on the
// previous result.
func ExampleChain() {
	// Kleisli arrow: scale the input by the configured multiplier.
	scale := func(n int) Effect[exCoreConfig, int] {
		return F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(n)))
	}

	eff := F.Pipe1(Of[exCoreConfig](5), Chain(scale))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))
	// Output:
	// 15 <nil>
}

// ExampleChain_shortCircuits demonstrates that a failing step skips everything
// after it.
func ExampleChain_shortCircuits() {
	reject := func(int) Effect[exCoreConfig, int] {
		return Fail[exCoreConfig, int](exCoreErr("rejected"))
	}

	eff := F.Pipe2(
		Of[exCoreConfig](5),
		Chain(reject),
		Map[exCoreConfig](N.Mul(100)), // never reached
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 0 rejected
}

// ExampleChainFirst demonstrates running an effect for its side effect while
// keeping the original value.
func ExampleChainFirst() {
	audit := func(n int) Effect[exCoreConfig, string] {
		return FromIO[exCoreConfig](func() string {
			fmt.Println("audited:", n)
			return "logged"
		})
	}

	eff := F.Pipe1(Of[exCoreConfig](42), ChainFirst(audit))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 42
	// 42 <nil>
}

// ExampleTap demonstrates that Tap is ChainFirst under a name that reads better
// when the intent is observation rather than sequencing.
func ExampleTap() {
	observe := func(n int) Effect[exCoreConfig, string] {
		return FromIO[exCoreConfig](func() string {
			fmt.Println("saw:", n)
			return ""
		})
	}

	eff := F.Pipe2(Of[exCoreConfig](7), Tap(observe), Map[exCoreConfig](N.Mul(2)))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// saw: 7
	// 14 <nil>
}

// ExampleAp demonstrates applicative application: a function inside an effect is
// applied to a value inside another effect.  Unlike Chain, neither side depends
// on the other's result.
func ExampleAp() {
	// An effect producing a function, built point-free from the dependency.
	fnEff := F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul[int]))

	eff := F.Pipe1(fnEff, Ap[int](Of[exCoreConfig](5)))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// 20 <nil>
}

// ExampleTernary demonstrates branching between two Kleisli arrows on a
// predicate, without writing an if statement inside the pipeline.
func ExampleTernary() {
	classify := Ternary(
		N.MoreThan(0),
		F.Flow2(exCoreLabel, Of[exCoreConfig, string]),
		F.Flow2(exCoreLabel, F.Flow2(S.Format[string]("negative(%s)"), Of[exCoreConfig, string])),
	)

	eff := F.Pipe1(Of[exCoreConfig](5), Chain(classify))
	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	eff = F.Pipe1(Of[exCoreConfig](-5), Chain(classify))
	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// n=5 <nil>
	// negative(n=-5) <nil>
}

// ---------------------------------------------------------------------------
// Chain*K — lifting foreign Kleisli arrows
// ---------------------------------------------------------------------------

// ExampleChainResultK demonstrates chaining a pure function that may fail,
// without wrapping it in an Effect by hand.
func ExampleChainResultK() {
	eff := F.Pipe2(
		Of[exCoreConfig](10),
		ChainResultK[exCoreConfig](exCoreRequirePositive),
		ChainResultK[exCoreConfig](exCoreHalve),
	)
	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	failing := F.Pipe1(Of[exCoreConfig](-1), ChainResultK[exCoreConfig](exCoreRequirePositive))
	fmt.Println(exCoreRun(exCoreConfig{}, failing))
	// Output:
	// 5 <nil>
	// 0 not positive
}

// ExampleChainIOK demonstrates chaining a side effect that cannot fail.
func ExampleChainIOK() {
	eff := F.Pipe1(
		Of[exCoreConfig](21),
		ChainIOK[exCoreConfig](F.Flow2(N.Mul(2), io.Of[int])),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 42 <nil>
}

// ExampleChainFirstIOK demonstrates running an IO for its side effect while
// keeping the upstream value.
func ExampleChainFirstIOK() {
	logIt := func(n int) io.IO[string] {
		return func() string {
			fmt.Println("logged:", n)
			return ""
		}
	}

	eff := F.Pipe1(Of[exCoreConfig](42), ChainFirstIOK[exCoreConfig](logIt))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// logged: 42
	// 42 <nil>
}

// ExampleTapIOK demonstrates that TapIOK is ChainFirstIOK under an
// observation-flavoured name.
func ExampleTapIOK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("payload"),
		TapIOK[exCoreConfig](func(s string) io.IO[string] {
			return func() string {
				fmt.Println("tapped:", s)
				return ""
			}
		}),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// tapped: payload
	// payload <nil>
}

// ExampleChainThunkK demonstrates chaining a Thunk — an IO that may fail and can
// read the runtime context.Context — into the pipeline.
func ExampleChainThunkK() {
	eff := F.Pipe1(
		Of[exCoreConfig](21),
		ChainThunkK[exCoreConfig](F.Flow2(N.Mul(2), ctxthunk.Of[int])),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 42 <nil>
}

// ExampleChainFirstThunkK demonstrates running a Thunk for its side effect while
// keeping the upstream value.
func ExampleChainFirstThunkK() {
	audit := func(n int) ctxthunk.ReaderIOResult[string] {
		return ctxthunk.FromIO(func() string {
			fmt.Println("audited:", n)
			return ""
		})
	}

	eff := F.Pipe1(Of[exCoreConfig](42), ChainFirstThunkK[exCoreConfig](audit))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 42
	// 42 <nil>
}

// ExampleTapThunkK demonstrates that TapThunkK is ChainFirstThunkK under an
// observation-flavoured name.
func ExampleTapThunkK() {
	eff := F.Pipe1(
		Of[exCoreConfig](7),
		TapThunkK[exCoreConfig](func(n int) ctxthunk.ReaderIOResult[string] {
			return ctxthunk.FromIO(func() string {
				fmt.Println("tapped:", n)
				return ""
			})
		}),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// tapped: 7
	// 7 <nil>
}

// ExampleChainReaderK demonstrates chaining a pure computation that reads the
// dependency but cannot fail.
func ExampleChainReaderK() {
	// int -> Reader[exCoreConfig, int]
	scale := func(n int) reader.Reader[exCoreConfig, int] {
		return F.Flow2(exCoreMultiplier, N.Mul(n))
	}

	eff := F.Pipe1(Of[exCoreConfig](5), ChainReaderK(scale))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))
	// Output:
	// 15 <nil>
}

// ExampleChainReaderIOK demonstrates chaining a computation that reads the
// dependency *and* performs IO, but cannot fail.
func ExampleChainReaderIOK() {
	render := func(n int) readerio.ReaderIO[exCoreConfig, string] {
		return func(cfg exCoreConfig) io.IO[string] {
			return func() string { return fmt.Sprintf("%s: %d", cfg.Prefix, n) }
		}
	}

	eff := F.Pipe1(Of[exCoreConfig](42), ChainReaderIOK(render))

	fmt.Println(exCoreRun(exCoreConfig{Prefix: "RESULT"}, eff))
	// Output:
	// RESULT: 42 <nil>
}

// ---------------------------------------------------------------------------
// Read / Paired
// ---------------------------------------------------------------------------

// ExampleRead demonstrates supplying the dependency, producing a Thunk.  Provide
// is the same function under the name used by the run pipeline.
func ExampleRead() {
	eff := F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(2)))

	thunk := Read[int](exCoreConfig{Multiplier: 21})(eff)

	fmt.Println(RunSync(thunk)(context.Background()))
	// Output:
	// 42 <nil>
}

// ExamplePaired demonstrates repackaging an Effect so that the runtime context
// and the dependency arrive together as a single Pair argument.
func ExamplePaired() {
	eff := F.Pipe1(Asks(exCoreMultiplier), Map[exCoreConfig](N.Mul(3)))

	// Paired turns Effect[C, A] into a function of Pair[context.Context, C].
	run := Paired(eff)

	fmt.Println(run(pair.MakePair(context.Background(), exCoreConfig{Multiplier: 5}))())
	// Output:
	// Right[int](15)
}
