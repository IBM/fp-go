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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	G "github.com/IBM/fp-go/v2/internal/apply"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBA ~func() ET.Either[E, func(B) A], E, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadAp[GA, GBA, GB],
		MonadMap[GA, GBA, E, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBA ~func() ET.Either[E, func(B) A], E, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		MonadAp[GA, GBA, GB],
		MonadMap[GA, GBA, E, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBB ~func() ET.Either[E, func(B) B], E, A, B any](first GA, second GB) GB {
	return G.MonadApSecond(
		MonadAp[GB, GBB, GB],
		MonadMap[GA, GBB, E, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GBB ~func() ET.Either[E, func(B) B], E, A, B any](second GB) func(GA) GB {
	return G.ApSecond(
		MonadAp[GB, GBB, GB],
		MonadMap[GA, GBB, E, A, func(B) B],

		second,
	)
}
