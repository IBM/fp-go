package io

import (
	"iter"

	"github.com/IBM/fp-go/v2/consumer"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/semigroup"
)

type (

	// IO represents a synchronous computation that cannot fail
	// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioltagt] for more details
	IO[A any]      = func() A
	Pair[L, R any] = pair.Pair[L, R]

	Kleisli[A, B any]  = reader.Reader[A, IO[B]]
	Operator[A, B any] = Kleisli[IO[A], B]
	Monoid[A any]      = M.Monoid[IO[A]]
	Semigroup[A any]   = S.Semigroup[IO[A]]

	Consumer[A any] = consumer.Consumer[A]

	Seq[T any] = iter.Seq[T]
)
