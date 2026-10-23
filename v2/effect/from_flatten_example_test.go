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

	ctxthunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// exFromLabel turns the configured prefix into an IO action that returns it
// with an exclamation mark.
func exFromLabel(c exCoreConfig) IO[string] {
	return io.Of(c.Prefix + "!")
}

// ExampleFromReaderIO demonstrates lifting a ReaderIO, a context-dependent IO
// action that cannot fail, into an Effect.
func ExampleFromReaderIO() {
	fmt.Println(exCoreRun(exCoreConfig{Prefix: "hello"}, FromReaderIO(exFromLabel)))
	// Output:
	// hello! <nil>
}

// ExampleFromReaderIOResult demonstrates lifting a context-dependent IO action
// that may fail into an Effect.
func ExampleFromReaderIOResult() {
	parsePrefix := F.Flow2(exCorePrefix, ioresult.Eitherize1(strconv.Atoi))
	eff := FromReaderIOResult(parsePrefix)

	fmt.Println(exCoreRun(exCoreConfig{Prefix: "42"}, eff))
	_, err := exCoreRun(exCoreConfig{Prefix: "abc"}, eff)
	fmt.Println(err != nil)
	// Output:
	// 42 <nil>
	// true
}

// ExampleFromReaderIOResult_failure demonstrates that a failed ReaderIOResult
// becomes a failed Effect.
func ExampleFromReaderIOResult_failure() {
	fmt.Println(exCoreRun(exCoreConfig{}, FromReaderIOResult(RIOR.Left[exCoreConfig, int](exCoreErr("boom")))))
	// Output:
	// 0 boom
}

// ExampleFromReaderThunk demonstrates lifting a function from the context to a
// Thunk into an Effect.
func ExampleFromReaderThunk() {
	eff := FromReaderThunk(F.Flow2(exCoreMultiplier, ctxthunk.Of[int]))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 7}, eff))
	// Output:
	// 7 <nil>
}

// ExampleFlatten demonstrates removing one level of nesting from an Effect
// that produces another Effect.
func ExampleFlatten() {
	nested := Of[exCoreConfig](Asks(exCoreMultiplier))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 5}, Flatten(nested)))
	// Output:
	// 5 <nil>
}

// ExampleFlattenIO demonstrates running an IO action that was produced by a
// Map step in the middle of a pipeline.
func ExampleFlattenIO() {
	eff := F.Pipe2(
		Of[exCoreConfig](21),
		Map[exCoreConfig](F.Flow2(strconv.Itoa, io.Of[string])),
		FlattenIO[exCoreConfig, string],
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 21 <nil>
}

// ExampleFlattenIOResult demonstrates running an IOResult that was produced
// by a Map step; its failure fails the pipeline.
func ExampleFlattenIOResult() {
	parse := F.Flow2(
		Map[exCoreConfig](ioresult.Eitherize1(strconv.Atoi)),
		FlattenIOResult[exCoreConfig, int],
	)

	fmt.Println(exCoreRun(exCoreConfig{}, parse(Of[exCoreConfig]("42"))))
	_, err := exCoreRun(exCoreConfig{}, parse(Of[exCoreConfig]("abc")))
	fmt.Println(err != nil)
	// Output:
	// 42 <nil>
	// true
}

// ExampleFlattenThunk demonstrates running a Thunk that was produced in the
// middle of a pipeline, e.g. by reading a field of the context.
func ExampleFlattenThunk() {
	eff := F.Pipe2(
		Asks(exCoreMultiplier),
		Map[exCoreConfig](F.Flow2(strconv.Itoa, ctxthunk.Of[string])),
		FlattenThunk[exCoreConfig, string],
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))
	// Output:
	// 3 <nil>
}
