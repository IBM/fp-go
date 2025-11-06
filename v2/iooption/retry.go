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

package iooption

import (
	R "github.com/IBM/fp-go/v2/retry"
	G "github.com/IBM/fp-go/v2/retry/generic"
)

// Retrying will retry the actions according to the check policy
func Retrying[A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) IOOption[A],
	check func(A) bool,
) IOOption[A] {
	// get an implementation for the types
	return G.Retrying(
		Chain[A, A],
		Chain[R.RetryStatus, A],
		Of[A],
		Of[R.RetryStatus],
		Delay[R.RetryStatus],

		policy, action, check)
}
