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
	C "github.com/IBM/fp-go/internal/chain"
	"github.com/IBM/fp-go/internal/eithert"
	FE "github.com/IBM/fp-go/internal/fromeither"
	FI "github.com/IBM/fp-go/internal/fromio"
	FC "github.com/IBM/fp-go/internal/functor"
	IO "github.com/IBM/fp-go/io/generic"
	O "github.com/IBM/fp-go/option"
)

// type IOEither[E, A any] = func() Either[E, A]

func MakeIO[GA ~func() ET.Either[E, A], E, A any](f GA) GA {
	return f
}

func Left[GA ~func() ET.Either[E, A], E, A any](l E) GA {
	return MakeIO(eithert.Left(IO.MonadOf[GA, ET.Either[E, A]], l))
}

func Right[GA ~func() ET.Either[E, A], E, A any](r A) GA {
	return MakeIO(eithert.Right(IO.MonadOf[GA, ET.Either[E, A]], r))
}

func Of[GA ~func() ET.Either[E, A], E, A any](r A) GA {
	return Right[GA](r)
}

func MonadOf[GA ~func() ET.Either[E, A], E, A any](r A) GA {
	return Of[GA](r)
}

func LeftIO[GA ~func() ET.Either[E, A], GE ~func() E, E, A any](ml GE) GA {
	return MakeIO(eithert.LeftF(IO.MonadMap[GE, GA, E, ET.Either[E, A]], ml))
}

func RightIO[GA ~func() ET.Either[E, A], GR ~func() A, E, A any](mr GR) GA {
	return MakeIO(eithert.RightF(IO.MonadMap[GR, GA, A, ET.Either[E, A]], mr))
}

func FromEither[GA ~func() ET.Either[E, A], E, A any](e ET.Either[E, A]) GA {
	return IO.Of[GA](e)
}

func FromOption[GA ~func() ET.Either[E, A], E, A any](onNone func() E) func(o O.Option[A]) GA {
	return FE.FromOption(
		FromEither[GA, E, A],
		onNone,
	)
}

func ChainOptionK[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](onNone func() E) func(func(A) O.Option[B]) func(GA) GB {
	return FE.ChainOptionK(
		MonadChain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		onNone,
	)
}

func FromIO[GA ~func() ET.Either[E, A], GR ~func() A, E, A any](mr GR) GA {
	return RightIO[GA](mr)
}

func MonadMap[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](fa GA, f func(A) B) GB {
	return eithert.MonadMap(IO.MonadMap[GA, GB, ET.Either[E, A], ET.Either[E, B]], fa, f)
}

func Map[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](f func(A) B) func(GA) GB {
	return eithert.Map(IO.Map[GA, GB, ET.Either[E, A], ET.Either[E, B]], f)
}

func MonadMapTo[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](fa GA, b B) GB {
	return MonadMap[GA, GB](fa, F.Constant1[A](b))
}

func MapTo[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](b B) func(GA) GB {
	return Map[GA, GB](F.Constant1[A](b))
}

func MonadChain[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](fa GA, f func(A) GB) GB {
	return eithert.MonadChain(IO.MonadChain[GA, GB, ET.Either[E, A], ET.Either[E, B]], IO.MonadOf[GB, ET.Either[E, B]], fa, f)
}

func Chain[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](f func(A) GB) func(GA) GB {
	return eithert.Chain(IO.Chain[GA, GB, ET.Either[E, A], ET.Either[E, B]], IO.Of[GB, ET.Either[E, B]], f)
}

func MonadChainTo[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](fa GA, fb GB) GB {
	return MonadChain(fa, F.Constant1[A](fb))
}

func ChainTo[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](fb GB) func(GA) GB {
	return Chain[GA, GB, E, A, B](F.Constant1[A](fb))
}

func MonadChainEitherK[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](ma GA, f func(A) ET.Either[E, B]) GB {
	return FE.MonadChainEitherK(
		MonadChain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		ma,
		f,
	)
}

func MonadChainIOK[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GR ~func() B, E, A, B any](ma GA, f func(A) GR) GB {
	return FI.MonadChainIOK(
		MonadChain[GA, GB, E, A, B],
		FromIO[GB, GR, E, B],
		ma,
		f,
	)
}

func ChainIOK[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], GR ~func() B, E, A, B any](f func(A) GR) func(GA) GB {
	return FI.ChainIOK(
		Chain[GA, GB, E, A, B],
		FromIO[GB, GR, E, B],
		f,
	)
}

func ChainEitherK[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](f func(A) ET.Either[E, B]) func(GA) GB {
	return FE.ChainEitherK(
		Chain[GA, GB, E, A, B],
		FromEither[GB, E, B],
		f,
	)
}

func MonadAp[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadAp[GA, GB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, A], ET.Either[E, B]],
		IO.MonadMap[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		mab, ma)
}

func Ap[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.Ap[GB, func() func(ET.Either[E, A]) ET.Either[E, B], GA, ET.Either[E, B], ET.Either[E, A]],
		IO.Map[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		ma)
}

func MonadApSeq[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadApSeq[GA, GB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, A], ET.Either[E, B]],
		IO.MonadMap[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		mab, ma)
}

func ApSeq[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.ApSeq[GB, func() func(ET.Either[E, A]) ET.Either[E, B], GA, ET.Either[E, B], ET.Either[E, A]],
		IO.Map[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		ma)
}

func MonadApPar[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](mab GAB, ma GA) GB {
	return eithert.MonadAp(
		IO.MonadApPar[GA, GB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, A], ET.Either[E, B]],
		IO.MonadMap[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		mab, ma)
}

func ApPar[GB ~func() ET.Either[E, B], GAB ~func() ET.Either[E, func(A) B], GA ~func() ET.Either[E, A], E, A, B any](ma GA) func(GAB) GB {
	return eithert.Ap(
		IO.ApPar[GB, func() func(ET.Either[E, A]) ET.Either[E, B], GA, ET.Either[E, B], ET.Either[E, A]],
		IO.Map[GAB, func() func(ET.Either[E, A]) ET.Either[E, B], ET.Either[E, func(A) B], func(ET.Either[E, A]) ET.Either[E, B]],
		ma)
}

func Flatten[GA ~func() ET.Either[E, A], GAA ~func() ET.Either[E, GA], E, A any](mma GAA) GA {
	return MonadChain(mma, F.Identity[GA])
}

func TryCatch[GA ~func() ET.Either[E, A], E, A any](f func() (A, error), onThrow func(error) E) GA {
	return MakeIO(func() ET.Either[E, A] {
		a, err := f()
		return ET.TryCatch(a, err, onThrow)
	})
}

func TryCatchError[GA ~func() ET.Either[error, A], A any](f func() (A, error)) GA {
	return MakeIO(func() ET.Either[error, A] {
		return ET.TryCatchError(f())
	})
}

// Memoize computes the value of the provided IO monad lazily but exactly once
func Memoize[GA ~func() ET.Either[E, A], E, A any](ma GA) GA {
	return IO.Memoize(ma)
}

func MonadMapLeft[GA1 ~func() ET.Either[E1, A], GA2 ~func() ET.Either[E2, A], E1, E2, A any](fa GA1, f func(E1) E2) GA2 {
	return eithert.MonadMapLeft(
		IO.MonadMap[GA1, GA2, ET.Either[E1, A], ET.Either[E2, A]],
		fa,
		f,
	)
}

func MapLeft[GA1 ~func() ET.Either[E1, A], GA2 ~func() ET.Either[E2, A], E1, E2, A any](f func(E1) E2) func(GA1) GA2 {
	return eithert.MapLeft(
		IO.Map[GA1, GA2, ET.Either[E1, A], ET.Either[E2, A]],
		f,
	)
}

// Delay creates an operation that passes in the value after some [time.Duration]
func Delay[GA ~func() ET.Either[E, A], E, A any](delay time.Duration) func(GA) GA {
	return IO.Delay[GA](delay)
}

// After creates an operation that passes after the given [time.Time]
func After[GA ~func() ET.Either[E, A], E, A any](timestamp time.Time) func(GA) GA {
	return IO.After[GA](timestamp)
}

func MonadBiMap[GA ~func() ET.Either[E1, A], GB ~func() ET.Either[E2, B], E1, E2, A, B any](fa GA, f func(E1) E2, g func(A) B) GB {
	return eithert.MonadBiMap(IO.MonadMap[GA, GB, ET.Either[E1, A], ET.Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[GA ~func() ET.Either[E1, A], GB ~func() ET.Either[E2, B], E1, E2, A, B any](f func(E1) E2, g func(A) B) func(GA) GB {
	return eithert.BiMap(IO.Map[GA, GB, ET.Either[E1, A], ET.Either[E2, B]], f, g)
}

// Fold convers an IOEither into an IO
func Fold[GA ~func() ET.Either[E, A], GB ~func() B, E, A, B any](onLeft func(E) GB, onRight func(A) GB) func(GA) GB {
	return eithert.MatchE(IO.MonadChain[GA, GB, ET.Either[E, A], B], onLeft, onRight)
}

func MonadFold[GA ~func() ET.Either[E, A], GB ~func() B, E, A, B any](ma GA, onLeft func(E) GB, onRight func(A) GB) GB {
	return eithert.FoldE(IO.MonadChain[GA, GB, ET.Either[E, A], B], ma, onLeft, onRight)
}

// GetOrElse extracts the value or maps the error
func GetOrElse[GA ~func() ET.Either[E, A], GB ~func() A, E, A any](onLeft func(E) GB) func(GA) GB {
	return eithert.GetOrElse(IO.MonadChain[GA, GB, ET.Either[E, A], A], IO.MonadOf[GB, A], onLeft)
}

// MonadChainFirst runs the monad returned by the function but returns the result of the original monad
func MonadChainFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](ma GA, f func(A) GB) GA {
	return C.MonadChainFirst(
		MonadChain[GA, GA, E, A, A],
		MonadMap[GB, GA, E, B, A],
		ma,
		f,
	)
}

// ChainFirst runs the monad returned by the function but returns the result of the original monad
func ChainFirst[GA ~func() ET.Either[E, A], GB ~func() ET.Either[E, B], E, A, B any](f func(A) GB) func(GA) GA {
	return C.ChainFirst(
		Chain[GA, GA, E, A, A],
		Map[GB, GA, E, B, A],
		f,
	)
}

// MonadChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[GA ~func() ET.Either[E, A], GIOB ~func() B, E, A, B any](first GA, f func(A) GIOB) GA {
	return FI.MonadChainFirstIOK(
		MonadChain[GA, GA, E, A, A],
		MonadMap[func() ET.Either[E, B], GA, E, B, A],
		FromIO[func() ET.Either[E, B], GIOB, E, B],
		first,
		f,
	)
}

// ChainFirstIOK runs the monad returned by the function but returns the result of the original monad
func ChainFirstIOK[GA ~func() ET.Either[E, A], GIOB ~func() B, E, A, B any](f func(A) GIOB) func(GA) GA {
	return FI.ChainFirstIOK(
		Chain[GA, GA, E, A, A],
		Map[func() ET.Either[E, B], GA, E, B, A],
		FromIO[func() ET.Either[E, B], GIOB, E, B],
		f,
	)
}

// MonadChainFirstEitherK runs the monad returned by the function but returns the result of the original monad
func MonadChainFirstEitherK[GA ~func() ET.Either[E, A], E, A, B any](first GA, f func(A) ET.Either[E, B]) GA {
	return FE.MonadChainFirstEitherK(
		MonadChain[GA, GA, E, A, A],
		MonadMap[func() ET.Either[E, B], GA, E, B, A],
		FromEither[func() ET.Either[E, B], E, B],
		first,
		f,
	)
}

// ChainFirstEitherK runs the monad returned by the function but returns the result of the original monad
func ChainFirstEitherK[GA ~func() ET.Either[E, A], E, A, B any](f func(A) ET.Either[E, B]) func(GA) GA {
	return FE.ChainFirstEitherK(
		Chain[GA, GA, E, A, A],
		Map[func() ET.Either[E, B], GA, E, B, A],
		FromEither[func() ET.Either[E, B], E, B],
		f,
	)
}

// Swap changes the order of type parameters
func Swap[GEA ~func() ET.Either[E, A], GAE ~func() ET.Either[A, E], E, A any](val GEA) GAE {
	return MonadFold(val, Right[GAE], Left[GAE])
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[GA ~func() ET.Either[E, any], IMP ~func(), E any](f IMP) GA {
	return F.Pipe2(
		f,
		IO.FromImpure[func() any, IMP],
		FromIO[GA, func() any],
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[GEA ~func() ET.Either[E, A], E, A any](gen func() GEA) GEA {
	return IO.Defer[GEA](gen)
}

func MonadAlt[LAZY ~func() GIOA, GIOA ~func() ET.Either[E, A], E, A any](first GIOA, second LAZY) GIOA {
	return eithert.MonadAlt(
		IO.Of[GIOA],
		IO.MonadChain[GIOA, GIOA],

		first,
		second,
	)
}

func Alt[LAZY ~func() GIOA, GIOA ~func() ET.Either[E, A], E, A any](second LAZY) func(GIOA) GIOA {
	return F.Bind2nd(MonadAlt[LAZY], second)
}

func MonadFlap[GEAB ~func() ET.Either[E, func(A) B], GEB ~func() ET.Either[E, B], E, B, A any](fab GEAB, a A) GEB {
	return FC.MonadFlap(MonadMap[GEAB, GEB], fab, a)
}

func Flap[GEAB ~func() ET.Either[E, func(A) B], GEB ~func() ET.Either[E, B], E, B, A any](a A) func(GEAB) GEB {
	return FC.Flap(Map[GEAB, GEB], a)
}

func ToIOOption[GA ~func() O.Option[A], GEA ~func() ET.Either[E, A], E, A any](ioe GEA) GA {
	return F.Pipe1(
		ioe,
		IO.Map[GEA, GA](ET.ToOption[E, A]),
	)
}

// OrElse returns the original IOEither if it is a Right, otherwise it applies the given function to the error and returns the result.
func OrElse[GA ~func() ET.Either[E, A], E, A any](onLeft func(E) GA) func(GA) GA {
	return eithert.OrElse(IO.MonadChain[GA, GA, ET.Either[E, A], ET.Either[E, A]], IO.Of[GA, ET.Either[E, A]], onLeft)
}

func FromIOOption[GEA ~func() ET.Either[E, A], GA ~func() O.Option[A], E, A any](onNone func() E) func(ioo GA) GEA {
	return IO.Map[GA, GEA](ET.FromOption[A](onNone))
}
