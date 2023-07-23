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

package readerioeither

import (
	G "github.com/IBM/fp-go/context/readerioeither/generic"
)

// TraverseArray uses transforms an array [[]A] into [[]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[[]B]]
func TraverseArray[A, B any](f func(A) ReaderIOEither[B]) func([]A) ReaderIOEither[[]B] {
	return G.TraverseArray[[]A, ReaderIOEither[[]B]](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIOEither[A]) ReaderIOEither[[]A] {
	return G.SequenceArray[[]A, []ReaderIOEither[A], ReaderIOEither[[]A]](ma)
}

// TraverseRecord uses transforms a record [map[K]A] into [map[K]ReaderIOEither[B]] and then resolves that into a [ReaderIOEither[map[K]B]]
func TraverseRecord[K comparable, A, B any](f func(A) ReaderIOEither[B]) func(map[K]A) ReaderIOEither[map[K]B] {
	return G.TraverseRecord[K, map[K]A, ReaderIOEither[map[K]B]](f)
}

// SequenceRecord converts a homogeneous sequence of either into an either of sequence
func SequenceRecord[K comparable, A any](ma map[K]ReaderIOEither[A]) ReaderIOEither[map[K]A] {
	return G.SequenceRecord[K, map[K]A, map[K]ReaderIOEither[A], ReaderIOEither[map[K]A]](ma)
}
