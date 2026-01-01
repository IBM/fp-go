// Copyright (c) 2025 IBM Corp.
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

package readerio

import (
	"github.com/IBM/fp-go/v2/retry"
	RG "github.com/IBM/fp-go/v2/retry/generic"
)

//go:inline
func Retrying[A any](
	policy retry.RetryPolicy,
	action Kleisli[retry.RetryStatus, A],
	check func(A) bool,
) ReaderIO[A] {
	// get an implementation for the types
	return RG.RetryingWithTrampoline(
		Chain[A, Trampoline[retry.RetryStatus, A]],
		Map[retry.RetryStatus, Trampoline[retry.RetryStatus, A]],
		Of[Trampoline[retry.RetryStatus, A]],
		Of[retry.RetryStatus],
		Delay[retry.RetryStatus],

		TailRec,

		policy,
		action,
		check,
	)
}
