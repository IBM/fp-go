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

package readerioeither

import (
	"time"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
)

//go:inline
func FromReaderOption[R, A, E any](onNone Lazy[E]) Kleisli[R, E, ReaderOption[R, A], A] {
	return function.Bind2nd(function.Flow2[ReaderOption[R, A], IOE.Kleisli[E, Option[A], A]], IOE.FromOption[A](onNone))
}

//go:inline
func FromReaderIO[E, R, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return RightReaderIO[E](ma)
}

// RightReaderIO lifts a ReaderIO into a ReaderIOEither, placing the result in the Right side.
//
//go:inline
func RightReaderIO[E, R, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return eithert.RightF(
		readerio.MonadMap[R, A, either.Either[E, A]],
		ma,
	)
}

// LeftReaderIO lifts a ReaderIO into a ReaderIOEither, placing the result in the Left (error) side.
//
//go:inline
func LeftReaderIO[A, R, E any](me ReaderIO[R, E]) ReaderIOEither[R, E, A] {
	return eithert.LeftF(
		readerio.MonadMap[R, E, either.Either[E, A]],
		me,
	)
}

// MonadMap applies a function to the value inside a ReaderIOEither context.
// If the computation is successful (Right), the function is applied to the value.
// If it's an error (Left), the error is propagated unchanged.
//
//go:inline
func MonadMap[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) B) ReaderIOEither[R, E, B] {
	return eithert.MonadMap(readerio.MonadMap[R, either.Either[E, A], either.Either[E, B]], fa, f)
}

// Map returns a function that applies a transformation to the success value of a ReaderIOEither.
// This is the curried version of MonadMap, useful for function composition.
//
//go:inline
func Map[R, E, A, B any](f func(A) B) Operator[R, E, A, B] {
	return eithert.Map(readerio.Map[R, either.Either[E, A], either.Either[E, B]], f)
}

// MonadMapTo replaces the success value with a constant value.
// Useful when you want to discard the result but keep the effect.
func MonadMapTo[R, E, A, B any](fa ReaderIOEither[R, E, A], b B) ReaderIOEither[R, E, B] {
	return MonadMap(fa, function.Constant1[A](b))
}

// MapTo returns a function that replaces the success value with a constant.
// This is the curried version of MonadMapTo.
func MapTo[R, E, A, B any](b B) Operator[R, E, A, B] {
	return Map[R, E](function.Constant1[A](b))
}

// MonadChain sequences two computations where the second depends on the result of the first.
// This is the fundamental operation for composing dependent effectful computations.
// If the first computation fails, the second is not executed.
//
//go:inline
func MonadChain[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, B] {
	return eithert.MonadChain(
		readerio.MonadChain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		fa,
		f)
}

// MonadChainFirst sequences two computations but keeps the result of the first.
// Useful for performing side effects while preserving the original value.
//
//go:inline
func MonadChainFirst[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, A] {
	return chain.MonadChainFirst(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		fa,
		f)
}

//go:inline
func MonadTap[R, E, A, B any](fa ReaderIOEither[R, E, A], f Kleisli[R, E, A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirst(fa, f)
}

// MonadChainEitherK chains a computation that returns an Either into a ReaderIOEither.
// The Either is automatically lifted into the ReaderIOEither context.
//
//go:inline
func MonadChainEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		ma,
		f,
	)
}

// ChainEitherK returns a function that chains an Either-returning function into ReaderIOEither.
// This is the curried version of MonadChainEitherK.
//
//go:inline
func ChainEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, E, A, B],
		FromEither[R, E, B],
		f,
	)
}

// MonadChainFirstEitherK chains an Either-returning computation but keeps the original value.
// Useful for validation or side effects that return Either.
//
//go:inline
func MonadChainFirstEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromEither[R, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f either.Kleisli[E, A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirstEitherK(ma, f)
}

// ChainFirstEitherK returns a function that chains an Either computation while preserving the original value.
// This is the curried version of MonadChainFirstEitherK.
//
//go:inline
func ChainFirstEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromEither[R, E, B],
		f,
	)
}

//go:inline
func TapEitherK[R, E, A, B any](f either.Kleisli[E, A, B]) Operator[R, E, A, A] {
	return ChainFirstEitherK[R](f)
}

// MonadChainReaderK chains a Reader-returning computation into a ReaderIOEither.
// The Reader is automatically lifted into the ReaderIOEither context.
//
//go:inline
func MonadChainReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, E, A, B],
		FromReader[E, R, B],
		ma,
		f,
	)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOEither.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, E, A, B],
		FromReader[E, R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, E, A, B],
		FromReader[E, R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderK[E, R, A, B any](ma ReaderIOEither[R, E, A], f reader.Kleisli[R, A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirstReaderK(ma, f)
}

// ChainFirstReaderK returns a function that chains a Reader-returning function into ReaderIOEither
// while preserving the original value. This is the curried version of MonadChainFirstReaderK.
//
//go:inline
func ChainFirstReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, E, A, B],
		FromReader[E, R, B],
		f,
	)
}

//go:inline
func TapReaderK[E, R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, E, A, A] {
	return ChainFirstReaderK[E](f)
}

//go:inline
func MonadChainReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, E, A, B],
		FromReaderIO[E, R, B],
		ma,
		f,
	)
}

//go:inline
func ChainReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, E, A, B],
		FromReaderIO[E, R, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, E, A, B],
		FromReaderIO[E, R, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderIOK[E, R, A, B any](ma ReaderIOEither[R, E, A], f readerio.Kleisli[R, A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, E, A, B],
		FromReaderIO[E, R, B],
		f,
	)
}

//go:inline
func TapReaderIOK[E, R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, E, A, A] {
	return ChainFirstReaderIOK[E](f)
}

//go:inline
func MonadChainReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, E, A, B],
		FromReaderEither[R, E, B],
		ma,
		f,
	)
}

// ChainReaderEitherK returns a function that chains a ReaderEither-returning function into ReaderIOEither.
// This is the curried version of MonadChainReaderEitherK.
//
//go:inline
func ChainReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, E, A, B],
		FromReaderEither[R, E, B],
		f,
	)
}

//go:inline
func MonadChainFirstReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, E, A, B],
		FromReaderEither[R, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapReaderEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f RE.Kleisli[R, E, A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirstReaderEitherK(ma, f)
}

// ChainFirstReaderEitherK returns a function that chains a ReaderEither-returning function into ReaderIOEither
// while preserving the original value. This is the curried version of MonadChainFirstReaderEitherK.
//
//go:inline
func ChainFirstReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, E, A, B],
		FromReaderEither[R, E, B],
		f,
	)
}

//go:inline
func TapReaderEitherK[E, R, A, B any](f RE.Kleisli[R, E, A, B]) Operator[R, E, A, A] {
	return ChainFirstReaderEitherK(f)
}

//go:inline
func ChainReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, B] {
	fro := FromReaderOption[R, B](onNone)
	return func(f readeroption.Kleisli[R, A, B]) Operator[R, E, A, B] {
		return fromreader.ChainReaderK(
			Chain[R, E, A, B],
			fro,
			f,
		)
	}
}

//go:inline
func ChainFirstReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, A] {
	fro := FromReaderOption[R, B](onNone)
	return func(f readeroption.Kleisli[R, A, B]) Operator[R, E, A, A] {
		return fromreader.ChainFirstReaderK(
			ChainFirst[R, E, A, B],
			fro,
			f,
		)
	}
}

//go:inline
func TapReaderOptionK[R, A, B, E any](onNone Lazy[E]) func(readeroption.Kleisli[R, A, B]) Operator[R, E, A, A] {
	return ChainFirstReaderOptionK[R, A, B](onNone)
}

// MonadChainIOEitherK chains an IOEither-returning computation into a ReaderIOEither.
// The IOEither is automatically lifted into the ReaderIOEither context.
//
//go:inline
func MonadChainIOEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f IOE.Kleisli[E, A, B]) ReaderIOEither[R, E, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, E, A, B],
		FromIOEither[R, E, B],
		ma,
		f,
	)
}

// ChainIOEitherK returns a function that chains an IOEither-returning function into ReaderIOEither.
// This is the curried version of MonadChainIOEitherK.
//
//go:inline
func ChainIOEitherK[R, E, A, B any](f IOE.Kleisli[E, A, B]) Operator[R, E, A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, E, A, B],
		FromIOEither[R, E, B],
		f,
	)
}

// MonadChainIOK chains an IO-returning computation into a ReaderIOEither.
// The IO is automatically lifted into the ReaderIOEither context (always succeeds).
//
//go:inline
func MonadChainIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, E, A, B],
		FromIO[R, E, B],
		ma,
		f,
	)
}

// ChainIOK returns a function that chains an IO-returning function into ReaderIOEither.
// This is the curried version of MonadChainIOK.
//
//go:inline
func ChainIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, B] {
	return fromio.ChainIOK(
		Chain[R, E, A, B],
		FromIO[R, E, B],
		f,
	)
}

// MonadChainFirstIOK chains an IO computation but keeps the original value.
// Useful for performing IO side effects while preserving the original value.
//
//go:inline
func MonadChainFirstIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromIO[R, E, B],
		ma,
		f,
	)
}

//go:inline
func MonadTapIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f io.Kleisli[A, B]) ReaderIOEither[R, E, A] {
	return MonadChainFirstIOK(ma, f)
}

// ChainFirstIOK returns a function that chains an IO computation while preserving the original value.
// This is the curried version of MonadChainFirstIOK.
//
//go:inline
func ChainFirstIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromIO[R, E, B],
		f,
	)
}

//go:inline
func TapIOK[R, E, A, B any](f io.Kleisli[A, B]) Operator[R, E, A, A] {
	return ChainFirstIOK[R, E](f)
}

// ChainOptionK returns a function that chains an Option-returning function into ReaderIOEither.
// If the Option is None, the provided error function is called to produce the error value.
//
//go:inline
func ChainOptionK[R, A, B, E any](onNone Lazy[E]) func(func(A) Option[B]) Operator[R, E, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		onNone,
	)
}

// MonadAp applies a function wrapped in a context to a value wrapped in a context.
// Both computations are executed (default behavior may be sequential or parallel depending on implementation).
//
//go:inline
func MonadAp[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadAp[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

// MonadApSeq applies a function in a context to a value in a context, executing them sequentially.
//
//go:inline
func MonadApSeq[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadApSeq[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

// MonadApPar applies a function in a context to a value in a context, executing them in parallel.
//
//go:inline
func MonadApPar[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadApPar[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

// Ap returns a function that applies a function in a context to a value in a context.
// This is the curried version of MonadAp.
func Ap[B, R, E, A any](fa ReaderIOEither[R, E, A]) func(fab ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return function.Bind2nd(MonadAp[R, E, A, B], fa)
}

// Chain returns a function that sequences computations where the second depends on the first.
// This is the curried version of MonadChain.
//
//go:inline
func Chain[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, B] {
	return eithert.Chain(
		readerio.Chain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		f)
}

// ChainFirst returns a function that sequences computations but keeps the first result.
// This is the curried version of MonadChainFirst.
//
//go:inline
func ChainFirst[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, A] {
	return chain.ChainFirst(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		f)
}

//go:inline
func Tap[R, E, A, B any](f Kleisli[R, E, A, B]) Operator[R, E, A, A] {
	return ChainFirst(f)
}

// Right creates a successful ReaderIOEither with the given value.
//
//go:inline
func Right[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return eithert.Right(readerio.Of[R, Either[E, A]], a)
}

// Left creates a failed ReaderIOEither with the given error.
//
//go:inline
func Left[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return eithert.Left(readerio.Of[R, Either[E, A]], e)
}

// ThrowError creates a failed ReaderIOEither with the given error.
// This is an alias for Left, following the naming convention from other functional libraries.
func ThrowError[R, A, E any](e E) ReaderIOEither[R, E, A] {
	return Left[R, A](e)
}

// Of creates a successful ReaderIOEither with the given value.
// This is the pointed functor operation, lifting a pure value into the ReaderIOEither context.
func Of[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return Right[R, E](a)
}

// Flatten removes one level of nesting from a nested ReaderIOEither.
// Converts ReaderIOEither[R, E, ReaderIOEither[R, E, A]] to ReaderIOEither[R, E, A].
func Flatten[R, E, A any](mma ReaderIOEither[R, E, ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return MonadChain(mma, function.Identity[ReaderIOEither[R, E, A]])
}

// FromEither lifts an Either into a ReaderIOEither context.
// The Either value is independent of any context or IO effects.
//
//go:inline
func FromEither[R, E, A any](t either.Either[E, A]) ReaderIOEither[R, E, A] {
	return readerio.Of[R](t)
}

// RightReader lifts a Reader into a ReaderIOEither, placing the result in the Right side.
func RightReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, IOE.Right[E, A])
}

// LeftReader lifts a Reader into a ReaderIOEither, placing the result in the Left (error) side.
func LeftReader[A, R, E any](ma Reader[R, E]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, IOE.Left[A, E])
}

// FromReader lifts a Reader into a ReaderIOEither context.
// The Reader result is placed in the Right side (success).
func FromReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return RightReader[E](ma)
}

// RightIO lifts an IO into a ReaderIOEither, placing the result in the Right side.
func RightIO[R, E, A any](ma IO[A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, IOE.RightIO[E, A], FromIOEither[R, E, A])
}

// LeftIO lifts an IO into a ReaderIOEither, placing the result in the Left (error) side.
func LeftIO[R, A, E any](ma IO[E]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, IOE.LeftIO[A, E], FromIOEither[R, E, A])
}

// FromIO lifts an IO into a ReaderIOEither context.
// The IO result is placed in the Right side (success).
func FromIO[R, E, A any](ma IO[A]) ReaderIOEither[R, E, A] {
	return RightIO[R, E](ma)
}

// FromIOEither lifts an IOEither into a ReaderIOEither context.
// The computation becomes independent of any reader context.
//
//go:inline
func FromIOEither[R, E, A any](ma IOEither[E, A]) ReaderIOEither[R, E, A] {
	return reader.Of[R](ma)
}

// FromReaderEither lifts a ReaderEither into a ReaderIOEither context.
// The Either result is lifted into an IO effect.
func FromReaderEither[R, E, A any](ma RE.ReaderEither[R, E, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, IOE.FromEither[E, A])
}

// Ask returns a ReaderIOEither that retrieves the current context.
// Useful for accessing configuration or dependencies.
//
//go:inline
func Ask[R, E any]() ReaderIOEither[R, E, R] {
	return fromreader.Ask(FromReader[E, R, R])()
}

// Asks returns a ReaderIOEither that retrieves a value derived from the context.
// This is useful for extracting specific fields from a configuration object.
//
//go:inline
func Asks[E, R, A any](r Reader[R, A]) ReaderIOEither[R, E, A] {
	return fromreader.Asks(FromReader[E, R, A])(r)
}

// FromOption converts an Option to a ReaderIOEither.
// If the Option is None, the provided function is called to produce the error.
//
//go:inline
func FromOption[R, A, E any](onNone Lazy[E]) func(Option[A]) ReaderIOEither[R, E, A] {
	return fromeither.FromOption(FromEither[R, E, A], onNone)
}

// FromPredicate creates a ReaderIOEither from a predicate.
// If the predicate returns false, the onFalse function is called to produce the error.
//
//go:inline
func FromPredicate[R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOEither[R, E, A] {
	return fromeither.FromPredicate(FromEither[R, E, A], pred, onFalse)
}

// Fold handles both success and error cases, producing a ReaderIO.
// This is useful for converting a ReaderIOEither into a ReaderIO by handling all cases.
//
//go:inline
func Fold[R, E, A, B any](onLeft readerio.Kleisli[R, E, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOEither[R, E, A]) ReaderIO[R, B] {
	return eithert.MatchE(readerio.MonadChain[R, either.Either[E, A], B], onLeft, onRight)
}

//go:inline
func MonadFold[R, E, A, B any](ma ReaderIOEither[R, E, A], onLeft readerio.Kleisli[R, E, B], onRight func(A) ReaderIO[R, B]) ReaderIO[R, B] {
	return eithert.FoldE(readerio.MonadChain[R, either.Either[E, A], B], ma, onLeft, onRight)
}

// GetOrElse provides a default value in case of error.
// The default is computed lazily via a ReaderIO.
//
//go:inline
func GetOrElse[R, E, A any](onLeft readerio.Kleisli[R, E, A]) func(ReaderIOEither[R, E, A]) ReaderIO[R, A] {
	return eithert.GetOrElse(readerio.MonadChain[R, either.Either[E, A], A], readerio.Of[R, A], onLeft)
}

// OrLeft transforms the error using a ReaderIO if the computation fails.
// The success value is preserved unchanged.
//
//go:inline
func OrLeft[A, E1, R, E2 any](onLeft func(E1) ReaderIO[R, E2]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.OrLeft(
		readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]],
		readerio.MonadMap[R, E2, either.Either[E2, A]],
		readerio.Of[R, either.Either[E2, A]],
		onLeft,
	)
}

// MonadBiMap applies two functions: one to the error, one to the success value.
// This allows transforming both channels simultaneously.
//
//go:inline
func MonadBiMap[R, E1, E2, A, B any](fa ReaderIOEither[R, E1, A], f func(E1) E2, g func(A) B) ReaderIOEither[R, E2, B] {
	return eithert.MonadBiMap(
		readerio.MonadMap[R, either.Either[E1, A], either.Either[E2, B]],
		fa, f, g,
	)
}

// BiMap returns a function that maps over both the error and success channels.
// This is the curried version of MonadBiMap.
//
//go:inline
func BiMap[R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, B] {
	return eithert.BiMap(readerio.Map[R, either.Either[E1, A], either.Either[E2, B]], f, g)
}

// Swap exchanges the error and success types.
// Left becomes Right and Right becomes Left.
//
//go:inline
func Swap[R, E, A any](val ReaderIOEither[R, E, A]) ReaderIOEither[R, A, E] {
	return reader.MonadMap(val, IOE.Swap[E, A])
}

// Defer creates a ReaderIOEither lazily via a generator function.
// The generator is called each time the ReaderIOEither is executed.
//
//go:inline
func Defer[R, E, A any](gen L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return readerio.Defer(gen)
}

// TryCatch wraps a function that returns (value, error) into a ReaderIOEither.
// The onThrow function converts the error into the desired error type.
func TryCatch[R, E, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOEither[R, E, A] {
	return func(r R) IOEither[E, A] {
		return IOE.TryCatch(f(r), onThrow)
	}
}

// MonadAlt tries the first computation, and if it fails, tries the second.
// This implements the Alternative pattern for error recovery.
//
//go:inline
func MonadAlt[R, E, A any](first ReaderIOEither[R, E, A], second L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return eithert.MonadAlt(
		readerio.Of[R, Either[E, A]],
		readerio.MonadChain[R, Either[E, A], Either[E, A]],

		first,
		second,
	)
}

// Alt returns a function that tries an alternative computation if the first fails.
// This is the curried version of MonadAlt.
//
//go:inline
func Alt[R, E, A any](second L.Lazy[ReaderIOEither[R, E, A]]) Operator[R, E, A, A] {
	return eithert.Alt(
		readerio.Of[R, Either[E, A]],
		readerio.Chain[R, Either[E, A], Either[E, A]],

		second,
	)
}

// Memoize computes the value of the ReaderIOEither lazily but exactly once.
// The context used is from the first call. Do not use if the value depends on the context.
//
//go:inline
func Memoize[
	R, E, A any](rdr ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return readerio.Memoize(rdr)
}

// MonadFlap applies a value to a function wrapped in a context.
// This is the reverse of Ap - the value is fixed and the function varies.
//
//go:inline
func MonadFlap[R, E, B, A any](fab ReaderIOEither[R, E, func(A) B], a A) ReaderIOEither[R, E, B] {
	return functor.MonadFlap(MonadMap[R, E, func(A) B, B], fab, a)
}

// Flap returns a function that applies a fixed value to a function in a context.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[R, E, B, A any](a A) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return functor.Flap(Map[R, E, func(A) B, B], a)
}

// MonadMapLeft applies a function to the error value, leaving success unchanged.
//
//go:inline
func MonadMapLeft[R, E1, E2, A any](fa ReaderIOEither[R, E1, A], f func(E1) E2) ReaderIOEither[R, E2, A] {
	return eithert.MonadMapLeft(readerio.MonadMap[R, Either[E1, A], Either[E2, A]], fa, f)
}

// MapLeft returns a function that transforms the error channel.
// This is the curried version of MonadMapLeft.
//
//go:inline
func MapLeft[R, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.MapLeft(readerio.Map[R, Either[E1, A], Either[E2, A]], f)
}

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
// This is similar to Contravariant's contramap operation.
//
//go:inline
func Local[E, A, R1, R2 any](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A] {
	return reader.Local[IOEither[E, A]](f)
}

//go:inline
func Read[E, A, R any](r R) func(ReaderIOEither[R, E, A]) IOEither[E, A] {
	return reader.Read[IOEither[E, A]](r)
}

// ReadIOEither executes a ReaderIOEither computation by providing it with an environment
// obtained from an IOEither computation. This is useful when the environment itself needs
// to be computed with side effects and error handling.
//
// The function first executes the IOEither[E, R] to get the environment R (or fail with error E),
// then uses that environment to run the ReaderIOEither[R, E, A] computation.
//
// Type parameters:
//   - A: The success value type of the ReaderIOEither computation
//   - R: The environment/context type required by the ReaderIOEither
//   - E: The error type
//
// Parameters:
//   - r: An IOEither[E, R] that produces the environment (or an error)
//
// Returns:
//   - A function that takes a ReaderIOEither[R, E, A] and returns IOEither[E, A]
//
// Behavior:
//   - If the IOEither[E, R] fails (Left), the error is propagated without running the ReaderIOEither
//   - If the IOEither[E, R] succeeds (Right), the resulting environment is used to execute the ReaderIOEither
//
// Example:
//
//	// Load configuration from a file (may fail)
//	loadConfig := func() IOEither[error, Config] {
//	    return Lazy[E]ither[error, Config] {
//	        // Read config file with error handling
//	        return either.Right(Config{BaseURL: "https://api.example.com"})
//	    }
//	}
//
//	// A computation that needs the config
//	fetchUser := func(id int) ReaderIOEither[Config, error, User] {
//	    return func(cfg Config) IOEither[error, User] {
//	        // Use cfg.BaseURL to fetch user
//	        return ioeither.Right[error](User{ID: id})
//	    }
//	}
//
//	// Execute the computation with dynamically loaded config
//	result := ReadIOEither[User](loadConfig())(fetchUser(123))()
//
//go:inline
func ReadIOEither[A, R, E any](r IOEither[E, R]) func(ReaderIOEither[R, E, A]) IOEither[E, A] {
	return function.Flow2(
		IOE.Chain[E, R, A],
		Read[E, A](r),
	)
}

// ReadIO executes a ReaderIOEither computation by providing it with an environment
// obtained from an IO computation. This is useful when the environment needs to be
// computed with side effects but cannot fail.
//
// The function first executes the IO[R] to get the environment R,
// then uses that environment to run the ReaderIOEither[R, E, A] computation.
//
// Type parameters:
//   - E: The error type of the ReaderIOEither computation
//   - A: The success value type of the ReaderIOEither computation
//   - R: The environment/context type required by the ReaderIOEither
//
// Parameters:
//   - r: An IO[R] that produces the environment
//
// Returns:
//   - A function that takes a ReaderIOEither[R, E, A] and returns IOEither[E, A]
//
// Behavior:
//   - The IO[R] is always executed successfully to obtain the environment
//   - The resulting environment is then used to execute the ReaderIOEither
//   - Only the ReaderIOEither computation can fail with error type E
//
// Example:
//
//	// Get current timestamp (cannot fail)
//	getCurrentTime := func() IO[time.Time] {
//	    return func() time.Time {
//	        return time.Now()
//	    }
//	}
//
//	// A computation that needs the timestamp
//	logWithTimestamp := func(msg string) ReaderIOEither[time.Time, error, string] {
//	    return func(t time.Time) IOEither[error, string] {
//	        logged := fmt.Sprintf("[%s] %s", t.Format(time.RFC3339), msg)
//	        return ioeither.Right[error](logged)
//	    }
//	}
//
//	// Execute the computation with current time
//	result := ReadIO[error, string](getCurrentTime())(logWithTimestamp("Hello"))()
//
//go:inline
func ReadIO[E, A, R any](r IO[R]) func(ReaderIOEither[R, E, A]) IOEither[E, A] {
	return function.Flow2(
		io.Chain[R, Either[E, A]],
		Read[E, A](r),
	)
}

// MonadChainLeft chains a computation on the left (error) side of a ReaderIOEither.
// If the input is a Left value, it applies the function f to transform the error and potentially
// change the error type from EA to EB. If the input is a Right value, it passes through unchanged.
//
// This is useful for error recovery or error transformation scenarios where you want to handle
// errors by performing another computation that may also fail, with access to configuration context.
//
// Note: This is functionally identical to the uncurried form of [OrElse]. Use [ChainLeft] when
// emphasizing the monadic chaining perspective, and [OrElse] for error recovery semantics.
//
// Parameters:
//   - fa: The input ReaderIOEither that may contain an error of type EA
//   - f: A Kleisli function that takes an error of type EA and returns a ReaderIOEither with error type EB
//
// Returns:
//   - A ReaderIOEither with the potentially transformed error type EB
//
// Example:
//
//	type Config struct{ retryCount int }
//	type NetworkError struct{ msg string }
//	type SystemError struct{ code int }
//
//	// Recover from network errors by retrying with config
//	result := MonadChainLeft(
//	    Left[Config, string](NetworkError{"connection failed"}),
//	    func(ne NetworkError) readerioeither.ReaderIOEither[Config, SystemError, string] {
//	        return readerioeither.Asks[SystemError](func(cfg Config) ioeither.IOEither[SystemError, string] {
//	            if cfg.retryCount > 0 {
//	                return ioeither.Right[SystemError]("recovered")
//	            }
//	            return ioeither.Left[string](SystemError{500})
//	        })
//	    },
//	)
//
//go:inline
func MonadChainLeft[R, EA, EB, A any](fa ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, A]) ReaderIOEither[R, EB, A] {
	return readert.MonadChain(
		IOE.MonadChainLeft[EA, EB, A],
		fa,
		f,
	)
}

// ChainLeft is the curried version of [MonadChainLeft].
// It returns a function that chains a computation on the left (error) side of a ReaderIOEither.
//
// This is particularly useful in functional composition pipelines where you want to handle
// errors by performing another computation that may also fail, with access to configuration context.
//
// Note: This is functionally identical to [OrElse]. They are different names for the same operation.
// Use [ChainLeft] when emphasizing the monadic chaining perspective on the error channel,
// and [OrElse] when emphasizing error recovery/fallback semantics.
//
// Parameters:
//   - f: A Kleisli function that takes an error of type EA and returns a ReaderIOEither with error type EB
//
// Returns:
//   - A function that transforms a ReaderIOEither with error type EA to one with error type EB
//
// Example:
//
//	type Config struct{ fallbackService string }
//
//	// Create a reusable error handler with config access
//	recoverFromNetworkError := ChainLeft(func(err string) readerioeither.ReaderIOEither[Config, string, int] {
//	    if strings.Contains(err, "network") {
//	        return readerioeither.Asks[string](func(cfg Config) ioeither.IOEither[string, int] {
//	            return ioeither.TryCatch(
//	                func() (int, error) { return callService(cfg.fallbackService) },
//	                func(e error) string { return e.Error() },
//	            )
//	        })
//	    }
//	    return readerioeither.Left[Config, int](err)
//	})
//
//	result := F.Pipe1(
//	    readerioeither.Left[Config, int]("network timeout"),
//	    recoverFromNetworkError,
//	)(Config{fallbackService: "backup"})()
//
//go:inline
func ChainLeft[R, EA, EB, A any](f Kleisli[R, EB, EA, A]) func(ReaderIOEither[R, EA, A]) ReaderIOEither[R, EB, A] {
	return readert.Chain[ReaderIOEither[R, EA, A]](
		IOE.ChainLeft[EA, EB, A],
		f,
	)
}

// MonadChainFirstLeft chains a computation on the left (error) side but always returns the original error.
// If the input is a Left value, it applies the function f to the error and executes the resulting computation,
// but always returns the original Left error regardless of what f returns (Left or Right).
// If the input is a Right value, it passes through unchanged without calling f.
//
// This is useful for side effects on errors (like logging or metrics) where you want to perform an action
// when an error occurs but always propagate the original error, ensuring the error path is preserved.
//
// Parameters:
//   - ma: The input ReaderIOEither that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a ReaderIOEither (typically for side effects)
//
// Returns:
//   - A ReaderIOEither with the original error preserved if input was Left, or the original Right value
//
//go:inline
func MonadChainFirstLeft[A, R, EA, EB, B any](ma ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderIOEither[R, EA, A] {
	return eithert.MonadChainFirstLeft(
		readerio.MonadChain[R, Either[EA, A], Either[EA, A]],
		readerio.MonadMap[R, Either[EB, B], Either[EA, A]],
		readerio.Of[R, Either[EA, A]],
		ma,
		f,
	)
}

//go:inline
func MonadTapLeft[A, R, EA, EB, B any](ma ReaderIOEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderIOEither[R, EA, A] {
	return MonadChainFirstLeft(ma, f)
}

// ChainFirstLeft is the curried version of [MonadChainFirstLeft].
// It returns a function that chains a computation on the left (error) side while always preserving the original error.
//
// This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// ensuring the error path is preserved.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a ReaderIOEither (typically for side effects)
//
// Returns:
//   - An Operator that performs the side effect but always returns the original error if input was Left
//
//go:inline
func ChainFirstLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A] {
	return eithert.ChainFirstLeft(
		readerio.Chain[R, Either[EA, A], Either[EA, A]],
		readerio.Map[R, Either[EB, B], Either[EA, A]],
		readerio.Of[R, Either[EA, A]],
		f,
	)
}

func ChainFirstLeftIOK[A, R, EA, B any](f io.Kleisli[EA, B]) Operator[R, EA, A, A] {
	return ChainFirstLeft[A](function.Flow2(
		f,
		FromIO[R, EA],
	))
}

//go:inline
func TapLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A] {
	return ChainFirstLeft[A](f)
}

//go:inline
func TapLeftIOK[A, R, EA, B any](f io.Kleisli[EA, B]) Operator[R, EA, A, A] {
	return ChainFirstLeftIOK[A, R](f)
}

// Delay creates an operation that passes in the value after some delay
//
//go:inline
func Delay[R, E, A any](delay time.Duration) Operator[R, E, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIOEither[R, E, A]], io.Delay[Either[E, A]](delay))
}

// After creates an operation that passes after the given [time.Time]
//
//go:inline
func After[R, E, A any](timestamp time.Time) Operator[R, E, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIOEither[R, E, A]], io.After[Either[E, A]](timestamp))
}

// OrElse recovers from a Left (error) by providing an alternative IO computation with access to the reader context.
// If the ReaderIOEither is Right, it returns the value unchanged.
// If the ReaderIOEither is Left, it applies the provided function to the error value,
// which returns a new ReaderIOEither that replaces the original.
//
// Note: OrElse is identical to [ChainLeft] - both provide the same functionality for error recovery.
//
// This is useful for error recovery, fallback logic, or chaining alternative IO computations
// that need access to configuration or dependencies. The error type can be widened from E1 to E2.
//
// Example:
//
//	type Config struct{ retryLimit int }
//
//	// Recover with IO operation using config
//	recover := readerioeither.OrElse(func(err error) readerioeither.ReaderIOEither[Config, error, int] {
//	    if err.Error() == "retryable" {
//	        return readerioeither.Asks[error](func(cfg Config) ioeither.IOEither[error, int] {
//	            if cfg.retryLimit > 0 {
//	                return ioeither.Right[error](42)
//	            }
//	            return ioeither.Left[int](err)
//	        })
//	    }
//	    return readerioeither.Left[Config, int](err)
//	})
//
//go:inline
func OrElse[R, E1, E2, A any](onLeft Kleisli[R, E2, E1, A]) Kleisli[R, E2, ReaderIOEither[R, E1, A], A] {
	return Fold(onLeft, Of[R, E2, A])
}
