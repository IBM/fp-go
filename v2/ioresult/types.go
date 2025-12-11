package ioresult

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	IO[A any]           = io.IO[A]
	Lazy[A any]         = lazy.Lazy[A]
	Either[E, A any]    = either.Either[E, A]
	Result[A any]       = result.Result[A]
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// IOEither represents a synchronous computation that may fail
	// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details
	IOResult[A any]  = IO[Result[A]]
	Monoid[A any]    = monoid.Monoid[IOResult[A]]
	Semigroup[A any] = semigroup.Semigroup[IOResult[A]]

	Kleisli[A, B any]  = reader.Reader[A, IOResult[B]]
	Operator[A, B any] = Kleisli[IOResult[A], B]

	Consumer[A any] = consumer.Consumer[A]
)
