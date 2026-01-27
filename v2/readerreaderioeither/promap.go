package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerioeither"
)

//go:inline
func Local[C, E, A, R1, R2 any](f func(R2) R1) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return reader.Local[ReaderIOEither[C, E, A]](f)
}

//go:inline
func LocalIOK[C, E, A, R1, R2 any](f io.Kleisli[R2, R1]) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return func(rri ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
		return F.Flow4(
			f,
			io.Map(rri),
			readerioeither.FromIO[C],
			readerioeither.Flatten,
		)
	}
}

//go:inline
func LocalIOEitherK[C, A, R1, R2, E any](f ioeither.Kleisli[E, R2, R1]) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return func(rri ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
		return F.Flow4(
			f,
			ioeither.Map[E](rri),
			readerioeither.FromIOEither[C],
			readerioeither.Flatten,
		)
	}
}

//go:inline
func LocalEitherK[C, A, R1, R2, E any](f either.Kleisli[E, R2, R1]) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return func(rri ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
		return F.Flow4(
			f,
			either.Map[E](rri),
			readerioeither.FromEither[C],
			readerioeither.Flatten,
		)
	}
}

//go:inline
func LocalReaderIOEitherK[A, C, E, R1, R2 any](f readerioeither.Kleisli[C, E, R2, R1]) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return func(rri ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
		return F.Flow3(
			f,
			readerioeither.Map[C, E](rri),
			readerioeither.Flatten,
		)
	}
}

//go:inline
func LocalReaderReaderIOEitherK[A, C, E, R1, R2 any](f Kleisli[R2, C, E, R2, R1]) func(ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
	return func(rri ReaderReaderIOEither[R1, C, E, A]) ReaderReaderIOEither[R2, C, E, A] {
		return F.Flow2(
			reader.AsksReader(f),
			readerioeither.Chain(rri),
		)
	}
}
