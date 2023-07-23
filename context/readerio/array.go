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

package readerio

import (
	IO "github.com/IBM/fp-go/io"
	R "github.com/IBM/fp-go/readerio/generic"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f func(A) ReaderIO[B]) func([]A) ReaderIO[[]B] {
	return R.TraverseArray[ReaderIO[B], ReaderIO[[]B], IO.IO[B], IO.IO[[]B], []A](f)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []ReaderIO[A]) ReaderIO[[]A] {
	return R.SequenceArray[ReaderIO[A], ReaderIO[[]A]](ma)
}
