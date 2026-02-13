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

package readerioresult

import (
	"time"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
)

//go:inline
func FromReaderOption[R, A any](onNone Lazy[error]) Kleisli[R, ReaderOption[R, A], A] {
	return RIOE.FromReaderOption[R, A](onNone)
}

// FromReaderIO creates a function that lifts a ReaderIO-producing function into ReaderIOResult.
// The ReaderIO result is placed in the Right side of the Either.
//
//go:inline
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A] {
	return RIOE.FromReaderIO[error](ma)
}

// RightReaderIO lifts a ReaderIO into a ReaderIOResult, placing the result in the Right side.
//
//go:inline
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderIOResult[R, A] {
	return RIOE.RightReaderIO[error](ma)
}

// LeftReaderIO lifts a ReaderIO into a ReaderIOResult, placing the result in the Left (error) side.
//
//go:inline
func LeftReaderIO[A, R any](me ReaderIO[R, error]) ReaderIOResult[R, A] {
	return RIOE.LeftReaderIO[A](me)
}

// MonadMap applies a function to the value inside a ReaderIOResult context.
// If the computation is successful (Right), the function is applied to the value.
// If it's an error (Left), the error is propagated unchanged.
//
//go:inline
func MonadMap[R, A, B any](fa ReaderIOResult[R, A], f func(A) B) ReaderIOResult[R, B] {
	return RIOE.MonadMap(fa, f)
}

// Map returns a function that applies a transformation to the success value of a ReaderIOResult.
// This is the curried version of MonadMap, useful for function composition.
//
//go:inline
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return RIOE.Map[R, error](f)
}

// MonadMapTo replaces the success value with a constant value.
// Useful when you want to discard the result but keep the effect.
//
//go:inline
func MonadMapTo[R, A, B any](fa ReaderIOResult[R, A], b B) ReaderIOResult[R, B] {
	return RIOE.MonadMapTo(fa, b)
}

// MapTo returns a function that replaces the success value with a constant.
// This is the curried version of MonadMapTo.
//
//go:inline
func MapTo[R, A, B any](b B) Operator[R, A, B] {
	return RIOE.MapTo[R, error, A](b)
}

// MonadChain sequences two computations where the second depends on the result of the first.
// This is the fundamental operation for composing dependent effectful computations.
// If the first computation fails, the second is not executed.
//
//go:inline
func MonadChain[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChain(fa, f)
}

// MonadChainFirst sequences two computations but keeps the result of the first.
// Useful for performing side effects while preserving the original value.
//
//go:inline
func MonadChainFirst[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirst(fa, f)
}

//go:inline
func MonadTap[R, A, B any](fa ReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTap(fa, f)
}

// MonadChainEitherK chains a computation that returns an Either into a ReaderIOResult.
// The Either is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainEitherK(ma, f)
}

// MonadChainEitherK chains a computation that returns an Either into a ReaderIOResult.
// The Either is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainEitherK(ma, f)
}

// ChainEitherK returns a function that chains an Either-returning function into ReaderIOResult.
// This is the curried version of MonadChainEitherK.
//
//go:inline
func ChainEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B] {
	return RIOE.ChainEitherK[R](f)
}

// ChainResultK returns a function that chains an Either-returning function into ReaderIOResult.
// This is the curried version of MonadChainEitherK.
//
//go:inline
func ChainResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B] {
	return RIOE.ChainEitherK[R](f)
}

// MonadChainFirstEitherK chains an Either-returning computation but keeps the original value.
// Useful for validation or side effects that return Either.
//
//go:inline
func MonadChainFirstEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstEitherK(ma, f)
}

//go:inline
func MonadTapEitherK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapEitherK(ma, f)
}

// ChainFirstEitherK returns a function that chains an Either computation while preserving the original value.
// This is the curried version of MonadChainFirstEitherK.
//
//go:inline
func ChainFirstEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstEitherK[R](f)
}

//go:inline
func TapEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A] {
	return RIOE.TapEitherK[R](f)
}

// MonadChainFirstEitherK chains an Either-returning computation but keeps the original value.
// Useful for validation or side effects that return Either.
//
//go:inline
func MonadChainFirstResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstEitherK(ma, f)
}

//go:inline
func MonadTapResultK[R, A, B any](ma ReaderIOResult[R, A], f result.Kleisli[A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapEitherK(ma, f)
}

// ChainFirstEitherK returns a function that chains an Either computation while preserving the original value.
// This is the curried version of MonadChainFirstEitherK.
//
//go:inline
func ChainFirstResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstEitherK[R](f)
}

//go:inline
func TapResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, A] {
	return RIOE.TapEitherK[R](f)
}

// MonadChainReaderK chains a Reader-returning computation into a ReaderIOResult.
// The Reader is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainReaderK(ma, f)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return RIOE.ChainReaderK[error](f)
}

//go:inline
func MonadChainFirstReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstReaderK(ma, f)
}

//go:inline
func MonadTapReaderK[R, A, B any](ma ReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapReaderK(ma, f)
}

//go:inline
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstReaderK[error](f)
}

//go:inline
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.TapReaderK[error](f)
}

//go:inline
func ChainReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, B] {
	return RIOE.ChainReaderOptionK[R, A, B](onNone)
}

//go:inline
func ChainFirstReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstReaderOptionK[R, A, B](onNone)
}

//go:inline
func TapReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.TapReaderOptionK[R, A, B](onNone)
}

// MonadChainReaderK chains a Reader-returning computation into a ReaderIOResult.
// The Reader is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainReaderEitherK(ma, f)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B] {
	return RIOE.ChainReaderEitherK(f)
}

//go:inline
func MonadChainFirstReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstReaderEitherK(ma, f)
}

//go:inline
func MonadTapReaderEitherK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapReaderEitherK(ma, f)
}

//go:inline
func ChainFirstReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstReaderEitherK(f)
}

//go:inline
func TapReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return RIOE.TapReaderEitherK(f)
}

//go:inline
func MonadChainReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainReaderEitherK(ma, f)
}

// ChainReaderK returns a function that chains a Reader-returning function into ReaderIOResult.
// This is the curried version of MonadChainReaderK.
//
//go:inline
func ChainReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B] {
	return RIOE.ChainReaderEitherK(f)
}

//go:inline
func MonadChainFirstReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstReaderEitherK(ma, f)
}

//go:inline
func MonadTapReaderResultK[R, A, B any](ma ReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapReaderEitherK(ma, f)
}

//go:inline
func ChainFirstReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstReaderEitherK(f)
}

//go:inline
func TapReaderResultK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return RIOE.TapReaderEitherK(f)
}

//go:inline
func MonadChainReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainReaderIOK(ma, f)
}

//go:inline
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B] {
	return RIOE.ChainReaderIOK[error](f)
}

//go:inline
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstReaderIOK(ma, f)
}

//go:inline
func MonadTapReaderIOK[R, A, B any](ma ReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapReaderIOK(ma, f)
}

//go:inline
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirstReaderIOK[error](f)
}

//go:inline
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.TapReaderIOK[error](f)
}

// MonadChainIOEitherK chains an IOEither-returning computation into a ReaderIOResult.
// The IOEither is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainIOEitherK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IOResult[B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainIOEitherK(ma, f)
}

// MonadChainIOEitherK chains an IOEither-returning computation into a ReaderIOResult.
// The IOEither is automatically lifted into the ReaderIOResult context.
//
//go:inline
func MonadChainIOResultK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IOResult[B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainIOEitherK(ma, f)
}

// ChainIOEitherK returns a function that chains an IOEither-returning function into ReaderIOResult.
// This is the curried version of MonadChainIOEitherK.
//
//go:inline
func ChainIOEitherK[R, A, B any](f func(A) IOResult[B]) Operator[R, A, B] {
	return RIOE.ChainIOEitherK[R](f)
}

// ChainIOEitherK returns a function that chains an IOEither-returning function into ReaderIOResult.
// This is the curried version of MonadChainIOEitherK.
//
//go:inline
func ChainIOResultK[R, A, B any](f func(A) IOResult[B]) Operator[R, A, B] {
	return RIOE.ChainIOEitherK[R](f)
}

// MonadChainIOK chains an IO-returning computation into a ReaderIOResult.
// The IO is automatically lifted into the ReaderIOResult context (always succeeds).
//
//go:inline
func MonadChainIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, B] {
	return RIOE.MonadChainIOK(ma, f)
}

// ChainIOK returns a function that chains an IO-returning function into ReaderIOResult.
// This is the curried version of MonadChainIOK.
//
//go:inline
func ChainIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, B] {
	return RIOE.ChainIOK[R, error](f)
}

// MonadChainFirstIOK chains an IO computation but keeps the original value.
// Useful for performing IO side effects while preserving the original value.
//
//go:inline
func MonadChainFirstIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstIOK(ma, f)
}

//go:inline
func MonadTapIOK[R, A, B any](ma ReaderIOResult[R, A], f func(A) IO[B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapIOK(ma, f)
}

// ChainFirstIOK returns a function that chains an IO computation while preserving the original value.
// This is the curried version of MonadChainFirstIOK.
//
//go:inline
func ChainFirstIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, A] {
	return RIOE.ChainFirstIOK[R, error](f)
}

//go:inline
func TapIOK[R, A, B any](f func(A) IO[B]) Operator[R, A, A] {
	return RIOE.TapIOK[R, error](f)
}

// ChainOptionK returns a function that chains an Option-returning function into ReaderIOResult.
// If the Option is None, the provided error function is called to produce the error value.
//
//go:inline
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(func(A) Option[B]) Operator[R, A, B] {
	return RIOE.ChainOptionK[R, A, B](onNone)
}

// MonadAp applies a function wrapped in a context to a value wrapped in a context.
// Both computations are executed (default behavior may be sequential or parallel depending on implementation).
//
//go:inline
func MonadAp[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return RIOE.MonadAp(fab, fa)
}

// MonadApSeq applies a function in a context to a value in a context, executing them sequentially.
//
//go:inline
func MonadApSeq[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return RIOE.MonadApSeq(fab, fa)
}

// MonadApPar applies a function in a context to a value in a context, executing them in parallel.
//
//go:inline
func MonadApPar[R, A, B any](fab ReaderIOResult[R, func(A) B], fa ReaderIOResult[R, A]) ReaderIOResult[R, B] {
	return RIOE.MonadApPar(fab, fa)
}

// Ap returns a function that applies a function in a context to a value in a context.
// This is the curried version of MonadAp.
//
//go:inline
func Ap[B, R, A any](fa ReaderIOResult[R, A]) Operator[R, func(A) B, B] {
	return RIOE.Ap[B](fa)

}

// Chain returns a function that sequences computations where the second depends on the first.
// This is the curried version of MonadChain.
//
//go:inline
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return RIOE.Chain(f)
}

// ChainFirst returns a function that sequences computations but keeps the first result.
// This is the curried version of MonadChainFirst.
//
//go:inline
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.ChainFirst(f)
}

//go:inline
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return RIOE.Tap(f)
}

// Right creates a successful ReaderIOResult with the given value.
//
//go:inline
func Right[R, A any](a A) ReaderIOResult[R, A] {
	return RIOE.Right[R, error](a)
}

// Left creates a failed ReaderIOResult with the given error.
//
//go:inline
func Left[R, A any](e error) ReaderIOResult[R, A] {
	return RIOE.Left[R, A](e)
}

// ThrowError creates a failed ReaderIOResult with the given error.
// This is an alias for Left, following the naming convention from other functional libraries.
//
//go:inline
func ThrowError[R, A any](e error) ReaderIOResult[R, A] {
	return RIOE.ThrowError[R, A](e)
}

// Of creates a successful ReaderIOResult with the given value.
// This is the pointed functor operation, lifting a pure value into the ReaderIOResult context.
//
//go:inline
func Of[R, A any](a A) ReaderIOResult[R, A] {
	return RIOE.Of[R, error](a)
}

// Flatten removes one level of nesting from a nested ReaderIOResult.
// Converts ReaderIOResult[R, ReaderIOResult[R, A]] to ReaderIOResult[R, A].
//
//go:inline
func Flatten[R, A any](mma ReaderIOResult[R, ReaderIOResult[R, A]]) ReaderIOResult[R, A] {
	return RIOE.Flatten(mma)
}

// FromEither lifts an Either into a ReaderIOResult context.
// The Either value is independent of any context or IO effects.
//
//go:inline
func FromEither[R, A any](t Result[A]) ReaderIOResult[R, A] {
	return RIOE.FromEither[R](t)
}

// FromResult lifts an Either into a ReaderIOResult context.
// The Either value is independent of any context or IO effects.
//
//go:inline
func FromResult[R, A any](t Result[A]) ReaderIOResult[R, A] {
	return RIOE.FromEither[R](t)
}

// RightReader lifts a Reader into a ReaderIOResult, placing the result in the Right side.
//
//go:inline
func RightReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A] {
	return RIOE.RightReader[error](ma)
}

// LeftReader lifts a Reader into a ReaderIOResult, placing the result in the Left (error) side.
//
//go:inline
func LeftReader[A, R any](ma Reader[R, error]) ReaderIOResult[R, A] {
	return RIOE.LeftReader[A](ma)
}

// FromReader lifts a Reader into a ReaderIOResult context.
// The Reader result is placed in the Right side (success).
//
//go:inline
func FromReader[R, A any](ma Reader[R, A]) ReaderIOResult[R, A] {
	return RIOE.FromReader[error](ma)
}

// RightIO lifts an IO into a ReaderIOResult, placing the result in the Right side.
//
//go:inline
func RightIO[R, A any](ma IO[A]) ReaderIOResult[R, A] {
	return RIOE.RightIO[R, error](ma)
}

// LeftIO lifts an IO into a ReaderIOResult, placing the result in the Left (error) side.
//
//go:inline
func LeftIO[R, A any](ma IO[error]) ReaderIOResult[R, A] {
	return RIOE.LeftIO[R, A](ma)
}

// FromIO lifts an IO into a ReaderIOResult context.
// The IO result is placed in the Right side (success).
//
//go:inline
func FromIO[R, A any](ma IO[A]) ReaderIOResult[R, A] {
	return RIOE.FromIO[R, error](ma)
}

// FromIOEither lifts an IOEither into a ReaderIOResult context.
// The computation becomes independent of any reader context.
//
//go:inline
func FromIOEither[R, A any](ma IOResult[A]) ReaderIOResult[R, A] {
	return RIOE.FromIOEither[R](ma)
}

// FromIOEither lifts an IOEither into a ReaderIOResult context.
// The computation becomes independent of any reader context.
//
//go:inline
func FromIOResult[R, A any](ma IOResult[A]) ReaderIOResult[R, A] {
	return RIOE.FromIOEither[R](ma)
}

// FromReaderEither lifts a ReaderEither into a ReaderIOResult context.
// The Either result is lifted into an IO effect.
//
//go:inline
func FromReaderEither[R, A any](ma RE.ReaderEither[R, error, A]) ReaderIOResult[R, A] {
	return RIOE.FromReaderEither(ma)
}

// Ask returns a ReaderIOResult that retrieves the current context.
// Useful for accessing configuration or dependencies.
//
//go:inline
func Ask[R any]() ReaderIOResult[R, R] {
	return RIOE.Ask[R, error]()
}

// Asks returns a ReaderIOResult that retrieves a value derived from the context.
// This is useful for extracting specific fields from a configuration object.
//
//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderIOResult[R, A] {
	return RIOE.Asks[error](r)
}

// FromOption converts an Option to a ReaderIOResult.
// If the Option is None, the provided function is called to produce the error.
//
//go:inline
func FromOption[R, A any](onNone Lazy[error]) Kleisli[R, Option[A], A] {
	return RIOE.FromOption[R, A](onNone)
}

// FromPredicate creates a ReaderIOResult from a predicate.
// If the predicate returns false, the onFalse function is called to produce the error.
//
//go:inline
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A] {
	return RIOE.FromPredicate[R](pred, onFalse)
}

// Fold handles both success and error cases, producing a ReaderIO.
// This is useful for converting a ReaderIOResult into a ReaderIO by handling all cases.
//
//go:inline
func Fold[R, A, B any](onLeft readerio.Kleisli[R, error, B], onRight func(A) ReaderIO[R, B]) func(ReaderIOResult[R, A]) ReaderIO[R, B] {
	return RIOE.Fold(onLeft, onRight)
}

// GetOrElse provides a default value in case of error.
// The default is computed lazily via a ReaderIO.
//
//go:inline
func GetOrElse[R, A any](onLeft readerio.Kleisli[R, error, A]) func(ReaderIOResult[R, A]) ReaderIO[R, A] {
	return RIOE.GetOrElse(onLeft)
}

// OrElse tries an alternative computation if the first one fails.
//
//go:inline
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A] {
	return RIOE.OrElse(onLeft)
}

// OrLeft transforms the error using a ReaderIO if the computation fails.
// The success value is preserved unchanged.
//
//go:inline
func OrLeft[A, R, E any](onLeft readerio.Kleisli[R, error, E]) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, A] {
	return RIOE.OrLeft[A](onLeft)
}

// MonadBiMap applies two functions: one to the error, one to the success value.
// This allows transforming both channels simultaneously.
//
//go:inline
func MonadBiMap[R, E, A, B any](fa ReaderIOResult[R, A], f func(error) E, g func(A) B) RIOE.ReaderIOEither[R, E, B] {
	return RIOE.MonadBiMap(fa, f, g)
}

// BiMap returns a function that maps over both the error and success channels.
// This is the curried version of MonadBiMap.
//
//go:inline
func BiMap[R, E, A, B any](f func(error) E, g func(A) B) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, B] {
	return RIOE.BiMap[R](f, g)
}

// Swap exchanges the error and success types.
// Left becomes Right and Right becomes Left.
//
//go:inline
func Swap[R, A any](val ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, A, error] {
	return RIOE.Swap(val)
}

// Defer creates a ReaderIOResult lazily via a generator function.
// The generator is called each time the ReaderIOResult is executed.
//
//go:inline
func Defer[R, A any](gen Lazy[ReaderIOResult[R, A]]) ReaderIOResult[R, A] {
	return RIOE.Defer(gen)
}

// TryCatch wraps a function that returns (value, error) into a ReaderIOResult.
// The onThrow function converts the error into the desired error type.
//
//go:inline
func TryCatch[R, A any](f func(R) func() (A, error), onThrow Endomorphism[error]) ReaderIOResult[R, A] {
	return RIOE.TryCatch(f, onThrow)
}

// MonadAlt tries the first computation, and if it fails, tries the second.
// This implements the Alternative pattern for error recovery.
//
//go:inline
func MonadAlt[R, A any](first ReaderIOResult[R, A], second Lazy[ReaderIOResult[R, A]]) ReaderIOResult[R, A] {
	return RIOE.MonadAlt(first, second)
}

// Alt returns a function that tries an alternative computation if the first fails.
// This is the curried version of MonadAlt.
//
//go:inline
func Alt[R, A any](second Lazy[ReaderIOResult[R, A]]) Operator[R, A, A] {
	return RIOE.Alt(second)
}

// Memoize computes the value of the ReaderIOResult lazily but exactly once.
// The context used is from the first call. Do not use if the value depends on the context.
//
//go:inline
func Memoize[R, A any](rdr ReaderIOResult[R, A]) ReaderIOResult[R, A] {
	return RIOE.Memoize(rdr)
}

// MonadFlap applies a value to a function wrapped in a context.
// This is the reverse of Ap - the value is fixed and the function varies.
//
//go:inline
func MonadFlap[R, B, A any](fab ReaderIOResult[R, func(A) B], a A) ReaderIOResult[R, B] {
	return RIOE.MonadFlap(fab, a)
}

// Flap returns a function that applies a fixed value to a function in a context.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return RIOE.Flap[R, error, B](a)
}

// MonadMapLeft applies a function to the error value, leaving success unchanged.
//
//go:inline
func MonadMapLeft[R, E, A any](fa ReaderIOResult[R, A], f func(error) E) RIOE.ReaderIOEither[R, E, A] {
	return RIOE.MonadMapLeft(fa, f)
}

// MapLeft returns a function that transforms the error channel.
// This is the curried version of MonadMapLeft.
//
//go:inline
func MapLeft[R, A, E any](f func(error) E) func(ReaderIOResult[R, A]) RIOE.ReaderIOEither[R, E, A] {
	return RIOE.MapLeft[R, A](f)
}

// Local runs a computation with a modified context.
// The function f transforms the context before passing it to the computation.
// This is similar to Contravariant's contramap operation.
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOResult[R1, A]) ReaderIOResult[R2, A] {
	return RIOE.Local[error, A](f)
}

//go:inline
func Read[A, R any](r R) func(ReaderIOResult[R, A]) IOResult[A] {
	return RIOE.Read[error, A](r)
}

//go:inline
func MonadChainLeft[R, A any](fa ReaderIOResult[R, A], f Kleisli[R, error, A]) ReaderIOResult[R, A] {
	return RIOE.MonadChainLeft(fa, f)
}

//go:inline
func ChainLeft[R, A any](f Kleisli[R, error, A]) func(ReaderIOResult[R, A]) ReaderIOResult[R, A] {
	return RIOE.ChainLeft(f)
}

//go:inline
func MonadChainFirstLeft[A, R, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A] {
	return RIOE.MonadChainFirstLeft(ma, f)
}

//go:inline
func MonadTapLeft[A, R, B any](ma ReaderIOResult[R, A], f Kleisli[R, error, B]) ReaderIOResult[R, A] {
	return RIOE.MonadTapLeft(ma, f)
}

//go:inline
func ChainFirstLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return RIOE.ChainFirstLeft[A](f)
}

//go:inline
func ChainFirstLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return RIOE.ChainFirstLeftIOK[A, R](f)
}

//go:inline
func TapLeft[A, R, B any](f Kleisli[R, error, B]) Operator[R, A, A] {
	return RIOE.TapLeft[A](f)
}

//go:inline
func TapLeftIOK[A, R, B any](f io.Kleisli[error, B]) Operator[R, A, A] {
	return RIOE.TapLeftIOK[A, R](f)
}

// Delay creates an operation that passes in the value after some delay
//
//go:inline
func Delay[R, A any](delay time.Duration) Operator[R, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIOResult[R, A]], io.Delay[Result[A]](delay))
}

// After creates an operation that passes after the given [time.Time]
//
//go:inline
func After[R, A any](timestamp time.Time) Operator[R, A, A] {
	return function.Bind2nd(function.Flow2[ReaderIOResult[R, A]], io.After[Result[A]](timestamp))
}

// ReadIOEither executes a ReaderIOResult computation by providing an environment
// obtained from an IOResult. This function bridges the gap between IOResult-based
// environment acquisition and ReaderIOResult-based computations.
//
// The function first executes the IOResult[R] to obtain the environment (or an error),
// then uses that environment to run the ReaderIOResult[R, A] computation.
//
// Type parameters:
//   - A: The success value type of the ReaderIOResult computation
//   - R: The environment/context type required by the ReaderIOResult
//
// Parameters:
//   - r: An IOResult[R] that produces the environment (or an error)
//
// Returns:
//   - A function that takes a ReaderIOResult[R, A] and returns IOResult[A]
//
// Example:
//
//	type Config struct { BaseURL string }
//
//	// Get config from environment with potential error
//	getConfig := func() IOResult[Config] {
//	    return func() Result[Config] {
//	        // Load config, may fail
//	        return result.Of(Config{BaseURL: "https://api.example.com"})
//	    }
//	}
//
//	// A computation that needs config
//	fetchUser := func(id int) ReaderIOResult[Config, User] {
//	    return func(cfg Config) IOResult[User] {
//	        return func() Result[User] {
//	            // Use cfg.BaseURL to fetch user
//	            return result.Of(User{ID: id})
//	        }
//	    }
//	}
//
//	// Execute the computation with the config
//	result := ReadIOEither[User](getConfig())(fetchUser(123))()
//
//go:inline
func ReadIOEither[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A] {
	return RIOE.ReadIOEither[A](r)
}

// ReadIOResult executes a ReaderIOResult computation by providing an environment
// obtained from an IOResult. This is an alias for ReadIOEither with more explicit naming.
//
// The function first executes the IOResult[R] to obtain the environment (or an error),
// then uses that environment to run the ReaderIOResult[R, A] computation.
//
// Type parameters:
//   - A: The success value type of the ReaderIOResult computation
//   - R: The environment/context type required by the ReaderIOResult
//
// Parameters:
//   - r: An IOResult[R] that produces the environment (or an error)
//
// Returns:
//   - A function that takes a ReaderIOResult[R, A] and returns IOResult[A]
//
// Example:
//
//	type Database struct { ConnectionString string }
//
//	// Get database connection with potential error
//	getDB := func() IOResult[Database] {
//	    return func() Result[Database] {
//	        return result.Of(Database{ConnectionString: "localhost:5432"})
//	    }
//	}
//
//	// Query that needs database
//	queryUsers := ReaderIOResult[Database, []User] {
//	    return func(db Database) IOResult[[]User] {
//	        return func() Result[[]User] {
//	            // Execute query using db
//	            return result.Of([]User{})
//	        }
//	    }
//	}
//
//	// Execute query with database
//	users := ReadIOResult[[]User](getDB())(queryUsers)()
//
//go:inline
func ReadIOResult[A, R any](r IOResult[R]) func(ReaderIOResult[R, A]) IOResult[A] {
	return RIOE.ReadIOEither[A](r)
}

// ReadIO executes a ReaderIOResult computation by providing an environment
// obtained from an IO computation. Unlike ReadIOEither/ReadIOResult, the environment
// acquisition cannot fail (it's a pure IO, not IOResult).
//
// The function first executes the IO[R] to obtain the environment,
// then uses that environment to run the ReaderIOResult[R, A] computation.
//
// Type parameters:
//   - A: The success value type of the ReaderIOResult computation
//   - R: The environment/context type required by the ReaderIOResult
//
// Parameters:
//   - r: An IO[R] that produces the environment (cannot fail)
//
// Returns:
//   - A function that takes a ReaderIOResult[R, A] and returns IOResult[A]
//
// Example:
//
//	type Logger struct { Level string }
//
//	// Get logger (always succeeds)
//	getLogger := func() IO[Logger] {
//	    return func() Logger {
//	        return Logger{Level: "INFO"}
//	    }
//	}
//
//	// Log operation that may fail
//	logMessage := func(msg string) ReaderIOResult[Logger, string] {
//	    return func(logger Logger) IOResult[string] {
//	        return func() Result[string] {
//	            // Log with logger, may fail
//	            return result.Of(fmt.Sprintf("[%s] %s", logger.Level, msg))
//	        }
//	    }
//	}
//
//	// Execute logging with logger
//	logged := ReadIO[string](getLogger())(logMessage("Hello"))()
//
//go:inline
func ReadIO[A, R any](r IO[R]) func(ReaderIOResult[R, A]) IOResult[A] {
	return RIOE.ReadIO[error, A](r)
}
