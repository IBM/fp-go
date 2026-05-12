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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return itereither.TraverseArray(f)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[A, B any](f func(int, A) SeqResult[B]) Kleisli[[]A, []B] {
	return itereither.TraverseArrayWithIndex(f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []SeqResult[A]) SeqResult[[]A] {
	return itereither.SequenceArray(ma)
}

// TraverseRecord transforms a record
func TraverseRecord[K comparable, A, B any](f Kleisli[A, B]) Kleisli[map[K]A, map[K]B] {
	return itereither.TraverseRecord[K](f)
}

// TraverseRecordWithIndex transforms a record
func TraverseRecordWithIndex[K comparable, A, B any](f func(K, A) SeqResult[B]) Kleisli[map[K]A, map[K]B] {
	return itereither.TraverseRecordWithIndex(f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]SeqResult[A]) SeqResult[map[K]A] {
	return itereither.SequenceRecord(ma)
}
