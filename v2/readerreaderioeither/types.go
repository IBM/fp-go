package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	Option[A any]               = option.Option[A]
	Lazy[A any]                 = lazy.Lazy[A]
	Reader[R, A any]            = reader.Reader[R, A]
	ReaderOption[R, A any]      = readeroption.ReaderOption[R, A]
	ReaderIO[R, A any]          = readerio.ReaderIO[R, A]
	ReaderIOEither[R, E, A any] = readerioeither.ReaderIOEither[R, E, A]
	Either[E, A any]            = either.Either[E, A]
	IOEither[E, A any]          = ioeither.IOEither[E, A]
	IO[A any]                   = io.IO[A]

	ReaderReaderIOEither[R, C, E, A any] = Reader[R, ReaderIOEither[C, E, A]]

	Kleisli[R, C, E, A, B any]  = Reader[A, ReaderReaderIOEither[R, C, E, B]]
	Operator[R, C, E, A, B any] = Kleisli[R, C, E, ReaderReaderIOEither[R, C, E, A], B]
	Lens[S, T any]              = lens.Lens[S, T]
	Trampoline[L, B any]        = tailrec.Trampoline[L, B]
	Predicate[A any]            = predicate.Predicate[A]
)
