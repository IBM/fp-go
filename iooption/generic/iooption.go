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

package generic

import (
	"time"

	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	FI "github.com/IBM/fp-go/internal/fromio"
	"github.com/IBM/fp-go/internal/optiont"
	IO "github.com/IBM/fp-go/io/generic"
	O "github.com/IBM/fp-go/option"
)

// type IOOption[A any] = func() Option[A]

func MakeIO[GA ~func() O.Option[A], A any](f GA) GA {
	return f
}

func Of[GA ~func() O.Option[A], A any](r A) GA {
	return MakeIO(optiont.Of(IO.MonadOf[GA, O.Option[A]], r))
}

func Some[GA ~func() O.Option[A], A any](r A) GA {
	return Of[GA](r)
}

func None[GA ~func() O.Option[A], A any]() GA {
	return MakeIO(optiont.None(IO.MonadOf[GA, O.Option[A]]))
}

func MonadOf[GA ~func() O.Option[A], A any](r A) GA {
	return Of[GA](r)
}

func FromIO[GA ~func() O.Option[A], GR ~func() A, A any](mr GR) GA {
	return MakeIO(optiont.OfF(IO.MonadMap[GR, GA, A, O.Option[A]], mr))
}

func FromOption[GA ~func() O.Option[A], A any](o O.Option[A]) GA {
	return IO.Of[GA](o)
}

func FromEither[GA ~func() O.Option[A], E, A any](e ET.Either[E, A]) GA {
	return F.Pipe2(
		e,
		ET.ToOption[E, A],
		FromOption[GA],
	)
}

func MonadMap[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](fa GA, f func(A) B) GB {
	return optiont.MonadMap(IO.MonadMap[GA, GB, O.Option[A], O.Option[B]], fa, f)
}

func Map[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](f func(A) B) func(GA) GB {
	return F.Bind2nd(MonadMap[GA, GB, A, B], f)
}

func MonadChain[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](fa GA, f func(A) GB) GB {
	return optiont.MonadChain(IO.MonadChain[GA, GB, O.Option[A], O.Option[B]], IO.MonadOf[GB, O.Option[B]], fa, f)
}

func Chain[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](f func(A) GB) func(GA) GB {
	return F.Bind2nd(MonadChain[GA, GB, A, B], f)
}

func MonadChainOptionK[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](ma GA, f func(A) O.Option[B]) GB {
	return optiont.MonadChainOptionK(
		IO.MonadChain[GA, GB, O.Option[A], O.Option[B]],
		FromOption[GB, B],
		ma,
		f,
	)
}

func ChainOptionK[GA ~func() O.Option[A], GB ~func() O.Option[B], A, B any](f func(A) O.Option[B]) func(GA) GB {
	return F.Bind2nd(MonadChainOptionK[GA, GB, A, B], f)
}

func MonadChainIOK[GA ~func() O.Option[A], GB ~func() O.Option[B], GR ~func() B, A, B any](ma GA, f func(A) GR) GB {
	return FI.MonadChainIOK(
		MonadChain[GA, GB, A, B],
		FromIO[GB, GR, B],
		ma,
		f,
	)
}

func ChainIOK[GA ~func() O.Option[A], GB ~func() O.Option[B], GR ~func() B, A, B any](f func(A) GR) func(GA) GB {
	return FI.ChainIOK(
		MonadChain[GA, GB, A, B],
		FromIO[GB, GR, B],
		f,
	)
}

func MonadAp[GB ~func() O.Option[B], GAB ~func() O.Option[func(A) B], GA ~func() O.Option[A], A, B any](mab GAB, ma GA) GB {
	return optiont.MonadAp(
		IO.MonadAp[GA, GB, func() func(O.Option[A]) O.Option[B], O.Option[A], O.Option[B]],
		IO.MonadMap[GAB, func() func(O.Option[A]) O.Option[B], O.Option[func(A) B], func(O.Option[A]) O.Option[B]],
		mab, ma)
}

func Ap[GB ~func() O.Option[B], GAB ~func() O.Option[func(A) B], GA ~func() O.Option[A], A, B any](ma GA) func(GAB) GB {
	return F.Bind2nd(MonadAp[GB, GAB, GA, A, B], ma)
}

func Flatten[GA ~func() O.Option[A], GAA ~func() O.Option[GA], A any](mma GAA) GA {
	return MonadChain(mma, F.Identity[GA])
}

func Optionize0[GA ~func() O.Option[A], A any](f func() (A, bool)) func() GA {
	ef := O.Optionize0(f)
	return func() GA {
		return MakeIO[GA](ef)
	}
}

func Optionize1[GA ~func() O.Option[A], T1, A any](f func(t1 T1) (A, bool)) func(T1) GA {
	ef := O.Optionize1(f)
	return func(t1 T1) GA {
		return MakeIO[GA](func() O.Option[A] {
			return ef(t1)
		})
	}
}

func Optionize2[GA ~func() O.Option[A], T1, T2, A any](f func(t1 T1, t2 T2) (A, bool)) func(T1, T2) GA {
	ef := O.Optionize2(f)
	return func(t1 T1, t2 T2) GA {
		return MakeIO[GA](func() O.Option[A] {
			return ef(t1, t2)
		})
	}
}

func Optionize3[GA ~func() O.Option[A], T1, T2, T3, A any](f func(t1 T1, t2 T2, t3 T3) (A, bool)) func(T1, T2, T3) GA {
	ef := O.Optionize3(f)
	return func(t1 T1, t2 T2, t3 T3) GA {
		return MakeIO[GA](func() O.Option[A] {
			return ef(t1, t2, t3)
		})
	}
}

func Optionize4[GA ~func() O.Option[A], T1, T2, T3, T4, A any](f func(t1 T1, t2 T2, t3 T3, t4 T4) (A, bool)) func(T1, T2, T3, T4) GA {
	ef := O.Optionize4(f)
	return func(t1 T1, t2 T2, t3 T3, t4 T4) GA {
		return MakeIO[GA](func() O.Option[A] {
			return ef(t1, t2, t3, t4)
		})
	}
}

// Memoize computes the value of the provided IO monad lazily but exactly once
func Memoize[GA ~func() O.Option[A], A any](ma GA) GA {
	return IO.Memoize(ma)
}

// Delay creates an operation that passes in the value after some delay
func Delay[GA ~func() O.Option[A], A any](delay time.Duration) func(GA) GA {
	return IO.Delay[GA](delay)
}

// Fold convers an IOOption into an IO
func Fold[GA ~func() O.Option[A], GB ~func() B, A, B any](onNone func() GB, onSome func(A) GB) func(GA) GB {
	return optiont.MatchE(IO.MonadChain[GA, GB, O.Option[A], B], onNone, onSome)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[GA ~func() O.Option[A], A any](gen func() GA) GA {
	return IO.Defer[GA](gen)
}

func MonadAlt[LAZY ~func() GIOA, GIOA ~func() O.Option[A], A any](first GIOA, second LAZY) GIOA {
	return optiont.MonadAlt(
		IO.Of[GIOA],
		IO.MonadChain[GIOA, GIOA],

		first,
		second,
	)
}

func Alt[LAZY ~func() GIOA, GIOA ~func() O.Option[A], A any](second LAZY) func(GIOA) GIOA {
	return F.Bind2nd(MonadAlt[LAZY], second)
}
