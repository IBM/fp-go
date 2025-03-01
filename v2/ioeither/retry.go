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

package ioeither

import (
	ET "github.com/IBM/fp-go/v2/either"
	G "github.com/IBM/fp-go/v2/ioeither/generic"
	R "github.com/IBM/fp-go/v2/retry"
)

// Retrying will retry the actions according to the check policy
//
// policy - refers to the retry policy
// action - converts a status into an operation to be executed
// check  - checks if the result of the action needs to be retried
func Retrying[E, A any](
	policy R.RetryPolicy,
	action func(R.RetryStatus) IOEither[E, A],
	check func(ET.Either[E, A]) bool,
) IOEither[E, A] {
	return G.Retrying(policy, action, check)
}
