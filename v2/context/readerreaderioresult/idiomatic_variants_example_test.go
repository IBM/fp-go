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
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
)

// ---------------------------------------------------------------------------
// Variants
//
// Examples for the uncurried (Monad*), value-preserving (ChainFirst*, Tap*) and
// Either-named forms of the layer-specific idiomatic operators. They reuse the
// fixtures from idiomatic_example_test.go. The side-effect fixtures below
// print, so each example's output shows both that the effect ran and which
// value survived.
// ---------------------------------------------------------------------------

// exVarValidate parses s and reports success - the shape of a validation step.
func exVarValidate(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err == nil {
		fmt.Println("valid:", n)
	}
	return n, err
}

// exVarCheckBudget reads the config and reports whether n is affordable.
func exVarCheckBudget(n int) func(exConfig) (int, error) {
	return func(cfg exConfig) (int, error) {
		if n > cfg.Multiplier*10 {
			return 0, fmt.Errorf("%d exceeds budget", n)
		}
		fmt.Println("within budget:", n)
		return n, nil
	}
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

// ExampleMonadChainEitherIK demonstrates the Either-named alias of MonadChainResultIK.
func ExampleMonadChainEitherIK() {
	fmt.Println(exRun(exConfig{}, MonadChainEitherIK(Of[exConfig]("42"), exParse)))

	// Output:
	// 42 <nil>
}

// ExampleChainEitherIK demonstrates the Either-named alias of ChainResultIK.
func ExampleChainEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("42"), ChainEitherIK[exConfig](exParse))))

	// Output:
	// 42 <nil>
}

// ExampleMonadChainFirstEitherIK demonstrates a validation step in uncurried form.
func ExampleMonadChainFirstEitherIK() {
	fmt.Println(exRun(exConfig{}, MonadChainFirstEitherIK(Of[exConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadChainFirstResultIK demonstrates the Result-named alias of
// MonadChainFirstEitherIK.
func ExampleMonadChainFirstResultIK() {
	fmt.Println(exRun(exConfig{}, MonadChainFirstResultIK(Of[exConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadTapEitherIK demonstrates the Tap alias of MonadChainFirstEitherIK.
func ExampleMonadTapEitherIK() {
	fmt.Println(exRun(exConfig{}, MonadTapEitherIK(Of[exConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadTapResultIK demonstrates the Tap alias of MonadChainFirstResultIK.
func ExampleMonadTapResultIK() {
	fmt.Println(exRun(exConfig{}, MonadTapResultIK(Of[exConfig]("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleChainFirstEitherIK demonstrates a validation step that keeps the
// original value.
func ExampleChainFirstEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("42"), ChainFirstEitherIK[exConfig](exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleChainFirstResultIK demonstrates that a failed validation fails the
// computation.
func ExampleChainFirstResultIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("42"), ChainFirstResultIK[exConfig](exVarValidate))))
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("nope"), ChainFirstResultIK[exConfig](exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
	//  strconv.Atoi: parsing "nope": invalid syntax
}

// ExampleTapEitherIK demonstrates the Tap alias of ChainFirstEitherIK.
func ExampleTapEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig]("42"), TapEitherIK[exConfig](exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainReaderEitherIK demonstrates the Either-named, uncurried form
// of ChainReaderResultIK.
func ExampleMonadChainReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadChainReaderEitherIK(Of[exConfig](5), exScaleByConfig)))

	// Output:
	// 15 <nil>
}

// ExampleMonadChainReaderResultIK demonstrates the uncurried form of
// ChainReaderResultIK.
func ExampleMonadChainReaderResultIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadChainReaderResultIK(Of[exConfig](5), exScaleByConfig)))

	// Output:
	// 15 <nil>
}

// ExampleChainReaderEitherIK demonstrates the Either-named alias of
// ChainReaderResultIK.
func ExampleChainReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, F.Pipe1(Of[exConfig](5), ChainReaderEitherIK(exScaleByConfig))))

	// Output:
	// 15 <nil>
}

// ExampleMonadChainFirstReaderEitherIK demonstrates a config-dependent check in
// uncurried form that keeps the original value.
func ExampleMonadChainFirstReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadChainFirstReaderEitherIK(Of[exConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleMonadChainFirstReaderResultIK demonstrates the Result-named alias of
// MonadChainFirstReaderEitherIK.
func ExampleMonadChainFirstReaderResultIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadChainFirstReaderResultIK(Of[exConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleMonadTapReaderEitherIK demonstrates the Tap alias of
// MonadChainFirstReaderEitherIK.
func ExampleMonadTapReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadTapReaderEitherIK(Of[exConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleMonadTapReaderResultIK demonstrates the Tap alias of
// MonadChainFirstReaderResultIK.
func ExampleMonadTapReaderResultIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, MonadTapReaderResultIK(Of[exConfig](25), exVarCheckBudget)))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleChainFirstReaderEitherIK demonstrates a config-dependent check that
// keeps the original value.
func ExampleChainFirstReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, F.Pipe1(Of[exConfig](25), ChainFirstReaderEitherIK(exVarCheckBudget))))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ExampleChainFirstReaderResultIK demonstrates that a failed config-dependent
// check fails the computation.
func ExampleChainFirstReaderResultIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, F.Pipe1(Of[exConfig](25), ChainFirstReaderResultIK(exVarCheckBudget))))
	fmt.Println(exRun(exConfig{Multiplier: 3}, F.Pipe1(Of[exConfig](40), ChainFirstReaderResultIK(exVarCheckBudget))))

	// Output:
	// within budget: 25
	// 25 <nil>
	// 0 40 exceeds budget
}

// ExampleTapReaderEitherIK demonstrates the Tap alias of ChainFirstReaderEitherIK.
func ExampleTapReaderEitherIK() {
	fmt.Println(exRun(exConfig{Multiplier: 3}, F.Pipe1(Of[exConfig](25), TapReaderEitherIK(exVarCheckBudget))))

	// Output:
	// within budget: 25
	// 25 <nil>
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainIOEitherIK demonstrates the Either-named, uncurried form of
// ChainIOResultIK.
func ExampleMonadChainIOEitherIK() {
	fmt.Println(exRun(exConfig{}, MonadChainIOEitherIK(Of[exConfig](7), exLoad)))

	// Output:
	// record-7 <nil>
}

// ExampleMonadChainIOResultIK demonstrates the uncurried form of ChainIOResultIK.
func ExampleMonadChainIOResultIK() {
	fmt.Println(exRun(exConfig{}, MonadChainIOResultIK(Of[exConfig](7), exLoad)))

	// Output:
	// record-7 <nil>
}

// ExampleChainIOEitherIK demonstrates the Either-named alias of ChainIOResultIK,
// including its failure path.
func ExampleChainIOEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig](7), ChainIOEitherIK[exConfig](exLoad))))
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig](-1), ChainIOEitherIK[exConfig](exLoad))))

	// Output:
	// record-7 <nil>
	//  cannot load -1
}

// ExampleChainFirstIOEitherIK demonstrates running a deferred side effect while
// keeping the original value.
func ExampleChainFirstIOEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig](21), ChainFirstIOEitherIK[exConfig](exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ExampleChainFirstIOResultIK demonstrates the Result-named alias of
// ChainFirstIOEitherIK.
func ExampleChainFirstIOResultIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig](21), ChainFirstIOResultIK[exConfig](exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ExampleTapIOEitherIK demonstrates the Tap alias of ChainFirstIOEitherIK.
func ExampleTapIOEitherIK() {
	fmt.Println(exRun(exConfig{}, F.Pipe1(Of[exConfig](21), TapIOEitherIK[exConfig](exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}
