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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

//go:inline
func Do[C, S any](
	empty S,
) Effect[C, S] {
	return readerreaderioresult.Of[C](empty)
}

//go:inline
func Bind[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.Bind(setter, f)
}

//go:inline
func Let[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[C, S1, S2] {
	return readerreaderioresult.Let[C](setter, f)
}

//go:inline
func LetTo[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[C, S1, S2] {
	return readerreaderioresult.LetTo[C](setter, b)
}

//go:inline
func BindTo[C, S1, T any](
	setter func(T) S1,
) Operator[C, T, S1] {
	return readerreaderioresult.BindTo[C](setter)
}

//go:inline
func ApS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Effect[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApS(setter, fa)
}

//go:inline
func ApSL[C, S, T any](
	lens Lens[S, T],
	fa Effect[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApSL(lens, fa)
}

//go:inline
func BindL[C, S, T any](
	lens Lens[S, T],
	f func(T) Effect[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindL(lens, f)
}

//go:inline
func LetL[C, S, T any](
	lens Lens[S, T],
	f func(T) T,
) Operator[C, S, S] {
	return readerreaderioresult.LetL[C](lens, f)
}

//go:inline
func LetToL[C, S, T any](
	lens Lens[S, T],
	b T,
) Operator[C, S, S] {
	return readerreaderioresult.LetToL[C](lens, b)
}

//go:inline
func BindIOEitherK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioeither.Kleisli[error, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOEitherK[C](setter, f)
}

//go:inline
func BindIOResultK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOResultK[C](setter, f)
}

//go:inline
func BindIOK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOK[C](setter, f)
}

//go:inline
func BindReaderK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindReaderK(setter, f)
}

//go:inline
func BindReaderIOK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindReaderIOK(setter, f)
}

//go:inline
func BindEitherK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f either.Kleisli[error, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindEitherK[C](setter, f)
}

//go:inline
func BindIOEitherKL[C, S, T any](
	lens Lens[S, T],
	f ioeither.Kleisli[error, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindIOEitherKL[C](lens, f)
}

//go:inline
func BindIOKL[C, S, T any](
	lens Lens[S, T],
	f io.Kleisli[T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindIOKL[C](lens, f)
}

//go:inline
func BindReaderKL[C, S, T any](
	lens Lens[S, T],
	f reader.Kleisli[C, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindReaderKL(lens, f)
}

//go:inline
func BindReaderIOKL[C, S, T any](
	lens Lens[S, T],
	f readerio.Kleisli[C, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindReaderIOKL(lens, f)
}

//go:inline
func ApIOEitherS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[error, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApIOEitherS[C](setter, fa)
}

//go:inline
func ApIOS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApIOS[C](setter, fa)
}

//go:inline
func ApReaderS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApReaderS(setter, fa)
}

//go:inline
func ApReaderIOS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApReaderIOS(setter, fa)
}

//go:inline
func ApEitherS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[error, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApEitherS[C](setter, fa)
}

//go:inline
func ApIOEitherSL[C, S, T any](
	lens Lens[S, T],
	fa IOEither[error, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApIOEitherSL[C](lens, fa)
}

//go:inline
func ApIOSL[C, S, T any](
	lens Lens[S, T],
	fa IO[T],
) Operator[C, S, S] {
	return readerreaderioresult.ApIOSL[C](lens, fa)
}

//go:inline
func ApReaderSL[C, S, T any](
	lens Lens[S, T],
	fa Reader[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApReaderSL(lens, fa)
}

//go:inline
func ApReaderIOSL[C, S, T any](
	lens Lens[S, T],
	fa ReaderIO[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApReaderIOSL(lens, fa)
}

//go:inline
func ApEitherSL[C, S, T any](
	lens Lens[S, T],
	fa Either[error, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApEitherSL[C](lens, fa)
}
