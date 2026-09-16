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
// Variants
//
// Examples for the uncurried (Monad*) and value-preserving (ChainFirst*, Tap*)
// forms of the layer-specific idiomatic operators. They reuse exCoreConfig,
// exCoreRun and the exIK* fixtures. The side-effect fixtures below print, so
// each example's output shows both that the effect ran and which value
// survived.
// ---------------------------------------------------------------------------

// exVarValidate parses s and reports success - the shape of a validation step.
func exVarValidate(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err == nil {
		fmt.Println("valid:", n)
	}
	return n, err
}

// exVarCheckBudget reads the dependency and reports whether n is affordable.
func exVarCheckBudget(n int) func(exCoreConfig) (int, error) {
	return func(cfg exCoreConfig) (int, error) {
		if n > cfg.Multiplier*10 {
			return 0, fmt.Errorf("%d exceeds budget", n)
		}
		fmt.Println("within budget:", n)
		return n, nil
	}
}

// exVarFetchLogged sees the context, performs IO and reports what it fetched.
func exVarFetchLogged(path string) func(context.Context) func() (string, error) {
	return func(ctx context.Context) func() (string, error) {
		return func() (string, error) {
			if err := ctx.Err(); err != nil {
				return "", err
			}
			fmt.Println("fetched:", path)
			return "fetched:" + path, nil
		}
	}
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

// ExampleMonadChainFirstResultIK demonstrates a validation step in uncurried form.
func ExampleMonadChainFirstResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainFirstResultIK(Of[exCoreConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadTapResultIK demonstrates the Tap alias of MonadChainFirstResultIK.
func ExampleMonadTapResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadTapResultIK(Of[exCoreConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleChainFirstResultIK demonstrates that a failed validation fails the
// effect, while a successful one leaves the value untouched.
func ExampleChainFirstResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("42"), ChainFirstResultIK[exCoreConfig](exVarValidate))))
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("nope"), ChainFirstResultIK[exCoreConfig](exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
	//  strconv.Atoi: parsing "nope": invalid syntax
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainReaderResultIK demonstrates the uncurried form of
// ChainReaderResultIK.
func ExampleMonadChainReaderResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, MonadChainReaderResultIK(Of[exCoreConfig](5), exIKScale)))

	// Output:
	// 15 <nil>
}

// ExampleMonadChainFirstReaderResultIK demonstrates a dependency-based check in
// uncurried form that keeps the original value.
func ExampleMonadChainFirstReaderResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, MonadChainFirstReaderResultIK(Of[exCoreConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleMonadTapReaderResultIK demonstrates the Tap alias of
// MonadChainFirstReaderResultIK.
func ExampleMonadTapReaderResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, MonadTapReaderResultIK(Of[exCoreConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleChainFirstReaderResultIK demonstrates that a failed dependency-based
// check fails the effect.
func ExampleChainFirstReaderResultIK() {
	cfg := exCoreConfig{Multiplier: 3}

	fmt.Println(exCoreRun(cfg, F.Pipe1(Of[exCoreConfig](25), ChainFirstReaderResultIK(exVarCheckBudget))))
	fmt.Println(exCoreRun(cfg, F.Pipe1(Of[exCoreConfig](40), ChainFirstReaderResultIK(exVarCheckBudget))))

	// Output:
	// within budget: 25
	// 25 <nil>
	// 0 40 exceeds budget
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainIOResultIK demonstrates the uncurried form of ChainIOResultIK.
func ExampleMonadChainIOResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainIOResultIK(Of[exCoreConfig](7), exIKLoad)))

	// Output:
	// record-7 <nil>
}

// ExampleChainFirstIOResultIK demonstrates running a deferred side effect while
// keeping the original value.
func ExampleChainFirstIOResultIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig](21), ChainFirstIOResultIK[exCoreConfig](exIKAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ---------------------------------------------------------------------------
// Thunk level
// ---------------------------------------------------------------------------

// ExampleMonadChainThunkIK demonstrates the uncurried form of ChainThunkIK.
func ExampleMonadChainThunkIK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainThunkIK(Of[exCoreConfig]("users"), exIKFetch)))

	// Output:
	// fetched:users <nil>
}

// ExampleChainFirstThunkIK demonstrates running a context-aware side effect
// while keeping the original value.
func ExampleChainFirstThunkIK() {
	eff := F.Pipe1(
		Of[exCoreConfig]("users"),
		ChainFirstThunkIK[exCoreConfig, string](exVarFetchLogged),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))

	// Output:
	// fetched: users
	// users <nil>
}

// ExampleChainFirstLeftThunkIK demonstrates reporting an error through a
// context-aware handler while leaving the failure itself untouched, and that
// the handler is skipped on success.
func ExampleChainFirstLeftThunkIK() {
	report := func(err error) func(context.Context) func() (string, error) {
		return func(context.Context) func() (string, error) {
			return func() (string, error) {
				fmt.Println("reporting:", err)
				return "", nil
			}
		}
	}

	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(
		Fail[exCoreConfig, int](errors.New("boom")),
		ChainFirstLeftThunkIK[exCoreConfig, int](report),
	)))
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(
		Of[exCoreConfig](42),
		ChainFirstLeftThunkIK[exCoreConfig, int](report),
	)))

	// Output:
	// reporting: boom
	// 0 boom
	// 42 <nil>
}
