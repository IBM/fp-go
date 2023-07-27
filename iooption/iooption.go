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
	I "github.com/IBM/fp-go/io"
	G "github.com/IBM/fp-go/iooption/generic"
	O "github.com/IBM/fp-go/option"
)

// IO represents a synchronous computation that may fail
// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details
type IOOption[A any] I.IO[O.Option[A]]

func MakeIO[A any](f IOOption[A]) IOOption[A] {
	return G.MakeIO(f)
}

func Of[A any](r A) IOOption[A] {
	return G.Of[IOOption[A]](r)
}

func Some[A any](r A) IOOption[A] {
	return G.Some[IOOption[A]](r)
}

func None[A any]() IOOption[A] {
	return G.None[IOOption[A]]()
}

func MonadOf[A any](r A) IOOption[A] {
	return G.MonadOf[IOOption[A]](r)
}

func FromOption[A any](o O.Option[A]) IOOption[A] {
	return G.FromOption[IOOption[A]](o)
}

func ChainOptionK[A, B any](f func(A) O.Option[B]) func(IOOption[A]) IOOption[B] {
	return G.ChainOptionK[IOOption[A], IOOption[B]](f)
}

func MonadChainIOK[A, B any](ma IOOption[A], f func(A) I.IO[B]) IOOption[B] {
	return G.MonadChainIOK[IOOption[A], IOOption[B]](ma, f)
}

func ChainIOK[A, B any](f func(A) I.IO[B]) func(IOOption[A]) IOOption[B] {
	return G.ChainIOK[IOOption[A], IOOption[B]](f)
}

func FromIO[A any](mr I.IO[A]) IOOption[A] {
	return G.FromIO[IOOption[A]](mr)
}

func MonadMap[A, B any](fa IOOption[A], f func(A) B) IOOption[B] {
	return G.MonadMap[IOOption[A], IOOption[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(IOOption[A]) IOOption[B] {
	return G.Map[IOOption[A], IOOption[B]](f)
}

func MonadChain[A, B any](fa IOOption[A], f func(A) IOOption[B]) IOOption[B] {
	return G.MonadChain(fa, f)
}

func Chain[A, B any](f func(A) IOOption[B]) func(IOOption[A]) IOOption[B] {
	return G.Chain[IOOption[A]](f)
}

func MonadAp[B, A any](mab IOOption[func(A) B], ma IOOption[A]) IOOption[B] {
	return G.MonadAp[IOOption[B]](mab, ma)
}

func Ap[B, A any](ma IOOption[A]) func(IOOption[func(A) B]) IOOption[B] {
	return G.Ap[IOOption[B], IOOption[func(A) B]](ma)
}

func Flatten[A any](mma IOOption[IOOption[A]]) IOOption[A] {
	return G.Flatten(mma)
}

func Optionize0[A any](f func() (A, bool)) func() IOOption[A] {
	return G.Optionize0[IOOption[A]](f)
}

func Optionize1[T1, A any](f func(t1 T1) (A, bool)) func(T1) IOOption[A] {
	return G.Optionize1[IOOption[A]](f)
}

func Optionize2[T1, T2, A any](f func(t1 T1, t2 T2) (A, bool)) func(T1, T2) IOOption[A] {
	return G.Optionize2[IOOption[A]](f)
}

func Optionize3[T1, T2, T3, A any](f func(t1 T1, t2 T2, t3 T3) (A, bool)) func(T1, T2, T3) IOOption[A] {
	return G.Optionize3[IOOption[A]](f)
}

func Optionize4[T1, T2, T3, T4, A any](f func(t1 T1, t2 T2, t3 T3, t4 T4) (A, bool)) func(T1, T2, T3, T4) IOOption[A] {
	return G.Optionize4[IOOption[A]](f)
}

func Memoize[A any](ma IOOption[A]) IOOption[A] {
	return G.Memoize(ma)
}

// Fold convers an IOOption into an IO
func Fold[A, B any](onNone func() I.IO[B], onSome func(A) I.IO[B]) func(IOOption[A]) I.IO[B] {
	return G.Fold[IOOption[A]](onNone, onSome)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() IOOption[A]) IOOption[A] {
	return G.Defer[IOOption[A]](gen)
}
