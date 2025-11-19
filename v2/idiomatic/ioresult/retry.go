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

package ioresult

import (
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/result"
	R "github.com/IBM/fp-go/v2/retry"
)

// Retrying retries an IOResult computation according to a retry policy.
// The action receives retry status information on each attempt.
// The check function determines if the result warrants another retry.
func Retrying[A any](
	policy R.RetryPolicy,
	action Kleisli[R.RetryStatus, A],
	check func(A, error) bool,
) IOResult[A] {
	fromResult := io.Retrying(policy,
		func(rs R.RetryStatus) IO[Result[A]] {
			return func() Result[A] {
				return result.TryCatchError(action(rs)())
			}
		},
		func(a Result[A]) bool {
			return check(result.Unwrap(a))
		},
	)
	return func() (A, error) {
		return result.Unwrap(fromResult())
	}
}
