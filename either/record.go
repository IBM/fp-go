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
	RR "github.com/IBM/fp-go/internal/record"
)

// TraverseRecord transforms a record of options into an option of a record
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, E, A, B any](f func(A) Either[E, B]) func(GA) Either[E, GB] {
	return RR.Traverse[GA](
		Of[E, GB],
		Map[E, GB, func(B) GB],
		Ap[GB, E, B],
		f,
	)
}

// TraverseRecord transforms a record of eithers into an either of a record
func TraverseRecord[K comparable, E, A, B any](f func(A) Either[E, B]) func(map[K]A) Either[E, map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Either[E, A], K comparable, E, A any](ma GOA) Either[E, GA] {
	return TraverseRecordG[GOA, GA](F.Identity[Either[E, A]])(ma)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]Either[E, A]) Either[E, map[K]A] {
	return SequenceRecordG[map[K]A](ma)
}
