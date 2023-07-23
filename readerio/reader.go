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
	R "github.com/IBM/fp-go/reader"
	G "github.com/IBM/fp-go/readerio/generic"
)

type ReaderIO[E, A any] R.Reader[E, IO.IO[A]]

func FromIO[E, A any](t IO.IO[A]) ReaderIO[E, A] {
	return G.FromIO[ReaderIO[E, A]](t)
}

func MonadMap[E, A, B any](fa ReaderIO[E, A], f func(A) B) ReaderIO[E, B] {
	return G.MonadMap[ReaderIO[E, A], ReaderIO[E, B]](fa, f)
}

func Map[E, A, B any](f func(A) B) func(ReaderIO[E, A]) ReaderIO[E, B] {
	return G.Map[ReaderIO[E, A], ReaderIO[E, B]](f)
}

func MonadChain[E, A, B any](ma ReaderIO[E, A], f func(A) ReaderIO[E, B]) ReaderIO[E, B] {
	return G.MonadChain(ma, f)
}

func Chain[E, A, B any](f func(A) ReaderIO[E, B]) func(ReaderIO[E, A]) ReaderIO[E, B] {
	return G.Chain[ReaderIO[E, A]](f)
}

func Of[E, A any](a A) ReaderIO[E, A] {
	return G.Of[ReaderIO[E, A]](a)
}

func MonadAp[B, E, A any](fab ReaderIO[E, func(A) B], fa ReaderIO[E, A]) ReaderIO[E, B] {
	return G.MonadAp[ReaderIO[E, A], ReaderIO[E, B]](fab, fa)
}

func Ap[B, E, A any](fa ReaderIO[E, A]) func(ReaderIO[E, func(A) B]) ReaderIO[E, B] {
	return G.Ap[ReaderIO[E, A], ReaderIO[E, B], ReaderIO[E, func(A) B]](fa)
}

func Ask[E any]() ReaderIO[E, E] {
	return G.Ask[ReaderIO[E, E]]()
}

func Asks[E, A any](r R.Reader[E, A]) ReaderIO[E, A] {
	return G.Asks[R.Reader[E, A], ReaderIO[E, A]](r)
}

func MonadChainIOK[E, A, B any](ma ReaderIO[E, A], f func(A) IO.IO[B]) ReaderIO[E, B] {
	return G.MonadChainIOK[ReaderIO[E, A], ReaderIO[E, B]](ma, f)
}

func ChainIOK[E, A, B any](f func(A) IO.IO[B]) func(ReaderIO[E, A]) ReaderIO[E, B] {
	return G.ChainIOK[ReaderIO[E, A], ReaderIO[E, B]](f)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[E, A any](gen func() ReaderIO[E, A]) ReaderIO[E, A] {
	return G.Defer[ReaderIO[E, A]](gen)
}
