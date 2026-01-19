package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/traversal/result"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/readerreaderioeither"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	Option[A any]            = option.Option[A]
	Lazy[A any]              = lazy.Lazy[A]
	Reader[R, A any]         = reader.Reader[R, A]
	ReaderOption[R, A any]   = readeroption.ReaderOption[R, A]
	ReaderIO[R, A any]       = readerio.ReaderIO[R, A]
	ReaderIOResult[R, A any] = readerioresult.ReaderIOResult[R, A]
	Either[E, A any]         = either.Either[E, A]
	Result[A any]            = result.Result[A]
	IOEither[E, A any]       = ioeither.IOEither[E, A]
	IOResult[A any]          = ioresult.IOResult[A]
	IO[A any]                = io.IO[A]

	ReaderReaderIOEither[R, C, E, A any] = readerreaderioeither.ReaderReaderIOEither[R, C, E, A]

	ReaderReaderIOResult[R, A any] = ReaderReaderIOEither[R, context.Context, error, A]

	Kleisli[R, A, B any]  = Reader[A, ReaderReaderIOResult[R, B]]
	Operator[R, A, B any] = Kleisli[R, ReaderReaderIOResult[R, A], B]
	Lens[S, T any]        = lens.Lens[S, T]
	Trampoline[L, B any]  = tailrec.Trampoline[L, B]
	Predicate[A any]      = predicate.Predicate[A]

	Endmorphism[A any] = endomorphism.Endomorphism[A]
)
