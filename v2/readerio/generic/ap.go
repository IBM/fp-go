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

package generic

import (
	G "github.com/IBM/fp-go/v2/internal/apply"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[GRA ~func(R) GA, GRB ~func(R) GB, GRBA ~func(R) GBA, GA ~func() A, GB ~func() B, GBA ~func() func(B) A, R, A, B any](first GRA, second GRB) GRA {
	return G.MonadApFirst(
		MonadAp[GRB, GRA, GRBA, GB, GA, GBA, R, B, A],
		MonadMap[GRA, GRBA, GA, GBA, R, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[GRA ~func(R) GA, GRB ~func(R) GB, GRBA ~func(R) GBA, GA ~func() A, GB ~func() B, GBA ~func() func(B) A, R, A, B any](second GRB) func(GRA) GRA {
	return G.ApFirst(
		Ap[GRB, GRA, GRBA, GB, GA, GBA, R, B, A],
		Map[GRA, GRBA, GA, GBA, R, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[GRA ~func(R) GA, GRB ~func(R) GB, GRBB ~func(R) GBB, GA ~func() A, GB ~func() B, GBB ~func() func(B) B, R, A, B any](first GRA, second GRB) GRB {
	return G.MonadApSecond(
		MonadAp[GRB, GRB, GRBB, GB, GB, GBB, R, B, B],
		MonadMap[GRA, GRBB, GA, GBB, R, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[GRA ~func(R) GA, GRB ~func(R) GB, GRBB ~func(R) GBB, GA ~func() A, GB ~func() B, GBB ~func() func(B) B, R, A, B any](second GRB) func(GRA) GRB {
	return G.ApSecond(
		Ap[GRB, GRB, GRBB, GB, GB, GBB, R, B, B],
		Map[GRA, GRBB, GA, GBB, R, A, func(B) B],

		second,
	)
}
