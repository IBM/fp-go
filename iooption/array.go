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

package iooption

import (
	G "github.com/IBM/fp-go/iooption/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) IOOption[B]) func([]A) IOOption[[]B] {
	return G.TraverseArray[IOOption[B], IOOption[[]B], []A](f)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[A, B any](f func(int, A) IOOption[B]) func([]A) IOOption[[]B] {
	return G.TraverseArrayWithIndex[IOOption[B], IOOption[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []IOOption[A]) IOOption[[]A] {
	return G.SequenceArray[IOOption[A], IOOption[[]A], []IOOption[A], []A, A](ma)
}
