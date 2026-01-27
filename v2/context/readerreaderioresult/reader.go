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

package readerreaderioresult

import (
	"context"
	"time"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	RE "github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readeroption"
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
	"github.com/IBM/fp-go/v2/result"
)

// FromReaderOption converts a ReaderOption to a ReaderReaderIOResult.
// If the option is None, it uses the provided onNone function to generate an error.
//
//go:inline
func FromReaderOption[R, A any](onNone Lazy[error]) Kleisli[R, ReaderOption[R, A], A] {
	return RRIOE.FromReaderOption[R, context.Context, A](onNone)
}

// FromReaderIOResult lifts a ReaderIOResult into a ReaderReaderIOResult.
// This adds an additional reader layer to the computation.
//
//go:inline
func FromReaderIOResult[R, A any](ma ReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderIOEither[context.Context](ma)
}

// FromReaderIO lifts a ReaderIO into a ReaderReaderIOResult.
// The IO computation is wrapped in a Right (success) value.
//
//go:inline
func FromReaderIO[R, A any](ma ReaderIO[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderIO[context.Context, error](ma)
}

// RightReaderIO lifts a ReaderIO into a ReaderReaderIOResult as a Right (success) value.
// Alias for FromReaderIO.
//
//go:inline
func RightReaderIO[R, A any](ma ReaderIO[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightReaderIO[context.Context, error](ma)
}

// LeftReaderIO lifts a ReaderIO that produces an error into a ReaderReaderIOResult as a Left (failure) value.
//
//go:inline
func LeftReaderIO[A, R any](me ReaderIO[R, error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftReaderIO[context.Context, A](me)
}

// MonadMap applies a function to the value inside a ReaderReaderIOResult (Functor operation).
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadMap[R, A, B any](fa ReaderReaderIOResult[R, A], f func(A) B) ReaderReaderIOResult[R, B] {
	return reader.MonadMap(fa, RIOE.Map(f))
}

// Map applies a function to the value inside a ReaderReaderIOResult (Functor operation).
// This is the curried version that returns an operator.
//
//go:inline
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return reader.Map[R](RIOE.Map(f))
}

// MonadMapTo replaces the value inside a ReaderReaderIOResult with a constant value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadMapTo[R, A, B any](fa ReaderReaderIOResult[R, A], b B) ReaderReaderIOResult[R, B] {
	return reader.MonadMap(fa, RIOE.MapTo[A](b))
}

// MapTo replaces the value inside a ReaderReaderIOResult with a constant value.
// This is the curried version that returns an operator.
//
//go:inline
func MapTo[R, A, B any](b B) Operator[R, A, B] {
	return reader.Map[R](RIOE.MapTo[A](b))
}

// MonadChain sequences two computations, where the second depends on the result of the first (Monad operation).
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChain[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return readert.MonadChain(
		RIOE.MonadChain[A, B],
		fa,
		f,
	)
}

// MonadChainFirst sequences two computations but returns the result of the first.
// Useful for performing side effects while preserving the original value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirst[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return chain.MonadChainFirst(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		fa,
		f)
}

// MonadTap is an alias for MonadChainFirst.
// Executes a side effect while preserving the original value.
//
//go:inline
func MonadTap[R, A, B any](fa ReaderReaderIOResult[R, A], f Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirst(fa, f)
}

// MonadChainEitherK chains a computation that returns an Either.
// The Either is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, A, B],
		FromEither[R, B],
		ma,
		f,
	)
}

// ChainEitherK chains a computation that returns an Either.
// The Either is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, A, B],
		FromEither[R, B],
		f,
	)
}

//go:inline
func ChainResultK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, A, B],
		FromEither[R, B],
		f,
	)
}

// MonadChainFirstEitherK chains a computation that returns an Either but preserves the original value.
// Useful for validation or side effects that may fail.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirstEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, A] {
	return fromeither.MonadChainFirstEitherK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromEither[R, B],
		ma,
		f,
	)
}

// MonadTapEitherK is an alias for MonadChainFirstEitherK.
// Executes an Either-returning side effect while preserving the original value.
//
//go:inline
func MonadTapEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f either.Kleisli[error, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstEitherK(ma, f)
}

// ChainFirstEitherK chains a computation that returns an Either but preserves the original value.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirstEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return fromeither.ChainFirstEitherK(
		Chain[R, A, A],
		Map[R, B, A],
		FromEither[R, B],
		f,
	)
}

// TapEitherK is an alias for ChainFirstEitherK.
// Executes an Either-returning side effect while preserving the original value.
//
//go:inline
func TapEitherK[R, A, B any](f either.Kleisli[error, A, B]) Operator[R, A, A] {
	return ChainFirstEitherK[R](f)
}

// MonadChainReaderK chains a computation that returns a Reader.
// The Reader is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

// ChainReaderK chains a computation that returns a Reader.
// The Reader is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReader[R, B],
		f,
	)
}

// MonadChainFirstReaderK chains a computation that returns a Reader but preserves the original value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirstReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReader[R, B],
		ma,
		f,
	)
}

// MonadTapReaderK is an alias for MonadChainFirstReaderK.
// Executes a Reader-returning side effect while preserving the original value.
//
//go:inline
func MonadTapReaderK[R, A, B any](ma ReaderReaderIOResult[R, A], f reader.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderK(ma, f)
}

// ChainFirstReaderK chains a computation that returns a Reader but preserves the original value.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirstReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReader[R, B],
		f,
	)
}

// TapReaderK is an alias for ChainFirstReaderK.
// Executes a Reader-returning side effect while preserving the original value.
//
//go:inline
func TapReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderK(f)
}

// MonadChainReaderIOK chains a computation that returns a ReaderIO.
// The ReaderIO is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

// ChainReaderIOK chains a computation that returns a ReaderIO.
// The ReaderIO is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

// MonadChainFirstReaderIOK chains a computation that returns a ReaderIO but preserves the original value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirstReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReaderIO[R, B],
		ma,
		f,
	)
}

// MonadTapReaderIOK is an alias for MonadChainFirstReaderIOK.
// Executes a ReaderIO-returning side effect while preserving the original value.
//
//go:inline
func MonadTapReaderIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f readerio.Kleisli[R, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderIOK(ma, f)
}

// ChainFirstReaderIOK chains a computation that returns a ReaderIO but preserves the original value.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirstReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReaderIO[R, B],
		f,
	)
}

// TapReaderIOK is an alias for ChainFirstReaderIOK.
// Executes a ReaderIO-returning side effect while preserving the original value.
//
//go:inline
func TapReaderIOK[R, A, B any](f readerio.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderIOK(f)
}

// MonadChainReaderEitherK chains a computation that returns a ReaderEither.
// The ReaderEither is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, B] {
	return fromreader.MonadChainReaderK(
		MonadChain[R, A, B],
		FromReaderEither[R, B],
		ma,
		f,
	)
}

// ChainReaderEitherK chains a computation that returns a ReaderEither.
// The ReaderEither is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, B] {
	return fromreader.ChainReaderK(
		Chain[R, A, B],
		FromReaderEither[R, B],
		f,
	)
}

// MonadChainFirstReaderEitherK chains a computation that returns a ReaderEither but preserves the original value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirstReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, A] {
	return fromreader.MonadChainFirstReaderK(
		MonadChainFirst[R, A, B],
		FromReaderEither[R, B],
		ma,
		f,
	)
}

// MonadTapReaderEitherK is an alias for MonadChainFirstReaderEitherK.
// Executes a ReaderEither-returning side effect while preserving the original value.
//
//go:inline
func MonadTapReaderEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f RE.Kleisli[R, error, A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstReaderEitherK(ma, f)
}

// ChainFirstReaderEitherK chains a computation that returns a ReaderEither but preserves the original value.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirstReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return fromreader.ChainFirstReaderK(
		ChainFirst[R, A, B],
		FromReaderEither[R, B],
		f,
	)
}

// TapReaderEitherK is an alias for ChainFirstReaderEitherK.
// Executes a ReaderEither-returning side effect while preserving the original value.
//
//go:inline
func TapReaderEitherK[R, A, B any](f RE.Kleisli[R, error, A, B]) Operator[R, A, A] {
	return ChainFirstReaderEitherK(f)
}

// ChainReaderOptionK chains a computation that returns a ReaderOption.
// If the option is None, it uses the provided onNone function to generate an error.
// Returns a function that takes a ReaderOption Kleisli and returns an operator.
//
//go:inline
func ChainReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, B] {
	return RRIOE.ChainReaderOptionK[R, context.Context, A, B](onNone)
}

// ChainFirstReaderOptionK chains a computation that returns a ReaderOption but preserves the original value.
// If the option is None, it uses the provided onNone function to generate an error.
// Returns a function that takes a ReaderOption Kleisli and returns an operator.
func ChainFirstReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return RRIOE.ChainFirstReaderOptionK[R, context.Context, A, B](onNone)
}

// TapReaderOptionK is an alias for ChainFirstReaderOptionK.
// Executes a ReaderOption-returning side effect while preserving the original value.
//
//go:inline
func TapReaderOptionK[R, A, B any](onNone Lazy[error]) func(readeroption.Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirstReaderOptionK[R, A, B](onNone)
}

// MonadChainIOEitherK chains a computation that returns an IOEither.
// The IOEither is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainIOEitherK[R, A, B any](ma ReaderReaderIOResult[R, A], f IOE.Kleisli[error, A, B]) ReaderReaderIOResult[R, B] {
	return fromioeither.MonadChainIOEitherK(
		MonadChain[R, A, B],
		FromIOEither[R, B],
		ma,
		f,
	)
}

// ChainIOEitherK chains a computation that returns an IOEither.
// The IOEither is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainIOEitherK[R, A, B any](f IOE.Kleisli[error, A, B]) Operator[R, A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[R, A, B],
		FromIOEither[R, B],
		f,
	)
}

// MonadChainIOK chains a computation that returns an IO.
// The IO is automatically lifted into ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, B] {
	return fromio.MonadChainIOK(
		MonadChain[R, A, B],
		FromIO[R, B],
		ma,
		f,
	)
}

// ChainIOK chains a computation that returns an IO.
// The IO is automatically lifted into ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func ChainIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, B] {
	return fromio.ChainIOK(
		Chain[R, A, B],
		FromIO[R, B],
		f,
	)
}

// MonadChainFirstIOK chains a computation that returns an IO but preserves the original value.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainFirstIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[R, A, A],
		MonadMap[R, B, A],
		FromIO[R, B],
		ma,
		f,
	)
}

// MonadTapIOK is an alias for MonadChainFirstIOK.
// Executes an IO-returning side effect while preserving the original value.
//
//go:inline
func MonadTapIOK[R, A, B any](ma ReaderReaderIOResult[R, A], f io.Kleisli[A, B]) ReaderReaderIOResult[R, A] {
	return MonadChainFirstIOK(ma, f)
}

// ChainFirstIOK chains a computation that returns an IO but preserves the original value.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirstIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return fromio.ChainFirstIOK(
		Chain[R, A, A],
		Map[R, B, A],
		FromIO[R, B],
		f,
	)
}

// TapIOK is an alias for ChainFirstIOK.
// Executes an IO-returning side effect while preserving the original value.
//
//go:inline
func TapIOK[R, A, B any](f io.Kleisli[A, B]) Operator[R, A, A] {
	return ChainFirstIOK[R](f)
}

// ChainOptionK chains a computation that returns an Option.
// If the option is None, it uses the provided onNone function to generate an error.
// Returns a function that takes an Option Kleisli and returns an operator.
//
//go:inline
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[R, A, B] {
	return fromeither.ChainOptionK(
		MonadChain[R, A, B],
		FromEither[R, B],
		onNone,
	)
}

// MonadAp applies a function wrapped in a ReaderReaderIOResult to a value wrapped in a ReaderReaderIOResult (Applicative operation).
// This is the monadic version that takes both computations as parameters.
//
//go:inline
func MonadAp[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadAp[B, A],
		fab,
		fa,
	)
}

// MonadApSeq is like MonadAp but evaluates effects sequentially.
//
//go:inline
func MonadApSeq[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadApSeq[B, A],
		fab,
		fa,
	)
}

// MonadApPar is like MonadAp but evaluates effects in parallel.
//
//go:inline
func MonadApPar[R, A, B any](fab ReaderReaderIOResult[R, func(A) B], fa ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, B] {
	return readert.MonadAp[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.MonadApPar[B, A],
		fab,
		fa,
	)
}

// Ap applies a function wrapped in a ReaderReaderIOResult to a value wrapped in a ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func Ap[B, R, A any](fa ReaderReaderIOResult[R, A]) Operator[R, func(A) B, B] {
	return readert.Ap[
		ReaderReaderIOResult[R, A],
		ReaderReaderIOResult[R, B],
		ReaderReaderIOResult[R, func(A) B], R, A](
		RIOE.Ap[B, A],
		fa,
	)
}

// Chain sequences two computations, where the second depends on the result of the first (Monad operation).
// This is the curried version that returns an operator.
//
//go:inline
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return readert.Chain[ReaderReaderIOResult[R, A]](
		RIOE.Chain[A, B],
		f,
	)
}

// ChainFirst sequences two computations but returns the result of the first.
// This is the curried version that returns an operator.
//
//go:inline
func ChainFirst[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return chain.ChainFirst(
		Chain[R, A, A],
		Map[R, B, A],
		f)
}

// Tap is an alias for ChainFirst.
// Executes a side effect while preserving the original value.
//
//go:inline
func Tap[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, A] {
	return ChainFirst(f)
}

// Right creates a ReaderReaderIOResult that succeeds with the given value.
// This is the success constructor for the Result type.
//
//go:inline
func Right[R, A any](a A) ReaderReaderIOResult[R, A] {
	return RRIOE.Right[R, context.Context, error](a)
}

// Left creates a ReaderReaderIOResult that fails with the given error.
// This is the failure constructor for the Result type.
//
//go:inline
func Left[R, A any](e error) ReaderReaderIOResult[R, A] {
	return RRIOE.Left[R, context.Context, A](e)
}

// Of creates a ReaderReaderIOResult that succeeds with the given value (Pointed operation).
// Alias for Right.
//
//go:inline
func Of[R, A any](a A) ReaderReaderIOResult[R, A] {
	return RRIOE.Of[R, context.Context, error](a)
}

// Flatten removes one level of nesting from a nested ReaderReaderIOResult.
// Converts ReaderReaderIOResult[R, ReaderReaderIOResult[R, A]] to ReaderReaderIOResult[R, A].
//
//go:inline
func Flatten[R, A any](mma ReaderReaderIOResult[R, ReaderReaderIOResult[R, A]]) ReaderReaderIOResult[R, A] {
	return MonadChain(mma, function.Identity[ReaderReaderIOResult[R, A]])
}

// FromEither lifts an Either into a ReaderReaderIOResult.
//
//go:inline
func FromEither[R, A any](t Either[error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromEither[R, context.Context](t)
}

// FromResult lifts a Result into a ReaderReaderIOResult.
// Alias for FromEither since Result is Either[error, A].
//
//go:inline
func FromResult[R, A any](t Result[A]) ReaderReaderIOResult[R, A] {
	return FromEither[R](t)
}

// RightReader lifts a Reader into a ReaderReaderIOResult as a Right (success) value.
//
//go:inline
func RightReader[R, A any](ma Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightReader[context.Context, error](ma)
}

// LeftReader lifts a Reader that produces an error into a ReaderReaderIOResult as a Left (failure) value.
//
//go:inline
func LeftReader[A, R any](ma Reader[R, error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftReader[context.Context, A](ma)
}

// FromReader lifts a Reader into a ReaderReaderIOResult.
// The Reader's result is wrapped in a Right (success) value.
//
//go:inline
func FromReader[R, A any](ma Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReader[context.Context, error](ma)
}

// RightIO lifts an IO into a ReaderReaderIOResult as a Right (success) value.
//
//go:inline
func RightIO[R, A any](ma IO[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.RightIO[R, context.Context, error](ma)
}

// LeftIO lifts an IO that produces an error into a ReaderReaderIOResult as a Left (failure) value.
//
//go:inline
func LeftIO[R, A any](ma IO[error]) ReaderReaderIOResult[R, A] {
	return RRIOE.LeftIO[R, context.Context, A](ma)
}

// FromIO lifts an IO into a ReaderReaderIOResult.
// The IO's result is wrapped in a Right (success) value.
//
//go:inline
func FromIO[R, A any](ma IO[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIO[R, context.Context, error](ma)
}

// FromIOEither lifts an IOEither into a ReaderReaderIOResult.
//
//go:inline
func FromIOEither[R, A any](ma IOEither[error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIOEither[R, context.Context](ma)
}

// FromIOResult lifts an IOResult into a ReaderReaderIOResult.
// Alias for FromIOEither since IOResult is IOEither[error, A].
//
//go:inline
func FromIOResult[R, A any](ma IOResult[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromIOEither[R, context.Context](ma)
}

// FromReaderEither lifts a ReaderEither into a ReaderReaderIOResult.
//
//go:inline
func FromReaderEither[R, A any](ma RE.ReaderEither[R, error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromReaderEither[R, context.Context](ma)
}

// Ask retrieves the outer environment R.
// Returns a ReaderReaderIOResult that succeeds with the environment value.
//
//go:inline
func Ask[R any]() ReaderReaderIOResult[R, R] {
	return RRIOE.Ask[R, context.Context, error]()
}

// Asks retrieves a value derived from the outer environment R using the provided function.
//
//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.Asks[context.Context, error](r)
}

// FromOption converts an Option to a ReaderReaderIOResult.
// If the option is None, it uses the provided onNone function to generate an error.
// Returns a function that takes an Option and returns a ReaderReaderIOResult.
//
//go:inline
func FromOption[R, A any](onNone Lazy[error]) func(Option[A]) ReaderReaderIOResult[R, A] {
	return RRIOE.FromOption[R, context.Context, A](onNone)
}

// FromPredicate creates a ReaderReaderIOResult from a predicate.
// If the predicate returns true, the value is wrapped in Right.
// If false, onFalse is called to generate an error wrapped in Left.
//
//go:inline
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A] {
	return RRIOE.FromPredicate[R, context.Context](pred, onFalse)
}

// MonadAlt provides alternative/fallback behavior.
// If the first computation fails, it tries the second (lazy-evaluated).
// This is the monadic version that takes both computations as parameters.
//
//go:inline
func MonadAlt[R, A any](first ReaderReaderIOResult[R, A], second Lazy[ReaderReaderIOResult[R, A]]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadAlt(first, second)
}

// Alt provides alternative/fallback behavior.
// If the first computation fails, it tries the second (lazy-evaluated).
// This is the curried version that returns an operator.
//
//go:inline
func Alt[R, A any](second Lazy[ReaderReaderIOResult[R, A]]) Operator[R, A, A] {
	return RRIOE.Alt(second)
}

// MonadFlap applies a value to a function wrapped in a ReaderReaderIOResult.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadFlap[R, B, A any](fab ReaderReaderIOResult[R, func(A) B], a A) ReaderReaderIOResult[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap applies a value to a function wrapped in a ReaderReaderIOResult.
// This is the curried version that returns an operator.
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

// MonadMapLeft transforms the error value if the computation fails.
// Has no effect if the computation succeeds.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadMapLeft[R, A any](fa ReaderReaderIOResult[R, A], f Endmorphism[error]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadMapLeft(fa, f)
}

// MapLeft transforms the error value if the computation fails.
// Has no effect if the computation succeeds.
// This is the curried version that returns an operator.
//
//go:inline
func MapLeft[R, A any](f Endmorphism[error]) Operator[R, A, A] {
	return RRIOE.MapLeft[R, context.Context, A](f)
}

// Read provides a specific outer environment value to a computation.
// Converts ReaderReaderIOResult[R, A] to ReaderIOResult[context.Context, A].
//
//go:inline
func Read[A, R any](r R) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.Read[context.Context, error, A](r)
}

// ReadIOEither provides an outer environment value from an IOEither to a computation.
//
//go:inline
func ReadIOEither[A, R any](rio IOEither[error, R]) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.ReadIOEither[A, R, context.Context](rio)
}

// ReadIO provides an outer environment value from an IO to a computation.
//
//go:inline
func ReadIO[A, R any](rio IO[R]) func(ReaderReaderIOResult[R, A]) ReaderIOResult[context.Context, A] {
	return RRIOE.ReadIO[context.Context, error, A](rio)
}

// MonadChainLeft handles errors by chaining a recovery computation.
// If the computation fails, the error is passed to f for recovery.
// This is the monadic version that takes the computation as the first parameter.
//
//go:inline
func MonadChainLeft[R, A any](fa ReaderReaderIOResult[R, A], f Kleisli[R, error, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.MonadChainLeft(fa, f)
}

// ChainLeft handles errors by chaining a recovery computation.
// If the computation fails, the error is passed to f for recovery.
// This is the curried version that returns an operator.
//
//go:inline
func ChainLeft[R, A any](f Kleisli[R, error, A]) func(ReaderReaderIOResult[R, A]) ReaderReaderIOResult[R, A] {
	return RRIOE.ChainLeft(f)
}

// Delay adds a time delay before executing the computation.
// Useful for rate limiting, retry backoff, or scheduled execution.
//
//go:inline
func Delay[R, A any](delay time.Duration) Operator[R, A, A] {
	return reader.Map[R](RIOE.Delay[A](delay))
}

//go:inline
func Defer[R, A any](fa Lazy[ReaderReaderIOResult[R, A]]) ReaderReaderIOResult[R, A] {
	return RRIOE.Defer(fa)
}
