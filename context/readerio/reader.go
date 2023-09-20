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
	"context"

	L "github.com/IBM/fp-go/lazy"
	R "github.com/IBM/fp-go/readerio/generic"
)

func MonadMap[A, B any](fa ReaderIO[A], f func(A) B) ReaderIO[B] {
	return R.MonadMap[ReaderIO[A], ReaderIO[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(ReaderIO[A]) ReaderIO[B] {
	return R.Map[ReaderIO[A], ReaderIO[B]](f)
}

func MonadChain[A, B any](ma ReaderIO[A], f func(A) ReaderIO[B]) ReaderIO[B] {
	return R.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) ReaderIO[B]) func(ReaderIO[A]) ReaderIO[B] {
	return R.Chain[ReaderIO[A]](f)
}

func Of[A any](a A) ReaderIO[A] {
	return R.Of[ReaderIO[A]](a)
}

func MonadAp[A, B any](fab ReaderIO[func(A) B], fa ReaderIO[A]) ReaderIO[B] {
	return R.MonadAp[ReaderIO[A], ReaderIO[B]](fab, fa)
}

func Ap[A, B any](fa ReaderIO[A]) func(ReaderIO[func(A) B]) ReaderIO[B] {
	return R.Ap[ReaderIO[A], ReaderIO[B], ReaderIO[func(A) B]](fa)
}

func Ask() ReaderIO[context.Context] {
	return R.Ask[ReaderIO[context.Context]]()
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen L.Lazy[ReaderIO[A]]) ReaderIO[A] {
	return R.Defer[ReaderIO[A]](gen)
}

// Memoize computes the value of the provided [ReaderIO] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[A any](rdr ReaderIO[A]) ReaderIO[A] {
	return R.Memoize[ReaderIO[A]](rdr)
}

// Flatten converts a nested [ReaderIO] into a [ReaderIO]
func Flatten[A any](mma ReaderIO[ReaderIO[A]]) ReaderIO[A] {
	return R.Flatten[ReaderIO[A]](mma)
}
