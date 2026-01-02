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

// Package generic provides generic retry combinators that work with any monadic type.
//
// This package implements retry logic for monads that don't raise exceptions but
// signal failure in their type, such as Option, Either, and EitherT. The retry
// logic is parameterized over the monad operations, making it highly composable
// and reusable across different effect types.
//
// # Key Concepts
//
// The retry combinator takes:
//   - A retry policy that determines when and how long to wait between retries
//   - An action that produces a monadic value
//   - A check function that determines if the result should be retried
//
// # Usage with Different Monads
//
// The generic retry function can be used with any monad by providing the
// appropriate monad operations (Chain, Of, and Delay). This allows the same
// retry logic to work with IO, IOEither, ReaderIO, and other monadic types.
//
// Example conceptual usage:
//
//	// For IOEither[E, A]
//	result := Retrying(
//		IOE.Chain[E, A],
//		IOE.Chain[E, R.RetryStatus],
//		IOE.Of[E, A],
//		IOE.Of[E, R.RetryStatus],
//		IOE.Delay[E, R.RetryStatus],
//		policy,
//		action,
//		shouldRetry,
//	)
package generic

import (
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/retry"
	"github.com/IBM/fp-go/v2/tailrec"
)

// applyAndDelay applies a retry policy to the current status and delays by the
// returned amount if the policy returns a delay. This is an internal helper
// function used by the Retrying combinator.
//
// The function:
//  1. Applies the policy to get the next status
//  2. If the policy returned a delay, waits for that duration
//  3. Returns the updated status wrapped in the monad
//
// Type parameters:
//   - HKTSTATUS: The higher-kinded type representing the monadic status (e.g., IO[RetryStatus])
//
// Parameters:
//   - monadOf: Lifts a RetryStatus into the monad
//   - monadDelay: Delays execution by a duration within the monad
//
// Returns:
//   - A function that takes a policy and status and returns the delayed status in the monad
func applyAndDelay[HKTSTATUS any](
	monadOf func(R.RetryStatus) HKTSTATUS,
	monadDelay func(time.Duration) func(HKTSTATUS) HKTSTATUS,

	policy R.RetryPolicy,
) func(status R.RetryStatus) HKTSTATUS {

	apStatus := F.Bind1st(R.ApplyPolicy, policy)

	return func(status R.RetryStatus) HKTSTATUS {
		newStatus := apStatus(status)
		ofNewStatus := monadOf(newStatus)
		return F.Pipe2(
			newStatus,
			R.PreviousDelayLens.Get,
			O.Fold(
				F.Constant(ofNewStatus),
				func(delay time.Duration) HKTSTATUS {
					return monadDelay(delay)(ofNewStatus)
				},
			),
		)
	}
}

// Retrying implements a generic retry combinator that works with any monadic type.
// It repeatedly executes an action until either the check function returns false
// (indicating success) or the retry policy terminates (indicating failure).
//
// This function is the core retry implementation that can be specialized for different
// monadic types (IO, IOEither, ReaderIO, etc.) by providing the appropriate monad
// operations. It uses tail recursion via trampolining to avoid stack overflow on
// deep retry chains.
//
// # How It Works
//
// The retry logic follows these steps:
//  1. Execute the action with the current retry status
//  2. Check the result using the check function
//  3. If check returns false, return the result (success case)
//  4. If check returns true, apply the retry policy to get the next delay
//  5. If the policy returns None, stop retrying and return the last result
//  6. If the policy returns Some(delay), wait for that duration and retry (step 1)
//
// # Type Parameters
//
//   - HKTTRAMPOLINE: The higher-kinded type for the trampoline monad (e.g., IO[Trampoline[RetryStatus, A]])
//   - HKTA: The higher-kinded type for the action result (e.g., IO[A])
//   - HKTSTATUS: The higher-kinded type for the status monad (e.g., IO[RetryStatus])
//   - A: The result type of the action
//
// # Parameters
//
// Monad operations for the result type:
//   - monadChain: Chains computations in the result monad (flatMap/bind operation)
//   - monadMapStatus: Maps over the status monad to produce a trampoline
//   - monadOf: Lifts a trampoline value into the result monad
//   - monadOfStatus: Lifts a RetryStatus into the status monad
//   - monadDelay: Delays execution by a duration within the status monad
//
// Tail recursion support:
//   - tailRec: Executes a tail-recursive function using trampolining to avoid stack overflow
//
// Retry configuration:
//   - policy: The retry policy that determines delays and when to stop retrying
//   - action: The action to retry, which receives the current RetryStatus
//   - check: A predicate that returns true if the result should be retried, false otherwise
//
// # Returns
//
// The result of the action wrapped in the monad HKTA. This will be either:
//   - The first result where check returns false (success)
//   - The last result when the policy terminates (exhausted retries)
//
// # Example Usage Pattern
//
// For IOEither[E, A]:
//
//	result := Retrying(
//		IOE.Chain[E, A],                    // Chain for IOEither[E, A]
//		IOE.Chain[E, R.RetryStatus],        // Chain for IOEither[E, RetryStatus]
//		IOE.Of[E, A],                       // Lift trampoline to IOEither
//		IOE.Of[E, R.RetryStatus],           // Lift status to IOEither
//		IOE.Delay[E, R.RetryStatus],        // Delay in IOEither
//		IOE.TailRec[R.RetryStatus, A],      // Tail recursion for IOEither
//		policy,                              // Retry policy
//		action,                              // Action to retry
//		E.IsLeft[A],                         // Retry on Left (error)
//	)
//
// # Implementation Notes
//
// The function uses trampolining to implement tail recursion, which prevents stack
// overflow when retrying many times. The trampoline can be in one of two states:
//   - Land: Indicates completion with a final result
//   - Bounce: Indicates another iteration is needed with an updated status
//
// The retry logic checks if the policy returned a delay (Some) or termination (None).
// If a delay is present, it bounces to the next iteration. If None, it lands with
// the current result.
func Retrying[HKTTRAMPOLINE, HKTA, HKTSTATUS, A any](
	monadChain func(func(A) HKTTRAMPOLINE) func(HKTA) HKTTRAMPOLINE,
	monadMapStatus func(func(R.RetryStatus) tailrec.Trampoline[R.RetryStatus, A]) func(HKTSTATUS) HKTTRAMPOLINE,
	monadOf func(tailrec.Trampoline[R.RetryStatus, A]) HKTTRAMPOLINE,
	monadOfStatus func(R.RetryStatus) HKTSTATUS,
	monadDelay func(time.Duration) func(HKTSTATUS) HKTSTATUS,

	tailRec func(func(R.RetryStatus) HKTTRAMPOLINE) func(R.RetryStatus) HKTA,

	policy R.RetryPolicy,
	action func(R.RetryStatus) HKTA,
	check func(A) bool,
) HKTA {
	// delay callback
	applyDelay := applyAndDelay(monadOfStatus, monadDelay, policy)

	// function to check if we need to retry or not
	checkForRetry := O.FromPredicate(check)

	// need some lazy init because we reference it in the chain
	retryFct := func(status R.RetryStatus) HKTTRAMPOLINE {
		return F.Pipe2(
			status,
			action,
			monadChain(func(a A) HKTTRAMPOLINE {
				return F.Pipe3(
					a,
					checkForRetry,
					O.Map(func(a A) HKTTRAMPOLINE {
						return F.Pipe1(
							applyDelay(status),
							monadMapStatus(func(status R.RetryStatus) tailrec.Trampoline[R.RetryStatus, A] {
								return F.Pipe2(
									status,
									R.PreviousDelayLens.Get,
									O.Fold(
										F.Constant(tailrec.Land[R.RetryStatus](a)),
										F.Constant1[time.Duration](tailrec.Bounce[A](status)),
									),
								)
							}),
						)
					}),
					O.GetOrElse(F.Constant(monadOf(tailrec.Land[R.RetryStatus](a)))),
				)
			}),
		)
	}

	// seed
	return tailRec(retryFct)(R.DefaultRetryStatus)
}
