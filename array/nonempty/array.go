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

package nonempty

import (
	G "github.com/IBM/fp-go/array/generic"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/array"
)

// NonEmptyArray represents an array with at least one element
type NonEmptyArray[A any] []A

// Of constructs a single element array
func Of[A any](first A) NonEmptyArray[A] {
	return G.Of[NonEmptyArray[A]](first)
}

// From constructs an array from a set of variadic arguments
func From[A any](first A, data ...A) NonEmptyArray[A] {
	count := len(data)
	if count == 0 {
		return Of(first)
	}
	// allocate the requested buffer
	buffer := make(NonEmptyArray[A], count+1)
	buffer[0] = first
	copy(buffer[1:], data)
	return buffer
}

func IsEmpty[A any](as NonEmptyArray[A]) bool {
	return false
}

func IsNonEmpty[A any](as NonEmptyArray[A]) bool {
	return true
}

func MonadMap[A, B any](as NonEmptyArray[A], f func(a A) B) NonEmptyArray[B] {
	return G.MonadMap[NonEmptyArray[A], NonEmptyArray[B]](as, f)
}

func Map[A, B any](f func(a A) B) func(NonEmptyArray[A]) NonEmptyArray[B] {
	return F.Bind2nd(MonadMap[A, B], f)
}

func Reduce[A, B any](f func(B, A) B, initial B) func(NonEmptyArray[A]) B {
	return func(as NonEmptyArray[A]) B {
		return array.Reduce(as, f, initial)
	}
}

func Tail[A any](as NonEmptyArray[A]) []A {
	return as[1:]
}

func Head[A any](as NonEmptyArray[A]) A {
	return as[0]
}

func First[A any](as NonEmptyArray[A]) A {
	return as[0]
}

func Last[A any](as NonEmptyArray[A]) A {
	return as[len(as)-1]
}

func Size[A any](as NonEmptyArray[A]) int {
	return G.Size(as)
}

func Flatten[A any](mma NonEmptyArray[NonEmptyArray[A]]) NonEmptyArray[A] {
	return G.Flatten(mma)
}

func MonadChain[A, B any](fa NonEmptyArray[A], f func(a A) NonEmptyArray[B]) NonEmptyArray[B] {
	return G.MonadChain[NonEmptyArray[A], NonEmptyArray[B]](fa, f)
}

func Chain[A, B any](f func(A) NonEmptyArray[B]) func(NonEmptyArray[A]) NonEmptyArray[B] {
	return G.Chain[NonEmptyArray[A], NonEmptyArray[B]](f)
}

func MonadAp[B, A any](fab NonEmptyArray[func(A) B], fa NonEmptyArray[A]) NonEmptyArray[B] {
	return G.MonadAp[NonEmptyArray[B]](fab, fa)
}

func Ap[B, A any](fa NonEmptyArray[A]) func(NonEmptyArray[func(A) B]) NonEmptyArray[B] {
	return G.Ap[NonEmptyArray[B], NonEmptyArray[func(A) B]](fa)
}
