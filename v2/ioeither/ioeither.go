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

package ioeither

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/file"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	IOO "github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	R "github.com/IBM/fp-go/v2/reader"
)

type (
	IO[A any]        = io.IO[A]
	Either[E, A any] = either.Either[E, A]

	// IOEither represents a synchronous computation that may fail
	// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details
	IOEither[E, A any] = IO[Either[E, A]]

	Mapper[E, A, B any] = R.Reader[IOEither[E, A], IOEither[E, B]]
)

func Left[A, E any](l E) IOEither[E, A] {
	return eithert.Left(io.MonadOf[Either[E, A]], l)
}

func Right[E, A any](r A) IOEither[E, A] {
	return eithert.Right(io.MonadOf[Either[E, A]], r)
}

func Of[E, A any](r A) IOEither[E, A] {
	return Right[E](r)
}

func MonadOf[E, A any](r A) IOEither[E, A] {
	return Of[E](r)
}

func LeftIO[A, E any](ml IO[E]) IOEither[E, A] {
	return eithert.LeftF(io.MonadMap[E, Either[E, A]], ml)
}

func RightIO[E, A any](mr IO[A]) IOEither[E, A] {
	return eithert.RightF(io.MonadMap[A, Either[E, A]], mr)
}

func FromEither[E, A any](e Either[E, A]) IOEither[E, A] {
	return io.Of(e)
}

func FromOption[A, E any](onNone func() E) func(o O.Option[A]) IOEither[E, A] {
	return fromeither.FromOption(
		FromEither[E, A],
		onNone,
	)
}

func FromIOOption[A, E any](onNone func() E) func(o IOO.IOOption[A]) IOEither[E, A] {
	return io.Map(either.FromOption[A](onNone))
}

func ChainOptionK[A, B, E any](onNone func() E) func(func(A) O.Option[B]) func(IOEither[E, A]) IOEither[E, B] {
	return fromeither.ChainOptionK(
		MonadChain[E, A, B],
		FromEither[E, B],
		onNone,
	)
}

func MonadChainIOK[E, A, B any](ma IOEither[E, A], f func(A) IO[B]) IOEither[E, B] {
	return fromio.MonadChainIOK(
		MonadChain[E, A, B],
		FromIO[E, B],
		ma,
		f,
	)
}

func ChainIOK[E, A, B any](f func(A) IO[B]) func(IOEither[E, A]) IOEither[E, B] {
	return fromio.ChainIOK(
		Chain[E, A, B],
		FromIO[E, B],
		f,
	)
}

func ChainLazyK[E, A, B any](f func(A) lazy.Lazy[B]) func(IOEither[E, A]) IOEither[E, B] {
	return ChainIOK[E](f)
}

// FromIO creates an [IOEither] from an [IO] instance, invoking [IO] for each invocation of [IOEither]
func FromIO[E, A any](mr IO[A]) IOEither[E, A] {
	return RightIO[E](mr)
}

// FromLazy creates an [IOEither] from a [Lazy] instance, invoking [Lazy] for each invocation of [IOEither]
func FromLazy[E, A any](mr lazy.Lazy[A]) IOEither[E, A] {
	return FromIO[E](mr)
}

func MonadMap[E, A, B any](fa IOEither[E, A], f func(A) B) IOEither[E, B] {
	return eithert.MonadMap(io.MonadMap[Either[E, A], Either[E, B]], fa, f)
}

func Map[E, A, B any](f func(A) B) Mapper[E, A, B] {
	return eithert.Map(io.Map[Either[E, A], Either[E, B]], f)
}

func MonadMapTo[E, A, B any](fa IOEither[E, A], b B) IOEither[E, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

func MapTo[E, A, B any](b B) Mapper[E, A, B] {
	return Map[E](function.Constant1[A](b))
}

func MonadChain[E, A, B any](fa IOEither[E, A], f func(A) IOEither[E, B]) IOEither[E, B] {
	return eithert.MonadChain(io.MonadChain[Either[E, A], Either[E, B]], io.MonadOf[Either[E, B]], fa, f)
}

func Chain[E, A, B any](f func(A) IOEither[E, B]) Mapper[E, A, B] {
	return eithert.Chain(io.Chain[Either[E, A], Either[E, B]], io.Of[Either[E, B]], f)
}

func MonadChainEitherK[E, A, B any](ma IOEither[E, A], f func(A) Either[E, B]) IOEither[E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[E, A, B],
		FromEither[E, B],
		ma,
		f,
	)
}

func ChainEitherK[E, A, B any](f func(A) Either[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return fromeither.ChainEitherK(
		Chain[E, A, B],
		FromEither[E, B],
		f,
	)
}

func MonadAp[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return eithert.MonadAp(
		io.MonadAp[Either[E, A], Either[E, B]],
		io.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// Ap is an alias of [ApPar]
func Ap[B, E, A any](ma IOEither[E, A]) Mapper[E, func(A) B, B] {
	return eithert.Ap(
		io.Ap[Either[E, B], Either[E, A]],
		io.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

func MonadApPar[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return eithert.MonadAp(
		io.MonadApPar[Either[E, A], Either[E, B]],
		io.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// ApPar applies function and value in parallel
func ApPar[B, E, A any](ma IOEither[E, A]) Mapper[E, func(A) B, B] {
	return eithert.Ap(
		io.ApPar[Either[E, B], Either[E, A]],
		io.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

func MonadApSeq[B, E, A any](mab IOEither[E, func(A) B], ma IOEither[E, A]) IOEither[E, B] {
	return eithert.MonadAp(
		io.MonadApSeq[Either[E, A], Either[E, B]],
		io.MonadMap[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		mab, ma)
}

// ApSeq applies function and value sequentially
func ApSeq[B, E, A any](ma IOEither[E, A]) func(IOEither[E, func(A) B]) IOEither[E, B] {
	return eithert.Ap(
		io.ApSeq[Either[E, B], Either[E, A]],
		io.Map[Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		ma)
}

func Flatten[E, A any](mma IOEither[E, IOEither[E, A]]) IOEither[E, A] {
	return MonadChain(mma, function.Identity[IOEither[E, A]])
}

func TryCatch[E, A any](f func() (A, error), onThrow func(error) E) IOEither[E, A] {
	return func() Either[E, A] {
		a, err := f()
		return either.TryCatch(a, err, onThrow)
	}
}

func TryCatchError[A any](f func() (A, error)) IOEither[error, A] {
	return func() Either[error, A] {
		return either.TryCatchError(f())
	}
}

func Memoize[E, A any](ma IOEither[E, A]) IOEither[E, A] {
	return io.Memoize(ma)
}

func MonadMapLeft[E1, E2, A any](fa IOEither[E1, A], f func(E1) E2) IOEither[E2, A] {
	return eithert.MonadMapLeft(
		io.MonadMap[Either[E1, A], Either[E2, A]],
		fa,
		f,
	)
}

func MapLeft[A, E1, E2 any](f func(E1) E2) func(IOEither[E1, A]) IOEither[E2, A] {
	return eithert.MapLeft(
		io.Map[Either[E1, A], Either[E2, A]],
		f,
	)
}

func MonadBiMap[E1, E2, A, B any](fa IOEither[E1, A], f func(E1) E2, g func(A) B) IOEither[E2, B] {
	return eithert.MonadBiMap(io.MonadMap[Either[E1, A], Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E1, E2, A, B any](f func(E1) E2, g func(A) B) func(IOEither[E1, A]) IOEither[E2, B] {
	return eithert.BiMap(io.Map[Either[E1, A], Either[E2, B]], f, g)
}

// Fold converts an IOEither into an IO
func Fold[E, A, B any](onLeft func(E) IO[B], onRight func(A) IO[B]) func(IOEither[E, A]) IO[B] {
	return eithert.MatchE(io.MonadChain[Either[E, A], B], onLeft, onRight)
}

// GetOrElse extracts the value or maps the error
func GetOrElse[E, A any](onLeft func(E) IO[A]) func(IOEither[E, A]) IO[A] {
	return eithert.GetOrElse(io.MonadChain[Either[E, A], A], io.MonadOf[A], onLeft)
}

// MonadChainTo composes to the second monad ignoring the return value of the first
func MonadChainTo[A, E, B any](fa IOEither[E, A], fb IOEither[E, B]) IOEither[E, B] {
	return MonadChain(fa, function.Constant1[A](fb))
}

// ChainTo composes to the second [IOEither] monad ignoring the return value of the first
func ChainTo[A, E, B any](fb IOEither[E, B]) Mapper[E, A, B] {
	return Chain(function.Constant1[A](fb))
}

// MonadChainFirst runs the [IOEither] monad returned by the function but returns the result of the original monad
func MonadChainFirst[E, A, B any](ma IOEither[E, A], f func(A) IOEither[E, B]) IOEither[E, A] {
	return chain.MonadChainFirst(
		MonadChain[E, A, A],
		MonadMap[E, B, A],
		ma,
		f,
	)
}

// ChainFirst runs the [IOEither] monad returned by the function but returns the result of the original monad
func ChainFirst[E, A, B any](f func(A) IOEither[E, B]) Mapper[E, A, A] {
	return chain.ChainFirst(
		Chain[E, A, A],
		Map[E, B, A],
		f,
	)
}

func MonadChainFirstEitherK[A, E, B any](ma IOEither[E, A], f func(A) Either[E, B]) IOEither[E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[E, A, A],
		MonadMap[E, B, A],
		FromEither[E, B],
		ma,
		f,
	)
}

func ChainFirstEitherK[A, E, B any](f func(A) Either[E, B]) Mapper[E, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[E, A, A],
		Map[E, B, A],
		FromEither[E, B],
		f,
	)
}

// MonadChainFirstIOK runs [IO] the monad returned by the function but returns the result of the original monad
func MonadChainFirstIOK[E, A, B any](ma IOEither[E, A], f func(A) IO[B]) IOEither[E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[E, A, A],
		MonadMap[E, B, A],
		FromIO[E, B],
		ma,
		f,
	)
}

// ChainFirstIOK runs the [IO] monad returned by the function but returns the result of the original monad
func ChainFirstIOK[E, A, B any](f func(A) IO[B]) func(IOEither[E, A]) IOEither[E, A] {
	return fromio.ChainFirstIOK(
		Chain[E, A, A],
		Map[E, B, A],
		FromIO[E, B],
		f,
	)
}

func MonadFold[E, A, B any](ma IOEither[E, A], onLeft func(E) IO[B], onRight func(A) IO[B]) IO[B] {
	return eithert.FoldE(io.MonadChain[Either[E, A], B], ma, onLeft, onRight)
}

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[A, E, R, ANY any](onCreate IOEither[E, R], onRelease func(R) IOEither[E, ANY]) func(func(R) IOEither[E, A]) IOEither[E, A] {
	return file.WithResource(
		MonadChain[E, R, A],
		MonadFold[E, A, Either[E, A]],
		MonadFold[E, ANY, Either[E, A]],
		MonadMap[E, ANY, A],
		Left[A, E],
	)(function.Constant(onCreate), onRelease)
}

// Swap changes the order of type parameters
func Swap[E, A any](val IOEither[E, A]) IOEither[A, E] {
	return MonadFold(val, Right[A, E], Left[E, A])
}

// FromImpure converts a side effect without a return value into a side effect that returns any
func FromImpure[E any](f func()) IOEither[E, any] {
	return function.Pipe2(
		f,
		io.FromImpure,
		FromIO[E, any],
	)
}

// Defer creates an IO by creating a brand new IO via a generator function, each time
func Defer[E, A any](gen lazy.Lazy[IOEither[E, A]]) IOEither[E, A] {
	return io.Defer(gen)
}

// MonadAlt identifies an associative operation on a type constructor
func MonadAlt[E, A any](first IOEither[E, A], second lazy.Lazy[IOEither[E, A]]) IOEither[E, A] {
	return eithert.MonadAlt(
		io.Of[Either[E, A]],
		io.MonadChain[Either[E, A], Either[E, A]],

		first,
		second,
	)
}

// Alt identifies an associative operation on a type constructor
func Alt[E, A any](second lazy.Lazy[IOEither[E, A]]) Mapper[E, A, A] {
	return function.Bind2nd(MonadAlt[E, A], second)
}

func MonadFlap[E, B, A any](fab IOEither[E, func(A) B], a A) IOEither[E, B] {
	return functor.MonadFlap(MonadMap[E, func(A) B, B], fab, a)
}

func Flap[E, B, A any](a A) Mapper[E, func(A) B, B] {
	return functor.Flap(Map[E, func(A) B, B], a)
}

// ToIOOption converts an [IOEither] to an [IOO.IOOption]
func ToIOOption[E, A any](ioe IOEither[E, A]) IOO.IOOption[A] {
	return function.Pipe1(
		ioe,
		io.Map(either.ToOption[E, A]),
	)
}

// Delay creates an operation that passes in the value after some delay
func Delay[E, A any](delay time.Duration) Mapper[E, A, A] {
	return io.Delay[Either[E, A]](delay)
}

// After creates an operation that passes after the given [time.Time]
func After[E, A any](timestamp time.Time) Mapper[E, A, A] {
	return io.After[Either[E, A]](timestamp)
}
