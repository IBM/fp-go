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

package option

import (
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// TraverseArray transforms an array
func TraverseArrayG[GA ~[]A, GB ~[]B, A, B any](f func(A) Option[B]) func(GA) Option[GB] {
	return RA.Traverse[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) Option[B]) func([]A) Option[[]B] {
	return TraverseArrayG[[]A, []B](f)
}

// TraverseArrayWithIndexG transforms an array
func TraverseArrayWithIndexG[GA ~[]A, GB ~[]B, A, B any](f func(int, A) Option[B]) func(GA) Option[GB] {
	return RA.TraverseWithIndex[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[A, B any](f func(int, A) Option[B]) func([]A) Option[[]B] {
	return TraverseArrayWithIndexG[[]A, []B](f)
}

func SequenceArrayG[GA ~[]A, GOA ~[]Option[A], A any](ma GOA) Option[GA] {
	return TraverseArrayG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []Option[A]) Option[[]A] {
	return SequenceArrayG[[]A](ma)
}

// CompactArrayG discards the none values and keeps the some values
func CompactArrayG[A1 ~[]Option[A], A2 ~[]A, A any](fa A1) A2 {
	return RA.Reduce(fa, func(out A2, value Option[A]) A2 {
		return MonadFold(value, F.Constant(out), F.Bind1st(RA.Append[A2, A], out))
	}, make(A2, 0, len(fa)))
}

// CompactArray discards the none values and keeps the some values
func CompactArray[A any](fa []Option[A]) []A {
	return CompactArrayG[[]Option[A], []A](fa)
}
