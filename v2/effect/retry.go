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
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/retry"
)

// Retrying executes an effect with retry logic based on a policy and check predicate.
// The effect is retried according to the policy until either:
//   - The effect succeeds and the check predicate returns false
//   - The retry policy is exhausted
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - A: The type of the success value
//
// # Parameters
//
//   - policy: The retry policy defining retry limits and delays
//   - action: An effectful computation that receives retry status and produces a value
//   - check: A predicate that determines if the result should trigger a retry
//
// # Returns
//
//   - Effect[C, A]: An effect that retries according to the policy
//
// # Example
//
//	policy := retry.LimitRetries(3)
//	eff := effect.Retrying[MyContext, string](
//		policy,
//		func(status retry.RetryStatus) Effect[MyContext, string] {
//			return fetchData() // may fail
//		},
//		func(result Result[string]) bool {
//			return result.IsLeft() // retry on error
//		},
//	)
//	// Retries up to 3 times if fetchData fails
func Retrying[C, A any](
	policy retry.RetryPolicy,
	action Kleisli[C, retry.RetryStatus, A],
	check Predicate[Result[A]],
) Effect[C, A] {
	return readerreaderioresult.Retrying(policy, action, check)
}
