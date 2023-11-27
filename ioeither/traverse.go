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

package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[E, A, B any](f func(A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArray[IOEither[E, B], IOEither[E, []B], []A](f)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[E, A, B any](f func(int, A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArrayWithIndex[IOEither[E, B], IOEither[E, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return G.SequenceArray[IOEither[E, A], IOEither[E, []A]](ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, E, A, B any](f func(A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecord[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// TraverseRecordWithIndex transforms a record
func TraverseRecordWithIndex[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecordWithIndex[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return G.SequenceRecord[IOEither[E, A], IOEither[E, map[K]A]](ma)
}

// TraverseArraySeq transforms an array
func TraverseArraySeq[E, A, B any](f func(A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArraySeq[IOEither[E, B], IOEither[E, []B], []A](f)
}

// TraverseArrayWithIndexSeq transforms an array
func TraverseArrayWithIndexSeq[E, A, B any](f func(int, A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArrayWithIndexSeq[IOEither[E, B], IOEither[E, []B], []A](f)
}

// SequenceArraySeq converts a homogeneous sequence of either into an either of sequence
func SequenceArraySeq[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return G.SequenceArraySeq[IOEither[E, A], IOEither[E, []A]](ma)
}

// TraverseRecordSeq transforms a record
func TraverseRecordSeq[K comparable, E, A, B any](f func(A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecordSeq[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// TraverseRecordWithIndexSeq transforms a record
func TraverseRecordWithIndexSeq[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecordWithIndexSeq[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecordSeq converts a homogeneous sequence of either into an either of sequence
func SequenceRecordSeq[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return G.SequenceRecordSeq[IOEither[E, A], IOEither[E, map[K]A]](ma)
}

// TraverseArrayPar transforms an array
func TraverseArrayPar[E, A, B any](f func(A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArrayPar[IOEither[E, B], IOEither[E, []B], []A](f)
}

// TraverseArrayWithIndexPar transforms an array
func TraverseArrayWithIndexPar[E, A, B any](f func(int, A) IOEither[E, B]) func([]A) IOEither[E, []B] {
	return G.TraverseArrayWithIndexPar[IOEither[E, B], IOEither[E, []B], []A](f)
}

// SequenceArrayPar converts a homogeneous Paruence of either into an either of Paruence
func SequenceArrayPar[E, A any](ma []IOEither[E, A]) IOEither[E, []A] {
	return G.SequenceArrayPar[IOEither[E, A], IOEither[E, []A]](ma)
}

// TraverseRecordPar transforms a record
func TraverseRecordPar[K comparable, E, A, B any](f func(A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecordPar[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// TraverseRecordWithIndexPar transforms a record
func TraverseRecordWithIndexPar[K comparable, E, A, B any](f func(K, A) IOEither[E, B]) func(map[K]A) IOEither[E, map[K]B] {
	return G.TraverseRecordWithIndexPar[IOEither[E, B], IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecordPar converts a homogeneous Paruence of either into an either of Paruence
func SequenceRecordPar[K comparable, E, A any](ma map[K]IOEither[E, A]) IOEither[E, map[K]A] {
	return G.SequenceRecordPar[IOEither[E, A], IOEither[E, map[K]A]](ma)
}
