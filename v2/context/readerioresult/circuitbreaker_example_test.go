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
	"time"

	"github.com/IBM/fp-go/v2/circuitbreaker"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/retry"
	S "github.com/IBM/fp-go/v2/string"
)

// exampleRender turns a Result into a single readable line.
var exampleRender = result.Fold(
	S.Format[error]("blocked or failed: %s"),
	function.Identity[string],
)

// examplePolicy opens the circuit for one second, doubling the delay on every failed
// canary request, for at most five attempts.
func examplePolicy() retry.RetryPolicy {
	return retry.Monoid.Concat(
		retry.LimitRetries(5),
		retry.ExponentialBackoff(time.Second),
	)
}

// ExampleMakeSingletonBreaker protects a ReaderIOResult with a circuit breaker that keeps
// its state internally. The example uses a controllable clock so that its output is
// reproducible; in production you would pass time.Now.
func ExampleMakeSingletonBreaker() {
	// a clock the example can move forward at will
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }

	// the circuit opens on the second consecutive failure
	protect := MakeSingletonBreaker[string](
		clock,
		circuitbreaker.MakeClosedStateCounter(2),
		circuitbreaker.AnyError,
		examplePolicy(),
		circuitbreaker.MakeVoidMetrics(),
	)

	unhealthy := Left[string](errors.New("service unavailable"))
	healthy := Of("pong")

	ctx := context.Background()
	call := func(label string, op ReaderIOResult[string]) {
		fmt.Printf("%-18s %s\n", label+":", exampleRender(protect(op)(ctx)()))
	}

	call("first failure", unhealthy)
	call("second failure", unhealthy)

	// the circuit is open now, the service is not called at all
	call("while open", healthy)

	// once the backoff has elapsed, the next request becomes the canary
	now = now.Add(2 * time.Second)
	call("canary", healthy)

	// the canary succeeded, so the circuit is closed again
	call("after recovery", healthy)

	// Output:
	// first failure:     blocked or failed: service unavailable
	// second failure:    blocked or failed: service unavailable
	// while open:        blocked or failed: circuit breaker is open [Generic Circuit Breaker], will close at 2024-01-01 12:00:01 +0000 UTC
	// canary:            pong
	// after recovery:    pong
}

// ExampleMakeCircuitBreaker uses the breaker with an explicit state reference. This is the
// form to pick when several call sites have to share one breaker, or when the state has to
// be inspected from the outside.
func ExampleMakeCircuitBreaker() {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }

	closedState := circuitbreaker.MakeClosedStateCounter(2)

	cb := MakeCircuitBreaker[string](
		clock,
		closedState,
		circuitbreaker.AnyError,
		examplePolicy(),
		circuitbreaker.MakeVoidMetrics(),
	)

	// the breaker state lives in an IORef that the caller owns
	stateRef := circuitbreaker.MakeClosedIORef(closedState)()

	ctx := context.Background()
	call := func(op ReaderIOResult[string]) string {
		return exampleRender(pair.Tail(cb(pair.MakePair(stateRef, op)))(ctx)())
	}

	unhealthy := Left[string](errors.New("service unavailable"))

	fmt.Println("open before:", circuitbreaker.IsOpen(ioref.Read(stateRef)()))
	call(unhealthy)
	call(unhealthy)
	fmt.Println("open after: ", circuitbreaker.IsOpen(ioref.Read(stateRef)()))

	// Output:
	// open before: false
	// open after:  true
}

// ExampleMakeSingletonBreaker_errorFilter shows how checkError keeps application level
// failures from opening the circuit. Only errors the filter maps to Some are counted, so
// a validation error is treated like a success by the breaker while it is still reported
// to the caller.
func ExampleMakeSingletonBreaker_errorFilter() {
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	clock := func() time.Time { return now }

	protect := MakeSingletonBreaker[string](
		clock,
		circuitbreaker.MakeClosedStateCounter(2),
		// only infrastructure failures count towards the threshold
		circuitbreaker.InfrastructureError,
		examplePolicy(),
		circuitbreaker.MakeVoidMetrics(),
	)

	invalid := Left[string](errors.New("invalid request payload"))

	ctx := context.Background()
	for i := range 3 {
		fmt.Printf("attempt %d: %s\n", i+1, exampleRender(protect(invalid)(ctx)()))
	}

	// Output:
	// attempt 1: blocked or failed: invalid request payload
	// attempt 2: blocked or failed: invalid request payload
	// attempt 3: blocked or failed: invalid request payload
}
