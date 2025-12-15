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

package iooption

import (
	"time"

	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/optiont"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
)

func Of[A any](r A) IOOption[A] {
	return optiont.Of(io.Of[Option[A]], r)
}

func Some[A any](r A) IOOption[A] {
	return Of(r)
}

func None[A any]() IOOption[A] {
	return optiont.None(io.Of[Option[A]])
}

func MonadOf[A any](r A) IOOption[A] {
	return Of(r)
}

func FromOption[A any](o Option[A]) IOOption[A] {
	return io.Of(o)
}

func ChainOptionK[A, B any](f func(A) Option[B]) Operator[A, B] {
	return optiont.ChainOptionK(
		io.Chain[Option[A], Option[B]],
		FromOption[B],
		f,
	)
}

func MonadChainIOK[A, B any](ma IOOption[A], f io.Kleisli[A, B]) IOOption[B] {
	return fromio.MonadChainIOK(
		MonadChain[A, B],
		FromIO[B],
		ma,
		f,
	)
}

func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B] {
	return fromio.ChainIOK(
		Chain[A, B],
		FromIO[B],
		f,
	)
}

func FromIO[A any](mr IO[A]) IOOption[A] {
	return optiont.OfF(io.MonadMap[A, Option[A]], mr)
}

func MonadMap[A, B any](fa IOOption[A], f func(A) B) IOOption[B] {
	return optiont.MonadMap(io.MonadMap[Option[A], Option[B]], fa, f)
}

func Map[A, B any](f func(A) B) Operator[A, B] {
	return optiont.Map(io.Map[Option[A], Option[B]], f)
}

func MonadChain[A, B any](fa IOOption[A], f Kleisli[A, B]) IOOption[B] {
	return optiont.MonadChain(io.MonadChain[Option[A], Option[B]], io.MonadOf[Option[B]], fa, f)
}

func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return optiont.Chain(io.Chain[Option[A], Option[B]], io.Of[Option[B]], f)
}

func MonadAp[B, A any](mab IOOption[func(A) B], ma IOOption[A]) IOOption[B] {
	return optiont.MonadAp(
		io.MonadAp[Option[A], Option[B]],
		io.MonadMap[Option[func(A) B], func(Option[A]) Option[B]],
		mab, ma)
}

func Ap[B, A any](ma IOOption[A]) Operator[func(A) B, B] {
	return optiont.Ap(
		io.Ap[Option[B], Option[A]],
		io.Map[Option[func(A) B], func(Option[A]) Option[B]],
		ma)
}

func ApSeq[B, A any](ma IOOption[A]) Operator[func(A) B, B] {
	return optiont.Ap(
		io.ApSeq[Option[B], Option[A]],
		io.Map[Option[func(A) B], func(Option[A]) Option[B]],
		ma)
}

func ApPar[B, A any](ma IOOption[A]) Operator[func(A) B, B] {
	return optiont.Ap(
		io.ApPar[Option[B], Option[A]],
		io.Map[Option[func(A) B], func(Option[A]) Option[B]],
		ma)
}

func Flatten[A any](mma IOOption[IOOption[A]]) IOOption[A] {
	return MonadChain(mma, function.Identity[IOOption[A]])
}

func Optionize0[A any](f func() (A, bool)) Lazy[IOOption[A]] {
	ef := option.Optionize0(f)
	return func() IOOption[A] {
		return ef
	}
}

func Optionize1[T1, A any](f func(t1 T1) (A, bool)) Kleisli[T1, A] {
	ef := option.Optionize1(f)
	return func(t1 T1) IOOption[A] {
		return func() Option[A] {
			return ef(t1)
		}
	}
}

func Optionize2[T1, T2, A any](f func(t1 T1, t2 T2) (A, bool)) func(T1, T2) IOOption[A] {
	ef := option.Optionize2(f)
	return func(t1 T1, t2 T2) IOOption[A] {
		return func() Option[A] {
			return ef(t1, t2)
		}
	}
}

func Optionize3[T1, T2, T3, A any](f func(t1 T1, t2 T2, t3 T3) (A, bool)) func(T1, T2, T3) IOOption[A] {
	ef := option.Optionize3(f)
	return func(t1 T1, t2 T2, t3 T3) IOOption[A] {
		return func() Option[A] {
			return ef(t1, t2, t3)
		}
	}
}

func Optionize4[T1, T2, T3, T4, A any](f func(t1 T1, t2 T2, t3 T3, t4 T4) (A, bool)) func(T1, T2, T3, T4) IOOption[A] {
	ef := option.Optionize4(f)
	return func(t1 T1, t2 T2, t3 T3, t4 T4) IOOption[A] {
		return func() Option[A] {
			return ef(t1, t2, t3, t4)
		}
	}
}

func Memoize[A any](ma IOOption[A]) IOOption[A] {
	return io.Memoize(ma)
}

// Fold convers an [IOOption] into an [IO]
func Fold[A, B any](onNone IO[B], onSome io.Kleisli[A, B]) func(IOOption[A]) IO[B] {
	return optiont.MatchE(io.Chain[Option[A], B], function.Constant(onNone), onSome)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[A any](gen func() IOOption[A]) IOOption[A] {
	return io.Defer(gen)
}

// FromEither converts an [Either] into an [IOOption]
func FromEither[E, A any](e Either[E, A]) IOOption[A] {
	return function.Pipe2(
		e,
		ET.ToOption[E, A],
		FromOption[A],
	)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[A any](first, second IOOption[A]) IOOption[A] {
	return optiont.MonadAlt(
		io.MonadOf[Option[A]],
		io.MonadChain[Option[A], Option[A]],

		first,
		lazy.Of(second),
	)
}

// Alt identifies an associative operation on a type constructor
func Alt[A any](second IOOption[A]) Operator[A, A] {
	return optiont.Alt(
		io.Of[Option[A]],
		io.Chain[Option[A], Option[A]],

		lazy.Of(second),
	)
}

// MonadChainFirst runs the monad returned by the function but returns the result of the original monad
func MonadChainFirst[A, B any](ma IOOption[A], f Kleisli[A, B]) IOOption[A] {
	return chain.MonadChainFirst(
		MonadChain[A, A],
		MonadMap[B, A],
		ma,
		f,
	)
}

// ChainFirst runs the monad returned by the function but returns the result of the original monad
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(
		Chain[A, A],
		Map[B, A],
		f,
	)
}

// MonadChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[A, B any](first IOOption[A], f io.Kleisli[A, B]) IOOption[A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[A, A],
		MonadMap[B, A],
		FromIO[B],
		first,
		f,
	)
}

// ChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return fromio.ChainFirstIOK(
		Chain[A, A],
		Map[B, A],
		FromIO[B],
		f,
	)
}

// Delay creates an operation that passes in the value after some delay
func Delay[A any](delay time.Duration) Operator[A, A] {
	return io.Delay[Option[A]](delay)
}

// After creates an operation that passes after the given [time.Time]
func After[A any](timestamp time.Time) Operator[A, A] {
	return io.After[Option[A]](timestamp)
}
