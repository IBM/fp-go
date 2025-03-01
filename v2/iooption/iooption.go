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
	"time"

	ET "github.com/IBM/fp-go/v2/either"
	I "github.com/IBM/fp-go/v2/io"
	IO "github.com/IBM/fp-go/v2/io"
	G "github.com/IBM/fp-go/v2/iooption/generic"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
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

// Fold convers an [IOOption] into an [IO]
func Fold[A, B any](onNone func() I.IO[B], onSome func(A) I.IO[B]) func(IOOption[A]) I.IO[B] {
	return G.Fold[IOOption[A]](onNone, onSome)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() IOOption[A]) IOOption[A] {
	return G.Defer[IOOption[A]](gen)
}

// FromEither converts an [Either] into an [IOOption]
func FromEither[E, A any](e ET.Either[E, A]) IOOption[A] {
	return G.FromEither[IOOption[A]](e)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[A any](first IOOption[A], second L.Lazy[IOOption[A]]) IOOption[A] {
	return G.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second L.Lazy[IOOption[A]]) func(IOOption[A]) IOOption[A] {
	return G.Alt(second)
}

// MonadChainFirst runs the monad returned by the function but returns the result of the original monad
func MonadChainFirst[A, B any](ma IOOption[A], f func(A) IOOption[B]) IOOption[A] {
	return G.MonadChainFirst[IOOption[A], IOOption[B]](ma, f)
}

// ChainFirst runs the monad returned by the function but returns the result of the original monad
func ChainFirst[A, B any](f func(A) IOOption[B]) func(IOOption[A]) IOOption[A] {
	return G.ChainFirst[IOOption[A], IOOption[B]](f)
}

// MonadChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[A, B any](first IOOption[A], f func(A) IO.IO[B]) IOOption[A] {
	return G.MonadChainFirstIOK[IOOption[A], IO.IO[B]](first, f)
}

// ChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func ChainFirstIOK[A, B any](f func(A) IO.IO[B]) func(IOOption[A]) IOOption[A] {
	return G.ChainFirstIOK[IOOption[A], IO.IO[B]](f)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) func(IOOption[A]) IOOption[A] {
	return G.Delay[IOOption[A]](delay)
}

// After creates an operation that passes after the given [time.Time]
func After[A any](timestamp time.Time) func(IOOption[A]) IOOption[A] {
	return G.After[IOOption[A]](timestamp)
}
