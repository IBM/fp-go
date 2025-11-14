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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// FindFirstWithIndex finds the first element which satisfies a predicate (or a refinement) function
func FindFirstWithIndex[AS ~[]A, PRED ~func(int, A) bool, A any](pred PRED) func(AS) O.Option[A] {
	none := O.None[A]()
	return func(as AS) O.Option[A] {
		for i, a := range as {
			if pred(i, a) {
				return O.Some(a)
			}
		}
		return none
	}
}

// FindFirst finds the first element which satisfies a predicate (or a refinement) function
func FindFirst[AS ~[]A, PRED ~func(A) bool, A any](pred PRED) func(AS) O.Option[A] {
	return FindFirstWithIndex[AS](F.Ignore1of2[int](pred))
}

// FindFirstMapWithIndex finds the first element returned by an [O.Option] based selector function
func FindFirstMapWithIndex[AS ~[]A, PRED ~func(int, A) O.Option[B], A, B any](pred PRED) func(AS) O.Option[B] {
	none := O.None[B]()
	return func(as AS) O.Option[B] {
		for i := range len(as) {
			out := pred(i, as[i])
			if O.IsSome(out) {
				return out
			}
		}
		return none
	}
}

// FindFirstMap finds the first element returned by an [O.Option] based selector function
func FindFirstMap[AS ~[]A, PRED ~func(A) O.Option[B], A, B any](pred PRED) func(AS) O.Option[B] {
	return FindFirstMapWithIndex[AS](F.Ignore1of2[int](pred))
}

// FindLastWithIndex finds the first element which satisfies a predicate (or a refinement) function
func FindLastWithIndex[AS ~[]A, PRED ~func(int, A) bool, A any](pred PRED) func(AS) O.Option[A] {
	none := O.None[A]()
	return func(as AS) O.Option[A] {
		for i := len(as) - 1; i >= 0; i-- {
			a := as[i]
			if pred(i, a) {
				return O.Some(a)
			}
		}
		return none
	}
}

// FindLast finds the first element which satisfies a predicate (or a refinement) function
func FindLast[AS ~[]A, PRED ~func(A) bool, A any](pred PRED) func(AS) O.Option[A] {
	return FindLastWithIndex[AS](F.Ignore1of2[int](pred))
}

// FindLastMapWithIndex finds the first element returned by an [O.Option] based selector function
func FindLastMapWithIndex[AS ~[]A, PRED ~func(int, A) O.Option[B], A, B any](pred PRED) func(AS) O.Option[B] {
	none := O.None[B]()
	return func(as AS) O.Option[B] {
		for i := len(as) - 1; i >= 0; i-- {
			out := pred(i, as[i])
			if O.IsSome(out) {
				return out
			}
		}
		return none
	}
}

// FindLastMap finds the first element returned by an [O.Option] based selector function
func FindLastMap[AS ~[]A, PRED ~func(A) O.Option[B], A, B any](pred PRED) func(AS) O.Option[B] {
	return FindLastMapWithIndex[AS](F.Ignore1of2[int](pred))
}
