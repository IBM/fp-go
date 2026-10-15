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

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
)

// exPrintInt prints an int as a side effect that cannot fail.
var exPrintInt = io.Printf[int]("log: %d\n")

// ExampleMonadChainIOK is the non curried version of ChainIOK.
func ExampleMonadChainIOK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainIOK(Of[exCoreConfig](21), F.Flow2(N.Mul(2), io.Of[int]))))
	// Output:
	// 42 <nil>
}

// ExampleMonadChainFirstIOK is the non curried version of ChainFirstIOK.
func ExampleMonadChainFirstIOK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainFirstIOK(Of[exCoreConfig](7), exPrintInt)))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleMonadTapIOK is the non curried version of TapIOK; the side effect is
// skipped on failure.
func ExampleMonadTapIOK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadTapIOK(Fail[exCoreConfig, int](exCoreErr("boom")), exPrintInt)))
	// Output:
	// 0 boom
}

// ExampleMonadChainResultK is the non curried version of ChainResultK.
func ExampleMonadChainResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainResultK(Of[exCoreConfig](10), exCoreHalve)))
	// Output:
	// 5 <nil>
}

// ExampleChainFirstResultK demonstrates validating a value without changing it.
func ExampleChainFirstResultK() {
	validate := ChainFirstResultK[exCoreConfig](exCoreRequirePositive)

	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig](3), validate)))
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig](-3), validate)))
	// Output:
	// 3 <nil>
	// 0 not positive
}

// ExampleMonadChainFirstResultK is the non curried version of ChainFirstResultK.
func ExampleMonadChainFirstResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadChainFirstResultK(Of[exCoreConfig](10), exCoreHalve)))
	// Output:
	// 10 <nil>
}

// ExampleTapResultK is an alias for ChainFirstResultK.
func ExampleTapResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, F.Pipe1(Of[exCoreConfig](0), TapResultK[exCoreConfig](exCoreRequirePositive))))
	// Output:
	// 0 not positive
}

// ExampleMonadTapResultK is the non curried version of TapResultK.
func ExampleMonadTapResultK() {
	fmt.Println(exCoreRun(exCoreConfig{}, MonadTapResultK(Of[exCoreConfig](4), exCoreRequirePositive)))
	// Output:
	// 4 <nil>
}
