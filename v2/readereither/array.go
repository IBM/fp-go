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

package readereither

import (
	G "github.com/IBM/fp-go/v2/readereither/generic"
)

// TraverseArray transforms an array
func TraverseArray[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B] {
	return G.TraverseArray[ReaderEither[E, L, B], ReaderEither[E, L, []B], []A](f)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[E, L, A, B any](f func(int, A) ReaderEither[E, L, B]) func([]A) ReaderEither[E, L, []B] {
	return G.TraverseArrayWithIndex[ReaderEither[E, L, B], ReaderEither[E, L, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[E, L, A any](ma []ReaderEither[E, L, A]) ReaderEither[E, L, []A] {
	return G.SequenceArray[ReaderEither[E, L, A], ReaderEither[E, L, []A]](ma)
}
