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
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	R "github.com/IBM/fp-go/v2/result"
)

// exParse parses a decimal string as an IOResult, e.g. a value read from a file.
var exParse = ioresult.Eitherize1(strconv.Atoi)

// exAudit prints a value as a side effect that may fail.
func exAudit(n int) IOResult[string] {
	return func() Result[string] {
		fmt.Println("audit:", n)
		return R.Of("ok")
	}
}

// ExampleFromIOResult demonstrates lifting an IOResult into an Effect. The
// IOResult is not executed until the effect is run.
func ExampleFromIOResult() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromIOResult[exCoreConfig](exParse("42"))))
	fmt.Println(exCoreRun(exCoreConfig{}, FromIOResult[exCoreConfig](ioresult.Left[int](exCoreErr("io failed")))))
	// Output:
	// 42 <nil>
	// 0 io failed
}

// ExampleChainIOResultK demonstrates chaining a side effect that may fail; its
// error short-circuits the pipeline.
func ExampleChainIOResultK() {
	parse := ChainIOResultK[exCoreConfig](exParse)

	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("21"), parse)))
	_, err := exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig]("abc"), parse))
	fmt.Println(err != nil)
	// Output:
	// 21 <nil>
	// true
}

// ExampleMonadChainIOResultK is the non curried version of ChainIOResultK.
func ExampleMonadChainIOResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainIOResultK(Of[exCoreConfig]("7"), exParse)))
	// Output:
	// 7 <nil>
}

// ExampleChainFirstIOResultK demonstrates running a fallible side effect while
// keeping the upstream value.
func ExampleChainFirstIOResultK() {
	eff := F.Pipe1(Of[exCoreConfig](42), ChainFirstIOResultK[exCoreConfig](exAudit))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audit: 42
	// 42 <nil>
}

// ExampleChainFirstIOResultK_failure demonstrates that a failing side effect
// replaces the upstream value with its error.
func ExampleChainFirstIOResultK_failure() {
	reject := func(int) IOResult[string] { return ioresult.Left[string](exCoreErr("rejected")) }
	eff := F.Pipe1(Of[exCoreConfig](42), ChainFirstIOResultK[exCoreConfig](reject))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 0 rejected
}

// ExampleMonadChainFirstIOResultK is the non curried version of ChainFirstIOResultK.
func ExampleMonadChainFirstIOResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainFirstIOResultK(Of[exCoreConfig](7), exAudit)))
	// Output:
	// audit: 7
	// 7 <nil>
}

// ExampleTapIOResultK demonstrates that TapIOResultK is ChainFirstIOResultK
// under an observation-flavoured name; it is skipped on failure.
func ExampleTapIOResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig](7), TapIOResultK[exCoreConfig](exAudit))))
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Fail[exCoreConfig, int](exCoreErr("boom")), TapIOResultK[exCoreConfig](exAudit))))
	// Output:
	// audit: 7
	// 7 <nil>
	// 0 boom
}

// ExampleMonadTapIOResultK is the non curried version of TapIOResultK.
func ExampleMonadTapIOResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadTapIOResultK(Of[exCoreConfig](7), exAudit)))
	// Output:
	// audit: 7
	// 7 <nil>
}
