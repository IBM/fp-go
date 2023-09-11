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
	F "github.com/IBM/fp-go/function"
	RA "github.com/IBM/fp-go/internal/array"
	RR "github.com/IBM/fp-go/internal/record"
)

func MonadTraverseArray[GB ~func() B, GBS ~func() BBS, AAS ~[]A, BBS ~[]B, A, B any](tas AAS, f func(A) GB) GBS {
	return RA.MonadTraverse(
		Of[GBS, BBS],
		Map[GBS, func() func(B) BBS, BBS, func(B) BBS],
		Ap[GBS, func() func(B) BBS, GB],

		tas,
		f,
	)
}

func TraverseArray[GB ~func() B, GBS ~func() BBS, AAS ~[]A, BBS ~[]B, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, BBS],
		Map[GBS, func() func(B) BBS, BBS, func(B) BBS],
		Ap[GBS, func() func(B) BBS, GB],

		f,
	)
}

func TraverseArrayWithIndex[GB ~func() B, GBS ~func() BBS, AAS ~[]A, BBS ~[]B, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, BBS],
		Map[GBS, func() func(B) BBS, BBS, func(B) BBS],
		Ap[GBS, func() func(B) BBS, GB],

		f,
	)
}

func SequenceArray[GA ~func() A, GAS ~func() AAS, AAS ~[]A, GAAS ~[]GA, A any](tas GAAS) GAS {
	return MonadTraverseArray[GA, GAS](tas, F.Identity[GA])
}

// MonadTraverseRecord transforms a record using an IO transform an IO of a record
func MonadTraverseRecord[GB ~func() B, GBS ~func() MB, MA ~map[K]A, MB ~map[K]B, K comparable, A, B any](ma MA, f func(A) GB) GBS {
	return RR.MonadTraverse[MA](
		Of[GBS, MB],
		Map[GBS, func() func(B) MB, MB, func(B) MB],
		Ap[GBS, func() func(B) MB, GB],
		ma, f,
	)
}

// TraverseRecord transforms a record using an IO transform an IO of a record
func TraverseRecord[GB ~func() B, GBS ~func() MB, MA ~map[K]A, MB ~map[K]B, K comparable, A, B any](f func(A) GB) func(MA) GBS {
	return RR.Traverse[MA](
		Of[GBS, MB],
		Map[GBS, func() func(B) MB, MB, func(B) MB],
		Ap[GBS, func() func(B) MB, GB],
		f,
	)
}

// TraverseRecordWithIndex transforms a record using an IO transform an IO of a record
func TraverseRecordWithIndex[GB ~func() B, GBS ~func() MB, MA ~map[K]A, MB ~map[K]B, K comparable, A, B any](f func(K, A) GB) func(MA) GBS {
	return RR.TraverseWithIndex[MA](
		Of[GBS, MB],
		Map[GBS, func() func(B) MB, MB, func(B) MB],
		Ap[GBS, func() func(B) MB, GB],
		f,
	)
}

func SequenceRecord[GA ~func() A, GAS ~func() AAS, AAS ~map[K]A, GAAS ~map[K]GA, K comparable, A any](tas GAAS) GAS {
	return MonadTraverseRecord[GA, GAS](tas, F.Identity[GA])
}
