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
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
)

// MonadFromReaderIO creates a ReaderIOEither from a value and a function that produces a ReaderIO.
// The ReaderIO result is lifted into the Right side of the Either.
func MonadFromReaderIO[R, E, A any](a A, f func(A) ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(
		a,
		f,
		RightReaderIO[R, E, A],
	)
}

// FromReaderIO creates a function that lifts a ReaderIO-producing function into ReaderIOEither.
// The ReaderIO result is placed in the Right side of the Either.
func FromReaderIO[R, E, A any](f func(A) ReaderIO[R, A]) func(A) ReaderIOEither[R, E, A] {
	return function.Bind2nd(MonadFromReaderIO[R, E, A], f)
}

// RightReaderIO lifts a ReaderIO into a ReaderIOEither, placing the result in the Right side.
func RightReaderIO[R, E, A any](ma ReaderIO[R, A]) ReaderIOEither[R, E, A] {
	return eithert.RightF(
		readerio.MonadMap[R, A, either.Either[E, A]],
		ma,
	)
}

// LeftReaderIO lifts a ReaderIO into a ReaderIOEither, placing the result in the Left (error) side.
func LeftReaderIO[A, R, E any](me ReaderIO[R, E]) ReaderIOEither[R, E, A] {
	return eithert.LeftF(
		readerio.MonadMap[R, E, either.Either[E, A]],
		me,
	)
}

// MonadMap applies a function to the value inside a ReaderIOEither context.
// If the computation is successful (Right), the function is applied to the value.
// If it's an error (Left), the error is propagated unchanged.
func MonadMap[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) B) ReaderIOEither[R, E, B] {
	return eithert.MonadMap(readerio.MonadMap[R, either.Either[E, A], either.Either[E, B]], fa, f)
}

// Map returns a function that applies a transformation to the success value of a ReaderIOEither.
// This is the curried version of MonadMap, useful for function composition.
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
func MonadChain[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, B] {
	return eithert.MonadChain(
		readerio.MonadChain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		fa,
		f)
}

// MonadChainFirst sequences two computations but keeps the result of the first.
// Useful for performing side effects while preserving the original value.
func MonadChainFirst[R, E, A, B any](fa ReaderIOEither[R, E, A], f func(A) ReaderIOEither[R, E, B]) ReaderIOEither[R, E, A] {
	return chain.MonadChainFirst(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		fa,
		f)
}

// MonadChainEitherK chains a computation that returns an Either into a ReaderIOEither.
// The Either is automatically lifted into the ReaderIOEither context.
func MonadChainEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		ma,
		f,
	)
}

// ChainEitherK returns a function that chains an Either-returning function into ReaderIOEither.
// This is the curried version of MonadChainEitherK.
func ChainEitherK[R, E, A, B any](f func(A) either.Either[E, B]) Operator[R, E, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, E, A, B],
		FromEither[R, E, B],
		f,
	)
}

// MonadChainFirstEitherK chains an Either-returning computation but keeps the original value.
// Useful for validation or side effects that return Either.
func MonadChainFirstEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) either.Either[E, B]) ReaderIOEither[R, E, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromEither[R, E, B],
		ma,
		f,
	)
}

// ChainFirstEitherK returns a function that chains an Either computation while preserving the original value.
// This is the curried version of MonadChainFirstEitherK.
func ChainFirstEitherK[R, E, A, B any](f func(A) either.Either[E, B]) Operator[R, E, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromEither[R, E, B],
		f,
	)
}

// MonadChainReaderK chains a Reader-returning computation into a ReaderIOEither.
// The Reader is automatically lifted into the ReaderIOEither context.
func MonadChainReaderK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) Reader[R, B]) ReaderIOEither[R, E, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, E, A, B],
		FromReader[E, R, B],
		ma,
		f,
	)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOEither.
// This is the curried version of MonadChainReaderK.
func ChainReaderK[E, R, A, B any](f func(A) Reader[R, B]) Operator[R, E, A, B] {
	return fromreader.ChainReaderK(
		MonadChain[R, E, A, B],
		FromReader[E, R, B],
		f,
	)
}

// MonadChainIOEitherK chains an IOEither-returning computation into a ReaderIOEither.
// The IOEither is automatically lifted into the ReaderIOEither context.
func MonadChainIOEitherK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) IOE.IOEither[E, B]) ReaderIOEither[R, E, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, E, A, B],
		FromIOEither[R, E, B],
		ma,
		f,
	)
}

// ChainIOEitherK returns a function that chains an IOEither-returning function into ReaderIOEither.
// This is the curried version of MonadChainIOEitherK.
func ChainIOEitherK[R, E, A, B any](f func(A) IOE.IOEither[E, B]) Operator[R, E, A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, E, A, B],
		FromIOEither[R, E, B],
		f,
	)
}

// MonadChainIOK chains an IO-returning computation into a ReaderIOEither.
// The IO is automatically lifted into the ReaderIOEither context (always succeeds).
func MonadChainIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, E, A, B],
		FromIO[R, E, B],
		ma,
		f,
	)
}

// ChainIOK returns a function that chains an IO-returning function into ReaderIOEither.
// This is the curried version of MonadChainIOK.
func ChainIOK[R, E, A, B any](f func(A) io.IO[B]) Operator[R, E, A, B] {
	return fromio.ChainIOK(
		Chain[R, E, A, B],
		FromIO[R, E, B],
		f,
	)
}

// MonadChainFirstIOK chains an IO computation but keeps the original value.
// Useful for performing IO side effects while preserving the original value.
func MonadChainFirstIOK[R, E, A, B any](ma ReaderIOEither[R, E, A], f func(A) io.IO[B]) ReaderIOEither[R, E, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, E, A, A],
		MonadMap[R, E, B, A],
		FromIO[R, E, B],
		ma,
		f,
	)
}

// ChainFirstIOK returns a function that chains an IO computation while preserving the original value.
// This is the curried version of MonadChainFirstIOK.
func ChainFirstIOK[R, E, A, B any](f func(A) io.IO[B]) Operator[R, E, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		FromIO[R, E, B],
		f,
	)
}

// ChainOptionK returns a function that chains an Option-returning function into ReaderIOEither.
// If the Option is None, the provided error function is called to produce the error value.
func ChainOptionK[R, A, B, E any](onNone func() E) func(func(A) O.Option[B]) Operator[R, E, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, E, A, B],
		FromEither[R, E, B],
		onNone,
	)
}

// MonadAp applies a function wrapped in a context to a value wrapped in a context.
// Both computations are executed (default behavior may be sequential or parallel depending on implementation).
func MonadAp[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadAp[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

// MonadApSeq applies a function in a context to a value in a context, executing them sequentially.
func MonadApSeq[R, E, A, B any](fab ReaderIOEither[R, E, func(A) B], fa ReaderIOEither[R, E, A]) ReaderIOEither[R, E, B] {
	return eithert.MonadAp(
		readerio.MonadApSeq[Either[E, B], R, Either[E, A]],
		readerio.MonadMap[R, Either[E, func(A) B], func(Either[E, A]) Either[E, B]],
		fab,
		fa,
	)
}

// MonadApPar applies a function in a context to a value in a context, executing them in parallel.
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
func Chain[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) Operator[R, E, A, B] {
	return eithert.Chain(
		readerio.Chain[R, either.Either[E, A], either.Either[E, B]],
		readerio.Of[R, either.Either[E, B]],
		f)
}

// ChainFirst returns a function that sequences computations but keeps the first result.
// This is the curried version of MonadChainFirst.
func ChainFirst[R, E, A, B any](f func(A) ReaderIOEither[R, E, B]) Operator[R, E, A, A] {
	return chain.ChainFirst(
		Chain[R, E, A, A],
		Map[R, E, B, A],
		f)
}

// Right creates a successful ReaderIOEither with the given value.
func Right[R, E, A any](a A) ReaderIOEither[R, E, A] {
	return eithert.Right(readerio.Of[R, Either[E, A]], a)
}

// Left creates a failed ReaderIOEither with the given error.
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
func FromEither[R, E, A any](t either.Either[E, A]) ReaderIOEither[R, E, A] {
	return readerio.Of[R](t)
}

// RightReader lifts a Reader into a ReaderIOEither, placing the result in the Right side.
func RightReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.Right[E, A])
}

// LeftReader lifts a Reader into a ReaderIOEither, placing the result in the Left (error) side.
func LeftReader[A, R, E any](ma Reader[R, E]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.Left[A, E])
}

// FromReader lifts a Reader into a ReaderIOEither context.
// The Reader result is placed in the Right side (success).
func FromReader[E, R, A any](ma Reader[R, A]) ReaderIOEither[R, E, A] {
	return RightReader[E](ma)
}

// RightIO lifts an IO into a ReaderIOEither, placing the result in the Right side.
func RightIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, ioeither.RightIO[E, A], FromIOEither[R, E, A])
}

// LeftIO lifts an IO into a ReaderIOEither, placing the result in the Left (error) side.
func LeftIO[R, A, E any](ma io.IO[E]) ReaderIOEither[R, E, A] {
	return function.Pipe2(ma, ioeither.LeftIO[A, E], FromIOEither[R, E, A])
}

// FromIO lifts an IO into a ReaderIOEither context.
// The IO result is placed in the Right side (success).
func FromIO[R, E, A any](ma io.IO[A]) ReaderIOEither[R, E, A] {
	return RightIO[R, E](ma)
}

// FromIOEither lifts an IOEither into a ReaderIOEither context.
// The computation becomes independent of any reader context.
func FromIOEither[R, E, A any](ma IOE.IOEither[E, A]) ReaderIOEither[R, E, A] {
	return reader.Of[R](ma)
}

// FromReaderEither lifts a ReaderEither into a ReaderIOEither context.
// The Either result is lifted into an IO effect.
func FromReaderEither[R, E, A any](ma RE.ReaderEither[R, E, A]) ReaderIOEither[R, E, A] {
	return function.Flow2(ma, ioeither.FromEither[E, A])
}

// Ask returns a ReaderIOEither that retrieves the current context.
// Useful for accessing configuration or dependencies.
func Ask[R, E any]() ReaderIOEither[R, E, R] {
	return fromreader.Ask(FromReader[E, R, R])()
}

// Asks returns a ReaderIOEither that retrieves a value derived from the context.
// This is useful for extracting specific fields from a configuration object.
func Asks[E, R, A any](r Reader[R, A]) ReaderIOEither[R, E, A] {
	return fromreader.Asks(FromReader[E, R, A])(r)
}

// FromOption converts an Option to a ReaderIOEither.
// If the Option is None, the provided function is called to produce the error.
func FromOption[R, A, E any](onNone func() E) func(O.Option[A]) ReaderIOEither[R, E, A] {
	return fromeither.FromOption(FromEither[R, E, A], onNone)
}

// FromPredicate creates a ReaderIOEither from a predicate.
// If the predicate returns false, the onFalse function is called to produce the error.
func FromPredicate[R, E, A any](pred func(A) bool, onFalse func(A) E) func(A) ReaderIOEither[R, E, A] {
	return fromeither.FromPredicate(FromEither[R, E, A], pred, onFalse)
}

// Fold handles both success and error cases, producing a ReaderIO.
// This is useful for converting a ReaderIOEither into a ReaderIO by handling all cases.
func Fold[R, E, A, B any](onLeft func(E) ReaderIO[R, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOEither[R, E, A]) ReaderIO[R, B] {
	return eithert.MatchE(readerio.MonadChain[R, either.Either[E, A], B], onLeft, onRight)
}

// GetOrElse provides a default value in case of error.
// The default is computed lazily via a ReaderIO.
func GetOrElse[R, E, A any](onLeft func(E) ReaderIO[R, A]) func(ReaderIOEither[R, E, A]) ReaderIO[R, A] {
	return eithert.GetOrElse(readerio.MonadChain[R, either.Either[E, A], A], readerio.Of[R, A], onLeft)
}

// OrElse tries an alternative computation if the first one fails.
// The alternative can produce a different error type.
func OrElse[R, E1, A, E2 any](onLeft func(E1) ReaderIOEither[R, E2, A]) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.OrElse(readerio.MonadChain[R, either.Either[E1, A], either.Either[E2, A]], readerio.Of[R, either.Either[E2, A]], onLeft)
}

// OrLeft transforms the error using a ReaderIO if the computation fails.
// The success value is preserved unchanged.
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
func MonadBiMap[R, E1, E2, A, B any](fa ReaderIOEither[R, E1, A], f func(E1) E2, g func(A) B) ReaderIOEither[R, E2, B] {
	return eithert.MonadBiMap(
		readerio.MonadMap[R, either.Either[E1, A], either.Either[E2, B]],
		fa, f, g,
	)
}

// BiMap returns a function that maps over both the error and success channels.
// This is the curried version of MonadBiMap.
func BiMap[R, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, B] {
	return eithert.BiMap(readerio.Map[R, either.Either[E1, A], either.Either[E2, B]], f, g)
}

// Swap exchanges the error and success types.
// Left becomes Right and Right becomes Left.
func Swap[R, E, A any](val ReaderIOEither[R, E, A]) ReaderIOEither[R, A, E] {
	return reader.MonadMap(val, ioeither.Swap[E, A])
}

// Defer creates a ReaderIOEither lazily via a generator function.
// The generator is called each time the ReaderIOEither is executed.
func Defer[R, E, A any](gen L.Lazy[ReaderIOEither[R, E, A]]) ReaderIOEither[R, E, A] {
	return readerio.Defer(gen)
}

// TryCatch wraps a function that returns (value, error) into a ReaderIOEither.
// The onThrow function converts the error into the desired error type.
func TryCatch[R, E, A any](f func(R) func() (A, error), onThrow func(error) E) ReaderIOEither[R, E, A] {
	return func(r R) IOEither[E, A] {
		return ioeither.TryCatch(f(r), onThrow)
	}
}

// MonadAlt tries the first computation, and if it fails, tries the second.
// This implements the Alternative pattern for error recovery.
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
func Alt[R, E, A any](second L.Lazy[ReaderIOEither[R, E, A]]) Operator[R, E, A, A] {
	return eithert.Alt(
		readerio.Of[R, Either[E, A]],
		readerio.MonadChain[R, Either[E, A], Either[E, A]],

		second,
	)
}

// Memoize computes the value of the ReaderIOEither lazily but exactly once.
// The context used is from the first call. Do not use if the value depends on the context.
func Memoize[
	R, E, A any](rdr ReaderIOEither[R, E, A]) ReaderIOEither[R, E, A] {
	return readerio.Memoize(rdr)
}

// MonadFlap applies a value to a function wrapped in a context.
// This is the reverse of Ap - the value is fixed and the function varies.
func MonadFlap[R, E, B, A any](fab ReaderIOEither[R, E, func(A) B], a A) ReaderIOEither[R, E, B] {
	return functor.MonadFlap(MonadMap[R, E, func(A) B, B], fab, a)
}

// Flap returns a function that applies a fixed value to a function in a context.
// This is the curried version of MonadFlap.
func Flap[R, E, B, A any](a A) func(ReaderIOEither[R, E, func(A) B]) ReaderIOEither[R, E, B] {
	return functor.Flap(Map[R, E, func(A) B, B], a)
}

// MonadMapLeft applies a function to the error value, leaving success unchanged.
func MonadMapLeft[R, E1, E2, A any](fa ReaderIOEither[R, E1, A], f func(E1) E2) ReaderIOEither[R, E2, A] {
	return eithert.MonadMapLeft(readerio.MonadMap[R, Either[E1, A], Either[E2, A]], fa, f)
}

// MapLeft returns a function that transforms the error channel.
// This is the curried version of MonadMapLeft.
func MapLeft[R, A, E1, E2 any](f func(E1) E2) func(ReaderIOEither[R, E1, A]) ReaderIOEither[R, E2, A] {
	return eithert.MapLeft(readerio.Map[R, Either[E1, A], Either[E2, A]], f)
}

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
// This is similar to Contravariant's contramap operation.
func Local[R1, R2, E, A any](f func(R2) R1) func(ReaderIOEither[R1, E, A]) ReaderIOEither[R2, E, A] {
	return reader.Local[R2, R1, IOEither[E, A]](f)
}
