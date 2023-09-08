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

package either

import (
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
)

// TraverseArray transforms an array
func TraverseArrayG[GA ~[]A, GB ~[]B, E, A, B any](f func(A) Either[E, B]) func(GA) Either[E, GB] {
	return RA.Traverse[GA](
		Of[E, GB],
		Map[E, GB, func(B) GB],
		Ap[GB, E, B],

		f,
	)
}

// TraverseArray transforms an array
func TraverseArray[E, A, B any](f func(A) Either[E, B]) func([]A) Either[E, []B] {
	return TraverseArrayG[[]A, []B](f)
}

func SequenceArrayG[GA ~[]A, GOA ~[]Either[E, A], E, A any](ma GOA) Either[E, GA] {
	return TraverseArrayG[GOA, GA](F.Identity[Either[E, A]])(ma)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, A any](ma []Either[E, A]) Either[E, []A] {
	return SequenceArrayG[[]A](ma)
}

// CompactArrayG discards the none values and keeps the some values
func CompactArrayG[A1 ~[]Either[E, A], A2 ~[]A, E, A any](fa A1) A2 {
	return RA.Reduce(fa, func(out A2, value Either[E, A]) A2 {
		return MonadFold(value, F.Constant1[E](out), F.Bind1st(RA.Append[A2, A], out))
	}, make(A2, len(fa)))
}

// CompactArray discards the none values and keeps the some values
func CompactArray[E, A any](fa []Either[E, A]) []A {
	return CompactArrayG[[]Either[E, A], []A](fa)
}
