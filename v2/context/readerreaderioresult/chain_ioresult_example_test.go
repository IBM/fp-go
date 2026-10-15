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

package readerreaderioresult

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
)

// rxPositiveIO fails for numbers that are not positive, as an IOResult.
var rxPositiveIO = F.Flow2(rxPositive, ioresult.FromResult[int])

// ExampleChainIOResultK chains a side effect that may fail.
func ExampleChainIOResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainIOResultK[exConfig](rxPositiveIO))))
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](-1), ChainIOResultK[exConfig](rxPositiveIO))))
	// Output:
	// 7 <nil>
	// 0 not positive
}

// ExampleMonadChainIOResultK is the non curried version of ChainIOResultK.
func ExampleMonadChainIOResultK() {
	fmt.Println(exRun(rxCfg, MonadChainIOResultK(Of[exConfig](7), rxPositiveIO)))
	// Output: 7 <nil>
}

// ExampleChainFirstIOResultK runs a side effect that may fail and keeps the original value.
func ExampleChainFirstIOResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainFirstIOResultK[exConfig](F.Flow2(rxPrintInt, ioresult.FromIO[int])))))
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](-1), ChainFirstIOResultK[exConfig](rxPositiveIO))))
	// Output:
	// log: 7
	// 7 <nil>
	// 0 not positive
}

// ExampleMonadChainFirstIOResultK is the non curried version of ChainFirstIOResultK.
func ExampleMonadChainFirstIOResultK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstIOResultK(Of[exConfig](7), rxPositiveIO)))
	// Output: 7 <nil>
}

// ExampleTapIOResultK is an alias for ChainFirstIOResultK; the side effect is skipped on failure.
func ExampleTapIOResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Left[exConfig, int](rxErr), TapIOResultK[exConfig](rxPositiveIO))))
	// Output: 0 boom
}

// ExampleMonadTapIOResultK is the non curried version of TapIOResultK.
func ExampleMonadTapIOResultK() {
	fmt.Println(exRun(rxCfg, MonadTapIOResultK(Of[exConfig](7), F.Flow2(rxPrintInt, ioresult.FromIO[int]))))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleChainFirstIOEitherK runs an IOEither side effect and keeps the original value.
func ExampleChainFirstIOEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainFirstIOEitherK[exConfig](rxPositiveIO))))
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](-1), ChainFirstIOEitherK[exConfig](rxPositiveIO))))
	// Output:
	// 7 <nil>
	// 0 not positive
}

// ExampleMonadChainFirstIOEitherK is the non curried version of ChainFirstIOEitherK.
func ExampleMonadChainFirstIOEitherK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstIOEitherK(Of[exConfig](7), rxPositiveIO)))
	// Output: 7 <nil>
}

// ExampleTapIOEitherK is an alias for ChainFirstIOEitherK.
func ExampleTapIOEitherK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), TapIOEitherK[exConfig](F.Flow2(rxPrintInt, ioresult.FromIO[int])))))
	// Output:
	// log: 7
	// 7 <nil>
}

// ExampleMonadTapIOEitherK is the non curried version of TapIOEitherK.
func ExampleMonadTapIOEitherK() {
	fmt.Println(exRun(rxCfg, MonadTapIOEitherK(Left[exConfig, int](rxErr), rxPositiveIO)))
	// Output: 0 boom
}

// ExampleMonadChainResultK is the non curried version of ChainResultK.
func ExampleMonadChainResultK() {
	fmt.Println(exRun(rxCfg, MonadChainResultK(Of[exConfig](-1), rxPositive)))
	// Output: 0 not positive
}

// ExampleChainFirstResultK validates a value without changing it.
func ExampleChainFirstResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](7), ChainFirstResultK[exConfig](rxPositive))))
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](-7), ChainFirstResultK[exConfig](rxPositive))))
	// Output:
	// 7 <nil>
	// 0 not positive
}

// ExampleMonadChainFirstResultK is the non curried version of ChainFirstResultK.
func ExampleMonadChainFirstResultK() {
	fmt.Println(exRun(rxCfg, MonadChainFirstResultK(Of[exConfig](7), rxPositive)))
	// Output: 7 <nil>
}

// ExampleTapResultK is an alias for ChainFirstResultK.
func ExampleTapResultK() {
	fmt.Println(exRun(rxCfg, F.Pipe1(Of[exConfig](0), TapResultK[exConfig](rxPositive))))
	// Output: 0 not positive
}

// ExampleMonadTapResultK is the non curried version of TapResultK.
func ExampleMonadTapResultK() {
	fmt.Println(exRun(rxCfg, MonadTapResultK(Left[exConfig, int](rxErr), rxPositive)))
	// Output: 0 boom
}
