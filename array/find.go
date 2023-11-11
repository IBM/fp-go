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

package array

import (
	G "github.com/IBM/fp-go/array/generic"
	O "github.com/IBM/fp-go/option"
)

// FindFirst finds the first element which satisfies a predicate (or a refinement) function
func FindFirst[A any](pred func(A) bool) func([]A) O.Option[A] {
	return G.FindFirst[[]A](pred)
}

// FindFirstWithIndex finds the first element which satisfies a predicate (or a refinement) function
func FindFirstWithIndex[A any](pred func(int, A) bool) func([]A) O.Option[A] {
	return G.FindFirstWithIndex[[]A](pred)
}

// FindFirstMap finds the first element returned by an [O.Option] based selector function
func FindFirstMap[A, B any](sel func(A) O.Option[B]) func([]A) O.Option[B] {
	return G.FindFirstMap[[]A](sel)
}

// FindFirstMapWithIndex finds the first element returned by an [O.Option] based selector function
func FindFirstMapWithIndex[A, B any](sel func(int, A) O.Option[B]) func([]A) O.Option[B] {
	return G.FindFirstMapWithIndex[[]A](sel)
}

// FindLast finds the Last element which satisfies a predicate (or a refinement) function
func FindLast[A any](pred func(A) bool) func([]A) O.Option[A] {
	return G.FindLast[[]A](pred)
}

// FindLastWithIndex finds the Last element which satisfies a predicate (or a refinement) function
func FindLastWithIndex[A any](pred func(int, A) bool) func([]A) O.Option[A] {
	return G.FindLastWithIndex[[]A](pred)
}

// FindLastMap finds the Last element returned by an [O.Option] based selector function
func FindLastMap[A, B any](sel func(A) O.Option[B]) func([]A) O.Option[B] {
	return G.FindLastMap[[]A](sel)
}

// FindLastMapWithIndex finds the Last element returned by an [O.Option] based selector function
func FindLastMapWithIndex[A, B any](sel func(int, A) O.Option[B]) func([]A) O.Option[B] {
	return G.FindLastMapWithIndex[[]A](sel)
}
