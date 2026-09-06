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
	"time"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// A KleisliI is the shape a plain Go service method already has once the
// receiver is bound: it takes the input, then a context.Context plus the
// dependency, and returns (value, error).  Nothing in this file wraps a value
// in a monad by hand — that is the point of the idiomatic bridge.
// ---------------------------------------------------------------------------

// exScaleI scales its input by the configured multiplier.
func exScaleI(n int) func(context.Context, exCoreConfig) (int, error) {
	return func(_ context.Context, cfg exCoreConfig) (int, error) {
		return n * cfg.Multiplier, nil
	}
}

// exRejectI always fails, the way a service method reports a problem.
func exRejectI(n int) func(context.Context, exCoreConfig) (int, error) {
	return func(context.Context, exCoreConfig) (int, error) {
		return 0, fmt.Errorf("cannot handle %d", n)
	}
}

// exAuditI performs a side effect and reports nothing interesting.
func exAuditI(n int) func(context.Context, exCoreConfig) (string, error) {
	return func(context.Context, exCoreConfig) (string, error) {
		fmt.Println("audited:", n)
		return "", nil
	}
}

// exRecoverI turns any error into a fallback value.
func exRecoverI(err error) func(context.Context, exCoreConfig) (int, error) {
	return func(context.Context, exCoreConfig) (int, error) {
		fmt.Println("recovering from:", err)
		return -1, nil
	}
}

// ---------------------------------------------------------------------------
// FromIdiomatic
// ---------------------------------------------------------------------------

// ExampleFromIdiomatic demonstrates lifting an ordinary Go method signature into
// a Kleisli arrow, so it composes with Chain, Bind and the rest.
func ExampleFromIdiomatic() {
	scale := FromIdiomatic(exScaleI)

	eff := F.Pipe1(Of[exCoreConfig](5), Chain(scale))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 3}, eff))
	// Output:
	// 15 <nil>
}

// ExampleFromIdiomatic_failure demonstrates that the returned error becomes a
// failed effect.
func ExampleFromIdiomatic_failure() {
	eff := F.Pipe2(
		Of[exCoreConfig](5),
		Chain(FromIdiomatic(exRejectI)),
		Map[exCoreConfig](N.Mul(100)), // never reached
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// 0 cannot handle 5
}

// ---------------------------------------------------------------------------
// ChainI / ChainFirstI / TapI
// ---------------------------------------------------------------------------

// ExampleChainI demonstrates chaining an idiomatic function directly, without
// the intermediate FromIdiomatic call.
func ExampleChainI() {
	eff := F.Pipe1(Of[exCoreConfig](5), ChainI(exScaleI))

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// 20 <nil>
}

// ExampleMonadChainI demonstrates the uncurried form of ChainI.
func ExampleMonadChainI() {
	eff := MonadChainI(Of[exCoreConfig](5), exScaleI)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// 20 <nil>
}

// ExampleChainFirstI demonstrates running an idiomatic function for its side
// effect while keeping the upstream value.
func ExampleChainFirstI() {
	eff := F.Pipe2(
		Of[exCoreConfig](21),
		ChainFirstI(exAuditI),
		Map[exCoreConfig](N.Mul(2)),
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 21
	// 42 <nil>
}

// ExampleMonadChainFirstI demonstrates the uncurried form of ChainFirstI.
func ExampleMonadChainFirstI() {
	eff := MonadChainFirstI(Of[exCoreConfig](21), exAuditI)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 21
	// 21 <nil>
}

// ExampleTapI demonstrates that TapI is ChainFirstI under an
// observation-flavoured name.
func ExampleTapI() {
	eff := F.Pipe1(Of[exCoreConfig](7), TapI(exAuditI))

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 7
	// 7 <nil>
}

// ExampleMonadTapI demonstrates the uncurried form of TapI.
func ExampleMonadTapI() {
	eff := MonadTapI(Of[exCoreConfig](7), exAuditI)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// audited: 7
	// 7 <nil>
}

// ---------------------------------------------------------------------------
// ChainLeftI
// ---------------------------------------------------------------------------

// ExampleChainLeftI demonstrates recovering from a failure with an idiomatic
// handler.  A successful effect passes through untouched.
func ExampleChainLeftI() {
	recovered := F.Pipe1(
		Fail[exCoreConfig, int](fmt.Errorf("primary failed")),
		ChainLeftI(exRecoverI),
	)
	fmt.Println(exCoreRun(exCoreConfig{}, recovered))

	untouched := F.Pipe1(Of[exCoreConfig](42), ChainLeftI(exRecoverI))
	fmt.Println(exCoreRun(exCoreConfig{}, untouched))
	// Output:
	// recovering from: primary failed
	// -1 <nil>
	// 42 <nil>
}

// ExampleMonadChainLeftI demonstrates the uncurried form of ChainLeftI.
func ExampleMonadChainLeftI() {
	eff := MonadChainLeftI(Fail[exCoreConfig, int](fmt.Errorf("primary failed")), exRecoverI)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	// Output:
	// recovering from: primary failed
	// -1 <nil>
}

// ---------------------------------------------------------------------------
// TraverseArrayI
// ---------------------------------------------------------------------------

// ExampleTraverseArrayI demonstrates running an idiomatic function once per
// element and collecting the results.
func ExampleTraverseArrayI() {
	eff := TraverseArrayI(exScaleI)([]int{1, 2, 3})

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 10}, eff))
	// Output:
	// [10 20 30] <nil>
}

// ExampleTraverseArrayI_failFast demonstrates that the first failing element
// fails the whole traversal.
func ExampleTraverseArrayI_failFast() {
	guard := func(n int) func(context.Context, exCoreConfig) (int, error) {
		if n < 0 {
			return exRejectI(n)
		}
		return exScaleI(n)
	}

	eff := TraverseArrayI(guard)([]int{1, -2, 3})

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 10}, eff))
	// Output:
	// [] cannot handle -2
}

// ---------------------------------------------------------------------------
// BindI / BindIL
// ---------------------------------------------------------------------------

// ExampleBindI demonstrates using an idiomatic function inside a do-notation
// block.  The function receives the accumulated state.
func ExampleBindI() {
	// exOrder -> (int, error): price the order.
	price := func(o exOrder) func(context.Context, exCoreConfig) (int, error) {
		return func(_ context.Context, cfg exCoreConfig) (int, error) {
			return o.Quantity * cfg.Multiplier, nil
		}
	}

	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApS(exSetQuantity, Of[exCoreConfig](3)),
		BindI(exSetTotal, price),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 25}, eff))
	// Output:
	// { 3 75} <nil>
}

// ExampleBindIL demonstrates the lens form: the idiomatic function receives the
// current field value and refines it in place.
func ExampleBindIL() {
	eff := F.Pipe2(
		Do[exCoreConfig](exOrder{}),
		ApSL(exQuantityLens, Of[exCoreConfig](3)),
		BindIL(exQuantityLens, exScaleI),
	)

	fmt.Println(exCoreRun(exCoreConfig{Multiplier: 4}, eff))
	// Output:
	// { 12 0} <nil>
}

// ---------------------------------------------------------------------------
// Retrying / RetryingI
// ---------------------------------------------------------------------------

// ExampleRetrying demonstrates retrying a flaky effect until it succeeds or the
// policy is exhausted.  The check predicate decides what counts as "try again".
func ExampleRetrying() {
	attempts := 0

	// The action receives the retry status, so it can report or adapt.
	action := func(status retry.RetryStatus) Effect[exCoreConfig, string] {
		return FromIO[exCoreConfig](func() string {
			attempts++
			return fmt.Sprintf("attempt %d", status.IterNumber)
		})
	}

	// Retry while the result is still an error; succeed as soon as it is not.
	eff := Retrying(
		retry.LimitRetries(3),
		action,
		R.IsLeft[string],
	)

	value, err := exCoreRun(exCoreConfig{}, eff)
	fmt.Println(value, err)
	fmt.Println("attempts:", attempts)
	// Output:
	// attempt 0 <nil>
	// attempts: 1
}

// ExampleRetrying_untilPolicyExhausted demonstrates an action that keeps
// failing: the policy caps the number of attempts and the last error surfaces.
func ExampleRetrying_untilPolicyExhausted() {
	attempts := 0

	action := func(retry.RetryStatus) Effect[exCoreConfig, string] {
		return Suspend(func() Effect[exCoreConfig, string] {
			attempts++
			return Fail[exCoreConfig, string](fmt.Errorf("still failing"))
		})
	}

	eff := Retrying(
		retry.CapDelay(time.Millisecond, retry.LimitRetries(2)),
		action,
		R.IsLeft[string],
	)

	value, err := exCoreRun(exCoreConfig{}, eff)
	fmt.Printf("%q %v\n", value, err)
	fmt.Println("attempts:", attempts)
	// Output:
	// "" still failing
	// attempts: 3
}

// ExampleRetryingI demonstrates the idiomatic spelling: the retried action is an
// ordinary Go function rather than an Effect.
func ExampleRetryingI() {
	attempts := 0

	action := func(status retry.RetryStatus) func(context.Context, exCoreConfig) (string, error) {
		return func(context.Context, exCoreConfig) (string, error) {
			attempts++
			if attempts < 3 {
				return "", fmt.Errorf("transient failure %d", attempts)
			}
			return fmt.Sprintf("succeeded on iteration %d", status.IterNumber), nil
		}
	}

	eff := RetryingI(
		retry.CapDelay(time.Millisecond, retry.LimitRetries(5)),
		action,
		R.IsLeft[string],
	)

	fmt.Println(exCoreRun(exCoreConfig{}, eff))
	fmt.Println("attempts:", attempts)
	// Output:
	// succeeded on iteration 2 <nil>
	// attempts: 3
}
