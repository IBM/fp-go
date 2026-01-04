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
	"time"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	FE "github.com/IBM/fp-go/v2/internal/fromeither"
	FI "github.com/IBM/fp-go/v2/internal/fromio"
	FC "github.com/IBM/fp-go/v2/internal/functor"
	IO "github.com/IBM/fp-go/v2/io/generic"
	O "github.com/IBM/fp-go/v2/option"
)

// type IOEither[E, A any] = func() Either[E, A]

// Deprecated:
func MakeIO[GA ~func() either.Either[E, A], E, A any](f GA) GA {
	return f
}

// Deprecated:
func Left[GA ~func() either.Either[E, A], E, A any](l E) GA {
	return MakeIO(eithert.Left(IO.MonadOf[GA, either.Either[E, A]], l))
}

// Deprecated:
func Right[GA ~func() either.Either[E, A], E, A any](r A) GA {
	return MakeIO(eithert.Right(IO.MonadOf[GA, either.Either[E, A]], r))
}

// Deprecated:
func Of[GA ~func() either.Either[E, A], E, A any](r A) GA {
	return Right[GA](r)
}

// Deprecated:
func MonadOf[GA ~func() either.Either[E, A], E, A any](r A) GA {
	return Of[GA](r)
}

// Deprecated:
func LeftIO[GA ~func() either.Either[E, A], GE ~func() E, E, A any](ml GE) GA {
	return MakeIO(eithert.LeftF(IO.MonadMap[GE, GA, E, either.Either[E, A]], ml))
}

// Deprecated:
func RightIO[GA ~func() either.Either[E, A], GR ~func() A, E, A any](mr GR) GA {
	return MakeIO(eithert.RightF(IO.MonadMap[GR, GA, A, either.Either[E, A]], mr))
}

func FromEither[GA ~func() either.Either[E, A], E, A any](e either.Either[E, A]) GA {
	return IO.Of[GA](e)
}

func FromOption[GA ~func() either.Either[E, A], E, A any](onNone func() E) func(o O.Option[A]) GA {
	return FE.FromOption(
		FromEither[GA, E, A],
		onNone,
	)
}

func ChainOptionK[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](onNone func() E) func(func(A) O.Option[B]) func(GA) GB {
	return FE.ChainOptionK(
		MonadChain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		onNone,
	)
}

// Deprecated:
func FromIO[GA ~func() either.Either[E, A], GR ~func() A, E, A any](mr GR) GA {
	return RightIO[GA](mr)
}

// Deprecated:
func MonadMap[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](fa GA, f func(A) B) GB {
	return eithert.MonadMap(IO.MonadMap[GA, GB, either.Either[E, A], either.Either[E, B]], fa, f)
}

// Deprecated:
func Map[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](f func(A) B) func(GA) GB {
	return eithert.Map(IO.Map[GA, GB, either.Either[E, A], either.Either[E, B]], f)
}

// Deprecated:
func MonadMapTo[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](fa GA, b B) GB {
	return MonadMap[GA, GB](fa, F.Constant1[A](b))
}

// Deprecated:
func MapTo[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](b B) func(GA) GB {
	return Map[GA, GB](F.Constant1[A](b))
}

// Deprecated:
func MonadChain[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](fa GA, f func(A) GB) GB {
	return eithert.MonadChain(IO.MonadChain[GA, GB, either.Either[E, A], either.Either[E, B]], IO.MonadOf[GB, either.Either[E, B]], fa, f)
}

// Deprecated:
func Chain[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](f func(A) GB) func(GA) GB {
	return eithert.Chain(IO.Chain[GA, GB, either.Either[E, A], either.Either[E, B]], IO.Of[GB, either.Either[E, B]], f)
}

func MonadChainTo[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](fa GA, fb GB) GB {
	return MonadChain(fa, F.Constant1[A](fb))
}

func ChainTo[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](fb GB) func(GA) GB {
	return Chain[GA](F.Constant1[A](fb))
}

// Deprecated:
func MonadChainEitherK[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](ma GA, f func(A) either.Either[E, B]) GB {
	return FE.MonadChainEitherK(
		MonadChain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		ma,
		f,
	)
}

// Deprecated:
func MonadChainIOK[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], GR ~func() B, E, A, B any](ma GA, f func(A) GR) GB {
	return FI.MonadChainIOK(
		MonadChain[GA, GB, E, A, B],
		FromIO[GB, GR, E, B],
		ma,
		f,
	)
}

func ChainIOK[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], GR ~func() B, E, A, B any](f func(A) GR) func(GA) GB {
	return FI.ChainIOK(
		Chain[GA, GB, E, A, B],
		FromIO[GB, GR, E, B],
		f,
	)
}

// Deprecated:
func ChainEitherK[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](f func(A) either.Either[E, B]) func(GA) GB {
	return FE.ChainEitherK(
		Chain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		f,
	)
}

// Deprecated:
func MonadAp[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadAp[GA, GB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, A], either.Either[E, B]],
		IO.MonadMap[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		mab, ma)
}

// Deprecated:
func Ap[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.Ap[GB, func() func(either.Either[E, A]) either.Either[E, B], GA, either.Either[E, B], either.Either[E, A]],
		IO.Map[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		ma)
}

// Deprecated:
func MonadApSeq[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadApSeq[GA, GB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, A], either.Either[E, B]],
		IO.MonadMap[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		mab, ma)
}

// Deprecated:
func ApSeq[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.ApSeq[GB, func() func(either.Either[E, A]) either.Either[E, B], GA, either.Either[E, B], either.Either[E, A]],
		IO.Map[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		ma)
}

// Deprecated:
func MonadApPar[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadApPar[GA, GB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, A], either.Either[E, B]],
		IO.MonadMap[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		mab, ma)
}

// Deprecated:
func ApPar[GB ~func() either.Either[E, B], GAB ~func() either.Either[E, func(A) B], GA ~func() either.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.ApPar[GB, func() func(either.Either[E, A]) either.Either[E, B], GA, either.Either[E, B], either.Either[E, A]],
		IO.Map[GAB, func() func(either.Either[E, A]) either.Either[E, B], either.Either[E, func(A) B], func(either.Either[E, A]) either.Either[E, B]],
		ma)
}

// Deprecated:
func Flatten[GA ~func() either.Either[E, A], GAA ~func() either.Either[E, GA], E, A any](mma GAA) GA {
	return MonadChain(mma, F.Identity[GA])
}

// Deprecated:
func TryCatch[GA ~func() either.Either[E, A], E, A any](f func() (A, error), onThrow func(error) E) GA {
	return MakeIO(func() either.Either[E, A] {
		a, err := f()
		return either.TryCatch(a, err, onThrow)
	})
}

// Deprecated:
func TryCatchError[GA ~func() either.Either[error, A], A any](f func() (A, error)) GA {
	return MakeIO(func() either.Either[error, A] {
		return either.TryCatchError(f())
	})
}

// Memoize computes the value of the provided IO monad lazily but exactly once
//
// Deprecated:
func Memoize[GA ~func() either.Either[E, A], E, A any](ma GA) GA {
	return IO.Memoize(ma)
}

// Deprecated:
func MonadMapLeft[GA1 ~func() either.Either[E1, A], GA2 ~func() either.Either[E2, A], E1, E2, A any](fa GA1, f func(E1) E2) GA2 {
	return eithert.MonadMapLeft(
		IO.MonadMap[GA1, GA2, either.Either[E1, A], either.Either[E2, A]],
		fa,
		f,
	)
}

// Deprecated:
func MapLeft[GA1 ~func() either.Either[E1, A], GA2 ~func() either.Either[E2, A], E1, E2, A any](f func(E1) E2) func(GA1) GA2 {
	return eithert.MapLeft(
		IO.Map[GA1, GA2, either.Either[E1, A], either.Either[E2, A]],
		f,
	)
}

// Delay creates an operation that passes in the value after some [time.Duration]
//
// Deprecated:
func Delay[GA ~func() either.Either[E, A], E, A any](delay time.Duration) func(GA) GA {
	return IO.Delay[GA](delay)
}

// After creates an operation that passes after the given [time.Time]
//
// Deprecated:
func After[GA ~func() either.Either[E, A], E, A any](timestamp time.Time) func(GA) GA {
	return IO.After[GA](timestamp)
}

// Deprecated:
func MonadBiMap[GA ~func() either.Either[E1, A], GB ~func() either.Either[E2, B], E1, E2, A, B any](fa GA, f func(E1) E2, g func(A) B) GB {
	return eithert.MonadBiMap(IO.MonadMap[GA, GB, either.Either[E1, A], either.Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
//
// Deprecated:
func BiMap[GA ~func() either.Either[E1, A], GB ~func() either.Either[E2, B], E1, E2, A, B any](f func(E1) E2, g func(A) B) func(GA) GB {
	return eithert.BiMap(IO.Map[GA, GB, either.Either[E1, A], either.Either[E2, B]], f, g)
}

// Fold convers an IOEither into an IO
//
// Deprecated:
func Fold[GA ~func() either.Either[E, A], GB ~func() B, E, A, B any](onLeft func(E) GB, onRight func(A) GB) func(GA) GB {
	return eithert.MatchE(IO.MonadChain[GA, GB, either.Either[E, A], B], onLeft, onRight)
}

func MonadFold[GA ~func() either.Either[E, A], GB ~func() B, E, A, B any](ma GA, onLeft func(E) GB, onRight func(A) GB) GB {
	return eithert.FoldE(IO.MonadChain[GA, GB, either.Either[E, A], B], ma, onLeft, onRight)
}

// GetOrElse extracts the value or maps the error
func GetOrElse[GA ~func() either.Either[E, A], GB ~func() A, E, A any](onLeft func(E) GB) func(GA) GB {
	return eithert.GetOrElse(IO.MonadChain[GA, GB, either.Either[E, A], A], IO.MonadOf[GB, A], onLeft)
}

// MonadChainFirst runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func MonadChainFirst[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](ma GA, f func(A) GB) GA {
	return C.MonadChainFirst(
		MonadChain[GA, GA, E, A, A],
		MonadMap[GB, GA, E, B, A],
		ma,
		f,
	)
}

// ChainFirst runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func ChainFirst[GA ~func() either.Either[E, A], GB ~func() either.Either[E, B], E, A, B any](f func(A) GB) func(GA) GA {
	return C.ChainFirst(
		Chain[GA, GA, E, A, A],
		Map[GB, GA, E, B, A],
		f,
	)
}

// MonadChainFirstIOK runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func MonadChainFirstIOK[GA ~func() either.Either[E, A], GIOB ~func() B, E, A, B any](first GA, f func(A) GIOB) GA {
	return FI.MonadChainFirstIOK(
		MonadChain[GA, GA, E, A, A],
		MonadMap[func() either.Either[E, B], GA, E, B, A],
		FromIO[func() either.Either[E, B], GIOB, E, B],
		first,
		f,
	)
}

// ChainFirstIOK runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func ChainFirstIOK[GA ~func() either.Either[E, A], GIOB ~func() B, E, A, B any](f func(A) GIOB) func(GA) GA {
	return FI.ChainFirstIOK(
		Chain[GA, GA, E, A, A],
		Map[func() either.Either[E, B], GA, E, B, A],
		FromIO[func() either.Either[E, B], GIOB, E, B],
		f,
	)
}

// MonadChainFirstEitherK runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func MonadChainFirstEitherK[GA ~func() either.Either[E, A], E, A, B any](first GA, f func(A) either.Either[E, B]) GA {
	return FE.MonadChainFirstEitherK(
		MonadChain[GA, GA, E, A, A],
		MonadMap[func() either.Either[E, B], GA, E, B, A],
		FromEither[func() either.Either[E, B], E, B],
		first,
		f,
	)
}

// ChainFirstEitherK runs the monad returned by the function but returns the result of the original monad
//
// Deprecated:
func ChainFirstEitherK[GA ~func() either.Either[E, A], E, A, B any](f func(A) either.Either[E, B]) func(GA) GA {
	return FE.ChainFirstEitherK(
		Chain[GA, GA, E, A, A],
		Map[func() either.Either[E, B], GA, E, B, A],
		FromEither[func() either.Either[E, B], E, B],
		f,
	)
}

// Swap changes the order of type parameters
//
// Deprecated:
func Swap[GEA ~func() either.Either[E, A], GAE ~func() either.Either[A, E], E, A any](val GEA) GAE {
	return MonadFold(val, Right[GAE], Left[GAE])
}

// FromImpure converts a side effect without a return value into a side effect that returns any
//
// Deprecated:
func FromImpure[GA ~func() either.Either[E, any], IMP ~func(), E any](f IMP) GA {
	return F.Pipe2(
		f,
		IO.FromImpure[func() any, IMP],
		FromIO[GA, func() any],
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
//
// Deprecated:
func Defer[GEA ~func() either.Either[E, A], E, A any](gen func() GEA) GEA {
	return IO.Defer(gen)
}

// Deprecated:
func MonadAlt[LAZY ~func() GIOA, GIOA ~func() either.Either[E, A], E, A any](first GIOA, second LAZY) GIOA {
	return eithert.MonadAlt(
		IO.Of[GIOA],
		IO.MonadChain[GIOA, GIOA],

		first,
		second,
	)
}

// Deprecated:
func Alt[LAZY ~func() GIOA, GIOA ~func() either.Either[E, A], E, A any](second LAZY) func(GIOA) GIOA {
	return F.Bind2nd(MonadAlt[LAZY], second)
}

// Deprecated:
func MonadFlap[GEAB ~func() either.Either[E, func(A) B], GEB ~func() either.Either[E, B], E, B, A any](fab GEAB, a A) GEB {
	return FC.MonadFlap(MonadMap[GEAB, GEB], fab, a)
}

// Deprecated:
func Flap[GEAB ~func() either.Either[E, func(A) B], GEB ~func() either.Either[E, B], E, B, A any](a A) func(GEAB) GEB {
	return FC.Flap(Map[GEAB, GEB], a)
}

// Deprecated:
func ToIOOption[GA ~func() O.Option[A], GEA ~func() either.Either[E, A], E, A any](ioe GEA) GA {
	return F.Pipe1(
		ioe,
		IO.Map[GEA, GA](either.ToOption[E, A]),
	)
}

// Deprecated:
func FromIOOption[GEA ~func() either.Either[E, A], GA ~func() O.Option[A], E, A any](onNone func() E) func(ioo GA) GEA {
	return IO.Map[GA, GEA](either.FromOption[A](onNone))
}
