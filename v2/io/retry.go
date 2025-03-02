// Copyright (c) 2023 IBM Corp.
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
	RetryStatus = IO[R.RetryStatus]
)

// Retrying will retry the actions according to the check policy
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) IO[A],
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
