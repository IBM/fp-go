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

package io

import (
	R "github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

type (
	// RetryStatus is an IO computation that returns retry status information.
	RetryStatus = IO[R.RetryStatus]
)

// Retrying retries an IO action according to a retry policy until it succeeds or the policy gives up.
//
// Parameters:
//   - policy: The retry policy that determines delays and maximum attempts
//   - action: A function that takes retry status and returns an IO computation
//   - check: A predicate that determines if the result should trigger a retry (true = retry)
//
// The action receives retry status information (attempt number, cumulative delay, etc.)
// which can be used for logging or conditional behavior.
//
// Example:
//
//	result := io.Retrying(
//	    retry.ExponentialBackoff(time.Second, 5),
//	    func(status retry.RetryStatus) io.IO[Response] {
//	        log.Printf("Attempt %d", status.IterNumber)
//	        return fetchData()
//	    },
//	    func(r Response) bool { return r.StatusCode >= 500 },
//	)
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check func(A) bool,
) IO[A] {
	// get an implementation for the types
	return RG.Retrying(
		Chain[A, A],
		Chain[R.RetryStatus, A],
		Of[A],
		Of[R.RetryStatus],
		Delay[R.RetryStatus],

		policy,
		action,
		check,
	)
}
