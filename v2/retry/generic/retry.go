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
) func(policy R.RetryPolicy, status R.RetryStatus) HKTSTATUS {
	return func(policy R.RetryPolicy, status R.RetryStatus) HKTSTATUS {
		newStatus := R.ApplyPolicy(policy, status)
		return F.Pipe1(
			newStatus.PreviousDelay,
			O.Fold(
				F.Nullary2(F.Constant(newStatus), monadOf),
				func(delay time.Duration) HKTSTATUS {
					return monadDelay(delay)(monadOf(newStatus))
				},
			),
		)
	}
}

// Retrying is a generic retry combinator for actions that don't raise exceptions,
// but signal failure in their type. This works with monads like Option, Either,
// and EitherT where the type itself indicates success or failure.
//
// The function repeatedly executes an action until either:
//  1. The action succeeds (check returns false)
//  2. The retry policy returns None (retry limit reached)
//  3. The action fails in a way that shouldn't be retried
//
// Type parameters:
//   - HKTA: The higher-kinded type for the action result (e.g., IO[A], Either[E, A])
//   - HKTSTATUS: The higher-kinded type for the retry status (e.g., IO[RetryStatus])
//   - A: The result type of the action
//
// Monad operations (first 5 parameters):
//   - monadChain: Chains operations on HKTA (flatMap/bind for the result type)
//   - monadChainStatus: Chains operations from HKTSTATUS to HKTA
//   - monadOf: Lifts a value A into HKTA (pure/return for the result type)
//   - monadOfStatus: Lifts a RetryStatus into HKTSTATUS
//   - monadDelay: Delays execution by a duration within the monad
//
// Retry configuration (last 3 parameters):
//   - policy: The retry policy that determines delays and limits
//   - action: The action to retry, which receives the current RetryStatus
//   - check: A predicate that returns true if the result should be retried
//
// Returns:
//   - HKTA: The monadic result after retrying according to the policy
//
// Example conceptual usage with IOEither:
//
//	policy := R.Monoid.Concat(
//		R.LimitRetries(3),
//		R.ExponentialBackoff(100*time.Millisecond),
//	)
//
//	action := func(status R.RetryStatus) IOEither[error, string] {
//		return fetchData() // some IO operation that might fail
//	}
//
//	shouldRetry := func(result string) bool {
//		return result == "" // retry if empty
//	}
//
//	result := Retrying(
//		IOE.Chain[error, string],
//		IOE.Chain[error, R.RetryStatus],
//		IOE.Of[error, string],
//		IOE.Of[error, R.RetryStatus],
//		IOE.Delay[error, R.RetryStatus],
//		policy,
//		action,
//		shouldRetry,
//	)
func Retrying[HKTA, HKTSTATUS, A any](
	monadChain func(func(A) HKTA) func(HKTA) HKTA,
	monadChainStatus func(func(R.RetryStatus) HKTA) func(HKTSTATUS) HKTA,
	monadOf func(A) HKTA,
	monadOfStatus func(R.RetryStatus) HKTSTATUS,
	monadDelay func(time.Duration) func(HKTSTATUS) HKTSTATUS,

	policy R.RetryPolicy,
	action func(R.RetryStatus) HKTA,
	check func(A) bool,
) HKTA {
	// delay callback
	applyDelay := applyAndDelay(monadOfStatus, monadDelay)

	// function to check if we need to retry or not
	checkForRetry := O.FromPredicate(check)

	var f func(status R.RetryStatus) HKTA

	// need some lazy init because we reference it in the chain
	f = func(status R.RetryStatus) HKTA {
		return F.Pipe2(
			status,
			action,
			monadChain(func(a A) HKTA {
				return F.Pipe3(
					a,
					checkForRetry,
					O.Map(func(a A) HKTA {
						return F.Pipe1(
							applyDelay(policy, status),
							monadChainStatus(func(status R.RetryStatus) HKTA {
								return F.Pipe1(
									status.PreviousDelay,
									O.Fold(F.Constant(monadOf(a)), func(_ time.Duration) HKTA {
										return f(status)
									}),
								)
							}),
						)
					}),
					O.GetOrElse(F.Constant(monadOf(a))),
				)
			}),
		)
	}
	// seed
	return f(R.DefaultRetryStatus)
}
