// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"github.com/IBM/fp-go/v2/ioeither"
)

// TraverseArray transforms an array
//
//go:inline
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArray(f)
}

// TraverseArrayWithIndex transforms an array
//
//go:inline
func TraverseArrayWithIndex[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArrayWithIndex(f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
//
//go:inline
func SequenceArray[A any](ma []IOResult[A]) IOResult[[]A] {
	return ioeither.SequenceArray(ma)
}

// TraverseRecord transforms a record
//
//go:inline
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecord[K](f)
}

// TraverseRecordWithIndex transforms a record
//
//go:inline
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecordWithIndex(f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
//
//go:inline
func SequenceRecord[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return ioeither.SequenceRecord(ma)
}

// TraverseArraySeq transforms an array
//
//go:inline
func TraverseArraySeq[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArraySeq(f)
}

// TraverseArrayWithIndexSeq transforms an array
//
//go:inline
func TraverseArrayWithIndexSeq[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArrayWithIndexSeq(f)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
//
//go:inline
func SequenceArraySeq[A any](ma []IOResult[A]) IOResult[[]A] {
	return ioeither.SequenceArraySeq(ma)
}

// TraverseRecordSeq transforms a record
//
//go:inline
func TraverseRecordSeq[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecordSeq[K](f)
}

// TraverseRecordWithIndexSeq transforms a record
//
//go:inline
func TraverseRecordWithIndexSeq[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecordWithIndexSeq(f)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
//
//go:inline
func SequenceRecordSeq[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return ioeither.SequenceRecordSeq(ma)
}

// TraverseArrayPar transforms an array
//
//go:inline
func TraverseArrayPar[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArrayPar(f)
}

// TraverseArrayWithIndexPar transforms an array
//
//go:inline
func TraverseArrayWithIndexPar[A, B any](f func(int, A) IOResult[B]) Kleisli[[]A, []B] {
	return ioeither.TraverseArrayWithIndexPar(f)
}

// SequenceArrayPar converts a homogeneous Paruence of either into an either of Paruence
//
//go:inline
func SequenceArrayPar[A any](ma []IOResult[A]) IOResult[[]A] {
	return ioeither.SequenceArrayPar(ma)
}

// TraverseRecordPar transforms a record
//
//go:inline
func TraverseRecordPar[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecordPar[K](f)
}

// TraverseRecordWithIndexPar transforms a record
//
//go:inline
func TraverseRecordWithIndexPar[K comparable, A, B any](f func(K, A) IOResult[B]) Kleisli[map[K]A, map[K]B] {
	return ioeither.TraverseRecordWithIndexPar(f)
}

// SequenceRecordPar converts a homogeneous Paruence of either into an either of Paruence
//
//go:inline
func SequenceRecordPar[K comparable, A any](ma map[K]IOResult[A]) IOResult[map[K]A] {
	return ioeither.SequenceRecordPar(ma)
}
