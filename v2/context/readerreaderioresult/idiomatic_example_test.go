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

package readerreaderioresult

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// This monad stacks four things: the outer environment R, a context.Context,
// IO, and error handling. Real Go code rarely needs all four at once, so the
// examples below show one adapter per layer - each taking only what it needs
// and returning the ordinary (value, error) pair.
// ---------------------------------------------------------------------------

// exConfig is the outer environment the examples read from.
type exConfig struct {
	Multiplier int
	Table      map[string]string
}

// exRun executes an effect with the given config and a background context, and
// returns the idiomatic (value, error) pair so examples can print it directly.
func exRun[A any](cfg exConfig, eff ReaderReaderIOResult[exConfig, A]) (A, error) {
	return result.Unwrap(eff(cfg)(context.Background())())
}

// exParse only needs the value: no config, no context, no IO.
func exParse(s string) (int, error) {
	return strconv.Atoi(s)
}

// exScaleByConfig needs the config, but no context and no IO.
func exScaleByConfig(n int) func(exConfig) (int, error) {
	return func(cfg exConfig) (int, error) {
		if cfg.Multiplier == 0 {
			return 0, errors.New("multiplier is not configured")
		}
		return n * cfg.Multiplier, nil
	}
}

// exLoad needs to perform a side effect, but neither config nor context.
func exLoad(n int) func() (string, error) {
	return func() (string, error) {
		if n < 0 {
			return "", fmt.Errorf("cannot load %d", n)
		}
		return fmt.Sprintf("record-%d", n), nil
	}
}

// exLookup needs the config and performs a side effect, but not the context.
func exLookup(key string) func(exConfig) func() (string, error) {
	return func(cfg exConfig) func() (string, error) {
		return func() (string, error) {
			v, ok := cfg.Table[key]
			if !ok {
				return "", fmt.Errorf("no entry for %q", key)
			}
			return v, nil
		}
	}
}

// exAudit performs a side effect and reports nothing interesting - the typical
// shape of a function passed to Tap.
func exAudit(n int) func() (int, error) {
	return func() (int, error) {
		fmt.Println("audited:", n)
		return n, nil
	}
}

// exFullStack needs everything: the value, the context and the config. This is
// the shape [FromIdiomatic] and [ChainI] already handle.
func exFullStack(n int) func(context.Context, exConfig) (string, error) {
	return func(ctx context.Context, cfg exConfig) (string, error) {
		if err := ctx.Err(); err != nil {
			return "", err
		}
		return fmt.Sprintf("%d x %d", n, cfg.Multiplier), nil
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

// ExampleFromResultI demonstrates lifting an already evaluated (value, error)
// pair, as produced by any ordinary Go call.
func ExampleFromResultI() {
	fmt.Println(exRun(exConfig{}, FromResultI[exConfig](strconv.Atoi("42"))))
	fmt.Println(exRun(exConfig{}, FromResultI[exConfig](strconv.Atoi("nope"))))

	// Output:
	// 42 <nil>
	// 0 strconv.Atoi: parsing "nope": invalid syntax
}

// ExampleFromIOResultI demonstrates lifting a deferred `func() (A, error)`.
func ExampleFromIOResultI() {
	fmt.Println(exRun(exConfig{}, FromIOResultI[exConfig](exLoad(7))))

	// Output:
	// record-7 <nil>
}

// ExampleFromReaderResultI demonstrates lifting a `func(R) (A, error)`, i.e. a
// computation that reads the environment but performs no IO.
func ExampleFromReaderResultI() {
	fmt.Println(exRun(exConfig{Multiplier: 4}, FromReaderResultI(exScaleByConfig(5))))

	// Output:
	// 20 <nil>
}

// ExampleFromReaderIOResultI demonstrates lifting a
// `func(R) func() (A, error)`: the environment plus IO, but no context.
func ExampleFromReaderIOResultI() {
	cfg := exConfig{Table: map[string]string{"host": "localhost"}}

	fmt.Println(exRun(cfg, FromReaderIOResultI(exLookup("host"))))
	fmt.Println(exRun(cfg, FromReaderIOResultI(exLookup("port"))))

	// Output:
	// localhost <nil>
	//  no entry for "port"
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

// ExampleChainResultIK demonstrates chaining a plain `func(A) (B, error)`, the
// most common shape in existing Go code.
func ExampleChainResultIK() {
	eff := F.Pipe1(
		Of[exConfig]("42"),
		ChainResultIK[exConfig](exParse),
	)

	fmt.Println(exRun(exConfig{}, eff))

	// Output:
	// 42 <nil>
}

// ExampleChainResultIK_failure demonstrates that the returned error
// short-circuits the remainder of the pipeline.
func ExampleChainResultIK_failure() {
	reached := false

	eff := F.Pipe2(
		Of[exConfig]("not a number"),
		ChainResultIK[exConfig](exParse),
		ChainResultIK[exConfig](func(n int) (int, error) {
			reached = true // never happens
			return n, nil
		}),
	)

	value, err := exRun(exConfig{}, eff)
	fmt.Println(value, err, reached)

	// Output:
	// 0 strconv.Atoi: parsing "not a number": invalid syntax false
}

// ExampleMonadChainResultIK demonstrates the uncurried form of ChainResultIK.
func ExampleMonadChainResultIK() {
	fmt.Println(exRun(exConfig{}, MonadChainResultIK(Of[exConfig]("42"), exParse)))

	// Output:
	// 42 <nil>
}

// ExampleTapResultIK demonstrates a validation step that keeps the original
// value on success.
func ExampleTapResultIK() {
	eff := F.Pipe1(
		Of[exConfig]("42"),
		TapResultIK[exConfig](exParse),
	)

	fmt.Println(exRun(exConfig{}, eff))

	// Output:
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

// ExampleChainReaderResultIK demonstrates chaining a
// `func(A) func(R) (B, error)`, i.e. a step that reads the environment.
func ExampleChainReaderResultIK() {
	eff := F.Pipe1(
		Of[exConfig](5),
		ChainReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, eff))

	// Output:
	// 15 <nil>
}

// ExampleChainReaderResultIK_failure demonstrates the error path, here caused by
// a missing configuration value.
func ExampleChainReaderResultIK_failure() {
	eff := F.Pipe1(
		Of[exConfig](5),
		ChainReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{}, eff))

	// Output:
	// 0 multiplier is not configured
}

// ExampleTapReaderResultIK demonstrates a check against the environment that
// preserves the original value.
func ExampleTapReaderResultIK() {
	eff := F.Pipe1(
		Of[exConfig](5),
		TapReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, eff))

	// Output:
	// 5 <nil>
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

// ExampleChainIOResultIK demonstrates chaining a deferred
// `func(A) func() (B, error)`.
func ExampleChainIOResultIK() {
	eff := F.Pipe1(
		Of[exConfig](7),
		ChainIOResultIK[exConfig](exLoad),
	)

	fmt.Println(exRun(exConfig{}, eff))

	// Output:
	// record-7 <nil>
}

// ExampleTapIOResultIK demonstrates an audit step: the side effect runs, its
// result is discarded and the original value flows on.
func ExampleTapIOResultIK() {
	eff := F.Pipe1(
		Of[exConfig](21),
		TapIOResultIK[exConfig](exAudit),
	)

	fmt.Println(exRun(exConfig{}, eff))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ---------------------------------------------------------------------------
// Option level
// ---------------------------------------------------------------------------

// ExampleChainOptionIK demonstrates chaining a comma-ok lookup, turning a false
// flag into an error.
func ExampleChainOptionIK() {
	table := map[string]string{"host": "localhost"}
	lookup := func(k string) (string, bool) {
		v, ok := table[k]
		return v, ok
	}
	chain := ChainOptionIK[exConfig, string, string](func() error {
		return errors.New("key not found")
	})

	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("host"), chain(lookup))))
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("port"), chain(lookup))))

	// Output:
	// localhost <nil>
	//  key not found
}

// ---------------------------------------------------------------------------
// Putting it together
// ---------------------------------------------------------------------------

// Example_idiomaticLayers mixes every layer in one pipeline: a pure fallible
// step, an environment-reading step, a deferred side effect, and finally the
// full-stack ChainI. Each function is plain Go.
func Example_idiomaticLayers() {
	eff := F.Pipe3(
		Of[exConfig]("7"),
		ChainResultIK[exConfig](exParse),     // "7" -> 7
		ChainReaderResultIK(exScaleByConfig), // 7 -> 21
		ChainI(exFullStack),                  // 21 -> "21 x 3"
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, eff))

	// Output:
	// 21 x 3 <nil>
}
