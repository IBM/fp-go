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
	"context"
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
)

// exkPositive fails for numbers that are not positive.
func exkPositive(n int) result.Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

// exkParse parses a decimal string as an IOResult.
var exkParse = ioresult.Eitherize1(strconv.Atoi)

// exkRun runs a ReaderIOResult and prints value and error.
func exkRun[A any](rr ReaderIOResult[A]) {
	fmt.Println(result.Unwrap(rr(context.Background())()))
}

// ExampleChainIOResultK chains a side effect that may fail.
func ExampleChainIOResultK() {
	exkRun(F.Pipe1(Of("42"), ChainIOResultK(exkParse)))
	// Output: 42 <nil>
}

// ExampleMonadChainIOResultK is the non curried version of ChainIOResultK.
func ExampleMonadChainIOResultK() {
	exkRun(MonadChainIOResultK(Of("7"), exkParse))
	// Output: 7 <nil>
}

// ExampleChainFirstIOResultK runs a fallible side effect and keeps the original value.
func ExampleChainFirstIOResultK() {
	exkRun(F.Pipe1(Of("42"), ChainFirstIOResultK(exkParse)))
	// Output: 42 <nil>
}

// ExampleTapIOResultK is an alias for ChainFirstIOResultK; a failing side
// effect replaces the value.
func ExampleTapIOResultK() {
	_, err := result.Unwrap(F.Pipe1(Of("x"), TapIOResultK(exkParse))(context.Background())())
	fmt.Println(err != nil)
	// Output: true
}

// ExampleMonadChainResultK is the non curried version of ChainResultK.
func ExampleMonadChainResultK() {
	exkRun(MonadChainResultK(Of(-1), exkPositive))
	// Output: 0 not positive
}

// ExampleChainFirstResultK validates a value without changing it.
func ExampleChainFirstResultK() {
	exkRun(F.Pipe1(Of(3), ChainFirstResultK(exkPositive)))
	exkRun(F.Pipe1(Of(-3), ChainFirstResultK(exkPositive)))
	// Output:
	// 3 <nil>
	// 0 not positive
}

// ExampleMonadTapResultK is the non curried version of TapResultK.
func ExampleMonadTapResultK() {
	exkRun(MonadTapResultK(Of(3), exkPositive))
	// Output: 3 <nil>
}

// ExampleMonadChainFirstResultK demonstrates MonadChainFirstResultK.
func ExampleMonadChainFirstResultK() {
	exkRun(MonadChainFirstResultK(Of(3), exkPositive))
	// Output: 3 <nil>
}

// ExampleTapResultK demonstrates TapResultK.
func ExampleTapResultK() {
	exkRun(F.Pipe1(Of(-3), TapResultK(exkPositive)))
	// Output: 0 not positive
}

// ExampleMonadChainIOEitherK demonstrates MonadChainIOEitherK.
func ExampleMonadChainIOEitherK() {
	exkRun(MonadChainIOEitherK(Of("7"), exkParse))
	// Output: 7 <nil>
}

// ExampleMonadChainFirstIOEitherK demonstrates MonadChainFirstIOEitherK.
func ExampleMonadChainFirstIOEitherK() {
	exkRun(MonadChainFirstIOEitherK(Of("7"), exkParse))
	// Output: 7 <nil>
}

// ExampleMonadTapIOEitherK demonstrates MonadTapIOEitherK.
func ExampleMonadTapIOEitherK() {
	exkRun(MonadTapIOEitherK(Of("7"), exkParse))
	// Output: 7 <nil>
}

// ExampleChainFirstIOEitherK demonstrates ChainFirstIOEitherK.
func ExampleChainFirstIOEitherK() {
	exkRun(F.Pipe1(Of("7"), ChainFirstIOEitherK(exkParse)))
	// Output: 7 <nil>
}

// ExampleTapIOEitherK demonstrates TapIOEitherK.
func ExampleTapIOEitherK() {
	exkRun(F.Pipe1(Of("7"), TapIOEitherK(exkParse)))
	// Output: 7 <nil>
}

// ExampleMonadChainFirstIOResultK demonstrates MonadChainFirstIOResultK.
func ExampleMonadChainFirstIOResultK() {
	exkRun(MonadChainFirstIOResultK(Of("7"), exkParse))
	// Output: 7 <nil>
}

// ExampleMonadTapIOResultK demonstrates MonadTapIOResultK.
func ExampleMonadTapIOResultK() {
	exkRun(MonadTapIOResultK(Of("7"), exkParse))
	// Output: 7 <nil>
}
