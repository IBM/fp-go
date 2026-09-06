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
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
)

// exRunConfig is the dependency used by the run examples.
type exRunConfig struct {
	Multiplier int
}

// exRunMultiplier reads the Multiplier field; naming the getter keeps the
// pipelines below point-free.
func exRunMultiplier(c exRunConfig) int { return c.Multiplier }

// ExampleProvide demonstrates eliminating the context dependency of an Effect,
// turning it into a Thunk that no longer mentions the config type.
func ExampleProvide() {
	eff := Of[exRunConfig](42)

	// Provide fixes the dependency; the result is a ReaderIOResult[int].
	thunk := Provide[int](exRunConfig{Multiplier: 3})(eff)

	// A Thunk still has to be run: apply a context.Context, then the IO.
	fmt.Println(thunk(context.Background())())
	// Output:
	// Right[int](42)
}

// ExampleProvide_dependencyIsVisibleToTheEffect demonstrates that the value
// passed to Provide is what Asks and Ask observe.
func ExampleProvide_dependencyIsVisibleToTheEffect() {
	// Asks reads a field out of the context, point-free.
	eff := F.Pipe1(
		Asks(exRunMultiplier),
		Map[exRunConfig](N.Mul(10)),
	)

	fmt.Println(RunSync(Provide[int](exRunConfig{Multiplier: 7})(eff))(context.Background()))
	// Output:
	// 70 <nil>
}

// ExampleRunSync demonstrates executing a Thunk and receiving the idiomatic Go
// (value, error) pair.
func ExampleRunSync() {
	thunk := Provide[int](exRunConfig{Multiplier: 3})(Of[exRunConfig](42))

	value, err := RunSync(thunk)(context.Background())

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 42
	// <nil>
}

// ExampleRunSync_failure demonstrates that a failed effect surfaces as a
// non-nil error and the zero value of A.
func ExampleRunSync_failure() {
	eff := Fail[exRunConfig, int](fmt.Errorf("boom"))

	value, err := RunSync(Provide[int](exRunConfig{})(eff))(context.Background())

	fmt.Println(value)
	fmt.Println(err)
	// Output:
	// 0
	// boom
}

// ExampleRunSync_completePipeline demonstrates the canonical end-to-end shape:
// build an effect, provide its dependency, then run it.
func ExampleRunSync_completePipeline() {
	// Point-free pipeline: read the multiplier, triple it, render it.
	eff := F.Pipe2(
		Asks(exRunMultiplier),
		Map[exRunConfig](N.Mul(3)),
		Map[exRunConfig](func(n int) string { return fmt.Sprintf("total=%d", n) }),
	)

	value, err := F.Pipe2(
		eff,
		Provide[string](exRunConfig{Multiplier: 5}),
		RunSync,
	)(context.Background())

	fmt.Println(value, err)
	// Output:
	// total=15 <nil>
}
