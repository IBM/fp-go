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
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
	RR "github.com/IBM/fp-go/v2/internal/record"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// MonadTraverseArrayWithIndex transforms an array
func MonadTraverseArrayWithIndex[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(int, A) GB) GBS {
	return RA.MonadTraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~[]A, GAAS ~[]GA, E, A any](tas GAAS) GAS {
	return MonadTraverseArray[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseRecord transforms an array
func MonadTraverseRecord[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseRecord transforms an array
func TraverseRecord[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// TraverseRecordWithIndex transforms an array
func TraverseRecordWithIndex[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(K, A) GB) func(AAS) GBS {
	return RR.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		Ap[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, E, A any](tas GAAS) GAS {
	return MonadTraverseRecord[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseArraySeq transforms an array
func MonadTraverseArraySeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArraySeq transforms an array
func TraverseArraySeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// MonadTraverseArrayWithIndexSeq transforms an array
func MonadTraverseArrayWithIndexSeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(int, A) GB) GBS {
	return RA.MonadTraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArrayWithIndexSeq transforms an array
func TraverseArrayWithIndexSeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
func SequenceArraySeq[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~[]A, GAAS ~[]GA, E, A any](tas GAAS) GAS {
	return MonadTraverseArraySeq[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseRecordSeq transforms an array
func MonadTraverseRecordSeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseRecordSeq transforms an array
func TraverseRecordSeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// TraverseRecordWithIndexSeq transforms an array
func TraverseRecordWithIndexSeq[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(K, A) GB) func(AAS) GBS {
	return RR.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApSeq[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
func SequenceRecordSeq[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, E, A any](tas GAAS) GAS {
	return MonadTraverseRecordSeq[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseArrayPar transforms an array
func MonadTraverseArrayPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArrayPar transforms an array
func TraverseArrayPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// MonadTraverseArrayWithIndexPar transforms an array
func MonadTraverseArrayWithIndexPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](tas AAS, f func(int, A) GB) GBS {
	return RA.MonadTraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseArrayWithIndexPar transforms an array
func TraverseArrayWithIndexPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~[]A, BBS ~[]B, E, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceArrayPar converts a homogeneous sequence of either into an either of sequence
func SequenceArrayPar[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~[]A, GAAS ~[]GA, E, A any](tas GAAS) GAS {
	return MonadTraverseArrayPar[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseRecordPar transforms an array
func MonadTraverseRecordPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](tas AAS, f func(A) GB) GBS {
	return RR.MonadTraverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		tas,
		f,
	)
}

// TraverseRecordPar transforms an array
func TraverseRecordPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RR.Traverse[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// TraverseRecordWithIndexPar transforms an array
func TraverseRecordWithIndexPar[GB ~func() either.Either[E, B], GBS ~func() either.Either[E, BBS], AAS ~map[K]A, BBS ~map[K]B, K comparable, E, A, B any](f func(K, A) GB) func(AAS) GBS {
	return RR.TraverseWithIndex[AAS](
		Of[GBS, E, BBS],
		Map[GBS, func() either.Either[E, func(B) BBS], E, BBS, func(B) BBS],
		ApPar[GBS, func() either.Either[E, func(B) BBS], GB],

		f,
	)
}

// SequenceRecordPar converts a homogeneous sequence of either into an either of sequence
func SequenceRecordPar[GA ~func() either.Either[E, A], GAS ~func() either.Either[E, AAS], AAS ~map[K]A, GAAS ~map[K]GA, K comparable, E, A any](tas GAAS) GAS {
	return MonadTraverseRecordPar[GA, GAS](tas, F.Identity[GA])
}
