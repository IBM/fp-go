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

package reader

import (
	G "github.com/IBM/fp-go/v2/reader/generic"
)

// TraverseArray transforms an array
func TraverseArray[R, A, B any](f func(A) Reader[R, B]) func([]A) Reader[R, []B] {
	return G.TraverseArray[Reader[R, B], Reader[R, []B], []A](f)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[R, A, B any](f func(int, A) Reader[R, B]) func([]A) Reader[R, []B] {
	return G.TraverseArrayWithIndex[Reader[R, B], Reader[R, []B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[R, A any](ma []Reader[R, A]) Reader[R, []A] {
	return G.SequenceArray[Reader[R, A], Reader[R, []A]](ma)
}
