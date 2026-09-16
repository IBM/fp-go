// Copyright (c) 2024 IBM Corp.
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
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// An Effect stacks four things: the dependency C, a context.Context, IO, and
// error handling. Real Go code rarely needs all four at once, so the functions
// below each take only what they need and return the ordinary (value, error)
// pair. The `...IK` operators lift them at the right layer.
//
// They reuse exCoreConfig and exCoreRun from effect_core_example_test.go.
// ---------------------------------------------------------------------------

// exIKParse only needs the value: no dependency, no context, no IO.
func exIKParse(s string) (int, error) {
	return strconv.Atoi(s)
}

// exIKScale needs the dependency, but no context and no IO.
func exIKScale(n int) func(exCoreConfig) (int, error) {
	return func(cfg exCoreConfig) (int, error) {
		if cfg.Multiplier == 0 {
			return 0, errors.New("multiplier is not configured")
		}
		return n * cfg.Multiplier, nil
	}
}

// exIKLoad performs a side effect, but needs neither dependency nor context.
func exIKLoad(n int) func() (string, error) {
	return func() (string, error) {
		if n < 0 {
			return "", fmt.Errorf("cannot load %d", n)
		}
		return fmt.Sprintf("record-%d", n), nil
	}
}

// exIKFetch needs the context and performs a side effect, but not the
// dependency - the shape of most context-aware client calls.
func exIKFetch(path string) func(context.Context) func() (string, error) {
	return func(ctx context.Context) func() (string, error) {
		return func() (string, error) {
			if err := ctx.Err(); err != nil {
				return "", err
			}
			return "fetched:" + path, nil
		}
	}
}

// exIKAudit performs a side effect and reports nothing interesting - the
// typical shape of a function passed to Tap.
func exIKAudit(n int) func() (int, error) {
	return func() (int, error) {
		fmt.Println("audited:", n)
		return n, nil
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

// ExampleFromResultI demonstrates lifting an already evaluated (value, error)
// pair, as produced by any ordinary Go call.
func ExampleFromResultI() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromResultI[exCoreConfig](strconv.Atoi("42"))))
	fmt.Println(exCoreRun(exCoreConfig{}, FromResultI[exCoreConfig](strconv.Atoi("nope"))))

	// Output:
	// 42 <nil>
	// 0 strconv.Atoi: parsing "nope": invalid syntax
}

// ExampleFromIOResultI demonstrates lifting a deferred `func() (A, error)`.
func ExampleFromIOResultI() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromIOResultI[exCoreConfig](exIKLoad(7))))

	// Output:
	// record-7 <nil>
}

// ExampleFromReaderResultI demonstrates lifting a `func(C) (A, error)`, i.e. a
// computation that reads the dependency but performs no IO.
func ExampleFromReaderResultI() {
	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, FromReaderResultI(exIKScale(5))))

	// Output:
	// 20 <nil>
}

// ExampleFromThunkI demonstrates lifting a
// `func(context.Context) func() (A, error)`: the context plus IO, but no
// dependency.
func ExampleFromThunkI() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromThunkI[exCoreConfig](exIKFetch("users"))))

	// Output:
	// fetched:users <nil>
}

// ---------------------------------------------------------------------------
// ChainResultIK
// ---------------------------------------------------------------------------

// ExampleChainResultIK demonstrates chaining a plain `func(A) (B, error)`, the
// most common shape in existing Go code.
func ExampleChainResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("42"),
		ChainResultIK[exCoreConfig](exIKParse),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// 42 <nil>
}

// ExampleChainResultIK_failure demonstrates that the returned error
// short-circuits the remainder of the pipeline.
func ExampleChainResultIK_failure() {
	reached := false

	eff := F.Pipe2(
		Of[exCoreConfig]("not a number"),
		ChainResultIK[exCoreConfig](exIKParse),
		ChainResultIK[exCoreConfig](func(n int) (int, error) {
			reached = true // never happens
			return n, nil
		}),
	)

	value, err := exCoreRun(exCoreConfig{}, eff)
	fmt.Println(value, err, reached)

	// Output:
	// 0 strconv.Atoi: parsing "not a number": invalid syntax false
}

// ExampleMonadChainResultIK demonstrates the uncurried form of ChainResultIK.
func ExampleMonadChainResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainResultIK(Of[exCoreConfig]("42"), exIKParse)))

	// Output:
	// 42 <nil>
}

// ExampleTapResultIK demonstrates a validation step that keeps the original
// value on success.
func ExampleTapResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("42"),
		TapResultIK[exCoreConfig](exIKParse),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ChainReaderResultIK
// ---------------------------------------------------------------------------

// ExampleChainReaderResultIK demonstrates chaining a
// `func(A) func(C) (B, error)`, i.e. a step that reads the dependency.
func ExampleChainReaderResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig](5),
		ChainReaderResultIK(exIKScale),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))

	// Output:
	// 15 <nil>
}

// ExampleChainReaderResultIK_failure demonstrates the error path, here caused by
// a missing configuration value.
func ExampleChainReaderResultIK_failure() {
	eff := F.Pipe1(
		Of[exCoreConfig](5),
		ChainReaderResultIK(exIKScale),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// 0 multiplier is not configured
}

// ExampleTapReaderResultIK demonstrates a check against the dependency that
// preserves the original value.
func ExampleTapReaderResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig](5),
		TapReaderResultIK(exIKScale),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))

	// Output:
	// 5 <nil>
}

// ---------------------------------------------------------------------------
// ChainIOResultIK
// ---------------------------------------------------------------------------

// ExampleChainIOResultIK demonstrates chaining a deferred
// `func(A) func() (B, error)`.
func ExampleChainIOResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig](7),
		ChainIOResultIK[exCoreConfig](exIKLoad),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// record-7 <nil>
}

// ExampleTapIOResultIK demonstrates an audit step: the side effect runs, its
// result is discarded and the original value flows on.
func ExampleTapIOResultIK() {
	eff := F.Pipe1(
		Of[exCoreConfig](21),
		TapIOResultIK[exCoreConfig](exIKAudit),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ---------------------------------------------------------------------------
// ChainThunkIK
// ---------------------------------------------------------------------------

// ExampleChainThunkIK demonstrates chaining a
// `func(A) func(context.Context) func() (B, error)`, the shape of a
// context-aware call that may fail but needs no dependency.
func ExampleChainThunkIK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("users"),
		ChainThunkIK[exCoreConfig](exIKFetch),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// fetched:users <nil>
}

// ExampleChainThunkIK_cancellation demonstrates that a cancelled context
// surfaces as an ordinary error.
func ExampleChainThunkIK_cancellation() {
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()

	eff := F.Pipe1(
		Of[exCoreConfig]("users"),
		ChainThunkIK[exCoreConfig](exIKFetch),
	)

	fmt.Println(RunSync(Provide[string](exCoreConfig{})(eff))(cancelled))

	// Output:
	//  context canceled
}

// ExampleTapThunkIK demonstrates a context-aware side effect that preserves the
// original value.
func ExampleTapThunkIK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("users"),
		TapThunkIK[exCoreConfig, string](exIKFetch),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// users <nil>
}

// ExampleTapLeftThunkIK demonstrates reporting an error through a context-aware
// handler while leaving the failure itself untouched.
func ExampleTapLeftThunkIK() {
	report := func(err error) func(context.Context) func() (string, error) {
		return func(context.Context) func() (string, error) {
			return func() (string, error) {
				fmt.Println("reporting:", err)
				return "", nil
			}
		}
	}

	eff := F.Pipe1(
		Fail[exCoreConfig, int](errors.New("boom")),
		TapLeftThunkIK[exCoreConfig, int](report),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// reporting: boom
	// 0 boom
}

// ---------------------------------------------------------------------------
// ChainOptionIK
// ---------------------------------------------------------------------------

// ExampleChainOptionIK demonstrates chaining a comma-ok lookup, turning a false
// flag into an error.
func ExampleChainOptionIK() {
	table := map[string]string{"host": "localhost"}
	lookup := func(k string) (string, bool) {
		v, ok := table[k]
		return v, ok
	}
	chain := ChainOptionIK[exCoreConfig, string, string](func() error {
		return errors.New("key not found")
	})

	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("host"), chain(lookup))))
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("port"), chain(lookup))))

	// Output:
	// localhost <nil>
	//  key not found
}

// ---------------------------------------------------------------------------
// Putting it together
// ---------------------------------------------------------------------------

// Example_idiomaticLayers mixes every layer in one pipeline: a pure fallible
// step, a dependency-reading step, a deferred side effect and a context-aware
// call. Every function involved is plain Go.
func Example_idiomaticLayers() {
	eff := F.Pipe3(
		Of[exCoreConfig]("7"),
		ChainResultIK[exCoreConfig](exIKParse),  // "7" -> 7
		ChainReaderResultIK(exIKScale),          // 7 -> 21
		ChainIOResultIK[exCoreConfig](exIKLoad), // 21 -> "record-21"
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))

	// Output:
	// record-21 <nil>
}
