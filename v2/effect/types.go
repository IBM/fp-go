package effect

import (
	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Either[E, A any]      = either.Either[E, A]
	Reader[R, A any]      = reader.Reader[R, A]
	ReaderIO[R, A any]    = readerio.ReaderIO[R, A]
	IO[A any]             = io.IO[A]
	IOEither[E, A any]    = ioeither.IOEither[E, A]
	Lazy[A any]           = lazy.Lazy[A]
	IOResult[A any]       = ioresult.IOResult[A]
	ReaderIOResult[A any] = readerioresult.ReaderIOResult[A]
	Monoid[A any]         = monoid.Monoid[A]
	Effect[C, A any]      = readerreaderioresult.ReaderReaderIOResult[C, A]
	Thunk[A any]          = ReaderIOResult[A]
	Predicate[A any]      = predicate.Predicate[A]
	Result[A any]         = result.Result[A]
	Lens[S, T any]        = lens.Lens[S, T]

	Kleisli[C, A, B any]  = readerreaderioresult.Kleisli[C, A, B]
	Operator[C, A, B any] = readerreaderioresult.Operator[C, A, B]
)
