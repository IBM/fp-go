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
	F "github.com/IBM/fp-go/v2/function"
	G "github.com/IBM/fp-go/v2/internal/apply"
)

const (
	// useParallel is the feature flag to control if we use the parallel or the sequential implementation of ap
	useParallel = true
)

// Deprecated:  MonadApSeq implements the applicative on a single thread by first executing mab and the ma
func MonadApSeq[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	return MonadChain(mab, F.Bind1st(MonadMap[GA, GB], ma))
}

// MonadApPar implements the applicative on two threads, the main thread executes mab and the actuall
// apply operation and the second thread computes ma. Communication between the threads happens via a channel
//
// Deprecated:
func MonadApPar[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	return MakeIO[GB](func() B {
		c := make(chan A)
		go func() {
			c <- ma()
			close(c)
		}()
		return mab()(<-c)
	})
}

// MonadAp implements the `ap` operation. Depending on a feature flag this will be sequential or parallel, the preferred implementation
// is parallel
//
// Deprecated:
func MonadAp[GA ~func() A, GB ~func() B, GAB ~func() func(A) B, A, B any](mab GAB, ma GA) GB {
	if useParallel {
		return MonadApPar[GA, GB](mab, ma)
	}
	return MonadApSeq[GA, GB](mab, ma)
}

// MonadApFirst combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func MonadApFirst[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadAp[GB, GA, GBA, B, A],
		MonadMap[GA, GBA, A, func(B) A],

		first,
		second,
	)
}

// MonadApFirstPar combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func MonadApFirstPar[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadApPar[GB, GA, GBA, B, A],
		MonadMap[GA, GBA, A, func(B) A],

		first,
		second,
	)
}

// MonadApFirstSeq combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func MonadApFirstSeq[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](first GA, second GB) GA {
	return G.MonadApFirst(
		MonadApSeq[GB, GA, GBA, B, A],
		MonadMap[GA, GBA, A, func(B) A],

		first,
		second,
	)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func ApFirst[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		Ap[GA, GBA, GB, A, B],
		Map[GA, GBA, A, func(B) A],

		second,
	)
}

// ApFirstPar combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func ApFirstPar[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		ApPar[GA, GBA, GB, A, B],
		Map[GA, GBA, A, func(B) A],

		second,
	)
}

// ApFirstSeq combines two effectful actions, keeping only the result of the first.
//
// Deprecated:
func ApFirstSeq[GA ~func() A, GB ~func() B, GBA ~func() func(B) A, A, B any](second GB) func(GA) GA {
	return G.ApFirst(
		ApSeq[GA, GBA, GB, A, B],
		Map[GA, GBA, A, func(B) A],

		second,
	)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
//
// Deprecated:
func MonadApSecond[GA ~func() A, GB ~func() B, GBB ~func() func(B) B, A, B any](first GA, second GB) GB {
	return G.MonadApSecond(
		MonadAp[GB, GB, GBB, B, B],
		MonadMap[GA, GBB, A, func(B) B],

		first,
		second,
	)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
//
// Deprecated:
func ApSecond[GA ~func() A, GB ~func() B, GBB ~func() func(B) B, A, B any](second GB) func(GA) GB {
	return G.ApSecond(
		ApSeq[GB, GBB, GB, B, B],
		Map[GA, GBB, A, func(B) B],

		second,
	)
}
