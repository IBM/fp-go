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
	"github.com/IBM/fp-go/v2/internal/apply"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
//
//go:inline
func MonadApFirst[A, B any](first IOResult[A], second IOResult[B]) IOResult[A] {
	return apply.MonadApFirst(
		MonadAp[A, B],
		MonadMap[A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
//
//go:inline
func ApFirst[A, B any](second IOResult[B]) Operator[A, A] {
	return apply.ApFirst(
		Ap[A, B],
		Map[A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
//
//go:inline
func MonadApSecond[A, B any](first IOResult[A], second IOResult[B]) IOResult[B] {
	return apply.MonadApSecond(
		MonadAp[B, B],
		MonadMap[A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
//
//go:inline
func ApSecond[A, B any](second IOResult[B]) Operator[A, B] {
	return apply.ApSecond(
		Ap[B, B],
		Map[A, func(B) B],

		second,
	)
}
