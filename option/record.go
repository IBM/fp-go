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

package option

import (
	F "github.com/IBM/fp-go/function"
	RR "github.com/IBM/fp-go/internal/record"
)

// TraverseRecordG transforms a record of options into an option of a record
func TraverseRecordG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(A) Option[B]) func(GA) Option[GB] {
	return RR.Traverse[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseRecord transforms a record of options into an option of a record
func TraverseRecord[K comparable, A, B any](f func(A) Option[B]) func(map[K]A) Option[map[K]B] {
	return TraverseRecordG[map[K]A, map[K]B](f)
}

// TraverseRecordWithIndexG transforms a record of options into an option of a record
func TraverseRecordWithIndexG[GA ~map[K]A, GB ~map[K]B, K comparable, A, B any](f func(K, A) Option[B]) func(GA) Option[GB] {
	return RR.TraverseWithIndex[GA](
		Of[GB],
		Map[GB, func(B) GB],
		Ap[GB, B],

		f,
	)
}

// TraverseRecordWithIndex transforms a record of options into an option of a record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) Option[B]) func(map[K]A) Option[map[K]B] {
	return TraverseRecordWithIndexG[map[K]A, map[K]B](f)
}

func SequenceRecordG[GA ~map[K]A, GOA ~map[K]Option[A], K comparable, A any](ma GOA) Option[GA] {
	return TraverseRecordG[GOA, GA](F.Identity[Option[A]])(ma)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]Option[A]) Option[map[K]A] {
	return SequenceRecordG[map[K]A](ma)
}

func upsertAtReadWrite[M ~map[K]V, K comparable, V any](r M, k K, v V) M {
	r[k] = v
	return r
}

// CompactRecordG discards the noe values and keeps the some values
func CompactRecordG[M1 ~map[K]Option[A], M2 ~map[K]A, K comparable, A any](m M1) M2 {
	bnd := F.Bind12of3(upsertAtReadWrite[M2])
	return RR.ReduceWithIndex(m, func(key K, m M2, value Option[A]) M2 {
		return MonadFold(value, F.Constant(m), bnd(m, key))
	}, make(M2))
}

// CompactRecord discards the noe values and keeps the some values
func CompactRecord[K comparable, A any](m map[K]Option[A]) map[K]A {
	return CompactRecordG[map[K]Option[A], map[K]A](m)
}
