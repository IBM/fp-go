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
	"context"

	CRRI "github.com/IBM/fp-go/v2/idiomatic/context/readerresult"
	IORI "github.com/IBM/fp-go/v2/idiomatic/ioresult"
	OI "github.com/IBM/fp-go/v2/idiomatic/option"
	RIORI "github.com/IBM/fp-go/v2/idiomatic/readerioresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

type (
	// KleisliI is the idiomatic counterpart of [Kleisli]. It is a plain Go function
	// of the form
	//
	//	func(A) func(context.Context) func() (B, error)
	//
	// i.e. a function that reads the context, performs IO and may fail.
	KleisliI[A, B any] = RIORI.Kleisli[context.Context, A, B]
)

// FromIdiomatic lifts a plain Go function of the form
// `func(context.Context) func() (A, error)` into a [ReaderIOResult].
//
//go:inline
func FromIdiomatic[A any](f func(context.Context) func() (A, error)) ReaderIOResult[A] {
	return RIOR.FromIdiomatic[context.Context](f)
}

// FromResultI lifts an idiomatic Go (value, error) pair into a [ReaderIOResult]
// that ignores the context.
//
//go:inline
func FromResultI[A any](a A, err error) ReaderIOResult[A] {
	return RIOR.FromResultI[context.Context](a, err)
}

// FromIOResultI converts an idiomatic IOResult (a `func() (A, error)`) into a
// [ReaderIOResult] that ignores the context.
//
//go:inline
func FromIOResultI[A any](mr IORI.IOResult[A]) ReaderIOResult[A] {
	return RIOR.FromIOResultI[context.Context](mr)
}

// FromReaderResultI converts an idiomatic ReaderResult (a
// `func(context.Context) (A, error)`) into a [ReaderIOResult].
//
//go:inline
func FromReaderResultI[A any](rr CRRI.ReaderResult[A]) ReaderIOResult[A] {
	return RIOR.FromReaderResultI[context.Context](rr)
}

// FromReaderIOResultI converts an idiomatic ReaderIOResult (a
// `func(context.Context) func() (A, error)`) into the functional [ReaderIOResult].
//
//go:inline
func FromReaderIOResultI[A any](rr RIORI.ReaderIOResult[context.Context, A]) ReaderIOResult[A] {
	return RIOR.FromReaderIOResultI[context.Context](rr)
}

// MonadChainI sequences a [ReaderIOResult] with an idiomatic [KleisliI] arrow.
//
//go:inline
func MonadChainI[A, B any](ma ReaderIOResult[A], f KleisliI[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainI(ma, f)
}

// ChainI is the curried version of [MonadChainI].
//
//go:inline
func ChainI[A, B any](f KleisliI[A, B]) Operator[A, B] {
	return RIOR.ChainI(f)
}

// MonadChainFirstI sequences ma with an idiomatic [KleisliI] arrow for its side
// effects but returns the original value of ma.
//
//go:inline
func MonadChainFirstI[A, B any](ma ReaderIOResult[A], f KleisliI[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstI(ma, f)
}

// ChainFirstI is the curried version of [MonadChainFirstI].
//
//go:inline
func ChainFirstI[A, B any](f KleisliI[A, B]) Operator[A, A] {
	return RIOR.ChainFirstI(f)
}

// MonadTapI is an alias for [MonadChainFirstI].
//
//go:inline
func MonadTapI[A, B any](ma ReaderIOResult[A], f KleisliI[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapI(ma, f)
}

// TapI is an alias for [ChainFirstI].
//
//go:inline
func TapI[A, B any](f KleisliI[A, B]) Operator[A, A] {
	return RIOR.TapI(f)
}

// MonadChainEitherIK is the idiomatic version of [MonadChainEitherK].
// It chains a [ReaderIOResult] with a plain Go function returning (B, error).
//
//go:inline
func MonadChainEitherIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainEitherIK(ma, f)
}

// MonadChainResultIK is an alias for [MonadChainEitherIK] with more explicit naming.
//
//go:inline
func MonadChainResultIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainResultIK(ma, f)
}

// ChainEitherIK is the curried version of [MonadChainEitherIK].
//
//go:inline
func ChainEitherIK[A, B any](f RI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainEitherIK[context.Context](f)
}

// ChainResultIK is an alias for [ChainEitherIK] with more explicit naming.
// It is the idiomatic counterpart of [ChainResultK].
//
//go:inline
func ChainResultIK[A, B any](f RI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainResultIK[context.Context](f)
}

// MonadChainFirstEitherIK is the idiomatic version of [MonadChainFirstEitherK].
//
//go:inline
func MonadChainFirstEitherIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstEitherIK(ma, f)
}

// MonadChainFirstResultIK is an alias for [MonadChainFirstEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstResultIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstResultIK(ma, f)
}

// MonadTapEitherIK is an alias for [MonadChainFirstEitherIK].
//
//go:inline
func MonadTapEitherIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapEitherIK(ma, f)
}

// MonadTapResultIK is an alias for [MonadChainFirstResultIK].
//
//go:inline
func MonadTapResultIK[A, B any](ma ReaderIOResult[A], f RI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapResultIK(ma, f)
}

// ChainFirstEitherIK is the curried version of [MonadChainFirstEitherIK].
//
//go:inline
func ChainFirstEitherIK[A, B any](f RI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstEitherIK[context.Context](f)
}

// ChainFirstResultIK is an alias for [ChainFirstEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstResultIK[A, B any](f RI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstResultIK[context.Context](f)
}

// TapEitherIK is an alias for [ChainFirstEitherIK].
//
//go:inline
func TapEitherIK[A, B any](f RI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapEitherIK[context.Context](f)
}

// TapResultIK is an alias for [ChainFirstResultIK].
//
//go:inline
func TapResultIK[A, B any](f RI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapResultIK[context.Context](f)
}

// MonadChainIOEitherIK is the idiomatic version of [MonadChainIOEitherK].
// It chains a plain Go function returning a deferred `func() (B, error)`.
//
//go:inline
func MonadChainIOEitherIK[A, B any](ma ReaderIOResult[A], f IORI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainIOEitherIK(ma, f)
}

// MonadChainIOResultIK is an alias for [MonadChainIOEitherIK] with more explicit naming.
//
//go:inline
func MonadChainIOResultIK[A, B any](ma ReaderIOResult[A], f IORI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainIOResultIK(ma, f)
}

// ChainIOEitherIK is the curried version of [MonadChainIOEitherIK].
//
//go:inline
func ChainIOEitherIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainIOEitherIK[context.Context](f)
}

// ChainIOResultIK is an alias for [ChainIOEitherIK] with more explicit naming.
//
//go:inline
func ChainIOResultIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainIOResultIK[context.Context](f)
}

// ChainFirstIOEitherIK is the idiomatic version of [ChainFirstIOEitherK].
//
//go:inline
func ChainFirstIOEitherIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstIOEitherIK[context.Context](f)
}

// ChainFirstIOResultIK is an alias for [ChainFirstIOEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstIOResultIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstIOResultIK[context.Context](f)
}

// TapIOEitherIK is an alias for [ChainFirstIOEitherIK].
//
//go:inline
func TapIOEitherIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapIOEitherIK[context.Context](f)
}

// TapIOResultIK is an alias for [ChainFirstIOResultIK].
//
//go:inline
func TapIOResultIK[A, B any](f IORI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapIOResultIK[context.Context](f)
}

// MonadChainReaderEitherIK is the idiomatic version of [MonadChainReaderEitherK].
// It chains a plain Go function returning `func(context.Context) (B, error)`.
//
//go:inline
func MonadChainReaderEitherIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainReaderEitherIK(ma, f)
}

// MonadChainReaderResultIK is an alias for [MonadChainReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainReaderResultIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[B] {
	return RIOR.MonadChainReaderResultIK(ma, f)
}

// ChainReaderEitherIK is the curried version of [MonadChainReaderEitherIK].
//
//go:inline
func ChainReaderEitherIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainReaderEitherIK(f)
}

// ChainReaderResultIK is an alias for [ChainReaderEitherIK] with more explicit naming.
// It is the idiomatic counterpart of [ChainReaderResultK].
//
//go:inline
func ChainReaderResultIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainReaderResultIK(f)
}

// MonadChainFirstReaderEitherIK is the idiomatic version of [MonadChainFirstReaderEitherK].
//
//go:inline
func MonadChainFirstReaderEitherIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstReaderEitherIK(ma, f)
}

// MonadChainFirstReaderResultIK is an alias for [MonadChainFirstReaderEitherIK] with more explicit naming.
//
//go:inline
func MonadChainFirstReaderResultIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadChainFirstReaderResultIK(ma, f)
}

// MonadTapReaderEitherIK is an alias for [MonadChainFirstReaderEitherIK].
//
//go:inline
func MonadTapReaderEitherIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapReaderEitherIK(ma, f)
}

// MonadTapReaderResultIK is an alias for [MonadChainFirstReaderResultIK].
//
//go:inline
func MonadTapReaderResultIK[A, B any](ma ReaderIOResult[A], f CRRI.Kleisli[A, B]) ReaderIOResult[A] {
	return RIOR.MonadTapReaderResultIK(ma, f)
}

// ChainFirstReaderEitherIK is the curried version of [MonadChainFirstReaderEitherIK].
//
//go:inline
func ChainFirstReaderEitherIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderEitherIK(f)
}

// ChainFirstReaderResultIK is an alias for [ChainFirstReaderEitherIK] with more explicit naming.
//
//go:inline
func ChainFirstReaderResultIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.ChainFirstReaderResultIK(f)
}

// TapReaderEitherIK is an alias for [ChainFirstReaderEitherIK].
//
//go:inline
func TapReaderEitherIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapReaderEitherIK(f)
}

// TapReaderResultIK is an alias for [ChainFirstReaderResultIK].
//
//go:inline
func TapReaderResultIK[A, B any](f CRRI.Kleisli[A, B]) Operator[A, A] {
	return RIOR.TapReaderResultIK(f)
}

// ChainOptionIK is the idiomatic version of [ChainOptionK].
// It chains a plain Go function returning the comma-ok pair (B, bool), converting a
// false flag into an error produced by onNone.
//
//go:inline
func ChainOptionIK[A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[A, B] {
	return RIOR.ChainOptionIK[context.Context, A, B](onNone)
}
