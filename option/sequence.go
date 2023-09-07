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
)

// Sequence converts an [Option] of some higher kinded type into the higher kinded type of an [Option]
func Sequence[A, HKTA, HKTOA any](
	mof func(Option[A]) HKTOA,
	mmap func(func(A) Option[A]) func(HKTA) HKTOA,
) func(Option[HKTA]) HKTOA {
	return Fold(F.Nullary2(None[A], mof), mmap(Some[A]))
}

// Traverse converts an [Option] of some higher kinded type into the higher kinded type of an [Option]
func Traverse[A, B, HKTB, HKTOB any](
	mof func(Option[B]) HKTOB,
	mmap func(func(B) Option[B]) func(HKTB) HKTOB,
) func(func(A) HKTB) func(Option[A]) HKTOB {
	onNone := F.Nullary2(None[B], mof)
	onSome := mmap(Some[B])
	return func(f func(A) HKTB) func(Option[A]) HKTOB {
		return Fold(onNone, F.Flow2(f, onSome))
	}
}
