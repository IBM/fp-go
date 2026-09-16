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

package readerioresult

import (
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// Every function below has the shape that ordinary Go code already has: it
// returns (value, error), possibly behind a context parameter or behind a
// thunk. Nothing wraps a value into a monad by hand - that is the point of the
// idiomatic bridge.
// ---------------------------------------------------------------------------

// exConfig is the environment the examples read from.
type exConfig struct {
	Multiplier int
	Table      map[string]string
}

// exRun executes a ReaderIOResult and returns the idiomatic (value, error) pair,
// so the examples can print it the way plain Go code would.
func exRun[A any](cfg exConfig, rio ReaderIOResult[exConfig, A]) (A, error) {
	return result.Unwrap(rio(cfg)())
}

// exParse is an idiomatic Result Kleisli arrow: it neither reads the
// environment nor performs IO, it just may fail.
func exParse(s string) (int, error) {
	return strconv.Atoi(s)
}

// exScaleByConfig is an idiomatic ReaderResult Kleisli arrow: it reads the
// environment and may fail, but performs no IO.
func exScaleByConfig(n int) func(exConfig) (int, error) {
	return func(cfg exConfig) (int, error) {
		if cfg.Multiplier == 0 {
			return 0, errors.New("multiplier is not configured")
		}
		return n * cfg.Multiplier, nil
	}
}

// exLoad is an idiomatic IOResult Kleisli arrow: it performs a (here simulated)
// side effect and may fail, but does not read the environment.
func exLoad(n int) func() (string, error) {
	return func() (string, error) {
		if n < 0 {
			return "", fmt.Errorf("cannot load %d", n)
		}
		return fmt.Sprintf("record-%d", n), nil
	}
}

// exFetch is an idiomatic ReaderIOResult Kleisli arrow: it reads the
// environment, performs a side effect and may fail - the full stack.
func exFetch(key string) func(exConfig) func() (string, error) {
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

// ---------------------------------------------------------------------------
// FromIdiomatic and friends
// ---------------------------------------------------------------------------

// ExampleFromIdiomatic demonstrates lifting a plain `func(R) func() (A, error)`
// into a ReaderIOResult so it composes with the rest of the package.
func ExampleFromIdiomatic() {
	rio := FromIdiomatic(exFetch("host"))

	fmt.Println(exRun(exConfig{Table: map[string]string{"host": "localhost"}}, rio))

	// Output:
	// localhost <nil>
}

// ExampleFromResultI demonstrates lifting an already evaluated (value, error)
// pair, as produced by any ordinary Go call.
func ExampleFromResultI() {
	rio := FromResultI[exConfig](strconv.Atoi("42"))

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 42 <nil>
}

// ExampleFromResultI_failure demonstrates that a non-nil error becomes a failed
// computation.
func ExampleFromResultI_failure() {
	rio := FromResultI[exConfig](strconv.Atoi("not a number"))

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 0 strconv.Atoi: parsing "not a number": invalid syntax
}

// ExampleFromIOResultI demonstrates lifting a deferred `func() (A, error)`.
func ExampleFromIOResultI() {
	rio := FromIOResultI[exConfig](exLoad(7))

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// record-7 <nil>
}

// ExampleFromReaderResultI demonstrates lifting a `func(R) (A, error)`, i.e. a
// computation that reads the environment but performs no IO.
func ExampleFromReaderResultI() {
	rio := FromReaderResultI(exScaleByConfig(5))

	fmt.Println(exRun(exConfig{Multiplier: 4}, rio))

	// Output:
	// 20 <nil>
}

// ExampleFromReaderIOResultI demonstrates lifting the full idiomatic stack.
func ExampleFromReaderIOResultI() {
	rio := FromReaderIOResultI(exFetch("missing"))

	fmt.Println(exRun(exConfig{Table: map[string]string{}}, rio))

	// Output:
	//  no entry for "missing"
}

// ---------------------------------------------------------------------------
// ChainI / ChainFirstI / TapI
// ---------------------------------------------------------------------------

// ExampleChainI demonstrates chaining an idiomatic ReaderIOResult function
// directly, without an intermediate FromIdiomatic call.
func ExampleChainI() {
	rio := F.Pipe1(
		Of[exConfig]("host"),
		ChainI(exFetch),
	)

	fmt.Println(exRun(exConfig{Table: map[string]string{"host": "localhost"}}, rio))

	// Output:
	// localhost <nil>
}

// ExampleMonadChainI demonstrates the uncurried form of ChainI.
func ExampleMonadChainI() {
	rio := MonadChainI(Of[exConfig]("host"), exFetch)

	fmt.Println(exRun(exConfig{Table: map[string]string{"host": "localhost"}}, rio))

	// Output:
	// localhost <nil>
}

// ExampleChainFirstI demonstrates running an idiomatic function for its effect
// while keeping the original value.
func ExampleChainFirstI() {
	rio := F.Pipe1(
		Of[exConfig]("host"),
		ChainFirstI(exFetch),
	)

	fmt.Println(exRun(exConfig{Table: map[string]string{"host": "localhost"}}, rio))

	// Output:
	// host <nil>
}

// ---------------------------------------------------------------------------
// ChainResultIK / ChainEitherIK
// ---------------------------------------------------------------------------

// ExampleChainResultIK demonstrates chaining a plain `func(A) (B, error)`, the
// most common shape in existing Go code.
func ExampleChainResultIK() {
	rio := F.Pipe2(
		Of[exConfig]("42"),
		ChainResultIK[exConfig](exParse),
		Map[exConfig](N.Mul(2)),
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 84 <nil>
}

// ExampleChainResultIK_failure demonstrates that the returned error
// short-circuits the remainder of the pipeline.
func ExampleChainResultIK_failure() {
	rio := F.Pipe2(
		Of[exConfig]("not a number"),
		ChainResultIK[exConfig](exParse),
		Map[exConfig](N.Mul(2)), // never reached
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 0 strconv.Atoi: parsing "not a number": invalid syntax
}

// ExampleMonadChainResultIK demonstrates the uncurried form of ChainResultIK.
func ExampleMonadChainResultIK() {
	rio := MonadChainResultIK(Of[exConfig]("42"), exParse)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 42 <nil>
}

// ExampleChainFirstResultIK demonstrates a validation step that keeps the
// original value on success.
func ExampleChainFirstResultIK() {
	validate := func(s string) (int, error) { return strconv.Atoi(s) }

	rio := F.Pipe1(
		Of[exConfig]("42"),
		ChainFirstResultIK[exConfig](validate),
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ChainIOResultIK
// ---------------------------------------------------------------------------

// ExampleChainIOResultIK demonstrates chaining a deferred
// `func(A) func() (B, error)`, the shape of a lazily executed side effect.
func ExampleChainIOResultIK() {
	rio := F.Pipe1(
		Of[exConfig](7),
		ChainIOResultIK[exConfig](exLoad),
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// record-7 <nil>
}

// ExampleTapIOResultIK demonstrates an audit step: the side effect runs, its
// result is discarded and the original value flows on.
func ExampleTapIOResultIK() {
	rio := F.Pipe2(
		Of[exConfig](21),
		TapIOResultIK[exConfig](exAudit),
		Map[exConfig](N.Mul(2)),
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// audited: 21
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ChainReaderResultIK
// ---------------------------------------------------------------------------

// ExampleChainReaderResultIK demonstrates chaining a `func(A) func(R) (B, error)`,
// i.e. a computation that reads the environment and may fail.
func ExampleChainReaderResultIK() {
	rio := F.Pipe1(
		Of[exConfig](5),
		ChainReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, rio))

	// Output:
	// 15 <nil>
}

// ExampleChainReaderResultIK_failure demonstrates the error path, here caused by
// a missing configuration value.
func ExampleChainReaderResultIK_failure() {
	rio := F.Pipe1(
		Of[exConfig](5),
		ChainReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{}, rio))

	// Output:
	// 0 multiplier is not configured
}

// ExampleMonadChainReaderResultIK demonstrates the uncurried form of
// ChainReaderResultIK.
func ExampleMonadChainReaderResultIK() {
	rio := MonadChainReaderResultIK(Of[exConfig](5), exScaleByConfig)

	fmt.Println(exRun(exConfig{Multiplier: 3}, rio))

	// Output:
	// 15 <nil>
}

// ExampleTapReaderResultIK demonstrates a check against the environment that
// preserves the original value.
func ExampleTapReaderResultIK() {
	rio := F.Pipe1(
		Of[exConfig](5),
		TapReaderResultIK(exScaleByConfig),
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, rio))

	// Output:
	// 5 <nil>
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

// Example_idiomaticPipeline mixes all four levels - Result, ReaderResult,
// IOResult and ReaderIOResult - in one pipeline built entirely from plain Go
// functions.
func Example_idiomaticPipeline() {
	rio := F.Pipe3(
		Of[exConfig]("7"),
		ChainResultIK[exConfig](exParse),     // "7" -> 7
		ChainReaderResultIK(exScaleByConfig), // 7 -> 7 * 3 = 21
		ChainIOResultIK[exConfig](exLoad),    // 21 -> "record-21"
	)

	fmt.Println(exRun(exConfig{Multiplier: 3}, rio))

	// Output:
	// record-21 <nil>
}
