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
	"sync"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

// FromIO converts an [IO] to a [ReaderIO]
func FromIO[R, A any](t IO[A]) ReaderIO[R, A] {
	return reader.Of[R](t)
}

func FromReader[R, A any](r Reader[R, A]) ReaderIO[R, A] {
	return readert.MonadFromReader[Reader[R, A], ReaderIO[R, A]](io.Of[A], r)
}

func MonadMap[R, A, B any](fa ReaderIO[R, A], f func(A) B) ReaderIO[R, B] {
	return readert.MonadMap[ReaderIO[R, A], ReaderIO[R, B]](io.MonadMap[A, B], fa, f)
}

func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return readert.Map[ReaderIO[R, A], ReaderIO[R, B]](io.Map[A, B], f)
}

func MonadChain[R, A, B any](ma ReaderIO[R, A], f func(A) ReaderIO[R, B]) ReaderIO[R, B] {
	return readert.MonadChain(io.MonadChain[A, B], ma, f)
}

func Chain[R, A, B any](f func(A) ReaderIO[R, B]) Operator[R, A, B] {
	return readert.Chain[ReaderIO[R, A]](io.Chain[A, B], f)
}

func Of[R, A any](a A) ReaderIO[R, A] {
	return readert.MonadOf[ReaderIO[R, A]](io.Of[A], a)
}

func MonadAp[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadAp[A, B], fab, fa)
}

func MonadApSeq[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadApSeq[A, B], fab, fa)
}

func MonadApPar[B, R, A any](fab ReaderIO[R, func(A) B], fa ReaderIO[R, A]) ReaderIO[R, B] {
	return readert.MonadAp[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(A) B], R, A](io.MonadApPar[A, B], fab, fa)
}

func Ap[B, R, A any](fa ReaderIO[R, A]) Operator[R, func(A) B, B] {
	return function.Bind2nd(MonadAp[B, R, A], fa)
}

func Ask[R any]() ReaderIO[R, R] {
	return fromreader.Ask(FromReader[R, R])()
}

func Asks[R, A any](r Reader[R, A]) ReaderIO[R, A] {
	return fromreader.Asks(FromReader[R, A])(r)
}

func MonadChainIOK[R, A, B any](ma ReaderIO[R, A], f func(A) IO[B]) ReaderIO[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma, f,
	)
}

func ChainIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[R, A any](gen func() ReaderIO[R, A]) ReaderIO[R, A] {
	return func(r R) IO[A] {
		return func() A {
			return gen()(r)()
		}
	}
}

// Memoize computes the value of the provided [ReaderIO] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[R, A any](rdr ReaderIO[R, A]) ReaderIO[R, A] {
	// synchronization primitives
	var once sync.Once
	var result A
	// callback
	gen := func(r R) func() {
		return func() {
			result = rdr(r)()
		}
	}
	// returns our memoized wrapper
	return func(r R) IO[A] {
		io := gen(r)
		return func() A {
			once.Do(io)
			return result
		}
	}
}

func Flatten[R, A any](mma ReaderIO[R, ReaderIO[R, A]]) ReaderIO[R, A] {
	return MonadChain(mma, function.Identity[ReaderIO[R, A]])
}

func MonadFlap[R, A, B any](fab ReaderIO[R, func(A) B], a A) ReaderIO[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

func Flap[R, A, B any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}
