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
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	RA "github.com/IBM/fp-go/v2/internal/array"
)

// MonadTraverseArray transforms an array
func MonadTraverseArray[GB ~func(E) ET.Either[L, B], GBS ~func(E) ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, L, E, A, B any](ma AAS, f func(A) GB) GBS {
	return RA.MonadTraverse(
		Of[GBS, L, E, BBS],
		Map[GBS, func(E) ET.Either[L, func(B) BBS], L, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) ET.Either[L, func(B) BBS], L, E, B, BBS],

		ma, f,
	)
}

// TraverseArray transforms an array
func TraverseArray[GB ~func(E) ET.Either[L, B], GBS ~func(E) ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, L, E, A, B any](f func(A) GB) func(AAS) GBS {
	return RA.Traverse[AAS](
		Of[GBS, L, E, BBS],
		Map[GBS, func(E) ET.Either[L, func(B) BBS], L, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) ET.Either[L, func(B) BBS], L, E, B, BBS],

		f,
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[GB ~func(E) ET.Either[L, B], GBS ~func(E) ET.Either[L, BBS], AAS ~[]A, BBS ~[]B, L, E, A, B any](f func(int, A) GB) func(AAS) GBS {
	return RA.TraverseWithIndex[AAS](
		Of[GBS, L, E, BBS],
		Map[GBS, func(E) ET.Either[L, func(B) BBS], L, E, BBS, func(B) BBS],
		Ap[GB, GBS, func(E) ET.Either[L, func(B) BBS], L, E, B, BBS],

		f,
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[GA ~func(E) ET.Either[L, A], GAS ~func(E) ET.Either[L, AAS], AAS ~[]A, GAAS ~[]GA, L, E, A any](ma GAAS) GAS {
	return MonadTraverseArray[GA, GAS](ma, F.Identity[GA])
}
