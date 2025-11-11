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
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	RR "github.com/IBM/fp-go/v2/internal/record"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() either.Either[L, B], GIOBS ~func() either.Either[L, BBS], AAS ~[]A, BBS ~[]B, E, L, A, B any](ma AAS, f func(A) GB) GBS {
	return RA.MonadTraverse(
		Of[GBS, GIOBS, E, L, BBS],
		Map[GBS, func(E) func() either.Either[L, func(B) BBS], GIOBS, func() either.Either[L, func(B) BBS], E, L, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() either.Either[L, func(B) BBS], GIOB, GIOBS, func() either.Either[L, func(B) BBS], E, L, B, BBS],

		ma, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() either.Either[L, B], GIOBS ~func() either.Either[L, BBS], AAS ~[]A, BBS ~[]B, E, L, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, GIOBS, E, L, BBS],
		Map[GBS, func(E) func() either.Either[L, func(B) BBS], GIOBS, func() either.Either[L, func(B) BBS], E, L, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() either.Either[L, func(B) BBS], GIOB, GIOBS, func() either.Either[L, func(B) BBS], E, L, B, BBS],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[GB ~func(E) GIOB, GBS ~func(E) GIOBS, GIOB ~func() either.Either[L, B], GIOBS ~func() either.Either[L, BBS], AAS ~[]A, BBS ~[]B, E, L, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, GIOBS, E, L, BBS],
		Map[GBS, func(E) func() either.Either[L, func(B) BBS], GIOBS, func() either.Either[L, func(B) BBS], E, L, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) func() either.Either[L, func(B) BBS], GIOB, GIOBS, func() either.Either[L, func(B) BBS], E, L, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(E) GIOA, GAS ~func(E) GIOAS, GIOA ~func() either.Either[L, A], GIOAS ~func() either.Either[L, AAS], AAS ~[]A, GAAS ~[]GA, E, L, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}

// MonadTraverseRecord transforms an array
func MonadTraverseRecord[GB ~func(C) GIOB, GBS ~func(C) GIOBS, GIOB ~func() either.Either[E, B], GIOBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, C, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse(
		Of[GBS, GIOBS, C, E, BBS],
		Map[GBS, func(C) func() either.Either[E, func(B) BBS], GIOBS, func() either.Either[E, func(B) BBS], C, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(C) func() either.Either[E, func(B) BBS], GIOB, GIOBS, func() either.Either[E, func(B) BBS], C, E, B, BBS],

		tas,
		f,
	)
}

// TraverseRecord transforms a record
func TraverseRecord[GB ~func(C) GIOB, GBS ~func(C) GIOBS, GIOB ~func() either.Either[E, B], GIOBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, C, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, GIOBS, C, E, BBS],
		Map[GBS, func(C) func() either.Either[E, func(B) BBS], GIOBS, func() either.Either[E, func(B) BBS], C, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(C) func() either.Either[E, func(B) BBS], GIOB, GIOBS, func() either.Either[E, func(B) BBS], C, E, B, BBS],

		f,
	)
}

// TraverseRecordWithIndex transforms a record
func TraverseRecordWithIndex[GB ~func(C) GIOB, GBS ~func(C) GIOBS, GIOB ~func() either.Either[E, B], GIOBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, C, E, A, B any](f func(K, A) GB) func(AAS) GBS {
	return RR.TraverseWithIndex[AAS](
		Of[GBS, GIOBS, C, E, BBS],
		Map[GBS, func(C) func() either.Either[E, func(B) BBS], GIOBS, func() either.Either[E, func(B) BBS], C, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(C) func() either.Either[E, func(B) BBS], GIOB, GIOBS, func() either.Either[E, func(B) BBS], C, E, B, BBS],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[GA ~func(C) GIOA, GAS ~func(C) GIOAS, GIOA ~func() either.Either[E, A], GIOAS ~func() either.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, C, E, A any](tas GAAS) GAS {
	return MonadTraverseRecord[GA, GAS](tas, F.Identity[GA])
}
