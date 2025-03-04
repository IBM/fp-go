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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

type (
	ReaderIOEither[R, E, A any] = reader.Reader[R, IOE.IOEither[E, A]]

	Mapper[R, E, A, B any] = reader.Reader[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]]
)

func MonadFromReaderIO[R, E, A any](a A, f func(A) readerio.ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(
		a,
		f,
		RightReaderIO[R, E, A],
	)
}

func FromReaderIO[R, E, A any](f func(A) readerio.ReaderIO[R, A]) func(A) ReaderIOEither[R, E, A] {
	return function.Bind2nd(MonadFromReaderIO[R, E, A], f)
}

func RightReaderIO[R, E, A any](ma readerio.ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return eithert.RightF(
		readerio.MonadMap[R, A, either.Either[E, A]],
		ma,
	)
}

func LeftReaderIO[A, R, E any](me readerio.ReaderIO[R, E]) ReaderIOEither[R, E, A] {
	return eithert.LeftF(
		readerio.MonadMap[R, E, either.Either[E, A]],
		me,
	)
}

func MonadMap[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) B) ReaderIOEither[R, E, B] {
	return eithert.MonadMap(readerio.MonadMap[R, either.Either[E, A], either.Either[E, B]], fa, f)
}

func Map[R, E, A, B any](f func(A) B) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.Map(readerio.Map[R, either.Either[E, A], either.Either[E, B]], f)
}

func MonadMapTo[R, E, A, B any](fa ReaderIOEither[R, E, A], b B) ReaderIOEither[R, E, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

func MapTo[R, E, A, B any](b B) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return Map[R, E](function.Constant1[A](b))
}

func MonadChain[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, B] {
	return eithert.MonadChain(
		readerio.MonadChain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		fa,
		f)
}

func MonadChainFirst[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, A] {
	return chain.MonadChainFirst(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		fa,
		f)
}

func MonadChainEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		ma,
		f,
	)
}

func ChainEitherK[R, E, A, B any](f func(A) either.Either[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return fromeither.ChainEitherK(
		Chain[R, E, A, B],
		FromEither[R, E, B],
		f,
	)
}

func MonadChainFirstEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, A] {
	return G.MonadChainFirstEitherK[ReaderIOEither[R, E, A]](ma, f)
}

func ChainFirstEitherK[R, E, A, B any](f func(A) either.Either[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.ChainFirstEitherK[ReaderIOEither[R, E, A]](f)
}

func MonadChainReaderK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) reader.Reader[R, B]) ReaderIOEither[R, E, B] {
	return G.MonadChainReaderK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](ma, f)
}

func ChainReaderK[E, R, A, B any](f func(A) reader.Reader[R, B]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.ChainReaderK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](f)
}

func MonadChainIOEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) IOE.IOEither[E, B]) ReaderIOEither[R, E, B] {
	return G.MonadChainIOEitherK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](ma, f)
}

func ChainIOEitherK[R, E, A, B any](f func(A) IOE.IOEither[E, B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.ChainIOEitherK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](f)
}

func MonadChainIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, B] {
	return G.MonadChainIOK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](ma, f)
}

func ChainIOK[R, E, A, B any](f func(A) io.IO[B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.ChainIOK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](f)
}

func MonadChainFirstIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, A] {
	return G.MonadChainFirstIOK[ReaderIOEither[R, E, A]](ma, f)
}

func ChainFirstIOK[R, E, A, B any](f func(A) io.IO[B]) func(ma ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.ChainFirstIOK[ReaderIOEither[R, E, A]](f)
}

func ChainOptionK[R, A, B, E any](onNone func() E) func(func(A) O.Option[B]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.ChainOptionK[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](onNone)
}

func MonadAp[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.MonadAp[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]](fab, fa)
}

func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(fab ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return G.Ap[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B], ReaderIOEither[R, E, func(A) B]](fa)
}

func Chain[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return G.Chain[ReaderIOEither[R, E, A]](f)
}

func ChainFirst[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) func(fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.ChainFirst[ReaderIOEither[R, E, A]](f)
}

func Right[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return G.Right[ReaderIOEither[R, E, A]](a)
}

func Left[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return G.Left[ReaderIOEither[R, E, A]](e)
}

func ThrowError[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return G.ThrowError[ReaderIOEither[R, E, A]](e)
}

// Of returns a Reader with a fixed value
func Of[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return G.Of[ReaderIOEither[R, E, A]](a)
}

func Flatten[R, E, A any](mma ReaderIOEither[R, E, ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return G.Flatten(mma)
}

func FromEither[R, E, A any](t either.Either[E, A]) ReaderIOEither[R, E, A] {
	return G.FromEither[ReaderIOEither[R, E, A]](t)
}

func RightReader[E, R, A any](ma reader.Reader[R, A]) ReaderIOEither[R, E, A] {
	return G.RightReader[reader.Reader[R, A], ReaderIOEither[R, E, A]](ma)
}

func LeftReader[A, R, E any](ma reader.Reader[R, E]) ReaderIOEither[R, E, A] {
	return G.LeftReader[reader.Reader[R, E], ReaderIOEither[R, E, A]](ma)
}

func FromReader[E, R, A any](ma reader.Reader[R, A]) ReaderIOEither[R, E, A] {
	return G.FromReader[reader.Reader[R, A], ReaderIOEither[R, E, A]](ma)
}

func RightIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return G.RightIO[ReaderIOEither[R, E, A]](ma)
}

func LeftIO[R, A, E any](ma io.IO[E]) ReaderIOEither[R, E, A] {
	return G.LeftIO[ReaderIOEither[R, E, A]](ma)
}

func FromIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return G.FromIO[ReaderIOEither[R, E, A]](ma)
}

func FromIOEither[R, E, A any](ma IOE.IOEither[E, A]) ReaderIOEither[R, E, A] {
	return G.FromIOEither[ReaderIOEither[R, E, A]](ma)
}

func FromReaderEither[R, E, A any](ma RE.ReaderEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.FromReaderEither[RE.ReaderEither[R, E, A], ReaderIOEither[R, E, A]](ma)
}

func Ask[R, E any]() ReaderIOEither[R, E, R] {
	return G.Ask[ReaderIOEither[R, E, R]]()
}

func Asks[E, R, A any](r reader.Reader[R, A]) ReaderIOEither[R, E, A] {
	return G.Asks[reader.Reader[R, A], ReaderIOEither[R, E, A]](r)
}

func FromOption[R, A, E any](onNone func() E) func(O.Option[A]) ReaderIOEither[R, E, A] {
	return G.FromOption[ReaderIOEither[R, E, A]](onNone)
}

func FromPredicate[R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOEither[R, E, A] {
	return G.FromPredicate[ReaderIOEither[R, E, A]](pred, onFalse)
}

func Fold[R, E, A, B any](onLeft func(E) readerio.ReaderIO[R, B], onRight func(A) readerio.ReaderIO[R, B]) func(ReaderIOEither[R, E, A]) readerio.ReaderIO[R, B] {
	return G.Fold[readerio.ReaderIO[R, B], ReaderIOEither[R, E, A]](onLeft, onRight)
}

func GetOrElse[R, E, A any](onLeft func(E) readerio.ReaderIO[R, A]) func(ReaderIOEither[R, E, A]) readerio.ReaderIO[R, A] {
	return G.GetOrElse[readerio.ReaderIO[R, A], ReaderIOEither[R, E, A]](onLeft)
}

func OrElse[R, E1, A, E2 any](onLeft func(E1) ReaderIOEither[R, E2, A]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return G.OrElse[ReaderIOEither[R, E1, A]](onLeft)
}

func OrLeft[A, E1, R, E2 any](onLeft func(E1) readerio.ReaderIO[R, E2]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return G.OrLeft[ReaderIOEither[R, E1, A], readerio.ReaderIO[R, E2], ReaderIOEither[R, E2, A]](onLeft)
}

func MonadBiMap[R, E1, E2, A, B any](fa ReaderIOEither[R, E1, A], f func(E1) E2, g func(A) B) ReaderIOEither[R, E2, B] {
	return G.MonadBiMap[ReaderIOEither[R, E1, A], ReaderIOEither[R, E2, B]](fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, B] {
	return G.BiMap[ReaderIOEither[R, E1, A], ReaderIOEither[R, E2, B]](f, g)
}

// Swap changes the order of type parameters
func Swap[R, E, A any](val ReaderIOEither[R, E, A]) ReaderIOEither[R, A, E] {
	return G.Swap[ReaderIOEither[R, E, A], ReaderIOEither[R, A, E]](val)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[R, E, A any](gen L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return G.Defer[ReaderIOEither[R, E, A]](gen)
}

// TryCatch wraps a reader returning a tuple as an error into ReaderIOEither
func TryCatch[R, E, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOEither[R, E, A] {
	return G.TryCatch[ReaderIOEither[R, E, A]](f, onThrow)
}

// MonadAlt identifies an associative operation on a type constructor.
func MonadAlt[R, E, A any](first ReaderIOEither[R, E, A], second L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return G.MonadAlt(first, second)
}

// Alt identifies an associative operation on a type constructor.
func Alt[R, E, A any](second L.Lazy[ReaderIOEither[R, E, A]]) func(ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.Alt(second)
}

// Memoize computes the value of the provided [ReaderIOEither] monad lazily but exactly once
// The context used to compute the value is the context of the first call, so do not use this
// method if the value has a functional dependency on the content of the context
func Memoize[
	R, E, A any](rdr ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return G.Memoize[ReaderIOEither[R, E, A]](rdr)
}

func MonadFlap[R, E, B, A any](fab ReaderIOEither[R, E, func(A) B], a A) ReaderIOEither[R, E, B] {
	return G.MonadFlap[ReaderIOEither[R, E, func(A) B], ReaderIOEither[R, E, B]](fab, a)
}

func Flap[R, E, B, A any](a A) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return G.Flap[ReaderIOEither[R, E, func(A) B], ReaderIOEither[R, E, B]](a)
}

func MonadMapLeft[R, E1, E2, A any](fa ReaderIOEither[R, E1, A], f func(E1) E2) ReaderIOEither[R, E2, A] {
	return G.MonadMapLeft[ReaderIOEither[R, E1, A], ReaderIOEither[R, E2, A]](fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[R, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return G.MapLeft[ReaderIOEither[R, E1, A], ReaderIOEither[R, E2, A]](f)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[R1, R2, E, A any](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A] {
	return G.Local[ReaderIOEither[R1, E, A], ReaderIOEither[R2, E, A]](f)
}
