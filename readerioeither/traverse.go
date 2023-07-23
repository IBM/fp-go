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
	IOE "github.com/IBM/fp-go/ioeither"
	G "github.com/IBM/fp-go/readerioeither/generic"
)

// TraverseArray transforms an array
func TraverseArray[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func([]A) ReaderIOEither[R, E, []B] {
	return G.TraverseArray[ReaderIOEither[R, E, B], ReaderIOEither[R, E, []B], IOE.IOEither[E, B], IOE.IOEither[E, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of Readers into a Reader of a sequence
func SequenceArray[R, E, A any](ma []ReaderIOEither[R, E, A]) ReaderIOEither[R, E, []A] {
	return G.SequenceArray[ReaderIOEither[R, E, A], ReaderIOEither[R, E, []A]](ma)
}

// TraverseRecord transforms an array
func TraverseRecord[R any, K comparable, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(map[K]A) ReaderIOEither[R, E, map[K]B] {
	return G.TraverseRecord[ReaderIOEither[R, E, B], ReaderIOEither[R, E, map[K]B], IOE.IOEither[E, B], IOE.IOEither[E, map[K]B], map[K]A](f)
}

// SequenceRecord converts a homogeneous sequence of Readers into a Reader of a sequence
func SequenceRecord[R any, K comparable, E, A any](ma map[K]ReaderIOEither[R, E, A]) ReaderIOEither[R, E, map[K]A] {
	return G.SequenceRecord[ReaderIOEither[R, E, A], ReaderIOEither[R, E, map[K]A]](ma)
}
